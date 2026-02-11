# Sidebar Admin Link - Verification Checklist

## Pre-Verification Setup

### 1. Create an Admin User
```bash
cd /Users/CT303853/Projects/Other_Projects/shadow-nova/backend
./scripts/make_user_admin.sh your-test-email@example.com
```

Or directly in database:
```sql
UPDATE users SET user_role = 'admin' WHERE email = 'your-test-email@example.com';
```

### 2. Start Backend Server
```bash
cd /Users/CT303853/Projects/Other_Projects/shadow-nova/backend
go run main.go
```

### 3. Start Frontend Development Server
```bash
cd /Users/CT303853/Projects/Other_Projects/shadow-nova/frontend
npm run dev
```

## Verification Steps

### Test 1: Admin User Login and Sidebar Visibility

**Steps:**
1. Navigate to http://localhost:5173 (or your frontend port)
2. Login with admin user credentials
3. Wait for authentication to complete

**Expected Results:**
- [ ] Login successful
- [ ] Redirected to `/dashboard`
- [ ] Sidebar is visible on the left
- [ ] Sidebar contains "PLATFORM" section with: Dashboard, Learning Paths, Projects, Community, Profile
- [ ] Sidebar contains "ACCOUNT" section with: Settings
- [ ] **Sidebar contains "ADMIN" section with: User Management link**
- [ ] User Management link has Shield icon
- [ ] User Management link is styled consistently with other nav items

**Browser Console Checks:**
```javascript
// Check user object in store
const userStore = useUserStore()
console.log('User:', userStore.user)
// Should show: { id: '...', email: '...', name: '...', role: 'admin' }

// Check localStorage
console.log('LocalStorage User:', localStorage.getItem('user'))
// Should include "role":"admin"
```

### Test 2: Admin Link Navigation

**Steps:**
1. While logged in as admin, click "User Management" link in sidebar
2. Observe URL and page content

**Expected Results:**
- [ ] URL changes to `/admin/users`
- [ ] User Management page loads successfully
- [ ] Admin interface is displayed
- [ ] No console errors

### Test 3: Regular User Cannot See Admin Section

**Steps:**
1. Logout from admin account
2. Login with regular (non-admin) user account
3. Check sidebar

**Expected Results:**
- [ ] Login successful
- [ ] Redirected to `/dashboard`
- [ ] Sidebar is visible
- [ ] Sidebar contains PLATFORM section
- [ ] Sidebar contains ACCOUNT section
- [ ] **Sidebar does NOT contain ADMIN section**
- [ ] No User Management link visible

**Browser Console Checks:**
```javascript
const userStore = useUserStore()
console.log('User:', userStore.user)
// Should show: { id: '...', email: '...', name: '...', role: 'user' }
```

### Test 4: Regular User Cannot Access Admin Route

**Steps:**
1. While logged in as regular user, manually navigate to http://localhost:5173/admin/users
2. Observe behavior

**Expected Results:**
- [ ] Browser redirects to `/dashboard`
- [ ] Admin page is NOT accessible
- [ ] User remains on dashboard

### Test 5: Backend API Returns Role

**Steps:**
1. Login and capture auth token
2. Test profile endpoint

**Using curl:**
```bash
# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"yourpassword"}' \
  -c cookies.txt -v

# Get profile
curl http://localhost:8080/api/v1/user/profile \
  -b cookies.txt
```

**Expected Results:**
- [ ] Login response includes `"role":"admin"`
- [ ] Profile endpoint returns user object with role field
- [ ] Response structure matches:
```json
{
  "id": 123,
  "email": "admin@example.com",
  "username": "admin",
  "role": "admin",
  "created_at": "...",
  "updated_at": "..."
}
```

### Test 6: JWT Token Contains Role Claim

**Steps:**
1. Login as admin
2. Extract auth token from cookies
3. Decode JWT payload

**Browser Console:**
```javascript
// Get token from cookie
const cookies = document.cookie.split('; ')
const authCookie = cookies.find(c => c.startsWith('auth_token='))
const token = authCookie?.split('=')[1]

// Decode JWT (base64)
const payload = JSON.parse(atob(token.split('.')[1]))
console.log('JWT Claims:', payload)
```

**Expected Results:**
- [ ] JWT payload contains `"role": "admin"`
- [ ] JWT payload contains other claims: `user_id`, `email`, `name`, `exp`, `iat`, `iss`

### Test 7: Sidebar State Persistence

**Steps:**
1. Login as admin
2. Verify Admin section is visible
3. Refresh the page (F5)
4. Wait for page to reload

