package servers

import (
	"github.com/gofiber/fiber/v2"
	handlerMid "github.com/sing3demons/shop/modules/middlewares/handlers"
	repoMid "github.com/sing3demons/shop/modules/middlewares/repositories"
	usecaseMid "github.com/sing3demons/shop/modules/middlewares/usecases"
	"github.com/sing3demons/shop/modules/monitor/handlers"
)

type IModuleFactory interface {
	MonitorModule()
}

type moduleFactory struct {
	router     fiber.Router
	server     *server
	middleware handlerMid.IMiddlewareHandler
}

func InitModule(r fiber.Router, server *server, m handlerMid.IMiddlewareHandler) IModuleFactory {
	return &moduleFactory{
		router:     r,
		server:     server,
		middleware: m,
	}
}

func InitMiddleware(server *server) handlerMid.IMiddlewareHandler {
	repo := repoMid.NewMiddlewareRepository(server.db)
	usecase := usecaseMid.NewMiddlewareUsecase(repo)
	return handlerMid.NewMiddlewareHandler(server.cfg, usecase)
}

func (m *moduleFactory) MonitorModule() {
	handler := handlers.NewMonitorHandlers(m.server.cfg)
	m.router.Get("/health", handler.HealthCheck)
}
