//
// main.go
// Copyright (C) 2019 Grigorii Sokolik <g.sokol99@g-sokol.info>
//
// Distributed under terms of the MIT license.
//

package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"

	configUtil "github.com/GSokol/go-aviasales-task/internal/config/util"
	"github.com/GSokol/go-aviasales-task/internal/server"
	"github.com/GSokol/go-aviasales-task/internal/storage"
	httpClient "github.com/GSokol/go-aviasales-task/pkg/aviasales/places/client"
	"github.com/buaazp/fasthttprouter"
	"github.com/kelseyhightower/envconfig"
	"github.com/panjf2000/ants"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

func main() {
	cfgPath := os.Getenv("AV_CONFIG_PATH")

	cfg, err := configUtil.Parse(cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	if err := envconfig.Process("config", &cfg); err != nil {
		log.Fatal(err)
	}

	logger, err := cfg.Logger.Build()
	if err != nil {
		log.Fatalf("failed to init logger: %s", err)
	}

	pool, err := ants.NewTimingPool(cfg.Pool.Size, cfg.Pool.ExpiritySec)
	if err != nil {
		logger.Fatal("failed to init pool", zap.Error(err))
	}

	client := httpClient.NewClient(
		httpClient.Host(cfg.Client.Host),
		httpClient.TimeoutMs(cfg.Client.TimeoutMs),
	)

	cache, err := configUtil.NewCache(cfg.Cache)
	if err != nil {
		logger.Fatal("failed to init cache", zap.Error(err))
	}

	srv := server.NewServer(
		storage.NewClient(client),
		storage.NewCached(cache),
		pool,
		logger,
		cfg.Client.TimeoutMs,
		cfg.Server.TimeoutMs,
	)

	router := fasthttprouter.New()

	router.GET("/search", srv.RequestHandler())
	router.GET("/healthz", server.HealthCheckHandler)

	httpServer := fasthttp.Server{Handler: router.Handler}

	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := httpServer.Shutdown(); err != nil {
			logger.Error("Failed to shutdowm", zap.Error(err))
			logger.Sync()
		}
		close(idleConnsClosed)
	}()

	go func() {
		if err := http.ListenAndServe(":18080", nil); err != nil {
			logger.Error("Failed to start pprof server", zap.Error(err))
			logger.Sync()
		}
	}()

	if err := httpServer.ListenAndServe(":8080"); err != nil {
		logger.Error("Failed to shutdowm", zap.Error(err))
		logger.Sync()
	}

	<-idleConnsClosed
}
