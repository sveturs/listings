# Фаза 1: Разработка микросервиса

**Срок**: Week 1-2

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
