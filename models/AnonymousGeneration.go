package models

import (
	"time"

	"gorm.io/gorm"
)

type AnonymousGeneration struct {
	ID                 uint `gorm:"primaryKey"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          gorm.DeletedAt `gorm:"index"`
	AnonymousID        string         `gorm:"unique;column:anonymous_id"`
	GenerationCount    int            `gorm:"column:generation_count"`
	LastGenerationTime time.Time      `gorm:"column:last_generation_time"`
}
