package models

import (
	"time"

	"gorm.io/gorm"
)

type ReadingActivity struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	UserBookID  uint           `json:"user_book_id" gorm:"not null"`
	PagesRead   int            `json:"pages_read" gorm:"not null"`
	StartPage   int            `json:"start_page" gorm:"not null"`
	EndPage     int            `json:"end_page" gorm:"not null"`
	Notes       string         `json:"notes" gorm:"type:text"`
	ReadingDate time.Time      `json:"reading_date" gorm:"not null"`
	UserBook    UserBook       `json:"user_book" gorm:"foreignKey:UserBookID"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
