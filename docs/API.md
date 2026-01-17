# Listings Service API Documentation

## Overview

Listings Service предоставляет gRPC API для управления товарами, избранным, заказами и корзиной.

**gRPC Endpoint:** `localhost:50053`
**HTTP Endpoint:** `localhost:8086`
**Database:** `listings_dev_db` (PostgreSQL port 35434)

---

## gRPC Services

### ListingsService

Управление товарами (listings).

#### Core Methods

- `CreateListing` - создание нового товара
- `GetListing` - получение товара по ID
- `UpdateListing` - обновление товара
- `DeleteListing` - удаление товара
- `ListListings` - список товаров с фильтрацией

#### Search & Filtering

- `SearchListings` - полнотекстовый поиск через OpenSearch
- `GetListingsByCategory` - товары по категории
- `GetListingsBySeller` - товары продавца

#### Images

- `UploadImage` - загрузка изображения товара
- `DeleteImage` - удаление изображения
- `ReorderImages` - изменение порядка изображений
- `SetPrimaryImage` - установка главного изображения

#### Attributes

- `UpdateListingAttributes` - обновление атрибутов товара
- `GetCategoryAttributes` - получение атрибутов категории

---

### FavoritesService

Управление избранным пользователя.

#### Methods

- `AddToFavorites` - добавить в избранное
- `RemoveFromFavorites` - удалить из избранного
- `GetFavorites` - получить список избранного
- `IsFavorite` - проверить наличие в избранном

**Request Example:**
```protobuf
message AddToFavoritesRequest {
  int64 user_id = 1;
  int64 listing_id = 2;
}
```

**Response Example:**
```protobuf
message AddToFavoritesResponse {
  bool success = 1;
  string message = 2;
}
```

---

### CartService

Управление корзиной покупок.

#### Methods

- `AddToCart` - добавить товар в корзину
- `RemoveFromCart` - удалить из корзины
- `GetCart` - получить содержимое корзины
- `UpdateCartItem` - обновить количество товара
- `ClearCart` - очистить корзину

---

### OrdersService

Управление заказами.

#### Methods

- `CreateOrder` - создание заказа из корзины
- `GetOrder` - получение заказа по ID
- `ListOrders` - список заказов пользователя
- `UpdateOrderStatus` - обновление статуса заказа
- `CancelOrder` - отмена заказа

#### Order Statuses

- `pending` - ожидает оплаты
- `confirmed` - подтверждён
- `processing` - в обработке
- `shipped` - отправлен
- `delivered` - доставлен
- `cancelled` - отменён
- `refunded` - возврат средств

---

### ChatService

Чаты по объявлениям.

#### Methods

- `CreateChat` - создание чата по товару
- `GetChat` - получение чата по ID
- `ListChats` - список чатов пользователя
- `SendMessage` - отправка сообщения в чат
- `GetMessages` - получение истории сообщений
- `MarkAsRead` - пометить сообщения как прочитанные

---

## OpenSearch Integration

### Index Structure

**Index:** `listings_microservice`

**Mapping:**
```json
{
  "title": "text",
  "description": "text",
  "category_id": "keyword",
  "seller_id": "keyword",
  "price": "float",
  "quantity": "integer",
  "status": "keyword",
  "attributes": "object",
  "created_at": "date"
}
```

### Search Query Example

```bash
curl -X POST "localhost:9200/listings_microservice/_search" \
  -H "Content-Type: application/json" \
  -d '{
    "query": {
      "multi_match": {
        "query": "laptop",
        "fields": ["title^3", "description"]
      }
    },
    "filter": {
      "term": {"status": "active"}
    },
    "size": 20
  }'
```

---

## REST API (через монолит)

Frontend использует BFF proxy через монолит:

### Favorites

```
GET    /api/v1/marketplace/favorites
POST   /api/v1/marketplace/favorites/:listing_id
DELETE /api/v1/marketplace/favorites/:listing_id
```

### Listings

```
GET    /api/v1/marketplace/listings
GET    /api/v1/marketplace/listings/:id
POST   /api/v1/marketplace/listings
PUT    /api/v1/marketplace/listings/:id
DELETE /api/v1/marketplace/listings/:id
```

### Cart

