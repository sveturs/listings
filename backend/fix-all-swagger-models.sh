#!/bin/bash

echo "Исправление всех оставшихся ссылок на модели в swagger..."

# Находим все файлы с неправильными ссылками и исправляем их
find /data/hostel-booking-system/backend -name "*.go" -type f | while read file; do
    # Пропускаем vendor и другие системные директории
    if [[ "$file" == *"/vendor/"* ]] || [[ "$file" == *"/.git/"* ]]; then
        continue
    fi

    # Находим модуль из пути файла
    if [[ "$file" == */internal/proj/*/handler/*.go ]]; then
        # Извлекаем имя модуля из пути
        module=$(echo "$file" | sed 's|.*/internal/proj/\([^/]*\)/.*|\1|')

        # Заменяем backend_internal_domain_models на правильный путь к локальным моделям
        sed -i "s|backend_internal_domain_models\.|backend_internal_proj_${module}_models.|g" "$file"

        # Специальные случаи для logistics и других общих моделей
        sed -i "s|backend_internal_proj_${module}_models\.DeliveryProblem|backend_internal_domain_logistics.DeliveryProblem|g" "$file"
        sed -i "s|backend_internal_proj_${module}_models\.AnalyticsData|backend_internal_domain_logistics.AnalyticsData|g" "$file"
        sed -i "s|backend_internal_proj_${module}_models\.DashboardData|backend_internal_domain_logistics.DashboardData|g" "$file"

        echo "Обработан файл: $file (модуль: $module)"
    fi
done

echo "Исправления завершены!"