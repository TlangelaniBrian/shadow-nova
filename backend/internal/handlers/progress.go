package handlers

import (
	"net/http"
	"shadow-nova/backend/internal/database"
	"shadow-nova/backend/internal/httputil"
	"shadow-nova/backend/internal/metrics"
	"shadow-nova/backend/internal/middleware"
	"shadow-nova/backend/internal/models"
	"shadow-nova/backend/internal/validator"

	"github.com/go-chi/chi/v5"
)

type ProgressHandler struct {
	db database.Service
}

func NewProgressHandler(db database.Service) *ProgressHandler {
	return &ProgressHandler{db: db}
}

func (h *ProgressHandler) UpdateProgress(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		httputil.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req models.UpdateProgressRequest
	if err := validator.ValidateRequest(r, &req); err != nil {
		validator.WriteValidationError(w, err)
		return
	}

	if err := h.db.UpdateUserProgress(r.Context(), userID, req); err != nil {
		httputil.HandleError(w, err)
		return
	}

	// Track lesson completions
	if req.Completed {
		lesson, err := h.db.GetLesson(r.Context(), req.LessonID)
		if err == nil && lesson != nil {
			metrics.LessonCompletions.WithLabelValues(lesson.ContentType).Inc()
		}
	}

	httputil.WriteSuccess(w, "Progress updated successfully", nil)
}

func (h *ProgressHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		httputil.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	stats, err := h.db.GetUserStats(r.Context(), userID)
	if err != nil {
		httputil.HandleError(w, err)
		return
	}

	httputil.WriteSuccess(w, "User stats retrieved successfully", stats)
}

func (h *ProgressHandler) GetPathProgress(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		httputil.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}
	pathID := chi.URLParam(r, "id")

	// Return array of user progress records for each lesson in the path
	progress, err := h.db.GetUserProgressForPath(r.Context(), userID, pathID)
	if err != nil {
		httputil.HandleError(w, err)
		return
	}

	httputil.WriteSuccess(w, "Path progress retrieved successfully", progress)
}
