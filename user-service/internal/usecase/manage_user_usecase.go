package usecase

import (
	"context"
	"gaman-microservice/user-service/internal/domain/entity"
	"gaman-microservice/user-service/internal/port/in"
	"gaman-microservice/user-service/internal/port/out"
)

type manageUserUseCase struct {
	userRepo out.UserRepository
}

func NewManageUserUseCase(userRepo out.UserRepository) in.ManageUserUseCase {
	return &manageUserUseCase{
		userRepo: userRepo,
	}
}

func (u *manageUserUseCase) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	err := u.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *manageUserUseCase) GetUser(ctx context.Context, id string) (*entity.User, error) {
	return u.userRepo.FindByID(ctx, id)
}

func (u *manageUserUseCase) UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	err := u.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *manageUserUseCase) DeleteUser(ctx context.Context, id string) error {
	return u.userRepo.Delete(ctx, id)
}

func (u *manageUserUseCase) ListUsers(ctx context.Context, limit uint64, cursor string) ([]*entity.User, error) {
	return u.userRepo.List(ctx, limit, cursor)
}
