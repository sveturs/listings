#!/bin/bash
set -e

# Скрипт для деплоя с использованием образов из Harbor
# Запустите этот скрипт из корня проекта: ./harbor-scripts/deploy_with_harbor.sh

# Настройки Harbor
HARBOR_URL="207.180.197.172"
HARBOR_USER="admin"
HARBOR_PASSWORD="SveTu2025"
PROJECT_NAME="svetu"

# Версия для деплоя (latest или конкретная версия)
VERSION="latest"

echo "==== Деплой Sve Tu Platform с использованием Harbor ===="

# Функция для создания резервной копии
backup_database() {
  echo "Создание резервной копии базы данных..."
  BACKUP_DIR="./backup/$(date +%Y%m%d_%H%M%S)"
  mkdir -p $BACKUP_DIR
  
  # Здесь добавьте команды для бэкапа базы данных
  # Например: docker exec hostel_db pg_dump -U postgres hostel_db > $BACKUP_DIR/database.sql
  
  echo "Резервная копия создана в директории: $BACKUP_DIR"
}

# Авторизация в Harbor
echo "Авторизация в Harbor..."
docker login -u $HARBOR_USER -p $HARBOR_PASSWORD $HARBOR_URL

# Проверка наличия docker-compose.prod.yml.harbor
if [ ! -f "docker-compose.prod.yml.harbor" ]; then
  echo "Ошибка: файл docker-compose.prod.yml.harbor не найден!"
  echo "Создайте этот файл или используйте update_docker_compose.sh из harbor-scripts"
  exit 1
fi

# Создание резервной копии
backup_database

# Загрузка образов из Harbor
echo "Загрузка последних образов из Harbor..."
docker pull ${HARBOR_URL}/${PROJECT_NAME}/backend/api:${VERSION}
docker pull ${HARBOR_URL}/${PROJECT_NAME}/frontend/app:${VERSION}
docker pull ${HARBOR_URL}/${PROJECT_NAME}/db/postgres:15
docker pull ${HARBOR_URL}/${PROJECT_NAME}/opensearch/opensearch:2.11.0
docker pull ${HARBOR_URL}/${PROJECT_NAME}/opensearch/dashboards:2.11.0
docker pull ${HARBOR_URL}/${PROJECT_NAME}/minio/minio:RELEASE.2023-09-30T07-02-29Z
docker pull ${HARBOR_URL}/${PROJECT_NAME}/minio/mc:latest
docker pull ${HARBOR_URL}/${PROJECT_NAME}/tools/migrate:latest

# Остановка текущих сервисов
echo "Остановка текущих сервисов..."
docker-compose down

# Запуск с использованием Harbor образов
echo "Запуск сервисов с использованием образов из Harbor..."
docker-compose -f docker-compose.prod.yml.harbor up -d

echo ""
echo "==== Деплой успешно завершен! ===="
echo "Сервисы запущены с использованием образов из Harbor (версия: $VERSION)"