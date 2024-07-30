package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/tusmasoma/go-tech-dojo/domain/model"
	"github.com/tusmasoma/go-tech-dojo/domain/repository/mock"
)

func TestRankingUseCase_ListRankings(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockRankingRepository,
		)
		arg struct {
			ctx   context.Context
			start int
		}
		wantErr error
	}{
		{
			name: "success",
			setup: func(m *mock.MockRankingRepository) {
				m.EXPECT().List(
					gomock.Any(),
					model.ScoreBoardKey,
					1).Return([]*model.Ranking{
					{
						UserName: "user1",
						Score:    100,
						Rank:     1,
					},
				}, nil)
			},
			arg: struct {
				ctx   context.Context
				start int
			}{
				ctx:   context.Background(),
				start: 1,
			},
			wantErr: nil,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			rr := mock.NewMockRankingRepository(ctrl)
			if tt.setup != nil {
				tt.setup(rr)
			}

			ruc := NewRankingUseCase(rr)

			_, err := ruc.ListRankings(tt.arg.ctx, tt.arg.start)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("wantErr: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}
