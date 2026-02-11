package database

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	dbConnectionsActive = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "db_connections_active",
		Help: "Number of active database connections",
	})

	dbConnectionsIdle = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "db_connections_idle",
		Help: "Number of idle database connections in the pool",
	})

	dbConnectionsMax = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "db_connections_max",
		Help: "Maximum number of database connections allowed",
	})

	dbConnectionsTotal = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "db_connections_total",
		Help: "Total number of database connections in the pool",
	})

	dbConnectionsAcquireCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "db_connections_acquire_total",
		Help: "Total number of successful connection acquisitions",
	})

	dbConnectionsAcquireDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "db_connections_acquire_duration_seconds",
		Help:    "Histogram of connection acquisition durations",
		Buckets: prometheus.DefBuckets,
	})

	dbConnectionsCanceledAcquireCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "db_connections_canceled_acquire_total",
		Help: "Total number of canceled connection acquisitions",
	})

	dbConnectionsConstructingCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "db_connections_constructing",
		Help: "Number of connections currently being constructed",
	})

	dbConnectionsNewCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "db_connections_new_total",
		Help: "Total number of new connections opened",
	})
)

// StartMetricsCollection starts collecting connection pool metrics
func (s *service) StartMetricsCollection(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				stats := s.db.Stat()

				// Update gauge metrics
				dbConnectionsActive.Set(float64(stats.AcquiredConns()))
				dbConnectionsIdle.Set(float64(stats.IdleConns()))
				dbConnectionsMax.Set(float64(stats.MaxConns()))
				dbConnectionsTotal.Set(float64(stats.TotalConns()))
				dbConnectionsConstructingCount.Set(float64(stats.ConstructingConns()))

				// Update counter metrics
				dbConnectionsAcquireCount.Add(float64(stats.AcquireCount()))
				dbConnectionsCanceledAcquireCount.Add(float64(stats.CanceledAcquireCount()))
				dbConnectionsNewCount.Add(float64(stats.NewConnsCount()))

				// Update histogram for acquire duration
				dbConnectionsAcquireDuration.Observe(stats.AcquireDuration().Seconds())

			case <-ctx.Done():
				return
			}
		}
	}()
}
