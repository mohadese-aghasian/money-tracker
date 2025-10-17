package migrations

import (
	"fmt"
	"money-tracker/internal/repository"

	"github.com/go-gormigrate/gormigrate/v2"

	"gorm.io/gorm"
)

func CreateUserTokenTable(tx *gorm.DB) error {
	if !tx.Migrator().HasTable(&repository.UserToken{}) {
		fmt.Println("creating table user Token..")
		if err := tx.Migrator().CreateTable(&repository.UserToken{}); err != nil {
			return err
		}
		fmt.Println("User Token table created successfully!")
	}
	return nil
}
func dropUserTokenTable(tx *gorm.DB) error {
	if tx.Migrator().HasTable(&repository.UserToken{}) {
		fmt.Println("Dropping 'User Token.' table...")
		if err := tx.Migrator().DropTable(&repository.UserToken{}); err != nil {
			return err
		}
		fmt.Println("'User Token' table dropped successfully!")
	}
	return nil
}

func CreateUserTokenMigrate() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "202509021201_create_user_token",
		Migrate: func(tx *gorm.DB) error {
			return CreateUserTokenTable(tx)
		},
		Rollback: func(tx *gorm.DB) error {
			return dropUserTokenTable(tx)
		},
	}
}
