package handlers

import (
	"encoding/json"
	"net/http"
	"shadow-nova/backend/internal/database"
	"shadow-nova/backend/internal/models"
	"strconv"
)

type AdminHandler struct {
	db database.Service
}

func NewAdminHandler(db database.Service) *AdminHandler {
	return &AdminHandler{db: db}
}

type UpdateCollectorFrequencyRequest struct {
	RunsPerDay int `json:"runs_per_day"`
}

func (h *AdminHandler) UpdateCollectorFrequency(w http.ResponseWriter, r *http.Request) {
	// TODO: Verify user is admin (Middleware should handle this, but we can double check)
	
	var req UpdateCollectorFrequencyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.RunsPerDay < 1 || req.RunsPerDay > 24 {
		http.Error(w, "Runs per day must be between 1 and 24", http.StatusBadRequest)
		return
	}

	value := strconv.Itoa(req.RunsPerDay)
	if err := h.db.UpdateSystemSetting(r.Context(), "collector_runs_per_day", value); err != nil {
		http.Error(w, "Failed to update setting", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.SuccessResponse{
		Message: "Collector frequency updated successfully",
		Data:    map[string]int{"runs_per_day": req.RunsPerDay},
	})
}
