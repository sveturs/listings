# Паспорт таблицы `storefront_products`

## Назначение
Каталог товаров для витрин магазинов. Содержит продукты, которые продаются через персональные витрины пользователей с расширенными возможностями управления инвентарем и атрибутами.

## Полная структура таблицы

```sql
CREATE TABLE IF NOT EXISTS storefront_products (
    id SERIAL PRIMARY KEY,
    storefront_id INTEGER NOT NULL REFERENCES storefronts(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    price DECIMAL(10, 2) NOT NULL CHECK (price >= 0),
    currency CHAR(3) NOT NULL DEFAULT 'USD',
    category_id INTEGER NOT NULL REFERENCES marketplace_categories(id),
    sku VARCHAR(100),
    barcode VARCHAR(100),
    stock_quantity INTEGER NOT NULL DEFAULT 0 CHECK (stock_quantity >= 0),
    stock_status VARCHAR(20) NOT NULL DEFAULT 'in_stock' CHECK (stock_status IN ('in_stock', 'low_stock', 'out_of_stock')),
    is_active BOOLEAN NOT NULL DEFAULT true,
    attributes JSONB DEFAULT '{}',
    view_count INTEGER NOT NULL DEFAULT 0,
    sold_count INTEGER NOT NULL DEFAULT 0,
    external_id VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

## Описание полей

| Поле | Тип | Обязательность | Описание |
|------|-----|---------------|----------|
| `id` | SERIAL | PRIMARY KEY | Уникальный ID товара |
| `storefront_id` | INTEGER | NOT NULL FK | ID витрины магазина |
| `name` | VARCHAR(255) | NOT NULL | Название товара |
| `description` | TEXT | NOT NULL | Подробное описание |
| `price` | DECIMAL(10,2) | NOT NULL | Цена товара (>= 0) |
| `currency` | CHAR(3) | NOT NULL | Валюта (USD, EUR, RSD) |
| `category_id` | INTEGER | NOT NULL FK | ID категории товара |
| `sku` | VARCHAR(100) | NULLABLE | Артикул товара |
| `barcode` | VARCHAR(100) | NULLABLE | Штрих-код товара |
| `stock_quantity` | INTEGER | NOT NULL | Количество на складе (>= 0) |
| `stock_status` | VARCHAR(20) | NOT NULL | Статус склада |
| `is_active` | BOOLEAN | NOT NULL | Активность товара |
| `attributes` | JSONB | DEFAULT '{}' | Дополнительные атрибуты |
| `view_count` | INTEGER | NOT NULL | Количество просмотров |
| `sold_count` | INTEGER | NOT NULL | Количество продаж |
| `external_id` | VARCHAR(100) | NULLABLE | Внешний ID (для импорта) |
| `created_at` | TIMESTAMP WITH TIME ZONE | DEFAULT NOW | Дата создания |
| `updated_at` | TIMESTAMP WITH TIME ZONE | DEFAULT NOW | Дата обновления |

## Статусы склада

| Статус | Описание |
|--------|----------|
| `in_stock` | Товар в наличии |
| `low_stock` | Мало товара |
| `out_of_stock` | Товар закончился |

## Индексы

```sql
-- Основные индексы
CREATE INDEX idx_storefront_products_storefront_id ON storefront_products(storefront_id);
CREATE INDEX idx_storefront_products_category_id ON storefront_products(category_id);
CREATE INDEX idx_storefront_products_stock_status ON storefront_products(stock_status);
CREATE INDEX idx_storefront_products_is_active ON storefront_products(is_active);

-- Поиск по артикулам и штрих-кодам
CREATE INDEX idx_storefront_products_sku ON storefront_products(sku) WHERE sku IS NOT NULL;
CREATE INDEX idx_storefront_products_barcode ON storefront_products(barcode) WHERE barcode IS NOT NULL;
CREATE INDEX idx_storefront_products_external_id ON storefront_products(external_id) WHERE external_id IS NOT NULL;

-- Полнотекстовый поиск по названию
CREATE INDEX idx_storefront_products_name_gin ON storefront_products USING gin(to_tsvector('simple', name));

