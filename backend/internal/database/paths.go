package database

import (
	"context"
	"fmt"
	"time"

	"shadow-nova/backend/internal/errors"
	"shadow-nova/backend/internal/models"

	"github.com/jackc/pgx/v5"
)

func (s *service) GetLearningPaths(ctx context.Context, limit, offset int) ([]models.LearningPath, error) {
	query := `
		SELECT id, title, description, difficulty, created_at
		FROM learning_paths
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := s.db.Query(ctx, query, limit, offset)
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
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating learning paths: %w", err)
	}

	return paths, nil
}

func (s *service) GetLearningPathsCount(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM learning_paths WHERE deleted_at IS NULL`
	err := s.db.QueryRow(ctx, query).Scan(&count)
	return count, err
}

func (s *service) GetLearningPath(ctx context.Context, id string) (*models.LearningPath, error) {
	query := `
		SELECT id, title, description, difficulty, created_at
		FROM learning_paths
		WHERE id = $1 AND deleted_at IS NULL
	`

	var p models.LearningPath
	err := s.db.QueryRow(ctx, query, id).Scan(&p.ID, &p.Title, &p.Description, &p.Difficulty, &p.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound(fmt.Sprintf("learning path %s not found", id))
		}
		return nil, errors.DatabaseError(err, "failed to get learning path")
	}

	// Fetch all modules and their lessons in a single query using LEFT JOIN
	modulesLessonsQuery := `
		SELECT
			m.id, m.path_id, m.title, m.description, m.order_index, m.created_at,
			l.id, l.module_id, l.title, l.content_type, l.content_url, l.content_body, l.duration_minutes, l.order_index, l.created_at
		FROM modules m
		LEFT JOIN lessons l ON m.id = l.module_id AND l.deleted_at IS NULL
		WHERE m.path_id = $1 AND m.deleted_at IS NULL
		ORDER BY m.order_index ASC, l.order_index ASC
	`

	rows, err := s.db.Query(ctx, modulesLessonsQuery, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get modules and lessons: %w", err)
	}
	defer rows.Close()

	// Map to group lessons by module ID
	moduleMap := make(map[int]*models.Module)
	var moduleOrder []int // Track insertion order

	for rows.Next() {
		var m models.Module
		var lessonID, lessonModuleID, lessonDurationMinutes, lessonOrderIndex *int
		var lessonTitle, lessonContentType, lessonContentURL, lessonContentBody *string
		var lessonCreatedAt *time.Time

		err := rows.Scan(
			&m.ID, &m.PathID, &m.Title, &m.Description, &m.OrderIndex, &m.CreatedAt,
			&lessonID, &lessonModuleID, &lessonTitle, &lessonContentType, &lessonContentURL, &lessonContentBody, &lessonDurationMinutes, &lessonOrderIndex, &lessonCreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan module and lesson: %w", err)
		}

		// Check if module already exists in map
		if _, exists := moduleMap[m.ID]; !exists {
			m.Lessons = []models.Lesson{}
			moduleMap[m.ID] = &m
			moduleOrder = append(moduleOrder, m.ID)
		}

		// If lesson exists (LEFT JOIN may return NULL for modules without lessons)
		if lessonID != nil {
			lesson := models.Lesson{
				ID:              *lessonID,
				ModuleID:        *lessonModuleID,
				Title:           *lessonTitle,
				ContentType:     *lessonContentType,
				ContentURL:      *lessonContentURL,
				ContentBody:     *lessonContentBody,
				DurationMinutes: *lessonDurationMinutes,
				OrderIndex:      *lessonOrderIndex,
				CreatedAt:       *lessonCreatedAt,
			}
			moduleMap[m.ID].Lessons = append(moduleMap[m.ID].Lessons, lesson)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating modules and lessons: %w", err)
	}

	// Build modules slice in correct order
	p.Modules = make([]models.Module, 0, len(moduleOrder))
	for _, moduleID := range moduleOrder {
		p.Modules = append(p.Modules, *moduleMap[moduleID])
	}

	return &p, nil
}

