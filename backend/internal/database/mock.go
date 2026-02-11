package database

import (
	"context"
	"time"

	"shadow-nova/backend/internal/models"

	"github.com/jackc/pgx/v5"
)

type MockService struct {
	GetProjectsFunc             func(ctx context.Context, limit, offset int) ([]models.Project, error)
	GetProjectsCountFunc        func(ctx context.Context) (int, error)
	CreateProjectFunc           func(ctx context.Context, project *models.Project) error
	SubmitProjectFunc           func(ctx context.Context, sub *models.ProjectSubmission) error
	GetUserSubmissionsFunc      func(ctx context.Context, userID int, limit, offset int) ([]models.ProjectSubmission, error)
	GetUserSubmissionsCountFunc func(ctx context.Context, userID int) (int, error)
	SaveGitHubTokenFunc         func(ctx context.Context, integration *models.GitHubIntegration) error
	GetGitHubIntegrationFunc    func(ctx context.Context, userID int) (*models.GitHubIntegration, error)
	DeleteGitHubIntegrationFunc func(ctx context.Context, userID int) error

	// Add other fields as needed for other tests
	GetUserByEmailFunc          func(ctx context.Context, email string) (*models.User, error)
	GetUserByIDFunc             func(ctx context.Context, userID int) (*models.User, error)
	UpdateUserFunc              func(ctx context.Context, userID int, user *models.User) error
	UpdateUserPasswordFunc      func(ctx context.Context, userID int, hashedPassword string) error
	GetUsersFunc                func(ctx context.Context, limit, offset int) ([]models.User, error)
	GetUsersCountFunc           func(ctx context.Context) (int, error)
	DeleteUserFunc              func(ctx context.Context, userID int) error
	GetContentSourcesFunc       func(ctx context.Context, limit, offset int) ([]models.ContentSource, error)
	GetContentSourcesCountFunc  func(ctx context.Context) (int, error)
	CreateContentItemFunc       func(ctx context.Context, item *models.ContentItem) error
	GetUnprocessedItemsFunc     func(ctx context.Context, limit int) ([]models.ContentItem, error)
	UpdateContentItemAIFunc     func(ctx context.Context, item *models.ContentItem) error
	GetSystemSettingFunc        func(ctx context.Context, key string) (string, error)
	UpdateSystemSettingFunc     func(ctx context.Context, key, value string) error
	GetLearningPathsFunc        func(ctx context.Context, limit, offset int) ([]models.LearningPath, error)
	GetLearningPathsCountFunc   func(ctx context.Context) (int, error)

	// Ownership validation
	UserHasAccessToPathFunc func(ctx context.Context, userID int, pathID string) (bool, error)
	UserOwnsSubmissionFunc  func(ctx context.Context, userID int, submissionID int) (bool, error)
	UserOwnsProgressFunc    func(ctx context.Context, userID int, progressID int) (bool, error)
}

func (m *MockService) GetProjects(ctx context.Context, limit, offset int) ([]models.Project, error) {
	if m.GetProjectsFunc != nil {
		return m.GetProjectsFunc(ctx, limit, offset)
	}
	return nil, nil
}

func (m *MockService) GetProjectsCount(ctx context.Context) (int, error) {
	if m.GetProjectsCountFunc != nil {
		return m.GetProjectsCountFunc(ctx)
	}
	return 0, nil
}

func (m *MockService) CreateProject(ctx context.Context, project *models.Project) error {
	if m.CreateProjectFunc != nil {
		return m.CreateProjectFunc(ctx, project)
	}
	return nil
}

func (m *MockService) SubmitProject(ctx context.Context, sub *models.ProjectSubmission) error {
	if m.SubmitProjectFunc != nil {
		return m.SubmitProjectFunc(ctx, sub)
	}
	return nil
}

func (m *MockService) GetUserSubmissions(ctx context.Context, userID int, limit, offset int) ([]models.ProjectSubmission, error) {
	if m.GetUserSubmissionsFunc != nil {
		return m.GetUserSubmissionsFunc(ctx, userID, limit, offset)
	}
	return nil, nil
}

func (m *MockService) GetUserSubmissionsCount(ctx context.Context, userID int) (int, error) {
	if m.GetUserSubmissionsCountFunc != nil {
		return m.GetUserSubmissionsCountFunc(ctx, userID)
	}
	return 0, nil
}

func (m *MockService) SaveGitHubToken(ctx context.Context, integration *models.GitHubIntegration) error {
	if m.SaveGitHubTokenFunc != nil {
		return m.SaveGitHubTokenFunc(ctx, integration)
	}
	return nil
}

