# Shadow Nova - Comprehensive Architectural Audit Report

**Date**: February 11, 2026
**Auditors**: 6 Parallel AI Agents (Backend, Frontend, Database, Security, API, Observability)
**Codebase**: 44 Go files, 70 Vue/TypeScript files
**Analysis Method**: Recursive Language Model (RLM) parallel analysis

---

## Executive Summary

This comprehensive audit identified **78 architectural issues** across 6 categories. The codebase has solid foundations but contains **critical security vulnerabilities** and **performance bottlenecks** that will cause production incidents if not addressed.

### Severity Breakdown

| Severity | Count | Categories |
|----------|-------|------------|
| **CRITICAL** | 13 | Security (3), Backend (5), Database (3), Frontend (2) |
| **HIGH** | 21 | Security (6), Database (5), Frontend (4), API (3), Observability (3) |
| **MEDIUM** | 32 | All categories |
| **LOW** | 12 | All categories |

### Top 5 Critical Issues

1. **🔒 Secrets Leaked in Git** - All OAuth secrets, API keys exposed in commit history
2. **🚨 Unbounded Goroutine** - Background collector runs forever, no graceful shutdown
3. **💥 Missing Admin Authorization** - Any authenticated user can access admin endpoints
4. **⚡ N+1 Query Bug** - `GetLearningPath` makes 1+N database queries
5. **🎯 localStorage Token Chaos** - Two different token keys, scattered across 9+ files

---

## 1. SECURITY VULNERABILITIES

### CRITICAL (Priority 0 - Fix Today)

#### 1.1 Secrets in Git History
**Impact**: Complete authentication system compromise
**Files**: `.env`, `backend/.env`, `frontend/.env`
**Evidence**:
```
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
GOOGLE_CLIENT_ID=54718351338-tbied6l38ldbiqj4nmghqgi9l4lf1ge9.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=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
GITHUB_CLIENT_ID=Ov23liZmK8EPzEymFddm
GITHUB_CLIENT_SECRET=6a416aec07a071e39a2b2cf5fa00c50fc9d8b269
GEMINI_API_KEY=AIzaSyDWQ5k3_eBmA1qLWhnfSV3OtAglgD7X_f0
```
**Action**:
1. Rotate ALL secrets immediately (revoke at provider consoles)
2. Remove .env files from git: `git filter-repo --path .env --invert-paths`
3. Add `.env` to `.gitignore` (already done)
4. Use environment variables or AWS Secrets Manager in production

#### 1.2 Missing Admin Authorization Middleware
**Impact**: Any authenticated user can modify system settings
**File**: `backend/internal/server/routes.go:144-148`
```go
// r.Use(adminMiddleware) // TODO: Implement Admin Middleware
r.Post("/admin/settings/collector", adminHandler.UpdateCollectorFrequency)
```
**Action**:
1. Add `role` field to users table: `ALTER TABLE users ADD COLUMN role VARCHAR(20) DEFAULT 'user';`
2. Create admin middleware:
```go
func AdminOnly(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        userRole := r.Context().Value("user_role").(string)
        if userRole != "admin" {
            http.Error(w, "Forbidden", http.StatusForbidden)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```
3. Apply to admin routes immediately

#### 1.3 Weak JWT Secret with Unsafe Fallback
**Impact**: Token forgery if JWT_SECRET not set
**File**: `backend/internal/auth/auth.go:26-29`
```go
if jwtSecret == "" {
    jwtSecret = "default-secret-change-in-production"
}
```
**Action**:
```go
if jwtSecret == "" {
    log.Fatal("JWT_SECRET environment variable is required")
}
if len(jwtSecret) < 32 {
    log.Fatal("JWT_SECRET must be at least 32 characters")
}
```

### HIGH (Priority 1 - This Week)

#### 1.4 Client Secrets Exposed to Frontend
**File**: `frontend/.env:4-5`
```
VITE_GOOGLE_CLIENT_SECRET=...
VITE_GITHUB_CLIENT_SECRET=...
```
**Action**: Remove all `VITE_*_SECRET` variables. Secrets should NEVER be in frontend. Use backend proxy for OAuth flows.

