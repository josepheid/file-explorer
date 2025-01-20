package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goteleport-interview/fs4/api/internal/sessions"
)

func TestLogoutHandler(t *testing.T) {
	tests := []struct {
		name          string
		method        string
		sessionCookie *http.Cookie
		wantStatus    int
		wantSetCookie bool
	}{
		{
			name:   "successful logout",
			method: http.MethodPost,
			sessionCookie: &http.Cookie{
				Name:  "session_id",
				Value: "test-session",
			},
			wantStatus:    http.StatusOK,
			wantSetCookie: true,
		},
		{
			name:          "already logged out",
			method:        http.MethodPost,
			sessionCookie: nil,
			wantStatus:    http.StatusOK,
			wantSetCookie: false,
		},
		{
			name:          "wrong method",
			method:        http.MethodGet,
			sessionCookie: nil,
			wantStatus:    http.StatusMethodNotAllowed,
			wantSetCookie: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sessions := sessions.New()
			handler := NewLogoutHandler(sessions)

			req := httptest.NewRequest(tt.method, "/api/v1/auth/logout", nil)
			if tt.sessionCookie != nil {
				req.AddCookie(tt.sessionCookie)
			}

			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			// Check status code
			if rec.Code != tt.wantStatus {
				t.Errorf("status code = %v, want %v", rec.Code, tt.wantStatus)
			}

			// Check cookie deletion
			cookies := rec.Result().Cookies()
			hasCookie := len(cookies) > 0
			if hasCookie != tt.wantSetCookie {
				t.Errorf("cookie present = %v, want %v", hasCookie, tt.wantSetCookie)
			}

			if hasCookie {
				cookie := cookies[0]
				if cookie.MaxAge != -1 {
					t.Error("cookie not marked for deletion")
				}
			}
		})
	}
}
