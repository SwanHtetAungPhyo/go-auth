package goauth

import (
	"context"
	"os"

	"github.com/SwanHtetAungPhyo/go-auth/initialization"
	"github.com/SwanHtetAungPhyo/go-auth/third-party/email"
	"github.com/redis/go-redis/v9"
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
type EmailConfig struct {
	Type string
}

type Payment struct {
	Type string
}
type ThirdPartyConfig struct {
	Email   bool
	Payment bool
}

type Config struct {
	DNS       string
	JwtAuth   bool
	PestoAuth bool
	//EmailSend           bool
	Session             bool
	SessionStoreAsRedis bool
	TokenMetaData       map[string]string
	ThirdParty          ThirdPartyConfig
	GithubOauth         *GithubOauth
	GoogleOauth         *GoogleOauth
	EmailConfig         *EmailConfig
	Payment             *Payment
	redisClient         *redis.Client
	EmailService        *email.EmailService
	IsProduction        bool
}

type Option func(*Config)

func NewGoAuth(dsn string, config *Config, opts ...Option) *Config {
	if dsn == "" {
		log.Fatal().Msg("dsn is required for the goauth service")
	}

	cfg := config
	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.SessionStoreAsRedis {
		initialization.ValidateRedis()

		cfg.redisClient = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("GOAUTH_REDIS_ADDRESS"),
			Password: os.Getenv("GOAUTH_REDIS_PASSWORD"),
		})

		cfg.redisClient.Ping(context.Background())
	}
	// Database setup
	_ = initialization.Database(cfg.DNS, cfg.SessionStoreAsRedis)

	// Auth setup
	if cfg.JwtAuth {
		initialization.ValidateJwtAuth()
	}
	if cfg.PestoAuth {
		initialization.ValidatePestoAuth()
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
func WithRedisAsSessionStore(redis bool) Option {
	return func(cfg *Config) {
		cfg.SessionStoreAsRedis = redis
	}
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

func WithJwtAuth(jwtAuth bool) Option {
	return func(cfg *Config) {
		cfg.JwtAuth = jwtAuth
	}
}
