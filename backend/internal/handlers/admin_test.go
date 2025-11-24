package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"shadow-nova/backend/internal/database"
	"testing"
)

func TestAdminHandler_UpdateCollectorFrequency(t *testing.T) {
	// Setup
	mockDB := &database.MockService{
		UpdateSystemSettingFunc: func(ctx context.Context, key, value string) error {
			if key != "collector_runs_per_day" {
				t.Errorf("expected key 'collector_runs_per_day', got '%s'", key)
			}
			if value != "4" {
				t.Errorf("expected value '4', got '%s'", value)
			}
			return nil
		},
	}
	handler := NewAdminHandler(mockDB)

	payload := map[string]int{
		"runs_per_day": 4,
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/admin/settings/collector", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// Execute
	handler.UpdateCollectorFrequency(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestAdminHandler_UpdateCollectorFrequency_Invalid(t *testing.T) {
	mockDB := &database.MockService{}
	handler := NewAdminHandler(mockDB)

	// Test with invalid value (too high)
	payload := map[string]int{
		"runs_per_day": 25,
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/admin/settings/collector", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler.UpdateCollectorFrequency(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("expected BadRequest for value 25, got %v", status)
	}
}
