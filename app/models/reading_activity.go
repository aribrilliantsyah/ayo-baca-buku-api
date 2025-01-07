package models

import (
	"time"

	"gorm.io/gorm"
)

type ReadingActivity struct {
	gorm.Model
	UserBookID  uint      `gorm:"not null"`
	PagesRead   int       `gorm:"not null"`
	StartPage   int       `gorm:"not null"`
	EndPage     int       `gorm:"not null"`
	Notes       string    `gorm:"type:text"`
	ReadingDate time.Time `gorm:"not null"`
	UserBook    UserBook  `gorm:"foreignKey:UserBookID"`
}
