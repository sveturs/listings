# Паспорт таблицы: marketplace_listings

## Назначение
Основная таблица объявлений маркетплейса. Хранит все товары и услуги, размещенные пользователями на платформе Sve Tu.

## Структура таблицы

```sql
CREATE TABLE marketplace_listings (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    category_id INT REFERENCES marketplace_categories(id),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(12,2),
    condition VARCHAR(50),
    status VARCHAR(20) DEFAULT 'active',
    location VARCHAR(255),
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    address_city VARCHAR(100),
    address_country VARCHAR(100),
    views_count INT DEFAULT 0,
    show_on_map BOOLEAN NOT NULL DEFAULT true,
    original_language VARCHAR(10) DEFAULT 'sr',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    storefront_id INT REFERENCES user_storefronts(id) ON DELETE SET NULL,
    external_id VARCHAR(255),
    metadata JSONB,
    needs_reindex BOOLEAN DEFAULT false
);
```

## Поля таблицы

### Основные поля
- `id` - уникальный идентификатор объявления (SERIAL)
- `user_id` - владелец объявления (FK к users)
- `category_id` - категория товара (FK к marketplace_categories)
- `title` - заголовок объявления, обязательный (до 255 символов)
- `description` - подробное описание товара/услуги

### Цена и состояние
- `price` - цена товара (до 12 цифр, 2 знака после запятой)
- `condition` - состояние товара: 'new', 'used', 'refurbished'
- `status` - статус объявления: 'active', 'inactive', 'sold', 'draft'

### Местоположение
- `location` - текстовое описание местоположения
- `latitude` - широта (10 цифр, 8 после запятой)
- `longitude` - долгота (11 цифр, 8 после запятой)
- `address_city` - город
- `address_country` - страна
- `show_on_map` - показывать ли на карте (по умолчанию true)

### Статистика и метаданные
- `views_count` - количество просмотров
- `original_language` - язык оригинала (по умолчанию 'sr')
- `metadata` - дополнительные данные в JSON формате

### Интеграция
- `storefront_id` - привязка к витрине магазина
- `external_id` - идентификатор из внешней системы
- `needs_reindex` - флаг необходимости переиндексации в OpenSearch

### Системные поля
- `created_at` - дата создания
- `updated_at` - дата последнего обновления

## Индексы

### Основные индексы
1. **idx_marketplace_listings_status** - по статусу
2. **idx_marketplace_listings_storefront** - по витрине
3. **idx_marketplace_listings_external_id** - по внешнему ID
4. **idx_marketplace_listings_external_id_storefront_id** - составной индекс

### Оптимизационные индексы
5. **idx_listings_metadata_discount** - GIN индекс по скидкам в metadata
6. **idx_marketplace_listings_status_created** - для активных по дате создания
7. **idx_marketplace_listings_location** - геопространственный индекс
8. **idx_marketplace_listings_category_status** - по категории и статусу
9. **idx_marketplace_listings_user_status** - по пользователю и статусу
10. **idx_marketplace_listings_city** - по городу
11. **idx_marketplace_listings_title_gin** - полнотекстовый поиск по заголовку
12. **idx_marketplace_listings_price** - по цене для активных

## Триггеры

1. **trg_new_listing_price_history** - создает запись в price_history при создании
2. **trg_update_listing_price_history** - обновляет price_history при изменении цены
3. **update_marketplace_listings_updated_at** - обновляет updated_at

## Связи с другими таблицами

### Прямые связи (эта таблица ссылается на)
- `user_id` → `users.id` - владелец объявления
- `category_id` → `marketplace_categories.id` - категория
- `storefront_id` → `user_storefronts.id` - витрина (может быть NULL)

### Обратные связи (другие таблицы ссылаются на marketplace_listings)
- `marketplace_images.listing_id` - изображения объявления
- `marketplace_favorites.listing_id` - избранные объявления
- `marketplace_chats.listing_id` - чаты по объявлению
- `listing_attribute_values.listing_id` - значения атрибутов
- `listing_views.listing_id` - просмотры объявления
- `price_history.listing_id` - история цен
- `reviews` - отзывы о сделках

## Бизнес-правила

### Статусы объявлений
- `active` - активное, видно всем
- `inactive` - временно скрыто владельцем
- `sold` - продано
- `draft` - черновик, не опубликовано

### Валидация
1. **Обязательные поля** - title обязателен
2. **Цена** - может быть NULL (договорная)
3. **Геолокация** - latitude и longitude должны быть заполнены вместе
4. **Статус по умолчанию** - 'active' для новых объявлений

### Metadata структура
```json
{
  "discount": {
    "percentage": 20,
    "old_price": 100.00
  },
  "shipping": {
    "available": true,
    "price": 5.00
  },
  "warranty": "12 months",
  "brand": "Samsung",
  "model": "Galaxy S21"
}
```

## Примеры использования

### Создание объявления
```sql
INSERT INTO marketplace_listings (
    user_id, category_id, title, description, price, 
    condition, address_city, address_country
) VALUES (
    1, 5, 'iPhone 13 Pro Max', 'Отличное состояние...', 
    899.99, 'used', 'Белград', 'Сербия'
);
```

### Поиск активных объявлений
```sql
SELECT l.*, u.name as seller_name, c.name as category_name
FROM marketplace_listings l
JOIN users u ON l.user_id = u.id
JOIN marketplace_categories c ON l.category_id = c.id
WHERE l.status = 'active'
  AND l.price BETWEEN 100 AND 1000
  AND l.address_city = 'Белград'
ORDER BY l.created_at DESC;
```

### Обновление статуса при продаже
```sql
UPDATE marketplace_listings 
SET status = 'sold', updated_at = CURRENT_TIMESTAMP
WHERE id = 123;
```

### Полнотекстовый поиск
```sql
SELECT * FROM marketplace_listings
WHERE to_tsvector('simple', title) @@ to_tsquery('simple', 'iphone')
  AND status = 'active';
```

## Известные особенности

1. **Автоматическая история цен** - триггеры отслеживают изменения цены
2. **Переиндексация** - флаг needs_reindex для синхронизации с OpenSearch
3. **Мягкое удаление витрин** - ON DELETE SET NULL для storefront_id
4. **JSONB metadata** - гибкое хранение дополнительных атрибутов
5. **Полнотекстовый поиск** - GIN индекс по заголовку
6. **Оптимизированные индексы** - множество составных индексов для производительности

## Миграции

- **000001** - создание таблицы
- **000005** - добавление storefront_id
- **000012** - добавление external_id
- **000013** - добавление metadata
- **000017** - добавление needs_reindex
- **000039** - добавление оптимизационных индексов

## Интеграция с другими компонентами

1. **OpenSearch** - основной поиск и фильтрация
2. **MinIO** - хранение изображений
3. **Price History** - автоматическое отслеживание цен
4. **Import System** - импорт через external_id
5. **Chat System** - коммуникация покупатель-продавец