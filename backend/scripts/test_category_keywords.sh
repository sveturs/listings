#!/bin/bash

# Quick Category Keywords Test Script
# Быстрое тестирование ключевых слов для конкретных категорий

set -e

# Configuration
BASE_URL="http://localhost:3000"
API_URL="$BASE_URL/api/v1/marketplace/ai"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo "🔍 Category Keywords Test Tool"
echo "=============================="

# Function to test a product and show detailed results
test_product_detailed() {
    local title="$1"
    local expected_category="$2"

    echo -e "\n${BLUE}Testing:${NC} $title"
    if [ ! -z "$expected_category" ]; then
        echo -e "${BLUE}Expected:${NC} $expected_category"
    else
        echo -e "${BLUE}Mode:${NC} Auto-detection (no expected category)"
    fi

    # Test detection
    local response=$(curl -s -X POST "$API_URL/detect-category" \
        -H "Content-Type: application/json" \
        -d "{\"title\": \"$title\", \"description\": \"$title\", \"language\": \"ru\"}")

    if [ $? -ne 0 ]; then
        echo -e "${RED}❌ API Error${NC}"
        return 1
    fi

    local detected=$(echo "$response" | jq -r '.data.categoryName // "unknown"')
    local confidence=$(echo "$response" | jq -r '.data.confidenceScore // 0')
    local time_ms=$(echo "$response" | jq -r '.data.processingTimeMs // 0')
    local method=$(echo "$response" | jq -r '.data.algorithm // "unknown"')

    echo -e "${GREEN}Detected:${NC} $detected"
    echo -e "${GREEN}Confidence:${NC} $confidence"
    echo -e "${GREEN}Method:${NC} $method"
    echo -e "${GREEN}Time:${NC} ${time_ms}ms"

    # Test AI validation
    local validation=$(curl -s -X POST "$API_URL/validate-category" \
        -H "Content-Type: application/json" \
        -d "{\"title\": \"$title\", \"categoryName\": \"$detected\"}")

    local is_correct=$(echo "$validation" | jq -r '.data.isCorrect // false')
    local ai_confidence=$(echo "$validation" | jq -r '.data.confidence // 0')
    local reasoning=$(echo "$validation" | jq -r '.data.reasoning // ""')

    echo -e "${GREEN}AI Validation:${NC} $is_correct (confidence: $ai_confidence)"
    echo -e "${GREEN}Reasoning:${NC} $reasoning"

    # Show keywords that matched (if available)
    local keywords=$(echo "$response" | jq -r '.data.keywords[]? // empty' 2>/dev/null | tr '\n' ', ' | sed 's/,$//')
    if [ ! -z "$keywords" ]; then
        echo -e "${GREEN}Matched Keywords:${NC} $keywords"
    fi

    # Final assessment
    local overall_assessment=""
    if [ "$is_correct" = "true" ]; then
        echo -e "${GREEN}✅ AI VALIDATION PASSED${NC}"
        overall_assessment="✅ EXCELLENT"
    else
        echo -e "${RED}❌ AI VALIDATION FAILED${NC}"
        if (( $(echo "$ai_confidence > 0.7" | bc -l) )); then
            overall_assessment="⚠️  QUESTIONABLE (high AI confidence disagreement)"
        else
            overall_assessment="❓ UNCERTAIN (low AI confidence)"
        fi
    fi

    # Compare with expected category if provided
    if [ ! -z "$expected_category" ]; then
        echo -e "\n${YELLOW}Expected vs Detected Comparison:${NC}"
        if [[ "$detected" == *"$expected_category"* ]] || [[ "$expected_category" == *"$detected"* ]]; then
            echo -e "${GREEN}✅ MATCHES EXPECTATION${NC}"
            overall_assessment="✅ PERFECT"
        else
            echo -e "${RED}❌ DOESN'T MATCH EXPECTATION${NC}"
            if [ "$is_correct" = "true" ]; then
                overall_assessment="⚠️  AI SAYS CORRECT BUT DIFFERS FROM EXPECTATION"
            else
                overall_assessment="❌ FAILED (both AI and expectation disagree)"
            fi
        fi
    fi

    echo -e "\n${YELLOW}🎯 OVERALL ASSESSMENT: ${overall_assessment}${NC}"

    # Add smart suggestions for common products
    suggest_common_category "$title" "$detected"
}

