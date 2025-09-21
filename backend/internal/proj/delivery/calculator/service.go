package calculator

import (
	"context"
	"encoding/json"
	"fmt"
	"math"

	"github.com/jmoiron/sqlx"

	"backend/internal/proj/delivery/interfaces"
	"backend/internal/proj/delivery/models"
)

// Service - сервис калькулятора доставки
type Service struct {
	db *sqlx.DB
}

// NewService создает новый экземпляр сервиса калькулятора
func NewService(db *sqlx.DB) *Service {
	return &Service{
		db: db,
	}
}

// CalculationRequest - запрос расчета стоимости доставки
type CalculationRequest struct {
	FromLocation   Location        `json:"from_location"`
	ToLocation     Location        `json:"to_location"`
	Items          []ItemWithAttrs `json:"items"`
	ProviderID     *int            `json:"provider_id,omitempty"`
	InsuranceValue float64         `json:"insurance_value,omitempty"`
	CODAmount      float64         `json:"cod_amount,omitempty"`
	DeliveryType   string          `json:"delivery_type,omitempty"`
}

// Location - местоположение
type Location struct {
	City       string  `json:"city"`
	PostalCode string  `json:"postal_code"`
	Country    string  `json:"country"`
	Latitude   float64 `json:"latitude,omitempty"`
	Longitude  float64 `json:"longitude,omitempty"`
}

// ItemWithAttrs - товар с атрибутами доставки
type ItemWithAttrs struct {
	ProductID   int                        `json:"product_id"`
	ProductType string                     `json:"product_type"` // "listing" или "storefront_product"
	Quantity    int                        `json:"quantity"`
	Attributes  *models.DeliveryAttributes `json:"attributes,omitempty"`
}

// Calculate - рассчитывает стоимость доставки для всех доступных провайдеров
func (s *Service) Calculate(ctx context.Context, req *CalculationRequest) (*CalculationResponse, error) {
	// Загружаем атрибуты товаров если не переданы
	if err := s.loadItemAttributes(ctx, req.Items); err != nil {
		return &CalculationResponse{
			Success: false,
			Message: "Failed to load item attributes",
		}, fmt.Errorf("failed to load item attributes: %w", err)
	}

	// Оптимизируем упаковку товаров
	packages := s.optimizePackaging(req.Items)

	// Определяем зону доставки
	zone, err := s.determineZone(ctx, req.FromLocation, req.ToLocation)
	if err != nil {
		return &CalculationResponse{
			Success: false,
			Message: "Failed to determine delivery zone",
		}, fmt.Errorf("failed to determine zone: %w", err)
	}

	// Получаем активных провайдеров
	providers, err := s.getActiveProviders(ctx, req.ProviderID)
	if err != nil {
		return &CalculationResponse{
			Success: false,
			Message: "Failed to get delivery providers",
		}, fmt.Errorf("failed to get providers: %w", err)
	}

	// Рассчитываем стоимость для каждого провайдера
	var quotes []ProviderQuote
	for _, provider := range providers {
		quote, err := s.calculateProviderQuote(ctx, provider, packages, zone, req)
		if err != nil {
			// Логируем ошибку, но продолжаем с другими провайдерами
			continue
		}
		quotes = append(quotes, *quote)
	}

	// Создаем данные расчета
	data := &CalculationData{
		Providers: quotes,
	}

	// Выбираем лучшие предложения
	s.selectBestQuotes(data)

	// Формируем ответ
	response := &CalculationResponse{
		Success: true,
		Data:    data,
	}

	return response, nil
}

// loadItemAttributes - загружает атрибуты доставки для товаров
func (s *Service) loadItemAttributes(ctx context.Context, items []ItemWithAttrs) error {
	for i := range items {
		if items[i].Attributes == nil {
			attrs, err := s.getProductAttributes(ctx, items[i].ProductID, items[i].ProductType)
			if err != nil {
				return err
			}
			items[i].Attributes = attrs
		}
	}
	return nil
}

