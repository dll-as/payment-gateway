package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

// func generateRefreshToken(days int) (string, error) {
// 	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"exp": time.Now().AddDate(0, 0, days).Unix(),
// 	})

// 	return refresh.SignedString(SecretKey)
// }
