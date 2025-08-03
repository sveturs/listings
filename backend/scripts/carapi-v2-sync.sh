#!/bin/bash

# Синхронизация через API v2

set -e

echo "=== CarAPI V2 Full Sync ==="

# JWT из предыдущего запуска
JWT="eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJjYXJhcGkuYXBwIiwic3ViIjoiMGIwNzVmYjctMTliZS00NzgzLWFiODMtMWVhZjIxMzRhM2RjIiwiYXVkIjoiMGIwNzVmYjctMTliZS00NzgzLWFiODMtMWVhZjIxMzRhM2RjIiwiZXhwIjoxNzU0NzY1NTUxLCJpYXQiOjE3NTQxNjA3NTEsImp0aSI6ImRmZTI0MjIzLWZiODUtNDViOC1iMzFkLTA2MDA2M2FhN2E1ZiIsInVzZXIiOnsic3Vic2NyaXB0aW9ucyI6WyJiYXNlIl0sInJhdGVfbGltaXRfdHlwZSI6ImhhcmQiLCJhZGRvbnMiOnsiYW50aXF1ZV92ZWhpY2xlcyI6ZmFsc2UsImRhdGFfZmVlZCI6ZmFsc2V9fX0.S8OeXqudcoqDIzW3f00tfmx_sq889pCTDyG_9A-Sb18"

DATA_DIR="/data/hostel-booking-system/backend/data/carapi-v2-$(date +%Y%m%d-%H%M%S)"
mkdir -p "$DATA_DIR"

# Проверяем какие эндпоинты доступны на v2
echo "Testing v2 endpoints..."

# Пробуем разные варианты
echo "1. Testing /api/v2/makes..."
curl -s -H "Authorization: Bearer $JWT" \
     -H "Accept: application/json" \
     "https://carapi.app/api/v2/makes" \
     -o "$DATA_DIR/v2_makes_test.json"

echo "2. Testing /api/models/v2..."
curl -s -H "Authorization: Bearer $JWT" \
     -H "Accept: application/json" \
     "https://carapi.app/api/models/v2?make_id=43" \
     -o "$DATA_DIR/v2_models_test.json"

echo "3. Testing /api/v2/models..."
curl -s -H "Authorization: Bearer $JWT" \
     -H "Accept: application/json" \
     "https://carapi.app/api/v2/models?make_id=43" \
     -o "$DATA_DIR/v2_models_test2.json"

# Проверяем что вернулось
echo -e "\n=== Results ==="
for file in "$DATA_DIR"/*.json; do
    echo -e "\n--- $file ---"
    head -c 500 "$file" | jq . 2>/dev/null || cat "$file" | head -5
done

# Пробуем получить все модели без фильтра
echo -e "\n4. Testing all models..."
curl -s -H "Authorization: Bearer $JWT" \
     -H "Accept: application/json" \
     "https://carapi.app/api/models" \
     -o "$DATA_DIR/all_models.json"

# Размер файла
ls -lh "$DATA_DIR/all_models.json"

# Пробуем с пагинацией
echo -e "\n5. Testing with pagination..."
curl -s -H "Authorization: Bearer $JWT" \
     -H "Accept: application/json" \
     "https://carapi.app/api/models?limit=100&page=1" \
     -o "$DATA_DIR/models_page1.json"

# VIN тест
echo -e "\n6. Testing VIN decoder..."
curl -s -H "Authorization: Bearer $JWT" \
     -H "Accept: application/json" \
     "https://carapi.app/api/vin/decode?vin=WVWZZZ1JZ3W386752" \
     -o "$DATA_DIR/vin_test.json"

cat "$DATA_DIR/vin_test.json" | jq .

echo -e "\nData saved to: $DATA_DIR"