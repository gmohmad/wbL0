package postgres

import (
	"fmt"
	"log/slog"

	"gihub.com/gmohmad/wb_l0/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(cfg *config.DB, log *slog.Logger) error {

	db := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s", 
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
	)

	source := fmt.Sprintf("file://%s", cfg.MigrationsPath)

	m, err := migrate.New(source, db)

	if err != nil {
		return fmt.Errorf("Error while instantiating migrate.Migrate: %w", err)
	}

	log.Info("Applying UP migrations...")
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Info("No change in migrations to apply")
			return nil
		}
		return fmt.Errorf("Error while applying UP migrations: %w", err)
	}

	return nil
}
