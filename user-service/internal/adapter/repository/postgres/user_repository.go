package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gaman-microservice/user-service/internal/domain/entity"
	"gaman-microservice/user-service/internal/port/out"

	"github.com/Masterminds/squirrel"
)

type userRepository struct {
	db      *sql.DB
	builder squirrel.StatementBuilderType
}

func NewUserRepository(db *sql.DB) out.UserRepository {
	return &userRepository{
		db:      db,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	query, args, err := r.builder.Select("id", "username", "password", "email").
		From("users").
		Where(squirrel.Eq{"username": username}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	user := &entity.User{}
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	return user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	query, args, err := r.builder.Select("id", "username", "password", "email").
		From("users").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	user := &entity.User{}
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	return user, nil
}
