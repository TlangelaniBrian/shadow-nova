# Structured Logging in Shadow Nova Backend

This document describes the structured logging implementation using Go's `log/slog` package.

## Overview

Shadow Nova uses structured logging to provide consistent, queryable, and machine-readable logs. All logs are output in JSON format with contextual fields for efficient debugging and monitoring.

## Benefits

1. **Machine-Readable**: JSON logs can be easily parsed and queried by log aggregation tools
2. **Contextual Information**: Each log entry includes relevant metadata (request_id, user_id, etc.)
3. **Consistent Format**: All logs follow the same structure across the application
4. **Performance**: slog is optimized for high-performance logging
5. **Integration Ready**: Works seamlessly with tools like Grafana Loki, ELK Stack, Datadog, etc.

## Configuration

Logging is configured via environment variables:

```env
# Logging configuration
LOG_LEVEL=info  # debug, info, warn, error
```

- `LOG_LEVEL`: Controls the minimum log level to output
  - `debug`: Most verbose, includes all logs
  - `info`: Informational messages (default)
  - `warn`: Warning messages
  - `error`: Error messages only

## Usage

### Basic Logging

```go
import "shadow-nova/backend/internal/logging"

// Info level
logging.Info("user login successful", "user_id", 123, "email", "user@example.com")

// Warning level
logging.Warn("failed to send email", "error", err, "user_id", 123)

// Error level
logging.Error("database query failed", err, "query", "SELECT * FROM users")

// Debug level (only shown when LOG_LEVEL=debug)
logging.Debug("cache hit", "key", "user:123", "ttl", 3600)
```

### Contextual Logging with Request ID

All HTTP requests automatically receive a unique request ID. Extract it for contextual logging:

```go
func (h *Handler) MyHandler(w http.ResponseWriter, r *http.Request) {
    log := logging.WithContext(r.Context()).With(
        "handler", "MyHandler",
        "user_id", middleware.GetUserID(r),
    )

    log.Info("processing request")

    // ... your logic

    if err != nil {
        log.Error("operation failed", "error", err)
        return
    }

    log.Info("request completed successfully")
}
```

### Adding Additional Context

Use `With()` to create a logger with additional context fields:

```go
// Create a logger with service-specific context
serviceLog := logging.With(
    "service", "collector",
    "source_id", sourceID,
)

serviceLog.Info("starting collection")
// ... processing
serviceLog.Info("collection complete", "items_collected", count)
```

## Log Structure

Each log entry includes:

- `time`: Timestamp in RFC3339 format
- `level`: Log level (DEBUG, INFO, WARN, ERROR)
- `msg`: Human-readable message
- `source`: File and line number where the log was generated
- `request_id`: Unique ID for HTTP requests (when applicable)
- Custom fields: Any additional contextual data

Example JSON log entry:

```json
{
  "time": "2026-02-12T10:30:45Z",
  "level": "INFO",
  "source": {
    "function": "shadow-nova/backend/internal/handlers.(*PathsHandler).Get",
    "file": "/app/internal/handlers/paths.go",
    "line": 42
  },
  "msg": "fetching learning path",
  "request_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "handler": "paths.Get",
  "path_id": "frontend-beginner"
}
```

## Best Practices

### 1. Use Structured Fields

**Don't:**
```go
logging.Info(fmt.Sprintf("User %d logged in from %s", userID, ip))
```

**Do:**
```go
logging.Info("user login", "user_id", userID, "ip_address", ip)
```

### 2. Log Actions, Not Data

**Don't:**
```go
logging.Info("user data", "user", userObject)
```

**Do:**
```go
logging.Info("user created successfully", "user_id", user.ID, "email", user.Email)
```

### 3. Use Appropriate Log Levels

- **Debug**: Detailed diagnostic information (disabled in production)
- **Info**: General informational messages (user actions, system events)
- **Warn**: Potentially harmful situations (retries, fallbacks, degraded functionality)
- **Error**: Error events that might still allow the application to continue

### 4. Include Error Context

**Don't:**
```go
logging.Error("database error", err)
```

**Do:**
```go
logging.Error("failed to create user", err, "email", email, "operation", "CreateUser")
```

### 5. Use Consistent Field Names

- `user_id` (not `userId` or `uid`)
- `request_id` (not `requestId` or `req_id`)
- `error` (for error messages)
- Use snake_case for field names

## Querying Logs

### With Grafana Loki