// getProductAttributes - получает атрибуты доставки товара из БД
func (s *Service) getProductAttributes(ctx context.Context, productID int, productType string) (*models.DeliveryAttributes, error) {
	var attrs models.DeliveryAttributes
	var query string

	if productType == "listing" {
		query = `
			SELECT
				COALESCE(
					ml.metadata->'delivery_attributes',
					(SELECT jsonb_build_object(
						'weight_kg', dcd.default_weight_kg,
						'dimensions', jsonb_build_object(
							'length_cm', dcd.default_length_cm,
							'width_cm', dcd.default_width_cm,
							'height_cm', dcd.default_height_cm
						),
						'packaging_type', dcd.default_packaging_type,
						'is_fragile', dcd.is_typically_fragile
					)
					FROM delivery_category_defaults dcd
					WHERE dcd.category_id = ml.category_id)
				) as attributes
			FROM marketplace_listings ml
			WHERE ml.id = $1`
	} else {
		query = `
			SELECT
				COALESCE(
					sp.attributes->'delivery_attributes',
					(SELECT jsonb_build_object(
						'weight_kg', dcd.default_weight_kg,
						'dimensions', jsonb_build_object(
							'length_cm', dcd.default_length_cm,
							'width_cm', dcd.default_width_cm,
							'height_cm', dcd.default_height_cm
						),
						'packaging_type', dcd.default_packaging_type,
						'is_fragile', dcd.is_typically_fragile
					)
					FROM delivery_category_defaults dcd
					WHERE dcd.category_id = sp.category_id)
				) as attributes
			FROM storefront_products sp
			WHERE sp.id = $1`
	}

	var jsonData []byte
	if err := s.db.GetContext(ctx, &jsonData, query, productID); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(jsonData, &attrs); err != nil {
		return nil, err
	}

	// Рассчитываем объем если не указан
	if attrs.VolumeM3 == 0 && attrs.Dimensions != nil {
		attrs.VolumeM3 = attrs.Dimensions.CalculateVolume()
	}

	return &attrs, nil
}

// Package - упакованная посылка
type Package struct {
	Items         []ItemWithAttrs
	TotalWeight   float64
	Dimensions    models.Dimensions
	Volume        float64
	IsFragile     bool
	PackagingType string
}

// optimizePackaging - оптимизирует упаковку товаров
func (s *Service) optimizePackaging(items []ItemWithAttrs) []Package {
	// Простая реализация - группируем товары по типу упаковки
	// В продакшене здесь может быть сложный алгоритм оптимизации

	packageMap := make(map[string]*Package)

	for _, item := range items {
		if item.Attributes == nil {
			continue
		}

		packType := item.Attributes.PackagingType
		if packType == "" {
			packType = "box"
		}

		pkg, exists := packageMap[packType]
		if !exists {
			pkg = &Package{
				PackagingType: packType,
				IsFragile:     item.Attributes.IsFragile,
			}
			packageMap[packType] = pkg
		}

		// Добавляем товар в упаковку
		for i := 0; i < item.Quantity; i++ {
			pkg.Items = append(pkg.Items, item)
			pkg.TotalWeight += item.Attributes.WeightKg
			if item.Attributes.IsFragile {
				pkg.IsFragile = true
			}
		}
	}

	// Конвертируем в слайс и рассчитываем габариты
	var packages []Package
	for _, pkg := range packageMap {
		// Упрощенный расчет габаритов - берем максимальные размеры
		var maxLength, maxWidth, totalHeight float64
		for _, item := range pkg.Items {
			if item.Attributes.Dimensions != nil {
				if item.Attributes.Dimensions.LengthCm > maxLength {
					maxLength = item.Attributes.Dimensions.LengthCm
				}
				if item.Attributes.Dimensions.WidthCm > maxWidth {
					maxWidth = item.Attributes.Dimensions.WidthCm
				}
				totalHeight += item.Attributes.Dimensions.HeightCm
			}
		}

		pkg.Dimensions = models.Dimensions{
			LengthCm: maxLength,
			WidthCm:  maxWidth,
			HeightCm: totalHeight,
		}
		pkg.Volume = pkg.Dimensions.CalculateVolume()

		packages = append(packages, *pkg)
	}

	return packages
}

