#!/bin/bash

# Скрипт для полного сохранения сайта с помощью wget

SITE_URL="http://localhost:3001"
OUTPUT_DIR="./site-offline-copy"

echo "🌐 Начинаем полное сохранение сайта Sve Tu..."
echo "📁 Директория для сохранения: $OUTPUT_DIR"

# Создаем директорию
mkdir -p "$OUTPUT_DIR"

# Основные параметры wget:
# -r, --recursive - рекурсивная загрузка
# -l 3 - глубина рекурсии (3 уровня)
# -k, --convert-links - конвертировать ссылки для оффлайн просмотра
# -p, --page-requisites - загрузить все необходимые ресурсы (CSS, JS, изображения)
# -E, --adjust-extension - добавить .html расширение где нужно
# -K, --backup-converted - сохранить оригинальные файлы
# -np, --no-parent - не подниматься выше начальной директории
# -N, --timestamping - не перезагружать файлы если не изменились
# --no-host-directories - не создавать директорию с именем хоста
# --restrict-file-names=windows - совместимость имен файлов с Windows
# --user-agent - идентификация как браузер

wget \
  --recursive \
  --level=3 \
  --convert-links \
  --page-requisites \
  --adjust-extension \
  --backup-converted \
  --no-parent \
  --timestamping \
  --no-host-directories \
  --directory-prefix="$OUTPUT_DIR" \
  --restrict-file-names=windows \
  --user-agent="Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36" \
  --wait=1 \
  --random-wait \
  --accept="html,htm,css,js,json,jpg,jpeg,png,gif,svg,woff,woff2,ttf,eot" \
  --reject="pdf,zip,tar,gz" \
  --execute robots=off \
  "$SITE_URL" \
  "$SITE_URL/marketplace" \
  "$SITE_URL/create-listing-choice" \
  "$SITE_URL/storefronts" \
  "$SITE_URL/ideal-homepage" \
  "$SITE_URL/ideal-homepage-v2" \
  "$SITE_URL/auth/login" \
  "$SITE_URL/auth/register"

echo ""
echo "✅ Загрузка завершена!"
echo ""
echo "📋 Дополнительные действия:"
echo "1. Проверьте сохраненные файлы в директории: $OUTPUT_DIR"
echo "2. Откройте index.html в браузере для проверки"
echo "3. Заархивируйте папку для отправки:"
echo "   tar -czf sve-tu-offline.tar.gz $OUTPUT_DIR"
echo ""
echo "⚠️  Примечание: Некоторые динамические функции могут не работать в оффлайн версии"