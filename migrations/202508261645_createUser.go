package migrations

import (
	"fmt"
	"money-tracker/internal/repository"

	"github.com/go-gormigrate/gormigrate/v2"

	"gorm.io/gorm"
)

func CreateUserTable(tx *gorm.DB) error {
	if !tx.Migrator().HasTable(&repository.User{}) {
		fmt.Println("creating table user..")
		if err := tx.Migrator().CreateTable(&repository.User{}); err != nil {
			return err
		}
		fmt.Println("User table created successfully!")
	}
	return nil
}
func dropUserTable(tx *gorm.DB) error {
	if tx.Migrator().HasTable(&repository.User{}) {
		fmt.Println("Dropping 'User.' table...")
		if err := tx.Migrator().DropTable(&repository.User{}); err != nil {
			return err
		}
		fmt.Println("'User' table dropped successfully!")
	}
	return nil
}

func CreateUserMigrate() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "202508261645_createUser",
		Migrate: func(tx *gorm.DB) error {
			return CreateUserTable(tx)
		},
		Rollback: func(tx *gorm.DB) error {
			return dropUserTable(tx)
		},
	}
}
