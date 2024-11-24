#!/bin/bash

echo "Starting deployment..."

# Переходим в папку проекта
cd /opt/hostel-booking-system

# Сохраняем важные файлы
echo "Backing up environment files..."
mkdir -p /tmp/hostel-backup
cp -f backend/.env* /tmp/hostel-backup/ 2>/dev/null || true
cp -f frontend/hostel-frontend/.env* /tmp/hostel-backup/ 2>/dev/null || true

# Сохраняем загруженные изображения
echo "Backing up uploads..."
cp -r backend/uploads /tmp/hostel-backup/ 2>/dev/null || true

# Сбрасываем все локальные изменения
echo "Resetting local changes..."
git reset --hard
git clean -fdx -e "*.env*" -e "uploads/"

# Получаем последние изменения из git
echo "Pulling latest changes..."
git pull

# Восстанавливаем env файлы и загрузки
echo "Restoring environment files and uploads..."
cp -f /tmp/hostel-backup/.env* backend/ 2>/dev/null || true
cp -f /tmp/hostel-backup/.env* frontend/hostel-frontend/ 2>/dev/null || true
mkdir -p backend/uploads
cp -r /tmp/hostel-backup/uploads/* backend/uploads/ 2>/dev/null || true

# Удаляем временные файлы
rm -rf /tmp/hostel-backup

# Собираем фронтенд
echo "Building frontend..."
cd frontend/hostel-frontend
NODE_ENV=production docker run -v $(pwd):/app -w /app node:18 sh -c "npm install --legacy-peer-deps && npm install @babel/plugin-proposal-private-property-in-object && npm run build"

# Возвращаемся в корневую папку
cd ../..

# Перезапускаем контейнеры
echo "Restarting containers..."
docker-compose -f docker-compose.prod.yml down
docker-compose -f docker-compose.prod.yml up -d --build

# Проверяем статус
echo "Checking status..."
docker-compose -f docker-compose.prod.yml ps

echo "Deployment completed!"