package dto

import "time"

type PurchaseFindAll struct {
	ID            uint   `form:"id"`
	CategoryID    *uint  `form:"category_id"`
	SubCategoryID *uint  `form:"sub_category_id"`
	Reason        string `form:"reason"`
	Note          string `form:"note"`
	Color         string `form:"color"`
	Method        int8   `form:"method"`
	Amount        int64  `form:"amount"`
	StatusID      uint   `form:"status_id"`
	TagIDs        []uint `form:"tag_ids"`
	Start         int    `form:"start" default:"0"`
	Limit         int    `form:"limit" default:"10"`
	OrderBy       string `form:"order_by" default:"id"`
	Sort          string `form:"sort" default:"desc"`
	OtherFields   bool   `form:"other_fields"`
}

type AddPurchaseInput struct {
	CategoryId    *uint     `json:"category_id" binding:"required"`
	SubCategoryId *uint     `json:"sub_category_id`
	Reason        string    `json:"reason"`
	Date          time.Time `json:"date" binding:"required"`
	Note          string    `json:"note"`
	Color         string    `json:"color"`
	Method        int8      `json:"method"`
	Amount        int64     `json:"amount"`
	StatusID      uint      `json:"status_id" binding:"oneof=1 0"`
	TagIDs        string    `json:"tag_ids"`
}

type PurchaseResponse struct {
	ID          uint              `json:"id"`
	Category    *CategoryResponse `json:"category"`
	SubCategory *CategoryResponse `json:"sub_category"`
	Reason      string            `json:"reason"`
	Date        time.Time         `json:"date"`
	Note        string            `json:"note"`
	Color       string            `json:"color"`
	Method      int8              `json:"method"`
	Amount      int64             `json:"amount"`
	StatusID    uint              `json:"status_id"`
	Tags        []FetchedTag      `json:"tags"`
	CreatedAt   time.Time         `json:"created_at"`
}

type UpdatePurchaseInput struct {
	ID            uint      `json:"id" binding:"required"`
	CategoryId    *uint     `json:"category_id"`
	SubCategoryId *uint     `json:"sub_category_id"`
	Reason        string    `json:"reason"`
	Date          time.Time `json:"date"`
	Note          string    `json:"note"`
	Color         string    `json:"color"`
	Method        int8      `json:"method"`
	Amount        int64     `json:"amount"`
	StatusID      uint      `json:"status_id" binding:"oneof=1 0"`
	TagIDs        string    `json:"tag_ids"`
}
