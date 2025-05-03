#!/bin/bash
set -e

# Скрипт для загрузки образа mailserver в Harbor с увеличенным таймаутом и механизмом повторных попыток

echo "Начинаем загрузку образа mailserver в Harbor..."

# Параметры
MAX_ATTEMPTS=5
WAIT_TIME=60 # секунд между попытками
HARBOR_URL="harbor.svetu.rs"
HARBOR_USER="admin"
HARBOR_PASSWORD="SveTu2025"
SOURCE_IMAGE="mailserver/docker-mailserver:latest"
TARGET_IMAGE="$HARBOR_URL/svetu/mail/server:latest"

# Авторизация в Harbor
echo "Авторизация в Harbor..."
docker login -u $HARBOR_USER -p $HARBOR_PASSWORD $HARBOR_URL

# Проверяем наличие образа локально
if ! docker image inspect $SOURCE_IMAGE &>/dev/null; then
  echo "Образ $SOURCE_IMAGE не найден локально. Пожалуйста, убедитесь, что он существует."
  exit 1
fi

# Тегирование образа для Harbor
echo "Тегирование образа для Harbor..."
docker tag $SOURCE_IMAGE $TARGET_IMAGE

# Функция загрузки с повторными попытками
upload_with_retry() {
  local attempt=1
  local success=false

  while [ $attempt -le $MAX_ATTEMPTS ] && [ "$success" = false ]; do
    echo "Попытка $attempt из $MAX_ATTEMPTS загрузить образ в Harbor..."
    
    # Увеличиваем таймаут Docker для push операции
    if DOCKER_CLIENT_TIMEOUT=1800 COMPOSE_HTTP_TIMEOUT=1800 docker push $TARGET_IMAGE; then
      success=true
      echo "Образ успешно загружен в Harbor!"
    else
      echo "Попытка $attempt не удалась."
      if [ $attempt -lt $MAX_ATTEMPTS ]; then
        echo "Ожидание 30 секунд перед следующей попыткой..."
        sleep 30
      fi
      attempt=$((attempt+1))
    fi
  done

  if [ "$success" = false ]; then
    echo "Не удалось загрузить образ после $MAX_ATTEMPTS попыток."
    return 1
  fi
  
  return 0
}

# Запускаем загрузку с повторными попытками
upload_with_retry

# Проверяем, что образ доступен в Harbor
echo "Проверка доступности образа в Harbor..."
if curl -u "$HARBOR_USER:$HARBOR_PASSWORD" -X GET -s "https://$HARBOR_URL/api/v2.0/projects/svetu/repositories/mail%2Fserver/artifacts/latest" | grep -q "digest"; then
  echo "✅ Образ mailserver успешно загружен и доступен в Harbor!"
else
  echo "❌ Не удалось проверить наличие образа в Harbor. Возможно, загрузка не завершена или произошла ошибка."
  exit 1
fi

echo "Загрузка mailserver завершена!"