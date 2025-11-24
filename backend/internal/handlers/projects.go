package handlers

import (
	"encoding/json"
	"net/http"
	"shadow-nova/backend/internal/database"
	"shadow-nova/backend/internal/models"
	"shadow-nova/backend/internal/validator"
)

type ProjectsHandler struct {
	db database.Service
}

func NewProjectsHandler(db database.Service) *ProjectsHandler {
	return &ProjectsHandler{db: db}
}

func (h *ProjectsHandler) List(w http.ResponseWriter, r *http.Request) {
	projects, err := h.db.GetProjects(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch projects", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.SuccessResponse{
		Message: "Projects retrieved successfully",
		Data:    projects,
	})
}

func (h *ProjectsHandler) Create(w http.ResponseWriter, r *http.Request) {
	// TODO: Add Admin check here
	var req models.CreateProjectRequest
	if err := validator.ValidateRequest(r, &req); err != nil {
		validator.WriteValidationError(w, err)
		return
	}

	project := &models.Project{
		ID:          req.ID,
		Title:       req.Title,
		Description: req.Description,
		Difficulty:  req.Difficulty,
		TechStack:   req.TechStack,
	}

	if err := h.db.CreateProject(r.Context(), project); err != nil {
		http.Error(w, "Failed to create project", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.SuccessResponse{
		Message: "Project created successfully",
		Data:    project,
	})
}

func (h *ProjectsHandler) Submit(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	var req models.SubmitProjectRequest
	if err := validator.ValidateRequest(r, &req); err != nil {
		validator.WriteValidationError(w, err)
		return
	}

	submission := &models.ProjectSubmission{
		UserID:        userID,
		ProjectID:     req.ProjectID,
		GithubRepoURL: req.GithubRepoURL,
		PRURL:         req.PRURL,
		DemoURL:       req.DemoURL,
	}

	if err := h.db.SubmitProject(r.Context(), submission); err != nil {
		http.Error(w, "Failed to submit project", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.SuccessResponse{
		Message: "Project submitted successfully",
		Data:    submission,
	})
}
