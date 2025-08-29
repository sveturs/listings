#!/bin/bash

# Простая синхронизация MinIO через rsync
set -e

echo "🔄 Синхронизация MinIO: rsync метод"

LOCAL_CONTAINER="minio"
DEV_SERVER="root@svetu.rs"
DEV_CONTAINER="svetu-dev_minio_1"

# Создать временные директории
LOCAL_TEMP="/tmp/minio-local-$(date +%s)"
DEV_TEMP="/tmp/minio-dev-sync"

# Функция очистки
cleanup() {
    echo "🧹 Очистка..."
    rm -rf "$LOCAL_TEMP"
    ssh "$DEV_SERVER" "rm -rf $DEV_TEMP"
}
trap cleanup EXIT

echo "📦 Экспорт данных из локального MinIO..."
mkdir -p "$LOCAL_TEMP"

# Скопировать данные из локального контейнера
docker cp "$LOCAL_CONTAINER:/data/listings" "$LOCAL_TEMP/" 2>/dev/null || echo "⚠️ listings пропущен"
docker cp "$LOCAL_CONTAINER:/data/chat-files" "$LOCAL_TEMP/" 2>/dev/null || echo "⚠️ chat-files пропущен"
docker cp "$LOCAL_CONTAINER:/data/review-photos" "$LOCAL_TEMP/" 2>/dev/null || echo "⚠️ review-photos пропущен"

echo "📁 Локальные данные для синхронизации:"
ls -la "$LOCAL_TEMP/"

echo "📤 Отправка данных на dev.svetu.rs..."

# Синхронизация с dev сервером
rsync -avz --delete "$LOCAL_TEMP/" "$DEV_SERVER:$DEV_TEMP/"

echo "🔄 Импорт данных в MinIO на dev.svetu.rs..."

# Остановить MinIO, скопировать данные, запустить
ssh "$DEV_SERVER" << EOF
    echo "🛑 Остановка MinIO..."
    docker stop svetu-dev_minio_1
    
    echo "📁 Резервное копирование..."
    docker run --rm -v svetu-dev_minio-data:/data -v /tmp:/backup alpine tar czf /backup/minio-backup-\$(date +%Y%m%d-%H%M%S).tar.gz -C /data . || true
    
    echo "🗑️ Очистка старых данных..."
    docker run --rm -v svetu-dev_minio-data:/data alpine sh -c "rm -rf /data/listings /data/chat-files /data/review-photos" || true
    
    echo "📦 Копирование новых данных..."
    for folder in listings chat-files review-photos; do
        if [ -d "$DEV_TEMP/\$folder" ]; then
            echo "  Копирование \$folder..."
            docker run --rm -v svetu-dev_minio-data:/data -v $DEV_TEMP:/source alpine cp -r "/source/\$folder" /data/ || echo "   ⚠️ Ошибка копирования \$folder"
        fi
    done
    
    echo "🔧 Установка прав доступа..."
    docker run --rm -v svetu-dev_minio-data:/data alpine chown -R 1000:1000 /data
    
    echo "▶️ Запуск MinIO..."
    docker start svetu-dev_minio_1
    
    echo "⏳ Ожидание готовности..."
    sleep 10
    
    echo "📊 Проверка результата..."
    docker exec svetu-dev_minio_1 sh -c "ls -la /data/" || echo "Ошибка проверки"
    docker exec svetu-dev_minio_1 sh -c "ls -la /data/listings/ | head -5" || echo "listings пуст"
EOF

echo "✅ Синхронизация через rsync завершена!"