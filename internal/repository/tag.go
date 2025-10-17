package repository

import (
	"money-tracker/internal/entity"
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"size:100"`
	StatusID  uint   `gorm:"default:1;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type TagRepoGorm struct {
	db *gorm.DB
}

func NewTagRepoGorm(db *gorm.DB) *TagRepoGorm {
	return &TagRepoGorm{db: db}
}

func (rep *TagRepoGorm) Insert(tag *entity.Tag) error {
	return rep.db.Create(tag).Error
}

func (rep TagRepoGorm) Delete(id uint) error {
	return rep.db.Model(&Tag{}).Where("id = ?", id).Update("status_id", 0).Error
}

func (rep *TagRepoGorm) Update(tag *entity.Tag) (*entity.Tag, error) {
	tag.UpdatedAt = time.Now()
	if err := rep.db.Save(tag).Error; err != nil {
		return nil, err
	}
	return tag, nil
}

func (rep *TagRepoGorm) FindById(id uint) (*entity.Tag, error) {
	var tag entity.Tag
	if err := rep.db.Where("status_id = ?", 1).First(&tag, id).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (rep TagRepoGorm) FindByTitle(title string) (*entity.Tag, error) {
	var tag entity.Tag
	if err := rep.db.Where("status_id = ?", 1).Where("title = ?", title).First(&tag).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (rep TagRepoGorm) FindAll(start int, limit int, orderBy string, sort string, id uint, status_id uint, title string) ([]entity.Tag, int, error) {

	query := rep.db.Model(&entity.Tag{})

	if title != "" {
		query = query.Where("title ILIKE ?", "%"+title+"%")
	}
	// Filter by ID
	if id > 0 {
		query = query.Where("id = ?", id)
	}

	if status_id == 0 {
		query = query.Where("status_id = ?", 1)
	} else {
		query = query.Where("status_id = ?", status_id)
	}

	var items []entity.Tag
	var count int64

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	result :=
		query.
			Distinct().
			Offset(start).
			Limit(limit).
			Order(orderBy + " " + sort).
			Find(&items)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return items, int(count), nil

}

func (rep TagRepoGorm) FindByIDs(ids []uint) ([]entity.Tag, error) {
	var tags []entity.Tag
	err := rep.db.Where("id IN ?", ids).Find(&tags).Error
	return tags, err
}