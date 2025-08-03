#!/bin/bash

# Финальная синхронизация всех доступных данных CarAPI

set -e

echo "=== CarAPI Final Full Sync ==="
echo "Starting at: $(date)"

# JWT токен
JWT="eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJjYXJhcGkuYXBwIiwic3ViIjoiMGIwNzVmYjctMTliZS00NzgzLWFiODMtMWVhZjIxMzRhM2RjIiwiYXVkIjoiMGIwNzVmYjctMTliZS00NzgzLWFiODMtMWVhZjIxMzRhM2RjIiwiZXhwIjoxNzU0NzY1NTUxLCJpYXQiOjE3NTQxNjA3NTEsImp0aSI6ImRmZTI0MjIzLWZiODUtNDViOC1iMzFkLTA2MDA2M2FhN2E1ZiIsInVzZXIiOnsic3Vic2NyaXB0aW9ucyI6WyJiYXNlIl0sInJhdGVfbGltaXRfdHlwZSI6ImhhcmQiLCJhZGRvbnMiOnsiYW50aXF1ZV92ZWhpY2xlcyI6ZmFsc2UsImRhdGFfZmVlZCI6ZmFsc2V9fX0.S8OeXqudcoqDIzW3f00tfmx_sq889pCTDyG_9A-Sb18"

DATA_DIR="/data/hostel-booking-system/backend/data/carapi-final-$(date +%Y%m%d-%H%M%S)"
mkdir -p "$DATA_DIR"/{makes,models,trims,years}

# Статистика
STATS_FILE="$DATA_DIR/statistics.txt"
echo "CarAPI Final Sync Statistics" > "$STATS_FILE"
echo "============================" >> "$STATS_FILE"
echo "Start time: $(date)" >> "$STATS_FILE"

# 1. Копируем марки из предыдущего запуска
echo "Copying makes data..."
cp /data/hostel-booking-system/backend/data/carapi-full-20250802-205316/all_makes.json "$DATA_DIR/makes/"
MAKES_COUNT=$(jq '.data | length' "$DATA_DIR/makes/all_makes.json")
echo "Makes: $MAKES_COUNT" >> "$STATS_FILE"

# 2. Загружаем модели для всех марок через v2 API
echo "=== Fetching models for all makes ==="
MAKE_IDS=$(jq -r '.data[].id' "$DATA_DIR/makes/all_makes.json")

MODELS_TOTAL=0
for make_id in $MAKE_IDS; do
    make_name=$(jq -r ".data[] | select(.id == $make_id) | .name" "$DATA_DIR/makes/all_makes.json" | tr ' ' '_')
    echo -n "Fetching models for $make_name (ID: $make_id)... "
    
    if curl -s -H "Authorization: Bearer $JWT" \
            -H "Accept: application/json" \
            "https://carapi.app/api/models/v2?make_id=$make_id" \
            -o "$DATA_DIR/models/make_${make_id}_${make_name}.json"; then
        
        count=$(jq '.data | length' "$DATA_DIR/models/make_${make_id}_${make_name}.json" 2>/dev/null || echo 0)
        echo "found $count models"
        MODELS_TOTAL=$((MODELS_TOTAL + count))
    else
        echo "failed"
    fi
    sleep 0.2
done

echo "Total models: $MODELS_TOTAL" >> "$STATS_FILE"

# 3. Загружаем комплектации для популярных моделей последних лет
echo -e "\n=== Fetching trims for popular models ==="
CURRENT_YEAR=$(date +%Y)

# Популярные модели ID (из Volkswagen, Toyota, etc)
POPULAR_MODELS=(
    # Volkswagen
    4891  # Golf
    4563  # Jetta
    5343  # Passat
    6765  # Tiguan
    6358  # Touareg
    7464  # ID.4
    # Toyota (нужно найти ID)
)

TRIMS_TOTAL=0
for model_id in "${POPULAR_MODELS[@]}"; do
    for year in $((CURRENT_YEAR-2)) $((CURRENT_YEAR-1)) $CURRENT_YEAR $((CURRENT_YEAR+1)); do
        echo -n "Fetching trims for model $model_id, year $year... "
        
        if curl -s -H "Authorization: Bearer $JWT" \
                -H "Accept: application/json" \
                "https://carapi.app/api/trims?model_id=$model_id&year=$year" \
                -o "$DATA_DIR/trims/model_${model_id}_year_${year}.json"; then
            
            if [ -s "$DATA_DIR/trims/model_${model_id}_year_${year}.json" ]; then
                count=$(jq 'if .data then .data | length else 0 end' "$DATA_DIR/trims/model_${model_id}_year_${year}.json" 2>/dev/null || echo 0)
                if [ "$count" -gt 0 ]; then
                    echo "found $count trims"
                    TRIMS_TOTAL=$((TRIMS_TOTAL + count))
                else
                    echo "no data"
                    rm -f "$DATA_DIR/trims/model_${model_id}_year_${year}.json"
                fi
            else
                echo "empty"
                rm -f "$DATA_DIR/trims/model_${model_id}_year_${year}.json"
            fi
        else
            echo "failed"
        fi
        sleep 0.1
    done
done

echo "Total trims: $TRIMS_TOTAL" >> "$STATS_FILE"

# 4. Загружаем годы выпуска
echo -e "\n=== Fetching years data ==="
for year in $(seq 2020 2026); do
    echo -n "Fetching data for year $year... "
    curl -s -H "Authorization: Bearer $JWT" \
         -H "Accept: application/json" \
         "https://carapi.app/api/years" \
         -o "$DATA_DIR/years/year_${year}.json"
    echo "done"
done

# 5. Финальная статистика
echo -e "\n=== Final Statistics ==="
echo "End time: $(date)" >> "$STATS_FILE"
echo "Files created:" >> "$STATS_FILE"
echo "- Make files: $(ls -1 "$DATA_DIR/makes" | wc -l)" >> "$STATS_FILE"
echo "- Model files: $(ls -1 "$DATA_DIR/models" | wc -l)" >> "$STATS_FILE"
echo "- Trim files: $(ls -1 "$DATA_DIR/trims" 2>/dev/null | wc -l)" >> "$STATS_FILE"
echo "- Year files: $(ls -1 "$DATA_DIR/years" | wc -l)" >> "$STATS_FILE"
echo "Total size: $(du -sh "$DATA_DIR" | cut -f1)" >> "$STATS_FILE"

cat "$STATS_FILE"

# 6. Создаем архив
echo -e "\nCreating archive..."
cd "$(dirname "$DATA_DIR")"
tar -czf "$(basename "$DATA_DIR").tar.gz" "$(basename "$DATA_DIR")"

echo -e "\n=== SYNC COMPLETED ==="
echo "Data directory: $DATA_DIR"
echo "Archive: $(dirname "$DATA_DIR")/$(basename "$DATA_DIR").tar.gz"

# 7. Сохраняем в базу данных
echo -e "\nSaving to database..."
cd /data/hostel-booking-system/backend
export $(grep -v '^#' .env | xargs)

# Компилируем и запускаем синхронизацию
go run ./cmd/carapi-sync/main.go || echo "Database sync will be done manually"