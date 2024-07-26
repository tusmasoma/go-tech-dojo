//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package repository

import (
	"context"

	"github.com/tusmasoma/go-tech-dojo/domain/model"
)

type CollectionRepository interface {
	Get(ctx context.Context, id string) (*model.Collection, error)
	List(ctx context.Context) (model.Collections, error)
	Create(ctx context.Context, collection model.Collection) error
	BatchCreate(ctx context.Context, collections model.Collections) error
	Update(ctx context.Context, collection model.Collection) error
	Delete(ctx context.Context, id string) error
}

type CollectionCacheRepository interface {
	Get(ctx context.Context, key string) (model.Collections, error)
	Create(ctx context.Context, key string, collection model.Collections) error
	Delete(ctx context.Context, key string) error
}
