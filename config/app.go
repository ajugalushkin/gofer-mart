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
	NumOfWorkers         int    `env:"NUM_WORKERS"`
	LogLevel             string `env:"LOG_LEVEL"`
}

func init() {
	viper.SetDefault("RunAddr", ":8080")
	viper.SetDefault("DataBase_URI", "")
	viper.SetDefault("AccrualSystemAddress", "")
	viper.SetDefault("TokenKey", "")
	viper.SetDefault("NumOfWorkers", 10)
	viper.SetDefault("LogLevel", "info")
}

func bindToEnv() {
	//viper.SetEnvPrefix("test")
	_ = viper.BindEnv("RunAddr")
	_ = viper.BindEnv("DataBase_URI")
	_ = viper.BindEnv("AccrualSystemAddress")
	_ = viper.BindEnv("TokenKey")
	_ = viper.BindEnv("NumOfWorkers")
	_ = viper.BindEnv("LogLevel")
}

func bindToFlag() {
	flag.String("a", "localhost:8080", "address and port to run server")
	flag.String("d", "postgres://postgres:postgres@localhost/postgres?sslmode=disable",
		"database connection address")
	flag.String("r", "", "address of the accrualclient calculation system")
	flag.String("l", "info", "log level")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)
}

func ReadConfig() *AppConfig {
	bindToFlag()
	bindToEnv()

	result := &AppConfig{
		RunAddr:              viper.GetString("RunAddr"),
		DataBaseURI:          viper.GetString("DataBase_URI"),
		AccrualSystemAddress: viper.GetString("AccrualSystemAddress"),
		TokenKey:             viper.GetString("TokenKey"),
		NumOfWorkers:         viper.GetInt("NumOfWorkers"),
		LogLevel:             viper.GetString("LogLevel"),
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
