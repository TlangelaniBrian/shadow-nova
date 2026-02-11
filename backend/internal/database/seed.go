package database

import (
	"context"
	"shadow-nova/backend/internal/logging"
	"shadow-nova/backend/internal/models"

	"github.com/jackc/pgx/v5"
)

func (s *service) SeedLearningPaths(ctx context.Context) error {
	// Check if paths already exist
	paths, err := s.GetLearningPaths(ctx, 1, 0)
	if err != nil {
		return err
	}
	if len(paths) > 0 {
		logging.Info("learning paths already exist, skipping seeding")
		return nil
	}

	logging.Info("seeding learning paths")

	// --- FRONTEND PATHS ---

	// 1. Frontend Beginner
	feBeginner := &models.LearningPath{
		ID:          "frontend-beginner",
		Title:       "Frontend Development Basics",
		Description: "Start your journey here. Learn HTML, CSS, and how the web works.",
		Difficulty:  "Beginner",
	}
	if err := s.seedPathWithModules(ctx, feBeginner, s.seedFrontendBeginnerModules); err != nil {
		return err
	}

	// 2. Frontend Intermediate
	feInter := &models.LearningPath{
		ID:          "frontend-intermediate",
		Title:       "Frontend Mastery with Vue",
		Description: "Deep dive into modern JavaScript and the Vue.js framework.",
		Difficulty:  "Intermediate",
	}
	if err := s.seedPathWithModules(ctx, feInter, s.seedFrontendIntermediateModules); err != nil {
		return err
	}

	// 3. Frontend Advanced
	feAdv := &models.LearningPath{
		ID:          "frontend-advanced",
		Title:       "Advanced Frontend Architecture",
		Description: "Performance, testing, and large-scale application design.",
		Difficulty:  "Advanced",
	}
	if err := s.seedPathWithModules(ctx, feAdv, s.seedFrontendAdvancedModules); err != nil {
		return err
	}

	// --- BACKEND PATHS ---

	// 4. Backend Beginner
	beBeginner := &models.LearningPath{
		ID:          "backend-beginner",
		Title:       "Backend Basics with Go",
		Description: "Introduction to server-side programming, HTTP, and Go syntax.",
		Difficulty:  "Beginner",
	}
	if err := s.seedPathWithModules(ctx, beBeginner, s.seedBackendBeginnerModules); err != nil {
		return err
	}

	// 5. Backend Intermediate
	beInter := &models.LearningPath{
		ID:          "backend-intermediate",
		Title:       "Building REST APIs",
		Description: "Database integration, authentication, and API design patterns.",
		Difficulty:  "Intermediate",
	}
	if err := s.seedPathWithModules(ctx, beInter, s.seedBackendIntermediateModules); err != nil {
		return err
	}

	// 6. Backend Advanced
	beAdv := &models.LearningPath{
		ID:          "backend-advanced",
		Title:       "Distributed Systems & Microservices",
		Description: "Scaling applications, concurrency, and cloud deployment.",
		Difficulty:  "Advanced",
	}
	if err := s.seedPathWithModules(ctx, beAdv, s.seedBackendAdvancedModules); err != nil {
		return err
	}

	logging.Info("learning paths seeded successfully")
	return nil
}

// seedPathWithModules creates a learning path with its modules and lessons in a single transaction.
func (s *service) seedPathWithModules(ctx context.Context, path *models.LearningPath, seedFunc func(context.Context, pgx.Tx, string) error) error {
	return s.WithTx(ctx, func(tx pgx.Tx) error {
		// Create the learning path
		query := `
			INSERT INTO learning_paths (id, title, description, difficulty)
			VALUES ($1, $2, $3, $4)
			RETURNING created_at
		`
		err := tx.QueryRow(ctx, query, path.ID, path.Title, path.Description, path.Difficulty).Scan(&path.CreatedAt)
		if err != nil {
			return err
		}

		// Seed modules and lessons for this path
		return seedFunc(ctx, tx, path.ID)
	})
}

// --- Helper Functions ---

func (s *service) seedFrontendBeginnerModules(ctx context.Context, tx pgx.Tx, pathID string) error {
	m1 := &models.Module{PathID: pathID, Title: "The Web Ecosystem", Description: "How browsers and servers talk", OrderIndex: 1}
	if err := s.createModuleTx(ctx, tx, m1); err != nil {
		return err
	}
	if err := s.createLessonTx(ctx, tx, &models.Lesson{ModuleID: m1.ID, Title: "HTTP & DNS", ContentType: "video", DurationMinutes: 10, OrderIndex: 1}); err != nil {
		return err
	}

	m2 := &models.Module{PathID: pathID, Title: "HTML5 Essentials", Description: "Structuring content", OrderIndex: 2}
	if err := s.createModuleTx(ctx, tx, m2); err != nil {
		return err
	}
	return s.createLessonTx(ctx, tx, &models.Lesson{ModuleID: m2.ID, Title: "Semantic HTML", ContentType: "article", DurationMinutes: 15, OrderIndex: 1})
}

