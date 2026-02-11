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
		fmt.Printf("Failed to initialize schema: %v\n", err)
	}

	// Seed Learning Paths
	if err := NewServer.db.SeedLearningPaths(ctx); err != nil {
		fmt.Printf("Failed to seed learning paths: %v\n", err)
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
					fmt.Printf("Failed to hash admin password: %v\n", err)
					return
				}
				user := &models.User{
					Email:        email,
					Username:     "SuperAdmin",
					PasswordHash: string(hashedPassword),
					Role:         "admin",
				}
				if err := NewServer.db.CreateUser(ctx, user); err != nil {
					fmt.Printf("Failed to seed super user: %v\n", err)
				} else {
					fmt.Println("Super user seeded successfully")
				}
			}
		}()
	}

	// Start database metrics collection
	NewServer.db.StartMetricsCollection(collectorCtx)
	fmt.Println("Database metrics collection started")

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
	fmt.Println("Initiating graceful shutdown...")

	// Cancel collector goroutine
	if s.collectorCancel != nil {
		fmt.Println("Stopping collector goroutine...")
		s.collectorCancel()
	}

	// Close database connections
	fmt.Println("Closing database connections...")
	s.db.Close()

	// Close flags service if it has cleanup
	if s.flags != nil {
		fmt.Println("Closing flags service...")
		s.flags.Close()
	}

	fmt.Println("Graceful shutdown complete")
	return nil
}
