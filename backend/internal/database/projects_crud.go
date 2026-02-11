package database

import (
	"context"
	"fmt"

	"shadow-nova/backend/internal/errors"
	"shadow-nova/backend/internal/models"

	"github.com/jackc/pgx/v5"
)

// GetProject retrieves a single project by ID.
func (s *service) GetProject(ctx context.Context, id string) (*models.Project, error) {
	query := `
		SELECT id, title, description, difficulty, tech_stack, created_at
		FROM projects
		WHERE id = $1 AND deleted_at IS NULL
	`

	var p models.Project
	err := s.db.QueryRow(ctx, query, id).Scan(&p.ID, &p.Title, &p.Description, &p.Difficulty, &p.TechStack, &p.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound(fmt.Sprintf("project %s not found", id))
		}
		return nil, errors.DatabaseError(err, "failed to get project")
	}

	return &p, nil
}

// UpdateProject updates an existing project.
func (s *service) UpdateProject(ctx context.Context, id string, updates *models.Project) error {
	query := `
		UPDATE projects
		SET title = $1, description = $2, tech_stack = $3, difficulty = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5 AND deleted_at IS NULL
	`
	result, err := s.db.Exec(ctx, query, updates.Title, updates.Description, updates.TechStack, updates.Difficulty, id)
	if err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}
	if result.RowsAffected() == 0 {
		return errors.NotFound(fmt.Sprintf("project %s not found", id))
	}
	return nil
}

// DeleteProject performs a soft delete on a project.
func (s *service) DeleteProject(ctx context.Context, id string) error {
	query := `UPDATE projects SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1 AND deleted_at IS NULL`
	result, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}
	if result.RowsAffected() == 0 {
		return errors.NotFound(fmt.Sprintf("project %s not found", id))
	}
	return nil
}
