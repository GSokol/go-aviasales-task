//
// parser.go
// Copyright (C) 2019 Grigorii Sokolik <g.sokol99@g-sokol.info>
//
// Distributed under terms of the MIT license.
//

package util

import (
	"encoding/json"
	"io/ioutil"

	"github.com/GSokol/go-aviasales-task/internal/config"
)

func Parse(fn string) (config.Config, error) {
	var cfg config.Config
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return cfg, err
	}

	err = json.Unmarshal(data, &cfg)
	return cfg, err
}
