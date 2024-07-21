package postgresql

import (
	"chat-backend/internal/apperror"
	"chat-backend/internal/model"
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
)

type FriendshipPostgres struct {
	db *sqlx.DB
}

func NewFriendshipPostgres(db *sqlx.DB) *FriendshipPostgres {
	return &FriendshipPostgres{db: db}
}

func (f *FriendshipPostgres) Create(ctx context.Context, friendship model.CreateFriendshipDTO) (model.Friendship, error) {
	var createdFriendship model.Friendship

	query, args, err := squirrel.
		Insert(FriendshipsTable).
		Columns("user_id", "friend_id", "created_at").
		Values(friendship.UserID, friendship.FriendID, friendship.CreatedAt).
		SuffixExpr(squirrel.Expr("RETURNING *")).
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return model.Friendship{}, err
	}

	err = f.db.GetContext(ctx, &createdFriendship, query, args...)
	if err != nil {
		if IsPgErrorCode(err, pgerrcode.UniqueViolation) {
			return model.Friendship{}, apperror.ErrFriendshipAlreadyExists
		}
		if IsPgErrorCode(err, pgerrcode.ForeignKeyViolation) {
			return model.Friendship{}, apperror.ErrUserNotFound
		}
		return model.Friendship{}, fmt.Errorf("FriendshipPostgres.Create: %w", err)
	}
	return createdFriendship, nil
}

func (f *FriendshipPostgres) GetAllByUserID(ctx context.Context, pagination model.Pagination, userID uint64) ([]model.User, uint64, error) {
	var (
		total   uint64
		friends = make([]model.User, 0)
	)

	query, args, err := squirrel.Select("u.*").
		From(UsersTable + " u").
		Join(FriendshipsTable + " f ON u.id = f.user_id OR u.id = f.friend_id").
		Where(squirrel.And{
			squirrel.Or{squirrel.Eq{"f.user_id": userID}, squirrel.Eq{"f.friend_id": userID}},
			squirrel.NotEq{"u.id": userID},
		}).
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return []model.User{}, 0, err
	}
	err = f.db.SelectContext(ctx, &friends, query, args...)
	if err != nil {
		return []model.User{}, 0, fmt.Errorf("FriendshipPostgres.GetAllByUserID: getting friends: %w", err)
	}

	query, _, err = squirrel.Select("COUNT(u.*)").
		From(UsersTable + " u").
		Join(FriendshipsTable + " f ON u.id = f.user_id OR u.id = f.friend_id").
		Where(squirrel.And{
			squirrel.Or{squirrel.Eq{"f.user_id": userID}, squirrel.Eq{"f.friend_id": userID}},
			squirrel.NotEq{"u.id": userID},
		}).
		PlaceholderFormat(squirrel.Dollar).ToSql()
	err = f.db.GetContext(ctx, &total, query, args...)
	if err != nil {
		return []model.User{}, 0, fmt.Errorf("FriendshipPostgres.GetAllByUserID: getting friends count: %w", err)
	}
	return friends, total, err
}

func (f *FriendshipPostgres) Delete(ctx context.Context, userID, friendID uint64) error {
	query, args, err := squirrel.Delete(FriendshipsTable).Where(squirrel.And{
		squirrel.Eq{"user_id": userID},
		squirrel.Eq{"friend_id": friendID},
	}).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("FriendshipPostgres.Delete: %w", err)
	}
	_, err = f.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("FriendshipPostgres.Delete: %w", err)
	}
	return nil
}
