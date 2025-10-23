# Delivery Microservice Migration: Clean Cut Plan

**Дата**: 2025-10-22
**Статус продукта**: Pre-production (не в продакшне)
**Подход**: Полный переход без обратной совместимости
**Срок**: 3-4 недели

> **📚 ВАЖНО**: Этот файл стал очень большим (3337 строк).
> Для удобства создана **модульная документация** в директории `delivery-migration/`.
>
> **Переходите по ссылке**: [delivery-migration/README.md](delivery-migration/README.md)
>
> Модульная структура содержит:
> - 📋 Навигацию по всем разделам
> - 🚀 Инфраструктурную документацию (Docker, Nginx, порты)
> - 📝 Пошаговые инструкции развертывания
> - ✅ Детальные чеклисты по фазам
> - 🔍 Troubleshooting guide

---

## 🎯 Цель

Вынести функциональность доставки из монолита в отдельный gRPC микросервис **БЕЗ** промежуточных состояний, feature flags и canary deployment.

**Принцип**: Clean Cut - удаляем старое, внедряем новое, никаких компромиссов.

---

## 📊 Что есть сейчас

### Монолит: `backend/internal/proj/delivery/`

**~2500 строк кода, полная функциональность**:
- ✅ Универсальный интерфейс DeliveryProvider
- ✅ Factory для создания провайдеров
- ✅ Post Express адаптер (B2B интеграция)
- ✅ Mock провайдеры (6+ штук)
- ✅ Service layer (расчет, создание, трекинг)
- ✅ Storage layer (PostgreSQL)
- ✅ Calculator (оптимизация упаковки)
- ✅ Attributes (атрибуты товаров)
- ✅ Admin (управление, аналитика)
- ✅ Notifications (уведомления)

### Микросервис: `github.com/sveturs/delivery`

**~35% готовности, только скелет**:
- ✅ Proto API (gRPC спецификация)
- ✅ Database + migrations
- ✅ Config management
- ✅ Logging infrastructure
- ✅ Makefile (build/lint/test)
- ❌ Domain models - НЕТ
- ❌ Service layer - НЕТ
- ❌ Repository - НЕТ
- ❌ Gateway (провайдеры) - НЕТ
- ❌ Tests - НЕТ

---

## 📐 Архитектура: Текущая vs Целевая

### ТЕКУЩЕЕ СОСТОЯНИЕ: Монолит

```
┌─────────────────────────────────────────────────────────────────┐
│                     МОНОЛИТ (backend/)                          │
│                                                                 │
│  ┌────────────────────────────────────────────────────────┐   │
│  │ internal/proj/                                          │   │
│  │                                                          │   │
│  │  ├─ marketplace/      # Marketplace объявления         │   │
│  │  ├─ storefronts/      # Витрины продавцов              │   │
│  │  ├─ users/            # Пользователи                   │   │
│  │  ├─ payments/         # Платежи                        │   │
│  │  ├─ notifications/    # Уведомления                    │   │
│  │  ├─ chat/             # Чат                            │   │
│  │  ├─ search/           # Поиск (OpenSearch)             │   │
│  │  ├─ admin/            # Админка                        │   │
│  │  │                                                       │   │
│  │  └─ delivery/         # ⚠️ DELIVERY (2500 строк)       │   │
│  │      ├─ calculator/   # Расчет стоимости               │   │
│  │      ├─ factory/      # Provider factory               │   │
│  │      ├─ handler/      # REST handlers                  │   │
│  │      ├─ interfaces/   # DeliveryProvider interface     │   │
│  │      ├─ models/       # Domain models                  │   │
│  │      ├─ service/      # Business logic                 │   │
│  │      ├─ storage/      # PostgreSQL repos               │   │
│  │      ├─ attributes/   # Товарные атрибуты              │   │
│  │      └─ notifications/# Интеграция с уведомлениями     │   │
│  │                                                          │   │
│  └──────────────────────────────────────────────────────────┘   │
│                                                                 │
│  ┌────────────────────────────────────────────────────────┐   │
│  │ internal/proj/postexpress/  # Post Express интеграция  │   │
│  │  ├─ client.go               # HTTP клиент API          │   │
│  │  ├─ service.go              # Бизнес-логика           │   │
│  │  ├─ types.go                # Request/Response типы    │   │
│  │  └─ config.go               # Конфигурация            │   │
│  └──────────────────────────────────────────────────────────┘   │
│                                                                 │
│  ┌────────────────────────────────────────────────────────┐   │
│  │ storage/postgres/                                       │   │
│  │  └─ svetubd (БД монолита)                             │   │
│  │     ├─ marketplace_listings                            │   │
│  │     ├─ users                                           │   │
│  │     ├─ delivery_shipments          ⚠️                  │   │
│  │     ├─ delivery_providers          ⚠️                  │   │
│  │     ├─ delivery_tracking_events    ⚠️                  │   │
│  │     └─ delivery_category_defaults  ⚠️                  │   │
│  └──────────────────────────────────────────────────────────┘   │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘

⚠️ = Будет перенесено в микросервис
```

### ЦЕЛЕВОЕ СОСТОЯНИЕ: Монолит + Микросервис

```
┌──────────────────────────────────────┐   ┌─────────────────────────────────┐
│   МОНОЛИТ (backend/)                 │   │  DELIVERY MICROSERVICE          │
│                                      │   │  (github.com/sveturs/delivery)  │
│  ┌────────────────────────────────┐ │   │                                 │
│  │ internal/proj/                 │ │   │  ┌──────────────────────────┐  │
│  │                                 │ │   │  │ cmd/server/              │  │
│  │  ├─ marketplace/                │ │   │  │  └─ main.go              │  │
│  │  ├─ storefronts/                │ │   │  └──────────────────────────┘  │
│  │  ├─ users/                      │ │   │                                 │
│  │  ├─ payments/                   │ │   │  ┌──────────────────────────┐  │
│  │  ├─ notifications/              │ │   │  │ internal/                │  │
│  │  ├─ chat/                       │ │   │  │                          │  │
│  │  ├─ search/                     │ │   │  │  ├─ domain/              │  │
│  │  ├─ admin/                      │ │   │  │  │  └─ models.go         │  │
│  │  │                               │ │   │  │                          │  │
│  │  └─ delivery/  ✅ ТОНКИЙ СЛОЙ   │ │   │  │  ├─ repository/         │  │
│  │      ├─ client.go  (gRPC wrap) │ │   │  │  │  └─ postgres/         │  │
│  │      ├─ handler.go (proxy)     │ │   │  │                          │  │
│  │      ├─ module.go  (init)      │ │   │  │  ├─ service/             │  │
│  │      └─ types.go   (req/resp)  │ │   │  │  │  ├─ delivery.go       │  │
│  │                                 │ │   │  │  │  ├─ calculator.go     │  │
│  │      (~150 строк вместо 2500)  │ │   │  │  │  └─ tracking.go       │  │
│  └────────────────────────────────┘ │   │  │                          │  │
│                                      │   │  │  ├─ gateway/             │  │
│  ┌────────────────────────────────┐ │   │  │  │  └─ provider/         │  │
│  │ go.mod                          │ │   │  │      ├─ interface.go    │  │
│  │  require (                      │ │   │  │      ├─ factory.go      │  │
│  │   github.com/sveturs/delivery  │←─┼───┼──│      ├─ postexpress/    │  │
│  │     v1.0.0                      │ │   │  │      ├─ dex/            │  │
│  │  )                              │ │   │  │      └─ mock/           │  │
│  └────────────────────────────────┘ │   │  │                          │  │
│                                      │   │  │  └─ server/              │  │
│  ┌────────────────────────────────┐ │   │  │      └─ grpc/            │  │
│  │ PostgreSQL: svetubd             │ │   │  │          └─ delivery.go │  │
│  │  ├─ marketplace_listings        │ │   │  └──────────────────────────┘  │
│  │  ├─ users                       │ │   │                                 │
│  │  ├─ storefronts                 │ │   │  ┌──────────────────────────┐  │
│  │  └─ ... (все кроме delivery)   │ │   │  │ PostgreSQL: delivery_db  │  │
│  └────────────────────────────────┘ │   │  │  ├─ shipments            │  │
│                                      │   │  │  ├─ tracking_events      │  │
│                         gRPC Request │   │  │  └─ providers            │  │
│  handlers ───────────────────────────┼───┼──┤                          │  │
│       │                              │   │  └──────────────────────────┘  │
│       │ CreateShipment()             │   │                                 │
│       │ TrackShipment()              │   │  ┌──────────────────────────┐  │
│       └──> deliveryService.xxx() ───┼───┼─>│ gRPC Server :50052       │  │
│            (pkg/service wrapper)     │   │  │                          │  │
│                                      │   │  │  CreateShipment()        │  │
│                                      │   │  │  GetShipment()           │  │
│                                      │   │  │  TrackShipment()         │  │
│                                      │   │  │  CalculateRate()         │  │
│                                      │   │  │  CancelShipment()        │  │
│                                      │   │  └──────────────────────────┘  │
│                                      │   │                                 │
│  Port: 3000 (HTTP REST)              │   │  Port: 50052 (gRPC)             │
└──────────────────────────────────────┘   └─────────────────────────────────┘
```

---

## 📊 Детальное сравнение компонентов

### Что переносится из монолита в микросервис

| Компонент | Текущее место | Целевое место | Строк кода | Статус |
|-----------|---------------|---------------|------------|--------|
| **Domain Models** | `backend/internal/proj/delivery/models/` | `delivery/internal/domain/` | ~300 | ✅ Переносится |
| **Repository** | `backend/internal/proj/delivery/storage/` | `delivery/internal/repository/` | ~400 | ✅ Переносится |
| **Service Logic** | `backend/internal/proj/delivery/service/` | `delivery/internal/service/` | ~700 | ✅ Переносится |
| **Calculator** | `backend/internal/proj/delivery/calculator/` | `delivery/internal/service/calculator.go` | ~300 | ✅ Переносится |
| **Provider Interface** | `backend/internal/proj/delivery/interfaces/` | `delivery/internal/gateway/provider/` | ~200 | ✅ Переносится |
| **Factory** | `backend/internal/proj/delivery/factory/` | `delivery/internal/gateway/provider/factory.go` | ~150 | ✅ Переносится |
| **Post Express** | `backend/internal/proj/postexpress/` | `delivery/internal/gateway/provider/postexpress/` | ~600 | ✅ Переносится |
| **Attributes** | `backend/internal/proj/delivery/attributes/` | `delivery/internal/service/attributes.go` | ~200 | ✅ Переносится |
| **Notifications** | `backend/internal/proj/delivery/notifications/` | ❌ Удаляется | ~150 | ❌ Не нужно в микросервисе |

**Итого переносится**: ~2850 строк

### Что остается в монолите (новый тонкий слой)

| Компонент | Файл | Назначение | Строк кода |
|-----------|------|------------|------------|
| **gRPC Client Wrapper** | `backend/internal/proj/delivery/client.go` | Обертка над gRPC клиентом | ~30 |
| **HTTP Handlers** | `backend/internal/proj/delivery/handler.go` | Proxy запросы в микросервис | ~100 |
| **Module** | `backend/internal/proj/delivery/module.go` | Инициализация и routes | ~50 |
| **Types** | `backend/internal/proj/delivery/types.go` | Request/Response типы | ~50 |

**Итого в монолите**: ~230 строк (было 2500!)

### Библиотека микросервиса для монолита

