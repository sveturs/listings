# Event Publisher для Listings Service

## Описание

Event Publisher публикует события заказов в Redis Streams для межсервисной коммуникации. WMS (Warehouse Management Service) слушает эти события и создаёт соответствующие задачи.

## Архитектура

```
Listings Service → Redis Stream (listings:events:orders) → WMS Consumer
```

## События

### 1. order.confirmed

Публикуется когда заказ подтверждён и готов к фулфилменту.

**Формат:**
```json
{
  "type": "order.confirmed",
  "order_id": "12345",
  "storefront_id": "999",
  "items": "[{\"listing_id\":100,\"quantity\":2,\"warehouse_id\":1}]",
  "timestamp": "1701709200"
}
```

### 2. order.cancelled

Публикуется когда заказ отменён.

**Формат:**
```json
{
  "type": "order.cancelled",
  "order_id": "12345",
  "reason": "customer_request",
  "timestamp": "1701709200"
}
```

## Использование

### Инициализация

```go
import (
    "github.com/vondi-global/listings/internal/events"
    "github.com/redis/go-redis/v9"
    "github.com/rs/zerolog"
)

// Создать Redis клиент
redisClient := redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
})

// Создать publisher
publisher := events.NewRedisOrderEventPublisher(
    redisClient,
    logger,
    events.OrdersStream,
    1, // default warehouse ID
)
defer publisher.Close()
```

### Публикация события order.confirmed

```go
items := []events.OrderItem{
    {
        ListingID:   100,
        Quantity:    2,
        WarehouseID: 1, // если 0 - будет использован default warehouse
    },
    {
        ListingID:   101,
        Quantity:    1,
        WarehouseID: 0, // будет заменён на default warehouse
    },
}

err := publisher.PublishOrderConfirmed(ctx, orderID, storefrontID, items)
if err != nil {
    log.Error().Err(err).Msg("failed to publish order confirmed event")
}
```

### Публикация события order.cancelled

```go
err := publisher.PublishOrderCancelled(ctx, orderID, "customer_request")
if err != nil {
    log.Error().Err(err).Msg("failed to publish order cancelled event")
}
```

## Интеграция с gRPC сервисом

В `cmd/server/main.go`:

```go
// После инициализации Redis клиента
eventPublisher := events.NewRedisOrderEventPublisher(
    redisClient,
    logger,
    events.OrdersStream,
    cfg.DefaultWarehouseID,
)
defer eventPublisher.Close()

// Передать publisher в service layer
orderService := service.NewOrderService(
    orderRepo,
    cartRepo,
    listingRepo,
    eventPublisher, // <-- добавить как зависимость
    logger,
)
```

В `internal/service/order_service.go`:

```go
type OrderService struct {
    orderRepo      repository.OrderRepository
    cartRepo       repository.CartRepository
    listingRepo    repository.ListingRepository
    eventPublisher events.OrderEventPublisher // <-- добавить поле
    logger         zerolog.Logger
}

func (s *OrderService) ConfirmOrder(ctx context.Context, orderID int64) error {
    // ... логика подтверждения заказа ...

    // Опубликовать событие
    items := make([]events.OrderItem, len(order.Items))
    for i, item := range order.Items {
        items[i] = events.OrderItem{
            ListingID:   item.ListingID,
            Quantity:    item.Quantity,
            WarehouseID: item.WarehouseID,
        }
    }

    if err := s.eventPublisher.PublishOrderConfirmed(ctx, orderID, order.StorefrontID, items); err != nil {
        s.logger.Error().Err(err).Int64("order_id", orderID).Msg("failed to publish order confirmed event")
        // НЕ возвращаем ошибку - событие не критично для подтверждения заказа
    }

    return nil
}
```

## Конфигурация

Добавить в `.env`:

```bash
# Redis для событий (может быть тот же что и для кеша)
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# Default warehouse ID для новых заказов
DEFAULT_WAREHOUSE_ID=1
```

## Мониторинг

### Проверка событий в Redis

```bash
# Посмотреть все события в стриме
redis-cli XRANGE listings:events:orders - +

# Посмотреть длину стрима
redis-cli XLEN listings:events:orders

# Посмотреть последние 10 событий
redis-cli XREVRANGE listings:events:orders + - COUNT 10

# Посмотреть consumer groups (если WMS использует groups)
redis-cli XINFO GROUPS listings:events:orders
```

### Логирование

Publisher логирует все публикации:

```
{"level":"info","component":"order_event_publisher","message_id":"1701709200-0","order_id":12345,"storefront_id":999,"items_count":2,"stream":"listings:events:orders","message":"published order.confirmed event"}
```

## Обработка ошибок

- Publisher НЕ использует retry механизм - ответственность за retry на стороне вызывающего кода
- Ошибки Redis логируются и возвращаются caller'у
- Рекомендуется НЕ блокировать основной flow при ошибке публикации события

## Тестирование

```bash
# Запустить unit тесты
cd /p/github.com/sveturs/listings
go test -v ./internal/events

# Проверить покрытие
go test -cover ./internal/events
```

## WMS Consumer

WMS слушает стрим `listings:events:orders` и обрабатывает события:

- `order.confirmed` → создаёт PickingTask
- `order.cancelled` → отменяет PickingTask (если ещё не выполнена)

Подробности см. в `/p/github.com/sveturs/warehouse/internal/consumer/`
