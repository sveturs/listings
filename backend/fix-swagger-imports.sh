#!/bin/bash

echo "Исправляю Swagger импорты в Go файлах..."

# Исправляем backend_internal_domain_models на models
echo "Исправляю backend_internal_domain_models..."
rg -l "backend_internal_domain_models" --type go | while read file; do
    echo "  Обрабатываю: $file"
    sed -i 's/backend_internal_domain_models\./models./g' "$file"
done

# Исправляем internal_proj_marketplace_handler
echo "Исправляю internal_proj_marketplace_handler..."
rg -l "internal_proj_marketplace_handler" --type go | while read file; do
    echo "  Обрабатываю: $file"
    sed -i 's/internal_proj_marketplace_handler\./handler./g' "$file"
done

# Исправляем backend_internal_proj_marketplace_service
echo "Исправляю backend_internal_proj_marketplace_service..."
rg -l "backend_internal_proj_marketplace_service" --type go | while read file; do
    echo "  Обрабатываю: $file"
    sed -i 's/backend_internal_proj_marketplace_service\./service./g' "$file"
done

# Исправляем backend_pkg_utils
echo "Исправляю backend_pkg_utils..."
rg -l "backend_pkg_utils" --type go | while read file; do
    echo "  Обрабатываю: $file"
    sed -i 's/backend_pkg_utils\./utils./g' "$file"
done

echo "Готово! Теперь регенерируем Swagger документацию..."
cd /data/hostel-booking-system/backend && make generate-types

echo "Исправление завершено!"