package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"github.com/tusmasoma/go-tech-dojo/config"
	"github.com/tusmasoma/go-tech-dojo/domain/model"
	"github.com/tusmasoma/go-tech-dojo/domain/repository/mock"
)

func TestUserUseCase_GetUser(t *testing.T) {
	t.Parallel()

	userID := uuid.New().String()
	ctx := context.WithValue(context.Background(), config.ContextUserIDKey, userID)

	user := model.User{
		ID:        userID,
		Name:      "test",
		Email:     "test@gmail.com",
		Coins:     100,
		HighScore: 1000,
	}

	patterns := []struct {
		name  string
		ctx   context.Context
		setup func(
			m *mock.MockUserRepository,
			m1 *mock.MockTransactionRepository,
		)
		wantErr error
	}{
		{
			name: "success",
			ctx:  ctx,
			setup: func(m *mock.MockUserRepository, m1 *mock.MockTransactionRepository) {
				m.EXPECT().Get(
					ctx,
					userID,
				).Return(&user, nil)
			},
			wantErr: nil,
		},
		{
			name:    "Fail: User ID not found in request context",
			ctx:     context.Background(),
			wantErr: fmt.Errorf("user name not found in request context"),
		},
	}
	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			ur := mock.NewMockUserRepository(ctrl)
			tr := mock.NewMockTransactionRepository(ctrl)
			ucr := mock.NewMockUserCollectionRepository(ctrl)
			cr := mock.NewMockCollectionRepository(ctrl)
			ccr := mock.NewMockCollectionCacheRepository(ctrl)

			if tt.setup != nil {
				tt.setup(ur, tr)
			}

			usecase := NewUserUseCase(ur, tr, ucr, cr, ccr)
			_, err := usecase.GetUser(tt.ctx)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type CreateUserAndTokenArg struct {
	ctx      context.Context
	email    string
	passward string
}

func TestUserUseCase_CreateUserAndToken(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockUserRepository,
			m1 *mock.MockTransactionRepository,
		)
		arg     CreateUserAndTokenArg
		wantErr error
	}{
		{
			name: "success",
			setup: func(m *mock.MockUserRepository, m1 *mock.MockTransactionRepository) {
				m1.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				})
				m.EXPECT().LockUserByEmail(
					gomock.Any(),
					"test@gmail.com",
				).Return(false, nil)
				m.EXPECT().Create(
					gomock.Any(),
					gomock.Any(),
				).Return(nil)
			},
			arg: CreateUserAndTokenArg{
				ctx:      context.Background(),
				email:    "test@gmail.com",
				passward: "password123",
			},
			wantErr: nil,
		},
		{
			name: "Fail: Username already exists",
			setup: func(m *mock.MockUserRepository, m2 *mock.MockTransactionRepository) {
				m2.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				})
				m.EXPECT().LockUserByEmail(
					gomock.Any(),
					"test@gmail.com",
				).Return(true, nil)
			},
			arg: CreateUserAndTokenArg{
				ctx:      context.Background(),
				email:    "test@gmail.com",
				passward: "password123",
			},
			wantErr: fmt.Errorf("user with this email already exists"),
		},
	}
	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			ur := mock.NewMockUserRepository(ctrl)
			tr := mock.NewMockTransactionRepository(ctrl)
			ucr := mock.NewMockUserCollectionRepository(ctrl)
			cr := mock.NewMockCollectionRepository(ctrl)
			ccr := mock.NewMockCollectionCacheRepository(ctrl)

			if tt.setup != nil {
				tt.setup(ur, tr)
			}

			usecase := NewUserUseCase(ur, tr, ucr, cr, ccr)
			jwt, err := usecase.CreateUserAndToken(tt.arg.ctx, tt.arg.email, tt.arg.passward)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("CreateUserAndToken() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("CreateUserAndToken() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr == nil && jwt == "" {
				t.Error("Failed to generate token")
			}
		})
	}
}

type UpdateUserArg struct {
	ctx       context.Context
	coins     int
	highscore int
}

func TestUserUseCase_UpdateUser(t *testing.T) {
	t.Parallel()

	userID := uuid.New().String()
	ctx := context.WithValue(context.Background(), config.ContextUserIDKey, userID)

	user := model.User{
		ID:        userID,
		Name:      "test",
		Email:     "test@gmail.com",
		Coins:     100,
		HighScore: 1000,
	}

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockUserRepository,
			m1 *mock.MockTransactionRepository,
		)
		arg     UpdateUserArg
		wantErr error
	}{
		{
			name: "success",
			setup: func(m *mock.MockUserRepository, m1 *mock.MockTransactionRepository) {
				m.EXPECT().Get(
					ctx,
					userID,
				).Return(&user, nil)
				user.Coins = 120
				user.HighScore = 1100
				m.EXPECT().Update(
					gomock.Any(),
					user,
				).Return(nil)
			},
			arg: UpdateUserArg{
				ctx:       ctx,
				coins:     120,
				highscore: 1100,
			},
			wantErr: nil,
		},
		{
			name: "Fail: User ID not found in request context",
			arg: UpdateUserArg{
				ctx:       context.Background(),
				coins:     100,
				highscore: 1000,
			},
			wantErr: fmt.Errorf("user name not found in request context"),
		},
	}
	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			ur := mock.NewMockUserRepository(ctrl)
			tr := mock.NewMockTransactionRepository(ctrl)
			ucr := mock.NewMockUserCollectionRepository(ctrl)
			cr := mock.NewMockCollectionRepository(ctrl)
			ccr := mock.NewMockCollectionCacheRepository(ctrl)

			if tt.setup != nil {
				tt.setup(ur, tr)
			}

			usecase := NewUserUseCase(ur, tr, ucr, cr, ccr)
			updateUser, err := usecase.UpdateUser(tt.arg.ctx, tt.arg.coins, tt.arg.highscore)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr == nil && updateUser == nil {
				t.Error("Failed to update user")
			}
		})
	}
}

