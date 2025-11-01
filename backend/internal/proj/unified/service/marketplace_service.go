// backend/internal/proj/unified/service/marketplace_service.go
package service

import (
	"context"
	"fmt"
	"time"

	"backend/internal/domain/adapters"
	"backend/internal/domain/models"
	"backend/internal/domain/search"
	"backend/internal/metrics"

	"github.com/rs/zerolog"
)

const (
	// SourceTypeC2C тип источника для C2C listings
	SourceTypeC2C = "c2c"
	// SourceTypeB2C тип источника для B2C products
	SourceTypeB2C = "b2c"
)

// C2CRepository определяет интерфейс для C2C репозитория
type C2CRepository interface {
	CreateListing(ctx context.Context, listing *models.MarketplaceListing) (int, error)
	GetListing(ctx context.Context, id int) (*models.MarketplaceListing, error)
	UpdateListing(ctx context.Context, listing *models.MarketplaceListing) error
	DeleteListing(ctx context.Context, id int) error
}

// B2CRepository определяет интерфейс для B2C репозитория
type B2CRepository interface {
	GetStorefrontProductByID(ctx context.Context, productID int) (*models.StorefrontProduct, error)
	CreateStorefrontProduct(ctx context.Context, storefrontID int, req *models.CreateProductRequest) (*models.StorefrontProduct, error)
	UpdateStorefrontProduct(ctx context.Context, storefrontID, productID int, req *models.UpdateProductRequest) error
	DeleteStorefrontProduct(ctx context.Context, storefrontID, productID int) error
	GetStorefrontProducts(ctx context.Context, filter models.ProductFilter) ([]*models.StorefrontProduct, error)
	GetStorefrontByID(ctx context.Context, id int) (*models.Storefront, error)
}

// SearchParams параметры поиска для unified listings
type SearchParams struct {
	Query        string
	CategoryID   int
	MinPrice     float64
	MaxPrice     float64
	Condition    string
	SourceType   string // "c2c", "b2c", "all"
	StorefrontID int
	UserID       int
	Limit        int
	Offset       int
	SortBy       string
	SortOrder    string
}

// RoutingContext контекст для routing решений
type RoutingContext struct {
	UserID  int  // User ID для routing
	IsAdmin bool // Является ли пользователь админом
}

// MarketplaceServiceInterface интерфейс для unified marketplace service
type MarketplaceServiceInterface interface {
	CreateListing(ctx context.Context, unified *models.UnifiedListing, routingCtx *RoutingContext) (int64, error)
	GetListing(ctx context.Context, id int64, sourceType string, routingCtx *RoutingContext) (*models.UnifiedListing, error)
	UpdateListing(ctx context.Context, unified *models.UnifiedListing, routingCtx *RoutingContext) error
	DeleteListing(ctx context.Context, id int64, sourceType string, routingCtx *RoutingContext) error
	SearchListings(ctx context.Context, params *SearchParams, routingCtx *RoutingContext) ([]*models.UnifiedListing, int64, error)
}

// ListingsGRPCClient определяет интерфейс для gRPC клиента listings микросервиса
type ListingsGRPCClient interface {
	GetListing(ctx context.Context, id int64) (*models.UnifiedListing, error)
	CreateListing(ctx context.Context, unified *models.UnifiedListing) (*models.UnifiedListing, error)
	UpdateListing(ctx context.Context, unified *models.UnifiedListing) (*models.UnifiedListing, error)
	DeleteListing(ctx context.Context, id int64, userID int64) error
}

// TrafficRouter интерфейс для routing решений (monolith vs microservice)
// Минимальный интерфейс для service layer
type TrafficRouter interface {
	ShouldUseMicroservice(userID string, isAdmin bool) *TrafficRoutingDecision
	ValidateConfig() error
}

