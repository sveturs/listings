#!/bin/bash

echo "=== Testing Full Logout Flow ==="

# 1. Login
echo "1. Login..."
curl -s -c cookies.txt -X POST http://localhost:3001/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"info@svetu.rs","password":"info@svetu.rs"}' | jq '.'

echo -e "\n2. Check cookies after login:"
cat cookies.txt

# 3. Check session
echo -e "\n3. Check session..."
curl -s -b cookies.txt http://localhost:3000/auth/session | jq '.'

# 4. Logout
echo -e "\n4. Logout..."
curl -s -b cookies.txt -c cookies.txt -X POST http://localhost:3001/api/auth/logout | jq '.'

echo -e "\n5. Check cookies after logout:"
cat cookies.txt

# 6. Try to check session after logout
echo -e "\n6. Check session after logout (should be unauthenticated)..."
curl -s -b cookies.txt http://localhost:3000/auth/session | jq '.'

# Cleanup
rm -f cookies.txt