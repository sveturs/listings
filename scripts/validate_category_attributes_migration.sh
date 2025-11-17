#!/bin/bash

# Validation script for category_attributes migration
# Проверяет корректность миграции данных

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Database connection strings
SOURCE_DB="postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5433/svetubd?sslmode=disable"
DEST_DB="postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable"

echo -e "${BLUE}╔════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║   Валидация миграции category_attributes                      ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════════╝${NC}"
echo ""

# 1. Проверка количества записей
echo -e "${YELLOW}[1/7] Проверка количества записей...${NC}"

SOURCE_COUNT=$(psql "$SOURCE_DB" -t -c "SELECT COUNT(*) FROM unified_category_attributes;")
DEST_COUNT=$(psql "$DEST_DB" -t -c "SELECT COUNT(*) FROM category_attributes;")

echo "  📊 Источник (монолит):     $SOURCE_COUNT записей"
echo "  📊 Получатель (микросервис): $DEST_COUNT записей"

if [ "$SOURCE_COUNT" -eq "$DEST_COUNT" ]; then
    echo -e "  ${GREEN}✅ Количество записей совпадает${NC}"
else
    DIFF=$((SOURCE_COUNT - DEST_COUNT))
    echo -e "  ${YELLOW}⚠️  Разница: $DIFF записей${NC}"
fi
echo ""

# 2. Проверка уникальных категорий
echo -e "${YELLOW}[2/7] Проверка уникальных категорий...${NC}"

SOURCE_CATEGORIES=$(psql "$SOURCE_DB" -t -c "SELECT COUNT(DISTINCT category_id) FROM unified_category_attributes;")
DEST_CATEGORIES=$(psql "$DEST_DB" -t -c "SELECT COUNT(DISTINCT category_id) FROM category_attributes;")

echo "  📂 Источник: $SOURCE_CATEGORIES уникальных категорий"
echo "  📂 Получатель: $DEST_CATEGORIES уникальных категорий"

if [ "$SOURCE_CATEGORIES" -eq "$DEST_CATEGORIES" ]; then
    echo -e "  ${GREEN}✅ Количество категорий совпадает${NC}"
else
    echo -e "  ${YELLOW}⚠️  Разница в категориях${NC}"
fi
echo ""

# 3. Проверка уникальных атрибутов
echo -e "${YELLOW}[3/7] Проверка уникальных атрибутов...${NC}"

SOURCE_ATTRIBUTES=$(psql "$SOURCE_DB" -t -c "SELECT COUNT(DISTINCT attribute_id) FROM unified_category_attributes;")
DEST_ATTRIBUTES=$(psql "$DEST_DB" -t -c "SELECT COUNT(DISTINCT attribute_id) FROM category_attributes;")

echo "  🏷️  Источник: $SOURCE_ATTRIBUTES уникальных атрибутов"
echo "  🏷️  Получатель: $DEST_ATTRIBUTES уникальных атрибутов"

if [ "$SOURCE_ATTRIBUTES" -eq "$DEST_ATTRIBUTES" ]; then
    echo -e "  ${GREEN}✅ Количество атрибутов совпадает${NC}"
else
    echo -e "  ${YELLOW}⚠️  Разница в атрибутах${NC}"
fi
echo ""

# 4. Проверка связей category_id + attribute_id (должны быть уникальными)
echo -e "${YELLOW}[4/7] Проверка уникальности пар (category_id, attribute_id)...${NC}"

