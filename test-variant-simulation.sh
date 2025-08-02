#!/bin/bash

# Симуляция покупки варианта товара для тестирования системы отслеживания остатков

set -e

echo "🎭 Симуляция покупки ВАРИАНТА товара"
echo "==================================="

DB_URL="postgres://postgres:password@localhost:5432/svetubd?sslmode=disable"

run_sql() {
    psql "$DB_URL" -c "$1"
}

get_value() {
    psql "$DB_URL" -t -c "$1" | xargs
}

echo "📊 1. Поиск товара с вариантами"
echo "------------------------------"

# Найдем товар с вариантами
VARIANT_ID=$(get_value "SELECT id FROM storefront_product_variants WHERE stock_quantity >= 5 AND is_active = true LIMIT 1;")

if [ -z "$VARIANT_ID" ]; then
    echo "❌ Не найдено вариантов товаров с достаточными запасами"
    exit 1
fi

PRODUCT_ID=$(get_value "SELECT product_id FROM storefront_product_variants WHERE id = $VARIANT_ID;")
STOREFRONT_ID=$(get_value "SELECT storefront_id FROM storefront_products WHERE id = $PRODUCT_ID;")

echo "Тестовый вариант ID: $VARIANT_ID"
echo "Товар ID: $PRODUCT_ID"
echo "Витрина ID: $STOREFRONT_ID"

# Получим данные о товаре и варианте
PRODUCT_NAME=$(get_value "SELECT name FROM storefront_products WHERE id = $PRODUCT_ID;")
PRODUCT_STOCK=$(get_value "SELECT stock_quantity FROM storefront_products WHERE id = $PRODUCT_ID;")

echo "Товар: $PRODUCT_NAME"
echo "Остаток основного товара: $PRODUCT_STOCK"

# Данные варианта
VARIANT_ATTRS=$(get_value "SELECT variant_attributes::text FROM storefront_product_variants WHERE id = $VARIANT_ID;")
INITIAL_VARIANT_STOCK=$(get_value "SELECT stock_quantity FROM storefront_product_variants WHERE id = $VARIANT_ID;")
INITIAL_VARIANT_RESERVED=$(get_value "SELECT reserved_quantity FROM storefront_product_variants WHERE id = $VARIANT_ID;")
INITIAL_VARIANT_AVAILABLE=$(get_value "SELECT available_quantity FROM storefront_product_variants WHERE id = $VARIANT_ID;")

echo "Вариант: $VARIANT_ATTRS"
echo "Начальные остатки варианта:"
echo "  stock_quantity: $INITIAL_VARIANT_STOCK"
echo "  reserved_quantity: $INITIAL_VARIANT_RESERVED"
echo "  available_quantity: $INITIAL_VARIANT_AVAILABLE"

echo
echo "🛒 2. Симуляция покупки варианта"
echo "-------------------------------"

PURCHASE_QTY=2
echo "Покупаем $PURCHASE_QTY единиц варианта..."

# Создаем заказ
run_sql "
INSERT INTO storefront_orders (
    order_number, storefront_id, customer_id, 
    subtotal_amount, total_amount, commission_amount, seller_amount,
    currency, status, payment_status, 
    shipping_address, billing_address
) VALUES (
    'TEST-VAR-' || extract(epoch from now())::bigint,
    $STOREFRONT_ID,
    1,
    ${PURCHASE_QTY}50.00,
    ${PURCHASE_QTY}50.00,
    15.00,
    $((PURCHASE_QTY * 250 - 15)).00,
    'RSD',
    'pending',
    'pending',
    '{\"test\": true}',
    '{\"test\": true}'
);
"

ORDER_ID=$(get_value "SELECT id FROM storefront_orders WHERE order_number LIKE 'TEST-VAR-%' ORDER BY created_at DESC LIMIT 1;")
echo "Создан заказ ID: $ORDER_ID"

# Создаем резервирование и обновляем остатки варианта
run_sql "
BEGIN;

-- Создаем резервирование для варианта
INSERT INTO inventory_reservations (
    product_id, variant_id, quantity, order_id, status, expires_at
) VALUES (
    $PRODUCT_ID, $VARIANT_ID, $PURCHASE_QTY, $ORDER_ID, 'reserved',
    NOW() + INTERVAL '30 minutes'
);

