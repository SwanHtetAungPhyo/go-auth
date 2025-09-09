package auth

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/SwanHtetAungPhyo/go-auth/framework/utils"
	"github.com/SwanHtetAungPhyo/go-auth/initialization"
)

func (s Service) generateToken(userId string, role string, meta map[string]interface{}) (*utils.TokenContextContainer, error) {
	duration := initialization.GetEnvDuration("GOAUTH_TOKEN_DURATION", 3*time.Minute)

	claims := utils.Claims{
		UserID:   userId,
		Role:     role,
		MetaData: meta,
	}

	if s.cfg.JwtAuth {
		return utils.GenerateToken(claims, utils.JWT, duration)
		//} else if s.cfg.PestoAuth {
		//	return utils.GenerateToken(claims, utils.PASETO, duration)
		//}
	}

	return nil, nil
}

func (s Service) generateSixDigit() string {
	rand.NewSource(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}
func (s Service) buildDeepLink(email, code string) string {
	frontendURL := os.Getenv("GOAUTH_FRONTEND_URL") // e.g. https://myapp.com/verify
	return fmt.Sprintf("%s/verify?email=%s&code=%s", frontendURL, email, code)
}
