package database

import (
	"context"
	"shadow-nova/backend/internal/models"
)

type MockService struct {
	GetProjectsFunc        func(ctx context.Context) ([]models.Project, error)
	CreateProjectFunc      func(ctx context.Context, project *models.Project) error
	SubmitProjectFunc      func(ctx context.Context, sub *models.ProjectSubmission) error
	GetUserSubmissionsFunc func(ctx context.Context, userID int) ([]models.ProjectSubmission, error)
	SaveGitHubTokenFunc    func(ctx context.Context, integration *models.GitHubIntegration) error
	
	// Add other fields as needed for other tests
	// Add other fields as needed for other tests
	GetUserByEmailFunc      func(ctx context.Context, email string) (*models.User, error)
	GetContentSourcesFunc   func(ctx context.Context) ([]models.ContentSource, error)
	CreateContentItemFunc   func(ctx context.Context, item *models.ContentItem) error
	GetUnprocessedItemsFunc func(ctx context.Context, limit int) ([]models.ContentItem, error)
	UpdateContentItemAIFunc func(ctx context.Context, item *models.ContentItem) error
	GetSystemSettingFunc    func(ctx context.Context, key string) (string, error)
	UpdateSystemSettingFunc func(ctx context.Context, key, value string) error
}

func (m *MockService) GetProjects(ctx context.Context) ([]models.Project, error) {
	if m.GetProjectsFunc != nil {
		return m.GetProjectsFunc(ctx)
	}
	return nil, nil
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

func (m *MockService) GetUserSubmissions(ctx context.Context, userID int) ([]models.ProjectSubmission, error) {
	if m.GetUserSubmissionsFunc != nil {
		return m.GetUserSubmissionsFunc(ctx, userID)
	}
	return nil, nil
}

func (m *MockService) SaveGitHubToken(ctx context.Context, integration *models.GitHubIntegration) error {
	if m.SaveGitHubTokenFunc != nil {
		return m.SaveGitHubTokenFunc(ctx, integration)
	}
	return nil
}

func (m *MockService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if m.GetUserByEmailFunc != nil {
		return m.GetUserByEmailFunc(ctx, email)
	}
	return nil, nil
}

// Stub implementations for interface compliance
func (m *MockService) Health() map[string]string { return nil }
func (m *MockService) InitSchema(ctx context.Context) error { return nil }
func (m *MockService) Close() {}
func (m *MockService) CreateUser(ctx context.Context, user *models.User) error { return nil }
func (m *MockService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) { return nil, nil }
func (m *MockService) GetLearningPaths(ctx context.Context) ([]models.LearningPath, error) { return nil, nil }
func (m *MockService) GetLearningPath(ctx context.Context, id string) (*models.LearningPath, error) { return nil, nil }
func (m *MockService) CreateLearningPath(ctx context.Context, path *models.LearningPath) error { return nil }
func (m *MockService) CreateModule(ctx context.Context, module *models.Module) error { return nil }
func (m *MockService) CreateLesson(ctx context.Context, lesson *models.Lesson) error { return nil }
func (m *MockService) SeedLearningPaths(ctx context.Context) error { return nil }
func (m *MockService) UpdateUserProgress(ctx context.Context, userID int, req models.UpdateProgressRequest) error { return nil }
func (m *MockService) GetUserStats(ctx context.Context, userID int) (*models.UserStats, error) { return nil, nil }
func (m *MockService) GetPathProgress(ctx context.Context, userID int, pathID string) (*models.PathProgress, error) { return nil, nil }
func (m *MockService) CreateContentSource(ctx context.Context, source *models.ContentSource) error { return nil }

func (m *MockService) GetContentSources(ctx context.Context) ([]models.ContentSource, error) {
	if m.GetContentSourcesFunc != nil {
		return m.GetContentSourcesFunc(ctx)
	}
	return nil, nil
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
