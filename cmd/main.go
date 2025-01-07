package main

import (
	"ayo-baca-buku/app/config"
	"ayo-baca-buku/app/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connectToPostgreSQL() (*gorm.DB, error) {
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

func InitializeMigration(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.UserBook{},
		&models.ReadingActivity{},
	)
}

func main() {
	db, err := connectToPostgreSQL()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	err = InitializeMigration(db)
	if err != nil {
		log.Fatal(err)
	}
}