#### 1.5 JWT Stored in localStorage (XSS Risk)
**Files**: `frontend/src/stores/user.ts:41`, `frontend/src/api/client.ts:14`
**Action**: Migrate to HttpOnly cookies
```go
// Backend
http.SetCookie(w, &http.Cookie{
    Name:     "auth_token",
    Value:    jwtToken,
    HttpOnly: true,
    Secure:   true,
    SameSite: http.SameSiteStrictMode,
    MaxAge:   86400,
})
```

#### 1.6 IDOR Vulnerabilities - No Resource Ownership Validation
**Files**: `backend/internal/handlers/paths.go:31-40`, `progress.go:59-74`
**Action**: Add ownership checks:
```go
func (h *PathsHandler) Get(w http.ResponseWriter, r *http.Request) {
    pathID := chi.URLParam(r, "id")
    userID := r.Context().Value("user_id").(int)

    // Verify user owns or has access to this path
    if !h.db.UserHasAccessToPath(r.Context(), userID, pathID) {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }
    // ... proceed
}
```

#### 1.7 GitHub Tokens Stored in Plain Text
**File**: `backend/internal/database/schema.sql:70`
```sql
access_token VARCHAR(255) NOT NULL,
```
**Action**: Encrypt tokens before storage
```go
import "github.com/gtank/cryptopasta"

encryptedToken, err := cryptopasta.Encrypt([]byte(token), key)
// Store encryptedToken
```

#### 1.8 No CSRF Protection
**File**: `backend/internal/middleware/security.go:101`
**Impact**: State-changing operations vulnerable to CSRF
**Action**: Implement CSRF middleware using `gorilla/csrf`

#### 1.9 No Token Revocation Mechanism
**Impact**: Compromised tokens cannot be invalidated
**Action**:
1. Create token blacklist table
2. Add logout endpoint that blacklists token
3. Check blacklist in auth middleware

---

## 2. BACKEND ARCHITECTURE

### CRITICAL (Priority 0)

#### 2.1 Unbounded Goroutine - Production Killer
**File**: `backend/internal/server/routes.go:63-93`
```go
go func() {
    ctx := context.Background()  // Never cancelled!
    for {
        time.Sleep(interval)
        collectorService.CollectAll(ctx)
    }
}()
```
**Impact**:
- Cannot gracefully shutdown
- Goroutine leak
- DB connections held forever
- Panic crashes entire app

**Action**:
```go
// Add to Server struct
ctx, cancel := context.WithCancel(context.Background())
s.collectorCancel = cancel

go func() {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()
    for {
        select {
        case <-ticker.C:
            collectorService.CollectAll(s.ctx)
        case <-s.ctx.Done():
            return
        }
    }
}()

// In shutdown
func (s *Server) Shutdown() {
    s.collectorCancel()
    s.db.Close()
}
```

#### 2.2 Singleton Database Pattern
**File**: `backend/internal/database/database.go:58-66`
```go
var dbInstance *service
func New() Service {
    if dbInstance != nil {
        return dbInstance
    }
```
**Impact**: Cannot test, hidden global state, race condition on init
**Action**: Use dependency injection
```go
func New(databaseURL string) (Service, error) {
    // Create new instance every time
    // Pass to handlers via constructor
}
```

#### 2.3 Panics in Library Code
**File**: `backend/internal/database/database.go:76,81`
```go
if err != nil {
    panic(fmt.Sprintf("Unable to parse database URL: %v", err))
}
```
**Action**: Return errors, handle in main()
```go
func New(databaseURL string) (Service, error) {
    config, err := pgxpool.ParseConfig(databaseURL)
    if err != nil {
        return nil, fmt.Errorf("invalid database URL: %w", err)
    }
    // ...
}
```

#### 2.4 No Transaction Support
**File**: `backend/internal/database/projects.go:107-133`
**Impact**: Partial failures leave inconsistent data
**Action**: Wrap multi-step operations
```go
func (s *service) SaveGitHubToken(ctx context.Context, integration *models.GitHubIntegration) error {
    tx, err := s.db.Begin(ctx)
    if err != nil {
        return err
    }
    defer tx.Rollback(ctx)

    // Insert integration
    // Update user

    return tx.Commit(ctx)
}
```

