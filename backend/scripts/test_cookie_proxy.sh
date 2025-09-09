#!/bin/bash

# Тест проксирования cookies через auth proxy

echo "=== Тест проксирования Set-Cookie заголовков ==="
echo ""

# URL сервера
BASE_URL="http://localhost:3000"
AUTH_SERVICE_URL="http://localhost:28080"

# Временный файл для сохранения полного ответа
RESPONSE_FILE="/tmp/auth_proxy_response.txt"
HEADERS_FILE="/tmp/auth_proxy_headers.txt"

# Проверяем доступность Auth Service напрямую
echo "1. Проверяем Auth Service напрямую на порту 28080..."
direct_response=$(curl -s -o /dev/null -w "%{http_code}" "$AUTH_SERVICE_URL/api/v1/auth/health" 2>/dev/null)
if [ "$direct_response" = "200" ]; then
    echo "   ✓ Auth Service доступен напрямую на порту 28080"
else
    echo "   ✗ Auth Service не отвечает на порту 28080 (код: $direct_response)"
    echo "   Возможно, Auth Service запущен в Docker контейнере"
fi

echo ""
echo "2. Проверяем проксирование через backend..."
# Делаем запрос через прокси и сохраняем все заголовки
curl -v -s "$BASE_URL/api/v1/auth/google" > $RESPONSE_FILE 2>&1

# Извлекаем заголовки
grep -E "^< " $RESPONSE_FILE > $HEADERS_FILE

echo "   Заголовки ответа:"
echo "   ----------------"
cat $HEADERS_FILE | sed 's/^< /   /'

echo ""
echo "3. Проверяем наличие Set-Cookie заголовков..."
cookie_count=$(grep -c "Set-Cookie:" $HEADERS_FILE 2>/dev/null || echo "0")
if [ "$cookie_count" -gt "0" ]; then
    echo "   ✓ Найдено $cookie_count Set-Cookie заголовков:"
    grep "Set-Cookie:" $HEADERS_FILE | sed 's/^< /   /'
else
    echo "   ⚠ Set-Cookie заголовки не найдены"
    echo "   Это нормально для OAuth редиректа, cookies устанавливаются после callback"
fi

echo ""
echo "4. Проверяем Location заголовок..."
location=$(grep "Location:" $HEADERS_FILE | cut -d' ' -f2- | tr -d '\r')
if [ ! -z "$location" ]; then
    echo "   ✓ Редирект на: $location"
    if [[ $location == *"accounts.google.com"* ]]; then
        echo "   ✓ Корректный редирект на Google OAuth"
    fi
else
    echo "   ✗ Location заголовок не найден"
fi

echo ""
echo "5. Тестируем эмуляцию callback с токенами..."
echo "   Создаём тестовый запрос к callback эндпоинту..."

# Эмулируем callback от Google с фейковым кодом
callback_response=$(curl -v -s -L "$BASE_URL/api/v1/auth/google/callback?code=fake_auth_code&state=fake_state" 2>&1)

# Проверяем Set-Cookie в callback ответе
echo "$callback_response" | grep -E "Set-Cookie:|< set-cookie:" > $HEADERS_FILE
callback_cookies=$(grep -c -i "set-cookie" $HEADERS_FILE 2>/dev/null || echo "0")

if [ "$callback_cookies" -gt "0" ]; then
    echo "   ✓ Callback установил $callback_cookies cookies:"
    grep -i "set-cookie" $HEADERS_FILE | sed 's/^[<]*/   /'
else
    echo "   ⚠ Callback не установил cookies (ожидаемо для фейкового кода)"
fi

echo ""
echo "6. Проверяем validate и refresh эндпоинты..."

# Проверяем validate
validate_response=$(curl -s -X POST "$BASE_URL/api/v1/auth/validate" \
    -H "Content-Type: application/json" \
    -d '{"token":"test_token"}' \
    -w "\nHTTP_STATUS:%{http_code}")
    
validate_status=$(echo "$validate_response" | grep "HTTP_STATUS:" | cut -d: -f2)
echo "   Validate эндпоинт вернул код: $validate_status"

# Проверяем refresh
refresh_response=$(curl -s -X POST "$BASE_URL/api/v1/auth/refresh" \
    -H "Cookie: refresh_token=test_refresh" \
    -w "\nHTTP_STATUS:%{http_code}")
    
refresh_status=$(echo "$refresh_response" | grep "HTTP_STATUS:" | cut -d: -f2)
echo "   Refresh эндпоинт вернул код: $refresh_status"

echo ""
echo "=== Итоги тестирования проксирования ==="
echo ""
if [ "$cookie_count" -eq "0" ] && [ "$callback_cookies" -eq "0" ]; then
    echo "⚠ Cookies не проксируются через auth proxy"
    echo ""
    echo "Возможные причины:"
    echo "1. Auth Service не устанавливает cookies на этих эндпоинтах"
    echo "2. Проблема с проксированием Set-Cookie заголовков"
    echo "3. Auth Service использует другой механизм для установки cookies"
else
    echo "✓ Проксирование работает корректно"
fi

echo ""
echo "Для полной проверки:"
echo "1. Откройте браузер в инкогнито режиме"
echo "2. Откройте DevTools -> Network"
echo "3. Перейдите на $BASE_URL/api/v1/auth/google"
echo "4. После авторизации проверьте Set-Cookie заголовки в callback ответе"

# Очистка
rm -f $RESPONSE_FILE $HEADERS_FILE