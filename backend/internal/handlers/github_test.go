package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"shadow-nova/backend/internal/auth"
	"shadow-nova/backend/internal/database"
	"testing"
)

func TestGitHubHandler_Disconnect(t *testing.T) {
	// Setup
	mockDB := &database.MockService{
		DeleteGitHubIntegrationFunc: func(ctx context.Context, userID int) error {
			return nil
		},
	}

	authService := auth.NewGitHubAuthService()
	handler := NewGitHubHandler(authService, mockDB)

	// Request
	req := httptest.NewRequest("DELETE", "/auth/github/disconnect", nil)

	// Mock Context with UserID (simulate auth middleware)
	ctx := context.WithValue(req.Context(), "user_id", 1)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	// Execute
	handler.Disconnect(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if message, ok := response["message"]; !ok || message != "GitHub account disconnected successfully" {
		t.Errorf("expected success message, got %v", response)
	}
}

func TestGitHubHandler_Disconnect_NoAuth(t *testing.T) {
	// Setup
	mockDB := &database.MockService{}
	authService := auth.NewGitHubAuthService()
	handler := NewGitHubHandler(authService, mockDB)

	// Request without user_id in context
	req := httptest.NewRequest("DELETE", "/auth/github/disconnect", nil)
	rr := httptest.NewRecorder()

	// Execute
	handler.Disconnect(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}
}
