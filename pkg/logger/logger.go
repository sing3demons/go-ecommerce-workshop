package logger

import (
	"encoding/json"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger() (*zap.Logger, func(), error) {
	cfg := zapcore.EncoderConfig{
		// TimeKey:        "ts",
		// LevelKey:       "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		FunctionKey:   zapcore.OmitKey,
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime:    zapcore.EpochTimeEncoder,
		// EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		// EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})

	core := zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(zapcore.NewCore(zapcore.NewConsoleEncoder(cfg), zapcore.AddSync(os.Stdout), highPriority))
	})

	log, err := zap.NewProduction(core, zap.AddStacktrace(zap.WarnLevel))
	if err != nil {
		return nil, nil, err
	}

	close := func() {
		log.Sync()
	}

	return log, close, nil
}

type Config struct {
	TimeFormat string
	UTC        bool
	SkipPaths  []string
	TraceID    bool // optionally log Open Telemetry TraceID
	Context    func(c *fiber.Ctx) []zapcore.Field
}

func ZapLogger(logger *zap.Logger) fiber.Handler {
	conf := &Config{
		UTC:        true,
		TimeFormat: time.RFC3339,
	}

	skipPaths := make(map[string]bool, len(conf.SkipPaths))
	for _, path := range conf.SkipPaths {
		skipPaths[path] = true
	}
	return func(c *fiber.Ctx) error {
		start := time.Now()
		path := c.Path()
		query := c.Request().URI().QueryArgs().String()
		c.Next()

		if _, ok := skipPaths[path]; !ok {
			end := time.Now()
			latency := end.Sub(start)
			if conf.UTC {
				end = end.UTC()
			}

			var result any
			json.Unmarshal(c.Response().Body(), &result)

			body := func(c *fiber.Ctx) map[string]any {
				var b map[string]any
				c.BodyParser(&b)

				if b["password"] != nil {
					delete(b, "password")
				}
				return b

			}

			fields := []zapcore.Field{
				zap.Int("status", c.Response().StatusCode()),
				zap.String("method", c.Method()),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.IP()),
				zap.String(fiber.HeaderUserAgent, string(c.Context().UserAgent())),
				zap.Duration("latency", latency),
				zap.Any("body", body(c)),
				zap.String(fiber.HeaderXRequestID, c.GetRespHeader(fiber.HeaderXRequestID)),
				zap.String(fiber.HeaderAccept, c.GetRespHeader(fiber.HeaderXRequestID)),
				zap.Any(fiber.HeaderContentType, c.GetRespHeader(fiber.HeaderContentType)),
				zap.Any("data", result),
			}

			if conf.TimeFormat != "" {
				fields = append(fields, zap.String("time", end.Format(conf.TimeFormat)))
			}

			if conf.Context != nil {
				fields = append(fields, conf.Context(c)...)
			}
			logger.Info(path, fields...)
		}
		return nil
	}
}
