package utils

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/o1egl/paseto"
)

// TokenType represents the type of token
type TokenType string

const (
	JWT    TokenType = "jwt"
	PASETO TokenType = "paseto"
)

// Claims represents the common claims for both JWT and PASETO
type Claims struct {
	UserID string
	Email  string
	Role   string
}

// ----------------------
// TOKEN GENERATION
// ----------------------

// GenerateToken generates JWT or PASETO token based on type
func GenerateToken(claims Claims, tokenType TokenType, duration time.Duration) (string, error) {
	switch tokenType {
	case JWT:
		return generateJWT(claims, duration)
	case PASETO:
		return generatePaseto(claims, duration)
	default:
		return "", errors.New("unsupported token type")
	}
}

// ----------------------
// JWT GENERATION
// ----------------------

func generateJWT(claims Claims, duration time.Duration) (string, error) {
	jwtSecret := os.Getenv("GOAUTH_JWT_SECRET")
	if jwtSecret == "" {
		return "", errors.New("JWT_SECRET not set in environment")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": claims.UserID,
		"email":   claims.Email,
		"role":    claims.Role,
		"exp":     time.Now().Add(duration).Unix(),
	})

	return token.SignedString([]byte(jwtSecret))
}

// ----------------------
// PASETO GENERATION
// ----------------------

func generatePaseto(claims Claims, duration time.Duration) (string, error) {
	pasetoKey := os.Getenv("GOAUTH_PASETO_KEY")
	if pasetoKey == "" {
		return "", errors.New("PASETO_KEY not set in environment")
	}

	v2 := paseto.NewV2()
	jsonToken := paseto.JSONToken{
		Audience:   "",
		Issuer:     "",
		Jti:        "",
		Subject:    "",
		Expiration: time.Now().Add(duration),
		IssuedAt:   time.Time{},
		NotBefore:  time.Time{},
	}

	return v2.Encrypt([]byte(pasetoKey), jsonToken, "")
}

// ----------------------
// TOKEN VALIDATION
// ----------------------

func ValidateToken(tokenString string, tokenType TokenType) (interface{}, error) {
	switch tokenType {
	case JWT:
		return validateJWT(tokenString)
	case PASETO:
		return validatePaseto(tokenString)
	default:
		return nil, errors.New("unsupported token type")
	}
}

// ----------------------
// JWT VALIDATION
// ----------------------

func validateJWT(tokenString string) (*jwt.Token, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET not set in environment")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

// ----------------------
// PASETO VALIDATION
// ----------------------

func validatePaseto(tokenString string) (*paseto.JSONToken, error) {
	pasetoKey := os.Getenv("PASETO_KEY")
	if pasetoKey == "" {
		return nil, errors.New("PASETO_KEY not set in environment")
	}

	var jsonToken paseto.JSONToken
	var footer string
	v2 := paseto.NewV2()

	if err := v2.Decrypt(tokenString, []byte(pasetoKey), &jsonToken, &footer); err != nil {
		return nil, err
	}

	if jsonToken.Expiration.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	return &jsonToken, nil
}

// ----------------------
// TOKEN EXTRACTION
// ----------------------

func ExtractToken(rHeader string) string {
	if strings.HasPrefix(rHeader, "Bearer ") {
		return strings.TrimPrefix(rHeader, "Bearer ")
	}
	return ""
}
