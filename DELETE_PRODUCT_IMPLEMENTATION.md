# DeleteProduct gRPC Handler Implementation

## Overview
Реализован полный функционал удаления продуктов в listings микросервисе с поддержкой soft и hard delete.

## Реализованные компоненты

### 1. Proto Definition (API Contract)
**Файл:** `/p/github.com/sveturs/listings/api/proto/listings/v1/listings.proto`

```protobuf
// RPC Method
rpc DeleteProduct(DeleteProductRequest) returns (DeleteProductResponse);

// Request
message DeleteProductRequest {
  int64 product_id = 1;
  int64 storefront_id = 2; // Ownership validation
  bool hard_delete = 3; // false = soft delete (default)
}

// Response
message DeleteProductResponse {
  bool success = 1;
  optional string message = 2;
  int32 variants_deleted = 3; // Cascade count
}
```

### 2. Repository Layer
**Файл:** `/p/github.com/sveturs/listings/internal/repository/postgres/products_repository.go`

**Метод:** `DeleteProduct(ctx, productID, storefrontID, hardDelete) (int32, error)`

**Функциональность:**
- ✅ **Ownership validation** - проверка что продукт принадлежит storefront
- ✅ **Soft delete** - `UPDATE deleted_at = NOW()` (fallback: `is_active = false`)
- ✅ **Hard delete** - `DELETE` с CASCADE на variants
- ✅ **Variants counting** - подсчет удаленных вариантов
- ✅ **Transaction safety** - все операции в одной транзакции
- ⏳ **Active orders check** - TODO (когда появится orders microservice)

**SQL операции:**

Soft delete (по умолчанию):
```sql
UPDATE b2c_products
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = $1 AND storefront_id = $2 AND deleted_at IS NULL
```

Fallback (если нет deleted_at):
```sql
UPDATE b2c_products
SET is_active = false, updated_at = NOW()
WHERE id = $1 AND storefront_id = $2
```

Hard delete:
```sql
DELETE FROM b2c_products
WHERE id = $1 AND storefront_id = $2
-- CASCADE автоматически удалит b2c_product_variants
```

### 3. Service Layer
**Файл:** `/p/github.com/sveturs/listings/internal/service/listings/service.go`

**Интерфейс:**
```go
DeleteProduct(ctx context.Context, productID, storefrontID int64, hardDelete bool) (int32, error)
```

**Функциональность:**
- Валидация входных параметров
- Делегирование в repository
- Логирование операций
- Сохранение error placeholders

### 4. gRPC Handler
**Файл:** `/p/github.com/sveturs/listings/internal/transport/grpc/handlers_products.go`

**Handler:** `DeleteProduct(ctx, req) (*DeleteProductResponse, error)`

**Валидация:**
- `product_id > 0`
- `storefront_id > 0`

**Error Mapping:**
| Repository Error | gRPC Code | Error Placeholder |
|-----------------|-----------|-------------------|
| `products.not_found` | `NotFound` | `products.not_found` |
| `products.has_active_orders` | `FailedPrecondition` | `products.has_active_orders` |
| Other | `Internal` | `products.delete_failed` |

**Response:**
```go
{
  Success: true,
  Message: "Product soft/hard deleted successfully",
  VariantsDeleted: <count>
}
```

## Безопасность

### Ownership Validation
Каждая операция удаления проверяет:
```sql
SELECT EXISTS(
  SELECT 1 FROM b2c_products
  WHERE id = $1 AND storefront_id = $2
)
```

Если продукт не принадлежит storefront → `products.not_found`

### Active Orders Protection (TODO)
Планируется добавить проверку:
```sql
-- Example (when orders microservice available)
SELECT EXISTS(
  SELECT 1 FROM orders
  WHERE product_id = $1
    AND status IN ('pending', 'processing', 'shipped')
)
```

Если есть активные заказы → `products.has_active_orders` (FailedPrecondition)

## Database Schema Requirements

### Минимальные требования (работает сейчас):
```sql
CREATE TABLE b2c_products (
  id BIGSERIAL PRIMARY KEY,
  storefront_id BIGINT NOT NULL,
  is_active BOOLEAN DEFAULT true,
  -- ... other fields
);
```

### Рекомендуемая схема (для полного функционала):
```sql
CREATE TABLE b2c_products (
  id BIGSERIAL PRIMARY KEY,
  storefront_id BIGINT NOT NULL,
  deleted_at TIMESTAMP,
  is_active BOOLEAN DEFAULT true,
  -- ... other fields
);

CREATE TABLE b2c_product_variants (
  id BIGSERIAL PRIMARY KEY,
  product_id BIGINT REFERENCES b2c_products(id) ON DELETE CASCADE,
  -- ... other fields
);
```

