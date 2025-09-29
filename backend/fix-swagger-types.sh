#!/bin/bash

# Скрипт для исправления неправильных ссылок на типы в Swagger аннотациях
# Исправляет backend_internal_domain_models.* на правильные пути к модулям

echo "🔧 Исправление Swagger аннотаций..."

# Директория backend
BACKEND_DIR="/data/hostel-booking-system/backend"

# Найти все файлы с неправильными ссылками на типы
FILES_TO_FIX=$(grep -r "backend_internal_domain_models\." "$BACKEND_DIR" --include="*.go" -l)

echo "📂 Найдено файлов для исправления: $(echo "$FILES_TO_FIX" | wc -l)"

# Функция для определения правильного пути типа
get_correct_type_path() {
    local type_name=$1

    # BEX Express типы
    if grep -q "^type $type_name struct" "$BACKEND_DIR/internal/proj/bexexpress/models/models.go" 2>/dev/null; then
        echo "backend_internal_proj_bexexpress_models.$type_name"
        return
    fi

    # Post Express типы
    if grep -q "^type $type_name struct" "$BACKEND_DIR/internal/proj/postexpress/models/models.go" 2>/dev/null; then
        echo "backend_internal_proj_postexpress_models.$type_name"
        return
    fi

    # Delivery типы
    if grep -q "^type $type_name struct" "$BACKEND_DIR/internal/proj/delivery/models/models.go" 2>/dev/null; then
        echo "backend_internal_proj_delivery_models.$type_name"
        return
    fi

    # Delivery admin типы
    if grep -q "^type $type_name struct" "$BACKEND_DIR/internal/proj/delivery/models/admin_types.go" 2>/dev/null; then
        echo "backend_internal_proj_delivery_models.$type_name"
        return
    fi

    # VIN типы
    if grep -q "^type $type_name struct" "$BACKEND_DIR/internal/proj/vin/models/models.go" 2>/dev/null; then
        echo "backend_internal_proj_vin_models.$type_name"
        return
    fi

    # Viber типы
    if grep -q "^type $type_name struct" "$BACKEND_DIR/internal/proj/viber/models/models.go" 2>/dev/null; then
        echo "backend_internal_proj_viber_models.$type_name"
        return
    fi

    # Если тип не найден в модулях, проверим domain/models
    if find "$BACKEND_DIR/internal/domain/models/" -name "*.go" -exec grep -q "^type $type_name struct" {} \; 2>/dev/null; then
        echo "backend_internal_domain_models.$type_name"
        return
    fi

    # Если не найден, возвращаем оригинальный путь
    echo "backend_internal_domain_models.$type_name"
}

# Создать карту замен
declare -A replacements

# Найти все типы, которые используются с неправильными путями
echo "🔍 Анализ типов..."

# Извлечь все типы из аннотаций
TYPES_USED=$(grep -r "backend_internal_domain_models\." "$BACKEND_DIR" --include="*.go" -o | sed 's/backend_internal_domain_models\.//' | sort -u)

for type_name in $TYPES_USED; do
    # Удалить возможные специальные символы
    clean_type=$(echo "$type_name" | sed 's/[^a-zA-Z0-9_].*$//')
    if [[ -n "$clean_type" ]]; then
        correct_path=$(get_correct_type_path "$clean_type")
        if [[ "$correct_path" != "backend_internal_domain_models.$clean_type" ]]; then
            replacements["backend_internal_domain_models.$clean_type"]="$correct_path"
            echo "  ✓ $clean_type -> $correct_path"
        fi
    fi
done

echo ""
echo "🔄 Выполнение замен..."

# Применить замены ко всем файлам
for file in $FILES_TO_FIX; do
    echo "📝 Обрабатываю: $file"

    # Создать временный файл
    temp_file=$(mktemp)
    cp "$file" "$temp_file"

    # Применить все замены
    for old_path in "${!replacements[@]}"; do
        new_path="${replacements[$old_path]}"
        sed -i "s|$old_path|$new_path|g" "$temp_file"
    done

    # Проверить, были ли изменения
    if ! diff -q "$file" "$temp_file" > /dev/null; then
        mv "$temp_file" "$file"
        echo "  ✅ Изменен"
    else
        rm "$temp_file"
        echo "  ➖ Без изменений"
    fi
done

echo ""
echo "✨ Исправление завершено!"

# Показать статистику
echo ""
echo "📊 Статистика замен:"
for old_path in "${!replacements[@]}"; do
    new_path="${replacements[$old_path]}"
    count=$(grep -r "$new_path" "$BACKEND_DIR" --include="*.go" | wc -l)
    echo "  • $old_path -> $new_path ($count использований)"
done

echo ""
echo "🔍 Проверка оставшихся проблем..."
remaining=$(grep -r "backend_internal_domain_models\." "$BACKEND_DIR" --include="*.go" -l | wc -l)
if [[ $remaining -eq 0 ]]; then
    echo "✅ Все ссылки исправлены!"
else
    echo "⚠️  Осталось файлов с проблемами: $remaining"
    echo "Файлы:"
    grep -r "backend_internal_domain_models\." "$BACKEND_DIR" --include="*.go" -l
fi