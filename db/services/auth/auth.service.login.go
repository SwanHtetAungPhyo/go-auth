package auth

import (
	"context"
	"errors"
	"time"

	"github.com/SwanHtetAungPhyo/go-auth/framework"
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

func (s Service) Login(req *framework.LoginRequest) (framework.AuthResponse, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := s.Store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Error().Err(err).Str("email", req.Email).Msg("user not found")
		return framework.AuthResponse{}, errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(req.Password)); err != nil {
		log.Error().Err(err).Str("email", req.Email).Msg("wrong password")
		return framework.AuthResponse{}, fiber.ErrInternalServerError
	}

	//userInfo := framework.GoAuthUserInfo{
	//	UserId:   user.ID.String(),
	//	Email:    user.Email,
	//	Username: user.Name.String,
	//	CreateAt: user.CreatedAt.Time,
	//}

	if s.cfg.JwtAuth {
		token, err := s.generateToken(user.Email, user.RoleName, map[string]interface{}{})
		if err != nil {
			log.Error().Err(err).Msg("failed to generate token")
			return framework.AuthResponse{}, fiber.ErrInternalServerError
		}

		return framework.AuthResponse{
			//UserInfo:     userInfo,
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
		}, nil
	}

	return framework.AuthResponse{
		//UserInfo:     userInfo,
		AccessToken:  "",
		RefreshToken: "",
	}, nil
}
