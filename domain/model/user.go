package model

import (
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Coins     int    `json:"coins"`
	HighScore int    `json:"highscore"`
}

func NewUser(email, password string) (*User, error) {
	if email == "" || password == "" {
		log.Error("Email or Password is empty", log.Fstring("email", email))
		return nil, fmt.Errorf("email or password is empty")
	}
	name := extractNameFromEmail(email)
	return &User{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		Password:  password,
		Coins:     0,
		HighScore: 0,
	}, nil
}

func extractNameFromEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) > 0 {
		return parts[0]
	}
	return "unknown"
}
