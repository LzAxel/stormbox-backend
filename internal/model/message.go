package model

import "time"

type Message struct {
	ID          uint64
	SenderID    uint64
	RecipientID uint64
	Content     string
	CreatedAt   time.Time
}

type CreateMessageInput struct {
	SenderID    uint64
	RecipientID uint64
	Content     string
}
