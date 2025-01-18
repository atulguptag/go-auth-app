package models

import "gorm.io/gorm"

type Prompt struct {
	gorm.Model
	UserID uint   `json:"user_id"`
	Text   string `json:"text"`
	User   User   `gorm:"foreignKey:UserID"`
}
