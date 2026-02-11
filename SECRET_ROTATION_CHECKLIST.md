# Secret Rotation Checklist

## Overview
This checklist guides you through rotating all compromised secrets that were accidentally committed to git history.

## Critical: Complete ALL Steps

### Prerequisites
- [ ] Backup current `.env` files
- [ ] Note current secret values for reference
- [ ] Have access to all admin consoles

---

## 1. Google OAuth Credentials

### Status: COMPROMISED
**Client ID**: `762615404849-rp6h0v1q3f1fv9b3kdkrq4s3kjhqe8lk.apps.googleusercontent.com`

### Rotation Steps:
1. [ ] Go to [Google Cloud Console - Credentials](https://console.cloud.google.com/apis/credentials)
2. [ ] Select project or create new OAuth client
3. [ ] Delete old OAuth 2.0 Client ID
4. [ ] Create new OAuth 2.0 Client ID:
   - Application type: Web application
   - Authorized redirect URIs: `http://localhost:5173/auth/google/callback`
5. [ ] Copy new credentials
6. [ ] Update environment variables:
   ```bash
   GOOGLE_CLIENT_ID=<new_value>
   GOOGLE_CLIENT_SECRET=<new_value>
   ```

### Verification:
```bash
# Test Google OAuth login flow
curl http://localhost:3000/api/auth/google
```

---

## 2. GitHub OAuth Application

### Status: COMPROMISED
**Client ID**: `Ov23liXTlWNKP2TQ13kb`

### Rotation Steps:
1. [ ] Go to [GitHub Developer Settings - OAuth Apps](https://github.com/settings/developers)
2. [ ] Find the compromised application
3. [ ] Delete the old OAuth application
4. [ ] Create new OAuth App:
   - Application name: ShadowNova (or your app name)
   - Homepage URL: `http://localhost:5173`
   - Authorization callback URL: `http://localhost:5173/auth/github/callback`
5. [ ] Generate new client secret
6. [ ] Update environment variables:
   ```bash
   GITHUB_CLIENT_ID=<new_value>
   GITHUB_CLIENT_SECRET=<new_value>
   ```

### Verification:
```bash
# Test GitHub OAuth login flow
curl http://localhost:3000/api/auth/github
```

---

## 3. Google Gemini API Key

### Status: COMPROMISED
**Partial Key**: `AIzaSyAKg...`

### Rotation Steps:
1. [ ] Go to [Google Cloud Console - API Keys](https://console.cloud.google.com/apis/credentials)
2. [ ] Find the compromised API key
3. [ ] Delete or disable the old key
4. [ ] Create new API key:
   - Click "Create Credentials" > "API Key"
   - Restrict key to Gemini API
   - Add application restrictions if needed
5. [ ] Copy new key
6. [ ] Update environment variable:
   ```bash
   GEMINI_API_KEY=<new_value>
   ```

### Verification:
```bash
# Test Gemini API
curl -X POST http://localhost:3000/api/ai/generate \
  -H "Content-Type: application/json" \
  -d '{"prompt": "Hello"}'
```

---

## 4. JWT Secret

### Status: COMPROMISED
**Value**: `your-secret-key-change-this-in-production`

### Rotation Steps:
1. [ ] Generate new secure random secret:
   ```bash
   # Option 1: Using OpenSSL
   openssl rand -base64 64

   # Option 2: Using Node.js
   node -e "console.log(require('crypto').randomBytes(64).toString('base64'))"

   # Option 3: Using /dev/urandom
   head -c 64 /dev/urandom | base64
   ```
2. [ ] Update environment variable:
   ```bash
   JWT_SECRET=<generated_value>
   ```
3. [ ] **WARNING**: This will invalidate all existing user sessions
4. [ ] Notify users they need to log in again

### Verification:
```bash
# Test JWT token generation
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "test123"}'
```

---

## 5. Database URL

### Status: POTENTIALLY COMPROMISED
**Pattern**: Contains database credentials

### Rotation Steps:
1. [ ] Access your database admin panel
2. [ ] Change database user password
3. [ ] Update `DATABASE_URL` with new credentials:
   ```bash
   DATABASE_URL=postgresql://user:NEW_PASSWORD@host:port/database
   ```
4. [ ] Test database connection

### Verification:
```bash
# Test database connection
node -e "const { Pool } = require('pg'); const pool = new Pool({ connectionString: process.env.DATABASE_URL }); pool.query('SELECT NOW()').then(r => console.log('Connected:', r.rows[0])).catch(e => console.error('Error:', e)).finally(() => pool.end());"
```

---

## 6. Update All Environment Files

### Files to Update:
- [ ] `/backend/.env`
- [ ] `/frontend/.env`
- [ ] Any `.env.local`, `.env.production` files
- [ ] CI/CD environment variables
- [ ] Production deployment environment variables

### Checklist:
```bash
# Backend .env
cd /Users/CT303853/Projects/Other_Projects/shadow-nova/backend
cat > .env << 'EOF'
# Server
PORT=3000
NODE_ENV=development

# Database
DATABASE_URL=postgresql://user:password@localhost:5432/shadownova

# JWT
JWT_SECRET=<NEW_SECURE_VALUE>

# OAuth - Google
GOOGLE_CLIENT_ID=<NEW_VALUE>
GOOGLE_CLIENT_SECRET=<NEW_VALUE>
GOOGLE_REDIRECT_URI=http://localhost:5173/auth/google/callback

# OAuth - GitHub
GITHUB_CLIENT_ID=<NEW_VALUE>
GITHUB_CLIENT_SECRET=<NEW_VALUE>
GITHUB_REDIRECT_URI=http://localhost:5173/auth/github/callback

# AI Services
GEMINI_API_KEY=<NEW_VALUE>

# CORS
CORS_ORIGIN=http://localhost:5173
EOF
```

---

## 7. Post-Rotation Actions

### Immediate Actions:
- [ ] Restart all backend services
- [ ] Clear Redis/session cache if applicable
- [ ] Test all authentication flows
- [ ] Test all API endpoints that use secrets

### Security Hardening:
- [ ] Add `.env` to `.gitignore` (verify it's there)
- [ ] Remove `.env` from git history (use cleanup script)
- [ ] Enable secret scanning on GitHub repository
- [ ] Set up environment variable management (e.g., 1Password, Vault)
- [ ] Document secret rotation procedure for team

### Verification Commands:
```bash
# Check .env is in .gitignore
grep -q "^\.env$" .gitignore && echo "✓ .env in .gitignore" || echo "✗ Missing from .gitignore"

# Verify no secrets in current staging area
git diff --cached | grep -i "secret\|api_key\|client_id" && echo "✗ Secrets in staging" || echo "✓ No secrets staged"

# Check git history for .env (should return nothing after cleanup)
git log --all --full-history --oneline -- '*.env'
```

---

## 8. Monitor for Unauthorized Access

### Immediate Monitoring (First 48 Hours):
- [ ] Check Google Cloud Console for unusual API usage
- [ ] Check GitHub OAuth app for unauthorized access attempts
- [ ] Monitor database logs for suspicious queries
- [ ] Review application logs for failed auth attempts

### Setup Alerts:
- [ ] Enable GitHub secret scanning alerts
- [ ] Set up Google Cloud billing alerts
- [ ] Configure database connection monitoring
- [ ] Enable rate limiting on auth endpoints

---

## Reference: Quick Command Summary

```bash
# Generate JWT Secret
openssl rand -base64 64

# Verify .env not in git
git ls-files | grep .env

# Test backend health
curl http://localhost:3000/health

# Test database connection
psql $DATABASE_URL -c "SELECT version();"

# Restart backend
cd backend && npm run dev
```

---

## Completion Checklist

- [ ] All secrets rotated
- [ ] All environment files updated
- [ ] Services restarted and tested
- [ ] Git history cleaned (run cleanup script)
- [ ] Monitoring enabled
- [ ] Team notified
- [ ] Documentation updated

**Date Completed**: _______________
**Completed By**: _______________
**Verification**: Run `./scripts/verify-urgent-fixes.sh`
