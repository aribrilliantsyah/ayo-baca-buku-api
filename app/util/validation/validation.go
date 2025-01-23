package validation

import (
	"ayo-baca-buku/app/models"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func UniqueUsername(db *gorm.DB, id int64) validator.Func {
	return func(fl validator.FieldLevel) bool {
		username := fl.Field().String()

		var count int64
		query := db.Model(&models.User{}).Where("username = ?", username)

		// Jika ID disediakan, kecualikan data dengan ID tersebut
		if id > 0 {
			query = query.Where("id <> ?", id)
		}

		query.Count(&count)
		return count == 0
	}
}

func UniqueEmail(db *gorm.DB, id int64) validator.Func {
	return func(fl validator.FieldLevel) bool {
		email := fl.Field().String()

		var count int64
		query := db.Model(&models.User{}).Where("email = ?", email)

		// Jika ID disediakan, kecualikan data dengan ID tersebut
		if id > 0 {
			query = query.Where("id <> ?", id)
		}

		query.Count(&count)
		return count == 0
	}
}
