#!/bin/bash
set -e

# Скрипт для загрузки оставшихся mail-образов в Harbor

echo "Начинаем загрузку оставшихся mail-образов в Harbor..."

# Запускаем скрипт для загрузки mailserver
echo "=== ЗАГРУЗКА MAILSERVER ==="
if ./upload_mailserver.sh; then
  echo "✅ Скрипт upload_mailserver.sh выполнен успешно"
else
  echo "❌ Ошибка при выполнении скрипта upload_mailserver.sh"
  echo "Продолжаем загрузку mail-webui..."
fi

# Пауза между загрузками
echo "Ожидание 30 секунд перед следующей загрузкой..."
sleep 30

# Запускаем скрипт для загрузки mail-webui
echo "=== ЗАГРУЗКА MAIL-WEBUI ==="
if ./upload_mail_webui.sh; then
  echo "✅ Скрипт upload_mail_webui.sh выполнен успешно"
else
  echo "❌ Ошибка при выполнении скрипта upload_mail_webui.sh"
fi

echo "Процесс загрузки оставшихся образов завершен!"

# Проверка статуса загрузки всех образов
echo "=== ИТОГОВЫЙ СТАТУС ==="
echo "Проверка наличия образов в Harbor:"

HARBOR_URL="harbor.svetu.rs"
HARBOR_USER="admin"
HARBOR_PASSWORD="SveTu2025"

# Функция проверки наличия образа
check_image() {
  local repo=$1
  local image_path=$2
  
  if curl -u "$HARBOR_USER:$HARBOR_PASSWORD" -X GET -s "https://$HARBOR_URL/api/v2.0/projects/svetu/repositories/$repo/artifacts/latest" | grep -q "digest"; then
    echo "✅ $image_path - успешно загружен"
    return 0
  else
    echo "❌ $image_path - отсутствует или недоступен"
    return 1
  fi
}

# Проверяем наличие образов
check_image "mail%2Fserver" "svetu/mail/server:latest"
check_image "mail%2Fwebui" "svetu/mail/webui:latest"

echo "Скрипт завершен!"