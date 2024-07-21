package model

import "time"

type Message struct {
	ID          uint64    `db:"id"`
	SenderID    uint64    `db:"sender_id"`
	RecipientID uint64    `db:"recipient_id"`
	Content     string    `db:"content"`
	CreatedAt   time.Time `db:"created_at"`
}

func (m *Message) ToView() ViewMessage {
	return ViewMessage{
		ID:          m.ID,
		SenderID:    m.SenderID,
		RecipientID: m.RecipientID,
		Content:     m.Content,
		CreatedAt:   m.CreatedAt,
	}
}

type ViewMessage struct {
	ID          uint64    `json:"id"`
	SenderID    uint64    `json:"sender_id"`
	RecipientID uint64    `json:"recipient_id"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateMessageInput struct {
	SenderID    uint64
	RecipientID uint64
	Content     string
}

type CreateMessageDTO struct {
	SenderID    uint64
	RecipientID uint64
	Content     string
	CreatedAt   time.Time
}
