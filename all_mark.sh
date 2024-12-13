#!/bin/bash

# Конфигурация
output_file="market_project_code.txt"

# Массив директорий для поиска
project_dirs=(
    "backend/migrations"
    
)

# Массив расширений файлов для поиска
file_extensions=(
    "go"
    "js"
    "sql"
#    "jsx"
#    "ts"
#    "tsx"
    "json"
    "yml"
#    "yaml"
)

# Массив конкретных файлов для добавления
specific_files=(
# общие серверное
"backend/cmd/api/main.go"
"backend/internal/config/config.go"
"backend/internal/domain/models/models.go"

"backend/internal/handlers/auth.go"
"backend/internal/handlers/handler.go"
 

"backend/internal/handlers/users.go"
"backend/internal/middleware/auth.go"
"backend/internal/middleware/cors.go"
"backend/internal/middleware/logger.go"
"backend/internal/middleware/middleware.go"
"backend/internal/server/server.go"
"backend/internal/services/auth.go"
"backend/internal/services/interfaces.go"
 
"backend/internal/services/services.go"
"backend/internal/services/user.go"
"backend/internal/storage/postgres/db.go"
"backend/internal/storage/postgres/user.go"
"backend/internal/storage/storage.go"
"backend/internal/types/auth.go"
"backend/pkg/logger/logger.go"
"backend/pkg/utils/utils.go"
"backend/.env"
"backend/.env.local"
"backend/Dockerfile"
".gitignore"
"deploy.sh"
"docker-compose.override.yml"
"docker-compose.prod.yml"
"docker-compose.yml"
"nginx.conf"
"nginx.local.conf"
# сервер маркет
"backend/internal/domain/models/review.go"
"backend/internal/handlers/marketplace.go"
"backend/internal/handlers/reviews.go"
"backend/internal/services/marketplace.go"
"backend/internal/services/review.go"
"backend/internal/storage/postgres/marketplace.go"
"backend/internal/storage/postgres/reviews.go"
# общие фронт
"frontend/hostel-frontend/package.json"
"frontend/hostel-frontend/.env"
"frontend/hostel-frontend/.env.local"
"frontend/hostel-frontend/Dockerfile"
"frontend/hostel-frontend/src/App.css"
"frontend/hostel-frontend/src/App.js"
"frontend/hostel-frontend/src/index.css"
"frontend/hostel-frontend/src/index.js"
"frontend/hostel-frontend/src/setupProxy.js"
"frontend/hostel-frontend/src/api/axios.js"
"frontend/hostel-frontend/src/components/Layout.js"
"frontend/hostel-frontend/src/contexts/AuthContext.js"

# фронт маркет
"frontend/hostel-frontend/src/components/marketplace/CategoryTree.js"
"frontend/hostel-frontend/src/components/marketplace/ItemDetails.js"
"frontend/hostel-frontend/src/components/marketplace/ListingCard.js"
"frontend/hostel-frontend/src/components/marketplace/MarketplaceFilters.js"
"frontend/hostel-frontend/src/pages/CreateListingPage.js"
"frontend/hostel-frontend/src/pages/MarketplacePage.js"
"frontend/hostel-frontend/src/components/reviews/ReviewCard.js"
"frontend/hostel-frontend/src/components/reviews/ReviewForm.js"
"frontend/hostel-frontend/src/components/reviews/ReviewsSection.js"
"frontend/hostel-frontend/src/components/Review.js"
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
