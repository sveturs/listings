# План миграции Delivery Service из монолита в микросервис

**Дата создания**: 2025-10-22
**Статус**: План разработки
**Версия**: 1.0

---

## 🎯 Цель миграции

Перенести функциональность доставки из монолита в отдельный gRPC микросервис для:
- Независимого масштабирования
- Изоляции логики работы с провайдерами доставки
- Переиспользования в других проектах
- Упрощения интеграции новых провайдеров

## 📊 Текущее состояние

### Монолит (backend/internal/proj/delivery)
**Реализовано**: ~2500 строк Go кода

**Компоненты**:
- ✅ Универсальный интерфейс DeliveryProvider
- ✅ Factory pattern для создания провайдеров
- ✅ Адаптер Post Express (полная интеграция)
- ✅ Mock провайдеры (bex_express, aks_express, dhl_express, etc.)
- ✅ Service layer (расчет стоимости, создание отправлений, трекинг)
- ✅ Storage layer (PostgreSQL)
- ✅ Handlers (REST API endpoints)
- ✅ Calculator (расчет стоимости с оптимизацией упаковки)
- ✅ Attributes service (управление атрибутами доставки товаров)
- ✅ Notifications integration (уведомления об изменении статуса)
- ✅ Admin functionality (управление провайдерами, аналитика)

### Микросервис (github.com/sveturs/delivery)
**Готовность**: ~35%

**Есть**:
- ✅ Proto API определения (gRPC)
- ✅ Database connection + migrations
- ✅ Config management (env-based)
- ✅ Logging infrastructure
- ✅ Makefile (build, lint, test, proto)
- ✅ Docker Compose (PostgreSQL)
- ✅ Схема БД (shipments, tracking_events)

**Отсутствует**:
- ❌ Domain models
- ❌ Service layer
- ❌ Repository layer
- ❌ Gateway integrations (Dex, Post RS)
- ❌ gRPC handlers implementation
- ❌ Tests (unit, integration)
- ❌ Metrics collection
- ❌ Provider factory
- ❌ Generated proto code

---

## 🏗️ Архитектура микросервиса

### Слои приложения

```
┌─────────────────────────────────────────────┐
│           gRPC Transport Layer              │
│  (internal/server/grpc/delivery.go)         │
└──────────────────┬──────────────────────────┘
                   │
┌──────────────────▼──────────────────────────┐
│           Service Layer                     │
│  (internal/service/delivery_service.go)     │
│  - Business logic                           │
│  - Orchestration                            │
└──────────────────┬──────────────────────────┘
                   │
        ┌──────────┴──────────┐
        │                     │
┌───────▼────────┐   ┌────────▼─────────┐
│  Repository    │   │   Gateway        │
│  (PostgreSQL)  │   │   (Providers)    │
│                │   │                  │
│ - Shipments    │   │ - Post Express   │
│ - Events       │   │ - Dex            │
│ - Providers    │   │ - Mock           │
└────────────────┘   └──────────────────┘
```

### Provider Pattern

```go
// Провайдеры доставки реализуют единый интерфейс
type DeliveryProvider interface {
    GetCode() string
    CalculateRate(ctx, *RateRequest) (*RateResponse, error)
    CreateShipment(ctx, *ShipmentRequest) (*ShipmentResponse, error)
    TrackShipment(ctx, trackingNumber) (*TrackingResponse, error)
    CancelShipment(ctx, shipmentID) error
    ValidateAddress(ctx, *Address) (*ValidationResponse, error)
}

// Фабрика создает провайдеров по коду
type ProviderFactory struct {
    providers map[string]DeliveryProvider
}

// Адаптеры конвертируют между универсальным API и специфичным API провайдера
type PostExpressAdapter struct {
    client *postexpress.Client
}
```

---

## 📋 План миграции (поэтапный)

### Фаза 1: Подготовка микросервиса (Неделя 1)

#### 1.1 Генерация proto кода
```bash
cd ~/delivery
make proto  # Генерация Go кода из proto/delivery/v1/delivery.proto
```

**Результат**: Директория `gen/go/delivery/v1/` с gRPC клиентом и сервером

