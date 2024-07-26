package redis

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/uuid"

	"github.com/tusmasoma/go-tech-dojo/domain/model"
)

func Test_CollectionRepository(t *testing.T) {
	ctx := context.Background()
	repo := NewCollectionRepository(client)

	collections := model.Collections{
		{
			ID:     uuid.New().String(),
			Name:   "collection1",
			Rarity: 1,
			Weight: 1,
		},
		{
			ID:     uuid.New().String(),
			Name:   "collection2",
			Rarity: 2,
			Weight: 2,
		},
	}

	// Create
	err := repo.Create(ctx, "collections", collections)
	ValidateErr(t, err, nil)

	// Get
	getCollections, err := repo.Get(ctx, "collections")
	ValidateErr(t, err, nil)
	if !reflect.DeepEqual(collections, getCollections) {
		t.Errorf("want: %v, got: %v", collections, getCollections)
	}

	// Delete
	err = repo.Delete(ctx, "collections")
	ValidateErr(t, err, nil)
	_, err = repo.Get(ctx, "collections")
	if err == nil {
		t.Errorf("want: %v, got: %v", nil, err)
	}
}
