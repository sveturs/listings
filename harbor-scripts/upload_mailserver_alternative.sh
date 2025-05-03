#!/bin/bash
set -e

# Скрипт для загрузки образа mailserver через архив

echo "Начинаем загрузку образа mailserver в Harbor (альтернативный метод)..."

# Параметры
HARBOR_URL="harbor.svetu.rs"
HARBOR_USER="admin"
HARBOR_PASSWORD="SveTu2025"
SOURCE_IMAGE="mailserver/docker-mailserver:latest"
TARGET_IMAGE="$HARBOR_URL/svetu/mail/server:latest"
ARCHIVE_DIR="/tmp/docker_archives"
ARCHIVE_FILE="$ARCHIVE_DIR/mailserver.tar"

# Создаем директорию для архивов
mkdir -p $ARCHIVE_DIR

# Авторизация в Harbor
echo "Авторизация в Harbor..."
docker login -u $HARBOR_USER -p $HARBOR_PASSWORD $HARBOR_URL

# Проверяем наличие образа локально
if ! docker image inspect $SOURCE_IMAGE &>/dev/null; then
  echo "Образ $SOURCE_IMAGE не найден локально. Пожалуйста, убедитесь, что он существует."
  exit 1
fi

# Создаем временный архив образа
echo "Создаем временный архив образа..."
docker save $SOURCE_IMAGE -o $ARCHIVE_FILE
echo "Архив успешно создан в $ARCHIVE_FILE (размер: $(du -h $ARCHIVE_FILE | cut -f1))"

# Тегирование образа для Harbor
echo "Тегирование образа для Harbor..."
docker tag $SOURCE_IMAGE $TARGET_IMAGE

# Загружаем образ в Harbor
echo "Загружаем образ в Harbor с увеличенным таймаутом..."
DOCKER_CLIENT_TIMEOUT=1800 COMPOSE_HTTP_TIMEOUT=1800 docker push $TARGET_IMAGE
if [ $? -eq 0 ]; then
  echo "✅ Образ успешно загружен в Harbor!"
else
  echo "❌ Ошибка при загрузке образа в Harbor."
  exit 1
fi

# Проверяем, что образ доступен в Harbor
echo "Проверка доступности образа в Harbor..."
if curl -k -u "$HARBOR_USER:$HARBOR_PASSWORD" "https://$HARBOR_URL/api/v2.0/projects/svetu/repositories/mail%2Fserver/artifacts" | grep -q "digest"; then
  echo "✅ Образ mailserver успешно загружен и доступен в Harbor!"
else
  echo "❌ Не удалось проверить наличие образа в Harbor. Повторная проверка..."
  sleep 10
  if curl -k -u "$HARBOR_USER:$HARBOR_PASSWORD" "https://$HARBOR_URL/api/v2.0/projects/svetu/repositories/mail%2Fserver/artifacts" | grep -q "digest"; then
    echo "✅ Образ mailserver успешно загружен и доступен в Harbor!"
  else
    echo "❌ Не удалось проверить наличие образа в Harbor. Возможно, API недоступен."
    echo "Проверьте наличие образа в веб-интерфейсе Harbor."
  fi
fi

# Удаляем временный архив
echo "Удаляем временный архив..."
rm -f $ARCHIVE_FILE

echo "Процесс загрузки mailserver завершен!"