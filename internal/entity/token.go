package entity

import (
	"errors"
	"time"
)

type UserToken struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

func NewUserToken(token string, user_id uint) (*UserToken, error) {
	if token == "" || user_id == 0 {
		return nil, errors.New("user_id and token are required")
	}

	return &UserToken{
		UserID:    user_id,
		Token:     token,
		CreatedAt: time.Now(),
	}, nil
}

type UserTokenRepository interface {
	Insert(user_token *UserToken) error
	FindByToken(token string, status_id []uint) (*UserToken, error)
	Delete(id uint) error
}
