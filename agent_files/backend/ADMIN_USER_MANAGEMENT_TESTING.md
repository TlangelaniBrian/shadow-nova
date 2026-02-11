# Admin User Management Testing Guide

This guide provides instructions for testing the complete admin user management system.

## Prerequisites

1. Backend server running on `http://localhost:8080`
2. Frontend server running on `http://localhost:5173`
3. PostgreSQL database with schema applied
4. At least one admin user created in the database

## Create Initial Admin User

If you don't have an admin user yet, create one directly in the database:

```sql
-- Create an admin user
INSERT INTO users (email, username, password_hash, user_role)
VALUES (
    'admin@shadownova.com',
    'admin',
    '$2a$10$...',  -- Use bcrypt hash of your password
    'admin'
);
```

Or use the registration endpoint and manually update the role:

```sql
UPDATE users SET user_role = 'admin' WHERE email = 'your-email@example.com';
```

## Backend API Testing

### 1. List Users (GET /api/v1/admin/users)

```bash
# Get first page of users
curl -X GET "http://localhost:8080/api/v1/admin/users?page=1&limit=20" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -b "auth_token=YOUR_AUTH_TOKEN"

# Expected response:
{
  "data": [
    {
      "id": 1,
      "email": "user@example.com",
      "username": "testuser",
      "role": "user",
      "github_username": null,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "page": 1,
  "limit": 20,
  "total": 1,
  "total_pages": 1,
  "has_next": false,
  "has_prev": false
}
```

### 2. Get Specific User (GET /api/v1/admin/users/{id})

```bash
curl -X GET "http://localhost:8080/api/v1/admin/users/1" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -b "auth_token=YOUR_AUTH_TOKEN"

# Expected response:
{
  "id": 1,
  "email": "user@example.com",
  "username": "testuser",
  "role": "user",
  "github_username": null,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### 3. Create User (POST /api/v1/admin/users)

```bash
curl -X POST "http://localhost:8080/api/v1/admin/users" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "X-CSRF-Token: YOUR_CSRF_TOKEN" \
  -b "auth_token=YOUR_AUTH_TOKEN" \
  -d '{
    "email": "newuser@example.com",
    "username": "newuser",
    "password": "securepassword123",
    "role": "user"
  }'

# Expected response: 201 Created
{
  "id": 2,
  "email": "newuser@example.com",
  "username": "newuser",
  "role": "user",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### 4. Update User (PUT /api/v1/admin/users/{id})

```bash
curl -X PUT "http://localhost:8080/api/v1/admin/users/2" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "X-CSRF-Token: YOUR_CSRF_TOKEN" \
  -b "auth_token=YOUR_AUTH_TOKEN" \
  -d '{
    "email": "updated@example.com",
    "role": "admin"
  }'

# Expected response: 200 OK
{
  "id": 2,
  "email": "updated@example.com",
  "username": "newuser",
  "role": "admin",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T01:00:00Z"
}
```

### 5. Delete User (DELETE /api/v1/admin/users/{id})

```bash
curl -X DELETE "http://localhost:8080/api/v1/admin/users/2" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "X-CSRF-Token: YOUR_CSRF_TOKEN" \
  -b "auth_token=YOUR_AUTH_TOKEN"

# Expected response: 204 No Content
```

## Frontend Testing

### Access Admin Panel

1. Log in as an admin user
2. Navigate to the sidebar
3. You should see an "Admin" section with "User Management" link
4. Click on "User Management"

### Test User List View

1. Verify all users are displayed in a table
2. Check pagination controls appear if more than 20 users exist
3. Verify user roles are displayed with colored badges (purple for admin, gray for user)
4. Check GitHub username column shows "-" for users without GitHub integration

### Test User Creation

1. Click "Create User" button
2. Fill in the form:
   - Email: test@example.com
   - Username: testuser
   - Password: testpassword123
   - Role: user
3. Click "Save"
4. Verify success toast appears
5. Verify new user appears in the list

### Test User Editing

1. Click "Edit" button on any user row
2. Modify the email or role
3. Click "Save"
4. Verify success toast appears
5. Verify changes are reflected in the list
6. Try editing your own account - should see error message

### Test User Deletion

1. Click "Delete" button on any user row
2. Confirm the deletion dialog
3. Verify success toast appears
4. Verify user is removed from the list
5. Try deleting your own account - button should be disabled

### Test Validation

1. Try creating a user with invalid email - should see validation error
2. Try creating a user with password less than 8 characters - should see error
3. Try creating a user with existing email - should see database error
4. Try updating with empty fields - form should handle gracefully

### Test Permissions

1. Log out as admin
2. Log in as a regular user
3. Try to access `/admin/users` directly
4. Should be redirected to `/dashboard`
5. Verify "Admin" section doesn't appear in sidebar

## Database Verification

### Check Soft Deletes

```sql
-- View all users including soft-deleted ones
SELECT id, email, username, user_role, deleted_at
FROM users;

-- View only active users (what the app shows)
SELECT id, email, username, user_role
FROM users
WHERE deleted_at IS NULL;

-- View only deleted users
SELECT id, email, username, user_role, deleted_at
FROM users
WHERE deleted_at IS NOT NULL;
```

### Verify Role Updates

```sql
-- Check user roles
SELECT id, username, user_role FROM users;
```

## Security Testing

### Test Admin-Only Access

1. Make API request to admin endpoints without admin role:

```bash
# As regular user (should fail with 403 Forbidden)
curl -X GET "http://localhost:8080/api/v1/admin/users" \
  -H "Authorization: Bearer REGULAR_USER_TOKEN" \
  -b "auth_token=REGULAR_USER_TOKEN"

# Expected response: 403 Forbidden
{
  "error": "Access denied: admin role required"
}
```

### Test Self-Modification Protection

1. Try to update your own account via admin endpoint (should fail)
2. Try to delete your own account (should fail)

### Test Password Handling

1. Verify password hash is never returned in API responses
2. Check that `password_hash` field is always empty in JSON responses
3. Confirm passwords are bcrypt hashed in database

## Performance Testing

### Test Pagination

1. Create 100+ test users
2. Navigate through pages
3. Verify pagination controls work correctly
4. Check page load times are acceptable

### Test Large User List

```bash
# Create multiple users for testing
for i in {1..100}; do
  curl -X POST "http://localhost:8080/api/v1/admin/users" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer YOUR_JWT_TOKEN" \
    -d "{
      \"email\": \"user${i}@example.com\",
      \"username\": \"user${i}\",
      \"password\": \"password123\",
      \"role\": \"user\"
    }"
done
```

## Expected Results

All tests should:
- Return appropriate HTTP status codes
- Display clear error messages for validation failures
- Maintain data consistency in the database
- Properly enforce admin-only access
- Handle edge cases gracefully
- Show loading states during API calls
- Display success/error toasts appropriately

## Common Issues

### Issue: 403 Forbidden on admin endpoints
**Solution**: Verify the user has `user_role = 'admin'` in the database

### Issue: CSRF token errors
**Solution**: Ensure CSRF token is fetched on app load and included in state-changing requests

### Issue: Users list is empty
**Solution**: Check that `deleted_at IS NULL` filter is applied in the query

### Issue: Cannot delete users
**Solution**: Verify the `deleted_at` column exists in the users table

### Issue: Password validation fails
**Solution**: Ensure password is at least 8 characters and meets requirements
