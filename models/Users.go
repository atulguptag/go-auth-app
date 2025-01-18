package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name              string   `json:"name"`
	Email             string   `gorm:"unique" json:"email"`
	Password          string   `json:"password"`
	IsVerified        bool     `json:"is_verified" gorm:"default:false"`
	VerificationToken string   `json:"verification_token"`
	Prompts           []Prompt `json:"prompts"`
}
