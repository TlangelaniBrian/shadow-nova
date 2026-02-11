# Sidebar Admin Link Fix - Implementation Summary

## Problem
The user management link in the Sidebar component was not visible to admin users because the frontend `User` interface was missing the `role` field.

## Changes Made

### 1. Frontend Type Definitions

#### `/Users/CT303853/Projects/Other_Projects/shadow-nova/frontend/src/api/auth.ts`
- Added `role?: string` field to the `User` interface (line 9)
- This allows the authentication response to include the user's role

#### `/Users/CT303853/Projects/Other_Projects/shadow-nova/frontend/src/api/user.ts`
- Added `role?: string` field to the `UserProfile` interface (line 10)
- This allows the user profile endpoint to return the user's role

### 2. Sidebar Component

#### `/Users/CT303853/Projects/Other_Projects/shadow-nova/frontend/src/components/layout/Sidebar.vue`
The Sidebar component is already properly configured:
- Imports user store correctly (line 14)
- Computes `isAdmin` based on `userStore.user?.role === 'admin'` (line 24)
- Conditionally adds Admin section to menu when `isAdmin.value` is true (lines 46-53)
- Admin section includes User Management link to `/admin/users` with Shield icon (line 50)
- Uses dynamic component rendering to show/hide admin section

### 3. Backend Configuration

#### Already Implemented:
- Database schema includes `user_role` column with CHECK constraint for 'user' or 'admin' values
- JWT tokens include role claim
- Authentication endpoints (`GoogleCallback`, `VerifyGoogleToken`, `Login`) return role in response
- `GetProfile` endpoint returns full user object including role
- Admin middleware checks role from JWT claims

## How It Works

### Authentication Flow:
1. User logs in via Google OAuth or email/password
2. Backend generates JWT with role claim
3. Backend returns user object including `role` field
4. Frontend stores user object in Pinia store and localStorage
5. Sidebar component reads user from store
6. If `user.role === 'admin'`, Admin section is rendered

### Sidebar Rendering:
```vue
<script setup>
const isAdmin = computed(() => userStore.user?.role === 'admin')

const menuItems = computed(() => {
  const items = [/* platform and account items */]

  if (isAdmin.value) {
    items.push({
      category: 'Admin',
      items: [
        { name: 'User Management', icon: Shield, path: '/admin/users' }
      ]
    })
  }

  return items
})
</script>
```

## Making a User an Admin

### Option 1: Database Direct Update
```sql
UPDATE users
SET user_role = 'admin'
WHERE email = 'user@example.com';
```

### Option 2: Using the Script
```bash
cd /Users/CT303853/Projects/Other_Projects/shadow-nova/backend
./scripts/make_user_admin.sh user@example.com
```

### Option 3: During User Creation
When creating a user via the database, set `user_role = 'admin'`:
```sql
INSERT INTO users (email, username, password_hash, user_role)
VALUES ('admin@example.com', 'admin', '$2a$...', 'admin');
```

## Testing the Fix

### 1. Create an Admin User
```bash
# From backend directory
./scripts/make_user_admin.sh your-email@example.com
```

### 2. Login as Admin
- Navigate to the frontend application
- Login with the admin user credentials
- The authentication response will include `"role": "admin"`

### 3. Verify Sidebar Shows Admin Section
- After login, check the Sidebar component
- You should see a new "ADMIN" section at the bottom
- It should contain "User Management" link with Shield icon
- Clicking it should navigate to `/admin/users`

### 4. Verify Regular Users Don't See Admin Section
- Logout and login as a regular user (without admin role)
- The Admin section should NOT appear in the sidebar

## Debugging

### If Admin Section Doesn't Show:

1. **Check User Object in Store**
```javascript
// In browser console
const userStore = useUserStore()
console.log(userStore.user)
// Should show: { id: '...', email: '...', role: 'admin', ... }
```

2. **Check localStorage**
```javascript
// In browser console
console.log(localStorage.getItem('user'))
// Should include "role":"admin"
```

3. **Check API Response**
```bash
# Login and check response
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"password"}' \
  -c cookies.txt

# Check profile endpoint
curl http://localhost:8080/api/v1/user/profile \
  -b cookies.txt
```

4. **Check Database**
```sql
SELECT id, email, username, user_role
FROM users
WHERE email = 'admin@example.com';
```

5. **Check JWT Claims**
```javascript
// In browser console, decode JWT
const token = document.cookie.split('auth_token=')[1].split(';')[0]
const payload = JSON.parse(atob(token.split('.')[1]))
console.log(payload)
// Should include "role":"admin"
```

## Key Files Modified

### Frontend:
- `/Users/CT303853/Projects/Other_Projects/shadow-nova/frontend/src/api/auth.ts`
- `/Users/CT303853/Projects/Other_Projects/shadow-nova/frontend/src/api/user.ts`

### Backend (Already Configured):
- `/Users/CT303853/Projects/Other_Projects/shadow-nova/backend/internal/handlers/auth.go`
- `/Users/CT303853/Projects/Other_Projects/shadow-nova/backend/internal/models/user.go`
- `/Users/CT303853/Projects/Other_Projects/shadow-nova/backend/internal/auth/auth.go`

### New Files:
- `/Users/CT303853/Projects/Other_Projects/shadow-nova/backend/scripts/make_user_admin.sh`

## Expected Behavior

### Admin User:
- Sees all regular navigation items (Dashboard, Learning Paths, Projects, Community, Profile)
- Sees Account section (Settings)
- Sees Admin section with User Management link
- Can access `/admin/users` page

### Regular User:
- Sees all regular navigation items
- Sees Account section
- Does NOT see Admin section
- Cannot access `/admin/users` (should be blocked by backend middleware)

## Security Notes

1. **Frontend Check is Not Security**: The `v-if="isAdmin"` in the Sidebar only controls visibility, not access
2. **Backend Enforcement**: The `/admin/*` routes are protected by `AdminOnly` middleware in the backend
3. **JWT Claims**: The role is stored in the JWT and validated on every protected request
4. **HttpOnly Cookies**: JWT tokens are stored in HttpOnly cookies, preventing XSS attacks

## Future Enhancements

Consider adding:
1. More granular role-based permissions (e.g., 'moderator', 'editor')
2. Role selection UI for admins to promote users
3. Audit logging when users are promoted to admin
4. Admin dashboard showing system statistics
5. Bulk user management operations
