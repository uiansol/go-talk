package main

import (
	"time"

	"github.com/go-chi/jwtauth"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ByTime int

const (
	newer ByTime = iota
	older
)

type Config struct {
	DB        *gorm.DB
	TokenAuth *jwtauth.JWTAuth
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

func (c *Config) CreateUser(name, password string) (*User, error) {
	hashedPassword, err := CreatePassword([]byte(password))
	if err != nil {
		return nil, err
	}

	user := User{Name: name, Password: string(hashedPassword)}
	if err := c.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *Config) CheckLogin(name, password string) error {
	var user User
	if err := c.DB.Where("name = ?", name).First(&user).Error; err != nil {
		return err
	}

	return CheckPassword([]byte(user.Password), []byte(password))
}

func (c *Config) CreateMessage(userName string, text string) (*Message, error) {
	message := Message{UserName: userName, Text: text}
	if err := c.DB.Create(&message).Error; err != nil {
		return nil, err
	}

	return &message, nil
}

func (c *Config) GetMessagesByTime(limit int, t time.Time, b ByTime) ([]Message, error) {
	var whereClause = "created_at > ?"
	var orderClause = "created_at ASC"
	if b == older {
		whereClause = "created_at < ?"
		orderClause = "created_at DESC"
	}

	query := c.DB.Table("messages").
		Where(whereClause, t).
		Order(orderClause).
		Limit(limit)

	var messages []Message
	if err := query.Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}
