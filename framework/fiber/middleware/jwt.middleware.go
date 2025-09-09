package middleware

import (
	"github.com/SwanHtetAungPhyo/go-auth/framework/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"

	"net/http"
)

// FiberAuthMiddleware returns a Fiber middleware that validates JWT or PASETO tokens
// and stores the user_id and full claims in locals.
func (m *Maker) FiberAuthMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		tokenString := utils.ExtractToken(authHeader)

		if tokenString == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing or invalid token",
			})
		}

		var claims map[string]interface{}
		var userID string

		if m.cfg.JwtAuth {
			token, err := utils.ValidateToken(tokenString, utils.JWT)
			if err != nil || !token.(*jwt.Token).Valid {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
					"error": "invalid or expired JWT token",
				})
			}

			jwtToken := token.(*jwt.Token)
			if mc, ok := jwtToken.Claims.(jwt.MapClaims); ok {
				claims = mc
				if id, exists := mc["user_id"].(string); exists {
					userID = id
				}
			}

		} else if m.cfg.PestoAuth {
			cClaims, err := utils.ValidateToken(tokenString, utils.PASETO)
			if err != nil {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
					"error": "invalid or expired PASETO token",
				})
			}
			claims = cClaims.(map[string]interface{})
			if id, exists := claims["user_id"].(string); exists {
				userID = id
			}
		}

		c.Locals("user_id", userID)
		c.Locals("user_claims", claims)

		return c.Next()
	}
}
