package repository

import (
	"money-tracker/internal/entity"
	"time"

	"gorm.io/gorm"
)

type UserToken struct {
	ID        uint   `gorm:"primaryKey"`
	Token     string `gorm:"type:text"`
	UserID    uint   `gorm:"index;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type UserTokenRepoGormPostgres struct {
	db *gorm.DB
}

func NewUserTokenRepositoryGorm(db *gorm.DB) *UserTokenRepoGormPostgres {
	return &UserTokenRepoGormPostgres{db: db}
}

func (rep UserTokenRepoGormPostgres) Insert(user_token *entity.UserToken) error {
	return rep.db.Create(user_token).Error
}

func (rep UserTokenRepoGormPostgres) Delete(id uint) error {
	return rep.db.Model(&UserToken{}).Where("id = ?", id).Update("status_id", 0).Error
}

func (rep UserTokenRepoGormPostgres) FindByToken(token string, status_id []uint) (*entity.UserToken, error) {
	var t entity.UserToken
	if err := rep.db.Where("token = ? AND status_id IN ?", token, status_id).First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}
