#!/bin/bash

# Тестовый скрипт для проверки CarAPI

echo "=== CarAPI Test Script ==="

# Данные из документации - нужно узнать правильные credentials
API_TOKEN="${CARAPI_TOKEN:-1a0bef70-bdd4-4b43-afea-5836c36d32b1}"
API_SECRET="${CARAPI_SECRET:-your_api_secret}"  # Нужен secret!

echo "Using API_TOKEN: $API_TOKEN"
echo "Using API_SECRET: ${API_SECRET:0:10}..."

# Шаг 1: Получаем JWT токен
echo "Getting JWT token..."
JWT_RESPONSE=$(curl -s -X 'POST' \
  'https://carapi.app/api/auth/login' \
  -H 'accept: text/plain' \
  -H 'Content-Type: application/json' \
  -d "{
    \"api_token\": \"$API_TOKEN\",
    \"api_secret\": \"$API_SECRET\"
  }")

echo "JWT Response: $JWT_RESPONSE"

# Проверяем что получили токен
if [[ $JWT_RESPONSE == *"exception"* ]]; then
    echo "ERROR: Failed to get JWT token"
    echo "Response: $JWT_RESPONSE"
    exit 1
fi

JWT_TOKEN="$JWT_RESPONSE"
echo "Got JWT token: ${JWT_TOKEN:0:50}..."

# Шаг 2: Тестируем API с JWT токеном
echo "Testing API with JWT..."

# Получаем марки
echo "Fetching makes..."
curl -X 'GET' \
  'https://carapi.app/api/makes' \
  -H 'accept: application/json' \
  -H "Authorization: Bearer $JWT_TOKEN" \
  | jq '.[0:5]'  # Показываем первые 5 марок

# Получаем модели для Volkswagen (ID обычно 184)
echo -e "\nFetching Volkswagen models..."
curl -X 'GET' \
  'https://carapi.app/api/models?make_id=184' \
  -H 'accept: application/json' \
  -H "Authorization: Bearer $JWT_TOKEN" \
  | jq '.[0:5]'  # Показываем первые 5 моделей