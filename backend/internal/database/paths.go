package database

import (
	"context"
	"fmt"
	"shadow-nova/backend/internal/models"
)

func (s *service) GetLearningPaths(ctx context.Context) ([]models.LearningPath, error) {
	query := `
		SELECT id, title, description, difficulty, created_at
		FROM learning_paths
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query learning paths: %w", err)
	}
	defer rows.Close()

	var paths []models.LearningPath
	for rows.Next() {
		var p models.LearningPath
		if err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.Difficulty, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan learning path: %w", err)
		}
		paths = append(paths, p)
	}

	return paths, nil
}

func (s *service) GetLearningPath(ctx context.Context, id string) (*models.LearningPath, error) {
	query := `
		SELECT id, title, description, difficulty, created_at
		FROM learning_paths
		WHERE id = $1
	`

	var p models.LearningPath
	err := s.db.QueryRow(ctx, query, id).Scan(&p.ID, &p.Title, &p.Description, &p.Difficulty, &p.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get learning path: %w", err)
	}

	// Fetch Modules
	moduleQuery := `
		SELECT id, path_id, title, description, order_index, created_at
		FROM modules
		WHERE path_id = $1
		ORDER BY order_index ASC
	`
	rows, err := s.db.Query(ctx, moduleQuery, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get modules: %w", err)
	}
	defer rows.Close()

	var modules []models.Module
	for rows.Next() {
		var m models.Module
		if err := rows.Scan(&m.ID, &m.PathID, &m.Title, &m.Description, &m.OrderIndex, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan module: %w", err)
		}
		modules = append(modules, m)
	}
	p.Modules = modules

	// Fetch Lessons for each module (N+1 query for now, can optimize later if needed)
	for i := range p.Modules {
		lessonQuery := `
			SELECT id, module_id, title, content_type, content_url, content_body, duration_minutes, order_index, created_at
			FROM lessons
			WHERE module_id = $1
			ORDER BY order_index ASC
		`
		lRows, err := s.db.Query(ctx, lessonQuery, p.Modules[i].ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get lessons for module %d: %w", p.Modules[i].ID, err)
		}
		defer lRows.Close()

		var lessons []models.Lesson
		for lRows.Next() {
			var l models.Lesson
			if err := lRows.Scan(&l.ID, &l.ModuleID, &l.Title, &l.ContentType, &l.ContentURL, &l.ContentBody, &l.DurationMinutes, &l.OrderIndex, &l.CreatedAt); err != nil {
				return nil, fmt.Errorf("failed to scan lesson: %w", err)
			}
			lessons = append(lessons, l)
		}
		p.Modules[i].Lessons = lessons
	}

	return &p, nil
}

func (s *service) CreateLearningPath(ctx context.Context, path *models.LearningPath) error {
	query := `
		INSERT INTO learning_paths (id, title, description, difficulty)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at
	`

	err := s.db.QueryRow(ctx, query, path.ID, path.Title, path.Description, path.Difficulty).Scan(&path.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create learning path: %w", err)
	}

	return nil
}

func (s *service) CreateModule(ctx context.Context, module *models.Module) error {
	query := `
		INSERT INTO modules (path_id, title, description, order_index)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err := s.db.QueryRow(ctx, query, module.PathID, module.Title, module.Description, module.OrderIndex).Scan(&module.ID, &module.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create module: %w", err)
	}

	return nil
}

func (s *service) CreateLesson(ctx context.Context, lesson *models.Lesson) error {
	query := `
		INSERT INTO lessons (module_id, title, content_type, content_url, content_body, duration_minutes, order_index)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`

	err := s.db.QueryRow(ctx, query, lesson.ModuleID, lesson.Title, lesson.ContentType, lesson.ContentURL, lesson.ContentBody, lesson.DurationMinutes, lesson.OrderIndex).Scan(&lesson.ID, &lesson.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create lesson: %w", err)
	}

	return nil
}
