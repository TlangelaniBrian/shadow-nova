package handlers

import (
	"encoding/json"
	"net/http"
	"shadow-nova/backend/internal/database"
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
	userID := r.Context().Value("user_id").(int) // Assuming auth middleware sets this

	var req models.UpdateProgressRequest
	if err := validator.ValidateRequest(r, &req); err != nil {
		validator.WriteValidationError(w, err)
		return
	}

	if err := h.db.UpdateUserProgress(r.Context(), userID, req); err != nil {
		http.Error(w, "Failed to update progress", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.SuccessResponse{
		Message: "Progress updated successfully",
	})
}

func (h *ProgressHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	stats, err := h.db.GetUserStats(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to fetch user stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.SuccessResponse{
		Message: "User stats retrieved successfully",
		Data:    stats,
	})
}

func (h *ProgressHandler) GetPathProgress(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	pathID := chi.URLParam(r, "id")

	progress, err := h.db.GetPathProgress(r.Context(), userID, pathID)
	if err != nil {
		http.Error(w, "Failed to fetch path progress", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.SuccessResponse{
		Message: "Path progress retrieved successfully",
		Data:    progress,
	})
}
