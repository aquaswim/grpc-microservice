package container

import (
	"database/sql"
	usergrpc "gaman-microservice/user-service/internal/adapter/handler/grpc"
	"gaman-microservice/user-service/internal/adapter/repository/postgres"
	"gaman-microservice/user-service/internal/infrastructure/config"
	"gaman-microservice/user-service/internal/infrastructure/pgsql"
	"gaman-microservice/user-service/internal/usecase"

	"github.com/golobby/container/v3"
)

func Init() container.Container {
	c := container.New()

	// Config
	container.MustSingleton(c, config.Load)

	// DB
	container.MustSingleton(c, func(cfg *config.Config) (*sql.DB, error) {
		return pgsql.Connect(cfg.DatabaseUrl)
	})

	// Repository
	container.MustSingleton(c, postgres.NewUserRepository)

	// UseCase
	container.MustSingleton(c, usecase.NewUserUseCase)

	// Handler
	container.MustSingleton(c, usergrpc.NewUserHandler)

	return c
}