# Function to suggest likely categories for common products
suggest_common_category() {
    local title="$1"
    local detected="$2"
    local title_lower=$(echo "$title" | tr '[:upper:]' '[:lower:]')

    echo -e "\n${BLUE}💡 Smart Analysis:${NC}"

    # Common product patterns
    case "$title_lower" in
        *"тыква"*|*"огурец"*|*"помидор"*|*"картофель"*|*"морковь"*)
            echo -e "${YELLOW}🥕 Likely category:${NC} Продукты питания / Овощи и фрукты"
            ;;
        *"молоко"*|*"хлеб"*|*"масло"*|*"сыр"*|*"мясо"*)
            echo -e "${YELLOW}🥛 Likely category:${NC} Продукты питания / Молочные/Мясные изделия"
            ;;
        *"iphone"*|*"samsung"*|*"xiaomi"*|*"телефон"*|*"смартфон"*)
            echo -e "${YELLOW}📱 Likely category:${NC} Электроника / Мобильные телефоны"
            ;;
        *"macbook"*|*"ноутбук"*|*"компьютер"*|*"laptop"*)
            echo -e "${YELLOW}💻 Likely category:${NC} Электроника / Компьютеры"
            ;;
        *"болгарка"*|*"дрель"*|*"молоток"*|*"отвертка"*)
            echo -e "${YELLOW}🔧 Likely category:${NC} Строительные инструменты"
            ;;
        *"bmw"*|*"mercedes"*|*"audi"*|*"автомобиль"*|*"машина"*)
            echo -e "${YELLOW}🚗 Likely category:${NC} Автомобили"
            ;;
        *"диван"*|*"кровать"*|*"стол"*|*"шкаф"*|*"кресло"*)
            echo -e "${YELLOW}🪑 Likely category:${NC} Мебель и интерьер"
            ;;
        *"куртка"*|*"джинсы"*|*"рубашка"*|*"платье"*)
            echo -e "${YELLOW}👕 Likely category:${NC} Одежда и обувь"
            ;;
        *"книга"*|*"роман"*|*"учебник"*)
            echo -e "${YELLOW}📚 Likely category:${NC} Книги и канцелярия"
            ;;
        *"велосипед"*|*"лыжи"*|*"мяч"*|*"гантели"*)
            echo -e "${YELLOW}⚽ Likely category:${NC} Спорт и отдых"
            ;;
        *)
            echo -e "${YELLOW}🤔 Product type:${NC} Analyzing '$title_lower'..."

            # Check if it's likely food
            if [[ "$title_lower" == *"ягод"* || "$title_lower" == *"фрукт"* || "$title_lower" == *"овощ"* ]]; then
                echo -e "${YELLOW}🍎 Suggestion:${NC} Likely food/agricultural product"
            # Check if it's likely electronic
            elif [[ "$title_lower" == *"телевизор"* || "$title_lower" == *"наушник"* || "$title_lower" == *"планшет"* ]]; then
                echo -e "${YELLOW}📱 Suggestion:${NC} Likely electronics"
            # Check if it's likely tool
            elif [[ "$title_lower" == *"инструмент"* || "$title_lower" == *"пила"* || "$title_lower" == *"ключ"* ]]; then
                echo -e "${YELLOW}🔧 Suggestion:${NC} Likely construction tool"
            else
                echo -e "${YELLOW}❓ Suggestion:${NC} Unusual product - AI should handle this"
            fi
            ;;
    esac

    # Compare with what was actually detected
    if [[ "$detected" == *"General"* ]] || [[ "$detected" == *"Общие"* ]] || [ "$detected" = "unknown" ]; then
        echo -e "${RED}⚠️  Warning:${NC} Detected generic category - may need keyword expansion"
    else
        echo -e "${GREEN}✨ Good:${NC} Specific category detected: '$detected'"
    fi
}

# Function to generate keywords for a category
generate_keywords() {
    local category_id="$1"
    local category_name="$2"

    echo -e "\n${BLUE}Generating keywords for:${NC} $category_name (ID: $category_id)"

    local response=$(curl -s -X POST "$API_URL/generate-keywords" \
        -H "Content-Type: application/json" \
        -d "{\"categoryId\": $category_id, \"categoryName\": \"$category_name\", \"minKeywords\": 50}")

    if [ $? -ne 0 ]; then
        echo -e "${RED}❌ Keyword generation failed${NC}"
        return 1
    fi

    local generated_count=$(echo "$response" | jq -r '.data.generatedCount // 0')
    local processing_time=$(echo "$response" | jq -r '.data.processingTimeMs // 0')

    echo -e "${GREEN}Generated:${NC} $generated_count keywords in ${processing_time}ms"

    # Show some examples
    local keywords=$(echo "$response" | jq -r '.data.keywords[0:5][] | .keyword' 2>/dev/null)
    if [ ! -z "$keywords" ]; then
        echo -e "${GREEN}Examples:${NC} $(echo "$keywords" | tr '\n' ', ' | sed 's/,$//')"
    fi
}

# Function to check keyword statistics
check_keyword_stats() {
    local category_id="$1"

    echo -e "\n${BLUE}Checking keyword statistics${NC}"

    local response=$(curl -s "$API_URL/keyword-stats?categoryId=$category_id")

    if [ $? -ne 0 ]; then
        echo -e "${RED}❌ Stats retrieval failed${NC}"
        return 1
    fi

    local keyword_count=$(echo "$response" | jq -r '.data.categoryKeywordCount // 0')
    echo -e "${GREEN}Total keywords:${NC} $keyword_count"

    # Show keyword breakdown by type
    local main_count=$(echo "$response" | jq -r '.data.keywordsByType.main | length' 2>/dev/null || echo "0")
    local brand_count=$(echo "$response" | jq -r '.data.keywordsByType.brand | length' 2>/dev/null || echo "0")
    local synonym_count=$(echo "$response" | jq -r '.data.keywordsByType.synonym | length' 2>/dev/null || echo "0")

    echo -e "${GREEN}Breakdown:${NC} Main: $main_count, Brands: $brand_count, Synonyms: $synonym_count"
}