func (m *MockService) GetGitHubIntegration(ctx context.Context, userID int) (*models.GitHubIntegration, error) {
	if m.GetGitHubIntegrationFunc != nil {
		return m.GetGitHubIntegrationFunc(ctx, userID)
	}
	return nil, nil
}

func (m *MockService) DeleteGitHubIntegration(ctx context.Context, userID int) error {
	if m.DeleteGitHubIntegrationFunc != nil {
		return m.DeleteGitHubIntegrationFunc(ctx, userID)
	}
	return nil
}

func (m *MockService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if m.GetUserByEmailFunc != nil {
		return m.GetUserByEmailFunc(ctx, email)
	}
	return nil, nil
}

func (m *MockService) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	if m.GetUserByIDFunc != nil {
		return m.GetUserByIDFunc(ctx, userID)
	}
	return nil, nil
}

func (m *MockService) UpdateUser(ctx context.Context, userID int, user *models.User) error {
	if m.UpdateUserFunc != nil {
		return m.UpdateUserFunc(ctx, userID, user)
	}
	return nil
}

func (m *MockService) UpdateUserPassword(ctx context.Context, userID int, hashedPassword string) error {
	if m.UpdateUserPasswordFunc != nil {
		return m.UpdateUserPasswordFunc(ctx, userID, hashedPassword)
	}
	return nil
}

func (m *MockService) GetUsers(ctx context.Context, limit, offset int) ([]models.User, error) {
	if m.GetUsersFunc != nil {
		return m.GetUsersFunc(ctx, limit, offset)
	}
	return nil, nil
}

func (m *MockService) GetUsersCount(ctx context.Context) (int, error) {
	if m.GetUsersCountFunc != nil {
		return m.GetUsersCountFunc(ctx)
	}
	return 0, nil
}

func (m *MockService) DeleteUser(ctx context.Context, userID int) error {
	if m.DeleteUserFunc != nil {
		return m.DeleteUserFunc(ctx, userID)
	}
	return nil
}

// Stub implementations for interface compliance
func (m *MockService) Health() map[string]string { return nil }
func (m *MockService) InitSchema(ctx context.Context) error { return nil }
func (m *MockService) Close() {}
func (m *MockService) CreateUser(ctx context.Context, user *models.User) error { return nil }
func (m *MockService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) { return nil, nil }
func (m *MockService) GetLearningPaths(ctx context.Context, limit, offset int) ([]models.LearningPath, error) {
	if m.GetLearningPathsFunc != nil {
		return m.GetLearningPathsFunc(ctx, limit, offset)
	}
	return nil, nil
}

func (m *MockService) GetLearningPathsCount(ctx context.Context) (int, error) {
	if m.GetLearningPathsCountFunc != nil {
		return m.GetLearningPathsCountFunc(ctx)
	}
	return 0, nil
}
func (m *MockService) GetLearningPath(ctx context.Context, id string) (*models.LearningPath, error) { return nil, nil }
func (m *MockService) GetLesson(ctx context.Context, lessonID int) (*models.Lesson, error) { return nil, nil }
func (m *MockService) CreateLearningPath(ctx context.Context, path *models.LearningPath) error { return nil }
func (m *MockService) UpdateLearningPath(ctx context.Context, id string, updates *models.LearningPath) error { return nil }
func (m *MockService) DeleteLearningPath(ctx context.Context, id string) error { return nil }
func (m *MockService) CreateModule(ctx context.Context, module *models.Module) error { return nil }
func (m *MockService) UpdateModule(ctx context.Context, id int, updates *models.Module) error { return nil }
func (m *MockService) DeleteModule(ctx context.Context, id int) error { return nil }
func (m *MockService) CreateLesson(ctx context.Context, lesson *models.Lesson) error { return nil }
func (m *MockService) UpdateLesson(ctx context.Context, id int, updates *models.Lesson) error { return nil }
func (m *MockService) DeleteLesson(ctx context.Context, id int) error { return nil }
func (m *MockService) SeedLearningPaths(ctx context.Context) error { return nil }
func (m *MockService) UpdateUserProgress(ctx context.Context, userID int, req models.UpdateProgressRequest) error { return nil }
func (m *MockService) GetUserStats(ctx context.Context, userID int) (*models.UserStats, error) { return nil, nil }
func (m *MockService) GetPathProgress(ctx context.Context, userID int, pathID string) (*models.PathProgress, error) { return nil, nil }
func (m *MockService) GetUserProgressForPath(ctx context.Context, userID int, pathID string) ([]models.UserProgress, error) { return nil, nil }
func (m *MockService) CreateContentSource(ctx context.Context, source *models.ContentSource) error { return nil }

