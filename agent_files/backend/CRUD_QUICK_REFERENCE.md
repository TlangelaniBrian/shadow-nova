# CRUD Operations Quick Reference

## API Endpoints Summary

### Learning Paths

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/paths` | Public | List all paths (paginated) |
| GET | `/api/paths/{id}` | Public | Get single path with modules/lessons |
| POST | `/api/paths` | Admin | Create new path |
| PUT/PATCH | `/api/paths/{id}` | Admin | Update path |
| DELETE | `/api/paths/{id}` | Admin | Soft delete path |

### Modules

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/api/paths/{id}/modules` | Admin | Add module to path |
| PUT/PATCH | `/api/modules/{id}` | Admin | Update module |
| DELETE | `/api/modules/{id}` | Admin | Soft delete module |

### Lessons

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/api/lessons` | Admin | Create lesson |
| PUT/PATCH | `/api/lessons/{id}` | Admin | Update lesson |
| DELETE | `/api/lessons/{id}` | Admin | Soft delete lesson |

### Projects

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/projects` | Public | List all projects (paginated) |
| GET | `/api/projects/{id}` | Public | Get single project |
| POST | `/api/projects` | Admin | Create new project |
| PUT/PATCH | `/api/projects/{id}` | Admin | Update project |
| DELETE | `/api/projects/{id}` | Admin | Soft delete project |

## Database Methods

### Learning Paths
```go
// Create
CreateLearningPath(ctx, path *models.LearningPath) error

// Read
GetLearningPaths(ctx, limit, offset int) ([]models.LearningPath, error)
GetLearningPathsCount(ctx) (int, error)
GetLearningPath(ctx, id string) (*models.LearningPath, error)

// Update
UpdateLearningPath(ctx, id string, updates *models.LearningPath) error

// Delete (soft)
DeleteLearningPath(ctx, id string) error
```

### Modules
```go
CreateModule(ctx, module *models.Module) error
UpdateModule(ctx, id int, updates *models.Module) error
DeleteModule(ctx, id int) error
```

### Lessons
```go
CreateLesson(ctx, lesson *models.Lesson) error
UpdateLesson(ctx, id int, updates *models.Lesson) error
DeleteLesson(ctx, id int) error
GetLesson(ctx, lessonID int) (*models.Lesson, error)
```

### Projects
```go
// Create
CreateProject(ctx, project *models.Project) error

// Read
GetProjects(ctx, limit, offset int) ([]models.Project, error)
GetProjectsCount(ctx) (int, error)
GetProject(ctx, id string) (*models.Project, error)

// Update
UpdateProject(ctx, id string, updates *models.Project) error

// Delete (soft)
DeleteProject(ctx, id string) error
```

## Soft Delete Queries

All read queries filter soft-deleted records:
```sql
SELECT * FROM learning_paths WHERE deleted_at IS NULL
SELECT * FROM modules WHERE deleted_at IS NULL
SELECT * FROM lessons WHERE deleted_at IS NULL
SELECT * FROM projects WHERE deleted_at IS NULL
```

Soft delete sets timestamp:
```sql
UPDATE learning_paths
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
```

## Migration

Apply migration to add soft delete support:
```bash
psql $DATABASE_URL -f backend/internal/database/migrations/005_add_soft_delete.sql
```

Columns added:
- `deleted_at TIMESTAMP DEFAULT NULL`
- `updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP`

Triggers created:
- Auto-update `updated_at` on any UPDATE

## Testing Examples

### Create Learning Path
```bash
curl -X POST http://localhost:8080/api/paths \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "go-advanced",
    "title": "Advanced Go",
    "description": "Advanced Go concepts",
    "difficulty": "advanced"
  }'
```

### Update Learning Path
```bash
curl -X PUT http://localhost:8080/api/paths/go-advanced \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Advanced Go Programming",
    "description": "Updated description",
    "difficulty": "advanced"
  }'
```

### Delete Learning Path
```bash
curl -X DELETE http://localhost:8080/api/paths/go-advanced \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

### List Learning Paths (Public)
```bash
curl http://localhost:8080/api/paths?page=1&limit=10
```

### Get Single Project (Public)
```bash
curl http://localhost:8080/api/projects/rest-api-project
```

## Response Examples

### Success (200 OK)
```json
{
  "success": true,
  "message": "Learning path updated successfully",
  "data": null
}
```

### Created (201 Created)
```json
{
  "success": true,
  "message": "Learning path created successfully",
  "data": {
    "id": "go-advanced",
    "title": "Advanced Go",
    "created_at": "2026-02-12T10:30:00Z"
  }
}
```

### No Content (204)
- Empty response body
- Used for DELETE operations

### Not Found (404)
```json
{
  "error": "Resource not found",
  "message": "learning path go-advanced not found"
}
```

### Forbidden (403)
```json
{
  "error": "Forbidden",
  "message": "Admin role required for this operation"
}
```

## Files Reference

### Database Layer
- `/backend/internal/database/paths.go` - Learning paths, modules, lessons CRUD
- `/backend/internal/database/projects_crud.go` - Projects CRUD operations
- `/backend/internal/database/database.go` - Service interface definitions
- `/backend/internal/database/mock.go` - Mock implementations for testing

### Handlers
- `/backend/internal/handlers/paths.go` - Learning path HTTP handlers
- `/backend/internal/handlers/projects.go` - Project HTTP handlers

### Routes
- `/backend/internal/server/routes.go` - API route definitions

### Migrations
- `/backend/internal/database/migrations/005_add_soft_delete.sql` - Soft delete migration

### Documentation
- `/backend/CRUD_OPERATIONS.md` - Comprehensive guide
- `/backend/CRUD_IMPLEMENTATION_SUMMARY.md` - Implementation details
- `/backend/CRUD_QUICK_REFERENCE.md` - This file

## Authorization Headers

Get admin token:
```bash
# Login as admin
RESPONSE=$(curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"password"}')

# Extract token
export ADMIN_TOKEN=$(echo $RESPONSE | jq -r '.data.access_token')

# Use token
curl -H "Authorization: Bearer $ADMIN_TOKEN" http://localhost:8080/api/paths
```

## Common Operations

### Cascading Delete
When you delete a learning path, it cascades:
```
DELETE /api/paths/go-advanced
  ├─ Soft deletes the path
  ├─ Soft deletes all modules in the path
  └─ Soft deletes all lessons in those modules
```

### Pagination
All list endpoints support pagination:
- `?page=1` (default: 1)
- `?limit=10` (default: 10, max: 100)

Response includes:
```json
{
  "data": [...],
  "page": 1,
  "limit": 10,
  "total": 25
}
```

## Status Codes

| Code | Meaning | Usage |
|------|---------|-------|
| 200 | OK | Successful GET, PUT, PATCH |
| 201 | Created | Successful POST |
| 204 | No Content | Successful DELETE |
| 400 | Bad Request | Validation error |
| 401 | Unauthorized | Missing/invalid token |
| 403 | Forbidden | Insufficient permissions |
| 404 | Not Found | Resource not found |
| 500 | Server Error | Internal error |

## Next Steps

1. Apply migration: `005_add_soft_delete.sql`
2. Test endpoints with admin token
3. Update frontend API clients
4. Implement restore functionality (optional)
5. Add audit logging (optional)

For detailed information, see `/backend/CRUD_OPERATIONS.md`