// TrafficRoutingDecision содержит информацию о routing решении
type TrafficRoutingDecision struct {
	UseМicroservice bool
	Reason          string
	UserID          string
	IsAdmin         bool
	IsCanary        bool
	Hash            uint32
}

// MarketplaceService унифицированный сервис для работы с C2C и B2C listings
type MarketplaceService struct {
	c2cRepo                 C2CRepository
	b2cRepo                 B2CRepository
	osClient                OpenSearchRepository
	listingsGRPCClient      ListingsGRPCClient // gRPC клиент для listings микросервиса (опционально)
	useListingsMicroservice bool               // Feature flag для переключения между local DB и микросервисом
	trafficRouter           TrafficRouter      // Traffic router для marketplace microservice migration (Sprint 6.1)
	logger                  zerolog.Logger
}

// OpenSearchRepository определяет интерфейс для OpenSearch репозитория
type OpenSearchRepository interface {
	SearchListings(ctx context.Context, params *search.ServiceParams) (*search.ServiceResult, error)
	Index(ctx context.Context, listing *models.MarketplaceListing) error
	Delete(ctx context.Context, listingID int) error
}

// NewMarketplaceService создает новый unified service
func NewMarketplaceService(
	c2cRepo C2CRepository,
	b2cRepo B2CRepository,
	osClient OpenSearchRepository,
	logger zerolog.Logger,
) *MarketplaceService {
	return &MarketplaceService{
		c2cRepo:                 c2cRepo,
		b2cRepo:                 b2cRepo,
		osClient:                osClient,
		listingsGRPCClient:      nil,   // Будет установлен через SetListingsGRPCClient если нужно
		useListingsMicroservice: false, // По умолчанию используем local DB
		logger:                  logger.With().Str("service", "unified_marketplace").Logger(),
	}
}

// SetListingsGRPCClient устанавливает gRPC клиент для listings микросервиса
func (s *MarketplaceService) SetListingsGRPCClient(client ListingsGRPCClient, enabled bool) {
	s.listingsGRPCClient = client
	s.useListingsMicroservice = enabled
	if enabled && client != nil {
		s.logger.Info().Msg("Listings microservice integration enabled")
	} else {
		s.logger.Info().Msg("Using local database for listings")
	}
}

// SetTrafficRouter устанавливает traffic router для marketplace microservice migration
func (s *MarketplaceService) SetTrafficRouter(router TrafficRouter) {
	s.trafficRouter = router
	if router != nil {
		s.logger.Info().Msg("Traffic router configured for marketplace microservice migration")
	} else {
		s.logger.Info().Msg("No traffic router configured - using monolith only")
	}
}

// CreateListing создает listing (C2C или B2C в зависимости от source_type)
func (s *MarketplaceService) CreateListing(ctx context.Context, unified *models.UnifiedListing, routingCtx *RoutingContext) (int64, error) {
	s.logger.Info().
		Str("source_type", unified.SourceType).
		Str("title", unified.Title).
		Int("user_id", routingCtx.UserID).
		Bool("is_admin", routingCtx.IsAdmin).
		Msg("Creating unified listing")

	// Traffic routing decision (только для C2C)
	if unified.SourceType == SourceTypeC2C && s.trafficRouter != nil {
		decision := s.trafficRouter.ShouldUseMicroservice(fmt.Sprintf("%d", routingCtx.UserID), routingCtx.IsAdmin)
		s.logger.Info().
			Bool("use_microservice", decision.UseМicroservice).
			Str("reason", decision.Reason).
			Str("user_id", decision.UserID).
			Bool("is_admin", decision.IsAdmin).
			Bool("is_canary", decision.IsCanary).
			Uint32("hash", decision.Hash).
			Msg("Traffic routing decision for CreateListing")

		if decision.UseМicroservice {
			// TODO: Route to marketplace microservice via gRPC
			s.logger.Info().Msg("Routing to marketplace microservice (not implemented yet, fallback to monolith)")
		}
	}

	switch unified.SourceType {
	case SourceTypeC2C:
		return s.createC2CListing(ctx, unified)
	case SourceTypeB2C:
		return s.createB2CListing(ctx, unified)
	default:
		return 0, fmt.Errorf("invalid source_type: %s (must be '%s' or '%s')", unified.SourceType, SourceTypeC2C, SourceTypeB2C)
	}
}

