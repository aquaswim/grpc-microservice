package postgres

import (
	"context"
	"errors"
	appError "gaman-microservice/user-service/internal/domain/app_error"
	"gaman-microservice/user-service/internal/domain/entity"
	"gaman-microservice/user-service/internal/port/out"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db      *pgxpool.Pool
	builder squirrel.StatementBuilderType
}

func NewUserRepository(db *pgxpool.Pool) out.UserRepository {
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
		return nil, appError.ErrInternal.WrapWithNoMessage(err)
	}

	user := &entity.User{}
	err = r.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appError.ErrNotFound.Wrap(err, "user not found")
		}
		return nil, appError.ErrInternal.WrapWithNoMessage(err)
	}

	return user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	query, args, err := r.builder.Select("id", "username", "password", "email").
		From("users").
		Where(squirrel.Eq{"email": email}).
		ToSql()
	if err != nil {
		return nil, appError.ErrInternal.WrapWithNoMessage(err)
	}

	user := &entity.User{}
	err = r.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appError.ErrNotFound.Wrap(err, "user not found")
		}
		return nil, appError.ErrInternal.WrapWithNoMessage(err)
	}

	return user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	query, args, err := r.builder.Select("id", "username", "password", "email").
		From("users").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, appError.ErrInternal.WrapWithNoMessage(err)
	}

	user := &entity.User{}
	err = r.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appError.ErrNotFound.Wrap(err, "user not found")
		}
		return nil, appError.ErrInternal.WrapWithNoMessage(err)
	}

	return user, nil
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	query, args, err := r.builder.Insert("users").
		Columns("id", "username", "password", "email").
		Values(user.ID, user.Username, user.Password, user.Email).
		ToSql()
	if err != nil {
		return appError.ErrInternal.WrapWithNoMessage(err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return appError.ErrInternal.Wrap(err, "failed to insert user")
	}

	return nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	query, args, err := r.builder.Update("users").
		Set("username", user.Username).
		Set("email", user.Email).
		Where(squirrel.Eq{"id": user.ID}).
		ToSql()
	if err != nil {
		return appError.ErrInternal.WrapWithNoMessage(err)
	}

	res, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return appError.ErrInternal.WrapWithNoMessage(err)
	}

	rows := res.RowsAffected()
	if rows == 0 {
		return appError.ErrNotFound.New("user not found")
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	query, args, err := r.builder.Delete("users").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return appError.ErrInternal.WrapWithNoMessage(err)
	}

	res, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return appError.ErrInternal.WrapWithNoMessage(err)
	}

	rows := res.RowsAffected()
	if rows == 0 {
		return appError.ErrNotFound.New("user not found")
	}

	return nil
}

func (r *userRepository) List(ctx context.Context, limit uint64, cursor string) ([]*entity.User, error) {
	builder := r.builder.Select("id", "username", "password", "email").
		From("users").
		OrderBy("id ASC").
		Limit(limit)

	if cursor != "" {
		builder = builder.Where(squirrel.Gt{"id": cursor})
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, appError.ErrInternal.WrapWithNoMessage(err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, appError.ErrInternal.WrapWithNoMessage(err)
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		user := &entity.User{}
		err = rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email)
		if err != nil {
			return nil, appError.ErrInternal.WrapWithNoMessage(err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, appError.ErrInternal.WrapWithNoMessage(err)
	}

	return users, nil
}
