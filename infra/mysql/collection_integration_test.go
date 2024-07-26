package mysql

import (
	"context"
	"reflect"
	"testing"

	"github.com/tusmasoma/go-tech-dojo/domain/model"
)

func Test_CollectionRepository(t *testing.T) {
	ctx := context.Background()
	repo := NewCollectionRepository(db)

	collection1, _ := model.NewCollection(
		"collection1",
		1,
		1,
	)
	collection2, _ := model.NewCollection(
		"collection2",
		2,
		2,
	)
	collection3, _ := model.NewCollection(
		"collection3",
		3,
		3,
	)

	collections := model.Collections{
		collection2,
		collection3,
	}

	// Create
	err := repo.Create(ctx, *collection1)
	ValidateErr(t, err, nil)

	// BatchCreate
	err = repo.BatchCreate(ctx, collections)
	ValidateErr(t, err, nil)

	// Get
	getCollection, err := repo.Get(ctx, collection1.ID)
	ValidateErr(t, err, nil)
	if !reflect.DeepEqual(collection1, getCollection) {
		t.Errorf("want: %v, got: %v", collection1, getCollection)
	}

	// List
	listCollections, err := repo.List(ctx)
	ValidateErr(t, err, nil)
	if len(listCollections) != 3 {
		t.Errorf("want: %v, got: %v", 3, len(listCollections))
	}

	// Update
	getCollection.Name = "updatedName"
	err = repo.Update(ctx, *getCollection)
	ValidateErr(t, err, nil)

	updatedCollection, err := repo.Get(ctx, collection1.ID)
	ValidateErr(t, err, nil)
	if !reflect.DeepEqual(getCollection, updatedCollection) {
		t.Errorf("want: %v, got: %v", getCollection, updatedCollection)
	}

	// Delete
	err = repo.Delete(ctx, collection1.ID)
	ValidateErr(t, err, nil)

	_, err = repo.Get(ctx, collection1.ID)
	if err == nil {
		t.Errorf("want: %v, got: %v", nil, err)
	}
}
