//
// cache.go
// Copyright (C) 2019 Grigorii Sokolik <g.sokol99@g-sokol.info>
//
// Distributed under terms of the MIT license.
//

package config

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
)

type cacheType string

const (
	CacheTypeRedis    cacheType = "redis"
	CacheTypeInMemory cacheType = "inMemory"
)

var UnsupportedCacheTypeError = fmt.Errorf(
	"unsupported type error; supported are ('%s', '%s')",
	CacheTypeRedis,
	CacheTypeInMemory,
)

func (ct *cacheType) UnmarshalJSON(b []byte) error {
	var buf string

	if err := json.Unmarshal(b, &buf); err != nil {
		return err
	}

	c := cacheType(buf)
	switch c {
	case CacheTypeRedis, CacheTypeInMemory:
		(*ct) = c
	default:
		return UnsupportedCacheTypeError
	}
	return nil
}

type Cache struct {
	Type             cacheType   `json:"type" ignored:"true"`
	MaxIdle          int         `json:"maxIdle,omitempty" ignored:"true"`
	Addr             string      `json:"addr,omitempty" ignored:"true"`
	ConnectTimeoutMs ValidInt64  `json:"connectTimeoutMs,omitempty" ignored:"true"`
	Database         ValidInt    `json:"database,omitempty" ignored:"true"`
	KeepAliveMs      ValidInt64  `json:"keepAlive,omitempty" ignored:"true"`
	Password         ValidString `json:"password,omitempty" envconfig:"password"`
	ReadTimeoutMs    ValidInt64  `json:"readTimeout" ignored:"true"`
	TLS              *tls.Config `json:"tls" ignored:"true"`
	TLSSkipVerify    ValidBool   `json:"tlsSkipVerify,omitempty" ignored:"true"`
	UseTLS           ValidBool   `json:"useTLS,omitempty" ignored:"true"`
	WriteTimeoutMs   ValidInt64  `json:"writeTimeoutMs,omitempty" ignored:"true"`
}
