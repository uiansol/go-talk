package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	DB *gorm.DB
}

type User struct {
	gorm.Model
	Name     string `gorm:"uniqueIndex"`
	Password string
	Messages []Message `gorm:"foreignKey:UserName"`
}

type Message struct {
	ID        uint   `json:"-"`
	UserName  string `json:"user_name"`
	Text      string `json:"text"`
	CreatedAt int64  `json:"created_at" gorm:"index,autoCreateTime:milli"`
}

func DBConnect(dsn string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
}

func (c *Config) EnsureDBSetup() error {
	err := c.DB.AutoMigrate(&User{})
	if err != nil {
		return err
	}

	return c.DB.AutoMigrate(&Message{})
}
