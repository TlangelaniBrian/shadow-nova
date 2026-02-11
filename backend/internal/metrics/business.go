package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	UserRegistrations = promauto.NewCounter(prometheus.CounterOpts{
		Name: "user_registrations_total",
		Help: "Total number of user registrations",
	})

	UserLogins = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "user_logins_total",
		Help: "Total number of user logins by method",
	}, []string{"method"}) // method: email, google

	LessonCompletions = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "lesson_completions_total",
		Help: "Total lessons completed by content type",
	}, []string{"content_type"}) // content_type: video, article, quiz

	PathCompletions = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "learning_path_completions_total",
		Help: "Learning paths completed by difficulty",
	}, []string{"difficulty"}) // difficulty: beginner, intermediate, advanced

	ProjectSubmissions = promauto.NewCounter(prometheus.CounterOpts{
		Name: "project_submissions_total",
		Help: "Total project submissions",
	})

	ContentItemsCollected = promauto.NewCounter(prometheus.CounterOpts{
		Name: "content_items_collected_total",
		Help: "Total content items collected from sources",
	})

	AIProcessingDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "ai_processing_duration_seconds",
		Help: "Time to process content with AI",
		Buckets: []float64{0.1, 0.5, 1, 2, 5, 10, 30},
	})

	AIProcessingErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ai_processing_errors_total",
		Help: "Total AI processing errors",
	})

	ProjectStatusUpdates = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "project_status_updates_total",
		Help: "Total project status updates by new status",
	}, []string{"status"}) // status: pending, approved, rejected
)
