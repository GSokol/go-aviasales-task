//
// cached.go
// Copyright (C) 2019 Grigorii Sokolik <g.sokol99@g-sokol.info>
//
// Distributed under terms of the MIT license.
//

package storage

import (
	"errors"
	"fmt"
	"time"
)

var WrongDataTypeErr = errors.New("wrong data type")

type Cache interface {
	Set(string, interface{}, time.Duration) error
	Get(string) (interface{}, error)
}

type Cached struct {
	cache Cache
}

func NewCached(c Cache) *Cached {
	return &Cached{cache: c}
}

func (c *Cached) Get(term, locale []byte) ([]byte, error) {
	data, err := c.cache.Get(getKey(term, locale))
	if err != nil {
		return nil, err
	}

	bytes, ok := data.([]byte)
	if !ok {
		return nil, WrongDataTypeErr
	}
	return bytes, nil
}

func (c *Cached) Put(term, locale, data []byte) error {
	return c.cache.Set(getKey(term, locale), data, -1)
}

func getKey(term, locale []byte) string {
	return fmt.Sprintf("%s:%s", locale, term)
}
