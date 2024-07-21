package service

import (
	"chat-backend/internal/jwt"
	"chat-backend/internal/model"
	"chat-backend/internal/repository"
	"context"
)

type Authorization interface {
	RefreshTokens(ctx context.Context, refreshToken string) (jwt.TokenPair, error)
	Login(ctx context.Context, input model.LoginInput) (jwt.TokenPair, error)
	Register(ctx context.Context, input model.CreateUserInput) error
}

type User interface {
	GetByID(ctx context.Context, id uint64) (model.ViewUser, error)
	GetByUsername(ctx context.Context, username string) (model.ViewUser, error)
	GetAll(ctx context.Context, pagination model.Pagination) ([]model.ViewUser, model.FullPagination, error)
}

type Services struct {
	User
	Authorization
}

func New(
	repository *repository.Repository,
	jwt *jwt.JWT,
) *Services {
	return &Services{
		User:          NewUserService(repository.User),
		Authorization: NewAuthorizationService(jwt, repository.User),
	}
}
