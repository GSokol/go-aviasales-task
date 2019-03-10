//
// cache.go
// Copyright (C) 2019 Grigorii Sokolik <g.sokol99@g-sokol.info>
//
// Distributed under terms of the MIT license.
//

package util

import (
	"errors"
	"time"

	"github.com/GSokol/go-aviasales-task/internal/config"
	"github.com/GSokol/go-aviasales-task/internal/storage"
	"github.com/GSokol/go-aviasales-task/internal/storage/inmemory"
	"github.com/GSokol/go-aviasales-task/internal/storage/redis"
	redigo "github.com/gomodule/redigo/redis"
	cache "github.com/patrickmn/go-cache"
)

var BadCacheType = errors.New("Bad cache type")

func NewCache(cfg config.Cache) (storage.Cache, error) {
	switch cfg.Type {
	case config.CacheTypeRedis:
		return newRedisCache(cfg)
	case config.CacheTypeInMemory:
		return newInMemoryConfig(cfg)
	default:
		return nil, BadCacheType
	}
}

func newRedisCache(cfg config.Cache) (storage.Cache, error) {
	opts := []redigo.DialOption{}

	if cfg.ConnectTimeoutMs.Valid {
		opts = append(opts, redigo.DialConnectTimeout(time.Duration(cfg.ConnectTimeoutMs.Value)*time.Millisecond))
	}
	if cfg.Database.Valid {
		opts = append(opts, redigo.DialDatabase(cfg.Database.Value))
	}
	if cfg.KeepAliveMs.Valid {
		opts = append(opts, redigo.DialKeepAlive(time.Duration(cfg.KeepAliveMs.Value)*time.Millisecond))
	}
	if cfg.Password.Valid {
		opts = append(opts, redigo.DialPassword(cfg.Password.Value))
	}
	if cfg.ReadTimeoutMs.Valid {
		opts = append(opts, redigo.DialReadTimeout(time.Duration(cfg.ReadTimeoutMs.Value)*time.Millisecond))
	}
	if cfg.TLS != nil {
		opts = append(opts, redigo.DialTLSConfig(cfg.TLS))
	}
	if cfg.TLSSkipVerify.Valid {
		opts = append(opts, redigo.DialTLSSkipVerify(cfg.TLSSkipVerify.Value))
	}
	if cfg.UseTLS.Valid {
		opts = append(opts, redigo.DialUseTLS(cfg.UseTLS.Value))
	}
	if cfg.WriteTimeoutMs.Valid {
		opts = append(opts, redigo.DialWriteTimeout(time.Duration(cfg.WriteTimeoutMs.Value)*time.Millisecond))
	}

	addr := cfg.Addr
	return redis.New(&redigo.Pool{
		MaxIdle: cfg.MaxIdle,
		Dial:    func() (redigo.Conn, error) { return redigo.DialURL(addr, opts...) },
	}), nil
}

func newInMemoryConfig(cfg config.Cache) (storage.Cache, error) {
	return inmemory.New(cache.New(cache.NoExpiration, cache.DefaultExpiration)), nil
}
