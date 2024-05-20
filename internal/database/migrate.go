package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/pressly/goose/v3"

	"github.com/ajugalushkin/gofer-mart/migrations"
)

func Migrate(dataSourceName string) {
	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		log.Fatalf("sql.Open(): %v", err)
	}
	defer db.Close()

	goose.SetBaseFS(migrations.Migrations)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const cmd = "up"

	err = goose.RunContext(ctx, cmd, db, ".")
	if err != nil {
		log.Fatalf("goose.Status(): %v", err)
	}
}
