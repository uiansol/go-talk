package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Returns a database connection.
func DBConnect(dsn string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
}

// Handles database migrations.
func (c *Config) EnsureDBSetup() error {
	err := c.DB.AutoMigrate(&User{})
	if err != nil {
		return err
	}

	return c.DB.AutoMigrate(&Message{})
}

// Saves a user to the database after bcrypting the user's password.
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

// Checks a user's password is correct.
func (c *Config) CheckLogin(name, password string) error {
	var user User
	if err := c.DB.Where("name = ?", name).First(&user).Error; err != nil {
		return err
	}

	return CheckPassword([]byte(user.Password), []byte(password))
}

// Saves a message in the database.
func (c *Config) CreateMessage(userName string, text string) (*Message, error) {
	message := Message{UserName: userName, Text: text}
	if err := c.DB.Create(&message).Error; err != nil {
		return nil, err
	}

	return &message, nil
}

// Used to determine the query order for messages.
type ByTime int

const (
	newer ByTime = iota
	older
)

// Returns a given number of messages from the database.
// If newer, the order will be oldest to newest.
// If older, the order be newest to oldest.
func (c *Config) GetMessagesByTime(limit int, unixTime int64, b ByTime) ([]Message, error) {
	var whereClause = "created_at > ?"
	var orderClause = "created_at asc"
	if b == older {
		whereClause = "created_at < ?"
		orderClause = "created_at desc"
	}

	query := c.DB.Table("messages").
		Where(whereClause, unixTime).
		Order(orderClause).
		Limit(limit)

	var messages []Message
	if err := query.Find(&messages).Error; err != nil {
		return nil, err
	}

	return messages, nil
}
