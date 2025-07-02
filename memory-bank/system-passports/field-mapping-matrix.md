# Матрица соответствия полей между слоями

## Принципы маппинга

### Соглашения о наименованиях по слоям
- **PostgreSQL**: `snake_case` (user_id, created_at)
- **Backend Go**: `PascalCase` поля + `snake_case` JSON теги
- **Frontend TS**: `snake_case` для API данных (сохраняется как есть)
- **OpenSearch**: `snake_case` с дополнительными полями для поиска
- **MinIO**: `kebab-case` для paths, `snake_case` для metadata

## User Entity

| PostgreSQL | Backend Go | JSON API | Frontend TS | OpenSearch | MinIO |
|------------|------------|----------|-------------|------------|-------|
| `id` | `ID int` | `id` | `id: number` | N/A | N/A |
| `name` | `Name string` | `name` | `name: string` | N/A | N/A |
| `email` | `Email string` | `email` | `email: string` | N/A | N/A |
| `google_id` | `GoogleID *string` | `google_id` | `google_id?: string` | N/A | N/A |
| `picture_url` | `PictureURL *string` | `picture_url` | `picture_url?: string` | N/A | N/A |
| `phone` | `Phone *string` | `phone` | `phone?: string` | N/A | N/A |
| `created_at` | `CreatedAt time.Time` | `created_at` | `created_at: string` | N/A | N/A |
| `updated_at` | `UpdatedAt time.Time` | `updated_at` | `updated_at: string` | N/A | N/A |

## Marketplace Listing Entity

| PostgreSQL | Backend Go | JSON API | Frontend TS | OpenSearch | MinIO |
|------------|------------|----------|-------------|------------|-------|
| `id` | `ID int` | `id` | `id: number` | `id` | N/A |
| `user_id` | `UserID int` | `user_id` | `user_id: number` | `user_id` | N/A |
| `category_id` | `CategoryID *int` | `category_id` | `category_id?: number` | `category_id` | N/A |
| `storefront_id` | `StorefrontID *int` | `storefront_id` | `storefront_id?: number` | `storefront_id` | N/A |
| `title` | `Title string` | `title` | `title: string` | `title` + `title_*` (multilang) | N/A |
| `description` | `Description string` | `description` | `description: string` | `description` + `description_*` | N/A |
| `price` | `Price *decimal.Decimal` | `price` | `price?: string` | `price` (double) | N/A |
| `old_price` | `OldPrice *decimal.Decimal` | `old_price` | `old_price?: string` | `old_price` (double) | N/A |
| `currency` | `Currency string` | `currency` | `currency: string` | `currency` | N/A |
| `status` | `Status string` | `status` | `status: string` | `status` | N/A |
| `location` | `Location *string` | `location` | `location?: string` | `location` + `location.keyword` | N/A |
| `city` | `City *string` | `city` | `city?: string` | `city` + `city.keyword` | N/A |
| `country` | `Country *string` | `country` | `country?: string` | `country` + `country.keyword` | N/A |
| `latitude` | `Latitude *float64` | `latitude` | `latitude?: number` | `coordinates.lat` | N/A |
| `longitude` | `Longitude *float64` | `longitude` | `longitude?: number` | `coordinates.lon` | N/A |
| `views_count` | `ViewsCount int` | `views_count` | `views_count: number` | `views_count` | N/A |
| `created_at` | `CreatedAt time.Time` | `created_at` | `created_at: string` | `created_at` | N/A |
| `updated_at` | `UpdatedAt time.Time` | `updated_at` | `updated_at: string` | `updated_at` | N/A |
| `deleted_at` | `DeletedAt *time.Time` | `deleted_at` | `deleted_at?: string` | N/A (не индексируется) | N/A |

## Marketplace Images

| PostgreSQL | Backend Go | JSON API | Frontend TS | OpenSearch | MinIO |
|------------|------------|----------|-------------|------------|-------|
| `id` | `ID int` | `id` | `id: number` | N/A | N/A |
| `listing_id` | `ListingID int` | `listing_id` | `listing_id: number` | `images[].listing_id` | N/A |
| `file_path` | `FilePath string` | `file_path` | `file_path: string` | `images[].file_path` | `object_path` |
| `file_name` | `FileName string` | `file_name` | `file_name: string` | `images[].file_name` | `original_name` |
| `file_size` | `FileSize int` | `file_size` | `file_size: number` | `images[].file_size` | `size` (metadata) |
| `content_type` | `ContentType string` | `content_type` | `content_type: string` | `images[].content_type` | `content_type` (metadata) |
| `storage_type` | `StorageType string` | `storage_type` | `storage_type: string` | `images[].storage_type` | N/A |
| `public_url` | `PublicURL string` | `public_url` | `public_url: string` | `images[].public_url` | `public_url` |

## Chat System

