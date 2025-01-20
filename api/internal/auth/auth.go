package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// Service handles user authentication operations
type Service struct {
	// In a real app, this would be a database
	// For this challenge, we'll use an in-memory map
	users map[string]string // username -> hashed password
}

// Authenticator defines the interface for authentication operations
type Authenticator interface {
	// ValidateCredentials checks if the provided username and password are valid
	ValidateCredentials(username, password string) error
}

// It has been generated outside of the ValidateCredentials function to avoid creating a new hash for each request and potentially highlighting a timing difference
var dummyHash, _ = bcrypt.GenerateFromPassword(make([]byte, 60), bcrypt.DefaultCost)

func New() *Service {
	s := &Service{
		users: make(map[string]string),
	}

	// Add a test user (in production, this would be in a database)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	s.users["testuser"] = string(hashedPassword)

	return s
}

func (s *Service) ValidateCredentials(username, password string) error {
	hashedPassword, exists := s.users[username]
	if !exists {
		// If the user doesn't exist, compare the password with the dummy hash and return an error anyway
		// This will take the same amount of time as if the user existed
		// This prevents timing attacks
		_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(password))
		return bcrypt.ErrMismatchedHashAndPassword
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err
}
