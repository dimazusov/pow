package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type ClientConfig struct {
	ClientAddress string           `env:"CLIENT_ADDRESS" envDefault:"0.0.0.0:8080"`
	RequestCount  int              `env:"REQUEST_COUNT" envDefault:"5"`
	LoggerLvl     string           `envPrefix:"LOGGER_LVL" envDefault:"info"`
	Connection    ConnectionConfig `envPrefix:"CONN_"`
}

func NewClientCfg() (*ClientConfig, error) {
	cfg := &ClientConfig{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse configuration: %w", err)
	}

	return cfg, nil
}
