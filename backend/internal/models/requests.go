package models

// User registration request
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=20,alphanum"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

// User login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Learning path progress update
type UpdateProgressRequest struct {
	PathID      string `json:"path_id" validate:"required"`
	ModuleIndex int    `json:"module_index" validate:"required,gte=0"`
	Completed   bool   `json:"completed"`
}

// Project submission request
type SubmitProjectRequest struct {
	ProjectID string `json:"project_id" validate:"required"`
	GithubURL string `json:"github_url" validate:"omitempty,url"`
	DemoURL   string `json:"demo_url" validate:"omitempty,url"`
}

// Generic response structures
type ErrorResponse struct {
	Error   string            `json:"error"`
	Details map[string]string `json:"details,omitempty"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
