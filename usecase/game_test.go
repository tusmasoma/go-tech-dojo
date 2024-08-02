package usecase

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"github.com/tusmasoma/go-tech-dojo/config"
	"github.com/tusmasoma/go-tech-dojo/domain/model"
	"github.com/tusmasoma/go-tech-dojo/domain/repository/mock"
)

func TestUserUseCase_FinishGame(t *testing.T) {
	t.Parallel()

	userID := uuid.New().String()
	ctx := context.WithValue(context.Background(), config.ContextUserIDKey, userID)

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockTransactionRepository,
			m1 *mock.MockUserRepository,
			m2 *mock.MockScoreRepository,
			m3 *mock.MockRankingRepository,
		)
		arg struct {
			ctx   context.Context
			score int
		}
		want struct {
			coin int
			err  error
		}
	}{
		{
			name: "Success: with high score",
			setup: func(tr *mock.MockTransactionRepository, ur *mock.MockUserRepository, sr *mock.MockScoreRepository, rr *mock.MockRankingRepository) {
				user := model.User{
					ID:        userID,
					Name:      "test",
					Email:     "test@gmail.com",
					Coins:     100,
					HighScore: 1000,
				}
				ur.EXPECT().Get(ctx, userID).Return(&user, nil)
				tr.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				})
				sr.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				ur.EXPECT().Update(
					gomock.Any(),
					model.User{
						ID:        user.ID,
						Name:      user.Name,
						Email:     user.Email,
						Coins:     2600,
						HighScore: 1200,
					},
				).Return(nil)
				rr.EXPECT().Create(
					gomock.Any(),
					model.ScoreBoardKey,
					&model.Ranking{
						UserName: user.Name,
						Score:    1200,
					},
				).Return(nil)
			},
			arg: struct {
				ctx   context.Context
				score int
			}{
				ctx:   ctx,
				score: 1200,
			},
			want: struct {
				coin int
				err  error
			}{
				coin: func() int {
					game := &model.Game{}
					return game.Reward(1200)
				}(),
				err: nil,
			},
		},
		{
			name: "Success: with low score",
			setup: func(tr *mock.MockTransactionRepository, ur *mock.MockUserRepository, sr *mock.MockScoreRepository, rr *mock.MockRankingRepository) {
				user := model.User{
					ID:        userID,
					Name:      "test",
					Email:     "test@gmail.com",
					Coins:     100,
					HighScore: 1000,
				}
				ur.EXPECT().Get(ctx, userID).Return(&user, nil)
				tr.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				})
				sr.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				ur.EXPECT().Update(
					gomock.Any(),
					model.User{
						ID:        user.ID,
						Name:      user.Name,
						Email:     user.Email,
						Coins:     400,
						HighScore: 1000,
					},
				).Return(nil)
				rr.EXPECT().Create(
					gomock.Any(),
					model.ScoreBoardKey,
					&model.Ranking{
						UserName: user.Name,
						Score:    100,
					},
				).Return(nil)
			},
			arg: struct {
				ctx   context.Context
				score int
			}{
				ctx:   ctx,
				score: 100,
			},
			want: struct {
				coin int
				err  error
			}{
				coin: func() int {
					game := &model.Game{}
					return game.Reward(100)
				}(),
				err: nil,
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			tr := mock.NewMockTransactionRepository(ctrl)
			ur := mock.NewMockUserRepository(ctrl)
			ucr := mock.NewMockUserCollectionRepository(ctrl)
			sr := mock.NewMockScoreRepository(ctrl)
			rr := mock.NewMockRankingRepository(ctrl)
			cr := mock.NewMockCollectionRepository(ctrl)
			ccr := mock.NewMockCollectionCacheRepository(ctrl)

			if tt.setup != nil {
				tt.setup(tr, ur, sr, rr)
			}

			usecase := NewGameUseCase(tr, ur, ucr, sr, rr, cr, ccr)
			coin, err := usecase.FinishGame(tt.arg.ctx, tt.arg.score)

			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("FinishGame() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("FinishGame() error = %v, wantErr %v", err, tt.want.err)
			}

			if coin != tt.want.coin {
				t.Errorf("FinishGame() coin = %v, want %v", coin, tt)
			}
		})
	}
}

