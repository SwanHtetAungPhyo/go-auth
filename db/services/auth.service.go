package services

import (
	"context"
	"time"

	db "github.com/SwanHtetAungPhyo/go-auth/db/sqlc"
	"github.com/SwanHtetAungPhyo/go-auth/framework"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(req *framework.LoginRequest) (error, framework.AuthResponse)
	Register(req *framework.RegisterRequest) (error, framework.AuthResponse)
}

type Service struct {
	Store *db.Store
}

func NewAuthService(conn *pgxpool.Pool) Service {
	store := db.NewStore(conn)
	return Service{
		Store: store,
	}
}

func (s Service) Login(req *framework.LoginRequest) (error, framework.AuthResponse) {
	databaseCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	existingGoAuthUser, err := s.Store.GetUserByEmail(databaseCtx, req.Email)
	if err != nil {
		log.Err(err).Str("GOAUTH", "goauth_auth_service").Str("email", req.Email).Msg("can't get user by email")
		return err, framework.AuthResponse{}
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingGoAuthUser.HashPassword), []byte(req.Password))
	if err != nil {
		log.Err(err).Str("GOAUTH", "goauth_auth_service").Str("email", req.Email).Msg("wrong password")
		return err, framework.AuthResponse{}
	}

	return nil, framework.AuthResponse{
		AccessToken:  "",
		UserClaims:   "",
		RefreshToken: "",
	}
}
