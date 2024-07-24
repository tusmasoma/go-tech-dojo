//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package usecase

import (
	"context"
	"fmt"

	"github.com/tusmasoma/go-tech-dojo/domain/model"
	"github.com/tusmasoma/go-tech-dojo/domain/repository"
	"github.com/tusmasoma/go-tech-dojo/pkg/auth"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

type UserUseCase interface {
	CreateUserAndToken(ctx context.Context, email string, passward string) (string, error)
}

type userUseCase struct {
	ur repository.UserRepository
	tr repository.TransactionRepository
}

func NewUserUseCase(ur repository.UserRepository, tr repository.TransactionRepository) UserUseCase {
	return &userUseCase{
		ur: ur,
		tr: tr,
	}
}

func (uuc *userUseCase) CreateUserAndToken(ctx context.Context, email string, password string) (string, error) {
	var user *model.User
	err := uuc.tr.Transaction(ctx, func(ctx context.Context) error {
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
	})
	if err != nil {
		return "", err
	}

	jwt, _ := auth.GenerateToken(user.ID, user.Email)
	return jwt, nil
}