**Ключевой момент:** `ON DELETE CASCADE` для автоматического удаления вариантов при hard delete.

## Использование

### Soft Delete (по умолчанию):
```go
// gRPC call
resp, err := client.DeleteProduct(ctx, &pb.DeleteProductRequest{
  ProductId:    123,
  StorefrontId: 456,
  HardDelete:   false, // или не указывать (default)
})
```

### Hard Delete:
```go
// gRPC call
resp, err := client.DeleteProduct(ctx, &pb.DeleteProductRequest{
  ProductId:    123,
  StorefrontId: 456,
  HardDelete:   true, // FULL DELETE
})
```

### Response:
```go
if resp.Success {
  fmt.Printf("Deleted! Variants: %d\n", resp.VariantsDeleted)
  fmt.Println(*resp.Message) // "Product soft/hard deleted successfully"
}
```

## Логирование

### Debug уровень:
```
deleting product product_id=123 storefront_id=456 hard_delete=false
```

### Info уровень (успех):
```
product soft deleted successfully product_id=123 variants_count=5
product hard deleted successfully product_id=123 variants_deleted=5
```

### Error уровень:
```
failed to delete product product_id=123 error="products.not_found"
failed to check product ownership error="connection refused"
```

## Testing

### Тестовые сценарии:

#### 1. Успешный soft delete:
```bash
# Продукт существует, принадлежит storefront, нет активных заказов
product_id=100, storefront_id=10 → Success, variants_deleted=3
```

#### 2. Успешный hard delete:
```bash
# Продукт существует, принадлежит storefront, нет активных заказов
product_id=100, storefront_id=10, hard_delete=true → Success, variants_deleted=3
```

#### 3. Продукт не найден:
```bash
# Продукт не существует ИЛИ не принадлежит storefront
product_id=999, storefront_id=10 → NotFound: products.not_found
```

#### 4. Неверные параметры:
```bash
product_id=0 → InvalidArgument: "product ID must be greater than 0"
storefront_id=0 → InvalidArgument: "storefront ID must be greater than 0"
```

#### 5. Активные заказы (TODO):
```bash
# Когда будет orders microservice
product_id=100 (has pending orders) → FailedPrecondition: products.has_active_orders
```

## TODO / Future Enhancements

1. **Active Orders Check:**
   - Интеграция с orders microservice
   - Проверка pending/processing/shipped заказов
   - Блокировка удаления при активных заказах

2. **Migration for deleted_at:**
   ```sql
   ALTER TABLE b2c_products ADD COLUMN deleted_at TIMESTAMP;
   CREATE INDEX idx_b2c_products_deleted_at ON b2c_products(deleted_at)
   WHERE deleted_at IS NULL;
   ```

3. **Audit Log:**
   - Логирование кто, когда и какие продукты удалил
   - История восстановления

4. **Bulk Delete:**
   - Удаление нескольких продуктов за раз
   - Batch операции

5. **Restore Function:**
   - Восстановление soft-deleted продуктов
   - `RestoreProduct(product_id, storefront_id)`

## Error Placeholders для Frontend

```javascript
// В frontend/svetu/src/messages/ru/products.json
{
  "not_found": "Продукт не найден",
  "delete_failed": "Не удалось удалить продукт",
  "has_active_orders": "Невозможно удалить продукт с активными заказами"
}
```

## Dependencies

- PostgreSQL >= 12 (для JSONB, CASCADE constraints)
- Go >= 1.21
- gRPC protobuf definitions
- zerolog для логирования
- lib/pq для PostgreSQL driver

## Компиляция

Всё компилируется успешно:
```bash
cd /p/github.com/sveturs/listings
go build ./internal/transport/grpc/...   # ✅ OK
go build ./internal/service/listings/... # ✅ OK
go build ./internal/repository/postgres/... # ✅ OK
```

## Status

- ✅ Proto definitions
- ✅ Repository implementation
- ✅ Service layer
- ✅ gRPC handler
- ✅ Error handling with placeholders
- ✅ Ownership validation
- ✅ Soft delete support
- ✅ Hard delete with CASCADE
- ✅ Compilation successful
- ⏳ Active orders check (pending orders microservice)
- ⏳ Migration for deleted_at column
- ⏳ Unit tests

---

**Автор:** Claude Code
**Дата:** 2025-11-04
**Микросервис:** listings v1
**Статус:** ✅ Ready for Testing
