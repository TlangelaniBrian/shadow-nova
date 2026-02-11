package logging

import (
	"context"
	"log/slog"
	"os"
)

var logger *slog.Logger

// Init initializes the structured logger
func Init() {
	logLevel := slog.LevelInfo
	if os.Getenv("LOG_LEVEL") == "debug" {
		logLevel = slog.LevelDebug
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true, // Add file:line information
	})
	logger = slog.New(handler)
	slog.SetDefault(logger)
}

// With returns a logger with additional context fields
func With(args ...any) *slog.Logger {
	return logger.With(args...)
}

// WithContext returns a logger with context fields extracted from the context
func WithContext(ctx context.Context) *slog.Logger {
	// Extract request_id from context if present
	if reqID := ctx.Value("request_id"); reqID != nil {
		return logger.With("request_id", reqID)
	}
	return logger
}

// Error logs an error message with structured fields
func Error(msg string, err error, args ...any) {
	allArgs := append([]any{"error", err}, args...)
	logger.Error(msg, allArgs...)
}

// Info logs an informational message
func Info(msg string, args ...any) {
	logger.Info(msg, args...)
}

// Warn logs a warning message
func Warn(msg string, args ...any) {
	logger.Warn(msg, args...)
}

// Debug logs a debug message
func Debug(msg string, args ...any) {
	logger.Debug(msg, args...)
}
