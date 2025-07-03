# Паспорт Marketplace Models (Backend)

## MarketplaceListing (основная модель товара)

### Go структура
```go
type MarketplaceListing struct {
    ID            int                    `json:"id" db:"id"`
    UserID        int                    `json:"user_id" db:"user_id"`
    CategoryID    *int                   `json:"category_id,omitempty" db:"category_id"`
    StorefrontID  *int                   `json:"storefront_id,omitempty" db:"storefront_id"`
    Title         string                 `json:"title" db:"title"`
    Description   string                 `json:"description" db:"description"`
    Price         *decimal.Decimal       `json:"price,omitempty" db:"price"`
    OldPrice      *decimal.Decimal       `json:"old_price,omitempty" db:"old_price"`
    Currency      string                 `json:"currency" db:"currency"`
    Status        string                 `json:"status" db:"status"`
    Location      *string                `json:"location,omitempty" db:"location"`
    City          *string                `json:"city,omitempty" db:"city"`
    Country       *string                `json:"country,omitempty" db:"country"`
    Latitude      *float64               `json:"latitude,omitempty" db:"latitude"`
    Longitude     *float64               `json:"longitude,omitempty" db:"longitude"`
    ViewsCount    int                    `json:"views_count" db:"views_count"`
    CreatedAt     time.Time              `json:"created_at" db:"created_at"`
    UpdatedAt     time.Time              `json:"updated_at" db:"updated_at"`
    DeletedAt     *time.Time             `json:"deleted_at,omitempty" db:"deleted_at"`
    
    // Связанные данные (загружаются отдельно)
    Images        []MarketplaceImage     `json:"images,omitempty"`
    Attributes    []ListingAttributeValue `json:"attributes,omitempty"`
    Category      *MarketplaceCategory   `json:"category,omitempty"`
    User          *User                  `json:"user,omitempty"`
}
```

### Статусы listing
- `active` - активный
- `inactive` - неактивный  
- `sold` - продан
- `draft` - черновик

### Валюты
- `RSD` - сербский динар (по умолчанию)
- `EUR` - евро
- `USD` - доллар США

## MarketplaceImage (изображения товаров)

### Go структура
```go
type MarketplaceImage struct {
    ID          int       `json:"id" db:"id"`
    ListingID   int       `json:"listing_id" db:"listing_id"`
    FilePath    string    `json:"file_path" db:"file_path"`
    FileName    string    `json:"file_name" db:"file_name"`
    FileSize    int       `json:"file_size" db:"file_size"`
    ContentType string    `json:"content_type" db:"content_type"`
    StorageType string    `json:"storage_type" db:"storage_type"`
    PublicURL   string    `json:"public_url" db:"public_url"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