#### 1.2 Создание domain models

**Файл**: `internal/domain/models.go`

```go
package domain

import (
    "time"
    pb "github.com/sveturs/delivery/gen/go/delivery/v1"
)

type Shipment struct {
    ID                uuid.UUID
    TrackingNumber    string
    Status            ShipmentStatus
    Provider          DeliveryProvider
    FromAddress       Address
    ToAddress         Address
    Package           Package
    Cost              Money
    ProviderShipmentID string
    ProviderMetadata   json.RawMessage
    EstimatedDelivery time.Time
    ActualDelivery    *time.Time
    CreatedAt         time.Time
    UpdatedAt         time.Time
}

type Address struct { /* ... */ }
type Package struct { /* ... */ }
type TrackingEvent struct { /* ... */ }
```

**Маппинг**: Функции конвертации domain ↔ protobuf

#### 1.3 Реализация Repository layer

**Файл**: `internal/repository/postgres/shipment_repository.go`

```go
package postgres

type ShipmentRepository struct {
    db *sql.DB
}

func (r *ShipmentRepository) Create(ctx, *domain.Shipment) error
func (r *ShipmentRepository) GetByID(ctx, uuid.UUID) (*domain.Shipment, error)
func (r *ShipmentRepository) GetByTracking(ctx, string) (*domain.Shipment, error)
func (r *ShipmentRepository) Update(ctx, *domain.Shipment) error
func (r *ShipmentRepository) List(ctx, *ListFilter) ([]*domain.Shipment, error)
```

**Миграция кода**: Адаптировать из `backend/internal/proj/delivery/storage/storage.go`

#### 1.4 Создание Provider Factory

**Файл**: `internal/gateway/provider/factory.go`

```go
package provider

type Factory struct {
    providers map[string]DeliveryProvider
}

func NewFactory() *Factory {
    f := &Factory{providers: make(map[string]DeliveryProvider)}

    // Регистрация провайдеров
    f.Register("post_express", NewPostExpressProvider())
    f.Register("dex", NewDexProvider())
    f.Register("mock", NewMockProvider())

    return f
}

func (f *Factory) GetProvider(code string) (DeliveryProvider, error)
```

#### 1.5 Интеграция Post Express

**Файл**: `internal/gateway/provider/postexpress/adapter.go`

**Задачи**:
1. Скопировать `backend/internal/proj/postexpress/` → `internal/gateway/provider/postexpress/`
2. Адаптировать под новые интерфейсы
3. Сохранить всю B2B логику (маппинг полей, валидация, расчет стоимости)

**Конфигурация**:
```env
SVETUDELIVERY_GATEWAYS_POSTRS_ENABLED=true
SVETUDELIVERY_GATEWAYS_POSTRS_API_KEY=xxx
SVETUDELIVERY_GATEWAYS_POSTRS_BASE_URL=https://api.postexpress.rs
```

### Фаза 2: Реализация Service Layer (Неделя 2)

#### 2.1 Delivery Service

**Файл**: `internal/service/delivery_service.go`

```go
package service

type DeliveryService struct {
    repo    repository.ShipmentRepository
    factory *provider.Factory
    logger  *logger.Logger
}

func (s *DeliveryService) CreateShipment(ctx, *CreateShipmentInput) (*domain.Shipment, error) {
    // 1. Валидация входных данных
    // 2. Выбор провайдера через factory
    // 3. Создание отправления через provider
    // 4. Сохранение в БД через repository
    // 5. Возврат результата
}

func (s *DeliveryService) GetShipment(ctx, uuid.UUID) (*domain.Shipment, error)
func (s *DeliveryService) TrackShipment(ctx, string) (*TrackingInfo, error)
func (s *DeliveryService) CalculateRate(ctx, *RateRequest) (*RateResponse, error)
func (s *DeliveryService) CancelShipment(ctx, uuid.UUID) error
```

**Миграция логики**: Из `backend/internal/proj/delivery/service/service.go`

#### 2.2 Rate Calculator Service

**Файл**: `internal/service/calculator_service.go`

