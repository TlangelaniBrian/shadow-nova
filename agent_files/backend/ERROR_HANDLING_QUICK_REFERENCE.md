# Error Handling Quick Reference

## Quick Start

### Database Layer

```go
import (
    "shadow-nova/backend/internal/errors"
    "github.com/jackc/pgx/v5"
)

func (s *service) GetResource(ctx context.Context, id string) (*Resource, error) {
    var resource Resource
    err := s.db.QueryRow(ctx, query, id).Scan(&resource.ID, &resource.Name)

    if err != nil {
        // Expected: resource not found
        if err == pgx.ErrNoRows {
            return nil, errors.NotFound("resource not found")
        }
        // Unexpected: database error
        return nil, errors.DatabaseError(err, "failed to get resource")
    }

    return &resource, nil
}
```

### Handler Layer

```go
import "shadow-nova/backend/internal/httputil"

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
    data, err := h.service.GetData(r.Context(), id)
    if err != nil {
        httputil.HandleError(w, err) // Automatically maps error to HTTP status
        return
    }
    httputil.WriteSuccess(w, "Success", data)
}
```

## Error Constructors

| Constructor | HTTP Status | Use Case |
|-------------|-------------|----------|
| `errors.NotFound("msg")` | 404 | Resource doesn't exist |
| `errors.Unauthorized("msg")` | 401 | Auth required/failed |
| `errors.Forbidden("msg")` | 403 | Insufficient permissions |
| `errors.InvalidInput("msg")` | 400 | Validation failed |
| `errors.DuplicateEntry("msg")` | 409 | Unique constraint violation |
| `errors.DatabaseError(err, "msg")` | 500 | Database operation failed |

## Error Checking

```go
import "shadow-nova/backend/internal/errors"

// Check error type
if errors.IsNotFound(err) {
    // Handle not found
}

// Check AppError
var appErr *errors.AppError
if errors.As(err, &appErr) {
    log.Printf("Code: %s, HTTP: %d", appErr.Code, appErr.HTTPCode)
}
```

## Logging

```go
import "shadow-nova/backend/internal/logging"

// Error with context
logging.Error("Operation failed", err,
    "user_id", userID,
    "resource_id", resourceID,
)

// Info
logging.Info("User action", "action", "login", "user_id", userID)

// Warn
logging.Warn("Rate limit approaching", "user_id", userID, "count", count)
```

## Common Patterns

### Resource Ownership Check

```go
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
    userID := middleware.GetUserID(r)
    resourceID := chi.URLParam(r, "id")

    owns, err := h.db.UserOwnsResource(r.Context(), userID, resourceID)
    if err != nil {
        httputil.HandleError(w, err)
        return
    }

    if !owns {
        httputil.HandleError(w, errors.Forbidden("access denied"))
        return
    }

    // Proceed with update...
}
```

### Transaction Error Handling

```go
func (s *service) ComplexOperation(ctx context.Context) error {
    return s.WithTx(ctx, func(tx pgx.Tx) error {
        if err := s.stepOne(ctx, tx); err != nil {
            return fmt.Errorf("step one failed: %w", err)
        }

        if err := s.stepTwo(ctx, tx); err != nil {
            return fmt.Errorf("step two failed: %w", err)
        }

        return nil // Commits automatically
    })
}
```

### Security: Don't Leak Info

```go
// BAD - reveals user existence
user, err := h.db.GetUserByEmail(ctx, email)
if err != nil {
    httputil.HandleError(w, err) // Returns "user not found"
    return
}

// GOOD - generic message
user, err := h.db.GetUserByEmail(ctx, email)
if err != nil {
    httputil.WriteError(w, http.StatusUnauthorized, "invalid credentials")
    return
}
```

## Error Wrapping

```go
// Good - preserves error chain
if err != nil {
    return fmt.Errorf("failed to process: %w", err)
}

// Bad - loses error chain
if err != nil {
    return fmt.Errorf("failed to process: %v", err)
}
```

## HTTP Status Code Mapping

The centralized error handler automatically maps errors:

```go
httputil.HandleError(w, err)
```

Maps to:

- `ErrNotFound` → 404
- `ErrUnauthorized` → 401
- `ErrForbidden` → 403
- `ErrInvalidInput` → 400
- `ErrDuplicateEntry` → 409
- `ErrDatabaseError` → 500
- Other → 500

## Testing

```go
func TestGetResource_NotFound(t *testing.T) {
    err := service.GetResource(ctx, "nonexistent")

    // Check error type
    if !errors.IsNotFound(err) {
        t.Errorf("expected NotFound error, got %v", err)
    }
}

func TestErrorChain(t *testing.T) {
    err := service.DoSomething(ctx)

    // Check wrapped error
    if !errors.Is(err, errors.ErrDatabaseError) {
        t.Error("expected database error in chain")
    }
}
```

## Checklist

When adding new endpoints:

- [ ] Database methods return typed errors (NotFound, DatabaseError)
- [ ] Handlers use `httputil.HandleError(w, err)`
- [ ] Auth failures return generic "invalid credentials"
- [ ] 5xx errors are logged with context
- [ ] Error chains preserved with `%w`
- [ ] Tests cover error scenarios

## See Also

- Full documentation: `ERROR_HANDLING_GUIDE.md`
- Implementation summary: `ERROR_HANDLING_IMPLEMENTATION_SUMMARY.md`
