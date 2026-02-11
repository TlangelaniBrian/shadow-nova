package handlers

import (
	"fmt"
	"net/http"
	"shadow-nova/backend/internal/database"
	"shadow-nova/backend/internal/httputil"
	"shadow-nova/backend/internal/models"
	"shadow-nova/backend/internal/validator"

	"github.com/go-chi/chi/v5"
)

type PathsHandler struct {
	db database.Service
}

func NewPathsHandler(db database.Service) *PathsHandler {
	return &PathsHandler{db: db}
}

func (h *PathsHandler) List(w http.ResponseWriter, r *http.Request) {
	pagination := models.ParsePagination(r)
	offset := (pagination.Page - 1) * pagination.Limit

	paths, err := h.db.GetLearningPaths(r.Context(), pagination.Limit, offset)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "Failed to fetch learning paths")
		return
	}

	total, err := h.db.GetLearningPathsCount(r.Context())
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "Failed to count learning paths")
		return
	}

	response := models.NewPaginatedResponse(paths, pagination.Page, pagination.Limit, total)
	httputil.WriteJSON(w, http.StatusOK, response)
}

func (h *PathsHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	path, err := h.db.GetLearningPath(r.Context(), id)
	if err != nil {
		httputil.HandleError(w, err)
		return
	}

	httputil.WriteSuccess(w, "Learning path retrieved successfully", path)
}

func (h *PathsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePathRequest
	if err := validator.ValidateRequest(r, &req); err != nil {
		validator.WriteValidationError(w, err)
		return
	}

	path := &models.LearningPath{
		ID:          req.ID,
		Title:       req.Title,
		Description: req.Description,
		Difficulty:  req.Difficulty,
	}

	if err := h.db.CreateLearningPath(r.Context(), path); err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "Failed to create learning path")
		return
	}

	httputil.WriteCreated(w, "Learning path created successfully", path)
}

func (h *PathsHandler) AddModule(w http.ResponseWriter, r *http.Request) {
	pathID := chi.URLParam(r, "id")
	var req models.CreateModuleRequest
	if err := validator.ValidateRequest(r, &req); err != nil {
		validator.WriteValidationError(w, err)
		return
	}

	module := &models.Module{
		PathID:      pathID,
		Title:       req.Title,
		Description: req.Description,
		OrderIndex:  req.OrderIndex,
	}

	if err := h.db.CreateModule(r.Context(), module); err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "Failed to create module")
		return
	}

	httputil.WriteCreated(w, "Module added successfully", module)
}

func (h *PathsHandler) AddLesson(w http.ResponseWriter, r *http.Request) {
	var req models.Lesson
	if err := validator.ValidateRequest(r, &req); err != nil {
		validator.WriteValidationError(w, err)
		return
	}

	if req.Title == "" || req.ContentType == "" {
		httputil.WriteError(w, http.StatusBadRequest, "Title and ContentType are required")
		return
	}

	if err := h.db.CreateLesson(r.Context(), &req); err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "Failed to create lesson")
		return
	}

	httputil.WriteCreated(w, "Lesson added successfully", req)
}

func (h *PathsHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req models.LearningPath
	if err := validator.ValidateRequest(r, &req); err != nil {
		validator.WriteValidationError(w, err)
		return
	}

	if err := h.db.UpdateLearningPath(r.Context(), id, &req); err != nil {
		httputil.HandleError(w, err)
		return
	}

	httputil.WriteSuccess(w, "Learning path updated successfully", nil)
}

func (h *PathsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.db.DeleteLearningPath(r.Context(), id); err != nil {
		httputil.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *PathsHandler) UpdateModule(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := parseIntID(idStr)
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "Invalid module ID")
		return
	}

	var req models.Module
	if err := validator.ValidateRequest(r, &req); err != nil {
		validator.WriteValidationError(w, err)
		return
	}

	if err := h.db.UpdateModule(r.Context(), id, &req); err != nil {
		httputil.HandleError(w, err)
		return
	}

	httputil.WriteSuccess(w, "Module updated successfully", nil)
}

func (h *PathsHandler) DeleteModule(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := parseIntID(idStr)
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "Invalid module ID")
		return
	}

	if err := h.db.DeleteModule(r.Context(), id); err != nil {
		httputil.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *PathsHandler) UpdateLesson(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := parseIntID(idStr)
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "Invalid lesson ID")
		return
	}

	var req models.Lesson
	if err := validator.ValidateRequest(r, &req); err != nil {
		validator.WriteValidationError(w, err)
		return
	}

	if err := h.db.UpdateLesson(r.Context(), id, &req); err != nil {
		httputil.HandleError(w, err)
		return
	}

	httputil.WriteSuccess(w, "Lesson updated successfully", nil)
}

func (h *PathsHandler) DeleteLesson(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := parseIntID(idStr)
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "Invalid lesson ID")
		return
	}

	if err := h.db.DeleteLesson(r.Context(), id); err != nil {
		httputil.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// parseIntID is a helper to parse integer IDs from URL parameters
func parseIntID(idStr string) (int, error) {
	var id int
	_, err := fmt.Sscanf(idStr, "%d", &id)
	return id, err
}
