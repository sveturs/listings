#!/bin/bash

echo "Умное исправление типов в Swagger аннотациях..."

# Функция для поиска правильного пути к типу
find_correct_type() {
    local type_name=$1
    local found_path=""

    # Ищем определение типа
    found_path=$(grep -r "type $type_name struct" --include="*.go" /data/hostel-booking-system/backend 2>/dev/null | head -1 | cut -d: -f1)

    if [ -n "$found_path" ]; then
        # Конвертируем путь файла в путь пакета для swagger
        echo "$found_path" | sed 's|/data/hostel-booking-system/backend/||' | sed 's|/[^/]*\.go$||' | sed 's|/|_|g' | sed 's|^|backend_|'
    fi
}

# Исправляем analytics handler
sed -i 's/backend_internal_proj_marketplace_service\.SearchMetrics/backend_internal_proj_analytics_service.SearchMetrics/g' \
    /data/hostel-booking-system/backend/internal/proj/analytics/handler/analytics_handler.go
sed -i 's/backend_internal_proj_marketplace_service\.ItemPerformance/backend_internal_proj_analytics_service.ItemPerformance/g' \
    /data/hostel-booking-system/backend/internal/proj/analytics/handler/analytics_handler.go

# Исправляем delivery handler
# CreateShipmentRequest определен в delivery/models, не в service
sed -i 's/backend_internal_proj_marketplace_service\.CreateShipmentRequest/backend_internal_proj_delivery_models.CreateShipmentRequest/g' \
    /data/hostel-booking-system/backend/internal/proj/delivery/handler/handler.go

# Исправляем неправильные ссылки на service в коде (не в swagger)
find /data/hostel-booking-system/backend/internal/proj -name "*.go" -type f | while read file; do
    # В коде Go заменяем backend_internal_proj_marketplace_service. на service. только в вызовах методов
    sed -i 's/h\.backend_internal_proj_marketplace_service\./h.service./g' "$file"
    sed -i 's/backend_internal_proj_marketplace_service\./service./g' "$file" 2>/dev/null
done

echo "Готово!"