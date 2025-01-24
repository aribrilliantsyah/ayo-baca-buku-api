package routes

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type UserBookController struct {
	DB       *gorm.DB
	Validate *validator.Validate
}

func NewUserBookController(DB *gorm.DB) *UserBookController {
	return &UserBookController{
		DB:       DB,
		Validate: validator.New(),
	}
}
