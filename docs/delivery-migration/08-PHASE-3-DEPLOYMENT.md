# Фаза 3: Развертывание и миграция монолита

**Срок**: Week 4


func (h *Handler) GetShipment(c *fiber.Ctx) error {
    id := c.Params("id")
    shipmentID, err := uuid.Parse(id)
    if err != nil {
        return utils.SendErrorResponse(c, fiber.StatusBadRequest, "error.invalid_shipment_id", nil)
    }

    shipment, err := h.client.grpc.GetShipment(c.Context(), shipmentID)
    if err != nil {
        return utils.SendErrorResponse(c, fiber.StatusNotFound, "error.shipment_not_found", nil)
    }

    return utils.SendSuccessResponse(c, shipment, "Информация об отправлении")
}

func (h *Handler) TrackShipment(c *fiber.Ctx) error {
    trackingNumber := c.Params("tracking")

    tracking, err := h.client.grpc.TrackShipment(c.Context(), trackingNumber)
    if err != nil {
        return utils.SendErrorResponse(c, fiber.StatusNotFound, "error.shipment_not_found", nil)
    }

    return utils.SendSuccessResponse(c, tracking, "Информация об отслеживании")
}

func (h *Handler) CalculateRate(c *fiber.Ctx) error {
    var req CalculateRateRequest
    if err := c.BodyParser(&req); err != nil {
        return utils.SendErrorResponse(c, fiber.StatusBadRequest, "error.invalid_request", nil)
    }

    rates, err := h.client.grpc.CalculateRate(c.Context(), &deliveryClient.CalculateRateRequest{
        FromAddress: req.FromAddress,
        ToAddress:   req.ToAddress,
        Package:     req.Package,
        Type:        req.Type,
    })

    if err != nil {
        return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "error.calculation_failed", nil)
    }

    return utils.SendSuccessResponse(c, rates, "Стоимость доставки рассчитана")
}

func (h *Handler) CancelShipment(c *fiber.Ctx) error {
    id := c.Params("id")
    shipmentID, err := uuid.Parse(id)
    if err != nil {
        return utils.SendErrorResponse(c, fiber.StatusBadRequest, "error.invalid_shipment_id", nil)
    }

    if err := h.client.grpc.CancelShipment(c.Context(), shipmentID); err != nil {
        return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "error.failed_to_cancel", nil)
    }

    return utils.SendSuccessResponse(c, nil, "Отправление отменено")
}
```

**Файл**: `backend/internal/proj/delivery/module.go`

```go
package delivery

import (
    "github.com/gofiber/fiber/v2"
    authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"
    "backend/internal/config"
    "backend/internal/middleware"
)

type Module struct {
    client  *Client
    handler *Handler
}

func NewModule(cfg *config.Config) (*Module, error) {
    // Подключение к микросервису delivery
    client, err := NewClient(cfg.DeliveryServiceAddress)
    if err != nil {
        return nil, fmt.Errorf("failed to create delivery client: %w", err)
    }

    handler := NewHandler(client)

    return &Module{
        client:  client,
        handler: handler,
    }, nil
}

func (m *Module) RegisterRoutes(app *fiber.App, mw *middleware.Middleware) error {
    // Защищенные роуты
    api := app.Group("/api/v1", mw.JWTParser(), authMiddleware.RequireAuth())

    delivery := api.Group("/delivery")
    delivery.Post("/calculate", m.handler.CalculateRate)

    shipments := api.Group("/shipments")
    shipments.Post("/", m.handler.CreateShipment)
    shipments.Get("/:id", m.handler.GetShipment)
    shipments.Get("/track/:tracking", m.handler.TrackShipment)
    shipments.Delete("/:id", m.handler.CancelShipment)

    return nil
}

func (m *Module) Close() error {
    return m.client.Close()
}
```

**Файл**: `backend/internal/config/config.go` (добавить)

```go
type Config struct {
    // ... существующие поля

    DeliveryServiceAddress string `env:"DELIVERY_SERVICE_ADDRESS" envDefault:"localhost:50052"`
}
```

#### 3.3 Обновление server.go

**Файл**: `backend/cmd/api/main.go`

```go
func main() {
    // ... существующая инициализация

    // Delivery module (теперь gRPC клиент)
    deliveryModule, err := delivery.NewModule(cfg)
    if err != nil {
        log.Fatal().Err(err).Msg("Failed to initialize delivery module")
    }
    defer deliveryModule.Close()

    if err := deliveryModule.RegisterRoutes(app, mw); err != nil {
        log.Fatal().Err(err).Msg("Failed to register delivery routes")
    }

    // ... остальное
}
```

#### 3.4 Миграция данных (если есть существующие shipments)

**Скрипт**: `backend/scripts/migrate_delivery_data.sql`

```sql
-- Подключаемся к обеим БД через dblink
CREATE EXTENSION IF NOT EXISTS dblink;

