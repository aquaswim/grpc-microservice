package in

import (
	"context"
	"gaman-microservice/user-service/internal/domain/model"
)

type UserUseCase interface {
	Login(ctx context.Context, username, password string) (*model.User, string, error)
	GetProfile(ctx context.Context, userID string) (*model.User, error)
}
