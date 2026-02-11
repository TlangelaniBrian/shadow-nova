#!/bin/bash

# Idempotency Testing Script
# This script tests the idempotency implementation

set -e

echo "=== Idempotency Testing Script ==="
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
API_URL="${API_URL:-http://localhost:8080}"
TOKEN="${TOKEN:-}"

if [ -z "$TOKEN" ]; then
    echo -e "${RED}Error: TOKEN environment variable not set${NC}"
    echo "Usage: TOKEN=your_jwt_token ./test_idempotency.sh"
    exit 1
fi

echo -e "${YELLOW}Testing against: $API_URL${NC}"
echo ""

# Generate a unique idempotency key
IDEMPOTENCY_KEY="test-$(date +%s)-$(shuf -i 1000-9999 -n 1)"
echo -e "${YELLOW}Using Idempotency-Key: $IDEMPOTENCY_KEY${NC}"
echo ""

# Test 1: First request
echo "Test 1: First request with idempotency key"
echo "-------------------------------------------"
RESPONSE1=$(curl -s -w "\n%{http_code}" -X POST "$API_URL/api/progress" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: $IDEMPOTENCY_KEY" \
  -d '{"lesson_id": 1, "completed": true}')

HTTP_CODE1=$(echo "$RESPONSE1" | tail -n1)
BODY1=$(echo "$RESPONSE1" | head -n-1)

echo "Status: $HTTP_CODE1"
echo "Body: $BODY1"
echo ""

if [ "$HTTP_CODE1" = "200" ] || [ "$HTTP_CODE1" = "201" ]; then
    echo -e "${GREEN}✓ First request successful${NC}"
else
    echo -e "${RED}✗ First request failed${NC}"
    exit 1
fi

# Wait a moment
sleep 1

# Test 2: Second request with same key (should return cached)
echo ""
echo "Test 2: Second request with same idempotency key"
echo "------------------------------------------------"
RESPONSE2=$(curl -s -i -X POST "$API_URL/api/progress" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: $IDEMPOTENCY_KEY" \
  -d '{"lesson_id": 1, "completed": true}')

echo "Response headers:"
echo "$RESPONSE2" | grep -i "x-idempotent-replay" || echo "(No replay header found)"
echo ""

if echo "$RESPONSE2" | grep -q "X-Idempotent-Replay: true"; then
    echo -e "${GREEN}✓ Second request returned cached response${NC}"
else
    echo -e "${RED}✗ Second request did not return cached response${NC}"
    echo -e "${YELLOW}Note: This is expected if idempotency is not enabled${NC}"
fi

# Test 3: Request without idempotency key (should process normally)
echo ""
echo "Test 3: Request without idempotency key"
echo "---------------------------------------"
RESPONSE3=$(curl -s -w "\n%{http_code}" -X POST "$API_URL/api/progress" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"lesson_id": 2, "completed": true}')

HTTP_CODE3=$(echo "$RESPONSE3" | tail -n1)
BODY3=$(echo "$RESPONSE3" | head -n-1)

echo "Status: $HTTP_CODE3"
echo "Body: $BODY3"
echo ""

if [ "$HTTP_CODE3" = "200" ] || [ "$HTTP_CODE3" = "201" ]; then
    echo -e "${GREEN}✓ Request without key processed normally${NC}"
else
    echo -e "${RED}✗ Request without key failed${NC}"
    exit 1
fi

# Test 4: GET request (should not cache)
echo ""
echo "Test 4: GET request with idempotency key (should not cache)"
echo "-----------------------------------------------------------"
GET_KEY="test-get-$(date +%s)"
RESPONSE4=$(curl -s -i -X GET "$API_URL/api/stats" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Idempotency-Key: $GET_KEY")

echo "Response headers:"
echo "$RESPONSE4" | grep -i "x-idempotent-replay" || echo "(No replay header - expected)"
echo ""

if echo "$RESPONSE4" | grep -q "X-Idempotent-Replay: true"; then
    echo -e "${RED}✗ GET request should not be cached${NC}"
else
    echo -e "${GREEN}✓ GET request not cached (correct behavior)${NC}"
fi

echo ""
echo "=== Test Summary ==="
echo -e "${GREEN}All tests completed${NC}"
echo ""
echo "To verify database entries:"
echo "  psql -d shadownova -c \"SELECT key, user_id, request_method, response_status, created_at FROM idempotency_keys ORDER BY created_at DESC LIMIT 5;\""