// createC2CListing создает C2C listing через адаптер или микросервис
func (s *MarketplaceService) createC2CListing(ctx context.Context, unified *models.UnifiedListing) (int64, error) {
	// Feature flag: если включен listings микросервис - используем его
	if s.useListingsMicroservice && s.listingsGRPCClient != nil {
		s.logger.Info().Msg("Creating C2C listing via listings microservice")

		// Measure latency
		start := time.Now()
		created, err := s.listingsGRPCClient.CreateListing(ctx, unified)
		duration := time.Since(start).Seconds()

		if err != nil {
			s.logger.Error().Err(err).Msg("Failed to create C2C listing via microservice")

			// Record error metrics
			errorType := metrics.ClassifyGRPCError(err)
			metrics.RecordMicroserviceError(errorType)
			metrics.RecordFallback("microservice_error")

			// Graceful degradation: fallback to local DB if microservice fails
			s.logger.Warn().Msg("Falling back to local database")

			// Measure fallback latency
			fallbackStart := time.Now()
			result, fallbackErr := s.createC2CListingLocal(ctx, unified)
			fallbackDuration := time.Since(fallbackStart).Seconds()
			metrics.ObserveRouteDuration("monolith", "create", fallbackDuration)

			return result, fallbackErr
		}

		// Record success metrics
		metrics.ObserveRouteDuration("microservice", "create", duration)
		s.logger.Info().Int("listing_id", created.ID).Msg("C2C listing created via microservice successfully")
		return int64(created.ID), nil
	}

	// Иначе используем локальную БД
	start := time.Now()
	result, err := s.createC2CListingLocal(ctx, unified)
	duration := time.Since(start).Seconds()
	metrics.ObserveRouteDuration("monolith", "create", duration)
	return result, err
}

// createC2CListingLocal создает C2C listing через локальную БД
func (s *MarketplaceService) createC2CListingLocal(ctx context.Context, unified *models.UnifiedListing) (int64, error) {
	// Конвертируем unified → C2C
	c2c, err := adapters.UnifiedToC2C(unified)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to convert unified to C2C")
		return 0, fmt.Errorf("failed to convert unified to C2C: %w", err)
	}

	if c2c == nil {
		return 0, fmt.Errorf("conversion returned nil C2C listing")
	}

	// Создаем через C2C репозиторий
	id, err := s.c2cRepo.CreateListing(ctx, c2c)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to create C2C listing")
		return 0, fmt.Errorf("failed to create C2C listing: %w", err)
	}

	s.logger.Info().Int("listing_id", id).Msg("C2C listing created successfully (local DB)")

	// Индексируем в OpenSearch
	c2c.ID = id
	if err := s.osClient.Index(ctx, c2c); err != nil {
		s.logger.Warn().Err(err).Int("listing_id", id).Msg("Failed to index C2C listing in OpenSearch")
		// Не прерываем процесс, т.к. listing уже создан
	}

	return int64(id), nil
}

