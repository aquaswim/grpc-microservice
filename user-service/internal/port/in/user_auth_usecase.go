package in

import (
	"context"
	"gaman-microservice/user-service/internal/domain/entity"
)

type UserAuthUseCase interface {
	Login(ctx context.Context, username, password string) (*entity.User, *entity.TokenWithExpiry, error)
	ValidateToken(ctx context.Context, token string) (*entity.TokenData, error)
}
