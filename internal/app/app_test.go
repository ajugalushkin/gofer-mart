package app

import (
	"context"
	"github.com/labstack/echo/v4"
	"testing"
)

func Test_app_register(t *testing.T) {
	type fields struct {
		ctx context.Context
	}
	type args struct {
		echoCtx echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := app{
				ctx: tt.fields.ctx,
			}
			if err := a.register(tt.args.echoCtx); (err != nil) != tt.wantErr {
				t.Errorf("register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_app_login(t *testing.T) {
	type fields struct {
		ctx context.Context
	}
	type args struct {
		echoCtx echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := app{
				ctx: tt.fields.ctx,
			}
			if err := a.login(tt.args.echoCtx); (err != nil) != tt.wantErr {
				t.Errorf("login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
