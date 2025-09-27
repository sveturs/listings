#!/bin/bash

# Скрипт для исправления типов в Swagger аннотациях

echo "Исправляю типы в Swagger аннотациях..."

# Находим все Go файлы с проблемными аннотациями
find /data/hostel-booking-system/backend/internal/proj -name "*.go" -type f | while read file; do
    # Проверяем, есть ли в файле проблемные паттерны
    if grep -q "@Success.*models\." "$file" 2>/dev/null || \
       grep -q "@Param.*models\." "$file" 2>/dev/null || \
       grep -q "@Failure.*models\." "$file" 2>/dev/null || \
       grep -q "data=models\." "$file" 2>/dev/null; then

        echo "Обрабатываю: $file"

        # Создаем временный файл
        cp "$file" "$file.bak"

        # Заменяем models. на backend_internal_domain_models.
        sed -i 's/models\./backend_internal_domain_models./g' "$file"

        # Также исправляем другие распространенные паттерны
        sed -i 's/handler\.\([A-Z]\)/internal_proj_marketplace_handler.\1/g' "$file"
        sed -i 's/service\.\([A-Z]\)/backend_internal_proj_marketplace_service.\1/g' "$file"
        sed -i 's/domain\.\([A-Z]\)/backend_internal_domain.\1/g' "$file"

        # Проверяем, изменился ли файл
        if ! diff -q "$file" "$file.bak" > /dev/null; then
            echo "  ✓ Исправлен"
            rm "$file.bak"
        else
            echo "  - Без изменений"
            rm "$file.bak"
        fi
    fi
done

echo "Готово!"