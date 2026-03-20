package in

import (
	"context"
	"gaman-microservice/user-service/internal/domain/entity"
)

type UserUseCase interface {
	Login(ctx context.Context, username, password string) (*entity.User, string, error)
	ValidateToken(ctx context.Context, token string) (*entity.TokenData, error)
	CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	GetUser(ctx context.Context, id string) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context, limit uint64, cursor string) ([]*entity.User, error)
}
