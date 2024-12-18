#!/bin/bash

# Конфигурация
output_file="project_code.txt"

# Массив директорий для поиска
project_dirs=(
    "/data/proj/hostel-booking-system/"
)

# Массив расширений файлов для поиска
file_extensions=(
    "go"
#    "js"
    "sql"
#    "jsx"
    "env"
#    "tsx"
    "json"
    "yml"
#    "yaml"
)

# Массив конкретных файлов для добавления
specific_files=(
    "/data/proj/hostel-booking-system/backend/.env.local"
    "/data/proj/hostel-booking-system/backend/Dockerfile"
#   "/data/proj/hostel-booking-system/backend/docker-compose.yml"
    "nginx.conf"
    "deploy.sh"
#    "/data/proj/hostel-booking-system/frontend/hostel-frontend/.env"
#    "/data/proj/hostel-booking-system/frontend/hostel-frontend/.env.local"
    "/data/proj/hostel-booking-system/backend/.env"
)

# Массив исключаемых директорий
exclude_dirs=(
    "node_modules"
    "vendor"
    "dist"
    "build"
    ".git"
    "uploads"
    "/data/proj/hostel-booking-system/frontend/"
    "/data/proj/hostel-booking-system/frontend/hostel-frontend/"
#    "/data/proj/hostel-booking-system/frontend/node_modules/"
#    "/data/proj/hostel-booking-system/frontend/hostel-frontend/build/"
    "/data/proj/hostel-booking-system/node_modules/"
#   "/data/proj/hostel-booking-system/frontend/"

)

# Массив исключаемых файлов (по имени файла или полному пути)
exclude_files=(
    "check_imports.sh"
    "collect_code.sh"
    "/data/proj/hostel-booking-system/frontend/hostel-frontend/build/static/js/main.2d72e40b.js"
    ".DS_Store"
    "/data/proj/hostel-booking-system/backend/some/specific/file.go"
    "package-lock.json"
    "/data/proj/hostel-booking-system/frontend/package-lock.json"
    "/data/proj/hostel-booking-system/backend/main.go"
    "/data/proj/hostel-booking-system/backend/database/db.go"
    "/data/proj/hostel-booking-system/package-lock.json"
    "/data/proj/hostel-booking-system/test/main_test.go"
    


    
)

# Очистка выходного файла
> "$output_file"

# Функция для проверки исключаемых директорий
should_exclude_dir() {
    local path="$1"
    for exclude_dir in "${exclude_dirs[@]}"; do
        if [[ $path == *"/$exclude_dir/"* ]] || [[ $path == *"/$exclude_dir" ]]; then
            return 0
        fi
    done
    return 1
}

# Функция для проверки исключаемых файлов
should_exclude_file() {
    local path="$1"
    local filename=$(basename "$path")
    for exclude_file in "${exclude_files[@]}"; do
        if [[ "$path" == "$exclude_file" ]] || [[ "$filename" == "$exclude_file" ]]; then
            return 0
        fi
    done
    return 1
}

# Обработка конкретных файлов
echo "Processing specific files..."
for file in "${specific_files[@]}"; do
    if [ -f "$file" ] && ! should_exclude_file "$file"; then
        echo "=== $file ===" >> "$output_file"
        echo "" >> "$output_file"
        cat "$file" >> "$output_file"
        echo "" >> "$output_file"
        echo "Added specific file: $file"
    else
        echo "Skipped or not found: $file"
    fi
done

# Проходим по всем указанным директориям
for search_dir in "${project_dirs[@]}"; do
    echo "Processing directory: $search_dir"
    
    # Поиск и обработка файлов по расширениям
    for ext in "${file_extensions[@]}"; do
        while IFS= read -r -d '' file; do
            if ! should_exclude_dir "$file" && ! should_exclude_file "$file"; then
                echo "=== $file ===" >> "$output_file"
                echo "" >> "$output_file"
                cat "$file" >> "$output_file"
                echo "" >> "$output_file"
                echo "Added: $file"
            else
                echo "Excluded: $file"
            fi
        done < <(find "$search_dir" -type f -name "*.$ext" -print0)
    done
done

# Вывод статистики
echo "==============================================="
echo "Files have been collected in $output_file"
echo "Total size: $(wc -l < "$output_file") lines"
echo "File size: $(du -h "$output_file" | cut -f1)"
