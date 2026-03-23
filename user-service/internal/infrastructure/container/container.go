package container

import (
	"context"
	"gaman-microservice/user-service/internal/adapter/auth"
	usergrpc "gaman-microservice/user-service/internal/adapter/handler/grpc"
	"gaman-microservice/user-service/internal/adapter/repository/postgres"
	"gaman-microservice/user-service/internal/infrastructure/config"
	"gaman-microservice/user-service/internal/infrastructure/pgsql"
	"gaman-microservice/user-service/internal/port/out"
	"gaman-microservice/user-service/internal/usecase"
	"os"

	"github.com/golobby/container/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Init() container.Container {
	c := container.New()

	// Config
	container.MustSingleton(c, config.Load)

	// setup logging
	container.MustCall(c, func(cfg *config.Config) {
		if cfg.PrettyLog {
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
			log.Warn().Msgf("Pretty logging is enabled, this must only be used in local!")
		}
	})

	// DB
	container.MustSingleton(c, func(cfg *config.Config) (*pgxpool.Pool, error) {
		return pgsql.Connect(context.Background(), cfg.DatabaseUrl)
	})

	// Token Manager
	container.MustSingleton(c, func(cfg *config.Config) (out.TokenManager, error) {
		return auth.NewPasetoManager(cfg.TokenSecret, cfg.GetTokenExpiryDuration())
	})

	// Repository
	container.MustSingleton(c, postgres.NewUserRepository)

	// UseCase
	container.MustSingleton(c, usecase.NewUserUseCase)

	// Handler
	container.MustSingleton(c, usergrpc.NewUserHandler)

	return c
}
