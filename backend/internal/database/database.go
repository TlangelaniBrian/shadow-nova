package database

import (
	"context"
	"fmt"
	"os"
	"shadow-nova/backend/internal/logging"
	"shadow-nova/backend/internal/models"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service interface {
	Health() map[string]string
	Close()
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, userID int) (*models.User, error)
	UpdateUser(ctx context.Context, userID int, user *models.User) error
	UpdateUserPassword(ctx context.Context, userID int, hashedPassword string) error
	InitSchema(ctx context.Context) error

	// Learning Paths
	GetLearningPaths(ctx context.Context, limit, offset int) ([]models.LearningPath, error)
	GetLearningPathsCount(ctx context.Context) (int, error)
	GetLearningPath(ctx context.Context, id string) (*models.LearningPath, error)
	CreateLearningPath(ctx context.Context, path *models.LearningPath) error
	UpdateLearningPath(ctx context.Context, id string, updates *models.LearningPath) error
	DeleteLearningPath(ctx context.Context, id string) error
	CreateModule(ctx context.Context, module *models.Module) error
	UpdateModule(ctx context.Context, id int, updates *models.Module) error
	DeleteModule(ctx context.Context, id int) error
	CreateLesson(ctx context.Context, lesson *models.Lesson) error
	UpdateLesson(ctx context.Context, id int, updates *models.Lesson) error
	DeleteLesson(ctx context.Context, id int) error
	GetLesson(ctx context.Context, lessonID int) (*models.Lesson, error)

	// Seeding
	SeedLearningPaths(ctx context.Context) error

	// User Progress
	UpdateUserProgress(ctx context.Context, userID int, req models.UpdateProgressRequest) error
	GetUserStats(ctx context.Context, userID int) (*models.UserStats, error)
	GetPathProgress(ctx context.Context, userID int, pathID string) (*models.PathProgress, error)
	GetUserProgressForPath(ctx context.Context, userID int, pathID string) ([]models.UserProgress, error)

	// Projects & GitHub
	GetProjects(ctx context.Context, limit, offset int) ([]models.Project, error)
	GetProjectsCount(ctx context.Context) (int, error)
	GetProject(ctx context.Context, id string) (*models.Project, error)
	CreateProject(ctx context.Context, project *models.Project) error
	UpdateProject(ctx context.Context, id string, updates *models.Project) error
	DeleteProject(ctx context.Context, id string) error
	SubmitProject(ctx context.Context, sub *models.ProjectSubmission) error
	GetUserSubmissions(ctx context.Context, userID int, limit, offset int) ([]models.ProjectSubmission, error)
	GetUserSubmissionsCount(ctx context.Context, userID int) (int, error)
	GetSubmission(ctx context.Context, submissionID int) (*models.ProjectSubmission, error)
	UpdateSubmission(ctx context.Context, submissionID int, status, feedback string) error
	SaveGitHubToken(ctx context.Context, integration *models.GitHubIntegration) error
	GetGitHubIntegration(ctx context.Context, userID int) (*models.GitHubIntegration, error)
	DeleteGitHubIntegration(ctx context.Context, userID int) error

	// FeedSourceEngine & Content
	CreateContentSource(ctx context.Context, source *models.ContentSource) error
	GetContentSources(ctx context.Context, limit, offset int) ([]models.ContentSource, error)
	GetContentSourcesCount(ctx context.Context) (int, error)
	CreateContentItem(ctx context.Context, item *models.ContentItem) error
	GetUnprocessedItems(ctx context.Context, limit int) ([]models.ContentItem, error)
	UpdateContentItemAI(ctx context.Context, item *models.ContentItem) error

	// System Settings
	GetSystemSetting(ctx context.Context, key string) (string, error)
	UpdateSystemSetting(ctx context.Context, key, value string) error

	// Token Blacklist
	BlacklistToken(ctx context.Context, jti string, userID int, expiresAt time.Time, reason string) error
	IsTokenBlacklisted(ctx context.Context, jti string) (bool, error)
	BlacklistAllUserTokens(ctx context.Context, userID int, reason string) error
	DeleteExpiredBlacklistedTokens(ctx context.Context) (int64, error)

	// Idempotency
	StoreIdempotentResponse(ctx context.Context, key string, userID int, path, method string, status int, body string, expiresAt time.Time) error
	GetIdempotentResponse(ctx context.Context, key string, userID int) (*IdempotentResponse, error)
	DeleteExpiredIdempotencyKeys(ctx context.Context) (int64, error)

	// Ownership Validation (IDOR Prevention)
	UserHasAccessToPath(ctx context.Context, userID int, pathID string) (bool, error)
	UserOwnsSubmission(ctx context.Context, userID int, submissionID int) (bool, error)
	UserOwnsProgress(ctx context.Context, userID int, progressID int) (bool, error)

	// Transaction support
	BeginTx(ctx context.Context) (pgx.Tx, error)
	WithTx(ctx context.Context, fn func(pgx.Tx) error) error

	// Metrics
	StartMetricsCollection(ctx context.Context)
}

