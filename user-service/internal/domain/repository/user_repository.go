package repository

import (
	"context"
	"gaman-microservice/user-service/internal/domain/model"
)

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	FindByID(ctx context.Context, id string) (*model.User, error)
}
