# Фаза 2: Тестирование

**Срок**: Week 3

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
