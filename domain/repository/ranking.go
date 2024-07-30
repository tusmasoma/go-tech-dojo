//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package repository

import (
	"context"

	"github.com/tusmasoma/go-tech-dojo/domain/model"
)

type RankingRepository interface {
	List(ctx context.Context, key string, start int) ([]*model.Ranking, error)
	Create(ctx context.Context, key string, ranking *model.Ranking) error
}
