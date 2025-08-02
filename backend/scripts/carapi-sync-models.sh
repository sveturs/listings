#!/bin/bash

# Синхронизация моделей для всех марок

set -e

echo "=== CarAPI Models Sync ==="

# JWT из предыдущего запуска
JWT="eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJjYXJhcGkuYXBwIiwic3ViIjoiMGIwNzVmYjctMTliZS00NzgzLWFiODMtMWVhZjIxMzRhM2RjIiwiYXVkIjoiMGIwNzVmYjctMTliZS00NzgzLWFiODMtMWVhZjIxMzRhM2RjIiwiZXhwIjoxNzU0NzY1NTUxLCJpYXQiOjE3NTQxNjA3NTEsImp0aSI6ImRmZTI0MjIzLWZiODUtNDViOC1iMzFkLTA2MDA2M2FhN2E1ZiIsInVzZXIiOnsic3Vic2NyaXB0aW9ucyI6WyJiYXNlIl0sInJhdGVfbGltaXRfdHlwZSI6ImhhcmQiLCJhZGRvbnMiOnsiYW50aXF1ZV92ZWhpY2xlcyI6ZmFsc2UsImRhdGFfZmVlZCI6ZmFsc2V9fX0.S8OeXqudcoqDIzW3f00tfmx_sq889pCTDyG_9A-Sb18"

DATA_DIR="/data/hostel-booking-system/backend/data/carapi-complete-$(date +%Y%m%d-%H%M%S)"
mkdir -p "$DATA_DIR/models"

# Копируем файл марок
cp /data/hostel-booking-system/backend/data/carapi-full-20250802-205316/all_makes.json "$DATA_DIR/"

# Извлекаем ID всех марок правильно
MAKE_IDS=$(cat "$DATA_DIR/all_makes.json" | jq -r '.data[].id')

echo "Found $(echo "$MAKE_IDS" | wc -l) makes to process"

# Функция для загрузки
fetch_models() {
    local make_id=$1
    local make_name=$2
    local output="$DATA_DIR/models/make_${make_id}_${make_name}.json"
    
    echo "Fetching models for $make_name (ID: $make_id)..."
    
    if curl -s -H "Authorization: Bearer $JWT" \
            -H "Accept: application/json" \
            "https://carapi.app/api/models?make_id=$make_id" \
            -o "$output"; then
        
        # Проверяем что получили валидные данные
        if [ -s "$output" ] && jq -e '.data | length' "$output" >/dev/null 2>&1; then
            local count=$(jq '.data | length' "$output")
            echo "✓ Found $count models for $make_name"
            return 0
        else
            echo "✗ No models found for $make_name"
            rm -f "$output"
            return 1
        fi
    else
        echo "✗ Failed to fetch models for $make_name"
        return 1
    fi
}

# Загружаем модели для каждой марки
for make_id in $MAKE_IDS; do
    make_name=$(cat "$DATA_DIR/all_makes.json" | jq -r ".data[] | select(.id == $make_id) | .name" | tr ' ' '_')
    fetch_models "$make_id" "$make_name" || true
    sleep 0.3  # Rate limiting
done

# Статистика
echo "=== Statistics ==="
echo "Total model files: $(ls -1 "$DATA_DIR/models" | wc -l)"
echo "Total models: $(cat "$DATA_DIR/models"/*.json 2>/dev/null | jq -s 'map(.data // [] | length) | add' || echo 0)"

# Сохраняем популярные для дальнейшей обработки
mkdir -p "$DATA_DIR/popular"
for brand in "Volkswagen" "Toyota" "BMW" "Mercedes-Benz" "Audi" "Ford" "Honda" "Nissan" "Hyundai" "Kia"; do
    brand_file=$(ls "$DATA_DIR/models/"*"$(echo $brand | tr ' ' '_')"*.json 2>/dev/null | head -1)
    if [ -n "$brand_file" ] && [ -f "$brand_file" ]; then
        cp "$brand_file" "$DATA_DIR/popular/"
        echo "Copied $brand to popular"
    fi
done

echo "Data saved to: $DATA_DIR"