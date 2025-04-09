#!/bin/bash
cd /opt/hostel-booking-system

# Скачиваем скрипт управления
curl -o mailserver/config/setup.sh https://raw.githubusercontent.com/docker-mailserver/docker-mailserver/master/setup.sh
chmod +x mailserver/config/setup.sh

# Создаем учетную запись info@svetu.rs
docker-compose -f docker-compose.prod.yml exec mailserver setup email add info@svetu.rs password123

# Создаем учетную запись postmaster@svetu.rs
docker-compose -f docker-compose.prod.yml exec mailserver setup email add postmaster@svetu.rs password123

# Создаем DKIM ключи
docker-compose -f docker-compose.prod.yml exec mailserver setup config dkim
