#!/bin/bash

# Список файлов с дублированными импортами
FILES=(
    "internal/middleware/middleware.go"
    "internal/proj/balance/handler/balance.go"
    "internal/proj/gis/handler/advanced_filters.go"
    "internal/proj/gis/handler/district_handler.go"
    "internal/proj/gis/handler/spatial_search.go"
    "internal/proj/global/handler/handler.go"
    "internal/proj/global/handler/unified_search.go"
    "internal/proj/marketplace/handler/admin_attributes.go"
    "internal/proj/marketplace/handler/admin_translations.go"
    "internal/proj/marketplace/handler/admin_variant_attributes.go"
    "internal/proj/marketplace/handler/categories.go"
    "internal/proj/marketplace/handler/custom_components.go"
    "internal/proj/marketplace/handler/favorites.go"
    "internal/proj/marketplace/handler/fuzzy_search.go"
    "internal/proj/marketplace/handler/handler.go"
    "internal/proj/marketplace/handler/images.go"
    "internal/proj/marketplace/handler/indexing.go"
    "internal/proj/marketplace/handler/saved_searches.go"
    "internal/proj/marketplace/handler/search.go"
    "internal/proj/marketplace/handler/translation.go"
    "internal/proj/marketplace/handler/unified_attributes_test.go"
    "internal/proj/marketplace/handler/variant_mappings.go"
    "internal/proj/notifications/handler/handler.go"
    "internal/proj/payments/handler/handler.go"
    "internal/proj/postexpress/handler/handler.go"
    "internal/proj/search_admin/handler/handler.go"
    "internal/proj/storefronts/handler/staff_handler.go"
    "internal/proj/storefronts/handler/storefront_handler.go"
    "internal/proj/users/handler/users.go"
    "internal/proj/viber/handler.go"
)

for file in "${FILES[@]}"; do
    echo "Processing $file..."

    # Создаём временный файл
    tmpfile=$(mktemp)

    # Флаг для отслеживания, встретили ли мы первый импорт utils
    utils_found=false

    # Читаем файл построчно
    while IFS= read -r line; do
        # Если строка содержит импорт utils
        if [[ "$line" == *'"backend/pkg/utils"' ]]; then
            if [ "$utils_found" = false ]; then
                # Это первый импорт - оставляем его
                echo "$line" >> "$tmpfile"
                utils_found=true
            else
                # Это дублированный импорт - пропускаем его
                echo "  Removing duplicate import: $line"
            fi
        else
            # Обычная строка - копируем как есть
            echo "$line" >> "$tmpfile"
        fi
    done < "$file"

    # Заменяем оригинальный файл
    mv "$tmpfile" "$file"
    echo "  Fixed $file"
done

echo "All duplicate imports removed!"