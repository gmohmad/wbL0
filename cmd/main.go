package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"gihub.com/gmohmad/wb_l0/internal/config"
	"gihub.com/gmohmad/wb_l0/internal/http/handlers/orders"
	"gihub.com/gmohmad/wb_l0/internal/nats/subscribers"
	"gihub.com/gmohmad/wb_l0/internal/storage"
	"gihub.com/gmohmad/wb_l0/internal/storage/cache"
	"gihub.com/gmohmad/wb_l0/internal/storage/postgres"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting the app", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	postgres, err := postgres.NewClient(ctx, cfg, 5, 3, log)
	if err != nil {
		log.Error(err.Error(), slog.Any("smth", postgres))
		os.Exit(1)
	}

	storage := storage.NewStorage(postgres)
	cache := cache.NewCache()

	if err := cache.FillUpCache(ctx, storage); err != nil {
		log.Info(fmt.Sprintf("Error filling cache from db: %s", err))
	}

	go func() {
		ordSub := subscribers.NewOrderSubscriber(cache, storage, log)
		if err := ordSub.Start(ctx, *cfg); err != nil {
			log.Error(err.Error())
		}
	}()

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Get("/order/{id}", orders.GetOrder(ctx, log, cache, storage))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	log.Info("Starting http server", slog.String("addr", cfg.Address))
	if err := srv.ListenAndServe(); err != nil {
		log.Error("Error starting http server")
	}

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
