#!/bin/bash

# Скрипт для исправления путей изображений в БД
# Синхронизирует записи в БД с реальными файлами в MinIO

echo "🔍 Начинаем анализ и исправление путей изображений..."

# Параметры подключения к БД
DB_HOST="localhost"
DB_PORT="5433"
DB_USER="svetu_dev_user"
DB_PASS="svetu_dev_user"
DB_NAME="svetu_dev_db"

# Получаем список объявлений с изображениями
echo "📋 Получаем список объявлений..."

LISTINGS=$(docker exec svetu-dev_db_1 sh -c "PGPASSWORD=$DB_PASS psql -U $DB_USER -d $DB_NAME -t -c \"SELECT DISTINCT listing_id FROM marketplace_images ORDER BY listing_id;\"")

echo "📊 Найдено объявлений с изображениями: $(echo "$LISTINGS" | wc -l)"

# Обрабатываем каждое объявление
for listing_id in $LISTINGS; do
    listing_id=$(echo $listing_id | tr -d ' ')
    
    if [ -z "$listing_id" ]; then
        continue
    fi
    
    echo ""
    echo "🔄 Обрабатываем объявление #$listing_id..."
    
    # Получаем файлы из MinIO
    MINIO_FILES=$(docker exec svetu-dev_minio_1 ls /data/listings/$listing_id/ 2>/dev/null | grep -E '\.(jpg|jpeg|png)$' | awk '{print $NF}')
    
    if [ -z "$MINIO_FILES" ]; then
        echo "  ⚠️  Папка /data/listings/$listing_id/ не найдена в MinIO"
        continue
    fi
    
    # Получаем записи из БД
    DB_FILES=$(docker exec svetu-dev_db_1 sh -c "PGPASSWORD=$DB_PASS psql -U $DB_USER -d $DB_NAME -t -c \"SELECT file_name FROM marketplace_images WHERE listing_id = $listing_id;\"")
    
    echo "  📁 Файлы в MinIO: $(echo "$MINIO_FILES" | wc -l)"
    echo "  📝 Записи в БД: $(echo "$DB_FILES" | wc -l)"
    
    # Проверяем соответствие
    DB_FILE=$(echo "$DB_FILES" | head -1 | tr -d ' ')
    MINIO_FILE=$(echo "$MINIO_FILES" | head -1)
    
    if [ "$DB_FILE" != "$MINIO_FILE" ]; then
        echo "  ❌ Несоответствие: БД='$DB_FILE', MinIO='$MINIO_FILE'"
        
        # Удаляем старые записи
        echo "  🗑️  Удаляем старые записи из БД..."
        docker exec svetu-dev_db_1 sh -c "PGPASSWORD=$DB_PASS psql -U $DB_USER -d $DB_NAME -c \"DELETE FROM marketplace_images WHERE listing_id = $listing_id;\""
        
        # Добавляем новые записи для каждого файла в MinIO
        IS_FIRST=true
        for file in $MINIO_FILES; do
            IS_MAIN="false"
            if [ "$IS_FIRST" = true ] || [ "$file" = "main.jpg" ]; then
                IS_MAIN="true"
                IS_FIRST=false
            fi
            
            echo "  ➕ Добавляем запись: $file (is_main=$IS_MAIN)"
            
            # Получаем размер файла
            FILE_SIZE=$(docker exec svetu-dev_minio_1 stat /data/listings/$listing_id/$file 2>/dev/null | grep "Size:" | awk '{print $2}')
            if [ -z "$FILE_SIZE" ]; then
                FILE_SIZE="0"
            fi
            
            # Вставляем запись в БД
            docker exec svetu-dev_db_1 sh -c "PGPASSWORD=$DB_PASS psql -U $DB_USER -d $DB_NAME -c \"
                INSERT INTO marketplace_images (
                    listing_id, 
                    file_path, 
                    file_name, 
                    file_size, 
                    content_type, 
                    is_main, 
                    storage_type, 
                    storage_bucket, 
                    public_url
                ) VALUES (
                    $listing_id,
                    '$listing_id/$file',
                    '$file',
                    $FILE_SIZE,
                    'image/jpeg',
                    $IS_MAIN,
                    'minio',
                    'listings',
                    '/listings/$listing_id/$file'
                );
            \""
        done
        
        echo "  ✅ Исправлено!"
    else
        echo "  ✅ Пути совпадают"
    fi
done

echo ""
echo "🎉 Обработка завершена!"
echo ""
echo "📌 Проверка объявления #268:"
docker exec svetu-dev_db_1 sh -c "PGPASSWORD=$DB_PASS psql -U $DB_USER -d $DB_NAME -c \"SELECT file_name, public_url FROM marketplace_images WHERE listing_id = 268;\""