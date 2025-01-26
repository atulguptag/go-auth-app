package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name              string   `json:"name"`
	Email             string   `gorm:"uniqueIndex" json:"email"`
	Password          string   `json:"-"`
	IsVerified        bool     `json:"is_verified" gorm:"default:false"`
	VerificationToken string   `json:"verification_token"`
	Provider          string   `json:"provider"`
	GoogleID          string   `json:"google_id" gorm:"uniqueIndex;null"`
	ImageURL          string   `json:"image_url"`
	Prompts           []Prompt `json:"prompts"`
}
