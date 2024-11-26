#!/bin/bash

set -e  # Останавливаем выполнение при ошибках

echo "Starting deployment..."

# Переходим в папку проекта
cd /opt/hostel-booking-system

# Создаем необходимые директории
echo "Creating required directories..."
mkdir -p backend/uploads
mkdir -p frontend/hostel-frontend/build
mkdir -p certbot/conf
mkdir -p certbot/www

# Функция для ожидания готовности базы данных
wait_for_postgres() {
    echo "Waiting for PostgreSQL to start..."
    for i in {1..30}; do
        if docker exec hostel_db pg_isready -U postgres > /dev/null 2>&1; then
            echo "PostgreSQL is ready!"
            sleep 5  # Даем дополнительное время для полной инициализации
            return 0
        fi
        echo "Waiting... ($i/30)"
        sleep 2
    done
    echo "Failed to connect to PostgreSQL"
    return 1
}

# Функция для выполнения миграций
run_migrations() {
    echo "Running database migrations..."
    docker run --rm \
        --network hostel-booking-system_hostel_network \
        -v $(pwd)/backend/migrations:/migrations \
        migrate/migrate \
        -path=/migrations/ \
        -database="postgres://postgres:c9XWc7Cm@hostel_db:5432/hostel_db?sslmode=disable" \
        up

    if [ $? -eq 0 ]; then
        echo "Migrations completed successfully!"
        return 0
    else
        echo "Migration failed!"
        return 1
    fi
}

# Сохраняем важные файлы
echo "Backing up environment files and SSL certificates..."
mkdir -p /tmp/hostel-backup
cp -f backend/.env /tmp/hostel-backup/backend.env 2>/dev/null || true
cp -f frontend/hostel-frontend/.env /tmp/hostel-backup/frontend.env 2>/dev/null || true

# Сохраняем SSL сертификаты
if [ -d "certbot/conf" ]; then
    cp -r certbot/conf /tmp/hostel-backup/ 2>/dev/null || true
fi

# Сохраняем загруженные изображения
cp -r backend/uploads /tmp/hostel-backup/ 2>/dev/null || true

# Сбрасываем все локальные изменения
echo "Resetting local changes..."
git reset --hard
git clean -fdx -e "*.env*" -e "uploads/" -e "certbot/"

# Получаем последние изменения из git
echo "Pulling latest changes..."
git pull

# Восстанавливаем файлы
echo "Restoring backups..."
cp -f /tmp/hostel-backup/backend.env backend/.env 2>/dev/null || true
cp -f /tmp/hostel-backup/frontend.env frontend/hostel-frontend/.env 2>/dev/null || true
if [ -d "/tmp/hostel-backup/conf" ]; then
    rm -rf certbot/conf
    cp -r /tmp/hostel-backup/conf certbot/ 2>/dev/null || true
fi

# Устанавливаем права
chmod -R 755 backend/uploads
chmod -R 755 frontend/hostel-frontend/build
chmod -R 755 certbot/conf

# Удаляем бэкапы
rm -rf /tmp/hostel-backup

# Собираем фронтенд
echo "Building frontend..."
cd frontend/hostel-frontend
NODE_ENV=production docker run -v $(pwd):/app -w /app node:18 sh -c "npm install --legacy-peer-deps && npm install @babel/plugin-proposal-private-property-in-object && npm run build"

# Возвращаемся в корневую папку
cd ../..

# Останавливаем все контейнеры и удаляем volumes
echo "Stopping all containers..."
docker-compose -f docker-compose.prod.yml down -v

echo "Starting database..."
docker-compose -f docker-compose.prod.yml up -d hostel_db

# Ждем готовности базы данных
if wait_for_postgres; then
    if run_migrations; then
        echo "Starting remaining services..."
        docker-compose -f docker-compose.prod.yml up -d --build
        
        echo "Verifying database setup..."
        sleep 5  # Даем время на запуск всех сервисов
        
        echo "Database tables:"
        docker exec hostel_db psql -U postgres -d hostel_db -c "\dt"
        
        echo "Checking demo data..."
        docker exec hostel_db psql -U postgres -d hostel_db -c "SELECT COUNT(*) FROM rooms;"
        
        echo "Checking container status and logs..."
        docker-compose -f docker-compose.prod.yml ps
        docker logs hostel_nginx
        docker logs hostel_backend
    else
        echo "Failed to run migrations. Stopping deployment."
        exit 1
    fi
else
    echo "Database failed to start. Stopping deployment."
    exit 1
fi

echo "Deployment completed!"