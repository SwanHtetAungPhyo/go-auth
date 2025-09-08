package initialization

import (
	"os"

	"github.com/rs/zerolog/log"
)

func ValidateJwtAuth() {
	jwtSecret := getEnv("GOAUTH_JWT_SECRET", "")
	if jwtSecret == "" {
		log.Fatal().Msg("goauth: jwt secret is required")
	}
	log.Info().Msg("goauth: using jwt authentication")
}

func ValidateGoogleOauth() {
	clientId := getEnv("GOAUTH_GOOGLE_CLIENT_ID", "")
	clientSecret := getEnv("GOAUTH_GOOGLE_CLIENT_SECRET", "")
	redirectURL := getEnv("GOAUTH_GOOGLE_REDIRECT_URL", "")

	if clientId == "" || clientSecret == "" || redirectURL == "" {
		log.Fatal().Msg("goauth: missing Google OAuth credentials")
	}
	log.Info().Msg("goauth: using google authentication")
}

func ValidateGithubOauth() {
	clientId := getEnv("GOAUTH_GITHUB_CLIENT_ID", "")
	clientSecret := getEnv("GOAUTH_GITHUB_CLIENT_SECRET", "")
	redirectURL := getEnv("GOAUTH_GITHUB_REDIRECT_URL", "")

	if clientId == "" || clientSecret == "" || redirectURL == "" {
		log.Fatal().Msg("goauth: missing Github OAuth credentials")
	}
	log.Info().Msg("goauth: using github authentication")
}
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func ValidatePestoAuth() {
	pestoKey := getEnv("GOAUTH_PESTO_KEY", "")
	pestoSecret := getEnv("GOAUTH_PESTO_SECRET", "")
	if pestoKey == "" || pestoSecret == "" {
		log.Fatal().Msg("goauth: missing Pesto OAuth credentials")
	}
	log.Info().Msg("goauth: using Pesto authentication")
}
