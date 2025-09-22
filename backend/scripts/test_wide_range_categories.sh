#!/bin/bash

# Тест широкого спектра товаров для проверки покрытия AI маппингов
# Цель: выявить товары, которые определяются неправильно

echo "=== AI Category Detection - Wide Range Test ==="
echo "Time: $(date)"
echo ""

API_URL="http://localhost:3000/api/v1/marketplace/ai/detect-category"
FAILED_TESTS=()
TOTAL_TESTS=0
PASSED_TESTS=0

# Функция для тестирования одного кейса
test_category() {
    local title="$1"
    local expected_domain="$2"
    local expected_type="$3"
    local min_confidence="$4"

    TOTAL_TESTS=$((TOTAL_TESTS + 1))

    echo "Testing: $title"

    # Формируем JSON запрос
    json_payload=$(cat <<EOF
{
  "title": "$title",
  "description": "Качественный товар в хорошем состоянии",
  "keywords": [],
  "language": "ru",
  "categoryHints": {
    "domain": "$expected_domain",
    "productType": "$expected_type",
    "keywords": []
  }
}
EOF
    )

    # Отправляем запрос
    response=$(curl -s -X POST "$API_URL" \
        -H "Content-Type: application/json" \
        -d "$json_payload")

    # Парсим ответ
    category_id=$(echo "$response" | jq -r '.data.categoryId')
    category_name=$(echo "$response" | jq -r '.data.categoryName')
    confidence=$(echo "$response" | jq -r '.data.confidenceScore')
    processing_time=$(echo "$response" | jq -r '.data.processingTimeMs')

    # Проверяем результат
    if [[ "$category_id" == "1001" || "$confidence" == "0.1" ]]; then
        echo "  Result: Category=$category_name (ID=$category_id)"
        echo "  Confidence: $confidence (expected >$min_confidence)"
        echo "  Processing time: ${processing_time}ms"
        echo "  ❌ FAILED - попал в General категорию!"
        FAILED_TESTS+=("$title")
    else
        echo "  Result: Category=$category_name (ID=$category_id)"
        echo "  Confidence: $confidence (expected >$min_confidence)"
        echo "  Processing time: ${processing_time}ms"

        # Проверяем confidence
        if (( $(echo "$confidence >= $min_confidence" | bc -l) )); then
            echo "  ✅ PASSED"
            PASSED_TESTS=$((PASSED_TESTS + 1))
        else
            echo "  ❌ FAILED - низкий confidence!"
            FAILED_TESTS+=("$title")
        fi
    fi
    echo ""
}

echo "=== Строительные инструменты ==="
test_category "Болгарка Makita 125мм 800W" "construction tools" "angle grinder" "0.7"
test_category "Дрель ударная Bosch GSB 600W" "construction tools" "drill" "0.7"
test_category "Циркулярная пила Hitachi 190мм" "construction tools" "circular saw" "0.7"
test_category "Перфоратор Hilti TE 6-A36" "construction tools" "hammer drill" "0.7"
test_category "Шуруповерт аккумуляторный Makita" "construction tools" "screwdriver" "0.7"
test_category "Лобзик электрический Bosch PST" "construction tools" "jigsaw" "0.7"
test_category "Рубанок электрический DeWalt" "construction tools" "planer" "0.7"
test_category "Фрезер ручной Festool OF" "construction tools" "router" "0.7"
test_category "Углошлифовальная машина" "construction tools" "angle grinder" "0.7"
test_category "Торцовочная пила Metabo" "construction tools" "miter saw" "0.7"

echo "=== Автомобили и запчасти ==="
test_category "BMW X5 2020 дизель автомат" "automotive" "car" "0.7"
test_category "Масло моторное Castrol 5W30 4л" "automotive" "car parts" "0.6"
test_category "Шины зимние Michelin R16 205/55" "automotive" "tires" "0.6"
test_category "Аккумулятор Bosch S4 60Ah" "automotive" "battery" "0.6"
test_category "Фары передние BMW E39" "automotive" "headlights" "0.6"
test_category "Тормозные колодки Brembo" "automotive" "brake parts" "0.6"
test_category "Диски литые R17 BMW оригинал" "automotive" "wheels" "0.6"
test_category "Глушитель Volkswagen Golf" "automotive" "exhaust" "0.6"

