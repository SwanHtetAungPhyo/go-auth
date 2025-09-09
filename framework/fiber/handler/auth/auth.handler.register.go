package auth

import (
	"github.com/SwanHtetAungPhyo/go-auth/framework"
	"github.com/gofiber/fiber/v3"
)

func (g *GoAuthFiber) Register(c fiber.Ctx) error {
	var req framework.RegisterRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.JSON(fiber.Map{
			"error": err,
		})
	}
	if err := framework.ValidateStruct(req); err != nil {
		return c.JSON(fiber.Map{
			"error": err,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}