```go
type CalculatorService struct {
    factory *provider.Factory
}

func (s *CalculatorService) CalculateForMultipleProviders(ctx, *RateRequest) ([]RateResponse, error) {
    // Запрашивает расчет у всех активных провайдеров параллельно
    // Возвращает сортированный список опций (по цене, по скорости)
}

func (s *CalculatorService) OptimizePackaging(items []Item) []Package {
    // Оптимизация упаковки товаров в посылки (bin packing)
}
```

**Миграция**: Из `backend/internal/proj/delivery/calculator/service.go`

#### 2.3 gRPC Handlers

**Файл**: `internal/server/grpc/delivery.go`

```go
func (s *DeliveryServer) CreateShipment(ctx, *pb.CreateShipmentRequest) (*pb.CreateShipmentResponse, error) {
    // 1. Парсинг и валидация protobuf
    // 2. Конвертация pb → domain
    // 3. Вызов service.CreateShipment()
    // 4. Конвертация domain → pb
    // 5. Возврат ответа
}

// Аналогично для остальных методов
```

### Фаза 3: Интеграция с монолитом (Неделя 3)

#### 3.1 Создание gRPC клиента в монолите

**Файл**: `backend/pkg/delivery/client.go`

```go
package delivery

import (
    pb "github.com/sveturs/delivery/gen/go/delivery/v1"
    "google.golang.org/grpc"
)

type Client struct {
    conn   *grpc.ClientConn
    client pb.DeliveryServiceClient
}

func NewClient(addr string) (*Client, error) {
    conn, err := grpc.Dial(addr, grpc.WithInsecure())
    if err != nil {
        return nil, err
    }

    return &Client{
        conn:   conn,
        client: pb.NewDeliveryServiceClient(conn),
    }, nil
}

func (c *Client) CreateShipment(ctx, *CreateShipmentRequest) (*Shipment, error) {
    // Маппинг request → protobuf
    resp, err := c.client.CreateShipment(ctx, pbReq)
    // Маппинг protobuf → response
    return shipment, nil
}
```

#### 3.2 Создание адаптера в монолите

**Файл**: `backend/internal/proj/delivery/client/adapter.go`

```go
package client

import deliveryClient "backend/pkg/delivery"

// Адаптер сохраняет текущий интерфейс, но делегирует вызовы микросервису
type MicroserviceAdapter struct {
    client *deliveryClient.Client
}

func (a *MicroserviceAdapter) CreateShipment(ctx, *interfaces.ShipmentRequest) (*interfaces.ShipmentResponse, error) {
    // Конвертация interfaces → client types
    resp, err := a.client.CreateShipment(ctx, req)
    // Конвертация обратно
    return response, err
}
```

#### 3.3 Feature flag для переключения

**Конфигурация**:
```env
# Если true - используется микросервис, иначе локальная реализация
DELIVERY_USE_MICROSERVICE=false
DELIVERY_GRPC_ADDRESS=localhost:50052
```

**Инициализация**:
```go
var deliveryService interfaces.DeliveryProvider

if config.UseDeliveryMicroservice {
    client := deliveryClient.NewClient(config.DeliveryGRPCAddress)
    deliveryService = client.NewAdapter(client)
} else {
    deliveryService = factory.NewProviderFactory(db)
}
```

### Фаза 4: Тестирование и развертывание (Неделя 4)

#### 4.1 Unit тесты микросервиса

**Файлы**:
- `internal/service/delivery_service_test.go` - тесты сервиса
- `internal/repository/postgres/shipment_repository_test.go` - тесты репозитория
- `internal/gateway/provider/postexpress/adapter_test.go` - тесты адаптера

**Запуск**:
```bash
make test-unit
make test-integration  # Использует testcontainers для PostgreSQL
```

#### 4.2 Integration тесты

**Файл**: `tests/integration/delivery_flow_test.go`

**Сценарии**:
1. ✅ Создание отправления → Получение tracking number → Трекинг → Delivery
2. ✅ Расчет стоимости для нескольких провайдеров
3. ✅ Отмена отправления
4. ✅ Webhook от провайдера → Обновление статуса
5. ✅ Валидация адреса

