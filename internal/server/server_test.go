//
// server_test.go
// Copyright (C) 2019 Grigorii Sokolik <g.sokol99@g-sokol.info>
//
// Distributed under terms of the MIT license.
//

package server

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/GSokol/go-aviasales-task/internal/server/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestServerProcessRequestPoolError(t *testing.T) {
	err := errors.New("foo")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	primarySource := mock.NewMockSource(ctrl)
	reserveSource := mock.NewMockPutSource(ctrl)
	pool := mock.NewMockWorkerPool(ctrl)

	pool.EXPECT().Submit(gomock.Any()).Return(err)

	target := NewServer(
		primarySource,
		reserveSource,
		pool,
		zap.NewNop(),
		10,
		10,
	)

	data, actErr := target.processRequest([]byte("foo"), []byte("bar"), nil, nil)
	assert.Nil(t, data)
	assert.Equal(t, err, actErr)
}

func TestServerProcessRequestPrimarySourceOk(t *testing.T) {
	fooTerm := []byte("fooTerm")
	fooLocale := []byte("fooLocale")
	fooRes := []byte("fooRes")

	cancel := make(chan struct{})
	defer func() {
		close(cancel)
	}()
	timeout := make(chan time.Time)
	defer func() {
		close(timeout)
	}()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var wg sync.WaitGroup
	defer func() {
		wg.Wait()
	}()

	primarySource := mock.NewMockSource(ctrl)
	reserveSource := mock.NewMockPutSource(ctrl)
	pool := mock.NewMockWorkerPool(ctrl)

	pool.EXPECT().
		Submit(gomock.Any()).
		DoAndReturn(func(f func()) error {
			go func() {
				f()
			}()
			return nil
		})

	primarySource.EXPECT().
		Get(fooTerm, fooLocale).
		Return(fooRes, nil)

	pool.EXPECT().
		Submit(gomock.Any()).
		DoAndReturn(func(f func()) error {
			wg.Add(1)
			go func() {
				f()
				wg.Done()
			}()
			return nil
		})

	reserveSource.EXPECT().
		Put(fooTerm, fooLocale, fooRes).
		Return(nil)

	target := NewServer(
		primarySource,
		reserveSource,
		pool,
		zap.NewNop(),
		2000,
		100,
	)

	data, actErr := target.processRequest(fooTerm, fooLocale, cancel, timeout)
	assert.NoError(t, actErr)
	assert.Equal(t, fooRes, data)
}

func TestServerProcessRequestPrimarySourceError(t *testing.T) {
	fooTerm := []byte("fooTerm")
	fooLocale := []byte("fooLocale")
	fooRes := []byte("fooRes")
	fooErr := errors.New("fooErr")

	cancel := make(chan struct{})
	defer func() {
		close(cancel)
	}()
	timeout := make(chan time.Time)
	defer func() {
		close(timeout)
	}()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	primarySource := mock.NewMockSource(ctrl)
	reserveSource := mock.NewMockPutSource(ctrl)
	pool := mock.NewMockWorkerPool(ctrl)

	pool.EXPECT().
		Submit(gomock.Any()).
		DoAndReturn(func(f func()) error {
			go func() {
				f()
			}()
			return nil
		})

	primarySource.EXPECT().
		Get(fooTerm, fooLocale).
		Return(nil, fooErr)

	pool.EXPECT().
		Submit(gomock.Any()).
		DoAndReturn(func(f func()) error {
			go func() {
				f()
			}()
			return nil
		})

	reserveSource.EXPECT().
		Get(fooTerm, fooLocale).
		Return(fooRes, nil)

	target := NewServer(
		primarySource,
		reserveSource,
		pool,
		zap.NewNop(),
		2000,
		100,
	)

	data, actErr := target.processRequest(fooTerm, fooLocale, cancel, timeout)
	assert.NoError(t, actErr)
	assert.Equal(t, fooRes, data)
}

func TestServerProcessRequestPrimarySourceTimeout(t *testing.T) {
	fooTerm := []byte("fooTerm")
	fooLocale := []byte("fooLocale")
	fooRes := []byte("fooRes")
	barRes := []byte("barRes")

	cancel := make(chan struct{})
	defer func() {
		close(cancel)
	}()
	timeout := time.NewTimer(2 * time.Millisecond)
	defer func() {
		timeout.Stop()
	}()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var wg sync.WaitGroup
	defer func() {
		wg.Wait()
	}()

	primarySource := mock.NewMockSource(ctrl)
	reserveSource := mock.NewMockPutSource(ctrl)
	pool := mock.NewMockWorkerPool(ctrl)

	pool.EXPECT().
		Submit(gomock.Any()).
		DoAndReturn(func(f func()) error {
			wg.Add(1)
			go func() {
				time.Sleep(3 * time.Millisecond)
				f()
				wg.Done()
			}()
			return nil
		})

	pool.EXPECT().
		Submit(gomock.Any()).
		DoAndReturn(func(f func()) error {
			go func() {
				f()
			}()
			return nil
		})

	reserveSource.EXPECT().
		Get(fooTerm, fooLocale).
		Return(fooRes, nil)

	primarySource.EXPECT().
		Get(fooTerm, fooLocale).
		Return(barRes, nil)

	target := NewServer(
		primarySource,
		reserveSource,
		pool,
		zap.NewNop(),
		2000,
		100,
	)

	data, actErr := target.processRequest(fooTerm, fooLocale, cancel, timeout.C)
	assert.NoError(t, actErr)
	assert.Equal(t, fooRes, data)

}
