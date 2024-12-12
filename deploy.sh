#!/bin/bash
set -e  # Останавливаем выполнение при ошибках

echo "Starting deployment..."

cd /opt/hostel-booking-system

# Настраиваем git pull strategy
git config pull.rebase false

# Создаем необходимые директории
mkdir -p backend/uploads
mkdir -p frontend/hostel-frontend/build
mkdir -p certbot/conf
mkdir -p certbot/www

# Сохраняем важные файлы
mkdir -p /tmp/hostel-backup
cp -f backend/.env /tmp/hostel-backup/backend.env 2>/dev/null || true
cp -f frontend/hostel-frontend/.env /tmp/hostel-backup/frontend.env 2>/dev/null || true

# Сохраняем SSL сертификаты
if [ -d "certbot/conf" ]; then
    cp -r certbot/conf /tmp/hostel-backup/ 2>/dev/null || true
fi

# Сохраняем загруженные изображения
cp -r backend/uploads /tmp/hostel-backup/ 2>/dev/null || true

# Обеспечиваем чистое состояние git
git fetch origin
git reset --hard origin/main
git clean -fdx -e "*.env*" -e "uploads/" -e "certbot/"

# Восстанавливаем файлы
cp -f /tmp/hostel-backup/backend.env backend/.env 2>/dev/null || true
cp -f /tmp/hostel-backup/frontend.env frontend/hostel-frontend/.env 2>/dev/null || true
if [ -d "/tmp/hostel-backup/conf" ]; then
    rm -rf certbot/conf
    cp -r /tmp/hostel-backup/conf certbot/ 2>/dev/null || true
fi

# Удаляем старые образы
docker image prune -f

# Очищаем сети и осиротевшие контейнеры
echo "Cleaning up orphan containers and networks..."
docker-compose -f docker-compose.prod.yml down -v --remove-orphans || true
docker network prune -f || true

# Собираем фронтенд
echo "Building frontend..."
cd frontend/hostel-frontend
NODE_ENV=production docker run -v $(pwd):/app -w /app node:18 sh -c "npm install && npm run build"
cd ../..

# Запускаем только базу данных
echo "Starting database..."
docker-compose -f docker-compose.prod.yml up --build -d db

# Проверяем базу данных
echo "Checking database readiness..."
RETRY_COUNT=30
for i in $(seq 1 $RETRY_COUNT); do
    if docker-compose -f docker-compose.prod.yml exec -T db pg_isready -U postgres > /dev/null 2>&1; then
        echo "Database is ready!"
        break
    fi
    echo "Waiting for database to be ready... Attempt $i/$RETRY_COUNT"
    sleep 2
done

if ! docker-compose -f docker-compose.prod.yml exec -T db pg_isready -U postgres > /dev/null 2>&1; then
    echo "Database failed to start"
    exit 1
fi

# Запускаем миграции
echo "Running migrations..."
docker run --rm --network hostel-booking-system_hostel_network -v $(pwd)/backend/migrations:/migrations migrate/migrate -path=/migrations/ -database="postgres://postgres:password@db:5432/hostel_db?sslmode=disable" up

# Запускаем остальные сервисы
echo "Starting services..."
docker-compose -f docker-compose.prod.yml up --build -d

echo "Deployment completed!"
