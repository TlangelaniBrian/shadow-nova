# Dependency Injection Architecture

## Overview

Shadow Nova has been refactored to use proper dependency injection instead of the singleton pattern. This document explains the benefits, architecture, and testing improvements.

## What Changed

### Before (Singleton Pattern)

```go
// database/database.go
var dbInstance *service

func New() Service {
    if dbInstance != nil {
        return dbInstance  // Return cached instance
    }
    // Create and cache instance
    dbInstance = &service{db: db}
    return dbInstance
}

// server/server.go
func NewServer() (*http.Server, *Server) {
    server := &Server{
        db: database.New(),  // Hidden dependency
    }
    return httpServer, server
}

// main.go
func main() {
    httpServer, appServer := server.NewServer()
    // Database created internally, no control
}
```

**Problems:**
- Hidden dependencies (not visible in function signatures)
- Difficult to test (can't inject mock database)
- Global state (singleton can cause race conditions)
- Hard to track lifetime (when is it created/closed?)
- Tight coupling (components directly depend on concrete implementation)

### After (Dependency Injection)

```go
// database/database.go
func New() (Service, error) {
    // Create new instance, no global state
    db, err := pgxpool.NewWithConfig(ctx, config)
    if err != nil {
        return nil, fmt.Errorf("unable to create connection pool: %w", err)
    }
    return &service{db: db}, nil
}

// server/server.go
func NewServer(db database.Service, flagsService flags.Service) (*http.Server, *Server, error) {
    server := &Server{
        db:    db,           // Injected dependency
        flags: flagsService, // Injected dependency
    }
    return httpServer, server, nil
}

// main.go
func main() {
    // Create dependencies
    db, err := database.New()
    if err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    defer db.Close()

    flagsService, err := flags.New()
    if err != nil {
        log.Printf("Warning: Failed to initialize feature flags: %v", err)
    }

    // Inject dependencies
    httpServer, appServer, err := server.NewServer(db, flagsService)
    if err != nil {
        log.Fatalf("Failed to create server: %v", err)
    }
}
```

**Benefits:**
- Explicit dependencies (visible in function signatures)
- Easy to test (can inject mocks)
- No global state (thread-safe by design)
- Clear lifetime management (created in main, passed down)
- Loose coupling (components depend on interfaces)

## Architecture

### Dependency Flow

```
main.go
  │
  ├─> database.New() ──> Database Service (interface)
  │                         │
  ├─> flags.New() ────────> Flags Service (interface)
  │                         │
  └─> server.NewServer(db, flags)
         │
         ├─> handlers.NewAuthHandler(googleAuth, db)
         ├─> handlers.NewGitHubHandler(githubAuth, db)
         ├─> handlers.NewPathsHandler(db)
         ├─> handlers.NewProgressHandler(db)
         ├─> handlers.NewProjectsHandler(db)
         └─> handlers.NewAdminHandler(db)
```

### Key Principles

1. **Single Responsibility**: Each component has one clear purpose
2. **Dependency Inversion**: Depend on interfaces, not concrete types
3. **Explicit Dependencies**: All dependencies visible in constructors
4. **Lifetime Management**: Dependencies created once, passed down
5. **Interface Segregation**: Small, focused interfaces

## Testing Improvements

### Before (Singleton)

```go
func TestHandler(t *testing.T) {
    // Can't control database, uses global singleton
    handler := NewAuthHandler()

    // Test runs against real database or shared test database
    // Race conditions possible if tests run in parallel
}
```

### After (Dependency Injection)

```go
// Create a mock database for testing
type MockDatabase struct {
    mock.Mock
}

func (m *MockDatabase) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
    args := m.Called(ctx, email)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*models.User), args.Error(1)
}

// Test with mock
func TestHandler(t *testing.T) {
    mockDB := new(MockDatabase)
    mockDB.On("GetUserByEmail", mock.Anything, "test@example.com").
        Return(&models.User{ID: 1, Email: "test@example.com"}, nil)

    handler := NewAuthHandler(nil, mockDB)

    // Test with controlled behavior
    // No database connection needed
    // Tests run fast and isolated
}
```

### Mock Example

```go
package handlers_test

import (
    "context"
    "testing"
    "shadow-nova/backend/internal/models"
    "github.com/stretchr/testify/mock"
)

// MockDatabaseService implements database.Service for testing
type MockDatabaseService struct {
    mock.Mock
}

func (m *MockDatabaseService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
    args := m.Called(ctx, email)
    if user := args.Get(0); user != nil {
        return user.(*models.User), args.Error(1)
    }
    return nil, args.Error(1)
}

func (m *MockDatabaseService) CreateUser(ctx context.Context, user *models.User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}

// ... implement other interface methods as needed

func TestAuthHandler_Login(t *testing.T) {
    // Arrange
    mockDB := new(MockDatabaseService)
    expectedUser := &models.User{
        ID:           1,
        Email:        "test@example.com",
        Username:     "testuser",
        PasswordHash: "$2a$10$...", // bcrypt hash
        Role:         "user",
    }

    mockDB.On("GetUserByEmail", mock.Anything, "test@example.com").
        Return(expectedUser, nil)

    handler := handlers.NewAuthHandler(nil, mockDB)

    // Act
    req := httptest.NewRequest("POST", "/api/login", nil)
    w := httptest.NewRecorder()
    handler.Login(w, req)

    // Assert
    assert.Equal(t, http.StatusOK, w.Code)
    mockDB.AssertExpectations(t)
}
```

## Connection Pool Configuration

The database now includes proper connection pool configuration:

```go
config.MaxConns = 25              // Maximum connections in pool
config.MinConns = 5               // Minimum idle connections
config.MaxConnLifetime = time.Hour        // Recycle connections after 1 hour
config.MaxConnIdleTime = 30 * time.Minute // Close idle connections after 30 mins
```

Benefits:
- **Performance**: Reuse connections instead of creating new ones
- **Resource Management**: Limit total connections to database
- **Connection Health**: Recycle stale connections automatically
- **Efficiency**: Keep minimum connections ready for quick requests

## Error Handling

All initialization now returns errors instead of panicking:

```go
// Before
func New() Service {
    db, err := pgxpool.NewWithConfig(ctx, config)
    if err != nil {
        panic(fmt.Sprintf("Unable to create pool: %v", err))  // Crashes app
    }
    return &service{db: db}
}

// After
func New() (Service, error) {
    db, err := pgxpool.NewWithConfig(ctx, config)
    if err != nil {
        return nil, fmt.Errorf("unable to create pool: %w", err)  // Returns error
    }
    return &service{db: db}, nil
}
```

Benefits:
- **Graceful Degradation**: App can handle initialization failures
- **Better Logging**: Errors propagate with context
- **Testability**: Can test error paths
- **Production Ready**: No unexpected crashes

## Migration Guide

### Adding a New Handler

```go
// 1. Define handler with dependencies
type MyHandler struct {
    db database.Service
}

// 2. Create constructor that accepts dependencies
func NewMyHandler(db database.Service) *MyHandler {
    return &MyHandler{db: db}
}

// 3. In routes.go, create handler with injected db
func (s *Server) RegisterRoutes() (http.Handler, error) {
    myHandler := handlers.NewMyHandler(s.db)

    r.Get("/my-endpoint", myHandler.MyMethod)

    return r, nil
}
```

### Adding a New Dependency

```go
// 1. Create the service in main.go
cache, err := cache.New()
if err != nil {
    log.Fatalf("Failed to initialize cache: %v", err)
}
defer cache.Close()

// 2. Add to Server struct
type Server struct {
    db    database.Service
    flags flags.Service
    cache cache.Service  // New dependency
}

// 3. Update NewServer to accept it
func NewServer(db database.Service, flags flags.Service, cache cache.Service) (*http.Server, *Server, error) {
    server := &Server{
        db:    db,
        flags: flags,
        cache: cache,
    }
    // ...
}

// 4. Update main.go to pass it
httpServer, appServer, err := server.NewServer(db, flagsService, cache)
```

## Best Practices

1. **Always return errors**: Never panic in initialization code
2. **Accept interfaces**: Depend on `database.Service`, not `*service`
3. **Inject all dependencies**: Don't create them internally
4. **One constructor per type**: Keep it simple and predictable
5. **Close resources**: Use `defer` to close database, files, etc.
6. **Test with mocks**: Use mock implementations for unit tests
7. **Keep interfaces small**: Easier to mock and maintain

## Benefits Summary

| Aspect | Before (Singleton) | After (DI) |
|--------|-------------------|------------|
| **Testability** | Hard (global state) | Easy (inject mocks) |
| **Thread Safety** | Risky (shared state) | Safe (no shared state) |
| **Lifetime** | Hidden | Explicit |
| **Dependencies** | Implicit | Explicit |
| **Error Handling** | Panic | Return errors |
| **Coupling** | Tight | Loose |
| **Maintainability** | Low | High |

## References

- [Effective Go - Dependency Injection](https://golang.org/doc/effective_go)
- [Go Proverbs - Clear is better than clever](https://go-proverbs.github.io/)
- [testify/mock documentation](https://github.com/stretchr/testify)
