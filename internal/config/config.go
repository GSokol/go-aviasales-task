//
// config.go
// Copyright (C) 2019 Grigorii Sokolik <g.sokol99@g-sokol.info>
//
// Distributed under terms of the MIT license.
//

package config

import "go.uber.org/zap"

type Client struct {
	Host      string `json:"host"`
	TimeoutMs int64  `json:"timeoutMs"`
}

type Pool struct {
	Size        int `json:"size"`
	ExpiritySec int `json:"expiritySec"`
}

type Server struct {
	TimeoutMs int64 `json:"timeoutMs"`
}

type Config struct {
	Logger zap.Config `json:"logger" ignored:"true"`
	Client Client     `json:"client" ignored:"true"`
	Cache  Cache      `json:"cache" envconfig:"cache"`
	Pool   Pool       `json:"pool" ignored:"true"`
	Server Server     `json:"server" ignored:"true"`
}
