#!/bin/bash

# Финальный скрипт для исправления всех оставшихся ссылок в Swagger аннотациях

echo "🔧 Финальное исправление Swagger аннотаций..."

BACKEND_DIR="/data/hostel-booking-system/backend"

# Исправить utils.SuccessResponse и utils.ErrorResponse во всех файлах
echo "🔄 Исправляю utils.SuccessResponse и utils.ErrorResponse..."
find "$BACKEND_DIR" -name "*.go" -exec sed -i 's|utils\.SuccessResponse|backend_pkg_utils.SuccessResponseSwag|g' {} \;
find "$BACKEND_DIR" -name "*.go" -exec sed -i 's|utils\.ErrorResponse|backend_pkg_utils.ErrorResponseSwag|g' {} \;

echo "✅ Исправления завершены!"

# Проверим, что осталось
echo ""
echo "🔍 Проверка оставшихся проблем..."

# Проверить utils ссылки
utils_issues=$(find "$BACKEND_DIR" -name "*.go" -exec grep -l "utils\..*Response" {} \; 2>/dev/null | wc -l)
if [[ $utils_issues -gt 0 ]]; then
    echo "⚠️  Осталось файлов с utils.*Response: $utils_issues"
    find "$BACKEND_DIR" -name "*.go" -exec grep -l "utils\..*Response" {} \; 2>/dev/null | head -3
else
    echo "✅ Все utils.*Response исправлены"
fi

# Проверить backend_internal_domain_models ссылки
domain_issues=$(find "$BACKEND_DIR" -name "*.go" -exec grep -l "backend_internal_domain_models\." {} \; 2>/dev/null | wc -l)
echo "📊 Файлов с backend_internal_domain_models.*: $domain_issues (это нормально для общих типов)"

echo ""
echo "🚀 Пробуем сгенерировать типы..."