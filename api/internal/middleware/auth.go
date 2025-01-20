package middleware

import (
	"net/http"

	"github.com/goteleport-interview/fs4/api/internal/respond"
	"github.com/goteleport-interview/fs4/api/internal/sessions"
)

func RequireAuth(ss *sessions.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_id")
			if err != nil {
				respond.WithError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			session := ss.Get(cookie.Value)
			if session == (sessions.Session{}) {
				respond.WithError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
