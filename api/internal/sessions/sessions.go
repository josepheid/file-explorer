package sessions

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
	"time"
)

// Session represents a user session with expiration
type Session struct {
	// ID uniquely identifies the session
	ID string
	// UserID identifies the user this session belongs to
	UserID string
	// CreatedAt indicates when the session was created
	CreatedAt time.Time
	// ExpiresAt indicates when the session expires
	ExpiresAt time.Time
}

// Service manages user sessions including creation, retrieval, and deletion
type Service struct {
	sessions map[string]*Session
	mu       sync.RWMutex
}

func New() *Service {
	return &Service{
		sessions: make(map[string]*Session),
	}
}

func (s *Service) Create(userID string) (Session, error) {
	// Generate random session ID
	b := make([]byte, 256)
	if _, err := rand.Read(b); err != nil {
		return Session{}, err
	}
	sessionID := base64.URLEncoding.EncodeToString(b)

	session := Session{
		ID:        sessionID,
		UserID:    userID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	s.mu.Lock()
	s.sessions[sessionID] = &session
	s.mu.Unlock()

	return session, nil
}

func (s *Service) Get(sessionID string) Session {
	s.mu.RLock()
	defer s.mu.RUnlock()
	session, exists := s.sessions[sessionID]
	if !exists {
		return Session{}
	}
	// Check if session is expired
	if time.Now().After(session.ExpiresAt) {
		// use a go routine to delete the session without blocking the read operation
		go s.Delete(sessionID)
		return Session{}
	}

	return *session
}

func (s *Service) Delete(sessionID string) {
	s.mu.Lock()
	delete(s.sessions, sessionID)
	s.mu.Unlock()
}