#### 4.3 E2E тесты монолит → микросервис

**Файл**: `backend/tests/delivery_microservice_integration_test.go`

**Сценарии**:
1. Монолит вызывает CreateShipment через gRPC
2. Микросервис создает отправление через Post Express
3. Монолит получает tracking number
4. Монолит вызывает TrackShipment
5. Проверка корректности данных

#### 4.4 Развертывание на dev

**Docker Compose**: `docker-compose.dev.yml`

```yaml
services:
  delivery-db:
    image: postgres:17-alpine
    ports: ["5433:5432"]
    environment:
      POSTGRES_DB: delivery_db
      POSTGRES_USER: delivery_user
      POSTGRES_PASSWORD: ${DELIVERY_DB_PASSWORD}

  delivery-service:
    build: ./delivery
    ports: ["50052:50052", "9091:9091"]
    environment:
      SVETUDELIVERY_DATABASE_HOST: delivery-db
      SVETUDELIVERY_GATEWAYS_POSTRS_API_KEY: ${POSTEXPRESS_API_KEY}
    depends_on:
      - delivery-db

  backend:
    build: ./backend
    environment:
      DELIVERY_USE_MICROSERVICE: "true"
      DELIVERY_GRPC_ADDRESS: "delivery-service:50052"
```

**Миграция данных**:
```sql
-- Копирование существующих отправлений из монолита в микросервис
INSERT INTO delivery_db.shipments (...)
SELECT ... FROM svetubd.delivery_shipments;
```

#### 4.5 Постепенный переход

**Week 1**: Feature flag = false (используется монолит), микросервис в standby
**Week 2**: Feature flag = true для 10% пользователей (canary deployment)
**Week 3**: Feature flag = true для 50% пользователей
**Week 4**: Feature flag = true для 100% пользователей
**Week 5**: Удаление старой реализации из монолита

---

## 🔄 Миграция Post Express интеграции

### Что переносится

**1. Client HTTP** (`backend/internal/proj/postexpress/client.go`)
→ `internal/gateway/provider/postexpress/client.go`

**2. Service** (`backend/internal/proj/postexpress/service/service.go`)
→ `internal/gateway/provider/postexpress/service.go`

**3. Types & Models** (`backend/internal/proj/postexpress/types.go`, `models/models.go`)
→ `internal/gateway/provider/postexpress/types.go`

**4. Config** (`backend/internal/proj/postexpress/config.go`)
→ Интеграция в `internal/config/config.go` (секция Gateways.PostRS)

### Адаптация

**Было**:
```go
// В монолите
postExpressService := postexpress.NewService(config)
adapter := factory.NewPostExpressAdapter(postExpressService)
shipment, err := adapter.CreateShipment(ctx, req)
```

**Стало**:
```go
// В микросервисе
postExpressProvider := postexpress.NewProvider(cfg.Gateways.PostRS)
shipment, err := postExpressProvider.CreateShipment(ctx, req)
```

### Сохранение функциональности

✅ **Все B2B поля** - ExtBrend, ExtMagacin, NacinPrijema, etc.
✅ **Валидация** - ValidateShipment перед отправкой
✅ **Маппинг статусов** - mapPostExpressStatus()
✅ **Расчет стоимости** - CalculateRate с весовыми диапазонами
✅ **COD логика** - Otkupnina, банковские реквизиты
✅ **Tracking** - GetTrackingInfo, события
✅ **Webhooks** - Обработка callback от Post Express

---

## 🗂️ Провайдер поставщиков услуг доставки (Provider Pattern)

### Архитектура

