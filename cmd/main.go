package main

import (
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
	if err := config.LoadConfig(); err != nil {
		return nil, err
	}

	// Initialize logger
	logger.Init()

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Payment Gateway",
		ErrorHandler: func(c fiber.Ctx, err error) error {
			logger.Error("Request error", err, "path", c.Path(), "method", c.Method())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal Server Error",
				"code":  fiber.StatusInternalServerError,
			})
		},
	})

	// Connect to database
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	// Create repositories
	userRepo := repository.NewUserRepository(db)
	invoiceRepo := repository.NewInvoiceRepository(db)

	// Create blockchain services
	blockchainFactory := blockchain.NewFactory()

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

	// Start server
	port := config.GetEnv("PORT", "8080")
	logger.Info("Starting server", "port", port)
	if err := app.Listen(":" + port); err != nil {
		logger.Fatal("Failed to start server", err, "port", port)
	}
}
