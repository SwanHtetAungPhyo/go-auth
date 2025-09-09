package auth

import (
	"context"
	"time"

	"github.com/SwanHtetAungPhyo/go-auth/framework/utils"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func (s Service) Me(userId uuid.UUID) (utils.GeneralResponse, error) {
	databaseCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	foundUser, err := s.Store.GetUserByID(databaseCtx, userId)
	if err != nil {
		log.Err(err).Str("userId", userId.String()).Msg("error getting user")
		return utils.GeneralResponse{}, err
	}
	return utils.GeneralResponse{
		Message: "User successfully authenticated for the user" + userId.String(),
		Data:    foundUser,
	}, nil
}
