package framework

import "github.com/go-playground/validator/v10"

type (
	RegisterRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"user_name,omitempty"`
		Image    string `json:"image,omitempty"`
	}
	LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	AuthResponse struct {
		AccessToken  string `json:"access_token"`
		UserClaims   string `json:"user_claims"`
		RefreshToken string `json:"refresh_token"`
	}
	RegisterResponse struct{}
)

var validate = validator.New()

func ValidateStruct(v interface{}) error {
	return validate.Struct(v)
}
