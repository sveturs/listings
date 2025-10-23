# 📖 Delivery Microservice - Руководство по использованию

**Версия**: 1.0
**Дата**: 2025-10-23
**Аудитория**: Backend разработчики, DevOps, QA

---

## 🎯 Содержание

1. [Быстрый старт](#быстрый-старт)
2. [Установка и настройка](#установка-и-настройка)
3. [Примеры использования API](#примеры-использования-api)
4. [Интеграция в Marketplace](#интеграция-в-marketplace)
5. [Работа с провайдерами](#работа-с-провайдерами)
6. [Troubleshooting](#troubleshooting)
7. [FAQ](#faq)

---

## 🚀 Быстрый старт

### Предварительные требования

- Go 1.23+
- PostgreSQL 17+ с PostGIS
- Redis 7+
- Docker & Docker Compose (опционально)
- grpcurl (для тестирования)

### Запуск через Docker Compose

```bash
# 1. Клонировать репозиторий
git clone https://github.com/sveturs/delivery.git
cd delivery

# 2. Настроить переменные окружения
cp .env.example .env
# Отредактировать .env файл

# 3. Запустить все сервисы
docker-compose up -d

# 4. Проверить статус
docker-compose ps
docker-compose logs delivery-service

# 5. Применить миграции (автоматически при старте)
# Или вручную:
docker-compose exec delivery-service ./migrator up

# 6. Проверить health
curl http://localhost:8081/health

# 7. Тестовый gRPC запрос
grpcurl -plaintext localhost:50052 list
```

### Запуск локально (для разработки)

```bash
# 1. Запустить зависимости
docker-compose up -d postgres redis

# 2. Применить миграции
export DATABASE_URL="postgres://delivery_user:password@localhost:5432/delivery_db?sslmode=disable"
./migrator up

# 3. Запустить сервис
go run cmd/api/main.go

# Сервис запущен на:
# - gRPC: localhost:50052
# - HTTP (metrics): localhost:8081
# - Prometheus metrics: localhost:9091
```

---

## ⚙️ Установка и настройка

### 1. Настройка PostgreSQL

```bash
# Создать базу данных
createdb delivery_db

# Создать пользователя
createuser -P delivery_user

# Дать права
psql delivery_db
GRANT ALL PRIVILEGES ON DATABASE delivery_db TO delivery_user;

# Установить PostGIS
CREATE EXTENSION postgis;
CREATE EXTENSION postgis_topology;
```

### 2. Настройка переменных окружения

Создать `.env` файл:

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=delivery_db
DB_USER=delivery_user
DB_PASSWORD=secure_password
DB_SSLMODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# gRPC Server
GRPC_PORT=50052

# HTTP Server
HTTP_PORT=8081
METRICS_PORT=9091

# Logging
LOG_LEVEL=info

# Provider API Keys
POST_EXPRESS_API_KEY=your_api_key_here
POST_EXPRESS_API_URL=https://api.postexpress.rs/v1
```

### 3. Применение миграций

```bash
# Установить migrate tool
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Применить миграции
migrate -path db/migrations -database "${DATABASE_URL}" up

# Или использовать встроенный migrator
./migrator up

# Откатить последнюю миграцию
./migrator down

# Проверить версию
./migrator version
```

### 4. Загрузка начальных данных

```bash
# Загрузить провайдеров доставки
psql $DATABASE_URL < db/fixtures/providers.sql

# Загрузить pricing rules
psql $DATABASE_URL < db/fixtures/pricing_rules.sql

# Загрузить зоны доставки
psql $DATABASE_URL < db/fixtures/zones.sql
```

---

## 📡 Примеры использования API

### Go Client

#### Установка

```bash
go get github.com/sveturs/delivery
```

#### Создание клиента

```go
package main

import (
    "context"
    "log"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"

    pb "github.com/sveturs/delivery/api/proto"
)

func main() {
    // Подключение к серверу
    conn, err := grpc.Dial("localhost:50052",
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    )
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer conn.Close()

    // Создание клиента
    client := pb.NewDeliveryServiceClient(conn)

    // Теперь можно вызывать методы
    ctx := context.Background()
    // ... использовать client
}
```

#### Пример 1: Расчет стоимости доставки

```go
func calculateDeliveryRate(client pb.DeliveryServiceClient) {
    ctx := context.Background()

    req := &pb.CalculateRateRequest{
        Provider: pb.DeliveryProvider_DELIVERY_PROVIDER_POST_EXPRESS,
        FromAddress: &pb.Address{
            Street:       "Kneza Milosa 10",
            City:         "Belgrade",
            State:        "Belgrade",
            PostalCode:   "11000",
            Country:      "RS",
            ContactName:  "John Doe",
            ContactPhone: "+381611234567",
        },
        ToAddress: &pb.Address{
            Street:       "Bulevar Oslobodjenja 1",
            City:         "Novi Sad",
            State:        "Vojvodina",
            PostalCode:   "21000",
            Country:      "RS",
            ContactName:  "Jane Smith",
            ContactPhone: "+381621234567",
        },
        Package: &pb.Package{
            Weight:      "1.5",  // kg
            Length:      "40",   // cm
            Width:       "30",   // cm
            Height:      "20",   // cm
            Description: "Books and magazines",
        },
        DeliveryType: pb.DeliveryType_DELIVERY_TYPE_STANDARD,
    }

    resp, err := client.CalculateRate(ctx, req)
    if err != nil {
        log.Fatalf("CalculateRate failed: %v", err)
    }

    log.Printf("Cost: %s %s", resp.Cost, resp.Currency)
    log.Printf("Estimated delivery: %s", resp.EstimatedDelivery)

    // Output:
    // Cost: 450.00 RSD
    // Estimated delivery: 2025-10-28T12:00:00Z
}
```

#### Пример 2: Создание отправки

```go
func createShipment(client pb.DeliveryServiceClient, orderID int32) string {
    ctx := context.Background()

    req := &pb.CreateShipmentRequest{
        Provider: pb.DeliveryProvider_DELIVERY_PROVIDER_POST_EXPRESS,
        OrderId:  orderID,
        FromAddress: &pb.Address{
            Street:       "Kneza Milosa 10",
            City:         "Belgrade",
            PostalCode:   "11000",
            Country:      "RS",
            ContactName:  "Shop Owner",
            ContactPhone: "+381611234567",
            ContactEmail: "shop@example.com",
        },
        ToAddress: &pb.Address{
            Street:       "Bulevar Oslobodjenja 1",
            City:         "Novi Sad",
            PostalCode:   "21000",
            Country:      "RS",
            ContactName:  "Customer Name",
            ContactPhone: "+381621234567",
            ContactEmail: "customer@example.com",
        },
        Package: &pb.Package{
            Weight:        "2.0",
            Length:        "50",
            Width:         "40",
            Height:        "30",
            Description:   "Order #12345 - Electronics",
            DeclaredValue: "15000", // RSD
        },
        DeliveryType:   pb.DeliveryType_DELIVERY_TYPE_EXPRESS,
        InsuranceValue: "15000",
        UserId:         "550e8400-e29b-41d4-a716-446655440000",
        Reference:      "ORDER-12345",
        Notes:          "Handle with care - fragile items",
    }

    resp, err := client.CreateShipment(ctx, req)
    if err != nil {
        log.Fatalf("CreateShipment failed: %v", err)
    }

    log.Printf("Shipment created!")
    log.Printf("  ID: %s", resp.Shipment.Id)
    log.Printf("  Tracking: %s", resp.Shipment.TrackingNumber)
    log.Printf("  Status: %s", resp.Shipment.Status)
    log.Printf("  Cost: %s %s", resp.Shipment.Cost, resp.Shipment.Currency)

    // Сохранить tracking URL для клиента
    log.Printf("  Tracking URL: %s", resp.TrackingUrl)

    return resp.Shipment.TrackingNumber
}
```

#### Пример 3: Отслеживание отправки

```go
func trackShipment(client pb.DeliveryServiceClient, trackingNumber string) {
    ctx := context.Background()

    req := &pb.TrackShipmentRequest{
        TrackingNumber: trackingNumber,
        ForceSync:      true, // Принудительная синхронизация с провайдером
    }

    resp, err := client.TrackShipment(ctx, req)
    if err != nil {
        log.Fatalf("TrackShipment failed: %v", err)
    }

    log.Printf("Current status: %s", resp.Shipment.Status)
    log.Printf("Last updated: %s", resp.Shipment.UpdatedAt)

    log.Println("\nTracking history:")
    for i, event := range resp.Events {
        log.Printf("%d. %s - %s", i+1, event.EventTime, event.Description)
        if event.Location != "" {
            log.Printf("   Location: %s", event.Location)
        }
    }

    // Output example:
    // Current status: SHIPMENT_STATUS_OUT_FOR_DELIVERY
    // Last updated: 2025-10-23T14:30:00Z
    //
    // Tracking history:
    // 1. 2025-10-22T10:00:00Z - Заказ создан
    //    Location: Белград
    // 2. 2025-10-22T12:00:00Z - Посылка забрана курьером
    //    Location: Белград, Склад
    // 3. 2025-10-23T08:00:00Z - Посылка в пути
    //    Location: Нови-Сад
    // 4. 2025-10-23T14:00:00Z - На доставке
    //    Location: Нови-Сад, Центр
}
```

#### Пример 4: Отмена отправки

```go
func cancelShipment(client pb.DeliveryServiceClient, shipmentID string) {
    ctx := context.Background()

    req := &pb.CancelShipmentRequest{
        Id:     shipmentID,
        Reason: "Customer requested cancellation",
    }

    resp, err := client.CancelShipment(ctx, req)
    if err != nil {
        log.Fatalf("CancelShipment failed: %v", err)
    }

    log.Printf("Shipment cancelled successfully")
    log.Printf("  New status: %s", resp.Shipment.Status)

    if resp.RefundEligible {
        log.Printf("  Refund eligible: YES")
        log.Printf("  Refund amount: %s %s", resp.RefundAmount, resp.Shipment.Currency)
    } else {
        log.Printf("  Refund eligible: NO")
    }
}
```

#### Пример 5: Получение списка провайдеров

```go
func listProviders(client pb.DeliveryServiceClient) {
    ctx := context.Background()

    req := &pb.ListProvidersRequest{
        ActiveOnly: true,
        Country:    "RS",
    }

    resp, err := client.ListProviders(ctx, req)
    if err != nil {
        log.Fatalf("ListProviders failed: %v", err)
    }

    log.Printf("Available providers:")
    for _, provider := range resp.Providers {
        log.Printf("\n%s (%s)", provider.Name, provider.Code)
        log.Printf("  COD: %v, Insurance: %v, Tracking: %v",
            provider.SupportsCod,
            provider.SupportsInsurance,
            provider.SupportsTracking,
        )
        log.Printf("  Countries: %v", provider.Countries)
        log.Printf("  Delivery types: %v", provider.DeliveryTypes)
    }
}
```

---

### cURL Examples (для тестирования)

#### Через grpcurl

```bash
# 1. Список методов
grpcurl -plaintext localhost:50052 list

# 2. Описание метода
grpcurl -plaintext localhost:50052 describe delivery.v1.DeliveryService.CalculateRate

# 3. CalculateRate
grpcurl -plaintext localhost:50052 delivery.v1.DeliveryService/CalculateRate -d '{
  "provider": "DELIVERY_PROVIDER_POST_EXPRESS",
  "from_address": {
    "street": "Kneza Milosa 10",
    "city": "Belgrade",
    "postal_code": "11000",
    "country": "RS",
    "contact_name": "John Doe",
    "contact_phone": "+381611234567"
  },
  "to_address": {
    "street": "Bulevar Oslobodjenja 1",
    "city": "Novi Sad",
    "postal_code": "21000",
    "country": "RS",
    "contact_name": "Jane Smith",
    "contact_phone": "+381621234567"
  },
  "package": {
    "weight": "1.0",
    "length": "30",
    "width": "20",
    "height": "10",
    "description": "Test package"
  },
  "delivery_type": "DELIVERY_TYPE_STANDARD"
}'

# 4. CreateShipment
grpcurl -plaintext localhost:50052 delivery.v1.DeliveryService/CreateShipment -d '{
  "provider": "DELIVERY_PROVIDER_POST_EXPRESS",
  "order_id": 12345,
  "from_address": {...},
  "to_address": {...},
  "package": {...},
  "user_id": "550e8400-e29b-41d4-a716-446655440000"
}'

# 5. GetShipment
grpcurl -plaintext localhost:50052 delivery.v1.DeliveryService/GetShipment -d '{
  "id": "5"
}'

# 6. TrackShipment
grpcurl -plaintext localhost:50052 delivery.v1.DeliveryService/TrackShipment -d '{
  "tracking_number": "post_express-1761215005-6768",
  "force_sync": true
}'

# 7. CancelShipment
grpcurl -plaintext localhost:50052 delivery.v1.DeliveryService/CancelShipment -d '{
  "id": "5",
  "reason": "Testing cancellation"
}'
```

---

## 🔌 Интеграция в Marketplace

### Архитектура интеграции

```
User → Frontend → BFF → Marketplace Backend → gRPC Client → Delivery Microservice
```

### 1. Создание gRPC клиента в Marketplace

```go
// backend/internal/delivery/client.go
package delivery

import (
    "context"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"

    pb "github.com/sveturs/delivery/api/proto"
)

type Client struct {
    conn   *grpc.ClientConn
    client pb.DeliveryServiceClient
}

func NewClient(address string) (*Client, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    conn, err := grpc.DialContext(ctx, address,
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithBlock(),
    )
    if err != nil {
        return nil, err
    }

    return &Client{
        conn:   conn,
        client: pb.NewDeliveryServiceClient(conn),
    }, nil
}

func (c *Client) Close() error {
    return c.conn.Close()
}

func (c *Client) CalculateRate(ctx context.Context, req *pb.CalculateRateRequest) (*pb.CalculateRateResponse, error) {
    return c.client.CalculateRate(ctx, req)
}

func (c *Client) CreateShipment(ctx context.Context, req *pb.CreateShipmentRequest) (*pb.CreateShipmentResponse, error) {
    return c.client.CreateShipment(ctx, req)
}

// ... остальные методы
```

### 2. Добавление в Checkout Flow

```go
// backend/internal/proj/marketplace/handler/checkout.go
func (h *Handler) Checkout(c *fiber.Ctx) error {
    var req CheckoutRequest
    if err := c.BodyParser(&req); err != nil {
        return err
    }

    // ... валидация корзины, создание заказа

    // Расчет доставки
    if req.DeliveryRequired {
        deliveryReq := &pb.CalculateRateRequest{
            Provider:     req.DeliveryProvider,
            FromAddress:  convertToProtoAddress(order.SellerAddress),
            ToAddress:    convertToProtoAddress(req.DeliveryAddress),
            Package:      buildPackageFromOrder(order),
            DeliveryType: req.DeliveryType,
        }

        deliveryResp, err := h.deliveryClient.CalculateRate(c.Context(), deliveryReq)
        if err != nil {
            return fiber.NewError(fiber.StatusBadRequest, "delivery.calculation_failed")
        }

        // Добавить стоимость доставки к заказу
        deliveryCost, _ := decimal.NewFromString(deliveryResp.Cost)
        order.DeliveryCost = deliveryCost
        order.TotalAmount = order.TotalAmount.Add(deliveryCost)
    }

    // Создать заказ
    if err := h.orderRepo.Create(c.Context(), order); err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "order_id": order.ID,
        "total":    order.TotalAmount,
        "delivery": fiber.Map{
            "cost":               order.DeliveryCost,
            "estimated_delivery": deliveryResp.EstimatedDelivery,
        },
    })
}
```

### 3. Создание отправки после оплаты

```go
// backend/internal/proj/marketplace/handler/order.go
func (h *Handler) ProcessPayment(c *fiber.Ctx) error {
    orderID := c.Params("id")

    // ... обработка платежа

    if paymentSuccess {
        // Создать отправку
        shipmentReq := &pb.CreateShipmentRequest{
            Provider:    order.DeliveryProvider,
            OrderId:     int32(order.ID),
            FromAddress: convertToProtoAddress(order.SellerAddress),
            ToAddress:   convertToProtoAddress(order.DeliveryAddress),
            Package:     buildPackageFromOrder(order),
            UserId:      order.BuyerID.String(),
            Reference:   fmt.Sprintf("ORDER-%d", order.ID),
        }

        shipmentResp, err := h.deliveryClient.CreateShipment(c.Context(), shipmentReq)
        if err != nil {
            // Логировать ошибку, но не отменять заказ
            log.Error().Err(err).Msg("Failed to create shipment")
        } else {
            // Сохранить tracking number
            order.TrackingNumber = shipmentResp.Shipment.TrackingNumber
            h.orderRepo.Update(c.Context(), order)
        }
    }

    return c.JSON(fiber.Map{"status": "success"})
}
```

### 4. Отображение tracking в Order Details

```go
func (h *Handler) GetOrderDetails(c *fiber.Ctx) error {
    orderID := c.Params("id")
    order, err := h.orderRepo.GetByID(c.Context(), orderID)
    if err != nil {
        return err
    }

    var trackingInfo *pb.TrackShipmentResponse
    if order.TrackingNumber != "" {
        trackingReq := &pb.TrackShipmentRequest{
            TrackingNumber: order.TrackingNumber,
            ForceSync:      false, // Использовать кэш
        }

        trackingInfo, err = h.deliveryClient.TrackShipment(c.Context(), trackingReq)
        if err != nil {
            log.Error().Err(err).Msg("Failed to get tracking info")
        }
    }

    return c.JSON(fiber.Map{
        "order":    order,
        "tracking": trackingInfo,
    })
}
```

---

## 🚚 Работа с провайдерами

### Добавление нового провайдера

#### 1. Реализовать Provider Interface

```go
// internal/provider/new_provider.go
package provider

