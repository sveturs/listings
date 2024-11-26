#!/bin/bash

# Создаем необходимые директории
mkdir -p certbot/conf
mkdir -p certbot/www

# Останавливаем существующие контейнеры
docker-compose -f docker-compose.prod.yml down

# Запускаем nginx для первоначальной проверки домена
docker-compose -f docker-compose.prod.yml up -d nginx

# Ждем запуска nginx
sleep 5

# Получаем SSL сертификат
docker-compose -f docker-compose.prod.yml run --rm certbot certbot certonly \
    --webroot \
    --webroot-path=/var/www/certbot \
    --email admin@landhub.rs \
    --agree-tos \
    --no-eff-email \
    -d landhub.rs \
    -d www.landhub.rs

docker-compose -f docker-compose.prod.yml down
docker-compose -f docker-compose.prod.yml up -d