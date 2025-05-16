package config

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() error {
	return godotenv.Load()
}

func GetEnv(key, fallback string) string {
	if value := os.Getenv(key); value == "" {
		return fallback
	} else {
		return value
	}
}
