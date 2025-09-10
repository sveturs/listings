#!/bin/bash

echo "=== Testing Auth Service Refresh Token Flow ==="

# 1. First login to get tokens
echo -e "\n1. Login to get initial tokens:"
LOGIN_RESPONSE=$(curl -s -c /tmp/auth_cookies.txt -X POST http://localhost:28080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}' 2>/dev/null)

echo "Login response: $LOGIN_RESPONSE"

# 2. Check cookies
echo -e "\n2. Cookies saved:"
cat /tmp/auth_cookies.txt | grep -v "^#"

# 3. Try refresh with cookies
echo -e "\n3. Trying refresh with cookies:"
REFRESH_RESPONSE=$(curl -s -b /tmp/auth_cookies.txt -c /tmp/auth_cookies2.txt \
  -X POST http://localhost:28080/api/v1/auth/refresh \
  -H "Content-Type: application/json")

echo "Refresh response: $REFRESH_RESPONSE"

# 4. Check new cookies
echo -e "\n4. New cookies after refresh:"
cat /tmp/auth_cookies2.txt | grep -v "^#"

# 5. Try refresh through main backend proxy
echo -e "\n5. Trying refresh through main backend proxy:"
PROXY_REFRESH=$(curl -s -b /tmp/auth_cookies.txt -c /tmp/auth_cookies3.txt \
  -X POST http://localhost:3000/api/v1/auth/refresh \
  -H "Content-Type: application/json")

echo "Proxy refresh response: $PROXY_REFRESH"

# 6. Check cookies from proxy
echo -e "\n6. Cookies from proxy:"
cat /tmp/auth_cookies3.txt | grep -v "^#"

# Cleanup
rm -f /tmp/auth_cookies*.txt