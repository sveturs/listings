#!/bin/bash
# scripts/cleanup_opensearch.sh
# 🔴 CRITICAL: Удаление всех OpenSearch индексов перед миграцией категорий
# Выполнить ПЕРЕД миграциями PostgreSQL

set -e

OPENSEARCH_URL="${OPENSEARCH_URL:-http://localhost:9200}"

echo "🗑️  Удаление старых индексов OpenSearch..."
echo "OpenSearch URL: $OPENSEARCH_URL"
echo ""

# Проверка доступности OpenSearch
if ! curl -s "$OPENSEARCH_URL" >/dev/null 2>&1; then
    echo "❌ OpenSearch недоступен по адресу $OPENSEARCH_URL"
    echo "Проверьте что OpenSearch запущен и доступен"
    exit 1
fi

echo "✅ OpenSearch доступен"
echo ""

# Показываем текущие индексы
echo "=== Текущие индексы ==="
curl -s "$OPENSEARCH_URL/_cat/indices?v" 2>/dev/null | grep -E "(marketplace|categories|products)" || echo "Нет индексов для удаления"
echo ""

# Удаляем все marketplace индексы
echo "=== Удаление индексов ==="

# 1. Marketplace индексы
echo "1. Удаление marketplace_* индексов..."
curl -s -X DELETE "$OPENSEARCH_URL/marketplace_*" 2>/dev/null
echo ""

# 2. Categories индексы
echo "2. Удаление categories_* индексов..."
curl -s -X DELETE "$OPENSEARCH_URL/categories_*" 2>/dev/null
echo ""

# 3. Products индексы
echo "3. Удаление products_* индексов..."
curl -s -X DELETE "$OPENSEARCH_URL/products_*" 2>/dev/null
echo ""

# Удаляем aliases
echo "=== Удаление aliases ==="
curl -s -X POST "$OPENSEARCH_URL/_aliases" -H 'Content-Type: application/json' -d '{
  "actions": [
    { "remove": { "index": "*", "alias": "marketplace" } },
    { "remove": { "index": "*", "alias": "categories" } },
    { "remove": { "index": "*", "alias": "products" } }
  ]
}' 2>/dev/null
echo ""

# Проверка результата
echo ""
echo "=== Результат ==="
REMAINING=$(curl -s "$OPENSEARCH_URL/_cat/indices?v" 2>/dev/null | grep -cE "(marketplace|categories|products)" || echo "0")

if [ "$REMAINING" -eq "0" ]; then
    echo "✅ OpenSearch полностью очищен"
    echo "   Все marketplace/categories/products индексы удалены"
else
    echo "⚠️  Внимание: найдено $REMAINING индексов"
    curl -s "$OPENSEARCH_URL/_cat/indices?v" 2>/dev/null | grep -E "(marketplace|categories|products)"
fi

echo ""
echo "=== Следующий шаг ==="
echo "Теперь выполните миграции PostgreSQL:"
echo "  cd /p/github.com/vondi-global/listings"
echo "  # Выполнить миграцию 000_cleanup"
echo ""
