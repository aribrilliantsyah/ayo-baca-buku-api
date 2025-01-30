package controllers

import "gorm.io/gorm"

type ReadingActivityController struct {
	DB *gorm.DB
}

func NewReadingActivityController(DB *gorm.DB) *ReadingActivityController {
	return &ReadingActivityController{
		DB: DB,
	}
}
