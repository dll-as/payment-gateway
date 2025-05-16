package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

// RateLimitMiddleware applies rate limiting to API requests
func RateLimitMiddleware() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        100, // Max requests per IP
		Expiration: 60,  // Seconds
		KeyGenerator: func(c fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(map[string]string{
				"error": "Too many requests",
			})
		},
	})
}
