package servers

import (
	"github.com/gofiber/fiber/v2"
	handlerMid "github.com/sing3demons/shop/modules/middlewares/handlers"
	repoMid "github.com/sing3demons/shop/modules/middlewares/repositories"
	usecaseMid "github.com/sing3demons/shop/modules/middlewares/usecases"
	"github.com/sing3demons/shop/modules/monitor/handlers"
	userHandlers "github.com/sing3demons/shop/modules/users/handlers"
	userRepository "github.com/sing3demons/shop/modules/users/repositories"
	userUsecases "github.com/sing3demons/shop/modules/users/usecases"
)

type IModuleFactory interface {
	MonitorModule()
	UserModule()
}

type moduleFactory struct {
	router     fiber.Router
	server     *server
	middleware handlerMid.IMiddlewareHandler
}

func InitMiddleware(server *server) handlerMid.IMiddlewareHandler {
	repo := repoMid.NewMiddlewareRepository(server.db)
	usecase := usecaseMid.NewMiddlewareUsecase(repo)
	return handlerMid.NewMiddlewareHandler(server.cfg, usecase)
}

func InitModule(r fiber.Router, server *server, m handlerMid.IMiddlewareHandler) IModuleFactory {
	return &moduleFactory{
		router:     r,
		server:     server,
		middleware: m,
	}
}
func (m *moduleFactory) MonitorModule() {
	handler := handlers.NewMonitorHandlers(m.server.cfg)
	m.router.Get("/health", handler.HealthCheck)
}

func (m *moduleFactory) UserModule() {
	repo := userRepository.NewUserRepository(m.server.db)
	usecase := userUsecases.NewUserUsecases(m.server.cfg, repo)
	handler := userHandlers.NewUserHandlers(m.server.cfg, usecase)

	r := m.router.Group("/users")
	r.Post("/signup", handler.SignUpCustomer)
	r.Post("/signin", handler.SignIn)
}
