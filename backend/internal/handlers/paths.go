package handlers

import (
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
	paths, err := h.db.GetLearningPaths(r.Context())
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "Failed to fetch learning paths")
		return
	}

	httputil.WriteSuccess(w, "Learning paths retrieved successfully", paths)
}

func (h *PathsHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	path, err := h.db.GetLearningPath(r.Context(), id)
	if err != nil {
		httputil.WriteError(w, http.StatusNotFound, "Learning path not found")
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
