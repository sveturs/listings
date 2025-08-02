# Система вариативных атрибутов

## Обзор

Система вариативных атрибутов позволяет создавать товары с различными вариантами (например, разные цвета и размеры), каждый из которых может иметь свою цену и остаток на складе.

## Архитектура

### База данных

#### Таблица `product_variant_attributes`
Хранит определения доступных вариативных атрибутов:

```sql
CREATE TABLE product_variant_attributes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    display_name VARCHAR(200) NOT NULL,
    type VARCHAR(50) DEFAULT 'select',
    is_required BOOLEAN DEFAULT FALSE,
    affects_stock BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

#### Таблица `storefront_product_variants`
Хранит варианты товаров с комбинациями атрибутов:

```sql
CREATE TABLE storefront_product_variants (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL,
    sku VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    stock_quantity INTEGER DEFAULT 0,
    variant_attributes JSONB, -- Хранит комбинации: {"color": "Red", "size": "XL"}
    is_default BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (product_id) REFERENCES storefront_products(id) ON DELETE CASCADE,
    UNIQUE(product_id, sku)
);
```

### Доступные вариативные атрибуты

1. **color** - Цвет товара
2. **size** - Размер
3. **material** - Материал
4. **pattern** - Узор/рисунок
5. **style** - Стиль
6. **memory** - Оперативная память (для электроники)
7. **storage** - Объем памяти/хранилища
8. **connectivity** - Тип подключения
9. **bundle** - Комплектация
10. **capacity** - Объем/вместимость
11. **power** - Мощность

### Маппинг категорий к атрибутам

| Категория | Вариативные атрибуты |
|-----------|---------------------|
| smartphones | color, memory, storage |
| womens-clothing | color, size, material, pattern, style |
| mens-clothing | color, size, material, pattern, style |
| kids-clothing | color, size, material, pattern |
| sports-clothing | color, size, material |
| shoes | color, size, material, style |
| bags | color, size, material, style, pattern |
| accessories | color, size, material, style, pattern |
| computers | color, memory, storage, connectivity |
| gaming-consoles | color, storage, bundle |
| electronics-accessories | color, connectivity, bundle |
| home-appliances | color, capacity, power |
| furniture | color, material, style |
| kitchenware | color, capacity, material |

## Backend API

### Эндпоинты

#### GET /api/v1/marketplace/product-variant-attributes
Возвращает список всех доступных вариативных атрибутов.

**Ответ:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "color",
      "display_name": "Color",
      "type": "select",
      "is_required": false,
      "affects_stock": false
    },
    // ...
  ]
}
```

#### GET /api/v1/marketplace/categories/{slug}/variant-attributes
Возвращает вариативные атрибуты для конкретной категории.

**Параметры:**
- `slug` - slug категории (например, "smartphones")

**Ответ:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "color",
      "display_name": "Color",
      "type": "select",
      "is_required": false,
      "affects_stock": false
    },
    {
      "id": 6,
      "name": "memory",
      "display_name": "Memory (RAM)",
      "type": "select",
      "is_required": false,
      "affects_stock": true
    },
    {
      "id": 7,
      "name": "storage",
      "display_name": "Storage",
      "type": "select",
      "is_required": false,
      "affects_stock": true
    }
  ]
}
```

### Сервисный слой

Файл: `/backend/internal/proj/marketplace/service/variant_attributes.go`

```go
// GetCategoryVariantAttributes возвращает вариативные атрибуты для категории
func (s *Service) GetCategoryVariantAttributes(ctx context.Context, categorySlug string) ([]models.ProductVariantAttribute, error) {
    // Получаем список атрибутов для категории из карты
    attributeNames, exists := categoryVariantAttributesMap[categorySlug]
    if !exists {
        return []models.ProductVariantAttribute{}, nil
    }
    
    // Получаем полную информацию об атрибутах из БД
    return s.storage.GetProductVariantAttributesByNames(ctx, attributeNames)
}
```

## Frontend компоненты

### SimplifiedVariantGenerator

Компонент для генерации вариантов товара на основе выбранных атрибутов.

**Расположение:** `/frontend/svetu/src/components/products/SimplifiedVariantGenerator.tsx`

**Основные функции:**
1. Загрузка доступных вариативных атрибутов для категории
2. Фильтрация атрибутов категории для отображения только вариативных
3. Генерация всех возможных комбинаций вариантов
4. Управление ценами и остатками для каждого варианта

**Процесс работы:**
1. При загрузке компонент получает slug категории
2. Загружает список вариативных атрибутов для категории
3. Фильтрует атрибуты товара, оставляя только те, которые являются вариативными
4. При нажатии "Сгенерировать варианты" создает все комбинации
5. Позволяет настроить цену и остаток для каждого варианта

## Известные проблемы и ограничения

### 1. Проблема с отображением атрибутов color и storage для смартфонов

**Симптомы:** При создании товара в категории smartphones в вариантах отображается только атрибут Memory.

**Причина:** Атрибуты color и storage не были заполнены на шаге заполнения атрибутов товара.

**Решение:** Необходимо убедиться, что на шаге "Атрибуты" при создании товара выбраны значения для всех вариативных атрибутов категории.

### 2. Временное решение с захардкоженными данными

В текущей реализации frontend использует захардкоженные данные вместо API из-за проблем с аутентификацией публичных эндпоинтов.

**TODO:** 
- Исправить конфигурацию middleware для публичных эндпоинтов вариативных атрибутов
- Удалить захардкоженные данные из SimplifiedVariantGenerator
- Реализовать полноценную загрузку через API

### 3. Отсутствие множественного выбора в полях атрибутов

**Проблема:** Пользователь не может выбрать несколько значений для атрибута (например, несколько цветов).

**TODO:** Реализовать поддержку multiselect для вариативных атрибутов.

## Рекомендации по дальнейшему развитию

1. **Динамическое управление атрибутами**
   - Создать админ-панель для управления вариативными атрибутами
   - Возможность добавления новых атрибутов без изменения кода
   - Настройка маппинга категорий к атрибутам через UI

2. **Улучшение UX**
   - Добавить превью вариантов при создании
   - Массовое редактирование цен и остатков
   - Импорт/экспорт вариантов в CSV

3. **Оптимизация производительности**
   - Кэширование списка атрибутов
   - Пагинация для большого количества вариантов
   - Lazy loading для компонентов управления вариантами

4. **Расширение функциональности**
   - Поддержка изображений для каждого варианта
   - История изменения цен и остатков
   - Автоматическая генерация SKU
   - Поддержка вложенных вариантов (например, размер + цвет + материал)

## Миграции

Для добавления новых вариативных атрибутов используйте миграцию:

```sql
-- Пример добавления нового атрибута
INSERT INTO product_variant_attributes (name, display_name, type, is_required, affects_stock)
VALUES ('new_attribute', 'New Attribute', 'select', false, false);

-- Добавление переводов
INSERT INTO translations (entity_type, entity_id, entity_field, language_code, translation)
VALUES 
  ('product_variant_attribute', (SELECT id FROM product_variant_attributes WHERE name = 'new_attribute'), 'display_name', 'ru', 'Новый атрибут'),
  ('product_variant_attribute', (SELECT id FROM product_variant_attributes WHERE name = 'new_attribute'), 'display_name', 'sr', 'Нови атрибут');
```