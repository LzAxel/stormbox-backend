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

type Repository struct {
	logger logger.Logger
	User
}

func New(psql *sqlx.DB, logger logger.Logger) *Repository {
	return &Repository{
		logger: logger,
		User:   postgresql.NewUserPostgres(psql),
	}
}
