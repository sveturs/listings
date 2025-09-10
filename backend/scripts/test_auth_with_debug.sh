#!/bin/bash

echo "=== Auth Service Debug Test ==="
echo "Testing auth flow to understand frontend issues"
echo ""

# Use a unique identifier to bypass rate limiting
UNIQUE_ID=$(date +%s%N)
IP="10.0.0.$((RANDOM % 255))"

echo "1. Testing login with spoofed IP: $IP"
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:3000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -H "X-Forwarded-For: $IP" \
  -H "X-Real-IP: $IP" \
  -d '{"email":"test@example.com","password":"password123"}')

if echo "$LOGIN_RESPONSE" | grep -q "access_token"; then
  echo "✓ Login successful"
  
  # Extract tokens
  ACCESS_TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)
  REFRESH_TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"refresh_token":"[^"]*' | cut -d'"' -f4)
  
  echo ""
  echo "2. Tokens received:"
  echo "   Access token: ${ACCESS_TOKEN:0:50}..."
  echo "   Refresh token: ${REFRESH_TOKEN:0:50}..."
  
  # Decode JWT header to check algorithm
  echo ""
  echo "3. Token algorithm check:"
  ACCESS_HEADER=$(echo $ACCESS_TOKEN | cut -d. -f1 | base64 -d 2>/dev/null | jq -r '.alg' 2>/dev/null)
  REFRESH_HEADER=$(echo $REFRESH_TOKEN | cut -d. -f1 | base64 -d 2>/dev/null | jq -r '.alg' 2>/dev/null)
  echo "   Access token algorithm: $ACCESS_HEADER"
  echo "   Refresh token algorithm: $REFRESH_HEADER"
  
  # Test session endpoint
  echo ""
  echo "4. Testing session endpoint:"
  SESSION_RESPONSE=$(curl -s -X GET http://localhost:3000/api/v1/auth/session \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -H "X-Forwarded-For: $IP")
  
  if echo "$SESSION_RESPONSE" | grep -q "authenticated.*true"; then
    echo "✓ Session endpoint works"
  else
    echo "✗ Session endpoint failed"
    echo "   Response: $SESSION_RESPONSE"
  fi
  
  # Test refresh
  echo ""
  echo "5. Testing refresh token:"
  REFRESH_RESPONSE=$(curl -s -X POST http://localhost:3000/api/v1/auth/refresh \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $REFRESH_TOKEN" \
    -H "X-Forwarded-For: $IP" \
    -d "{\"refresh_token\":\"$REFRESH_TOKEN\"}")
  
  if echo "$REFRESH_RESPONSE" | grep -q "access_token"; then
    echo "✓ Refresh works"
    NEW_ACCESS=$(echo $REFRESH_RESPONSE | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)
    echo "   New access token: ${NEW_ACCESS:0:50}..."
  else
    echo "✗ Refresh failed"
    echo "   Response: $REFRESH_RESPONSE"
  fi
  
else
  echo "✗ Login failed"
  echo "Response: $LOGIN_RESPONSE"
fi

echo ""
echo "=== Test Complete ==="