#### 2.5 Hardcoded Admin Email
**File**: `backend/internal/server/server.go:58`
```go
email := "mrbtmkhabela@gmail.com"
```
**Action**: Use environment variable

### HIGH (Priority 1)

#### 2.6 Large Database Interface (25+ Methods)
**File**: `backend/internal/database/database.go:13-52`
**Action**: Split into focused interfaces
```go
type UserService interface {
    CreateUser(...) error
    GetUserByEmail(...) (*User, error)
}

type PathService interface {
    CreateLearningPath(...) error
    GetLearningPaths(...) ([]Path, error)
}
```

#### 2.7 Race Condition in Rate Limiter
**File**: `backend/internal/middleware/security.go:48-74`
```go
rl.mu.Lock()  // Full lock for read operation
```
**Action**: Use RWMutex correctly
```go
rl.mu.RLock()  // Use read lock for lookups
v, exists := rl.visitors[ip]
rl.mu.RUnlock()
```

#### 2.8 Connection Pool Not Configured
**File**: `backend/internal/database/database.go:74-82`
**Action**:
```go
config.MaxConns = 25
config.MinConns = 5
config.MaxConnLifetime = time.Hour
config.MaxConnIdleTime = 30 * time.Minute
config.HealthCheckPeriod = time.Minute
```

---

## 3. DATABASE DESIGN

### CRITICAL (Priority 0)

#### 3.1 N+1 Query Bug
**File**: `backend/internal/database/paths.go:76-104`
```go
for i := range p.Modules {
    lessonQuery := `SELECT ... FROM lessons WHERE module_id = $1`
    lRows, err := s.db.Query(ctx, lessonQuery, p.Modules[i].ID)
```
**Impact**: 11 queries for path with 10 modules (severe performance degradation)
**Action**: Use single JOIN query
```go
query := `
    SELECT
        m.id, m.title, m.description, m.order_index,
        l.id, l.title, l.content_type, l.content_url, l.order_index
    FROM modules m
    LEFT JOIN lessons l ON l.module_id = m.id
    WHERE m.path_id = $1
    ORDER BY m.order_index, l.order_index
`
// Group results by module
```

#### 3.2 Schema Mismatch - Missing Column
**File**: `backend/internal/database/projects.go:124`
```go
updateUserQuery := `UPDATE users SET github_username = $1 WHERE id = $2`
```
**Issue**: `github_username` column doesn't exist in schema
**Action**:
```sql
ALTER TABLE users ADD COLUMN github_username VARCHAR(100);
```

#### 3.3 No Migration System
**File**: `backend/internal/database/database.go:114-138`
**Impact**: Running raw schema.sql with DROP TABLE, no versioning, no rollback
**Action**: Use golang-migrate
```bash
migrate create -ext sql -dir migrations -seq initial_schema
migrate create -ext sql -dir migrations -seq add_github_username
```

### HIGH (Priority 1)

#### 3.4 Missing Indexes (8 Critical)
**File**: `backend/internal/database/schema.sql`
```sql
-- Add these indexes
CREATE INDEX idx_modules_path_id ON modules(path_id);
CREATE INDEX idx_lessons_module_id ON lessons(module_id);
CREATE INDEX idx_user_progress_user_completed ON user_progress(user_id, completed);
CREATE INDEX idx_content_items_processed ON content_items(processed_by_ai) WHERE processed_by_ai = FALSE;
CREATE INDEX idx_content_items_source_id ON content_items(source_id);
CREATE INDEX idx_project_submissions_project_id ON project_submissions(project_id);
CREATE INDEX idx_project_submissions_status ON project_submissions(status, submitted_at DESC);
CREATE INDEX idx_github_integrations_github_user_id ON github_integrations(github_user_id);
```

