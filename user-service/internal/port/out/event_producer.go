package out

import (
	"context"
	"gaman-microservice/user-service/internal/domain/entity"
)

type EventProducer interface {
	ForgotPassword(ctx context.Context, data *entity.UserForgotPasswordData) error
	UserResetPasswordDone(ctx context.Context, data *entity.UserResetPasswordDoneData) error
}