// createB2CListing создает B2C product через адаптер
func (s *MarketplaceService) createB2CListing(ctx context.Context, unified *models.UnifiedListing) (int64, error) {
	// B2C требует storefront_id
	if unified.StorefrontID == nil || *unified.StorefrontID == 0 {
		return 0, fmt.Errorf("storefront_id is required for B2C listings")
	}

	// Конвертируем unified → B2C
	b2c, err := adapters.UnifiedToB2C(unified)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to convert unified to B2C")
		return 0, fmt.Errorf("failed to convert unified to B2C: %w", err)
	}

	if b2c == nil {
		return 0, fmt.Errorf("conversion returned nil B2C product")
	}

	// Создаем CreateProductRequest
	hasIndividualLocation := b2c.HasIndividualLocation
	showOnMap := b2c.ShowOnMap
	req := &models.CreateProductRequest{
		Name:                  b2c.Name,
		Description:           b2c.Description,
		Price:                 b2c.Price,
		CategoryID:            b2c.CategoryID,
		SKU:                   b2c.SKU,
		Barcode:               b2c.Barcode,
		StockQuantity:         b2c.StockQuantity,
		IsActive:              b2c.IsActive,
		Attributes:            b2c.Attributes,
		Currency:              b2c.Currency,
		HasIndividualLocation: &hasIndividualLocation,
		IndividualAddress:     b2c.IndividualAddress,
		IndividualLatitude:    b2c.IndividualLatitude,
		IndividualLongitude:   b2c.IndividualLongitude,
		ShowOnMap:             &showOnMap,
	}

	// Создаем через B2C репозиторий
	product, err := s.b2cRepo.CreateStorefrontProduct(ctx, *unified.StorefrontID, req)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to create B2C product")
		return 0, fmt.Errorf("failed to create B2C product: %w", err)
	}

	s.logger.Info().Int("product_id", product.ID).Msg("B2C product created successfully")

	// TODO: Индексация B2C в OpenSearch (требует адаптацию для B2C)

	return int64(product.ID), nil
}

// GetListing получает listing по ID и типу источника
func (s *MarketplaceService) GetListing(ctx context.Context, id int64, sourceType string, routingCtx *RoutingContext) (*models.UnifiedListing, error) {
	s.logger.Info().
		Int64("id", id).
		Str("source_type", sourceType).
		Int("user_id", routingCtx.UserID).
		Bool("is_admin", routingCtx.IsAdmin).
		Msg("Getting unified listing")

	// Traffic routing decision (только для C2C)
	if sourceType == SourceTypeC2C && s.trafficRouter != nil {
		decision := s.trafficRouter.ShouldUseMicroservice(fmt.Sprintf("%d", routingCtx.UserID), routingCtx.IsAdmin)
		s.logger.Info().
			Bool("use_microservice", decision.UseМicroservice).
			Str("reason", decision.Reason).
			Str("user_id", decision.UserID).
			Msg("Traffic routing decision for GetListing")

		if decision.UseМicroservice {
			s.logger.Info().Msg("Routing to marketplace microservice (not implemented yet, fallback to monolith)")
		}
	}

	switch sourceType {
	case SourceTypeC2C:
		return s.getC2CListing(ctx, int(id))
	case SourceTypeB2C:
		return s.getB2CListing(ctx, int(id))
	default:
		return nil, fmt.Errorf("invalid source_type: %s (must be '%s' or '%s')", sourceType, SourceTypeC2C, SourceTypeB2C)
	}
}

// getC2CListing получает C2C listing через микросервис или локальную БД
func (s *MarketplaceService) getC2CListing(ctx context.Context, id int) (*models.UnifiedListing, error) {
	// Feature flag: если включен listings микросервис - используем его
	if s.useListingsMicroservice && s.listingsGRPCClient != nil {
		s.logger.Debug().Int("id", id).Msg("Getting C2C listing via listings microservice")

		// Measure latency
		start := time.Now()
		unified, err := s.listingsGRPCClient.GetListing(ctx, int64(id))
		duration := time.Since(start).Seconds()

		if err != nil {
			s.logger.Error().Err(err).Int("id", id).Msg("Failed to get C2C listing via microservice")

			// Record error metrics
			errorType := metrics.ClassifyGRPCError(err)
			metrics.RecordMicroserviceError(errorType)
			metrics.RecordFallback("microservice_error")

			// Graceful degradation: fallback to local DB if microservice fails
			s.logger.Warn().Msg("Falling back to local database")

			// Measure fallback latency
			fallbackStart := time.Now()
			result, fallbackErr := s.getC2CListingLocal(ctx, id)
			fallbackDuration := time.Since(fallbackStart).Seconds()
			metrics.ObserveRouteDuration("monolith", "get", fallbackDuration)

			return result, fallbackErr
		}

		// Record success metrics
		metrics.ObserveRouteDuration("microservice", "get", duration)
		return unified, nil
	}

	// Иначе используем локальную БД
	start := time.Now()
	result, err := s.getC2CListingLocal(ctx, id)
	duration := time.Since(start).Seconds()
	metrics.ObserveRouteDuration("monolith", "get", duration)
	return result, err
}

