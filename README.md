# Shadow Nova

> A comprehensive learning platform for developers to master modern web technologies through structured learning paths, hands-on projects, and AI-curated content.

## Overview

Shadow Nova helps developers navigate the overwhelming landscape of web development resources by providing:

- **Structured Learning Paths** - Curated courses for Vue.js, Go, and DevOps
- **AI-Powered Content** - Automated aggregation and summarization from YouTube, RSS feeds, and blogs
- **Project-Based Learning** - Real-world projects with GitHub integration
- **Progress Tracking** - Gamified learning with stats and achievements
- **Google OAuth** - Seamless authentication

## Tech Stack

### Frontend
- **Vue 3.5** - Modern composition API with TypeScript
- **Vite 7** - Lightning-fast build tooling
- **Tailwind CSS 4** - Utility-first styling
- **Pinia** - State management
- **Radix Vue** - Accessible UI components
- **Axios** - HTTP client with interceptors

### Backend
- **Go 1.24** - High-performance API server
- **PostgreSQL** - Relational database with pgx driver
- **Chi Router** - Lightweight HTTP routing
- **JWT Authentication** - Token-based auth with Google OAuth 2.0
- **Prometheus** - Metrics instrumentation
- **Unleash** - Feature flag management

### Infrastructure
- **Docker** - Container orchestration
- **GitHub Actions** - CI/CD pipeline
- **Grafana Stack** - PLG (Prometheus, Loki, Grafana) observability
- **AWS Ready** - Amplify, App Runner, RDS deployment guides

## Quick Start

### Prerequisites
- **Go 1.24+**
- **Node.js 20+**
- **pnpm**
- **PostgreSQL 16+**
- **Docker** (optional)

### Option 1: Local Development

```bash
# Clone the repository
git clone https://github.com/TlangelaniBrian/shadow-nova.git
cd shadow-nova

# Set up environment variables
cp .env.example .env
# Edit .env with your Google OAuth credentials

# Install dependencies
cd frontend && pnpm install
cd ../backend && go mod download

# Run development servers
cd ..
npm run dev  # Runs frontend (5173) and backend (3000)
```

### Option 2: Docker Compose

```bash
# Start all services (frontend, backend, database, observability)
docker-compose up -d

# Access services:
# Frontend: http://localhost:8080
# Backend: http://localhost:3000
# Grafana: http://localhost:3001 (admin/admin)
# Prometheus: http://localhost:9090
```

## Project Structure

```
shadow-nova/
├── frontend/               # Vue 3 TypeScript application
│   ├── src/
│   │   ├── views/         # Page components
│   │   ├── components/    # Reusable UI components
│   │   ├── composables/   # Vue composition functions
│   │   ├── stores/        # Pinia state stores
│   │   ├── api/           # API client and services
│   │   ├── router/        # Vue Router config
│   │   └── types/         # TypeScript definitions
│   └── observability/     # PLG stack configs
├── backend/               # Go API server
│   ├── main.go           # Application entry point
│   └── internal/
│       ├── server/       # HTTP server setup
│       ├── handlers/     # Request handlers
│       ├── database/     # PostgreSQL data layer
│       ├── middleware/   # Auth, CORS, metrics
│       ├── auth/         # Google OAuth implementation
│       ├── ai/           # Gemini AI integration
│       ├── collector/    # Content aggregation service
│       └── models/       # Data models
├── .github/workflows/    # CI/CD pipelines
└── docker-compose.yml    # Local orchestration
```

## Core Features

### 1. Authentication
- **Google OAuth 2.0** - Passwordless authentication
- **JWT Tokens** - 24-hour expiry with HS256 signing
- **GitHub Integration** - Connect for project tracking

### 2. Learning Paths
- **Hierarchical Structure** - Paths → Modules → Lessons
- **Multiple Content Types** - Videos, articles, quizzes
- **Progress Tracking** - Track completion per lesson
- **Difficulty Levels** - Beginner to advanced

### 3. AI-Powered Content Curation
- **Automated Collection** - Fetch from YouTube, RSS feeds, Twitter
- **AI Summarization** - Gemini-powered summaries and tagging
- **Difficulty Rating** - Automatic content classification
- **Personalized Feed** - Curated content based on progress

### 4. Project Management
- **Hands-On Projects** - Aligned with learning paths
- **GitHub Integration** - Submit via repository link
- **Pull Request Tracking** - Automatic PR detection
- **Instructor Feedback** - Review and approval workflow

