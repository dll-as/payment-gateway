package payment

import (
	"github.com/gofiber/fiber/v3"
)

// RegisterRoutes registers payment-related routes
func RegisterRoutes(router fiber.Router, handler *PaymentHandler) {
	paymentGroup := router.Group("/payment")
	// paymentGroup.Use(middleware.AuthMiddleware(handler.service.GetAuthService()))
	paymentGroup.Post("/check-balance", handler.CheckBalance)
}