| Пакет | Файлы | Назначение | Строк кода |
|-------|-------|------------|------------|
| **pkg/client** | `client.go`, `types.go`, `converter.go` | Низкоуровневый gRPC клиент | ~400 |
| **pkg/service** | `delivery.go`, `validator.go`, `retry.go`, `cache.go` | Высокоуровневая обертка | ~600 |

**Итого библиотека**: ~1000 строк

---

## 🗄️ База данных: Текущая vs Целевая

### ТЕКУЩЕЕ: Одна БД (svetubd)

```sql
-- PostgreSQL: svetubd (монолит)

-- Все таблицы вместе:
marketplace_listings
marketplace_categories
marketplace_orders
users
user_profiles
storefronts
storefront_products
payments
payment_transactions
delivery_shipments              ⚠️ → микросервис
delivery_providers              ⚠️ → микросервис
delivery_tracking_events        ⚠️ → микросервис
delivery_category_defaults      ⚠️ → микросервис
delivery_pricing_rules          ⚠️ → микросервис
delivery_zones                  ⚠️ → микросервис
chat_messages
notifications
```

### ЦЕЛЕВОЕ: Две БД

```sql
-- PostgreSQL: svetubd (монолит)
marketplace_listings
marketplace_categories
marketplace_orders
users
user_profiles
storefronts
storefront_products
payments
payment_transactions
chat_messages
notifications
-- delivery таблицы УДАЛЕНЫ ❌
```

```sql
-- PostgreSQL: delivery_db (микросервис)
shipments                       ✅ НОВАЯ
tracking_events                 ✅ НОВАЯ
providers                       ✅ НОВАЯ (опционально)
-- Простая схема, только essential данные
```

**Миграция данных**: Одноразовый скрипт копирования `delivery_*` таблиц → `delivery_db`

---

## 🔄 Взаимодействие: Текущее vs Целевое

### ТЕКУЩЕЕ: Прямые вызовы

```go
// В монолите
deliveryService := service.NewService(db, providerFactory)
shipment, err := deliveryService.CreateShipment(ctx, req)
```

**Проблема**: Вся логика в монолите, невозможно масштабировать отдельно

### ЦЕЛЕВОЕ: gRPC через библиотеку

```go
// В монолите
import "github.com/sveturs/delivery/pkg/service"

deliveryService := service.NewDeliveryService(&service.Config{
    GRPCAddress: "delivery-service:50052",
    RetryAttempts: 3,
    CacheEnabled: true,
})

// Вызов идентичный, но идет в микросервис
shipment, err := deliveryService.CreateShipment(ctx, req)
```

**Преимущества**:
- ✅ Независимое развертывание
- ✅ Независимое масштабирование
- ✅ Изоляция сбоев
- ✅ Переиспользование в других проектах

---

## 🏗️ Архитектура микросервиса (детально)

```
┌─────────────────────────────────────────┐
│    gRPC Server (port 50052)             │
│    internal/server/grpc/delivery.go     │
└──────────────────┬──────────────────────┘
                   │
┌──────────────────▼──────────────────────┐
│         Service Layer                   │
│  internal/service/                      │
│  ├─ delivery_service.go                 │
│  ├─ calculator_service.go               │
│  └─ tracking_service.go                 │
└──────┬────────────────────┬─────────────┘
       │                    │
┌──────▼─────────┐   ┌──────▼──────────┐
│  Repository    │   │   Gateway       │
│  internal/     │   │   internal/     │
│  repository/   │   │   gateway/      │
│                │   │   provider/     │
│ - shipments    │   │  ├─postexpress/ │
│ - events       │   │  ├─dex/         │
│ - providers    │   │  └─mock/        │
└────────────────┘   └─────────────────┘
       │
┌──────▼─────────┐
│  PostgreSQL    │
│  delivery_db   │
└────────────────┘
```

---

## 📋 План миграции (3 фазы)

### ФАЗА 1: Реализация микросервиса (Week 1-2)

#### 1.1 Генерация proto кода

```bash
cd ~/delivery
make proto
```

**Результат**: `gen/go/delivery/v1/` с gRPC клиентом/сервером

#### 1.2 Domain Layer

**Файл**: `internal/domain/models.go`

```go
package domain

type Shipment struct {
    ID                 uuid.UUID
    TrackingNumber     string
    Status             ShipmentStatus
    Provider           DeliveryProvider
    UserID             uuid.UUID
    FromAddress        Address
    ToAddress          Address
    Package            Package
    Cost               Money
    ProviderShipmentID *string
    ProviderMetadata   json.RawMessage
    EstimatedDelivery  *time.Time
    ActualDelivery     *time.Time
    CreatedAt          time.Time
    UpdatedAt          time.Time
}

type Address struct {
    Street     string
    City       string
    State      string
    PostalCode string
    Country    string
    Phone      string
    Email      string
    Name       string
}

type Package struct {
    WeightKg    float64
    LengthCm    float64
    WidthCm     float64
    HeightCm    float64
    Description string
    Value       float64
}

type TrackingEvent struct {
    ID         uuid.UUID
    ShipmentID uuid.UUID
    Status     ShipmentStatus
    Location   string
    Details    string
    Timestamp  time.Time
    CreatedAt  time.Time
}

type ShipmentStatus string

const (
    StatusPending          ShipmentStatus = "pending"
    StatusConfirmed        ShipmentStatus = "confirmed"
    StatusInTransit        ShipmentStatus = "in_transit"
    StatusOutForDelivery   ShipmentStatus = "out_for_delivery"
    StatusDelivered        ShipmentStatus = "delivered"
    StatusFailed           ShipmentStatus = "failed"
    StatusCancelled        ShipmentStatus = "cancelled"
    StatusReturned         ShipmentStatus = "returned"
)

type DeliveryProvider string

const (
    ProviderPostExpress DeliveryProvider = "post_express"
    ProviderDex         DeliveryProvider = "dex"
)
```

**Файл**: `internal/domain/converter.go`

```go
package domain

import pb "github.com/sveturs/delivery/gen/go/delivery/v1"

// ToProto конвертирует domain модель в protobuf
func (s *Shipment) ToProto() *pb.Shipment {
    return &pb.Shipment{
        Id:             s.ID.String(),
        TrackingNumber: s.TrackingNumber,
        Status:         pb.ShipmentStatus(pb.ShipmentStatus_value[string(s.Status)]),
        // ... остальные поля
    }
}

// FromProto конвертирует protobuf в domain модель
func ShipmentFromProto(pb *pb.Shipment) (*Shipment, error) {
    id, err := uuid.Parse(pb.Id)
    if err != nil {
        return nil, err
    }
    return &Shipment{
        ID:             id,
        TrackingNumber: pb.TrackingNumber,
        // ... остальные поля
    }, nil
}
```

#### 1.3 Repository Layer

**Файл**: `internal/repository/shipment_repository.go`

```go
package repository

import (
    "context"
    "database/sql"
    "github.com/sveturs/delivery/internal/domain"
)

type ShipmentRepository interface {
    Create(ctx context.Context, shipment *domain.Shipment) error
    GetByID(ctx context.Context, id uuid.UUID) (*domain.Shipment, error)
    GetByTracking(ctx context.Context, trackingNumber string) (*domain.Shipment, error)
    UpdateStatus(ctx context.Context, id uuid.UUID, status domain.ShipmentStatus, deliveredAt *time.Time) error
    List(ctx context.Context, filter ListFilter) ([]*domain.Shipment, error)
}

type PostgresShipmentRepository struct {
    db *sql.DB
}

func NewPostgresShipmentRepository(db *sql.DB) *PostgresShipmentRepository {
    return &PostgresShipmentRepository{db: db}
}

func (r *PostgresShipmentRepository) Create(ctx context.Context, shipment *domain.Shipment) error {
    query := `
        INSERT INTO shipments (
            id, tracking_number, status, provider, user_id,
            from_street, from_city, from_state, from_postal_code, from_country, from_phone, from_email, from_name,
            to_street, to_city, to_state, to_postal_code, to_country, to_phone, to_email, to_name,
            weight_kg, length_cm, width_cm, height_cm, package_description, package_value,
            cost, currency, provider_shipment_id, provider_metadata,
            estimated_delivery_at
        ) VALUES (
            $1, $2, $3, $4, $5,
            $6, $7, $8, $9, $10, $11, $12, $13,
            $14, $15, $16, $17, $18, $19, $20, $21,
            $22, $23, $24, $25, $26, $27,
            $28, $29, $30, $31, $32
        )
    `

    _, err := r.db.ExecContext(ctx, query,
        shipment.ID,
        shipment.TrackingNumber,
        shipment.Status,
        shipment.Provider,
        shipment.UserID,
        // from address
        shipment.FromAddress.Street,
        shipment.FromAddress.City,
        shipment.FromAddress.State,
        shipment.FromAddress.PostalCode,
        shipment.FromAddress.Country,
        shipment.FromAddress.Phone,
        shipment.FromAddress.Email,
        shipment.FromAddress.Name,
        // to address
        shipment.ToAddress.Street,
        shipment.ToAddress.City,
        shipment.ToAddress.State,
        shipment.ToAddress.PostalCode,
        shipment.ToAddress.Country,
        shipment.ToAddress.Phone,
        shipment.ToAddress.Email,
        shipment.ToAddress.Name,
        // package
        shipment.Package.WeightKg,
        shipment.Package.LengthCm,
        shipment.Package.WidthCm,
        shipment.Package.HeightCm,
        shipment.Package.Description,
        shipment.Package.Value,
        // cost
        shipment.Cost.Amount,
        shipment.Cost.Currency,
        shipment.ProviderShipmentID,
        shipment.ProviderMetadata,
        shipment.EstimatedDelivery,
    )

    return err
}

func (r *PostgresShipmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Shipment, error) {
    query := `SELECT * FROM shipments WHERE id = $1`
    // ... реализация
}

func (r *PostgresShipmentRepository) GetByTracking(ctx context.Context, trackingNumber string) (*domain.Shipment, error) {
    query := `SELECT * FROM shipments WHERE tracking_number = $1`
    // ... реализация
}
```

**Источник кода**: Адаптировать из `backend/internal/proj/delivery/storage/storage.go`

#### 1.4 Gateway Layer (Provider Pattern)

**Файл**: `internal/gateway/provider/interface.go`

```go
package provider

type Provider interface {
    GetCode() string
    GetName() string
    IsAvailable() bool
    GetCapabilities() *Capabilities

    CalculateRate(ctx context.Context, req *RateRequest) (*RateResponse, error)
    CreateShipment(ctx context.Context, req *ShipmentRequest) (*ShipmentResponse, error)
    TrackShipment(ctx context.Context, trackingNumber string) (*TrackingResponse, error)
    CancelShipment(ctx context.Context, shipmentID string) error
    ValidateAddress(ctx context.Context, address *Address) (*AddressValidation, error)
}

type Capabilities struct {
    MaxWeightKg       float64
    MaxVolumeM3       float64
    SupportedZones    []string // local, national, international
    SupportedTypes    []string // standard, express
    SupportsCOD       bool
    SupportsInsurance bool
    SupportsTracking  bool
}

type RateRequest struct {
    FromAddress *Address
    ToAddress   *Address
    Package     *Package
    Type        string // standard, express
}

type RateResponse struct {
    Options []RateOption
}

type RateOption struct {
    Type          string  // standard, express
    Cost          float64
    Currency      string
    EstimatedDays int
}
```

**Файл**: `internal/gateway/provider/factory.go`

