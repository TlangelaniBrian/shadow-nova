package database

import (
	"context"
	"fmt"
	"os"
	"shadow-nova/backend/internal/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service interface {
	Health() map[string]string
	Close()
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	InitSchema(ctx context.Context) error
}

type service struct {
	db *pgxpool.Pool
}

var (
	dbInstance *service
)

func New() Service {
	// Reuse instance if it exists
	if dbInstance != nil {
		return dbInstance
	}

	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		// Default for local development if not set
		databaseUrl = "postgres://user:password@localhost:5432/shadownova?sslmode=disable"
	}

	config, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse database URL: %v", err))
	}

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		panic(fmt.Sprintf("Unable to create connection pool: %v", err))
	}

	dbInstance = &service{
		db: db,
	}

	return dbInstance
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.Ping(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		return stats
	}

	stats["status"] = "up"
	stats["message"] = "It's healthy"
	return stats
}

func (s *service) Close() {
	s.db.Close()
}

func (s *service) InitSchema(ctx context.Context) error {
	// Read schema file
	// Note: In a real production app, we'd use embed or a migration tool
	// For now, we'll assume the file is relative to the execution directory or use an absolute path
	// Since we run from root, let's try the relative path
	schemaPath := "backend/internal/database/schema.sql"
	
	// If running from backend dir
	if _, err := os.Stat(schemaPath); os.IsNotExist(err) {
		schemaPath = "internal/database/schema.sql"
	}

	content, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	// Execute schema
	_, err = s.db.Exec(ctx, string(content))
	if err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	return nil
}