#### 3.5 No Pagination
**Files**: `backend/internal/database/paths.go:9-34`, `projects.go:11-36`
**Impact**: Unbounded result sets risk DoS
**Action**: Add LIMIT and OFFSET to all list queries
```go
func (s *service) GetLearningPaths(ctx context.Context, limit, offset int) ([]models.LearningPath, error) {
    query := `SELECT ... FROM learning_paths ORDER BY created_at DESC LIMIT $1 OFFSET $2`
    rows, err := s.db.Query(ctx, query, limit, offset)
```

#### 3.6 Missing Unique Constraints
**File**: `backend/internal/database/schema.sql:21-41`
```sql
ALTER TABLE modules ADD CONSTRAINT unique_module_order UNIQUE(path_id, order_index);
ALTER TABLE lessons ADD CONSTRAINT unique_lesson_order UNIQUE(module_id, order_index);
```

---

## 4. FRONTEND ARCHITECTURE

### CRITICAL (Priority 0)

#### 4.1 localStorage Token Inconsistency (Security Risk)
**Files**: 9+ locations
**Issue**: Two different token keys used
- `frontend/src/router/index.ts:96` - Uses `'token'`
- `frontend/src/components/GoogleSignIn.vue:66` - Uses `'auth_token'`
**Impact**: Authentication state corruption, security confusion
**Action**: Centralize in Pinia store, use single key

#### 4.2 No Lazy Loading (Performance)
**File**: `frontend/src/router/index.ts:2-12`
```typescript
import DashboardView from '../views/DashboardView.vue'
import LoginView from '../views/LoginView.vue'
// ... 9 more static imports
```
**Impact**: 150-200KB initial bundle bloat, slow first load
**Action**:
```typescript
component: () => import('../views/DashboardView.vue')
```

### HIGH (Priority 1)

#### 4.3 Projects State Lost (No Pinia Store)
**File**: `frontend/src/composables/useProjects.ts:16`
**Impact**: State lost on unmount, refetched unnecessarily
**Action**: Create `projectsStore.ts`

#### 4.4 Weak Auth Guard
**File**: `frontend/src/router/index.ts:95-106`
```typescript
if (requiresAuth && !token) {
    next('/login')
}
```
**Issue**: Only checks token existence, not JWT expiry
**Action**: Add JWT validation
```typescript
import jwtDecode from 'jwt-decode'

const decoded = jwtDecode(token)
if (decoded.exp * 1000 < Date.now()) {
    localStorage.removeItem('token')
    next('/login')
}
```

#### 4.5 Large Components Need Decomposition
- `ProfileView.vue` - 238 lines (split into ProfileCard, GitHubIntegration, AccountSettings)
- `LoginView.vue` - 230 lines (split into LoginForm, AuthProviders, LoginFeatures)

#### 4.6 Unsafe Type Assertions
**Files**: Multiple
```typescript
inject('unleash') as any
window.google: any
catch (error: any)
```
**Action**: Define proper types
```typescript
interface UnleashClient {
    isEnabled(flag: string): boolean
}
const unleash = inject<UnleashClient>('unleash')
```

---

## 5. API DESIGN

### CRITICAL (Priority 0)

#### 5.1 No API Versioning
**File**: `backend/internal/server/routes.go:40`
**Impact**: Breaking changes affect all clients
**Action**: Add `/api/v1/` prefix to all routes

### HIGH (Priority 1)

#### 5.2 No Idempotency Keys
**Files**: `handlers/progress.go:22-40`, `projects.go:53-80`
**Impact**: Duplicate submissions on network retry
**Action**:
```go
func (h *ProgressHandler) UpdateProgress(w http.ResponseWriter, r *http.Request) {
    idempotencyKey := r.Header.Get("Idempotency-Key")
    if idempotencyKey != "" {
        // Check if already processed
        if cached, exists := h.cache.Get(idempotencyKey); exists {
            httputil.WriteJSON(w, http.StatusOK, cached)
            return
        }
    }
    // Process and cache result
}
```

#### 5.3 Missing CRUD Operations
**Issue**: No UPDATE or DELETE for paths, modules, lessons, projects
**Action**: Implement full REST CRUD

