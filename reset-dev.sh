#!/bin/bash

# Скрипт для полного перезапуска dev-окружения SveTu

echo "🔄 Полный перезапуск dev-окружения..."

# Загружаем переменные окружения
set -a
source .env
set +a

# Останавливаем все контейнеры
echo "📦 Останавливаем контейнеры..."
docker-compose -f docker-compose.dev.yml down

# Удаляем volumes для полной очистки (опционально)
read -p "❓ Удалить все данные (volumes)? [y/N]: " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "🗑️ Удаляем volumes..."
    docker volume rm svetu-dev_postgres_data_dev 2>/dev/null || true
    docker volume rm svetu-dev_redis_data_dev 2>/dev/null || true
    docker volume rm svetu-dev_opensearch-data_dev 2>/dev/null || true
    docker volume rm svetu-dev_minio_data_dev 2>/dev/null || true
    rm -rf data/minio_dev/* 2>/dev/null || true
    rm -rf backend/uploads/* 2>/dev/null || true
fi

# Пересобираем контейнеры
echo "🔨 Пересобираем контейнеры..."
docker-compose -f docker-compose.dev.yml build

# Запускаем все сервисы
echo "🚀 Запускаем сервисы..."
docker-compose -f docker-compose.dev.yml up -d

# Ждем готовности базы данных
echo "⏳ Ждем готовности PostgreSQL..."
until docker exec svetu-dev_db_1 pg_isready -U postgres > /dev/null 2>&1; do
    sleep 2
done
echo "✅ PostgreSQL готов"

# Ждем готовности backend
echo "⏳ Ждем готовности backend..."
until curl -s http://localhost:3002/health > /dev/null 2>&1; do
    sleep 2
done
echo "✅ Backend готов"

# Реиндексация OpenSearch (если нужно)
echo "🔍 Реиндексация OpenSearch..."
docker exec svetu-dev_backend_1 ./reindex || echo "⚠️ Реиндексация пропущена"

# Показываем статус
echo ""
echo "✅ Dev-окружение запущено!"
echo ""
echo "📍 Доступные сервисы:"
echo "   Frontend: https://dev.svetu.rs"
echo "   Backend API: https://devapi.svetu.rs"
echo "   Swagger: https://devapi.svetu.rs/swagger/index.html"
echo "   MinIO S3: https://devs3.svetu.rs"
echo "   MinIO Console: http://svetu.rs:9003"
echo ""
echo "📊 Статус контейнеров:"
docker-compose -f docker-compose.dev.yml ps

echo ""
echo "💡 Полезные команды:"
echo "   Логи: docker-compose -f docker-compose.dev.yml logs -f"
echo "   Перезапуск backend: docker-compose -f docker-compose.dev.yml restart backend"
echo "   Очистка Redis: docker exec svetu-dev_redis_1 redis-cli FLUSHALL"