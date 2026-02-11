# Idempotency in Shadow Nova Backend

## Overview

Idempotency ensures that performing the same operation multiple times has the same effect as performing it once. This is critical for preventing duplicate submissions, double charges, or inconsistent state when network issues cause request retries.

## What is Idempotency?

In distributed systems, network failures can cause clients to retry requests. Without idempotency protection:
- A user might submit the same project twice
- Progress updates could be duplicated
- Learning path enrollments might be processed multiple times

With idempotency, the server recognizes duplicate requests and returns the cached response from the first successful execution.

## How It Works

### Architecture

1. **Client Generates Key**: Client generates a unique `Idempotency-Key` (typically a UUID) and includes it in the request header
2. **Server Checks Cache**: Server checks if this key has been seen before for this user
3. **Cache Hit**: If found, server returns the cached response without re-executing the operation
4. **Cache Miss**: If not found, server processes the request and caches the response
4. **Expiration**: Cached responses expire after 24 hours and are cleaned up automatically

### Supported Operations

Idempotency is automatically applied to all protected routes with these HTTP methods:
- `POST` - Creating new resources
- `PUT` - Updating entire resources
- `PATCH` - Partial updates

The following endpoints benefit from idempotency:
- `/api/progress` - Progress updates
- `/api/submissions` - Project submissions
- `/api/paths` - Learning path creation
- `/api/auth/logout` - Logout operations
- All other mutating operations on protected routes

## Usage

### Client-Side Implementation

#### Basic Usage

```typescript
import { v4 as uuidv4 } from 'uuid';

// Generate a unique idempotency key
const idempotencyKey = uuidv4();

// Include it in the request headers
fetch('/api/submissions', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer <token>',
    'Idempotency-Key': idempotencyKey
  },
  body: JSON.stringify({
    project_id: 'web-fundamentals',
    github_repo_url: 'https://github.com/user/project'
  })
});
```

#### Axios Integration

The frontend API client automatically adds idempotency keys to all mutating requests:

```typescript
// In frontend/src/api/client.ts
apiClient.interceptors.request.use((config) => {
  if (['post', 'put', 'patch'].includes(config.method?.toLowerCase() || '')) {
    if (!config.headers['Idempotency-Key']) {
      config.headers['Idempotency-Key'] = uuidv4();
    }
  }
  return config;
});
```

#### Retry Logic with Same Key

When implementing retry logic, use the same idempotency key:

```typescript
async function submitProjectWithRetry(data, maxRetries = 3) {
  const idempotencyKey = uuidv4();

  for (let attempt = 0; attempt < maxRetries; attempt++) {
    try {
      const response = await apiClient.post('/submissions', data, {
        headers: { 'Idempotency-Key': idempotencyKey }
      });

      // Check if this was a replayed response
      if (response.headers['x-idempotent-replay'] === 'true') {
        console.log('Received cached response from previous attempt');
      }

      return response.data;
    } catch (error) {
      if (attempt === maxRetries - 1) throw error;
      await sleep(Math.pow(2, attempt) * 1000); // Exponential backoff
    }
  }
}
```

### Server-Side Behavior

#### Response Headers

When a cached response is returned, the server adds this header:
```
X-Idempotent-Replay: true
```

This allows clients to distinguish between a fresh response and a replayed one.

#### Key Uniqueness

Idempotency keys are scoped to individual users. Two users can use the same key without conflicts.

```go
// Keys are unique per (key, user_id) combination
CREATE TABLE idempotency_keys (
    key VARCHAR(100) PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    ...
);
```

#### Storage and Expiration

- Responses are cached for 24 hours
- A daily cleanup job removes expired entries
- Storage uses the `idempotency_keys` table in PostgreSQL

## Testing Idempotency

### Manual Testing

```bash
# First request - processes normally
curl -X POST http://localhost:8080/api/progress \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: 550e8400-e29b-41d4-a716-446655440000" \
  -d '{"lesson_id": 1, "completed": true}'

# Second request with same key - returns cached response
curl -X POST http://localhost:8080/api/progress \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: 550e8400-e29b-41d4-a716-446655440000" \
  -d '{"lesson_id": 1, "completed": true}'

# Response will include: X-Idempotent-Replay: true
```

### Integration Testing

```go
func TestIdempotency(t *testing.T) {
    // Setup test server and user
    token := loginTestUser(t)
    idempotencyKey := "test-key-12345"

    // First request
    resp1 := makeRequest(t, "POST", "/api/progress",
        map[string]string{
            "Authorization": "Bearer " + token,
            "Idempotency-Key": idempotencyKey,
        },
        map[string]interface{}{
            "lesson_id": 1,
            "completed": true,
        },
    )
    assert.Equal(t, 200, resp1.StatusCode)
    assert.Empty(t, resp1.Header.Get("X-Idempotent-Replay"))

    // Second request with same key
    resp2 := makeRequest(t, "POST", "/api/progress",
        map[string]string{
            "Authorization": "Bearer " + token,
            "Idempotency-Key": idempotencyKey,
        },
        map[string]interface{}{
            "lesson_id": 1,
            "completed": true,
        },
    )
    assert.Equal(t, 200, resp2.StatusCode)
    assert.Equal(t, "true", resp2.Header.Get("X-Idempotent-Replay"))

    // Responses should be identical
    assert.JSONEq(t, resp1.Body, resp2.Body)
}
```

