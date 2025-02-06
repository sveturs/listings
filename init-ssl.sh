#!/bin/bash

# Определяем пути
PROJECT_DIR="/opt/hostel-booking-system"
cd $PROJECT_DIR

# Создаем необходимые директории
mkdir -p certbot/conf
mkdir -p certbot/www

# Создаем временный nginx.conf для первоначальной настройки
cat > nginx.initial.conf << 'EOF'
upstream api_backend {
    server hostel_backend:3000;
}

server {
    listen 80;
    server_name SveTu.rs www.SveTu.rs;
    
    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
        try_files $uri =404;
    }

    location / {
        return 200 'OK';
        add_header Content-Type text/plain;
    }
}
EOF

# Останавливаем существующие контейнеры
docker-compose -f docker-compose.prod.yml down

# Временно подменяем nginx.conf
cp nginx.conf nginx.conf.backup || true
cp nginx.initial.conf nginx.conf

# Запускаем nginx для первоначальной проверки домена
docker-compose -f docker-compose.prod.yml up -d nginx

# Ждем запуска nginx
sleep 10

echo "Attempting to get SSL certificate..."

# Получаем SSL сертификат с подробным выводом
docker-compose -f docker-compose.prod.yml run --rm certbot certbot certonly \
    --webroot \
    --webroot-path=/var/www/certbot \
    --email admin@SveTu.rs \
    --agree-tos \
    --no-eff-email \
    -d SveTu.rs \
    -d www.SveTu.rs \
    --force-renewal \
    --verbose

# Проверяем, были ли созданы сертификаты
if [ -d "certbot/conf/live/SveTu.rs" ]; then
    echo "Certificates successfully obtained!"
    
    # Восстанавливаем оригинальный nginx.conf
    mv nginx.conf.backup nginx.conf
    
    # Перезапускаем все контейнеры
    docker-compose -f docker-compose.prod.yml down
    docker-compose -f docker-compose.prod.yml up -d
else
    echo "Failed to obtain certificates!"
    # Восстанавливаем конфигурацию
    if [ -f nginx.conf.backup ]; then
        mv nginx.conf.backup nginx.conf
    fi
fi

# Удаляем временный конфиг
rm -f nginx.initial.conf

echo "SSL initialization completed!"