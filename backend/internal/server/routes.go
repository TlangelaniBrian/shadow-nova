package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"shadow-nova/backend/internal/ai"
	"shadow-nova/backend/internal/auth"
	"shadow-nova/backend/internal/collector"
	"shadow-nova/backend/internal/handlers"
	"shadow-nova/backend/internal/middleware"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (s *Server) RegisterRoutes() (http.Handler, error) {
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

	// Validate required configuration before setting up routes
	csrfMiddleware, err := middleware.CSRF()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize CSRF: %w", err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}
	authMiddleware := middleware.NewAuthMiddleware(jwtSecret, s.db)

	r.Route("/api", func(r chi.Router) {
		// Apply CSRF protection with exemptions for auth endpoints
		r.Use(csrfMiddleware)
		
		// Initialize Google Auth service
		googleAuth, err := auth.NewGoogleAuthService()
		if err != nil {
			fmt.Printf("Warning: Failed to initialize Google Auth: %v\n", err)
		}
		authHandler := handlers.NewAuthHandler(googleAuth, s.db)
		
		// Initialize GitHub Auth service
		githubAuth := auth.NewGitHubAuthService()
		githubHandler := handlers.NewGitHubHandler(githubAuth, s.db)
		
		// Initialize AI & Collector
		aiService := ai.NewAIService()
		collectorService := collector.New(s.db, aiService)
		
		// Start background collector (simple goroutine for now)
		go func() {
			// Wait for server to start
			select {
			case <-time.After(5 * time.Second):
				// Continue
			case <-s.collectorCtx.Done():
				log.Println("Collector goroutine stopped before initial run")
				return
			}

			// Initial run
			log.Println("Running initial content collection...")
			collectorService.CollectAll(s.collectorCtx)
			collectorService.ProcessUnprocessedItems(s.collectorCtx)

			// Dynamic periodic run
			for {
				// Get frequency from DB
				runsPerDayStr, err := s.db.GetSystemSetting(s.collectorCtx, "collector_runs_per_day")
				runsPerDay := 1 // Default
				if err == nil && runsPerDayStr != "" {
					fmt.Sscanf(runsPerDayStr, "%d", &runsPerDay)
				}
				if runsPerDay < 1 { runsPerDay = 1 }

				// Calculate interval
				interval := 24 * time.Hour / time.Duration(runsPerDay)
				log.Printf("Next collection in %v (Runs per day: %d)", interval, runsPerDay)

				// Wait for interval or shutdown signal
				select {
				case <-time.After(interval):
					log.Println("Running scheduled content collection...")
					collectorService.CollectAll(s.collectorCtx)
					collectorService.ProcessUnprocessedItems(s.collectorCtx)
				case <-s.collectorCtx.Done():
					log.Println("Collector goroutine stopped gracefully")
					return
				}
			}
		}()

		// Start token blacklist cleanup (runs daily)
		go func() {
			ticker := time.NewTicker(24 * time.Hour)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					ctx := context.Background()
					deleted, err := s.db.DeleteExpiredBlacklistedTokens(ctx)
					if err != nil {
						log.Printf("Failed to clean expired blacklisted tokens: %v", err)
					} else {
						log.Printf("Cleaned %d expired blacklisted tokens", deleted)
					}
				case <-s.collectorCtx.Done():
					log.Println("Token cleanup goroutine stopped gracefully")
					return
				}
			}
		}()

		adminHandler := handlers.NewAdminHandler(s.db)
		pathsHandler := handlers.NewPathsHandler(s.db)
		progressHandler := handlers.NewProgressHandler(s.db)
		projectsHandler := handlers.NewProjectsHandler(s.db)
		
		// GitHub handler already initialized on line 52, just removing duplicate
		
		// CSRF token endpoint (GET is safe, exempt from CSRF validation)
		r.Get("/csrf-token", authHandler.GetCSRFToken)

		// Public routes (no auth required)
		r.Group(func(r chi.Router) {
			// Google OAuth endpoints
			r.Get("/auth/google", authHandler.GoogleLogin)
			r.Get("/auth/google/callback", authHandler.GoogleCallback)
			r.Post("/auth/google/verify", authHandler.VerifyGoogleToken)

			// GitHub OAuth endpoints
			r.Get("/auth/github/login", githubHandler.Login)
			r.Get("/auth/github/callback", githubHandler.Callback)

			// Traditional auth (optional)
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)

			// Public project listing
			r.Get("/projects", projectsHandler.List)
		})

		// Protected routes (auth required, CSRF already applied above)
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.VerifyToken)

			// Logout endpoint (protected to ensure token blacklisting)
			r.Post("/auth/logout", authHandler.Logout)

			// Progress & Stats
			r.Post("/progress", progressHandler.UpdateProgress)
			r.Get("/stats", progressHandler.GetStats)

			// Path progress with ownership validation
			r.With(middleware.ValidatePathOwnership(s.db)).Get("/paths/{id}/progress", progressHandler.GetPathProgress)
			
			// Learning Paths Routes
			r.Get("/paths", pathsHandler.List)
			r.Get("/paths/{id}", pathsHandler.Get)
			r.Post("/paths", pathsHandler.Create)
			r.Post("/paths/{id}/modules", pathsHandler.AddModule)
			r.Post("/lessons", pathsHandler.AddLesson) // keeping it simple, expects module_id in body
			
			// Projects Routes (submit is protected, create is admin only)
			r.Post("/submissions", projectsHandler.Submit)

			// Submission endpoints with ownership validation
			r.With(middleware.ValidateSubmissionOwnership(s.db)).Get("/submissions/{id}", projectsHandler.GetSubmission)
			r.With(middleware.ValidateSubmissionOwnership(s.db)).Patch("/submissions/{id}", projectsHandler.UpdateSubmission)
			
			// GitHub Connect (Protected)
			r.Get("/auth/github/connect", githubHandler.Connect)
			
			// Admin Routes (require admin role)
			r.Group(func(r chi.Router) {
				r.Use(middleware.AdminOnly)
				r.Post("/projects", projectsHandler.Create)
				r.Post("/admin/settings/collector", adminHandler.UpdateCollectorFrequency)
			})
		})
	})

	return r, nil
}


