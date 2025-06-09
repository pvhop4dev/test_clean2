package entities

import (
	"time"

	"gorm.io/gorm"
)

type Translation struct {
	gorm.Model
	SourceText     string    `gorm:"not null"`
	TargetLang     string    `gorm:"not null"`
	TranslatedText string    `gorm:"not null"`
	LastAccessed   time.Time `gorm:"index"`
}
