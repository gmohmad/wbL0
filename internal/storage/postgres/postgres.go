package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"gihub.com/gmohmad/wb_l0/internal/config"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func NewClient(ctx context.Context, cfg *config.Config, maxRetries int, delay time.Duration, log *slog.Logger) (pool Client, err error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	for i := range maxRetries {

		attemptCtx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()

		pool, err = pgxpool.New(attemptCtx, dsn)

		if err == nil {
			return pool, nil
		}

		log.Debug(fmt.Sprintf("Failed to connect to the database (attempt: %d/%d): %v", i+1, maxRetries, err))
		time.Sleep(time.Second * delay)

	}
	return nil, fmt.Errorf("Failed to connect to the database after %d attempts: %w", maxRetries, err)

}