// getC2CListingLocal получает C2C listing из локальной БД
func (s *MarketplaceService) getC2CListingLocal(ctx context.Context, id int) (*models.UnifiedListing, error) {
	c2c, err := s.c2cRepo.GetListing(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).Int("id", id).Msg("Failed to get C2C listing")
		return nil, fmt.Errorf("failed to get C2C listing: %w", err)
	}

	if c2c == nil {
		return nil, fmt.Errorf("C2C listing not found: id=%d", id)
	}

	// Конвертируем C2C → unified
	unified, err := adapters.C2CToUnified(c2c)
	if err != nil {
		s.logger.Error().Err(err).Int("id", id).Msg("Failed to convert C2C to unified")
		return nil, fmt.Errorf("failed to convert C2C to unified: %w", err)
	}

	return unified, nil
}

// getB2CListing получает B2C product
func (s *MarketplaceService) getB2CListing(ctx context.Context, id int) (*models.UnifiedListing, error) {
	b2c, err := s.b2cRepo.GetStorefrontProductByID(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).Int("id", id).Msg("Failed to get B2C product")
		return nil, fmt.Errorf("failed to get B2C product: %w", err)
	}

	if b2c == nil {
		return nil, fmt.Errorf("B2C product not found: id=%d", id)
	}

	// Получаем storefront для полной конвертации
	var storefront *models.Storefront
	if b2c.StorefrontID != 0 {
		storefront, err = s.b2cRepo.GetStorefrontByID(ctx, b2c.StorefrontID)
		if err != nil {
			s.logger.Warn().Err(err).Int("storefront_id", b2c.StorefrontID).Msg("Failed to get storefront")
			// Продолжаем без storefront
		}
	}

	// Конвертируем B2C → unified
	unified, err := adapters.B2CToUnified(b2c, storefront)
	if err != nil {
		s.logger.Error().Err(err).Int("id", id).Msg("Failed to convert B2C to unified")
		return nil, fmt.Errorf("failed to convert B2C to unified: %w", err)
	}

	return unified, nil
}

// UpdateListing обновляет listing
func (s *MarketplaceService) UpdateListing(ctx context.Context, unified *models.UnifiedListing, routingCtx *RoutingContext) error {
	s.logger.Info().
		Int("id", unified.ID).
		Str("source_type", unified.SourceType).
		Int("user_id", routingCtx.UserID).
		Bool("is_admin", routingCtx.IsAdmin).
		Msg("Updating unified listing")

	// Traffic routing decision (только для C2C)
	if unified.SourceType == SourceTypeC2C && s.trafficRouter != nil {
		decision := s.trafficRouter.ShouldUseMicroservice(fmt.Sprintf("%d", routingCtx.UserID), routingCtx.IsAdmin)
		s.logger.Info().
			Bool("use_microservice", decision.UseМicroservice).
			Str("reason", decision.Reason).
			Str("user_id", decision.UserID).
			Msg("Traffic routing decision for UpdateListing")

		if decision.UseМicroservice {
			s.logger.Info().Msg("Routing to marketplace microservice (not implemented yet, fallback to monolith)")
		}
	}

	switch unified.SourceType {
	case SourceTypeC2C:
		return s.updateC2CListing(ctx, unified)
	case SourceTypeB2C:
		return s.updateB2CListing(ctx, unified)
	default:
		return fmt.Errorf("invalid source_type: %s (must be '%s' or '%s')", unified.SourceType, SourceTypeC2C, SourceTypeB2C)
	}
}

