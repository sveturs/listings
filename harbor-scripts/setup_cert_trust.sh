#!/bin/bash
set -e

# Скрипт для настройки доверия к самоподписанным сертификатам Harbor
# Запустите этот скрипт на всех серверах, которые должны работать с Harbor
# Например: ./setup_cert_trust.sh 207.180.197.172

# Проверка наличия аргумента
if [ -z "$1" ]; then
  echo "Ошибка: не указан IP-адрес или домен Harbor"
  echo "Использование: $0 <harbor_ip_or_domain>"
  echo "Пример: $0 207.180.197.172"
  exit 1
fi

HARBOR_HOST=$1
CERT_DIR="/etc/docker/certs.d/$HARBOR_HOST"
HARBOR_USER="dima"

echo "==== Настройка доверия к сертификатам Harbor ($HARBOR_HOST) ===="

# Создание директории для сертификатов
echo "Создание директории для сертификатов..."
sudo mkdir -p $CERT_DIR

# Проверка, можно ли скопировать сертификат с сервера Harbor
echo "Попытка скопировать сертификат с сервера Harbor..."
if ssh -q -o BatchMode=yes -o ConnectTimeout=5 $HARBOR_USER@$HARBOR_HOST exit; then
  echo "Подключение к серверу Harbor доступно, копирование сертификата..."
  sudo scp $HARBOR_USER@$HARBOR_HOST:/opt/harbor/certs/server.crt $CERT_DIR/ca.crt
  echo "Сертификат успешно скопирован"
else
  echo "Подключение к серверу Harbor недоступно"
  echo "Вам нужно вручную скопировать сертификат с сервера Harbor"
  echo "Выполните на сервере Harbor:"
  echo "  cat /opt/harbor/certs/server.crt"
  echo "Скопируйте вывод и создайте файл $CERT_DIR/ca.crt на этом сервере"
  
  # Создание пустого файла для сертификата
  sudo touch $CERT_DIR/ca.crt
  sudo chmod 644 $CERT_DIR/ca.crt
  
  echo "Файл $CERT_DIR/ca.crt создан, но его необходимо заполнить содержимым сертификата"
  echo "Используйте команду: sudo nano $CERT_DIR/ca.crt"
  exit 1
fi

# Проверка наличия сертификата
if [ ! -f "$CERT_DIR/ca.crt" ]; then
  echo "Ошибка: сертификат не был скопирован"
  exit 1
fi

# Перезапуск Docker
echo "Перезапуск Docker для применения изменений..."
sudo systemctl restart docker

# Проверка, работает ли Docker
echo "Проверка работы Docker..."
docker info > /dev/null
if [ $? -eq 0 ]; then
  echo "Docker успешно перезапущен"
else
  echo "Ошибка: Docker не запустился после перезапуска"
  echo "Проверьте логи Docker: sudo journalctl -u docker"
  exit 1
fi

# Проверка доступности Harbor
echo "Проверка доступности Harbor..."
if curl -s -f https://$HARBOR_HOST/api/v2.0/health > /dev/null; then
  echo "Harbor доступен по HTTPS"
else
  echo "Предупреждение: Harbor недоступен по HTTPS"
  echo "Это может быть связано с тем, что Harbor не настроен на использование HTTPS"
  echo "или другими проблемами с сетью"
fi

# Проверка авторизации в Harbor
echo "Проверка авторизации в Harbor..."
echo "Введите учетные данные для Harbor (по умолчанию: admin/SveTu2025)"
docker login $HARBOR_HOST

echo ""
echo "==== Настройка доверия к сертификатам Harbor завершена ===="
echo "Теперь Docker должен доверять сертификатам Harbor"
echo "Вы можете проверить это, выполнив: docker pull $HARBOR_HOST/svetu/backend/api:latest"