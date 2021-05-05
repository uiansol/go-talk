package main

import (
	"context"

	"github.com/go-chi/jwtauth"
	"golang.org/x/crypto/bcrypt"
)

func CreatePassword(password []byte) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return hashedPassword, err
}

func CheckPassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}

func (c *Config) MakeToken(name string) string {
	_, tokenString, _ := c.TokenAuth.Encode(map[string]interface{}{"user_name": name})
	return tokenString
}

func GetUserNameFromContext(ctx context.Context) string {
	_, c, _ := jwtauth.FromContext(ctx)
	return c["user_name"].(string)
}
