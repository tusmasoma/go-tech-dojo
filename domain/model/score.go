package model

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

type Score struct {
	ID     string `db:"id" json:"id"`
	UserID string `db:"user_id" json:"user_id"`
	Value  int    `db:"value" json:"value"`
}

func NewScore(userID string, value int) (*Score, error) {
	if userID == "" {
		log.Error("userID is required")
		return nil, fmt.Errorf("userID is required")
	}
	if value < 0 {
		log.Error("value is less than 0")
		return nil, fmt.Errorf("value is less than 0")
	}
	return &Score{
		ID:     uuid.New().String(),
		UserID: userID,
		Value:  value,
	}, nil
}
