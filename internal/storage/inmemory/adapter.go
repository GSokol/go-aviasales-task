//
// adapter.go
// Copyright (C) 2019 Grigorii Sokolik <g.sokol99@g-sokol.info>
//
// Distributed under terms of the MIT license.
//

package inmemory

import (
	"errors"
	"time"

	cache "github.com/patrickmn/go-cache"
)

var NotFound = errors.New("Not found")

type Adapter struct {
	c *cache.Cache
}

func New(c *cache.Cache) *Adapter {
	return &Adapter{c: c}
}

func (a *Adapter) Set(key string, value interface{}, d time.Duration) error {
	a.c.Set(key, value, d)
	return nil
}

func (a *Adapter) Get(key string) (interface{}, error) {
	ret, ok := a.c.Get(key)
	if !ok {
		return nil, NotFound
	}
	return ret, nil
}
