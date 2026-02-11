# Admin User Management - Quick Start Guide

## Setup (One-Time)

### 1. Apply Database Migration

```bash
cd /Users/CT303853/Projects/Other_Projects/shadow-nova/backend
psql $DATABASE_URL -f migrations/add_deleted_at_to_users.sql
```

Or if using docker-compose:
```bash
docker-compose exec postgres psql -U user -d shadownova -f /path/to/migrations/add_deleted_at_to_users.sql
```

### 2. Create Your First Admin User

Option A: Update existing user to admin:
```sql
UPDATE users SET user_role = 'admin' WHERE email = 'your-email@example.com';
```

Option B: Create new admin user:
```sql
INSERT INTO users (email, username, password_hash, user_role)
VALUES (
    'admin@shadownova.com',
    'admin',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', -- password: admin123
    'admin'
);
```

## Daily Usage

### Access Admin Panel

1. Start the servers (if not running):
   ```bash
   # Terminal 1 - Backend
   cd backend && go run main.go

   # Terminal 2 - Frontend
   cd frontend && pnpm dev
   ```

2. Open browser to: `http://localhost:5173`

3. Log in with admin credentials

4. Click "User Management" in the Admin section of the sidebar

5. You should see: `http://localhost:5173/admin/users`

### Common Tasks

**Create a User**
1. Click "Create User" button
2. Fill in email, username, password, and role
3. Click "Save"

**Edit a User**
1. Click "Edit" next to any user
2. Modify email, username, or role
3. Click "Save"

**Delete a User**
1. Click "Delete" next to any user
2. Confirm deletion

**View All Users**
- Use pagination at the bottom to navigate pages
- Each page shows 20 users

## API Testing (curl)

### Get CSRF Token First
```bash
curl -c cookies.txt http://localhost:8080/api/v1/csrf-token
CSRF_TOKEN=$(cat cookies.txt | grep csrf | awk '{print $7}')
```

### Login to Get Auth Token
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -H "X-CSRF-Token: $CSRF_TOKEN" \
  -c cookies.txt \
  -d '{"email":"admin@shadownova.com","password":"admin123"}'
```

### List Users
```bash
curl -X GET "http://localhost:8080/api/v1/admin/users?page=1&limit=20" \
  -b cookies.txt
```

### Create User
```bash
curl -X POST http://localhost:8080/api/v1/admin/users \
  -H "Content-Type: application/json" \
  -H "X-CSRF-Token: $CSRF_TOKEN" \
  -b cookies.txt \
  -d '{
    "email":"newuser@example.com",
    "username":"newuser",
    "password":"password123",
    "role":"user"
  }'
```

### Update User
```bash
curl -X PUT http://localhost:8080/api/v1/admin/users/2 \
  -H "Content-Type: application/json" \
  -H "X-CSRF-Token: $CSRF_TOKEN" \
  -b cookies.txt \
  -d '{"role":"admin"}'
```

### Delete User
```bash
curl -X DELETE http://localhost:8080/api/v1/admin/users/2 \
  -H "X-CSRF-Token: $CSRF_TOKEN" \
  -b cookies.txt
```

## Troubleshooting

### "Access denied: admin role required"
- Check your user role in database: `SELECT user_role FROM users WHERE email = 'your-email';`
- Update to admin if needed: `UPDATE users SET user_role = 'admin' WHERE email = 'your-email';`

### "Cannot access /admin/users route"
- Ensure you're logged in as an admin user
- Check browser console for errors
- Verify localStorage has user object: `localStorage.getItem('user')`

### "Users list is empty"
- Check database has users: `SELECT * FROM users WHERE deleted_at IS NULL;`
- Check network tab for API errors
- Verify backend is running on port 8080

### "CSRF token error"
- Refresh the page to get new CSRF token
- Check that cookies are enabled in browser

### Backend won't compile
- Run: `go mod tidy`
- Check Go version is 1.21+
- Verify all imports are correct

### Frontend won't load
- Run: `pnpm install`
- Check Node version is 18+
- Clear pnpm cache if needed: `pnpm store prune`

## Security Notes

- Never commit passwords in plain text
- Always use bcrypt for password hashing
- Admin users cannot delete themselves
- Admin users cannot edit themselves via admin panel
- All endpoints require authentication + admin role
- CSRF protection is enabled for all state-changing operations
- Users are soft-deleted (can be recovered)

## File Locations

- **Backend Handler**: `/backend/internal/handlers/admin_users.go`
- **Frontend View**: `/frontend/src/views/admin/UsersView.vue`
- **API Client**: `/frontend/src/api/admin/users.ts`
- **Database Schema**: `/backend/internal/database/schema.sql`
- **Migration**: `/backend/migrations/add_deleted_at_to_users.sql`

## Support

For detailed testing instructions, see: `ADMIN_USER_MANAGEMENT_TESTING.md`

For complete implementation details, see: `ADMIN_USER_MANAGEMENT_SUMMARY.md`
