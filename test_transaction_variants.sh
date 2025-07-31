#!/bin/bash

# Тестирование создания товара с вариантами в транзакции

API_URL="http://localhost:3000/api/v1"
STOREFRONT_ID=1
USER_ID=1

# Получение токена аутентификации (замените на ваш способ получения токена)
# Для тестирования используем заглушку
AUTH_TOKEN="test-token"

echo "=== Тестирование создания товара с вариантами ==="
echo ""

# 1. Тест: Создание товара с вариантами (успешный сценарий)
echo "1. Создание товара с вариантами (успешный сценарий)..."
curl -X POST "$API_URL/storefronts/$STOREFRONT_ID/products" \
  -H "Content-Type: application/json" \
  -H "X-User-ID: $USER_ID" \
  -d '{
    "name": "Тестовый товар с вариантами",
    "description": "Описание товара",
    "price": 1000,
    "currency": "USD",
    "category_id": 1,
    "stock_quantity": 100,
    "is_active": true,
    "has_variants": true,
    "variants": [
      {
        "sku": "TEST-VAR-001",
        "price": 900,
        "stock_quantity": 50,
        "variant_attributes": {
          "size": "S",
          "color": "Red"
        },
        "is_default": true
      },
      {
        "sku": "TEST-VAR-002",
        "price": 1100,
        "stock_quantity": 30,
        "variant_attributes": {
          "size": "M",
          "color": "Blue"
        },
        "is_default": false
      },
      {
        "sku": "TEST-VAR-003",
        "price": 1200,
        "stock_quantity": 20,
        "variant_attributes": {
          "size": "L",
          "color": "Green"
        },
        "is_default": false
      }
    ]
  }' | jq '.'

echo ""
echo "2. Тест с дублированным SKU (должен вызвать ошибку и откат транзакции)..."
curl -X POST "$API_URL/storefronts/$STOREFRONT_ID/products" \
  -H "Content-Type: application/json" \
  -H "X-User-ID: $USER_ID" \
  -d '{
    "name": "Товар с дублированным SKU",
    "description": "Должен вызвать ошибку",
    "price": 2000,
    "currency": "USD",
    "category_id": 1,
    "stock_quantity": 50,
    "is_active": true,
    "has_variants": true,
    "variants": [
      {
        "sku": "DUPLICATE-001",
        "price": 1900,
        "stock_quantity": 25,
        "variant_attributes": {
          "size": "S"
        },
        "is_default": true
      },
      {
        "sku": "DUPLICATE-001",
        "price": 2100,
        "stock_quantity": 25,
        "variant_attributes": {
          "size": "M"
        },
        "is_default": false
      }
    ]
  }' | jq '.'

echo ""
echo "3. Проверка, что товар из теста 2 не был создан..."
echo "Для проверки используйте SQL:"
echo "psql \"postgres://postgres:password@localhost:5432/svetubd?sslmode=disable\" -c \"SELECT id, name FROM storefront_products WHERE name = 'Товар с дублированным SKU';\""