# User Profile Management Implementation

## Overview

This document describes the complete implementation of user profile management endpoints for Shadow Nova, allowing users to view and update their profile information and change their passwords.

## Implementation Summary

### Backend Components

#### 1. User Handler (`internal/handlers/user.go`)

Created a new handler with three main endpoints:

- **GetProfile** - Retrieves the authenticated user's profile
- **UpdateProfile** - Updates username and/or email
- **UpdatePassword** - Changes the user's password with verification

Key features:
- Password hash is never sent to the client
- Current password verification before allowing password changes
- Proper error handling with centralized error responses
- Request validation using the validator package

#### 2. Database Methods (`internal/database/users.go`)

Added three new database methods:

```go
GetUserByID(ctx context.Context, userID int) (*models.User, error)
UpdateUser(ctx context.Context, userID int, user *models.User) error
UpdateUserPassword(ctx context.Context, userID int, hashedPassword string) error
```

Features:
- Uses pgx for efficient database operations
- Proper error handling with custom error types
- Returns NotFound error when user doesn't exist
- Updates timestamp automatically on changes

#### 3. Routes (`internal/server/routes.go`)

Added protected routes under authentication middleware:

```go
// User profile management
r.Get("/user/profile", userHandler.GetProfile)
r.Patch("/user/profile", userHandler.UpdateProfile)
r.Put("/user/password", userHandler.UpdatePassword)
```

All routes require:
- JWT authentication via auth middleware
- CSRF protection
- Idempotency support

#### 4. Mock Service (`internal/database/mock.go`)

Added stub implementations for testing:

```go
GetUserByIDFunc
UpdateUserFunc
UpdateUserPasswordFunc
```

### Frontend Components

#### 1. API Client (`frontend/src/api/user.ts`)

Already existed with the correct interface:

```typescript
export const userApi = {
    getUserProfile(): Promise<UserProfile>
    updateUserProfile(data: UpdateProfileData): Promise<UserProfile>
    updatePassword(data: UpdatePasswordData): Promise<void>
}
```

#### 2. Profile View (`frontend/src/views/ProfileView.vue`)

Updated with complete functionality:

**Features:**
- Fetches real user data on mount from `/api/v1/user/profile`
- Profile edit form with username and email fields
- Password change form with current and new password
- Real-time validation
- Loading states for all async operations
- Toast notifications for success/error feedback
- Updates local storage when profile changes
- Integration with existing GitHub connection features

**User Experience:**
- Clean, card-based UI matching existing design
- Disabled buttons during operations
- Clear error messages
- Form validation before submission
- Prevents unnecessary API calls (no-op if no changes)

## API Endpoints

### GET /api/v1/user/profile

Retrieves the authenticated user's profile.

**Authentication:** Required (JWT)

