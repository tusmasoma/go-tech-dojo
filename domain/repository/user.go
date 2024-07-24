package repository

import (
	"context"

	"github.com/tusmasoma/go-tech-dojo/domain/model"
)

type UserRepository interface {
	Get(ctx context.Context, id string) (*model.User, error)
	Create(ctx context.Context, user model.User) error
	Update(ctx context.Context, user model.User) error
	Delete(ctx context.Context, id string) error
}
