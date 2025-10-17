package entity

import (
	"errors"
	"money-tracker/internal/dto"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          uint      `json:"id"`
	UserName    string    `json:"username"`
	Password    string    `json:"-"`
	LevelManage int8      `json:"level_manage"`
	StatusID    uint      `json:"status_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at,omitempty"`
}

func NewUser(user_name string, password string, level_manage int8, status_id uint) (*User, error) {
	if user_name == "" || level_manage == 0 {
		return nil, errors.New("name and level manage are required")
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	return &User{
		UserName:    user_name,
		LevelManage: level_manage,
		Password:    hashedPassword,
		StatusID:    status_id,
		CreatedAt:   time.Now(),
	}, nil
}

type UserRepository interface {
	Insert(user *User) error
	FindById(id uint) (*User, error)
	FindByUserName(username string) (*User, error)
	FindAll(input dto.ListUsersInput) ([]User, int, error)
	Update(user *User) (*User, error)
	Delete(id uint) error
}

// -----------------functions
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
