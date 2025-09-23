#!/bin/bash

# Скрипт для тестирования новой системы AI категоризации

API_URL="http://localhost:3000/api/v1/marketplace/ai"

# Цвета для вывода
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Тестирование новой системы AI категоризации ===${NC}"

# Тест 1: Керамическая ваза
echo -e "\n${BLUE}Тест 1: Керамическая ваза-осьминог${NC}"
curl -s -X POST "$API_URL/select-category" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Керамическая ваза-осьминог ручной работы",
    "description": "Уникальная керамическая ваза в форме осьминога. Ручная работа, высота 30см. Идеально подходит для украшения интерьера. Цвет: морская волна с золотистыми акцентами."
  }' | jq .

# Тест 2: iPhone 15 Pro
echo -e "\n${BLUE}Тест 2: iPhone 15 Pro${NC}"
curl -s -X POST "$API_URL/select-category" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "iPhone 15 Pro Max 256GB",
    "description": "Новый iPhone 15 Pro Max, цвет титановый синий, 256GB памяти. В комплекте зарядка и чехол."
  }' | jq .

# Тест 3: Женское платье
echo -e "\n${BLUE}Тест 3: Вечернее платье${NC}"
curl -s -X POST "$API_URL/detect-category" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Вечернее платье с пайетками",
    "description": "Красивое вечернее платье черного цвета с пайетками. Размер M (44-46). Длина миди, открытая спина."
  }' | jq .

# Тест 4: Mercedes-Benz
echo -e "\n${BLUE}Тест 4: Mercedes-Benz E-Class${NC}"
curl -s -X POST "$API_URL/detect-category" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Mercedes-Benz E220d 2019",
    "description": "Mercedes E220d, 2019 год, дизель 2.0, автомат, пробег 45000км. Полная комплектация, состояние отличное."
  }' | jq .

# Тест 5: Квартира (проверка что НЕ попадет в автомобили)
echo -e "\n${BLUE}Тест 5: Квартира в Новом Саде${NC}"
curl -s -X POST "$API_URL/detect-category" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Квартира 2-комнатная в центре Нового Сада",
    "description": "Продается 2-комнатная квартира 65м2 в центре Нового Сада. 3 этаж, лифт, балкон, центральное отопление."
  }' | jq .

echo -e "\n${GREEN}✅ Тестирование завершено!${NC}"