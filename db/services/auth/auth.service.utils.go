package auth

import (
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
