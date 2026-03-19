package container

import (
	usergrpc "gaman-microservice/user-service/internal/adapter/handler/grpc"
	"gaman-microservice/user-service/internal/adapter/repository/memory"
	"gaman-microservice/user-service/internal/infrastructure/config"
	"gaman-microservice/user-service/internal/usecase"

	"github.com/golobby/container/v3"
)

func Init() container.Container {
	c := container.New()

	// Config
	container.MustSingleton(c, func() (*config.Config, error) {
		return config.Load()
	})

	// Repository
	container.MustSingleton(c, memory.NewUserMemoryRepository)

	// UseCase
	container.MustSingleton(c, usecase.NewUserUseCase)

	// Handler
	container.MustSingleton(c, usergrpc.NewUserHandler)

	return c
}
