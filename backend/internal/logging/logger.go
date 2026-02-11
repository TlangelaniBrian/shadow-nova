package logging

import (
	"log/slog"
	"os"
)

var logger *slog.Logger

// Init initializes the structured logger
func Init() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger = slog.New(handler)
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