### 5. Observability
- **Prometheus Metrics** - Request rates, latency, custom metrics
- **Loki Logs** - Centralized log aggregation
- **Grafana Dashboards** - Unified visualization
- **Health Checks** - `/health` and `/metrics` endpoints

## API Endpoints

### Public Routes
```
GET  /health                          # Health check
GET  /metrics                         # Prometheus metrics
GET  /api/auth/google                 # Google OAuth login
GET  /api/auth/google/callback        # OAuth callback
POST /api/auth/google/verify          # Verify Google token
POST /api/register                    # Traditional signup
POST /api/login                       # Traditional login
GET  /api/projects                    # List projects
```

### Protected Routes (Require JWT)
```
GET  /api/paths                       # List learning paths
GET  /api/paths/{id}                  # Get path details
POST /api/progress                    # Update progress
GET  /api/stats                       # User statistics
POST /api/submissions                 # Submit project
GET  /api/auth/github/connect         # Connect GitHub
```

## Development

### Running Tests
```bash
# Backend tests
cd backend && go test -v ./...

# Frontend tests (when configured)
cd frontend && pnpm test
```

### Code Quality
```bash
# Frontend linting
cd frontend
pnpm lint        # Check for issues
pnpm lint:fix    # Auto-fix issues
pnpm format      # Prettier formatting
pnpm type-check  # TypeScript validation
```

### Database Migrations
The backend automatically runs migrations on startup from `backend/internal/database/schema.sql`.

### Hot Reload
- **Frontend**: Vite provides instant HMR
- **Backend**: Use [Air](https://github.com/air-verse/air) for Go hot reload
  ```bash
  go install github.com/air-verse/air@latest
  cd backend && air
  ```

## Deployment

### AWS (Student-Friendly)
- **Frontend**: AWS Amplify (auto-detects Vue/Vite)
- **Backend**: AWS App Runner or EC2 t2.micro (free tier)
- **Database**: Amazon RDS PostgreSQL db.t3.micro

See [aws-deployment.md](aws-deployment.md) for detailed instructions.

### Docker Production
```bash
# Build images
docker build -t shadow-nova-frontend -f Dockerfile.frontend ./frontend
docker build -t shadow-nova-backend -f Dockerfile.backend ./backend

# Run with docker-compose
docker-compose -f docker-compose.yml up -d
```

### CI/CD Pipeline
GitHub Actions automatically:
1. Builds frontend (pnpm) and backend (Go)
2. Runs tests
3. Creates Docker images
4. Publishes to GitHub Container Registry
5. Deploys to DEV → Staging → Production

**Note**: The workflow currently fails due to GitHub account billing issues. Fix by updating billing in GitHub Settings.

## Configuration

### Frontend Environment Variables
```env
VITE_API_URL=http://localhost:3000
VITE_GOOGLE_CLIENT_ID=your-google-client-id
```

### Backend Environment Variables
```env
PORT=3000
DATABASE_URL=postgres://user:password@localhost:5432/shadownova
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
JWT_SECRET=your-secret-key
GEMINI_API_KEY=your-gemini-api-key
```

See [.env.example](.env.example) for full configuration.

## Observability

Access monitoring dashboards:
- **Grafana**: http://localhost:3001 (admin/admin)
- **Prometheus**: http://localhost:9090
- **Backend Metrics**: http://localhost:3000/metrics

Custom metrics can be added using the Prometheus Go client. See [frontend/observability/README.md](frontend/observability/README.md) for details.

## Documentation

- [DEV.md](DEV.md) - Local development setup
- [DOCKER.md](DOCKER.md) - Container orchestration
- [GOOGLE_AUTH.md](GOOGLE_AUTH.md) - OAuth implementation
- [CICD.md](CICD.md) - GitHub Actions workflow
- [VALIDATION.md](VALIDATION.md) - Request validation
- [IMPROVEMENTS.md](IMPROVEMENTS.md) - Roadmap

## Security

- **CORS** - Configurable origin whitelist
- **Rate Limiting** - 100 req/min per IP
- **Security Headers** - CSP, X-Frame-Options, HSTS
- **JWT Validation** - HMAC-signed tokens
- **Input Validation** - Struct-based validation with custom rules

### Production Hardening
- Store secrets in AWS Secrets Manager
- Use httpOnly cookies for tokens
- Add refresh token rotation
- Enable HTTPS with valid certificates
- Implement CSRF protection

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.

## Support

- **Issues**: [GitHub Issues](https://github.com/TlangelaniBrian/shadow-nova/issues)
- **Documentation**: See `/docs` directory
- **Email**: your-email@example.com

---

Built with ❤️ using Vue.js, Go, and modern DevOps practices.