```

### Особенности
- **StorageType**: всегда `"minio"`
- **PublicURL**: полный URL для доступа к изображению
- **FilePath**: путь в MinIO bucket

## MarketplaceCategory (категории)

### Go структура
```go
type MarketplaceCategory struct {
    ID          int                    `json:"id" db:"id"`
    Name        string                 `json:"name" db:"name"`
    Slug        string                 `json:"slug" db:"slug"`
    ParentID    *int                   `json:"parent_id,omitempty" db:"parent_id"`
    Icon        *string                `json:"icon,omitempty" db:"icon"`
    IsActive    bool                   `json:"is_active" db:"is_active"`
    CreatedAt   time.Time              `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time              `json:"updated_at" db:"updated_at"`
    
    // Связанные данные
    Children    []MarketplaceCategory  `json:"children,omitempty"`
    Attributes  []CategoryAttribute    `json:"attributes,omitempty"`
    Translations []Translation         `json:"translations,omitempty"`
}
```

### Особенности иерархии
- Поддержка неограниченной вложенности через `parent_id`
- `slug` используется для URL-friendly идентификации
- Мультиязычная поддержка через `translations`

## CategoryAttribute (атрибуты категорий)
  
### Go структура
```go
type CategoryAttribute struct {
    ID              int                 `json:"id" db:"id"`
    CategoryID      int                 `json:"category_id" db:"category_id"`
    Name            string              `json:"name" db:"name"`
    DisplayName     string              `json:"display_name" db:"display_name"`
    AttributeType   string              `json:"attribute_type" db:"attribute_type"`
    IsRequired      bool                `json:"is_required" db:"is_required"`
    Options         json.RawMessage     `json:"options,omitempty" db:"options"`
    CreatedAt       time.Time           `json:"created_at" db:"created_at"`
    UpdatedAt       time.Time           `json:"updated_at" db:"updated_at"`
}
```

### Типы атрибутов
- `text` - текстовое поле
- `number` - числовое значение
- `select` - выбор из списка (options в JSON)
- `boolean` - да/нет
- `multiselect` - множественный выбор

## ListingAttributeValue (значения атрибутов)

### Go структура
```go
type ListingAttributeValue struct {
    ID             int             `json:"id" db:"id"`
    ListingID      int             `json:"listing_id" db:"listing_id"`
    AttributeID    int             `json:"attribute_id" db:"attribute_id"`
    TextValue      *string         `json:"text_value,omitempty" db:"text_value"`
    NumericValue   *decimal.Decimal `json:"numeric_value,omitempty" db:"numeric_value"`
    BooleanValue   *bool           `json:"boolean_value,omitempty" db:"boolean_value"`
    JsonValue      json.RawMessage `json:"json_value,omitempty" db:"json_value"`
    CreatedAt      time.Time       `json:"created_at" db:"created_at"`
    UpdatedAt      time.Time       `json:"updated_at" db:"updated_at"`
    
    // Связанные данные
    Attribute      *CategoryAttribute `json:"attribute,omitempty"`
}
```

### Хранение значений
- **TextValue**: для `text` и `select` типов
- **NumericValue**: для `number` типа
- **BooleanValue**: для `boolean` типа  
- **JsonValue**: для `multiselect` и сложных структур

## MarketplaceChat & MarketplaceMessage (чаты)

### MarketplaceChat
```go
type MarketplaceChat struct {
    ID            int                    `json:"id" db:"id"`
    ListingID     int                    `json:"listing_id" db:"listing_id"`
    BuyerID       int                    `json:"buyer_id" db:"buyer_id"`  
    SellerID      int                    `json:"seller_id" db:"seller_id"`
    Status        string                 `json:"status" db:"status"`
    CreatedAt     time.Time              `json:"created_at" db:"created_at"`
    UpdatedAt     time.Time              `json:"updated_at" db:"updated_at"`
    
    // Дополнительные поля для API
    UnreadCount   int                    `json:"unread_count,omitempty"`
    LastMessage   *MarketplaceMessage    `json:"last_message,omitempty"`
    Listing       *MarketplaceListing    `json:"listing,omitempty"`
    Buyer         *User                  `json:"buyer,omitempty"`
    Seller        *User                  `json:"seller,omitempty"`
}
```

### MarketplaceMessage
```go
type MarketplaceMessage struct {
    ID              int               `json:"id" db:"id"`
    ChatID          int               `json:"chat_id" db:"chat_id"`
    SenderID        int               `json:"sender_id" db:"sender_id"`
    Content         string            `json:"content" db:"content"`
    MessageType     string            `json:"message_type" db:"message_type"`
    IsRead          bool              `json:"is_read" db:"is_read"`
    HasAttachments  bool              `json:"has_attachments" db:"has_attachments"`
    CreatedAt       time.Time         `json:"created_at" db:"created_at"`
    
    // Связанные данные
    Attachments     []ChatAttachment  `json:"attachments,omitempty"`
    Sender          *User            `json:"sender,omitempty"`
}
```

### Типы сообщений
- `text` - обычное текстовое сообщение
- `image` - изображение
- `file` - файл

## MarketplaceFavorite (избранное)

### Go структура
```go
type MarketplaceFavorite struct {
    ID        int       `json:"id" db:"id"`
    UserID    int       `json:"user_id" db:"user_id"`
    ListingID int       `json:"listing_id" db:"listing_id"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}
```

## API Endpoints для Marketplace

### Listings
- `GET /marketplace/listings` - список товаров с фильтрами
- `POST /marketplace/listings` - создать товар
- `GET /marketplace/listings/{id}` - получить товар
- `PUT /marketplace/listings/{id}` - обновить товар
- `DELETE /marketplace/listings/{id}` - удалить товар
- `POST /marketplace/listings/{id}/images` - загрузить изображения

### Categories  
- `GET /marketplace/categories` - дерево категорий
- `GET /marketplace/categories/{id}/attributes` - атрибуты категории

### Chat
- `GET /marketplace/chats` - список чатов пользователя
- `POST /marketplace/chats` - создать чат
- `GET /marketplace/chats/{id}/messages` - сообщения чата
- `POST /marketplace/chats/{id}/messages` - отправить сообщение

### Favorites
- `GET /marketplace/favorites` - избранные товары
- `POST /marketplace/favorites` - добавить в избранное
- `DELETE /marketplace/favorites/{listing_id}` - удалить из избранного

## Соответствие с другими слоями

### PostgreSQL
- Прямое соответствие с таблицами `marketplace_*`
- FK constraints обеспечивают целостность данных

### Frontend  
- Типы генерируются автоматически из OpenAPI
- Snake_case сохраняется в frontend без изменений

### OpenSearch
- Listing индексируется с денормализованными данными
- Атрибуты встраиваются в документ для быстрого поиска
- Координаты преобразуются в geo_point тип

### MinIO
- Изображения хранятся в bucket `listings`
- PublicURL формируется для прямого доступа

## Ключевые особенности

1. **Гибкие атрибуты**: EAV модель для категориальных атрибутов
2. **Мультимедиа**: Поддержка множественных изображений
3. **Геолокация**: Координаты для карт и местоположения  
4. **Чаты**: Встроенная система обмена сообщениями
5. **Версионность**: Soft delete через `deleted_at`