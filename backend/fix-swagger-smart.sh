#!/bin/bash

# Умный скрипт для исправления Swagger аннотаций
# Заменяет типы в контексте соответствующих модулей

echo "🔧 Исправление Swagger аннотаций (умная версия)..."

BACKEND_DIR="/data/hostel-booking-system/backend"

# Функция для выполнения замены в конкретных файлах модуля
replace_in_module() {
    local module_path=$1
    local old_prefix=$2
    local new_prefix=$3

    echo "🔄 Исправляю типы в модуле: $module_path"

    files_changed=0
    while IFS= read -r -d '' file; do
        if grep -q "$old_prefix" "$file"; then
            # Заменить все типы с данным префиксом
            sed -i "s|${old_prefix}|${new_prefix}|g" "$file"
            echo "  ✅ Изменен: $file"
            ((files_changed++))
        fi
    done < <(find "$BACKEND_DIR/internal/proj/$module_path" -name "*.go" -print0)

    if [[ $files_changed -eq 0 ]]; then
        echo "  ➖ Не найдено файлов для изменения в модуле $module_path"
    else
        echo "  📊 Изменено файлов в модуле $module_path: $files_changed"
    fi
    echo ""
}

# Исправить типы в модуле bexexpress
replace_in_module "bexexpress" "backend_internal_domain_models." "backend_internal_proj_bexexpress_models."

# Исправить типы в модуле postexpress
replace_in_module "postexpress" "backend_internal_domain_models." "backend_internal_proj_postexpress_models."

# Исправить типы в модуле delivery
replace_in_module "delivery" "backend_internal_domain_models." "backend_internal_proj_delivery_models."

# Исправить типы в модуле vin
replace_in_module "vin" "backend_internal_domain_models." "backend_internal_proj_vin_models."

# Исправить типы в модуле viber
replace_in_module "viber" "backend_internal_domain_models." "backend_internal_proj_viber_models."

# Для storefronts нужно проверить, какие типы действительно есть в отдельных файлах
echo "🔍 Анализ типов в storefronts..."
storefront_types_found=0

# Найти все типы, определенные в storefronts handlers
while IFS= read -r -d '' file; do
    if grep -q "type.*struct" "$file"; then
        echo "  📄 Найдены типы в: $file"
        grep "^type.*struct" "$file" | while read line; do
            type_name=$(echo "$line" | awk '{print $2}')
            echo "    • $type_name"
        done
        ((storefront_types_found++))
    fi
done < <(find "$BACKEND_DIR/internal/proj/storefronts" -name "*.go" -print0)

if [[ $storefront_types_found -gt 0 ]]; then
    echo "  ⚠️  В storefronts найдены локальные типы, но нет отдельной models директории"
    echo "  💡 Рекомендуется создать storefronts/models/ и перенести туда типы"
fi

echo ""
echo "✨ Исправление завершено!"

# Теперь нужно обновить сгенерированный docs.go файл
echo "🔄 Обновление сгенерированной документации..."
cd "$BACKEND_DIR"
make generate-types

echo ""
echo "🔍 Проверка оставшихся проблем..."
remaining_files=$(find "$BACKEND_DIR" -name "*.go" -exec grep -l "backend_internal_domain_models\." {} \; 2>/dev/null | wc -l)
if [[ $remaining_files -eq 0 ]]; then
    echo "✅ Все ссылки исправлены!"
else
    echo "⚠️  Осталось файлов с проблемами: $remaining_files"
    echo ""
    echo "Файлы с оставшимися проблемами:"
    find "$BACKEND_DIR" -name "*.go" -exec grep -l "backend_internal_domain_models\." {} \; 2>/dev/null | head -5
    echo ""
    echo "Первые 5 оставшихся типов:"
    find "$BACKEND_DIR" -name "*.go" -exec grep -o "backend_internal_domain_models\.[A-Za-z0-9_]*" {} \; 2>/dev/null | sort -u | head -5
fi