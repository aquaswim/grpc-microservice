package memory

import (
	"context"
	appError "gaman-microservice/user-service/internal/domain/app_error"
	"gaman-microservice/user-service/internal/domain/entity"
	"gaman-microservice/user-service/internal/port/out"
	"slices"
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

func (r *userMemoryRepository) FindByUsername(_ context.Context, username string) (*entity.User, error) {
	for _, user := range r.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, appError.ErrNotFound.New("user not found")
}

func (r *userMemoryRepository) FindByID(_ context.Context, id string) (*entity.User, error) {
	user, ok := r.users[id]
	if !ok {
		return nil, appError.ErrNotFound.New("user not found")
	}
	return user, nil
}

func (r *userMemoryRepository) Create(_ context.Context, user *entity.User) error {
	r.users[user.ID] = user
	return nil
}

func (r *userMemoryRepository) Update(_ context.Context, user *entity.User) error {
	if _, ok := r.users[user.ID]; !ok {
		return appError.ErrNotFound.New("user not found")
	}
	r.users[user.ID] = user
	return nil
}

func (r *userMemoryRepository) Delete(_ context.Context, id string) error {
	if _, ok := r.users[id]; !ok {
		return appError.ErrNotFound.New("user not found")
	}
	delete(r.users, id)
	return nil
}

func (r *userMemoryRepository) List(_ context.Context, limit uint64, cursor string) ([]*entity.User, error) {
	var users []*entity.User
	for _, user := range r.users {
		users = append(users, user)
	}

	slices.SortFunc(users, func(a, b *entity.User) int {
		if a.ID < b.ID {
			return -1
		}
		if a.ID > b.ID {
			return 1
		}
		return 0
	})

	var result []*entity.User
	for _, user := range users {
		if cursor != "" && user.ID <= cursor {
			continue
		}
		result = append(result, user)
		if uint64(len(result)) >= limit {
			break
		}
	}
	return result, nil
}
