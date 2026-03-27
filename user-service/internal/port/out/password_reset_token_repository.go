package out

import (
	"context"
	"gaman-microservice/user-service/internal/domain/entity"
)

type PasswordResetTokenRepository interface {
	Create(ctx context.Context, token *entity.PasswordResetToken) error
	FindByToken(ctx context.Context, token string) (*entity.PasswordResetToken, error)
	DeleteByUserID(ctx context.Context, userID string) error
	Delete(ctx context.Context, token string) error
}
