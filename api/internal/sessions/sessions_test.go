package sessions

import (
	"testing"
	"time"
)

func TestSessionService(t *testing.T) {
	service := New()

	// Test session creation
	t.Run("create session", func(t *testing.T) {
		session, err := service.Create("testuser")
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}

		if session.UserID != "testuser" {
			t.Errorf("Expected UserID %q, got %q", "testuser", session.UserID)
		}

		if session.ID == "" {
			t.Error("Session ID should not be empty")
		}
	})

	// Test session retrieval
	t.Run("get session", func(t *testing.T) {
		session, err := service.Create("testuser")

		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}

		got := service.Get(session.ID)
		if got == (Session{}) {
			t.Fatal("Expected to get session, got nil")
		}

		if got.UserID != "testuser" {
			t.Errorf("Expected UserID %q, got %q", "testuser", got.UserID)
		}
	})

	// Test session expiration
	t.Run("expired session", func(t *testing.T) {
		session, err := service.Create("testuser")

		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}

		// Manually expire the session
		service.sessions[session.ID].ExpiresAt = time.Now().Add(-time.Hour)

		got := service.Get(session.ID)
		if got != (Session{}) {
			t.Error("Expected nil for expired session, got session")
		}
	})

	// Test session deletion
	t.Run("delete session", func(t *testing.T) {
		session, err := service.Create("testuser")
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}

		service.Delete(session.ID)

		if got := service.Get(session.ID); got != (Session{}) {
			t.Error("Expected nil after deletion, got session")
		}
	})
}
