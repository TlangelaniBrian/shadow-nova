#!/bin/bash

# Business Metrics Verification Script
# This script verifies that all business metrics are properly exposed

set -e

BACKEND_URL="${BACKEND_URL:-http://localhost:8080}"
METRICS_ENDPOINT="${BACKEND_URL}/metrics"

echo "🔍 Verifying Business Metrics for Shadow Nova"
echo "================================================"
echo ""

# Check if backend is running
echo "1. Checking if backend is accessible..."
if ! curl -s -f "${BACKEND_URL}/health" > /dev/null 2>&1; then
    echo "❌ Backend is not accessible at ${BACKEND_URL}"
    echo "Please start the backend server first:"
    echo "  cd backend && go run main.go"
    exit 1
fi
echo "✅ Backend is accessible"
echo ""

# Check if metrics endpoint is accessible
echo "2. Checking metrics endpoint..."
if ! curl -s -f "${METRICS_ENDPOINT}" > /dev/null 2>&1; then
    echo "❌ Metrics endpoint is not accessible"
    exit 1
fi
echo "✅ Metrics endpoint is accessible"
echo ""

# Define expected metrics
EXPECTED_METRICS=(
    "user_registrations_total"
    "user_logins_total"
    "lesson_completions_total"
    "learning_path_completions_total"
    "project_submissions_total"
    "project_status_updates_total"
    "content_items_collected_total"
    "ai_processing_duration_seconds"
    "ai_processing_errors_total"
)

echo "3. Verifying business metrics..."
MISSING_METRICS=()

for metric in "${EXPECTED_METRICS[@]}"; do
    if curl -s "${METRICS_ENDPOINT}" | grep -q "^# HELP ${metric}"; then
        echo "  ✅ ${metric}"
    else
        echo "  ❌ ${metric} - NOT FOUND"
        MISSING_METRICS+=("${metric}")
    fi
done

echo ""

if [ ${#MISSING_METRICS[@]} -eq 0 ]; then
    echo "🎉 All business metrics are properly exposed!"
    echo ""
    echo "Next steps:"
    echo "  1. Import Grafana dashboard from frontend/observability/dashboards/business-metrics.json"
    echo "  2. Trigger some actions to see metrics in action:"
    echo "     - Register a user: POST /api/v1/register"
    echo "     - Login: POST /api/v1/login"
    echo "     - Complete a lesson: POST /api/v1/progress"
    echo "  3. View metrics: ${METRICS_ENDPOINT}"
    echo ""
    exit 0
else
    echo "❌ Missing metrics:"
    for metric in "${MISSING_METRICS[@]}"; do
        echo "  - ${metric}"
    done
    echo ""
    echo "Please check:"
    echo "  1. Metrics are properly initialized in backend/internal/metrics/business.go"
    echo "  2. Metrics are imported and used in handlers"
    echo "  3. Application was restarted after adding metrics"
    exit 1
fi
