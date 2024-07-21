package postgresql

import (
	"chat-backend/internal/model"
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type MessagePostgres struct {
	db *sqlx.DB
}

func NewMessagePostgres(db *sqlx.DB) *MessagePostgres {
	return &MessagePostgres{db: db}
}

func (m *MessagePostgres) Create(ctx context.Context, message model.CreateMessageDTO) (model.Message, error) {
	var createdMessage model.Message

	query, args, err := squirrel.Insert(MessagesTable).
		Columns("sender_id", "recipient_id", "content", "created_at").
		Values(message.SenderID, message.RecipientID, message.Content, message.CreatedAt).
		Suffix("RETURNING *").
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return model.Message{}, fmt.Errorf("MessagePostgres.Create: %w", err)
	}

	err = m.db.GetContext(ctx, &createdMessage, query, args...)
	if err != nil {
		return model.Message{}, fmt.Errorf("MessagePostgres.Create: %w", err)
	}
	return createdMessage, nil
}
func (m *MessagePostgres) GetAllWithFriend(ctx context.Context, pagination model.Pagination, userID, friendID uint64) ([]model.Message, uint64, error) {
	var (
		total    uint64
		messages = make([]model.Message, 0)
	)

	query, args, err := squirrel.Select("*").
		From(MessagesTable).
		Where(
			squirrel.Or{
				squirrel.And{
					squirrel.Eq{"sender_id": userID},
					squirrel.Eq{"recipient_id": friendID},
				},
				squirrel.And{
					squirrel.Eq{"sender_id": friendID},
					squirrel.Eq{"recipient_id": userID},
				}}).
		OrderBy("created_at DESC").
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return []model.Message{}, 0, err
	}
	err = m.db.SelectContext(ctx, &messages, query, args...)
	if err != nil {
		return []model.Message{}, 0, fmt.Errorf("MessagePostgres.GetAllWithFriend: getting messages: %w", err)
	}

	query, _, err = squirrel.Select("COUNT(*)").
		From(MessagesTable).
		Where(
			squirrel.Or{
				squirrel.And{
					squirrel.Eq{"sender_id": userID},
					squirrel.Eq{"recipient_id": friendID},
				},
				squirrel.And{
					squirrel.Eq{"sender_id": friendID},
					squirrel.Eq{"recipient_id": userID},
				}}).
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return []model.Message{}, 0, err
	}
	err = m.db.GetContext(ctx, &total, query, args...)
	if err != nil {
		return []model.Message{}, 0, fmt.Errorf("MessagePostgres.GetAllWithFriend: getting total: %w", err)
	}

	return messages, total, nil
}
