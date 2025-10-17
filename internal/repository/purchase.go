package repository

import (
	"encoding/json"
	"fmt"
	"money-tracker/internal/constants"
	"money-tracker/internal/dto"
	"money-tracker/internal/entity"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Purchase struct {
	ID            uint `gorm:"primaryKey"`
	Date          time.Time
	Amount        int64
	Reason        string         `gorm:"size:255"`
	StatusID      uint           `gorm:"default:1;not null"`
	Color         string         `gorm:"size:100"`
	Method        int8           `gorm:"size:255"`
	TagIDs        string         `gorm:"size:255"`
	Note          string         `gorm:"type:text"`
	Category      *Category      `gorm:"foreignKey:CategoryId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	CategoryId    *uint          `gorm:"index"`
	SubCategoryId *uint          `gorm:"index"`
	SubCategory   *Category      `gorm:"foreignKey:SubCategoryId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Details       datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time `gorm:"index"`
}

// /-------------------------------------------

type PurchaseRepo struct {
	db *gorm.DB
}

func NewPurchaseRepo(db *gorm.DB) *PurchaseRepo {
	return &PurchaseRepo{db: db}
}

func (rep PurchaseRepo) Insert(purchase *entity.Purchase) error {
	p := ToRepoPurchase(purchase)
	return rep.db.Create(p).Error
}

func (rep PurchaseRepo) FindById(id uint, status_id []uint) (*entity.Purchase, error) {
	var purchase Purchase
	if err := rep.db.Where("status_id IN ?", status_id).First(&purchase, id).Error; err != nil {
		return nil, err
	}
	return purchase.ToEntityPurchase(), nil
}

func (rep PurchaseRepo) Delete(id uint) error {
	now := time.Now()

	return rep.db.Model(&Purchase{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status_id":  0,
			"deleted_at": now,
		}).Error
}
func (rep PurchaseRepo) Update(purchase *entity.Purchase) (*entity.Purchase, error) {
	purchase.UpdatedAt = time.Now()
	dbQ := ToRepoPurchase(purchase)
	if err := rep.db.Save(dbQ).Error; err != nil {
		return nil, err
	}
	return purchase, nil
}

func (rep PurchaseRepo) FindAll(input dto.PurchaseFindAll) ([]entity.Purchase, int, error) {
	query := rep.db.Model(&Purchase{})

	// --- Filters ---
	if input.CategoryID != nil && *input.CategoryID > 0 {
		query = query.Where("category_id = ?", *input.CategoryID)
	}
	if input.SubCategoryID != nil && *input.SubCategoryID > 0 {
		query = query.Where("sub_category_id = ?", *input.SubCategoryID)
	}
	if input.Reason != "" {
		query = query.Where("reason ILIKE ?", "%"+input.Reason+"%")
	}
	if input.Note != "" {
		query = query.Where("note ILIKE ?", "%"+input.Note+"%")
	}
	if input.Color != "" {
		query = query.Where("color = ?", input.Color)
	}
	if input.Method != 0 {
		query = query.Where("method = ?", input.Method)
	}
	if input.ID > 0 {
		query = query.Where("id = ?", input.ID)
	}
	if input.Amount > 0 {
		query = query.Where("amount = ?", input.Amount)
	}

	// --- Default active status ---
	if input.StatusID == 0 {
		query = query.Where("status_id = ?", 1)
	} else {
		query = query.Where("status_id = ?", input.StatusID)
	}

	// --- Tags filtering (if applicable) ---
	for _, id := range input.TagIDs {
		query = query.Where("? = ANY(string_to_array(tag_ids, ','))", fmt.Sprintf("%d", id))
	}

	// --- Count before pagination ---
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// --- Columns to select ---
	columns := []string{
		"id", "date", "amount", "reason", "status_id", "color",
		"method", "tag_ids", "note", "category_id", "sub_category_id", "details",
	}

	if input.OtherFields {
		columns = append(columns, "created_at", "updated_at")
	}

	// --- Execute query with preloads ---
	var items []Purchase
	result := query.
		Select(columns).
		Offset(input.Start).
		Limit(input.Limit).
		Order(input.OrderBy+" "+input.Sort).
		Preload("Category", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "slug", "color").Where("status_id = ?", 1)
		}).
		Preload("SubCategory", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "slug", "color").Where("status_id = ?", 1)
		}).
		Find(&items)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	// --- Convert to entities ---
	var purchases []entity.Purchase
	for _, dbF := range items {
		p := dbF.ToEntityPurchase()
		purchases = append(purchases, *p)
	}

	return purchases, int(count), nil
}

// ///------------------------------------------------------------
func (m *Purchase) ToEntityPurchase() *entity.Purchase {
	var det constants.JSONMap
	if len(m.Details) > 0 {
		_ = json.Unmarshal(m.Details, &det)
	}

	var category *entity.Category
	if m.Category != nil {
		category = m.Category.ToEntityCategory()
	}

	var subcat *entity.Category
	if m.SubCategory != nil {
		subcat = m.SubCategory.ToEntityCategory()
	}

	return &entity.Purchase{
		ID:            m.ID,
		Date:          m.Date,
		StatusID:      m.StatusID,
		Amount:        m.Amount,
		Color:         m.Color,
		Reason:        m.Reason,
		Method:        m.Method,
		Note:          m.Note,
		CategoryId:    m.CategoryId,
		Category:      category,
		SubCategoryId: m.SubCategoryId,
		SubCategory:   subcat,
		Details:       det,
		TagIDs:        m.TagIDs,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
		DeletedAt:     m.DeletedAt,
	}
}

// Convert entity → repository model
func ToRepoPurchase(e *entity.Purchase) *Purchase {
	var det datatypes.JSON
	if e.Details != nil {
		b, _ := json.Marshal(e.Details)
		det = b
	}
	var category *Category
	if e.Category != nil {
		category = ToRepoCategory(e.Category)
	}

	var subcat *Category
	if e.SubCategory != nil {
		subcat = ToRepoCategory(e.SubCategory)
	}
	return &Purchase{
		ID:            e.ID,
		Reason:        e.Reason,
		StatusID:      e.StatusID,
		Date:          e.Date,
		Color:         e.Color,
		TagIDs:        e.TagIDs,
		Note:          e.Note,
		Amount:        e.Amount,
		Method:        e.Method,
		Category:      category,
		CategoryId:    e.CategoryId,
		SubCategory:   subcat,
		SubCategoryId: e.SubCategoryId,
		Details:       det,
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
		DeletedAt:     e.DeletedAt,
	}
}
