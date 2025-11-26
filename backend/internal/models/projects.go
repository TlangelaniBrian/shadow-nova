package models

import "time"

type Project struct {
	ID          string    `json:"id"`
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description"`
	Difficulty  string    `json:"difficulty"`
	TechStack   []string  `json:"tech_stack"`
	CreatedAt   time.Time `json:"created_at"`
}

type ProjectSubmission struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	ProjectID     string    `json:"project_id"`
	GithubRepoURL string    `json:"github_repo_url"`
	PRURL         string    `json:"pr_url"`
	DemoURL       string    `json:"demo_url"`
	Status        string    `json:"status"` // pending, approved, rejected
	Feedback      string    `json:"feedback"`
	SubmittedAt   time.Time `json:"submitted_at"`
}

type GitHubIntegration struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	GithubUserID string    `json:"github_user_id"`
	Username     string    `json:"username"` // GitHub username
	AccessToken  string    `json:"-"` // Never expose token in JSON
	RefreshToken string    `json:"-"`
	TokenExpiry  time.Time `json:"token_expiry"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateProjectRequest struct {
	ID          string   `json:"id" validate:"required,alphanum"`
	Title       string   `json:"title" validate:"required"`
	Description string   `json:"description"`
	Difficulty  string   `json:"difficulty"`
	TechStack   []string `json:"tech_stack"`
}

