package servers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sing3demons/shop/modules/monitor/handlers"
)

type IModuleFactory interface {
	MonitorModule()
}

type moduleFactory struct {
	router fiber.Router
	server *server
}

func InitModule(r fiber.Router, server *server) IModuleFactory {
	return &moduleFactory{
		router: r,
		server: server,
	}
}

func (m *moduleFactory) MonitorModule() {
	handler := handlers.NewMonitorHandlers(m.server.cfg)
	m.router.Get("/health", handler.HealthCheck)
}
