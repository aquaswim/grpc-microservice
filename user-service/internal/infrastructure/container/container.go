package container

import (
	usergrpc "gaman-microservice/user-service/internal/adapter/handler/grpc"
	"gaman-microservice/user-service/internal/adapter/repository/memory"
	"gaman-microservice/user-service/internal/application/usecase"
	"gaman-microservice/user-service/internal/domain/repository"
	"gaman-microservice/user-service/internal/port/in"

	"github.com/golobby/container/v3"
)

func Init() container.Container {
	c := container.New()

	// Repository
	c.Singleton(func() repository.UserRepository {
		return memory.NewUserMemoryRepository()
	})

	// UseCase
	c.Singleton(func(repo repository.UserRepository) in.UserUseCase {
		return usecase.NewUserUseCase(repo)
	})

	// Handler
	c.Singleton(func(uc in.UserUseCase) *usergrpc.UserHandler {
		return usergrpc.NewUserHandler(uc)
	})

	return c
}