// updateC2CListing обновляет C2C listing через микросервис или локальную БД
func (s *MarketplaceService) updateC2CListing(ctx context.Context, unified *models.UnifiedListing) error {
	// Feature flag: если включен listings микросервис - используем его
	if s.useListingsMicroservice && s.listingsGRPCClient != nil {
		s.logger.Info().Int("id", unified.ID).Msg("Updating C2C listing via listings microservice")

		_, err := s.listingsGRPCClient.UpdateListing(ctx, unified)
		if err != nil {
			s.logger.Error().Err(err).Int("id", unified.ID).Msg("Failed to update C2C listing via microservice")

			// Graceful degradation: fallback to local DB if microservice fails
			s.logger.Warn().Msg("Falling back to local database")
			return s.updateC2CListingLocal(ctx, unified)
		}

		s.logger.Info().Int("listing_id", unified.ID).Msg("C2C listing updated via microservice successfully")
		return nil
	}

	// Иначе используем локальную БД
	return s.updateC2CListingLocal(ctx, unified)
}

// updateC2CListingLocal обновляет C2C listing в локальной БД
func (s *MarketplaceService) updateC2CListingLocal(ctx context.Context, unified *models.UnifiedListing) error {
	// Конвертируем unified → C2C
	c2c, err := adapters.UnifiedToC2C(unified)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to convert unified to C2C")
		return fmt.Errorf("failed to convert unified to C2C: %w", err)
	}

	if c2c == nil {
		return fmt.Errorf("conversion returned nil C2C listing")
	}

	// Обновляем через C2C репозиторий
	if err := s.c2cRepo.UpdateListing(ctx, c2c); err != nil {
		s.logger.Error().Err(err).Int("id", unified.ID).Msg("Failed to update C2C listing")
		return fmt.Errorf("failed to update C2C listing: %w", err)
	}

	s.logger.Info().Int("listing_id", unified.ID).Msg("C2C listing updated successfully (local DB)")

	// Переиндексируем в OpenSearch
	if err := s.osClient.Index(ctx, c2c); err != nil {
		s.logger.Warn().Err(err).Int("listing_id", unified.ID).Msg("Failed to reindex C2C listing in OpenSearch")
	}

	return nil
}

// updateB2CListing обновляет B2C product
func (s *MarketplaceService) updateB2CListing(ctx context.Context, unified *models.UnifiedListing) error {
	// B2C требует storefront_id
	if unified.StorefrontID == nil || *unified.StorefrontID == 0 {
		return fmt.Errorf("storefront_id is required for B2C listings")
	}

	// Конвертируем unified → B2C
	b2c, err := adapters.UnifiedToB2C(unified)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to convert unified to B2C")
		return fmt.Errorf("failed to convert unified to B2C: %w", err)
	}

	if b2c == nil {
		return fmt.Errorf("conversion returned nil B2C product")
	}

	// Создаем UpdateProductRequest
	hasIndividualLocation := b2c.HasIndividualLocation
	showOnMap := b2c.ShowOnMap
	req := &models.UpdateProductRequest{
		Name:                  &b2c.Name,
		Description:           &b2c.Description,
		Price:                 &b2c.Price,
		CategoryID:            &b2c.CategoryID,
		SKU:                   b2c.SKU,
		Barcode:               b2c.Barcode,
		StockQuantity:         &b2c.StockQuantity,
		IsActive:              &b2c.IsActive,
		Attributes:            b2c.Attributes,
		HasIndividualLocation: &hasIndividualLocation,
		IndividualAddress:     b2c.IndividualAddress,
		IndividualLatitude:    b2c.IndividualLatitude,
		IndividualLongitude:   b2c.IndividualLongitude,
		ShowOnMap:             &showOnMap,
	}

	// Обновляем через B2C репозиторий
	if err := s.b2cRepo.UpdateStorefrontProduct(ctx, *unified.StorefrontID, unified.ID, req); err != nil {
		s.logger.Error().Err(err).Int("id", unified.ID).Msg("Failed to update B2C product")
		return fmt.Errorf("failed to update B2C product: %w", err)
	}

	s.logger.Info().Int("product_id", unified.ID).Msg("B2C product updated successfully")

	// TODO: Переиндексация B2C в OpenSearch

	return nil
}

