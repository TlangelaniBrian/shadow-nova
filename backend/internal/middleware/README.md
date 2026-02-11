# Middleware Package

This package contains HTTP middleware for the Shadow Nova application.

## Available Middleware

### Authentication & Authorization

#### `auth.go` - JWT Authentication
- `VerifyToken` - Validates JWT tokens and sets user context
- `GetUserID(r *http.Request)` - Extract user ID from request context
- `GetUserRole(r *http.Request)` - Extract user role from request context

#### `admin.go` - Admin Role Check
- `AdminOnly` - Ensures the user has admin role
- Must be used after `VerifyToken` middleware

### Security

#### `security.go` - Security Headers
- `SecurityHeaders` - Adds security headers (CSP, HSTS, etc.)
- `CORSMiddleware` - Handles CORS for cross-origin requests

#### `ownership.go` - IDOR Prevention
- `ValidatePathOwnership` - Validates user access to learning paths
- `ValidateSubmissionOwnership` - Validates ownership of project submissions
- `ValidateProgressOwnership` - Validates ownership of progress records

#### `csrf.go` - CSRF Protection
- `CSRFMiddleware` - Validates CSRF tokens for state-changing operations

#### `rate_limiter.go` - Rate Limiting
- `NewRateLimiter(limit int)` - Creates a rate limiter
- `Limit` - Limits requests per IP address

### Observability

#### `prometheus.go` - Metrics
- `PrometheusMiddleware` - Collects HTTP metrics for monitoring

## Usage Examples

### Basic Authentication
```go
r.Group(func(r chi.Router) {
    r.Use(authMiddleware.VerifyToken)

    // Protected routes here
    r.Get("/profile", profileHandler)
})
```

### Admin-Only Routes
```go
r.Group(func(r chi.Router) {
    r.Use(authMiddleware.VerifyToken)
    r.Use(middleware.AdminOnly)

    // Admin routes here
    r.Post("/admin/settings", adminHandler)
})
```

### Resource Ownership Validation
```go
// Single resource protection
r.With(middleware.ValidateSubmissionOwnership(db)).Get("/submissions/{id}", handler)

// Multiple routes with same protection
r.Group(func(r chi.Router) {
    r.Use(authMiddleware.VerifyToken)
    r.Use(middleware.ValidateSubmissionOwnership(db))

    r.Get("/submissions/{id}", getHandler)
    r.Patch("/submissions/{id}", updateHandler)
    r.Delete("/submissions/{id}", deleteHandler)
})
```

### Combined Middleware
```go
// Order matters: Auth → Ownership → Handler
r.Route("/submissions/{id}", func(r chi.Router) {
    r.Use(authMiddleware.VerifyToken)  // First: Authenticate
    r.Use(middleware.ValidateSubmissionOwnership(db))  // Second: Authorize

    r.Get("/", handler.Get)
    r.Patch("/", handler.Update)
})
```

## Middleware Order

The order of middleware execution is important:

1. **Security Headers** - Applied globally
2. **CORS** - Applied globally
3. **Rate Limiting** - Applied globally
4. **Logging** - Applied globally
5. **Metrics** - Applied globally
6. **Authentication** (`VerifyToken`) - Applied to protected routes
7. **CSRF** - Applied to state-changing endpoints
8. **Authorization** (`AdminOnly`, Ownership checks) - Applied to specific routes
9. **Handler** - Your actual route handler

Example:
```go
r := chi.NewRouter()

// Global middleware (order matters)
r.Use(middleware.SecurityHeaders)
r.Use(middleware.CORSMiddleware())
r.Use(rateLimiter.Limit)
r.Use(chimiddleware.Logger)
r.Use(middleware.PrometheusMiddleware)

r.Route("/api", func(r chi.Router) {
    // Protected routes
    r.Group(func(r chi.Router) {
        r.Use(authMiddleware.VerifyToken)  // Auth required
        r.Use(middleware.CSRFMiddleware)   // CSRF protection

        // Resource with ownership check
        r.With(middleware.ValidateSubmissionOwnership(db)).Get("/submissions/{id}", handler)
    })
})
```

