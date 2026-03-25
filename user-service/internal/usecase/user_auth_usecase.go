package usecase

import (
	"context"
	"fmt"
	appError "gaman-microservice/user-service/internal/domain/app_error"
	"gaman-microservice/user-service/internal/domain/entity"
	"gaman-microservice/user-service/internal/port/in"
	"gaman-microservice/user-service/internal/port/out"
)

type userAuthUseCase struct {
	userRepo     out.UserRepository
	tokenManager out.TokenManager
}

func NewUserAuthUseCase(userRepo out.UserRepository, tokenManager out.TokenManager) in.UserAuthUseCase {
	return &userAuthUseCase{
		userRepo:     userRepo,
		tokenManager: tokenManager,
	}
}

func (u *userAuthUseCase) Login(ctx context.Context, username, password string) (*entity.User, *entity.TokenWithExpiry, error) {
	user, err := u.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, nil, appError.ErrUnauthorized.Wrap(err, "invalid credential")
	}

	if err := user.ValidatePassword(password); err != nil {
		return nil, nil, appError.ErrUnauthorized.New("invalid credential")
	}

	token, expTime, err := u.tokenManager.Generate(ctx, &entity.TokenData{
		Id:       user.ID,
		Username: user.Username,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate token: %w", err)
	}
	return user, &entity.TokenWithExpiry{Token: token, Expiry: expTime}, nil
}

func (u *userAuthUseCase) ValidateToken(ctx context.Context, token string) (*entity.TokenData, error) {
	return u.tokenManager.Validate(ctx, token)
}
