#!/bin/bash
set -e

# Скрипт для обновления docker-compose.yml для использования Harbor
# Измените путь к вашему файлу docker-compose.yml
DOCKER_COMPOSE_PATH="/data/hostel-booking-system/docker-compose.yml"
DOCKER_COMPOSE_PROD_PATH="/data/hostel-booking-system/docker-compose.prod.yml"
HARBOR_URL="207.180.197.172"

# Функция для создания резервной копии файла
backup_file() {
  local file_path=$1
  cp "$file_path" "${file_path}.bak-$(date +%Y%m%d-%H%M%S)"
  echo "Создана резервная копия файла $file_path"
}

# Обновление docker-compose.yml
if [ -f "$DOCKER_COMPOSE_PATH" ]; then
  echo "Обновление файла docker-compose.yml..."
  backup_file "$DOCKER_COMPOSE_PATH"
  
  # Замена образов в docker-compose.yml
  sed -i "s|image: postgres:15|image: ${HARBOR_URL}/svetu/db/postgres:15|g" "$DOCKER_COMPOSE_PATH"
  sed -i "s|image: opensearchproject/opensearch:2.11.0|image: ${HARBOR_URL}/svetu/opensearch/opensearch:2.11.0|g" "$DOCKER_COMPOSE_PATH"
  sed -i "s|image: opensearchproject/opensearch-dashboards:2.11.0|image: ${HARBOR_URL}/svetu/opensearch/dashboards:2.11.0|g" "$DOCKER_COMPOSE_PATH"
  sed -i "s|image: minio/minio:RELEASE.2023-09-30T07-02-29Z|image: ${HARBOR_URL}/svetu/minio/minio:RELEASE.2023-09-30T07-02-29Z|g" "$DOCKER_COMPOSE_PATH"
  sed -i "s|image: minio/mc|image: ${HARBOR_URL}/svetu/minio/mc:latest|g" "$DOCKER_COMPOSE_PATH"
  sed -i "s|image: migrate/migrate|image: ${HARBOR_URL}/svetu/tools/migrate:latest|g" "$DOCKER_COMPOSE_PATH"
  
  echo "Файл docker-compose.yml обновлен для использования Harbor!"
else
  echo "Файл docker-compose.yml не найден по пути $DOCKER_COMPOSE_PATH"
fi

# Обновление docker-compose.prod.yml
if [ -f "$DOCKER_COMPOSE_PROD_PATH" ]; then
  echo "Обновление файла docker-compose.prod.yml..."
  backup_file "$DOCKER_COMPOSE_PROD_PATH"
  
  # Замена образов в docker-compose.prod.yml
  sed -i "s|image: postgres:15|image: ${HARBOR_URL}/svetu/db/postgres:15|g" "$DOCKER_COMPOSE_PROD_PATH"
  sed -i "s|image: opensearchproject/opensearch:2.11.0|image: ${HARBOR_URL}/svetu/opensearch/opensearch:2.11.0|g" "$DOCKER_COMPOSE_PROD_PATH"
  sed -i "s|image: opensearchproject/opensearch-dashboards:2.11.0|image: ${HARBOR_URL}/svetu/opensearch/dashboards:2.11.0|g" "$DOCKER_COMPOSE_PROD_PATH"
  sed -i "s|image: minio/minio:RELEASE.2023-09-30T07-02-29Z|image: ${HARBOR_URL}/svetu/minio/minio:RELEASE.2023-09-30T07-02-29Z|g" "$DOCKER_COMPOSE_PROD_PATH"
  sed -i "s|image: minio/mc|image: ${HARBOR_URL}/svetu/minio/mc:latest|g" "$DOCKER_COMPOSE_PROD_PATH"
  sed -i "s|image: migrate/migrate|image: ${HARBOR_URL}/svetu/tools/migrate:latest|g" "$DOCKER_COMPOSE_PROD_PATH"
  
  echo "Файл docker-compose.prod.yml обновлен для использования Harbor!"
else
  echo "Файл docker-compose.prod.yml не найден по пути $DOCKER_COMPOSE_PROD_PATH"
fi

# Заменяем блоки build на image для бэкенда и фронтенда
update_build_to_image() {
  local file_path=$1
  
  # Замена для бэкенда
  sed -i '/backend:/,/dockerfile: Dockerfile/ s/build:/image: '"${HARBOR_URL}"'\/svetu\/backend\/api:latest\n    # build:/' "$file_path"
  
  # Замена для фронтенда
  sed -i '/frontend:/,/dockerfile: Dockerfile/ s/build:/image: '"${HARBOR_URL}"'\/svetu\/frontend\/app:latest\n    # build:/' "$file_path"
}

# Применяем замену для обоих файлов
if [ -f "$DOCKER_COMPOSE_PATH" ]; then
  update_build_to_image "$DOCKER_COMPOSE_PATH"
fi

if [ -f "$DOCKER_COMPOSE_PROD_PATH" ]; then
  update_build_to_image "$DOCKER_COMPOSE_PROD_PATH"
fi

echo "Docker-compose файлы обновлены для использования образов из Harbor!"
echo "Не забудьте загрузить образы в Harbor с помощью скрипта migrate_images.sh"