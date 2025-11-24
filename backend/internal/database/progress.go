package database

import (
	"context"
	"fmt"
	"shadow-nova/backend/internal/models"
	"time"
)

func (s *service) UpdateUserProgress(ctx context.Context, userID int, req models.UpdateProgressRequest) error {
	query := `
		INSERT INTO user_progress (user_id, lesson_id, completed, completed_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, lesson_id) 
		DO UPDATE SET completed = $3, completed_at = $4
	`

	var completedAt time.Time
	if req.Completed {
		completedAt = time.Now()
	}

	_, err := s.db.Exec(ctx, query, userID, req.LessonID, req.Completed, completedAt)
	if err != nil {
		return fmt.Errorf("failed to update user progress: %w", err)
	}

	return nil
}

func (s *service) GetUserStats(ctx context.Context, userID int) (*models.UserStats, error) {
	stats := &models.UserStats{
		Rank: 42, // Placeholder logic for rank
		CurrentStreak: 5, // Placeholder logic for streak
	}

	// Calculate completed courses (paths where all lessons are completed)
	// This is a complex query, for now let's just count completed lessons as a proxy or keep it simple
	// Real implementation would check if all lessons in a path are in user_progress with completed=true
	
	// Count total completed lessons
	var completedLessons int
	err := s.db.QueryRow(ctx, "SELECT COUNT(*) FROM user_progress WHERE user_id = $1 AND completed = true", userID).Scan(&completedLessons)
	if err != nil {
		return nil, fmt.Errorf("failed to count completed lessons: %w", err)
	}

	// Estimate hours learned (sum of duration of completed lessons)
	queryHours := `
		SELECT COALESCE(SUM(l.duration_minutes), 0)
		FROM user_progress up
		JOIN lessons l ON up.lesson_id = l.id
		WHERE up.user_id = $1 AND up.completed = true
	`
	var totalMinutes int
	if err := s.db.QueryRow(ctx, queryHours, userID).Scan(&totalMinutes); err != nil {
		return nil, fmt.Errorf("failed to calculate hours learned: %w", err)
	}
	stats.HoursLearned = totalMinutes / 60
	stats.TotalXP = completedLessons * 100 // Simple XP formula

	// For courses completed, let's just return a mock number or 0 for now until we have a proper query
	stats.CoursesCompleted = 0 

	return stats, nil
}

func (s *service) GetPathProgress(ctx context.Context, userID int, pathID string) (*models.PathProgress, error) {
	// Get total lessons in path
	var totalLessons int
	queryTotal := `
		SELECT COUNT(*)
		FROM lessons l
		JOIN modules m ON l.module_id = m.id
		WHERE m.path_id = $1
	`
	if err := s.db.QueryRow(ctx, queryTotal, pathID).Scan(&totalLessons); err != nil {
		return nil, fmt.Errorf("failed to count total lessons: %w", err)
	}

	if totalLessons == 0 {
		return &models.PathProgress{PathID: pathID, Percentage: 0}, nil
	}

	// Get completed lessons in path
	var completedLessons int
	queryCompleted := `
		SELECT COUNT(*)
		FROM user_progress up
		JOIN lessons l ON up.lesson_id = l.id
		JOIN modules m ON l.module_id = m.id
		WHERE up.user_id = $1 AND m.path_id = $2 AND up.completed = true
	`
	if err := s.db.QueryRow(ctx, queryCompleted, userID, pathID).Scan(&completedLessons); err != nil {
		return nil, fmt.Errorf("failed to count completed lessons: %w", err)
	}

	percentage := (float64(completedLessons) / float64(totalLessons)) * 100

	return &models.PathProgress{
		PathID:           pathID,
		TotalLessons:     totalLessons,
		CompletedLessons: completedLessons,
		Percentage:       percentage,
	}, nil
}
