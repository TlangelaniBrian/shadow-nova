# JWT Token Revocation Implementation

## Overview
This document describes the JWT token revocation mechanism implemented for Shadow Nova. The implementation provides secure token blacklisting to invalidate JWTs before their natural expiration.

## Changes Made

### 1. Database Schema (`backend/internal/database/schema.sql`)
- Added `token_blacklist` table to store revoked JWT IDs (JTI)
- Includes columns: `jti`, `user_id`, `expires_at`, `blacklisted_at`, `reason`
- Added indexes on `expires_at` and `user_id` for efficient lookups

### 2. JWT Claims Update (`backend/internal/auth/auth.go`)
- Added `github.com/google/uuid` package import
- Updated `GenerateJWT()` to include a unique JTI (JWT ID) in token claims
- JTI is generated using `uuid.New().String()` and stored in `RegisteredClaims.ID`

### 3. Database Interface (`backend/internal/database/database.go`)
- Added four new methods to the `Service` interface:
  - `BlacklistToken()` - Add a token to the blacklist
  - `IsTokenBlacklisted()` - Check if a token is blacklisted
  - `BlacklistAllUserTokens()` - Placeholder for future implementation
  - `DeleteExpiredBlacklistedTokens()` - Cleanup expired tokens

### 4. Database Implementation (`backend/internal/database/users.go`)
- Implemented all token blacklist methods
- `IsTokenBlacklisted()` checks if JTI exists and hasn't expired
- `DeleteExpiredBlacklistedTokens()` removes expired entries to keep table size manageable

### 5. Auth Middleware (`backend/internal/middleware/auth.go`)
- Updated `AuthMiddleware` struct to include database service
- Updated `NewAuthMiddleware()` constructor to accept database parameter
- Modified `VerifyToken()` to:
  - Extract JTI from token claims
  - Check if token is blacklisted before allowing access
  - Return "Token has been revoked" error for blacklisted tokens

### 6. Logout Handler (`backend/internal/handlers/auth.go`)
- Updated `Logout()` handler to:
  - Extract token from cookie or Authorization header
  - Validate token and extract claims (including JTI)
  - Blacklist the token with reason "user_logout"
  - Clear the authentication cookie
- Added helper function `extractTokenFromRequest()` for token extraction

### 7. Routes Configuration (`backend/internal/server/routes.go`)
- Updated `NewAuthMiddleware()` call to pass database service
- Moved `/auth/logout` to protected routes section (requires authentication)
- Added background cleanup goroutine that runs daily to remove expired blacklisted tokens

### 8. Migration Script (`backend/internal/database/migrations/002_add_token_blacklist.sql`)
- Created migration script for easy deployment
- Can be applied to existing databases

### 9. Dependencies (`backend/go.mod`)
- Added `github.com/google/uuid v1.6.0` for JTI generation
- Added `github.com/gorilla/csrf v1.7.2` (was already in use)

## How It Works

### Token Generation Flow
1. User logs in (Google OAuth, GitHub OAuth, or traditional)
2. `GenerateJWT()` creates a new token with a unique JTI
3. Token is returned to user (in cookie or JSON response)

### Token Validation Flow
1. User makes authenticated request with JWT
2. `VerifyToken()` middleware extracts and validates token
3. JTI is extracted from token claims
4. Database is checked for blacklisted JTI
5. If blacklisted, request is rejected with 401 Unauthorized
6. If valid, request proceeds

### Logout Flow
1. User calls `/api/auth/logout` with valid token
2. Token is extracted from request
3. Claims are validated and JTI extracted
4. Token is added to blacklist table with expiration time
5. Cookie is cleared
6. User is logged out

### Cleanup Flow
1. Background goroutine runs every 24 hours
2. Queries for tokens where `expires_at <= NOW()`
3. Deletes expired entries
4. Logs number of deleted entries

## Security Benefits

