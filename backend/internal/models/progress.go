package models

import "time"

type UserProgress struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	LessonID    int       `json:"lesson_id"`
	Completed   bool      `json:"completed"`
	CompletedAt time.Time `json:"completed_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type UpdateProgressRequest struct {
	LessonID  int  `json:"lesson_id" validate:"required"`
	Completed bool `json:"completed"`
}

type UserStats struct {
	CoursesCompleted int    `json:"courses_completed"`
	HoursLearned     int    `json:"hours_learned"`
	Rank             int    `json:"rank"`
	CurrentStreak    int    `json:"current_streak"`
	TotalXP          int    `json:"total_xp"`
}

type PathProgress struct {
	PathID          string  `json:"path_id"`
	TotalLessons    int     `json:"total_lessons"`
	CompletedLessons int    `json:"completed_lessons"`
	Percentage      float64 `json:"percentage"`
}
