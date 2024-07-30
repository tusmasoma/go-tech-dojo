package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/tusmasoma/go-tech-dojo/domain/model"
	"github.com/tusmasoma/go-tech-dojo/usecase/mock"
)

func TestRankingHandler_ListRankings(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockRankingUseCase,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(m *mock.MockRankingUseCase) {
				m.EXPECT().ListRankings(
					gomock.Any(),
					1,
				).Return(
					[]*model.Ranking{
						{
							UserName: "test",
							Score:    1000,
							Rank:     1,
						},
						{
							UserName: "test2",
							Score:    900,
							Rank:     2,
						},
					},
					nil,
				)
			},
			in: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/api/ranking/list?start=1", nil)
				return req
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Fail: Invalid start parameter",
			in: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/api/ranking/list", nil)
				return req
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			ruc := mock.NewMockRankingUseCase(ctrl)

			if tt.setup != nil {
				tt.setup(ruc)
			}

			handler := NewRankingHandler(ruc)
			recorder := httptest.NewRecorder()
			handler.ListRankings(recorder, tt.in())

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
		})
	}
}
