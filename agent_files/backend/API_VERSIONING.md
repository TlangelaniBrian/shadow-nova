# Shadow Nova API Versioning Strategy

## Overview

Shadow Nova uses URI-based API versioning to ensure backward compatibility and smooth transitions between API versions. All API endpoints are prefixed with `/api/v{version}` (e.g., `/api/v1`).

## Current Version

**Current API Version**: `v1`
**Base URL**: `http://localhost:8080/api/v1` (development)
**Production Base URL**: `https://your-domain.com/api/v1`

## Versioning Philosophy

### Why Version APIs?

1. **Backward Compatibility**: Existing clients continue to work when you introduce breaking changes
2. **Gradual Migration**: Users can migrate to new versions at their own pace
3. **Clear Communication**: Version numbers signal breaking changes vs. non-breaking changes
4. **Deprecation Strategy**: Old versions can be deprecated and eventually removed

### URI Versioning Benefits

- **Simple and Explicit**: Version is clear in the URL
- **Easy to Route**: Different versions can use different handlers/implementations
- **Client-Friendly**: No custom headers required
- **Cache-Friendly**: CDNs and proxies handle versioned URLs naturally

## Version Identification

### Version Info Endpoint

The `/version` endpoint (unversioned) provides discoverability:

```bash
curl http://localhost:8080/version
```

**Response**:
```json
{
  "version": "1.0.0",
  "api_version": "v1"
}
```

This endpoint:
- Is **not versioned** (no `/api/v1` prefix) for easy discovery
- Returns application version (semantic versioning)
- Returns current API version
- Never changes URL (permanent endpoint)

## Semantic Versioning

Shadow Nova follows semantic versioning for the application:

**Format**: `MAJOR.MINOR.PATCH` (e.g., `1.0.0`)

- **MAJOR**: Breaking changes (e.g., `1.0.0` → `2.0.0`)
- **MINOR**: New features, backward compatible (e.g., `1.0.0` → `1.1.0`)
- **PATCH**: Bug fixes, backward compatible (e.g., `1.0.0` → `1.0.1`)

API versions (v1, v2, etc.) increment only on **breaking changes**.

## What Constitutes a Breaking Change?

### Breaking Changes (Require New API Version)

These require incrementing the API version (v1 → v2):

1. **Removing endpoints**
   ```
   DELETE /api/v1/old-endpoint  # Removed in v2
   ```

2. **Removing request/response fields**
   ```json
   // v1 response
   {"id": 1, "name": "Test", "deprecated_field": "value"}

   // v2 response (breaking)
   {"id": 1, "name": "Test"}  // deprecated_field removed
   ```

3. **Changing field types**
   ```json
   // v1
   {"id": "123"}  // string

   // v2 (breaking)
   {"id": 123}  // number
   ```

4. **Changing authentication schemes**
   ```
   v1: Bearer token in Authorization header
   v2: API key in X-API-Key header (breaking)
   ```

5. **Changing URL structure**
   ```
   v1: GET /api/v1/users/123/posts
   v2: GET /api/v2/posts?user_id=123  (breaking)
   ```

6. **Changing error response format**
   ```json
   // v1
   {"error": "Not found"}

   // v2 (breaking)
   {"status": 404, "message": "Not found", "code": "NOT_FOUND"}
   ```

### Non-Breaking Changes (Same API Version)

These can be added to the current version:

1. **Adding new endpoints**
   ```
   POST /api/v1/new-feature  # New endpoint in v1
   ```

2. **Adding optional fields to requests**
   ```json
   // Existing: {"name": "Test"}
   // Enhanced: {"name": "Test", "description": "Optional field"}
   ```

3. **Adding fields to responses**
   ```json
   // Before
   {"id": 1, "name": "Test"}

   // After (non-breaking)
   {"id": 1, "name": "Test", "created_at": "2024-01-01"}
   ```

4. **Making required fields optional**
   ```json
   // Before: "email" required
   // After: "email" optional (non-breaking)
   ```

5. **Bug fixes and performance improvements**

6. **Adding new query parameters (optional)**
   ```
   GET /api/v1/paths?sort=asc  # New optional parameter
   ```

## Version Lifecycle

### Stage 1: Active Development
- **Status**: Current version (v1)
- **Support**: Full feature development and bug fixes
- **Documentation**: Complete and up-to-date
- **Deprecation**: None