```logql
# All errors in the last hour
{job="shadow-nova-backend"} | json | level="ERROR"

# All requests to a specific handler
{job="shadow-nova-backend"} | json | handler="paths.Get"

# Trace a specific request
{job="shadow-nova-backend"} | json | request_id="f47ac10b-58cc-4372-a567-0e02b2c3d479"

# Find slow requests
{job="shadow-nova-backend"} | json | duration_ms > 1000
```

### With jq (command line)

```bash
# Filter errors only
cat logs.json | jq 'select(.level == "ERROR")'

# Get all logs for a specific user
cat logs.json | jq 'select(.user_id == 123)'

# Find logs with specific message
cat logs.json | jq 'select(.msg | contains("database"))'
```

## Integration with Log Aggregation Systems

### Grafana Loki

Shadow Nova's JSON logs are ready for Loki ingestion:

```yaml
# promtail-config.yml
server:
  http_listen_port: 9080

positions:
  filename: /tmp/positions.yaml

clients:
  - url: http://loki:3100/loki/api/v1/push

scrape_configs:
  - job_name: shadow-nova
    static_configs:
      - targets:
          - localhost
        labels:
          job: shadow-nova-backend
          __path__: /var/log/shadow-nova/*.log
    pipeline_stages:
      - json:
          expressions:
            level: level
            msg: msg
            request_id: request_id
```

### ELK Stack (Elasticsearch, Logstash, Kibana)

```ruby
# logstash.conf
input {
  file {
    path => "/var/log/shadow-nova/*.log"
    codec => json
  }
}

filter {
  json {
    source => "message"
  }
}

output {
  elasticsearch {
    hosts => ["localhost:9200"]
    index => "shadow-nova-%{+YYYY.MM.dd}"
  }
}
```

## Request Tracking

Every HTTP request receives a unique request ID that appears in:
- All logs generated during request processing
- The `X-Request-ID` response header

Use this ID to trace a request across your entire stack.

## Performance Considerations

1. **Minimal Overhead**: slog is optimized for performance
2. **Lazy Evaluation**: Log arguments are only evaluated if the log level is enabled
3. **No String Formatting**: Structured fields avoid expensive string concatenation
4. **Batch Processing**: Logs can be buffered for batch processing in high-throughput scenarios

## Migration Notes

All previous logging methods have been replaced:

- `log.Println()` → `logging.Info()`
- `log.Printf()` → `logging.Info()` with structured fields
- `fmt.Printf()` → `logging.Info()` with structured fields
- `log.Fatal()` → `logging.Error()` + `os.Exit(1)`

## Example Scenarios

### Database Operations

```go
logging.Info("executing database query",
    "query", "GetLearningPath",
    "path_id", pathID,
)

// On error
logging.Error("database query failed", err,
    "query", "GetLearningPath",
    "path_id", pathID,
)
```

### Background Jobs

```go
logging.Info("starting content collection",
    "source_count", len(sources),
    "scheduled_at", time.Now(),
)

// Progress updates
logging.Info("fetched items from source",
    "source_name", source.Name,
    "item_count", len(items),
)

// Completion
logging.Info("content collection complete",
    "duration_ms", time.Since(start).Milliseconds(),
    "items_processed", totalItems,
)
```

### Authentication

```go
logging.Info("user authentication attempt",
    "email", email,
    "provider", "google",
)

// On success
logging.Info("user authenticated successfully",
    "user_id", user.ID,
    "provider", "google",
)

// On failure
logging.Warn("authentication failed",
    "email", email,
    "provider", "google",
    "reason", "invalid_credentials",
)
```

## Troubleshooting

### Logs Not Appearing

1. Check `LOG_LEVEL` environment variable
2. Ensure `logging.Init()` is called early in `main()`
3. Verify logs are being written to stdout

### Too Verbose

Set `LOG_LEVEL=warn` or `LOG_LEVEL=error` to reduce verbosity.

### Missing Context

Ensure handlers use `logging.WithContext(r.Context())` to include request ID.

## Related Files

- `/internal/logging/logger.go` - Core logging implementation
- `/internal/middleware/request_id.go` - Request ID middleware
- `/internal/middleware/logging.go` - HTTP request logging middleware
- `/.env` - Logging configuration

## References

- [Go slog Package](https://pkg.go.dev/log/slog)
- [Grafana Loki Documentation](https://grafana.com/docs/loki/)
- [Structured Logging Best Practices](https://www.honeycomb.io/blog/structured-logging-and-your-team)
