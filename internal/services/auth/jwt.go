package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rezatg/payment-gateway/config"
	"github.com/rezatg/payment-gateway/pkg/errors"
)

var SecretKey = []byte(os.Getenv("JWT_SECRET"))

func generateJWTToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretKey)
}

// validateJWTToken validates a JWT token and returns the user ID
func validateJWTToken(token string) (string, error) {
	secret := config.GetEnv("JWT_SECRET", "your-secret-key")
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.NewUnauthorizedError("Unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil || !parsedToken.Valid {
		return "", errors.NewUnauthorizedError("Invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return "", errors.NewUnauthorizedError("Invalid token claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.NewUnauthorizedError("Invalid user ID in token")
	}

	return userID, nil
}

// func generateRefreshToken(days int) (string, error) {
// 	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"exp": time.Now().AddDate(0, 0, days).Unix(),
// 	})

// 	return refresh.SignedString(SecretKey)
// }
