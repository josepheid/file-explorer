package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/josepheid/file-explorer/api/internal/auth"
	"github.com/josepheid/file-explorer/api/internal/respond"
	"github.com/josepheid/file-explorer/api/internal/sessions"
)

// LoginHandler defines the login handler and the dependencies it needs
type LoginHandler struct {
	auth     *auth.Service
	sessions *sessions.Service
}

// NewLoginHandler creates a new LoginHandler, it takes an auth service and a sessions service as parameters
func NewLoginHandler(auth *auth.Service, sessions *sessions.Service) *LoginHandler {
	return &LoginHandler{
		auth:     auth,
		sessions: sessions,
	}
}

// LoginRequest represents the request body for the login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ServeHTTP handles the login request
func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respond.WithError(w, "Method not allowed, method: "+r.Method, http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.WithError(w, "Invalid request body, error: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		respond.WithError(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Validate credentials
	if err := h.auth.ValidateCredentials(req.Username, req.Password); err != nil {
		respond.WithError(w, "Invalid credentials, error: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Create new session
	session, err := h.sessions.Create(req.Username)
	if err != nil {
		respond.WithError(w, "Failed to create session, error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   86400, // 24 hours
	})

	// Return success response
	respond.WithJSON(w, nil, http.StatusOK)
}
