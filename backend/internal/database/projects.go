package database

import (
	"context"
	"fmt"
	"shadow-nova/backend/internal/models"
)

// --- Projects ---

func (s *service) GetProjects(ctx context.Context) ([]models.Project, error) {
	query := `
		SELECT id, title, description, difficulty, tech_stack, created_at
		FROM projects
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query projects: %w", err)
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var p models.Project
		if err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.Difficulty, &p.TechStack, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}
		projects = append(projects, p)
	}

	return projects, nil
}

func (s *service) CreateProject(ctx context.Context, project *models.Project) error {
	query := `
		INSERT INTO projects (id, title, description, difficulty, tech_stack)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at
	`

	err := s.db.QueryRow(ctx, query, project.ID, project.Title, project.Description, project.Difficulty, project.TechStack).Scan(&project.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}

	return nil
}

// --- Submissions ---

func (s *service) SubmitProject(ctx context.Context, sub *models.ProjectSubmission) error {
	query := `
		INSERT INTO project_submissions (user_id, project_id, github_repo_url, pr_url, demo_url, status)
		VALUES ($1, $2, $3, $4, $5, 'pending')
		RETURNING id, submitted_at
	`

	err := s.db.QueryRow(ctx, query, sub.UserID, sub.ProjectID, sub.GithubRepoURL, sub.PRURL, sub.DemoURL).Scan(&sub.ID, &sub.SubmittedAt)
	if err != nil {
		return fmt.Errorf("failed to submit project: %w", err)
	}

	return nil
}

func (s *service) GetUserSubmissions(ctx context.Context, userID int) ([]models.ProjectSubmission, error) {
	query := `
		SELECT id, user_id, project_id, github_repo_url, pr_url, demo_url, status, feedback, submitted_at
		FROM project_submissions
		WHERE user_id = $1
		ORDER BY submitted_at DESC
	`

	rows, err := s.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query submissions: %w", err)
	}
	defer rows.Close()

	var submissions []models.ProjectSubmission
	for rows.Next() {
		var s models.ProjectSubmission
		// Handle nullable feedback
		var feedback *string
		if err := rows.Scan(&s.ID, &s.UserID, &s.ProjectID, &s.GithubRepoURL, &s.PRURL, &s.DemoURL, &s.Status, &feedback, &s.SubmittedAt); err != nil {
			return nil, fmt.Errorf("failed to scan submission: %w", err)
		}
		if feedback != nil {
			s.Feedback = *feedback
		}
		submissions = append(submissions, s)
	}

	return submissions, nil
}

// --- GitHub Integration ---

func (s *service) SaveGitHubToken(ctx context.Context, integration *models.GitHubIntegration) error {
	// Save the GitHub integration
	query := `
		INSERT INTO github_integrations (user_id, github_user_id, access_token, refresh_token, token_expiry)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id) 
		DO UPDATE SET access_token = $3, refresh_token = $4, token_expiry = $5, updated_at = CURRENT_TIMESTAMP
		RETURNING id, created_at
	`

	err := s.db.QueryRow(ctx, query, integration.UserID, integration.GithubUserID, integration.AccessToken, integration.RefreshToken, integration.TokenExpiry).Scan(&integration.ID, &integration.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to save github token: %w", err)
	}
	
	// Update the user's github_username field
	if integration.Username != "" {
		updateUserQuery := `UPDATE users SET github_username = $1 WHERE id = $2`
		_, err = s.db.Exec(ctx, updateUserQuery, integration.Username, integration.UserID)
		if err != nil {
			// Log but don't fail - the integration was saved successfully
			fmt.Printf("Warning: Failed to update user's github_username: %v\n", err)
		}
	}

	return nil
}
