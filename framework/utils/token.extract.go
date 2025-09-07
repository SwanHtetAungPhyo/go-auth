package utils

import (
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		jwtSecret := os.Getenv("JWT_SECRET")
		return jwtSecret, nil
	})
}

func ExtractToken(rHeader string) string {
	if strings.HasPrefix(rHeader, "Bearer ") {
		return strings.TrimPrefix(rHeader, "Bearer ")
	}
	return ""
}