DEST_DUPLICATES=$(psql "$DEST_DB" -t -c "
    SELECT COUNT(*)
    FROM (
        SELECT category_id, attribute_id, COUNT(*)
        FROM category_attributes
        GROUP BY category_id, attribute_id
        HAVING COUNT(*) > 1
    ) AS duplicates;
")

if [ "$DEST_DUPLICATES" -eq 0 ]; then
    echo -e "  ${GREEN}✅ Дубликатов не найдено${NC}"
else
    echo -e "  ${RED}❌ Найдено $DEST_DUPLICATES дубликатов!${NC}"
fi
echo ""

# 5. Проверка is_enabled распределения
echo -e "${YELLOW}[5/7] Проверка распределения is_enabled...${NC}"

SOURCE_ENABLED=$(psql "$SOURCE_DB" -t -c "SELECT COUNT(*) FROM unified_category_attributes WHERE is_enabled = true;")
DEST_ENABLED=$(psql "$DEST_DB" -t -c "SELECT COUNT(*) FROM category_attributes WHERE is_enabled = true;")

echo "  ✓ Источник (enabled=true):     $SOURCE_ENABLED"
echo "  ✓ Получатель (enabled=true):   $DEST_ENABLED"

if [ "$SOURCE_ENABLED" -eq "$DEST_ENABLED" ]; then
    echo -e "  ${GREEN}✅ Распределение enabled совпадает${NC}"
else
    echo -e "  ${YELLOW}⚠️  Разница в enabled записях${NC}"
fi
echo ""

# 6. Проверка is_required распределения
echo -e "${YELLOW}[6/7] Проверка распределения is_required...${NC}"

SOURCE_REQUIRED=$(psql "$SOURCE_DB" -t -c "SELECT COUNT(*) FROM unified_category_attributes WHERE is_required = true;")
DEST_REQUIRED=$(psql "$DEST_DB" -t -c "SELECT COUNT(*) FROM category_attributes WHERE is_required = true;")

echo "  ⚡ Источник (required=true):    $SOURCE_REQUIRED"
echo "  ⚡ Получатель (required=true):  $DEST_REQUIRED"

if [ "$SOURCE_REQUIRED" -eq "$DEST_REQUIRED" ]; then
    echo -e "  ${GREEN}✅ Распределение required совпадает${NC}"
else
    echo -e "  ${YELLOW}⚠️  Разница в required записях${NC}"
fi
echo ""

# 7. Проверка конкретных примеров
echo -e "${YELLOW}[7/7] Проверка конкретных примеров...${NC}"

echo "  Сравнение записей для category_id=1001:"
psql "$SOURCE_DB" -c "
    SELECT category_id, attribute_id, is_enabled, is_required, sort_order
    FROM unified_category_attributes
    WHERE category_id = 1001
    ORDER BY sort_order
    LIMIT 5;
" | head -n 10

echo ""
echo "  В микросервисе:"
psql "$DEST_DB" -c "
    SELECT category_id, attribute_id, is_enabled, is_required, sort_order
    FROM category_attributes
    WHERE category_id = 1001
    ORDER BY sort_order
    LIMIT 5;
" | head -n 10

echo ""

# 8. Проверка foreign key ссылок
echo -e "${YELLOW}[8/8] Проверка целостности foreign key...${NC}"

INVALID_CATEGORIES=$(psql "$DEST_DB" -t -c "
    SELECT COUNT(*)
    FROM category_attributes ca
    LEFT JOIN categories c ON ca.category_id = c.id
    WHERE c.id IS NULL;
")

INVALID_ATTRIBUTES=$(psql "$DEST_DB" -t -c "
    SELECT COUNT(*)
    FROM category_attributes ca
    LEFT JOIN attributes a ON ca.attribute_id = a.id
    WHERE a.id IS NULL;
")

if [ "$INVALID_CATEGORIES" -eq 0 ]; then
    echo -e "  ${GREEN}✅ Все category_id ссылки валидны${NC}"
else
    echo -e "  ${RED}❌ Найдено $INVALID_CATEGORIES невалидных category_id!${NC}"
fi

if [ "$INVALID_ATTRIBUTES" -eq 0 ]; then
    echo -e "  ${GREEN}✅ Все attribute_id ссылки валидны${NC}"
else
    echo -e "  ${RED}❌ Найдено $INVALID_ATTRIBUTES невалидных attribute_id!${NC}"
fi

echo ""

# Summary
echo -e "${BLUE}╔════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║   ИТОГОВАЯ СТАТИСТИКА                                         ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════════╝${NC}"

if [ "$SOURCE_COUNT" -eq "$DEST_COUNT" ] && \
   [ "$DEST_DUPLICATES" -eq 0 ] && \
   [ "$INVALID_CATEGORIES" -eq 0 ] && \
   [ "$INVALID_ATTRIBUTES" -eq 0 ]; then
    echo -e "${GREEN}✅ ВСЕ ПРОВЕРКИ ПРОЙДЕНЫ УСПЕШНО!${NC}"
    exit 0
else
    echo -e "${YELLOW}⚠️  ОБНАРУЖЕНЫ НЕКОТОРЫЕ РАСХОЖДЕНИЯ${NC}"
    echo -e "    Проверьте детали выше"
    exit 1
fi