// determineZone - определяет зону доставки
func (s *Service) determineZone(ctx context.Context, from, to Location) (string, error) {
	// Простая реализация на основе городов
	// В продакшене используется GIS и полигоны

	if from.Country != to.Country {
		return models.ZoneTypeInternational, nil
	}

	if from.City == to.City {
		return models.ZoneTypeLocal, nil
	}

	// Проверяем по списку городов в зонах
	var zone string
	query := `
		SELECT type
		FROM delivery_zones
		WHERE $1 = ANY(countries)
		AND ($2 = ANY(cities) OR cities IS NULL)
		ORDER BY
			CASE type
				WHEN 'local' THEN 1
				WHEN 'regional' THEN 2
				WHEN 'national' THEN 3
				ELSE 4
			END
		LIMIT 1`

	err := s.db.GetContext(ctx, &zone, query, to.Country, to.City)
	if err != nil {
		// По умолчанию национальная доставка
		return models.ZoneTypeNational, nil
	}

	return zone, nil
}

// getActiveProviders - получает активных провайдеров
func (s *Service) getActiveProviders(ctx context.Context, providerID *int) ([]models.Provider, error) {
	var providers []models.Provider
	query := `SELECT * FROM delivery_providers WHERE is_active = true`
	args := []interface{}{}

	if providerID != nil {
		query += " AND id = $1"
		args = append(args, *providerID)
	}

	query += " ORDER BY id"

	if err := s.db.SelectContext(ctx, &providers, query, args...); err != nil {
		return nil, err
	}

	return providers, nil
}

// calculateProviderQuote - рассчитывает стоимость для провайдера
func (s *Service) calculateProviderQuote(ctx context.Context, provider models.Provider, packages []Package, zone string, req *CalculationRequest) (*ProviderQuote, error) {
	// Проверяем доступность провайдера для зоны
	var capabilities map[string]interface{}
	if provider.Capabilities != nil {
		if err := json.Unmarshal(*provider.Capabilities, &capabilities); err != nil {
			return nil, err
		}
	}

	// Получаем правила расчета для провайдера
	var rules []models.PricingRule
	query := `SELECT * FROM delivery_pricing_rules WHERE provider_id = $1 AND is_active = true ORDER BY priority DESC`
	if err := s.db.SelectContext(ctx, &rules, query, provider.ID); err != nil {
		return nil, err
	}

	if len(rules) == 0 {
		return &ProviderQuote{
			ProviderID:        provider.ID,
			ProviderCode:      provider.Code,
			ProviderName:      provider.Name,
			IsAvailable:       false,
			UnavailableReason: "No pricing rules configured",
		}, nil
	}

	// Рассчитываем стоимость
	totalWeight := 0.0
	totalVolume := 0.0
	hasFragile := false
	hasOversized := false

	for _, pkg := range packages {
		totalWeight += pkg.TotalWeight
		totalVolume += pkg.Volume
		if pkg.IsFragile {
			hasFragile = true
		}
		// Проверяем на негабарит (любая сторона > 100см)
		if pkg.Dimensions.LengthCm > 100 || pkg.Dimensions.WidthCm > 100 || pkg.Dimensions.HeightCm > 100 {
			hasOversized = true
		}
	}

	// Применяем правила расчета
	breakdown := models.CostBreakdown{}

	for _, rule := range rules {
		if rule.RuleType == models.RuleTypeWeightBased {
			// Расчет по весу
			var weightRanges []models.WeightRange
			if err := json.Unmarshal(rule.WeightRanges, &weightRanges); err != nil {
				continue
			}

			for _, r := range weightRanges {
				if totalWeight >= r.From && totalWeight <= r.To {
					breakdown.BasePrice = r.BasePrice
					if r.PricePerKg > 0 {
						breakdown.WeightSurcharge = (totalWeight - r.From) * r.PricePerKg
					}
					break
				}
			}
		}

		// Добавляем доплаты
		if hasFragile {
			breakdown.FragileSurcharge = rule.FragileSurcharge
		}
		if hasOversized {
			breakdown.OversizeSurcharge = rule.OversizedSurcharge
		}

		// Страховка
		if req.InsuranceValue > 0 && provider.SupportsInsurance {
			breakdown.InsuranceFee = req.InsuranceValue * 0.01 // 1% от суммы
		}

		// Наложенный платеж
		if req.CODAmount > 0 && provider.SupportsCOD {
			breakdown.CODFee = math.Max(50, req.CODAmount*0.02) // 2% но не менее 50
		}

		// Применяем минимальную и максимальную стоимость
		breakdown.CalculateTotal()
		if rule.MinPrice != nil && breakdown.Total < *rule.MinPrice {
			breakdown.Total = *rule.MinPrice
		}
		if rule.MaxPrice != nil && breakdown.Total > *rule.MaxPrice {
			breakdown.Total = *rule.MaxPrice
		}

		break // Используем первое подходящее правило
	}

	// Оцениваем время доставки
	estimatedDays := s.estimateDeliveryDays(zone, req.DeliveryType)

	return &ProviderQuote{
		ProviderID:    provider.ID,
		ProviderCode:  provider.Code,
		ProviderName:  provider.Name,
		DeliveryType:  req.DeliveryType,
		TotalPrice:    breakdown.Total,
		DeliveryCost:  breakdown.BasePrice,
		InsuranceCost: breakdown.InsuranceFee,
		CODFee:        breakdown.CODFee,
		CostBreakdown: breakdown,
		EstimatedDays: estimatedDays,
		IsAvailable:   true,
	}, nil
}

