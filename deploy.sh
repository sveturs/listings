#!/bin/bash

echo "Starting deployment..."

# Переходим в папку проекта
cd /opt/hostel-booking-system

# Создаем необходимые директории
echo "Creating required directories..."
mkdir -p backend/uploads
mkdir -p frontend/hostel-frontend/build

# Сохраняем важные файлы
echo "Backing up environment files..."
mkdir -p /tmp/hostel-backup
cp -f backend/.env /tmp/hostel-backup/backend.env 2>/dev/null || true
cp -f frontend/hostel-frontend/.env /tmp/hostel-backup/frontend.env 2>/dev/null || true


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
echo "Restoring environment files..."
cp -f /tmp/hostel-backup/backend.env backend/.env 2>/dev/null || true
cp -f /tmp/hostel-backup/frontend.env frontend/hostel-frontend/.env 2>/dev/null || true

# Устанавливаем правильные права
echo "Setting permissions..."
chmod -R 755 backend/uploads
chmod -R 755 frontend/hostel-frontend/build

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

# Проверяем статус и логи
echo "Checking container status and logs..."
docker-compose -f docker-compose.prod.yml ps
docker logs hostel_nginx
docker logs hostel_backend

echo "Deployment completed!"