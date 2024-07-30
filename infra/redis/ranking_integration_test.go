package redis

import (
	"context"
	"testing"

	"github.com/tusmasoma/go-tech-dojo/domain/model"
)

func Test_RankingRepository(t *testing.T) {
	ctx := context.Background()
	repo := NewRankingRepository(client)

	ranking1 := model.Ranking{
		UserName: "user1",
		Score:    100,
	}
	ranking2 := model.Ranking{
		UserName: "user2",
		Score:    200,
	}

	// Create
	err := repo.Create(ctx, "ranking", &ranking1)
	ValidateErr(t, err, nil)
	err = repo.Create(ctx, "ranking", &ranking2)
	ValidateErr(t, err, nil)

	// List
	rankings, err := repo.List(ctx, "ranking", 1)
	ValidateErr(t, err, nil)
	if len(rankings) != 2 {
		t.Errorf("want: %d, got: %d", 2, len(rankings))
	}
	if rankings[0].UserName != "user2" || rankings[1].UserName != "user1" {
		t.Errorf("want: %v, got: %v", []string{"user2", "user1"}, []string{rankings[0].UserName, rankings[1].UserName})
	}
	if rankings[0].Rank != 1 || rankings[1].Rank != 2 {
		t.Errorf("want: %v, got: %v", []int{1, 2}, []int{rankings[0].Rank, rankings[1].Rank})
	}
}
