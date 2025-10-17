package entity

import (
	"errors"

	"money-tracker/internal/dto"
	"time"
)

type Category struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	StatusID  uint      `json:"status_id"`
	Slug      string    `json:"slug"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
}

func NewCategory(title string, slug string, status_id uint, color string) (*Category, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}

	return &Category{
		Title:     title,
		StatusID:  status_id,
		Slug:      slug,
		Color:     color,
		CreatedAt: time.Now(),
	}, nil
}

type CategoryRepository interface {
	Insert(category *Category) error
	FindById(id uint) (*Category, error)
	FindBySlug(slug string, status_id []uint) (*Category, error)
	FindAll(input dto.CategoryFindAll) ([]Category, int, error)
	Update(category *Category) (*Category, error)
	Delete(id uint) error
	// Search(input dto.SearchWithTagInput) ([]dto.FetchedSearchCategory, int64, error)
}
