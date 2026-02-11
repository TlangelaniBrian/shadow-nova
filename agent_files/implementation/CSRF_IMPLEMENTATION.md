# CSRF Protection Implementation

## Overview
This document describes the CSRF (Cross-Site Request Forgery) protection implementation for Shadow Nova using gorilla/csrf.

## Installation Steps

### 1. Install Dependencies

```bash
cd backend
go get github.com/gorilla/csrf
go mod tidy
```

### 2. Generate CSRF Key

Generate a 32-byte key for CSRF protection:

```bash
openssl rand -base64 32
```

Add this key to your `.env` file:

```env
CSRF_KEY=<generated-key-here>
```

### 3. Set Environment

For development (non-HTTPS):
```env
ENV=development
```

For production (HTTPS only):
```env
ENV=production
```

## Architecture

### Backend Components

#### 1. CSRF Middleware (`backend/internal/middleware/csrf.go`)
- Protects all state-changing operations (POST, PUT, PATCH, DELETE)
- Uses gorilla/csrf with configurable security settings
- Returns JSON error response on CSRF validation failure
- Requires 32-byte CSRF_KEY from environment

#### 2. CSRF Token Endpoint (`backend/internal/handlers/auth.go`)
- `GET /api/csrf-token` - Fetches a CSRF token
- Returns token in both response body and `X-CSRF-Token` header
- Called by frontend on app initialization

#### 3. Routes Configuration (`backend/internal/server/routes.go`)
- CSRF middleware applied to all `/api/*` routes
- GET requests are exempt from CSRF validation (safe methods)
- All POST/PUT/PATCH/DELETE requests require valid CSRF token

### Frontend Components

#### 1. App Initialization (`frontend/src/main.ts`)
- Fetches CSRF token before mounting the Vue app
- Stores token in memory via `window.__CSRF_TOKEN__`
- Does NOT use localStorage (XSS vulnerable)

#### 2. API Client (`frontend/src/api/client.ts`)
- Request interceptor adds CSRF token to state-changing requests
- Response interceptor handles CSRF errors (403)
- Auto-retry mechanism: refetches token on CSRF failure and retries request

#### 3. useCSRF Composable (`frontend/src/composables/useCSRF.ts`)
- Vue composable for manual CSRF token refresh
- Useful for long-running sessions

#### 4. TypeScript Types (`frontend/src/types/global.d.ts`)
- Global Window interface extension for `__CSRF_TOKEN__`

## Security Features

### 1. Token Storage
- **Memory only**: Token stored in `window.__CSRF_TOKEN__`
- **No localStorage**: Prevents XSS attacks from stealing CSRF tokens
- **HttpOnly cookies**: Auth tokens in HttpOnly cookies (separate concern)

### 2. SameSite Protection
- `SameSiteStrictMode` prevents cross-site request attacks
- HTTPS-only in production (`csrf.Secure(secure)`)

### 3. Token Lifecycle
- Tokens fetched on app initialization
- Automatic refresh on CSRF validation failure
- Retry mechanism for failed requests

### 4. CORS Configuration
- `X-CSRF-Token` allowed in request headers
- `X-CSRF-Token` exposed in response headers
- Credentials enabled for cookie support

## Testing

### Run Backend Tests

```bash
cd backend
go test ./internal/middleware/csrf_test.go -v
```

### Test Coverage

1. ✅ POST without CSRF token returns 403
2. ✅ POST with valid CSRF token succeeds
3. ✅ GET requests work without CSRF token
4. ✅ PUT/DELETE without CSRF token returns 403
5. ✅ Middleware panics on invalid CSRF_KEY

### Manual Testing

#### 1. Test Token Endpoint
```bash
curl http://localhost:8080/api/csrf-token
```

Expected response:
```json
{
  "csrf_token": "abcd1234..."
}
```

#### 2. Test Protected Endpoint Without Token
```bash
curl -X POST http://localhost:8080/api/progress \
  -H "Content-Type: application/json" \
  -d '{"lesson_id": 1, "completed": true}'
```

Expected response:
```json
{
  "error": "CSRF token validation failed",
  "status": 403
}
```

#### 3. Test Protected Endpoint With Token
```bash
# Get token
TOKEN=$(curl -s http://localhost:8080/api/csrf-token | jq -r '.csrf_token')

# Use token
curl -X POST http://localhost:8080/api/progress \
  -H "Content-Type: application/json" \
  -H "X-CSRF-Token: $TOKEN" \
  -d '{"lesson_id": 1, "completed": true}'
```

## Protected Endpoints

All state-changing operations are now protected:

### Progress & Stats
- `POST /api/progress` - Update learning progress
- `POST /api/admin/settings/collector` - Update collector frequency

### Learning Paths
- `POST /api/paths` - Create learning path
- `POST /api/paths/{id}/modules` - Add module to path
- `POST /api/lessons` - Add lesson

### Projects & Submissions
- `POST /api/projects` - Create project (admin only)
- `POST /api/submissions` - Submit project solution

### Authentication
- `POST /api/register` - User registration
- `POST /api/login` - User login
- `POST /api/auth/logout` - User logout
- `POST /api/auth/google/verify` - Google token verification

## Troubleshooting

### Issue: CSRF_KEY must be 32 bytes
**Solution**: Generate a new key with `openssl rand -base64 32`

### Issue: CSRF validation fails on valid requests
**Solutions**:
1. Verify `CSRF_KEY` is set in `.env`
2. Check `ENV` variable (development vs production)
3. Ensure frontend is calling `/api/csrf-token` on init
4. Verify `withCredentials: true` in axios config

### Issue: Token not included in requests
**Solutions**:
1. Check `window.__CSRF_TOKEN__` is set
2. Verify interceptor is adding token to headers
3. Ensure request method is POST/PUT/PATCH/DELETE

### Issue: CORS errors with CSRF token
**Solution**: Verify `X-CSRF-Token` is in `AllowedHeaders` and `ExposedHeaders` in CORS middleware

## Migration Notes

### Existing Sessions
- Existing users will need to fetch CSRF token on next app load
- No database migration required
- Auth cookies remain valid

### Backward Compatibility
- GET requests unchanged (no CSRF required)
- OAuth callbacks unchanged (GET requests)
- Health checks and metrics unchanged

## Future Enhancements

1. **Rate Limiting**: Add per-IP rate limiting for CSRF token endpoint
2. **Token Rotation**: Implement periodic token rotation for long sessions
3. **Metrics**: Add Prometheus metrics for CSRF failures
4. **Logging**: Enhanced logging for CSRF validation failures

## References

- [gorilla/csrf Documentation](https://github.com/gorilla/csrf)
- [OWASP CSRF Prevention Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html)
- [Double Submit Cookie Pattern](https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html#double-submit-cookie)
