package auth

import (
	"context"
	"sync"
	"time"

	db "github.com/SwanHtetAungPhyo/go-auth/db/sqlc"
	"github.com/SwanHtetAungPhyo/go-auth/framework"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

func (s Service) Register(req *framework.RegisterRequest) (framework.AuthResponse, error) {
	databaseCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Err(err).Msg("failed to hash password")
		return framework.AuthResponse{}, err
	}

	user, err := s.Store.GoAuthRegister(databaseCtx, db.GoAuthRegisterParams{
		Email:        req.Email,
		HashPassword: string(hash),
		RoleName:     req.RoleName,
		Name:         pgtype.Text{String: req.Name},
	})
	if err != nil {
		log.Err(err).Msg("failed to create user")
		return framework.AuthResponse{}, err
	}

	if s.cfg.IsProduction {
		var wg sync.WaitGroup
		wg.Go(func() {
			//s.emailType.SendNotificationEmail()

		})
	} else {
		err := s.Store.UpdateUserEmailVerified(databaseCtx, db.UpdateUserEmailVerifiedParams{
			ID:            user.ID,
			EmailVerified: pgtype.Bool{Bool: true},
		})
		if err != nil {
			log.Err(err).Str("GOAUTH", "register_service").Msg("failed to update user email verified")
			return framework.AuthResponse{}, err
		}
	}
	if s.cfg.TokenMetaData != nil {

	}
	token, err := s.generateToken(user.ID.String(), req.RoleName, nil)
	if err != nil {
		log.Err(err).Msg("failed to generate token")
		return framework.AuthResponse{}, err
	}

	return framework.AuthResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}
