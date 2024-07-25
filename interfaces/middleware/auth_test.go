package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/tusmasoma/go-tech-dojo/config"
	"github.com/tusmasoma/go-tech-dojo/pkg/auth"
)

func dummyTestHandler(w http.ResponseWriter, r *http.Request) {
	userIDValue := r.Context().Value(config.ContextUserIDKey)
	if userID, _ := userIDValue.(string); userID == "" {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func TestAuthMiddleware_Authenticate(t *testing.T) {
	t.Parallel()

	userID := uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2")
	email := "test@gmail.com"

	jwt, _ := auth.GenerateToken(userID.String(), email)

	patterns := []struct {
		name       string
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			in: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set("Authorization", "Bearer "+jwt)
				return req
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Fail: No Auth Header",
			in: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/", nil)
				return req
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "Fail: Invalid Auth Header Format",
			in: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set("Authorization", jwt)
				return req
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "Fail: Invalid Token",
			in: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set("Authorization", "Bearer "+"invalid Token")
				return req
			},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			am := NewAuthMiddleware()

			handler := am.Authenticate(http.HandlerFunc(dummyTestHandler))

			recoder := httptest.NewRecorder()
			handler.ServeHTTP(recoder, tt.in())

			// ステータスコードの検証
			if status := recoder.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
		})
	}
}
