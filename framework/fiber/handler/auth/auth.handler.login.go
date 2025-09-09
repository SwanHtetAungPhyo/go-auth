package auth

import (
	"os"
	"time"

	"github.com/SwanHtetAungPhyo/go-auth/framework"
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

func (g *GoAuthFiber) Login(c fiber.Ctx) error {
	var req framework.LoginRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := framework.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	authResponse, err := g.srv.Login(&req)
	if err != nil {
		log.Error().Err(err).Msg("Login failed")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid credentials",
		})
	}

	isProduction := os.Getenv("ENVIRONMENT") == "production"
	refreshCookies := fiber.Cookie{
		Name:        "refresh_token",
		Value:       authResponse.RefreshToken,
		Path:        "/",                                 // Set to root path so cookie is available site-wide
		Domain:      "",                                  // Leave empty for same-domain, or set specific domain in production
		Expires:     time.Now().Add(time.Hour * 24 * 7),  // 7 days for refresh token
		MaxAge:      int((time.Hour * 24 * 7).Seconds()), // 7 days in seconds
		Secure:      isProduction,                        // true in production (HTTPS only)
		HTTPOnly:    true,                                // Prevents XSS attacks by making cookie inaccessible to JavaScript
		SameSite:    fiber.CookieSameSiteLaxMode,         // CSRF protection
		Partitioned: false,
		SessionOnly: false,
	}
	c.Cookie(&refreshCookies)
	return c.JSON(authResponse)
}

func (g *GoAuthFiber) Logout(c fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (g *GoAuthFiber) GoogleLogin(c fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (g *GoAuthFiber) GoogleCallback(c fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (g *GoAuthFiber) GithubLogin(c fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (g *GoAuthFiber) GithubCallback(c fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}
