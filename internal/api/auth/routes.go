package auth

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rezatg/payment-gateway/internal/services/auth"
)

func RegisterRoutes(router fiber.Router, service auth.AuthService) {
	handler := NewAuthHandler(service)

	router.Post("/login", handler.Login)
	router.Post("/register", handler.Register)
}
