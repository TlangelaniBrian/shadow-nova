#!/bin/bash

# Test script for User Profile Management endpoints
# Requires a running backend server and a valid auth token

API_URL="${API_URL:-http://localhost:8080/api/v1}"
AUTH_TOKEN="${AUTH_TOKEN:-}"

echo "=========================================="
echo "User Profile Management Test"
echo "=========================================="
echo ""

if [ -z "$AUTH_TOKEN" ]; then
    echo "Error: AUTH_TOKEN environment variable is not set"
    echo "Usage: AUTH_TOKEN=your_token_here $0"
    exit 1
fi

echo "1. Testing GET /user/profile"
echo "----------------------------"
curl -s -X GET "$API_URL/user/profile" \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -w "\nStatus: %{http_code}\n\n" | jq '.' 2>/dev/null || echo "Response received"

echo ""
echo "2. Testing PATCH /user/profile (update username)"
echo "------------------------------------------------"
curl -s -X PATCH "$API_URL/user/profile" \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"username": "updated_user"}' \
  -w "\nStatus: %{http_code}\n\n" | jq '.' 2>/dev/null || echo "Response received"

echo ""
echo "3. Testing PATCH /user/profile (update email)"
echo "---------------------------------------------"
curl -s -X PATCH "$API_URL/user/profile" \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"email": "newemail@example.com"}' \
  -w "\nStatus: %{http_code}\n\n" | jq '.' 2>/dev/null || echo "Response received"

echo ""
echo "4. Testing PUT /user/password"
echo "-----------------------------"
echo "Note: This will fail with 401 if current_password is incorrect"
curl -s -X PUT "$API_URL/user/password" \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"current_password": "wrongpassword", "new_password": "NewPassword123!"}' \
  -w "\nStatus: %{http_code}\n\n" | jq '.' 2>/dev/null || echo "Response received"

echo ""
echo "=========================================="
echo "Tests completed!"
echo "=========================================="
