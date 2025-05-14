package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rezatg/payment-gateway/internal/api/auth"
	"github.com/rezatg/payment-gateway/internal/api/payment"
	services "github.com/rezatg/payment-gateway/internal/services/auth"
)

func RegisterAPIRoutes(app *fiber.App, authService services.AuthService) {
	api := app.Group("/api")

	// Register module routes
	auth.RegisterRoutes(api.Group("/auth"), authService)
	payment.RegisterRoutes(api.Group("/payment"))
}