func (m *MockService) GetContentSources(ctx context.Context, limit, offset int) ([]models.ContentSource, error) {
	if m.GetContentSourcesFunc != nil {
		return m.GetContentSourcesFunc(ctx, limit, offset)
	}
	return nil, nil
}

func (m *MockService) GetContentSourcesCount(ctx context.Context) (int, error) {
	if m.GetContentSourcesCountFunc != nil {
		return m.GetContentSourcesCountFunc(ctx)
	}
	return 0, nil
}

func (m *MockService) CreateContentItem(ctx context.Context, item *models.ContentItem) error {
	if m.CreateContentItemFunc != nil {
		return m.CreateContentItemFunc(ctx, item)
	}
	return nil
}

func (m *MockService) GetUnprocessedItems(ctx context.Context, limit int) ([]models.ContentItem, error) {
	if m.GetUnprocessedItemsFunc != nil {
		return m.GetUnprocessedItemsFunc(ctx, limit)
	}
	return nil, nil
}

func (m *MockService) UpdateContentItemAI(ctx context.Context, item *models.ContentItem) error {
	if m.UpdateContentItemAIFunc != nil {
		return m.UpdateContentItemAIFunc(ctx, item)
	}
	return nil
}
func (m *MockService) GetSystemSetting(ctx context.Context, key string) (string, error) {
	if m.GetSystemSettingFunc != nil {
		return m.GetSystemSettingFunc(ctx, key)
	}
	return "", nil
}

func (m *MockService) UpdateSystemSetting(ctx context.Context, key, value string) error {
	if m.UpdateSystemSettingFunc != nil {
		return m.UpdateSystemSettingFunc(ctx, key, value)
	}
	return nil
}

func (m *MockService) UserHasAccessToPath(ctx context.Context, userID int, pathID string) (bool, error) {
	if m.UserHasAccessToPathFunc != nil {
		return m.UserHasAccessToPathFunc(ctx, userID, pathID)
	}
	return true, nil
}

func (m *MockService) UserOwnsSubmission(ctx context.Context, userID int, submissionID int) (bool, error) {
	if m.UserOwnsSubmissionFunc != nil {
		return m.UserOwnsSubmissionFunc(ctx, userID, submissionID)
	}
	return true, nil
}

func (m *MockService) UserOwnsProgress(ctx context.Context, userID int, progressID int) (bool, error) {
	if m.UserOwnsProgressFunc != nil {
		return m.UserOwnsProgressFunc(ctx, userID, progressID)
	}
	return true, nil
}

// Token Blacklist stubs
func (m *MockService) BlacklistToken(ctx context.Context, jti string, userID int, expiresAt time.Time, reason string) error {
	return nil
}

func (m *MockService) IsTokenBlacklisted(ctx context.Context, jti string) (bool, error) {
	return false, nil
}

func (m *MockService) BlacklistAllUserTokens(ctx context.Context, userID int, reason string) error {
	return nil
}

func (m *MockService) DeleteExpiredBlacklistedTokens(ctx context.Context) (int64, error) {
	return 0, nil
}

// Admin stubs (if needed)
func (m *MockService) GetSubmission(ctx context.Context, submissionID int) (*models.ProjectSubmission, error) {
	return nil, nil
}

func (m *MockService) UpdateSubmission(ctx context.Context, submissionID int, status, feedback string) error {
	return nil
}

func (m *MockService) GetProject(ctx context.Context, id string) (*models.Project, error) {
	return nil, nil
}

func (m *MockService) UpdateProject(ctx context.Context, id string, updates *models.Project) error {
	return nil
}

func (m *MockService) DeleteProject(ctx context.Context, id string) error {
	return nil
}

// Transaction support stubs
func (m *MockService) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return nil, nil
}

func (m *MockService) WithTx(ctx context.Context, fn func(pgx.Tx) error) error {
	return fn(nil)
}

// Metrics stub
func (m *MockService) StartMetricsCollection(ctx context.Context) {
	// No-op for mock
}

// Idempotency stubs
func (m *MockService) StoreIdempotentResponse(ctx context.Context, key string, userID int, path, method string, status int, body string, expiresAt time.Time) error {
	return nil
}

func (m *MockService) GetIdempotentResponse(ctx context.Context, key string, userID int) (*IdempotentResponse, error) {
	return nil, nil
}

func (m *MockService) DeleteExpiredIdempotencyKeys(ctx context.Context) (int64, error) {
	return 0, nil
}
