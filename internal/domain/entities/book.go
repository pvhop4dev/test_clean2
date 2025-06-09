package entities

import (
	"time"


	"gorm.io/gorm"
)

type Book struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"size:255;not null"`
	Description string    `json:"description" gorm:"type:text"`
	Author      string    `json:"author" gorm:"size:100;not null"`
	UserID      string    `json:"user_id" gorm:"not null"`
	User        User      `json:"-" gorm:"foreignKey:UserID"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Book) TableName() string {
	return "books"
}
