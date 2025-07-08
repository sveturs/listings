#!/bin/bash

# Скрипт для нагрузочного тестирования поисковой системы
# Использует Apache Bench (ab) для тестирования различных endpoint'ов

BASE_URL="http://localhost:3000"
LOG_DIR="/tmp/load_test_logs"
RESULTS_DIR="/tmp/load_test_results"

# Создаем директории для логов
mkdir -p "$LOG_DIR"
mkdir -p "$RESULTS_DIR"

echo "=== Нагрузочное тестирование поисковой системы ==="
echo "Время начала: $(date)"
echo "URL: $BASE_URL"

# Функция для выполнения GET запроса с параметрами
test_get_endpoint() {
    local endpoint="$1"
    local params="$2"
    local concurrent="$3"
    local requests="$4"
    local test_name="$5"
    
    local url="${BASE_URL}${endpoint}"
    if [ -n "$params" ]; then
        url="${url}?${params}"
    fi
    
    echo "Тестирование: $test_name"
    echo "URL: $url"
    echo "Concurrent: $concurrent, Requests: $requests"
    
    # Запускаем Apache Bench
    ab -n "$requests" -c "$concurrent" -g "$RESULTS_DIR/${test_name}_${concurrent}c_${requests}r.tsv" \
       -e "$RESULTS_DIR/${test_name}_${concurrent}c_${requests}r.csv" \
       "$url" > "$LOG_DIR/${test_name}_${concurrent}c_${requests}r.log" 2>&1
    
    # Показываем краткие результаты
    echo "--- Результаты для $test_name ---"
    grep -E "Time per request|Requests per second|Transfer rate|Failed requests" "$LOG_DIR/${test_name}_${concurrent}c_${requests}r.log"
    echo ""
}

# Функция для выполнения POST запроса
test_post_endpoint() {
    local endpoint="$1"
    local data="$2"
    local concurrent="$3"
    local requests="$4"
    local test_name="$5"
    
    local url="${BASE_URL}${endpoint}"
    
    echo "Тестирование: $test_name"
    echo "URL: $url"
    echo "Concurrent: $concurrent, Requests: $requests"
    
    # Создаем временный файл с данными
    local temp_file="/tmp/post_data_${test_name}.json"
    echo "$data" > "$temp_file"
    
    # Запускаем Apache Bench для POST
    ab -n "$requests" -c "$concurrent" -p "$temp_file" -T "application/json" \
       -g "$RESULTS_DIR/${test_name}_${concurrent}c_${requests}r.tsv" \
       -e "$RESULTS_DIR/${test_name}_${concurrent}c_${requests}r.csv" \
       "$url" > "$LOG_DIR/${test_name}_${concurrent}c_${requests}r.log" 2>&1
    
    # Показываем краткие результаты
    echo "--- Результаты для $test_name ---"
    grep -E "Time per request|Requests per second|Transfer rate|Failed requests" "$LOG_DIR/${test_name}_${concurrent}c_${requests}r.log"
    echo ""
    
    # Удаляем временный файл
    rm -f "$temp_file"
}

# Функция для мониторинга ресурсов
start_monitoring() {
    echo "Запуск мониторинга ресурсов..."
    
    # Мониторинг CPU и памяти
    while true; do
        echo "$(date),$(ps -p $1 -o %cpu,%mem --no-headers)" >> "$LOG_DIR/resource_monitoring.csv"
        sleep 1
    done &
    
    MONITOR_PID=$!
    echo "Мониторинг запущен с PID: $MONITOR_PID"
}

stop_monitoring() {
    if [ -n "$MONITOR_PID" ]; then
        kill $MONITOR_PID 2>/dev/null
        echo "Мониторинг остановлен"
    fi
}

# Получаем PID процесса Go
GO_PID=$(lsof -ti :3000)
if [ -z "$GO_PID" ]; then
    echo "Ошибка: Backend не запущен на порту 3000"
    exit 1
fi

echo "PID процесса Go: $GO_PID"

# Заголовок для CSV файла мониторинга
echo "timestamp,cpu_percent,memory_percent" > "$LOG_DIR/resource_monitoring.csv"

# Тесты 1: Простой поиск
echo "=== ТЕСТ 1: Простой поиск ==="
start_monitoring $GO_PID

test_get_endpoint "/api/v1/marketplace/search" "query=телефон&size=20&sort=relevance" 50 500 "simple_search_50"
test_get_endpoint "/api/v1/marketplace/search" "query=телефон&size=20&sort=relevance" 100 1000 "simple_search_100"
test_get_endpoint "/api/v1/marketplace/search" "query=телефон&size=20&sort=relevance" 200 2000 "simple_search_200"

stop_monitoring

# Тесты 2: Сложный поиск с фильтрами
echo "=== ТЕСТ 2: Сложный поиск с фильтрами ==="
start_monitoring $GO_PID

test_get_endpoint "/api/v1/marketplace/search" "query=iphone&category_id=1&min_price=50000&max_price=150000&size=20&sort=price_asc" 50 500 "complex_search_50"
test_get_endpoint "/api/v1/marketplace/search" "query=iphone&category_id=1&min_price=50000&max_price=150000&size=20&sort=price_asc" 100 1000 "complex_search_100"
test_get_endpoint "/api/v1/marketplace/search" "query=iphone&category_id=1&min_price=50000&max_price=150000&size=20&sort=price_asc" 200 2000 "complex_search_200"

stop_monitoring

# Тесты 3: Автодополнение
echo "=== ТЕСТ 3: Автодополнение ==="
start_monitoring $GO_PID

test_get_endpoint "/api/v1/marketplace/suggestions" "prefix=теле&size=10" 50 500 "suggestions_50"
test_get_endpoint "/api/v1/marketplace/suggestions" "prefix=теле&size=10" 100 1000 "suggestions_100"
test_get_endpoint "/api/v1/marketplace/suggestions" "prefix=теле&size=10" 200 2000 "suggestions_200"

stop_monitoring

# Тесты 4: Похожие объявления
echo "=== ТЕСТ 4: Похожие объявления ==="
start_monitoring $GO_PID

test_get_endpoint "/api/v1/marketplace/listings/1/similar" "limit=5" 50 500 "similar_50"
test_get_endpoint "/api/v1/marketplace/listings/1/similar" "limit=5" 100 1000 "similar_100"
test_get_endpoint "/api/v1/marketplace/listings/1/similar" "limit=5" 200 2000 "similar_200"

stop_monitoring

# Тесты 5: Разные размеры результатов
echo "=== ТЕСТ 5: Разные размеры результатов ==="
start_monitoring $GO_PID

test_get_endpoint "/api/v1/marketplace/search" "query=автомобиль&size=10" 50 500 "size_10"
test_get_endpoint "/api/v1/marketplace/search" "query=автомобиль&size=50" 50 500 "size_50"
test_get_endpoint "/api/v1/marketplace/search" "query=автомобиль&size=100" 50 500 "size_100"

stop_monitoring

echo "=== Тестирование завершено ==="
echo "Время окончания: $(date)"
echo "Логи сохранены в: $LOG_DIR"
echo "Результаты сохранены в: $RESULTS_DIR"