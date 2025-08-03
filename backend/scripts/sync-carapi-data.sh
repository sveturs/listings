#!/bin/bash

# Скрипт для полной синхронизации данных из CarAPI
# ВАЖНО: Запустить до истечения подписки!

set -e

echo "=== CarAPI Data Sync Script ==="
echo "Starting at: $(date)"

# Проверяем наличие токена
if [ -z "$CARAPI_TOKEN" ]; then
    echo "ERROR: CARAPI_TOKEN environment variable is not set!"
    echo "Please set it: export CARAPI_TOKEN='your-token-here'"
    exit 1
fi

# Путь к корню проекта
PROJECT_ROOT="/data/hostel-booking-system"
cd "$PROJECT_ROOT/backend"

# Применяем миграцию если еще не применена
echo "Applying database migrations..."
./bin/migrator up

# Компилируем скрипт синхронизации
echo "Building sync tool..."
go build -o /tmp/carapi-sync ./cmd/carapi-sync/main.go

# Запускаем синхронизацию
echo "Starting data synchronization..."
echo "This may take a while due to API rate limits..."

# Экспортируем переменные окружения для подключения к БД
export $(grep -v '^#' .env | xargs)

# Запускаем синхронизацию с логированием
/tmp/carapi-sync 2>&1 | tee /tmp/carapi-sync-$(date +%Y%m%d-%H%M%S).log

echo "=== Sync completed at: $(date) ==="

# Показываем статистику
echo "Checking sync results..."
psql "$DATABASE_URL" << EOF
-- Статистика по маркам
SELECT 'Makes synced' as metric, COUNT(*) as count 
FROM car_makes WHERE external_id IS NOT NULL
UNION ALL
-- Статистика по моделям  
SELECT 'Models synced', COUNT(*) 
FROM car_models WHERE external_id IS NOT NULL
UNION ALL
-- Статистика по комплектациям
SELECT 'Trims synced', COUNT(*) 
FROM car_trims
UNION ALL
-- Использование API сегодня
SELECT 'API requests today', COALESCE(requests_count, 0)
FROM carapi_usage WHERE date = CURRENT_DATE;
EOF

echo "Logs saved to: /tmp/carapi-sync-*.log"