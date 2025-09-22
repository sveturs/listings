#!/bin/bash

# Тестовый скрипт для проверки точности AI определения категорий
# Выполняет тестовые запросы и проверяет качество работы

API_URL="http://localhost:3000/api/v1/marketplace/ai/detect-category"
RESULTS_FILE="/tmp/ai_test_results.log"

echo "=== AI Category Detection Accuracy Test ===" | tee $RESULTS_FILE
echo "Time: $(date)" | tee -a $RESULTS_FILE
echo "" | tee -a $RESULTS_FILE

# Счетчики
TOTAL=0
CORRECT=0

# Функция для тестирования
test_detection() {
    local title="$1"
    local domain="$2"
    local product_type="$3"
    local expected_confidence="$4"

    echo "Testing: $title" | tee -a $RESULTS_FILE

    RESPONSE=$(curl -s -X POST $API_URL \
        -H "Content-Type: application/json" \
        -d "{
            \"title\": \"$title\",
            \"aiHints\": {
                \"domain\": \"$domain\",
                \"productType\": \"$product_type\",
                \"keywords\": []
            }
        }")

    # Извлекаем данные из ответа
    CATEGORY_ID=$(echo "$RESPONSE" | jq -r '.data.categoryId')
    CATEGORY_NAME=$(echo "$RESPONSE" | jq -r '.data.categoryName')
    CONFIDENCE=$(echo "$RESPONSE" | jq -r '.data.confidenceScore')
    ALGORITHM=$(echo "$RESPONSE" | jq -r '.data.algorithm')
    TIME_MS=$(echo "$RESPONSE" | jq -r '.data.processingTimeMs')

    echo "  Result: Category=$CATEGORY_NAME (ID=$CATEGORY_ID)" | tee -a $RESULTS_FILE
    echo "  Confidence: $CONFIDENCE (expected >$expected_confidence)" | tee -a $RESULTS_FILE
    echo "  Processing time: ${TIME_MS}ms" | tee -a $RESULTS_FILE

    TOTAL=$((TOTAL + 1))

    # Проверяем, что confidence достаточно высокий
    if (( $(echo "$CONFIDENCE > $expected_confidence" | bc -l) )); then
        echo "  ✅ PASSED" | tee -a $RESULTS_FILE
        CORRECT=$((CORRECT + 1))
    else
        echo "  ❌ FAILED (confidence too low)" | tee -a $RESULTS_FILE
    fi

    echo "" | tee -a $RESULTS_FILE
}

# Тестовые кейсы
echo "=== Test Cases ===" | tee -a $RESULTS_FILE
echo "" | tee -a $RESULTS_FILE

# Пазлы
test_detection "Пазл с изображением природы 1000 деталей" "entertainment" "puzzle" "0.7"

# Строительные материалы
test_detection "Мешок с песком строительный 50кг" "construction" "sand" "0.7"
test_detection "Цемент М500 мешок 25кг" "construction" "cement" "0.6"

# Природные материалы
test_detection "Желудь дубовый для поделок" "nature" "acorn" "0.6"
test_detection "Семена подсолнуха отборные" "nature" "seeds" "0.5"

# Авиация
test_detection "Модель самолета Boeing 747" "aviation" "aircraft-model" "0.7"

# Антиквариат
test_detection "Старинная монета времен СССР" "antiques" "coin" "0.6"
test_detection "Коллекция почтовых марок" "antiques" "stamp" "0.6"

# Военные товары
test_detection "Военная форма камуфляж размер L" "military" "uniform" "0.6"
test_detection "Медаль За отвагу оригинал" "military" "medal" "0.6"

# Рукоделие
test_detection "Набор для вышивания крестиком" "crafts" "craft-kit" "0.5"

# Электроника
test_detection "iPhone 13 Pro Max 256GB" "electronics" "smartphone" "0.7"
test_detection "MacBook Pro M1 13 inch" "electronics" "laptop" "0.7"

# Автомобили
test_detection "Volkswagen Golf 2015 дизель" "automotive" "car" "0.7"

# Недвижимость
test_detection "Квартира 2 комнаты центр города" "real-estate" "apartment" "0.5"

# Вычисляем точность
echo "=== RESULTS ===" | tee -a $RESULTS_FILE
ACCURACY=$(echo "scale=2; $CORRECT * 100 / $TOTAL" | bc)
echo "Total tests: $TOTAL" | tee -a $RESULTS_FILE
echo "Passed: $CORRECT" | tee -a $RESULTS_FILE
echo "Failed: $((TOTAL - CORRECT))" | tee -a $RESULTS_FILE
echo "Accuracy: ${ACCURACY}%" | tee -a $RESULTS_FILE
echo "" | tee -a $RESULTS_FILE

# Тестируем производительность
echo "=== Performance Test ===" | tee -a $RESULTS_FILE
echo "Testing 10 rapid requests..." | tee -a $RESULTS_FILE
START_TIME=$(date +%s%N)

for i in {1..10}; do
    curl -s -X POST $API_URL \
        -H "Content-Type: application/json" \
        -d "{
            \"title\": \"Test product $i\",
            \"aiHints\": {
                \"domain\": \"electronics\",
                \"productType\": \"laptop\",
                \"keywords\": [\"test\"]
            }
        }" > /dev/null
done

END_TIME=$(date +%s%N)
ELAPSED=$((($END_TIME - $START_TIME) / 1000000))
AVG_TIME=$(($ELAPSED / 10))

echo "Total time for 10 requests: ${ELAPSED}ms" | tee -a $RESULTS_FILE
echo "Average time per request: ${AVG_TIME}ms" | tee -a $RESULTS_FILE
echo "" | tee -a $RESULTS_FILE

# Тестируем метрики
echo "=== Fetching Metrics ===" | tee -a $RESULTS_FILE
METRICS=$(curl -s http://localhost:3000/api/v1/marketplace/ai/metrics?days=1)
echo "$METRICS" | jq '.' | tee -a $RESULTS_FILE

echo "" | tee -a $RESULTS_FILE
echo "Test completed. Full results saved to: $RESULTS_FILE" | tee -a $RESULTS_FILE

# Вывод итогового статуса
if (( $(echo "$ACCURACY >= 70" | bc -l) )); then
    echo "" | tee -a $RESULTS_FILE
    echo "✅ SYSTEM IS WORKING WELL (Accuracy: ${ACCURACY}%)" | tee -a $RESULTS_FILE
else
    echo "" | tee -a $RESULTS_FILE
    echo "⚠️ SYSTEM NEEDS IMPROVEMENT (Accuracy: ${ACCURACY}%)" | tee -a $RESULTS_FILE
fi