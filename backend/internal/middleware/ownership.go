package middleware

import (
	"net/http"
	"strconv"

	"shadow-nova/backend/internal/database"
	"shadow-nova/backend/internal/httputil"

	"github.com/go-chi/chi/v5"
)

// ValidatePathOwnership validates that the user has access to the learning path.
// For now, all authenticated users can access all paths.
// In the future, this can be extended to check enrollment or purchase status.
func ValidatePathOwnership(db database.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := GetUserID(r)
			if !ok {
				httputil.WriteError(w, http.StatusUnauthorized, "User not authenticated")
				return
			}

			pathID := chi.URLParam(r, "id")
			if pathID == "" {
				httputil.WriteError(w, http.StatusBadRequest, "Path ID is required")
				return
			}

			hasAccess, err := db.UserHasAccessToPath(r.Context(), userID, pathID)
			if err != nil {
				httputil.WriteError(w, http.StatusInternalServerError, "Failed to verify path access")
				return
			}

			if !hasAccess {
				httputil.WriteError(w, http.StatusForbidden, "You do not have access to this path")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// ValidateSubmissionOwnership validates that the user owns the project submission.
// This prevents users from viewing or modifying other users' submissions.
func ValidateSubmissionOwnership(db database.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := GetUserID(r)
			if !ok {
				httputil.WriteError(w, http.StatusUnauthorized, "User not authenticated")
				return
			}

			submissionIDStr := chi.URLParam(r, "id")
			if submissionIDStr == "" {
				httputil.WriteError(w, http.StatusBadRequest, "Submission ID is required")
				return
			}

			submissionID, err := strconv.Atoi(submissionIDStr)
			if err != nil {
				httputil.WriteError(w, http.StatusBadRequest, "Invalid submission ID")
				return
			}

			ownsSubmission, err := db.UserOwnsSubmission(r.Context(), userID, submissionID)
			if err != nil {
				httputil.WriteError(w, http.StatusInternalServerError, "Failed to verify submission ownership")
				return
			}

			if !ownsSubmission {
				httputil.WriteError(w, http.StatusForbidden, "You do not have access to this submission")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// ValidateProgressOwnership validates that the user owns the progress record.
// This prevents users from viewing or modifying other users' progress.
func ValidateProgressOwnership(db database.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := GetUserID(r)
			if !ok {
				httputil.WriteError(w, http.StatusUnauthorized, "User not authenticated")
				return
			}

			progressIDStr := chi.URLParam(r, "id")
			if progressIDStr == "" {
				httputil.WriteError(w, http.StatusBadRequest, "Progress ID is required")
				return
			}

			progressID, err := strconv.Atoi(progressIDStr)
			if err != nil {
				httputil.WriteError(w, http.StatusBadRequest, "Invalid progress ID")
				return
			}

			ownsProgress, err := db.UserOwnsProgress(r.Context(), userID, progressID)
			if err != nil {
				httputil.WriteError(w, http.StatusInternalServerError, "Failed to verify progress ownership")
				return
			}

			if !ownsProgress {
				httputil.WriteError(w, http.StatusForbidden, "You do not have access to this progress record")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
