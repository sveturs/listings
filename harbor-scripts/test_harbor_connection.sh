#!/bin/bash
set -e

# Скрипт для тестирования соединения с Harbor на сервере svetu.rs
# Запустите этот скрипт на сервере svetu.rs для проверки работы Harbor

HARBOR_URL="harbor.svetu.rs"
HARBOR_USER="admin"
HARBOR_PASSWORD="SveTu2025"
PROJECT_NAME="svetu"
TEST_IMAGE="postgres:15"

echo "==== Тестирование соединения с Harbor ($HARBOR_URL) ===="

# Проверка доступности Harbor
echo "1. Проверка доступности Harbor..."
if curl -s -f https://$HARBOR_URL/api/v2.0/health > /dev/null; then
  echo "✓ Harbor доступен по HTTPS"
else
  echo "✗ Harbor недоступен по HTTPS"
  echo "Попытка проверки через HTTP..."
  if curl -s -f http://$HARBOR_URL/api/v2.0/health > /dev/null; then
    echo "✓ Harbor доступен по HTTP"
  else
    echo "✗ Harbor недоступен по HTTP"
    echo "Проверка ping..."
    ping -c 3 $HARBOR_URL
  fi
fi

# Проверка авторизации
echo ""
echo "2. Тестирование авторизации в Harbor..."
docker logout $HARBOR_URL 2>/dev/null || true
if docker login -u $HARBOR_USER -p $HARBOR_PASSWORD $HARBOR_URL; then
  echo "✓ Авторизация в Harbor успешна"
else
  echo "✗ Ошибка авторизации в Harbor"
  echo "Проверка наличия директории для сертификатов..."
  if [ -d "/etc/docker/certs.d/$HARBOR_URL" ]; then
    echo "✓ Директория для сертификатов Harbor существует"
    ls -la /etc/docker/certs.d/$HARBOR_URL
  else
    echo "✗ Директория для сертификатов Harbor не существует"
    echo "Создание директории для сертификатов..."
    sudo mkdir -p /etc/docker/certs.d/$HARBOR_URL
    echo "Попытка повторной авторизации..."
    docker login -u $HARBOR_USER -p $HARBOR_PASSWORD $HARBOR_URL
  fi
fi

# Проверка загрузки образа из Harbor
echo ""
echo "3. Тестирование загрузки образа из Harbor..."
if docker pull $HARBOR_URL/$PROJECT_NAME/db/postgres:15; then
  echo "✓ Образ postgres:15 успешно загружен из Harbor"
else
  echo "✗ Ошибка загрузки образа postgres:15 из Harbor"
  echo "Проверка доступности образа в Harbor..."
  curl -u $HARBOR_USER:$HARBOR_PASSWORD -X GET "https://$HARBOR_URL/api/v2.0/projects/$PROJECT_NAME/repositories/db%2Fpostgres/artifacts/15" 2>/dev/null || echo "Образ не найден в Harbor"
fi

# Проверка тестового деплоя
echo ""
echo "4. Тестирование деплоя с использованием Harbor..."
TEST_DIR=$(mktemp -d)
cat > $TEST_DIR/docker-compose.test.yml << EOF
version: '3'
services:
  db:
    image: $HARBOR_URL/$PROJECT_NAME/db/postgres:15
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
      POSTGRES_DB: test
    ports:
      - "5432:5432"
EOF

cd $TEST_DIR
echo "Запуск тестового деплоя..."
if docker compose -f docker-compose.test.yml up -d; then
  echo "✓ Тестовый деплой с использованием Harbor успешен"
  echo "Останавливаем тестовый контейнер..."
  docker compose -f docker-compose.test.yml down
else
  echo "✗ Ошибка тестового деплоя с использованием Harbor"
fi

# Очистка
rm -rf $TEST_DIR

echo ""
echo "==== Тестирование соединения с Harbor завершено ===="