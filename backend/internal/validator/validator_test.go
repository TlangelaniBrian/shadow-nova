package validator

import (
	"testing"

	"shadow-nova/backend/internal/models"
)

func TestValidateRegistration(t *testing.T) {
	tests := []struct {
		name      string
		request   models.RegisterRequest
		shouldErr bool
	}{
		{
			name: "valid registration",
			request: models.RegisterRequest{
				Email:    "user@example.com",
				Username: "johndoe",
				Password: "SecurePass123",
			},
			shouldErr: false,
		},
		{
			name: "invalid email",
			request: models.RegisterRequest{
				Email:    "not-an-email",
				Username: "johndoe",
				Password: "SecurePass123",
			},
			shouldErr: true,
		},
		{
			name: "username too short",
			request: models.RegisterRequest{
				Email:    "user@example.com",
				Username: "ab",
				Password: "SecurePass123",
			},
			shouldErr: true,
		},
		{
			name: "password too short",
			request: models.RegisterRequest{
				Email:    "user@example.com",
				Username: "johndoe",
				Password: "weak",
			},
			shouldErr: true,
		},
		{
			name: "username with special chars",
			request: models.RegisterRequest{
				Email:    "user@example.com",
				Username: "john@doe",
				Password: "SecurePass123",
			},
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.request)
			
			if tt.shouldErr && err == nil {
				t.Errorf("expected error but got none")
			}
			
			if !tt.shouldErr && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
		})
	}
}

func TestValidateLogin(t *testing.T) {
	tests := []struct {
		name      string
		request   models.LoginRequest
		shouldErr bool
	}{
		{
			name: "valid login",
			request: models.LoginRequest{
				Email:    "user@example.com",
				Password: "password123",
			},
			shouldErr: false,
		},
		{
			name: "missing email",
			request: models.LoginRequest{
				Email:    "",
				Password: "password123",
			},
			shouldErr: true,
		},
		{
			name: "missing password",
			request: models.LoginRequest{
				Email:    "user@example.com",
				Password: "",
			},
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.request)
			
			if tt.shouldErr && err == nil {
				t.Errorf("expected error but got none")
			}
			
			if !tt.shouldErr && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
		})
	}
}
