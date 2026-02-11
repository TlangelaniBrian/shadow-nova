# Business Metrics Documentation

This document outlines the business metrics implemented in Shadow Nova for tracking key user actions and system performance.

## Overview

Shadow Nova implements comprehensive business metrics using Prometheus for tracking user engagement, learning progress, project activity, content collection, and AI processing performance. These metrics are visualized in Grafana dashboards for real-time monitoring and business insights.

## Metrics Categories

### 1. User Acquisition & Authentication

#### `user_registrations_total`
- **Type**: Counter
- **Description**: Total number of user registrations
- **Location**: `backend/internal/handlers/auth.go` - `Register()` method
- **Business Value**: Track user growth and acquisition effectiveness

#### `user_logins_total`
- **Type**: Counter (with labels)
- **Labels**: `method` (email, google)
- **Description**: Total number of user logins by authentication method
- **Locations**:
  - `backend/internal/handlers/auth.go` - `Login()` method (email)
  - `backend/internal/handlers/auth.go` - `GoogleCallback()` method (google)
  - `backend/internal/handlers/auth.go` - `VerifyGoogleToken()` method (google)
- **Business Value**: Understand authentication preferences and detect authentication issues

**Example Queries**:
```promql
# Daily active users
sum(increase(user_logins_total[24h]))

# Login methods breakdown
sum(user_logins_total) by (method)

# Login rate per hour
rate(user_logins_total[1h]) * 3600
```

### 2. Learning Engagement

#### `lesson_completions_total`
- **Type**: Counter (with labels)
- **Labels**: `content_type` (video, article, quiz)
- **Description**: Total lessons completed by content type
- **Location**: `backend/internal/handlers/progress.go` - `UpdateProgress()` method
- **Business Value**: Measure learning engagement and content effectiveness

#### `learning_path_completions_total`
- **Type**: Counter (with labels)
- **Labels**: `difficulty` (beginner, intermediate, advanced)
- **Description**: Learning paths completed by difficulty level
- **Location**: To be implemented when path completion tracking is added
- **Business Value**: Track curriculum completion and identify difficulty bottlenecks

**Example Queries**:
```promql
# Lesson completion rate by type
rate(lesson_completions_total[5m]) * 300

# Most completed content type
topk(1, sum(lesson_completions_total) by (content_type))

# Daily lesson completions
sum(increase(lesson_completions_total[24h]))
```

### 3. Project Activity

#### `project_submissions_total`
- **Type**: Counter
- **Description**: Total project submissions
- **Location**: `backend/internal/handlers/projects.go` - `Submit()` method
- **Business Value**: Track practical application of skills and engagement with projects

#### `project_status_updates_total`
- **Type**: Counter (with labels)
- **Labels**: `status` (pending, approved, rejected)
- **Description**: Total project status updates by new status
- **Location**: `backend/internal/handlers/projects.go` - `UpdateSubmission()` method
- **Business Value**: Monitor review workflow and project quality

**Example Queries**:
```promql
# Project submission rate
rate(project_submissions_total[5m]) * 300

# Project approval ratio
sum(project_status_updates_total{status="approved"}) / sum(project_status_updates_total)

# Pending vs reviewed projects
sum(project_status_updates_total) by (status)
```

### 4. Content Collection

#### `content_items_collected_total`
- **Type**: Counter
- **Description**: Total content items collected from external sources
- **Location**: `backend/internal/collector/service.go` - `processSource()` method
- **Business Value**: Track content curation effectiveness and source reliability

**Example Queries**:
```promql
# Total items collected
content_items_collected_total

# Collection rate per hour
rate(content_items_collected_total[1h]) * 3600
```

### 5. AI Processing Performance

#### `ai_processing_duration_seconds`
- **Type**: Histogram
- **Buckets**: [0.1, 0.5, 1, 2, 5, 10, 30]
- **Description**: Time to process content with AI (in seconds)
- **Location**: `backend/internal/collector/service.go` - `ProcessUnprocessedItems()` method
- **Business Value**: Monitor AI API performance and cost optimization

