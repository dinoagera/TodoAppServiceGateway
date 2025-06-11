package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Address     string        `env:"SERVER_ADDRESS" env-default:"localhost:8080"`
	TimeOut     time.Duration `env:"SERVER_TIMEOUT"`
	IdleTimeout time.Duration `env:"SERVER_IDLE_TIMEOUT"`
	GRPCPorts   GRPCPorts
	SecretKey   string `env:"SECRET_KEY" env-required:"true"`
}
type GRPCPorts struct {
	GRPCApiDb   string `env:"GRPC_API_DB_PORT" env-required:"true"`
	GRPCApiAuth string `env:"GRPC_API_AUTH_PORT" env-required:"true"`
}

func InitConfig(log *slog.Logger) *Config {
	cfgPath := ".env"
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		log.Debug("config not found from", "Error:", err.Error(), "CfgPath", cfgPath)
		log.Info("not found config")
		os.Exit(1)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		log.Debug("cant read configfile", "err:", err.Error())
		log.Info("read config file to failed")
		os.Exit(1)
	}
	return &cfg
}
