#!/bin/bash

# Симуляция покупки товара для тестирования системы отслеживания остатков

set -e

echo "🎭 Симуляция покупки товара и проверка остатков"
echo "=============================================="

DB_URL="postgres://postgres:password@localhost:5432/svetubd?sslmode=disable"

run_sql() {
    psql "$DB_URL" -c "$1"
}

get_value() {
    psql "$DB_URL" -t -c "$1" | xargs
}

echo "📊 1. Подготовка тестовых данных"
echo "-------------------------------"

# Найдем товар с достаточными запасами
PRODUCT_ID=$(get_value "SELECT id FROM storefront_products WHERE stock_quantity >= 10 AND is_active = true LIMIT 1;")
STOREFRONT_ID=$(get_value "SELECT storefront_id FROM storefront_products WHERE id = $PRODUCT_ID;")

if [ -z "$PRODUCT_ID" ]; then
    echo "❌ Не найдено товаров с достаточными запасами"
    exit 1
fi

echo "Тестовый товар ID: $PRODUCT_ID"
echo "Витрина ID: $STOREFRONT_ID"

# Получим начальные данные
INITIAL_STOCK=$(get_value "SELECT stock_quantity FROM storefront_products WHERE id = $PRODUCT_ID;")
PRODUCT_NAME=$(get_value "SELECT name FROM storefront_products WHERE id = $PRODUCT_ID;")

echo "Товар: $PRODUCT_NAME"
echo "Начальный остаток: $INITIAL_STOCK"

# Проверим есть ли варианты
VARIANT_ID=$(get_value "SELECT id FROM storefront_product_variants WHERE product_id = $PRODUCT_ID AND stock_quantity >= 5 LIMIT 1;")
if [ -n "$VARIANT_ID" ]; then
    INITIAL_VARIANT_STOCK=$(get_value "SELECT stock_quantity FROM storefront_product_variants WHERE id = $VARIANT_ID;")
    INITIAL_VARIANT_RESERVED=$(get_value "SELECT reserved_quantity FROM storefront_product_variants WHERE id = $VARIANT_ID;")
    INITIAL_VARIANT_AVAILABLE=$(get_value "SELECT available_quantity FROM storefront_product_variants WHERE id = $VARIANT_ID;")
    
    echo "Найден вариант ID: $VARIANT_ID"
    echo "Запасы варианта: stock=$INITIAL_VARIANT_STOCK, reserved=$INITIAL_VARIANT_RESERVED, available=$INITIAL_VARIANT_AVAILABLE"
fi

echo
echo "🛒 2. Симуляция покупки (имитация CreateOrderWithTx)"
echo "---------------------------------------------------"

# Количество для "покупки"
PURCHASE_QTY=3

echo "Покупаем $PURCHASE_QTY единиц товара..."

# В транзакции делаем то же что делает CreateOrderWithTx:
# 1. Создаем фиктивный заказ
# 2. Создаем резервирование
# 3. Обновляем stock_quantity

run_sql "
BEGIN;

-- Создаем фиктивный заказ для тестирования
INSERT INTO storefront_orders (
    order_number, storefront_id, customer_id, 
    subtotal_amount, total_amount, commission_amount, seller_amount,
    currency, status, payment_status, 
    shipping_address, billing_address
) VALUES (
    'TEST-' || extract(epoch from now())::bigint,
    $STOREFRONT_ID,
    1, -- Фиктивный customer_id
    ${PURCHASE_QTY}00.00, -- Фиктивная цена
    ${PURCHASE_QTY}00.00,
    10.00,
    $((PURCHASE_QTY * 100 - 10)).00,
    'RSD',
    'pending',
    'pending',
    '{\"test\": true}',
    '{\"test\": true}'
);

-- Получаем ID созданного заказа
SELECT 'Создан заказ ID: ' || currval('storefront_orders_id_seq');

COMMIT;
"

ORDER_ID=$(get_value "SELECT id FROM storefront_orders WHERE order_number LIKE 'TEST-%' ORDER BY created_at DESC LIMIT 1;")
echo "Создан тестовый заказ ID: $ORDER_ID"

# Теперь создаем резервирование и обновляем остатки
if [ -n "$VARIANT_ID" ]; then
    echo "Обрабатываем вариант товара..."
    
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
    
    SELECT 'Резервирование создано ID: ' || currval('inventory_reservations_id_seq');
    
    COMMIT;
    "
    
    RESERVATION_ID=$(get_value "SELECT id FROM inventory_reservations WHERE order_id = $ORDER_ID ORDER BY created_at DESC LIMIT 1;")
    echo "Создано резервирование ID: $RESERVATION_ID"
    