#### `ai_processing_errors_total`
- **Type**: Counter
- **Description**: Total AI processing errors
- **Location**: `backend/internal/collector/service.go` - `ProcessUnprocessedItems()` method
- **Business Value**: Track AI service reliability and identify issues

**Example Queries**:
```promql
# AI processing p50, p95, p99 latency
histogram_quantile(0.50, rate(ai_processing_duration_seconds_bucket[5m]))
histogram_quantile(0.95, rate(ai_processing_duration_seconds_bucket[5m]))
histogram_quantile(0.99, rate(ai_processing_duration_seconds_bucket[5m]))

# AI error rate
rate(ai_processing_errors_total[5m]) * 300

# AI success rate
(1 - (rate(ai_processing_errors_total[5m]) / rate(ai_processing_duration_seconds_count[5m]))) * 100
```

## Grafana Dashboard

The business metrics dashboard is located at:
`frontend/observability/dashboards/business-metrics.json`

### Dashboard Sections

1. **User Acquisition**
   - User registrations graph (time series)
   - Login methods pie chart
   - Daily active users stat
   - Total registrations stat

2. **Learning Engagement**
   - Lesson completions by content type (time series)
   - Learning path completions by difficulty
   - Total lesson completions stat

3. **Project Activity**
   - Project submissions rate graph
   - Project status updates breakdown
   - Total submissions stat

4. **Content & AI**
   - Content items collected stat
   - AI processing duration percentiles (p50, p95, p99)
   - AI processing error rate with alerting
   - AI error rate graph

### Importing the Dashboard

1. Access Grafana at `http://localhost:3000`
2. Navigate to Dashboards → Import
3. Upload `business-metrics.json` or paste its contents
4. Select Prometheus as the data source
5. Click Import

## Alerting Recommendations

### Critical Alerts

1. **High AI Processing Error Rate**
   - **Condition**: Error rate > 5 errors per 5 minutes
   - **Action**: Check AI service status, API keys, and rate limits
   - **Impact**: Content processing delays

2. **Zero Registrations for 24h**
   - **Query**: `increase(user_registrations_total[24h]) == 0`
   - **Action**: Check registration endpoint, email verification, and marketing campaigns
   - **Impact**: Growth stagnation

3. **High AI Processing Latency**
   - **Condition**: p99 > 30 seconds
   - **Action**: Review AI service performance and consider caching
   - **Impact**: Poor user experience for content discovery

### Warning Alerts

1. **Low Lesson Completion Rate**
   - **Query**: `rate(lesson_completions_total[1h]) < 10`
   - **Action**: Review content difficulty and user engagement
   - **Impact**: Reduced learning effectiveness

2. **Project Submission Drop**
   - **Query**: `rate(project_submissions_total[24h]) < rate(project_submissions_total[24h] offset 7d) * 0.5`
   - **Action**: Check project requirements and user feedback
   - **Impact**: Lower practical skill application

## Business Insights to Derive

### User Engagement Metrics

1. **Activation Rate**
   - Formula: `lesson_completions / user_registrations`
   - Measures: Users who engage with content after registration

2. **Retention Rate**
   - Formula: `unique_logins_this_week / unique_logins_last_week`
   - Measures: User retention over time

3. **Completion Rate**
   - Formula: `learning_path_completions / learning_path_starts`
   - Measures: Path completion effectiveness

### Content Performance

1. **Content Preference**
   - Query: `sum(lesson_completions_total) by (content_type)`
   - Insight: Which content types users prefer

2. **Difficulty Distribution**
   - Query: `sum(learning_path_completions_total) by (difficulty)`
   - Insight: User skill level distribution

3. **Content Quality Score**
   - Formula: `(approved_projects / submitted_projects) * completion_rate`
   - Measures: Content effectiveness