```go
// Общий интерфейс для всех провайдеров
type DeliveryProvider interface {
    GetCode() string
    GetName() string
    IsAvailable() bool
    GetCapabilities() *Capabilities

    CalculateRate(ctx, *RateRequest) (*RateResponse, error)
    CreateShipment(ctx, *ShipmentRequest) (*ShipmentResponse, error)
    TrackShipment(ctx, trackingNumber) (*TrackingResponse, error)
    CancelShipment(ctx, shipmentID) error
    GetLabel(ctx, shipmentID) (*Label, error)
    ValidateAddress(ctx, *Address) (*AddressValidation, error)
    HandleWebhook(ctx, payload, headers) (*WebhookResult, error)
}

// Capabilities описывает возможности провайдера
type Capabilities struct {
    MaxWeightKg       float64
    MaxVolumeM3       float64
    MaxDimensions     Dimensions
    SupportedZones    []string  // local, regional, national, international
    SupportedTypes    []string  // standard, express, same_day
    SupportsCOD       bool
    SupportsInsurance bool
    SupportsTracking  bool
    SupportsReturn    bool
    AdditionalServices []string  // signature, photo_proof, weekend_delivery
}
```

### Реестр провайдеров

**Файл**: `internal/gateway/provider/registry.go`

```go
type Registry struct {
    providers map[string]ProviderFactory
}

type ProviderFactory func(config ProviderConfig) (DeliveryProvider, error)

func NewRegistry() *Registry {
    r := &Registry{providers: make(map[string]ProviderFactory)}

    // Регистрация провайдеров
    r.Register("post_express", postexpress.NewProvider)
    r.Register("dex", dex.NewProvider)
    r.Register("aks_express", aks.NewProvider)
    r.Register("bex_express", bex.NewProvider)
    r.Register("d_express", d.NewProvider)
    r.Register("city_express", city.NewProvider)
    r.Register("dhl_express", dhl.NewProvider)

    return r
}

func (r *Registry) GetProvider(code string, cfg ProviderConfig) (DeliveryProvider, error)
func (r *Registry) ListAvailableProviders() []string
```

### Добавление нового провайдера

**Шаг 1**: Создать пакет провайдера

```
internal/gateway/provider/
├── postexpress/
│   ├── provider.go
│   ├── client.go
│   └── types.go
└── dex/              ← Новый провайдер
    ├── provider.go
    ├── client.go
    └── types.go
```

**Шаг 2**: Реализовать интерфейс DeliveryProvider

```go
package dex

type Provider struct {
    client *Client
    config Config
}

func NewProvider(cfg ProviderConfig) (provider.DeliveryProvider, error) {
    return &Provider{
        client: NewClient(cfg.APIKey, cfg.BaseURL),
        config: cfg,
    }, nil
}

func (p *Provider) GetCode() string { return "dex" }
func (p *Provider) CreateShipment(ctx, *ShipmentRequest) (*ShipmentResponse, error) {
    // Конвертация универсального запроса в API Dex
    dexReq := p.mapToDexRequest(req)

    // Вызов API Dex
    dexResp, err := p.client.CreateShipment(ctx, dexReq)

    // Конвертация ответа Dex в универсальный формат
    return p.mapFromDexResponse(dexResp), nil
}
// ... остальные методы интерфейса
```

**Шаг 3**: Зарегистрировать в Registry

```go
func init() {
    globalRegistry.Register("dex", dex.NewProvider)
}
```

**Шаг 4**: Добавить конфигурацию

```env
SVETUDELIVERY_GATEWAYS_DEX_ENABLED=true
SVETUDELIVERY_GATEWAYS_DEX_API_KEY=xxx
SVETUDELIVERY_GATEWAYS_DEX_BASE_URL=https://api.dex.rs
```

### Мультипровайдерный расчет

```go
func (s *CalculatorService) CalculateForAllProviders(ctx, *RateRequest) ([]ProviderRateResponse, error) {
    providers := s.registry.ListAvailableProviders()

    // Параллельные запросы к провайдерам
    results := make(chan ProviderRateResponse, len(providers))

    for _, code := range providers {
        go func(providerCode string) {
            provider, _ := s.registry.GetProvider(providerCode)
            rate, err := provider.CalculateRate(ctx, req)
            results <- ProviderRateResponse{
                Provider: providerCode,
                Rate:     rate,
                Error:    err,
            }
        }(code)
    }

    // Сбор результатов
    var responses []ProviderRateResponse
    for i := 0; i < len(providers); i++ {
        responses = append(responses, <-results)
    }

    // Сортировка по цене или скорости
    sort.Slice(responses, func(i, j int) bool {
        return responses[i].Rate.TotalCost < responses[j].Rate.TotalCost
    })

    return responses, nil
}
```