```go
package provider

type Factory struct {
    providers map[string]Provider
    config    *config.Config
}

func NewFactory(cfg *config.Config) *Factory {
    f := &Factory{
        providers: make(map[string]Provider),
        config:    cfg,
    }

    // Регистрация провайдеров
    if cfg.Gateways.PostRS.Enabled {
        f.providers["post_express"] = postexpress.NewProvider(&cfg.Gateways.PostRS)
    }

    if cfg.Gateways.Dex.Enabled {
        f.providers["dex"] = dex.NewProvider(&cfg.Gateways.Dex)
    }

    // Mock провайдер всегда доступен для тестирования
    f.providers["mock"] = mock.NewProvider()

    return f
}

func (f *Factory) GetProvider(code string) (Provider, error) {
    provider, exists := f.providers[code]
    if !exists {
        return nil, fmt.Errorf("provider not found: %s", code)
    }
    return provider, nil
}

func (f *Factory) ListProviders() []Provider {
    providers := make([]Provider, 0, len(f.providers))
    for _, p := range f.providers {
        providers = append(providers, p)
    }
    return providers
}
```

#### 1.5 Post Express Integration

**Структура**:
```
internal/gateway/provider/postexpress/
├── provider.go      # Реализация интерфейса Provider
├── client.go        # HTTP клиент для API Post Express
├── types.go         # Типы запросов/ответов
├── mapper.go        # Маппинг domain ↔ Post Express API
└── validator.go     # Валидация B2B полей
```

**Файл**: `internal/gateway/provider/postexpress/provider.go`

```go
package postexpress

type Provider struct {
    client *Client
    config *Config
}

func NewProvider(cfg *Config) *Provider {
    return &Provider{
        client: NewClient(cfg.APIKey, cfg.BaseURL, cfg.Timeout),
        config: cfg,
    }
}

func (p *Provider) GetCode() string {
    return "post_express"
}

func (p *Provider) CreateShipment(ctx context.Context, req *provider.ShipmentRequest) (*provider.ShipmentResponse, error) {
    // 1. Валидация
    if err := p.validateRequest(req); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }

    // 2. Маппинг в формат Post Express B2B API
    peReq := p.mapToPostExpressRequest(req)

    // 3. Вызов API
    peResp, err := p.client.CreateShipment(ctx, peReq)
    if err != nil {
        return nil, fmt.Errorf("post express api error: %w", err)
    }

    // 4. Маппинг обратно
    return p.mapFromPostExpressResponse(peResp), nil
}
```

**Источник**: Полный перенос из `backend/internal/proj/postexpress/` и `backend/internal/proj/delivery/factory/postexpress_adapter.go`

**ВАЖНО**: Сохранить ВСЮ B2B логику:
- ExtBrend, ExtMagacin, ExtReferenca
- NacinPrijema, NacinPlacanja
- Otkupnina (COD) с банковскими реквизитами
- PosebneUsluge (PNA, SMS, OTK, VD)
- Валидация всех обязательных полей
- Маппинг статусов

#### 1.6 Service Layer

**Файл**: `internal/service/delivery_service.go`

```go
package service

type DeliveryService struct {
    repo     repository.ShipmentRepository
    eventRepo repository.TrackingEventRepository
    factory  *provider.Factory
    logger   *logger.Logger
}

func NewDeliveryService(
    repo repository.ShipmentRepository,
    eventRepo repository.TrackingEventRepository,
    factory *provider.Factory,
    logger *logger.Logger,
) *DeliveryService {
    return &DeliveryService{
        repo:      repo,
        eventRepo: eventRepo,
        factory:   factory,
        logger:    logger,
    }
}

func (s *DeliveryService) CreateShipment(ctx context.Context, input *CreateShipmentInput) (*domain.Shipment, error) {
    // 1. Получаем провайдера
    provider, err := s.factory.GetProvider(input.ProviderCode)
    if err != nil {
        return nil, fmt.Errorf("provider not found: %w", err)
    }

    // 2. Создаем shipment через провайдера
    providerResp, err := provider.CreateShipment(ctx, &provider.ShipmentRequest{
        FromAddress: input.FromAddress,
        ToAddress:   input.ToAddress,
        Package:     input.Package,
        Type:        input.Type,
    })
    if err != nil {
        return nil, fmt.Errorf("provider failed: %w", err)
    }

    // 3. Сохраняем в БД
    shipment := &domain.Shipment{
        ID:                 uuid.New(),
        TrackingNumber:     providerResp.TrackingNumber,
        Status:             domain.StatusConfirmed,
        Provider:           domain.DeliveryProvider(input.ProviderCode),
        UserID:             input.UserID,
        FromAddress:        input.FromAddress,
        ToAddress:          input.ToAddress,
        Package:            input.Package,
        Cost:               providerResp.Cost,
        ProviderShipmentID: &providerResp.ProviderShipmentID,
        EstimatedDelivery:  providerResp.EstimatedDelivery,
        CreatedAt:          time.Now(),
        UpdatedAt:          time.Now(),
    }

    if err := s.repo.Create(ctx, shipment); err != nil {
        return nil, fmt.Errorf("failed to save shipment: %w", err)
    }

    s.logger.Info().
        Str("shipment_id", shipment.ID.String()).
        Str("tracking_number", shipment.TrackingNumber).
        Str("provider", string(shipment.Provider)).
        Msg("Shipment created successfully")

    return shipment, nil
}

func (s *DeliveryService) GetShipment(ctx context.Context, id uuid.UUID) (*domain.Shipment, error) {
    return s.repo.GetByID(ctx, id)
}

func (s *DeliveryService) TrackShipment(ctx context.Context, trackingNumber string) (*TrackingInfo, error) {
    // 1. Получаем shipment из БД
    shipment, err := s.repo.GetByTracking(ctx, trackingNumber)
    if err != nil {
        return nil, fmt.Errorf("shipment not found: %w", err)
    }

    // 2. Получаем провайдера
    provider, err := s.factory.GetProvider(string(shipment.Provider))
    if err != nil {
        return nil, fmt.Errorf("provider not found: %w", err)
    }

    // 3. Запрашиваем актуальный статус у провайдера
    tracking, err := provider.TrackShipment(ctx, trackingNumber)
    if err != nil {
        // Провайдер недоступен - возвращаем последний известный статус
        s.logger.Warn().Err(err).Msg("Provider unavailable, returning cached status")
        events, _ := s.eventRepo.ListByShipment(ctx, shipment.ID)
        return &TrackingInfo{
            Shipment: shipment,
            Events:   events,
        }, nil
    }

    // 4. Обновляем статус если изменился
    if tracking.Status != string(shipment.Status) {
        newStatus := domain.ShipmentStatus(tracking.Status)
        if err := s.repo.UpdateStatus(ctx, shipment.ID, newStatus, tracking.DeliveredAt); err != nil {
            s.logger.Error().Err(err).Msg("Failed to update shipment status")
        }
        shipment.Status = newStatus
    }

    // 5. Сохраняем новые события
    for _, event := range tracking.Events {
        trackingEvent := &domain.TrackingEvent{
            ID:         uuid.New(),
            ShipmentID: shipment.ID,
            Status:     domain.ShipmentStatus(event.Status),
            Location:   event.Location,
            Details:    event.Details,
            Timestamp:  event.Timestamp,
            CreatedAt:  time.Now(),
        }
        if err := s.eventRepo.Create(ctx, trackingEvent); err != nil {
            s.logger.Error().Err(err).Msg("Failed to save tracking event")
        }
    }

    return &TrackingInfo{
        Shipment: shipment,
        Events:   tracking.Events,
    }, nil
}

func (s *DeliveryService) CancelShipment(ctx context.Context, id uuid.UUID) error {
    // ... реализация
}
```

**Файл**: `internal/service/calculator_service.go`

```go
package service

type CalculatorService struct {
    factory *provider.Factory
    logger  *logger.Logger
}

func (s *CalculatorService) CalculateRates(ctx context.Context, req *CalculateRatesInput) (*CalculateRatesOutput, error) {
    providers := s.factory.ListProviders()

    // Параллельный запрос ко всем провайдерам
    results := make(chan ProviderRateResult, len(providers))

    for _, p := range providers {
        go func(provider provider.Provider) {
            rate, err := provider.CalculateRate(ctx, &provider.RateRequest{
                FromAddress: req.FromAddress,
                ToAddress:   req.ToAddress,
                Package:     req.Package,
                Type:        req.Type,
            })
            results <- ProviderRateResult{
                Provider: provider.GetCode(),
                Rate:     rate,
                Error:    err,
            }
        }(p)
    }

    // Сбор результатов
    var rates []ProviderRateResult
    for i := 0; i < len(providers); i++ {
        result := <-results
        if result.Error == nil {
            rates = append(rates, result)
        } else {
            s.logger.Warn().
                Str("provider", result.Provider).
                Err(result.Error).
                Msg("Provider rate calculation failed")
        }
    }

    // Сортировка по цене
    sort.Slice(rates, func(i, j int) bool {
        return rates[i].Rate.Cost < rates[j].Rate.Cost
    })

    return &CalculateRatesOutput{Rates: rates}, nil
}
```

**Источник**: `backend/internal/proj/delivery/service/service.go` и `calculator/service.go`

#### 1.7 gRPC Handlers

**Файл**: `internal/server/grpc/delivery.go`

```go
package grpc

import (
    "context"
    pb "github.com/sveturs/delivery/gen/go/delivery/v1"
    "github.com/sveturs/delivery/internal/service"
    "github.com/sveturs/delivery/internal/domain"
)

type DeliveryServer struct {
    pb.UnimplementedDeliveryServiceServer
    deliveryService   *service.DeliveryService
    calculatorService *service.CalculatorService
}

func NewDeliveryServer(
    deliveryService *service.DeliveryService,
    calculatorService *service.CalculatorService,
) *DeliveryServer {
    return &DeliveryServer{
        deliveryService:   deliveryService,
        calculatorService: calculatorService,
    }
}

func (s *DeliveryServer) CreateShipment(ctx context.Context, req *pb.CreateShipmentRequest) (*pb.CreateShipmentResponse, error) {
    // 1. Валидация protobuf
    if err := validateCreateShipmentRequest(req); err != nil {
        return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
    }

    // 2. Конвертация pb → domain
    input := &service.CreateShipmentInput{
        ProviderCode: req.Provider.String(),
        UserID:       uuid.MustParse(req.UserId),
        FromAddress:  addressFromProto(req.FromAddress),
        ToAddress:    addressFromProto(req.ToAddress),
        Package:      packageFromProto(req.Package),
        Type:         req.Type,
    }

    // 3. Вызов service
    shipment, err := s.deliveryService.CreateShipment(ctx, input)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to create shipment: %v", err)
    }

    // 4. Конвертация domain → pb
    return &pb.CreateShipmentResponse{
        Shipment: shipment.ToProto(),
    }, nil
}

func (s *DeliveryServer) GetShipment(ctx context.Context, req *pb.GetShipmentRequest) (*pb.GetShipmentResponse, error) {
    id, err := uuid.Parse(req.Id)
    if err != nil {
        return nil, status.Errorf(codes.InvalidArgument, "invalid shipment id: %v", err)
    }

    shipment, err := s.deliveryService.GetShipment(ctx, id)
    if err != nil {
        return nil, status.Errorf(codes.NotFound, "shipment not found: %v", err)
    }

    return &pb.GetShipmentResponse{
        Shipment: shipment.ToProto(),
    }, nil
}

func (s *DeliveryServer) TrackShipment(ctx context.Context, req *pb.TrackShipmentRequest) (*pb.TrackShipmentResponse, error) {
    tracking, err := s.deliveryService.TrackShipment(ctx, req.TrackingNumber)
    if err != nil {
        return nil, status.Errorf(codes.NotFound, "tracking failed: %v", err)
    }

    events := make([]*pb.TrackingEvent, len(tracking.Events))
    for i, e := range tracking.Events {
        events[i] = e.ToProto()
    }

    return &pb.TrackShipmentResponse{
        Shipment: tracking.Shipment.ToProto(),
        Events:   events,
    }, nil
}

func (s *DeliveryServer) CalculateRate(ctx context.Context, req *pb.CalculateRateRequest) (*pb.CalculateRateResponse, error) {
    // ... реализация
}

func (s *DeliveryServer) CancelShipment(ctx context.Context, req *pb.CancelShipmentRequest) (*pb.CancelShipmentResponse, error) {
    // ... реализация
}
```

