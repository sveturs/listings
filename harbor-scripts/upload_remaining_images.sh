#!/bin/bash
set -e

# Скрипт для загрузки оставшихся образов в Harbor
HARBOR_URL="harbor.svetu.rs"
HARBOR_USER="admin"
HARBOR_PASSWORD="SveTu2025"
PROJECT_NAME="svetu"

echo "==== Загрузка оставшихся образов в Harbor ===="

# Авторизация в Harbor
echo "Авторизация в Harbor..."
docker login -u $HARBOR_USER -p $HARBOR_PASSWORD $HARBOR_URL

# Список оставшихся образов
declare -A images=(
  ["mailserver/docker-mailserver:latest"]="$PROJECT_NAME/mail/server:latest"
  ["roundcube/roundcubemail:latest"]="$PROJECT_NAME/mail/webui:latest"
  ["certbot/certbot:latest"]="$PROJECT_NAME/tools/certbot:latest"
)

# Загрузка, переименование и отправка образов
for src in "${!images[@]}"; do
  dst="${HARBOR_URL}/${images[$src]}"
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
echo "==== Все образы успешно загружены в Harbor ===="
echo "Обновите docker-compose.prod.yml.harbor для использования следующих образов:"
echo ""

for src in "${!images[@]}"; do
  dst="${HARBOR_URL}/${images[$src]}"
  echo "image: $src -> image: $dst"
done