#!/bin/bash

# Простой скрипт для исправления конкретных типов в Swagger аннотациях

echo "🔧 Исправление Swagger аннотаций (простая версия)..."

BACKEND_DIR="/data/hostel-booking-system/backend"

# Функция для выполнения замены во всех файлах
replace_type() {
    local old_type=$1
    local new_type=$2

    echo "🔄 Заменяю $old_type на $new_type"

    # Найти все файлы с этим типом и заменить
    files_changed=0
    while IFS= read -r -d '' file; do
        if grep -q "$old_type" "$file"; then
            sed -i "s|$old_type|$new_type|g" "$file"
            echo "  ✅ Изменен: $file"
            ((files_changed++))
        fi
    done < <(find "$BACKEND_DIR" -name "*.go" -print0)

    if [[ $files_changed -eq 0 ]]; then
        echo "  ➖ Не найдено файлов для изменения"
    else
        echo "  📊 Изменено файлов: $files_changed"
    fi
    echo ""
}

# BEX Express типы
replace_type "backend_internal_domain_models.CalculateRateRequest" "backend_internal_proj_bexexpress_models.CalculateRateRequest"
replace_type "backend_internal_domain_models.CalculateRateResponse" "backend_internal_proj_bexexpress_models.CalculateRateResponse"
replace_type "backend_internal_domain_models.SearchAddressRequest" "backend_internal_proj_bexexpress_models.SearchAddressRequest"
replace_type "backend_internal_domain_models.AddressSuggestion" "backend_internal_proj_bexexpress_models.AddressSuggestion"
replace_type "backend_internal_domain_models.BEXParcelShop" "backend_internal_proj_bexexpress_models.BEXParcelShop"

# Post Express типы
replace_type "backend_internal_domain_models.PostExpressSettings" "backend_internal_proj_postexpress_models.PostExpressSettings"
replace_type "backend_internal_domain_models.PostExpressLocation" "backend_internal_proj_postexpress_models.PostExpressLocation"
replace_type "backend_internal_domain_models.PostExpressOffice" "backend_internal_proj_postexpress_models.PostExpressOffice"
replace_type "backend_internal_domain_models.PostExpressRate" "backend_internal_proj_postexpress_models.PostExpressRate"
replace_type "backend_internal_domain_models.PostExpressShipment" "backend_internal_proj_postexpress_models.PostExpressShipment"
replace_type "backend_internal_domain_models.CreateShipmentRequest" "backend_internal_proj_postexpress_models.CreateShipmentRequest"
replace_type "backend_internal_domain_models.TrackingEvent" "backend_internal_proj_postexpress_models.TrackingEvent"
replace_type "backend_internal_domain_models.Warehouse" "backend_internal_proj_postexpress_models.Warehouse"
replace_type "backend_internal_domain_models.WarehousePickupOrder" "backend_internal_proj_postexpress_models.WarehousePickupOrder"
replace_type "backend_internal_domain_models.CreatePickupOrderRequest" "backend_internal_proj_postexpress_models.CreatePickupOrderRequest"

echo "✨ Исправление завершено!"

# Проверка оставшихся проблем
echo ""
echo "🔍 Проверка оставшихся проблем..."
remaining_files=$(find "$BACKEND_DIR" -name "*.go" -exec grep -l "backend_internal_domain_models\." {} \; | wc -l)
if [[ $remaining_files -eq 0 ]]; then
    echo "✅ Все ссылки исправлены!"
else
    echo "⚠️  Осталось файлов с проблемами: $remaining_files"
    echo ""
    echo "Оставшиеся типы:"
    find "$BACKEND_DIR" -name "*.go" -exec grep -o "backend_internal_domain_models\.[A-Za-z0-9_]*" {} \; | sort -u | head -10
fi