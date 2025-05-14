package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func ConnectPostgres() *sql.DB {
	connStr := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("cannot ping database: %v", err)
	}

	return db
}
