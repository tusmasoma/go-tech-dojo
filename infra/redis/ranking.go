package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"

	"github.com/tusmasoma/go-tech-dojo/domain/model"
	"github.com/tusmasoma/go-tech-dojo/domain/repository"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

type rankingRepository struct {
	client *redis.Client
}

func NewRankingRepository(client *redis.Client) repository.RankingRepository {
	return &rankingRepository{
		client: client,
	}
}

func (rr *rankingRepository) List(ctx context.Context, key string, start int) ([]*model.Ranking, error) {
	total, err := rr.client.ZCard(ctx, key).Result()
	if err != nil {
		log.Error("Failed to get ranking", log.Ferror(err))
		return nil, err
	}
	if int64(start) > total || start < 1 {
		log.Warn("Ranking is empty", log.Fint("start", start), log.Fint("count", int(total)))
		return nil, fmt.Errorf("ranking is empty")
	}

	results, err := rr.client.ZRevRangeWithScores(
		ctx,
		key,
		int64(start-1),
		int64(start+model.MaxRankingCount-1),
	).Result()
	if err != nil {
		log.Error("Failed to get ranking", log.Ferror(err))
		return nil, err
	}

	rankings := make([]*model.Ranking, 0, len(results))
	for i, result := range results {
		rank := start + i
		rankings = append(rankings, &model.Ranking{
			UserName: result.Member.(string),
			Rank:     rank,
			Score:    int(result.Score),
		})
	}
	return rankings, nil
}

func (rr *rankingRepository) Create(ctx context.Context, key string, ranking *model.Ranking) error {
	_, err := rr.client.ZAdd(ctx, key, &redis.Z{
		Score:  float64(ranking.Score),
		Member: ranking.UserName,
	}).Result()
	if err != nil {
		log.Error("Failed to set ranking", log.Ferror(err))
		return err
	}
	log.Info("Ranking set successfully", log.Fstring("user_name", ranking.UserName))
	return nil
}