## Best Practices

### For Frontend Developers

1. **Always Generate Keys**: Include idempotency keys on all mutating requests
2. **Use UUIDs**: Generate cryptographically random UUIDs (v4)
3. **Reuse on Retry**: Use the same key when retrying a failed request
4. **Don't Reuse**: Never reuse keys across different operations
5. **Check Replay Header**: Log when responses are replayed for debugging

### For Backend Developers

1. **Early Return**: Check idempotency before expensive operations
2. **Complete Response**: Cache the entire response body and status code
3. **Atomic Storage**: Use database transactions when storing responses
4. **Error Responses**: Also cache error responses (4xx, 5xx) to prevent retry storms
5. **Monitoring**: Track cache hit rates and storage growth

## Implementation Details

### Database Schema

```sql
CREATE TABLE idempotency_keys (
    key VARCHAR(100) PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    request_path VARCHAR(200) NOT NULL,
    request_method VARCHAR(10) NOT NULL,
    response_status INTEGER,
    response_body TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_idempotency_keys_user ON idempotency_keys(user_id);
CREATE INDEX idx_idempotency_keys_expires ON idempotency_keys(expires_at);
```

### Middleware Flow

```
Request with Idempotency-Key
    |
    v
Is method POST/PUT/PATCH?
    |
    v
Extract user_id from context
    |
    v
Check cache: GetIdempotentResponse(key, user_id)
    |
    +---> Cache Hit: Return cached response (with X-Idempotent-Replay header)
    |
    +---> Cache Miss: Continue to handler
              |
              v
          Process request
              |
              v
          Capture response (status + body)
              |
              v
          Store in cache: StoreIdempotentResponse(key, user_id, response, expires_at)
              |
              v
          Return response to client
```

### Cleanup Job

A daily background job removes expired entries:

```go
go func() {
    ticker := time.NewTicker(24 * time.Hour)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            deleted, err := s.db.DeleteExpiredIdempotencyKeys(context.Background())
            if err != nil {
                logging.Error("failed to clean expired idempotency keys", err)
            } else {
                logging.Info("cleaned expired idempotency keys", "count", deleted)
            }
        case <-s.collectorCtx.Done():
            return
        }
    }
}()
```

## Security Considerations

1. **User Isolation**: Keys are scoped per user to prevent cross-user exploits
2. **Key Length**: Maximum 100 characters to prevent DoS via large keys
3. **Expiration**: 24-hour TTL prevents indefinite storage growth
4. **Auth Required**: Idempotency only works on authenticated routes
5. **No Leakage**: Cached responses never leak to different users

## Limitations

1. **GET Requests**: Not applicable to GET (idempotent by nature)
2. **Public Routes**: Not enabled on public routes (no user context)
3. **Large Responses**: Very large responses consume more storage
4. **Cross-Instance**: Keys are not shared across multiple backend instances (consider Redis for distributed deployments)

## Monitoring

### Metrics to Track

- Cache hit rate: `idempotency_cache_hits / idempotency_total_requests`
- Storage size: `SELECT COUNT(*) FROM idempotency_keys`
- Oldest entry: `SELECT MIN(created_at) FROM idempotency_keys`
- Cleanup effectiveness: Track deleted rows per cleanup run

### Logs

The system logs:
- Daily cleanup: `"Cleaned N expired idempotency keys"`
- Errors during cleanup: `"Failed to clean expired idempotency keys"`

## Troubleshooting

### Issue: Duplicate operations still occurring

**Cause**: Client not sending or reusing idempotency keys

**Solution**: Verify the `Idempotency-Key` header is present in requests

### Issue: "Invalid token: missing user_id"

**Cause**: Idempotency middleware running before auth middleware

**Solution**: Ensure middleware order in routes.go is correct:
```go
r.Use(authMiddleware.VerifyToken)
r.Use(middleware.Idempotency(s.db))
```

### Issue: Clients always receiving stale data

**Cause**: Reusing the same key for different operations

**Solution**: Generate a fresh UUID for each distinct operation

### Issue: Storage growing too large

**Cause**: High request volume or cleanup job not running

**Solution**:
1. Verify cleanup job is running: Check logs for cleanup messages
2. Reduce TTL: Change from 24h to shorter duration if appropriate
3. Add monitoring: Track table size over time

## Migration Guide

To add idempotency to an existing Shadow Nova deployment:

1. Run migration: `psql -f backend/internal/database/migrations/004_add_idempotency.sql`
2. Deploy backend with idempotency middleware
3. Update frontend to include idempotency keys
4. Monitor cache hit rates and adjust TTL if needed

## Future Enhancements

Potential improvements for future versions:

1. **Distributed Cache**: Use Redis for multi-instance deployments
2. **Configurable TTL**: Allow per-endpoint TTL configuration
3. **Compression**: Compress large response bodies before storage
4. **Metrics Dashboard**: Real-time visualization of idempotency stats
5. **Smart Cleanup**: Clean expired keys more frequently during high load