#### 1.8 Инициализация в main.go

**Файл**: `cmd/server/main.go` (обновить)

```go
func main() {
    // Config
    cfg := config.Load()

    // Logger
    logger.Init(cfg.Service.Environment, cfg.Service.LogLevel, version.Version, true, true)

    // Database
    db, err := database.NewPostgresConnection(&cfg.Database)
    if err != nil {
        logger.Fatal().Err(err).Msg("Failed to connect to database")
    }

    // Migrations
    migrator := migrator.NewMigrator(db, cfg.Database.MigrationsPath)
    if err := migrator.Run(); err != nil {
        logger.Fatal().Err(err).Msg("Failed to run migrations")
    }

    // Repositories
    shipmentRepo := repository.NewPostgresShipmentRepository(db)
    eventRepo := repository.NewPostgresTrackingEventRepository(db)

    // Provider Factory
    providerFactory := provider.NewFactory(cfg)

    // Services
    deliveryService := service.NewDeliveryService(shipmentRepo, eventRepo, providerFactory, logger)
    calculatorService := service.NewCalculatorService(providerFactory, logger)

    // gRPC Server
    grpcServer := grpc.NewServer()
    deliveryServer := grpcServer.NewDeliveryServer(deliveryService, calculatorService)
    pb.RegisterDeliveryServiceServer(grpcServer, deliveryServer)

    // Start server
    lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.GRPCPort))
    if err != nil {
        logger.Fatal().Err(err).Msg("Failed to listen")
    }

    logger.Info().Int("port", cfg.Server.GRPCPort).Msg("Starting gRPC server")
    if err := grpcServer.Serve(lis); err != nil {
        logger.Fatal().Err(err).Msg("Failed to serve")
    }
}
```

#### 1.9 Client Library для монолита

Библиотека состоит из двух слоев:
1. **pkg/client** - низкоуровневый gRPC клиент (маппинг protobuf ↔ Go types)
2. **pkg/service** - высокоуровневая обертка с бизнес-логикой

##### 1.9.1 Low-level gRPC Client

**Файл**: `pkg/client/client.go`

```go
package client

import (
    "context"
    pb "github.com/sveturs/delivery/gen/go/delivery/v1"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

type Client struct {
    conn   *grpc.ClientConn
    client pb.DeliveryServiceClient
}

func NewClient(addr string) (*Client, error) {
    conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

func (c *Client) CreateShipment(ctx context.Context, req *CreateShipmentRequest) (*Shipment, error) {
    // Конвертация request → protobuf
    pbReq := &pb.CreateShipmentRequest{
        Provider: pb.DeliveryProvider(pb.DeliveryProvider_value[req.Provider]),
        UserId:   req.UserID.String(),
        FromAddress: &pb.Address{
            Street:     req.FromAddress.Street,
            City:       req.FromAddress.City,
            PostalCode: req.FromAddress.PostalCode,
            Country:    req.FromAddress.Country,
            Phone:      req.FromAddress.Phone,
            Email:      req.FromAddress.Email,
            Name:       req.FromAddress.Name,
        },
        ToAddress: &pb.Address{
            Street:     req.ToAddress.Street,
            City:       req.ToAddress.City,
            PostalCode: req.ToAddress.PostalCode,
            Country:    req.ToAddress.Country,
            Phone:      req.ToAddress.Phone,
            Email:      req.ToAddress.Email,
            Name:       req.ToAddress.Name,
        },
        Package: &pb.Package{
            WeightKg:    req.Package.WeightKg,
            LengthCm:    req.Package.LengthCm,
            WidthCm:     req.Package.WidthCm,
            HeightCm:    req.Package.HeightCm,
            Description: req.Package.Description,
            Value:       req.Package.Value,
        },
        Type: req.Type,
    }

    // Вызов gRPC
    resp, err := c.client.CreateShipment(ctx, pbReq)
    if err != nil {
        return nil, err
    }

    // Конвертация protobuf → response
    return shipmentFromProto(resp.Shipment), nil
}

func (c *Client) GetShipment(ctx context.Context, id uuid.UUID) (*Shipment, error) {
    resp, err := c.client.GetShipment(ctx, &pb.GetShipmentRequest{Id: id.String()})
    if err != nil {
        return nil, err
    }
    return shipmentFromProto(resp.Shipment), nil
}

func (c *Client) TrackShipment(ctx context.Context, trackingNumber string) (*TrackingInfo, error) {
    resp, err := c.client.TrackShipment(ctx, &pb.TrackShipmentRequest{TrackingNumber: trackingNumber})
    if err != nil {
        return nil, err
    }
    return trackingInfoFromProto(resp), nil
}

func (c *Client) CalculateRate(ctx context.Context, req *CalculateRateRequest) (*CalculateRateResponse, error) {
    // ... реализация
}

func (c *Client) CancelShipment(ctx context.Context, id uuid.UUID) error {
    _, err := c.client.CancelShipment(ctx, &pb.CancelShipmentRequest{Id: id.String()})
    return err
}
```

**Файл**: `pkg/client/types.go`

```go
package client

// Go структуры (НЕ protobuf) для удобного использования в монолите
type CreateShipmentRequest struct {
    Provider    string
    UserID      uuid.UUID
    FromAddress Address
    ToAddress   Address
    Package     Package
    Type        string
}

type Shipment struct {
    ID                 uuid.UUID
    TrackingNumber     string
    Status             string
    Provider           string
    Cost               float64
    Currency           string
    EstimatedDelivery  *time.Time
    ActualDelivery     *time.Time
    CreatedAt          time.Time
}

type Address struct {
    Street     string
    City       string
    PostalCode string
    Country    string
    Phone      string
    Email      string
    Name       string
}

type Package struct {
    WeightKg    float64
    LengthCm    float64
    WidthCm     float64
    HeightCm    float64
    Description string
    Value       float64
}

type TrackingInfo struct {
    Shipment *Shipment
    Events   []TrackingEvent
}

type TrackingEvent struct {
    Status    string
    Location  string
    Details   string
    Timestamp time.Time
}
```

##### 1.9.2 High-level Service Wrapper

**Структура pkg**:
```
pkg/
├── client/              # Низкоуровневый gRPC клиент
│   ├── client.go       # gRPC подключение
│   ├── types.go        # Go структуры (не protobuf)
│   └── converter.go    # Маппинг protobuf ↔ types
└── service/            # Высокоуровневая обертка
    ├── delivery.go     # DeliveryService с бизнес-логикой
    ├── calculator.go   # CalculatorService
    ├── validator.go    # Валидация входных данных
    ├── retry.go        # Retry логика
    └── cache.go        # Кеширование (опционально)
```

**Файл**: `pkg/service/delivery.go`

