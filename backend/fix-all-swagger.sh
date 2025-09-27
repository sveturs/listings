#!/bin/bash

echo "Исправляю все Swagger типы..."

# Находим все файлы в /internal/proj и исправляем типы
find /data/hostel-booking-system/backend/internal/proj -name "*.go" -type f | while read file; do
    # Получаем путь пакета из файла
    dir=$(dirname "$file")
    package_path=$(echo "$dir" | sed 's|/data/hostel-booking-system/backend/||')

    # Исправляем ссылки на неправильные типы
    sed -i 's/backend_internal_proj_marketplace_service\./service./g' "$file"

    # Для delivery модуля
    if [[ "$file" == *"/delivery/"* ]]; then
        sed -i 's/backend_internal_domain_models\./backend_internal_proj_delivery_models./g' "$file"
    fi

    # Для postexpress модуля
    if [[ "$file" == *"/postexpress/"* ]]; then
        sed -i 's/backend_internal_domain_models\./backend_internal_proj_postexpress_models./g' "$file"
    fi

    # Для storefronts модуля
    if [[ "$file" == *"/storefronts/"* ]]; then
        sed -i 's/backend_internal_domain_models\./backend_internal_domain_models./g' "$file"
    fi

    # Для orders модуля
    if [[ "$file" == *"/orders/"* ]]; then
        sed -i 's/backend_internal_domain_models\./backend_internal_domain_models./g' "$file"
    fi

    # Общие исправления
    sed -i 's/handler\.\([A-Z]\)/internal_proj_marketplace_handler.\1/g' "$file"
    sed -i 's/service\.\([A-Z]\)/backend_internal_proj_marketplace_service.\1/g' "$file"
done

echo "Готово!"