package model

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

type Collection struct {
	ID     string `db:"id" json:"id"`
	Name   string `db:"name" json:"name"`
	Rarity int    `db:"rarity" json:"rarity"`
	Weight int    `db:"weight" json:"weight"`
}

func NewCollection(name string, rarity, weight int) (*Collection, error) {
	if name == "" {
		log.Error("Name is empty", log.Fstring("name", name))
		return nil, fmt.Errorf("name is empty")
	}
	if rarity < 0 || rarity > 5 {
		log.Error("Rarity is invalid", log.Fint("rarity", rarity))
		return nil, fmt.Errorf("rarity is invalid")
	}
	if weight < 0 {
		log.Error("Weight is invalid", log.Fint("weight", weight))
		return nil, fmt.Errorf("weight is invalid")
	}
	return &Collection{
		ID:     uuid.New().String(),
		Name:   name,
		Rarity: rarity,
		Weight: weight,
	}, nil
}

type Collections []*Collection