func (s *service) CreateLearningPath(ctx context.Context, path *models.LearningPath) error {
	return s.WithTx(ctx, func(tx pgx.Tx) error {
		query := `
			INSERT INTO learning_paths (id, title, description, difficulty)
			VALUES ($1, $2, $3, $4)
			RETURNING created_at
		`

		err := tx.QueryRow(ctx, query, path.ID, path.Title, path.Description, path.Difficulty).Scan(&path.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to create learning path: %w", err)
		}

		// If the path has modules, create them in the same transaction
		for i := range path.Modules {
			module := &path.Modules[i]
			module.PathID = path.ID
			moduleQuery := `
				INSERT INTO modules (path_id, title, description, order_index)
				VALUES ($1, $2, $3, $4)
				RETURNING id, created_at
			`
			err = tx.QueryRow(ctx, moduleQuery, module.PathID, module.Title, module.Description, module.OrderIndex).Scan(&module.ID, &module.CreatedAt)
			if err != nil {
				return fmt.Errorf("failed to create module: %w", err)
			}

			// If the module has lessons, create them too
			for j := range module.Lessons {
				lesson := &module.Lessons[j]
				lesson.ModuleID = module.ID
				lessonQuery := `
					INSERT INTO lessons (module_id, title, content_type, content_url, content_body, duration_minutes, order_index)
					VALUES ($1, $2, $3, $4, $5, $6, $7)
					RETURNING id, created_at
				`
				err = tx.QueryRow(ctx, lessonQuery, lesson.ModuleID, lesson.Title, lesson.ContentType, lesson.ContentURL, lesson.ContentBody, lesson.DurationMinutes, lesson.OrderIndex).Scan(&lesson.ID, &lesson.CreatedAt)
				if err != nil {
					return fmt.Errorf("failed to create lesson: %w", err)
				}
			}
		}

		return nil
	})
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

// UpdateLearningPath updates an existing learning path.
func (s *service) UpdateLearningPath(ctx context.Context, id string, updates *models.LearningPath) error {
	query := `
		UPDATE learning_paths
		SET title = $1, description = $2, difficulty = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4 AND deleted_at IS NULL
	`
	result, err := s.db.Exec(ctx, query, updates.Title, updates.Description, updates.Difficulty, id)
	if err != nil {
		return fmt.Errorf("failed to update learning path: %w", err)
	}
	if result.RowsAffected() == 0 {
		return errors.NotFound(fmt.Sprintf("learning path %s not found", id))
	}
	return nil
}

// UpdateModule updates an existing module.
func (s *service) UpdateModule(ctx context.Context, id int, updates *models.Module) error {
	query := `
		UPDATE modules
		SET title = $1, description = $2, order_index = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4 AND deleted_at IS NULL
	`
	result, err := s.db.Exec(ctx, query, updates.Title, updates.Description, updates.OrderIndex, id)
	if err != nil {
		return fmt.Errorf("failed to update module: %w", err)
	}
	if result.RowsAffected() == 0 {
		return errors.NotFound(fmt.Sprintf("module %d not found", id))
	}
	return nil
}

// UpdateLesson updates an existing lesson.
func (s *service) UpdateLesson(ctx context.Context, id int, updates *models.Lesson) error {
	query := `
		UPDATE lessons
		SET title = $1, content_type = $2, content_url = $3,
		    content_body = $4, duration_minutes = $5, order_index = $6,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $7 AND deleted_at IS NULL
	`
	result, err := s.db.Exec(ctx, query, updates.Title, updates.ContentType,
		updates.ContentURL, updates.ContentBody, updates.DurationMinutes,
		updates.OrderIndex, id)
	if err != nil {
		return fmt.Errorf("failed to update lesson: %w", err)
	}
	if result.RowsAffected() == 0 {
		return errors.NotFound(fmt.Sprintf("lesson %d not found", id))
	}
	return nil
}

// DeleteLearningPath performs a soft delete on a learning path.
// Cascade behavior will soft delete modules and lessons due to triggers.
func (s *service) DeleteLearningPath(ctx context.Context, id string) error {
	return s.WithTx(ctx, func(tx pgx.Tx) error {
		// Soft delete the path
		query := `UPDATE learning_paths SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1 AND deleted_at IS NULL`
		result, err := tx.Exec(ctx, query, id)
		if err != nil {
			return fmt.Errorf("failed to delete learning path: %w", err)
		}
		if result.RowsAffected() == 0 {
			return errors.NotFound(fmt.Sprintf("learning path %s not found", id))
		}

		// Soft delete associated modules
		moduleQuery := `UPDATE modules SET deleted_at = CURRENT_TIMESTAMP WHERE path_id = $1 AND deleted_at IS NULL`
		_, err = tx.Exec(ctx, moduleQuery, id)
		if err != nil {
			return fmt.Errorf("failed to delete modules: %w", err)
		}

		// Soft delete associated lessons
		lessonQuery := `
			UPDATE lessons
			SET deleted_at = CURRENT_TIMESTAMP
			WHERE module_id IN (SELECT id FROM modules WHERE path_id = $1)
			AND deleted_at IS NULL
		`
		_, err = tx.Exec(ctx, lessonQuery, id)
		if err != nil {
			return fmt.Errorf("failed to delete lessons: %w", err)
		}

		return nil
	})
}

// DeleteModule performs a soft delete on a module and its lessons.
func (s *service) DeleteModule(ctx context.Context, id int) error {
	return s.WithTx(ctx, func(tx pgx.Tx) error {
		// Soft delete the module
		query := `UPDATE modules SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1 AND deleted_at IS NULL`
		result, err := tx.Exec(ctx, query, id)
		if err != nil {
			return fmt.Errorf("failed to delete module: %w", err)
		}
		if result.RowsAffected() == 0 {
			return errors.NotFound(fmt.Sprintf("module %d not found", id))
		}

		// Soft delete associated lessons
		lessonQuery := `UPDATE lessons SET deleted_at = CURRENT_TIMESTAMP WHERE module_id = $1 AND deleted_at IS NULL`
		_, err = tx.Exec(ctx, lessonQuery, id)
		if err != nil {
			return fmt.Errorf("failed to delete lessons: %w", err)
		}

		return nil
	})
}

// DeleteLesson performs a soft delete on a lesson.
func (s *service) DeleteLesson(ctx context.Context, id int) error {
	query := `UPDATE lessons SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1 AND deleted_at IS NULL`
	result, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete lesson: %w", err)
	}
	if result.RowsAffected() == 0 {
		return errors.NotFound(fmt.Sprintf("lesson %d not found", id))
	}
	return nil
}