else
    echo "Обрабатываем основной товар..."
    
    run_sql "
    BEGIN;
    
    -- Создаем резервирование для основного товара  
    INSERT INTO inventory_reservations (
        product_id, variant_id, quantity, order_id, status, expires_at
    ) VALUES (
        $PRODUCT_ID, NULL, $PURCHASE_QTY, $ORDER_ID, 'reserved',
        NOW() + INTERVAL '30 minutes'
    );
    
    -- Обновляем stock_quantity основного товара
    UPDATE storefront_products 
    SET stock_quantity = stock_quantity - $PURCHASE_QTY,
        updated_at = NOW()
    WHERE id = $PRODUCT_ID;
    
    COMMIT;
    "
    
    RESERVATION_ID=$(get_value "SELECT id FROM inventory_reservations WHERE order_id = $ORDER_ID ORDER BY created_at DESC LIMIT 1;")
    echo "Создано резервирование ID: $RESERVATION_ID"
fi

echo
echo "📈 3. Проверка результатов покупки"
echo "---------------------------------"

# Проверяем основной товар
CURRENT_STOCK=$(get_value "SELECT stock_quantity FROM storefront_products WHERE id = $PRODUCT_ID;")
echo "Остаток основного товара: $INITIAL_STOCK -> $CURRENT_STOCK (изменение: $((CURRENT_STOCK - INITIAL_STOCK)))"

# Проверяем вариант если есть
if [ -n "$VARIANT_ID" ]; then
    CURRENT_VARIANT_STOCK=$(get_value "SELECT stock_quantity FROM storefront_product_variants WHERE id = $VARIANT_ID;")
    CURRENT_VARIANT_RESERVED=$(get_value "SELECT reserved_quantity FROM storefront_product_variants WHERE id = $VARIANT_ID;")
    CURRENT_VARIANT_AVAILABLE=$(get_value "SELECT available_quantity FROM storefront_product_variants WHERE id = $VARIANT_ID;")
    
    echo "Остатки варианта:"
    echo "  stock_quantity: $INITIAL_VARIANT_STOCK -> $CURRENT_VARIANT_STOCK (изменение: $((CURRENT_VARIANT_STOCK - INITIAL_VARIANT_STOCK)))"
    echo "  reserved_quantity: $INITIAL_VARIANT_RESERVED -> $CURRENT_VARIANT_RESERVED (изменение: $((CURRENT_VARIANT_RESERVED - INITIAL_VARIANT_RESERVED)))"
    echo "  available_quantity: $INITIAL_VARIANT_AVAILABLE -> $CURRENT_VARIANT_AVAILABLE (изменение: $((CURRENT_VARIANT_AVAILABLE - INITIAL_VARIANT_AVAILABLE)))"
fi

echo
echo "📋 4. Проверка созданных резервирований"
echo "--------------------------------------"

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
echo "✅ Симуляция завершена!"
echo "======================"

if [ -n "$VARIANT_ID" ]; then
    if [ "$CURRENT_VARIANT_STOCK" -eq $((INITIAL_VARIANT_STOCK - PURCHASE_QTY)) ]; then
        echo "✅ УСПЕХ: stock_quantity варианта корректно уменьшился на $PURCHASE_QTY"
    else
        echo "❌ ОШИБКА: Некорректное изменение stock_quantity варианта"
    fi
    
    if [ "$CURRENT_VARIANT_RESERVED" -eq $((INITIAL_VARIANT_RESERVED + PURCHASE_QTY)) ]; then
        echo "✅ УСПЕХ: reserved_quantity варианта корректно увеличился на $PURCHASE_QTY"
    else
        echo "❌ ОШИБКА: Некорректное изменение reserved_quantity варианта"
    fi
    
    # available_quantity должен остаться примерно тем же (stock уменьшился, reserved увеличился)
    EXPECTED_AVAILABLE=$((INITIAL_VARIANT_STOCK - PURCHASE_QTY - CURRENT_VARIANT_RESERVED))
    if [ "$CURRENT_VARIANT_AVAILABLE" -eq "$EXPECTED_AVAILABLE" ]; then
        echo "✅ УСПЕХ: available_quantity варианта рассчитан корректно"
    else
        echo "⚠️  ВНИМАНИЕ: available_quantity = $CURRENT_VARIANT_AVAILABLE, ожидался $EXPECTED_AVAILABLE"
    fi
else
    if [ "$CURRENT_STOCK" -eq $((INITIAL_STOCK - PURCHASE_QTY)) ]; then
        echo "✅ УСПЕХ: stock_quantity корректно уменьшился на $PURCHASE_QTY"
    else
        echo "❌ ОШИБКА: Некорректное изменение stock_quantity"
    fi
fi

echo
echo "🧹 5. Очистка тестовых данных (опционально)"
echo "------------------------------------------"
echo "Для очистки выполните:"
echo "DELETE FROM inventory_reservations WHERE order_id = $ORDER_ID;"
echo "DELETE FROM storefront_orders WHERE id = $ORDER_ID;"
echo "-- И восстановите stock_quantity если нужно"