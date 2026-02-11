# Error Handling Guide

This guide documents the error handling philosophy and best practices for Shadow Nova backend.

## Philosophy

Shadow Nova follows these error handling principles:

1. **Errors are values**: Use Go's native error handling, not exceptions
2. **Wrap errors with context**: Use `fmt.Errorf` with `%w` to maintain error chains
3. **Sentinel errors for common cases**: Use predefined errors for expected conditions
4. **Type-safe error checking**: Use `errors.Is()` and `errors.As()` for error inspection
5. **Centralized HTTP error mapping**: Single source of truth for error-to-HTTP-status mapping
6. **Structured logging**: Use structured logs for errors that need investigation
7. **No panics in library code**: Only fatal initialization errors in main.go should panic

## Error Types

### Sentinel Errors

Located in `internal/errors/errors.go`:

```go
var (
    ErrNotFound       = errors.New("resource not found")
    ErrUnauthorized   = errors.New("unauthorized access")
    ErrForbidden      = errors.New("forbidden")
    ErrInvalidInput   = errors.New("invalid input")
    ErrDuplicateEntry = errors.New("duplicate entry")
    ErrDatabaseError  = errors.New("database error")
)
```

### AppError Type

Custom error type with HTTP context:

```go
type AppError struct {
    Err      error   // Underlying error
    Message  string  // User-friendly message
    Code     string  // Error code (e.g., "NOT_FOUND")
    HTTPCode int     // HTTP status code
}
```

## When to Use Each Error Type

### Database Layer

Use sentinel errors for expected conditions:

```go
// Not found - expected condition
if err == pgx.ErrNoRows {
    return nil, errors.NotFound("user with email example@test.com not found")
}

// Database error - unexpected condition
return nil, errors.DatabaseError(err, "failed to query users")
```

### Service Layer

Propagate database errors and add business logic errors:

```go
user, err := s.db.GetUserByEmail(ctx, email)
if err != nil {
    return err // Propagate as-is
}

if !user.IsActive {
    return errors.Forbidden("account is inactive")
}
```

### Handler Layer

Use centralized error handler:

```go
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
    data, err := h.service.GetData(r.Context(), id)
    if err != nil {
        httputil.HandleError(w, err) // Centralized mapping
        return
    }
    httputil.WriteSuccess(w, "Success", data)
}
```

## Error Constructors

Convenience constructors are provided in `internal/errors/errors.go`:

```go
// Returns 404
errors.NotFound("learning path not found")

// Returns 401
errors.Unauthorized("invalid credentials")

// Returns 403
errors.Forbidden("insufficient permissions")

// Returns 400
errors.InvalidInput("email format is invalid")

// Returns 409
errors.DuplicateEntry("email already exists")

// Returns 500
errors.DatabaseError(err, "failed to save user")
```

## Error Wrapping Best Practices

### Do Wrap Errors

```go
// Good - preserves error chain
if err != nil {
    return fmt.Errorf("failed to process payment: %w", err)
}

// Good - uses error constructor
if err == pgx.ErrNoRows {
    return errors.NotFound("order not found")
}
```

### Don't Wrap Errors

```go
// Bad - loses error chain
if err != nil {
    return fmt.Errorf("failed to process payment: %v", err)
}

// Bad - loses type information
if err != nil {
    return errors.New("something went wrong")
}
```

## HTTP Status Code Mapping

The centralized error handler maps errors to HTTP status codes:

| Error Type | HTTP Status | Use Case |
|------------|-------------|----------|
| `ErrNotFound` | 404 | Resource doesn't exist |
| `ErrUnauthorized` | 401 | Authentication required/failed |
| `ErrForbidden` | 403 | Authenticated but not authorized |
| `ErrInvalidInput` | 400 | Validation failed |
| `ErrDuplicateEntry` | 409 | Unique constraint violation |
| `ErrDatabaseError` | 500 | Database operation failed |
| Other | 500 | Unhandled error |

## Checking Error Types

Use `errors.Is()` for sentinel errors:

```go
if errors.IsNotFound(err) {
    // Handle not found
}
```

Use `errors.As()` for AppError:

