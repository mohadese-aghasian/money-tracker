package entity

import (
	"errors"
	"time"
)

type Tag struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	StatusID  uint      `json:"status_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"-"`
	DeletedAt time.Time `json:"-"`
}

func NewTag(title string, status_id uint) (*Tag, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}

	return &Tag{
		Title:     title,
		StatusID:  status_id,
		CreatedAt: time.Now(),
	}, nil
}

type TagRepository interface {
	Insert(tag *Tag) error
	Update(tag *Tag) (*Tag, error)
	FindById(id uint) (*Tag, error)
	FindByTitle(title string) (*Tag, error)
	FindAll(start int, limit int, orderBy string, sort string, id uint, status_id uint, title string) ([]Tag, int, error)
	Delete(id uint) error
	FindByIDs(ids []uint) ([]Tag, error)
}