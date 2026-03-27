package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	PrettyLog               bool   `env:"PRETTY_LOG" envDefault:"false"`
	TcpListenerUrl          string `env:"TCP_LISTENER_URL" envDefault:":50051"`
	DatabaseUrl             string `env:"DATABASE_URL,required"`
	TokenPrivateKey         string `env:"TOKEN_PRIVATE_KEY,required"`
	TokenPublicKey          string `env:"TOKEN_PUBLIC_KEY"`
	TokenExpiryMinutes      int    `env:"TOKEN_EXPIRY_MINUTES" envDefault:"60"`
	ResetTokenExpiryMinutes int    `env:"RESET_TOKEN_EXPIRY_MINUTES" envDefault:"10"`
}

func (c Config) GetTokenExpiryDuration() time.Duration {
	return time.Duration(c.TokenExpiryMinutes) * time.Minute
}

func (c Config) GetResetTokenExpiryDuration() time.Duration {
	return time.Duration(c.ResetTokenExpiryMinutes) * time.Minute
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return cfg, nil
}
