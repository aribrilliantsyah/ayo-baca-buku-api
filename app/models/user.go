package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name      string     `gorm:"type:varchar(255);not null"`
	Username  string     `gorm:"type:varchar(100);uniqueIndex;not null"`
	Email     string     `gorm:"type:varchar(255);uniqueIndex;not null"`
	Token     string     `gorm:"type:varchar(255)"`
	Password  string     `gorm:"type:varchar(255);not null"`
	UserBooks []UserBook `gorm:"foreignKey:UserID"`
}
