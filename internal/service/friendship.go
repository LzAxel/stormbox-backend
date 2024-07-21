package service

import (
	"chat-backend/internal/apperror"
	"chat-backend/internal/model"
	"chat-backend/internal/repository"
	"context"
)

type FriendshipService struct {
	repo repository.Friendship
}

func NewFriendshipService(repo repository.Friendship) *FriendshipService {
	return &FriendshipService{repo: repo}
}

func (f *FriendshipService) GetByUserID(ctx context.Context, pagination model.Pagination, userID uint64) ([]model.ViewUser, model.FullPagination, error) {
	friends, total, err := f.repo.GetAllByUserID(ctx, pagination, userID)
	if err != nil {
		if apperror.IsAppError(err) {
			return []model.ViewUser{}, model.FullPagination{}, err
		}

		return []model.ViewUser{}, model.FullPagination{}, apperror.NewDatabaseError("FriendshipService.GetByUserID", err)
	}

	var viewFriendships = make([]model.ViewUser, 0)
	for _, friend := range friends {
		viewFriendships = append(viewFriendships, friend.ToView())
	}

	return viewFriendships, pagination.ToFullPagination(total, len(viewFriendships)), nil
}

func (f *FriendshipService) Create(ctx context.Context, friendship model.CreateFriendshipDTO) error {
	if friendship.UserID == friendship.FriendID {
		return apperror.ErrCannotFriendSelf
	}
	_, err := f.repo.Create(ctx, friendship)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}

		return apperror.NewDatabaseError("FriendshipService.Create", err)
	}

	return nil
}