```
GET    /api/v1/marketplace/cart
POST   /api/v1/marketplace/cart
DELETE /api/v1/marketplace/cart/:item_id
PATCH  /api/v1/marketplace/cart/:item_id
```

### Orders

```
GET    /api/v1/marketplace/orders
GET    /api/v1/marketplace/orders/:id
POST   /api/v1/marketplace/orders
PATCH  /api/v1/marketplace/orders/:id/status
DELETE /api/v1/marketplace/orders/:id
```

---

## Database Schema

### Main Tables

- `listings` - товары (C2C + B2C unified)
- `listing_favorites` - избранное
- `listing_images` - изображения товаров
- `listing_locations` - геолокация
- `listing_attributes` - атрибуты товаров
- `categories` - категории (synced from monolith)
- `attributes` - атрибуты категорий (synced from monolith)
- `category_attributes` - связи категорий и атрибутов
- `chats` - чаты
- `messages` - сообщения
- `cart_items` - корзина
- `orders` - заказы
- `storefronts` - витрины магазинов

---

## Data Synchronization

### Categories & Attributes Sync

**Скрипт синхронизации:**
```bash
python3 /p/github.com/vondi-global/listings/scripts/sync_listings_data.py
```

**Что синхронизируется:**
- ✅ Categories: `c2c_categories` → `categories`
- ✅ Attributes: `unified_attributes` → `attributes` (VARCHAR → JSONB)
- ✅ Category Attributes: `unified_category_attributes` → `category_attributes`

**Когда запускать:**
- После добавления категорий в монолите
- После добавления атрибутов
- После изменения связей категория-атрибут

---

## OpenSearch Reindexing

### Create Index

```bash
python3 /p/github.com/vondi-global/listings/scripts/create_opensearch_index.py
```

### Reindex Listings

```bash
python3 /p/github.com/vondi-global/listings/scripts/reindex_listings.py \
  --target-port 35434 \
  --target-password listings_secret \
  --target-db listings_dev_db
```

**Важно:** Frontend ищет товары через OpenSearch, не напрямую из PostgreSQL!

---

## Image Upload

### MinIO S3 Integration

**Endpoint:** `s3.vondi.rs`
**Bucket:** `vondi-marketplace-listings`

**Upload Process:**
1. Frontend получает presigned URL от backend
2. Upload файла напрямую в MinIO
3. Backend сохраняет URL в `listing_images`

---

## Authentication

**Все методы требуют валидного JWT токена из Auth Service.**

gRPC metadata:
```go
md := metadata.Pairs("authorization", "Bearer "+token)
ctx := metadata.NewOutgoingContext(context.Background(), md)
```

---

## Error Handling

### gRPC Status Codes

- `INVALID_ARGUMENT` (3) - невалидные параметры
- `NOT_FOUND` (5) - listing не найден
- `PERMISSION_DENIED` (7) - нет прав на операцию
- `ALREADY_EXISTS` (6) - ресурс уже существует
- `INTERNAL` (13) - внутренняя ошибка

---

## Metrics

Prometheus метрики на `/metrics`:

```
listings_created_total
listings_updated_total
listings_deleted_total
favorites_added_total
favorites_removed_total
cart_operations_total
orders_created_total
```

---

## Examples

### Add to Favorites (gRPC)

```go
import pb "github.com/vondi-global/listings/api/gen/listings/v1"

conn, _ := grpc.Dial("localhost:50053", grpc.WithInsecure())
client := pb.NewFavoritesServiceClient(conn)

resp, err := client.AddToFavorites(ctx, &pb.AddToFavoritesRequest{
    UserId:    123,
    ListingId: 456,
})
```

### Search Listings (REST via monolith)

```bash
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:3000/api/v1/marketplace/listings?q=laptop&category=electronics&min_price=500&max_price=2000"
```

---

## Health Checks

```
GET /health
GET /ready
GET /metrics
```

**Response:**
```json
{
  "status": "ok",
  "database": "connected",
  "opensearch": "connected",
  "redis": "connected"
}
```

---

## See Also

- [Database Architecture](DATABASE_ARCHITECTURE.md)
- [Migration Plan](MIGRATION_TO_MICROSERVICE_EXECUTIVE_SUMMARY.md)
- [OpenSearch Monitoring](OPENSEARCH_MONITORING_QUICKSTART.md)
- [README](../README.md)