func (s *service) seedFrontendIntermediateModules(ctx context.Context, tx pgx.Tx, pathID string) error {
	m1 := &models.Module{PathID: pathID, Title: "JavaScript Deep Dive", Description: "ES6+, Async/Await", OrderIndex: 1}
	if err := s.createModuleTx(ctx, tx, m1); err != nil {
		return err
	}
	if err := s.createLessonTx(ctx, tx, &models.Lesson{ModuleID: m1.ID, Title: "Promises & Async/Await", ContentType: "video", DurationMinutes: 20, OrderIndex: 1}); err != nil {
		return err
	}

	m2 := &models.Module{PathID: pathID, Title: "Vue 3 Fundamentals", Description: "Components and Reactivity", OrderIndex: 2}
	if err := s.createModuleTx(ctx, tx, m2); err != nil {
		return err
	}
	return s.createLessonTx(ctx, tx, &models.Lesson{ModuleID: m2.ID, Title: "Composition API", ContentType: "article", DurationMinutes: 25, OrderIndex: 1})
}

func (s *service) seedFrontendAdvancedModules(ctx context.Context, tx pgx.Tx, pathID string) error {
	m1 := &models.Module{PathID: pathID, Title: "Performance Optimization", Description: "Lazy loading, code splitting", OrderIndex: 1}
	if err := s.createModuleTx(ctx, tx, m1); err != nil {
		return err
	}
	return s.createLessonTx(ctx, tx, &models.Lesson{ModuleID: m1.ID, Title: "Web Vitals", ContentType: "video", DurationMinutes: 30, OrderIndex: 1})
}

func (s *service) seedBackendBeginnerModules(ctx context.Context, tx pgx.Tx, pathID string) error {
	m1 := &models.Module{PathID: pathID, Title: "Go Language Basics", Description: "Variables, Loops, Functions", OrderIndex: 1}
	if err := s.createModuleTx(ctx, tx, m1); err != nil {
		return err
	}
	if err := s.createLessonTx(ctx, tx, &models.Lesson{ModuleID: m1.ID, Title: "Installing Go", ContentType: "article", DurationMinutes: 10, OrderIndex: 1}); err != nil {
		return err
	}
	return s.createLessonTx(ctx, tx, &models.Lesson{ModuleID: m1.ID, Title: "Your First Program", ContentType: "video", DurationMinutes: 15, OrderIndex: 2})
}

func (s *service) seedBackendIntermediateModules(ctx context.Context, tx pgx.Tx, pathID string) error {
	m1 := &models.Module{PathID: pathID, Title: "Working with Databases", Description: "SQL and Drivers", OrderIndex: 1}
	if err := s.createModuleTx(ctx, tx, m1); err != nil {
		return err
	}
	return s.createLessonTx(ctx, tx, &models.Lesson{ModuleID: m1.ID, Title: "Connecting to PostgreSQL", ContentType: "video", DurationMinutes: 20, OrderIndex: 1})
}

func (s *service) seedBackendAdvancedModules(ctx context.Context, tx pgx.Tx, pathID string) error {
	m1 := &models.Module{PathID: pathID, Title: "Concurrency Patterns", Description: "Pipelines, Fan-out/Fan-in", OrderIndex: 1}
	if err := s.createModuleTx(ctx, tx, m1); err != nil {
		return err
	}
	return s.createLessonTx(ctx, tx, &models.Lesson{ModuleID: m1.ID, Title: "Advanced Channels", ContentType: "article", DurationMinutes: 35, OrderIndex: 1})
}

// createModuleTx creates a module within a transaction.
func (s *service) createModuleTx(ctx context.Context, tx pgx.Tx, module *models.Module) error {
	query := `
		INSERT INTO modules (path_id, title, description, order_index)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	err := tx.QueryRow(ctx, query, module.PathID, module.Title, module.Description, module.OrderIndex).Scan(&module.ID, &module.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

// createLessonTx creates a lesson within a transaction.
func (s *service) createLessonTx(ctx context.Context, tx pgx.Tx, lesson *models.Lesson) error {
	query := `
		INSERT INTO lessons (module_id, title, content_type, content_url, content_body, duration_minutes, order_index)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`
	err := tx.QueryRow(ctx, query, lesson.ModuleID, lesson.Title, lesson.ContentType, lesson.ContentURL, lesson.ContentBody, lesson.DurationMinutes, lesson.OrderIndex).Scan(&lesson.ID, &lesson.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
