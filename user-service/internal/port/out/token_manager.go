package out

import (
	"context"
	"gaman-microservice/user-service/internal/domain/entity"
	"time"
)

type TokenManager interface {
	Generate(ctx context.Context, tokenData *entity.TokenData) (string, time.Time, error)
	Validate(ctx context.Context, token string) (*entity.TokenData, error)
}
