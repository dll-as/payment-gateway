package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"

	"github.com/rezatg/payment-gateway/config"
	"github.com/rezatg/payment-gateway/internal/api/router"
	"github.com/rezatg/payment-gateway/internal/core/database"
	"github.com/rezatg/payment-gateway/internal/repository"
	"github.com/rezatg/payment-gateway/internal/services/auth"
)

func main() {
	// Load .env config
	config.LoadConfig()

	// Create Fiber instance
	app := fiber.New()

	// Global middlewares
	app.Use(logger.New())

	// Connect to the database
	dbConn := database.ConnectPostgres()

	// Create UserRepository instance
	userRepo := repository.NewUserRepository(dbConn)

	// Create AuthService instance
	service := auth.NewAuthService(userRepo)

	// Register API routes
	router.RegisterAPIRoutes(app, service)

	// Start server
	port := config.GetEnv("PORT", "3000")
	log.Printf("🚀 Server running on port %s...", port)
	log.Fatal(app.Listen(":" + port))
}
