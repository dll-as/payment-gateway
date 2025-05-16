package middleware

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/rezatg/payment-gateway/pkg/logger"
)

// LoggerMiddleware logs HTTP requests with structured logging
func LoggerMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)

		logger.Info("HTTP request",
			"method", c.Method(),
			"path", c.Path(),
			"status", c.Response().StatusCode(),
			"duration_ms", duration.Milliseconds(),
			"ip", c.IP(),
		)

		return err
	}
}