echo "=== Электроника ==="
test_category "Телевизор Samsung 55 QLED 4K" "electronics" "tv" "0.7"
test_category "Наушники Sony WH-1000XM4" "electronics" "headphones" "0.7"
test_category "Роутер TP-Link AX6000 WiFi 6" "electronics" "router" "0.7"
test_category "Смартфон iPhone 14 Pro 128GB" "electronics" "smartphone" "0.8"
test_category "Ноутбук MacBook Pro M1 13" "electronics" "laptop" "0.8"
test_category "Планшет iPad Air 64GB" "electronics" "tablet" "0.7"
test_category "Камера Canon EOS R6" "electronics" "camera" "0.7"
test_category "Принтер HP LaserJet P1102" "electronics" "printer" "0.7"

echo "=== Дом и сад ==="
test_category "Газонокосилка электрическая" "home-garden" "lawn mower" "0.6"
test_category "Триммер бензиновый Stihl" "home-garden" "trimmer" "0.6"
test_category "Мотоблок Нева МБ-2" "home-garden" "cultivator" "0.6"
test_category "Теплица поликарбонат 3х6м" "home-garden" "greenhouse" "0.6"
test_category "Насос водяной центробежный" "home-garden" "water pump" "0.6"
test_category "Бензопила Husqvarna 445" "home-garden" "chainsaw" "0.7"

echo "=== Спорт и отдых ==="
test_category "Велосипед горный Trek 29" "sports-recreation" "bicycle" "0.6"
test_category "Лыжи горные Salomon 170см" "sports-recreation" "skis" "0.6"
test_category "Палатка туристическая 4 места" "sports-recreation" "tent" "0.6"
test_category "Удочка спиннинговая Shimano" "sports-recreation" "fishing rod" "0.6"
test_category "Гантели разборные 2х20кг" "sports-recreation" "weights" "0.6"

echo "=== Мода и одежда ==="
test_category "Куртка зимняя мужская L" "fashion" "jacket" "0.5"
test_category "Кроссовки Nike Air Max 42" "fashion" "shoes" "0.6"
test_category "Джинсы Levis 501 W32 L34" "fashion" "jeans" "0.5"
test_category "Платье вечернее размер M" "fashion" "dress" "0.5"
test_category "Сумка женская кожаная" "fashion" "bag" "0.5"

echo ""
echo "=== ИТОГОВЫЕ РЕЗУЛЬТАТЫ ==="
echo "Всего тестов: $TOTAL_TESTS"
echo "Пройдено: $PASSED_TESTS"
echo "Провалено: $((TOTAL_TESTS - PASSED_TESTS))"

if [ ${#FAILED_TESTS[@]} -eq 0 ]; then
    echo "✅ ВСЕ ТЕСТЫ ПРОШЛИ УСПЕШНО!"
else
    echo "❌ Провалившиеся тесты:"
    for failed_test in "${FAILED_TESTS[@]}"; do
        echo "  - $failed_test"
    done
fi

accuracy=$(echo "scale=2; $PASSED_TESTS * 100 / $TOTAL_TESTS" | bc)
echo ""
echo "Точность: ${accuracy}%"

if (( $(echo "$accuracy >= 80" | bc -l) )); then
    echo "✅ СИСТЕМА РАБОТАЕТ ХОРОШО (точность ≥80%)"
else
    echo "❌ СИСТЕМА ТРЕБУЕТ ДОРАБОТКИ (точность <80%)"
fi

echo ""
echo "Результаты сохранены в: /tmp/wide_range_test_results.log"

# Сохраняем подробные результаты
{
    echo "=== AI Category Detection Wide Range Test ==="
    echo "Time: $(date)"
    echo "Total: $TOTAL_TESTS, Passed: $PASSED_TESTS, Failed: $((TOTAL_TESTS - PASSED_TESTS))"
    echo "Accuracy: ${accuracy}%"
    echo ""
    echo "Failed tests:"
    for failed_test in "${FAILED_TESTS[@]}"; do
        echo "- $failed_test"
    done
} > /tmp/wide_range_test_results.log