-- Уникальные ограничения
CREATE UNIQUE INDEX unique_storefront_product_sku ON storefront_products (storefront_id, sku) WHERE sku IS NOT NULL;
CREATE UNIQUE INDEX unique_storefront_product_barcode ON storefront_products (storefront_id, barcode) WHERE barcode IS NOT NULL;
```

## Ограничения

- **CHECK**: `price >= 0` - цена не может быть отрицательной
- **CHECK**: `stock_quantity >= 0` - количество не может быть отрицательным
- **CHECK**: `stock_status IN ('in_stock', 'low_stock', 'out_of_stock')`
- **FOREIGN KEY**: `storefront_id` → `storefronts(id)` ON DELETE CASCADE
- **FOREIGN KEY**: `category_id` → `marketplace_categories(id)`
- **UNIQUE** (частичный): `(storefront_id, sku)` где `sku` не NULL
- **UNIQUE** (частичный): `(storefront_id, barcode)` где `barcode` не NULL

## Связи с другими таблицами

| Связь | Тип | Описание |
|-------|-----|----------|
| `storefront_id` → `storefronts.id` | Many-to-One | Витрина товара |
| `category_id` → `marketplace_categories.id` | Many-to-One | Категория товара |
| `storefront_product_images` | One-to-Many | Изображения товара |
| `storefront_product_variants` | One-to-Many | Варианты товара |
| `storefront_inventory_movements` | One-to-Many | Движения инвентаря |

## Бизнес-правила

1. **Уникальность артикулов**: В рамках одной витрины артикулы уникальны
2. **Управление складом**: Автоматическое обновление статуса по количеству
3. **Категоризация**: Товары должны быть привязаны к категориям
4. **Импорт**: Поддержка внешних ID для синхронизации
5. **Аналитика**: Отслеживание просмотров и продаж

## Примеры использования

### Добавление нового товара
```sql
INSERT INTO storefront_products (
    storefront_id, name, description, price, currency, category_id, sku, stock_quantity
) VALUES (
    42, 'iPhone 15 Pro', 'Новый iPhone 15 Pro 256GB', 1200.00, 'USD', 1, 'IPH15P256', 10
);
```

### Поиск товаров в наличии
```sql
SELECT * FROM storefront_products 
WHERE storefront_id = 42 
  AND is_active = true 
  AND stock_status = 'in_stock'
ORDER BY view_count DESC;
```

### Обновление статуса склада
```sql
UPDATE storefront_products 
SET stock_quantity = stock_quantity - 1,
    stock_status = CASE 
        WHEN stock_quantity - 1 = 0 THEN 'out_of_stock'
        WHEN stock_quantity - 1 <= 5 THEN 'low_stock'
        ELSE 'in_stock'
    END,
    sold_count = sold_count + 1
WHERE id = 123;
```

### Полнотекстовый поиск
```sql
SELECT * FROM storefront_products 
WHERE to_tsvector('simple', name) @@ to_tsquery('simple', 'iPhone')
  AND is_active = true;
```

## Атрибуты (JSONB поле)

Примеры структуры атрибутов:
```json
{
  "brand": "Apple",
  "model": "iPhone 15 Pro",
  "color": "Natural Titanium",
  "storage": "256GB",
  "condition": "new",
  "warranty": "1 year",
  "dimensions": {
    "width": 70.6,
    "height": 146.6,
    "depth": 8.25
  },
  "weight": 187
}
```

## Известные особенности

1. **Каскадное удаление**: При удалении витрины удаляются все товары
2. **Мультивалютность**: Поддержка разных валют
3. **Складской учет**: Автоматическое управление статусами
4. **Импорт-дружелюбность**: Поддержка внешних ID и массовых операций
5. **Аналитика**: Встроенные счетчики просмотров и продаж

## Использование в коде

**Backend**:
- Handler: `internal/proj/storefront_products/`
- Модели: `internal/domain/storefront_product.go`
- API: `/api/v1/storefronts/{id}/products`

**Frontend**:
- Компоненты: `src/components/storefront/products/`
- Страницы: `src/app/[locale]/storefront/[slug]/products/`
- Типы: `@/types/storefront-product.ts`

## Производительность

- **Индексирование**: Оптимизировано для поиска по витрине, категории, статусу
- **Полнотекстовый поиск**: GIN индекс для быстрого поиска по названию
- **Партиционирование**: Возможно по storefront_id для больших витрин

## Связанные компоненты

- **Storefront images**: Изображения товаров
- **Inventory management**: Управление складом
- **Category attributes**: Атрибуты категорий
- **Import system**: Система импорта товаров
- **Analytics**: Аналитика продаж и просмотров