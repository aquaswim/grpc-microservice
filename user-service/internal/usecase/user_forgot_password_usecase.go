package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	appError "gaman-microservice/user-service/internal/domain/app_error"
	"gaman-microservice/user-service/internal/domain/entity"
	"gaman-microservice/user-service/internal/infrastructure/config"
	"gaman-microservice/user-service/internal/port/in"
	"gaman-microservice/user-service/internal/port/out"
	"time"

	"github.com/joomcode/errorx"
	"github.com/rs/zerolog/log"
)

type userForgotPasswordUseCase struct {
	userRepo       out.UserRepository
	resetTokenRepo out.PasswordResetTokenRepository
	cfg            *config.Config
	eventProducer  out.EventProducer
}

func NewUserForgotPasswordUseCase(
	userRepo out.UserRepository,
	resetTokenRepo out.PasswordResetTokenRepository,
	cfg *config.Config,
	eventProducer out.EventProducer,
) in.UserForgotPasswordUseCase {
	return &userForgotPasswordUseCase{
		userRepo:       userRepo,
		resetTokenRepo: resetTokenRepo,
		cfg:            cfg,
		eventProducer:  eventProducer,
	}
}

func (u *userForgotPasswordUseCase) ForgotPassword(ctx context.Context, email string) error {
	l := log.Ctx(ctx)

	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errorx.IsOfType(err, appError.ErrNotFound) {
			// Don't leak user existence
			l.Warn().Str("email", email).Msg("Email not found (this still return ok to avoid security issue)")
			return nil
		}
		return err
	}

	// Generate random token
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return appError.ErrInternal.Wrap(err, "failed to generate reset token")
	}
	tokenString := hex.EncodeToString(b)

	// Clean up old tokens for this user
	err = u.resetTokenRepo.DeleteByUserID(ctx, user.ID)
	if err != nil {
		l.Warn().
			Err(err).
			Msg("Failed to delete old reset tokens")
	}

	resetToken := &entity.PasswordResetToken{
		UserID:    user.ID,
		Token:     tokenString,
		ExpiresAt: time.Now().Add(u.cfg.GetResetTokenExpiryDuration()),
		CreatedAt: time.Now(),
	}

	if err := u.resetTokenRepo.Create(ctx, resetToken); err != nil {
		return err
	}

	// forgot password event
	err = u.eventProducer.ForgotPassword(ctx, &entity.UserForgotPasswordData{
		User:      user,
		Token:     resetToken.Token,
		ExpiredAt: resetToken.ExpiresAt,
	})
	if err != nil {
		return appError.ErrInternal.Wrap(err, "failed to process forgot password event")
	}

	return nil
}

func (u *userForgotPasswordUseCase) ValidateResetToken(ctx context.Context, token string) (bool, error) {
	resetToken, err := u.resetTokenRepo.FindByToken(ctx, token)
	if err != nil {
		if errorx.IsOfType(err, appError.ErrNotFound) {
			return false, nil
		}
		return false, err
	}

	if resetToken.IsExpired() {
		return false, nil
	}

	return true, nil
}

func (u *userForgotPasswordUseCase) ResetPassword(ctx context.Context, token, newPassword string) error {
	l := log.Ctx(ctx)

	resetToken, err := u.resetTokenRepo.FindByToken(ctx, token)
	if err != nil {
		if errorx.IsOfType(err, appError.ErrNotFound) {
			return appError.ErrValidation.New("invalid or expired token")
		}
		return err
	}

	if resetToken.IsExpired() {
		return appError.ErrValidation.New("invalid or expired token")
	}

	user, err := u.userRepo.FindByID(ctx, resetToken.UserID)
	if err != nil {
		return err
	}

	if err := user.SetPassword(newPassword); err != nil {
		return appError.ErrInternal.Wrap(err, "failed to set new password")
	}

	if err := u.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// Invalidate token after use
	err2 := u.resetTokenRepo.Delete(ctx, token)
	if err2 != nil {
		l.Warn().Err(err2).Msg("failed to delete reset token")
	}

	err2 = u.eventProducer.UserResetPasswordDone(ctx, &entity.UserResetPasswordDoneData{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
	})
	if err2 != nil {
		l.Warn().Err(err2).Msg("failed to process user reset password done event")
	}

	return nil
}
