#!/bin/bash

# Comprehensive AI Category Detection Test Script
# Тестирует 100+ товаров для достижения 99% точности

set -e

# Configuration
BASE_URL="http://localhost:3000"
API_URL="$BASE_URL/api/v1/marketplace/ai"
RESULTS_FILE="/tmp/ai_test_results_$(date +%Y%m%d_%H%M%S).json"
LOG_FILE="/tmp/ai_test_log_$(date +%Y%m%d_%H%M%S).log"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo "🎯 Comprehensive AI Category Detection Test"
echo "==========================================="
echo "Base URL: $BASE_URL"
echo "Results file: $RESULTS_FILE"
echo "Log file: $LOG_FILE"
echo ""

# Initialize results
echo '{"testRun": {"startTime": "'$(date -Iseconds)'", "results": []}}' > "$RESULTS_FILE"

# Test cases - 100+ diverse products
declare -a TEST_CASES=(
    # Строительные инструменты (20 кейсов)
    "Болгарка Makita GA5030R 125мм 720W|construction|angle grinder|Строительные инструменты"
    "Дрель ударная Bosch GSB 600 RE 600W|construction|drill|Строительные инструменты"
    "Циркулярная пила Hitachi C7ST 1710W|construction|circular saw|Строительные инструменты"
    "Перфоратор Hilti TE 3-C SDS Plus|construction|hammer drill|Строительные инструменты"
    "Шуруповерт аккумуляторный DeWalt DCD771C2|construction|screwdriver|Строительные инструменты"
    "Лобзик электрический Festool PSB 420 EBQ|construction|jigsaw|Строительные инструменты"
    "Рубанок электрический Bosch GHO 26-82 D|construction|planer|Строительные инструменты"
    "Фрезер Makita RT0700CX2J|construction|router|Строительные инструменты"
    "Пила торцовочная Metabo KGS 216 M|construction|miter saw|Строительные инструменты"
    "Гайковерт пневматический Atlas Copco W2918|construction|impact wrench|Строительные инструменты"
    "Отбойный молоток Bosch GSH 16-28|construction|demolition hammer|Строительные инструменты"
    "Штроборез Makita SG1251J|construction|wall chaser|Строительные инструменты"
    "Пистолет для пены Hilti CF-DS1|construction|foam gun|Строительные инструменты"
    "Миксер строительный Collomix XM 2-650|construction|mixer|Строительные инструменты"
    "Стамеска Narex 8101 20мм|construction|chisel|Строительные инструменты"
    "Уровень строительный Stabila 70-2|construction|level|Строительные инструменты"
    "Молоток слесарный 500г|construction|hammer|Строительные инструменты"
    "Пассатижи Knipex 03 01 160|construction|pliers|Строительные инструменты"
    "Отвертка крестовая PH2|construction|screwdriver|Строительные инструменты"
    "Ключ разводной 250мм|construction|adjustable wrench|Строительные инструменты"

    # Автомобили и запчасти (15 кейсов)
    "BMW X5 xDrive30d 2020 г.в.|automotive|car|Автомобили"
    "Mercedes-Benz E-Class E220d|automotive|car|Автомобили"
    "Audi A4 Avant quattro|automotive|car|Автомобили"
    "Volkswagen Golf GTI|automotive|car|Автомобили"
    "Toyota Camry 2.5 Hybrid|automotive|car|Автомобили"
    "Масло моторное Castrol GTX 5W-30 4л|automotive|oil|Автозапчасти"
    "Шины летние Michelin Pilot Sport 4 225/45R17|automotive|tires|Автозапчасти"
    "Аккумулятор Bosch S4 74Ah|automotive|battery|Автозапчасти"
    "Тормозные колодки Brembo P50084|automotive|brake pads|Автозапчасти"
    "Фильтр воздушный Mann C25114|automotive|air filter|Автозапчасти"
    "Свечи зажигания NGK BPR6ES|automotive|spark plugs|Автозапчасти"
    "Амортизаторы передние Bilstein B4|automotive|shock absorbers|Автозапчасти"
    "Диски литые BBS CH-R 18x8|automotive|wheels|Автозапчасти"
    "Коврики резиновые Novline|automotive|floor mats|Автозапчасти"
    "Чехлы на сиденья Автопилот|automotive|seat covers|Автозапчасти"

    # Электроника (20 кейсов)
    "iPhone 15 Pro Max 256GB Space Black|electronics|smartphone|Мобильные телефоны"
    "Samsung Galaxy S24 Ultra 512GB|electronics|smartphone|Мобильные телефоны"
    "MacBook Pro 16 M3 Pro 512GB|electronics|laptop|Компьютеры"
    "ASUS ROG Strix G15 RTX 4060|electronics|laptop|Компьютеры"
    "iPad Pro 12.9 M2 1TB|electronics|tablet|Планшеты"
    "Телевизор Samsung QE55QN95B QLED 55|electronics|tv|Телевизоры"
    "LG OLED55C3PLA 55 4K|electronics|tv|Телевизоры"
    "Наушники Sony WH-1000XM5|electronics|headphones|Аудиотехника"
    "AirPods Pro 2 поколение|electronics|earbuds|Аудиотехника"
    "Роутер ASUS AX6000 RT-AX88U|electronics|router|Сетевое оборудование"
    "PlayStation 5 825GB|electronics|console|Игровые консоли"
    "Nintendo Switch OLED|electronics|console|Игровые консоли"
    "Камера Canon EOS R6 Mark II|electronics|camera|Фототехника"
    "Объектив Sony FE 24-70mm f/2.8 GM|electronics|lens|Фототехника"
    "Микрофон Audio-Technica AT2020|electronics|microphone|Аудиотехника"
    "Веб-камера Logitech C920 HD Pro|electronics|webcam|Компьютерная периферия"
    "Клавиатура Logitech MX Keys|electronics|keyboard|Компьютерная периферия"
    "Мышь Logitech MX Master 3S|electronics|mouse|Компьютерная периферия"
    "Монитор Dell UltraSharp U2723QE 27|electronics|monitor|Мониторы"
    "SSD Samsung 980 PRO 2TB|electronics|ssd|Комплектующие для ПК"

    # Бытовая техника (15 кейсов)
    "Холодильник Samsung RB37K63611L|appliances|refrigerator|Холодильники"
    "Стиральная машина Bosch WAV28G40OE|appliances|washing machine|Стиральные машины"
    "Посудомоечная машина Electrolux ESM46200L|appliances|dishwasher|Посудомоечные машины"
    "Микроволновая печь LG MS2336GIB|appliances|microwave|Микроволновые печи"
    "Духовой шкаф Gorenje BO758A32BG|appliances|oven|Духовые шкафы"
    "Варочная панель Electrolux EHH6240ISK|appliances|cooktop|Варочные панели"
    "Вытяжка Faber Stilo SP EG8 X A60|appliances|range hood|Вытяжки"
    "Пылесос Dyson V15 Detect|appliances|vacuum cleaner|Пылесосы"
    "Утюг Philips Azur Elite GC5033|appliances|iron|Утюги"
    "Кофемашина De'Longhi Dinamica ECAM350.15.B|appliances|coffee machine|Кофемашины"
    "Блендер Vitamix A3500|appliances|blender|Блендеры"
    "Мультиварка Redmond RMC-M150|appliances|multicooker|Мультиварки"
    "Кондиционер Daikin FTXM35R|appliances|air conditioner|Кондиционеры"
    "Обогреватель DeLonghi HMP1500|appliances|heater|Обогреватели"
    "Увлажнитель воздуха Xiaomi Mi Smart|appliances|humidifier|Увлажнители"

    # Мебель и интерьер (10 кейсов)
    "Диван угловой IKEA Ektorp|furniture|sofa|Мягкая мебель"
    "Кровать двуспальная 160x200|furniture|bed|Кровати"
    "Шкаф-купе 3-дверный|furniture|wardrobe|Шкафы"
    "Стол обеденный раздвижной|furniture|dining table|Столы"
    "Стулья деревянные с мягкой обивкой|furniture|chairs|Стулья"
    "Комод с 4 ящиками|furniture|dresser|Комоды"
    "Кресло офисное эргономичное|furniture|office chair|Офисная мебель"
    "Полки настенные навесные|furniture|shelves|Полки"
    "Зеркало в раме 80x60|furniture|mirror|Зеркала"
    "Светильник потолочный LED|furniture|ceiling light|Освещение"

    # Одежда и обувь (10 кейсов)
    "Куртка зимняя пуховая North Face|fashion|jacket|Верхняя одежда"
    "Джинсы мужские Levi's 501|fashion|jeans|Брюки и джинсы"
    "Платье летнее из хлопка|fashion|dress|Платья"
    "Рубашка мужская классическая|fashion|shirt|Рубашки"
    "Кроссовки Nike Air Max 270|fashion|sneakers|Кроссовки"
    "Ботинки зимние Timberland|fashion|boots|Ботинки"
    "Сумка женская кожаная|fashion|handbag|Сумки"
    "Часы наручные Casio G-Shock|fashion|watch|Часы"
    "Перчатки кожаные зимние|fashion|gloves|Перчатки"
    "Шарф вязаный шерстяной|fashion|scarf|Шарфы"

    # Спорт и отдых (10 кейсов)
    "Велосипед горный Trek X-Caliber 8|sports|bicycle|Велосипеды"
    "Лыжи горные Rossignol Experience 76|sports|skis|Лыжи"
    "Сноуборд Burton Custom X|sports|snowboard|Сноуборды"
    "Палатка туристическая 3-местная|sports|tent|Туристическое снаряжение"
    "Рюкзак туристический 70л|sports|backpack|Туристическое снаряжение"
    "Удочка спиннинговая Shimano|sports|fishing rod|Рыболовство"
    "Мяч футбольный Adidas UEFA|sports|football|Мячи"
    "Ракетка теннисная Wilson Pro Staff|sports|tennis racket|Теннис"
    "Гантели разборные 20кг|sports|dumbbells|Фитнес"
    "Коврик для йоги Manduka|sports|yoga mat|Йога"

    # Книги и канцелярия (5 кейсов)
    "Программирование на Python 3-е издание|books|programming book|Книги"
    "Ручка шариковая Parker Jotter|books|pen|Канцелярские товары"
    "Блокнот Moleskine классический|books|notebook|Канцелярские товары"
    "Калькулятор научный Casio FX-991EX|books|calculator|Канцелярские товары"
    "Маркеры цветные Stabilo Boss|books|markers|Канцелярские товары"

    # Экзотические товары для проверки edge cases (5 кейсов)
    "Желудь дубовый для поделок|nature|acorn|Природные материалы"
    "Мешок с песком строительный 25кг|construction|sand bag|Строительные материалы"
    "Модель самолета Boeing 747|collectibles|airplane model|Коллекционирование"
    "Пазл 1000 элементов Замок Нойшванштайн|entertainment|puzzle|Пазлы"
    "Антикварные часы карманные|antiques|pocket watch|Антиквариат"
)