**Expected Results:**
- [ ] User remains logged in
- [ ] Admin section is still visible in sidebar
- [ ] No additional login required

### Test 8: Mobile Sidebar Behavior

**Steps:**
1. Login as admin
2. Resize browser window to mobile size (< 768px width)
3. Click hamburger menu to open sidebar
4. Check if Admin section is visible

**Expected Results:**
- [ ] Sidebar slides in from left
- [ ] All sections visible including Admin
- [ ] User Management link is clickable
- [ ] Clicking link closes sidebar and navigates

### Test 9: Hover and Active States

**Steps:**
1. Login as admin
2. Hover over "User Management" link
3. Click "User Management" link
4. Verify active state

**Expected Results:**
- [ ] Hover state: background changes to light gray/purple tint
- [ ] Active state: background highlighted, text bold
- [ ] Icon and text maintain consistent styling
- [ ] Transition is smooth

### Test 10: Backend Admin Middleware Protection

**Steps:**
1. Login as regular user
2. Get auth token
3. Try to access admin endpoint directly

**Using curl:**
```bash
# Login as regular user
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"regular@example.com","password":"password"}' \
  -c regular-cookies.txt

# Try to access admin endpoint
curl http://localhost:8080/api/v1/admin/users \
  -b regular-cookies.txt -v
```

**Expected Results:**
- [ ] Response status: 403 Forbidden
- [ ] Error message: "Access denied: admin role required"
- [ ] Admin data is NOT returned

## Common Issues and Solutions

### Issue: Admin section not showing despite role being 'admin'

**Debug Steps:**
1. Check if user object in store has role field:
   ```javascript
   console.log(useUserStore().user)
   ```
2. Check if localStorage has role:
   ```javascript
   console.log(JSON.parse(localStorage.getItem('user')))
   ```
3. Clear localStorage and re-login:
   ```javascript
   localStorage.clear()
   // Then refresh and re-login
   ```

### Issue: Role is null or undefined

**Possible Causes:**
- Backend not returning role in auth response
- Frontend User interface missing role field (should be fixed now)
- Old cached data in localStorage

**Solution:**
1. Verify backend returns role (check Network tab in DevTools)
2. Clear browser cache and localStorage
3. Re-login

### Issue: TypeScript errors about role property

**Solution:**
- Ensure `/Users/CT303853/Projects/Other_Projects/shadow-nova/frontend/src/api/auth.ts` has `role?: string` in User interface
- Ensure `/Users/CT303853/Projects/Other_Projects/shadow-nova/frontend/src/api/user.ts` has `role?: string` in UserProfile interface
- Restart TypeScript server in IDE

### Issue: Route guard not working

**Solution:**
1. Check router configuration has `requiresAdmin: true` meta
2. Verify navigation guard logic in `/Users/CT303853/Projects/Other_Projects/shadow-nova/frontend/src/router/index.ts`
3. Clear router cache and reload

## Success Criteria

All checkboxes in Tests 1-10 should be checked (✓) for the implementation to be considered complete and working.

### Critical Tests (Must Pass):
- Test 1: Admin user sees Admin section
- Test 2: Admin link navigates correctly
- Test 3: Regular user does NOT see Admin section
- Test 4: Regular user cannot access admin route
- Test 5: Backend API returns role
- Test 10: Backend middleware protects admin routes

### Nice-to-Have Tests:
- Test 6: JWT contains role
- Test 7: State persistence
- Test 8: Mobile behavior
- Test 9: UI states

## Additional Verification

### Database Verification
```sql
-- Check user roles
SELECT id, email, username, user_role, created_at
FROM users
ORDER BY user_role DESC, created_at DESC;

-- Count users by role
SELECT user_role, COUNT(*) as count
FROM users
GROUP BY user_role;
```

### API Endpoint Verification
```bash
# List all admin routes
curl http://localhost:8080/api/v1/admin/users \
  -H "Cookie: auth_token=<admin-token>" | jq

# Verify admin user detail
curl http://localhost:8080/api/v1/admin/users/1 \
  -H "Cookie: auth_token=<admin-token>" | jq
```

## Sign-off

Once all tests pass:
- [ ] Admin sidebar link is visible and functional
- [ ] Regular users cannot see admin section
- [ ] Backend properly protects admin routes
- [ ] Role field is properly typed and returned
- [ ] No TypeScript errors
- [ ] No console errors
- [ ] Documentation is complete

**Tester Name:** _______________
**Date:** _______________
**Test Environment:** Development / Staging / Production
**Browser Tested:** _______________
**Result:** PASS / FAIL

**Notes:**
