package seeders

import (
	"ayo-baca-buku/app/models"
	"ayo-baca-buku/app/util/jwt"
	"log"
	"time"

	"gorm.io/gorm"
)

func SeedUser(db *gorm.DB) {
	password, err := jwt.HashPassword("rahasia")
	if err != nil {
		log.Fatal(err)
	}

	//check unique username & email
	if db.Where("username = ?", "sampleuser").Or("email = ?", "sampleuser@example.com").Find(&models.User{}).RowsAffected == 0 {
		db.Create(&models.User{
			Name:      "Sample User",
			Username:  "sampleuser",
			Email:     "sampleuser@example.com",
			Token:     "",
			Password:  password,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}
}
