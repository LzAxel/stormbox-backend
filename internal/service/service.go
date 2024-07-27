package service

import (
	"chat-backend/internal/jwt"
	"chat-backend/internal/model"
	"chat-backend/internal/repository"
	"context"
)

type Authorization interface {
	RefreshTokens(ctx context.Context, refreshToken string) (jwt.TokenPair, error)
	Login(ctx context.Context, input model.LoginInput) (model.ViewUser, jwt.TokenPair, error)
	Register(ctx context.Context, input model.CreateUserInput) (model.ViewUser, jwt.TokenPair, error)
}

type User interface {
	GetByID(ctx context.Context, id uint64) (model.ViewUser, error)
	GetByUsername(ctx context.Context, username string) (model.ViewUser, error)
	GetAll(ctx context.Context, pagination model.Pagination) ([]model.ViewUser, model.FullPagination, error)
}

type Friendship interface {
	GetByUserID(ctx context.Context, pagination model.Pagination, userID uint64) ([]model.ViewUser, model.FullPagination, error)
	Create(ctx context.Context, friendship model.CreateFriendshipDTO) error
}

type Message interface {
	Create(ctx context.Context, message model.CreateMessageInput) (model.Message, error)
	GetAllWithFriend(ctx context.Context, pagination model.Pagination, userID, friendID uint64) ([]model.ViewMessage, model.FullPagination, error)
}

type Services struct {
	User
	Authorization
	Friendship
	Message
}

func New(
	repository *repository.Repository,
	jwt *jwt.JWT,
) *Services {
	return &Services{
		User:          NewUserService(repository.User),
		Authorization: NewAuthorizationService(jwt, repository.User),
		Friendship:    NewFriendshipService(repository.Friendship),
		Message:       NewMessageService(repository.Message, repository.Friendship),
	}
}
