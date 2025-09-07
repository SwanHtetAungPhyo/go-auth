package goauth

import (
	"github.com/SwanHtetAungPhyo/go-auth/initialization"
	"github.com/rs/zerolog/log"
)

type GoogleOauth struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	CallBackURL  string
}

type GithubOauth struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	CallBackURL  string
}

type Config struct {
	DNS         string
	JwtAuth     bool
	GithubOauth *GithubOauth
	GoogleOauth *GoogleOauth
}

type Option func(*Config)

func WithJwtAuth(jwtAuth bool) Option {
	return func(cfg *Config) {
		cfg.JwtAuth = jwtAuth
	}
}

func NewGoAuth(dsn string, opts ...Option) *Config {
	if dsn == "" {
		log.Fatal().Msg("dsn is required for the goauth service")
	}

	cfg := &Config{DNS: dsn, JwtAuth: false}
	for _, opt := range opts {
		opt(cfg)
	}

	// Database setup
	_ = initialization.Database(cfg.DNS)

	// Auth setup
	if cfg.JwtAuth {
		initialization.ValidateJwtAuth()
	}
	if cfg.GoogleOauth != nil {
		initialization.ValidateGoogleOauth()
	}
	if cfg.GithubOauth != nil {
		initialization.ValidateGithubOauth()
	}

	log.Info().Msg("Prepare necessary environment variables")
	return cfg
}

func WithGithubOauth(clientID, clientSecret, redirect, callback string) Option {
	return func(c *Config) {
		c.GithubOauth = &GithubOauth{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirect,
			CallBackURL:  callback,
		}
	}
}

func WithGoogleOauth(clientID, clientSecret, redirect, callback string) Option {
	return func(c *Config) {
		c.GoogleOauth = &GoogleOauth{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirect,
			CallBackURL:  callback,
		}
	}
}