### Stage 2: Maintenance
- **Status**: Previous version (when v2 is released)
- **Support**: Critical bug fixes and security patches only
- **Documentation**: Maintained but marked as "legacy"
- **Deprecation**: Announced with timeline (e.g., 6 months)

### Stage 3: Deprecated
- **Status**: Old version (v1 when v3 is released)
- **Support**: None
- **Documentation**: Archived
- **Deprecation**: Active warnings in responses
- **Timeline**: 3 months until removal

### Stage 4: Removed
- **Status**: Completely removed
- **Support**: None
- **Documentation**: Removed (only in git history)

## Deprecation Strategy

### Announcing Deprecation

When deprecating v1 after releasing v2:

1. **Update documentation**
   ```markdown
   ## API v1 (Deprecated)

   **⚠️ Deprecated**: v1 will be removed on 2024-12-31
   **Migration Guide**: See [v1-to-v2-migration.md](v1-to-v2-migration.md)
   **Replacement**: Use v2 (`/api/v2`)
   ```

2. **Add deprecation headers**
   ```go
   w.Header().Set("X-API-Deprecated", "true")
   w.Header().Set("X-API-Sunset", "2024-12-31")
   w.Header().Set("X-API-Migration", "/docs/v1-to-v2-migration")
   ```

3. **Log deprecation warnings**
   ```go
   logging.Warn(r.Context(), "API v1 deprecated, please migrate to v2",
       "endpoint", r.URL.Path,
       "user_agent", r.UserAgent(),
   )
   ```

4. **Email notifications** to registered API consumers

### Deprecation Timeline

**Recommended timeline**:
- **T+0**: Announce deprecation
- **T+1 month**: Add deprecation headers
- **T+3 months**: Start logging warnings
- **T+6 months**: Send email notifications
- **T+9 months**: Final warning (loud logging)
- **T+12 months**: Remove old version

## Implementing a New API Version

### When to Create v2

Create v2 when you need to make **breaking changes** that cannot be handled with:
- Feature flags
- Optional parameters
- Graceful degradation

### Implementation Steps

#### 1. Create v2 Route Group

```go
// backend/internal/server/routes.go

// Existing v1
r.Route("/api/v1", func(r chi.Router) {
    // ... v1 routes
})

// New v2 (parallel to v1)
r.Route("/api/v2", func(r chi.Router) {
    // Apply same middleware
    r.Use(csrfMiddleware)

    // Initialize v2 handlers (can reuse or create new)
    pathsHandlerV2 := handlers.NewPathsHandlerV2(s.db)

    // Define v2 routes (can be different from v1)
    r.Get("/learning-paths", pathsHandlerV2.List)  // Note: different URL structure
    r.Get("/learning-paths/{id}", pathsHandlerV2.Get)
})
```

#### 2. Create v2 Handlers (if needed)

**Option A**: Reuse v1 handlers if logic is the same
```go
// No new code needed, just route to existing handlers
r.Get("/paths", pathsHandler.List)  // Same handler as v1
```

**Option B**: Create new handlers for breaking changes
```go
// backend/internal/handlers/paths_v2.go

type PathsHandlerV2 struct {
    db database.Service
}

func NewPathsHandlerV2(db database.Service) *PathsHandlerV2 {
    return &PathsHandlerV2{db: db}
}

func (h *PathsHandlerV2) List(w http.ResponseWriter, r *http.Request) {
    // v2 implementation with breaking changes
    // e.g., different response format
}
```

#### 3. Update Version Endpoint

```go
r.Get("/version", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{
        "version":"2.0.0",
        "api_version":"v2",
        "deprecated_versions":["v1"],
        "sunset_dates":{"v1":"2024-12-31"}
    }`))
})
```

#### 4. Add v1 Deprecation Headers

```go
r.Route("/api/v1", func(r chi.Router) {
    // Add deprecation middleware
    r.Use(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("X-API-Deprecated", "true")
            w.Header().Set("X-API-Sunset", "2024-12-31")
            w.Header().Set("X-API-Migration", "/docs/v1-to-v2")
            next.ServeHTTP(w, r)
        })
    })

    // ... v1 routes
})
```

#### 5. Update Frontend

```typescript
// frontend/src/api/client.ts

// Option 1: Gradual migration (feature flag)
const apiVersion = import.meta.env.VITE_API_VERSION || 'v1'
const client = axios.create({
    baseURL: `${import.meta.env.VITE_API_URL}/api/${apiVersion}`,
})