func TestUsecase_DrawGacha(t *testing.T) {
	t.Parallel()

	userID := uuid.New().String()
	collection1ID := uuid.New().String()
	collection2ID := uuid.New().String()
	collection3ID := uuid.New().String()
	collection4ID := uuid.New().String()

	user := &model.User{
		ID:    userID,
		Name:  "test",
		Email: "test@gmail.com",
		Coins: 10000,
	}
	ctx := context.WithValue(context.Background(), config.ContextUserIDKey, userID)
	collections := model.Collections{
		{
			ID:     collection1ID,
			Name:   "collection1",
			Rarity: 1,
			Weight: 10,
		},
		{
			ID:     collection2ID,
			Name:   "collection2",
			Rarity: 2,
			Weight: 10,
		},
		{
			ID:     collection3ID,
			Name:   "collection3",
			Rarity: 3,
			Weight: 10,
		},
		{
			ID:     collection4ID,
			Name:   "collection4",
			Rarity: 4,
			Weight: 10,
		},
	}
	userCollections := []*model.UserCollection{
		{
			UserID:       userID,
			CollectionID: collection1ID,
		},
		{
			UserID:       userID,
			CollectionID: collection2ID,
		},
	}

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockTransactionRepository,
			m1 *mock.MockUserRepository,
			m2 *mock.MockCollectionRepository,
			m3 *mock.MockCollectionCacheRepository,
			m4 *mock.MockUserCollectionRepository,
		)
		arg struct {
			ctx   context.Context
			times int
		}
		want struct {
			results []*GachaResult
			err     error
		}
	}{
		{
			name: "success",
			setup: func(
				tr *mock.MockTransactionRepository,
				ur *mock.MockUserRepository,
				cr *mock.MockCollectionRepository,
				ccr *mock.MockCollectionCacheRepository,
				ucr *mock.MockUserCollectionRepository,
			) {
				ur.EXPECT().Get(ctx, userID).Return(user, nil)
				ccr.EXPECT().Get(ctx, "collections").Return(
					collections,
					nil,
				)
				ucr.EXPECT().List(ctx, userID).Return(userCollections, nil)
				tr.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				})
				ur.EXPECT().Update(ctx, model.User{
					ID:    userID,
					Name:  "test",
					Email: "test@gmail.com",
					Coins: 10000 - 2*model.GachaCost,
				}).Return(nil)
				ucr.EXPECT().BatchCreate(
					ctx,
					gomock.Any(),
				).Return(nil)
			},
			arg: struct {
				ctx   context.Context
				times int
			}{
				ctx:   ctx,
				times: 2,
			},
			want: struct {
				results []*GachaResult
				err     error
			}{
				results: nil,
				err:     nil,
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			tr := mock.NewMockTransactionRepository(ctrl)
			ur := mock.NewMockUserRepository(ctrl)
			ucr := mock.NewMockUserCollectionRepository(ctrl)
			sr := mock.NewMockScoreRepository(ctrl)
			rr := mock.NewMockRankingRepository(ctrl)
			cr := mock.NewMockCollectionRepository(ctrl)
			ccr := mock.NewMockCollectionCacheRepository(ctrl)

			if tt.setup != nil {
				tt.setup(tr, ur, cr, ccr, ucr)
			}

			usecase := NewGameUseCase(tr, ur, ucr, sr, rr, cr, ccr)
			_, err := usecase.DrawGacha(tt.arg.ctx, tt.arg.times)

			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("DrawGacha() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("DrawGacha() error = %v, wantErr %v", err, tt.want.err)
			}

			// if !reflect.DeepEqual(getResults, tt.want.results) {
			//  	t.Errorf("DrawGacha() results = %v, want %v", getResults, tt.want.results)
			// }
		})
	}
}
