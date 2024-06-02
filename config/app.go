package config

import (
	"context"
	"flag"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type AppConfig struct {
	RunAddr              string `env:"RUN_ADDRESS"`
	DataBaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	TokenKey             string `env:"TOKEN_KEY"`
	NumWorkers           int    `env:"NUM_WORKERS"`
	LogLevel             string `env:"LOG_LEVEL"`
}

func init() {
	err := godotenv.Load("/docker/.env")
	if err != nil {
		log.Debug("Error loading .env file", "error", err)
	}

	viper.SetDefault("Run_Address", "")
	viper.SetDefault("DataBase_URI", "")
	viper.SetDefault("Accrual_System_Address", "")
	viper.SetDefault("Token_Key", "")
	viper.SetDefault("Num_Workers", 0)
	viper.SetDefault("Log_Level", "")
}

func bindToEnv() {
	_ = viper.BindEnv("Run_Address")
	_ = viper.BindEnv("DataBase_URI")
	_ = viper.BindEnv("Accrual_System_Address")
	_ = viper.BindEnv("Token_Key")
	_ = viper.BindEnv("Num_Workers")
	_ = viper.BindEnv("Log_Level")
}

func bindToFlag() {
	flag.String("a", "localhost:8080", "address and port to run server")
	flag.String("d", "postgres://postgres:postgres@localhost/postgres?sslmode=disable",
		"database connection address")
	flag.String("r", "", "address of the accrualclient calculation system")
	flag.String("l", "debug", "log level")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)
}

func ReadConfig() *AppConfig {
	bindToFlag()
	bindToEnv()

	result := &AppConfig{
		RunAddr:              viper.GetString("Run_Address"),
		DataBaseURI:          viper.GetString("DataBase_URI"),
		AccrualSystemAddress: viper.GetString("Accrual_System_Address"),
		TokenKey:             viper.GetString("Token_Key"),
		NumWorkers:           viper.GetInt("Num_Workers"),
		LogLevel:             viper.GetString("Log_Level"),
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
