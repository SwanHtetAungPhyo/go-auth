package auth

import "github.com/gofiber/fiber/v3"

func (g *GoAuthFiber) RefreshToken(c fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")

	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "refresh token not found",
		})
	}

	// TODO: validate refreshToken, generate new access token, etc.
	whoAreU := c.Cookies("session")
	if whoAreU == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "session not found",
		})
	}
	return c.JSON(fiber.Map{
		"message":       "refresh successful",
		"refresh_token": refreshToken,
	})
}
