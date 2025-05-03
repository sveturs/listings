#!/bin/bash
set -e

# Скрипт для обновления домена Harbor
# Запустите на сервере Harbor (harbor.svetu.rs / 207.180.197.172)

HARBOR_DOMAIN="harbor.svetu.rs"
HARBOR_DIR="/opt/harbor/harbor"  # Обновленный путь к директории Harbor
CONFIG_PATH="$HARBOR_DIR/harbor.yml"
EMAIL="info@svetu.rs"  # Email для Let's Encrypt

echo "==== Обновление домена Harbor на $HARBOR_DOMAIN ===="

# Проверка существования директории Harbor
if [ ! -d "$HARBOR_DIR" ]; then
  echo "Ошибка: директория Harbor не найдена ($HARBOR_DIR)"
  echo "Убедитесь, что вы запускаете скрипт на сервере Harbor"
  exit 1
fi

# Создание резервной копии конфигурации
echo "Создание резервной копии конфигурации..."
cp "$CONFIG_PATH" "${CONFIG_PATH}.bak-$(date +%Y%m%d-%H%M%S)"

# Обновление домена в harbor.yml
echo "Обновление домена в конфигурации Harbor..."
sed -i "s/^hostname:.*/hostname: $HARBOR_DOMAIN/" "$CONFIG_PATH"

# Настройка Let's Encrypt
echo "Установка certbot для получения сертификатов Let's Encrypt..."
apt-get update
apt-get install -y certbot

# Временная остановка Harbor
echo "Остановка Harbor для обновления сертификатов..."
cd "$HARBOR_DIR"
docker compose down || docker-compose down

# Получение сертификатов Let's Encrypt
echo "Получение сертификатов Let's Encrypt для $HARBOR_DOMAIN..."
mkdir -p $HARBOR_DIR/certs
certbot certonly --standalone -d $HARBOR_DOMAIN --email $EMAIL --agree-tos -n

# Копирование сертификатов
echo "Копирование сертификатов в директорию Harbor..."
cp /etc/letsencrypt/live/$HARBOR_DOMAIN/fullchain.pem $HARBOR_DIR/certs/server.crt
cp /etc/letsencrypt/live/$HARBOR_DOMAIN/privkey.pem $HARBOR_DIR/certs/server.key

# Обновление путей к сертификатам в harbor.yml
echo "Обновление путей к сертификатам в конфигурации..."
sed -i "s|^  certificate:.*|  certificate: $HARBOR_DIR/certs/server.crt|" "$CONFIG_PATH"
sed -i "s|^  private_key:.*|  private_key: $HARBOR_DIR/certs/server.key|" "$CONFIG_PATH"

# Перезапуск Harbor с новой конфигурацией
echo "Запуск Harbor с новой конфигурацией..."
./prepare
docker compose up -d || docker-compose up -d

# Настройка автообновления сертификатов
echo "Настройка автообновления сертификатов..."
cat > /etc/cron.d/certbot-renewal << EOF
0 0 * * * root certbot renew --quiet --post-hook "cp /etc/letsencrypt/live/$HARBOR_DOMAIN/fullchain.pem $HARBOR_DIR/certs/server.crt && cp /etc/letsencrypt/live/$HARBOR_DOMAIN/privkey.pem $HARBOR_DIR/certs/server.key && cd $HARBOR_DIR && (docker compose restart nginx || docker-compose restart nginx)"
EOF

echo ""
echo "==== Обновление домена Harbor завершено ===="
echo "Harbor теперь доступен по адресу: https://$HARBOR_DOMAIN"
echo "Сертификаты Let's Encrypt успешно установлены и будут автоматически обновляться"
echo ""
echo "ВАЖНО: Обновите URL в скриптах и docker-compose файлах с 207.180.197.172 на $HARBOR_DOMAIN"