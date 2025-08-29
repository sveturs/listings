#!/bin/bash

# Скрипт синхронизации MinIO из локального окружения в dev.svetu.rs
# Автор: Система автоматизации деплоя
# Дата: $(date)

set -e  # Остановить выполнение при любой ошибке

echo "🔄 Начинаем синхронизацию MinIO: local → dev.svetu.rs"

# Конфигурация
LOCAL_CONTAINER="minio"
DEV_SERVER="root@svetu.rs"
DEV_CONTAINER="minio"
TEMP_DIR="/tmp/minio-sync-$(date +%s)"

# Функция очистки временных файлов
cleanup() {
    echo "🧹 Очистка временных файлов..."
    rm -rf "$TEMP_DIR"
}

# Установить обработчик очистки при выходе
trap cleanup EXIT

# Создать временную директорию
mkdir -p "$TEMP_DIR"

echo "📦 Создание архива локальных данных MinIO..."

# Проверить какие папки есть и создать архив только существующих
FOLDERS_TO_SYNC=""
for folder in listings chat-files review-photos storefronts products; do
    if docker exec "$LOCAL_CONTAINER" sh -c "test -d /data/$folder"; then
        FOLDERS_TO_SYNC="$FOLDERS_TO_SYNC $folder"
    fi
done

echo "📂 Найденные папки для синхронизации: $FOLDERS_TO_SYNC"

# Создать архив только существующих папок
docker exec "$LOCAL_CONTAINER" tar -czf /tmp/minio-data.tar.gz \
    -C /data \
    --exclude='.minio.sys' \
    $FOLDERS_TO_SYNC

# Копировать архив из контейнера
docker cp "$LOCAL_CONTAINER:/tmp/minio-data.tar.gz" "$TEMP_DIR/"

echo "📤 Отправка архива на dev.svetu.rs..."

# Отправить архив на dev сервер
scp "$TEMP_DIR/minio-data.tar.gz" "$DEV_SERVER:/tmp/"

echo "🔄 Распаковка данных в MinIO на dev.svetu.rs..."

# Выполнить команды на dev сервере
ssh "$DEV_SERVER" << 'EOF'
    echo "🛑 Остановка MinIO контейнера..."
    docker stop minio
    
    echo "📁 Резервное копирование существующих данных..."
    docker run --rm -v svetu-dev_minio_data_dev:/data -v /tmp:/backup alpine tar czf /backup/minio-backup-$(date +%Y%m%d-%H%M%S).tar.gz -C /data .
    
    echo "🗑️ Очистка существующих данных..."
    docker run --rm -v svetu-dev_minio_data_dev:/data alpine sh -c "rm -rf /data/* /data/.* || true"
    
    echo "📦 Распаковка новых данных..."
    docker run --rm -v svetu-dev_minio_data_dev:/data -v /tmp:/backup alpine sh -c "cd /data && tar -xzf /backup/minio-data.tar.gz && ls -la /data/"
    
    echo "🔧 Установка правильных прав доступа..."
    docker run --rm -v svetu-dev_minio_data_dev:/data alpine chown -R 1000:1000 /data
    
    echo "▶️ Запуск MinIO контейнера..."
    docker start minio
    
    echo "⏳ Ожидание готовности MinIO..."
    sleep 10
    
    echo "🔄 Перезапуск backend для подключения к обновленному MinIO..."
    docker restart backend-final || docker restart backend-complete || echo "Backend контейнер не найден"
    
    echo "🔍 Проверка статуса MinIO..."
    docker logs minio --tail=20
    
    echo "🧹 Очистка временных файлов на сервере..."
    rm -f /tmp/minio-data.tar.gz
EOF

echo "🎉 Синхронизация завершена успешно!"

echo "📊 Статистика синхронизации:"
echo "   Локальных объявлений: $(docker exec "$LOCAL_CONTAINER" sh -c 'ls -1 /data/listings/ | wc -l')"
echo "   Размер архива: $(du -h "$TEMP_DIR/minio-data.tar.gz" | cut -f1)"

echo "🔗 Проверить результат можно по адресу: https://devs3.svetu.rs"
echo "   Админ панель: https://devs3.svetu.rs (логин: minioadmin)"

echo "✅ Готово! Все файлы MinIO синхронизированы с dev.svetu.rs"