```go
package service

import (
    "context"
    "fmt"
    "time"

    "github.com/sveturs/delivery/pkg/client"
)

// DeliveryService - высокоуровневая обертка над gRPC клиентом
// Добавляет валидацию, retry, логирование, кеширование
type DeliveryService struct {
    client    *client.Client
    validator *Validator
    retrier   *Retrier
    cache     *Cache // опционально
}

// Config для инициализации сервиса
type Config struct {
    GRPCAddress   string
    RetryAttempts int
    RetryTimeout  time.Duration
    CacheEnabled  bool
    CacheTTL      time.Duration
}

func NewDeliveryService(cfg *Config) (*DeliveryService, error) {
    // Создаем gRPC клиент
    grpcClient, err := client.NewClient(cfg.GRPCAddress)
    if err != nil {
        return nil, fmt.Errorf("failed to create grpc client: %w", err)
    }

    return &DeliveryService{
        client:    grpcClient,
        validator: NewValidator(),
        retrier:   NewRetrier(cfg.RetryAttempts, cfg.RetryTimeout),
        cache:     NewCache(cfg.CacheEnabled, cfg.CacheTTL),
    }, nil
}

func (s *DeliveryService) Close() error {
    return s.client.Close()
}

// CreateShipment с валидацией, retry и обработкой ошибок
func (s *DeliveryService) CreateShipment(ctx context.Context, req *CreateShipmentRequest) (*Shipment, error) {
    // 1. Валидация входных данных
    if err := s.validator.ValidateCreateShipmentRequest(req); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }

    // 2. Нормализация данных (приведение адресов к стандартному формату)
    req.FromAddress = s.normalizeAddress(req.FromAddress)
    req.ToAddress = s.normalizeAddress(req.ToAddress)

    // 3. Вызов gRPC с retry логикой
    var shipment *client.Shipment
    err := s.retrier.Do(ctx, func() error {
        var retryErr error
        shipment, retryErr = s.client.CreateShipment(ctx, &client.CreateShipmentRequest{
            Provider:    req.ProviderCode,
            UserID:      req.UserID,
            FromAddress: req.FromAddress,
            ToAddress:   req.ToAddress,
            Package:     req.Package,
            Type:        req.Type,
        })
        return retryErr
    })

    if err != nil {
        return nil, fmt.Errorf("failed to create shipment: %w", err)
    }

    // 4. Обогащение данных (добавление дополнительной информации)
    enrichedShipment := s.enrichShipment(shipment, req)

    return enrichedShipment, nil
}

// GetShipment с кешированием
func (s *DeliveryService) GetShipment(ctx context.Context, id uuid.UUID) (*Shipment, error) {
    // Проверяем кеш
    if s.cache.Enabled() {
        if cached, found := s.cache.Get(id.String()); found {
            return cached.(*Shipment), nil
        }
    }

    // Вызов gRPC
    shipment, err := s.client.GetShipment(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to get shipment: %w", err)
    }

    // Сохраняем в кеш
    if s.cache.Enabled() {
        s.cache.Set(id.String(), shipment)
    }

    return shipment, nil
}

// TrackShipment с обработкой различных статусов провайдеров
func (s *DeliveryService) TrackShipment(ctx context.Context, trackingNumber string) (*TrackingInfo, error) {
    // Вызов gRPC
    tracking, err := s.client.TrackShipment(ctx, trackingNumber)
    if err != nil {
        return nil, fmt.Errorf("failed to track shipment: %w", err)
    }

    // Обогащение информации о трекинге
    enrichedTracking := s.enrichTracking(tracking)

    return enrichedTracking, nil
}

// CalculateRateWithFallback - расчет с fallback на mock если все провайдеры недоступны
func (s *DeliveryService) CalculateRateWithFallback(ctx context.Context, req *CalculateRateRequest) (*CalculateRateResponse, error) {
    // Валидация
    if err := s.validator.ValidateCalculateRateRequest(req); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }

    // Вызов gRPC
    rates, err := s.client.CalculateRate(ctx, &client.CalculateRateRequest{
        FromAddress: req.FromAddress,
        ToAddress:   req.ToAddress,
        Package:     req.Package,
        Type:        req.Type,
    })

    // Если все провайдеры недоступны - используем mock расчет
    if err != nil || len(rates.Options) == 0 {
        return s.calculateMockRate(req), nil
    }

    return rates, nil
}

// Приватные хелперы

func (s *DeliveryService) normalizeAddress(addr client.Address) client.Address {
    // Приведение к верхнему регистру, удаление лишних пробелов, etc.
    return client.Address{
        Street:     strings.TrimSpace(addr.Street),
        City:       strings.Title(strings.ToLower(addr.City)),
        PostalCode: strings.ReplaceAll(addr.PostalCode, " ", ""),
        Country:    strings.ToUpper(addr.Country),
        Phone:      s.normalizePhone(addr.Phone),
        Email:      strings.ToLower(strings.TrimSpace(addr.Email)),
        Name:       strings.TrimSpace(addr.Name),
    }
}

func (s *DeliveryService) normalizePhone(phone string) string {
    // Удаление всех нецифровых символов кроме +
    phone = strings.TrimSpace(phone)
    if !strings.HasPrefix(phone, "+") {
        phone = "+" + phone
    }
    return phone
}

func (s *DeliveryService) enrichShipment(shipment *client.Shipment, req *CreateShipmentRequest) *Shipment {
    // Добавление дополнительной информации
    return &Shipment{
        ID:                shipment.ID,
        TrackingNumber:    shipment.TrackingNumber,
        Status:            shipment.Status,
        Provider:          shipment.Provider,
        Cost:              shipment.Cost,
        Currency:          shipment.Currency,
        EstimatedDelivery: shipment.EstimatedDelivery,
        CreatedAt:         shipment.CreatedAt,
        // Дополнительные поля
        EstimatedDeliveryFormatted: s.formatDeliveryDate(shipment.EstimatedDelivery),
        CostFormatted:              s.formatCost(shipment.Cost, shipment.Currency),
        TrackingURL:                s.generateTrackingURL(shipment.Provider, shipment.TrackingNumber),
    }
}

func (s *DeliveryService) enrichTracking(tracking *client.TrackingInfo) *TrackingInfo {
    return &TrackingInfo{
        Shipment: tracking.Shipment,
        Events:   tracking.Events,
        // Дополнительные поля
        CurrentStep:       s.calculateCurrentStep(tracking.Shipment.Status),
        ProgressPercent:   s.calculateProgress(tracking.Shipment.Status),
        IsDelivered:       tracking.Shipment.Status == "delivered",
        CanBeCancelled:    s.canBeCancelled(tracking.Shipment.Status),
        EstimatedTimeLeft: s.calculateTimeLeft(tracking.Shipment.EstimatedDelivery),
    }
}

func (s *DeliveryService) formatDeliveryDate(t *time.Time) string {
    if t == nil {
        return "Неизвестно"
    }
    return t.Format("02.01.2006")
}

func (s *DeliveryService) formatCost(cost float64, currency string) string {
    return fmt.Sprintf("%.2f %s", cost, currency)
}

func (s *DeliveryService) generateTrackingURL(provider, trackingNumber string) string {
    urls := map[string]string{
        "post_express": "https://postexpress.rs/tracking?number=%s",
        "dex":          "https://dex.rs/track/%s",
    }

    if template, ok := urls[provider]; ok {
        return fmt.Sprintf(template, trackingNumber)
    }
    return ""
}

func (s *DeliveryService) calculateCurrentStep(status string) int {
    steps := map[string]int{
        "pending":           1,
        "confirmed":         2,
        "picked_up":         3,
        "in_transit":        4,
        "out_for_delivery":  5,
        "delivered":         6,
    }
    if step, ok := steps[status]; ok {
        return step
    }
    return 1
}

func (s *DeliveryService) calculateProgress(status string) int {
    step := s.calculateCurrentStep(status)
    return step * 100 / 6
}

func (s *DeliveryService) canBeCancelled(status string) bool {
    nonCancellable := []string{"delivered", "cancelled", "returned", "out_for_delivery"}
    for _, s := range nonCancellable {
        if status == s {
            return false
        }
    }
    return true
}

func (s *DeliveryService) calculateTimeLeft(estimated *time.Time) string {
    if estimated == nil {
        return ""
    }

    duration := time.Until(*estimated)
    if duration < 0 {
        return "Просрочено"
    }

    hours := int(duration.Hours())
    if hours < 24 {
        return fmt.Sprintf("%d часов", hours)
    }

    days := hours / 24
    return fmt.Sprintf("%d дней", days)
}

func (s *DeliveryService) calculateMockRate(req *CalculateRateRequest) *CalculateRateResponse {
    // Простой mock расчет на случай недоступности провайдеров
    baseRate := 500.0
    weightFactor := req.Package.WeightKg * 50

    return &CalculateRateResponse{
        Options: []RateOption{
            {
                Type:          "standard",
                Cost:          baseRate + weightFactor,
                Currency:      "RSD",
                EstimatedDays: 3,
            },
            {
                Type:          "express",
                Cost:          (baseRate + weightFactor) * 1.5,
                Currency:      "RSD",
                EstimatedDays: 1,
            },
        },
    }
}
```

**Файл**: `pkg/service/validator.go`

```go
package service

import (
    "fmt"
    "regexp"
)

type Validator struct {
    emailRegex *regexp.Regexp
    phoneRegex *regexp.Regexp
}

func NewValidator() *Validator {
    return &Validator{
        emailRegex: regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
        phoneRegex: regexp.MustCompile(`^\+?[1-9]\d{1,14}$`),
    }
}

func (v *Validator) ValidateCreateShipmentRequest(req *CreateShipmentRequest) error {
    if req.ProviderCode == "" {
        return fmt.Errorf("provider code is required")
    }

    if err := v.validateAddress(req.FromAddress, "from"); err != nil {
        return err
    }

    if err := v.validateAddress(req.ToAddress, "to"); err != nil {
        return err
    }

    if err := v.validatePackage(req.Package); err != nil {
        return err
    }

    return nil
}

func (v *Validator) validateAddress(addr client.Address, prefix string) error {
    if addr.Street == "" {
        return fmt.Errorf("%s address: street is required", prefix)
    }
    if addr.City == "" {
        return fmt.Errorf("%s address: city is required", prefix)
    }
    if addr.PostalCode == "" {
        return fmt.Errorf("%s address: postal code is required", prefix)
    }
    if addr.Country == "" {
        return fmt.Errorf("%s address: country is required", prefix)
    }
    if addr.Phone == "" {
        return fmt.Errorf("%s address: phone is required", prefix)
    }
    if !v.phoneRegex.MatchString(addr.Phone) {
        return fmt.Errorf("%s address: invalid phone format", prefix)
    }
    if addr.Email != "" && !v.emailRegex.MatchString(addr.Email) {
        return fmt.Errorf("%s address: invalid email format", prefix)
    }
    if addr.Name == "" {
        return fmt.Errorf("%s address: name is required", prefix)
    }
    return nil
}

func (v *Validator) validatePackage(pkg client.Package) error {
    if pkg.WeightKg <= 0 {
        return fmt.Errorf("package weight must be positive")
    }
    if pkg.WeightKg > 30 {
        return fmt.Errorf("package weight exceeds maximum (30kg)")
    }
    if pkg.LengthCm <= 0 || pkg.WidthCm <= 0 || pkg.HeightCm <= 0 {
        return fmt.Errorf("package dimensions must be positive")
    }
    if pkg.Description == "" {
        return fmt.Errorf("package description is required")
    }
    return nil
}
```

**Файл**: `pkg/service/retry.go`

```go
package service

import (
    "context"
    "time"
)

type Retrier struct {
    maxAttempts int
    timeout     time.Duration
}

func NewRetrier(maxAttempts int, timeout time.Duration) *Retrier {
    if maxAttempts <= 0 {
        maxAttempts = 3
    }
    if timeout <= 0 {
        timeout = 5 * time.Second
    }
    return &Retrier{
        maxAttempts: maxAttempts,
        timeout:     timeout,
    }
}

func (r *Retrier) Do(ctx context.Context, fn func() error) error {
    var lastErr error

    for attempt := 1; attempt <= r.maxAttempts; attempt++ {
        // Проверяем контекст
        if ctx.Err() != nil {
            return ctx.Err()
        }

        // Пытаемся выполнить
        lastErr = fn()
        if lastErr == nil {
            return nil
        }

        // Если это последняя попытка - возвращаем ошибку
        if attempt == r.maxAttempts {
            break
        }

        // Exponential backoff
        backoff := time.Duration(attempt) * r.timeout
        select {
        case <-time.After(backoff):
            continue
        case <-ctx.Done():
            return ctx.Err()
        }
    }

    return fmt.Errorf("max retry attempts reached: %w", lastErr)
}
```

**Использование в монолите**:

```go
// backend/internal/proj/delivery/module.go

import (
    deliveryService "github.com/sveturs/delivery/pkg/service"
)

func NewModule(cfg *config.Config) (*Module, error) {
    // Используем высокоуровневый сервис вместо низкоуровневого клиента
    service, err := deliveryService.NewDeliveryService(&deliveryService.Config{
        GRPCAddress:   cfg.DeliveryServiceAddress,
        RetryAttempts: 3,
        RetryTimeout:  5 * time.Second,
        CacheEnabled:  true,
        CacheTTL:      5 * time.Minute,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to create delivery service: %w", err)
    }

    handler := NewHandler(service)

    return &Module{
        service: service,
        handler: handler,
    }, nil
}

// backend/internal/proj/delivery/handler.go

func (h *Handler) CreateShipment(c *fiber.Ctx) error {
    var req CreateShipmentRequest
    if err := c.BodyParser(&req); err != nil {
        return utils.SendErrorResponse(c, fiber.StatusBadRequest, "error.invalid_request", nil)
    }

    userID, _ := authmiddleware.GetUserID(c)

    // Вызов высокоуровневого сервиса (с валидацией, retry, обогащением)
    shipment, err := h.service.CreateShipment(c.Context(), &deliveryService.CreateShipmentRequest{
        ProviderCode: req.ProviderCode,
        UserID:       uuid.MustParse(userID),
        FromAddress:  req.FromAddress,
        ToAddress:    req.ToAddress,
        Package:      req.Package,
        Type:         req.DeliveryType,
    })

    if err != nil {
        return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "error.failed_to_create_shipment", nil)
    }

    // Shipment уже обогащен дополнительными полями
    return utils.SendSuccessResponse(c, shipment, "Отправление создано")
}
```

**Преимущества pkg/service обертки**:

1. ✅ **Валидация** - проверка данных перед отправкой в микросервис
2. ✅ **Retry** - автоматические повторы при временных ошибках
3. ✅ **Нормализация** - приведение данных к единому формату
4. ✅ **Обогащение** - добавление вычисляемых полей
5. ✅ **Кеширование** - снижение нагрузки на микросервис
6. ✅ **Fallback** - mock данные при недоступности провайдеров
7. ✅ **Удобный API** - высокоуровневые методы вместо protobuf
8. ✅ **Централизация** - вся бизнес-логика в одном месте

---

### ФАЗА 2: Тестирование (Week 3)

#### 2.1 Unit Tests

**Файл**: `internal/service/delivery_service_test.go`

