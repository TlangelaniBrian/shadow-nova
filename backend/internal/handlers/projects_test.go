package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"shadow-nova/backend/internal/database"
	"shadow-nova/backend/internal/models"
	"testing"
)

func TestProjectsHandler_List(t *testing.T) {
	// Setup
	mockDB := &database.MockService{
		GetProjectsFunc: func(ctx context.Context) ([]models.Project, error) {
			return []models.Project{
				{ID: "p1", Title: "Project 1"},
				{ID: "p2", Title: "Project 2"},
			}, nil
		},
	}
	handler := NewProjectsHandler(mockDB)

	// Request
	req := httptest.NewRequest("GET", "/projects", nil)
	rr := httptest.NewRecorder()

	// Execute
	handler.List(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response models.SuccessResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	projects := response.Data.([]interface{})
	if len(projects) != 2 {
		t.Errorf("expected 2 projects, got %d", len(projects))
	}
}

func TestProjectsHandler_Create(t *testing.T) {
	// Setup
	mockDB := &database.MockService{
		CreateProjectFunc: func(ctx context.Context, project *models.Project) error {
			return nil
		},
	}
	handler := NewProjectsHandler(mockDB)

	payload := models.CreateProjectRequest{
		ID:    "newproject",
		Title: "New Project",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/projects", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// Execute
	handler.Create(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}

func TestProjectsHandler_Submit(t *testing.T) {
	// Setup
	mockDB := &database.MockService{
		SubmitProjectFunc: func(ctx context.Context, sub *models.ProjectSubmission) error {
			sub.ID = 123
			return nil
		},
	}
	handler := NewProjectsHandler(mockDB)

	payload := models.SubmitProjectRequest{
		ProjectID:     "p1",
		GithubRepoURL: "https://github.com/user/repo",
		PRURL:         "https://github.com/user/repo/pull/1",
		DemoURL:       "https://demo.com",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/submissions", bytes.NewBuffer(body))
	
	// Mock Context with UserID (simulate auth middleware)
	ctx := context.WithValue(req.Context(), "user_id", 1)
	req = req.WithContext(ctx)
	
	rr := httptest.NewRecorder()

	// Execute
	handler.Submit(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}
