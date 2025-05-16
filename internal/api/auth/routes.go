package auth

import (
	"github.com/gofiber/fiber/v3"
)

// RegisterRoutes registers authentication-related routes
func RegisterRoutes(router fiber.Router, handler *AuthHandler) {
	authGroup := router.Group("/auth")
	authGroup.Post("/login", handler.Login)
	authGroup.Post("/register", handler.Register)
}
