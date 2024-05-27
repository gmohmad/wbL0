package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

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
	var wg sync.WaitGroup

	cfg := config.MustLoad()

	log := config.SetupLogger(cfg.Env)

	log.Info("starting the service", slog.String("env", cfg.Env))

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

	wg.Add(1)
	go func() {
		defer wg.Done()
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

	router.Mount("/", http.FileServer(http.Dir("./static")))

	router.Get("/order/{id}", orders.GetOrder(ctx, log, cache, storage))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	go func() {
		log.Info("Starting http server", slog.String("addr", cfg.Address))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Error starting http server: ", err)
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-stopChan

	cancel()

	wg.Wait()

	log.Info("Closing postgres connection pool")
	postgresClient.Close()

	log.Info("Shutting down the server...")

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Error shutting down the server: ", err)
	}

	log.Info("Server gracefully stopped")
}
