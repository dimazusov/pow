package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type ServerConfig struct {
	ServerAddress string           `env:"SERVER_ADDRESS"       envDefault:"0.0.0.0:8080"`
	PowDifficulty uint8            `envPrefix:"POW_DIFFICULTY" envDefault:"3"`
	LoggerLvl     string           `envPrefix:"LOGGER_LVL"     envDefault:"info"`
	DictPath      string           `env:"DICT_PATH"`
	Connection    ConnectionConfig `envPrefix:"CONN_"`
}

type ConnectionConfig struct {
	KeepAlive    time.Duration `env:"KEEP_ALIVE"    envDefault:"60s"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" envDefault:"10s"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT"  envDefault:"10s"`
	IdleTimeout  time.Duration `env:"IDLE_TIMEOUT"  envDefault:"30s"`
}

func NewServerCfg() (*ServerConfig, error) {
	cfg := &ServerConfig{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse configuration: %w", err)
	}

	return cfg, nil
}