**Response:**
```json
{
  "id": 1,
  "email": "user@example.com",
  "username": "johndoe",
  "role": "user",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

**Status Codes:**
- 200: Success
- 401: Unauthorized
- 404: User not found
- 500: Internal server error

### PATCH /api/v1/user/profile

Updates the user's profile information.

**Authentication:** Required (JWT)

**Request Body:**
```json
{
  "username": "newusername",  // optional, min 3 chars, max 100
  "email": "newemail@example.com"  // optional, must be valid email
}
```

**Validation:**
- Username: 3-100 characters if provided
- Email: Valid email format if provided
- At least one field should be provided

**Response:**
```json
{
  "id": 1,
  "email": "newemail@example.com",
  "username": "newusername",
  "role": "user",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

**Status Codes:**
- 200: Success
- 400: Validation error
- 401: Unauthorized
- 404: User not found
- 500: Internal server error

### PUT /api/v1/user/password

Changes the user's password.

**Authentication:** Required (JWT)

**Request Body:**
```json
{
  "current_password": "OldPassword123!",
  "new_password": "NewPassword123!"
}
```

**Validation:**
- current_password: Required
- new_password: Required, minimum 8 characters

**Response:**
```json
{
  "message": "Password updated successfully"
}
```

**Status Codes:**
- 200: Success
- 400: Validation error
- 401: Unauthorized or incorrect current password
- 404: User not found
- 500: Internal server error

## Security Considerations

### Password Security
- Passwords are hashed using bcrypt with default cost (10)
- Current password must be verified before allowing changes
- Password hash is never sent in API responses
- Minimum password length of 8 characters enforced

### Authentication & Authorization
- All endpoints require valid JWT authentication
- User can only modify their own profile (userID from JWT)
- CSRF protection applied to state-changing requests
- Idempotency support for profile updates

### Input Validation
- Server-side validation using go-playground/validator
- Email format validation
- Username length constraints (3-100 characters)
- Clear error messages without leaking sensitive information

## Testing

### Manual Testing

Use the provided test script:

```bash
# Set your auth token
export AUTH_TOKEN="your_jwt_token_here"

# Run the test script
./backend/scripts/test_user_profile.sh
```

### Testing Checklist

- [ ] User can view their profile
- [ ] User can update username
- [ ] User can update email
- [ ] User can update both username and email
- [ ] Profile update validates input (min/max length, email format)
- [ ] Password change requires correct current password
- [ ] Password change validates new password length
- [ ] Incorrect current password returns 401
- [ ] Unauthorized requests return 401
- [ ] Password hash is never returned to client
- [ ] Updated timestamp changes after updates
- [ ] Frontend shows loading states during operations
- [ ] Frontend displays success/error toasts
- [ ] Frontend updates local user state after changes

## Integration Points

### Existing Systems
- **Authentication**: Uses existing JWT middleware
- **CSRF**: Protected by existing CSRF middleware
- **Idempotency**: Supports idempotent requests
- **Error Handling**: Uses centralized error handling
- **Logging**: Structured logging throughout
- **Validation**: Consistent validation approach

### Frontend Integration
- **Auth Store**: Updates user data in localStorage
- **Toast System**: Shows notifications
- **API Client**: Uses existing axios client with interceptors
- **Components**: Integrates with ProfileCard component

## Database Schema

Uses existing `users` table:

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100),
    password_hash VARCHAR(255) NOT NULL,
    user_role VARCHAR(50) DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

No schema changes required.

## Future Enhancements

### Potential Improvements
1. **Email Verification**: Require email verification when changing email
2. **Password Strength Meter**: Add client-side password strength indicator
3. **Profile Pictures**: Add support for avatar uploads
4. **Account Deletion**: Allow users to delete their accounts
5. **Two-Factor Authentication**: Add 2FA support
6. **Activity Log**: Track profile changes
7. **Password History**: Prevent password reuse
8. **Rate Limiting**: Add specific rate limits for password changes

### Security Hardening
- Implement password complexity requirements
- Add password reset via email
- Log security-relevant events (password changes)
- Add session management (view/revoke active sessions)
- Implement account lockout after failed password attempts

## Files Modified/Created

### Backend
- ✅ Created: `backend/internal/handlers/user.go`
- ✅ Modified: `backend/internal/database/database.go` (added interface methods)
- ✅ Modified: `backend/internal/database/users.go` (added implementations)
- ✅ Modified: `backend/internal/server/routes.go` (added routes)
- ✅ Modified: `backend/internal/database/mock.go` (added stubs)
- ✅ Created: `backend/scripts/test_user_profile.sh`

### Frontend
- ✅ Already existed: `frontend/src/api/user.ts`
- ✅ Modified: `frontend/src/views/ProfileView.vue`

## Conclusion

The user profile management feature is now fully implemented with:
- Secure password handling
- Input validation
- Proper error handling
- Clean user interface
- Integration with existing authentication system
- Comprehensive documentation

Users can now view their profile, update their information, and change their passwords through a secure, well-designed interface.
