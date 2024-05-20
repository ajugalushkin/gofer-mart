package config

import (
	"context"
	"flag"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type AppConfig struct {
	RunAddr              string `env:"RUN_ADDRESS"`
	DataBaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	TokenKey             string `env:"TOKEN_KEY"`
}

func init() {
	viper.SetDefault("RunAddr", ":8080")
	viper.SetDefault("DataBaseURI",
		"postgres://postgres:postgres@localhost/postgres?sslmode=disable")
	viper.SetDefault("AccrualSystemAddress", "")
	viper.SetDefault("TokenKey", "")
}

func bindToEnv() {
	viper.SetEnvPrefix("test")
	_ = viper.BindEnv("RunAddr")
	_ = viper.BindEnv("DataBaseURI")
	_ = viper.BindEnv("AccrualSystemAddress")
	_ = viper.BindEnv("TokenKey")
}

func bindToFlag() {
	flag.String("a", "localhost:8080", "address and port to run server")
	flag.String("d", "postgres://postgres:postgres@localhost/postgres?sslmode=disable",
		"database connection address")
	flag.String("r", "", "address of the accrual calculation system")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)
}

func ReadConfig() *AppConfig {
	bindToFlag()
	bindToEnv()

	result := &AppConfig{
		RunAddr: viper.GetString("RunAddr"),
		//FlagLogLevel: viper.GetString("FlagLogLevel"),
		DataBaseURI:          viper.GetString("DataBaseURI"),
		AccrualSystemAddress: viper.GetString("AccrualSystemAddress"),
	}
	return result
}

type ctxConfig struct{}

func ContextWithFlags(ctx context.Context, config *AppConfig) context.Context {
	return context.WithValue(ctx, ctxConfig{}, config)
}

func FlagsFromContext(ctx context.Context) *AppConfig {
	if config, ok := ctx.Value(ctxConfig{}).(*AppConfig); ok {
		return config
	}
	return &AppConfig{}
}
