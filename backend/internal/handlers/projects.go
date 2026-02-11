package handlers

import (
	"net/http"
	"shadow-nova/backend/internal/database"
	"shadow-nova/backend/internal/httputil"
	"shadow-nova/backend/internal/middleware"
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
		httputil.WriteError(w, http.StatusInternalServerError, "Failed to fetch projects")
		return
	}

	httputil.WriteSuccess(w, "Projects retrieved successfully", projects)
}

func (h *ProjectsHandler) Create(w http.ResponseWriter, r *http.Request) {
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
		httputil.WriteError(w, http.StatusInternalServerError, "Failed to create project")
		return
	}

	httputil.WriteCreated(w, "Project created successfully", project)
}

func (h *ProjectsHandler) Submit(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		httputil.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

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
		httputil.WriteError(w, http.StatusInternalServerError, "Failed to submit project")
		return
	}

	httputil.WriteCreated(w, "Project submitted successfully", submission)
}
