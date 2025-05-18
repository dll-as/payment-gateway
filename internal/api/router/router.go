package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rezatg/payment-gateway/internal/api/auth"
	"github.com/rezatg/payment-gateway/internal/api/invoice"
	"github.com/rezatg/payment-gateway/internal/api/middleware"
	"github.com/rezatg/payment-gateway/internal/api/payment"
)

// Dependencies holds handler dependencies
type Dependencies struct {
	AuthHandler    *auth.AuthHandler
	PaymentHandler *payment.PaymentHandler
	InvoiceHandler *invoice.InvoiceHandler
}

// RegisterAPIRoutes registers all API routes
func RegisterAPIRoutes(app *fiber.App, deps Dependencies) {
	// Global middlewares
	app.Use(middleware.LoggerMiddleware())
	// app.Use(middleware.RateLimitMiddleware())

	// API group
	api := app.Group("/api")

	// Register module routes
	auth.RegisterRoutes(api, deps.AuthHandler)
	payment.RegisterRoutes(api, deps.PaymentHandler)
	invoice.RegisterRoutes(api, deps.InvoiceHandler)
}