func TestUserUseCase_ListUserCollections(t *testing.T) {
	t.Parallel()

	userID := uuid.New().String()
	collection1ID := uuid.New().String()
	collection2ID := uuid.New().String()
	collection3ID := uuid.New().String()

	ctx := context.WithValue(context.Background(), config.ContextUserIDKey, userID)

	collections := model.Collections{
		{
			ID:     collection1ID,
			Name:   "collection1",
			Rarity: 3,
			Weight: 3,
		},
		{
			ID:     collection2ID,
			Name:   "collection2",
			Rarity: 2,
			Weight: 2,
		},
		{
			ID:     collection3ID,
			Name:   "collection3",
			Rarity: 1,
			Weight: 1,
		},
	}

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockCollectionRepository,
			m1 *mock.MockCollectionCacheRepository,
			m2 *mock.MockUserCollectionRepository,
		)
		want struct {
			collections []*Collection
			err         error
		}
	}{
		{
			name: "Success: get collections in cache",
			setup: func(
				m *mock.MockCollectionRepository,
				m1 *mock.MockCollectionCacheRepository,
				m2 *mock.MockUserCollectionRepository,
			) {
				m1.EXPECT().Get(
					ctx,
					"collections",
				).Return(
					collections,
					nil,
				)
				m2.EXPECT().List(ctx, userID).Return(
					[]*model.UserCollection{
						{
							UserID:       userID,
							CollectionID: collection1ID,
						},
						{
							UserID:       userID,
							CollectionID: collection2ID,
						},
					},
					nil,
				)
			},
			want: struct {
				collections []*Collection
				err         error
			}{
				collections: []*Collection{
					{
						Collection: &model.Collection{
							ID:     collection1ID,
							Name:   "collection1",
							Rarity: 3,
							Weight: 3,
						},
						Has: true,
					},
					{
						Collection: &model.Collection{
							ID:     collection2ID,
							Name:   "collection2",
							Rarity: 2,
							Weight: 2,
						},
						Has: true,
					},
					{
						Collection: &model.Collection{
							ID:     collection3ID,
							Name:   "collection3",
							Rarity: 1,
							Weight: 1,
						},
						Has: false,
					},
				},
				err: nil,
			},
		},
		{
			name: "Success: get collections in DB",
			setup: func(
				m *mock.MockCollectionRepository,
				m1 *mock.MockCollectionCacheRepository,
				m2 *mock.MockUserCollectionRepository,
			) {
				m1.EXPECT().Get(
					ctx,
					"collections",
				).Return(
					nil,
					config.ErrCacheMiss,
				)
				m.EXPECT().List(ctx).Return(
					collections,
					nil,
				)
				m1.EXPECT().Create(
					ctx,
					"collections",
					collections,
				).Return(nil)
				m2.EXPECT().List(ctx, userID).Return(
					[]*model.UserCollection{
						{
							UserID:       userID,
							CollectionID: collection1ID,
						},
						{
							UserID:       userID,
							CollectionID: collection2ID,
						},
					},
					nil,
				)
			},
			want: struct {
				collections []*Collection
				err         error
			}{
				collections: []*Collection{
					{
						Collection: &model.Collection{
							ID:     collection1ID,
							Name:   "collection1",
							Rarity: 3,
							Weight: 3,
						},
						Has: true,
					},
					{
						Collection: &model.Collection{
							ID:     collection2ID,
							Name:   "collection2",
							Rarity: 2,
							Weight: 2,
						},
						Has: true,
					},
					{
						Collection: &model.Collection{
							ID:     collection3ID,
							Name:   "collection3",
							Rarity: 1,
							Weight: 1,
						},
						Has: false,
					},
				},
				err: nil,
			},
		},
	}
	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			ur := mock.NewMockUserRepository(ctrl)
			tr := mock.NewMockTransactionRepository(ctrl)
			ucr := mock.NewMockUserCollectionRepository(ctrl)
			cr := mock.NewMockCollectionRepository(ctrl)
			ccr := mock.NewMockCollectionCacheRepository(ctrl)

			if tt.setup != nil {
				tt.setup(cr, ccr, ucr)
			}

			usecase := NewUserUseCase(ur, tr, ucr, cr, ccr)
			collections, err := usecase.ListUserCollections(ctx)

			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("ListUserCollections() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("ListUserCollections() error = %v, wantErr %v", err, tt.want.err)
			}

			if tt.want.err == nil && collections == nil {
				t.Error("Failed to list user collections")
			}
		})
	}
}
