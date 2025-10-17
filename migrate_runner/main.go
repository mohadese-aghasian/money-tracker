package main

import (
	"flag"
	"fmt"
	"log"
	"money-tracker/internal/config"
	"money-tracker/migrations"
)

func init() {
	config.ConnectToDB()
}

func main() {
	rollback := flag.Bool("rollback", false, "rollback last migration")
	flag.Parse()

	// config.ConnectToDB()

	m := migrations.RunMigration(config.DB) // return *gormigrate.Gormigrate

	if *rollback {
		if err := m.RollbackLast(); err != nil {
			log.Fatalf("Rollback failed: %v", err)
		}
		fmt.Println("✅ Last migration rolled back successfully!")
		return
	}

	if err := m.Migrate(); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	fmt.Println("✅ All migrations applied successfully!")
}
