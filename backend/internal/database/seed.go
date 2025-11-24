package database

import (
	"context"
	"log"
	"shadow-nova/backend/internal/models"
)

func (s *service) SeedLearningPaths(ctx context.Context) error {
	// Check if paths already exist
	paths, err := s.GetLearningPaths(ctx)
	if err != nil {
		return err
	}
	if len(paths) > 0 {
		log.Println("Learning paths already exist, skipping seeding.")
		return nil
	}

	log.Println("Seeding learning paths...")

	// --- FRONTEND PATHS ---

	// 1. Frontend Beginner
	feBeginner := &models.LearningPath{
		ID:          "frontend-beginner",
		Title:       "Frontend Development Basics",
		Description: "Start your journey here. Learn HTML, CSS, and how the web works.",
		Difficulty:  "Beginner",
	}
	if err := s.CreateLearningPath(ctx, feBeginner); err != nil {
		return err
	}
	s.seedFrontendBeginnerModules(ctx, feBeginner.ID)

	// 2. Frontend Intermediate
	feInter := &models.LearningPath{
		ID:          "frontend-intermediate",
		Title:       "Frontend Mastery with Vue",
		Description: "Deep dive into modern JavaScript and the Vue.js framework.",
		Difficulty:  "Intermediate",
	}
	if err := s.CreateLearningPath(ctx, feInter); err != nil {
		return err
	}
	s.seedFrontendIntermediateModules(ctx, feInter.ID)

	// 3. Frontend Advanced
	feAdv := &models.LearningPath{
		ID:          "frontend-advanced",
		Title:       "Advanced Frontend Architecture",
		Description: "Performance, testing, and large-scale application design.",
		Difficulty:  "Advanced",
	}
	if err := s.CreateLearningPath(ctx, feAdv); err != nil {
		return err
	}
	s.seedFrontendAdvancedModules(ctx, feAdv.ID)

	// --- BACKEND PATHS ---

	// 4. Backend Beginner
	beBeginner := &models.LearningPath{
		ID:          "backend-beginner",
		Title:       "Backend Basics with Go",
		Description: "Introduction to server-side programming, HTTP, and Go syntax.",
		Difficulty:  "Beginner",
	}
	if err := s.CreateLearningPath(ctx, beBeginner); err != nil {
		return err
	}
	s.seedBackendBeginnerModules(ctx, beBeginner.ID)

	// 5. Backend Intermediate
	beInter := &models.LearningPath{
		ID:          "backend-intermediate",
		Title:       "Building REST APIs",
		Description: "Database integration, authentication, and API design patterns.",
		Difficulty:  "Intermediate",
	}
	if err := s.CreateLearningPath(ctx, beInter); err != nil {
		return err
	}
	s.seedBackendIntermediateModules(ctx, beInter.ID)

	// 6. Backend Advanced
	beAdv := &models.LearningPath{
		ID:          "backend-advanced",
		Title:       "Distributed Systems & Microservices",
		Description: "Scaling applications, concurrency, and cloud deployment.",
		Difficulty:  "Advanced",
	}
	if err := s.CreateLearningPath(ctx, beAdv); err != nil {
		return err
	}
	s.seedBackendAdvancedModules(ctx, beAdv.ID)

	log.Println("Learning paths seeded successfully!")
	return nil
}

// --- Helper Functions ---

func (s *service) seedFrontendBeginnerModules(ctx context.Context, pathID string) {
	m1 := &models.Module{PathID: pathID, Title: "The Web Ecosystem", Description: "How browsers and servers talk", OrderIndex: 1}
	s.CreateModule(ctx, m1)
	s.CreateLesson(ctx, &models.Lesson{ModuleID: m1.ID, Title: "HTTP & DNS", ContentType: "video", DurationMinutes: 10, OrderIndex: 1})

	m2 := &models.Module{PathID: pathID, Title: "HTML5 Essentials", Description: "Structuring content", OrderIndex: 2}
	s.CreateModule(ctx, m2)
	s.CreateLesson(ctx, &models.Lesson{ModuleID: m2.ID, Title: "Semantic HTML", ContentType: "article", DurationMinutes: 15, OrderIndex: 1})
}

func (s *service) seedFrontendIntermediateModules(ctx context.Context, pathID string) {
	m1 := &models.Module{PathID: pathID, Title: "JavaScript Deep Dive", Description: "ES6+, Async/Await", OrderIndex: 1}
	s.CreateModule(ctx, m1)
	s.CreateLesson(ctx, &models.Lesson{ModuleID: m1.ID, Title: "Promises & Async/Await", ContentType: "video", DurationMinutes: 20, OrderIndex: 1})

	m2 := &models.Module{PathID: pathID, Title: "Vue 3 Fundamentals", Description: "Components and Reactivity", OrderIndex: 2}
	s.CreateModule(ctx, m2)
	s.CreateLesson(ctx, &models.Lesson{ModuleID: m2.ID, Title: "Composition API", ContentType: "article", DurationMinutes: 25, OrderIndex: 1})
}

func (s *service) seedFrontendAdvancedModules(ctx context.Context, pathID string) {
	m1 := &models.Module{PathID: pathID, Title: "Performance Optimization", Description: "Lazy loading, code splitting", OrderIndex: 1}
	s.CreateModule(ctx, m1)
	s.CreateLesson(ctx, &models.Lesson{ModuleID: m1.ID, Title: "Web Vitals", ContentType: "video", DurationMinutes: 30, OrderIndex: 1})
}

func (s *service) seedBackendBeginnerModules(ctx context.Context, pathID string) {
	m1 := &models.Module{PathID: pathID, Title: "Go Language Basics", Description: "Variables, Loops, Functions", OrderIndex: 1}
	s.CreateModule(ctx, m1)
	s.CreateLesson(ctx, &models.Lesson{ModuleID: m1.ID, Title: "Installing Go", ContentType: "article", DurationMinutes: 10, OrderIndex: 1})
	s.CreateLesson(ctx, &models.Lesson{ModuleID: m1.ID, Title: "Your First Program", ContentType: "video", DurationMinutes: 15, OrderIndex: 2})
}

func (s *service) seedBackendIntermediateModules(ctx context.Context, pathID string) {
	m1 := &models.Module{PathID: pathID, Title: "Working with Databases", Description: "SQL and Drivers", OrderIndex: 1}
	s.CreateModule(ctx, m1)
	s.CreateLesson(ctx, &models.Lesson{ModuleID: m1.ID, Title: "Connecting to PostgreSQL", ContentType: "video", DurationMinutes: 20, OrderIndex: 1})
}

func (s *service) seedBackendAdvancedModules(ctx context.Context, pathID string) {
	m1 := &models.Module{PathID: pathID, Title: "Concurrency Patterns", Description: "Pipelines, Fan-out/Fan-in", OrderIndex: 1}
	s.CreateModule(ctx, m1)
	s.CreateLesson(ctx, &models.Lesson{ModuleID: m1.ID, Title: "Advanced Channels", ContentType: "article", DurationMinutes: 35, OrderIndex: 1})
}
