package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	LogPretty        bool   `env:"LOG_PRETTY" envDefault:"false"`
	LogLevel         string `env:"LOG_LEVEL" envDefault:"info"`
	RabbitMQUrl      string `env:"RABBITMQ_URL,required"`
	RabbitMqExchange string `env:"RABBITMQ_EXCHANGE,required"`

	UserForgotPasswordTopic    string `env:"TOPIC_USER_FORGOT_PASSWORD" envDefault:"user-forgot-password"`
	UserResetPasswordDoneTopic string `env:"TOPIC_USER_RESET_PASSWORD_DONE" envDefault:"user-reset-password-done"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return cfg, nil
}