| PostgreSQL | Backend Go | JSON API | Frontend TS | OpenSearch | MinIO (chat-files) |
|------------|------------|----------|-------------|------------|-------------------|
| `marketplace_chats.id` | `ID int` | `id` | `id: number` | N/A | N/A |
| `marketplace_chats.listing_id` | `ListingID int` | `listing_id` | `listing_id: number` | N/A | N/A |
| `marketplace_chats.buyer_id` | `BuyerID int` | `buyer_id` | `buyer_id: number` | N/A | N/A |
| `marketplace_chats.seller_id` | `SellerID int` | `seller_id` | `seller_id: number` | N/A | N/A |
| `marketplace_messages.id` | `ID int` | `id` | `id: number` | N/A | `{messageID}_*` |
| `marketplace_messages.chat_id` | `ChatID int` | `chat_id` | `chat_id: number` | N/A | N/A |
| `marketplace_messages.sender_id` | `SenderID int` | `sender_id` | `sender_id: number` | N/A | N/A |
| `marketplace_messages.content` | `Content string` | `content` | `content: string` | N/A | N/A |
| `marketplace_messages.has_attachments` | `HasAttachments bool` | `has_attachments` | `has_attachments: boolean` | N/A | N/A |

## Category Attributes (EAV Model)

| PostgreSQL | Backend Go | JSON API | Frontend TS | OpenSearch | MinIO |
|------------|------------|----------|-------------|------------|-------|
| `category_attributes.id` | `ID int` | `id` | `id: number` | N/A | N/A |
| `category_attributes.name` | `Name string` | `name` | `name: string` | `attributes[].name` | N/A |
| `category_attributes.display_name` | `DisplayName string` | `display_name` | `display_name: string` | `attributes[].display_name` | N/A |
| `category_attributes.attribute_type` | `AttributeType string` | `attribute_type` | `attribute_type: string` | `attributes[].type` | N/A |
| `listing_attribute_values.text_value` | `TextValue *string` | `text_value` | `text_value?: string` | `attributes[].text_value` | N/A |
| `listing_attribute_values.numeric_value` | `NumericValue *decimal.Decimal` | `numeric_value` | `numeric_value?: string` | `attributes[].numeric_value` | N/A |
| `listing_attribute_values.boolean_value` | `BooleanValue *bool` | `boolean_value` | `boolean_value?: boolean` | `attributes[].boolean_value` | N/A |
| `listing_attribute_values.json_value` | `JsonValue json.RawMessage` | `json_value` | `json_value?: any` | `attributes[].json_value` | N/A |

## Специальные поля OpenSearch

| Источник | OpenSearch Field | Назначение |
|----------|------------------|------------|
| Calculated | `make_lowercase` | Поиск марки без учета регистра |
| Calculated | `model_lowercase` | Поиск модели без учета регистра |
| Calculated | `brand_lowercase` | Поиск бренда без учета регистра |
| Calculated | `car_keywords[]` | Объединенные ключевые слова для авто |
| Calculated | `title_suggest` | Completion suggester |
| Calculated | `title_variations[]` | Вариации названий |
| Calculated | `all_attributes_text` | Объединенный текст всех атрибутов |
| lat + lon | `coordinates` (geo_point) | Геопоиск |
| Translations | `title_sr`, `title_ru`, `title_en` | Мультиязычные поля |

## MinIO Path Conventions

| Тип файла | Path Pattern | Пример |
|-----------|--------------|--------|
| Listing images | `listings/{filename}` | `listings/1704123456_product.jpg` |
| Chat attachments | `{type}/{year}/{month}/{day}/{messageId}_{timestamp}_{filename}` | `image/2024/01/15/123_1705334400_photo.jpg` |
| Review photos | `review-photos/{filename}` | `review-photos/1704123456_review.jpg` |

## Ключевые особенности маппинга

### 1. Консистентность
- **Snake_case**: Единый стандарт для всех API взаимодействий
- **ID Fields**: Всегда `{entity}_id` pattern
- **Timestamps**: ISO 8601 строки в JSON

### 2. Type Transformations
- **PostgreSQL Decimal** → **Go decimal.Decimal** → **JSON string** → **TS string**
- **PostgreSQL TIMESTAMP** → **Go time.Time** → **JSON ISO string** → **TS string**
- **PostgreSQL nullable** → **Go pointer** → **JSON omitempty** → **TS optional**

### 3. OpenSearch Enhancements
- **Denormalization**: Связанные данные встраиваются в документ
- **Search Fields**: Дополнительные `_lowercase` и `_keyword` поля
- **Geo Fields**: Координаты объединяются в `geo_point`
- **Multilang**: Отдельные поля для каждого языка

### 4. MinIO Organization
- **Bucket Separation**: Разные типы файлов в разных buckets
- **Path Structure**: Логическая иерархия для организации
- **Metadata Consistency**: Стандартные поля во всех объектах

### 5. Frontend Type Safety
- **Generated Types**: Автоматическая генерация из OpenAPI
- **Snake_case Preservation**: Никаких трансформаций в camelCase
- **Optional Fields**: Правильная обработка nullable полей

Эта матрица обеспечивает полную трассируемость данных через все слои системы и является основой для выявления и устранения несоответствий.