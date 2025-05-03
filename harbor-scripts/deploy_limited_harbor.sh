#!/bin/bash
set -e # Останавливаем выполнение при ошибках
echo "Начинаем частичный деплой с использованием Harbor..."
cd /opt/hostel-booking-system

# Останавливаем текущие службы
echo "Останавливаем текущие сервисы..."
docker-compose -f docker-compose.prod.yml down

# Авторизация в Harbor
echo "Авторизация в Harbor..."
docker login -u admin -p SveTu2025 harbor.svetu.rs

# Загрузка доступных образов из Harbor
echo "Загрузка образов из Harbor..."
docker pull harbor.svetu.rs/svetu/db/postgres:15
docker pull harbor.svetu.rs/svetu/minio/minio:RELEASE.2023-09-30T07-02-29Z
docker pull harbor.svetu.rs/svetu/minio/mc:latest
docker pull harbor.svetu.rs/svetu/tools/migrate:latest
docker pull harbor.svetu.rs/svetu/backend/api:latest

# Создаем необходимые директории
mkdir -p /tmp/hostel-backup/db
mkdir -p backend/uploads
mkdir -p /opt/hostel-data/{opensearch,uploads,minio}

# Делаем бэкап базы данных
echo "Создание бэкапа базы данных..."
if docker ps -a | grep -q hostel_db; then
  BACKUP_FILE="/tmp/hostel-backup/db/backup_$(date +%Y%m%d_%H%M%S).sql"
  docker exec -t hostel_db pg_dumpall -U postgres > "$BACKUP_FILE"
  if [ $? -eq 0 ]; then
    echo "Бэкап базы данных создан в $BACKUP_FILE"
  else
    echo "Ошибка создания бэкапа базы данных, но продолжаем деплой"
  fi
else
  echo "База данных не найдена, пропускаем создание бэкапа"
fi

# Запускаем базу данных из Harbor
echo "Запускаем базу данных из Harbor..."
docker-compose -f docker-compose.prod.yml.harbor up -d db

# Проверяем базу данных
echo "Проверяем готовность базы данных..."
RETRY_COUNT=30
for i in $(seq 1 $RETRY_COUNT); do
  if docker-compose -f docker-compose.prod.yml.harbor exec -T db pg_isready -U postgres > /dev/null 2>&1; then
    echo "База данных готова!"
    break
  fi
  echo "Ожидаем готовность базы данных... Попытка $i/$RETRY_COUNT"
  sleep 2
done

if ! docker-compose -f docker-compose.prod.yml.harbor exec -T db pg_isready -U postgres > /dev/null 2>&1; then
  echo "Не удалось запустить базу данных"
  exit 1
fi

# Запускаем миграции
echo "Запускаем миграции..."
docker-compose -f docker-compose.prod.yml.harbor up -d migrate

# Восстанавливаем данные из бэкапа, если он есть
if [ -n "$(ls -t /tmp/hostel-backup/db/*.sql 2>/dev/null | head -1)" ]; then
  LATEST_BACKUP=$(ls -t /tmp/hostel-backup/db/*.sql | head -1)
  echo "Восстанавливаем базу данных из $LATEST_BACKUP..."
  cat "$LATEST_BACKUP" | docker-compose -f docker-compose.prod.yml.harbor exec -T db psql -U postgres
  if [ $? -eq 0 ]; then
    echo "База данных успешно восстановлена"
  else
    echo "Ошибка восстановления базы данных, но продолжаем деплой"
  fi
else
  echo "Бэкап базы данных не найден, пропускаем восстановление"
fi

# Запускаем minio и createbuckets
echo "Запускаем minio..."
docker-compose -f docker-compose.prod.yml.harbor up -d minio createbuckets

# Запускаем backend с использованием образа из Harbor
echo "Запускаем backend..."
docker-compose -f docker-compose.prod.yml.harbor up -d backend

# Запускаем остальные сервисы из стандартного файла
echo "Запускаем остальные сервисы..."
docker-compose -f docker-compose.prod.yml up -d opensearch mailserver mail-webui nginx

echo "Частичный деплой с использованием Harbor завершен!"
echo "Используются образы Harbor для: db, backend, minio, createbuckets, migrate"
echo "Используются стандартные образы для: opensearch, mailserver, mail-webui, nginx"