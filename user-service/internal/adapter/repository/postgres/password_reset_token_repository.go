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

type passwordResetTokenRepository struct {
	db      *pgxpool.Pool
	builder squirrel.StatementBuilderType
}

func NewPasswordResetTokenRepository(db *pgxpool.Pool) out.PasswordResetTokenRepository {
	return &passwordResetTokenRepository{
		db:      db,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *passwordResetTokenRepository) Create(ctx context.Context, token *entity.PasswordResetToken) error {
	query, args, err := r.builder.Insert("user_password_reset_tokens").
		Columns("user_id", "token", "expires_at", "created_at").
		Values(token.UserID, token.Token, token.ExpiresAt, token.CreatedAt).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return appError.ErrInternal.WrapWithNoMessage(err)
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&token.ID)
	if err != nil {
		return appError.ErrInternal.Wrap(err, "failed to insert reset token")
	}

	return nil
}

func (r *passwordResetTokenRepository) FindByToken(ctx context.Context, token string) (*entity.PasswordResetToken, error) {
	query, args, err := r.builder.Select("id", "user_id", "token", "expires_at", "created_at").
		From("user_password_reset_tokens").
		Where(squirrel.Eq{"token": token}).
		ToSql()
	if err != nil {
		return nil, appError.ErrInternal.WrapWithNoMessage(err)
	}

	res := &entity.PasswordResetToken{}
	err = r.db.QueryRow(ctx, query, args...).Scan(&res.ID, &res.UserID, &res.Token, &res.ExpiresAt, &res.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appError.ErrNotFound.Wrap(err, "reset token not found")
		}
		return nil, appError.ErrInternal.WrapWithNoMessage(err)
	}

	return res, nil
}

func (r *passwordResetTokenRepository) DeleteByUserID(ctx context.Context, userID string) error {
	query, args, err := r.builder.Delete("user_password_reset_tokens").
		Where(squirrel.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return appError.ErrInternal.WrapWithNoMessage(err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return appError.ErrInternal.WrapWithNoMessage(err)
	}

	return nil
}

func (r *passwordResetTokenRepository) Delete(ctx context.Context, token string) error {
	query, args, err := r.builder.Delete("user_password_reset_tokens").
		Where(squirrel.Eq{"token": token}).
		ToSql()
	if err != nil {
		return appError.ErrInternal.WrapWithNoMessage(err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return appError.ErrInternal.WrapWithNoMessage(err)
	}

	return nil
}
