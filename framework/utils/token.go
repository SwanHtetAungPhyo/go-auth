package utils

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/o1egl/paseto"
	"github.com/rs/zerolog/log"
)

// ----------------------
// TOKEN GENERATION
// ----------------------

// GenerateToken generates JWT or PASETO token based on type
func GenerateToken(claims Claims, tokenType TokenType, duration time.Duration) (*TokenContextContainer, error) {
	switch tokenType {
	case JWT:
		return generateJWT(claims, duration)
	//case PASETO:
	//	return generatePaseto(claims, duration)
	default:
		return nil, errors.New("unsupported token type")
	}
}

// ----------------------
// JWT GENERATION
// ----------------------

func generateJWT(claims Claims, duration time.Duration) (*TokenContextContainer, error) {
	jwtSecret := os.Getenv("GOAUTH_JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET not set in environment")
	}

	accessDuration := time.Now().Add(duration).Unix()
	accesstoken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		Type:   JWT_ACCESS_TOKEN,
		UserId: claims.UserID,
		Role:   claims.Role,
		Exp:    accessDuration,
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		Type:   JWT_REFRESH_TOKEN,
		UserId: claims.UserID,
		Role:   claims.Role,
		Exp:    time.Now().Add(24 * time.Hour).Unix(),
	})

	signedAccessToken, err := accesstoken.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Err(err).Msg("error signing access token")
		return nil, err
	}
	signedRefreshToken, err := refreshToken.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Err(err).Msg("error signing refresh token")
		return nil, err
	}
	return &TokenContextContainer{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
	}, nil
}

// ----------------------
// PASETO GENERATION
// ----------------------

//func generatePaseto(claims Claims, duration time.Duration) (*TokenContextContainer, error) {
//	pasetoKey := os.Getenv("GOAUTH_PASETO_KEY")
//	if pasetoKey == "" {
//		return nil, errors.New("PASETO_KEY not set in environment")
//	}
//
//	v2 := paseto.NewV2()
//	jsonToken := paseto.JSONToken{
//		Audience:   "",
//		Issuer:     "",
//		Jti:        "",
//		Subject:    "",
//		Expiration: time.Now().Add(duration),
//		IssuedAt:   time.Time{},
//		NotBefore:  time.Time{},
//	}
//
//	return v2.Encrypt([]byte(pasetoKey), jsonToken, "")
//}

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
	jwtSecret := os.Getenv("GOAUTH_JWT_SECRET")
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
	pasetoKey := os.Getenv("GOAUTH_PASETO_KEY")
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
