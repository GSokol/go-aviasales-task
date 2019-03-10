//
// server.go
// Copyright (C) 2019 Grigorii Sokolik <g.sokol99@g-sokol.info>
//
// Distributed under terms of the MIT license.
//

package server

import (
	"errors"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

const (
	argTerm   = "term"
	argLocale = "locale"
)

var ContexCancelledError = errors.New("Context cancelled")

type Source interface {
	Get([]byte, []byte) ([]byte, error)
}

type PutSource interface {
	Source
	Put([]byte, []byte, []byte) error
}

type WorkerPool interface {
	Submit(func()) error
}

type Server struct {
	primarySource        Source
	reserveSource        PutSource
	pool                 WorkerPool
	logger               *zap.Logger
	primarySourceTimeout time.Duration
	requestTimeout       time.Duration
}

func NewServer(
	primarySource Source,
	reserveSource PutSource,
	pool WorkerPool,
	logger *zap.Logger,
	primarySourceTimeoutMs int64,
	requestTimeoutMs int64,
) *Server {
	return &Server{
		primarySource:        primarySource,
		reserveSource:        reserveSource,
		pool:                 pool,
		logger:               logger,
		primarySourceTimeout: time.Duration(primarySourceTimeoutMs) * time.Millisecond,
		requestTimeout:       time.Duration(requestTimeoutMs) * time.Millisecond,
	}
}

func (s *Server) RequestHandler() fasthttp.RequestHandler {
	return fasthttp.TimeoutHandler(
		func(ctx *fasthttp.RequestCtx) {
			func() {
				args := ctx.QueryArgs()
				if !args.Has(argTerm) || !args.Has(argLocale) {
					ctx.Error("Bad Request", fasthttp.StatusBadRequest)
					return
				}

				term := args.Peek(argTerm)
				locale := args.Peek(argLocale)

				timer := time.NewTimer(s.primarySourceTimeout)
				data, err := s.processRequest(term, locale, ctx.Done(), timer.C)
				timer.Stop()
				if err != nil {
					if err == ContexCancelledError {
						s.logger.Error("context cancelled", zap.Error(ctx.Err()))
					}

					ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
					return
				}

				if _, err := ctx.Write(data); err != nil {
					s.logger.Error("failed to write response", zap.Error(err))
				}
			}()
			s.logger.Sync()
		},
		s.requestTimeout,
		"Internal Server Error",
	)
}

func (s *Server) processRequest(
	term []byte,
	locale []byte,
	cancel <-chan struct{},
	primarySourceTimeout <-chan time.Time,
) ([]byte, error) {
	var (
		err  error
		data []byte
		mtx  sync.Mutex
	)

	doneCh := make(chan struct{})
	mtx.Lock()
	if err = s.pool.Submit(func() {
		mtx.Lock()
		data, err = s.primarySource.Get(term, locale)
		if err != nil {
			s.logger.Warn("failed to receive data from primary API", zap.Error(err))
		}
		mtx.Unlock()
		close(doneCh)
	}); err != nil {
		err1 := err
		mtx.Unlock()
		close(doneCh)
		s.logger.Error("failed to add task to pool", zap.Error(err))
		return data, err1
	}
	mtx.Unlock()

	select {
	case <-doneCh:
		mtx.Lock()
		if err != nil {
			data, err := s.useReserveSource(term, locale, cancel)
			mtx.Unlock()
			if err != nil {
				s.logger.Error("failed to receive data from reserve API", zap.Error(err))
			}
			return data, err
		}
		mtx.Unlock()
		if err := s.pool.Submit(func() {
			if err := s.reserveSource.Put(term, locale, data); err != nil {
				s.logger.Warn("failed to persist data", zap.Error(err))
			}
		}); err != nil {
			s.logger.Warn("failed to add task to pool", zap.Error(err))
		}
		return data, err
	case <-primarySourceTimeout:
		s.logger.Info("primary source API timeout")
		data, err := s.useReserveSource(term, locale, cancel)
		if err != nil {
			s.logger.Error("failed to receive data from reserve API", zap.Error(err))
		}
		return data, err

	case <-cancel:
		return nil, ContexCancelledError
	}
}

func (s *Server) useReserveSource(term, locale []byte, cancel <-chan struct{}) ([]byte, error) {
	var (
		err  error
		data []byte
		mtx  sync.Mutex
	)
	doneCh := make(chan struct{})
	mtx.Lock()
	if err = s.pool.Submit(func() {
		mtx.Lock()
		if data, err = s.reserveSource.Get(
			term,
			locale,
		); err != nil {
			s.logger.Error("failed to receive data from reserve API", zap.Error(err))
		}
		mtx.Unlock()
		close(doneCh)
	}); err != nil {
		err1 := err
		mtx.Unlock()
		close(doneCh)
		return data, err1
	}
	mtx.Unlock()
	select {
	case <-doneCh:
		return data, err
	case <-cancel:
		return nil, ContexCancelledError
	}
}
