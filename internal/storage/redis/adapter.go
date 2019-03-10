//
// adapter.go
// Copyright (C) 2019 Grigorii Sokolik <g.sokol99@g-sokol.info>
//
// Distributed under terms of the MIT license.
//

package redis

import (
	"errors"
	"time"

	redigo "github.com/gomodule/redigo/redis"
)

var NilConnectionError = errors.New("nil connection")

type Adapter struct {
	pool *redigo.Pool
}

func New(pool *redigo.Pool) *Adapter {
	return &Adapter{pool: pool}
}

func (a *Adapter) Set(key string, value interface{}, d time.Duration) error {
	conn := a.pool.Get()
	if conn == nil {
		return NilConnectionError
	}
	defer conn.Close()
	if d == -1 {
		return conn.Send("SET", key, value)
	}
	return conn.Send("SETEX", key, d, value)
}

func (a *Adapter) Get(key string) (interface{}, error) {
	conn := a.pool.Get()
	if conn == nil {
		return nil, NilConnectionError
	}
	defer conn.Close()
	return conn.Do("GET", key)
}
