package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	TcpListenerUrl string `env:"TCP_LISTENER_URL" envDefault:":50051"`
	DatabaseUrl    string `env:"DATABASE_URL,required"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return cfg, nil
}
