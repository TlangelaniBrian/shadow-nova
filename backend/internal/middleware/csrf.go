package middleware

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/csrf"
)

func CSRF() (func(http.Handler) http.Handler, error) {
	// Get CSRF key from environment (32 bytes, hex-encoded)
	csrfKeyHex := os.Getenv("CSRF_KEY")
	if csrfKeyHex == "" {
		return nil, fmt.Errorf("CSRF_KEY environment variable not set")
	}

	// Decode from hex
	csrfKey, err := hex.DecodeString(csrfKeyHex)
	if err != nil {
		return nil, fmt.Errorf("CSRF_KEY must be hex-encoded: %w", err)
	}

	if len(csrfKey) != 32 {
		return nil, fmt.Errorf("CSRF_KEY must be 32 bytes, got %d", len(csrfKey))
	}

	secure := os.Getenv("ENV") == "production"

	// Paths exempt from CSRF (authentication endpoints)
	exemptPaths := map[string]bool{
		"/api/login":                     true,
		"/api/register":                  true,
		"/api/auth/google":               true,
		"/api/auth/google/callback":      true,
		"/api/auth/google/verify":        true,
		"/api/auth/github/login":         true,
		"/api/auth/github/callback":      true,
		"/api/csrf-token":                true,
	}

	csrfProtect := csrf.Protect(
		csrfKey,
		csrf.Secure(secure),
		csrf.Path("/"),
		csrf.SameSite(csrf.SameSiteStrictMode),
		csrf.ErrorHandler(http.HandlerFunc(csrfErrorHandler)),
	)

	// Wrapper that exempts certain paths
	wrapper := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if exemptPaths[r.URL.Path] {
				// Skip CSRF for exempt paths
				next.ServeHTTP(w, r)
			} else {
				// Apply CSRF protection
				csrfProtect(next).ServeHTTP(w, r)
			}
		})
	}

	return wrapper, nil
}

func csrfErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte(`{"error": "CSRF token validation failed", "status": 403}`))
}
