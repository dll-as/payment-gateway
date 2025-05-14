package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func GetEnv(key, fallback string) string {
	if value := os.Getenv(key); value == "" {
		return fallback
	} else {
		return value
	}
}
