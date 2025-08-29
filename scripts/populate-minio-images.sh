#!/bin/bash

echo "🔄 Синхронизация изображений из файловой системы MinIO в бакет..."

# MinIO configuration
MINIO_HOST="localhost:9002"
MINIO_ACCESS_KEY="miniodevadmin"
MINIO_SECRET_KEY="h8Puk2qhcadCazC78J1"
BUCKET="listings"

# Configure MinIO client in container
echo "🔧 Настройка MinIO клиента в контейнере..."
docker exec svetu-dev_minio_1 mc alias set local http://localhost:9000 $MINIO_ACCESS_KEY $MINIO_SECRET_KEY

# Check if bucket exists, create if not
if ! docker exec svetu-dev_minio_1 mc ls local/$BUCKET &> /dev/null; then
    echo "📦 Creating bucket: $BUCKET"
    docker exec svetu-dev_minio_1 mc mb local/$BUCKET
fi

# Set bucket policy to public
echo "🔓 Setting bucket policy to public..."
docker exec svetu-dev_minio_1 mc anonymous set download local/$BUCKET

# Sync images from filesystem to bucket
echo "📦 Синхронизация изображений из /data/listings в бакет..."

# Get all directories
DIRS=$(docker exec svetu-dev_minio_1 find /data/listings -mindepth 1 -maxdepth 1 -type d | sort)
TOTAL_DIRS=$(echo "$DIRS" | wc -l)
CURRENT=0

for dir in $DIRS; do
    CURRENT=$((CURRENT + 1))
    LISTING_ID=$(basename $dir)
    
    echo "  📸 Обработка $CURRENT/$TOTAL_DIRS: объявление #$LISTING_ID"
    
    # Upload each image file to MinIO bucket
    docker exec svetu-dev_minio_1 sh -c "
        for file in /data/listings/$LISTING_ID/*.jpg /data/listings/$LISTING_ID/*.jpeg /data/listings/$LISTING_ID/*.png 2>/dev/null; do
            if [ -f \"\$file\" ]; then
                filename=\$(basename \"\$file\")
                mc cp -q \"\$file\" local/$BUCKET/$LISTING_ID/\$filename 2>/dev/null || true
            fi
        done
    "
done

# Verify upload
echo ""
echo "📊 Проверка загруженных изображений..."
UPLOADED_COUNT=$(docker exec svetu-dev_minio_1 mc ls --recursive local/$BUCKET/ | grep -c "\.jpg\|\.jpeg\|\.png" || echo "0")
echo "✅ Загружено изображений в MinIO: $UPLOADED_COUNT"

# Check specific listing
echo ""
echo "🔍 Проверка объявления #268:"
docker exec svetu-dev_minio_1 mc ls local/$BUCKET/268/

# Test HTTP access
echo ""
echo "🌐 Проверка доступности через HTTPS:"
for file in main.jpg image2.jpg image3.jpg; do
    if curl -s -o /dev/null -w "%{http_code}" https://devs3.svetu.rs/listings/268/$file | grep -q "200"; then
        echo "  ✅ https://devs3.svetu.rs/listings/268/$file - OK"
    else
        echo "  ❌ https://devs3.svetu.rs/listings/268/$file - ОШИБКА"
    fi
done

echo ""
echo "🎉 Готово! Изображения доступны по адресам:"
echo "   https://devs3.svetu.rs/listings/{id}/main.jpg"
echo "   https://devs3.svetu.rs/listings/{id}/image2.jpg"
echo "   и т.д."