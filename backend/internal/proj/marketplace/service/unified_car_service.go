package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"backend/internal/domain/models"
	"backend/internal/logger"
	"backend/internal/storage"

	"github.com/redis/go-redis/v9"
)

// UnifiedCarService объединяет всю логику работы с автомобилями
type UnifiedCarService struct {
	storage storage.Storage
	cache   *redis.Client
	config  *CarServiceConfig
}

// CarServiceConfig конфигурация для автомобильного сервиса
type CarServiceConfig struct {
	// API конфигурации для внешних сервисов (если понадобятся в будущем)
	VehicleDatabasesAPIKey string
	NHTSAEnabled           bool

	// Настройки кеширования
	CacheTTL     time.Duration
	CacheEnabled bool

	// Настройки VIN декодирования
	VINDecoderEnabled bool
}

// NewUnifiedCarService создает новый унифицированный сервис
func NewUnifiedCarService(storage storage.Storage, cache *redis.Client, config *CarServiceConfig) *UnifiedCarService {
	if config == nil {
		config = &CarServiceConfig{
			CacheTTL:     24 * time.Hour,
			CacheEnabled: true,
		}
	}

	return &UnifiedCarService{
		storage: storage,
		cache:   cache,
		config:  config,
	}
}

// GetMakesWithStats возвращает марки с дополнительной статистикой
func (s *UnifiedCarService) GetMakesWithStats(ctx context.Context, filters CarMakeFilters) ([]models.CarMakeWithStats, error) {
	// Попытка получить из кеша
	if s.config.CacheEnabled && s.cache != nil {
		_ = fmt.Sprintf("car:makes:stats:%s:%t:%t:%t", filters.Country, filters.IsDomestic, filters.IsMotorcycle, filters.ActiveOnly)
		// TODO: реализовать кеширование
	}

	// Получаем марки из storage
	makes, err := s.storage.GetCarMakes(ctx, filters.Country, filters.IsDomestic, filters.IsMotorcycle, filters.ActiveOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to get car makes: %w", err)
	}

	// Добавляем статистику для каждой марки
	result := make([]models.CarMakeWithStats, 0, len(makes))
	for _, make := range makes {
		stats, err := s.getCarMakeStats(ctx, make.ID)
		if err != nil {
			logger.Error().Err(err).Int("make_id", make.ID).Msg("Failed to get make stats")
			stats = &CarMakeStats{} // Пустая статистика в случае ошибки
		}

		result = append(result, models.CarMakeWithStats{
			CarMake:      make,
			ModelCount:   stats.ModelCount,
			ListingCount: stats.ListingCount,
		})
	}

	return result, nil
}

// SearchVehicles унифицированный поиск автомобилей
func (s *UnifiedCarService) SearchVehicles(ctx context.Context, criteria VehicleSearchCriteria) (*VehicleSearchResult, error) {
	// Валидация критериев поиска
	if err := criteria.Validate(); err != nil {
		return nil, fmt.Errorf("invalid search criteria: %w", err)
	}

	// TODO: Реализовать поиск через OpenSearch когда будет интеграция
	// Пока используем простой поиск через БД

	result := &VehicleSearchResult{
		Vehicles: []models.VehicleInfo{},
		Total:    0,
		Facets:   make(map[string][]FacetValue),
	}

	return result, nil
}

// DecodeVIN декодирует VIN номер автомобиля
func (s *UnifiedCarService) DecodeVIN(ctx context.Context, vin string) (*models.VINDecodeResult, error) {
	if !s.config.VINDecoderEnabled {
		return nil, fmt.Errorf("VIN decoder is disabled")
	}

	// Валидация VIN
	if len(vin) != 17 {
		return nil, fmt.Errorf("invalid VIN length: expected 17, got %d", len(vin))
	}

	// TODO: Проверить кеш
	// TODO: Использовать NHTSA API (бесплатный)
	// TODO: При необходимости использовать Vehicle Databases API (платный)

	return nil, fmt.Errorf("VIN decoding not implemented yet")
}

