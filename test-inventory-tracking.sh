#!/bin/bash

# Скрипт для тестирования отслеживания остатков товаров после покупки
# Использует прямые SQL запросы для проверки состояния инвентаря

set -e

echo "🧪 Тестирование системы отслеживания остатков товаров"
echo "=================================================="

# Настройки подключения к БД
DB_URL="postgres://postgres:password@localhost:5432/svetubd?sslmode=disable"

# Функция для выполнения SQL запросов
run_sql() {
    psql "$DB_URL" -c "$1"
}

# Функция для получения одного значения
get_value() {
    psql "$DB_URL" -t -c "$1" | xargs
}

echo
echo "📊 1. Проверка текущего состояния товаров и вариантов"
echo "---------------------------------------------------"

# Найдем товары с запасами для тестирования
echo "Товары с запасами > 0:"
run_sql "
SELECT 
    sp.id,
    sp.name,
    sp.stock_quantity,
    sp.price,
    s.name as storefront_name
FROM storefront_products sp
JOIN storefronts s ON s.id = sp.storefront_id
WHERE sp.stock_quantity > 0 AND sp.is_active = true
ORDER BY sp.stock_quantity DESC
LIMIT 5;
"

echo
echo "Варианты товаров с запасами > 0:"
run_sql "
SELECT 
    spv.id,
    sp.name as product_name,
    spv.variant_attributes,
    spv.stock_quantity,
    spv.available_quantity,
    spv.reserved_quantity,
    spv.price
FROM storefront_product_variants spv
JOIN storefront_products sp ON sp.id = spv.product_id
WHERE spv.stock_quantity > 0 AND spv.is_active = true
ORDER BY spv.stock_quantity DESC
LIMIT 5;
"

echo
echo "📋 2. Проверка активных резервирований"
echo "------------------------------------"

run_sql "
SELECT 
    ir.id,
    ir.product_id,
    ir.variant_id,
    ir.quantity,
    ir.status,
    ir.expires_at,
    so.status as order_status
FROM inventory_reservations ir
JOIN storefront_orders so ON so.id = ir.order_id
WHERE ir.status = 'reserved'
ORDER BY ir.created_at DESC
LIMIT 10;
"

echo
echo "💰 3. Проверка заказов витрин"
echo "----------------------------"

run_sql "
SELECT 
    so.id as order_id,
    so.order_number,
    so.status,
    so.payment_status,
    so.total_amount,
    so.currency,
    so.created_at
FROM storefront_orders so
ORDER BY so.created_at DESC
LIMIT 10;
"

echo
echo "🔍 4. Анализ консистентности данных"
echo "----------------------------------"

echo "Проверка товаров с отрицательными остатками:"
NEGATIVE_PRODUCTS=$(get_value "SELECT COUNT(*) FROM storefront_products WHERE stock_quantity < 0;")
if [ "$NEGATIVE_PRODUCTS" -gt 0 ]; then
    echo "⚠️  НАЙДЕНЫ товары с отрицательными остатками: $NEGATIVE_PRODUCTS"
    run_sql "SELECT id, name, stock_quantity FROM storefront_products WHERE stock_quantity < 0;"
else
    echo "✅ Нет товаров с отрицательными остатками"
fi

echo
echo "Проверка вариантов с отрицательными остатками:"
NEGATIVE_VARIANTS=$(get_value "SELECT COUNT(*) FROM storefront_product_variants WHERE stock_quantity < 0;")
if [ "$NEGATIVE_VARIANTS" -gt 0 ]; then
    echo "⚠️  НАЙДЕНЫ варианты с отрицательными остатками: $NEGATIVE_VARIANTS"
    run_sql "SELECT id, product_id, name, stock_quantity FROM storefront_product_variants WHERE stock_quantity < 0;"
else
    echo "✅ Нет вариантов с отрицательными остатками"
fi

echo
echo "🕐 5. Проверка истекших резервирований"
echo "------------------------------------"

EXPIRED_RESERVATIONS=$(get_value "SELECT COUNT(*) FROM inventory_reservations WHERE status = 'reserved' AND expires_at < NOW();")
if [ "$EXPIRED_RESERVATIONS" -gt 0 ]; then
    echo "⚠️  НАЙДЕНЫ истекшие резервирования: $EXPIRED_RESERVATIONS"
    run_sql "
    SELECT 
        ir.id,
        ir.product_id,
        ir.variant_id,
        ir.quantity,
        ir.expires_at,
        (NOW() - ir.expires_at) as expired_ago
    FROM inventory_reservations ir
    WHERE ir.status = 'reserved' AND ir.expires_at < NOW()
    ORDER BY ir.expires_at;
    "
else
    echo "✅ Нет истекших резервирований"
fi

echo
echo "📈 6. Статистика по резервированиям"
echo "----------------------------------"

run_sql "
SELECT 
    status,
    COUNT(*) as count,
    SUM(quantity) as total_quantity
FROM inventory_reservations
GROUP BY status
ORDER BY count DESC;
"

echo
echo "🎯 7. Детальный анализ конкретного товара (если есть)"
echo "--------------------------------------------------"

# Выберем первый активный товар для детального анализа
SAMPLE_PRODUCT_ID=$(get_value "SELECT id FROM storefront_products WHERE stock_quantity > 0 AND is_active = true LIMIT 1;")

if [ -n "$SAMPLE_PRODUCT_ID" ] && [ "$SAMPLE_PRODUCT_ID" != "" ]; then
    echo "Анализ товара ID: $SAMPLE_PRODUCT_ID"
    
    echo "Основная информация:"
    run_sql "
    SELECT 
        id,
        name,
        stock_quantity,
        price,
        created_at,
        updated_at
    FROM storefront_products 
    WHERE id = $SAMPLE_PRODUCT_ID;
    "
    
    echo "Варианты товара:"
    run_sql "
    SELECT 
        id,
        variant_attributes,
        stock_quantity,
        available_quantity,
        reserved_quantity,
        price,
        is_active
    FROM storefront_product_variants 
    WHERE product_id = $SAMPLE_PRODUCT_ID;
    "
    
    echo "Резервирования для этого товара:"
    run_sql "
    SELECT 
        ir.id,
        ir.variant_id,
        ir.quantity,
        ir.status,
        ir.expires_at,
        so.status as order_status
    FROM inventory_reservations ir
    JOIN storefront_orders so ON so.id = ir.order_id
    WHERE ir.product_id = $SAMPLE_PRODUCT_ID
    ORDER BY ir.created_at DESC;
    "
    
    echo "Связанные заказы (через резервирования):"
    run_sql "
    SELECT 
        ir.order_id,
        ir.variant_id,
        ir.quantity,
        ir.status as reservation_status,
        so.status as order_status,
        so.created_at
    FROM inventory_reservations ir
    JOIN storefront_orders so ON so.id = ir.order_id
    WHERE ir.product_id = $SAMPLE_PRODUCT_ID
    ORDER BY so.created_at DESC
    LIMIT 10;
    "
else
    echo "❌ Не найдено активных товаров для анализа"
fi

echo
echo "✅ Тестирование завершено!"
echo "========================="
echo
echo "📝 Следующие шаги:"
echo "1. Создать заказ через API и проверить изменение остатков"
echo "2. Протестировать поведение при недостатке товара"
echo "3. Проверить восстановление остатков при истечении резервирования"
echo "4. Тестировать конкурентные покупки"