package model

import "time"

type User struct {
	ID           uint64    `db:"id"`
	Username     string    `db:"username"`
	PasswordHash []byte    `db:"password_hash"`
	OnlineAt     time.Time `db:"online_at"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func (u *User) ToView() ViewUser {
	return ViewUser{
		ID:        u.ID,
		Username:  u.Username,
		OnlineAt:  u.OnlineAt,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

type ViewUser struct {
	ID        uint64    `json:"id"`
	Username  string    `json:"username"`
	OnlineAt  time.Time `json:"online_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserInput struct {
	Username string
	Password string
}

type CreateUserDTO struct {
	Username     string
	PasswordHash []byte
	OnlineAt     time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
