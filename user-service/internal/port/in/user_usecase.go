package in

import (
	"context"
	"gaman-microservice/user-service/internal/domain/entity"
)

type UserUseCase interface {
	Login(ctx context.Context, username, password string) (*entity.User, string, error)
	GetProfile(ctx context.Context, userID string) (*entity.User, error)
}