// Option 2: Direct cutover
const client = axios.create({
    baseURL: `${import.meta.env.VITE_API_URL}/api/v2`,
})
```

#### 6. Update Documentation

- Create `v2-migration-guide.md`
- Update `README.md` with v2 examples
- Mark v1 endpoints as deprecated in API docs
- Add breaking changes section

#### 7. Create Tests

```go
// backend/internal/handlers/paths_v2_test.go

func TestPathsHandlerV2_List(t *testing.T) {
    // Test v2 behavior
}

func TestV1ToV2Compatibility(t *testing.T) {
    // Ensure v1 still works during deprecation
}
```

## Testing API Versions

### Test Both Versions in Parallel

```bash
# Test v1 (should still work)
curl http://localhost:8080/api/v1/paths

# Test v2 (new version)
curl http://localhost:8080/api/v2/learning-paths

# Check deprecation headers
curl -I http://localhost:8080/api/v1/paths
# Should see: X-API-Deprecated: true
```

### Integration Tests

```go
// Test both versions
func TestAPIVersions(t *testing.T) {
    tests := []struct {
        version  string
        endpoint string
        expected int
    }{
        {"v1", "/api/v1/paths", 200},
        {"v2", "/api/v2/learning-paths", 200},
    }

    for _, tt := range tests {
        t.Run(tt.version, func(t *testing.T) {
            resp := makeRequest(tt.endpoint)
            assert.Equal(t, tt.expected, resp.StatusCode)
        })
    }
}
```

## Best Practices

### 1. Document Breaking Changes Clearly

Create a `CHANGELOG.md`:
```markdown
## v2.0.0 (2024-06-01) - BREAKING CHANGES

### Breaking Changes
- Changed `/api/v1/paths` to `/api/v2/learning-paths`
- Response format changed: `data` wrapper removed
- `user_id` now returns integer instead of string

### Migration Guide
See [v1-to-v2-migration.md](docs/v1-to-v2-migration.md)
```

### 2. Use Feature Flags for Gradual Rollout

```go
if flags.IsEnabled("api_v2") {
    // Route to v2 handler
} else {
    // Route to v1 handler (fallback)
}
```

### 3. Monitor API Version Usage

```go
// Metrics
prometheus.Counter("api_requests_total").
    With("version", "v1").
    Inc()
```

Track:
- Request counts per version
- Error rates per version
- Response times per version
- Unique users per version

### 4. Provide Migration Tools

Create CLI tools or scripts:
```bash
# Example: Validate v2 compatibility
./scripts/validate-v2-migration.sh
```

### 5. Avoid Over-Versioning

Don't create v2 for minor changes. Use versioning sparingly for **true breaking changes** only.

## Common Pitfalls

### ❌ Don't: Version Every Endpoint Independently

```
/api/users/v1
/api/projects/v2  # Confusing!
```

### ✅ Do: Version the Entire API

```
/api/v1/users
/api/v1/projects
/api/v2/users
/api/v2/projects
```

### ❌ Don't: Make Breaking Changes in the Same Version

```go
// v1 (before)
type User struct { ID string }

// v1 (after) - BREAKING CHANGE IN SAME VERSION!
type User struct { ID int }
```

### ✅ Do: Increment Version for Breaking Changes

```go
// v1
type User struct { ID string }

// v2 (new version)
type UserV2 struct { ID int }
```

### ❌ Don't: Remove Old Versions Immediately

### ✅ Do: Follow Deprecation Timeline (6-12 months)

## Monitoring and Metrics

### Key Metrics to Track

1. **Version distribution**
   ```
   api_requests_total{version="v1"} 1000
   api_requests_total{version="v2"} 5000
   ```

2. **Deprecation warnings**
   ```
   api_deprecated_requests{version="v1"} 1000
   ```

3. **Error rates by version**
   ```
   api_errors_total{version="v1",status="500"} 10
   api_errors_total{version="v2",status="500"} 2
   ```

### Dashboard Queries

**Grafana Query** (version adoption):
```promql
sum(rate(api_requests_total[5m])) by (version)
```

## Resources

- [Semantic Versioning](https://semver.org/)
- [API Versioning Best Practices](https://restfulapi.net/versioning/)
- [HTTP Deprecation Headers (RFC 8594)](https://www.rfc-editor.org/rfc/rfc8594.html)

## Questions?

For questions about API versioning:
- Check the [CHANGELOG.md](../CHANGELOG.md)
- See [API documentation](../README.md#api-endpoints)
- Open an issue on GitHub
