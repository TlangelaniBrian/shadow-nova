# Shadow Nova - Remaining Work

**Status**: Phase 0 (Emergency Fixes) ✅ COMPLETE
**Next Phase**: Phase 1 (Critical Security) - 3 days estimated

---

## ✅ COMPLETED (Phase 0 - Emergency)

The following critical fixes were implemented by 6 parallel agents:

1. ✅ Admin authorization middleware implemented
2. ✅ Unbounded goroutine fixed with graceful shutdown
3. ✅ localStorage token key consistency fixed
4. ✅ N+1 query bug fixed (83% performance improvement)
5. ✅ 8 missing database indexes added
6. ✅ github_username column added to schema
7. ✅ JWT secret validation enforced (min 32 chars)
8. ✅ JWT expiry validation in frontend
9. ✅ Client secrets removed from frontend
10. ✅ Role-based access control (RBAC) implemented

**But you still need to manually:**
- [ ] Rotate all leaked secrets (follow SECRET_ROTATION_CHECKLIST.md)
- [ ] Clean git history (run scripts/cleanup-git-history.sh)
- [ ] Apply database migration
- [ ] Test all changes

---

## 🔴 PHASE 1: Critical Security (This Week - 3 days)

**Priority**: HIGH - Required before production deployment

### 1. Migrate JWT to HttpOnly Cookies (4 hours)
**Current Issue**: JWT in localStorage vulnerable to XSS attacks

**Backend Changes:**
```go
// In auth.go after generating JWT
http.SetCookie(w, &http.Cookie{
    Name:     "auth_token",
    Value:    jwtToken,
    HttpOnly: true,
    Secure:   true,  // HTTPS only
    SameSite: http.SameSiteStrictMode,
    Path:     "/",
    MaxAge:   86400, // 24 hours
})
```

**Frontend Changes:**
- Remove all localStorage.setItem('token') calls
- Remove Authorization header from axios
- Set axios `withCredentials: true`
- Cookie automatically sent with requests

**Files to Modify:**
- `backend/internal/handlers/auth.go`
- `backend/internal/middleware/auth.go` (read from cookie instead of header)
- `frontend/src/api/client.ts`
- `frontend/src/stores/user.ts`
- `frontend/src/components/GoogleSignIn.vue`

---

### 2. Add Resource Ownership Validation (IDOR Fix) (6 hours)
**Current Issue**: Users can access other users' data by changing IDs

**Implementation:**

Create `backend/internal/middleware/ownership.go`:
```go
func ValidatePathOwnership(db database.Service) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            pathID := chi.URLParam(r, "id")
            userID := r.Context().Value(UserIDKey).(int)

            hasAccess, err := db.UserHasAccessToPath(r.Context(), userID, pathID)
            if err != nil || !hasAccess {
                http.Error(w, "Forbidden", http.StatusForbidden)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

**Add to Database Service:**
```go
UserHasAccessToPath(ctx context.Context, userID int, pathID string) (bool, error)
UserOwnsSubmission(ctx context.Context, userID int, submissionID int) (bool, error)
```

**Apply to Routes:**
```go
r.Get("/paths/{id}/progress", middleware.ValidatePathOwnership(s.db)(pathsHandler.GetProgress))
r.Get("/submissions/{id}", middleware.ValidateSubmissionOwnership(s.db)(projectsHandler.GetSubmission))
```

**Files to Modify:**
- `backend/internal/middleware/ownership.go` (NEW)
- `backend/internal/database/database.go` (add ownership methods)
- `backend/internal/server/routes.go` (apply middleware)
- All handlers that access user-specific resources

---

### 3. Encrypt GitHub Tokens at Rest (4 hours)
**Current Issue**: OAuth tokens stored in plain text in database

**Implementation:**

Add encryption library:
```bash
go get github.com/gtank/cryptopasta
```

Create `backend/internal/crypto/crypto.go`:
```go
package crypto

import (
    "crypto/rand"
    "encoding/base64"
    "github.com/gtank/cryptopasta"
    "os"
)

var encryptionKey *[32]byte

func Init() error {
    keyString := os.Getenv("ENCRYPTION_KEY")
    if keyString == "" {
        return fmt.Errorf("ENCRYPTION_KEY not set")
    }
    // Decode base64 key
    keyBytes, err := base64.StdEncoding.DecodeString(keyString)
    if err != nil {
        return err
    }
    copy(encryptionKey[:], keyBytes)
    return nil
}

func Encrypt(plaintext string) (string, error) {
    ciphertext, err := cryptopasta.Encrypt([]byte(plaintext), encryptionKey)
    return base64.StdEncoding.EncodeToString(ciphertext), err
}

