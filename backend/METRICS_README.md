# Shadow Nova Business Metrics - Quick Start

## Overview
Shadow Nova now tracks comprehensive business metrics for monitoring user engagement, learning progress, project activity, content collection, and AI processing performance.

## Quick Access

### Metrics Endpoint
```
http://localhost:8080/metrics
```

### Grafana Dashboard
Import from: `frontend/observability/dashboards/business-metrics.json`

### Documentation
- Full docs: `backend/BUSINESS_METRICS.md`
- Implementation details: `backend/BUSINESS_METRICS_IMPLEMENTATION.md`

## Available Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `user_registrations_total` | Counter | Total user registrations |
| `user_logins_total{method}` | Counter | Logins by auth method (email/google) |
| `lesson_completions_total{content_type}` | Counter | Lessons completed by type |
| `learning_path_completions_total{difficulty}` | Counter | Paths completed by difficulty |
| `project_submissions_total` | Counter | Total project submissions |
| `project_status_updates_total{status}` | Counter | Project status changes |
| `content_items_collected_total` | Counter | Content items collected |
| `ai_processing_duration_seconds` | Histogram | AI processing latency |
| `ai_processing_errors_total` | Counter | AI processing errors |

## Quick Queries

### User Metrics
```promql
# Daily Active Users
sum(increase(user_logins_total[24h]))

# Registration rate
rate(user_registrations_total[1h]) * 3600

# Login methods breakdown
sum(user_logins_total) by (method)
```

### Learning Metrics
```promql
# Lesson completion rate
rate(lesson_completions_total[5m]) * 300

# Most popular content type
topk(1, sum(lesson_completions_total) by (content_type))
```

### Project Metrics
```promql
# Project approval rate
sum(project_status_updates_total{status="approved"}) / sum(project_status_updates_total)

# Submission rate
rate(project_submissions_total[1h])
```

### AI Performance
```promql
# P95 latency
histogram_quantile(0.95, rate(ai_processing_duration_seconds_bucket[5m]))

# Error rate
rate(ai_processing_errors_total[5m]) * 300

# Success rate
(1 - (rate(ai_processing_errors_total[5m]) / rate(ai_processing_duration_seconds_count[5m]))) * 100
```

## Verification

Run the verification script:
```bash
./scripts/verify-metrics.sh
```

Or manually check:
```bash
# Check metrics endpoint
curl http://localhost:8080/metrics | grep -E "(user_|lesson_|project_|content_|ai_)"

# Test a metric increment
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","username":"testuser","password":"password123"}'

# Verify the increment
curl http://localhost:8080/metrics | grep user_registrations_total
```

## Grafana Setup

1. Start Grafana:
   ```bash
   docker-compose up -d grafana
   ```

2. Access Grafana at `http://localhost:3000`

3. Import dashboard:
   - Navigate to Dashboards → Import
   - Upload `frontend/observability/dashboards/business-metrics.json`
   - Select Prometheus data source
   - Click Import

4. View metrics in real-time!

## Alert Configuration

Add to `frontend/observability/alerting_rules.yml`:

```yaml
groups:
  - name: business_metrics
    interval: 1m
    rules:
      - alert: HighAIProcessingErrorRate
        expr: rate(ai_processing_errors_total[5m]) * 300 > 5
        for: 5m
        labels:
          severity: critical

      - alert: NoUserRegistrations
        expr: increase(user_registrations_total[24h]) == 0
        for: 1h
        labels:
          severity: warning
```

## Adding New Metrics

1. Define metric in `backend/internal/metrics/business.go`:
   ```go
   var NewMetric = promauto.NewCounter(prometheus.CounterOpts{
       Name: "new_metric_total",
       Help: "Description",
   })
   ```

2. Import and use in handler:
   ```go
   import "shadow-nova/backend/internal/metrics"

   func Handler(w http.ResponseWriter, r *http.Request) {
       metrics.NewMetric.Inc()
   }
   ```

3. Update Grafana dashboard with new panel

## Key Business Insights

### Activation Rate
```promql
sum(lesson_completions_total) / user_registrations_total
```

### Content Preference
```promql
sum(lesson_completions_total) by (content_type)
```

### AI Efficiency
```promql
rate(ai_processing_duration_seconds_count[1h]) * 3600
```

### Review Throughput
```promql
rate(project_status_updates_total[1h]) * 3600
```

## Troubleshooting

### Metrics not showing?
1. Check if backend is running: `curl http://localhost:8080/health`
2. Verify metrics endpoint: `curl http://localhost:8080/metrics`
3. Check Prometheus targets: `http://localhost:9090/targets`

### Dashboard not loading data?
1. Verify data source in Grafana
2. Check time range includes metric data
3. Validate PromQL query syntax

### Incorrect values?
1. Verify metric is incremented in correct location
2. Check for errors in application logs
3. Ensure metric type matches usage (Counter/Gauge/Histogram)

## Files Changed

**Created:**
- `backend/internal/metrics/business.go`
- `frontend/observability/dashboards/business-metrics.json`
- `backend/BUSINESS_METRICS.md`
- `scripts/verify-metrics.sh`

**Modified:**
- `backend/internal/handlers/auth.go`
- `backend/internal/handlers/progress.go`
- `backend/internal/handlers/projects.go`
- `backend/internal/collector/service.go`
- `backend/internal/database/database.go`
- `backend/internal/database/paths.go`

## Support

- Full documentation: `backend/BUSINESS_METRICS.md`
- Prometheus docs: https://prometheus.io/docs/
- Grafana docs: https://grafana.com/docs/

## Status

✅ All metrics implemented
✅ Code compiles successfully
✅ Dashboard created
✅ Documentation complete
✅ Ready for deployment
