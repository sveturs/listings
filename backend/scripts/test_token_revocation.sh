#!/bin/bash

# Test token revocation functionality

echo "=== Testing Token Revocation ==="

# 1. Create a test token
echo "1. Creating test JWT token..."
TOKEN=$(cd /data/hostel-booking-system/backend && JWT_SECRET=$(grep JWT_SECRET .env | cut -d '=' -f2) go run scripts/create_test_jwt.go)
echo "Token created: ${TOKEN:0:30}..."

# 2. Test that token works
echo -e "\n2. Testing token with protected endpoint..."
RESPONSE=$(curl -s -H "Authorization: Bearer $TOKEN" http://localhost:28080/api/v1/auth/validate)
echo "Validation response: $RESPONSE"

# 3. Logout with the token
echo -e "\n3. Logging out (should revoke token)..."
LOGOUT_RESPONSE=$(curl -s -X POST -H "Authorization: Bearer $TOKEN" http://localhost:28080/api/v1/auth/logout)
echo "Logout response: $LOGOUT_RESPONSE"

# 4. Test that token is now revoked
echo -e "\n4. Testing token after logout (should be revoked)..."
RESPONSE_AFTER=$(curl -s -H "Authorization: Bearer $TOKEN" http://localhost:28080/api/v1/auth/validate)
echo "Validation response after logout: $RESPONSE_AFTER"

# 5. Check database for revoked token
echo -e "\n5. Checking database for revoked tokens..."
psql "postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5432/auth_svetu?sslmode=disable" -c \
"SELECT id, jti, user_id, reason, created_at FROM auth.revoked_access_tokens ORDER BY created_at DESC LIMIT 5;"

echo -e "\n=== Test Complete ==="