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
	Role      string         `json:"role" gorm:"type:varchar(255)"`
	UserBooks []UserBook     `json:"user_books" gorm:"foreignKey:UserID"`
	CreatedAt time.Time      `json:"created_at"`
	CreatedBy int64          `json:"created_by"`
	UpdatedAt time.Time      `json:"updated_at"`
	UpdatedBy int64          `json:"updated_by"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	DeletedBy int64          `json:"deleted_by"`
}

type UserCreateRequest struct {
	Name                 string `json:"name" validate:"required"`
	Username             string `json:"username" validate:"required,alphanum,unique_username"`
	Email                string `json:"email" validate:"required,email,unique_email"`
	Password             string `json:"password" validate:"required,alphanum,min=6"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password,min=6"`
	Role                 string `json:"role" validate:"omitempty,oneof=admin user"` // Role is optional, defaults to 'user' if not provided or handled in controller
}

type UserUpdateRequest struct {
	Name                 string `json:"name" validate:"omitempty"`
	Username             string `json:"username" validate:"omitempty,alphanum,unique_username"`
	Email                string `json:"email" validate:"omitempty,email,unique_email"`
	Password             string `json:"password" validate:"omitempty,alphanum,min=6"`
	PasswordConfirmation string `json:"password_confirmation" validate:"omitempty,eqfield=Password,min=6"`
	Role                 string `json:"role" validate:"omitempty,oneof=admin user"`
}
