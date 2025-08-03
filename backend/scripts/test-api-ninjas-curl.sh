#!/bin/bash

# API-Ninjas Cars API Test Script
# Тестирование европейских марок автомобилей

echo "=== API-Ninjas Cars API Test ==="
echo "Для работы скрипта нужен API ключ"
echo "Получите бесплатный ключ на: https://api-ninjas.com/register"
echo ""

# Проверяем наличие API ключа
if [ -z "$API_NINJAS_KEY" ]; then
    echo "ERROR: API_NINJAS_KEY не установлен"
    echo "Запустите: export API_NINJAS_KEY='your-key-here'"
    echo ""
    echo "=== Демонстрация без API ключа ==="
    echo ""
    echo "Пример запроса для BMW 2023:"
    echo "curl -H 'X-Api-Key: YOUR_API_KEY' 'https://api.api-ninjas.com/v1/cars?make=BMW&year=2023'"
    echo ""
    echo "Ожидаемый ответ:"
    cat << 'EOF'
[
  {
    "city_mpg": 26,
    "class": "Compact Cars",
    "combination_mpg": 30,
    "cylinders": 4,
    "displacement": 2.0,
    "drive": "Rear-Wheel Drive",
    "fuel_type": "Premium Gasoline",
    "highway_mpg": 36,
    "make": "BMW",
    "model": "330i",
    "transmission": "Automatic",
    "year": 2023
  },
  {
    "city_mpg": 25,
    "class": "Compact Cars",
    "combination_mpg": 29,
    "cylinders": 4,
    "displacement": 2.0,
    "drive": "All-Wheel Drive",
    "fuel_type": "Premium Gasoline",
    "highway_mpg": 34,
    "make": "BMW",
    "model": "330i xDrive",
    "transmission": "Automatic",
    "year": 2023
  }
]
EOF
    exit 1
fi

# Функция для выполнения запроса
test_make() {
    local make=$1
    local year=$2
    echo ""
    echo "=== Тестирование $make ($year) ==="
    
    if [ -z "$year" ]; then
        url="https://api.api-ninjas.com/v1/cars?make=$make&limit=3"
    else
        url="https://api.api-ninjas.com/v1/cars?make=$make&year=$year&limit=3"
    fi
    
    response=$(curl -s -H "X-Api-Key: $API_NINJAS_KEY" "$url")
    
    if [ $? -eq 0 ]; then
        echo "$response" | jq '.' 2>/dev/null || echo "$response"
    else
        echo "ERROR: Не удалось выполнить запрос"
    fi
    
    sleep 1  # Rate limiting
}

# Тестируем европейские марки
echo ""
echo "=== Тестирование европейских марок ==="

# Немецкие марки
test_make "BMW" "2023"
test_make "Mercedes-Benz" "2023"
test_make "Audi" "2023"
test_make "Volkswagen" "2023"

# Французские марки
test_make "Peugeot" "2023"
test_make "Renault" "2023"
test_make "Citroën" ""  # Может не быть новых моделей

# Итальянские марки
test_make "Fiat" "2023"
test_make "Alfa Romeo" ""

# Испанские марки
test_make "Seat" ""

# Румынские марки
test_make "Dacia" ""

# Чешские марки
test_make "Škoda" ""

echo ""
echo "=== Тест завершен ==="
echo "Примечание: Пустые массивы [] означают, что марка/год не найдены в API"