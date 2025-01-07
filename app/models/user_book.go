package models

import (
	"time"

	"gorm.io/gorm"
)

type UserBook struct {
	gorm.Model
	UserID            uint      `gorm:"not null"`
	Title             string    `gorm:"type:varchar(255);not null"`
	Author            string    `gorm:"type:varchar(255);not null"`
	Publisher         string    `gorm:"type:varchar(255)"`
	Cover             string    `gorm:"type:varchar(255)"` // URL atau path ke gambar cover
	TotalPages        int       `gorm:"not null"`
	CurrentPage       int       `gorm:"default:0"`
	MotivationRead    string    `gorm:"type:text"`
	Status            string    `gorm:"type:varchar(20);check:status IN ('reading', 'finished');default:'reading'"`
	StartDate         time.Time `gorm:"not null"`
	EndDate           *time.Time
	ReadingActivities []ReadingActivity `gorm:"foreignKey:UserBookID"`
	User              User              `gorm:"foreignKey:UserID"`
}