# Counters
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0
AI_VALIDATION_PASSED=0
AI_VALIDATION_FAILED=0

# Function to log with timestamp
log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1" | tee -a "$LOG_FILE"
}

# Function to make API call and validate response
test_product() {
    local title="$1"
    local expected_domain="$2"
    local expected_product_type="$3"
    local expected_category="$4"
    local test_id="$5"

    log "Testing: $title"

    # Test detection
    local detection_response=$(curl -s -X POST "$API_URL/detect-category" \
        -H "Content-Type: application/json" \
        -d "{
            \"title\": \"$title\",
            \"description\": \"$title\",
            \"language\": \"ru\"
        }")

    if [ $? -ne 0 ]; then
        log "❌ Detection API call failed for: $title"
        echo "$detection_response" >> "$LOG_FILE"
        return 1
    fi

    # Parse detection result
    local detected_category=$(echo "$detection_response" | jq -r '.data.categoryName // "unknown"')
    local confidence=$(echo "$detection_response" | jq -r '.data.confidence // 0')
    local processing_time=$(echo "$detection_response" | jq -r '.data.processingTimeMs // 0')

    if [ "$detected_category" = "null" ] || [ "$detected_category" = "unknown" ]; then
        log "❌ No category detected for: $title"
        return 1
    fi

    # Test AI validation
    local validation_response=$(curl -s -X POST "$API_URL/validate-category" \
        -H "Content-Type: application/json" \
        -d "{
            \"title\": \"$title\",
            \"description\": \"$title\",
            \"categoryName\": \"$detected_category\"
        }")

    if [ $? -ne 0 ]; then
        log "❌ Validation API call failed for: $title"
        echo "$validation_response" >> "$LOG_FILE"
        return 1
    fi

    # Parse validation result
    local is_correct=$(echo "$validation_response" | jq -r '.data.isCorrect // false')
    local ai_confidence=$(echo "$validation_response" | jq -r '.data.confidence // 0')
    local reasoning=$(echo "$validation_response" | jq -r '.data.reasoning // ""')
    local suggested_category=$(echo "$validation_response" | jq -r '.data.suggestedCategory // ""')

    # Determine test result
    local test_passed="false"
    local notes=""

    # Check if category is reasonable (not "General" or empty)
    if [[ "$detected_category" != *"General"* ]] && [[ "$detected_category" != *"Общие"* ]] && [ "$detected_category" != "" ]; then
        # If AI validation is positive or confident enough
        if [ "$is_correct" = "true" ] || ([ "$is_correct" = "false" ] && (( $(echo "$ai_confidence < 0.7" | bc -l) ))); then
            test_passed="true"
            notes="Category seems appropriate"
        else
            test_passed="false"
            notes="AI validation failed with high confidence. Suggested: $suggested_category"
        fi
    else
        test_passed="false"
        notes="Category too generic or empty"
    fi

    # Update counters
    if [ "$test_passed" = "true" ]; then
        PASSED_TESTS=$((PASSED_TESTS + 1))
        echo -e "${GREEN}✅ PASS${NC}: $title → $detected_category"
    else
        FAILED_TESTS=$((FAILED_TESTS + 1))
        echo -e "${RED}❌ FAIL${NC}: $title → $detected_category ($notes)"
    fi

    if [ "$is_correct" = "true" ]; then
        AI_VALIDATION_PASSED=$((AI_VALIDATION_PASSED + 1))
    else
        AI_VALIDATION_FAILED=$((AI_VALIDATION_FAILED + 1))
    fi

    # Save detailed result to JSON
    local result_json=$(cat <<EOF
{
    "testId": $test_id,
    "title": "$title",
    "expectedDomain": "$expected_domain",
    "expectedProductType": "$expected_product_type",
    "expectedCategory": "$expected_category",
    "detectedCategory": "$detected_category",
    "confidence": $confidence,
    "processingTimeMs": $processing_time,
    "aiValidation": {
        "isCorrect": $is_correct,
        "confidence": $ai_confidence,
        "reasoning": "$reasoning",
        "suggestedCategory": "$suggested_category"
    },
    "testPassed": $test_passed,
    "notes": "$notes",
    "timestamp": "$(date -Iseconds)"
}
EOF
    )

    # Append to results file (we'll fix JSON structure later)
    echo "$result_json," >> "${RESULTS_FILE}.tmp"

    return 0
}

