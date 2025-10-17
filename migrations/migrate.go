package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func RunMigration(db *gorm.DB) *gormigrate.Gormigrate {

	return gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{

		CreateCatMigrate(),
		CreateUserMigrate(),
		CreateUserTokenMigrate(),
		CreateTagMigrate(),
		CreatePurchaseMigrate(),
	})

}