```go
var appErr *errors.AppError
if errors.As(err, &appErr) {
    log.Printf("Error code: %s, HTTP: %d", appErr.Code, appErr.HTTPCode)
}
```

## Logging Strategies

### When to Log

- **Always log 5xx errors**: These indicate bugs or infrastructure issues
- **Don't log 4xx errors**: These are expected user errors
- **Log with context**: Include relevant IDs, user info, etc.

### How to Log

Use structured logging from `internal/logging/logger.go`:

```go
// Error with context
logging.Error("Failed to process payment", err,
    "user_id", userID,
    "order_id", orderID,
)

// Info
logging.Info("User registered",
    "user_id", user.ID,
    "email", user.Email,
)
```

### What Not to Log

- User passwords (even hashed)
- API keys or tokens
- Credit card numbers
- PII (unless necessary and approved)

## Security Considerations

### Don't Leak Information

```go
// Bad - reveals user existence
if err != nil {
    return errors.NotFound("user with email user@example.com not found")
}

// Good - generic message for auth failures
if err != nil {
    return errors.Unauthorized("invalid credentials")
}
```

### Consistent Error Messages

Authentication failures should always return the same message:

```go
// User not found
return errors.Unauthorized("invalid credentials")

// Wrong password
return errors.Unauthorized("invalid credentials")
```

## Testing Error Handling

### Test Error Scenarios

```go
func TestGetUser_NotFound(t *testing.T) {
    err := service.GetUser(ctx, "nonexistent-id")

    // Check error type
    if !errors.IsNotFound(err) {
        t.Errorf("expected NotFound error, got %v", err)
    }
}
```

### Test Error Wrapping

```go
func TestErrorChain(t *testing.T) {
    err := service.DoSomething(ctx)

    // Check wrapped error
    if !errors.Is(err, ErrDatabaseError) {
        t.Error("expected database error in chain")
    }
}
```

## Migration Checklist

When updating existing code:

- [ ] Replace `panic()` calls with `return error`
- [ ] Use sentinel errors for expected conditions
- [ ] Wrap errors with context using `%w`
- [ ] Update handlers to use `httputil.HandleError()`
- [ ] Add structured logging for 5xx errors
- [ ] Write tests for error scenarios
- [ ] Document any new error types

## Examples

### Complete Handler Example

```go
func (h *PathsHandler) Get(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")

    // Validation
    if id == "" {
        httputil.HandleError(w, errors.InvalidInput("path ID is required"))
        return
    }

    // Database call (returns typed errors)
    path, err := h.db.GetLearningPath(r.Context(), id)
    if err != nil {
        httputil.HandleError(w, err) // Centralized handling
        return
    }

    httputil.WriteSuccess(w, "Learning path retrieved successfully", path)
}
```

### Complete Database Method Example

```go
func (s *service) GetLearningPath(ctx context.Context, id string) (*models.LearningPath, error) {
    query := `SELECT id, title, description FROM learning_paths WHERE id = $1`

    var path models.LearningPath
    err := s.db.QueryRow(ctx, query, id).Scan(&path.ID, &path.Title, &path.Description)

    if err != nil {
        // Expected condition
        if err == pgx.ErrNoRows {
            return nil, errors.NotFound(fmt.Sprintf("learning path %s not found", id))
        }
        // Unexpected condition
        return nil, errors.DatabaseError(err, "failed to get learning path")
    }

    return &path, nil
}
```

## Common Patterns

### Resource Ownership Validation

```go
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
    userID := middleware.GetUserID(r)
    resourceID := chi.URLParam(r, "id")

    // Check ownership
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
        if err := s.doStepOne(ctx, tx); err != nil {
            return fmt.Errorf("step one failed: %w", err)
        }

        if err := s.doStepTwo(ctx, tx); err != nil {
            return fmt.Errorf("step two failed: %w", err)
        }

        return nil // Transaction commits automatically
    })
}
```

## References

- [Go Blog: Error Handling and Go](https://go.dev/blog/error-handling-and-go)
- [Go Blog: Working with Errors in Go 1.13](https://go.dev/blog/go1.13-errors)
- [Effective Go: Errors](https://go.dev/doc/effective_go#errors)
