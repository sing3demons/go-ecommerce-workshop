package servers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/jmoiron/sqlx"
	"github.com/sing3demons/shop/config"
	logMiddleware "github.com/sing3demons/shop/pkg/logger"
	"go.uber.org/zap"
)

type server struct {
	app *fiber.App
	cfg config.IConfig
	db  *sqlx.DB // *sql.DB
}

type IServer interface {
	Start(logger *zap.Logger)
}

func NewServer(cfg config.IConfig, db *sqlx.DB) IServer {
	return &server{
		app: fiber.New(fiber.Config{
			AppName:      cfg.App().Name(),
			BodyLimit:    cfg.App().BodyLimit(),
			ReadTimeout:  cfg.App().ReadTimeout(),
			WriteTimeout: cfg.App().WriteTimeout(),
			JSONEncoder:  json.Marshal,
			JSONDecoder:  json.Unmarshal,
		}),
		cfg: cfg,
		db:  db,
	}
}

func (s *server) Start(logger *zap.Logger) {
	s.app.Use(logMiddleware.ZapLogger(logger))
	s.app.Use(requestid.New(requestid.ConfigDefault))

	v1 := s.app.Group("/api/v1")
	module := InitModule(v1, s)
	module.MonitorModule()
	//Graceful Shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := s.app.Listen(s.cfg.App().Url()); err != nil {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	stop()

	fmt.Println("shutting down gracefully, press Ctrl+C again to force")
	if err := s.app.Shutdown(); err != nil {
		log.Fatalf("shutdown: %s\n", err)
	}

	fmt.Println("Running cleanup tasks...")
	// Your cleanup tasks go here
	s.db.Close()
	fmt.Println("Fiber was successful shutdown.")
}
