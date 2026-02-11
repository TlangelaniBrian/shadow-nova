package handlers

import (
	"net/http"
	"shadow-nova/backend/internal/database"
	"shadow-nova/backend/internal/httputil"
	"shadow-nova/backend/internal/validator"
	"strconv"
)

type AdminHandler struct {
	db database.Service
}

func NewAdminHandler(db database.Service) *AdminHandler {
	return &AdminHandler{db: db}
}

type UpdateCollectorFrequencyRequest struct {
	RunsPerDay int `json:"runs_per_day" validate:"required,min=1,max=24"`
}

func (h *AdminHandler) UpdateCollectorFrequency(w http.ResponseWriter, r *http.Request) {
	var req UpdateCollectorFrequencyRequest
	if err := validator.ValidateRequest(r, &req); err != nil {
		validator.WriteValidationError(w, err)
		return
	}

	value := strconv.Itoa(req.RunsPerDay)
	if err := h.db.UpdateSystemSetting(r.Context(), "collector_runs_per_day", value); err != nil {
		httputil.HandleError(w, err)
		return
	}

	httputil.WriteSuccess(w, "Collector frequency updated successfully", map[string]int{"runs_per_day": req.RunsPerDay})
}
