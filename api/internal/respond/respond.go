package respond

import (
	"encoding/json"
	"net/http"
)

// Error represents an error response
type Error struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

func WithError(w http.ResponseWriter, msg string, status int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Error{
		Message:    msg,
		StatusCode: status,
	})
}

func WithJSON(w http.ResponseWriter, v any, status int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		WithError(w, "failed to encode", http.StatusInternalServerError)
	}
}
