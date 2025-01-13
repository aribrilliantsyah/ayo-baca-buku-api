package database

import (
	"ayo-baca-buku/app/config"
	"ayo-baca-buku/app/database/seeders"
	"ayo-baca-buku/app/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase() (*gorm.DB, error) {
	appConfig, err := config.LoadAppConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	dbSource := appConfig.DB_SOURCE
	db, err := gorm.Open(postgres.Open(dbSource), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func RunMigration(DB *gorm.DB) {
	err := DB.AutoMigrate(
		&models.User{},
		&models.UserBook{},
		&models.ReadingActivity{},
	)

	if err != nil {
		log.Fatal("Failed to migrate...")
	}

	fmt.Println("Migrated Successfully")
}

func RunSeeder(DB *gorm.DB) {
	seeders.SeedUser(DB)
	fmt.Println("Seeder Successfully")
}
