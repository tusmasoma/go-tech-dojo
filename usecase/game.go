//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package usecase

import (
	"context"
	"fmt"

	"github.com/tusmasoma/go-tech-dojo/config"
	"github.com/tusmasoma/go-tech-dojo/domain/model"
	"github.com/tusmasoma/go-tech-dojo/domain/repository"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

type GameUseCase interface {
	FinishGame(ctx context.Context, score int) (int, error)
}

type gameUseCase struct {
	tr repository.TransactionRepository
	ur repository.UserRepository
	sr repository.ScoreRepository
	rr repository.RankingRepository
}

func NewGameUseCase(
	tr repository.TransactionRepository,
	ur repository.UserRepository,
	sr repository.ScoreRepository,
	rr repository.RankingRepository,
) GameUseCase {
	return &gameUseCase{
		tr: tr,
		ur: ur,
		sr: sr,
		rr: rr,
	}
}

func (guc *gameUseCase) FinishGame(ctx context.Context, score int) (int, error) {
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
	err = guc.tr.Transaction(ctx, func(ctx context.Context) error {
		var game model.Game
		score, err := model.NewScore(user.ID, score) //nolint:govet // This is a valid code
		if err != nil {
			log.Error("Failed to create score", log.Ferror(err))
			return err
		}
		if err = guc.sr.Create(ctx, *score); err != nil {
			log.Error("Failed to create score", log.Ferror(err))
			return err
		}

		if user.HighScore < score.Value {
			user.HighScore = score.Value
		}
		coin = game.Reward(score.Value)
		user.Coins += coin
		if err = guc.ur.Update(ctx, *user); err != nil {
			log.Error("Failed to update user", log.Ferror(err))
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	if err = guc.rr.Create(
		ctx,
		model.ScoreBoardKey,
		&model.Ranking{
			UserName: user.Name,
			Score:    score,
		},
	); err != nil {
		log.Error("Failed to create ranking", log.Ferror(err))
		return 0, err
	}
	return coin, nil
}
