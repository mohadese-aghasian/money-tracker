package migrations

import (
	"fmt"
	"money-tracker/internal/repository"

	"github.com/go-gormigrate/gormigrate/v2"

	"gorm.io/gorm"
)

func CreateCatTable(tx *gorm.DB) error {
	if !tx.Migrator().HasTable(&repository.Category{}) {
		fmt.Println("creating table category..")
		if err := tx.Migrator().CreateTable(&repository.Category{}); err != nil {
			return err
		}
		fmt.Println("category table created successfully!")
	}
	return nil
}
func dropCatTable(tx *gorm.DB) error {
	if tx.Migrator().HasTable(&repository.Category{}) {
		fmt.Println("Dropping 'cat.' table...")
		if err := tx.Migrator().DropTable(&repository.Category{}); err != nil {
			return err
		}
		fmt.Println("'category' table dropped successfully!")
	}
	return nil
}

func CreateCatMigrate() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "202508261430_createCategory",
		Migrate: func(tx *gorm.DB) error {
			return CreateCatTable(tx)
		},
		Rollback: func(tx *gorm.DB) error {
			return dropCatTable(tx)
		},
	}
}