// GetCarMake получает марку по ID или slug
func (s *UnifiedCarService) GetCarMake(ctx context.Context, identifier string) (*models.CarMake, error) {
	// Пытаемся найти по slug
	make, err := s.storage.GetCarMakeBySlug(ctx, identifier)
	if err == nil {
		return make, nil
	}

	// Если не нашли по slug, пробуем по ID
	// TODO: реализовать поиск по ID

	return nil, fmt.Errorf("car make not found: %s", identifier)
}

// GetModelsByMake возвращает модели для конкретной марки
func (s *UnifiedCarService) GetModelsByMake(ctx context.Context, makeSlug string, activeOnly bool) ([]models.CarModel, error) {
	return s.storage.GetCarModelsByMake(ctx, makeSlug, activeOnly)
}

// GetGenerationsByModel возвращает поколения для конкретной модели
func (s *UnifiedCarService) GetGenerationsByModel(ctx context.Context, modelID int, activeOnly bool) ([]models.CarGeneration, error) {
	return s.storage.GetCarGenerationsByModel(ctx, modelID, activeOnly)
}

// SyncExternalData синхронизирует данные с внешними источниками
func (s *UnifiedCarService) SyncExternalData(ctx context.Context, source string) error {
	switch source {
	case "nhtsa":
		return s.syncNHTSAData(ctx)
	case "local":
		// Локальные данные уже в БД
		return nil
	default:
		return fmt.Errorf("unknown sync source: %s", source)
	}
}

// Helper методы

func (s *UnifiedCarService) getCarMakeStats(ctx context.Context, makeID int) (*CarMakeStats, error) {
	// TODO: Реализовать подсчет статистики
	// - Количество моделей
	// - Количество активных объявлений
	// - Средняя цена

	return &CarMakeStats{
		ModelCount:   0,
		ListingCount: 0,
	}, nil
}

func (s *UnifiedCarService) syncNHTSAData(ctx context.Context) error {
	// TODO: Реализовать синхронизацию с NHTSA API
	// Бесплатный API для VIN декодирования и данных о машинах

	client := &http.Client{Timeout: 30 * time.Second}
	_ = client // Временно, чтобы не было warning

	return fmt.Errorf("NHTSA sync not implemented yet")
}

// Структуры для работы с автомобилями

// CarMakeFilters фильтры для получения марок
type CarMakeFilters struct {
	Country      string
	IsDomestic   bool
	IsMotorcycle bool
	ActiveOnly   bool
}

// CarMakeStats статистика по марке
type CarMakeStats struct {
	ModelCount   int
	ListingCount int
	AvgPrice     float64
}

// VehicleSearchCriteria критерии поиска автомобилей
type VehicleSearchCriteria struct {
	MakeID       *int
	ModelID      *int
	GenerationID *int
	YearFrom     *int
	YearTo       *int
	PriceFrom    *float64
	PriceTo      *float64
	MileageFrom  *int
	MileageTo    *int
	FuelType     *string
	Transmission *string
	BodyType     *string

	// Пагинация
	Offset int
	Limit  int

	// Сортировка
	SortBy    string
	SortOrder string
}

// Validate проверяет корректность критериев поиска
func (c *VehicleSearchCriteria) Validate() error {
	if c.Limit <= 0 {
		c.Limit = 20
	}
	if c.Limit > 100 {
		c.Limit = 100
	}

	if c.YearFrom != nil && c.YearTo != nil && *c.YearFrom > *c.YearTo {
		return fmt.Errorf("year_from cannot be greater than year_to")
	}

	if c.PriceFrom != nil && c.PriceTo != nil && *c.PriceFrom > *c.PriceTo {
		return fmt.Errorf("price_from cannot be greater than price_to")
	}

	return nil
}

// VehicleSearchResult результат поиска автомобилей
type VehicleSearchResult struct {
	Vehicles []models.VehicleInfo
	Total    int
	Facets   map[string][]FacetValue
}

// FacetValue значение фасета для фильтрации
type FacetValue struct {
	Value string
	Count int
	Label string
}
