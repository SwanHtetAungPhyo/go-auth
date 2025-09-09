package framework

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type (
	GoAuthUserInfo struct {
		UserId   string    `json:"user_id"`
		Email    string    `json:"email"`
		Username string    `json:"username"`
		CreateAt time.Time `json:"create_at"`
	}
	RegisterRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"user_name,omitempty"`
		RoleName string `json:"role_name,omitempty"`
		Image    string `json:"image,omitempty"`
	}
	LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	AuthResponse struct {
		UserInfo     GoAuthUserInfo `json:"user_info,omitempty"`
		AccessToken  string         `json:"access_token"`
		RefreshToken string         `json:"refresh_token"`
		SessionToken string         `json:"session_token,omitempty"`
	}
	RegisterResponse struct{}
)

var validate = validator.New()

func ValidateStruct(v interface{}) error {
	return validate.Struct(v)
}
