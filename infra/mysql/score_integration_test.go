package mysql

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/uuid"

	"github.com/tusmasoma/go-tech-dojo/domain/model"
)

func Test_ScoreRepository(t *testing.T) {
	ctx := context.Background()
	repo := NewScoreRepository(db)

	userID := uuid.New().String()
	score1, _ := model.NewScore(userID, 100)

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

	// Create
	err = repo.Create(ctx, *score1)
	ValidateErr(t, err, nil)

	// Get
	getScore, err := repo.Get(ctx, score1.ID)
	ValidateErr(t, err, nil)
	if !reflect.DeepEqual(score1, getScore) {
		t.Errorf("want: %v, got: %v", score1, getScore)
	}

	// Delete
	err = repo.Delete(ctx, score1.ID)
	ValidateErr(t, err, nil)

	getScore, err = repo.Get(ctx, score1.ID)
	if err == nil {
		t.Errorf("want: %v, got: %v", "record not found", getScore)
	}
}
