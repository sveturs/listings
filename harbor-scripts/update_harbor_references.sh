#!/bin/bash
set -e

# Скрипт для обновления ссылок на Harbor в проекте
# Запустите этот скрипт из корня проекта

OLD_HARBOR_URL="207.180.197.172"
NEW_HARBOR_URL="harbor.svetu.rs"
BASE_DIR="/data/hostel-booking-system"

echo "==== Обновление ссылок на Harbor с $OLD_HARBOR_URL на $NEW_HARBOR_URL ===="

# Функция для обновления файла
update_file() {
  local file_path=$1
  
  if [ -f "$file_path" ]; then
    echo "Обновление файла: $file_path"
    # Создаем резервную копию файла
    cp "$file_path" "${file_path}.bak-$(date +%Y%m%d-%H%M%S)"
    # Заменяем URL
    sed -i "s|$OLD_HARBOR_URL|$NEW_HARBOR_URL|g" "$file_path"
  else
    echo "Файл не найден: $file_path"
  fi
}

# Обновление docker-compose файлов
update_file "$BASE_DIR/docker-compose.prod.yml.harbor"
update_file "$BASE_DIR/docker-compose.yml"

# Обновление скриптов в harbor-scripts
update_file "$BASE_DIR/harbor-scripts/build_and_push.sh"
update_file "$BASE_DIR/harbor-scripts/deploy_with_harbor.sh"
update_file "$BASE_DIR/harbor-scripts/fixed_migrate_images.sh"
update_file "$BASE_DIR/harbor-scripts/setup_cert_trust.sh"

# Обновление документации в harbor-docs
update_file "$BASE_DIR/harbor-docs/deployment_instructions.md"
update_file "$BASE_DIR/harbor-docs/final_harbor_report.md"
update_file "$BASE_DIR/harbor-docs/README.md"

# Обновление других возможных файлов
update_file "$BASE_DIR/deploy.sh"

echo ""
echo "==== Обновление ссылок завершено ===="
echo "Все ссылки на Harbor обновлены с $OLD_HARBOR_URL на $NEW_HARBOR_URL"
echo ""
echo "ВАЖНО: Проверьте файлы и убедитесь, что все ссылки обновлены корректно"
echo "Вы можете выполнить: grep -r \"$OLD_HARBOR_URL\" $BASE_DIR --include=\"*.sh\" --include=\"*.md\" --include=\"*.yml\""