package globalLogger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

type Config struct {
	LogPretty bool
	LogLevel  string
}

func Setup(cfg *Config) {
	// setup logging
	// src:  https://github.com/rs/zerolog/issues/174#issuecomment-516806803
	zerolog.LevelFieldName = "severity"
	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	if cfg.LogPretty {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		log.Warn().Msgf("Pretty logging is enabled, this must only be used in local!")
	}

	log.Info().Msgf("setting log level to %s", cfg.LogLevel)
	logLevel, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		logLevel = zerolog.InfoLevel
		log.Error().Err(err).Msgf("Invalid log level '%s', defaulting to INFO", cfg.LogLevel)
	}
	zerolog.SetGlobalLevel(logLevel)
}
