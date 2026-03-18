package usecase

import (
	"context"
	"errors"
	"gaman-microservice/user-service/internal/domain/model"
	"gaman-microservice/user-service/internal/domain/repository"
	"gaman-microservice/user-service/internal/port/in"
)

type userUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) in.UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
	}
}

func (u *userUseCase) Login(ctx context.Context, username, password string) (*model.User, string, error) {
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

func (u *userUseCase) GetProfile(ctx context.Context, userID string) (*model.User, error) {
	return u.userRepo.FindByID(ctx, userID)
}
