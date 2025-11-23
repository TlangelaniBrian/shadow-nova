# Shadow Nova - Docker Setup

## Quick Start

### Build and run all services:

```bash
docker-compose up --build
```

### Run in detached mode:

```bash
docker-compose up -d
```

### Stop all services:

```bash
docker-compose down
```

### View logs:

```bash
docker-compose logs -f
```

## Service URLs

- **Frontend**: http://localhost:8080
- **Backend API**: http://localhost:3000
- **Grafana**: http://localhost:3001 (admin/admin)
- **Prometheus**: http://localhost:9090
- **Unleash**: http://localhost:4242
- **PostgreSQL**: localhost:5432

## Troubleshooting

### Frontend build fails

If you see pnpm or package issues:

```bash
cd frontend
pnpm install
pnpm run build
```

### Backend build fails

```bash
cd backend
go mod tidy
go build
```

### Database connection issues

Ensure the database container is healthy:

```bash
docker-compose ps
docker-compose logs db
```

### Reset everything

```bash
docker-compose down -v  # Removes volumes
docker-compose up --build
```

## Development

### Frontend only (without Docker):

```bash
cd frontend
pnpm install
pnpm dev
```

### Backend only (without Docker):

```bash
cd backend
go run main.go
```

## Notes

- The frontend is built from the `/frontend` directory
- The backend is built from the `/backend` directory
- Both Dockerfiles are at the project root
- All config files are in their respective directories
