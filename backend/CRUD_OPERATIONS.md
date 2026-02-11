# Shadow Nova - CRUD Operations Documentation

## Overview

This document describes the complete CRUD (Create, Read, Update, Delete) operations available for Shadow Nova resources. All operations follow RESTful conventions and implement proper authorization, soft deletes, and audit trails.

## Table of Contents

- [Authentication & Authorization](#authentication--authorization)
- [Learning Paths](#learning-paths)
- [Modules](#modules)
- [Lessons](#lessons)
- [Projects](#projects)
- [Soft Delete Strategy](#soft-delete-strategy)
- [Cascading Behavior](#cascading-behavior)
- [Error Handling](#error-handling)

---

## Authentication & Authorization

### Authorization Matrix

| Resource | Create | Read | Update | Delete | Notes |
|----------|--------|------|--------|--------|-------|
| Learning Paths | Admin | All | Admin | Admin | Public read access |
| Modules | Admin | All | Admin | Admin | Public read access via parent path |
| Lessons | Admin | All | Admin | Admin | Public read access via parent module |
| Projects | Admin | All | Admin | Admin | Public read access |
| Submissions | User (own) | User (own) | Admin | - | Users can only see their own |

### Authorization Requirements

- **Public Routes**: No authentication required for reading resources
- **Protected Routes**: JWT token required for user-specific operations
- **Admin Routes**: JWT token with `user_role = 'admin'` required for CUD operations

---

## Learning Paths

### Create Learning Path

**Endpoint**: `POST /api/paths`
**Auth**: Admin only
**Request**:
```json
{
  "id": "go-fundamentals",
  "title": "Go Fundamentals",
  "description": "Learn Go programming from scratch",
  "difficulty": "beginner"
}
```

**Response**: `201 Created`
```json
{
  "success": true,
  "message": "Learning path created successfully",
  "data": {
    "id": "go-fundamentals",
    "title": "Go Fundamentals",
    "description": "Learn Go programming from scratch",
    "difficulty": "beginner",
    "created_at": "2026-02-12T10:30:00Z"
  }
}
```

### Read Learning Paths

**Endpoint**: `GET /api/paths`
**Auth**: None (public)
**Query Parameters**:
- `page` (default: 1)
- `limit` (default: 10, max: 100)

**Response**: `200 OK`
```json
{
  "data": [...],
  "page": 1,
  "limit": 10,
  "total": 25
}
```

### Read Single Learning Path

**Endpoint**: `GET /api/paths/{id}`
**Auth**: None (public)
**Response**: `200 OK` - Returns path with nested modules and lessons

### Update Learning Path

**Endpoint**: `PUT /api/paths/{id}` or `PATCH /api/paths/{id}`
**Auth**: Admin only
**Request**:
```json
{
  "title": "Advanced Go Programming",
  "description": "Updated description",
  "difficulty": "intermediate"
}
```

**Response**: `200 OK`
```json
{
  "success": true,
  "message": "Learning path updated successfully"
}
```

### Delete Learning Path

**Endpoint**: `DELETE /api/paths/{id}`
**Auth**: Admin only
**Response**: `204 No Content`

**Behavior**:
- Performs a **soft delete** (sets `deleted_at` timestamp)
- Cascades to all modules and lessons within the path
- Preserves data for audit trail and potential recovery

---

## Modules

### Create Module

**Endpoint**: `POST /api/paths/{id}/modules`
**Auth**: Admin only
**Request**:
```json
{
  "title": "Introduction to Go",
  "description": "Getting started with Go",
  "order_index": 1
}
```

### Update Module

**Endpoint**: `PUT /api/modules/{id}` or `PATCH /api/modules/{id}`
**Auth**: Admin only
**Request**:
```json
{
  "title": "Updated Module Title",
  "description": "Updated description",
  "order_index": 2
}
```

### Delete Module

**Endpoint**: `DELETE /api/modules/{id}`
**Auth**: Admin only
**Response**: `204 No Content`

**Behavior**:
- Performs a **soft delete**
- Cascades to all lessons within the module

---

## Lessons

### Create Lesson

**Endpoint**: `POST /api/lessons`
**Auth**: Admin only
**Request**:
```json
{
  "module_id": 123,
  "title": "Variables and Types",
  "content_type": "video",
  "content_url": "https://youtube.com/...",
  "content_body": "Markdown content here...",
  "duration_minutes": 15,
  "order_index": 1
}
```

**Validation**:
- `content_type` must be one of: `video`, `article`, `quiz`
- Either `content_url` or `content_body` should be provided

### Update Lesson

**Endpoint**: `PUT /api/lessons/{id}` or `PATCH /api/lessons/{id}`
**Auth**: Admin only
**Request**: Same as create, all fields optional

### Delete Lesson

**Endpoint**: `DELETE /api/lessons/{id}`
**Auth**: Admin only
**Response**: `204 No Content`

**Behavior**: Performs a **soft delete**

---

## Projects

### Create Project

**Endpoint**: `POST /api/projects`
**Auth**: Admin only
**Request**:
```json
{
  "id": "rest-api-project",
  "title": "Build a REST API",
  "description": "Create a production-ready REST API",
  "difficulty": "intermediate",
  "tech_stack": ["go", "postgresql", "docker"]
}
```

### Read Projects

**Endpoint**: `GET /api/projects`
**Auth**: None (public)
**Query Parameters**: `page`, `limit`

### Read Single Project

**Endpoint**: `GET /api/projects/{id}`
**Auth**: None (public)

### Update Project

**Endpoint**: `PUT /api/projects/{id}` or `PATCH /api/projects/{id}`
**Auth**: Admin only
**Request**:
```json
{
  "title": "Updated Project Title",
  "description": "Updated description",
  "difficulty": "advanced",
  "tech_stack": ["go", "postgresql", "docker", "kubernetes"]
}
```

### Delete Project

**Endpoint**: `DELETE /api/projects/{id}`
**Auth**: Admin only
**Response**: `204 No Content`

**Behavior**: Performs a **soft delete**

---

## Soft Delete Strategy

### Why Soft Deletes?

Shadow Nova uses soft deletes instead of hard deletes for the following reasons:

1. **Audit Trail**: Maintain a complete history of all content changes
2. **Recovery**: Allow restoration of accidentally deleted content
3. **Data Integrity**: Preserve relationships and references
4. **Analytics**: Track content lifecycle and usage patterns

### Implementation

- All resources have a `deleted_at` timestamp column (default: `NULL`)
- When deleted, `deleted_at` is set to `CURRENT_TIMESTAMP`
- All queries filter by `WHERE deleted_at IS NULL`
- Indexes on `deleted_at` ensure query performance

### Schema Changes

See migration: `backend/internal/database/migrations/005_add_soft_delete.sql`

Added columns:
- `deleted_at TIMESTAMP DEFAULT NULL`
- `updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP`

### Automatic Timestamp Updates

Triggers automatically update `updated_at` on any UPDATE operation:
```sql
CREATE TRIGGER update_learning_paths_updated_at
    BEFORE UPDATE ON learning_paths
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
```

---

## Cascading Behavior

### Soft Delete Cascades

When a parent resource is soft-deleted, the application automatically cascades the deletion:

```
Learning Path (soft deleted)
  └── Module 1 (soft deleted)
        ├── Lesson 1 (soft deleted)
        └── Lesson 2 (soft deleted)
  └── Module 2 (soft deleted)
        └── Lesson 3 (soft deleted)
```

**Implementation**: Transaction-based cascading in application layer
```go
func (s *service) DeleteLearningPath(ctx context.Context, id string) error {
    return s.WithTx(ctx, func(tx pgx.Tx) error {
        // Soft delete path
        // Soft delete modules
        // Soft delete lessons
        return nil
    })
}
```

### Hard Delete Cascades (Database Level)

Foreign key constraints ensure hard deletes cascade properly:
```sql
REFERENCES learning_paths(id) ON DELETE CASCADE
REFERENCES modules(id) ON DELETE CASCADE
```

If a hard delete is ever performed (e.g., GDPR compliance), the database automatically removes child records.

---

## Error Handling

### Standard Error Responses

**404 Not Found**:
```json
{
  "error": "Resource not found",
  "message": "learning path go-fundamentals not found"
}
```

**400 Bad Request**:
```json
{
  "error": "Validation failed",
  "fields": {
    "title": "title is required",
    "difficulty": "must be one of: beginner, intermediate, advanced"
  }
}
```

**401 Unauthorized**:
```json
{
  "error": "Unauthorized",
  "message": "No authentication token provided"
}
```

**403 Forbidden**:
```json
{
  "error": "Forbidden",
  "message": "Admin role required for this operation"
}
```

### IDOR Prevention

All ownership-sensitive operations validate ownership:
- Users can only access their own submissions
- Admin operations are role-gated
- Path/module/lesson access is public (read-only)

Middleware examples:
- `middleware.ValidateSubmissionOwnership`
- `middleware.ValidatePathOwnership`
- `middleware.AdminOnly`

---

## Future Enhancements

### Planned Features

1. **Bulk Operations**: Batch create/update/delete
2. **Versioning**: Track content version history
3. **Restore API**: Endpoint to restore soft-deleted resources
4. **Hard Delete Cleanup**: Scheduled job to permanently remove old soft-deleted records
5. **Ownership Transfer**: Admin ability to reassign content ownership
6. **Archive vs Delete**: Separate archival from deletion

### API Versioning

Future breaking changes will be versioned:
- Current: `/api/paths`
- Future: `/api/v2/paths`

---

## Testing CRUD Operations

### Using cURL

```bash
# Create a learning path (admin)
curl -X POST http://localhost:8080/api/paths \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "test-path",
    "title": "Test Path",
    "description": "Test description",
    "difficulty": "beginner"
  }'

# Update a learning path (admin)
curl -X PUT http://localhost:8080/api/paths/test-path \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Test Path",
    "description": "Updated description",
    "difficulty": "intermediate"
  }'

# Delete a learning path (admin)
curl -X DELETE http://localhost:8080/api/paths/test-path \
  -H "Authorization: Bearer $ADMIN_TOKEN"

# List all paths (public)
curl http://localhost:8080/api/paths?page=1&limit=10
```

### Using Frontend API Client

See: `frontend/src/api/paths.ts` for TypeScript client examples

---

## Database Schema Reference

### Key Tables

```sql
learning_paths (
    id VARCHAR(50) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    difficulty VARCHAR(50),
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
)

modules (
    id SERIAL PRIMARY KEY,
    path_id VARCHAR(50) REFERENCES learning_paths(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    order_index INTEGER,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
)

lessons (
    id SERIAL PRIMARY KEY,
    module_id INTEGER REFERENCES modules(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    content_type VARCHAR(50) NOT NULL,
    content_url VARCHAR(500),
    content_body TEXT,
    duration_minutes INTEGER,
    order_index INTEGER,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
)

projects (
    id VARCHAR(50) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    difficulty VARCHAR(50),
    tech_stack TEXT[],
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
)
```

---

## Support & Maintenance

For questions or issues related to CRUD operations:
1. Check this documentation
2. Review API handler code: `backend/internal/handlers/`
3. Review database layer: `backend/internal/database/`
4. Contact the development team

**Last Updated**: 2026-02-12
**Version**: 1.0.0
