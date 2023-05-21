package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/sing3demons/shop/config"
	"github.com/sing3demons/shop/modules/entities"
	"github.com/sing3demons/shop/modules/middlewares/usecases"
)

type IMiddlewareHandler interface {
	Cors() fiber.Handler
	RouterCheck() fiber.Handler
	Logger() fiber.Handler
}

type middlewareHandler struct {
	cfg     config.IConfig
	usecase usecases.IMiddlewareUsecase
}

func NewMiddlewareHandler(cfg config.IConfig, usecase usecases.IMiddlewareUsecase) IMiddlewareHandler {
	return &middlewareHandler{
		cfg:     cfg,
		usecase: usecase,
	}
}

func (h *middlewareHandler) Cors() fiber.Handler {
	return cors.New(cors.Config{
		Next:             cors.ConfigDefault.Next,
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
	})
}

func (h *middlewareHandler) RouterCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return entities.NewResponse(c).Error(fiber.ErrNotFound.Code, "router-001", fiber.ErrNotFound.Message).Response()
	}
}

func (h *middlewareHandler) Logger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "[${time}] ${ip} ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: "02/01/2006",
		TimeZone:   "Asia/Bangkok",
	})
}
