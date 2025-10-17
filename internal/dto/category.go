package dto

import (
	"time"
)

type ListCategoriesInput struct {
	Start    int    `form:"start"`
	Limit    int    `form:"limit"`
	OrderBy  string `form:"order_by"`
	Sort     string `form:"sort"`
	Title    string `form:"title"`
	ID       uint   `form:"id"`
	StatusID uint   `form:"status_id"`
	Slug     string `form:"slug"`
	Color    string `json:"color"`
}

type CreateCategoryRequest struct {
	Title    string `json:"title" binding:"required"`
	StatusID uint   `json:"status_id" binding:"oneof=1 0"`
	Color    string `json:"color"`
}

type UpdateCategoryRequest struct {
	ID       uint   `json:"id" binding:"required"`
	Title    string `json:"title"`
	StatusID uint   `json:"status_id" binding:"oneof=1 0"`
	TagIDs   string `json:"tag_ids"`
	// Slug     string `json:"slug"`
	Color string `json:"color"`
}
type UpdateCategoryInput struct {
	ID       uint   `json:"id" binding:"required"`
	Title    string `json:"title"`
	StatusID uint   `json:"status_id" binding:"oneof=1 0"`
	TagIDs   string `json:"tag_ids"`
	Slug     string `json:"slug"`
	Color    string `json:"color"`
}

type FetchedParent struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	StatusID  uint      `json:"status_id"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}

type CategoryResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	StatusID  uint      `json:"status_id"`
	CreatedAt time.Time `json:"created_at"`
}

type AddCategoryInput struct {
	Title    string `json:"title" binding:"required"`
	StatusID uint   `json:"status_id"`
	Slug     string `json:"slug" binding:"required"`
	Color    string `json;"color"`
}

type CategoryFindAll struct {
	Start    int    `form:"start" json:"start"`         // Pagination start index
	Limit    int    `form:"limit" json:"limit"`         // Number of articles to return
	OrderBy  string `form:"order_by" json:"order_by"`   // Field to order by
	Sort     string `form:"sort" json:"sort"`           // Sort direction (asc|desc)
	Title    string `form:"title" json:"title"`         // Filter by title
	ID       uint   `form:"id" json:"id"`               // Filter by ID
	StatusID uint   `form:"status_id" json:"status_id"` // Filter by status ID
	TagIDs   []int  `form:"tag_ids[]" json:"tag_ids"`   // Filter by tag IDs (array)
	Slug     string `form:"slug" json:"slug"`
	Color    string `form:"color"`
}
