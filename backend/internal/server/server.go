package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"shadow-nova/backend/internal/database"
	"shadow-nova/backend/internal/flags"
	"shadow-nova/backend/internal/httputil"
	"shadow-nova/backend/internal/logging"
	"shadow-nova/backend/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type Server struct {
	port            string
	db              database.Service
	flags           flags.Service
	collectorCtx    context.Context
	collectorCancel context.CancelFunc
}

func NewServer(db database.Service, flagsService flags.Service) (*http.Server, *Server, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Create context for collector goroutine
	collectorCtx, collectorCancel := context.WithCancel(context.Background())

	NewServer := &Server{
		port:            port,
		db:              db,
		flags:           flagsService,
		collectorCtx:    collectorCtx,
		collectorCancel: collectorCancel,
	}

	// Initialize Schema
	ctx := context.Background()
	if err := NewServer.db.InitSchema(ctx); err != nil {
		logging.Error("failed to initialize schema", err)
	}

	// Seed Learning Paths
	if err := NewServer.db.SeedLearningPaths(ctx); err != nil {
		logging.Error("failed to seed learning paths", err)
	}

	// Seed super user (only if ADMIN_DEFAULT_PASSWORD is set)
	adminPassword := os.Getenv("ADMIN_DEFAULT_PASSWORD")
	if adminPassword != "" {
		go func() {
			ctx := context.Background()
			email := "mrbtmkhabela@gmail.com"
			_, err := NewServer.db.GetUserByEmail(ctx, email)
			if err != nil {
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
				if err != nil {
					logging.Error("failed to hash admin password", err)
					return
				}
				user := &models.User{
					Email:        email,
					Username:     "SuperAdmin",
					PasswordHash: string(hashedPassword),
					Role:         "admin",
				}
				if err := NewServer.db.CreateUser(ctx, user); err != nil {
					logging.Error("failed to seed super user", err)
				} else {
					logging.Info("super user seeded successfully", "email", email)
				}
			}
		}()
	}

	// Start database metrics collection
	NewServer.db.StartMetricsCollection(collectorCtx)
	logging.Info("database metrics collection started")

	// Register routes
	handler, err := NewServer.RegisterRoutes()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to register routes: %w", err)
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", NewServer.port),
		Handler:      handler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server, NewServer, nil
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	httputil.WriteJSON(w, http.StatusOK, s.db.Health())
}

// Shutdown gracefully shuts down the server, stopping background tasks and closing connections
func (s *Server) Shutdown(ctx context.Context) error {
	logging.Info("initiating graceful shutdown")

	// Cancel collector goroutine
	if s.collectorCancel != nil {
		logging.Info("stopping collector goroutine")
		s.collectorCancel()
	}

	// Close database connections
	logging.Info("closing database connections")
	s.db.Close()

	// Close flags service if it has cleanup
	if s.flags != nil {
		logging.Info("closing flags service")
		s.flags.Close()
	}

	logging.Info("graceful shutdown complete")
	return nil
}
