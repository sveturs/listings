#!/bin/bash

domains=(svetu.rs www.svetu.rs mail.svetu.rs autodiscover.svetu.rs autoconfig.svetu.rs klimagrad.svetu.rs)
email="admin@svetu.rs"

# Остановка контейнеров, если они запущены
docker-compose -f docker-compose.prod.yml down

# Создание каталогов для certbot, если они не существуют
mkdir -p ./certbot/conf/live/svetu.rs
mkdir -p ./certbot/conf/live/mail.svetu.rs
mkdir -p ./certbot/conf/live/autodiscover.svetu.rs
mkdir -p ./certbot/conf/live/klimagrad.svetu.rs
mkdir -p ./certbot/www

# Создание временных самоподписанных сертификатов для первого запуска nginx
openssl req -x509 -nodes -newkey rsa:4096 -days 1 \
  -keyout ./certbot/conf/live/svetu.rs/privkey.pem \
  -out ./certbot/conf/live/svetu.rs/fullchain.pem \
  -subj "/CN=svetu.rs"

cp ./certbot/conf/live/svetu.rs/privkey.pem ./certbot/conf/live/mail.svetu.rs/privkey.pem
cp ./certbot/conf/live/svetu.rs/fullchain.pem ./certbot/conf/live/mail.svetu.rs/fullchain.pem
cp ./certbot/conf/live/svetu.rs/privkey.pem ./certbot/conf/live/autodiscover.svetu.rs/privkey.pem
cp ./certbot/conf/live/svetu.rs/fullchain.pem ./certbot/conf/live/autodiscover.svetu.rs/fullchain.pem
cp ./certbot/conf/live/svetu.rs/privkey.pem ./certbot/conf/live/klimagrad.svetu.rs/privkey.pem
cp ./certbot/conf/live/svetu.rs/fullchain.pem ./certbot/conf/live/klimagrad.svetu.rs/fullchain.pem

# Запуск только nginx
docker-compose -f docker-compose.prod.yml up -d nginx

# Получение реальных сертификатов
for domain in "${domains[@]}"; do
  docker-compose -f docker-compose.prod.yml run --rm certbot certonly --webroot --webroot-path=/var/www/certbot \
    --email $email --agree-tos --no-eff-email --force-renewal -d $domain
done

# Перезапуск всех сервисов
docker-compose -f docker-compose.prod.yml down
docker-compose -f docker-compose.prod.yml up -d