-- Копируем отправления из монолита в микросервис
INSERT INTO delivery_db.shipments (
    id,
    tracking_number,
    status,
    provider,
    user_id,
    from_street, from_city, from_state, from_postal_code, from_country, from_phone, from_email, from_name,
    to_street, to_city, to_state, to_postal_code, to_country, to_phone, to_email, to_name,
    weight_kg, length_cm, width_cm, height_cm, package_description, package_value,
    cost, currency,
    provider_shipment_id,
    provider_metadata,
    estimated_delivery_at,
    actual_delivery_at,
    created_at,
    updated_at
)
SELECT
    uuid_generate_v4(),  -- новый UUID
    tracking_number,
    status,
    provider_code,
    user_id,
    -- from address
    (sender_info->>'street')::text,
    (sender_info->>'city')::text,
    '',  -- state
    (sender_info->>'postal_code')::text,
    (sender_info->>'country')::text,
    (sender_info->>'phone')::text,
    (sender_info->>'email')::text,
    (sender_info->>'name')::text,
    -- to address
    (recipient_info->>'street')::text,
    (recipient_info->>'city')::text,
    '',  -- state
    (recipient_info->>'postal_code')::text,
    (recipient_info->>'country')::text,
    (recipient_info->>'phone')::text,
    (recipient_info->>'email')::text,
    (recipient_info->>'name')::text,
    -- package
    (package_info->>'weight_kg')::float,
    (package_info->>'length_cm')::float,
    (package_info->>'width_cm')::float,
    (package_info->>'height_cm')::float,
    (package_info->>'description')::text,
    (package_info->>'value')::float,
    -- cost
    delivery_cost,
    'RSD',
    external_id,
    provider_response,
    estimated_delivery,
    actual_delivery_date,
    created_at,
    updated_at
FROM svetubd.delivery_shipments;

-- Копируем tracking events
INSERT INTO delivery_db.tracking_events (
    id,
    shipment_id,
    status,
    location,
    details,
    timestamp,
    created_at
)
SELECT
    uuid_generate_v4(),
    -- найти новый shipment_id по tracking_number
    (SELECT id FROM delivery_db.shipments WHERE tracking_number = old_shipments.tracking_number),
    e.status,
    e.location,
    e.description,
    e.event_time,
    e.created_at
FROM svetubd.delivery_tracking_events e
JOIN svetubd.delivery_shipments old_shipments ON e.shipment_id = old_shipments.id;
```

**Запуск**:
```bash
psql "postgres://postgres:password@localhost:5432/delivery_db" -f backend/scripts/migrate_delivery_data.sql
```

#### 3.5 Deploy на dev

**Docker Compose**: `docker-compose.dev.yml` (обновить)

```yaml
version: '3.8'

services:
  # Новая БД для микросервиса delivery
  delivery-db:
    image: postgres:17-alpine
    container_name: delivery-db
    ports:
      - "5433:5432"
    environment:
      POSTGRES_DB: delivery_db
      POSTGRES_USER: delivery_user
      POSTGRES_PASSWORD: ${DELIVERY_DB_PASSWORD:-delivery_pass}
    volumes:
      - delivery_db_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U delivery_user -d delivery_db"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Микросервис delivery
  delivery-service:
    build:
      context: ./delivery
      dockerfile: Dockerfile
    container_name: delivery-service
    ports:
      - "50052:50052"  # gRPC
      - "9091:9091"    # Metrics
    environment:
      SVETUDELIVERY_DATABASE_HOST: delivery-db
      SVETUDELIVERY_DATABASE_PORT: 5432
      SVETUDELIVERY_DATABASE_NAME: delivery_db
      SVETUDELIVERY_DATABASE_USER: delivery_user
      SVETUDELIVERY_DATABASE_PASSWORD: ${DELIVERY_DB_PASSWORD:-delivery_pass}
      SVETUDELIVERY_GATEWAYS_POSTRS_ENABLED: "true"
      SVETUDELIVERY_GATEWAYS_POSTRS_API_KEY: ${POSTEXPRESS_API_KEY}
      SVETUDELIVERY_GATEWAYS_POSTRS_BASE_URL: "https://api.postexpress.rs"
    depends_on:
      delivery-db:
        condition: service_healthy
    restart: unless-stopped

  # Backend монолита (обновленный)
  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    container_name: backend
    ports:
      - "3000:3000"
    environment:
      # ... существующие env переменные
      DELIVERY_SERVICE_ADDRESS: "delivery-service:50052"
    depends_on:
      - delivery-service
    restart: unless-stopped

volumes:
  delivery_db_data:
```

**Запуск**:
```bash
# 1. Сборка и запуск
docker-compose -f docker-compose.dev.yml up -d --build

# 2. Проверка логов
docker-compose -f docker-compose.dev.yml logs -f delivery-service

# 3. Проверка через frontend
# Открыть https://dev.svetu.rs
# Создать тестовое объявление с доставкой
# Проверить что shipment создался и трекинг работает
```

#### 3.6 Удаление старых таблиц из монолита

**После** успешного тестирования на dev:

```sql
-- Подключаемся к БД монолита
psql "postgres://postgres:password@localhost:5432/svetubd"
