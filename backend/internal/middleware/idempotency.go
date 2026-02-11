package middleware

import (
	"bytes"
	"net/http"
	"time"

	"shadow-nova/backend/internal/database"
)

// Idempotency middleware prevents duplicate operations by caching responses
// for requests that include an Idempotency-Key header.
func Idempotency(db database.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Only apply to mutating methods
			if r.Method != "POST" && r.Method != "PUT" && r.Method != "PATCH" {
				next.ServeHTTP(w, r)
				return
			}

			// Check for idempotency key
			idempotencyKey := r.Header.Get("Idempotency-Key")
			if idempotencyKey == "" {
				// No key provided, continue without idempotency
				next.ServeHTTP(w, r)
				return
			}

			// Get user ID from context (set by auth middleware)
			userID, ok := r.Context().Value(UserIDKey).(int)
			if !ok {
				// If user ID is not in context, skip idempotency
				// This shouldn't happen on protected routes, but we handle it gracefully
				next.ServeHTTP(w, r)
				return
			}

			// Check if we've seen this key before
			cached, err := db.GetIdempotentResponse(r.Context(), idempotencyKey, userID)
			if err == nil && cached != nil {
				// Return cached response
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("X-Idempotent-Replay", "true")
				w.WriteHeader(cached.StatusCode)
				w.Write([]byte(cached.Body))
				return
			}

			// Capture response
			recorder := &idempotencyResponseRecorder{
				ResponseWriter: w,
				statusCode:     http.StatusOK, // Default to 200 if WriteHeader is not called
				body:           &bytes.Buffer{},
			}

			next.ServeHTTP(recorder, r)

			// Store response for future requests (24 hour TTL)
			expiresAt := time.Now().Add(24 * time.Hour)
			db.StoreIdempotentResponse(
				r.Context(),
				idempotencyKey,
				userID,
				r.URL.Path,
				r.Method,
				recorder.statusCode,
				recorder.body.String(),
				expiresAt,
			)
		})
	}
}

// idempotencyResponseRecorder captures both the status code and response body
type idempotencyResponseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (r *idempotencyResponseRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *idempotencyResponseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
