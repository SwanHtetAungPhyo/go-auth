package middleware

import (
	"github.com/SwanHtetAungPhyo/go-auth/framework/utils"
	"github.com/gofiber/fiber/v3"

	"net/http"
)

func FiberJWTMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		tokenString := utils.ExtractToken(authHeader)

		if tokenString == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing or invalid token",
			})
		}

		token, err := utils.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or expired token",
			})
		}
		c.Locals("user", token.Claims)
		return c.Next()
	}
}
