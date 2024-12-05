#!/bin/bash

OUTPUT_FILE="collected_code.txt"

# Очищаем или создаем файл
> "$OUTPUT_FILE"

# Функция для добавления разделителя и содержимого файла
add_file() {
    echo -e "\n=== $1 ===\n" >> "$OUTPUT_FILE"
    if [ -f "$1" ]; then
        cat "$1" >> "$OUTPUT_FILE"
    else
        echo "// TODO: File not found: $1" >> "$OUTPUT_FILE"
    fi
}

# Создаем структуру каталогов, если она не существует
mkdir -p cmd/api
mkdir -p internal/{config,server,handlers,middleware,domain/models,services,storage/postgres}
mkdir -p pkg/{logger,utils}

# Собираем файлы в нужном порядке
add_file "cmd/api/main.go"

add_file "internal/config/config.go"

add_file "internal/server/server.go"

add_file "internal/handlers/handler.go"
add_file "internal/handlers/auth.go"
add_file "internal/handlers/rooms.go"
add_file "internal/handlers/bookings.go"
add_file "internal/handlers/users.go"

add_file "internal/middleware/middleware.go"
add_file "internal/middleware/auth.go"
add_file "internal/middleware/cors.go"
add_file "internal/middleware/logger.go"

add_file "internal/domain/models/models.go"
add_file "internal/domain/models/booking.go"

add_file "internal/services/services.go"
add_file "internal/services/auth.go"
add_file "internal/services/room.go"
add_file "internal/services/booking.go"
add_file "internal/services/user.go"

add_file "internal/storage/postgres/database.go"
add_file "internal/storage/postgres/storage.go"
add_file "internal/storage/postgres/user.go"
add_file "internal/storage/postgres/room.go"
add_file "internal/storage/postgres/booking.go"
add_file "internal/storage/postgres/image.go"
add_file "internal/storage/postgres/bed.go"

add_file "pkg/logger/logger.go"
add_file "pkg/utils/utils.go"

echo "Files have been collected in $OUTPUT_FILE"