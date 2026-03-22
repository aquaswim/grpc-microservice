package config

import "github.com/caarlos0/env/v11"

type Config struct {
	PrettyLog   bool   `env:"PRETTY_LOG" envDefault:"false"`
	ListenAddr  string `env:"LISTEN_ADDR" envDefault:":8080"`
	UserSvcAddr string `env:"USER_SVC_ADDR,required"`
}

func Load() (*Config, error) {
	cfg, err := env.ParseAs[Config]()
	return &cfg, err
}
