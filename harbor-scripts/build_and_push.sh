#!/bin/bash
set -e

# Скрипт для сборки и отправки образов в Harbor
# Запустите этот скрипт из корня проекта: ./harbor-scripts/build_and_push.sh

# Настройки Harbor
HARBOR_URL="207.180.197.172"
HARBOR_USER="admin"
HARBOR_PASSWORD="SveTu2025"
PROJECT_NAME="svetu"

# Версионирование (используйте дату или git hash)
VERSION=$(date +%Y%m%d-%H%M%S)
# Альтернативный вариант с git hash
# VERSION=$(git rev-parse --short HEAD)

echo "==== Сборка и отправка образов в Harbor (версия: $VERSION) ===="

# Функция для сборки и отправки образа
build_and_push() {
  local context=$1
  local image_name=$2
  local dockerfile=$3
  local full_tag="${HARBOR_URL}/${PROJECT_NAME}/${image_name}:${VERSION}"
  local latest_tag="${HARBOR_URL}/${PROJECT_NAME}/${image_name}:latest"

  echo ""
  echo "Сборка образа: $image_name (контекст: $context)"
  docker build -t "$full_tag" -t "$latest_tag" -f "$dockerfile" "$context"
  
  echo "Отправка образа в Harbor: $full_tag"
  docker push "$full_tag"
  
  echo "Отправка образа в Harbor: $latest_tag"
  docker push "$latest_tag"
  
  echo "✓ Образ $image_name успешно собран и отправлен"
}

# Авторизация в Harbor
echo "Авторизация в Harbor..."
docker login -u $HARBOR_USER -p $HARBOR_PASSWORD $HARBOR_URL

# Сборка и отправка образа бэкенда
build_and_push "./backend" "backend/api" "./backend/Dockerfile"

# Сборка и отправка образа фронтенда
build_and_push "./frontend/hostel-frontend" "frontend/app" "./frontend/hostel-frontend/Dockerfile"

echo ""
echo "==== Все образы успешно собраны и отправлены в Harbor ===="
echo "Версия: $VERSION"
echo "Образы доступны по адресу: http://$HARBOR_URL"
echo ""
echo "Для деплоя в продакшн используйте:"
echo "docker-compose -f docker-compose.prod.yml.harbor up -d"