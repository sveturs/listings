#!/bin/bash

# Тестирование JWT аутентификации

API_URL="http://localhost:3000"

echo "=== JWT Authentication Test ==="
echo

# 1. Получаем CSRF токен
echo "1. Getting CSRF token..."
CSRF_TOKEN=$(curl -s -c cookies.txt "$API_URL/api/v1/csrf-token" | jq -r '.csrf_token')
echo "CSRF Token: $CSRF_TOKEN"
echo

# 2. Логинимся и получаем JWT
echo "2. Login with email/password..."
LOGIN_RESPONSE=$(curl -s -b cookies.txt -c cookies.txt \
  -X POST "$API_URL/api/v1/users/login" \
  -H "Content-Type: application/json" \
  -H "X-CSRF-Token: $CSRF_TOKEN" \
  -d '{
    "email": "info@svetu.rs",
    "password": "password"
  }')

echo "Login Response:"
echo "$LOGIN_RESPONSE" | jq '.'

# Извлекаем JWT токен
JWT_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.token')
echo
echo "JWT Token (first 50 chars): ${JWT_TOKEN:0:50}..."
echo "Token length: ${#JWT_TOKEN} chars"
echo

# 3. Тестируем защищенный эндпоинт с JWT
echo "3. Testing protected endpoint with JWT..."
echo "Making request to /api/v1/marketplace/listings"

# С Bearer токеном
echo "With Bearer token:"
curl -s -X GET "$API_URL/api/v1/marketplace/listings" \
  -H "Authorization: Bearer $JWT_TOKEN" | jq -c '.success'

# С cookie (должен быть fallback)
echo
echo "With cookie only (no Bearer):"
curl -s -b cookies.txt \
  -X GET "$API_URL/api/v1/marketplace/listings" | jq -c '.success'

# 4. Проверяем, что без токена не работает
echo
echo "4. Testing without authentication (should fail):"
curl -s -X GET "$API_URL/api/v1/marketplace/listings" | jq -c '.error'

# Cleanup
rm -f cookies.txt

echo
echo "=== Test completed ===