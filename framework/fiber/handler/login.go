package handler

import (
	"github.com/SwanHtetAungPhyo/go-auth/framework"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type GoAuthFiber struct {
	connPool *pgxpool.Pool
}

func NewGoAuthFiber(connPool *pgxpool.Pool) *GoAuthFiber {
	if connPool == nil {
		log.Fatal().Msg("connPool is nil in NewGoAuthFiber handler")
		return nil
	}
	return &GoAuthFiber{
		connPool: connPool,
	}
}

var _ framework.Fiber = (*GoAuthFiber)(nil)

func (g *GoAuthFiber) Login(c fiber.Ctx) error {
	panic("implement me")
}

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