---

## 🚀 Развертывание на production

### Infrastructure

**Kubernetes manifests**: `k8s/delivery-service/`

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: delivery-service
spec:
  replicas: 3
  template:
    spec:
      containers:
      - name: delivery
        image: registry.svetu.rs/delivery:v1.0.0
        ports:
        - containerPort: 50052  # gRPC
        - containerPort: 9091   # Metrics
        env:
        - name: SVETUDELIVERY_DATABASE_HOST
          valueFrom:
            secretKeyRef:
              name: delivery-db-secret
              key: host
        livenessProbe:
          grpc:
            port: 50052
        readinessProbe:
          grpc:
            port: 50052

---
# service.yaml
apiVersion: v1
kind: Service
metadata:
  name: delivery-service
spec:
  type: ClusterIP
  ports:
  - name: grpc
    port: 50052
    targetPort: 50052
  - name: metrics
    port: 9091
    targetPort: 9091
  selector:
    app: delivery-service
```

### Monitoring

**Prometheus metrics**:
- `delivery_requests_total{method, status}` - количество запросов
- `delivery_request_duration_seconds{method}` - латентность
- `delivery_provider_requests_total{provider, status}` - запросы к провайдерам
- `delivery_shipments_created_total{provider}` - созданные отправления
- `delivery_shipments_status{status}` - распределение по статусам

**Grafana dashboard**: `monitoring/delivery-dashboard.json`

### Logging

**Структурированные логи** (JSON):
```json
{
  "level": "info",
  "timestamp": "2025-10-22T20:00:00Z",
  "service": "delivery",
  "version": "1.0.0",
  "method": "CreateShipment",
  "shipment_id": "uuid",
  "provider": "post_express",
  "duration_ms": 250,
  "user_id": "uuid"
}
```

---

## ✅ Критерии успеха миграции

1. ✅ **Функциональность**: Все эндпоинты работают идентично монолиту
2. ✅ **Производительность**: Latency < 100ms (99th percentile)
3. ✅ **Надежность**: Uptime > 99.9%
4. ✅ **Тесты**: Coverage > 80%
5. ✅ **Документация**: README, API docs, runbooks
6. ✅ **Мониторинг**: Dashboards, алерты
7. ✅ **Обратная совместимость**: Нет breaking changes для монолита

---

## 📝 Чеклист перед удалением старого кода

- [ ] Микросервис работает в production > 2 недель без инцидентов
- [ ] Feature flag = 100% пользователей на микросервисе
- [ ] Все интеграционные тесты проходят
- [ ] Проведен load test (> 1000 RPS)
- [ ] Настроены алерты и мониторинг
- [ ] Создан runbook для on-call
- [ ] Проведен code review финальной версии
- [ ] Обновлена документация
- [ ] Создан план rollback
- [ ] Получено одобрение tech lead

---

## 🔄 Rollback план

Если возникают критические проблемы:

1. **Немедленный откат**: `DELIVERY_USE_MICROSERVICE=false` в монолите
2. **Откат развертывания**: `kubectl rollout undo deployment/delivery-service`
3. **Миграция данных обратно**: Скрипт копирования из delivery_db → svetubd
4. **Проверка**: Запуск smoke tests монолита
5. **Постмортем**: Анализ причин, plan remediation

---

## 📚 Ресурсы

**Репозитории**:
- Микросервис: https://github.com/sveturs/delivery
- Монолит: https://github.com/sveturs/svetu

**Документация**:
- API спецификация: `proto/delivery/v1/delivery.proto`
- Схема БД: `migrations/0001_create_shipments_table.up.sql`
- Runbook: `docs/DELIVERY_RUNBOOK.md` (создать)

**Контакты**:
- Tech Lead: @tech-lead
- DevOps: @devops-team

---

**Последнее обновление**: 2025-10-22
**Следующий review**: После завершения Фазы 1
