#!/bin/bash

# Скрипт для переиндексации товаров витрин в OpenSearch

echo "🔄 Начинаем переиндексацию товаров витрин..."

# Получаем все товары витрин из БД
PRODUCTS=$(psql "postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5432/svetubd?sslmode=disable" -t -c "
SELECT json_build_object(
    'id', sp.id,
    'storefront_id', sp.storefront_id,
    'category_id', sp.category_id,
    'name', sp.name,
    'description', COALESCE(sp.description, ''),
    'price', sp.price,
    'currency', sp.currency,
    'stock_quantity', sp.stock_quantity,
    'stock_status', sp.stock_status,
    'is_active', sp.is_active,
    'sku', sp.sku,
    'storefront_name', s.name,
    'storefront_slug', s.slug
)
FROM storefront_products sp
LEFT JOIN storefronts s ON s.id = sp.storefront_id
WHERE sp.is_active = true
ORDER BY sp.id;")

# Счетчик успешно проиндексированных товаров
SUCCESS_COUNT=0
FAIL_COUNT=0

# Обрабатываем каждый товар
echo "$PRODUCTS" | while IFS= read -r product_json; do
    # Пропускаем пустые строки
    if [ -z "$product_json" ]; then
        continue
    fi
    
    # Извлекаем данные из JSON
    ID=$(echo "$product_json" | jq -r '.id')
    STOREFRONT_ID=$(echo "$product_json" | jq -r '.storefront_id')
    CATEGORY_ID=$(echo "$product_json" | jq -r '.category_id // 0')
    NAME=$(echo "$product_json" | jq -r '.name')
    DESCRIPTION=$(echo "$product_json" | jq -r '.description // ""')
    PRICE=$(echo "$product_json" | jq -r '.price')
    CURRENCY=$(echo "$product_json" | jq -r '.currency // "RSD"')
    STOCK_QUANTITY=$(echo "$product_json" | jq -r '.stock_quantity')
    STOCK_STATUS=$(echo "$product_json" | jq -r '.stock_status')
    SKU=$(echo "$product_json" | jq -r '.sku // ""')
    STOREFRONT_NAME=$(echo "$product_json" | jq -r '.storefront_name // ""')
    STOREFRONT_SLUG=$(echo "$product_json" | jq -r '.storefront_slug // ""')
    
    # Определяем статус наличия
    if [ "$STOCK_QUANTITY" -le 0 ]; then
        STOCK_STATUS="out_of_stock"
    elif [ "$STOCK_QUANTITY" -le 5 ]; then
        STOCK_STATUS="low_stock"
    else
        STOCK_STATUS="in_stock"
    fi
    
    # Формируем документ для индексации
    DOC=$(cat <<EOFDOC
{
  "product_id": $ID,
  "product_type": "storefront",
  "storefront_id": $STOREFRONT_ID,
  "category_id": $CATEGORY_ID,
  "name": "$(echo "$NAME" | sed 's/"/\\"/g')",
  "name_lowercase": "$(echo "$NAME" | tr '[:upper:]' '[:lower:]' | sed 's/"/\\"/g')",
  "description": "$(echo "$DESCRIPTION" | sed 's/"/\\"/g')",
  "price": $PRICE,
  "currency": "$CURRENCY",
  "stock_quantity": $STOCK_QUANTITY,
  "stock_status": "$STOCK_STATUS",
  "sku": "$SKU",
  "is_active": true,
  "inventory": {
    "quantity": $STOCK_QUANTITY,
    "in_stock": $([ "$STOCK_QUANTITY" -gt 0 ] && echo "true" || echo "false"),
    "available": $STOCK_QUANTITY,
    "low_stock": $([ "$STOCK_QUANTITY" -gt 0 ] && [ "$STOCK_QUANTITY" -le 5 ] && echo "true" || echo "false"),
    "status": "$STOCK_STATUS"
  },
  "storefront": {
    "id": $STOREFRONT_ID,
    "name": "$(echo "$STOREFRONT_NAME" | sed 's/"/\\"/g')",
    "slug": "$(echo "$STOREFRONT_SLUG" | sed 's/"/\\"/g')"
  },
  "created_at": "$(date -Iseconds)",
  "updated_at": "$(date -Iseconds)"
}
EOFDOC
)
    
    # Индексируем документ в OpenSearch
    RESPONSE=$(curl -s -X POST "http://localhost:9200/marketplace/_doc/sp_$ID" \
        -H "Content-Type: application/json" \
        -d "$DOC")
    
    # Проверяем результат
    if echo "$RESPONSE" | grep -q '"result":"created"\|"result":"updated"'; then
        echo "✅ Товар ID=$ID успешно проиндексирован (остаток: $STOCK_QUANTITY)"
        SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
    else
        echo "❌ Ошибка при индексации товара ID=$ID"
        echo "   Ответ: $RESPONSE"
        FAIL_COUNT=$((FAIL_COUNT + 1))
    fi
done

echo ""
echo "📊 Результаты переиндексации:"
echo "   ✅ Успешно: $SUCCESS_COUNT товаров"
echo "   ❌ Ошибок: $FAIL_COUNT товаров"
echo ""
echo "✅ Переиндексация завершена\!"
