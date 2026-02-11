package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"shadow-nova/backend/internal/ai"
	"shadow-nova/backend/internal/auth"
	"shadow-nova/backend/internal/collector"
	"shadow-nova/backend/internal/handlers"
	"shadow-nova/backend/internal/logging"
	"shadow-nova/backend/internal/middleware"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (s *Server) RegisterRoutes() (http.Handler, error) {
	r := chi.NewRouter()

	// Request ID middleware (first, so all logs have request_id)
	r.Use(middleware.RequestID)

	// Structured logging middleware
	r.Use(middleware.RequestLogger)

	// Security middleware
	r.Use(middleware.SecurityHeaders)
	r.Use(middleware.CORSMiddleware())
	rateLimiter := middleware.NewRateLimiter(100) // 100 requests per minute
	r.Use(rateLimiter.Limit)

	// Standard middleware
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.PrometheusMiddleware)

	r.Get("/", s.HelloWorldHandler)
	r.Get("/health", s.healthHandler)
	r.Handle("/metrics", promhttp.Handler())

	// Version info endpoint (unversioned for discoverability)
	r.Get("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"version":"1.0.0","api_version":"v1"}`))
	})

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

	r.Route("/api/v1", func(r chi.Router) {
		// Apply CSRF protection with exemptions for auth endpoints
		r.Use(csrfMiddleware)
		
		// Initialize Google Auth service
		googleAuth, err := auth.NewGoogleAuthService()
		if err != nil {
			logging.Warn("failed to initialize google auth", "error", err)
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
				logging.Info("collector goroutine stopped before initial run")
				return
			}

			// Initial run
			logging.Info("starting initial content collection")
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
				logging.Info("scheduled next collection", "interval", interval, "runs_per_day", runsPerDay)

				// Wait for interval or shutdown signal
				select {
				case <-time.After(interval):
					logging.Info("running scheduled content collection")
					collectorService.CollectAll(s.collectorCtx)
					collectorService.ProcessUnprocessedItems(s.collectorCtx)
				case <-s.collectorCtx.Done():
					logging.Info("collector goroutine stopped gracefully")
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
						logging.Error("failed to clean expired blacklisted tokens", err)
					} else {
						logging.Info("cleaned expired blacklisted tokens", "count", deleted)
					}
				case <-s.collectorCtx.Done():
					logging.Info("token cleanup goroutine stopped gracefully")
					return
				}
			}
		}()

		// Start idempotency key cleanup (runs daily)
		go func() {
			ticker := time.NewTicker(24 * time.Hour)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					ctx := context.Background()
					deleted, err := s.db.DeleteExpiredIdempotencyKeys(ctx)
					if err != nil {
						logging.Error("failed to clean expired idempotency keys", err)
					} else {
						logging.Info("cleaned expired idempotency keys", "count", deleted)
					}
				case <-s.collectorCtx.Done():
					logging.Info("idempotency cleanup goroutine stopped gracefully")
					return
				}
			}
		}()

		adminHandler := handlers.NewAdminHandler(s.db)
		pathsHandler := handlers.NewPathsHandler(s.db)
		progressHandler := handlers.NewProgressHandler(s.db)
		projectsHandler := handlers.NewProjectsHandler(s.db)
		userHandler := handlers.NewUserHandler(s.db)

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

			// Public project listing and viewing
			r.Get("/projects", projectsHandler.List)
			r.Get("/projects/{id}", projectsHandler.Get)
		})

		// Protected routes (auth required, CSRF already applied above)
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.VerifyToken)
			r.Use(middleware.Idempotency(s.db))

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
			r.Delete("/auth/github/disconnect", githubHandler.Disconnect)

			// User profile management
			r.Get("/user/profile", userHandler.GetProfile)
			r.Patch("/user/profile", userHandler.UpdateProfile)
			r.Put("/user/password", userHandler.UpdatePassword)

			// Admin Routes (require admin role)
			r.Group(func(r chi.Router) {
				r.Use(middleware.AdminOnly)

				// User management
				adminUsersHandler := handlers.NewAdminUsersHandler(s.db)
				r.Get("/admin/users", adminUsersHandler.ListUsers)
				r.Get("/admin/users/{id}", adminUsersHandler.GetUser)
				r.Post("/admin/users", adminUsersHandler.CreateUser)
				r.Put("/admin/users/{id}", adminUsersHandler.UpdateUser)
				r.Delete("/admin/users/{id}", adminUsersHandler.DeleteUser)

				// Learning Paths CRUD (admin only)
				r.Put("/paths/{id}", pathsHandler.Update)
				r.Patch("/paths/{id}", pathsHandler.Update)
				r.Delete("/paths/{id}", pathsHandler.Delete)

				// Modules CRUD (admin only)
				r.Put("/modules/{id}", pathsHandler.UpdateModule)
				r.Patch("/modules/{id}", pathsHandler.UpdateModule)
				r.Delete("/modules/{id}", pathsHandler.DeleteModule)

				// Lessons CRUD (admin only)
				r.Put("/lessons/{id}", pathsHandler.UpdateLesson)
				r.Patch("/lessons/{id}", pathsHandler.UpdateLesson)
				r.Delete("/lessons/{id}", pathsHandler.DeleteLesson)

				// Projects CRUD (admin only)
				r.Post("/projects", projectsHandler.Create)
				r.Put("/projects/{id}", projectsHandler.Update)
				r.Patch("/projects/{id}", projectsHandler.Update)
				r.Delete("/projects/{id}", projectsHandler.Delete)

				// System settings
				r.Post("/admin/settings/collector", adminHandler.UpdateCollectorFrequency)
			})
		})
	})

	return r, nil
}


