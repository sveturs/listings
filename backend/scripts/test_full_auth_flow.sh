#!/bin/bash

echo "=== Testing Full Auth Flow with Auth Service ==="

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 1. Login to get tokens
echo -e "\n1. Login to get initial tokens:"
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:3000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}')

# Check if login was successful
if echo "$LOGIN_RESPONSE" | grep -q "access_token"; then
  echo -e "${GREEN}✓ Login successful${NC}"
  
  # Extract tokens
  ACCESS_TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)
  REFRESH_TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"refresh_token":"[^"]*' | cut -d'"' -f4)
  
  echo "Access token (first 50 chars): ${ACCESS_TOKEN:0:50}..."
  echo "Refresh token (first 50 chars): ${REFRESH_TOKEN:0:50}..."
else
  echo -e "${RED}✗ Login failed${NC}"
  echo "Response: $LOGIN_RESPONSE"
  exit 1
fi

# 2. Test session with access token
echo -e "\n2. Testing session with access token:"
SESSION_RESPONSE=$(curl -s -X GET http://localhost:3000/api/v1/auth/session \
  -H "Authorization: Bearer $ACCESS_TOKEN")

if echo "$SESSION_RESPONSE" | grep -q "authenticated"; then
  echo -e "${GREEN}✓ Session endpoint works with access token${NC}"
  echo "Session: $SESSION_RESPONSE"
else
  echo -e "${RED}✗ Session endpoint failed${NC}"
  echo "Response: $SESSION_RESPONSE"
fi

# 3. Test refresh token flow
echo -e "\n3. Testing refresh token flow:"
REFRESH_RESPONSE=$(curl -s -X POST http://localhost:3000/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $REFRESH_TOKEN" \
  -d "{\"refresh_token\":\"$REFRESH_TOKEN\"}")

if echo "$REFRESH_RESPONSE" | grep -q "access_token"; then
  echo -e "${GREEN}✓ Refresh token works${NC}"
  
  NEW_ACCESS_TOKEN=$(echo $REFRESH_RESPONSE | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)
  echo "New access token (first 50 chars): ${NEW_ACCESS_TOKEN:0:50}..."
  
  # 4. Test session with new token
  echo -e "\n4. Testing session with new access token:"
  NEW_SESSION_RESPONSE=$(curl -s -X GET http://localhost:3000/api/v1/auth/session \
    -H "Authorization: Bearer $NEW_ACCESS_TOKEN")
  
  if echo "$NEW_SESSION_RESPONSE" | grep -q "authenticated"; then
    echo -e "${GREEN}✓ Session works with refreshed token${NC}"
  else
    echo -e "${RED}✗ Session failed with refreshed token${NC}"
    echo "Response: $NEW_SESSION_RESPONSE"
  fi
else
  echo -e "${RED}✗ Refresh failed${NC}"
  echo "Response: $REFRESH_RESPONSE"
fi

# 5. Test logout
echo -e "\n5. Testing logout:"
LOGOUT_RESPONSE=$(curl -s -X POST http://localhost:3000/api/v1/auth/logout \
  -H "Authorization: Bearer $NEW_ACCESS_TOKEN")

echo "Logout response: $LOGOUT_RESPONSE"

# 6. Test that session is invalid after logout
echo -e "\n6. Testing session after logout (should fail):"
INVALID_SESSION=$(curl -s -X GET http://localhost:3000/api/v1/auth/session \
  -H "Authorization: Bearer $NEW_ACCESS_TOKEN")

if echo "$INVALID_SESSION" | grep -q "authenticated\":false\|unauthorized\|401"; then
  echo -e "${GREEN}✓ Session correctly invalidated after logout${NC}"
else
  echo -e "${RED}✗ Session still valid after logout (security issue!)${NC}"
  echo "Response: $INVALID_SESSION"
fi

echo -e "\n${GREEN}=== Auth Flow Test Complete ===${NC}"