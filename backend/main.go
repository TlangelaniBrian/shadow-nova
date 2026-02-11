package main

import (
	"context"
	"errors"
	"log/slog"
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
		// Can't use structured logging yet, so use slog directly
		slog.Warn("env file not found, using environment variables")
	}

	// Initialize structured logging
	logging.Init()

	// Initialize encryption
	if err := crypto.Init(); err != nil {
		logging.Error("failed to initialize encryption", err)
		os.Exit(1)
	}

	// Initialize database (single instance)
	db, err := database.New()
	if err != nil {
		logging.Error("failed to initialize database", err)
		os.Exit(1)
	}
	defer db.Close()

	// Initialize flags service
	flagsService, err := flags.New()
	if err != nil {
		logging.Warn("failed to initialize feature flags", "error", err)
		// Continue without feature flags - the system should handle nil gracefully
	}

	// Pass dependencies to server
	httpServer, appServer, err := server.NewServer(db, flagsService)
	if err != nil {
		logging.Error("failed to create server", err)
		os.Exit(1)
	}

	// Create a channel to listen for interrupt signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// Start server in a goroutine
	serverErrors := make(chan error, 1)
	go func() {
		logging.Info("server starting", "address", httpServer.Addr)
		serverErrors <- httpServer.ListenAndServe()
	}()

	// Block until we receive a signal or server error
	select {
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logging.Error("server error", err)
			os.Exit(1)
		}
	case sig := <-shutdown:
		logging.Info("received shutdown signal", "signal", sig)

		// Create context with timeout for shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Shutdown the HTTP server first
		if err := httpServer.Shutdown(ctx); err != nil {
			logging.Error("http server shutdown error", err)
			// Force close after timeout
			if err := httpServer.Close(); err != nil {
				logging.Error("http server force close error", err)
			}
		}

		// Shutdown background tasks and database
		if err := appServer.Shutdown(ctx); err != nil {
			logging.Error("application shutdown error", err)
		}

		logging.Info("server stopped")
	}
}
