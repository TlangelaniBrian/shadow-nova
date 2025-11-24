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
	
	// Learning Paths
	GetLearningPaths(ctx context.Context) ([]models.LearningPath, error)
	GetLearningPath(ctx context.Context, id string) (*models.LearningPath, error)
	CreateLearningPath(ctx context.Context, path *models.LearningPath) error
	CreateModule(ctx context.Context, module *models.Module) error
	CreateLesson(ctx context.Context, lesson *models.Lesson) error
	
	// Seeding
	SeedLearningPaths(ctx context.Context) error

	// User Progress
	UpdateUserProgress(ctx context.Context, userID int, req models.UpdateProgressRequest) error
	GetUserStats(ctx context.Context, userID int) (*models.UserStats, error)
	GetPathProgress(ctx context.Context, userID int, pathID string) (*models.PathProgress, error)

	// Projects & GitHub
	GetProjects(ctx context.Context) ([]models.Project, error)
	CreateProject(ctx context.Context, project *models.Project) error
	SubmitProject(ctx context.Context, sub *models.ProjectSubmission) error
	GetUserSubmissions(ctx context.Context, userID int) ([]models.ProjectSubmission, error)
	SaveGitHubToken(ctx context.Context, integration *models.GitHubIntegration) error

	// FeedSourceEngine & Content
	CreateContentSource(ctx context.Context, source *models.ContentSource) error
	GetContentSources(ctx context.Context) ([]models.ContentSource, error)
	CreateContentItem(ctx context.Context, item *models.ContentItem) error
	GetUnprocessedItems(ctx context.Context, limit int) ([]models.ContentItem, error)
	UpdateContentItemAI(ctx context.Context, item *models.ContentItem) error
	
	// System Settings
	GetSystemSetting(ctx context.Context, key string) (string, error)
	UpdateSystemSetting(ctx context.Context, key, value string) error
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
