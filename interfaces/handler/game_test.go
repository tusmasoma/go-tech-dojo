package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"

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
