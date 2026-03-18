package usecase

import (
	"context"
	"errors"
	"gaman-microservice/user-service/internal/domain/entity"
	"gaman-microservice/user-service/internal/port/in"
	"gaman-microservice/user-service/internal/port/out"
)

type userUseCase struct {
	userRepo out.UserRepository
}

func NewUserUseCase(userRepo out.UserRepository) in.UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
	}
}

func (u *userUseCase) Login(ctx context.Context, username, password string) (*entity.User, string, error) {
	user, err := u.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, "", err
	}

	if user.Password != password {
		return nil, "", errors.New("invalid credentials")
	}

	// Mock token
	token := "mock-token-for-" + user.ID
	return user, token, nil
}

func (u *userUseCase) GetProfile(ctx context.Context, userID string) (*entity.User, error) {
	return u.userRepo.FindByID(ctx, userID)
}