// estimateDeliveryDays - оценивает количество дней доставки
func (s *Service) estimateDeliveryDays(zone string, deliveryType string) int {
	baseDays := map[string]int{
		models.ZoneTypeLocal:         1,
		models.ZoneTypeRegional:      3,
		models.ZoneTypeNational:      5,
		models.ZoneTypeInternational: 10,
	}

	days := baseDays[zone]
	if days == 0 {
		days = 5
	}

	// Корректируем по типу доставки
	switch deliveryType {
	case interfaces.DeliveryTypeSameDay:
		return 0
	case interfaces.DeliveryTypeNextDay:
		return 1
	case interfaces.DeliveryTypeExpress:
		return days / 2
	case interfaces.DeliveryTypeEconomy:
		return days * 2
	default:
		return days
	}
}

// selectBestQuotes - выбирает лучшие предложения
func (s *Service) selectBestQuotes(data *CalculationData) {
	if len(data.Providers) == 0 {
		return
	}

	var cheapest, fastest, recommended *ProviderQuote

	for i := range data.Providers {
		quote := &data.Providers[i]

		if !quote.IsAvailable {
			continue
		}

		// Самый дешевый
		if cheapest == nil || quote.TotalPrice < cheapest.TotalPrice {
			cheapest = quote
		}

		// Самый быстрый
		if fastest == nil || quote.EstimatedDays < fastest.EstimatedDays {
			fastest = quote
		}

		// Рекомендованный (баланс цены и скорости)
		if recommended == nil {
			recommended = quote
		} else {
			// Простая формула: цена + дни * 100
			currentScore := recommended.TotalPrice + float64(recommended.EstimatedDays*100)
			quoteScore := quote.TotalPrice + float64(quote.EstimatedDays*100)
			if quoteScore < currentScore {
				recommended = quote
			}
		}
	}

	data.Cheapest = cheapest
	data.Fastest = fastest
	data.Recommended = recommended
}

// GetCategoryDefaults - получает дефолтные атрибуты для категории
func (s *Service) GetCategoryDefaults(ctx context.Context, categoryID int) (*models.CategoryDefaults, error) {
	var defaults models.CategoryDefaults
	query := `SELECT * FROM delivery_category_defaults WHERE category_id = $1`

	if err := s.db.GetContext(ctx, &defaults, query, categoryID); err != nil {
		return nil, err
	}

	return &defaults, nil
}

// SaveCategoryDefaults - сохраняет дефолтные атрибуты для категории
func (s *Service) SaveCategoryDefaults(ctx context.Context, defaults *models.CategoryDefaults) error {
	query := `
		INSERT INTO delivery_category_defaults (
			category_id, default_weight_kg, default_length_cm,
			default_width_cm, default_height_cm, default_packaging_type,
			is_typically_fragile
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (category_id) DO UPDATE SET
			default_weight_kg = EXCLUDED.default_weight_kg,
			default_length_cm = EXCLUDED.default_length_cm,
			default_width_cm = EXCLUDED.default_width_cm,
			default_height_cm = EXCLUDED.default_height_cm,
			default_packaging_type = EXCLUDED.default_packaging_type,
			is_typically_fragile = EXCLUDED.is_typically_fragile,
			updated_at = NOW()
		RETURNING id`

	return s.db.GetContext(ctx, &defaults.ID, query,
		defaults.CategoryID,
		defaults.DefaultWeightKg,
		defaults.DefaultLengthCm,
		defaults.DefaultWidthCm,
		defaults.DefaultHeightCm,
		defaults.DefaultPackagingType,
		defaults.IsTypicallyFragile,
	)
}
