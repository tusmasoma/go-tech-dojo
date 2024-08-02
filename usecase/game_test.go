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
