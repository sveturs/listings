# Спецификация системы вариантов товаров

## Обзор

Система вариантов позволяет создавать товары с различными модификациями (размер, цвет, материал и т.д.), где каждый вариант имеет свои характеристики: цену, остаток на складе, SKU и другие атрибуты.

## Архитектура

### Основные компоненты

1. **StorefrontProduct** - основной товар
2. **ProductVariant** - вариант товара с уникальными характеристиками
3. **VariantAttributes** - атрибуты варианта (размер, цвет и т.д.)

### База данных

#### Таблица `storefront_products`
- `has_variants` (boolean) - флаг наличия вариантов
- При `has_variants = true` основной товар становится "контейнером" для вариантов

#### Таблица `storefront_product_variants`
```sql
CREATE TABLE storefront_product_variants (
    id SERIAL PRIMARY KEY,
    product_id INTEGER REFERENCES storefront_products(id),
    sku VARCHAR(255) UNIQUE,
    barcode VARCHAR(255),
    price DECIMAL(10,2),
    compare_at_price DECIMAL(10,2),
    cost_price DECIMAL(10,2),
    stock_quantity INTEGER DEFAULT 0,
    reserved_quantity INTEGER DEFAULT 0,
    available_quantity INTEGER GENERATED ALWAYS AS (stock_quantity - reserved_quantity) STORED,
    stock_status VARCHAR(50),
    low_stock_threshold INTEGER,
    variant_attributes JSONB,
    weight DECIMAL(10,3),
    dimensions JSONB,
    is_active BOOLEAN DEFAULT true,
    is_default BOOLEAN DEFAULT false,
    view_count INTEGER DEFAULT 0,
    sold_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## API Endpoints

### Создание товара с вариантами

**POST** `/api/v1/storefronts/{storefront_id}/products`

```json
{
  "name": "Футболка базовая",
  "description": "Классическая хлопковая футболка",
  "price": 1500,
  "currency": "RUB",
  "category_id": 15,
  "stock_quantity": 150,
  "is_active": true,
  "has_variants": true,
  "variants": [
    {
      "sku": "TSHIRT-BLACK-S",
      "price": 1500,
      "stock_quantity": 50,
      "variant_attributes": {
        "color": "Черный",
        "size": "S"
      },
      "is_default": true
    },
    {
      "sku": "TSHIRT-BLACK-M",
      "price": 1500,
      "stock_quantity": 50,
      "variant_attributes": {
        "color": "Черный",
        "size": "M"
      }
    },
    {
      "sku": "TSHIRT-WHITE-L",
      "price": 1600,
      "compare_at_price": 2000,
      "stock_quantity": 30,
      "low_stock_threshold": 10,
      "variant_attributes": {
        "color": "Белый",
        "size": "L"
      },
      "weight": 0.25,
      "dimensions": {
        "length": 70,
        "width": 50,
        "height": 2
      }
    }
  ]
}
```

### Получение товара с вариантами

**GET** `/api/v1/storefronts/{storefront_id}/products/{product_id}`

Ответ включает массив `variants` с полной информацией о каждом варианте.

### Управление вариантами

#### Получение вариантов товара
**GET** `/api/v1/storefronts/products/{product_id}/variants`

#### Создание нового варианта
**POST** `/api/v1/storefronts/products/{product_id}/variants`

#### Обновление варианта
**PUT** `/api/v1/storefronts/products/{product_id}/variants/{variant_id}`

#### Удаление варианта
**DELETE** `/api/v1/storefronts/products/{product_id}/variants/{variant_id}`

#### Генерация вариантов
**POST** `/api/v1/storefronts/products/{product_id}/variants/generate`

Автоматическая генерация вариантов на основе матрицы атрибутов:
```json
{
  "base_sku": "PRODUCT",
  "base_price": 1000,
  "attribute_matrix": {
    "size": ["S", "M", "L", "XL"],
    "color": ["Черный", "Белый", "Серый"]
  },
  "pricing_rules": {
    "size": {
      "XL": 100
    },
    "color": {
      "Белый": 50
    }
  }
}
```

## Транзакционная целостность

Создание товара с вариантами происходит в единой транзакции:
1. Создается основной товар
2. Создаются все варианты
3. При ошибке на любом этапе - откат всей операции

## Особенности реализации

### Вариант по умолчанию
- Только один вариант может быть помечен как `is_default = true`
- При установке нового варианта по умолчанию, предыдущий автоматически сбрасывается

### Статус наличия
`stock_status` автоматически определяется на основе:
- `out_of_stock` - если `stock_quantity = 0`
- `low_stock` - если `stock_quantity <= low_stock_threshold`
- `in_stock` - в остальных случаях

### Доступное количество
`available_quantity` вычисляется автоматически:
```sql
available_quantity = stock_quantity - reserved_quantity
```

### Изображения вариантов
Каждый вариант может иметь свои изображения через таблицу `storefront_product_variant_images`

## Примеры использования

### 1. Одежда
```json
{
  "variant_attributes": {
    "size": "M",
    "color": "Синий",
    "material": "Хлопок 100%"
  }
}
```

### 2. Электроника
```json
{
  "variant_attributes": {
    "storage": "256GB",
    "color": "Space Gray",
    "model": "Pro"
  }
}
```

### 3. Мебель
```json
{
  "variant_attributes": {
    "width": "140см",
    "material": "Дуб",
    "configuration": "Левый угол"
  }
}
```

## Валидация

### При создании товара
- Если `has_variants = true`, массив `variants` обязателен и не может быть пустым
- Хотя бы один вариант должен быть помечен как `is_default`

### При создании варианта
- `SKU` должен быть уникальным в рамках всей системы
- `stock_quantity` не может быть отрицательным
- `variant_attributes` обязателен и должен содержать хотя бы один атрибут

## Поиск и фильтрация

Варианты индексируются в OpenSearch вместе с основным товаром, что позволяет:
- Искать по SKU варианта
- Фильтровать по атрибутам вариантов
- Показывать диапазон цен для товара с вариантами

## Рекомендации

1. **Используйте варианты когда:**
   - Товар имеет модификации с разными SKU
   - Необходим отдельный учет остатков для каждой модификации
   - Модификации имеют разные цены

2. **Не используйте варианты когда:**
   - Различия только в описании
   - Нет необходимости в отдельном учете остатков
   - Все модификации имеют одинаковую цену и характеристики

3. **Оптимизация:**
   - Не создавайте более 100 вариантов для одного товара
   - Используйте осмысленные SKU для удобства управления
   - Устанавливайте `low_stock_threshold` для автоматического отслеживания остатков