```go
func TestDeliveryService_CreateShipment(t *testing.T) {
    // Mock repository
    mockRepo := &MockShipmentRepository{}
    mockEventRepo := &MockTrackingEventRepository{}

    // Mock provider
    mockProvider := &MockProvider{
        CreateShipmentFunc: func(ctx, req) (*provider.ShipmentResponse, error) {
            return &provider.ShipmentResponse{
                TrackingNumber: "TRACK123",
                Cost: provider.Money{Amount: 500, Currency: "RSD"},
            }, nil
        },
    }

    factory := &MockFactory{
        GetProviderFunc: func(code string) (provider.Provider, error) {
            return mockProvider, nil
        },
    }

    service := NewDeliveryService(mockRepo, mockEventRepo, factory, logger)

    // Test
    shipment, err := service.CreateShipment(context.Background(), &CreateShipmentInput{
        ProviderCode: "mock",
        // ... остальные поля
    })

    assert.NoError(t, err)
    assert.NotNil(t, shipment)
    assert.Equal(t, "TRACK123", shipment.TrackingNumber)
}
```

**Запуск**:
```bash
make test-unit
```

**Coverage target**: > 80%

#### 2.2 Integration Tests (с testcontainers)

**Файл**: `tests/integration/delivery_test.go`

```go
func TestDeliveryIntegration(t *testing.T) {
    // Запуск PostgreSQL через testcontainers
    ctx := context.Background()
    postgresContainer, err := postgres.RunContainer(ctx,
        testcontainers.WithImage("postgres:17-alpine"),
        postgres.WithDatabase("delivery_test"),
        postgres.WithUsername("test"),
        postgres.WithPassword("test"),
    )
    require.NoError(t, err)
    defer postgresContainer.Terminate(ctx)

    // Подключение к БД
    connStr, _ := postgresContainer.ConnectionString(ctx)
    db, err := sql.Open("postgres", connStr)
    require.NoError(t, err)

    // Миграции
    migrator := migrator.NewMigrator(db, "../../migrations")
    require.NoError(t, migrator.Run())

    // Инициализация сервисов
    repo := repository.NewPostgresShipmentRepository(db)
    factory := provider.NewFactory(config)
    service := service.NewDeliveryService(repo, eventRepo, factory, logger)

    // Test: Создание отправления
    t.Run("CreateShipment", func(t *testing.T) {
        shipment, err := service.CreateShipment(ctx, &service.CreateShipmentInput{
            ProviderCode: "mock",
            // ...
        })

        assert.NoError(t, err)
        assert.NotEmpty(t, shipment.ID)

        // Проверка что сохранилось в БД
        saved, err := repo.GetByID(ctx, shipment.ID)
        assert.NoError(t, err)
        assert.Equal(t, shipment.TrackingNumber, saved.TrackingNumber)
    })

    // Test: Трекинг
    t.Run("TrackShipment", func(t *testing.T) {
        // ...
    })
}
```

**Запуск**:
```bash
make test-integration
```

#### 2.3 gRPC Client Test

**Файл**: `tests/grpc_client_test.go`

```go
func TestGRPCClient(t *testing.T) {
    // Подключение к локальному gRPC серверу
    client, err := client.NewClient("localhost:50052")
    require.NoError(t, err)
    defer client.Close()

    ctx := context.Background()

    t.Run("CreateShipment", func(t *testing.T) {
        shipment, err := client.CreateShipment(ctx, &client.CreateShipmentRequest{
            Provider: "mock",
            UserID:   uuid.New(),
            FromAddress: client.Address{
                Street:     "Test Street 1",
                City:       "Belgrade",
                PostalCode: "11000",
                Country:    "RS",
                Phone:      "+381641234567",
                Email:      "sender@test.com",
                Name:       "Test Sender",
            },
            ToAddress: client.Address{
                Street:     "Test Street 2",
                City:       "Novi Sad",
                PostalCode: "21000",
                Country:    "RS",
                Phone:      "+381651234567",
                Email:      "receiver@test.com",
                Name:       "Test Receiver",
            },
            Package: client.Package{
                WeightKg:    2.5,
                LengthCm:    30,
                WidthCm:     20,
                HeightCm:    15,
                Description: "Test package",
                Value:       5000,
            },
            Type: "standard",
        })

        assert.NoError(t, err)
        assert.NotEmpty(t, shipment.TrackingNumber)

        // Проверка трекинга
        tracking, err := client.TrackShipment(ctx, shipment.TrackingNumber)
        assert.NoError(t, err)
        assert.Equal(t, shipment.ID, tracking.Shipment.ID)
    })
}
```

#### 2.4 Локальный запуск

```bash
# 1. Запуск PostgreSQL
cd ~/delivery
docker-compose up -d

# 2. Применение миграций
make migrate-up

# 3. Конфигурация
export SVETUDELIVERY_GATEWAYS_POSTRS_ENABLED=true
export SVETUDELIVERY_GATEWAYS_POSTRS_API_KEY="your-key"
export SVETUDELIVERY_GATEWAYS_POSTRS_BASE_URL="https://api.postexpress.rs"

# 4. Запуск микросервиса
make run

# 5. Проверка через grpcurl
grpcurl -plaintext -d '{
  "provider": "PROVIDER_POST_EXPRESS",
  "user_id": "00000000-0000-0000-0000-000000000000",
  "from_address": {
    "street": "Bulevar kralja Aleksandra 73",
    "city": "Beograd",
    "postal_code": "11000",
    "country": "RS",
    "phone": "+381641234567",
    "email": "sender@test.com",
    "name": "Test Sender"
  },
  "to_address": {
    "street": "Bulevar oslobođenja 46",
    "city": "Novi Sad",
    "postal_code": "21000",
    "country": "RS",
    "phone": "+381651234567",
    "email": "receiver@test.com",
    "name": "Test Receiver"
  },
  "package": {
    "weight_kg": 2.5,
    "length_cm": 30,
    "width_cm": 20,
    "height_cm": 15,
    "description": "Test package",
    "value": 5000
  },
  "type": "standard"
}' localhost:50052 delivery.v1.DeliveryService/CreateShipment
```

---

### ФАЗА 3: Переход монолита на микросервис (Week 4)

#### 3.1 Удаление старого кода

```bash
cd /data/hostel-booking-system/backend

# Удаляем всю старую реализацию delivery
rm -rf internal/proj/delivery/

# Создаем новую директорию только для gRPC клиента
mkdir -p internal/proj/delivery
```

#### 3.2 Интеграция gRPC клиента в монолит

**Файл**: `backend/go.mod` (обновить)

```go
require (
    github.com/sveturs/delivery v1.0.0
    // ... остальные зависимости
)
```

**Файл**: `backend/internal/proj/delivery/client.go`

```go
package delivery

import (
    deliveryClient "github.com/sveturs/delivery/pkg/client"
)

type Client struct {
    grpc *deliveryClient.Client
}

func NewClient(addr string) (*Client, error) {
    grpcClient, err := deliveryClient.NewClient(addr)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to delivery service: %w", err)
    }

    return &Client{grpc: grpcClient}, nil
}

func (c *Client) Close() error {
    return c.grpc.Close()
}
```

**Файл**: `backend/internal/proj/delivery/handler.go`

```go
package delivery

import (
    "github.com/gofiber/fiber/v2"
    "backend/pkg/utils"
)

type Handler struct {
    client *Client
}

func NewHandler(client *Client) *Handler {
    return &Handler{client: client}
}

func (h *Handler) CreateShipment(c *fiber.Ctx) error {
    var req CreateShipmentRequest
    if err := c.BodyParser(&req); err != nil {
        return utils.SendErrorResponse(c, fiber.StatusBadRequest, "error.invalid_request", nil)
    }

    // Получаем user_id из JWT
    userID, _ := authmiddleware.GetUserID(c)

    // Вызов микросервиса
    shipment, err := h.client.grpc.CreateShipment(c.Context(), &deliveryClient.CreateShipmentRequest{
        Provider:    req.ProviderCode,
        UserID:      uuid.MustParse(userID),
        FromAddress: req.FromAddress,
        ToAddress:   req.ToAddress,
        Package:     req.Package,
        Type:        req.DeliveryType,
    })

    if err != nil {
        return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "error.failed_to_create_shipment", nil)
    }

    return utils.SendSuccessResponse(c, shipment, "Отправление создано")
}

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

-- Удаляем старые таблицы delivery
DROP TABLE IF EXISTS delivery_tracking_events CASCADE;
DROP TABLE IF EXISTS delivery_shipments CASCADE;
DROP TABLE IF EXISTS delivery_providers CASCADE;
DROP TABLE IF EXISTS delivery_pricing_rules CASCADE;
DROP TABLE IF EXISTS delivery_zones CASCADE;
DROP TABLE IF EXISTS delivery_category_defaults CASCADE;
```

---

## ✅ Критерии готовности

### После Фазы 1:
- [x] Микросервис запускается локально
- [x] Proto код сгенерирован
- [x] Все слои реализованы (domain, repo, service, gateway, grpc)
- [x] Post Express интеграция работает
- [x] Библиотека pkg/client готова

### После Фазы 2:
- [x] Unit tests coverage > 80%
- [x] Integration tests проходят
- [x] gRPC client test работает
- [x] Микросервис протестирован на dev окружении

### После Фазы 3:
- [x] Монолит использует gRPC клиент
- [x] Старый код удален
- [x] Данные мигрированы (если были)
- [x] Frontend работает без изменений
- [x] Все API endpoints работают через микросервис
- [x] Старые таблицы удалены из монолита

---

## 🚀 Production Deployment

### После успешного тестирования на dev:

**Неделя 5**: Развертывание на production

```bash
# 1. Deploy микросервиса
kubectl apply -f k8s/delivery-service/

# 2. Миграция данных на production
# Выполнить migrate_delivery_data.sql на production БД

# 3. Deploy обновленного монолита
kubectl apply -f k8s/backend/

# 4. Smoke tests
./scripts/smoke_test_delivery.sh

# 5. Мониторинг метрик
# Открыть Grafana → Delivery Service Dashboard

# 6. Удаление старых таблиц (через неделю)
```

**Kubernetes manifest**: `k8s/delivery-service/deployment.yaml`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: delivery-service
  namespace: production
spec:
  replicas: 3
  selector:
    matchLabels:
      app: delivery-service
  template:
    metadata:
      labels:
        app: delivery-service
    spec:
      containers:
      - name: delivery
        image: registry.svetu.rs/delivery:v1.0.0
        ports:
        - containerPort: 50052
          name: grpc
        - containerPort: 9091
          name: metrics
        env:
        - name: SVETUDELIVERY_DATABASE_HOST
          valueFrom:
            secretKeyRef:
              name: delivery-db-secret
              key: host
        - name: SVETUDELIVERY_DATABASE_PASSWORD
          valueFrom:
            secretKeyRef:
              name: delivery-db-secret
              key: password
        - name: SVETUDELIVERY_GATEWAYS_POSTRS_API_KEY
          valueFrom:
            secretKeyRef:
              name: postexpress-secret
              key: api-key
        livenessProbe:
          grpc:
            port: 50052
          initialDelaySeconds: 10
          periodSeconds: 10
        readinessProbe:
          grpc:
            port: 50052
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "256Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
```

---

## 📊 Мониторинг

### Prometheus Metrics

**Файл**: `internal/server/grpc/metrics.go`

```go
var (
    grpcRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "delivery_grpc_requests_total",
            Help: "Total number of gRPC requests",
        },
        []string{"method", "status"},
    )

    grpcRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "delivery_grpc_request_duration_seconds",
            Help:    "Duration of gRPC requests",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method"},
    )

    shipmentsCreatedTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "delivery_shipments_created_total",
            Help: "Total number of shipments created",
        },
        []string{"provider"},
    )
)
```

### Grafana Dashboard

**Панели**:
- Request rate (RPS)
- Request latency (p50, p95, p99)
- Error rate
- Shipments created by provider
- Active shipments by status

---

## 🔄 Rollback Plan

Если что-то пойдет не так:

```bash
# 1. Остановка микросервиса
docker-compose -f docker-compose.dev.yml stop delivery-service