#### 5.4 Inconsistent Response Structures
**File**: `backend/internal/handlers/auth.go:79-87`
```go
map[string]interface{}{
    "token": jwtToken,
    "user": map[string]string{...}
}
```
**Action**: Define typed response models
```go
type LoginResponse struct {
    Token string `json:"token"`
    User  User   `json:"user"`
}
```

#### 5.5 No OpenAPI/Swagger Documentation
**Action**: Generate OpenAPI spec using `swag` or manual YAML

---

## 6. OBSERVABILITY

### CRITICAL (Priority 0)

#### 6.1 No Structured Logging
**Files**: All `.go` files
```go
log.Println("Running initial content collection...")
fmt.Printf("Warning: Failed to update...")
```
**Action**: Implement slog (Go 1.21+)
```go
import "log/slog"

logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
logger.Info("content collection started",
    slog.Int("source_count", len(sources)),
    slog.String("run_id", runID),
)
```

#### 6.2 Missing Correlation IDs
**Impact**: Cannot trace single request across services
**Action**: Add RequestID middleware
```go
func RequestIDMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID := uuid.New().String()
        ctx := context.WithValue(r.Context(), "request_id", requestID)
        w.Header().Set("X-Request-ID", requestID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### HIGH (Priority 1)

#### 6.3 Only Technical Metrics (No Business Metrics)
**File**: `backend/internal/middleware/prometheus.go`
**Missing**:
- User registrations/logins
- Learning path completions
- Project submissions
- AI processing latency

**Action**:
```go
var userRegistrations = promauto.NewCounter(prometheus.CounterOpts{
    Name: "user_registrations_total",
})
```

#### 6.4 No Distributed Tracing
**Impact**: Cannot debug slow requests
**Action**: Implement OpenTelemetry + Jaeger

#### 6.5 No Error Tracking Service
**Impact**: Production errors only in logs
**Action**: Integrate Sentry
```go
sentry.Init(sentry.ClientOptions{
    Dsn: os.Getenv("SENTRY_DSN"),
})
```

#### 6.6 Incomplete Health Check
**File**: `backend/internal/server/server.go:96-98`
**Current**: Only checks database
**Missing**: Unleash, Gemini, disk space, memory
**Action**: Add comprehensive checks

#### 6.7 No Alert Rules
**File**: `frontend/observability/prometheus.yml`
**Action**: Create alert rules
```yaml
groups:
  - name: shadow_nova_alerts
    rules:
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.05
```

#### 6.8 PII in Logs Risk
**Files**: Multiple handlers
**Risk**: Emails, tokens logged without redaction
**Action**: Implement log sanitization

#### 6.9 Promtail Misconfigured
**File**: `frontend/observability/promtail-config.yml`
**Issue**: Only reads `/var/log`, not Docker container logs
**Action**:
```yaml
scrape_configs:
  - job_name: docker
    docker_sd_configs:
      - host: unix:///var/run/docker.sock
