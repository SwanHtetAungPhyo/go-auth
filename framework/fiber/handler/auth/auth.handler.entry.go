package auth

import (
	goauth "github.com/SwanHtetAungPhyo/go-auth"
	"github.com/SwanHtetAungPhyo/go-auth/db/services/auth"
	"github.com/SwanHtetAungPhyo/go-auth/framework"
	"github.com/SwanHtetAungPhyo/go-auth/third-party/email"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type Option func(authFiber *GoAuthFiber)
type GoAuthFiber struct {
	connPool     *pgxpool.Pool
	cfg          goauth.Config
	srv          auth.AuthService
	redis        *redis.Client
	session      *session.Session
	emailManager email.EmailManager
}

func NewGoAuthFiber(connPool *pgxpool.Pool, cfg goauth.Config, emailManager email.EmailManager, opts ...Option) *GoAuthFiber {
	if connPool == nil {
		log.Fatal().Msg("connPool is nil in NewGoAuthFiber handler")
		return nil
	}

	srv := auth.NewAuthService(connPool, cfg)
	goauthFiber := &GoAuthFiber{
		connPool:     connPool,
		cfg:          cfg,
		srv:          srv,
		emailManager: emailManager,
	}
	for _, opt := range opts {
		opt(goauthFiber)
	}

	if !cfg.Session {
		return goauthFiber
	}
	if !cfg.SessionStoreAsRedis {
		if goauthFiber.session == nil {
			log.Fatal().Msg("fiber session is needed to be provided in NewGoAuthFiber")
			return nil
		}
	} else {
		if goauthFiber.redis == nil {
			log.Fatal().Msg("goauthFiber.redis is nil in NewGoAuthFiber")
			return nil
		}
	}
	return goauthFiber
}

func WithFiberSessionStore(sess *session.Session) Option {
	return func(authFiber *GoAuthFiber) {
		authFiber.session = sess
	}
}
func WithRedisClient(client *redis.Client) Option {
	return func(authFiber *GoAuthFiber) {
		authFiber.redis = client

	}
}

var _ framework.Fiber = (*GoAuthFiber)(nil)
