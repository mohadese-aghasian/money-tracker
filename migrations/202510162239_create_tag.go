package migrations

import (
	"fmt"
	"money-tracker/internal/repository"

	"github.com/go-gormigrate/gormigrate/v2"

	"gorm.io/gorm"
)

func CreateTagTable(tx *gorm.DB) error {
	if !tx.Migrator().HasTable(&repository.Tag{}) {
		fmt.Println("creating table tag..")
		if err := tx.Migrator().CreateTable(&repository.Tag{}); err != nil {
			return err
		}
		fmt.Println("tag table created successfully!")
	}
	return nil
}
func dropTagTable(tx *gorm.DB) error {
	if tx.Migrator().HasTable(&repository.Tag{}) {
		fmt.Println("Dropping 'tag.' table...")
		if err := tx.Migrator().DropTable(&repository.Tag{}); err != nil {
			return err
		}
		fmt.Println("'tag' table dropped successfully!")
	}
	return nil
}

func CreateTagMigrate() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "202510162239_create_tag",
		Migrate: func(tx *gorm.DB) error {
			return CreateTagTable(tx)
		},
		Rollback: func(tx *gorm.DB) error {
			return dropTagTable(tx)
		},
	}
}
