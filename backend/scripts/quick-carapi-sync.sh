#!/bin/bash

# Быстрая синхронизация всех доступных данных из CarAPI

set -e

echo "=== Quick CarAPI Sync ==="
echo "Starting at: $(date)"

# Credentials
API_TOKEN="1a0bef70-bdd4-4b43-afea-5836c36d32b1"
API_SECRET="707db90e9d68a5a86a00ead3a5ae0663"

# Получаем JWT
echo "Getting JWT token..."
JWT=$(curl -s -X 'POST' \
  'https://carapi.app/api/auth/login' \
  -H 'accept: text/plain' \
  -H 'Content-Type: application/json' \
  -d "{\"api_token\": \"$API_TOKEN\", \"api_secret\": \"$API_SECRET\"}")

echo "JWT obtained: ${JWT:0:50}..."

# Создаем директорию для данных
DATA_DIR="/data/hostel-booking-system/backend/data/carapi-full-$(date +%Y%m%d-%H%M%S)"
mkdir -p "$DATA_DIR"

# Функция для скачивания с повторами
fetch_with_retry() {
    local url=$1
    local output=$2
    local retries=3
    
    for i in $(seq 1 $retries); do
        echo "Fetching: $url (attempt $i)"
        if curl -s -H "Authorization: Bearer $JWT" \
                -H "Accept: application/json" \
                "$url" > "$output"; then
            if [ -s "$output" ] && jq empty "$output" 2>/dev/null; then
                echo "✓ Saved to $output"
                return 0
            fi
        fi
        sleep 1
    done
    return 1
}

# 1. Скачиваем ВСЕ марки
echo "=== Fetching ALL makes ==="
fetch_with_retry "https://carapi.app/api/makes" "$DATA_DIR/all_makes.json"

# Получаем количество марок
MAKE_COUNT=$(jq 'length' "$DATA_DIR/all_makes.json" 2>/dev/null || echo "0")
echo "Found $MAKE_COUNT makes"

# 2. Скачиваем модели для ВСЕХ марок
echo "=== Fetching models for ALL makes ==="
mkdir -p "$DATA_DIR/models"

# Извлекаем ID всех марок
MAKE_IDS=$(jq -r '.[].id' "$DATA_DIR/all_makes.json" 2>/dev/null || echo "")

for make_id in $MAKE_IDS; do
    if [ -n "$make_id" ]; then
        fetch_with_retry "https://carapi.app/api/models?make_id=$make_id" \
                        "$DATA_DIR/models/make_${make_id}.json"
        sleep 0.2  # Небольшая задержка
    fi
done

# 3. Скачиваем комплектации для последних лет
echo "=== Fetching trims for recent years ==="
mkdir -p "$DATA_DIR/trims"

CURRENT_YEAR=$(date +%Y)
YEARS="$((CURRENT_YEAR-2)) $((CURRENT_YEAR-1)) $CURRENT_YEAR $((CURRENT_YEAR+1))"

# Для каждой модели и года
for model_file in "$DATA_DIR/models"/*.json; do
    if [ -f "$model_file" ]; then
        MODEL_IDS=$(jq -r '.[].id' "$model_file" 2>/dev/null || echo "")
        for model_id in $MODEL_IDS; do
            if [ -n "$model_id" ]; then
                for year in $YEARS; do
                    fetch_with_retry "https://carapi.app/api/trims?model_id=$model_id&year=$year" \
                                    "$DATA_DIR/trims/model_${model_id}_year_${year}.json" || true
                    sleep 0.1
                done
            fi
        done
    fi
done

# 4. Скачиваем дополнительные данные по годам
echo "=== Fetching year data ==="
mkdir -p "$DATA_DIR/years"
for year in $(seq 2020 2026); do
    fetch_with_retry "https://carapi.app/api/years" "$DATA_DIR/years/year_${year}.json" || true
done

# 5. Создаем статистику
echo "=== Creating statistics ==="
cat > "$DATA_DIR/statistics.txt" << EOF
CarAPI Full Data Sync
=====================
Date: $(date)
JWT Token: ${JWT:0:50}...

Data collected:
- Makes: $MAKE_COUNT
- Model files: $(ls -1 "$DATA_DIR/models" 2>/dev/null | wc -l)
- Trim files: $(ls -1 "$DATA_DIR/trims" 2>/dev/null | wc -l)
- Total size: $(du -sh "$DATA_DIR" | cut -f1)

Directory: $DATA_DIR
EOF

cat "$DATA_DIR/statistics.txt"

# 6. Создаем архив
echo "Creating archive..."
cd "$(dirname "$DATA_DIR")"
tar -czf "$(basename "$DATA_DIR").tar.gz" "$(basename "$DATA_DIR")"

echo "=== COMPLETED ==="
echo "Data directory: $DATA_DIR"
echo "Archive: $(dirname "$DATA_DIR")/$(basename "$DATA_DIR").tar.gz"