package entity

import (
	"errors"
	"money-tracker/internal/constants"
	"money-tracker/internal/dto"
	"time"
)

type Purchase struct {
	ID            uint              `json:"id"`
	Date          time.Time         `json:"date"`
	Reason        string            `json:"reason"`
	StatusID      uint              `json:"status_id"`
	Color         string            `json:"color"`
	Amount        int64             `json:"amount"`
	Method        int8              `json:"method"`
	TagIDs        string            `json:"tag_ids"`
	Note          string            `json:"note"`
	Category      *Category         `json:"category"`
	CategoryId    *uint             `json:"category_id"`
	SubCategoryId *uint             `json:"sub_category_id"`
	SubCategory   *Category         `json:"sub_category"`
	Details       constants.JSONMap `json:"details"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	DeletedAt     time.Time         `json:"deleted_at,omitempty"`
}

func NewPurchase(amount int64, date time.Time, category_id *uint, status_id uint) (*Purchase, error) {
	if amount == 0 || category_id == nil {
		return nil, errors.New("title is required")
	}

	return &Purchase{
		Amount:     amount,
		StatusID:   status_id,
		Date:       date,
		CategoryId: category_id,
		CreatedAt:  time.Now(),
	}, nil
}

type PurchaseRepository interface {
	Insert(purchase *Purchase) error
	FindById(id uint, status_id []uint) (*Purchase, error)
	FindAll(input dto.PurchaseFindAll) ([]Purchase, int, error)
	Update(p *Purchase) (*Purchase, error)
	Delete(id uint) error
	// Search(input dto.SearchWithTagInput) ([]dto.FetchedSearchCategory, int64, error)
}
