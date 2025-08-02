#!/bin/bash

# Скрипт для тестирования создания заказов и проверки изменения остатков товаров

set -e

echo "🛒 Тестирование создания заказов и изменения остатков"
echo "==================================================="

# Настройки
API_URL="http://localhost:3000/api/v1"
DB_URL="postgres://postgres:password@localhost:5432/svetubd?sslmode=disable"

# Функции
run_sql() {
    psql "$DB_URL" -c "$1"
}

get_value() {
    psql "$DB_URL" -t -c "$1" | xargs
}

# Проверим что backend работает
echo "🔍 1. Проверка работы backend API"
echo "--------------------------------"

if curl -s "$API_URL/health" >/dev/null 2>&1; then
    echo "✅ Backend API отвечает"
else
    echo "❌ Backend API не отвечает на $API_URL"
    echo "Проверьте что backend запущен на порту 3000"
    exit 1
fi

echo
echo "📊 2. Получение тестовых данных"
echo "------------------------------"

# Найдем товар с запасами для тестирования
PRODUCT_ID=$(get_value "SELECT id FROM storefront_products WHERE stock_quantity > 5 AND is_active = true LIMIT 1;")
STOREFRONT_ID=$(get_value "SELECT storefront_id FROM storefront_products WHERE id = $PRODUCT_ID;")

if [ -z "$PRODUCT_ID" ] || [ -z "$STOREFRONT_ID" ]; then
    echo "❌ Не найдено подходящих товаров для тестирования"
    exit 1
fi

echo "Тестовый товар ID: $PRODUCT_ID"
echo "Витрина ID: $STOREFRONT_ID"

# Получим текущие остатки
INITIAL_STOCK=$(get_value "SELECT stock_quantity FROM storefront_products WHERE id = $PRODUCT_ID;")
echo "Начальный остаток: $INITIAL_STOCK"

# Проверим есть ли варианты
VARIANT_ID=$(get_value "SELECT id FROM storefront_product_variants WHERE product_id = $PRODUCT_ID AND stock_quantity > 2 LIMIT 1;")
if [ -n "$VARIANT_ID" ]; then
    echo "Найден вариант ID: $VARIANT_ID"
    INITIAL_VARIANT_STOCK=$(get_value "SELECT stock_quantity FROM storefront_product_variants WHERE id = $VARIANT_ID;")
    echo "Начальный остаток варианта: $INITIAL_VARIANT_STOCK"
fi

echo
echo "🧪 3. Попытка создания заказа без авторизации"
echo "--------------------------------------------"

# Создадим тестовый запрос заказа
if [ -n "$VARIANT_ID" ]; then
    ORDER_JSON=$(cat <<EOF
{
    "storefront_id": $STOREFRONT_ID,
    "items": [
        {
            "product_id": $PRODUCT_ID,
            "variant_id": $VARIANT_ID,
            "quantity": 2
        }
    ],
    "shipping_method": "standard",
    "payment_method": "card",
    "customer_notes": "Test order",
    "shipping_address": {
        "name": "Test User",
        "address": "Test Address 123",
        "city": "Belgrade",
        "postal_code": "11000",
        "country": "Serbia"
    },
    "billing_address": {
        "name": "Test User",
        "address": "Test Address 123",
        "city": "Belgrade",
        "postal_code": "11000",
        "country": "Serbia"
    }
}
EOF
)
else
    ORDER_JSON=$(cat <<EOF
{
    "storefront_id": $STOREFRONT_ID,
    "items": [
        {
            "product_id": $PRODUCT_ID,
            "quantity": 2
        }
    ],
    "shipping_method": "standard",
    "payment_method": "card",
    "customer_notes": "Test order",
    "shipping_address": {
        "name": "Test User",
        "address": "Test Address 123",
        "city": "Belgrade",
        "postal_code": "11000",
        "country": "Serbia"
    },
    "billing_address": {
        "name": "Test User",
        "address": "Test Address 123",
        "city": "Belgrade",
        "postal_code": "11000",
        "country": "Serbia"
    }
}
EOF
)
fi

echo "Отправка заказа..."
RESPONSE=$(curl -s -X POST "$API_URL/storefront/orders" \
    -H "Content-Type: application/json" \
    -d "$ORDER_JSON")

echo "Ответ API: $RESPONSE"

# Проверим изменились ли остатки (даже если заказ не прошел из-за авторизации)
echo
echo "📈 4. Проверка остатков после попытки создания заказа"
echo "----------------------------------------------------"

CURRENT_STOCK=$(get_value "SELECT stock_quantity FROM storefront_products WHERE id = $PRODUCT_ID;")
echo "Остаток товара: $INITIAL_STOCK -> $CURRENT_STOCK"

if [ -n "$VARIANT_ID" ]; then
    CURRENT_VARIANT_STOCK=$(get_value "SELECT stock_quantity FROM storefront_product_variants WHERE id = $VARIANT_ID;")
    echo "Остаток варианта: $INITIAL_VARIANT_STOCK -> $CURRENT_VARIANT_STOCK"
fi

# Проверим резервирования
RESERVATIONS=$(get_value "SELECT COUNT(*) FROM inventory_reservations WHERE product_id = $PRODUCT_ID AND status = 'reserved';")
echo "Активных резервирований: $RESERVATIONS"

echo
echo "📋 5. Детальная информация о резервированиях"
echo "-------------------------------------------"

if [ "$RESERVATIONS" -gt 0 ]; then
    echo "Найденные резервирования:"
    run_sql "
    SELECT 
        ir.id,
        ir.product_id,
        ir.variant_id,
        ir.quantity,
        ir.status,
        ir.expires_at,
        ir.created_at
    FROM inventory_reservations ir
    WHERE ir.product_id = $PRODUCT_ID
    ORDER BY ir.created_at DESC;
    "
else
    echo "Резервирований не найдено"
fi

echo
echo "🔍 6. Проверка заказов в системе"
echo "-------------------------------"

ORDER_COUNT=$(get_value "SELECT COUNT(*) FROM storefront_orders WHERE storefront_id = $STOREFRONT_ID;")
echo "Всего заказов в витрине: $ORDER_COUNT"

if [ "$ORDER_COUNT" -gt 0 ]; then
    echo "Последние заказы:"
    run_sql "
    SELECT 
        id,
        order_number,
        status,
        payment_status,
        total_amount,
        created_at
    FROM storefront_orders 
    WHERE storefront_id = $STOREFRONT_ID
    ORDER BY created_at DESC
    LIMIT 5;
    "
fi

echo
echo "✅ Тестирование завершено!"
echo "========================="

if [ "$CURRENT_STOCK" != "$INITIAL_STOCK" ]; then
    echo "⚠️  ВНИМАНИЕ: Остатки товара изменились!"
    if [ -n "$VARIANT_ID" ] && [ "$CURRENT_VARIANT_STOCK" != "$INITIAL_VARIANT_STOCK" ]; then
        echo "⚠️  ВНИМАНИЕ: Остатки варианта изменились!"
    fi
    echo "🎯 Система отслеживания остатков РАБОТАЕТ"
else
    echo "ℹ️  Остатки товара не изменились (возможно заказ не прошел из-за авторизации)"
fi

echo
echo "📝 Рекомендации:"
echo "1. Если заказ не прошел из-за авторизации - это нормально"
echo "2. Если остатки изменились - проверьте логику отката при неудачном заказе"
echo "3. Проверьте что резервирования истекают через 30 минут"
echo "4. Протестируйте с правильной авторизацией"