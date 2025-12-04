# Интеграция EventPublisher в Listings Service

## Быстрый старт

### 1. Добавить в main.go

```go
import "github.com/vondi-global/listings/internal/events"

// После инициализации Redis клиента
eventPublisher := events.NewRedisOrderEventPublisher(
    redisClient,
    logger,
    events.OrdersStream,
    cfg.DefaultWarehouseID,
)
defer eventPublisher.Close()
```

### 2. Добавить в service

```go
type OrderService struct {
    // ... существующие поля ...
    eventPublisher events.OrderEventPublisher
}

func NewOrderService(
    orderRepo repository.OrderRepository,
    cartRepo repository.CartRepository,
    listingRepo repository.ListingRepository,
    eventPublisher events.OrderEventPublisher,
    logger zerolog.Logger,
) *OrderService {
    return &OrderService{
        orderRepo:      orderRepo,
        cartRepo:       cartRepo,
        listingRepo:    listingRepo,
        eventPublisher: eventPublisher,
        logger:         logger,
    }
}
```

### 3. Публиковать события

```go
// В методе подтверждения заказа
func (s *OrderService) ConfirmOrder(ctx context.Context, orderID int64) error {
    // ... бизнес-логика ...

    // Публикуем событие (НЕ блокируем при ошибке)
    items := s.prepareOrderItems(order)
    if err := s.eventPublisher.PublishOrderConfirmed(ctx, orderID, order.StorefrontID, items); err != nil {
        s.logger.Error().Err(err).Int64("order_id", orderID).Msg("failed to publish order confirmed event")
        // Продолжаем выполнение - событие не критично
    }

    return nil
}

// В методе отмены заказа
func (s *OrderService) CancelOrder(ctx context.Context, orderID int64, reason string) error {
    // ... бизнес-логика ...

    // Публикуем событие
    if err := s.eventPublisher.PublishOrderCancelled(ctx, orderID, reason); err != nil {
        s.logger.Error().Err(err).Int64("order_id", orderID).Msg("failed to publish order cancelled event")
    }

    return nil
}

func (s *OrderService) prepareOrderItems(order *Order) []events.OrderItem {
    items := make([]events.OrderItem, len(order.Items))
    for i, item := range order.Items {
        items[i] = events.OrderItem{
            ListingID:   item.ListingID,
            Quantity:    item.Quantity,
            WarehouseID: item.WarehouseID, // или 0 для default
        }
    }
    return items
}
```

## Конфигурация

Добавить в `.env`:

```bash
# Default warehouse для новых заказов (если не указан явно)
DEFAULT_WAREHOUSE_ID=1
```

## Тестирование

### Mock для тестов

```go
// Создать mock
type MockEventPublisher struct {
    mock.Mock
}

func (m *MockEventPublisher) PublishOrderConfirmed(ctx context.Context, orderID int64, storefrontID int64, items []events.OrderItem) error {
    args := m.Called(ctx, orderID, storefrontID, items)
    return args.Error(0)
}

func (m *MockEventPublisher) PublishOrderCancelled(ctx context.Context, orderID int64, reason string) error {
    args := m.Called(ctx, orderID, reason)
    return args.Error(0)
}

func (m *MockEventPublisher) Close() error {
    args := m.Called()
    return args.Error(0)
}

// Использовать в тестах
mockPublisher := new(MockEventPublisher)
mockPublisher.On("PublishOrderConfirmed", mock.Anything, int64(123), int64(999), mock.Anything).Return(nil)

service := NewOrderService(orderRepo, cartRepo, listingRepo, mockPublisher, logger)
```

## Проверка работы

### 1. Запустить Listings Service
```bash
cd /p/github.com/sveturs/listings
go run ./cmd/server/main.go
```

### 2. Проверить Redis Stream
```bash
# Посмотреть события
redis-cli XRANGE listings:events:orders - +

# Посмотреть длину стрима
redis-cli XLEN listings:events:orders
```

### 3. Запустить WMS Consumer
```bash
cd /p/github.com/sveturs/warehouse
go run ./cmd/server/main.go
```

## Troubleshooting

### События не публикуются

1. Проверить подключение к Redis:
```bash
redis-cli ping
```

2. Проверить логи Listings Service:
```bash
tail -f /tmp/listings-microservice.log | grep "order_event_publisher"
```

3. Проверить конфигурацию Redis в `.env`:
```bash
grep REDIS /p/github.com/sveturs/listings/.env
```

### WMS не получает события

1. Проверить что WMS использует правильный stream name:
```go
// Должно быть "listings:events:orders"
consumerStreamName := events.OrdersStream
```

2. Проверить логи WMS:
```bash
tail -f /tmp/warehouse.log | grep "consumer"
```

3. Проверить consumer group в Redis:
```bash
redis-cli XINFO GROUPS listings:events:orders
```

## Архитектура

```
┌─────────────────────┐
│  Listings Service   │
│                     │
│  ┌───────────────┐  │      ┌──────────────────┐
│  │ OrderService  │  │      │  Redis Streams   │      ┌──────────────┐
│  │               │─────────▶│                  │─────▶│     WMS      │
│  │ + eventPub    │  │      │ listings:events: │      │   Consumer   │
│  └───────────────┘  │      │    orders        │      └──────────────┘
│                     │      └──────────────────┘
└─────────────────────┘

События:
• order.confirmed → создать PickingTask
• order.cancelled → отменить PickingTask
```

## Мониторинг

### Prometheus метрики (TODO)

```go
// Добавить в будущем
var (
    eventsPublished = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "listings_events_published_total",
            Help: "Total number of events published",
        },
        []string{"event_type", "status"},
    )
)
```

### Health check

```go
// Добавить проверку Redis в health endpoint
func (s *Server) healthCheck(ctx *fiber.Ctx) error {
    // Проверить Redis
    if err := s.redisClient.Ping(ctx.Context()).Err(); err != nil {
        return ctx.Status(503).JSON(fiber.Map{
            "status": "unhealthy",
            "redis":  "down",
        })
    }

    return ctx.JSON(fiber.Map{
        "status": "healthy",
        "redis":  "up",
    })
}
```
