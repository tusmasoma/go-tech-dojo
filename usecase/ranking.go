//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package usecase

import (
	"context"

	"github.com/tusmasoma/go-tech-dojo/domain/model"
	"github.com/tusmasoma/go-tech-dojo/domain/repository"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

type RankingUseCase interface {
	ListRankings(ctx context.Context, start int) ([]*model.Ranking, error)
}

type rankingUseCase struct {
	rr repository.RankingRepository
}

func NewRankingUseCase(rr repository.RankingRepository) RankingUseCase {
	return &rankingUseCase{
		rr: rr,
	}
}

func (ruc *rankingUseCase) ListRankings(ctx context.Context, start int) ([]*model.Ranking, error) {
	rankings, err := ruc.rr.List(ctx, model.ScoreBoardKey, start)
	if err != nil {
		log.Error("Failed to list rankings", log.Ferror(err))
		return nil, err
	}
	return rankings, nil
}
