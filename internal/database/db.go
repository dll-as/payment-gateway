package database

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	// "github.com/rezatg/payment-gateway/config"
	"github.com/rezatg/payment-gateway/pkg/logger"
)

// Connect establishes a database connection with optimized settings
func Connect() (*sql.DB, error) {
	// connStr := config.GetEnv("DATABASE_URL", "postgres://user:password@localhost:5432/payment_gateway?sslmode=disable")
	db, err := sql.Open("postgres", "postgres://reza:rezatg15%23%23%23@localhost:5432/tg?sslmode=disable")
	if err != nil {
		logger.Error("Failed to open database connection", err)
		return nil, err
	}

	// Optimize connection pool
	db.SetMaxOpenConns(25) // Adjust based on load
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		logger.Error("Failed to ping database", err)
		return nil, err
	}

	logger.Info("Database connected successfully")
	return db, nil
}
