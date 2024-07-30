package model

const (
	BaseReward      = 100 // ゲーム基本報酬コイン
	ScoreMultiplier = 2   // ゲームスコア倍率
)

type Game struct{}

func NewGame() *Game {
	return &Game{}
}

func (g *Game) Reward(score int) int {
	if score > 0 {
		return BaseReward + score*ScoreMultiplier
	}
	return 0
}
