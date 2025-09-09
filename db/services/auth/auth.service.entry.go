package auth

import (
	goauth "github.com/SwanHtetAungPhyo/go-auth"
	db "github.com/SwanHtetAungPhyo/go-auth/db/sqlc"
	"github.com/SwanHtetAungPhyo/go-auth/framework"
	"github.com/SwanHtetAungPhyo/go-auth/framework/utils"
	"github.com/SwanHtetAungPhyo/go-auth/third-party/email"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type AuthService interface {
	Login(req *framework.LoginRequest) (framework.AuthResponse, error)
	Register(req *framework.RegisterRequest) (framework.AuthResponse, error)
	Me(userId uuid.UUID) (utils.GeneralResponse, error)
}

type Service struct {
	Store     *db.Store
	cfg       goauth.Config
	client    redis.Client
	emailType *email.EmailService
}

type Option func(*Service)

func NewAuthService(conn *pgxpool.Pool, cfg goauth.Config, opts ...Option) Service {
	store := db.NewStore(conn)
	service := Service{
		Store:     store,
		cfg:       cfg,
		emailType: cfg.EmailService,
	}
	for _, opt := range opts {
		opt(&service)
	}
	return service
}
func WithRedisClient(client redis.Client) Option {
	return func(s *Service) {
		s.client = client
	}
}

var _ AuthService = (*Service)(nil)
