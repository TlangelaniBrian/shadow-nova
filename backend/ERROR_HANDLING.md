# Error Handling Guide

This document outlines error handling patterns and best practices for the Shadow Nova backend.

## Core Principles

1. **Errors are Values**: All functions that can fail should return errors, never panic
2. **Error Wrapping**: Always wrap errors with context using `fmt.Errorf` with `%w` verb
3. **Fail Fast, Fail Clearly**: Return errors immediately with descriptive messages
4. **Graceful Degradation**: Non-critical failures should log warnings and continue

## Error Handling Patterns

### 1. Error Wrapping

Always wrap errors to preserve the error chain:

```go
config, err := pgxpool.ParseConfig(databaseUrl)
if err != nil {
    return nil, fmt.Errorf("unable to parse database URL: %w", err)
}
```

### 2. Error Inspection

Use `errors.Is` and `errors.As` to inspect wrapped errors:

```go
import (
    "errors"
    "github.com/jackc/pgx/v5"
)

// Check for specific error types
if errors.Is(err, pgx.ErrNoRows) {
    return nil, fmt.Errorf("record not found: %w", err)
}

// Extract specific error values
var pgErr *pgconn.PgError
if errors.As(err, &pgErr) {
    // Handle PostgreSQL-specific error
    if pgErr.Code == "23505" { // Unique violation
        return fmt.Errorf("duplicate entry: %w", err)
    }
}
```

### 3. Sentinel Errors

Define common errors as package-level variables for consistent checking:

```go
package database

import "errors"

var (
    ErrNotFound      = errors.New("record not found")
    ErrDuplicate     = errors.New("duplicate record")
    ErrUnauthorized  = errors.New("unauthorized access")
    ErrInvalidInput  = errors.New("invalid input")
)

// Usage
func GetUser(ctx context.Context, id int) (*User, error) {
    var user User
    err := db.QueryRow(ctx, query, id).Scan(&user)
    if errors.Is(err, pgx.ErrNoRows) {
        return nil, ErrNotFound
    }
    if err != nil {
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    return &user, nil
}
```

### 4. Error Context

Provide meaningful context in error messages:

```go
// Bad - no context
return nil, err

// Good - with context
return nil, fmt.Errorf("failed to create user with email %s: %w", email, err)

// Better - with structured context (for logging)
log.Printf("Failed to create user: email=%s, error=%v", email, err)
return nil, fmt.Errorf("failed to create user: %w", err)
```

## HTTP Status Code Mapping

Map application errors to appropriate HTTP status codes:

| Error Type | HTTP Status | Use Case |
|------------|-------------|----------|
| `ErrNotFound` | 404 | Resource doesn't exist |
| `ErrUnauthorized` | 401 | Invalid or missing authentication |
| `ErrForbidden` | 403 | Authenticated but lacks permission |
| `ErrInvalidInput` | 400 | Invalid request data |
| `ErrDuplicate` | 409 | Resource conflict (e.g., duplicate email) |
| `ErrInternal` | 500 | Unexpected server errors |
| `ErrServiceUnavailable` | 503 | Database or external service down |

### Example Handler Pattern

```go
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")

    user, err := h.db.GetUser(r.Context(), id)
    if err != nil {
        if errors.Is(err, database.ErrNotFound) {
            httputil.WriteError(w, http.StatusNotFound, "User not found")
            return
        }
        log.Printf("Failed to get user %s: %v", id, err)
        httputil.WriteError(w, http.StatusInternalServerError, "Internal server error")
        return
    }

    httputil.WriteJSON(w, http.StatusOK, user)
}
```

## Logging vs Returning Errors

### Return Errors When:
- The caller needs to handle the error
- The error should propagate up the stack
- The error represents a recoverable condition
- The error affects business logic

### Log Errors When:
- The error is handled at the current level
- Background operations fail non-critically
- You need debugging information but execution continues
- The error is unexpected but recoverable

### Example:

```go
// Return error - caller needs to handle
func (s *service) CreateUser(ctx context.Context, user *User) error {
    if err := s.db.Insert(ctx, user); err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }
    return nil
}

// Log warning - non-critical failure
func (s *service) NotifyUser(ctx context.Context, userID int, message string) {
    if err := s.notifier.Send(ctx, userID, message); err != nil {
        log.Printf("Warning: Failed to send notification to user %d: %v", userID, err)
        // Continue execution - notification failure is not critical
    }
}
```

