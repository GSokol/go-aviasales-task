//
// options.go
// Copyright (C) 2019 Grigorii Sokolik <g.sokol99@g-sokol.info>
//
// Distributed under terms of the MIT license.
//

package client

import "time"

type Options struct {
	host    string
	timeout time.Duration
}

func defaultOptions() Options {
	return Options{}
}

type OptionFunc func(*Options)

func Host(host string) OptionFunc {
	return func(o *Options) {
		o.host = host
	}
}

func TimeoutMs(timeout int64) OptionFunc {
	return func(o *Options) {
		o.timeout = time.Duration(timeout) * time.Millisecond
	}
}
