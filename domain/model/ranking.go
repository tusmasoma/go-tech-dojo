package model

const (
	ScoreBoardKey   = "score_board"
	MaxRankingCount = 10
)

type Ranking struct {
	UserName string `json:"user_name"`
	Rank     int    `json:"rank"`
	Score    int    `json:"score"`
}

func NewRanking(userName string, rank, score int) *Ranking {
	return &Ranking{
		UserName: userName,
		Rank:     rank,
		Score:    score,
	}
}
