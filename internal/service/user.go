package service

import (
	"chat-backend/internal/apperror"
	"chat-backend/internal/model"
	"chat-backend/internal/repository"
	"context"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repository repository.User) *UserService {
	return &UserService{
		repo: repository,
	}
}

func (u *UserService) GetByID(ctx context.Context, id uint64) (model.ViewUser, error) {
	user, err := u.repo.GetByID(ctx, id)
	if err != nil {
		if !apperror.IsAppError(err) {
			return model.ViewUser{}, apperror.NewDatabaseError("User.GetByID: %w", err)
		}
		return model.ViewUser{}, err
	}

	return user.ToView(), nil
}

func (u *UserService) GetByUsername(ctx context.Context, username string) (model.ViewUser, error) {
	user, err := u.repo.GetByUsername(ctx, username)
	if err != nil {
		if !apperror.IsAppError(err) {
			return model.ViewUser{}, apperror.NewDatabaseError("User.GetByUsername: %w", err)
		}
		return model.ViewUser{}, err
	}

	return user.ToView(), nil
}

func (u *UserService) GetAll(ctx context.Context, pagination model.Pagination) ([]model.ViewUser, model.FullPagination, error) {
	users, total, err := u.repo.GetAll(ctx, pagination)
	if err != nil {
		if !apperror.IsAppError(err) {
			return []model.ViewUser{}, model.FullPagination{}, apperror.NewDatabaseError("User.GetAll: %w", err)
		}
		return []model.ViewUser{}, model.FullPagination{}, err
	}

	var viewUsers = make([]model.ViewUser, len(users))
	for i, user := range users {
		viewUsers[i] = user.ToView()
	}

	return viewUsers, pagination.ToFullPagination(total, len(users)), nil
}
