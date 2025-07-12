package models

import (
	"time"

	"gorm.io/gorm"
)

type UserBook struct {
	ID                uint              `json:"id" gorm:"primarykey"`
	UserID            uint              `json:"user_id" gorm:"not null"`
	Title             string            `json:"title" gorm:"type:varchar(255);not null"`
	Author            string            `json:"author" gorm:"type:varchar(255);not null"`
	Publisher         string            `json:"publisher" gorm:"type:varchar(255)"`
	Cover             string            `json:"cover" gorm:"type:varchar(255)"` // URL atau path ke gambar cover
	TotalPages        int               `json:"total_pages" gorm:"not null"`
	CurrentPage       int               `json:"current_page" gorm:"default:0"`
	MotivationRead    string            `json:"motivation_read" gorm:"type:text"`
	Status            string            `json:"status" gorm:"type:varchar(20);check:status IN ('reading', 'finished');default:'reading'"`
	StartDate         time.Time         `json:"start_date" gorm:"not null"`
	EndDate           time.Time         `json:"end_date"`
	ReadingActivities []ReadingActivity `json:"reading_activities" gorm:"foreignKey:UserBookID"`
	User              User              `json:"user" gorm:"foreignKey:UserID"`
	CreatedAt         time.Time         `json:"created_at"`
	CreatedBy         int64             `json:"created_by"`
	UpdatedAt         time.Time         `json:"updated_at"`
	UpdatedBy         int64             `json:"updated_by"`
	DeletedAt         gorm.DeletedAt    `json:"deleted_at,omitempty" gorm:"index"`
	DeletedBy         int64             `json:"deleted_by,omitempty"`
}

// UserBookCreateRequest defines the structure for creating a new user book.
// UserID would typically come from the authenticated user context in a real app.
type UserBookCreateRequest struct {
	UserID         uint      `json:"user_id" validate:"required"`
	Title          string    `json:"title" validate:"required,min=1,max=255"`
	Author         string    `json:"author" validate:"required,min=1,max=255"`
	Publisher      string    `json:"publisher,omitempty" validate:"omitempty,max=255"`
	Cover          string    `json:"cover,omitempty" validate:"omitempty,url,max=255"`
	TotalPages     int       `json:"total_pages" validate:"required,gt=0"`
	MotivationRead string    `json:"motivation_read,omitempty"`
	StartDate      time.Time `json:"start_date" validate:"required"`
	// Status will default to 'reading' in the model or controller
}

// UserBookUpdateRequest defines the structure for updating an existing user book.
type UserBookUpdateRequest struct {
	Title          string    `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Author         string    `json:"author,omitempty" validate:"omitempty,min=1,max=255"`
	Publisher      string    `json:"publisher,omitempty" validate:"omitempty,max=255"`
	Cover          string    `json:"cover,omitempty" validate:"omitempty,url,max=255"`
	TotalPages     *int      `json:"total_pages,omitempty" validate:"omitempty,gt=0"` // Pointer to distinguish between 0 and not provided
	CurrentPage    *int      `json:"current_page,omitempty" validate:"omitempty,gte=0"`
	MotivationRead string    `json:"motivation_read,omitempty"`
	Status         string    `json:"status,omitempty" validate:"omitempty,oneof=reading finished"`
	StartDate      time.Time `json:"start_date,omitempty" validate:"omitempty,required"` // omitempty might not be ideal if you want to clear it, but for time.Time zero value is tricky.
	EndDate        time.Time `json:"end_date,omitempty"`                                 // omitempty is fine here. Consider *time.Time if clearing is needed and zero value is significant.
}
