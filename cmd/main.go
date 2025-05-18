package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/rezatg/payment-gateway/config"
	"github.com/rezatg/payment-gateway/internal/api/auth"
	"github.com/rezatg/payment-gateway/internal/api/invoice"
	"github.com/rezatg/payment-gateway/internal/api/payment"
	"github.com/rezatg/payment-gateway/internal/api/router"
	"github.com/rezatg/payment-gateway/internal/database"
	"github.com/rezatg/payment-gateway/internal/repository"
	authService "github.com/rezatg/payment-gateway/internal/services/auth"
	"github.com/rezatg/payment-gateway/internal/services/blockchain"
	invoiceService "github.com/rezatg/payment-gateway/internal/services/invoice"
	"github.com/rezatg/payment-gateway/pkg/logger"
)

// initApp initializes the Fiber app and dependencies
func initApp() (*fiber.App, error) {
	// Load configuration
	var err error = config.LoadConfig()
	if err != nil {
		return nil, err
	}

	// Initialize logger
	if err = logger.InitLogger(); err != nil {
		log.Fatalf("failed to initialize zap logger: %v", err)
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Payment Gateway",
		// ErrorHandler: func(c fiber.Ctx, err error) error {
		// 	logger.Error("Request error", err, "path", c.Path(), "method", c.Method())
		// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		// 		"error": "Internal Server Error",
		// 		"code":  fiber.StatusInternalServerError,
		// 	})
		// },
	})

	// // Add CORS middleware
	// app.Use(cors.New(cors.Config{
	// 	AllowOrigins: []string{"http://localhost:3000", "http://localhost:3001"},
	// 	AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	// 	AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// }))

	// Run database migrations
	if err = database.RunMigrations(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	// Connect to database
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	// Create repositories
	userRepo := repository.NewUserRepository(db)
	invoiceRepo := repository.NewInvoiceRepository(db)

	// Generate or load encryption key (32 bytes for AES-256)
	encryptionKey := make([]byte, 32)
	if key := os.Getenv("ENCRYPTION_KEY"); key != "" {
		keyBytes, err := hex.DecodeString(key)
		if err != nil || len(keyBytes) != 32 {
			return nil, fmt.Errorf("invalid ENCRYPTION_KEY: must be 32-byte hex string")
		}
		encryptionKey = keyBytes
	} else {
		if _, err := rand.Read(encryptionKey); err != nil {
			return nil, fmt.Errorf("failed to generate encryption key: %w", err)
		}
		logger.Info("Generated encryption key", "key", hex.EncodeToString(encryptionKey))
	}

	// Create blockchain services
	blockchainFactory := blockchain.NewFactory(encryptionKey)

	// Create services
	authService := authService.New(userRepo)
	invoiceService := invoiceService.New(invoiceRepo, authService)

	// Create handlers
	authHandler := auth.New(authService)
	paymentHandler := payment.New(blockchainFactory, authService)
	invoiceHandler := invoice.New(invoiceService)

	// Register API routes
	router.RegisterAPIRoutes(app, router.Dependencies{
		AuthHandler:    authHandler,
		PaymentHandler: paymentHandler,
		InvoiceHandler: invoiceHandler,
	})

	return app, nil
}

func main() {
	app, err := initApp()
	if err != nil {
		logger.Fatal("Failed to initialize app", err)
	}
	defer logger.Logger.Sync()

	// Start server
	port := config.GetEnv("PORT", "8080")
	logger.Info("Starting server", "port", port)
	if err := app.Listen(":" + port); err != nil {
		logger.Fatal("Failed to start server", err, "port", port)
	}
}
