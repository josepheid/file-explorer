package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/josepheid/file-explorer/api/internal/sessions"
)

func TestRequireAuth(t *testing.T) {
	session := sessions.New()

	// Create a test handler that we'll wrap with auth
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name           string
		cookie         *http.Cookie
		expectedStatus int
	}{
		{
			name:           "no cookie",
			cookie:         nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid session",
			cookie:         &http.Cookie{Name: "session_id", Value: "invalid"},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "valid session",
			cookie: func() *http.Cookie {
				s, _ := session.Create("testuser")
				return &http.Cookie{Name: "session_id", Value: s.ID}
			}(),
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			if tt.cookie != nil {
				req.AddCookie(tt.cookie)
			}

			rr := httptest.NewRecorder()
			handler := RequireAuth(session)(testHandler)
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}
