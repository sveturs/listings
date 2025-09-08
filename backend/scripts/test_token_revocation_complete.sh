#!/bin/bash

# Complete test of token revocation with actual login

echo "=== Complete Token Revocation Test ==="

# 1. Register/Login to get a real token
echo "1. Logging in to get a real token..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:28080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}')

# Extract access token
ACCESS_TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)

if [ -z "$ACCESS_TOKEN" ]; then
  echo "Failed to get access token. Response: $LOGIN_RESPONSE"
  echo "Trying to register first..."
  
  # Try to register first
  REG_RESPONSE=$(curl -s -X POST http://localhost:28080/api/v1/auth/register \
    -H "Content-Type: application/json" \
    -d '{"email":"test@example.com","password":"password123","name":"Test User","terms_accepted":true}')
  
  echo "Registration response: $REG_RESPONSE"
  
  # Try login again
  LOGIN_RESPONSE=$(curl -s -X POST http://localhost:28080/api/v1/auth/login \
    -H "Content-Type: application/json" \
    -d '{"email":"test@example.com","password":"password123"}')
  
  ACCESS_TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)
fi

if [ -z "$ACCESS_TOKEN" ]; then
  echo "Still no token. Login response: $LOGIN_RESPONSE"
  exit 1
fi

echo "Got access token: ${ACCESS_TOKEN:0:30}..."

# 2. Validate the token works
echo -e "\n2. Validating token before logout..."
VALIDATE_RESPONSE=$(curl -s -H "Authorization: Bearer $ACCESS_TOKEN" http://localhost:28080/api/v1/auth/validate)
echo "Validation response: $VALIDATE_RESPONSE"

# Extract JTI from validation response for later checking
JTI=$(echo $VALIDATE_RESPONSE | grep -o '"jti":"[^"]*' | cut -d'"' -f4)
echo "Token JTI: $JTI"

# 3. Logout with the token (should revoke it)
echo -e "\n3. Logging out (should revoke token)..."
LOGOUT_RESPONSE=$(curl -s -X POST \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  http://localhost:28080/api/v1/auth/logout)
echo "Logout response: $LOGOUT_RESPONSE"

# 4. Check auth service logs for debug output
echo -e "\n4. Checking auth service logs for revocation..."
cd /data/auth_svetu && docker-compose logs --tail=20 auth-service | grep -E "\[Logout Handler\]|\[Auth Service\]"

# 5. Validate the token after logout (should be invalid)
echo -e "\n5. Validating token after logout (should be revoked)..."
VALIDATE_AFTER=$(curl -s -H "Authorization: Bearer $ACCESS_TOKEN" http://localhost:28080/api/v1/auth/validate)
echo "Validation response after logout: $VALIDATE_AFTER"

# 6. Check the database for revoked tokens
echo -e "\n6. Checking database for revoked token..."
psql "postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5432/auth_db?sslmode=disable" -c \
"SELECT id, jti, user_id, reason, expires_at, created_at FROM auth.revoked_access_tokens WHERE jti = '$JTI' OR created_at > NOW() - INTERVAL '5 minutes' ORDER BY created_at DESC LIMIT 5;"

echo -e "\n=== Test Complete ===" 