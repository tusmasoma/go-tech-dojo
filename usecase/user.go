//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/tusmasoma/go-tech-dojo/config"
	"github.com/tusmasoma/go-tech-dojo/domain/model"
	"github.com/tusmasoma/go-tech-dojo/domain/repository"
	"github.com/tusmasoma/go-tech-dojo/pkg/auth"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

type UserUseCase interface {
	GetUser(ctx context.Context) (*model.User, error)
	ListUserCollections(ctx context.Context) ([]*Collection, error)
	CreateUserAndToken(ctx context.Context, email string, passward string) (string, error)
	UpdateUser(ctx context.Context, coins, highscore int) (*model.User, error)
}

type userUseCase struct {
	ur  repository.UserRepository
	tr  repository.TransactionRepository
	ucr repository.UserCollectionRepository
	cr  repository.CollectionRepository
	ccr repository.CollectionCacheRepository
}

func NewUserUseCase(
	ur repository.UserRepository,
	tr repository.TransactionRepository,
	ucr repository.UserCollectionRepository,
	cr repository.CollectionRepository,
	ccr repository.CollectionCacheRepository,
) UserUseCase {
	return &userUseCase{
		ur:  ur,
		tr:  tr,
		ucr: ucr,
		cr:  cr,
		ccr: ccr,
	}
}

func (uuc *userUseCase) GetUser(ctx context.Context) (*model.User, error) {
	userIDValue := ctx.Value(config.ContextUserIDKey)
	userID, ok := userIDValue.(string)
	if !ok {
		log.Error("User ID not found in request context")
		return nil, fmt.Errorf("user name not found in request context")
	}
	user, err := uuc.ur.Get(ctx, userID)
	if err != nil {
		log.Error("Error getting user", log.Fstring("user_id", userID))
		return nil, err
	}
	return user, nil
}

func (uuc *userUseCase) CreateUserAndToken(ctx context.Context, email string, password string) (string, error) {
	var user *model.User
	if err := uuc.tr.Transaction(ctx, func(ctx context.Context) error {
		exists, err := uuc.ur.LockUserByEmail(ctx, email)
		if err != nil {
			log.Error("Error retrieving user by email", log.Fstring("email", email))
			return err
		}
		if exists {
			log.Info("User with this email already exists", log.Fstring("email", email))
			return fmt.Errorf("user with this email already exists")
		}

		user, err = model.NewUser(email, password)
		if err != nil {
			log.Error("Error creating new user", log.Fstring("email", email))
			return err
		}

		if err = uuc.ur.Create(ctx, *user); err != nil {
			log.Error("Error creating new user", log.Fstring("email", email))
			return err
		}

		return nil
	}); err != nil {
		return "", err
	}

	jwt, _ := auth.GenerateToken(user.ID, user.Email)
	return jwt, nil
}

func (uuc *userUseCase) UpdateUser(ctx context.Context, coins, highscore int) (*model.User, error) {
	userIDValue := ctx.Value(config.ContextUserIDKey)
	userID, ok := userIDValue.(string)
	if !ok {
		log.Error("User ID not found in request context")
		return nil, fmt.Errorf("user name not found in request context")
	}
	user, err := uuc.ur.Get(ctx, userID)
	if err != nil {
		log.Error("Error getting user", log.Fstring("user_id", userID))
		return nil, err
	}

	// TODO: setter method for user
	user.Coins = coins
	user.HighScore = highscore

	if err = uuc.ur.Update(ctx, *user); err != nil {
		log.Error("Error updating user", log.Fstring("user_id", userID))
		return nil, err
	}
	return user, nil
}

type Collection struct {
	*model.Collection
	Has bool
}

func (uuc *userUseCase) ListUserCollections(ctx context.Context) ([]*Collection, error) {
	userIDValue := ctx.Value(config.ContextUserIDKey)
	userID, ok := userIDValue.(string)
	if !ok {
		log.Error("User ID not found in request context")
		return nil, fmt.Errorf("user name not found in request context")
	}

	collections, err := uuc.ccr.Get(ctx, "collections")
	if errors.Is(err, config.ErrCacheMiss) {
		log.Info("Cache miss", log.Fstring("key", "collections"))
		collections, err = uuc.cr.List(ctx)
		if err != nil {
			log.Error("Error getting collections", log.Ferror(err))
			return nil, err
		}

		if err = uuc.ccr.Create(ctx, "collections", collections); err != nil {
			log.Error("Error setting collections to cache", log.Ferror(err))
			return nil, err
		}
	} else if err != nil {
		log.Error("Error getting collections from cache", log.Ferror(err))
		return nil, err
	}

	userCollections, err := uuc.ucr.List(ctx, userID)
	if err != nil {
		log.Error("Error getting user collections", log.Fstring("user_id", userID))
		return nil, err
	}

	userCollectionMap := make(map[string]bool)
	for _, uc := range userCollections {
		userCollectionMap[uc.CollectionID] = true
	}

	var reusult []*Collection
	for _, c := range collections {
		has := userCollectionMap[c.ID]
		reusult = append(reusult, &Collection{
			Collection: c,
			Has:        has,
		})
	}

	return reusult, nil
}
