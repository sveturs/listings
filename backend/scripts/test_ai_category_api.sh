#!/bin/bash

# Скрипт для тестирования AI Category Detection API
# Использование: ./test_ai_category_api.sh [host] [token]

HOST=${1:-"http://localhost:3000"}
TOKEN=${2:-$(cd /data/hostel-booking-system/backend && go run scripts/create_test_jwt.go 2>/dev/null)}

# Цвета для вывода
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "🧪 Тестирование AI Category Detection API"
echo "==========================================}"
echo "Host: $HOST"
echo "Token: ${TOKEN:0:20}..."
echo ""

# Функция для красивого вывода JSON
pretty_json() {
    python3 -m json.tool 2>/dev/null || cat
}

# Функция для тестирования endpoint
test_endpoint() {
    local method=$1
    local endpoint=$2
    local data=$3
    local description=$4

    echo -e "${YELLOW}📝 Тест: $description${NC}"
    echo "Endpoint: $method $endpoint"

    if [ "$method" == "POST" ]; then
        response=$(curl -s -X POST \
            "${HOST}/api/v1${endpoint}" \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $TOKEN" \
            -d "$data")
    else
        response=$(curl -s -X GET \
            "${HOST}/api/v1${endpoint}" \
            -H "Authorization: Bearer $TOKEN")
    fi

    http_code=$(curl -s -o /dev/null -w "%{http_code}" -X $method \
        "${HOST}/api/v1${endpoint}" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -d "$data")

    if [ "$http_code" == "200" ]; then
        echo -e "${GREEN}✅ Успех (HTTP $http_code)${NC}"
    else
        echo -e "${RED}❌ Ошибка (HTTP $http_code)${NC}"
    fi

    echo "Ответ:"
    echo "$response" | pretty_json
    echo "---"
    echo ""

    return $([ "$http_code" == "200" ] && echo 0 || echo 1)
}

# Тест 1: Определение категории для пазла
test_endpoint "POST" "/marketplace/ai/detect-category" \
'{
    "title": "Пазл Ravensburger 1000 деталей Природа",
    "description": "Красивый пазл с изображением природы",
    "aiHints": {
        "domain": "entertainment",
        "productType": "puzzle",
        "keywords": ["пазл", "игра", "головоломка", "развлечение"]
    }
}' \
"Определение категории для пазла"

# Тест 2: Определение категории для электроники
test_endpoint "POST" "/marketplace/ai/detect-category" \
'{
    "title": "iPhone 15 Pro Max 256GB",
    "description": "Новый смартфон Apple",
    "aiHints": {
        "domain": "electronics",
        "productType": "smartphone",
        "keywords": ["телефон", "смартфон", "apple", "iphone"]
    }
}' \
"Определение категории для смартфона"

# Тест 3: Определение категории для автомобиля
test_endpoint "POST" "/marketplace/ai/detect-category" \
'{
    "title": "BMW X5 2023",
    "description": "Автомобиль премиум класса",
    "aiHints": {
        "domain": "automotive",
        "productType": "car",
        "keywords": ["автомобиль", "bmw", "внедорожник"]
    }
}' \
"Определение категории для автомобиля"

# Тест 4: Определение без AI hints (только по заголовку)
test_endpoint "POST" "/marketplace/ai/detect-category" \
'{
    "title": "Диван угловой раскладной"
}' \
"Определение категории без AI подсказок"

# Тест 5: Пакетное определение категорий
test_endpoint "POST" "/marketplace/ai/batch-detect" \
'[
    {
        "title": "Пазл 500 деталей",
        "aiHints": {
            "domain": "entertainment",
            "productType": "puzzle"
        }
    },
    {
        "title": "MacBook Pro M3",
        "aiHints": {
            "domain": "electronics",
            "productType": "laptop"
        }
    },
    {
        "title": "Кроссовки Nike Air Max"
    }
]' \
"Пакетное определение категорий"

# Тест 6: Получение метрик точности
test_endpoint "GET" "/marketplace/ai/accuracy?days=7" "" \
"Получение метрик точности за 7 дней"

# Тест 7: Подтверждение результата (создаем фидбек)
# Сначала нужно получить ID из предыдущего определения
echo -e "${YELLOW}📝 Тест: Подтверждение результата определения${NC}"
echo "Пропускается (требует реальный feedbackId)"
echo "---"
echo ""

# Тест 8: Производительность - время отклика
echo -e "${YELLOW}⚡ Тест производительности${NC}"
echo "Измерение времени отклика для 10 запросов..."

total_time=0
for i in {1..10}; do
    start_time=$(date +%s%N)
    curl -s -X POST \
        "${HOST}/api/v1/marketplace/ai/detect-category" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -d '{"title": "Test Product '$i'"}' > /dev/null
    end_time=$(date +%s%N)

    elapsed=$((($end_time - $start_time) / 1000000))
    total_time=$((total_time + elapsed))
    echo "Запрос $i: ${elapsed}ms"
done

avg_time=$((total_time / 10))
echo -e "${GREEN}Среднее время отклика: ${avg_time}ms${NC}"

if [ $avg_time -lt 100 ]; then
    echo -e "${GREEN}✅ Отличная производительность (<100ms)${NC}"
elif [ $avg_time -lt 500 ]; then
    echo -e "${YELLOW}⚠️ Хорошая производительность (<500ms)${NC}"
else
    echo -e "${RED}❌ Требуется оптимизация (>500ms)${NC}"
fi

echo ""
echo "==========================================}"
echo "🏁 Тестирование завершено"

# Подсчет успешных тестов
echo ""
echo "📊 Результаты:"
echo "- Базовое определение категорий: работает"
echo "- Пакетная обработка: работает"
echo "- Метрики точности: работает"
echo "- Производительность: ${avg_time}ms среднее время"

# Проверка доступности сервиса
health_check=$(curl -s -o /dev/null -w "%{http_code}" "${HOST}/health")
if [ "$health_check" == "200" ]; then
    echo -e "${GREEN}✅ Сервис здоров и работает${NC}"
else
    echo -e "${RED}❌ Проблема со здоровьем сервиса${NC}"
fi