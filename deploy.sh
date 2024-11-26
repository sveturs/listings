#!/bin/bash

echo "Starting deployment..."

# Переходим в папку проекта
cd /opt/hostel-booking-system

# Создаем необходимые директории
echo "Creating required directories..."
mkdir -p backend/uploads
mkdir -p frontend/hostel-frontend/build
mkdir -p certbot/conf
mkdir -p certbot/www

# Сохраняем важные файлы
echo "Backing up environment files and SSL certificates..."
mkdir -p /tmp/hostel-backup
cp -f backend/.env /tmp/hostel-backup/backend.env 2>/dev/null || true
cp -f frontend/hostel-frontend/.env /tmp/hostel-backup/frontend.env 2>/dev/null || true

# Сохраняем SSL сертификаты
echo "Backing up SSL certificates..."
if [ -d "certbot/conf" ]; then
    cp -r certbot/conf /tmp/hostel-backup/ 2>/dev/null || true
fi

# Сохраняем загруженные изображения
echo "Backing up uploads..."
cp -r backend/uploads /tmp/hostel-backup/ 2>/dev/null || true

# Сбрасываем все локальные изменения
echo "Resetting local changes..."
git reset --hard
git clean -fdx -e "*.env*" -e "uploads/" -e "certbot/"

# Получаем последние изменения из git
echo "Pulling latest changes..."
git pull

# Восстанавливаем env файлы, сертификаты и загрузки
echo "Restoring environment files and SSL certificates..."
cp -f /tmp/hostel-backup/backend.env backend/.env 2>/dev/null || true
cp -f /tmp/hostel-backup/frontend.env frontend/hostel-frontend/.env 2>/dev/null || true

# Восстанавливаем SSL сертификаты
echo "Restoring SSL certificates..."
if [ -d "/tmp/hostel-backup/conf" ]; then
    rm -rf certbot/conf
    cp -r /tmp/hostel-backup/conf certbot/ 2>/dev/null || true
    chmod -R 755 certbot/conf
fi

# Устанавливаем правильные права
echo "Setting permissions..."
chmod -R 755 backend/uploads
chmod -R 755 frontend/hostel-frontend/build
chmod -R 755 certbot/conf

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