## Error Response Format

Consistent JSON error responses:

```go
type ErrorResponse struct {
    Error   string `json:"error"`
    Status  int    `json:"status"`
    Details string `json:"details,omitempty"`
}

// Usage
httputil.WriteJSON(w, http.StatusBadRequest, ErrorResponse{
    Error:   "Invalid input",
    Status:  400,
    Details: "Email address is required",
})
```

## Initialization Errors

Critical initialization errors should prevent startup:

```go
func main() {
    // Critical - must succeed
    db, err := database.New()
    if err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    defer db.Close()

    // Non-critical - continue with degraded functionality
    flagsService, err := flags.New()
    if err != nil {
        log.Printf("Warning: Failed to initialize feature flags: %v", err)
        flagsService = flags.NewMockService()
    }

    // Critical - must succeed
    httpServer, appServer, err := server.NewServer(db, flagsService)
    if err != nil {
        log.Fatalf("Failed to create server: %v", err)
    }
}
```

## Database Error Handling

### Query Errors

```go
func (s *service) GetUser(ctx context.Context, id int) (*User, error) {
    var user User
    err := s.db.QueryRow(ctx, query, id).Scan(&user)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, ErrNotFound
        }
        return nil, fmt.Errorf("failed to query user: %w", err)
    }
    return &user, nil
}
```

### Transaction Errors

```go
func (s *service) CreateUserWithProfile(ctx context.Context, user *User, profile *Profile) error {
    tx, err := s.db.Begin(ctx)
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback(ctx) // Safe to call even after commit

    if err := createUser(ctx, tx, user); err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }

    if err := createProfile(ctx, tx, profile); err != nil {
        return fmt.Errorf("failed to create profile: %w", err)
    }

    if err := tx.Commit(ctx); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }

    return nil
}
```

## Validation Errors

Validate input before processing:

```go
func (s *service) CreateUser(ctx context.Context, user *User) error {
    // Validate input
    if user.Email == "" {
        return fmt.Errorf("email is required: %w", ErrInvalidInput)
    }
    if !isValidEmail(user.Email) {
        return fmt.Errorf("invalid email format: %w", ErrInvalidInput)
    }

    // Process
    if err := s.db.Insert(ctx, user); err != nil {
        return fmt.Errorf("failed to insert user: %w", err)
    }

    return nil
}
```

## External Service Errors

Handle external service failures gracefully:

```go
func (s *service) FetchExternalData(ctx context.Context) (*Data, error) {
    // Set timeout for external call
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    data, err := s.externalAPI.Fetch(ctx)
    if err != nil {
        // Log but don't expose internal details to client
        log.Printf("External API error: %v", err)
        return nil, fmt.Errorf("failed to fetch external data: %w", ErrServiceUnavailable)
    }

    return data, nil
}
```

## Testing Error Paths

Always test error handling:

```go
func TestCreateUser_DuplicateEmail(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()

    user := &User{Email: "test@example.com"}

    // First insert should succeed
    err := db.CreateUser(context.Background(), user)
    require.NoError(t, err)

    // Second insert should fail with duplicate error
    err = db.CreateUser(context.Background(), user)
    require.Error(t, err)
    assert.True(t, errors.Is(err, database.ErrDuplicate))
}
```

## Common Mistakes to Avoid

1. **Never ignore errors**: Always handle or explicitly document why it's safe to ignore
2. **Don't panic**: Use error returns instead of `panic()` for recoverable errors
3. **Don't lose context**: Always wrap errors with `%w` to preserve the error chain
4. **Don't expose internal details**: Return generic errors to clients, log specifics
5. **Don't return `nil, nil`**: If operation succeeds, return value; if fails, return error

## Migration Checklist

When replacing panics with error handling:

- [ ] Update function signature to return error
- [ ] Update all call sites to handle the new error
- [ ] Add error wrapping with context
- [ ] Update tests to verify error handling
- [ ] Document error conditions in function comments
- [ ] Map errors to appropriate HTTP status codes
- [ ] Add logging for unexpected errors
- [ ] Verify graceful degradation for non-critical failures
