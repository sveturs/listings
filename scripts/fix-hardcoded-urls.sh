#!/bin/bash

# Скрипт для исправления захардкоженных URL в frontend

echo "🔧 Исправление захардкоженных URL в frontend файлах..."

cd /data/hostel-booking-system/frontend/svetu

# Список файлов с захардкоженными URL (кроме уже исправленных)
FILES=(
  "src/app/api/admin/logistics/dashboard/route.ts"
  "src/app/[locale]/admin/storefronts/page.tsx"
  "src/app/[locale]/admin/storefronts/AdminStorefrontsTable.tsx"
  "src/app/api/admin/search/analytics/export/route.ts"
  "src/app/[locale]/admin/storefront-products/page.tsx"
  "src/app/[locale]/admin/storefront-products/AdminStorefrontProductsTable.tsx"
  "src/app/[locale]/admin/listings/page.tsx"
  "src/app/[locale]/admin/listings/AdminListingsTable.tsx"
  "src/app/api/admin/translations/costs/route.ts"
  "src/app/api/v1/admin/search/synonyms/route.ts"
  "src/app/api/v1/admin/search/synonyms/[id]/route.ts"
)

# Функция добавления импорта если его нет
add_config_import() {
  local file="$1"
  
  if ! grep -q "import configManager from '@/config'" "$file"; then
    # Добавляем импорт после существующих импортов
    sed -i "1i import configManager from '@/config';" "$file"
  fi
}

# Функция замены URL
fix_url_in_file() {
  local file="$1"
  
  echo "  📝 Обрабатываю $file"
  
  # Добавляем импорт configManager если нужно
  if grep -q "http://localhost:3000" "$file"; then
    add_config_import "$file"
  fi
  
  # Заменяем различные паттерны захардкоженных URL
  sed -i "s|process\.env\.NEXT_PUBLIC_API_URL || 'http://localhost:3000'|process.env.NEXT_PUBLIC_API_URL || configManager.getApiUrl()|g" "$file"
  sed -i "s|'http://localhost:3000'|configManager.getApiUrl()|g" "$file"
  sed -i "s|\"http://localhost:3000\"|configManager.getApiUrl()|g" "$file"
  sed -i "s|\`http://localhost:3000\`|configManager.getApiUrl()|g" "$file"
  
  echo "    ✅ Готово"
}

# Обрабатываем каждый файл
for file in "${FILES[@]}"; do
  if [ -f "$file" ]; then
    fix_url_in_file "$file"
  else
    echo "  ⚠️  Файл $file не найден, пропускаю"
  fi
done

echo "🎉 Исправление завершено!"
echo "📋 Проверьте следующие файлы и убедитесь что импорты корректны:"
for file in "${FILES[@]}"; do
  if [ -f "$file" ]; then
    echo "  - $file"
  fi
done