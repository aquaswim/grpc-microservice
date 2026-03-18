package memory

import (
	"context"
	"errors"
	"gaman-microservice/user-service/internal/domain/entity"
	"gaman-microservice/user-service/internal/port/out"
)

type userMemoryRepository struct {
	users map[string]*entity.User
}

func NewUserMemoryRepository() out.UserRepository {
	return &userMemoryRepository{
		users: map[string]*entity.User{
			"1": {ID: "1", Username: "john_doe", Password: "password123", Email: "john@example.com"},
		},
	}
}

func (r *userMemoryRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	for _, user := range r.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *userMemoryRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	user, ok := r.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}
