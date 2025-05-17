package database

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rezatg/payment-gateway/pkg/logger"
)

// RunMigrations executes database migrations
func RunMigrations() error {
	m, err := migrate.New("file://migrations", "postgres://reza:rezatg15%23%23%23@localhost:5432/tg?sslmode=disable")
	if err != nil {
		return fmt.Errorf("failed to initialize migration: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	logger.Info("Database migrations applied successfully")
	return nil
}
