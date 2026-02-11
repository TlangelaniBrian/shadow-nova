# Ownership Middleware Testing Guide

This document describes how to test the IDOR (Insecure Direct Object Reference) prevention middleware.

## Overview

The ownership middleware prevents unauthorized access to resources by validating that users can only access resources they own or have permission to access.

## Middleware Functions

### 1. ValidatePathOwnership
**Purpose**: Validates user access to learning paths
**Current Behavior**: All authenticated users can access all paths
**Future**: Will check enrollment/purchase status

**Protected Endpoints**:
- `GET /api/v1/paths/{id}/progress` - Get progress for a specific path

### 2. ValidateSubmissionOwnership
**Purpose**: Validates ownership of project submissions
**Behavior**: Users can only access their own submissions

**Protected Endpoints**:
- `GET /api/v1/submissions/{id}` - Get a specific submission
- `PATCH /api/v1/submissions/{id}` - Update a specific submission

### 3. ValidateProgressOwnership
**Purpose**: Validates ownership of progress records
**Behavior**: Users can only access their own progress records

**Protected Endpoints**:
- Currently no endpoints using this middleware (ready for future use)

## Test Scenarios

### Test 1: Owner Can Access Their Submission
```bash
# Login as user1
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user1@example.com","password":"password"}'

# Get token from response and use it
TOKEN="<your_token>"

# Submit a project
curl -X POST http://localhost:8080/api/v1/submissions \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "project_id":"web-dev",
    "github_repo_url":"https://github.com/user1/project",
    "pr_url":"https://github.com/user1/project/pull/1",
    "demo_url":"https://demo.com"
  }'

# Note the submission ID from response, e.g., id=1

# Access your own submission (should succeed - 200 OK)
curl http://localhost:8080/api/v1/submissions/1 \
  -H "Authorization: Bearer $TOKEN"
```

### Test 2: Non-Owner Gets 403 Forbidden
```bash
# Login as user2
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user2@example.com","password":"password"}'

TOKEN2="<user2_token>"

# Try to access user1's submission (should fail - 403 Forbidden)
curl http://localhost:8080/api/v1/submissions/1 \
  -H "Authorization: Bearer $TOKEN2"

# Expected response:
# {"error":"You do not have access to this submission","status":403}
```

### Test 3: Missing Authentication
```bash
# Try to access submission without token (should fail - 401 Unauthorized)
curl http://localhost:8080/api/v1/submissions/1

# Expected response:
# {"error":"User not authenticated","status":401}
```

### Test 4: Invalid Submission ID
```bash
# Try with invalid ID format (should fail - 400 Bad Request)
curl http://localhost:8080/api/v1/submissions/invalid \
  -H "Authorization: Bearer $TOKEN"

# Try with non-existent ID (should fail - 404 Not Found or 500)
curl http://localhost:8080/api/v1/submissions/99999 \
  -H "Authorization: Bearer $TOKEN"
```

### Test 5: Path Access (Currently Allows All)
```bash
# Any authenticated user can access path progress
curl http://localhost:8080/api/paths/web-dev/progress \
  -H "Authorization: Bearer $TOKEN"

# Should succeed - 200 OK
```

## Running Unit Tests

```bash
cd backend
go test -v ./internal/middleware/ownership_test.go ./internal/middleware/ownership.go ./internal/middleware/auth.go
```

## Expected Test Results

All tests should pass:
- `TestValidatePathOwnership` - Tests path access validation
- `TestValidatePathOwnership_MissingUserID` - Tests missing authentication
- `TestValidatePathOwnership_MissingPathID` - Tests missing path ID
- `TestValidateSubmissionOwnership` - Tests submission ownership validation
- `TestValidateSubmissionOwnership_InvalidID` - Tests invalid submission ID
- `TestValidateProgressOwnership` - Tests progress ownership validation
- `TestValidateProgressOwnership_InvalidID` - Tests invalid progress ID

## Database Methods

### UserHasAccessToPath
```go
UserHasAccessToPath(ctx context.Context, userID int, pathID string) (bool, error)
```
- Currently returns `true` for all authenticated users
- Verifies the path exists in the database
- Future: Check enrollment table

### UserOwnsSubmission
```go
UserOwnsSubmission(ctx context.Context, userID int, submissionID int) (bool, error)
```
- Checks if the user_id matches the submission's user_id
- Returns false if submission doesn't exist or IDs don't match

### UserOwnsProgress
```go
UserOwnsProgress(ctx context.Context, userID int, progressID int) (bool, error)
```
- Checks if the user_id matches the progress record's user_id
- Returns false if progress record doesn't exist or IDs don't match

## Security Considerations

1. **Always use these middleware on sensitive endpoints** that deal with user-specific data
2. **Place middleware after authentication** - ownership checks require a valid userID
3. **Extract IDs from URL params**, not from request body (to prevent manipulation)
4. **Return 403 Forbidden**, not 404 Not Found, to avoid information disclosure
5. **Log failed ownership checks** for security monitoring (future enhancement)

## Future Enhancements

1. Add enrollment system for learning paths
2. Add role-based access (admins can access all resources)
3. Add audit logging for failed access attempts
4. Add rate limiting for repeated failed access attempts
5. Add ownership middleware for other resources (comments, ratings, etc.)
