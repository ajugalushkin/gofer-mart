package storage

import (
	"context"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/ajugalushkin/gofer-mart/config"
	"github.com/ajugalushkin/gofer-mart/internal/database"
	"github.com/ajugalushkin/gofer-mart/internal/logger"
)

func Init(ctx context.Context) (*sqlx.DB, error) {
	cfg := config.FlagsFromContext(ctx)
	db, err := database.NewConnection("pgx", cfg.DataBaseURI)
	if err != nil {
		logger.LogFromContext(ctx).Debug("storage.Init", zap.Error(err))
		return nil, err
	}
	return db, nil
}

type ctxConnect struct{}

func ContextWithConnect(ctx context.Context, connect *sqlx.DB) context.Context {
	return context.WithValue(ctx, ctxConnect{}, connect)
}

func ConnectFromContext(ctx context.Context) *sqlx.DB {
	if connect, ok := ctx.Value(ctxConnect{}).(*sqlx.DB); ok {
		return connect
	}
	return nil
}
