package models

import "time"

type LearningPath struct {
	ID          string    `json:"id"`
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description"`
	Difficulty  string    `json:"difficulty"`
	CreatedAt   time.Time `json:"created_at"`
	Modules     []Module  `json:"modules,omitempty"` // For nesting in responses
}

type Module struct {
	ID          int       `json:"id"`
	PathID      string    `json:"path_id"`
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description"`
	OrderIndex  int       `json:"order_index"`
	CreatedAt   time.Time `json:"created_at"`
	Lessons     []Lesson  `json:"lessons,omitempty"`
}

type Lesson struct {
	ID              int       `json:"id"`
	ModuleID        int       `json:"module_id"`
	Title           string    `json:"title" validate:"required"`
	ContentType     string    `json:"content_type" validate:"required,oneof=video article quiz"`
	ContentURL      string    `json:"content_url"`
	ContentBody     string    `json:"content_body"`
	DurationMinutes int       `json:"duration_minutes"`
	OrderIndex      int       `json:"order_index"`
	CreatedAt       time.Time `json:"created_at"`
}

type CreatePathRequest struct {
	ID          string `json:"id" validate:"required,alphanum"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	Difficulty  string `json:"difficulty"`
}

type CreateModuleRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	OrderIndex  int    `json:"order_index"`
}