### Operational Efficiency

1. **AI Processing Efficiency**
   - Metrics: AI latency percentiles, error rate, throughput
   - Cost: Processing time * API rate
   - Optimization: Batch processing, caching strategies

2. **Content Curation ROI**
   - Formula: `lesson_completions_from_curated_content / content_items_collected`
   - Measures: Value of content curation efforts

3. **Review Throughput**
   - Query: `rate(project_status_updates_total[1h])`
   - Measures: Project review efficiency

## Integration with Application Code

### Adding New Metrics

1. Define the metric in `backend/internal/metrics/business.go`:
```go
var NewMetric = promauto.NewCounter(prometheus.CounterOpts{
    Name: "new_metric_total",
    Help: "Description of the metric",
})
```

2. Import metrics package in your handler:
```go
import "shadow-nova/backend/internal/metrics"
```

3. Increment the metric at the appropriate location:
```go
metrics.NewMetric.Inc()
```

### Metric Naming Conventions

Follow Prometheus naming best practices:
- Use `_total` suffix for counters
- Use `_seconds` suffix for durations
- Use snake_case for metric names
- Use descriptive labels for dimensional data
- Avoid high cardinality labels (e.g., user IDs, timestamps)

## Performance Considerations

1. **Metric Collection Overhead**
   - Prometheus metrics are extremely lightweight
   - Counter increments are ~10-20ns
   - Histogram observations are ~100-200ns

2. **Storage**
   - Metrics are stored in Prometheus TSDB
   - Default retention: 15 days
   - Configure retention in `prometheus.yml`

3. **Query Performance**
   - Use recording rules for complex queries
   - Pre-aggregate frequently used metrics
   - Avoid high cardinality in labels

## Testing Metrics

### Local Testing

1. Start the application with metrics enabled
2. Access Prometheus at `http://localhost:9090`
3. Query metrics using PromQL
4. Verify metric values after triggering actions

### Example Test Scenarios

```bash
# Test user registration
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","username":"testuser","password":"password123"}'

# Verify registration metric
curl http://localhost:9090/api/v1/query?query=user_registrations_total

# Test lesson completion
curl -X POST http://localhost:8080/api/progress \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"lesson_id":1,"completed":true}'

# Verify completion metric
curl http://localhost:9090/api/v1/query?query=lesson_completions_total
```

## Troubleshooting

### Metrics Not Appearing

1. Check if Prometheus is scraping the application:
   - Visit `http://localhost:9090/targets`
   - Ensure application target is UP

2. Verify metric endpoint:
   - `curl http://localhost:8080/metrics`
   - Search for your metric name

3. Check metric registration:
   - Ensure `promauto.New*` is used (auto-registers)
   - Verify import path is correct

### Incorrect Metric Values

1. Check metric type (Counter vs Gauge vs Histogram)
2. Verify metric is incremented in correct location
3. Check for errors in application logs
4. Ensure no race conditions on metric updates

### Dashboard Not Showing Data

1. Verify data source configuration in Grafana
2. Check PromQL query syntax
3. Ensure time range includes metric data
4. Check Prometheus retention settings

## Future Enhancements

1. **User Segmentation**
   - Add labels for user cohorts
   - Track feature adoption by segment

2. **Revenue Metrics**
   - Track subscription conversions
   - Monitor payment success rates

3. **Performance Metrics**
   - API endpoint latencies
   - Database query performance

4. **Engagement Scores**
   - Combined metrics for user engagement
   - Predictive churn indicators

## References

- [Prometheus Documentation](https://prometheus.io/docs/)
- [Grafana Dashboards](https://grafana.com/docs/grafana/latest/dashboards/)
- [Prometheus Best Practices](https://prometheus.io/docs/practices/naming/)
- [PromQL Cheat Sheet](https://promlabs.com/promql-cheat-sheet/)
