package httputil

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"shadow-nova/backend/internal/models"
)

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		log.Printf("httputil: failed to encode JSON response: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(buf.Bytes())
}

func WriteSuccess(w http.ResponseWriter, message string, data interface{}) {
	WriteJSON(w, http.StatusOK, models.SuccessResponse{Message: message, Data: data})
}

func WriteCreated(w http.ResponseWriter, message string, data interface{}) {
	WriteJSON(w, http.StatusCreated, models.SuccessResponse{Message: message, Data: data})
}

func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, models.ErrorResponse{Error: message, Status: status})
}
