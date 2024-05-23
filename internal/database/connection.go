package database

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
)

func NewConnection(driver, dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		slog.Error("failed to create a database connection", err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		slog.Error("failed to ping the database", err)
		return nil, err
	}

	Migrate(dsn)

	return db, nil
}
