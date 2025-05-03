#!/bin/bash
set -e

# Скрипт для обновления deploy.sh для интеграции с Harbor
# Настройки
DEPLOY_SCRIPT_PATH="/data/hostel-booking-system/deploy.sh"
HARBOR_URL="207.180.197.172"
HARBOR_USER="admin"
HARBOR_PASSWORD="SveTu2025"

# Проверка наличия файла
if [ ! -f "$DEPLOY_SCRIPT_PATH" ]; then
  echo "Ошибка: Файл deploy.sh не найден по пути $DEPLOY_SCRIPT_PATH"
  exit 1
fi

# Создание резервной копии
BACKUP_PATH="${DEPLOY_SCRIPT_PATH}.bak-$(date +%Y%m%d-%H%M%S)"
cp "$DEPLOY_SCRIPT_PATH" "$BACKUP_PATH"
echo "Создана резервная копия файла: $BACKUP_PATH"

# Добавление авторизации Harbor в начало скрипта
sed -i '2i\# Авторизация в Harbor' "$DEPLOY_SCRIPT_PATH"
sed -i '3i\echo "Авторизация в Harbor..."' "$DEPLOY_SCRIPT_PATH"
sed -i '4i\docker login -u '"$HARBOR_USER"' -p '"$HARBOR_PASSWORD"' '"$HARBOR_URL"'' "$DEPLOY_SCRIPT_PATH"
sed -i '5i\\' "$DEPLOY_SCRIPT_PATH"

echo "Скрипт deploy.sh обновлен - добавлена авторизация в Harbor."

# Заменяем docker-compose команды по необходимости
# Пример: обновление правил для использования --pull always
sed -i 's/docker-compose -f docker-compose.prod.yml up/docker-compose -f docker-compose.prod.yml pull \&\& docker-compose -f docker-compose.prod.yml up/g' "$DEPLOY_SCRIPT_PATH"

echo "Скрипт deploy.sh обновлен - добавлена принудительная загрузка актуальных образов."

echo "Обновление скрипта deploy.sh завершено!"
echo "Проверьте обновленный скрипт и внесите дополнительные изменения при необходимости."