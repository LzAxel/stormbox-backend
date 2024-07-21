package service

import (
	"chat-backend/internal/apperror"
	"chat-backend/internal/model"
	"chat-backend/internal/repository"
	"context"
	"time"
)

type MessageService struct {
	repo       repository.Message
	friendRepo repository.Friendship
}

func NewMessageService(repository repository.Message, friendRepo repository.Friendship) *MessageService {
	return &MessageService{
		repo:       repository,
		friendRepo: friendRepo,
	}
}

func (s *MessageService) Create(ctx context.Context, message model.CreateMessageInput) (model.Message, error) {
	isFriends, err := s.friendRepo.IsFriends(ctx, message.SenderID, message.RecipientID)
	if err != nil {
		return model.Message{}, apperror.NewDatabaseError("MessageService.Create: checking if friends", err)
	}
	if !isFriends {
		return model.Message{}, apperror.ErrCanSendMessagesOnlyFriends
	}

	return s.repo.Create(ctx, model.CreateMessageDTO{
		SenderID:    message.SenderID,
		RecipientID: message.RecipientID,
		Content:     message.Content,
		CreatedAt:   time.Now().UTC(),
	})
}

func (s *MessageService) GetAllWithFriend(ctx context.Context, pagination model.Pagination, userID, friendID uint64) ([]model.ViewMessage, model.FullPagination, error) {
	messages, total, err := s.repo.GetAllWithFriend(ctx, pagination, userID, friendID)
	if err != nil {
		return []model.ViewMessage{}, model.FullPagination{}, apperror.NewDatabaseError("MessageService.GetAllWithFriend: getting messages", err)
	}

	var viewMessages = make([]model.ViewMessage, 0)
	for _, message := range messages {
		viewMessages = append(viewMessages, message.ToView())
	}

	return viewMessages, pagination.ToFullPagination(total, len(viewMessages)), nil
}
