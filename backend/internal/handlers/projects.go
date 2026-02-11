package handlers

import (
	"net/http"
	"shadow-nova/backend/internal/database"
	"shadow-nova/backend/internal/httputil"
	"shadow-nova/backend/internal/metrics"
	"shadow-nova/backend/internal/middleware"
	"shadow-nova/backend/internal/models"
	"shadow-nova/backend/internal/validator"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ProjectsHandler struct {
	db database.Service
}

func NewProjectsHandler(db database.Service) *ProjectsHandler {
	return &ProjectsHandler{db: db}
}

func (h *ProjectsHandler) List(w http.ResponseWriter, r *http.Request) {
	pagination := models.ParsePagination(r)
	offset := (pagination.Page - 1) * pagination.Limit

	projects, err := h.db.GetProjects(r.Context(), pagination.Limit, offset)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "Failed to fetch projects")
		return
	}

	total, err := h.db.GetProjectsCount(r.Context())
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "Failed to count projects")
		return
	}

	response := models.NewPaginatedResponse(projects, pagination.Page, pagination.Limit, total)
	httputil.WriteJSON(w, http.StatusOK, response)
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

	metrics.ProjectSubmissions.Inc()

	httputil.WriteCreated(w, "Project submitted successfully", submission)
}

func (h *ProjectsHandler) GetSubmission(w http.ResponseWriter, r *http.Request) {
	submissionIDStr := chi.URLParam(r, "id")
	submissionID, err := strconv.Atoi(submissionIDStr)
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "Invalid submission ID")
		return
	}

	submission, err := h.db.GetSubmission(r.Context(), submissionID)
	if err != nil {
		httputil.HandleError(w, err)
		return
	}

	httputil.WriteSuccess(w, "Submission retrieved successfully", submission)
}

func (h *ProjectsHandler) UpdateSubmission(w http.ResponseWriter, r *http.Request) {
	submissionIDStr := chi.URLParam(r, "id")
	submissionID, err := strconv.Atoi(submissionIDStr)
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "Invalid submission ID")
		return
	}

	var req struct {
		Status   string `json:"status" validate:"omitempty,oneof=pending approved rejected"`
		Feedback string `json:"feedback"`
	}

	if err := validator.ValidateRequest(r, &req); err != nil {
		validator.WriteValidationError(w, err)
		return
	}

	// Get current submission to check status
	submission, err := h.db.GetSubmission(r.Context(), submissionID)
	if err != nil {
		httputil.HandleError(w, err)
		return
	}

	// Use existing values if not provided
	status := req.Status
	if status == "" {
		status = submission.Status
	}

	feedback := req.Feedback
	if feedback == "" {
		feedback = submission.Feedback
	}

	if err := h.db.UpdateSubmission(r.Context(), submissionID, status, feedback); err != nil {
		httputil.HandleError(w, err)
		return
	}

	// Track status changes
	if req.Status != "" {
		metrics.ProjectStatusUpdates.WithLabelValues(status).Inc()
	}

	httputil.WriteSuccess(w, "Submission updated successfully", nil)
}

func (h *ProjectsHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	project, err := h.db.GetProject(r.Context(), id)
	if err != nil {
		httputil.HandleError(w, err)
		return
	}

	httputil.WriteSuccess(w, "Project retrieved successfully", project)
}

func (h *ProjectsHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req models.Project
	if err := validator.ValidateRequest(r, &req); err != nil {
		validator.WriteValidationError(w, err)
		return
	}

	if err := h.db.UpdateProject(r.Context(), id, &req); err != nil {
		httputil.HandleError(w, err)
		return
	}

	httputil.WriteSuccess(w, "Project updated successfully", nil)
}

func (h *ProjectsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.db.DeleteProject(r.Context(), id); err != nil {
		httputil.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