# Main execution based on parameters
case "${1:-help}" in
    "test")
        if [ -z "$2" ]; then
            echo "Usage: $0 test \"Product Name\" [\"Expected Category\"]"
            echo "Examples:"
            echo "  $0 test \"тыква\"                           # Auto-detection mode"
            echo "  $0 test \"Болгарка Makita\" \"Строительные инструменты\"  # With expectation"
            exit 1
        fi
        test_product_detailed "$2" "$3"
        ;;

    "generate")
        if [ -z "$2" ] || [ -z "$3" ]; then
            echo "Usage: $0 generate <category_id> \"Category Name\""
            echo "Example: $0 generate 1007 \"Строительные инструменты\""
            exit 1
        fi
        generate_keywords "$2" "$3"
        ;;

    "stats")
        if [ -z "$2" ]; then
            echo "Usage: $0 stats <category_id>"
            echo "Example: $0 stats 1007"
            exit 1
        fi
        check_keyword_stats "$2"
        ;;

    "full-test")
        category_id="$2"
        category_name="$3"

        if [ -z "$category_id" ] || [ -z "$category_name" ]; then
            echo "Usage: $0 full-test <category_id> \"Category Name\""
            echo "Example: $0 full-test 1007 \"Строительные инструменты\""
            exit 1
        fi

        echo -e "${YELLOW}🔄 Full test for $category_name${NC}"

        # Check current stats
        check_keyword_stats "$category_id"

        # Generate keywords if needed
        echo -e "\n${YELLOW}Generating additional keywords...${NC}"
        generate_keywords "$category_id" "$category_name"

        # Test some products
        echo -e "\n${YELLOW}Testing products...${NC}"
        case "$category_name" in
            *"инструмент"*|*"Строительн"*)
                test_product_detailed "Болгарка Makita 125мм" "$category_name"
                test_product_detailed "Дрель ударная Bosch" "$category_name"
                test_product_detailed "Циркулярная пила" "$category_name"
                ;;
            *"Автомоб"*|*"машин"*)
                test_product_detailed "BMW X5 2020" "$category_name"
                test_product_detailed "Mercedes E-Class" "$category_name"
                ;;
            *"телефон"*|*"смартфон"*)
                test_product_detailed "iPhone 15 Pro" "$category_name"
                test_product_detailed "Samsung Galaxy S24" "$category_name"
                ;;
            *"компьютер"*|*"ноутбук"*)
                test_product_detailed "MacBook Pro M3" "$category_name"
                test_product_detailed "ASUS ROG Laptop" "$category_name"
                ;;
            *)
                echo "Добавьте конкретные тестовые продукты для этой категории"
                ;;
        esac

        # Final stats
        echo -e "\n${YELLOW}Final statistics:${NC}"
        check_keyword_stats "$category_id"
        ;;

    "bulk-generate")
        echo -e "${YELLOW}🚀 Starting bulk keyword generation...${NC}"

        response=$(curl -s -X POST "$API_URL/generate-keywords-all?minKeywords=50")

        if [ $? -ne 0 ]; then
            echo -e "${RED}❌ Bulk generation failed${NC}"
            exit 1
        fi

        local message=$(echo "$response" | jq -r '.data.message // "Unknown response"')
        local categories_found=$(echo "$response" | jq -r '.data.categoriesFound // 0')

        echo -e "${GREEN}Result:${NC} $message"
        echo -e "${GREEN}Categories to process:${NC} $categories_found"
        ;;

    "help"|*)
        echo "AI Category Testing Tool"
        echo ""
        echo "Commands:"
        echo "  test \"Product\" [\"Expected Category\"]  - Test single product detection"
        echo "  generate <id> \"Category\"               - Generate keywords for category"
        echo "  stats <id>                             - Show keyword statistics"
        echo "  full-test <id> \"Category\"              - Complete test (stats + generate + test)"
        echo "  bulk-generate                          - Generate keywords for all categories"
        echo ""
        echo "Examples:"
        echo "  $0 test \"тыква\"                                    # Auto-detection"
        echo "  $0 test \"Болгарка Makita\" \"Строительные инструменты\"  # With expectation"
        echo "  $0 test \"iPhone 15\"                               # Auto-detection"
        echo "  $0 generate 1007 \"Строительные инструменты\""
        echo "  $0 stats 1007"
        echo "  $0 full-test 1007 \"Строительные инструменты\""
        echo "  $0 bulk-generate"
        ;;
esac