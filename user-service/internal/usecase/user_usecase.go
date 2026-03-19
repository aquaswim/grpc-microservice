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
