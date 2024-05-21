package logger

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type ctxLogger struct{}

func ContextWithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxLogger{}, logger)
}

func LogFromContext(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(ctxLogger{}).(*zap.Logger); ok {
		return logger
	}
	return zap.L()
}

func Initialize(level string) (*zap.Logger, error) {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}
	cfg := zap.NewProductionConfig()

	cfg.Level = lvl

	zl, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return zl, nil
}

func MiddlewareLogger(ctx context.Context) func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		fn := func(context echo.Context) error {
			start := time.Now()

			if err := next(context); err != nil {
				context.Error(err)
			}

			duration := time.Since(start)

			log := LogFromContext(ctx)
			log.Debug("got incoming HTTP request",
				zap.String("method", context.Request().Method),
				zap.String("path", context.Request().URL.Path),
				zap.String("time", duration.String()),
			)
			return nil
		}
		return fn
	}
}