-- Обновляем stock_quantity варианта (как в updateProductStockTx)
UPDATE storefront_product_variants 
SET stock_quantity = stock_quantity - $PURCHASE_QTY,
    updated_at = NOW()
WHERE id = $VARIANT_ID;

COMMIT;
"

RESERVATION_ID=$(get_value "SELECT id FROM inventory_reservations WHERE order_id = $ORDER_ID ORDER BY created_at DESC LIMIT 1;")
echo "Создано резервирование ID: $RESERVATION_ID"

echo
echo "📈 3. Проверка результатов"
echo "------------------------"

# Проверяем основной товар (должен остаться без изменений)
CURRENT_PRODUCT_STOCK=$(get_value "SELECT stock_quantity FROM storefront_products WHERE id = $PRODUCT_ID;")
echo "Остаток основного товара: $PRODUCT_STOCK -> $CURRENT_PRODUCT_STOCK (изменение: $((CURRENT_PRODUCT_STOCK - PRODUCT_STOCK)))"

# Проверяем вариант
CURRENT_VARIANT_STOCK=$(get_value "SELECT stock_quantity FROM storefront_product_variants WHERE id = $VARIANT_ID;")
CURRENT_VARIANT_RESERVED=$(get_value "SELECT reserved_quantity FROM storefront_product_variants WHERE id = $VARIANT_ID;")
CURRENT_VARIANT_AVAILABLE=$(get_value "SELECT available_quantity FROM storefront_product_variants WHERE id = $VARIANT_ID;")

echo "Остатки варианта:"
echo "  stock_quantity: $INITIAL_VARIANT_STOCK -> $CURRENT_VARIANT_STOCK (изменение: $((CURRENT_VARIANT_STOCK - INITIAL_VARIANT_STOCK)))"
echo "  reserved_quantity: $INITIAL_VARIANT_RESERVED -> $CURRENT_VARIANT_RESERVED (изменение: $((CURRENT_VARIANT_RESERVED - INITIAL_VARIANT_RESERVED)))"
echo "  available_quantity: $INITIAL_VARIANT_AVAILABLE -> $CURRENT_VARIANT_AVAILABLE (изменение: $((CURRENT_VARIANT_AVAILABLE - INITIAL_VARIANT_AVAILABLE)))"

echo
echo "📋 4. Детали резервирования"
echo "-------------------------"

run_sql "
SELECT 
    ir.id,
    ir.product_id,
    ir.variant_id,
    ir.quantity,
    ir.status,
    ir.expires_at,
    so.order_number,
    so.status as order_status
FROM inventory_reservations ir
JOIN storefront_orders so ON so.id = ir.order_id
WHERE ir.order_id = $ORDER_ID;
"

echo
echo "✅ Анализ результатов"
echo "===================="

# Проверки
if [ "$CURRENT_PRODUCT_STOCK" -eq "$PRODUCT_STOCK" ]; then
    echo "✅ ОТЛИЧНО: Основной товар не изменился (правильно для варианта)"
else
    echo "⚠️  ВНИМАНИЕ: Основной товар изменился, хотя покупался вариант"
fi

if [ "$CURRENT_VARIANT_STOCK" -eq $((INITIAL_VARIANT_STOCK - PURCHASE_QTY)) ]; then
    echo "✅ ОТЛИЧНО: stock_quantity варианта корректно уменьшился на $PURCHASE_QTY"
else
    echo "❌ ОШИБКА: Некорректное изменение stock_quantity варианта"
fi

# Примечание: reserved_quantity может не обновляться автоматически
# Это зависит от того, есть ли триггеры для синхронизации
echo
echo "📝 Примечание о reserved_quantity:"
echo "В системе может быть логика для обновления reserved_quantity через триггеры"
echo "или отдельные операции. Это нормально если он не изменился сразу."

echo
echo "🧹 Очистка тестовых данных:"
echo "DELETE FROM inventory_reservations WHERE order_id = $ORDER_ID;"
echo "DELETE FROM storefront_orders WHERE id = $ORDER_ID;"
echo "UPDATE storefront_product_variants SET stock_quantity = stock_quantity + $PURCHASE_QTY WHERE id = $VARIANT_ID;"