package repository

import (
	"money-tracker/internal/dto"
	"money-tracker/internal/entity"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"size:255"`
	UserName    string `gorm:"size:255"`
	Email       string `gorm:"size:255"`
	Password    string `gorm:"size:255"`
	Mobile      string `gorm:"size:255"`
	LevelManage int8
	StatusID    uint `gorm:"default:1;not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time `gorm:"index"`
}

type UserRepoGormPostgres struct {
	db *gorm.DB
}

func NewUserRepositoryGorm(db *gorm.DB) *UserRepoGormPostgres {
	return &UserRepoGormPostgres{db: db}
}

func (rep UserRepoGormPostgres) Insert(user *entity.User) error {
	return rep.db.Create(user).Error
}

func (rep UserRepoGormPostgres) FindById(id uint) (*entity.User, error) {
	var user entity.User
	if err := rep.db.Where("status_id = ?", 1).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (rep UserRepoGormPostgres) Delete(id uint) error {
	return rep.db.Model(&entity.User{}).Where("id = ?", id).Update("status_id", 0).Error
}

func (rep UserRepoGormPostgres) Update(user *entity.User) (*entity.User, error) {
	if err := rep.db.Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (rep UserRepoGormPostgres) FindAll(input dto.ListUsersInput) ([]entity.User, int, error) {
	query := rep.db.Model(&entity.User{})

	// Filters

	if input.UserName != "" {
		query = query.Where("user_name ILIKE ?", "%"+input.UserName+"%")
	}
	if input.ID > 0 {
		query = query.Where("id = ?", input.ID)
	}
	if input.LevelManage > 0 {
		query = query.Where("level_manage = ?", input.LevelManage)
	}

	// Count
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Pagination and sorting defaults
	orderBy := input.OrderBy
	if orderBy == "" {
		orderBy = "id"
	}
	sort := input.Sort
	if sort == "" {
		sort = "ASC"
	}

	// Query results
	var users []entity.User
	result := query.
		Select("id", "user_name", "level_manage", "status_id").
		Offset(input.Start).
		Limit(input.Limit).
		Order(orderBy + " " + sort).
		Find(&users)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return users, int(count), nil
}

func (rep UserRepoGormPostgres) FindByUserName(username string) (*entity.User, error) {
	var user entity.User
	if err := rep.db.Where("status_id = ?", 1).Where("user_name = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// --------------------------------
func (m *User) ToEntityUser() *entity.User {
	return &entity.User{
		ID:          m.ID,
		UserName:    m.UserName,
		StatusID:    m.StatusID,
		LevelManage: m.LevelManage,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		// DeletedAt:   m.DeletedAt,
	}
}

// Convert entity → repository model
func ToRepoUser(e *entity.User) *User {
	return &User{
		ID:          e.ID,
		UserName:    e.UserName,
		StatusID:    e.StatusID,
		LevelManage: e.LevelManage,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
		// DeletedAt:   e.DeletedAt,
	}
}
