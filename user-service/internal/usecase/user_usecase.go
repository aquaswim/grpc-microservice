package usecase

import (
	"context"
	"fmt"
	appError "gaman-microservice/user-service/internal/domain/app_error"
	"gaman-microservice/user-service/internal/domain/entity"
	"gaman-microservice/user-service/internal/port/in"
	"gaman-microservice/user-service/internal/port/out"
)

type userUseCase struct {
	userRepo     out.UserRepository
	tokenManager out.TokenManager
}

func NewUserUseCase(userRepo out.UserRepository, tokenManager out.TokenManager) in.UserUseCase {
	return &userUseCase{
		userRepo:     userRepo,
		tokenManager: tokenManager,
	}
}

func (u *userUseCase) Login(ctx context.Context, username, password string) (*entity.User, *entity.TokenWithExpiry, error) {
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

func (u *userUseCase) ValidateToken(ctx context.Context, token string) (*entity.TokenData, error) {
	return u.tokenManager.Validate(ctx, token)
}

func (u *userUseCase) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	err := u.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUseCase) GetUser(ctx context.Context, id string) (*entity.User, error) {
	return u.userRepo.FindByID(ctx, id)
}

func (u *userUseCase) UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	err := u.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUseCase) DeleteUser(ctx context.Context, id string) error {
	return u.userRepo.Delete(ctx, id)
}

func (u *userUseCase) ListUsers(ctx context.Context, limit uint64, cursor string) ([]*entity.User, error) {
	return u.userRepo.List(ctx, limit, cursor)
}
