package config

import "github.com/caarlos0/env/v11"

type Config struct {
	PrettyLog   bool   `env:"PRETTY_LOG" envDefault:"false"`
	ListenAddr  string `env:"LISTEN_ADDR" envDefault:":8080"`
	UserSvcAddr string `env:"USER_SVC_ADDR,required"`

	RedisAddr string `env:"REDIS_ADDR,required"`
	RedisUser string `env:"REDIS_USER"`
	RedisPass string `env:"REDIS_PASS"`
	RedisDB   int    `env:"REDIS_DB"`
}

func Load() (*Config, error) {
	cfg, err := env.ParseAs[Config]()
	return &cfg, err
}