// DeleteListing удаляет listing
func (s *MarketplaceService) DeleteListing(ctx context.Context, id int64, sourceType string, routingCtx *RoutingContext) error {
	s.logger.Info().
		Int64("id", id).
		Str("source_type", sourceType).
		Int("user_id", routingCtx.UserID).
		Bool("is_admin", routingCtx.IsAdmin).
		Msg("Deleting unified listing")

	// Traffic routing decision (только для C2C)
	if sourceType == SourceTypeC2C && s.trafficRouter != nil {
		decision := s.trafficRouter.ShouldUseMicroservice(fmt.Sprintf("%d", routingCtx.UserID), routingCtx.IsAdmin)
		s.logger.Info().
			Bool("use_microservice", decision.UseМicroservice).
			Str("reason", decision.Reason).
			Str("user_id", decision.UserID).
			Msg("Traffic routing decision for DeleteListing")

		if decision.UseМicroservice {
			s.logger.Info().Msg("Routing to marketplace microservice (not implemented yet, fallback to monolith)")
		}
	}

	switch sourceType {
	case SourceTypeC2C:
		return s.deleteC2CListing(ctx, int(id))
	case SourceTypeB2C:
		return s.deleteB2CListing(ctx, int(id))
	default:
		return fmt.Errorf("invalid source_type: %s (must be '%s' or '%s')", sourceType, SourceTypeC2C, SourceTypeB2C)
	}
}

// deleteC2CListing удаляет C2C listing через микросервис или локальную БД
func (s *MarketplaceService) deleteC2CListing(ctx context.Context, id int) error {
	// Feature flag: если включен listings микросервис - используем его
	if s.useListingsMicroservice && s.listingsGRPCClient != nil {
		s.logger.Info().Int("id", id).Msg("Deleting C2C listing via listings microservice")

		// Для delete нужен userID - получим listing сначала
		listing, err := s.listingsGRPCClient.GetListing(ctx, int64(id))
		if err != nil {
			s.logger.Error().Err(err).Int("id", id).Msg("Failed to get C2C listing for deletion via microservice")
			// Fallback to local DB
			s.logger.Warn().Msg("Falling back to local database")
			return s.deleteC2CListingLocal(ctx, id)
		}

		err = s.listingsGRPCClient.DeleteListing(ctx, int64(id), int64(listing.UserID))
		if err != nil {
			s.logger.Error().Err(err).Int("id", id).Msg("Failed to delete C2C listing via microservice")
			// Fallback to local DB
			s.logger.Warn().Msg("Falling back to local database")
			return s.deleteC2CListingLocal(ctx, id)
		}

		s.logger.Info().Int("listing_id", id).Msg("C2C listing deleted via microservice successfully")
		return nil
	}

	// Иначе используем локальную БД
	return s.deleteC2CListingLocal(ctx, id)
}

