package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sing3demons/shop/config"
	"github.com/sing3demons/shop/modules/entities"
	"github.com/sing3demons/shop/modules/monitor"
)

type IMonitorHandlers interface {
	HealthCheck(c *fiber.Ctx) error
}

type monitorHandlers struct {
	cfg config.IConfig
}

func NewMonitorHandlers(cfg config.IConfig) IMonitorHandlers {
	return &monitorHandlers{
		cfg: cfg,
	}
}

func (h monitorHandlers) HealthCheck(c *fiber.Ctx) error {
	resp := monitor.Monitor{
		Name:    h.cfg.App().Name(),
		Version: h.cfg.App().Version(),
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, resp).Response()
	// return c.Status(fiber.StatusOK).JSON(resp)
}
