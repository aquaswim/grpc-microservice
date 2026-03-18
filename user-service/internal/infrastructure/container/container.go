package container

import (
	usergrpc "gaman-microservice/user-service/internal/adapter/handler/grpc"
	"gaman-microservice/user-service/internal/adapter/repository/memory"
	"gaman-microservice/user-service/internal/usecase"

	"github.com/golobby/container/v3"
)

func Init() container.Container {
	c := container.New()

	// Repository
	container.MustSingleton(c, memory.NewUserMemoryRepository)

	// UseCase
	container.MustSingleton(c, usecase.NewUserUseCase)

	// Handler
	container.MustSingleton(c, usergrpc.NewUserHandler)

	return c
}
