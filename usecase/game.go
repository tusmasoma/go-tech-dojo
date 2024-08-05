//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/tusmasoma/go-tech-dojo/config"
	"github.com/tusmasoma/go-tech-dojo/domain/model"
	"github.com/tusmasoma/go-tech-dojo/domain/repository"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

type GameUseCase interface {
	FinishGame(ctx context.Context, scoreValue int) (int, error)
	DrawGacha(ctx context.Context, times int) ([]*GachaResult, error)
}

type gameUseCase struct {
	tr  repository.TransactionRepository
	ur  repository.UserRepository
	ucr repository.UserCollectionRepository
	sr  repository.ScoreRepository
	rr  repository.RankingRepository
	cr  repository.CollectionRepository
	ccr repository.CollectionCacheRepository
}

func NewGameUseCase(
	tr repository.TransactionRepository,
	ur repository.UserRepository,
	ucr repository.UserCollectionRepository,
	sr repository.ScoreRepository,
	rr repository.RankingRepository,
	cr repository.CollectionRepository,
	ccr repository.CollectionCacheRepository,
) GameUseCase {
	return &gameUseCase{
		tr:  tr,
		ur:  ur,
		ucr: ucr,
		sr:  sr,
		rr:  rr,
		cr:  cr,
		ccr: ccr,
	}
}

func (guc *gameUseCase) FinishGame(ctx context.Context, scoreValue int) (int, error) {
	userIDValue := ctx.Value(config.ContextUserIDKey)
	userID, ok := userIDValue.(string)
	if !ok {
		log.Error("User ID not found in request context")
		return 0, fmt.Errorf("user name not found in request context")
	}
	user, err := guc.ur.Get(ctx, userID)
	if err != nil {
		log.Error("Error getting user", log.Fstring("user_id", userID))
		return 0, err
	}

	var coin int
	if err = guc.tr.Transaction(ctx, func(ctx context.Context) error {
		var game model.Game
		score, err := model.NewScore(user.ID, scoreValue) //nolint:govet // This is a valid code
		if err != nil {
			log.Error("Failed to create score", log.Ferror(err))
			return err
		}
		if err = guc.sr.Create(ctx, *score); err != nil {
			log.Error("Failed to create score", log.Ferror(err))
			return err
		}

		if user.HighScore < scoreValue {
			user.HighScore = scoreValue
		}
		coin = game.Reward(scoreValue)
		user.Coins += coin
		if err = guc.ur.Update(ctx, *user); err != nil {
			log.Error("Failed to update user", log.Ferror(err))
			return err
		}
		return nil
	}); err != nil {
		return 0, err
	}

	if err = guc.rr.Create(
		ctx,
		model.ScoreBoardKey,
		&model.Ranking{
			UserName: user.Name,
			Score:    scoreValue,
		},
	); err != nil {
		log.Error("Failed to create ranking", log.Ferror(err))
		return 0, err
	}
	return coin, nil
}

type GachaResult struct {
	*model.Collection
	Has bool `json:"has"`
}

func (guc *gameUseCase) DrawGacha(ctx context.Context, times int) ([]*GachaResult, error) { //nolint:gocognit // This is a valid code
	userIDValue := ctx.Value(config.ContextUserIDKey)
	userID, ok := userIDValue.(string)
	if !ok {
		log.Error("User ID not found in request context")
		return nil, fmt.Errorf("user name not found in request context")
	}
	user, err := guc.ur.Get(ctx, userID)
	if err != nil {
		log.Error("Error getting user", log.Fstring("user_id", userID))
		return nil, err
	}

	collections, err := guc.ccr.Get(ctx, "collections")
	if errors.Is(err, config.ErrCacheMiss) {
		log.Info("Cache miss", log.Fstring("key", "collections"))
		collections, err = guc.cr.List(ctx)
		if err != nil {
			log.Error("Error getting collections", log.Ferror(err))
			return nil, err
		}
		if err = guc.ccr.Create(ctx, "collections", collections); err != nil {
			log.Error("Error setting collections to cache", log.Ferror(err))
			return nil, err
		}
	} else if err != nil {
		log.Error("Error getting collections from cache", log.Ferror(err))
		return nil, err
	}

	var gacha model.Gacha
	results := make(model.Collections, 0, times)
	for i := 0; i < times; i++ {
		result, err := gacha.Draw(collections) //nolint:govet // This is a valid code
		if err != nil {
			log.Error("Failed to draw gacha", log.Ferror(err))
			return nil, err
		}
		results = append(results, result)
	}

	userCollections, err := guc.ucr.List(ctx, userID)
	if err != nil {
		log.Error("Error getting user collections", log.Fstring("user_id", userID))
		return nil, err
	}
	userCollectionMap := make(map[string]bool)
	for _, item := range userCollections {
		userCollectionMap[item.CollectionID] = true
	}
	var gachaResults []*GachaResult
	for _, item := range results {
		gachaResult := &GachaResult{
			Collection: item,
			Has:        userCollectionMap[item.ID],
		}
		gachaResults = append(gachaResults, gachaResult)
	}

	if err = guc.tr.Transaction(ctx, func(ctx context.Context) error {
		user.Coins -= gacha.Cost(times)
		if err = guc.ur.Update(ctx, *user); err != nil {
			log.Error("Failed to update user", log.Ferror(err))
			return err
		}

		var newUserCollections []*model.UserCollection
		for _, result := range results {
			newUserCollection, err := model.NewUserCollection(user.ID, result.ID) //nolint:govet // This is a valid code
			if err != nil {
				log.Error("Failed to create user collection", log.Ferror(err))
				return err
			}
			newUserCollections = append(newUserCollections, newUserCollection)
		}
		if err = guc.ucr.BatchCreate(ctx, newUserCollections); err != nil {
			log.Error("Failed to create user collections", log.Ferror(err))
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return gachaResults, nil
}