# 2. Откат монолита к предыдущей версии
git checkout HEAD~1
docker-compose -f docker-compose.dev.yml up -d --build backend

# 3. Проверка
./scripts/smoke_test.sh
```

**На production**:
```bash
kubectl rollout undo deployment/delivery-service -n production
kubectl rollout undo deployment/backend -n production
```

---

## 📝 Чеклист

### Фаза 1 (Week 1-2):
- [ ] Proto код сгенерирован
- [ ] Domain models созданы
- [ ] Repository реализован
- [ ] Provider factory создан
- [ ] Post Express интеграция перенесена
- [ ] Service layer реализован
- [ ] gRPC handlers реализованы
- [ ] pkg/client библиотека готова
- [ ] Микросервис запускается локально

### Фаза 2 (Week 3):
- [ ] Unit tests написаны (coverage > 80%)
- [ ] Integration tests написаны
- [ ] gRPC client test работает
- [ ] Локальное тестирование пройдено
- [ ] Все API endpoints работают

### Фаза 3 (Week 4):
- [ ] Старый код удален из монолита
- [ ] gRPC клиент интегрирован
- [ ] Handlers обновлены
- [ ] Routes обновлены
- [ ] Config обновлен
- [ ] Docker Compose обновлен
- [ ] Данные мигрированы
- [ ] Deploy на dev выполнен
- [ ] Frontend тестирование пройдено
- [ ] Старые таблицы удалены

### Production:
- [ ] Kubernetes manifests готовы
- [ ] Secrets настроены
- [ ] Мониторинг настроен
- [ ] Алерты настроены
- [ ] Runbook создан
- [ ] Deploy на production выполнен
- [ ] Smoke tests пройдены
- [ ] Метрики в норме

---

## 🎯 Итого

**Подход**: Clean Cut - полный переход без промежуточных состояний
**Срок**: 3-4 недели
**Результат**: Независимый микросервис, монолит использует gRPC клиент
**Обратная совместимость**: НЕ требуется
**Feature flags**: НЕ нужны
**Canary deployment**: НЕ нужен

**Принцип**: Делаем правильно с первого раза!

---

**Дата создания**: 2025-10-22
**Автор**: Claude Code
**Статус**: Ready for implementation

---

## 🚀 Инфраструктура: Развертывание на svetu.rs

> **Источник данных**: Реальный анализ сервера svetu.rs (2025-10-22)
> **Метод**: SSH анализ через Claude Code с полными правами доступа

### 📊 Текущая архитектура сервера

**Существующие окружения**:
```
/opt/
├── svetu-authpreprod/     # Auth микросервис (Go + gRPC)
├── svetu-dev/             # Dev окружение (монолит)
└── svetu-preprod/         # Preprod окружение (монолит)
```

**Паттерн развертывания**: Docker Compose с изолированными сервисами

### 🔌 Распределение портов

**Занятые порты по окружениям**:

| Окружение | PostgreSQL | Redis | OpenSearch | HTTP | gRPC | Metrics | Health |
|-----------|------------|-------|------------|------|------|---------|--------|
| **svetu-dev** | 5433 | 6380 | 9201 | - | - | - | - |
| **svetu-preprod** | 5489 | 6382 | 9203 | 3012 | - | - | - |
| **svetu-authpreprod** | 25432 | 26379 | - | 28080 | **20051** | 29090 | 28081 |

**Свободные gRPC порты** (диапазон 50050-50060):
- ✅ `50050, 50052, 50053, 54, 55, 56, 57, 58, 59, 60`
- ❌ `50051` (занят auth-service)

**Рекомендуемые порты для delivery-preprod**:

| Сервис | Порт | Назначение |
|--------|------|------------|
| PostgreSQL | `35432` | База данных delivery |
| Redis | `36379` | Кэш и очереди |
| HTTP API | `38080` | REST API (если нужен) |
| **gRPC API** | `30051` | **Основной gRPC сервис** |
| Health Check | `38081` | Healthcheck endpoint |
| Metrics | `39090` | Prometheus metrics |

> **Примечание**: Порты в диапазоне 30000-39999 выбраны для избежания конфликтов

### 📂 Структура директории (по образцу auth-service)

```
/opt/svetu-delivery-preprod/
├── cmd/
│   └── server/              # Точка входа gRPC сервера
│       └── main.go
├── internal/
│   ├── app/                 # Инициализация приложения
│   ├── transport/           # gRPC handlers
│   │   └── grpc/
│   ├── domain/              # Доменные модели
│   ├── repository/          # PostgreSQL repos
│   │   └── postgres/
│   ├── service/             # Бизнес-логика
│   │   ├── delivery.go
│   │   ├── calculator.go
│   │   └── tracking.go
│   ├── gateway/             # Интеграции с внешними API
│   │   └── provider/
│   │       ├── interface.go
│   │       ├── factory.go
│   │       ├── postexpress/
│   │       ├── dex/
│   │       └── mock/
│   └── config/              # Конфигурация
├── pkg/                     # Публичные библиотеки
│   ├── client/              # gRPC клиент для монолита
│   └── service/             # Высокоуровневая обертка
├── deployments/
│   └── docker/
│       └── Dockerfile
├── migrations/              # SQL миграции
├── fixtures/                # Тестовые данные
├── nginx/                   # Nginx конфигурация
├── .env                     # Переменные окружения
├── .env.example             # Шаблон .env
├── docker-compose.yml       # Для локальной разработки
└── docker-compose.preprod.yml  # Для production
```

### 🐳 Docker Compose конфигурация

**Файл**: `/opt/svetu-delivery-preprod/docker-compose.preprod.yml`

```yaml
version: '3.8'

volumes:
  svetudelivery_postgres_data:
    driver: local
  svetudelivery_redis_data:
    driver: local

networks:
  svetudelivery-network:
    driver: bridge

