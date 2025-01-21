package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	UID       string         `json:"uid" gorm:"type:uuid;default:gen_random_uuid()"`
	Name      string         `json:"name" gorm:"type:varchar(255);not null"`
	Username  string         `json:"username" gorm:"type:varchar(100);uniqueIndex;not null"`
	Email     string         `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	Token     string         `json:"token" gorm:"type:varchar(255)"`
	Password  string         `json:"-" gorm:"type:varchar(255);not null"`
	UserBooks []UserBook     `json:"user_books" gorm:"foreignKey:UserID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
