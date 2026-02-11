# Dependency Injection - Quick Reference

## Creating a New Handler

```go
// 1. Define handler struct with dependencies
type MyHandler struct {
    db    database.Service
    cache cache.Service  // Optional: add more dependencies
}

// 2. Create constructor
func NewMyHandler(db database.Service, cache cache.Service) *MyHandler {
    return &MyHandler{
        db:    db,
        cache: cache,
    }
}

// 3. Use dependencies in methods
func (h *MyHandler) GetData(w http.ResponseWriter, r *http.Request) {
    data, err := h.db.GetData(r.Context())
    if err != nil {
        httputil.WriteError(w, http.StatusInternalServerError, "Failed to get data")
        return
    }
    httputil.WriteSuccess(w, "Data retrieved", data)
}
```

## Registering Handler Routes

```go
// In server/routes.go
func (s *Server) RegisterRoutes() (http.Handler, error) {
    r := chi.NewRouter()

    // Create handler with injected dependencies
    myHandler := handlers.NewMyHandler(s.db, s.cache)

    // Register routes
    r.Get("/api/my-data", myHandler.GetData)

    return r, nil
}
```

## Adding a New Dependency

```go
// 1. Initialize in main.go
func main() {
    // ... existing code

    cache, err := cache.New()
    if err != nil {
        log.Fatalf("Failed to initialize cache: %v", err)
    }
    defer cache.Close()

    // Pass to server
    httpServer, appServer, err := server.NewServer(db, flagsService, cache)
    // ...
}

// 2. Add to Server struct (server/server.go)
type Server struct {
    port  string
    db    database.Service
    flags flags.Service
    cache cache.Service  // New dependency
}

// 3. Update NewServer signature
func NewServer(db database.Service, flags flags.Service, cache cache.Service) (*http.Server, *Server, error) {
    server := &Server{
        port:  port,
        db:    db,
        flags: flags,
        cache: cache,  // Assign new dependency
    }
    return httpServer, server, nil
}
```

## Testing with Mocks

```go
// 1. Create mock (can use testify/mock)
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

// 2. Use in test
func TestMyHandler_GetData(t *testing.T) {
    // Setup mock
    mockDB := new(MockDatabase)
    mockDB.On("GetUserByEmail", mock.Anything, "test@example.com").
        Return(&models.User{ID: 1, Email: "test@example.com"}, nil)

    // Create handler with mock
    handler := handlers.NewMyHandler(mockDB)

    // Create test request
    req := httptest.NewRequest("GET", "/api/data", nil)
    w := httptest.NewRecorder()

    // Execute
    handler.GetData(w, req)

    // Assert
    assert.Equal(t, http.StatusOK, w.Code)
    mockDB.AssertExpectations(t)
}
```

## Common Patterns

### Return Errors, Don't Panic
```go
// ❌ Bad
func New() Service {
    db, err := connect()
    if err != nil {
        panic(err)  // Don't do this
    }
    return &service{db: db}
}

// ✅ Good
func New() (Service, error) {
    db, err := connect()
    if err != nil {
        return nil, fmt.Errorf("failed to connect: %w", err)
    }
    return &service{db: db}, nil
}
```

### Depend on Interfaces
```go
// ❌ Bad - depends on concrete type
type Handler struct {
    db *database.service  // Concrete type
}

// ✅ Good - depends on interface
type Handler struct {
    db database.Service  // Interface
}
```

### Inject All Dependencies
```go
// ❌ Bad - creates dependency internally
func NewHandler() *Handler {
    return &Handler{
        db: database.New(),  // Hidden dependency
    }
}

// ✅ Good - receives dependency
func NewHandler(db database.Service) *Handler {
    return &Handler{
        db: db,  // Injected dependency
    }
}
```

### Single Responsibility Constructors
```go
// ❌ Bad - constructor does too much
func NewHandler(db database.Service) *Handler {
    h := &Handler{db: db}
    h.initialize()      // Setup
    h.loadConfig()      // More setup
    h.startBackground() // Side effects
    return h
}

// ✅ Good - constructor just creates
func NewHandler(db database.Service) *Handler {
    return &Handler{db: db}
}

// Separate initialization if needed
func (h *Handler) Initialize() error {
    // Do setup here
    return nil
}
```

## Environment Variables

### Database Configuration
```bash
# Connection Pool
DB_MAX_CONNS=25              # Max concurrent connections
DB_MIN_CONNS=5               # Min idle connections
DB_CONNECT_TIMEOUT=5         # Connection timeout (seconds)
DB_TX_TIMEOUT=30s            # Transaction timeout

# Connection String
DATABASE_URL=postgres://user:pass@host:5432/dbname?sslmode=disable
```

## Checklist for New Dependencies

- [ ] Create interface defining the service contract
- [ ] Implement service with `New() (Service, error)` constructor
- [ ] Initialize in `main.go` with error handling
- [ ] Add `defer service.Close()` if needed
- [ ] Add to `Server` struct
- [ ] Update `NewServer()` signature
- [ ] Pass to handlers that need it
- [ ] Create mock implementation for testing
- [ ] Add environment variable configuration
- [ ] Document in relevant `*.md` files

## Common Mistakes

### 1. Creating Dependencies Inside Handlers
```go
// ❌ Bad
func (h *Handler) Process(w http.ResponseWriter, r *http.Request) {
    db := database.New()  // Don't do this
    // ...
}

// ✅ Good
func (h *Handler) Process(w http.ResponseWriter, r *http.Request) {
    // Use h.db that was injected
    data, err := h.db.GetData(r.Context())
    // ...
}
```

### 2. Forgetting to Close Resources
```go
// ❌ Bad
func main() {
    db, err := database.New()
    // Missing defer db.Close()
    server.NewServer(db, flags)
}

// ✅ Good
func main() {
    db, err := database.New()
    defer db.Close()  // Always close
    server.NewServer(db, flags)
}
```

### 3. Ignoring Errors
```go
// ❌ Bad
func main() {
    db, _ := database.New()  // Ignoring error
}

// ✅ Good
func main() {
    db, err := database.New()
    if err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
}
```

### 4. Using Global Variables
```go
// ❌ Bad
var globalDB database.Service

func init() {
    globalDB = database.New()
}

// ✅ Good
func main() {
    db, err := database.New()
    if err != nil {
        log.Fatalf("Failed: %v", err)
    }
    defer db.Close()
    // Pass db to components
}
```

## Quick Reference Card

| Action | Command/Pattern |
|--------|----------------|
| Create service | `func New() (Service, error)` |
| Create handler | `func NewHandler(db Service) *Handler` |
| Use in routes | `handler := handlers.NewHandler(s.db)` |
| Test with mock | `mockDB := new(MockDatabase)` |
| Set mock expectation | `mockDB.On("Method", args).Return(result, err)` |
| Initialize in main | `db, err := database.New(); defer db.Close()` |
| Check connection | `if err := db.Ping(ctx); err != nil { ... }` |

## Further Reading

- [DEPENDENCY_INJECTION.md](./DEPENDENCY_INJECTION.md) - Full architecture guide
- [DEPENDENCY_INJECTION_CHANGES.md](./DEPENDENCY_INJECTION_CHANGES.md) - Implementation details
- [testify/mock](https://github.com/stretchr/testify) - Mocking framework
- [Go Proverbs](https://go-proverbs.github.io/) - Go best practices