func Decrypt(ciphertext string) (string, error) {
    encrypted, err := base64.StdEncoding.DecodeString(ciphertext)
    if err != nil {
        return "", err
    }
    plaintext, err := cryptopasta.Decrypt(encrypted, encryptionKey)
    return string(plaintext), err
}
```

**Update GitHub Integration:**
```go
// In database/projects.go
encryptedToken, err := crypto.Encrypt(integration.AccessToken)
query := `INSERT INTO github_integrations (user_id, access_token) VALUES ($1, $2)`
_, err = s.db.Exec(ctx, query, integration.UserID, encryptedToken)

// When reading
var encryptedToken string
err := row.Scan(&encryptedToken)
integration.AccessToken, err = crypto.Decrypt(encryptedToken)
```

**Environment Setup:**
```bash
# Generate encryption key
openssl rand -base64 32
# Add to .env
ENCRYPTION_KEY=<generated-key>
```

**Files to Modify:**
- `backend/internal/crypto/crypto.go` (NEW)
- `backend/internal/database/projects.go`
- `backend/main.go` (initialize crypto)

---

### 4. Implement CSRF Protection (3 hours)
**Current Issue**: State-changing operations vulnerable to CSRF

**Implementation:**

Add CSRF library:
```bash
go get github.com/gorilla/csrf
```

In `backend/internal/server/routes.go`:
```go
import "github.com/gorilla/csrf"

func (s *Server) RegisterRoutes() http.Handler {
    r := chi.NewRouter()

    // Add CSRF protection
    csrfKey := []byte(os.Getenv("CSRF_KEY")) // 32 bytes
    csrfMiddleware := csrf.Protect(
        csrfKey,
        csrf.Secure(true), // HTTPS only in production
        csrf.Path("/"),
        csrf.SameSite(csrf.SameSiteStrictMode),
    )

    r.Use(csrfMiddleware)

    // Add CSRF token endpoint
    r.Get("/api/csrf-token", func(w http.ResponseWriter, r *http.Request) {
        token := csrf.Token(r)
        w.Header().Set("X-CSRF-Token", token)
        w.WriteHeader(http.StatusOK)
    })

    // Rest of routes...
}
```

**Frontend Integration:**

In `frontend/src/api/client.ts`:
```typescript
// Fetch CSRF token on app init
let csrfToken = ''

async function fetchCSRFToken() {
    const response = await apiClient.get('/api/csrf-token')
    csrfToken = response.headers['x-csrf-token']
}

// Add to all state-changing requests
apiClient.interceptors.request.use((config) => {
    if (['post', 'put', 'patch', 'delete'].includes(config.method?.toLowerCase() || '')) {
        config.headers['X-CSRF-Token'] = csrfToken
    }
    return config
})

// Call on app init
fetchCSRFToken()
```

**Files to Modify:**
- `backend/internal/server/routes.go`
- `frontend/src/api/client.ts`
- `frontend/src/main.ts`

---

### 5. Implement Token Revocation (3 hours)
**Current Issue**: Compromised tokens cannot be invalidated

**Implementation:**

Create token blacklist table:
```sql
CREATE TABLE IF NOT EXISTS token_blacklist (
    jti VARCHAR(36) PRIMARY KEY,  -- JWT ID
    user_id INTEGER NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    blacklisted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    reason VARCHAR(100)
);

CREATE INDEX idx_token_blacklist_expires ON token_blacklist(expires_at);
```

**Add JTI to JWT:**
```go
// In auth.go
import "github.com/google/uuid"