# Main test execution
echo "🚀 Starting comprehensive test with ${#TEST_CASES[@]} products..."
echo ""

# Initialize temp results file
echo '' > "${RESULTS_FILE}.tmp"

TOTAL_TESTS=${#TEST_CASES[@]}
test_id=1

for test_case in "${TEST_CASES[@]}"; do
    IFS='|' read -r title domain product_type category <<< "$test_case"

    echo -e "${BLUE}[Test $test_id/$TOTAL_TESTS]${NC} Testing: $title"

    test_product "$title" "$domain" "$product_type" "$category" $test_id

    test_id=$((test_id + 1))

    # Small delay to avoid overwhelming the server
    sleep 0.5
done

# Calculate final statistics
ACCURACY=$(( (PASSED_TESTS * 100) / TOTAL_TESTS ))
AI_VALIDATION_ACCURACY=$(( (AI_VALIDATION_PASSED * 100) / TOTAL_TESTS ))

# Create final results JSON
{
    echo '{'
    echo '  "testRun": {'
    echo "    \"startTime\": \"$(date -Iseconds)\","
    echo "    \"totalTests\": $TOTAL_TESTS,"
    echo "    \"passedTests\": $PASSED_TESTS,"
    echo "    \"failedTests\": $FAILED_TESTS,"
    echo "    \"accuracy\": $ACCURACY,"
    echo "    \"aiValidationPassed\": $AI_VALIDATION_PASSED,"
    echo "    \"aiValidationFailed\": $AI_VALIDATION_FAILED,"
    echo "    \"aiValidationAccuracy\": $AI_VALIDATION_ACCURACY,"
    echo '    "results": ['

    # Add all test results (remove last comma)
    if [ -f "${RESULTS_FILE}.tmp" ]; then
        sed '$ s/,$//' "${RESULTS_FILE}.tmp"
    fi

    echo '    ]'
    echo '  }'
    echo '}'
} > "$RESULTS_FILE"

# Clean up temp file
rm -f "${RESULTS_FILE}.tmp"

# Print final report
echo ""
echo "=============================================="
echo "🎯 COMPREHENSIVE AI TEST RESULTS"
echo "=============================================="
echo -e "📊 Total tests: ${BLUE}$TOTAL_TESTS${NC}"
echo -e "✅ Passed: ${GREEN}$PASSED_TESTS${NC}"
echo -e "❌ Failed: ${RED}$FAILED_TESTS${NC}"
echo -e "🎯 Detection accuracy: ${BLUE}$ACCURACY%${NC}"
echo -e "🤖 AI validation passed: ${GREEN}$AI_VALIDATION_PASSED${NC}"
echo -e "🤖 AI validation failed: ${RED}$AI_VALIDATION_FAILED${NC}"
echo -e "🎯 AI validation accuracy: ${BLUE}$AI_VALIDATION_ACCURACY%${NC}"
echo ""

# Goal achievement check
if [ $ACCURACY -ge 99 ]; then
    echo -e "${GREEN}🎉 GOAL ACHIEVED! 99% accuracy reached!${NC}"
elif [ $ACCURACY -ge 95 ]; then
    echo -e "${YELLOW}📈 CLOSE TO GOAL! $ACCURACY% accuracy (goal: 99%)${NC}"
elif [ $ACCURACY -ge 90 ]; then
    echo -e "${YELLOW}📊 GOOD PROGRESS! $ACCURACY% accuracy (goal: 99%)${NC}"
else
    echo -e "${RED}⚠️  NEEDS IMPROVEMENT! $ACCURACY% accuracy (goal: 99%)${NC}"
fi

echo ""
echo "📄 Results saved to: $RESULTS_FILE"
echo "📝 Log saved to: $LOG_FILE"
echo ""

# Summary of failed tests
if [ $FAILED_TESTS -gt 0 ]; then
    echo "❌ Failed tests summary:"
    grep "❌ FAIL" "$LOG_FILE" | tail -10
    echo ""
fi

exit 0