package framework

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v3"
	"github.com/labstack/echo/v4"
	"github.com/valyala/fasthttp"
)

type (
	Fiber interface {
		Login(c fiber.Ctx) error
		Register(c fiber.Ctx) error
		RefreshToken(c fiber.Ctx) error
		Logout(c fiber.Ctx) error
		GoogleLogin(c fiber.Ctx) error
		GoogleCallback(c fiber.Ctx) error
		GithubLogin(c fiber.Ctx) error
		GithubCallback(c fiber.Ctx) error
		Me(c fiber.Ctx) error
	}

	Gin interface {
		Login(ctx *gin.Context)
		Register(ctx *gin.Context)
		Logout(ctx *gin.Context)
		GoogleLogin(ctx *gin.Context)
		GoogleCallback(ctx *gin.Context)
		GithubLogin(ctx *gin.Context)
		GithubCallback(ctx *gin.Context)
	}

	Echo interface {
		Login(c echo.Context) error
		Register(c echo.Context) error
		Logout(c echo.Context) error
		GoogleLogin(c echo.Context) error
		GoogleCallback(c echo.Context) error
		GithubLogin(c echo.Context) error
		GithubCallback(c echo.Context) error
	}

	HTTP interface {
		Login(w http.ResponseWriter, r *http.Request)
		Register(w http.ResponseWriter, r *http.Request)
		Logout(w http.ResponseWriter, r *http.Request)
		GoogleLogin(w http.ResponseWriter, r *http.Request)
		GoogleCallback(w http.ResponseWriter, r *http.Request)
		GithubLogin(w http.ResponseWriter, r *http.Request)
		GithubCallback(w http.ResponseWriter, r *http.Request)
	}

	FastHTTP interface {
		Login(ctx *fasthttp.RequestCtx)
		Register(ctx *fasthttp.RequestCtx)
		Logout(ctx *fasthttp.RequestCtx)
		GoogleLogin(ctx *fasthttp.RequestCtx)
		GoogleCallback(ctx *fasthttp.RequestCtx)
		GithubLogin(ctx *fasthttp.RequestCtx)
		GithubCallback(ctx *fasthttp.RequestCtx)
	}
)