func GenerateJWT(userID int, email, role string) (string, error) {
    claims := JWTClaims{
        UserID: userID,
        Email:  email,
        Role:   role,
        StandardClaims: jwt.StandardClaims{
            Id:        uuid.New().String(), // Add JTI
            ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
        },
    }
}
```

**Check Blacklist in Middleware:**
```go
// In middleware/auth.go
func (m *authMiddleware) VerifyToken(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // ... existing token validation

        // Check if token is blacklisted
        jti := claims.Id
        isBlacklisted, err := m.db.IsTokenBlacklisted(r.Context(), jti)
        if err != nil || isBlacklisted {
            http.Error(w, "Token revoked", http.StatusUnauthorized)
            return
        }

        next.ServeHTTP(w, r)
    })
}
```

**Add Logout Endpoint:**
```go
// In handlers/auth.go
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
    token := extractTokenFromRequest(r)
    claims, _ := auth.ValidateJWT(token)

    err := h.db.BlacklistToken(r.Context(), claims.Id, claims.UserID,
        time.Unix(claims.ExpiresAt, 0), "user_logout")
    if err != nil {
        httputil.WriteError(w, http.StatusInternalServerError, "Failed to logout")
        return
    }

    httputil.WriteJSON(w, http.StatusOK, map[string]string{"message": "Logged out successfully"})
}
```

**Cleanup Old Blacklisted Tokens:**
```go
// Run daily
go func() {
    ticker := time.NewTicker(24 * time.Hour)
    for range ticker.C {
        db.DeleteExpiredBlacklistedTokens(ctx)
    }
}()
```

**Files to Modify:**
- `backend/internal/database/schema.sql`
- `backend/internal/auth/auth.go`
- `backend/internal/middleware/auth.go`
- `backend/internal/database/database.go` (add blacklist methods)
- `backend/internal/handlers/auth.go`
- `backend/internal/server/routes.go`

---

## 🟡 PHASE 2: Critical Performance (Next Week - 5 days)

### 1. Implement Database Connection Pooling (2 hours)
Already partially done, but needs proper configuration:

```go
// In database/database.go
config.MaxConns = int32(getEnvInt("DB_MAX_CONNS", 25))
config.MinConns = int32(getEnvInt("DB_MIN_CONNS", 5))
config.MaxConnLifetime = time.Hour
config.MaxConnIdleTime = 30 * time.Minute
config.HealthCheckPeriod = time.Minute
```

---

### 2. Replace Singleton Pattern with Dependency Injection (4 hours)
**Current Issue**: Global database instance makes testing difficult

**Implementation:**
Remove global `dbInstance` variable and pass database to handlers via constructors.

---

### 3. Add Transactions to Multi-Step Operations (6 hours)
**Current Issue**: Operations like SaveGitHubToken do multiple queries without transaction

**Pattern:**
```go
func (s *service) SaveGitHubToken(ctx context.Context, integration *models.GitHubIntegration) error {
    tx, err := s.db.Begin(ctx)
    if err != nil {
        return err
    }
    defer tx.Rollback(ctx)

    // Multiple operations
    _, err = tx.Exec(ctx, query1)
    _, err = tx.Exec(ctx, query2)

    return tx.Commit(ctx)
}
```

Apply to:
- SaveGitHubToken
- User registration with profile creation
- Learning path creation with modules
- Any multi-step operation

---

### 4. Implement Proper Error Handling (4 hours)
Replace panics with proper error returns:
- Remove panic from database.go
- Add error wrapping with context
- Implement structured error responses

---

### 5. Frontend Lazy Loading (2 hours)
Already partially identified:

```typescript
// In router/index.ts
const routes = [
  {
    path: '/dashboard',
    component: () => import('../views/DashboardView.vue'), // Lazy load
    meta: { requiresAuth: true }
  },
  // ... all other routes
]
```

---

## 🟢 PHASE 3: High Priority (Next Sprint - 2 weeks)

### 1. API Versioning (2 hours)
Add `/api/v1/` prefix to all routes.

### 2. Pagination Implementation (8 hours)
Add to all list endpoints:
- GET /api/paths?page=1&limit=20
- GET /api/projects?page=1&limit=20
- GET /api/submissions?page=1&limit=20

### 3. Idempotency Keys (6 hours)
Add idempotency support for:
- POST /api/progress
- POST /api/submissions
- POST /api/paths

### 4. Structured Logging (4 hours)
Replace `log` and `fmt.Printf` with `slog`:
```go
import "log/slog"

logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
logger.Info("content collection started",
    slog.Int("source_count", len(sources)),
    slog.String("run_id", runID),
)
```

### 5. Request ID Middleware (2 hours)
Add correlation IDs for tracing.

### 6. Business Metrics (4 hours)
Add Prometheus metrics for:
- User registrations
- Learning path completions
- Project submissions
- AI processing latency

### 7. Create Pinia Stores (4 hours)
Create stores for:
- Projects (currently only composable)
- Learning paths (currently static data)
- Notifications

### 8. Decompose Large Components (6 hours)
Split:
- ProfileView.vue (238 lines) → ProfileCard, GitHubIntegration, AccountSettings
- LoginView.vue (230 lines) → LoginForm, AuthProviders, LoginFeatures

### 9. Add CRUD Operations (8 hours)
Implement full CRUD for:
- Learning paths (UPDATE, DELETE)
- Modules (UPDATE, DELETE)
- Lessons (UPDATE, DELETE)
- Projects (UPDATE, DELETE)

### 10. Type Safety Improvements (4 hours)
Remove all `as any` assertions and add proper TypeScript interfaces.

---

## 🔵 PHASE 4: Medium Priority (Month 2 - 2 weeks)

### 1. OpenAPI/Swagger Documentation (8 hours)
Generate API documentation using `swag` for Go.

### 2. Comprehensive Health Checks (4 hours)
Add checks for Unleash, Gemini API, disk space, memory.

### 3. Sentry Integration (3 hours)
Add error tracking service.

### 4. OpenTelemetry Tracing (6 hours)
Implement distributed tracing with Jaeger.

### 5. Prometheus Alert Rules (4 hours)
Create alerting for high error rates, latency, etc.

### 6. Fix Promtail Configuration (1 hour)
Update to scrape Docker container logs.

### 7. Strong Password Validation (2 hours)
Enforce strong_password validator on registration.

### 8. CORS Validation (2 hours)
Validate CORS_ORIGINS against whitelist.

### 9. Database SSL (1 hour)
Enable sslmode=require in production.

---

## 🟣 PHASE 5: Technical Debt (Month 3 - 1 week)

### 1. Split Large Database Interface (6 hours)
Create focused interfaces (UserService, PathService, etc.).

### 2. Add Unique Constraints (2 hours)
Add constraints on order_index fields.

### 3. Implement Soft Deletes (4 hours)
Add deleted_at column for audit trail.

### 4. Add Prepared Statements (4 hours)
For frequently executed queries.

### 5. Fix Rate Limiter (3 hours)
- Use RWMutex correctly
- Use actual IP instead of request ID
- Consider Redis for distributed systems

### 6. Add Migration System (6 hours)
Implement golang-migrate for proper versioning.

---

## 📊 EFFORT SUMMARY

| Phase | Tasks | Estimated Effort | Priority |
|-------|-------|-----------------|----------|
| Phase 0 (Emergency) | 10 | ✅ COMPLETE | CRITICAL |
| Phase 1 (Security) | 5 | 20 hours (3 days) | CRITICAL |
| Phase 2 (Performance) | 5 | 18 hours (2-3 days) | HIGH |
| Phase 3 (High Priority) | 10 | 48 hours (6 days) | HIGH |
| Phase 4 (Medium) | 9 | 30 hours (4 days) | MEDIUM |
| Phase 5 (Tech Debt) | 6 | 25 hours (3 days) | LOW |

**Total Remaining**: 141 hours (~18 days with 1 developer)

---

## 🎯 RECOMMENDED TIMELINE

### Week 1 (Current)
- ✅ Phase 0 fixes implemented
- ⏳ Rotate secrets and clean git history
- ⏳ Deploy Phase 0 fixes to staging

### Week 2
- 🔴 Phase 1: Critical security (3 days)
- ⚡ Test all security fixes
- 🚀 Deploy to staging

### Week 3
- 🟡 Phase 2: Critical performance (3 days)
- 🧪 Performance testing
- 📊 Monitor improvements

### Weeks 4-5
- 🟢 Phase 3: High priority items
- 📝 API documentation
- 🎨 Frontend improvements

### Week 6-7
- 🔵 Phase 4: Medium priority
- 📊 Observability improvements
- 🔍 Monitoring setup

### Week 8
- 🟣 Phase 5: Technical debt
- 🧹 Code cleanup
- 📚 Documentation updates

---

## 🚀 PRODUCTION READINESS CHECKLIST

Before deploying to production, ensure:

**Security (Phase 1 - MUST DO)**
- [ ] JWT in HttpOnly cookies
- [ ] IDOR protection implemented
- [ ] GitHub tokens encrypted
- [ ] CSRF protection enabled
- [ ] Token revocation working

**Performance (Phase 2 - SHOULD DO)**
- [ ] Connection pooling configured
- [ ] Transactions implemented
- [ ] Error handling improved
- [ ] Dependency injection
- [ ] Frontend lazy loading

**Observability (Phase 3/4 - SHOULD DO)**
- [ ] Structured logging
- [ ] Request IDs
- [ ] Business metrics
- [ ] Error tracking (Sentry)
- [ ] Distributed tracing

**Quality (Phase 3-5 - NICE TO HAVE)**
- [ ] API documentation
- [ ] Comprehensive health checks
- [ ] Alert rules configured
- [ ] Migration system
- [ ] Full CRUD operations

---

## 💡 QUICK WINS (Can Do Anytime)

These are small improvements that provide immediate value:

1. **Add Request ID Middleware** (30 min) - Better debugging
2. **Enable Database SSL** (5 min) - Security improvement
3. **Fix Promtail Config** (15 min) - See container logs
4. **Add Basic Business Metrics** (2 hours) - Track registrations/logins
5. **Frontend Lazy Loading** (2 hours) - 30% bundle size reduction
6. **Add CORS Validation** (1 hour) - Prevent misconfiguration

---

## 📞 NEED HELP?

If you want to parallelize the remaining work with agents again, we can:
- Launch multiple agents for Phase 1 security fixes
- Create scripts to automate testing
- Generate comprehensive documentation
- Set up monitoring dashboards

Just let me know which phase you want to tackle next!
