# Shadow Nova - Google OAuth Authentication Guide

## Overview

Shadow Nova uses **Google OAuth 2.0** for secure, passwordless authentication. Users can sign in with their Google accounts, and the backend issues JWT tokens for subsequent API requests.

## Architecture Flow

```
1. User clicks "Sign in with Google" → Frontend initiates OAuth flow
2. User authenticates with Google → Google returns ID token
3. Frontend sends ID token to backend → Backend verifies with Google
4. Backend generates JWT token → Returns JWT + user info
5. Frontend stores JWT → Uses for authenticated API calls
```

## Setup Instructions

### 1. Environment Configuration

Copy `.env.example` to `.env` and add your Google OAuth credentials:

**Backend (.env in project root):**

```env
GOOGLE_CLIENT_ID=54718351338-tbied6l38ldbiqj4nmghqgi9l4lf1ge9.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-actual-secret-here
GOOGLE_REDIRECT_URL=http://localhost:8080/auth/callback
JWT_SECRET=your-super-secret-jwt-key
```

**Frontend (.env in frontend directory):**

```env
VITE_API_URL=http://localhost:3000
VITE_GOOGLE_CLIENT_ID=54718351338-tbied6l38ldbiqj4nmghqgi9l4lf1ge9.apps.googleusercontent.com
```

### 2. Google Cloud Console Setup

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Select your project or create a new one
3. Navigate to **APIs & Services** → **Credentials**
4. Your Client ID is already created: `54718351338-tbied6l38ldbiqj4nmghqgi9l4lf1ge9`
5. Add authorized redirect URIs:
   - `http://localhost:8080/auth/callback`
   - `http://localhost:3000/api/auth/google/callback`
6. Add authorized JavaScript origins:
   - `http://localhost:8080`
   - `http://localhost:3000`

## API Endpoints

### Public Endpoints

#### GET `/api/auth/google`

Get Google OAuth URL for redirect flow (server-side)

**Response:**

```json
{
  "auth_url": "https://accounts.google.com/o/oauth2/v2/auth?..."
}
```

#### GET `/api/auth/google/callback`

OAuth callback endpoint (handles server-side flow)

**Query Parameters:**

- `code` - Authorization code from Google
- `state` - State token for CSRF protection

**Response:**

```json
{
  "token": "eyJhbG...",
  "user": {
    "id": "google-user-id",
    "email": "user@example.com",
    "name": "John Doe",
    "picture": "https://..."
  }
}
```

#### POST `/api/auth/google/verify`

Verify Google ID token (client-side flow) - **Recommended**

**Request:**

```json
{
  "id_token": "eyJhbG..."
}
```

**Response:**

```json
{
  "token": "our-jwt-token",
  "user": {
    "id": "google-user-id",
    "email": "user@example.com",
    "name": "John Doe",
    "picture": "https://..."
  }
}
```

### Protected Endpoints

All protected endpoints require the JWT token in the `Authorization` header:

```
Authorization: Bearer eyJhbG...
```

Example protected endpoints:

- `POST /api/progress` - Update learning progress
- `GET /api/paths` - Get learning paths
- `GET /api/projects` - Get projects

## Frontend Integration

### Using the GoogleSignIn Component

```vue
<script setup lang="ts">
import GoogleSignIn from '@/components/GoogleSignIn.vue';
</script>

<template>
  <div class="login-page">
    <h1>Welcome to Shadow Nova</h1>
    <GoogleSignIn />
  </div>
</template>
```

### Manual Integration

```typescript
// Load Google Sign-In script
const script = document.createElement('script');
script.src = 'https://accounts.google.com/gsi/client';
document.head.appendChild(script);

script.onload = () => {
  google.accounts.id.initialize({
    client_id: 'YOUR_CLIENT_ID',
    callback: handleGoogleResponse,
  });

  google.accounts.id.renderButton(document.getElementById('buttonDiv'), {
    theme: 'outline',
    size: 'large',
  });
};

// Handle response
async function handleGoogleResponse(response) {
  const res = await fetch('http://localhost:3000/api/auth/google/verify', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ id_token: response.credential }),
  });

  const data = await res.json();
  localStorage.setItem('auth_token', data.token);
  // Redirect to dashboard
}
```

### Making Authenticated Requests

```typescript
const token = localStorage.getItem('auth_token');

const response = await fetch('http://localhost:3000/api/progress', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    Authorization: `Bearer ${token}`,
  },
  body: JSON.stringify({
    path_id: 'frontend-vue',
    module_index: 1,
    completed: true,
  }),
});
```

## Security Features

### Backend

1. **Google Token Verification**: Validates tokens directly with Google's OIDC provider
2. **JWT Generation**: Creates short-lived JWT tokens (24-hour expiry)
3. **HMAC Signing**: Uses HS256 for JWT signing
4. **Token Validation**: Middleware verifies JWT on protected routes

### Frontend

1. **Secure Storage**: Tokens stored in localStorage (consider httpOnly cookies for production)
2. **Auto Prompt**: Google One Tap can be enabled for seamless login
3. **CORS Protection**: Backend validates origin headers

## Testing

### Test Google Sign-In Flow

1. Start the backend:

   ```bash
   cd backend
   go run main.go
   ```

2. Start the frontend:

   ```bash
   cd frontend
   pnpm dev
   ```

3. Navigate to the login page
4. Click "Sign in with Google"
5. Authenticate with your Google account
6. Verify you receive a JWT token

### Test API Authentication

```bash
# Get token from Google Sign-In, then:
TOKEN="your-jwt-token-here"

curl -X POST http://localhost:3000/api/progress \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "path_id": "frontend-vue",
    "module_index": 1,
    "completed": true
  }'
```

## Production Considerations

### Required Changes

1. **HTTPS**: Must use HTTPS in production

   - Update redirect URIs to `https://yourdomain.com/auth/callback`
   - Update JavaScript origins to `https://yourdomain.com`

2. **Secret Management**:

   - Store `GOOGLE_CLIENT_SECRET` in secure vault (AWS Secrets Manager, etc.)
   - Use strong `JWT_SECRET` (min 32 characters, random)

3. **Token Storage**:

   - Use httpOnly cookies instead of localStorage
   - Implement refresh token rotation
   - Add CSRF protection

4. **Database Integration**:

   - Save user records on first login
   - Update last login timestamp
   - Implement user session management

5. **Rate Limiting**:
   - Limit login attempts per IP
   - Implement account lockout after failed attempts

## Troubleshooting

### "redirect_uri_mismatch" Error

- Ensure redirect URI in Google Console exactly matches your backend URL
- Include protocol and port: `http://localhost:3000/api/auth/google/callback`

### "Invalid Token" Error

- Check JWT_SECRET is set correctly
- Verify token hasn't expired
- Ensure Authorization header format: `Bearer <token>`

### Google Sign-In Button Not Showing

- Check browser console for errors
- Verify VITE_GOOGLE_CLIENT_ID is set
- Ensure Google script is loaded before initialization

## Next Steps

- [ ] Implement user database storage
- [ ] Add refresh token mechanism
- [ ] Set up session management
- [ ] Add role-based access control (RBAC)
- [ ] Implement logout endpoint
- [ ] Add account deletion flow
