package model

import (
	"fmt"

	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

const (
	ScoreBoardKey   = "score_board"
	MaxRankingCount = 10
)

type Ranking struct {
	UserName string `json:"user_name"`
	Rank     int    `json:"rank"`
	Score    int    `json:"score"`
}

func NewRanking(userName string, rank, score int) (*Ranking, error) {
	if userName == "" {
		log.Error("userName is empty")
		return nil, fmt.Errorf("userName is empty")
	}
	if rank < 1 {
		log.Error("rank is less than 1")
		return nil, fmt.Errorf("rank is less than 1")
	}
	if score < 0 {
		log.Error("score is less than 0")
		return nil, fmt.Errorf("score is less than 0")
	}
	return &Ranking{
		UserName: userName,
		Rank:     rank,
		Score:    score,
	}, nil
}
