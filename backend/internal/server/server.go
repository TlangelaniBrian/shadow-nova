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
	port  string
	db    database.Service
	flags flags.Service
}

func NewServer() *http.Server {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Initialize Unleash
	flagsService, err := flags.New()
	if err != nil {
		fmt.Printf("Warning: Failed to initialize Unleash: %v\n", err)
	}

	NewServer := &Server{
		port:  port,
		db:    database.New(),
		flags: flagsService,
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
				}
				if err := NewServer.db.CreateUser(ctx, user); err != nil {
					fmt.Printf("Failed to seed super user: %v\n", err)
				} else {
					fmt.Println("Super user seeded successfully")
				}
			}
		}()
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	httputil.WriteJSON(w, http.StatusOK, s.db.Health())
}
