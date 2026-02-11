package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"shadow-nova/backend/internal/database"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestValidatePathOwnership(t *testing.T) {
	tests := []struct {
		name           string
		userID         int
		pathID         string
		mockFunc       func(ctx context.Context, userID int, pathID string) (bool, error)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Owner can access path",
			userID: 1,
			pathID: "web-dev",
			mockFunc: func(ctx context.Context, userID int, pathID string) (bool, error) {
				return true, nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
		},
		{
			name:   "Non-owner gets forbidden",
			userID: 2,
			pathID: "web-dev",
			mockFunc: func(ctx context.Context, userID int, pathID string) (bool, error) {
				return false, nil
			},
			expectedStatus: http.StatusForbidden,
			expectedBody:   "You do not have access to this path",
		},
		{
			name:   "Database error returns 500",
			userID: 1,
			pathID: "web-dev",
			mockFunc: func(ctx context.Context, userID int, pathID string) (bool, error) {
				return false, errors.New("database error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Failed to verify path access",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &database.MockService{
				UserHasAccessToPathFunc: tt.mockFunc,
			}

			// Create handler with middleware
			handler := ValidatePathOwnership(mockDB)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("success"))
			}))

			// Create request with context and path param
			req := httptest.NewRequest("GET", "/paths/"+tt.pathID, nil)
			ctx := context.WithValue(req.Context(), UserIDKey, tt.userID)
			req = req.WithContext(ctx)

			// Set up chi URL params
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.pathID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			// Execute
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			// Assert
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.expectedBody != "" && !contains(rr.Body.String(), tt.expectedBody) {
				t.Errorf("handler returned unexpected body: got %v want to contain %v", rr.Body.String(), tt.expectedBody)
			}
		})
	}
}

func TestValidatePathOwnership_MissingUserID(t *testing.T) {
	mockDB := &database.MockService{}
	handler := ValidatePathOwnership(mockDB)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/paths/web-dev", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}
}

func TestValidatePathOwnership_MissingPathID(t *testing.T) {
	mockDB := &database.MockService{}
	handler := ValidatePathOwnership(mockDB)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/paths/", nil)
	ctx := context.WithValue(req.Context(), UserIDKey, 1)
	req = req.WithContext(ctx)

	// Set up chi route context without path ID
	rctx := chi.NewRouteContext()
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestValidateSubmissionOwnership(t *testing.T) {
	tests := []struct {
		name           string
		userID         int
		submissionID   string
		mockFunc       func(ctx context.Context, userID int, submissionID int) (bool, error)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:         "Owner can access submission",
			userID:       1,
			submissionID: "123",
			mockFunc: func(ctx context.Context, userID int, submissionID int) (bool, error) {
				return true, nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
		},
		{
			name:         "Non-owner gets forbidden",
			userID:       2,
			submissionID: "123",
			mockFunc: func(ctx context.Context, userID int, submissionID int) (bool, error) {
				return false, nil
			},
			expectedStatus: http.StatusForbidden,
			expectedBody:   "You do not have access to this submission",
		},
		{
			name:         "Database error returns 500",
			userID:       1,
			submissionID: "123",
			mockFunc: func(ctx context.Context, userID int, submissionID int) (bool, error) {
				return false, errors.New("database error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Failed to verify submission ownership",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &database.MockService{
				UserOwnsSubmissionFunc: tt.mockFunc,
			}

			handler := ValidateSubmissionOwnership(mockDB)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("success"))
			}))

			req := httptest.NewRequest("GET", "/submissions/"+tt.submissionID, nil)
			ctx := context.WithValue(req.Context(), UserIDKey, tt.userID)
			req = req.WithContext(ctx)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.submissionID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.expectedBody != "" && !contains(rr.Body.String(), tt.expectedBody) {
				t.Errorf("handler returned unexpected body: got %v want to contain %v", rr.Body.String(), tt.expectedBody)
			}
		})
	}
}

func TestValidateSubmissionOwnership_InvalidID(t *testing.T) {
	mockDB := &database.MockService{}
	handler := ValidateSubmissionOwnership(mockDB)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/submissions/invalid", nil)
	ctx := context.WithValue(req.Context(), UserIDKey, 1)
	req = req.WithContext(ctx)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "invalid")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestValidateProgressOwnership(t *testing.T) {
	tests := []struct {
		name           string
		userID         int
		progressID     string
		mockFunc       func(ctx context.Context, userID int, progressID int) (bool, error)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:       "Owner can access progress",
			userID:     1,
			progressID: "456",
			mockFunc: func(ctx context.Context, userID int, progressID int) (bool, error) {
				return true, nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
		},
		{
			name:       "Non-owner gets forbidden",
			userID:     2,
			progressID: "456",
			mockFunc: func(ctx context.Context, userID int, progressID int) (bool, error) {
				return false, nil
			},
			expectedStatus: http.StatusForbidden,
			expectedBody:   "You do not have access to this progress record",
		},
		{
			name:       "Database error returns 500",
			userID:     1,
			progressID: "456",
			mockFunc: func(ctx context.Context, userID int, progressID int) (bool, error) {
				return false, errors.New("database error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Failed to verify progress ownership",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &database.MockService{
				UserOwnsProgressFunc: tt.mockFunc,
			}

			handler := ValidateProgressOwnership(mockDB)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("success"))
			}))

			req := httptest.NewRequest("GET", "/progress/"+tt.progressID, nil)
			ctx := context.WithValue(req.Context(), UserIDKey, tt.userID)
			req = req.WithContext(ctx)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.progressID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.expectedBody != "" && !contains(rr.Body.String(), tt.expectedBody) {
				t.Errorf("handler returned unexpected body: got %v want to contain %v", rr.Body.String(), tt.expectedBody)
			}
		})
	}
}

func TestValidateProgressOwnership_InvalidID(t *testing.T) {
	mockDB := &database.MockService{}
	handler := ValidateProgressOwnership(mockDB)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/progress/invalid", nil)
	ctx := context.WithValue(req.Context(), UserIDKey, 1)
	req = req.WithContext(ctx)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "invalid")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && indexOf(s, substr) >= 0))
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
