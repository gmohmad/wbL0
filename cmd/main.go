package main

import (
	"context"
	"log/slog"
	"os"

	"gihub.com/gmohmad/wb_l0/internal/config"
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

	ctx := context.Background()

	postgres, err := postgres.NewClient(ctx, cfg, 5, 3, log)

	if err != nil {
		log.Error(err.Error(), slog.Any("smth", postgres))
		os.Exit(1)
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
