package model

import "time"

type Friendship struct {
	ID        uint
	UserID    uint
	FriendID  uint
	CreatedAt time.Time
}

type CreateFriendshipInput struct {
	UserID   uint
	FriendID uint
}
