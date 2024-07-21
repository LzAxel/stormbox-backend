package postgresql

import (
	"chat-backend/internal/apperror"
	"chat-backend/internal/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{
		db: db,
	}
}

func (u *UserPostgres) Create(ctx context.Context, user model.CreateUserDTO) (model.User, error) {
	var createdUser model.User

	query, args, err := squirrel.
		Insert(UsersTable).
		Columns("username", "password_hash", "created_at", "updated_at", "online_at").
		Values(user.Username, user.PasswordHash, user.CreatedAt, user.UpdatedAt, user.UpdatedAt).
		SuffixExpr(squirrel.Expr("RETURNING *")).
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return model.User{}, err
	}

	err = u.db.GetContext(ctx, &createdUser, query, args...)
	if err != nil {
		if IsPgErrorCode(err, pgerrcode.UniqueViolation) {
			return model.User{}, apperror.ErrUsernameAlreadyExists
		}
		return model.User{}, fmt.Errorf("UserPostgres.Create: %w", err)
	}

	return createdUser, nil
}

func (u *UserPostgres) GetByID(ctx context.Context, id uint64) (model.User, error) {
	var user model.User

	query, args, err := squirrel.Select("*").From(UsersTable).Where(squirrel.Eq{"id": id}).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return model.User{}, err
	}
	err = u.db.GetContext(ctx, &user, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, apperror.ErrUserNotFound
		}
		return model.User{}, fmt.Errorf("UserPostgres.GetByID: %w", err)
	}

	return user, err
}

func (u *UserPostgres) GetByUsername(ctx context.Context, username string) (model.User, error) {
	var user model.User

	query, args, err := squirrel.Select("*").From(UsersTable).Where(squirrel.Eq{"username": username}).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return model.User{}, err
	}
	err = u.db.GetContext(ctx, &user, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, apperror.ErrUserNotFound
		}
		return model.User{}, fmt.Errorf("UserPostgres.GetByUsername: %w", err)
	}

	return user, err
}

func (u *UserPostgres) GetAll(ctx context.Context, pagination model.Pagination) ([]model.User, uint64, error) {
	var (
		total uint64
		users = make([]model.User, 0)
	)

	query, args, err := squirrel.Select("*").From(UsersTable).Limit(pagination.Limit).Offset(pagination.Offset).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return []model.User{}, 0, err
	}
	err = u.db.SelectContext(ctx, &users, query, args...)
	if err != nil {
		return []model.User{}, 0, fmt.Errorf("UserPostgres.GetAll: getting users: %w", err)
	}

	query, _, err = squirrel.Select("COUNT(*)").From(UsersTable).PlaceholderFormat(squirrel.Dollar).ToSql()
	err = u.db.GetContext(ctx, &total, query)

	if err != nil {
		return []model.User{}, 0, fmt.Errorf("UserPostgres.GetAll: getting users count: %w", err)
	}

	return users, total, err
}

func (u *UserPostgres) Delete(ctx context.Context, id uint64) error {
	query, args, err := squirrel.Delete(UsersTable).Where(squirrel.Eq{"id": id}).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}
	_, err = u.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("UserPostgres.Delete: %w", err)
	}
	return nil
}
