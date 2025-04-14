package models

import "time"

type User struct {
	ID uint `gorm:"primary_key" json:"id"`

	Username string `gorm:"unique" json:"username"`
	Email    string `gorm:"unique" json:"email"`

	Password string `json:"-"`

	IsAdmin bool `json:"is_admin" gorm:"default:false"`

	CreatedAt time.Time `json:"created_at"`
}
