package controllers

import "gorm.io/gorm"

type UserBookController struct {
	DB *gorm.DB
}

func NewUserBookController(DB *gorm.DB) *UserBookController {
	return &UserBookController{
		DB: DB,
	}
}