1. **Immediate Revocation**: Tokens can be invalidated immediately without waiting for expiration
2. **Logout Support**: Proper logout functionality with server-side validation
3. **Compromise Recovery**: Stolen tokens can be blacklisted
4. **Multi-Device Logout**: Foundation for future "logout all devices" feature
5. **Audit Trail**: Reason tracking for token revocation

## Performance Considerations

1. **Database Lookup**: Each request requires one additional database query
   - Query is indexed on JTI (primary key) for O(1) lookup
   - Query includes expiration check to avoid returning expired entries

2. **Table Growth**: Blacklist table grows with each logout
   - Daily cleanup job removes expired entries
   - Table size is bounded by: (daily logouts) × (token lifetime in days)
   - Example: 1000 logouts/day × 1 day = 1000 entries max

3. **Optimization Options**:
   - Consider Redis/Memcached for high-traffic applications
   - Implement connection pooling (already done with pgxpool)
   - Add caching layer for recently checked JTIs

## Testing Instructions

### Prerequisites
```bash
cd backend
go mod download  # Install google/uuid dependency
```

### 1. Apply Migration
```bash
# Option A: Re-run schema (for development)
# The schema.sql already includes token_blacklist table

# Option B: Run migration script
psql $DATABASE_URL -f internal/database/migrations/002_add_token_blacklist.sql
```

### 2. Test Login and Token Generation
```bash
# Register a user
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","username":"testuser","password":"password123"}'

# Login
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Save the token from response
```

### 3. Test Protected Endpoint
```bash
# Use token to access protected endpoint
curl http://localhost:8080/api/stats \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"

# Should return user stats
```

### 4. Test Logout (Token Blacklisting)
```bash
# Logout (blacklist the token)
curl -X POST http://localhost:8080/api/auth/logout \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"

# Should return: {"message": "Logged out successfully"}
```

### 5. Verify Token is Blacklisted
```bash
# Try to use the same token again
curl http://localhost:8080/api/stats \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"

# Should return 401 with: "Token has been revoked"
```

### 6. Verify Database
```bash
# Check blacklist table
psql $DATABASE_URL -c "SELECT * FROM token_blacklist;"

# Should show the blacklisted JTI
```

### 7. Test Cleanup Job
```bash
# Wait 24 hours or manually trigger by restarting server
# Check logs for: "Cleaned X expired blacklisted tokens"

# Or manually run cleanup:
psql $DATABASE_URL -c "DELETE FROM token_blacklist WHERE expires_at <= NOW();"
```

## Future Enhancements

1. **Logout All Devices**
   - Implement `BlacklistAllUserTokens()`
   - Requires tracking all issued tokens per user
   - Could use token versioning or user-level revocation list

2. **Redis Integration**
   - Move blacklist to Redis for faster lookups
   - Use TTL feature for automatic expiration
   - Fallback to database for persistence

3. **Token Refresh**
   - Implement refresh tokens for longer sessions
   - Short-lived access tokens (15 min)
   - Long-lived refresh tokens (30 days)

4. **Admin Panel**
   - View active tokens per user
   - Manually revoke tokens
   - View blacklist statistics

5. **Rate Limiting**
   - Limit logout attempts per IP
   - Prevent blacklist table abuse

## Troubleshooting

### Issue: "Invalid token: missing JTI"
- Old tokens issued before this update don't have JTI
- Users need to login again to get new tokens with JTI

### Issue: Database query slow
- Verify indexes exist: `\d token_blacklist` in psql
- Check query plan: `EXPLAIN SELECT EXISTS(SELECT 1 FROM token_blacklist WHERE jti = 'xxx' AND expires_at > NOW());`

### Issue: Cleanup job not running
- Check server logs for "Token cleanup goroutine"
- Verify `s.collectorCtx` is properly initialized
- Restart server to trigger cleanup goroutine

## Conclusion

The JWT revocation mechanism provides a robust foundation for token management in Shadow Nova. It balances security (immediate revocation) with performance (indexed lookups, automatic cleanup) and provides a clear path for future enhancements.
