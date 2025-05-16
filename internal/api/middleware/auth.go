// package middleware

// import (
// 	"os"
// 	"strings"

// 	"github.com/gofiber/fiber/v3"
// 	"github.com/golang-jwt/jwt/v5"
// 	"github.com/rezatg/payment-gateway/internal/services/auth"
// 	"github.com/rezatg/payment-gateway/pkg/errors"
// )

// func JWTProtected() fiber.Handler {
// 	return func(c fiber.Ctx) error {
// 		authHeader := c.Get("Authorization")
// 		if authHeader == "" {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing authorization header"})
// 		}

// 		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
// 		if tokenStr == authHeader {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token format"})
// 		}

// 		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
// 			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid signing method")
// 			}

// 			return []byte(os.Getenv("JWT_SECRET")), nil
// 		})

// 		if err != nil {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 				"error": "Invalid token",
// 			})
// 		}

// 		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 			c.Locals("user", claims)
// 			return c.Next()
// 		}

// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
// 	}
// }

// // AuthMiddleware validates JWT tokens
// func AuthMiddleware(authService auth.AuthService) fiber.Handler {
// 	return func(c fiber.Ctx) error {
// 		token := c.Get("Authorization")
// 		if token == "" {
// 			return errors.HandleAPIError(c, errors.NewUnauthorizedError("Missing authorization token"))
// 		}

// 		// Remove "Bearer " prefix if present
// 		if len(token) > 7 && token[:7] == "Bearer " {
// 			token = token[7:]
// 		}

// 		userID, err := authService.ValidateToken(c.Context(), token)
// 		if err != nil {
// 			return errors.HandleAPIError(c, err)
// 		}

// 		c.Locals("userID", userID)
// 		return c.Next()
// 	}
// }

package middleware

import (
	"context"

	"github.com/gofiber/fiber/v3"
	"github.com/rezatg/payment-gateway/pkg/errors"
)

// AuthService interface for token validation
type AuthService interface {
	ValidateToken(ctx context.Context, token string) (string, error)
}

// AuthMiddleware validates JWT tokens
func AuthMiddleware(authService AuthService) fiber.Handler {
	return func(c fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			return errors.HandleAPIError(c, errors.NewUnauthorizedError("Missing authorization token"))
		}

		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		userID, err := authService.ValidateToken(c.Context(), token)
		if err != nil {
			return errors.HandleAPIError(c, err)
		}

		c.Locals("userID", userID)
		return c.Next()
	}
}
