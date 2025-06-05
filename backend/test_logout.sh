#!/bin/bash

echo "=== Testing Logout Flow ==="

# 1. Login и получение токенов
echo "1. Logging in..."
LOGIN_RESPONSE=$(curl -s -c cookies.txt -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"info@svetu.rs","password":"svetu2024!"}')

echo "Login response: $LOGIN_RESPONSE"

# Извлекаем access_token
ACCESS_TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.access_token')
echo "Access token: ${ACCESS_TOKEN:0:20}..."

# Показываем cookies
echo -e "\n2. Cookies after login:"
cat cookies.txt | grep refresh_token

# 3. Проверяем refresh
echo -e "\n3. Testing refresh..."
REFRESH_RESPONSE=$(curl -s -b cookies.txt -c cookies.txt -X POST http://localhost:3000/api/auth/refresh)
echo "Refresh response: $REFRESH_RESPONSE"

# 4. Logout
echo -e "\n4. Logging out..."
LOGOUT_RESPONSE=$(curl -s -b cookies.txt -X POST http://localhost:3000/api/auth/logout)
echo "Logout response status: $?"

# 5. Пробуем refresh после logout
echo -e "\n5. Testing refresh after logout (should fail)..."
REFRESH_AFTER_LOGOUT=$(curl -s -b cookies.txt -X POST http://localhost:3000/api/auth/refresh)
echo "Refresh after logout: $REFRESH_AFTER_LOGOUT"

# Cleanup
rm -f cookies.txt

echo -e "\n=== Test completed ==="