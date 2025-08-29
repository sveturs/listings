#!/bin/bash

# Скрипт для аудита изображений в MinIO и PostgreSQL
# Проверяет синхронизацию между хранилищем и базой данных

set -e

# Переменные подключения
DB_URL="postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5432/svetubd?sslmode=disable"
MINIO_ALIAS="myminio"
MINIO_BUCKET="listings"
STOREFRONT_BUCKET="storefront-products"

echo "========================================="
echo "АУДИТ ИЗОБРАЖЕНИЙ MARKETPLACE LISTINGS"
echo "========================================="

# Создаем временный файл для результатов
REPORT="/tmp/images_audit_$(date +%Y%m%d_%H%M%S).txt"

echo "Начало аудита: $(date)" > "$REPORT"
echo "" >> "$REPORT"

# 1. Получаем все объявления с изображениями из БД
echo "Анализ marketplace_listings..." | tee -a "$REPORT"
echo "================================" >> "$REPORT"

# Получаем список всех объявлений
LISTINGS=$(psql "$DB_URL" -t -A -c "SELECT DISTINCT listing_id FROM marketplace_images ORDER BY listing_id")

TOTAL_LISTINGS=$(echo "$LISTINGS" | wc -l)
PROBLEMS_FOUND=0

echo "Найдено объявлений с изображениями: $TOTAL_LISTINGS" | tee -a "$REPORT"
echo "" >> "$REPORT"

# Проверяем каждое объявление
for LISTING_ID in $LISTINGS; do
    echo -n "Проверка объявления $LISTING_ID... "
    
    # Получаем изображения из БД
    DB_IMAGES=$(psql "$DB_URL" -t -A -c "
        SELECT file_path, file_name, is_main, public_url 
        FROM marketplace_images 
        WHERE listing_id = $LISTING_ID 
        ORDER BY is_main DESC, id
    ")
    
    DB_COUNT=$(echo "$DB_IMAGES" | grep -c "|" || echo 0)
    
    # Получаем изображения из MinIO
    MINIO_IMAGES=$(docker exec minio mc ls "$MINIO_ALIAS/$MINIO_BUCKET/$LISTING_ID/" 2>/dev/null | awk '{print $NF}' || echo "")
    MINIO_COUNT=$(echo "$MINIO_IMAGES" | grep -v "^$" | wc -l)
    
    # Сравниваем количество
    if [ "$DB_COUNT" -ne "$MINIO_COUNT" ]; then
        echo "ПРОБЛЕМА!" | tee -a "$REPORT"
        echo "  Объявление $LISTING_ID: БД=$DB_COUNT файлов, MinIO=$MINIO_COUNT файлов" >> "$REPORT"
        echo "  Файлы в БД:" >> "$REPORT"
        echo "$DB_IMAGES" | sed 's/^/    /' >> "$REPORT"
        echo "  Файлы в MinIO:" >> "$REPORT"
        echo "$MINIO_IMAGES" | sed 's/^/    /' >> "$REPORT"
        echo "" >> "$REPORT"
        ((PROBLEMS_FOUND++))
    else
        # Проверяем правильность путей
        HAS_WRONG_PATH=0
        while IFS='|' read -r file_path file_name is_main public_url; do
            # Проверяем наличие IP адресов в путях
            if echo "$public_url" | grep -q "100.88.44.15\|localhost:9000"; then
                if [ "$HAS_WRONG_PATH" -eq 0 ]; then
                    echo "НЕПРАВИЛЬНЫЕ ПУТИ!" | tee -a "$REPORT"
                    echo "  Объявление $LISTING_ID имеет неправильные URL:" >> "$REPORT"
                    HAS_WRONG_PATH=1
                    ((PROBLEMS_FOUND++))
                fi
                echo "    $public_url" >> "$REPORT"
            fi
        done <<< "$DB_IMAGES"
        
        if [ "$HAS_WRONG_PATH" -eq 0 ]; then
            echo "OK"
        fi
    fi
done

echo "" >> "$REPORT"
echo "=========================================" >> "$REPORT"
echo "АУДИТ STOREFRONT PRODUCTS" >> "$REPORT"
echo "=========================================" >> "$REPORT"

# 2. Проверяем изображения товаров витрин
echo "" | tee -a "$REPORT"
echo "Анализ storefront_products..." | tee -a "$REPORT"

# Получаем все товары с изображениями
PRODUCTS=$(psql "$DB_URL" -t -A -c "
    SELECT DISTINCT storefront_product_id 
    FROM storefront_product_images 
    ORDER BY storefront_product_id
")

TOTAL_PRODUCTS=$(echo "$PRODUCTS" | wc -l)
echo "Найдено товаров с изображениями: $TOTAL_PRODUCTS" | tee -a "$REPORT"

for PRODUCT_ID in $PRODUCTS; do
    echo -n "Проверка товара $PRODUCT_ID... "
    
    # Получаем изображения из БД
    PRODUCT_IMAGES=$(psql "$DB_URL" -t -A -c "
        SELECT image_url, thumbnail_url, is_default 
        FROM storefront_product_images 
        WHERE storefront_product_id = $PRODUCT_ID 
        ORDER BY display_order
    ")
    
    # Проверяем наличие IP адресов в URL
    HAS_PROBLEM=0
    while IFS='|' read -r image_url thumbnail_url is_default; do
        if echo "$image_url" | grep -q "100.88.44.15"; then
            if [ "$HAS_PROBLEM" -eq 0 ]; then
                echo "НЕПРАВИЛЬНЫЕ URL!" | tee -a "$REPORT"
                echo "  Товар $PRODUCT_ID имеет URL с IP адресами:" >> "$REPORT"
                HAS_PROBLEM=1
                ((PROBLEMS_FOUND++))
            fi
            echo "    $image_url" >> "$REPORT"
        fi
    done <<< "$PRODUCT_IMAGES"
    
    if [ "$HAS_PROBLEM" -eq 0 ]; then
        echo "OK"
    fi
done

echo "" >> "$REPORT"
echo "=========================================" >> "$REPORT"
echo "ИТОГОВАЯ СТАТИСТИКА" >> "$REPORT"
echo "=========================================" >> "$REPORT"
echo "Проверено объявлений: $TOTAL_LISTINGS" >> "$REPORT"
echo "Проверено товаров: $TOTAL_PRODUCTS" >> "$REPORT"
echo "Найдено проблем: $PROBLEMS_FOUND" >> "$REPORT"
echo "" >> "$REPORT"
echo "Завершение аудита: $(date)" >> "$REPORT"

echo ""
echo "========================================="
echo "РЕЗУЛЬТАТЫ АУДИТА"
echo "========================================="
echo "Проверено объявлений: $TOTAL_LISTINGS"
echo "Проверено товаров: $TOTAL_PRODUCTS"
echo "Найдено проблем: $PROBLEMS_FOUND"
echo ""
echo "Полный отчет сохранен в: $REPORT"

# Показываем первые проблемы
if [ "$PROBLEMS_FOUND" -gt 0 ]; then
    echo ""
    echo "Первые 10 проблем:"
    grep -A 3 "ПРОБЛЕМА\|НЕПРАВИЛЬНЫЕ" "$REPORT" | head -30
fi