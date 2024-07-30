//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package repository

import (
	"context"

	"github.com/tusmasoma/go-tech-dojo/domain/model"
)

type ScoreRepository interface {
	Get(ctx context.Context, id string) (*model.Score, error)
	Create(ctx context.Context, score model.Score) error
	Delete(ctx context.Context, id string) error
}