// deleteC2CListingLocal удаляет C2C listing из локальной БД
func (s *MarketplaceService) deleteC2CListingLocal(ctx context.Context, id int) error {
	// Удаляем через C2C репозиторий
	if err := s.c2cRepo.DeleteListing(ctx, id); err != nil {
		s.logger.Error().Err(err).Int("id", id).Msg("Failed to delete C2C listing")
		return fmt.Errorf("failed to delete C2C listing: %w", err)
	}

	s.logger.Info().Int("listing_id", id).Msg("C2C listing deleted successfully (local DB)")

	// Удаляем из OpenSearch
	if err := s.osClient.Delete(ctx, id); err != nil {
		s.logger.Warn().Err(err).Int("listing_id", id).Msg("Failed to delete C2C listing from OpenSearch")
	}

	return nil
}

// deleteB2CListing удаляет B2C product
func (s *MarketplaceService) deleteB2CListing(ctx context.Context, id int) error {
	// Получаем product чтобы узнать storefront_id
	b2c, err := s.b2cRepo.GetStorefrontProductByID(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).Int("id", id).Msg("Failed to get B2C product for deletion")
		return fmt.Errorf("failed to get B2C product: %w", err)
	}

	if b2c == nil {
		return fmt.Errorf("B2C product not found: id=%d", id)
	}

	// Удаляем через B2C репозиторий
	if err := s.b2cRepo.DeleteStorefrontProduct(ctx, b2c.StorefrontID, id); err != nil {
		s.logger.Error().Err(err).Int("id", id).Msg("Failed to delete B2C product")
		return fmt.Errorf("failed to delete B2C product: %w", err)
	}

	s.logger.Info().Int("product_id", id).Msg("B2C product deleted successfully")

	// TODO: Удаление B2C из OpenSearch

	return nil
}

// SearchListings выполняет поиск по unified listings через OpenSearch
func (s *MarketplaceService) SearchListings(ctx context.Context, params *SearchParams, routingCtx *RoutingContext) ([]*models.UnifiedListing, int64, error) {
	s.logger.Info().
		Str("query", params.Query).
		Int("category_id", params.CategoryID).
		Str("source_type", params.SourceType).
		Int("user_id", routingCtx.UserID).
		Bool("is_admin", routingCtx.IsAdmin).
		Msg("Searching unified listings")

	// Traffic routing decision (для C2C search)
	if s.trafficRouter != nil {
		decision := s.trafficRouter.ShouldUseMicroservice(fmt.Sprintf("%d", routingCtx.UserID), routingCtx.IsAdmin)
		s.logger.Info().
			Bool("use_microservice", decision.UseМicroservice).
			Str("reason", decision.Reason).
			Str("user_id", decision.UserID).
			Msg("Traffic routing decision for SearchListings")

		if decision.UseМicroservice {
			s.logger.Info().Msg("Routing to marketplace microservice (not implemented yet, fallback to monolith)")
		}
	}

	// Конвертируем SearchParams → search.ServiceParams
	categoryID := fmt.Sprintf("%d", params.CategoryID)
	searchParams := &search.ServiceParams{
		Query:         params.Query,
		CategoryID:    categoryID,
		PriceMin:      params.MinPrice,
		PriceMax:      params.MaxPrice,
		Condition:     params.Condition,
		Page:          params.Offset / params.Limit,
		Size:          params.Limit,
		Sort:          params.SortBy,
		SortDirection: params.SortOrder,
	}

	// Выполняем поиск в OpenSearch
	result, err := s.osClient.SearchListings(ctx, searchParams)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to search listings in OpenSearch")
		return nil, 0, fmt.Errorf("failed to search listings: %w", err)
	}

	// Конвертируем результаты в unified
	unified := make([]*models.UnifiedListing, 0, len(result.Items))
	for _, c2c := range result.Items {
		u, err := adapters.C2CToUnified(c2c)
		if err != nil {
			s.logger.Warn().Err(err).Int("listing_id", c2c.ID).Msg("Failed to convert C2C to unified")
			continue
		}

		unified = append(unified, u)
	}

	s.logger.Info().
		Int("count", len(unified)).
		Int("total", result.Total).
		Msg("Search completed successfully")

	return unified, int64(result.Total), nil
}
