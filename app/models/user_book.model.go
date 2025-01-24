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
	DeletedAt         gorm.DeletedAt    `json:"deleted_at" gorm:"index"`
	DeletedBy         int64             `json:"deleted_by"`
}