type service struct {
	db *pgxpool.Pool
}

// IdempotentResponse represents a cached response from a previous request
type IdempotentResponse struct {
	Key        string
	StatusCode int
	Body       string
}

func New() (Service, error) {
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		// Default for local development if not set
		databaseUrl = "postgres://user:password@localhost:5432/shadownova?sslmode=disable"
	}

	config, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %w", err)
	}

	// Configure connection pool for production
	config.MaxConns = int32(getEnvInt("DB_MAX_CONNS", 25))
	config.MinConns = int32(getEnvInt("DB_MIN_CONNS", 5))
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = time.Minute
	config.ConnConfig.ConnectTimeout = time.Duration(getEnvInt("DB_CONNECT_TIMEOUT", 5)) * time.Second

	// Add connection pool metrics callback
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		logging.Debug("new database connection established")
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Test the connection
	if err := db.Ping(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	logging.Info("database connection pool initialized", "max_conns", config.MaxConns, "min_conns", config.MinConns)

	return &service{db: db}, nil
}

// Helper function to get int from env with default
func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	poolStats := s.db.Stat()

	health := map[string]string{
		"status":       "up",
		"database":     "connected",
		"active_conns": fmt.Sprintf("%d", poolStats.AcquiredConns()),
		"idle_conns":   fmt.Sprintf("%d", poolStats.IdleConns()),
		"max_conns":    fmt.Sprintf("%d", poolStats.MaxConns()),
	}

	if err := s.db.Ping(ctx); err != nil {
		health["status"] = "down"
		health["database"] = "disconnected"
		health["error"] = err.Error()
	}

	return health
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

// BeginTx starts a new database transaction.
func (s *service) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return s.db.Begin(ctx)
}

// WithTx executes a function within a transaction.
// If the function returns an error, the transaction is rolled back.
// If the function completes successfully, the transaction is committed.
// Panics are caught, the transaction is rolled back, and the panic is re-raised.
func (s *service) WithTx(ctx context.Context, fn func(pgx.Tx) error) error {
	// Get transaction timeout from environment, default to 30 seconds
	timeoutStr := os.Getenv("DB_TX_TIMEOUT")
	timeout := 30 * time.Second
	if timeoutStr != "" {
		if duration, err := time.ParseDuration(timeoutStr); err == nil {
			timeout = duration
		}
	}

	// Create context with timeout for the transaction
	txCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	tx, err := s.db.Begin(txCtx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Ensure transaction is finalized
	defer func() {
		if p := recover(); p != nil {
			// Rollback on panic
			if rbErr := tx.Rollback(txCtx); rbErr != nil {
				logging.Error("failed to rollback transaction after panic", rbErr)
			}
			panic(p) // Re-raise panic after rollback
		} else if err != nil {
			// Rollback on error
			if rbErr := tx.Rollback(txCtx); rbErr != nil {
				logging.Error("failed to rollback transaction", rbErr)
			}
		} else {
			// Commit on success
			err = tx.Commit(txCtx)
			if err != nil {
				logging.Error("failed to commit transaction", err)
			}
		}
	}()

	err = fn(tx)
	return err
}
