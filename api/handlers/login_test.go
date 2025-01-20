package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/josepheid/file-explorer/api/internal/auth"
	"github.com/josepheid/file-explorer/api/internal/sessions"
)

func TestLoginHandler(t *testing.T) {
	tests := []struct {
		name       string
		request    LoginRequest
		method     string
		wantStatus int
		wantCookie bool
	}{
		{
			name: "valid credentials",
			request: LoginRequest{
				Username: "testuser",
				Password: "password123",
			},
			method:     http.MethodPost,
			wantStatus: http.StatusOK,
			wantCookie: true,
		},
		{
			name: "invalid password",
			request: LoginRequest{
				Username: "testuser",
				Password: "wrongpass",
			},
			method:     http.MethodPost,
			wantStatus: http.StatusUnauthorized,
			wantCookie: false,
		},
		{
			name: "user not found",
			request: LoginRequest{
				Username: "nonexistentuser",
				Password: "password123",
			},
			method:     http.MethodPost,
			wantStatus: http.StatusUnauthorized,
			wantCookie: false,
		},
		{
			name:       "empty request",
			request:    LoginRequest{},
			method:     http.MethodPost,
			wantStatus: http.StatusBadRequest,
			wantCookie: false,
		},
		{
			name: "wrong method",
			request: LoginRequest{
				Username: "testuser",
				Password: "password123",
			},
			method:     http.MethodGet,
			wantStatus: http.StatusMethodNotAllowed,
			wantCookie: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create handler with mock services
			auth := auth.New()
			sessions := sessions.New()
			handler := NewLoginHandler(auth, sessions)

			// Create request
			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest(tt.method, "/api/v1/auth/login", bytes.NewBuffer(body))
			rec := httptest.NewRecorder()

			// Handle request
			handler.ServeHTTP(rec, req)

			// Check status
			if rec.Code != tt.wantStatus {
				t.Errorf("want status %d, got %d", tt.wantStatus, rec.Code)
			}

			// Check cookie
			cookies := rec.Result().Cookies()
			hasCookie := len(cookies) > 0
			if hasCookie != tt.wantCookie {
				t.Errorf("want cookie: %v, got cookie: %v", tt.wantCookie, hasCookie)
			}
			if tt.wantCookie {
				if cookies[0].Value == "" {
					t.Errorf("cookie value is empty")
				}
				if cookies[0].Name != "session_id" {
					t.Errorf("cookie name is not 'session_id'")
				}
			}
		})
	}
}
