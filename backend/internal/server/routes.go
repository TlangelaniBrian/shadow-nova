package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"shadow-nova/backend/internal/auth"
	"shadow-nova/backend/internal/handlers"
	"shadow-nova/backend/internal/middleware"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	
	// Security middleware
	r.Use(middleware.SecurityHeaders)
	r.Use(middleware.CORSMiddleware())
	rateLimiter := middleware.NewRateLimiter(100) // 100 requests per minute
	r.Use(rateLimiter.Limit)
	
	// Standard middleware
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.PrometheusMiddleware)

	r.Get("/", s.HelloWorldHandler)
	r.Get("/health", s.healthHandler)
	r.Handle("/metrics", promhttp.Handler())

	r.Route("/api", func(r chi.Router) {
		// Initialize auth middleware with a secret (should be from env in real app)
		authMiddleware := middleware.NewAuthMiddleware("my-secret-key")
		
		// Initialize Google Auth service
		googleAuth, err := auth.NewGoogleAuthService()
		if err != nil {
			// Log error but don't fail - allow app to start
			fmt.Printf("Warning: Failed to initialize Google Auth: %v\n", err)
		}
		authHandler := handlers.NewAuthHandler(googleAuth, s.db)
		
		// Public routes (no auth required)
		r.Group(func(r chi.Router) {
			// Google OAuth endpoints
			r.Get("/auth/google", authHandler.GoogleLogin)
			r.Get("/auth/google/callback", authHandler.GoogleCallback)
			r.Post("/auth/google/verify", authHandler.VerifyGoogleToken)
			
			// Traditional auth (optional)
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
		})
		
		// Protected routes (auth required)
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.VerifyToken)
			
			r.Post("/progress", authHandler.UpdateProgress)
			r.Get("/paths", s.getPathsHandler)
			r.Get("/projects", s.getProjectsHandler)
		})
	})

	return r
}

func (s *Server) getPathsHandler(w http.ResponseWriter, r *http.Request) {
	// Mock data for now
	paths := []map[string]string{
		{"id": "frontend-vue", "title": "Frontend Mastery"},
		{"id": "backend-go", "title": "Backend Engineering"},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paths)
}

func (s *Server) getProjectsHandler(w http.ResponseWriter, r *http.Request) {
	// Mock data for now
	projects := []map[string]string{
		{"id": "portfolio", "title": "Personal Portfolio"},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}
