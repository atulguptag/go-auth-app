package models

import (
	"time"

	"gorm.io/gorm"
)

type AnonymousGeneration struct {
	gorm.Model
	IPAddress          string `gorm:"index"`
	GenerationCount    int    `gorm:"default:0"`
	LastGenerationTime time.Time
}