services:
  delivery-postgres:
    image: postgres:15-alpine
    container_name: svetudelivery-postgres
    environment:
      POSTGRES_DB: ${SVETUDELIVERY_DB_NAME:-delivery_db}
      POSTGRES_USER: ${SVETUDELIVERY_DB_USER:-delivery_user}
      POSTGRES_PASSWORD: ${SVETUDELIVERY_DB_PASSWORD}
      POSTGRES_INITDB_ARGS: "--encoding=UTF8 --lc-collate=C --lc-ctype=C"
    volumes:
      - svetudelivery_postgres_data:/var/lib/postgresql/data
    ports:
      - "35432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${SVETUDELIVERY_DB_USER:-delivery_user}"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped
    networks:
      - svetudelivery-network

  delivery-redis:
    image: redis:7-alpine
    container_name: svetudelivery-redis
    command: >
      redis-server
      --requirepass ${SVETUDELIVERY_REDIS_PASSWORD}
      --maxmemory 512mb
      --maxmemory-policy allkeys-lru
      --save 900 1
      --save 300 10
      --save 60 10000
    volumes:
      - svetudelivery_redis_data:/data
    ports:
      - "36379:6379"
    healthcheck:
      test: ["CMD", "redis-cli", "--no-auth-warning", "-a", "${SVETUDELIVERY_REDIS_PASSWORD}", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped
    networks:
      - svetudelivery-network

  delivery-service:
    build:
      context: .
      dockerfile: deployments/docker/Dockerfile
      args:
        GO_VERSION: "1.23"
    container_name: svetudelivery-service
    environment:
      # Service
      SVETUDELIVERY_SERVICE_NAME: ${SVETUDELIVERY_SERVICE_NAME:-delivery-service}
      SVETUDELIVERY_SERVICE_ENV: ${SVETUDELIVERY_SERVICE_ENV:-preprod}
      SVETUDELIVERY_SERVICE_LOG_LEVEL: ${SVETUDELIVERY_LOG_LEVEL:-info}

      # Server
      SVETUDELIVERY_SERVER_HTTP_PORT: 8080
      SVETUDELIVERY_SERVER_GRPC_PORT: 50052
      SVETUDELIVERY_SERVER_HEALTH_PORT: 8081
      SVETUDELIVERY_SERVER_METRICS_PORT: 9090

      # Database
      SVETUDELIVERY_DB_HOST: delivery-postgres
      SVETUDELIVERY_DB_PORT: 5432
      SVETUDELIVERY_DB_NAME: ${SVETUDELIVERY_DB_NAME:-delivery_db}
      SVETUDELIVERY_DB_USER: ${SVETUDELIVERY_DB_USER:-delivery_user}
      SVETUDELIVERY_DB_PASSWORD: ${SVETUDELIVERY_DB_PASSWORD}
      SVETUDELIVERY_DB_SSLMODE: disable

      # Redis
      SVETUDELIVERY_REDIS_HOST: delivery-redis
      SVETUDELIVERY_REDIS_PORT: 6379
      SVETUDELIVERY_REDIS_PASSWORD: ${SVETUDELIVERY_REDIS_PASSWORD}
      SVETUDELIVERY_REDIS_DB: 0

      # External APIs
      SVETUDELIVERY_POSTEXPRESS_API_URL: ${SVETUDELIVERY_POSTEXPRESS_API_URL:-https://api.postexpress.rs}
      SVETUDELIVERY_POSTEXPRESS_API_KEY: ${SVETUDELIVERY_POSTEXPRESS_API_KEY}
    ports:
      - "38080:8080"    # HTTP API (опционально)
      - "30051:50052"   # gRPC API (ОСНОВНОЙ!)
      - "38081:8081"    # Health Check
      - "39090:9090"    # Prometheus Metrics
    depends_on:
      delivery-postgres:
        condition: service_healthy
      delivery-redis:
        condition: service_healthy
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://127.0.0.1:8081/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
    restart: unless-stopped
    networks:
      - svetudelivery-network
```

### 🔐 Переменные окружения (.env)

**Файл**: `/opt/svetu-delivery-preprod/.env`

```bash
# Service Configuration
SVETUDELIVERY_SERVICE_NAME=delivery-service
SVETUDELIVERY_SERVICE_ENV=preprod
SVETUDELIVERY_LOG_LEVEL=info

# Database Configuration
SVETUDELIVERY_DB_NAME=delivery_db
SVETUDELIVERY_DB_USER=delivery_user
SVETUDELIVERY_DB_PASSWORD=GENERATE_STRONG_PASSWORD_HERE

# Redis Configuration
SVETUDELIVERY_REDIS_PASSWORD=GENERATE_STRONG_PASSWORD_HERE

# External APIs
SVETUDELIVERY_POSTEXPRESS_API_KEY=YOUR_POST_EXPRESS_API_KEY
SVETUDELIVERY_POSTEXPRESS_API_URL=https://api.postexpress.rs

# Monitoring (optional)
SVETUDELIVERY_PROMETHEUS_ENABLED=true
SVETUDELIVERY_JAEGER_ENABLED=false
```

### 🌐 Nginx конфигурация

**Файл**: `/etc/nginx/sites-available/deliverypreprod.svetu.rs`

```nginx
# Upstream для delivery gRPC service
upstream delivery_grpc_backend {
    server 127.0.0.1:30051;
    keepalive 32;
}

# HTTP/2 для gRPC (требуется SSL)
server {
    listen 443 ssl http2;
    server_name deliverypreprod.svetu.rs;

    # SSL сертификаты (Let's Encrypt)
    ssl_certificate /etc/letsencrypt/live/deliverypreprod.svetu.rs/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/deliverypreprod.svetu.rs/privkey.pem;

    # SSL оптимизация
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;

    # gRPC специфичные настройки
    grpc_read_timeout 300s;
    grpc_send_timeout 300s;
    client_body_timeout 300s;

    # Логирование
    access_log /var/log/nginx/deliverypreprod.access.log;
    error_log /var/log/nginx/deliverypreprod.error.log;

    # gRPC endpoint
    location / {
        grpc_pass grpc://delivery_grpc_backend;

        # Headers
        grpc_set_header Host $host;
        grpc_set_header X-Real-IP $remote_addr;
        grpc_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        grpc_set_header X-Forwarded-Proto $scheme;

        # Error handling
        error_page 502 = /error502grpc;
        error_page 503 = /error503grpc;
        error_page 504 = /error504grpc;
    }

    # Health check (HTTP, не gRPC)
    location /health {
        proxy_pass http://127.0.0.1:38081/health;
        access_log off;
    }

    # Metrics (HTTP, не gRPC) - для внутреннего использования
    location /metrics {
        proxy_pass http://127.0.0.1:39090/metrics;
        allow 127.0.0.1;
        deny all;
    }

    # gRPC error responses
    location = /error502grpc {
        internal;
        default_type application/grpc;
        add_header grpc-status 14;  # UNAVAILABLE
        add_header grpc-message "Bad Gateway";
        return 204;
    }

    location = /error503grpc {
        internal;
        default_type application/grpc;
        add_header grpc-status 14;  # UNAVAILABLE
        add_header grpc-message "Service Temporarily Unavailable";
        return 204;
    }

    location = /error504grpc {
        internal;
        default_type application/grpc;
        add_header grpc-status 4;   # DEADLINE_EXCEEDED
        add_header grpc-message "Gateway Timeout";
        return 204;
    }
}

# HTTP redirect to HTTPS
server {
    listen 80;
    server_name deliverypreprod.svetu.rs;
    return 301 https://$server_name$request_uri;
}
```

### 📝 Пошаговая инструкция развертывания

#### 1. Подготовка сервера

```bash
# SSH на сервер
ssh svetu@svetu.rs

# Создание директории
sudo mkdir -p /opt/svetu-delivery-preprod
sudo chown svetu:svetu /opt/svetu-delivery-preprod
cd /opt/svetu-delivery-preprod

# Клонирование репозитория
git clone git@github.com:sveturs/delivery.git .
git checkout main
```

#### 2. Настройка переменных окружения

```bash
# Копирование шаблона
cp .env.example .env

# Генерация паролей
DB_PASSWORD=$(openssl rand -base64 32)
REDIS_PASSWORD=$(openssl rand -base64 32)

# Обновление .env
sed -i "s/SVETUDELIVERY_DB_PASSWORD=.*/SVETUDELIVERY_DB_PASSWORD=$DB_PASSWORD/" .env
sed -i "s/SVETUDELIVERY_REDIS_PASSWORD=.*/SVETUDELIVERY_REDIS_PASSWORD=$REDIS_PASSWORD/" .env

# Добавление API ключей вручную
nano .env
```

#### 3. Запуск Docker Compose

```bash
# Сборка образа
docker-compose -f docker-compose.preprod.yml build

# Запуск сервисов
docker-compose -f docker-compose.preprod.yml up -d

# Проверка статуса
docker-compose -f docker-compose.preprod.yml ps

# Логи
docker-compose -f docker-compose.preprod.yml logs -f delivery-service
```

#### 4. Применение миграций

```bash
# Подключение к контейнеру
docker exec -it svetudelivery-service sh

# Применение миграций (из контейнера)
/app/migrator up

# Или через docker exec
docker exec svetudelivery-service /app/migrator up
```

#### 5. Настройка Nginx

```bash
# Копирование конфигурации
sudo cp nginx/deliverypreprod.svetu.rs.conf /etc/nginx/sites-available/
sudo ln -s /etc/nginx/sites-available/deliverypreprod.svetu.rs.conf /etc/nginx/sites-enabled/

# Получение SSL сертификата
sudo certbot certonly --nginx -d deliverypreprod.svetu.rs

# Проверка конфигурации
sudo nginx -t

# Перезагрузка Nginx
sudo systemctl reload nginx
```

#### 6. Проверка работоспособности

```bash
# Health check
curl http://localhost:38081/health

# Metrics
curl http://localhost:39090/metrics

# gRPC endpoint (через grpcurl)
grpcurl -plaintext localhost:30051 list
grpcurl -plaintext localhost:30051 delivery.v1.DeliveryService/GetShipment
```

#### 7. Настройка автозапуска

```bash
# Создание systemd service
sudo nano /etc/systemd/system/delivery-preprod.service
```

**Содержимое**:
```ini
[Unit]
Description=Delivery Microservice (Preprod)
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/opt/svetu-delivery-preprod
ExecStart=/usr/bin/docker-compose -f docker-compose.preprod.yml up -d
ExecStop=/usr/bin/docker-compose -f docker-compose.preprod.yml down
TimeoutStartSec=300

[Install]
WantedBy=multi-user.target
```

```bash
# Активация
sudo systemctl daemon-reload
sudo systemctl enable delivery-preprod.service
sudo systemctl start delivery-preprod.service
```

### 🔍 Мониторинг и отладка

#### Логи

```bash
# Все сервисы
docker-compose -f docker-compose.preprod.yml logs -f

# Только delivery-service
docker-compose -f docker-compose.preprod.yml logs -f delivery-service

# PostgreSQL
docker-compose -f docker-compose.preprod.yml logs -f delivery-postgres

# Redis
docker-compose -f docker-compose.preprod.yml logs -f delivery-redis
```

#### Проверка портов

```bash
# Занятые порты
sudo netstat -tlnp | grep -E "30051|35432|36379|38080|38081|39090"

# Процессы Docker
docker ps | grep svetudelivery
```

#### Подключение к базе данных

```bash
# Из хоста
psql "postgres://delivery_user:PASSWORD@localhost:35432/delivery_db"

# Или через docker exec
docker exec -it svetudelivery-postgres psql -U delivery_user -d delivery_db
```

#### Проверка Redis

```bash
# Ping
docker exec svetudelivery-redis redis-cli -a PASSWORD ping

# Мониторинг команд
docker exec svetudelivery-redis redis-cli -a PASSWORD monitor
```

### 🚨 Troubleshooting

#### Проблема: Порт 30051 занят

```bash
# Найти процесс
sudo lsof -i :30051

# Остановить конфликтующий сервис
docker-compose -f /opt/OTHER_SERVICE/docker-compose.yml stop
```

#### Проблема: БД не поднимается

```bash
# Проверка логов
docker logs svetudelivery-postgres

# Проверка прав доступа
docker exec svetudelivery-postgres ls -la /var/lib/postgresql/data

# Пересоздание volume
docker-compose -f docker-compose.preprod.yml down -v
docker-compose -f docker-compose.preprod.yml up -d
```

#### Проблема: gRPC недоступен

```bash
# Проверка Nginx конфигурации
sudo nginx -t

# Проверка SSL сертификата
sudo certbot certificates

# Проверка firewall
sudo ufw status
```

### 📊 Ресурсы сервера

**Текущее состояние** (2025-10-22):
- **Диск**: 22GB свободно из 193GB (90% использовано)
- **Docker**: версия 27.5.1
- **Go**: версия 1.25.0

**Рекомендации**:
1. ⚠️ Мониторить место на диске (осталось мало!)
2. Настроить ротацию логов Docker
3. Очистить старые образы: `docker system prune -a`

### 🔄 Интеграция с монолитом

После развертывания микросервиса, монолит будет обращаться к нему через:

**gRPC адрес (внутренний)**: `localhost:30051`
**gRPC адрес (внешний)**: `deliverypreprod.svetu.rs:443`

**Конфигурация в монолите** (`backend/internal/config/config.go`):
```go
type DeliveryConfig struct {
    GRPCAddress string `env:"DELIVERY_GRPC_ADDRESS" envDefault:"localhost:30051"`
    UseTLS      bool   `env:"DELIVERY_USE_TLS" envDefault:"false"`
}
```

**Для preprod окружения**:
```bash
# В .env монолита
DELIVERY_GRPC_ADDRESS=localhost:30051
DELIVERY_USE_TLS=false
```

**Для production**:
```bash
DELIVERY_GRPC_ADDRESS=deliverypreprod.svetu.rs:443
DELIVERY_USE_TLS=true
```

---

## 📋 Обновленный чеклист с учетом инфраструктуры

### Фаза 0: Подготовка инфраструктуры (Week 0)
- [ ] Создать директорию `/opt/svetu-delivery-preprod`
- [ ] Получить SSL сертификат для `deliverypreprod.svetu.rs`
- [ ] Настроить Nginx конфигурацию
- [ ] Проверить свободные порты (30051, 35432, 36379, 38080-81, 39090)
- [ ] Сгенерировать пароли для БД и Redis
- [ ] Создать `.env` файл с конфигурацией
- [ ] Настроить systemd service для автозапуска

### Фаза 1: Разработка (Week 1-2)
- [ ] Proto код сгенерирован
- [ ] Domain models созданы
- [ ] Repository реализован
- [ ] Provider factory создан
- [ ] Post Express интеграция перенесена
- [ ] Service layer реализован
- [ ] gRPC handlers реализованы
- [ ] pkg/client библиотека готова
- [ ] Dockerfile создан
- [ ] docker-compose.preprod.yml настроен
- [ ] Микросервис запускается локально

### Фаза 2: Тестирование (Week 3)
- [ ] Unit tests написаны (coverage > 80%)
- [ ] Integration tests написаны
- [ ] gRPC client test работает
- [ ] Локальное тестирование пройдено
- [ ] Health checks работают
- [ ] Metrics endpoint функционирует
- [ ] Docker образ собирается успешно

### Фаза 3: Развертывание (Week 4)
- [ ] Код выгружен на сервер `/opt/svetu-delivery-preprod`
- [ ] Docker Compose запущен
- [ ] Миграции применены
- [ ] Nginx перезагружен
- [ ] Health check доступен: `curl http://localhost:38081/health`
- [ ] gRPC доступен: `grpcurl localhost:30051 list`
- [ ] Metrics доступны: `curl http://localhost:39090/metrics`
- [ ] SSL работает: `curl https://deliverypreprod.svetu.rs/health`
- [ ] Systemd service активирован
- [ ] Мониторинг логов настроен

### Фаза 4: Миграция монолита (Week 4-5)
- [ ] Старый код удален из монолита
- [ ] gRPC клиент интегрирован
- [ ] Handlers обновлены (proxy в микросервис)
- [ ] Routes обновлены
- [ ] Config обновлен (DELIVERY_GRPC_ADDRESS)
- [ ] Монолит перезапущен
- [ ] Интеграционное тестирование пройдено
- [ ] Frontend работает без изменений
- [ ] Старые таблицы удалены из БД монолита

### Фаза 5: Финализация (Week 5)
- [ ] Документация обновлена
- [ ] Runbook создан
- [ ] Smoke tests пройдены
- [ ] Метрики в Prometheus настроены
- [ ] Grafana dashboard создан
- [ ] Алерты настроены
- [ ] Резервное копирование БД настроено

---

**Обновлено**: 2025-10-22 (добавлена реальная инфраструктура svetu.rs)
