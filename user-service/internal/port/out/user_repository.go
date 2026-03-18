package out

import (
	"context"
	"gaman-microservice/user-service/internal/domain/entity"
)

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	FindByID(ctx context.Context, id string) (*entity.User, error)
}
