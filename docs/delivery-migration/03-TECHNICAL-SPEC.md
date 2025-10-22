# Техническая спецификация микросервиса
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
