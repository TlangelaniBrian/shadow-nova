package database

import (
	"context"
	"fmt"

	"shadow-nova/backend/internal/crypto"
	"shadow-nova/backend/internal/errors"
	"shadow-nova/backend/internal/models"

	"github.com/jackc/pgx/v5"
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
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating projects: %w", err)
	}

	return projects, nil
}

func (s *service) CreateProject(ctx context.Context, project *models.Project) error {
	return s.WithTx(ctx, func(tx pgx.Tx) error {
		query := `
			INSERT INTO projects (id, title, description, difficulty, tech_stack)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING created_at
		`

		err := tx.QueryRow(ctx, query, project.ID, project.Title, project.Description, project.Difficulty, project.TechStack).Scan(&project.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to create project: %w", err)
		}

		// Future: Create initial project metadata, tags, requirements, etc.
		// All within the same transaction

		return nil
	})
}

// --- Submissions ---

func (s *service) SubmitProject(ctx context.Context, sub *models.ProjectSubmission) error {
	return s.WithTx(ctx, func(tx pgx.Tx) error {
		query := `
			INSERT INTO project_submissions (user_id, project_id, github_repo_url, pr_url, demo_url, status)
			VALUES ($1, $2, $3, $4, $5, 'pending')
			RETURNING id, submitted_at
		`

		err := tx.QueryRow(ctx, query, sub.UserID, sub.ProjectID, sub.GithubRepoURL, sub.PRURL, sub.DemoURL).Scan(&sub.ID, &sub.SubmittedAt)
		if err != nil {
			return fmt.Errorf("failed to submit project: %w", err)
		}

		// Future: Create audit log entry for submission
		// auditQuery := `INSERT INTO audit_log (user_id, action, resource_type, resource_id, timestamp) VALUES ($1, 'submit', 'project', $2, NOW())`
		// _, err = tx.Exec(ctx, auditQuery, sub.UserID, sub.ProjectID)
		// if err != nil {
		//     return fmt.Errorf("failed to create audit log: %w", err)
		// }

		return nil
	})
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
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating submissions: %w", err)
	}

	return submissions, nil
}

func (s *service) GetSubmission(ctx context.Context, submissionID int) (*models.ProjectSubmission, error) {
	query := `
		SELECT id, user_id, project_id, github_repo_url, pr_url, demo_url, status, feedback, submitted_at
		FROM project_submissions
		WHERE id = $1
	`

	var sub models.ProjectSubmission
	var feedback *string
	err := s.db.QueryRow(ctx, query, submissionID).Scan(
		&sub.ID, &sub.UserID, &sub.ProjectID, &sub.GithubRepoURL,
		&sub.PRURL, &sub.DemoURL, &sub.Status, &feedback, &sub.SubmittedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound(fmt.Sprintf("submission %d not found", submissionID))
		}
		return nil, errors.DatabaseError(err, "failed to get submission")
	}
	if feedback != nil {
		sub.Feedback = *feedback
	}

	return &sub, nil
}

func (s *service) UpdateSubmission(ctx context.Context, submissionID int, status, feedback string) error {
	query := `
		UPDATE project_submissions
		SET status = $1, feedback = $2
		WHERE id = $3
	`

	_, err := s.db.Exec(ctx, query, status, feedback, submissionID)
	if err != nil {
		return fmt.Errorf("failed to update submission: %w", err)
	}

	return nil
}

// --- GitHub Integration ---

func (s *service) SaveGitHubToken(ctx context.Context, integration *models.GitHubIntegration) error {
	return s.WithTx(ctx, func(tx pgx.Tx) error {
		// Encrypt tokens before storage
		encryptedAccessToken, err := crypto.Encrypt(integration.AccessToken)
		if err != nil {
			return fmt.Errorf("failed to encrypt access token: %w", err)
		}

		var encryptedRefreshToken *string
		if integration.RefreshToken != "" {
			encrypted, err := crypto.Encrypt(integration.RefreshToken)
			if err != nil {
				return fmt.Errorf("failed to encrypt refresh token: %w", err)
			}
			encryptedRefreshToken = &encrypted
		}

		// Save the GitHub integration
		query := `
			INSERT INTO github_integrations (user_id, github_user_id, access_token, refresh_token, token_expiry)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (user_id)
			DO UPDATE SET access_token = $3, refresh_token = $4, token_expiry = $5, updated_at = CURRENT_TIMESTAMP
			RETURNING id, created_at
		`

		err = tx.QueryRow(ctx, query, integration.UserID, integration.GithubUserID, encryptedAccessToken, encryptedRefreshToken, integration.TokenExpiry).Scan(&integration.ID, &integration.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to save github token: %w", err)
		}

		// Update the user's github_username field
		if integration.Username != "" {
			updateUserQuery := `UPDATE users SET github_username = $1 WHERE id = $2`
			_, err = tx.Exec(ctx, updateUserQuery, integration.Username, integration.UserID)
			if err != nil {
				return fmt.Errorf("failed to update github_username: %w", err)
			}
		}

		return nil
	})
}

func (s *service) GetGitHubIntegration(ctx context.Context, userID int) (*models.GitHubIntegration, error) {
	var integration models.GitHubIntegration
	var encryptedAccessToken string
	var encryptedRefreshToken *string

	query := `
		SELECT id, user_id, github_user_id, access_token, refresh_token, token_expiry, created_at
		FROM github_integrations
		WHERE user_id = $1
	`

	err := s.db.QueryRow(ctx, query, userID).Scan(
		&integration.ID,
		&integration.UserID,
		&integration.GithubUserID,
		&encryptedAccessToken,
		&encryptedRefreshToken,
		&integration.TokenExpiry,
		&integration.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound(fmt.Sprintf("github integration for user %d not found", userID))
		}
		return nil, errors.DatabaseError(err, "failed to get github integration")
	}

	// Decrypt tokens
	integration.AccessToken, err = crypto.Decrypt(encryptedAccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt access token: %w", err)
	}

	if encryptedRefreshToken != nil {
		decrypted, err := crypto.Decrypt(*encryptedRefreshToken)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt refresh token: %w", err)
		}
		integration.RefreshToken = decrypted
	}

	return &integration, nil
}

// UserOwnsSubmission checks if a user owns a specific project submission.
// This prevents IDOR vulnerabilities by ensuring users can only access their own submissions.
func (s *service) UserOwnsSubmission(ctx context.Context, userID int, submissionID int) (bool, error) {
	query := `SELECT user_id FROM project_submissions WHERE id = $1`
	var ownerID int
	err := s.db.QueryRow(ctx, query, submissionID).Scan(&ownerID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, errors.NotFound(fmt.Sprintf("submission %d not found", submissionID))
		}
		return false, errors.DatabaseError(err, "failed to get submission owner")
	}
	return ownerID == userID, nil
}