import (
    "context"
    "github.com/sveturs/delivery/internal/domain"
)

type NewProvider struct {
    apiKey    string
    apiURL    string
    httpClient *http.Client
}

func NewNewProvider(apiKey, apiURL string) *NewProvider {
    return &NewProvider{
        apiKey: apiKey,
        apiURL: apiURL,
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

func (p *NewProvider) CreateShipment(ctx context.Context, req *domain.ShipmentRequest) (*domain.ShipmentResponse, error) {
    // Реализация создания отправки через API провайдера
    // ...
}

func (p *NewProvider) TrackShipment(ctx context.Context, trackingNumber string) (*domain.TrackingResponse, error) {
    // Реализация отслеживания
    // ...
}

func (p *NewProvider) CancelShipment(ctx context.Context, externalID string) error {
    // Реализация отмены
    // ...
}

func (p *NewProvider) HandleWebhook(ctx context.Context, payload []byte, headers map[string]string) (*domain.WebhookResponse, error) {
    // Реализация webhook handler
    // ...
}
```

#### 2. Зарегистрировать в Factory

```go
// internal/provider/factory.go
func (f *Factory) CreateProvider(code string) (DeliveryProvider, error) {
    switch code {
    case "post_express":
        return NewPostExpressProvider(f.config.PostExpress), nil
    case "new_provider":  // Добавить новый провайдер
        return NewNewProvider(f.config.NewProvider.APIKey, f.config.NewProvider.APIURL), nil
    default:
        return nil, fmt.Errorf("unknown provider: %s", code)
    }
}
```

#### 3. Добавить в БД

```sql
INSERT INTO delivery_providers (code, name, logo_url, is_active, supports_cod, supports_insurance, supports_tracking, api_config)
VALUES (
    'new_provider',
    'New Provider Name',
    'https://example.com/logo.png',
    true,
    true,
    true,
    true,
    '{"api_url": "https://api.newprovider.com/v1", "timeout_seconds": 30}'::jsonb
);
```

#### 4. Добавить Pricing Rules

```sql
INSERT INTO delivery_pricing_rules (provider_id, rule_type, weight_ranges, is_active)
VALUES (
    (SELECT id FROM delivery_providers WHERE code = 'new_provider'),
    'weight_based',
    '[
        {"from": 0, "to": 1, "base_price": 250, "price_per_kg": 0},
        {"from": 1, "to": 5, "base_price": 400, "price_per_kg": 50},
        {"from": 5, "to": 30, "base_price": 600, "price_per_kg": 100}
    ]'::jsonb,
    true
);
```

---

## 🔧 Troubleshooting

### Проблема 1: gRPC connection refused

**Симптомы**:
```
Error: rpc error: code = Unavailable desc = connection error
```

**Решение**:
```bash
# Проверить что сервис запущен
docker-compose ps delivery-service

# Проверить логи
docker-compose logs delivery-service

# Проверить порт
netstat -tlnp | grep 50052

# Перезапустить сервис
docker-compose restart delivery-service
```

### Проблема 2: Database connection failed

**Симптомы**:
```
Error: pq: password authentication failed for user "delivery_user"
```

**Решение**:
```bash
# Проверить переменные окружения
docker-compose exec delivery-service env | grep DB_

# Проверить подключение вручную
psql "postgres://delivery_user:password@localhost:5432/delivery_db"

# Пересоздать контейнеры
docker-compose down -v
docker-compose up -d
```

### Проблема 3: Provider API timeout

**Симптомы**:
```
Error: provider API request timeout
```

**Решение**:
```bash
# Увеличить timeout в конфигурации
# .env
POST_EXPRESS_TIMEOUT=60s

# Проверить доступность API провайдера
curl -v https://api.postexpress.rs/v1/health

# Использовать mock provider для тестирования
MOCK_PROVIDER=true
```

### Проблема 4: JSONB marshaling error

**Симптомы**:
```
Error: pq: invalid input syntax for type json
```

**Решение**:
Эта проблема была исправлена в commit 4cc0b7d. Убедитесь что используете последнюю версию кода.

---

## ❓ FAQ

### Q: Как добавить новую страну доставки?

**A**: Добавить страну в capabilities провайдера:

```sql
UPDATE delivery_providers
SET capabilities = jsonb_set(
    capabilities,
    '{countries}',
    capabilities->'countries' || '["BA"]'::jsonb
)
WHERE code = 'post_express';
```

### Q: Как кэшировать результаты расчета?

**A**: Использовать Redis с TTL:

```go
// Calculate rate
cacheKey := fmt.Sprintf("rate:%s:%s:%s", provider, fromCity, toCity)
cached, err := redis.Get(ctx, cacheKey).Result()
if err == nil {
    return cached
}

// Call provider API
result := calculateRate(...)

// Cache for 5 minutes
redis.Set(ctx, cacheKey, result, 5*time.Minute)
return result
```

### Q: Как обрабатывать webhook от провайдеров?

**A**: Настроить webhook endpoint в marketplace:

```go
// backend/internal/proj/marketplace/handler/webhook.go
func (h *Handler) HandleDeliveryWebhook(c *fiber.Ctx) error {
    provider := c.Params("provider")

    webhookReq := &pb.ProcessWebhookRequest{
        Provider:  convertProviderCode(provider),
        Payload:   c.Body(),
        Headers:   convertHeaders(c.GetReqHeaders()),
        Signature: c.Get("X-Signature"),
    }

    resp, err := h.deliveryClient.ProcessWebhook(c.Context(), webhookReq)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{"status": "ok", "updated": resp.UpdatedShipments})
}
```

### Q: Как масштабировать сервис?

**A**:
1. Запустить несколько инстансов за load balancer
2. Использовать Redis для кэширования
3. Настроить connection pooling в БД
4. Добавить rate limiting per provider

```yaml
# docker-compose.scale.yml
version: '3.8'
services:
  delivery-service:
    image: delivery:latest
    deploy:
      replicas: 3
    depends_on:
      - postgres
      - redis
```

---

## 📞 Поддержка

**Документация**: `/data/hostel-booking-system/docs/`
**GitHub**: https://github.com/sveturs/delivery
**Issues**: https://github.com/sveturs/delivery/issues
