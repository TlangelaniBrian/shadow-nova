package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
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
			time.Sleep(5 * time.Second)
			ctx := context.Background()
			
			// Initial run
			log.Println("Running initial content collection...")
			collectorService.CollectAll(ctx)
			collectorService.ProcessUnprocessedItems(ctx)
			
			// Dynamic periodic run
			for {
				// Get frequency from DB
				runsPerDayStr, err := s.db.GetSystemSetting(ctx, "collector_runs_per_day")
				runsPerDay := 1 // Default
				if err == nil && runsPerDayStr != "" {
					fmt.Sscanf(runsPerDayStr, "%d", &runsPerDay)
				}
				if runsPerDay < 1 { runsPerDay = 1 }
				
				// Calculate interval
				interval := 24 * time.Hour / time.Duration(runsPerDay)
				log.Printf("Next collection in %v (Runs per day: %d)", interval, runsPerDay)
				
				time.Sleep(interval)
				
				log.Println("Running scheduled content collection...")
				collectorService.CollectAll(ctx)
				collectorService.ProcessUnprocessedItems(ctx)
			}
		}()
		
		adminHandler := handlers.NewAdminHandler(s.db)
		
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
		})
		
		pathsHandler := handlers.NewPathsHandler(s.db)
		progressHandler := handlers.NewProgressHandler(s.db)
		projectsHandler := handlers.NewProjectsHandler(s.db)

		// Protected routes (auth required)
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.VerifyToken)
			
			// Progress & Stats
			r.Post("/progress", progressHandler.UpdateProgress)
			r.Get("/stats", progressHandler.GetStats)
			r.Get("/paths/{id}/progress", progressHandler.GetPathProgress)
			
			// Learning Paths Routes
			r.Get("/paths", pathsHandler.List)
			r.Get("/paths/{id}", pathsHandler.Get)
			r.Post("/paths", pathsHandler.Create)
			r.Post("/paths/{id}/modules", pathsHandler.AddModule)
			r.Post("/lessons", pathsHandler.AddLesson) // keeping it simple, expects module_id in body
			
			// Projects Routes
			r.Get("/projects", projectsHandler.List)
			r.Post("/projects", projectsHandler.Create) // Should be admin only
			r.Post("/submissions", projectsHandler.Submit)
			
			// GitHub Connect (Protected)
			r.Get("/auth/github/connect", githubHandler.Connect)
			
			// Admin Routes (Should have extra middleware check for role)
			r.Group(func(r chi.Router) {
				// r.Use(adminMiddleware) // TODO: Implement Admin Middleware
				r.Post("/admin/settings/collector", adminHandler.UpdateCollectorFrequency)
			})
		})
	})

	return r
}