```

---

## PRIORITY ACTION PLAN

### Phase 0: Emergency (Today)

1. **Rotate all secrets** - Google, GitHub, Gemini API keys
2. **Remove .env from git history** - `git filter-repo --path .env --invert-paths`
3. **Implement admin middleware** - Protect admin endpoints
4. **Fix unbounded goroutine** - Add graceful shutdown to collector
5. **Fix localStorage token keys** - Consolidate to single key

### Phase 1: Critical Security (This Week)

6. Migrate JWT to HttpOnly cookies
7. Add IDOR resource ownership validation
8. Encrypt GitHub tokens at rest
9. Implement CSRF protection
10. Add token revocation mechanism

### Phase 2: Critical Performance (Next Week)

11. Fix N+1 query in GetLearningPath (use JOIN)
12. Add missing 8 database indexes
13. Implement lazy loading for frontend routes
14. Add database connection pooling config
15. Replace singleton pattern with dependency injection

### Phase 3: High Priority (Next Sprint)

16. Implement proper migration system (golang-migrate)
17. Add pagination to all list endpoints
18. Add structured logging (slog)
19. Implement request ID correlation
20. Add business metrics to Prometheus
21. Create Pinia stores for projects and paths
22. Decompose large components (ProfileView, LoginView)
23. Add API versioning (/api/v1/)
24. Implement idempotency keys
25. Wrap multi-step operations in transactions

### Phase 4: Medium Priority (Month 2)

26. Add OpenAPI/Swagger documentation
27. Implement full CRUD operations
28. Add comprehensive health checks
29. Integrate Sentry for error tracking
30. Implement OpenTelemetry tracing
31. Add Prometheus alert rules
32. Fix Promtail configuration
33. Add strong password validation
34. Implement proper CORS validation

### Phase 5: Technical Debt (Month 3)

35. Split large database interface
36. Add unique constraints on order fields
37. Implement soft deletes
38. Add prepared statements
39. Type safety improvements (remove `any`)
40. Add missing TypeScript interfaces

---

## RISK ASSESSMENT

### Production Blockers (Cannot Deploy)

1. ❌ Secrets in git history
2. ❌ Missing admin authorization
3. ❌ Unbounded goroutine (cannot shutdown)
4. ❌ Client secrets in frontend
5. ❌ N+1 query bug (will crash under load)

### Production Risks (Deploy with Caution)

6. ⚠️ JWT in localStorage (XSS risk)
7. ⚠️ IDOR vulnerabilities (data leakage)
8. ⚠️ No transaction support (data corruption risk)
9. ⚠️ No token revocation (compromised tokens valid forever)
10. ⚠️ Missing indexes (slow queries)

### Operational Debt (Fix Soon)

11. 📊 No structured logging (debugging nightmare)
12. 📊 No distributed tracing (cannot diagnose slow requests)
13. 📊 No error tracking (blind to production errors)
14. 📊 No business metrics (cannot measure success)
15. 📊 No alerting (reactive instead of proactive)

---

## ESTIMATED EFFORT

| Phase | Issues | Effort | Timeline |
|-------|--------|--------|----------|
| Phase 0 (Emergency) | 5 | 1 day | Today |
| Phase 1 (Security) | 5 | 3 days | This week |
| Phase 2 (Performance) | 5 | 5 days | Next week |
| Phase 3 (High Priority) | 10 | 2 weeks | Next sprint |
| Phase 4 (Medium Priority) | 9 | 2 weeks | Month 2 |
| Phase 5 (Tech Debt) | 6 | 1 week | Month 3 |

**Total**: 40 issues, ~6 weeks with 1 full-time developer

---

## CONCLUSION

Shadow Nova has a solid architectural foundation but contains **critical security vulnerabilities** and **performance bottlenecks** that make it unsuitable for production deployment in its current state.

### Strengths ✅
- Good use of parameterized queries (SQL injection protection)
- Modern tech stack (Go, Vue 3, PostgreSQL)
- Comprehensive observability setup (PLG stack configured)
- Clean package structure
- OAuth integration implemented

### Critical Weaknesses ❌
- Secrets leaked in git history
- Missing authorization on admin endpoints
- Goroutine management issues
- N+1 query bug
- No migration system
- Weak authentication token storage
- Missing business observability

### Recommendation

**DO NOT DEPLOY TO PRODUCTION** until at minimum Phase 0 and Phase 1 are complete. The leaked secrets alone represent a complete compromise of the authentication system.

With focused effort over the next 2-3 weeks, this codebase can be production-ready. The architectural foundations are sound—the issues are primarily in the implementation details and security hardening.

---

## APPENDIX: Full Agent Reports

Detailed findings from each agent are available in:
- `/private/tmp/claude-502/.../tasks/a94cafa.output` - Backend Architecture
- `/private/tmp/claude-502/.../tasks/a314b72.output` - Frontend Architecture
- `/private/tmp/claude-502/.../tasks/ad1a7cc.output` - Database Schema
- `/private/tmp/claude-502/.../tasks/af2eaa7.output` - Security Vulnerabilities
- `/private/tmp/claude-502/.../tasks/a4169ab.output` - API Design
- `/private/tmp/claude-502/.../tasks/a1e49f3.output` - Observability

**End of Report**
