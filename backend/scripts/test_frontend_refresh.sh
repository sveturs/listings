#!/bin/bash

echo "=== Testing Frontend Refresh Token Flow ==="

# 1. Login to get tokens
echo -e "\n1. Login to get tokens:"
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:3000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}')

echo "Login response: $LOGIN_RESPONSE"

# Extract tokens
ACCESS_TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)
REFRESH_TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"refresh_token":"[^"]*' | cut -d'"' -f4)

echo -e "\nAccess token (first 50 chars): ${ACCESS_TOKEN:0:50}..."
echo "Refresh token (first 50 chars): ${REFRESH_TOKEN:0:50}..."

# 2. Try refresh with refresh token in Authorization header
echo -e "\n2. Trying refresh with refresh token in Authorization header:"
REFRESH_RESPONSE=$(curl -s -X POST http://localhost:3000/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $REFRESH_TOKEN" \
  -d "{\"refresh_token\":\"$REFRESH_TOKEN\"}")

echo "Refresh response: $REFRESH_RESPONSE"

# 3. Try session with new access token
NEW_ACCESS_TOKEN=$(echo $REFRESH_RESPONSE | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)

if [ -n "$NEW_ACCESS_TOKEN" ]; then
  echo -e "\n3. Testing session with new access token:"
  SESSION_RESPONSE=$(curl -s -X GET http://localhost:3000/api/v1/auth/session \
    -H "Authorization: Bearer $NEW_ACCESS_TOKEN")
  
  echo "Session response: $SESSION_RESPONSE"
else
  echo -e "\n3. No new access token received from refresh"
fi

# 4. Test what happens when we send expired/invalid access token but valid refresh
echo -e "\n4. Testing auto-refresh when access token is invalid:"
INVALID_TOKEN="eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.invalid"

# Try session with invalid access token (should trigger refresh)
SESSION_WITH_INVALID=$(curl -s -X GET http://localhost:3000/api/v1/auth/session \
  -H "Authorization: Bearer $INVALID_TOKEN")

echo "Session with invalid token: $SESSION_WITH_INVALID"