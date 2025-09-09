#!/bin/bash

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Тестирование админской сессии ===${NC}"
echo ""

# Генерируем тестовый JWT токен для админа
echo -e "${YELLOW}1. Генерация JWT токена для админа...${NC}"
TOKEN=$(cd /data/hostel-booking-system/backend && go run scripts/create_admin_jwt_fixed.go 2>/dev/null)
if [ -z "$TOKEN" ]; then
    echo -e "${RED}Ошибка: не удалось создать токен${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Токен создан${NC}"
echo ""

# Проверяем Auth Service напрямую
echo -e "${YELLOW}2. Проверка Auth Service (порт 28080)...${NC}"
AUTH_RESPONSE=$(curl -s http://localhost:28080/api/v1/auth/session \
  -H "Authorization: Bearer $TOKEN")
echo "$AUTH_RESPONSE" | jq '.'

IS_ADMIN_AUTH=$(echo "$AUTH_RESPONSE" | jq -r '.user.is_admin')
if [ "$IS_ADMIN_AUTH" = "true" ]; then
    echo -e "${GREEN}✓ Auth Service возвращает is_admin: true${NC}"
else
    echo -e "${RED}✗ Auth Service НЕ возвращает is_admin: true${NC}"
fi
echo ""

# Проверяем через backend прокси
echo -e "${YELLOW}3. Проверка через Backend прокси (порт 3000)...${NC}"
BACKEND_RESPONSE=$(curl -s http://localhost:3000/api/v1/auth/session \
  -H "Authorization: Bearer $TOKEN")
echo "$BACKEND_RESPONSE" | jq '.'

IS_ADMIN_BACKEND=$(echo "$BACKEND_RESPONSE" | jq -r '.user.is_admin')
if [ "$IS_ADMIN_BACKEND" = "true" ]; then
    echo -e "${GREEN}✓ Backend прокси возвращает is_admin: true${NC}"
else
    echo -e "${RED}✗ Backend прокси НЕ возвращает is_admin: true${NC}"
fi
echo ""

# Итоговый результат
echo -e "${BLUE}=== Результат ===${NC}"
if [ "$IS_ADMIN_AUTH" = "true" ] && [ "$IS_ADMIN_BACKEND" = "true" ]; then
    echo -e "${GREEN}✅ Всё работает корректно! Админский флаг передается правильно.${NC}"
else
    echo -e "${RED}❌ Есть проблемы с передачей админского флага.${NC}"
fi