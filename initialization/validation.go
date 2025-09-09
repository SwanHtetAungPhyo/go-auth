package initialization

import (
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

func ValidateJwtAuth() {
	jwtSecret := GetEnv("GOAUTH_JWT_SECRET", "")
	if jwtSecret == "" {
		log.Fatal().Msg("goauth: jwt secret is required")
	}
	log.Info().Msg("goauth: using jwt authentication")
}

func ValidateGoogleOauth() {
	clientId := GetEnv("GOAUTH_GOOGLE_CLIENT_ID", "")
	clientSecret := GetEnv("GOAUTH_GOOGLE_CLIENT_SECRET", "")
	redirectURL := GetEnv("GOAUTH_GOOGLE_REDIRECT_URL", "")

	if clientId == "" || clientSecret == "" || redirectURL == "" {
		log.Fatal().Msg("goauth: missing Google OAuth credentials")
	}
	log.Info().Msg("goauth: using google authentication")
}

func ValidateGithubOauth() {
	clientId := GetEnv("GOAUTH_GITHUB_CLIENT_ID", "")
	clientSecret := GetEnv("GOAUTH_GITHUB_CLIENT_SECRET", "")
	redirectURL := GetEnv("GOAUTH_GITHUB_REDIRECT_URL", "")

	if clientId == "" || clientSecret == "" || redirectURL == "" {
		log.Fatal().Msg("goauth: missing Github OAuth credentials")
	}
	log.Info().Msg("goauth: using github authentication")
}
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func ValidatePestoAuth() {
	pestoKey := GetEnv("GOAUTH_PESTO_KEY", "")
	pestoSecret := GetEnv("GOAUTH_PESTO_SECRET", "")
	if pestoKey == "" || pestoSecret == "" {
		log.Fatal().Msg("goauth: missing Pesto OAuth credentials")
	}
	log.Info().Msg("goauth: using Pesto authentication")
}
func GetEnvDuration(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	d, err := time.ParseDuration(value)
	if err != nil {
		log.Printf("invalid duration for %s: %s, using fallback %s", key, value, fallback)
		return fallback
	}
	return d
}

func ValidateRedis() {
	redisHOST := GetEnv("GOAUTH_REDIS_HOST", "localhost")
	redisPORT := GetEnv("GOAUTH_REDIS_PORT", "6379")
	if redisHOST == "" || redisPORT == "" {
		log.Fatal().Msg("goauth: missing Redis host and port")
	}
	log.Info().Msg("goauth: using Redis authentication")
}
