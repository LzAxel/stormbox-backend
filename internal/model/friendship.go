package model

import "time"

type Friendship struct {
	ID        uint64    `db:"id"`
	UserID    uint64    `db:"user_id"`
	FriendID  uint64    `db:"friend_id"`
	CreatedAt time.Time `db:"created_at"`
}

func (f *Friendship) ToView() ViewFriendship {
	return ViewFriendship{
		ID:        f.ID,
		UserID:    f.UserID,
		FriendID:  f.FriendID,
		CreatedAt: f.CreatedAt,
	}
}

type ViewFriendship struct {
	ID        uint64    `json:"id"`
	UserID    uint64    `json:"user_id"`
	FriendID  uint64    `json:"friend_id"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateFriendshipInput struct {
	UserID   uint64
	FriendID uint64
}

type CreateFriendshipDTO struct {
	UserID    uint64
	FriendID  uint64
	CreatedAt time.Time
}
