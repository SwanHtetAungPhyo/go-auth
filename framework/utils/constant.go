package utils

// TokenType represents the type of token

type (
	TokenType             string
	TokenContextContainer struct {
		AccessToken  string
		RefreshToken string
	}

	// Claims represents the common claims for both JWT and PASETO
	Claims struct {
		UserID   string
		Role     string
		MetaData map[string]interface{}
	}
	GeneralResponse struct {
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
		Error   error       `json:"error,omitempty"`
	}
)

const (
	Type              string    = "type"
	UserId            string    = "user_id"
	Role              string    = "role"
	Exp               string    = "exp"
	JWT_ACCESS_TOKEN  TokenType = "access_token"
	JWT_REFRESH_TOKEN TokenType = "refresh_token"
	JWT               TokenType = "jwt"
	PASETO            TokenType = "paseto"
)