// GetLesson retrieves a specific lesson by ID.
func (s *service) GetLesson(ctx context.Context, lessonID int) (*models.Lesson, error) {
	query := `
		SELECT id, module_id, title, content_type, content_url, content_body, duration_minutes, order_index, created_at
		FROM lessons
		WHERE id = $1 AND deleted_at IS NULL
	`

	var lesson models.Lesson
	err := s.db.QueryRow(ctx, query, lessonID).Scan(
		&lesson.ID,
		&lesson.ModuleID,
		&lesson.Title,
		&lesson.ContentType,
		&lesson.ContentURL,
		&lesson.ContentBody,
		&lesson.DurationMinutes,
		&lesson.OrderIndex,
		&lesson.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound(fmt.Sprintf("lesson %d not found", lessonID))
		}
		return nil, errors.DatabaseError(err, "failed to get lesson")
	}

	return &lesson, nil
}

// UserHasAccessToPath checks if a user has access to a learning path.
// For now, all authenticated users can access all paths.
// In the future, this can be extended to check enrollment or purchase status.
func (s *service) UserHasAccessToPath(ctx context.Context, userID int, pathID string) (bool, error) {
	// First verify the path exists
	query := `SELECT id FROM learning_paths WHERE id = $1 AND deleted_at IS NULL`
	var id string
	err := s.db.QueryRow(ctx, query, pathID).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, errors.NotFound(fmt.Sprintf("learning path %s not found", pathID))
		}
		return false, errors.DatabaseError(err, "failed to verify path exists")
	}

	// For now, all authenticated users can access all paths
	// Future: Check enrollment table: SELECT EXISTS(SELECT 1 FROM enrollments WHERE user_id = $1 AND path_id = $2)
	return true, nil
}
