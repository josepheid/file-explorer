package handlers

import (
	"net/http"

	"github.com/goteleport-interview/fs4/api/internal/respond"
	"github.com/goteleport-interview/fs4/api/internal/sessions"
)

// LogoutHandler defines the logout handler and the dependencies it needs
type LogoutHandler struct {
	sessions *sessions.Service
}

// NewLogoutHandler creates a new LogoutHandler, it takes a sessions service as a parameter
func NewLogoutHandler(sessions *sessions.Service) *LogoutHandler {
	return &LogoutHandler{
		sessions: sessions,
	}
}

// ServeHTTP handles the logout request
func (h *LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respond.WithError(w, "Method not allowed, method: "+r.Method, http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		// Already logged out
		respond.WithJSON(w, nil, http.StatusOK)
		return
	}

	// Delete session
	h.sessions.Delete(cookie.Value)

	// Invalidate cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1, // Delete cookie
	})

	respond.WithJSON(w, nil, http.StatusOK)
}
