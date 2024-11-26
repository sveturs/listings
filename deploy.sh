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

# Функция для ожидания готовности базы данных
wait_for_postgres() {
    echo "Waiting for PostgreSQL to start..."
    for i in {1..30}; do
        if docker exec hostel_db pg_isready -U postgres > /dev/null 2>&1; then
            echo "PostgreSQL is ready!"
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

# Основной процесс деплоя...
# [предыдущий код до docker-compose down остается тем же]

# Перезапускаем контейнеры
echo "Restarting containers..."
docker-compose -f docker-compose.prod.yml down
docker-compose -f docker-compose.prod.yml up -d hostel_db

# Ждем готовности базы данных
wait_for_postgres

# Запускаем миграции после того, как база готова
if [ $? -eq 0 ]; then
    run_migrations
    
    if [ $? -eq 0 ]; then
        # Запускаем остальные сервисы только после успешных миграций
        docker-compose -f docker-compose.prod.yml up -d --build
        
        echo "Checking container status and logs..."
        docker-compose -f docker-compose.prod.yml ps
        docker logs hostel_nginx
        docker logs hostel_backend
        
        # Проверяем, что таблицы созданы
        echo "Checking database tables..."
        docker exec hostel_db psql -U postgres -d hostel_db -c "\dt"
    else
        echo "Failed to run migrations. Stopping deployment."
        exit 1
    fi
else
    echo "Database is not ready. Stopping deployment."
    exit 1
fi

echo "Deployment completed!"