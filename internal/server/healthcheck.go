//
// healthcheck.go
// Copyright (C) 2019 Grigorii Sokolik <g.sokol99@g-sokol.info>
//
// Distributed under terms of the MIT license.
//

package server

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

func HealthCheckHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprint(ctx, "OK")
}
