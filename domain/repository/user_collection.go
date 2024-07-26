//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package repository

import (
	"context"

	"github.com/tusmasoma/go-tech-dojo/domain/model"
)

type UserCollectionRepository interface {
	List(ctx context.Context, userID string) ([]*model.UserCollection, error)
	Get(ctx context.Context, userID, collectionID string) (*model.UserCollection, error)
	Create(ctx context.Context, userCollection model.UserCollection) error
	BatchCreate(ctx context.Context, userCollections []*model.UserCollection) error
	Delete(ctx context.Context, userID, collectionID string) error
}
