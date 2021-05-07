package main

import (
	"gorm.io/gorm"
)

// Database model for a user account.
type User struct {
	gorm.Model
	Name     string `gorm:"uniqueIndex"`
	Password string
	Messages []Message `gorm:"foreignKey:UserName"`
}

// Database model for a message.
type Message struct {
	ID        uint   `json:"-"`
	UserName  string `json:"user_name"`
	Text      string `json:"text"`
	CreatedAt int64  `json:"created_at" gorm:"index,autoCreateTime:milli"`
}
