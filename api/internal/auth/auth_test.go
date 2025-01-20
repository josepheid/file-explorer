package auth

import (
	"errors"
	"testing"
)

func TestValidateCredentials(t *testing.T) {
	service := New()

	tests := []struct {
		name        string
		username    string
		password    string
		expectedErr error
	}{
		{
			name:        "valid credentials",
			username:    "testuser",
			password:    "password123",
			expectedErr: nil,
		},
		{
			name:        "invalid password",
			username:    "testuser",
			password:    "wrongpassword",
			expectedErr: errors.New("crypto/bcrypt: hashedPassword is not the hash of the given password"),
		},
		{
			name:        "non-existent user, should behave the same as an invalid password",
			username:    "bruteforceattempt",
			password:    "password123",
			expectedErr: errors.New("crypto/bcrypt: hashedPassword is not the hash of the given password"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidateCredentials(tt.username, tt.password)
			if (err != nil) != (tt.expectedErr != nil) {
				t.Errorf("ValidateCredentials() error = %v, want error %v", err, tt.expectedErr)
				return
			}
			if err != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("ValidateCredentials() error = %v, want error %v", err, tt.expectedErr)
			}
		})
	}
}
