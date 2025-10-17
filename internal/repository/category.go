package repository

import (
	"money-tracker/internal/dto"
	"money-tracker/internal/entity"
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID       uint   `gorm:"primaryKey"`
	Title    string `gorm:"size:255"`
	StatusID uint   `gorm:"default:1;not null"`
	TagIDs    string    `gorm:"size:255"`
	Slug      string `gorm:"size:300";index`
	Color     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time `gorm:"index"`
}

type RepoGormPostgres struct {
	db *gorm.DB
}

func NewRepositoryGorm(db *gorm.DB) *RepoGormPostgres {
	return &RepoGormPostgres{db: db}
}

func (rep RepoGormPostgres) Insert(category *entity.Category) error {
	c := ToRepoCategory(category)
	return rep.db.Create(c).Error
}

func (rep RepoGormPostgres) FindById(id uint) (*entity.Category, error) {
	var category Category
	if err := rep.db.Where("status_id = ?", 1).First(&category, id).Error; err != nil {
		return nil, err
	}
	return category.ToEntityCategory(), nil
}
func (rep RepoGormPostgres) FindBySlug(slug string, status_id []uint) (*entity.Category, error) {
	var category Category
	if err := rep.db.Where("slug = ? AND status_id IN ?", slug, status_id).First(&category).Error; err != nil {
		return nil, err
	}
	return category.ToEntityCategory(), nil
}

func (rep RepoGormPostgres) FindAll(input dto.CategoryFindAll) ([]entity.Category, int, error) {
	query := rep.db.Model(&Category{})

	// Title filter
	if input.Title != "" {
		query = query.Where("title ILIKE ?", "%"+input.Title+"%")
	}

	// ID filter
	if input.ID > 0 {
		query = query.Where("id = ?", input.ID)
	}

	// Slug filter
	if input.Slug != "" {
		query = query.Where("slug = ?", input.Slug)
	}

	// Status filter (default to 1 if not provided)
	if input.StatusID == 0 {
		query = query.Where("status_id = ?", 1)
	} else {
		query = query.Where("status_id = ?", input.StatusID)
	}

	// if input.CoverId != nil && *input.CoverId > 0 {
	// 	query = query.Where("cover_id = ?", input.CoverId)
	// }

	// // Tag IDs filter
	// for _, id := range input.TagIDs {
	// 	// Assuming tag_ids is a comma-separated string column
	// 	query = query.Where("? = ANY(string_to_array(tag_ids, ','))", fmt.Sprintf("%d", id))
	// }

	// Count total
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Pagination + sorting
	orderBy := input.OrderBy
	if orderBy == "" {
		orderBy = "id" // default order
	}

	sortDir := input.Sort
	if sortDir == "" {
		sortDir = "ASC"
	}

	var items []Category
	result :=
		query.
			Distinct().
			Offset(input.Start).
			Limit(input.Limit).
			Order(orderBy + " " + sortDir).
			// Preload("Cover").
			Find(&items)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	var categories []entity.Category
	for _, dbF := range items {
		a := dbF.ToEntityCategory()
		categories = append(categories, *a)
	}

	return categories, int(count), nil
}

func (rep RepoGormPostgres) Delete(id uint) error {
	return rep.db.Model(&Category{}).Where("id = ?", id).Update("status_id", 0).Error
}

func (rep RepoGormPostgres) Update(category *entity.Category) (*entity.Category, error) {
	if err := rep.db.Save(ToRepoCategory(category)).Error; err != nil {
		return nil, err
	}
	return category, nil
}

// ///------------------------------------------------------------
func (m *Category) ToEntityCategory() *entity.Category {

	// var cover *entity.File
	// if m.Cover != nil {
	// 	cover = m.Cover.ToEntityFile()
	// }

	return &entity.Category{
		ID:        m.ID,
		Title:     m.Title,
		StatusID:  m.StatusID,
		Slug:      m.Slug,
		Color:     m.Color,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: m.DeletedAt,
	}
}

// Convert entity → repository model
func ToRepoCategory(e *entity.Category) *Category {

	// var cover *File
	// if e.Cover != nil {
	// 	cover = ToRepoFile(e.Cover)
	// }

	return &Category{
		ID:       e.ID,
		Title:    e.Title,
		StatusID: e.StatusID,
		Slug:     e.Slug,
		Color:    e.Color,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
		DeletedAt: e.DeletedAt,
	}
}
