#!/usr/bin/env bash
set -Eeuo pipefail

# --------- Defaults / Options ----------
NO_DOCKER=0
DOWN_ON_EXIT=0
SHOW_LOGS=1

usage() {
  cat <<EOF
Usage: ./dev.sh [options]

Options:
  --no-docker     Don't start docker services
  --down          Run 'compose down' on exit (instead of leaving containers running)
  --no-logs       Don't prefix/stream logs (processes still run)
  -h, --help      Show help
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --no-docker) NO_DOCKER=1; shift ;;
    --down) DOWN_ON_EXIT=1; shift ;;
    --no-logs) SHOW_LOGS=0; shift ;;
    -h|--help) usage; exit 0 ;;
    *) echo "Unknown option: $1"; usage; exit 1 ;;
  esac
done
# --------------------------------------

# Resolve project root relative to this script
SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$SCRIPT_DIR"  # put dev.sh at repo root

BACKEND_DIR="$PROJECT_ROOT/backend"
FRONTEND_DIR="$PROJECT_ROOT/frontend"

# Compose command detection: prefer docker compose (v2), fallback to docker-compose (v1)
if command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1; then
  COMPOSE=(docker compose)
elif command -v docker-compose >/dev/null 2>&1; then
  COMPOSE=(docker-compose)
else
  COMPOSE=()
fi

need_cmd() {
  command -v "$1" >/dev/null 2>&1 || { echo "Missing dependency: $1"; exit 1; }
}

echo "Starting Shadow Nova Development Environment..."
need_cmd go
need_cmd pnpm

if [[ $NO_DOCKER -eq 0 ]]; then
  need_cmd docker
  if [[ ${#COMPOSE[@]} -eq 0 ]]; then
    echo "Could not find docker compose (docker compose / docker-compose)."
    exit 1
  fi
fi

# Validate directories
[[ -d "$BACKEND_DIR" ]] || { echo "Backend dir not found: $BACKEND_DIR"; exit 1; }
[[ -d "$FRONTEND_DIR" ]] || { echo "Frontend dir not found: $FRONTEND_DIR"; exit 1; }

BACKEND_PID=""
FRONTEND_PID=""

cleanup() {
  echo ""
  echo "Stopping services..."

  # Stop child processes
  [[ -n "${BACKEND_PID}" ]] && kill "${BACKEND_PID}" 2>/dev/null || true
  [[ -n "${FRONTEND_PID}" ]] && kill "${FRONTEND_PID}" 2>/dev/null || true

  # Optionally stop docker services
  if [[ $NO_DOCKER -eq 0 ]]; then
    if [[ $DOWN_ON_EXIT -eq 1 ]]; then
      "${COMPOSE[@]}" down || true
    fi
  fi

  echo "All services stopped."
}

trap cleanup INT TERM EXIT

if [[ $NO_DOCKER -eq 0 ]]; then
  echo "Starting Docker containers (DB, Unleash)..."
  "${COMPOSE[@]}" up -d db unleash unleash-db
fi

prefix_logs() {
  local prefix="$1"
  # Read lines and prefix them
  while IFS= read -r line; do
    printf "[%s] %s\n" "$prefix" "$line"
  done
}

echo "Starting backend..."
if [[ $SHOW_LOGS -eq 1 ]]; then
  ( cd "$BACKEND_DIR" && go run main.go ) 2>&1 | prefix_logs "backend" &
  BACKEND_PID=$!
else
  ( cd "$BACKEND_DIR" && go run main.go ) >/dev/null 2>&1 &
  BACKEND_PID=$!
fi
echo "Backend started (PID: $BACKEND_PID)"

echo "Starting frontend..."
if [[ $SHOW_LOGS -eq 1 ]]; then
  ( cd "$FRONTEND_DIR" && pnpm dev ) 2>&1 | prefix_logs "frontend" &
  FRONTEND_PID=$!
else
  ( cd "$FRONTEND_DIR" && pnpm dev ) >/dev/null 2>&1 &
  FRONTEND_PID=$!
fi
echo "Frontend started (PID: $FRONTEND_PID)"

echo ""
echo "Press Ctrl+C to stop all services"
echo ""

# Wait for either to exit; if one dies, the script exits (and cleanup runs)
wait -n "$BACKEND_PID" "$FRONTEND_PID"
echo "One of the services exited. Shutting down..."
exit 1
