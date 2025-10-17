package migrations

import (
	"fmt"
	"money-tracker/internal/repository"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// ---------- Create Table ----------
func createPurchaseTable(tx *gorm.DB) error {
	if !tx.Migrator().HasTable(&repository.Purchase{}) {
		fmt.Println("Creating table 'purchase'...")
		if err := tx.Migrator().CreateTable(&repository.Purchase{}); err != nil {
			return err
		}
		fmt.Println("✅ 'purchase' table created successfully!")
	}
	return nil
}

// ---------- Drop Table ----------
func dropPurchaseTable(tx *gorm.DB) error {
	if tx.Migrator().HasTable(&repository.Purchase{}) {
		fmt.Println("Dropping table 'purchase'...")
		if err := tx.Migrator().DropTable(&repository.Purchase{}); err != nil {
			return err
		}
		fmt.Println("🗑️  'purchase' table dropped successfully!")
	}
	return nil
}

// ---------- Migration Definition ----------
func CreatePurchaseMigrate() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20251016_createPurchase",
		Migrate: func(tx *gorm.DB) error {
			return createPurchaseTable(tx)
		},
		Rollback: func(tx *gorm.DB) error {
			return dropPurchaseTable(tx)
		},
	}
}
