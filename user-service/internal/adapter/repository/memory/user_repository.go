package memory

import (
	"context"
	"errors"
	"gaman-microservice/user-service/internal/domain/model"
	"gaman-microservice/user-service/internal/domain/repository"
)

type userMemoryRepository struct {
	users map[string]*model.User
}

func NewUserMemoryRepository() repository.UserRepository {
	return &userMemoryRepository{
		users: map[string]*model.User{
			"1": {ID: "1", Username: "john_doe", Password: "password123", Email: "john@example.com"},
		},
	}
}

func (r *userMemoryRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	for _, user := range r.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *userMemoryRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	user, ok := r.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}
