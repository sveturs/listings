#!/bin/bash
set -e  # Останавливаем выполнение при ошибках

echo "Starting deployment..."

# Переходим в папку проекта
cd /opt/hostel-booking-system

# Настраиваем git pull strategy
echo "Configuring git..."
git config pull.rebase false

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
if [ -d "certbot/conf" ]; then
    cp -r certbot/conf /tmp/hostel-backup/ 2>/dev/null || true
fi

# Сохраняем загруженные изображения
cp -r backend/uploads /tmp/hostel-backup/ 2>/dev/null || true

# Обеспечиваем чистое состояние git
echo "Resetting git state..."
git fetch origin
git reset --hard origin/main
git clean -fdx -e "*.env*" -e "uploads/" -e "certbot/"

# Восстанавливаем файлы
echo "Restoring backups..."
cp -f /tmp/hostel-backup/backend.env backend/.env 2>/dev/null || true
cp -f /tmp/hostel-backup/frontend.env frontend/hostel-frontend/.env 2>/dev/null || true
if [ -d "/tmp/hostel-backup/conf" ]; then
    rm -rf certbot/conf
    cp -r /tmp/hostel-backup/conf certbot/ 2>/dev/null || true
fi

# Собираем фронтенд
echo "Building frontend..."
cd frontend/hostel-frontend
NODE_ENV=production docker run -v $(pwd):/app -w /app node:18 sh -c "npm install --legacy-peer-deps && npm install @babel/plugin-proposal-private-property-in-object && npm run build"
cd ../..

# Останавливаем все контейнеры
echo "Stopping all containers..."
docker-compose -f docker-compose.prod.yml down -v

# Создаем сеть если её нет
docker network create hostel-booking-system_hostel_network 2>/dev/null || true

# Запускаем только базу данных
echo "Starting database..."
docker-compose -f docker-compose.prod.yml up -d db

# Функция для проверки готовности базы данных
check_db() {
    echo "Checking database connection..."
    docker-compose -f docker-compose.prod.yml exec -T db pg_isready -U postgres
}

# Ждем готовности базы данных
echo "Waiting for database to start..."
for i in {1..30}; do
    if check_db > /dev/null 2>&1; then
        echo "Database is ready!"
        break
    fi
    echo "Waiting for database... Attempt $i/30"
    sleep 2
done

# Проверяем финальную готовность
if ! check_db > /dev/null 2>&1; then
    echo "Database failed to start"
    exit 1
fi

# Выполняем миграции
echo "Running migrations..."
docker run --rm \
    --network hostel-booking-system_hostel_network \
    -v $(pwd)/backend/migrations:/migrations \
    migrate/migrate \
    -path=/migrations/ \
    -database="postgres://postgres:c9XWc7Cm@db:5432/hostel_db?sslmode=disable" \
    up

# Проверяем результат миграций
if [ $? -eq 0 ]; then
    echo "Migrations successful! Starting other services..."
    
    # Запускаем остальные сервисы
    docker-compose -f docker-compose.prod.yml up -d

    # Проверяем структуру базы данных
    echo "Checking database structure..."
    sleep 5
    docker-compose -f docker-compose.prod.yml exec -T db psql -U postgres -d hostel_db -c "\dt"
    
    echo "Checking container status..."
    docker-compose -f docker-compose.prod.yml ps
else
    echo "Migration failed!"
    exit 1
fi

# Устанавливаем права на директории
echo "Setting permissions..."
chmod -R 755 backend/uploads || true
chmod -R 755 frontend/hostel-frontend/build || true
chmod -R 755 certbot/conf || true

# Удаляем бэкапы
rm -rf /tmp/hostel-backup

echo "Deployment completed!"