# Shadow Nova - Development Guide

## Running Both Frontend & Backend

### Option 1: PowerShell Script (Recommended for Windows)

```powershell
.\dev.ps1
```

This script:

- Starts backend on http://localhost:3000
- Starts frontend on http://localhost:5173
- Shows logs from both in one terminal
- Press Ctrl+C to stop both

### Option 2: Using npm-run-all

```bash
# First time setup
npm install

# Run dev servers
npm run dev
```

### Option 3: Separate Terminals

**Terminal 1 - Backend:**

```bash
cd backend
go run main.go
```

**Terminal 2 - Frontend:**

```bash
cd frontend
pnpm dev
```

### Option 4: Docker Compose (Production-like)

```bash
docker-compose up
```

## Debugging

### Frontend Debugging (VS Code)

Add to `.vscode/launch.json`:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "type": "chrome",
      "request": "launch",
      "name": "Debug Frontend",
      "url": "http://localhost:5173",
      "webRoot": "${workspaceFolder}/frontend/src"
    }
  ]
}
```

### Backend Debugging (VS Code)

Add to `.vscode/launch.json`:

```json
{
  "name": "Debug Backend",
  "type": "go",
  "request": "launch",
  "mode": "debug",
  "program": "${workspaceFolder}/backend",
  "env": {
    "PORT": "3000",
    "DATABASE_URL": "postgres://user:password@localhost:5432/shadownova"
  }
}
```

### Full Stack Debugging

Add a compound configuration:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "type": "chrome",
      "request": "launch",
      "name": "Debug Frontend",
      "url": "http://localhost:5173",
      "webRoot": "${workspaceFolder}/frontend/src"
    },
    {
      "name": "Debug Backend",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/backend"
    }
  ],
  "compounds": [
    {
      "name": "Debug Full Stack",
      "configurations": ["Debug Frontend", "Debug Backend"],
      "presentation": {
        "order": 1
      }
    }
  ]
}
```

Then press F5 and select "Debug Full Stack" to debug both!

## Environment Setup

Make sure you have `.env` files:

**Frontend (.env):**

```env
VITE_API_URL=http://localhost:3000
VITE_GOOGLE_CLIENT_ID=your-client-id
```

**Backend (.env in root):**

```env
PORT=3000
GOOGLE_CLIENT_ID=your-client-id
GOOGLE_CLIENT_SECRET=your-secret
JWT_SECRET=your-jwt-secret
```

## Hot Reload

- **Frontend**: Vite provides instant hot module replacement
- **Backend**: Use `air` for Go hot reload:

```bash
# Install air
go install github.com/air-verse/air@latest

# Run backend with hot reload
cd backend
air
```

## Ports

- Frontend: http://localhost:5173 (Vite dev server)
- Backend: http://localhost:3000 (Go API)
- Grafana: http://localhost:3001
- Prometheus: http://localhost:9090

## Troubleshooting

### Port Already in Use

```powershell
# Find process on port 3000
netstat -ano | findstr :3000

# Kill process
taskkill /PID <process_id> /F
```

### Database Connection

```bash
# Start PostgreSQL via Docker
docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=password postgres:16-alpine
```
