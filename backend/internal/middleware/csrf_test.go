package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/stretchr/testify/assert"
)

func TestCSRFMiddleware(t *testing.T) {
	// Set up test CSRF key (32 bytes)
	originalKey := os.Getenv("CSRF_KEY")
	testKey := "12345678901234567890123456789012" // 32 bytes
	os.Setenv("CSRF_KEY", testKey)
	os.Setenv("ENV", "development") // Disable secure flag for testing
	defer func() {
		os.Setenv("CSRF_KEY", originalKey)
		os.Unsetenv("ENV")
	}()

	t.Run("POST request without CSRF token returns 403", func(t *testing.T) {
		r := chi.NewRouter()
		r.Use(CSRF())

		r.Post("/test", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
		})

		req := httptest.NewRequest("POST", "/test", strings.NewReader("{}"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), "CSRF")
	})

	t.Run("POST request with valid CSRF token succeeds", func(t *testing.T) {
		r := chi.NewRouter()
		r.Use(CSRF())

		r.Get("/get-token", func(w http.ResponseWriter, r *http.Request) {
			token := csrf.Token(r)
			w.Header().Set("X-CSRF-Token", token)
			w.WriteHeader(http.StatusOK)
		})

		r.Post("/test", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
		})

		// First, get the CSRF token
		getReq := httptest.NewRequest("GET", "/get-token", nil)
		getW := httptest.NewRecorder()
		r.ServeHTTP(getW, getReq)

		assert.Equal(t, http.StatusOK, getW.Code)
		token := getW.Header().Get("X-CSRF-Token")
		assert.NotEmpty(t, token)

		// Extract cookie from response
		cookies := getW.Result().Cookies()
		var csrfCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "_gorilla_csrf" {
				csrfCookie = cookie
				break
			}
		}
		assert.NotNil(t, csrfCookie)

		// Now make a POST request with the token
		postReq := httptest.NewRequest("POST", "/test", strings.NewReader("{}"))
		postReq.Header.Set("Content-Type", "application/json")
		postReq.Header.Set("X-CSRF-Token", token)
		postReq.AddCookie(csrfCookie)
		postW := httptest.NewRecorder()

		r.ServeHTTP(postW, postReq)

		assert.Equal(t, http.StatusOK, postW.Code)
		assert.Equal(t, "success", postW.Body.String())
	})

	t.Run("GET requests work without CSRF token", func(t *testing.T) {
		r := chi.NewRouter()
		r.Use(CSRF())

		r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
		})

		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "success", w.Body.String())
	})

	t.Run("PUT request without CSRF token returns 403", func(t *testing.T) {
		r := chi.NewRouter()
		r.Use(CSRF())

		r.Put("/test", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
		})

		req := httptest.NewRequest("PUT", "/test", strings.NewReader("{}"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("DELETE request without CSRF token returns 403", func(t *testing.T) {
		r := chi.NewRouter()
		r.Use(CSRF())

		r.Delete("/test", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
		})

		req := httptest.NewRequest("DELETE", "/test", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

func TestCSRFMiddlewarePanicsOnInvalidKey(t *testing.T) {
	originalKey := os.Getenv("CSRF_KEY")
	defer os.Setenv("CSRF_KEY", originalKey)

	t.Run("panics when CSRF_KEY is not 32 bytes", func(t *testing.T) {
		os.Setenv("CSRF_KEY", "short")

		assert.Panics(t, func() {
			CSRF()
		})
	})

	t.Run("panics when CSRF_KEY is empty", func(t *testing.T) {
		os.Setenv("CSRF_KEY", "")

		assert.Panics(t, func() {
			CSRF()
		})
	})
}
