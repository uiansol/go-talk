package main

import (
	"context"

	"github.com/go-chi/jwtauth"
	"golang.org/x/crypto/bcrypt"
)

// Generates the bcrypt hash of the given password.
func CreatePassword(password []byte) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return hashedPassword, err
}

// Check if a given password is correct.
func CheckPassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}

// Generates a JSON Web Token (JWT) containing a user's name.
func (c *Config) MakeToken(name string) string {
	_, tokenString, _ := c.TokenAuth.Encode(map[string]interface{}{"user_name": name})
	return tokenString
}

// Retrieves a user's name from a JWT.
func GetUserNameFromContext(ctx context.Context) string {
	_, c, _ := jwtauth.FromContext(ctx)
	return c["user_name"].(string)
}
