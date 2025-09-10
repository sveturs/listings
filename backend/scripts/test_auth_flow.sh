#!/bin/bash

# Тестирование полного флоу аутентификации через Auth Service

echo "=== Тест аутентификации через Auth Service ==="
echo ""

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# URL сервера
BASE_URL="http://localhost:3000"
FRONTEND_URL="http://localhost:3001"

# Временный файл для cookies
COOKIE_JAR="/tmp/auth_test_cookies.txt"
rm -f $COOKIE_JAR

echo "1. Проверка доступности Auth Service..."
response=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/api/v1/auth/health")
if [ "$response" = "200" ]; then
    echo -e "${GREEN}✓ Auth Service доступен${NC}"
else
    echo -e "${RED}✗ Auth Service недоступен (HTTP $response)${NC}"
fi

echo ""
echo "2. Проверка публичных эндпоинтов без аутентификации..."

# Проверка storefronts
response=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/api/v1/storefronts?limit=10")
if [ "$response" = "200" ]; then
    echo -e "${GREEN}✓ /api/v1/storefronts доступен без токена${NC}"
else
    echo -e "${RED}✗ /api/v1/storefronts требует токен (HTTP $response)${NC}"
fi

# Проверка marketplace/search
response=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/api/v1/marketplace/search?q=test")
if [ "$response" = "200" ]; then
    echo -e "${GREEN}✓ /api/v1/marketplace/search доступен без токена${NC}"
else
    echo -e "${RED}✗ /api/v1/marketplace/search требует токен (HTTP $response)${NC}"
fi

echo ""
echo "3. Проверка защищённых эндпоинтов..."

# Проверка user/profile (должен требовать токен)
response=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/api/v1/user/profile")
if [ "$response" = "401" ]; then
    echo -e "${GREEN}✓ /api/v1/user/profile правильно требует токен${NC}"
else
    echo -e "${YELLOW}⚠ /api/v1/user/profile вернул неожиданный код: $response${NC}"
fi

echo ""
echo "4. Проверка OAuth редиректа..."
oauth_response=$(curl -s -I "$BASE_URL/api/v1/auth/google" | head -n 1)
if [[ $oauth_response == *"302"* ]]; then
    location=$(curl -s -I "$BASE_URL/api/v1/auth/google" | grep -i "location:" | cut -d' ' -f2)
    if [[ $location == *"accounts.google.com"* ]]; then
        echo -e "${GREEN}✓ OAuth редирект работает корректно${NC}"
        echo -e "   Редирект на: ${YELLOW}Google OAuth${NC}"
    else
        echo -e "${RED}✗ Неправильный OAuth редирект${NC}"
    fi
else
    echo -e "${RED}✗ OAuth редирект не работает${NC}"
fi

echo ""
echo "5. Проверка callback эндпоинта..."
# Проверяем что callback эндпоинт существует
response=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/api/v1/auth/google/callback?code=test&state=test")
if [ "$response" = "400" ] || [ "$response" = "401" ]; then
    echo -e "${GREEN}✓ OAuth callback эндпоинт существует${NC}"
else
    echo -e "${YELLOW}⚠ OAuth callback вернул код: $response${NC}"
fi

echo ""
echo "6. Проверка сохранения cookies..."
# Делаем запрос с сохранением cookies
curl -s -c $COOKIE_JAR -I "$BASE_URL/api/v1/auth/google" > /dev/null 2>&1
if [ -f $COOKIE_JAR ] && [ -s $COOKIE_JAR ]; then
    cookie_count=$(grep -c "Set-Cookie" $COOKIE_JAR 2>/dev/null || echo "0")
    if [ "$cookie_count" -gt 0 ]; then
        echo -e "${GREEN}✓ Cookies сохраняются ($cookie_count cookies)${NC}"
    else
        echo -e "${YELLOW}⚠ Cookies создаются, но пустые${NC}"
    fi
else
    echo -e "${YELLOW}⚠ Cookies не сохраняются${NC}"
fi

echo ""
echo "7. Проверка refresh эндпоинта..."
response=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/api/v1/auth/refresh")
if [ "$response" = "401" ] || [ "$response" = "400" ]; then
    echo -e "${GREEN}✓ Refresh эндпоинт существует${NC}"
else
    echo -e "${YELLOW}⚠ Refresh эндпоинт вернул код: $response${NC}"
fi

echo ""
echo "8. Проверка validate эндпоинта..."
response=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/api/v1/auth/validate" -H "Content-Type: application/json" -d '{"token":"test"}')
if [ "$response" = "401" ] || [ "$response" = "400" ]; then
    echo -e "${GREEN}✓ Validate эндпоинт существует${NC}"
else
    echo -e "${YELLOW}⚠ Validate эндпоинт вернул код: $response${NC}"
fi

echo ""
echo "=== Итоги тестирования ==="
echo ""
echo "Для полного тестирования аутентификации:"
echo "1. Откройте браузер в режиме инкогнито"
echo "2. Перейдите на $BASE_URL/api/v1/auth/google"
echo "3. Авторизуйтесь через Google"
echo "4. Проверьте что вы перенаправлены на $FRONTEND_URL"
echo "5. Откройте DevTools -> Application -> Cookies"
echo "6. Должны быть установлены cookies:"
echo "   - jwt_token (access token)"
echo "   - refresh_token (для обновления)"
echo ""

# Очистка
rm -f $COOKIE_JAR