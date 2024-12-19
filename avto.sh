#!/bin/bash

# Конфигурация
output_file="add_project_code.txt"

# Массив директорий для поиска
project_dirs=(
    "backend/ampty"
    
)

# Массив расширений файлов для поиска
file_extensions=(
    "go"
    "js"
    "sql"
    "jsx"
#    "ts"
#    "tsx"
    "json"
    "yml"
#    "yaml"
)

# Массив конкретных файлов для добавления
specific_files=(
# общие серверное
"backend/internal/handlers/car.go"
"backend/internal/services/car.go"
"backend/internal/storage/postgres/car_image.go"
"backend/internal/storage/postgres/car.go"

"backend/internal/storage/postgres/image.go"
"frontend/hostel-frontend/src/components/Car/BookingDialog.js"
"frontend/hostel-frontend/src/components/CarMapView.js"
"frontend/hostel-frontend/src/pages/car/AddCarPage.js"
"frontend/hostel-frontend/src/pages/car/CarListPage.js"

)

# Массив исключаемых директорий
exclude_dirs=(
    "node_modules"
    "vendor"
    "dist"
    "build"
    ".git"
    "uploads"
    "/data/proj/hostel-booking-system/frontend/node_modules/"
    "/data/proj/hostel-booking-system/frontend/hostel-frontend/build/"
    "/data/proj/hostel-booking-system/node_modules/"

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
