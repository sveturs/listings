#!/bin/bash

echo "Исправляем неправильные ссылки на модели logistics в Swagger аннотациях..."

# Список файлов для исправления
files=(
    "internal/proj/admin/logistics/handler/analytics.go"
    "internal/proj/admin/logistics/handler/dashboard.go"
    "internal/proj/delivery/handler/admin_handler.go"
)

# Исправления
for file in "${files[@]}"; do
    if [ -f "$file" ]; then
        echo "Обрабатываем файл: $file"

        # Заменяем неправильные ссылки на модели logistics
        sed -i 's/backend_internal_domain_logistics\./backend_internal_domain_logistics./g' "$file"
        sed -i 's/backend_internal_proj_admin_logistics_handler\./backend_internal_proj_admin_logistics_handler./g' "$file"

        # Особые случаи для конкретных моделей
        sed -i 's/backend_internal_proj_postexpress_models\./backend_internal_proj_postexpress_handler./g' "$file"
        sed -i 's/backend_internal_proj_bexexpress_models\./backend_internal_proj_bexexpress_handler./g' "$file"
    fi
done

echo "Исправление завершено!"