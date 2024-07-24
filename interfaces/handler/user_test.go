package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/tusmasoma/go-tech-dojo/usecase/mock"
)

func TestUserHandler_CreateUser(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockUserUseCase,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(m *mock.MockUserUseCase) {
				m.EXPECT().CreateUserAndToken(
					gomock.Any(),
					"test@gmail.com",
					"password123",
				).Return(
					"eyJhbGciOiJIUzI1NiIsI.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0Ijo.SflKxwRJSMeKKF2QT4fwpMeJf36P",
					nil,
				)
			},
			in: func() *http.Request {
				userCreateReq := CreateUserRequest{Email: "test@gmail.com", Password: "password123"}
				reqBody, _ := json.Marshal(userCreateReq)
				req, _ := http.NewRequest(http.MethodPost, "/api/user/create", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Fail: invalid request",
			in: func() *http.Request {
				userCreateReq := CreateUserRequest{Email: "test@gmail.com"}
				reqBody, _ := json.Marshal(userCreateReq)
				req, _ := http.NewRequest(http.MethodPost, "/api/user/create", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
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
			uuc := mock.NewMockUserUseCase(ctrl)

			if tt.setup != nil {
				tt.setup(uuc)
			}

			handler := NewUserHandler(uuc)
			recorder := httptest.NewRecorder()
			handler.CreateUser(recorder, tt.in())

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
			if tt.wantStatus == http.StatusOK {
				if token := recorder.Header().Get("Authorization"); token == "" || strings.TrimPrefix(token, "Bearer ") == "" {
					t.Fatalf("Expected Authorization header to be set")
				}
			}
		})
	}
}
