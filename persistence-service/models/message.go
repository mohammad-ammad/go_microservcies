package models

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	UserID    uint      `json:"userId"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}
