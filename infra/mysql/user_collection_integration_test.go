package mysql

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/uuid"

	"github.com/tusmasoma/go-tech-dojo/domain/model"
)

func Test_UserCollectionRepository(t *testing.T) {
	ctx := context.Background()
	repo := NewUserCollectionRepository(db)

	userID := uuid.New().String()
	collection1ID := uuid.New().String()
	collection2ID := uuid.New().String()
	collection3ID := uuid.New().String()

	userCollection1, _ := model.NewUserCollection(
		userID,
		collection1ID,
	)
	userCollection2, _ := model.NewUserCollection(
		userID,
		collection2ID,
	)
	userCollection3, _ := model.NewUserCollection(
		userID,
		collection3ID,
	)

	// serup
	userRepo := NewUserRepository(db)
	user := model.User{
		ID:       userID,
		Name:     "test",
		Email:    "test@gmail.com",
		Password: "password",
	}
	err := userRepo.Create(ctx, user)
	ValidateErr(t, err, nil)

	collectionRepo := NewCollectionRepository(db)
	collection1 := model.Collection{
		ID:     collection1ID,
		Name:   "collection1",
		Rarity: 1,
		Weight: 1,
	}
	collection2 := model.Collection{
		ID:     collection2ID,
		Name:   "collection2",
		Rarity: 2,
		Weight: 2,
	}
	collection3 := model.Collection{
		ID:     collection3ID,
		Name:   "collection3",
		Rarity: 3,
		Weight: 3,
	}

	err = collectionRepo.BatchCreate(ctx, []*model.Collection{&collection1, &collection2, &collection3})
	ValidateErr(t, err, nil)

	// Create
	err = repo.Create(ctx, *userCollection1)
	ValidateErr(t, err, nil)

	// BatchCreate
	err = repo.BatchCreate(ctx, []*model.UserCollection{userCollection2, userCollection3})
	ValidateErr(t, err, nil)

	// Get
	getUserCollection, err := repo.Get(ctx, userID, collection1ID)
	ValidateErr(t, err, nil)
	if !reflect.DeepEqual(userCollection1, getUserCollection) {
		t.Errorf("want: %v, got: %v", userCollection1, getUserCollection)
	}

	// List
	listUserCollections, err := repo.List(ctx, userID)
	ValidateErr(t, err, nil)
	if len(listUserCollections) != 3 {
		t.Errorf("want: %v, got: %v", 3, len(listUserCollections))
	}

	// Delete
	err = repo.Delete(ctx, userID, collection1ID)
	ValidateErr(t, err, nil)

	_, err = repo.Get(ctx, userID, collection1ID)
	if err == nil {
		t.Errorf("want: %v, got: %v", nil, err)
	}
}
