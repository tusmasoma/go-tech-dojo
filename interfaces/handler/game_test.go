package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"github.com/tusmasoma/go-tech-dojo/domain/model"
	"github.com/tusmasoma/go-tech-dojo/usecase"
	"github.com/tusmasoma/go-tech-dojo/usecase/mock"
)

func TestGameHandler_FinishGame(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockGameUseCase,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(m *mock.MockGameUseCase) {
				m.EXPECT().FinishGame(gomock.Any(), 100).Return(300, nil)
			},
			in: func() *http.Request {
				gameFinishReq := FinishGameRequest{Score: 100}
				reqBody, _ := json.Marshal(gameFinishReq)
				req, _ := http.NewRequest(http.MethodPut, "/api/game/finish", bytes.NewBuffer(reqBody))
				return req
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Fail: Invalid request",
			in: func() *http.Request {
				gameFinishReq := FinishGameRequest{Score: -100}
				reqBody, _ := json.Marshal(gameFinishReq)
				req, _ := http.NewRequest(http.MethodPut, "/api/game/finish", bytes.NewBuffer(reqBody))
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
			guc := mock.NewMockGameUseCase(ctrl)

			if tt.setup != nil {
				tt.setup(guc)
			}

			handler := NewGameHandler(guc)
			recorder := httptest.NewRecorder()
			handler.FinishGame(recorder, tt.in())

			if status := recorder.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
		})
	}
}

func TestGameHandler_DrawGacha(t *testing.T) {
	t.Parallel()

	collectionID := uuid.New().String()

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockGameUseCase,
		)
		in           func() *http.Request
		wantStatus   int
		wantResponse DrawGachaResponse
	}{
		{
			name: "success",
			setup: func(m *mock.MockGameUseCase) {
				m.EXPECT().DrawGacha(
					gomock.Any(),
					1,
				).Return(
					[]*usecase.GachaResult{
						{
							Collection: &model.Collection{
								ID:     collectionID,
								Name:   "collection1",
								Rarity: 1,
								Weight: 10,
							},
							Has: false,
						},
					},
					nil,
				)
			},
			in: func() *http.Request {
				drawGachaReq := DrawGachaRequest{
					Times: 1,
				}
				reqBody, _ := json.Marshal(drawGachaReq)
				req, _ := http.NewRequest(http.MethodPut, "/api/gacha/draw", bytes.NewBuffer(reqBody))
				return req
			},
			wantStatus: http.StatusOK,
			wantResponse: DrawGachaResponse{
				Results: []struct {
					ID     string `json:"id"`
					Name   string `json:"name"`
					Rarity int    `json:"rarity"`
					IsNew  bool   `json:"is_new"`
				}{
					{
						ID:     collectionID,
						Name:   "collection1",
						Rarity: 1,
						IsNew:  true,
					},
				},
			},
		},
		{
			name: "Fail: Invalid request",
			in: func() *http.Request {
				drawGachaReq := DrawGachaRequest{
					Times: -1,
				}
				reqBody, _ := json.Marshal(drawGachaReq)
				req, _ := http.NewRequest(http.MethodPut, "/api/gacha/draw", bytes.NewBuffer(reqBody))
				return req
			},
			wantStatus:   http.StatusBadRequest,
			wantResponse: DrawGachaResponse{},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			guc := mock.NewMockGameUseCase(ctrl)

			if tt.setup != nil {
				tt.setup(guc)
			}

			handler := NewGameHandler(guc)
			recorder := httptest.NewRecorder()
			handler.DrawGacha(recorder, tt.in())

			if status := recorder.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}

			if tt.wantStatus == http.StatusOK {
				var gotResponse DrawGachaResponse
				err := json.NewDecoder(recorder.Body).Decode(&gotResponse)
				if err != nil {
					t.Fatalf("failed to decode response body: %v", err)
				}
				if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
					t.Errorf("handler returned unexpected body: got %v want %v", gotResponse, tt.wantResponse)
				}
			}
		})
	}
}
