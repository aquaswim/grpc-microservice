package in

import (
	"context"
)

type UserForgotPasswordUseCase interface {
	ForgotPassword(ctx context.Context, email string) error
	ValidateResetToken(ctx context.Context, token string) (bool, error)
	ResetPassword(ctx context.Context, token, newPassword string) error
}
