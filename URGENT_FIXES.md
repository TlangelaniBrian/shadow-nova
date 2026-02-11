# 🚨 URGENT FIXES - Shadow Nova

**Status**: ❌ DO NOT DEPLOY TO PRODUCTION
**Critical Issues**: 13 found
**Estimated Fix Time**: 1-3 days

---

## IMMEDIATE ACTIONS (TODAY - Before Any Deployment)

### 1. 🔒 Rotate All Leaked Secrets (30 minutes)

**Issue**: All secrets exposed in git commits `fbc0f39` and `120da60`

**Compromised Credentials**:
```
Google OAuth: 54718351338-tbied6l38ldbiqj4nmghqgi9l4lf1ge9.apps.googleusercontent.com
GitHub OAuth: Ov23liZmK8EPzEymFddm
Gemini API: AIzaSyDWQ5k3_eBmA1qLWhnfSV3OtAglgD7X_f0
JWT Secret: your-super-secret-jwt-key-change-this-in-production
```

**Actions**:
1. Go to [Google Cloud Console](https://console.cloud.google.com/apis/credentials) → Delete old OAuth client → Create new
2. Go to [GitHub Settings > OAuth Apps](https://github.com/settings/developers) → Revoke client secret → Generate new
3. Go to [Google AI Studio](https://aistudio.google.com/app/apikey) → Delete old key → Create new
4. Generate new JWT secret: `openssl rand -base64 64`
5. Update environment variables in deployment environment

---

### 2. 🗑️ Remove .env Files from Git History (15 minutes)

```bash
# Install git-filter-repo
brew install git-filter-repo  # macOS
# or: pip3 install git-filter-repo

# Remove all .env files from history
git filter-repo --path .env --invert-paths
git filter-repo --path backend/.env --invert-paths
git filter-repo --path frontend/.env --invert-paths

# Force push (WARNING: This rewrites history)
git push origin --force --all
```

**Verify** `.env` is in `.gitignore` (already done ✅)

---

### 3. 🛡️ Implement Admin Authorization Middleware (1 hour)

**File**: `backend/internal/middleware/auth.go` (create new file)

```go
package middleware

import (
    "context"
    "net/http"
)

// AdminOnly restricts access to admin users only
func AdminOnly(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        userRole, ok := r.Context().Value("user_role").(string)
        if !ok || userRole != "admin" {
            http.Error(w, `{"error": "Forbidden - admin access required"}`, http.StatusForbidden)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

**Add role to JWT claims** in `backend/internal/auth/auth.go:32`:
```go
claims := jwt.MapClaims{
    "user_id": userID,
    "email":   email,
    "role":    role,  // Add this
    "exp":     jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
}
```

**Apply middleware** in `backend/internal/server/routes.go:144`:
```go
r.Group(func(r chi.Router) {
    r.Use(AdminOnly)  // Replace TODO comment
    r.Post("/admin/settings/collector", adminHandler.UpdateCollectorFrequency)
})
```

**Add role column** to users table:
```sql
ALTER TABLE users ADD COLUMN role VARCHAR(20) DEFAULT 'user';
UPDATE users SET role = 'admin' WHERE email = 'your-admin@email.com';
```

---

### 4. 🔧 Fix Unbounded Goroutine (45 minutes)

**File**: `backend/internal/server/routes.go:63-93`

**Current (BROKEN)**:
```go
go func() {
    ctx := context.Background()  // Never cancelled!
    for {
        time.Sleep(interval)
        collectorService.CollectAll(ctx)
    }
}()
```

**Fixed**:
```go
// Add to Server struct (server.go:20)
type Server struct {
    // ... existing fields
    collectorCtx    context.Context
    collectorCancel context.CancelFunc
}

// In routes.go, replace goroutine:
s.collectorCtx, s.collectorCancel = context.WithCancel(context.Background())

go func() {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            log.Println("Running scheduled content collection...")
            collectorService.CollectAll(s.collectorCtx)
            collectorService.ProcessUnprocessedItems(s.collectorCtx)
        case <-s.collectorCtx.Done():
            log.Println("Content collector shutting down...")
            return
        }
    }
}()

// Add shutdown handler in main.go:
sigCh := make(chan os.Signal, 1)
signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
go func() {
    <-sigCh
    log.Println("Shutting down gracefully...")
    server.collectorCancel()
    time.Sleep(2 * time.Second)  // Allow cleanup
    os.Exit(0)
}()
```

---

### 5. 🔑 Fix localStorage Token Inconsistency (30 minutes)

**Issue**: Two different keys used: `'token'` and `'auth_token'`

**Files to Fix**:
- `frontend/src/components/GoogleSignIn.vue:66` - Uses `'auth_token'`
- All other files use `'token'`

**Solution**: Standardize on `'token'`

**Fix GoogleSignIn.vue:66-67**:
```typescript
// BEFORE
localStorage.setItem('auth_token', response.token)
userStore.setToken(response.token)

// AFTER
localStorage.setItem('token', response.token)
userStore.setToken(response.token)
```

**Add cleanup migration** in `frontend/src/main.ts`:
```typescript
// Remove old auth_token if exists
if (localStorage.getItem('auth_token') && !localStorage.getItem('token')) {
    localStorage.setItem('token', localStorage.getItem('auth_token')!)
    localStorage.removeItem('auth_token')
}
```

---

## CRITICAL SECURITY FIXES (THIS WEEK)

### 6. Remove Client Secrets from Frontend (15 minutes)

**File**: `frontend/.env:4-5`

```env
# DELETE THESE LINES - NEVER put secrets in frontend!
VITE_GOOGLE_CLIENT_SECRET=...
VITE_GITHUB_CLIENT_SECRET=...
```

**Keep only**:
```env
VITE_API_URL=http://localhost:3000
VITE_GOOGLE_CLIENT_ID=...
VITE_GITHUB_CLIENT_ID=...
```

Secrets should ONLY be in backend.

---

### 7. Add JWT Validation to Frontend Auth Guard (20 minutes)

**File**: `frontend/src/router/index.ts:95-106`

```typescript
import { jwtDecode } from 'jwt-decode'

router.beforeEach((to, from, next) => {
    const token = localStorage.getItem('token')
    const requiresAuth = to.matched.some(record => record.meta.requiresAuth)

    if (requiresAuth && !token) {
        next('/login')
        return
    }

    // NEW: Validate JWT expiry
    if (requiresAuth && token) {
        try {
            const decoded = jwtDecode<{ exp: number }>(token)
            if (decoded.exp * 1000 < Date.now()) {
                localStorage.removeItem('token')
                next('/login')
                return
            }
        } catch (error) {
            localStorage.removeItem('token')
            next('/login')
            return
        }
    }

    next()
})
```

---

### 8. Fix N+1 Query Bug (2 hours)

**File**: `backend/internal/database/paths.go:76-104`

**Current (N+1 queries)**:
```go
for i := range p.Modules {
    lessonQuery := `SELECT ... FROM lessons WHERE module_id = $1`
    lRows, err := s.db.Query(ctx, lessonQuery, p.Modules[i].ID)
```

**Fixed (Single query with JOIN)**:
```go
func (s *service) GetLearningPath(ctx context.Context, id string) (*models.LearningPath, error) {
    // First get the path
    var p models.LearningPath
    pathQuery := `SELECT id, title, description, difficulty, created_at FROM learning_paths WHERE id = $1`
    err := s.db.QueryRow(ctx, pathQuery, id).Scan(&p.ID, &p.Title, &p.Description, &p.Difficulty, &p.CreatedAt)
    if err != nil {
        return nil, err
    }

    // Get all modules and lessons in ONE query with JOIN
    query := `
        SELECT
            m.id, m.title, m.description, m.order_index,
            l.id, l.title, l.content_type, l.content_url, l.content_body, l.duration_minutes, l.order_index
        FROM modules m
        LEFT JOIN lessons l ON l.module_id = m.id
        WHERE m.path_id = $1
        ORDER BY m.order_index, l.order_index
    `

    rows, err := s.db.Query(ctx, query, id)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    // Group results by module
    moduleMap := make(map[int]*models.Module)
    moduleOrder := []int{}

    for rows.Next() {
        var moduleID int
        var moduleTitle, moduleDesc string
        var moduleOrder int
        var lessonID sql.NullInt32
        var lessonTitle, lessonType, lessonURL sql.NullString
        var lessonBody sql.NullString
        var lessonDuration, lessonOrder sql.NullInt32

        err := rows.Scan(
            &moduleID, &moduleTitle, &moduleDesc, &moduleOrder,
            &lessonID, &lessonTitle, &lessonType, &lessonURL,
            &lessonBody, &lessonDuration, &lessonOrder,
        )
        if err != nil {
            return nil, err
        }

        // Create module if not exists
        if _, exists := moduleMap[moduleID]; !exists {
            moduleMap[moduleID] = &models.Module{
                ID:          moduleID,
                PathID:      p.ID,
                Title:       moduleTitle,
                Description: moduleDesc,
                OrderIndex:  moduleOrder,
                Lessons:     []models.Lesson{},
            }
            moduleOrder = append(moduleOrder, moduleID)
        }

        // Add lesson if exists
        if lessonID.Valid {
            lesson := models.Lesson{
                ID:              int(lessonID.Int32),
                ModuleID:        moduleID,
                Title:           lessonTitle.String,
                ContentType:     lessonType.String,
                ContentURL:      lessonURL.String,
                ContentBody:     lessonBody.String,
                DurationMinutes: int(lessonDuration.Int32),
                OrderIndex:      int(lessonOrder.Int32),
            }
            moduleMap[moduleID].Lessons = append(moduleMap[moduleID].Lessons, lesson)
        }
    }

    // Convert map to slice maintaining order
    for _, id := range moduleOrder {
        p.Modules = append(p.Modules, *moduleMap[id])
    }

    return &p, nil
}
```

---

### 9. Add Missing Database Indexes (5 minutes)

**File**: `backend/internal/database/schema.sql` (add to end of file)

```sql
-- Performance indexes (add these after existing CREATE INDEX statements)
CREATE INDEX IF NOT EXISTS idx_modules_path_id ON modules(path_id);
CREATE INDEX IF NOT EXISTS idx_lessons_module_id ON lessons(module_id);
CREATE INDEX IF NOT EXISTS idx_user_progress_user_completed ON user_progress(user_id, completed);
CREATE INDEX IF NOT EXISTS idx_content_items_processed ON content_items(processed_by_ai) WHERE processed_by_ai = FALSE;
CREATE INDEX IF NOT EXISTS idx_content_items_source_id ON content_items(source_id);
CREATE INDEX IF NOT EXISTS idx_project_submissions_project_id ON project_submissions(project_id);
CREATE INDEX IF NOT EXISTS idx_project_submissions_status ON project_submissions(status, submitted_at DESC);
CREATE INDEX IF NOT EXISTS idx_github_integrations_github_user_id ON github_integrations(github_user_id);
```

**Apply to database**:
```bash
psql -U user -d shadownova -f backend/internal/database/schema.sql
```

---

### 10. Add Missing github_username Column (2 minutes)

**File**: `backend/internal/database/schema.sql:9` (add after `password_hash`)

```sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100),
    password_hash VARCHAR(255),
    github_username VARCHAR(100),  -- ADD THIS LINE
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Apply migration**:
```sql
ALTER TABLE users ADD COLUMN IF NOT EXISTS github_username VARCHAR(100);
```

---

## VERIFICATION CHECKLIST

After completing all urgent fixes, verify:

- [ ] All secrets rotated and new secrets set in environment
- [ ] Git history cleaned of .env files
- [ ] New git push does not contain any secrets
- [ ] Admin endpoints return 403 for non-admin users
- [ ] Background collector shuts down gracefully on SIGTERM
- [ ] Token key consistent across all frontend files
- [ ] JWT expiry validated in router guard
- [ ] GetLearningPath makes only 2 queries (path + modules/lessons)
- [ ] Database indexes created successfully
- [ ] github_username column exists in users table

---

## TEST BEFORE DEPLOYMENT

```bash
# 1. Test admin authorization
curl -H "Authorization: Bearer <user-token>" http://localhost:3000/api/admin/settings/collector
# Should return 403 Forbidden

curl -H "Authorization: Bearer <admin-token>" http://localhost:3000/api/admin/settings/collector
# Should return 200 OK

# 2. Test graceful shutdown
docker-compose up -d
# Press Ctrl+C
# Check logs for "Content collector shutting down..."

# 3. Test N+1 fix
# Enable query logging in PostgreSQL
# Load a learning path
# Verify only 2 queries executed

# 4. Test JWT expiry
# Set token with past expiry in localStorage
# Navigate to protected route
# Should redirect to /login
```

---

## DEPLOYMENT READINESS

**After completing these 10 fixes**:
- ✅ Secrets secured
- ✅ Admin endpoints protected
- ✅ Goroutine managed properly
- ✅ Performance issue resolved
- ✅ Frontend auth working correctly

**Status**: ⚠️ CAUTION - Ready for staging deployment, NOT production

**Still TODO for production**:
- Migrate JWT to HttpOnly cookies (HIGH - prevents XSS)
- Add IDOR resource ownership validation (HIGH - prevents data leakage)
- Encrypt GitHub tokens at rest (HIGH)
- Implement CSRF protection (MEDIUM)
- Add structured logging (HIGH - for debugging)

See `ARCHITECTURAL_AUDIT.md` for complete roadmap.

---

## ESTIMATED TIME BREAKDOWN

| Task | Time | Difficulty |
|------|------|------------|
| 1. Rotate secrets | 30 min | Easy |
| 2. Clean git history | 15 min | Easy |
| 3. Admin middleware | 1 hour | Medium |
| 4. Fix goroutine | 45 min | Medium |
| 5. Token consistency | 30 min | Easy |
| 6. Remove client secrets | 15 min | Easy |
| 7. JWT validation | 20 min | Easy |
| 8. Fix N+1 query | 2 hours | Hard |
| 9. Add indexes | 5 min | Easy |
| 10. Add column | 2 min | Easy |

**Total**: ~5.5 hours with focused work

---

**Questions?** Review `ARCHITECTURAL_AUDIT.md` for detailed explanations.
