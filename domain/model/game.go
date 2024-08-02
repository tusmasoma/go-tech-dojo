package model

import (
	"errors"
	"math/rand"
)

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

type Gacha struct{}

func NewGacha() *Gacha {
	return &Gacha{}
}

func (g *Gacha) Draw(collections Collections) (*Collection, error) {
	// シャッフル 重み付け
	target := rand.Intn(collections.TotalWeight()) //nolint:gosec // Use math/rand
	for _, item := range collections {
		target -= item.Weight
		if target < 0 {
			return item, nil
		}
	}
	return nil, errors.New("failed to pick an item")
}
