#!/bin/bash
set -e # Останавливаем выполнение при ошибках
echo "Начинаем деплой только backend..."
cd /opt/hostel-booking-system

# Настраиваем git pull strategy
git config pull.rebase false

# Сохраняем .env файл бэкенда
mkdir -p /tmp/hostel-backup
cp -f backend/.env /tmp/hostel-backup/backend.env 2>/dev/null || true

# Сохраняем загруженные файлы, если они есть
mkdir -p /tmp/hostel-backup/uploads
cp -r backend/uploads/* /tmp/hostel-backup/uploads/ 2>/dev/null || true

# Получаем обновления для бэкенда
echo "Получаем обновления из git..."
git fetch origin
git checkout origin/main -- backend
# Сбрасываем изменения для .env и директории uploads, чтобы они не перезаписались
git reset HEAD backend/.env 2>/dev/null || true

# Восстанавливаем .env файл
cp -f /tmp/hostel-backup/backend.env backend/.env 2>/dev/null || true

# Проверяем, существует ли директория uploads в бэкенде
if [ ! -d "backend/uploads" ]; then
  mkdir -p backend/uploads
fi

# Восстанавливаем файлы из uploads, если они есть
cp -r /tmp/hostel-backup/uploads/* backend/uploads/ 2>/dev/null || true

# Перезапускаем только backend
echo "Перезапускаем backend..."
docker-compose -f docker-compose.prod.yml up --build -d backend

echo "Деплой backend завершен!"