#!/bin/bash

# Скрипт для сохранения сырых данных из CarAPI в JSON файлы
# Это резервный вариант на случай проблем с основной синхронизацией

set -e

echo "=== CarAPI Raw Data Backup Script ==="
echo "Starting at: $(date)"

# Проверяем токен
if [ -z "$CARAPI_TOKEN" ]; then
    echo "ERROR: CARAPI_TOKEN is not set!"
    exit 1
fi

# Создаем директорию для бекапов
BACKUP_DIR="/data/hostel-booking-system/backend/data/carapi-backup-$(date +%Y%m%d)"
mkdir -p "$BACKUP_DIR"

echo "Saving data to: $BACKUP_DIR"

# Функция для сохранения данных
save_data() {
    local endpoint=$1
    local filename=$2
    local url="https://carapi.app/api${endpoint}"
    
    echo "Fetching: $url"
    curl -s -H "Authorization: Bearer $CARAPI_TOKEN" \
         -H "Accept: application/json" \
         "$url" > "$BACKUP_DIR/$filename"
    
    # Проверяем что получили валидный JSON
    if jq empty "$BACKUP_DIR/$filename" 2>/dev/null; then
        echo "✓ Saved $filename ($(wc -c < "$BACKUP_DIR/$filename") bytes)"
    else
        echo "✗ ERROR: Invalid JSON in $filename"
        return 1
    fi
    
    # Небольшая задержка чтобы не превысить rate limit
    sleep 0.5
}

# 1. Сохраняем все марки
echo "Fetching all makes..."
save_data "/makes" "makes.json"

# 2. Сохраняем модели для популярных марок
echo "Fetching models for popular makes..."

# Получаем ID популярных марок из сохраненного файла
MAKE_IDS=$(cat "$BACKUP_DIR/makes.json" | jq -r '.[] | select(.name | test("Volkswagen|Toyota|Škoda|Fiat|Ford|Opel|Peugeot|Renault|Hyundai|Kia|Mercedes|BMW|Audi"; "i")) | .id' | head -20)

mkdir -p "$BACKUP_DIR/models"
for make_id in $MAKE_IDS; do
    make_name=$(cat "$BACKUP_DIR/makes.json" | jq -r ".[] | select(.id == $make_id) | .name")
    echo "Fetching models for $make_name (ID: $make_id)..."
    save_data "/models?make_id=$make_id" "models/make_${make_id}.json"
done

# 3. Сохраняем комплектации для последних лет популярных моделей
echo "Fetching trims for recent years..."

mkdir -p "$BACKUP_DIR/trims"
CURRENT_YEAR=$(date +%Y)

# Берем несколько популярных моделей для примера
MODEL_IDS=$(cat "$BACKUP_DIR/models/"*.json | jq -r '.[].id' | sort -u | head -50)

for model_id in $MODEL_IDS; do
    for year in $(seq $((CURRENT_YEAR-3)) $CURRENT_YEAR); do
        echo "Fetching trims for model $model_id, year $year..."
        save_data "/trims?model_id=$model_id&year=$year" "trims/model_${model_id}_year_${year}.json" || true
    done
done

# 4. Тестируем VIN декодер на нескольких примерах
echo "Testing VIN decoder..."
mkdir -p "$BACKUP_DIR/vin"

# Примеры VIN для тестирования
SAMPLE_VINS=(
    "WVWZZZ1JZ3W386752"  # VW Golf
    "JHMCM56557C404453"  # Honda Accord
    "1FTFW1ET5DFC10312"  # Ford F-150
)

for vin in "${SAMPLE_VINS[@]}"; do
    echo "Decoding VIN: $vin"
    save_data "/vin/decode?vin=$vin" "vin/${vin}.json" || true
done

# 5. Создаем сводный отчет
echo "Creating summary report..."
cat > "$BACKUP_DIR/summary.txt" << EOF
CarAPI Data Backup Summary
========================
Date: $(date)
Directory: $BACKUP_DIR

Files created:
- makes.json: $(cat "$BACKUP_DIR/makes.json" | jq '. | length') makes
- models/: $(ls -1 "$BACKUP_DIR/models/" | wc -l) files
- trims/: $(ls -1 "$BACKUP_DIR/trims/" | wc -l) files
- vin/: $(ls -1 "$BACKUP_DIR/vin/" 2>/dev/null | wc -l) files

Total size: $(du -sh "$BACKUP_DIR" | cut -f1)
EOF

cat "$BACKUP_DIR/summary.txt"

# 6. Создаем архив для долгосрочного хранения
echo "Creating archive..."
cd "$(dirname "$BACKUP_DIR")"
tar -czf "$(basename "$BACKUP_DIR").tar.gz" "$(basename "$BACKUP_DIR")"

echo "=== Backup completed at: $(date) ==="
echo "Data saved to: $BACKUP_DIR"
echo "Archive created: $(dirname "$BACKUP_DIR")/$(basename "$BACKUP_DIR").tar.gz"