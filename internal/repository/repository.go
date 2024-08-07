package repository

import (
	"chat-backend/internal/logger"
	"chat-backend/internal/model"
	"chat-backend/internal/repository/postgresql"
	"context"
	"github.com/jmoiron/sqlx"
)

type User interface {
	Create(ctx context.Context, user model.CreateUserDTO) (model.User, error)
	GetByID(ctx context.Context, id uint64) (model.User, error)
	GetByUsername(ctx context.Context, username string) (model.User, error)
	GetAll(ctx context.Context, pagination model.Pagination) ([]model.User, uint64, error)
	Delete(ctx context.Context, id uint64) error
}

type Friendship interface {
	Create(ctx context.Context, friendship model.CreateFriendshipDTO) (model.Friendship, error)
	GetAllByUserID(ctx context.Context, pagination model.Pagination, userID uint64) ([]model.User, uint64, error)
	IsFriends(ctx context.Context, userID, friendID uint64) (bool, error)
}

type Message interface {
	Create(ctx context.Context, message model.CreateMessageDTO) (model.Message, error)
	GetAllWithFriend(ctx context.Context, pagination model.Pagination, userID, friendID uint64) ([]model.Message, uint64, error)
}

type Repository struct {
	logger logger.Logger
	User
	Friendship
	Message
}

func New(psql *sqlx.DB, logger logger.Logger) *Repository {
	return &Repository{
		logger:     logger,
		User:       postgresql.NewUserPostgres(psql),
		Friendship: postgresql.NewFriendshipPostgres(psql),
		Message:    postgresql.NewMessagePostgres(psql),
	}
}
