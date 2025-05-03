#!/bin/bash
set -e

# Скрипт для миграции всех необходимых образов Docker в Harbor
# для проекта Sve Tu Platform

# Настройки - используем localhost для Docker на том же сервере
HARBOR_LOCAL="127.0.0.1"  # Локальный адрес для Docker
HARBOR_EXTERNAL="207.180.197.172"  # Внешний адрес для docker-compose
HARBOR_USER="admin"
HARBOR_PASSWORD="SveTu2025"
PROJECT_NAME="svetu"

echo "==== Миграция образов Docker в Harbor для Sve Tu Platform ===="

# Авторизация в Harbor
echo "Авторизация в Harbor..."
docker login -u $HARBOR_USER -p $HARBOR_PASSWORD $HARBOR_LOCAL

# Список базовых образов для миграции
declare -A images=(
  ["postgres:15"]="$PROJECT_NAME/db/postgres:15"
  ["opensearchproject/opensearch:2.11.0"]="$PROJECT_NAME/opensearch/opensearch:2.11.0"
  ["opensearchproject/opensearch-dashboards:2.11.0"]="$PROJECT_NAME/opensearch/dashboards:2.11.0"
  ["minio/minio:RELEASE.2023-09-30T07-02-29Z"]="$PROJECT_NAME/minio/minio:RELEASE.2023-09-30T07-02-29Z"
  ["minio/mc:latest"]="$PROJECT_NAME/minio/mc:latest"
  ["migrate/migrate:latest"]="$PROJECT_NAME/tools/migrate:latest"
  ["nginx:latest"]="$PROJECT_NAME/nginx/nginx:latest"
)

# Загрузка, переименование и отправка базовых образов
echo "Миграция базовых образов..."
for src in "${!images[@]}"; do
  dst="${HARBOR_LOCAL}/${images[$src]}"
  echo ""
  echo "Миграция $src -> $dst"
  
  echo "- Загрузка $src..."
  docker pull $src || { echo "Ошибка загрузки $src. Пропускаем."; continue; }
  
  echo "- Переименование в $dst..."
  docker tag $src $dst
  
  echo "- Отправка в Harbor..."
  docker push $dst || { echo "Ошибка отправки $dst. Пропускаем."; continue; }
  
  echo "✓ Завершено"
done

echo ""
echo "==== Образы успешно загружены в Harbor ===="
echo "Все образы теперь доступны в реестре: http://$HARBOR_EXTERNAL"
echo ""
echo "Для обновления docker-compose.yml создайте следующие замены:"
echo ""

for src in "${!images[@]}"; do
  dst="${HARBOR_EXTERNAL}/${images[$src]}"
  echo "image: $src -> image: $dst"
done

echo ""
echo "Для бэкенда и фронтенда замените блоки 'build:' на:"
echo "image: ${HARBOR_EXTERNAL}/${PROJECT_NAME}/backend/api:latest"
echo "image: ${HARBOR_EXTERNAL}/${PROJECT_NAME}/frontend/app:latest"