## Creating New Ownership Middleware

If you need to protect a new resource type:

1. **Add database method** to `database.Service` interface:
```go
UserOwnsResource(ctx context.Context, userID int, resourceID int) (bool, error)
```

2. **Implement the method** in appropriate database file:
```go
func (s *service) UserOwnsResource(ctx context.Context, userID int, resourceID int) (bool, error) {
    query := `SELECT user_id FROM resources WHERE id = $1`
    var ownerID int
    err := s.db.QueryRow(ctx, query, resourceID).Scan(&ownerID)
    if err != nil {
        return false, fmt.Errorf("failed to get resource owner: %w", err)
    }
    return ownerID == userID, nil
}
```

3. **Create middleware function** in `ownership.go`:
```go
func ValidateResourceOwnership(db database.Service) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            userID, ok := GetUserID(r)
            if !ok {
                httputil.WriteError(w, http.StatusUnauthorized, "User not authenticated")
                return
            }

            resourceIDStr := chi.URLParam(r, "id")
            if resourceIDStr == "" {
                httputil.WriteError(w, http.StatusBadRequest, "Resource ID is required")
                return
            }

            resourceID, err := strconv.Atoi(resourceIDStr)
            if err != nil {
                httputil.WriteError(w, http.StatusBadRequest, "Invalid resource ID")
                return
            }

            ownsResource, err := db.UserOwnsResource(r.Context(), userID, resourceID)
            if err != nil {
                httputil.WriteError(w, http.StatusInternalServerError, "Failed to verify resource ownership")
                return
            }

            if !ownsResource {
                httputil.WriteError(w, http.StatusForbidden, "You do not have access to this resource")
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}
```

4. **Add tests** in `ownership_test.go`

5. **Apply to routes** in `routes.go`

## Testing

Run all middleware tests:
```bash
go test ./internal/middleware/...
```

Run specific middleware tests:
```bash
go test ./internal/middleware/ownership_test.go
```

Run with coverage:
```bash
go test -cover ./internal/middleware/...
```

## Error Responses

All ownership middleware returns standardized error responses:

- **401 Unauthorized** - User not authenticated (missing/invalid token)
- **400 Bad Request** - Invalid or missing resource ID
- **403 Forbidden** - User doesn't own the resource
- **500 Internal Server Error** - Database or internal error

## Best Practices

1. **Always authenticate before authorizing**
   ```go
   r.Use(authMiddleware.VerifyToken)  // First
   r.Use(ownershipMiddleware)          // Second
   ```

2. **Use `With()` for single route protection**
   ```go
   r.With(middleware.ValidateSubmissionOwnership(db)).Get("/submissions/{id}", handler)
   ```

3. **Use `Use()` for protecting multiple routes**
   ```go
   r.Group(func(r chi.Router) {
       r.Use(middleware.ValidateSubmissionOwnership(db))
       r.Get("/submissions/{id}", getHandler)
       r.Patch("/submissions/{id}", updateHandler)
   })
   ```

4. **Extract IDs from URL, not body**
   - IDs in URLs are validated by middleware
   - IDs in body can be manipulated

5. **Return 403, not 404**
   - Don't reveal whether resource exists
   - Always return "access denied" message

6. **Test ownership checks thoroughly**
   - Owner can access
   - Non-owner gets 403
   - Invalid IDs get 400
   - Missing auth gets 401

## Related Documentation

- `OWNERSHIP_TESTING.md` - Testing guide for ownership middleware
- `../../../IDOR_PREVENTION_IMPLEMENTATION.md` - Full implementation details

## Security Notes

- Middleware order matters for security
- Always validate resource IDs from URL parameters
- Never trust user input, always validate
- Log failed authorization attempts (future enhancement)
- Consider rate limiting for failed access attempts
