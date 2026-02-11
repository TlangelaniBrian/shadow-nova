package httputil

import (
	"errors"
	"net/http"
	apperrors "shadow-nova/backend/internal/errors"
	"shadow-nova/backend/internal/logging"
)

// HandleError is a centralized error handler that maps application errors to HTTP responses
func HandleError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	// Check if it's an AppError with explicit HTTP code
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		WriteError(w, appErr.HTTPCode, appErr.Message)
		if appErr.HTTPCode >= 500 {
			logging.Error("Internal error", err)
		}
		return
	}

	// Check for specific error types using sentinel errors
	if apperrors.IsNotFound(err) {
		WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	if apperrors.IsUnauthorized(err) {
		WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if apperrors.IsForbidden(err) {
		WriteError(w, http.StatusForbidden, err.Error())
		return
	}

	if apperrors.IsInvalidInput(err) {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if apperrors.IsDuplicateEntry(err) {
		WriteError(w, http.StatusConflict, err.Error())
		return
	}

	// Default to 500 for unhandled errors
	logging.Error("Unhandled error", err)
	WriteError(w, http.StatusInternalServerError, "Internal server error")
}
