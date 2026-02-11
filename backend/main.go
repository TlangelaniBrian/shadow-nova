package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"shadow-nova/backend/internal/crypto"
	"shadow-nova/backend/internal/database"
	"shadow-nova/backend/internal/flags"
	"shadow-nova/backend/internal/logging"
	"shadow-nova/backend/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using environment variables")
	}

	// Initialize structured logging
	logging.Init()

	// Initialize encryption
	if err := crypto.Init(); err != nil {
		log.Fatalf("Failed to initialize encryption: %v", err)
	}

	// Initialize database (single instance)
	db, err := database.New()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize flags service
	flagsService, err := flags.New()
	if err != nil {
		log.Printf("Warning: Failed to initialize feature flags: %v", err)
		// Continue without feature flags - the system should handle nil gracefully
	}

	// Pass dependencies to server
	httpServer, appServer, err := server.NewServer(db, flagsService)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Create a channel to listen for interrupt signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// Start server in a goroutine
	serverErrors := make(chan error, 1)
	go func() {
		fmt.Printf("Server running on %s\n", httpServer.Addr)
		serverErrors <- httpServer.ListenAndServe()
	}()

	// Block until we receive a signal or server error
	select {
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	case sig := <-shutdown:
		log.Printf("Received signal: %v. Starting graceful shutdown...", sig)

		// Create context with timeout for shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Shutdown the HTTP server first
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
			// Force close after timeout
			if err := httpServer.Close(); err != nil {
				log.Printf("HTTP server force close error: %v", err)
			}
		}

		// Shutdown background tasks and database
		if err := appServer.Shutdown(ctx); err != nil {
			log.Printf("Application shutdown error: %v", err)
		}

		log.Println("Server stopped")
	}
}
