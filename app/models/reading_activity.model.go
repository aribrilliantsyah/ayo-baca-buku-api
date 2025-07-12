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
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// ReadingActivityCreateRequest defines the payload for creating a reading activity.
type ReadingActivityCreateRequest struct {
	UserBookID  uint      `json:"user_book_id" validate:"required"`
	PagesRead   int       `json:"pages_read" validate:"required,gt=0"`
	StartPage   int       `json:"start_page" validate:"required,gte=0"`
	EndPage     int       `json:"end_page" validate:"required,gtfield=StartPage"`
	Notes       string    `json:"notes,omitempty"`
	ReadingDate time.Time `json:"reading_date" validate:"required"`
}

// ReadingActivityUpdateRequest defines the payload for updating a reading activity.
type ReadingActivityUpdateRequest struct {
	PagesRead   *int      `json:"pages_read,omitempty" validate:"omitempty,gt=0"`
	StartPage   *int      `json:"start_page,omitempty" validate:"omitempty,gte=0"`
	EndPage     *int      `json:"end_page,omitempty" validate:"omitempty,gtfield=StartPage"`
	Notes       string    `json:"notes,omitempty"`
	ReadingDate time.Time `json:"reading_date,omitempty" validate:"omitempty,required"`
}
