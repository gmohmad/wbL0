package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"gihub.com/gmohmad/wb_l0/internal/config"
	"gihub.com/gmohmad/wb_l0/internal/http/handlers/orders"
	"gihub.com/gmohmad/wb_l0/internal/nats/subscribers"
	"gihub.com/gmohmad/wb_l0/internal/storage"
	"gihub.com/gmohmad/wb_l0/internal/storage/cache"
	"gihub.com/gmohmad/wb_l0/internal/storage/postgres"
	"gihub.com/gmohmad/wb_l0/internal/utils"
)

func main() {
	cfg := config.MustLoad()

	log := config.SetupLogger(cfg.Env)

	log.Info("starting the app", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	postgresClient, err := postgres.NewClient(ctx, &cfg.DB, 5, 3, log)
	if err != nil {
		utils.LogFatal(log, err)
	}

	if err := postgres.Migrate(&cfg.DB, log); err != nil {
		utils.LogFatal(log, err)
	}

	storage := storage.NewStorage(postgresClient)
	cache := cache.NewCache()

	if err := cache.FillUpCache(ctx, storage); err != nil {
		log.Info(fmt.Sprintf("Error filling cache from db: %s", err))
	}

	go func() {
		ordSub := subscribers.NewOrderSubscriber(cache, storage, log)
		if err := ordSub.Start(ctx, &cfg.Nats); err != nil {
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
