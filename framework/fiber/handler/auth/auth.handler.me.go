package auth

import (
	"errors"

	"github.com/SwanHtetAungPhyo/go-auth/framework/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (g *GoAuthFiber) Me(c fiber.Ctx) error {
	userId := c.Locals("userId")
	if userId == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.GeneralResponse{
			Message: "UserId was not found",
			Data:    nil,
			Error:   errors.New("UserId was not found in the token, malformed token "),
		})
	}
	userIdUUId := uuid.MustParse(userId.(string))
	if userIdUUId == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.GeneralResponse{
			Message: "UserId was not found",
			Data:    nil,
			Error:   errors.New("UserId was not found in the token, malformed token "),
		})
	}

	me, err := g.srv.Me(userIdUUId)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.GeneralResponse{
			Message: "UserId was not found",
			Data:    nil,
			Error:   errors.New("UserId was not found in the database "),
		})
	}
	return c.Status(fiber.StatusOK).JSON(utils.GeneralResponse{
		Data: me,
	})
}
