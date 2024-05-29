package main

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"

	"github.com/ajugalushkin/gofer-mart/config"
	"github.com/ajugalushkin/gofer-mart/internal/app"
	"github.com/ajugalushkin/gofer-mart/internal/storage"
)

func main() {
	cfg := config.ReadConfig()
	ctx := config.ContextWithFlags(context.Background(), cfg)

	db, err := storage.Init(ctx)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	a := app.NewApp(ctx, db)
	e := echo.New()
	a.Routes(e)

	err = e.Start(cfg.RunAddr)
	if err != nil {
		fmt.Println(err)
	}
}
