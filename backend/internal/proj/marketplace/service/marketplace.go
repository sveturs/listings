// backend/internal/proj/marketplace/service/marketplace.go
package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"backend/internal/cache"
	"backend/internal/common"
	"backend/internal/config"
	"backend/internal/domain/models"
	"backend/internal/domain/search"
	"backend/internal/logger"
	"backend/internal/storage"

	"github.com/rs/zerolog"
	//	"net/url"
)

const (
	// Attribute names
	attributeNameModel = "model"

	// Attribute types
	attributeTypeText = "text"

	// Languages
	languageAuto = "auto"

	// Field names
	fieldNameName = "name"

	// SQL queries
	insertTranslationQuery = `
        INSERT INTO translations (
            entity_type, entity_id, language, field_name,
            translated_text, is_machine_translated, is_verified, metadata,
            last_modified_by
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        ON CONFLICT (entity_type, entity_id, language, field_name)
        DO UPDATE SET
            translated_text = EXCLUDED.translated_text,
            is_machine_translated = EXCLUDED.is_machine_translated,
            is_verified = EXCLUDED.is_verified,
            metadata = EXCLUDED.metadata,
            last_modified_by = EXCLUDED.last_modified_by,
            updated_at = CURRENT_TIMESTAMP
    `
)

type MarketplaceService struct {
	storage            storage.Storage
	translationService TranslationServiceInterface
	OrderService       OrderServiceInterface
	searchWeights      *config.SearchWeights
	cache              CacheInterface
	logger             zerolog.Logger
}

func NewMarketplaceService(storage storage.Storage, translationService TranslationServiceInterface, searchWeights *config.SearchWeights, cache CacheInterface) MarketplaceServiceInterface {
	ms := &MarketplaceService{
		storage:            storage,
		translationService: translationService,
		searchWeights:      searchWeights,
		cache:              cache,
		logger:             logger.Get().With().Str("service", "marketplace").Logger(),
	}

	// Создаем сервис заказов напрямую, избегая циклической зависимости
	ms.OrderService = NewSimpleOrderService(storage)

	return ms
}

func (s *MarketplaceService) GetUserFavorites(ctx context.Context, userID int) ([]models.MarketplaceListing, error) {
	return s.storage.GetUserFavorites(ctx, userID)
}

// SetTranslationService allows injecting a translation service after creation
func (s *MarketplaceService) SetTranslationService(svc TranslationServiceInterface) {
	s.translationService = svc
}

func (s *MarketplaceService) CreateListing(ctx context.Context, listing *models.MarketplaceListing) (int, error) {
	listing.Status = "active"
	listing.ViewsCount = 0

	if listing.OriginalLanguage == "" {
		// Пытаемся получить язык из контекста
		if userLang, ok := ctx.Value("language").(string); ok && userLang != "" {
			listing.OriginalLanguage = userLang
			log.Printf("Using language from context: %s", userLang)
		} else if userLang, ok := ctx.Value("Accept-Language").(string); ok && userLang != "" {
			listing.OriginalLanguage = userLang
			log.Printf("Using language from Accept-Language header: %s", userLang)
		} else {
			// Используем русский по умолчанию, т.к. большинство пользователей русскоговорящие
			listing.OriginalLanguage = "ru"
			log.Printf("Using default language (ru)")
		}
	}

	// Если указана витрина, и нет данных о местоположении, получаем их из витрины
	if listing.StorefrontID != nil && (listing.City == "" || listing.Country == "" || listing.Location == "") {
		storefront, err := s.storage.GetStorefrontByID(ctx, *listing.StorefrontID)
		if err == nil && storefront != nil {
			// Применяем данные о местоположении из витрины
			if listing.City == "" && storefront.City != "" {
				listing.City = storefront.City
				log.Printf("Using city from storefront: %s", storefront.City)
			}

			if listing.Country == "" && storefront.Country != "" {
				listing.Country = storefront.Country
				log.Printf("Using country from storefront: %s", storefront.Country)
			}

			if listing.Location == "" && storefront.Address != "" {
				listing.Location = storefront.Address
				log.Printf("Using address from storefront: %s", storefront.Address)
			}

			// Если нет координат
			if (listing.Latitude == nil || listing.Longitude == nil) &&
				storefront.Latitude != nil && storefront.Longitude != nil {
				listing.Latitude = storefront.Latitude
				listing.Longitude = storefront.Longitude
				log.Printf("Using coordinates from storefront: Lat=%f, Lon=%f",
					*storefront.Latitude, *storefront.Longitude)
			}
		} else {
			log.Printf("Could not get storefront or storefront has no location info: %v", err)
		}
	}

	listingID, err := s.storage.CreateListing(ctx, listing)
	if err != nil {
		return 0, err
	}

	if len(listing.Attributes) > 0 {
		// Фильтрация дубликатов атрибутов
		uniqueAttrs := make(map[int]models.ListingAttributeValue)
		for _, attr := range listing.Attributes {
			uniqueAttrs[attr.AttributeID] = attr // Последнее значение перезапишет предыдущее
		}

		// Преобразуем карту обратно в срез
		filteredAttrs := make([]models.ListingAttributeValue, 0, len(uniqueAttrs))
		for _, attr := range uniqueAttrs {
			attr.ListingID = listingID
			filteredAttrs = append(filteredAttrs, attr)
		}

		if err := s.SaveListingAttributes(ctx, listingID, filteredAttrs); err != nil {
			log.Printf("Error saving attributes for listing %d: %v", listingID, err)
		} else {
			log.Printf("Successfully saved %d attributes for listing %d", len(filteredAttrs), listingID)
		}
	}
	listing.ID = listingID

	// Сохраняем переводы адресных полей
	addressFields := make(map[string]string)
	if listing.Location != "" {
		addressFields["location"] = listing.Location
	}
	if listing.City != "" {
		addressFields["city"] = listing.City
	}
	if listing.Country != "" {
		addressFields["country"] = listing.Country
	}

	// Если есть адресные поля, переводим их на другие языки
	if len(addressFields) > 0 {
		// Определяем целевые языки (все поддерживаемые языки кроме исходного)
		targetLanguages := []string{"en", "ru", "sr"}
		// Удаляем исходный язык из списка целевых
		filteredTargetLanguages := make([]string, 0, len(targetLanguages))
		for _, lang := range targetLanguages {
			if lang != listing.OriginalLanguage {
				filteredTargetLanguages = append(filteredTargetLanguages, lang)
			}
		}

		// Сохраняем переводы адресов асинхронно, чтобы не замедлять создание объявления
		go func(ctx context.Context) {
			if err := s.SaveAddressTranslations(ctx, listingID, addressFields, listing.OriginalLanguage, filteredTargetLanguages); err != nil {
				log.Printf("Error saving address translations for listing %d: %v", listingID, err)
			}
		}(ctx)
	}

	fullListing, err := s.storage.GetListingByID(ctx, listingID)
	if err != nil {
		log.Printf("Ошибка получения полного объявления для индексации: %v", err)
	} else {
		if err := s.storage.IndexListing(ctx, fullListing); err != nil {
			log.Printf("Ошибка индексации объявления в OpenSearch: %v", err)
		}
	}

	return listingID, nil
}

func (s *MarketplaceService) GetSimilarListings(ctx context.Context, listingID int, limit int) ([]*models.MarketplaceListing, error) {
	log.Printf("=== GetSimilarListings: начало поиска похожих объявлений для ID=%d, limit=%d ===", listingID, limit)

	// Получаем исходное объявление
	listing, err := s.GetListingByID(ctx, listingID)
	if err != nil {
		log.Printf("ERROR: не удалось получить объявление %d: %v", listingID, err)
		return nil, fmt.Errorf("ошибка получения объявления: %w", err)
	}

	log.Printf("Исходное объявление: ID=%d, Title=%s, CategoryID=%d, Price=%.2f, City=%s, Country=%s, StorefrontID=%v",
		listing.ID, listing.Title, listing.CategoryID, listing.Price, listing.City, listing.Country, listing.StorefrontID)

	// Определяем источник объявления и выбираем соответствующую стратегию поиска
	if listing.StorefrontID != nil && *listing.StorefrontID > 0 {
		log.Printf("Объявление %d принадлежит витрине %d - используем поиск по товарам витрин", listingID, *listing.StorefrontID)
		return s.getSimilarStorefrontProducts(ctx, listingID, limit)
	}

	log.Printf("Объявление %d является обычным объявлением маркетплейса - используем стандартный поиск", listingID)

	if len(listing.Attributes) > 0 {
		log.Printf("Атрибуты объявления %d:", listing.ID)
		for _, attr := range listing.Attributes {
			log.Printf("  - %s: %s", attr.AttributeName, attr.DisplayValue)
		}
	} else {
		log.Printf("У объявления %d нет атрибутов", listing.ID)
	}

	// Создаем калькулятор похожести
	calculator := NewSimilarityCalculator(s.searchWeights)

	// Пытаемся найти похожие объявления с разными уровнями строгости
	var similarListings []*models.MarketplaceListing
	triesCount := 0
	maxTries := 4

	for triesCount < maxTries && len(similarListings) < limit {
		// Формируем параметры поиска для получения кандидатов
		params := s.buildAdvancedSearchParams(listing, limit*5, triesCount) // Получаем больше для фильтрации

		// Выполняем поиск похожих объявлений
		results, err := s.SearchListingsAdvanced(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("ошибка поиска похожих объявлений: %w", err)
		}

		// Фильтруем и сортируем результаты по похожести
		for _, candidate := range results.Items {
			if candidate.ID != listingID {
				// Проверяем, что кандидат еще не добавлен
				found := false
				for _, existing := range similarListings {
					if existing.ID == candidate.ID {
						found = true
						break
					}
				}
				if found {
					continue
				}

				// Вычисляем похожесть
				score, _ := calculator.CalculateSimilarity(ctx, listing, candidate)

				// Добавляем информацию о скоре в метаданные (для отладки)
				if candidate.Metadata == nil {
					candidate.Metadata = make(map[string]interface{})
				}
				candidate.Metadata["similarity_score"] = map[string]interface{}{
					"total":      score.TotalScore,
					"category":   score.CategoryScore,
					"attributes": score.AttributeScore,
					"price":      score.PriceScore,
					"location":   score.LocationScore,
					"text":       score.TextScore,
					"search_try": triesCount,
				}

				similarListings = append(similarListings, candidate)
			}
		}

		triesCount++
		log.Printf("Попытка %d: найдено %d похожих объявлений, всего собрано %d", triesCount, len(results.Items), len(similarListings))

		// Если найдено достаточно результатов, прекращаем поиск
		if len(similarListings) >= limit {
			break
		}
	}

	// Сортируем по убыванию похожести
	sort.Slice(similarListings, func(i, j int) bool {
		scoreI := similarListings[i].Metadata["similarity_score"].(map[string]interface{})["total"].(float64)
		scoreJ := similarListings[j].Metadata["similarity_score"].(map[string]interface{})["total"].(float64)
		return scoreI > scoreJ
	})

	// Ограничиваем количество результатов
	if len(similarListings) > limit {
		similarListings = similarListings[:limit]
	}

	log.Printf("Найдено %d похожих объявлений для листинга %d после %d попыток", len(similarListings), listingID, triesCount)

	return similarListings, nil
}

// buildAdvancedSearchParams формирует параметры для поиска похожих объявлений
// tryNumber определяет уровень строгости поиска: 0 - самый строгий, 3 - самый широкий
func (s *MarketplaceService) buildAdvancedSearchParams(listing *models.MarketplaceListing, size int, tryNumber int) *search.ServiceParams {
	params := &search.ServiceParams{
		Size: size,
		Page: 1,
		Sort: "date_desc",
	}

	// Категория - расширяем поиск на последней попытке
	if tryNumber < 3 {
		// Первые 3 попытки - ищем в той же категории
		params.CategoryID = strconv.Itoa(listing.CategoryID)
		log.Printf("Попытка %d: поиск в категории %d", tryNumber, listing.CategoryID)
	} else {
		// Последняя попытка - поиск во всех категориях
		log.Printf("Попытка %d: поиск во всех категориях", tryNumber)
	}

	// Добавляем локацию в зависимости от попытки
	switch {
	case listing.City != "" && tryNumber < 2:
		// Первые 2 попытки - ищем в том же городе
		params.City = listing.City
		log.Printf("Попытка %d: поиск в городе %s", tryNumber, listing.City)
	case listing.Country != "" && tryNumber == 2:
		// Третья попытка - ищем в той же стране
		params.Country = listing.Country
		log.Printf("Попытка %d: поиск в стране %s", tryNumber, listing.Country)
	default:
		// Последняя попытка - без географических ограничений
		log.Printf("Попытка %d: поиск без географических ограничений", tryNumber)
	}

	// Добавляем диапазон цен в зависимости от попытки
	if listing.Price > 0 {
		switch tryNumber {
		case 0:
			// Первая попытка: ±50%
			params.PriceMin = listing.Price * 0.5
			params.PriceMax = listing.Price * 1.5
			log.Printf("Попытка %d: диапазон цен %.2f - %.2f (±50%%)", tryNumber, params.PriceMin, params.PriceMax)
		case 1:
			// Вторая попытка: ±100%
			params.PriceMin = listing.Price * 0.3
			params.PriceMax = listing.Price * 2.0
			log.Printf("Попытка %d: диапазон цен %.2f - %.2f (±100%%)", tryNumber, params.PriceMin, params.PriceMax)
		case 2:
			// Третья попытка: ±200%
			params.PriceMin = listing.Price * 0.1
			params.PriceMax = listing.Price * 3.0
			log.Printf("Попытка %d: диапазон цен %.2f - %.2f (±200%%)", tryNumber, params.PriceMin, params.PriceMax)
		default:
			// Последняя попытка: без ограничений по цене
			log.Printf("Попытка %d: без ограничений по цене", tryNumber)
		}
	}

	// Добавляем ключевые атрибуты для фильтрации в зависимости от попытки
	if len(listing.Attributes) > 0 && tryNumber < 3 {
		attributeFilters := make(map[string]string)

		// Приоритетные атрибуты для разных категорий
		priorityAttrs := []string{"make", attributeNameModel, "brand", "type", "rooms", "property_type", "body_type"}

		// В зависимости от попытки используем разное количество атрибутов
		var maxAttrs int
		switch tryNumber {
		case 1:
			maxAttrs = 2 // Во второй попытке используем меньше атрибутов
		case 2:
			maxAttrs = 1 // В третьей попытке используем только самые важные
		default:
			maxAttrs = 3
		}

		attrCount := 0
		for _, attr := range listing.Attributes {
			if attrCount >= maxAttrs {
				break
			}
			// Добавляем только приоритетные атрибуты
			for _, priority := range priorityAttrs {
				if attr.AttributeName == priority && attr.DisplayValue != "" {
					attributeFilters[attr.AttributeName] = attr.DisplayValue
					attrCount++
					break
				}
			}
		}

		if len(attributeFilters) > 0 {
			params.AttributeFilters = attributeFilters
			log.Printf("Попытка %d: фильтр по атрибутам %v", tryNumber, attributeFilters)
		}
	} else {
		log.Printf("Попытка %d: без фильтров по атрибутам", tryNumber)
	}

	return params
}

func (s *MarketplaceService) GetListings(ctx context.Context, filters map[string]string, limit int, offset int) ([]models.MarketplaceListing, int64, error) {
	return s.storage.GetListings(ctx, filters, limit, offset)
}

func (s *MarketplaceService) GetFavoritedUsers(ctx context.Context, listingID int) ([]int, error) {
	return s.storage.GetFavoritedUsers(ctx, listingID)
}

func (s *MarketplaceService) GetListingByID(ctx context.Context, id int) (*models.MarketplaceListing, error) {
	// Получаем объявление
	listing, err := s.storage.GetListingByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Проверяем, пришел ли запрос от пользовательской части приложения (через GetListing API)
	// или из административной/служебной части (другие API)
	if ctx.Value("increment_views") == true {
		// Увеличиваем счетчик просмотров только для просмотра деталей объявления
		if err := s.storage.IncrementViewsCount(ctx, id); err != nil {
			// Логируем ошибку, но не прерываем выполнение
			log.Printf("Ошибка при увеличении счетчика просмотров для объявления %d: %v", id, err)
		}
	}

	return listing, nil
}

// GetListingBySlug получает объявление по slug
func (s *MarketplaceService) GetListingBySlug(ctx context.Context, slug string) (*models.MarketplaceListing, error) {
	// Получаем объявление по slug
	listing, err := s.storage.GetListingBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	// Проверяем, пришел ли запрос от пользовательской части приложения
	if ctx.Value("increment_views") == true {
		// Увеличиваем счетчик просмотров только для просмотра деталей объявления
		if err := s.storage.IncrementViewsCount(ctx, listing.ID); err != nil {
			// Логируем ошибку, но не прерываем выполнение
			log.Printf("Ошибка при увеличении счетчика просмотров для объявления %d: %v", listing.ID, err)
		}
	}

	return listing, nil
}

// IsSlugAvailable проверяет доступность slug
func (s *MarketplaceService) IsSlugAvailable(ctx context.Context, slug string, excludeID int) (bool, error) {
	return s.storage.IsSlugUnique(ctx, slug, excludeID)
}

// GenerateUniqueSlug генерирует уникальный slug на основе базового
func (s *MarketplaceService) GenerateUniqueSlug(ctx context.Context, baseSlug string, excludeID int) (string, error) {
	// Сначала проверяем исходный slug
	isUnique, err := s.storage.IsSlugUnique(ctx, baseSlug, excludeID)
	if err != nil {
		return "", err
	}

	if isUnique {
		return baseSlug, nil
	}

	// Если не уникален, пробуем с числовыми суффиксами
	for i := 2; i <= 99; i++ {
		candidateSlug := fmt.Sprintf("%s-%d", baseSlug, i)
		isUnique, err := s.storage.IsSlugUnique(ctx, candidateSlug, excludeID)
		if err != nil {
			return "", err
		}

		if isUnique {
			return candidateSlug, nil
		}
	}

	// Если все числа от 2 до 99 заняты, используем timestamp
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s-%d", baseSlug, timestamp), nil
}

func (s *MarketplaceService) GetOpenSearchRepository() (interface {
	SearchListings(ctx context.Context, params *search.SearchParams) (*search.SearchResult, error)
}, bool,
) {
	// Пытаемся привести хранилище к типу с нужным методом
	if osRepo, ok := s.storage.(interface {
		SearchListings(ctx context.Context, params *search.SearchParams) (*search.SearchResult, error)
	}); ok {
		return osRepo, true
	}
	return nil, false
}

// Методы для работы с атрибутами перенесены в attribute_admin.go:
// - GetAttributeRanges
// - GetCategoryAttributes
// - SaveListingAttributes

// GetCategorySuggestions возвращает предложения категорий на основе поискового запроса
func (s *MarketplaceService) GetCategorySuggestions(ctx context.Context, query string, size int) ([]models.CategorySuggestion, error) {
	log.Printf("Запрос предложений категорий: '%s'", query)

	// Проверка входных параметров
	if query == "" {
		return []models.CategorySuggestion{}, nil
	}

	// Выполняем SQL-запрос для поиска категорий, связанных с запросом
	sqlQuery := `
        WITH RECURSIVE category_tree AS (
            SELECT c.id, c.name, c.parent_id
            FROM marketplace_categories c
            WHERE 1=1
            
            UNION
            
            SELECT c.id, c.name, c.parent_id
            FROM marketplace_categories c
            JOIN category_tree t ON c.parent_id = t.id
        ),
        matching_categories AS (
            SELECT 
                c.id,
                c.name,
                (SELECT COUNT(*) FROM marketplace_listings ml 
                 WHERE ml.category_id = c.id 
                 AND ml.status = 'active') as listing_count,
                CASE WHEN LOWER(c.name) LIKE LOWER($1) THEN 100 ELSE 0 END +
                (SELECT COUNT(*) FROM marketplace_listings ml 
                 WHERE ml.category_id = c.id 
                 AND (LOWER(ml.title) LIKE LOWER($1) OR LOWER(ml.description) LIKE LOWER($1)) 
                 AND ml.status = 'active') as relevance
            FROM marketplace_categories c
            WHERE LOWER(c.name) LIKE LOWER($1)
            OR EXISTS (
                SELECT 1 FROM marketplace_listings ml 
                WHERE ml.category_id = c.id 
                AND (LOWER(ml.title) LIKE LOWER($1) OR LOWER(ml.description) LIKE LOWER($1))
                AND ml.status = 'active'
            )
        )
        SELECT id, name, listing_count
        FROM matching_categories
        WHERE listing_count > 0
        ORDER BY relevance DESC, listing_count DESC
        LIMIT $2
    `

	rows, err := s.storage.Query(ctx, sqlQuery, "%"+query+"%", size)
	if err != nil {
		log.Printf("Ошибка при выполнении запроса категорий: %v", err)
		return []models.CategorySuggestion{}, nil
	}
	defer func() {
		if err := rows.Close(); err != nil {
			// Логирование ошибки закрытия rows
			_ = err // Explicitly ignore error
		}
	}()

	var results []models.CategorySuggestion
	for rows.Next() {
		var suggestion models.CategorySuggestion

		if err := rows.Scan(&suggestion.ID, &suggestion.Name, &suggestion.ListingCount); err != nil {
			log.Printf("Ошибка сканирования категории: %v", err)
			continue
		}

		results = append(results, suggestion)
	}

	log.Printf("Найдено %d релевантных категорий для запроса '%s'", len(results), query)

	return results, nil
}

func (s *MarketplaceService) UpdateListing(ctx context.Context, listing *models.MarketplaceListing) error {
	// Получаем текущую версию объявления перед обновлением
	currentListing, err := s.storage.GetListingByID(ctx, listing.ID)
	if err != nil {
		return fmt.Errorf("ошибка получения текущего объявления: %w", err)
	}

	// Проверяем изменение цены
	if currentListing.Price != listing.Price {
		// Если цена изменилась, и есть метаданные о скидке
		if currentListing.Metadata != nil {
			if discount, ok := currentListing.Metadata["discount"].(map[string]interface{}); ok {
				if prevPrice, ok := discount["previous_price"].(float64); ok {
					// Обновляем процент скидки на основе новой цены
					if prevPrice > 0 && prevPrice > listing.Price {
						// Пересчитываем скидку
						discountPercent := int((prevPrice - listing.Price) / prevPrice * 100)

						// Обновляем метаданные
						if listing.Metadata == nil {
							listing.Metadata = make(map[string]interface{})
						}

						// Копируем существующую информацию о скидке
						if listing.Metadata["discount"] == nil {
							listing.Metadata["discount"] = make(map[string]interface{})
						}

						discountMap, _ := listing.Metadata["discount"].(map[string]interface{})
						discountMap["discount_percent"] = discountPercent
						discountMap["previous_price"] = prevPrice
						discountMap["has_price_history"] = true

						log.Printf("Обновление процента скидки для объявления %d: %d%% (старая цена: %.2f, новая цена: %.2f)",
							listing.ID, discountPercent, prevPrice, listing.Price)
					}
				}
			}
		}
	}

	// Вызываем существующий метод для обновления объявления в БД
	if err := s.storage.UpdateListing(ctx, listing); err != nil {
		return err
	}

	// Проверяем, изменились ли адресные поля
	addressChanged := currentListing.Location != listing.Location ||
		currentListing.City != listing.City ||
		currentListing.Country != listing.Country

	// Если адресные поля изменились, обновляем переводы
	if addressChanged {
		addressFields := make(map[string]string)
		if listing.Location != "" {
			addressFields["location"] = listing.Location
		}
		if listing.City != "" {
			addressFields["city"] = listing.City
		}
		if listing.Country != "" {
			addressFields["country"] = listing.Country
		}

		if len(addressFields) > 0 {
			// Определяем язык обновленных данных
			sourceLanguage := listing.OriginalLanguage
			if sourceLanguage == "" {
				// Если язык не указан, используем язык из контекста или русский по умолчанию
				if userLang, ok := ctx.Value("language").(string); ok && userLang != "" {
					sourceLanguage = userLang
				} else {
					sourceLanguage = "ru"
				}
			}

			// Определяем целевые языки
			targetLanguages := []string{"en", "ru", "sr"}
			filteredTargetLanguages := make([]string, 0, len(targetLanguages))
			for _, lang := range targetLanguages {
				if lang != sourceLanguage {
					filteredTargetLanguages = append(filteredTargetLanguages, lang)
				}
			}

			// Обновляем переводы адресов асинхронно
			go func(ctx context.Context) {
				if err := s.SaveAddressTranslations(ctx, listing.ID, addressFields, sourceLanguage, filteredTargetLanguages); err != nil {
					log.Printf("Error updating address translations for listing %d: %v", listing.ID, err)
				}
			}(ctx)
		}
	}

	// Получаем полное объявление со всеми связанными данными после обновления
	fullListing, err := s.storage.GetListingByID(ctx, listing.ID)
	if err != nil {
		log.Printf("Ошибка получения полного объявления для обновления индекса: %v", err)
		return nil
	}

	// Обновляем объявление в OpenSearch
	if err := s.storage.IndexListing(ctx, fullListing); err != nil {
		log.Printf("Ошибка обновления объявления в OpenSearch: %v", err)
	}

	return nil
}

func (s *MarketplaceService) SynchronizeDiscountData(ctx context.Context, listingID int) error {
	// Получаем данные из PostgreSQL
	listing, err := s.storage.GetListingByID(ctx, listingID)
	if err != nil {
		return fmt.Errorf("ошибка получения объявления: %w", err)
	}

	// Получаем историю цен
	priceHistory, err := s.storage.GetPriceHistory(ctx, listingID)
	if err != nil {
		log.Printf("Ошибка получения истории цен: %v", err)
		// Продолжаем выполнение метода даже при ошибке получения истории
	}

	// Проверяем на манипуляции с ценой, если есть история
	if len(priceHistory) > 1 {
		// Сортируем историю цен по дате
		sort.Slice(priceHistory, func(i, j int) bool {
			return priceHistory[i].EffectiveFrom.Before(priceHistory[j].EffectiveFrom)
		})

		// Проверяем на манипуляции с ценой
		isManipulation := false
		for i := 0; i < len(priceHistory)-1; i++ {
			// Если цена была повышена более чем на 50%, а затем снижена в течение 3 дней
			var nextEffectiveTo time.Time
			if priceHistory[i+1].EffectiveTo == nil {
				nextEffectiveTo = time.Now()
			} else {
				nextEffectiveTo = *priceHistory[i+1].EffectiveTo
			}

			duration := nextEffectiveTo.Sub(priceHistory[i+1].EffectiveFrom)

			if priceHistory[i].Price*1.5 < priceHistory[i+1].Price &&
				duration < 3*24*time.Hour &&
				i+2 < len(priceHistory) &&
				priceHistory[i+2].Price < priceHistory[i+1].Price {

				isManipulation = true
				log.Printf("Обнаружена манипуляция с ценой для объявления %d: повышение с %.2f до %.2f на %.1f часов, затем снижение до %.2f",
					listingID, priceHistory[i].Price, priceHistory[i+1].Price,
					duration.Hours(), priceHistory[i+2].Price)
				break
			}
		}

		if isManipulation {
			// В случае обнаружения манипуляции удаляем информацию о скидке
			if listing.Metadata != nil && listing.Metadata["discount"] != nil {
				delete(listing.Metadata, "discount")

				// Обновляем объявление в БД
				_, err := s.storage.Exec(ctx, `
                    UPDATE marketplace_listings
                    SET metadata = $1
                    WHERE id = $2
                `, listing.Metadata, listingID)
				if err != nil {
					log.Printf("Ошибка удаления метаданных о скидке: %v", err)
					return err
				}

				log.Printf("Удалена информация о скидке из-за обнаружения манипуляций с ценой для объявления %d", listingID)

				// Переиндексируем объявление без скидки и возвращаемся
				return s.storage.IndexListing(ctx, listing)
			}
		}

		// Находим максимальную цену в истории с учетом минимальной продолжительности
		var maxPrice float64
		var maxPriceDate time.Time
		minDuration := 24 * time.Hour // Минимальная продолжительность - 1 день

		for _, entry := range priceHistory {
			// Рассчитываем продолжительность действия цены
			var duration time.Duration
			if entry.EffectiveTo == nil {
				duration = time.Since(entry.EffectiveFrom)
			} else {
				duration = entry.EffectiveTo.Sub(entry.EffectiveFrom)
			}

			// Учитываем только цены, которые действовали достаточно долго
			if duration >= minDuration && entry.Price > maxPrice {
				maxPrice = entry.Price
				maxPriceDate = entry.EffectiveFrom
			}
		}

		// Если текущая цена ниже максимальной, создаем скидку
		if maxPrice > listing.Price && maxPrice > 0 {
			discountPercent := int((maxPrice - listing.Price) / maxPrice * 100)

			if discountPercent >= 5 { // Если скидка составляет не менее 5%
				// Если у объекта нет метаданных, создаем их
				if listing.Metadata == nil {
					listing.Metadata = make(map[string]interface{})
				}

				// Создаем информацию о скидке с установленным флагом has_price_history = true
				discount := map[string]interface{}{
					"discount_percent":  discountPercent,
					"previous_price":    maxPrice,
					"effective_from":    maxPriceDate.Format(time.RFC3339),
					"has_price_history": true,
				}

				// Добавляем информацию о скидке в метаданные
				listing.Metadata["discount"] = discount

				// Обновляем объявление в БД
				_, err := s.storage.Exec(ctx, `
                    UPDATE marketplace_listings
                    SET metadata = $1
                    WHERE id = $2
                `, listing.Metadata, listingID)
				if err != nil {
					log.Printf("Ошибка обновления метаданных: %v", err)
					return err
				}

				log.Printf("Создана информация о скидке для объявления %d из истории цен: %d%%, старая цена: %.2f",
					listingID, discountPercent, maxPrice)
			} else if listing.Metadata != nil && listing.Metadata["discount"] != nil {
				// Если скидка меньше 5%, удаляем информацию о скидке
				delete(listing.Metadata, "discount")

				// Обновляем объявление в БД
				_, err := s.storage.Exec(ctx, `
                        UPDATE marketplace_listings
                        SET metadata = $1
                        WHERE id = $2
                    `, listing.Metadata, listingID)
				if err != nil {
					log.Printf("Ошибка удаления метаданных о малой скидке: %v", err)
					return err
				}

				log.Printf("Удалена информация о малой скидке (%.1f%%) для объявления %d",
					float64(discountPercent), listingID)
			}
		} else if listing.Metadata != nil && listing.Metadata["discount"] != nil {
			// Если максимальная цена не найдена или текущая цена не ниже максимальной,
			// удаляем информацию о скидке
			delete(listing.Metadata, "discount")

			// Обновляем объявление в БД
			_, err := s.storage.Exec(ctx, `
                UPDATE marketplace_listings
                SET metadata = $1
                WHERE id = $2
            `, listing.Metadata, listingID)
			if err != nil {
				log.Printf("Ошибка удаления метаданных о неактуальной скидке: %v", err)
				return err
			}

			log.Printf("Удалена неактуальная информация о скидке для объявления %d", listingID)
		}
	}

	// Проверяем наличие данных о скидке в тексте описания (если нет в метаданных)
	if (listing.Metadata == nil || listing.Metadata["discount"] == nil) &&
		(strings.Contains(listing.Description, "СКИДКА") || strings.Contains(listing.Description, "СКИДКА!")) {

		discountRegex := regexp.MustCompile(`(\d+)%\s*СКИДКА`)
		matches := discountRegex.FindStringSubmatch(listing.Description)

		priceRegex := regexp.MustCompile(`Старая цена:\s*(\d+[\.,]?\d*)\s*RSD`)
		priceMatches := priceRegex.FindStringSubmatch(listing.Description)

		if len(matches) > 1 && len(priceMatches) > 1 {
			discountPercent, _ := strconv.Atoi(matches[1])
			oldPriceStr := strings.ReplaceAll(priceMatches[1], ",", ".")
			oldPrice, _ := strconv.ParseFloat(oldPriceStr, 64)

			// Проверяем реальность скидки
			calculatedDiscountPercent := int((oldPrice - listing.Price) / oldPrice * 100)
			if calculatedDiscountPercent < 0 || abs(calculatedDiscountPercent-discountPercent) > 5 {
				log.Printf("Объявленная скидка (%d%%) не соответствует расчетной (%d%%) для объявления %d, игнорируем",
					discountPercent, calculatedDiscountPercent, listingID)
				return nil
			}

			// Если у объекта нет метаданных, создаем их
			if listing.Metadata == nil {
				listing.Metadata = make(map[string]interface{})
			}

			// Создаем записи в истории цен
			// 1. Закрываем все предыдущие открытые записи истории цен
			if err := s.storage.ClosePriceHistoryEntry(ctx, listing.ID); err != nil {
				log.Printf("Ошибка закрытия прошлых записей истории цен: %v", err)
			}

			// 2. Создаем запись с предыдущей ценой, датированную неделю назад
			effectiveFrom := time.Now().AddDate(0, 0, -7)

			oldPriceEntry := &models.PriceHistoryEntry{
				ListingID:     listing.ID,
				Price:         oldPrice,
				EffectiveFrom: effectiveFrom,
				ChangeSource:  "parsed_from_description",
			}
			if err := s.storage.AddPriceHistoryEntry(ctx, oldPriceEntry); err != nil {
				log.Printf("Ошибка добавления старой цены в историю: %v", err)
			}

			// 3. Создаем запись с текущей ценой
			currentTime := time.Now()
			newPriceEntry := &models.PriceHistoryEntry{
				ListingID:     listing.ID,
				Price:         listing.Price,
				EffectiveFrom: currentTime,
				ChangeSource:  "parsed_from_description",
			}
			if err := s.storage.AddPriceHistoryEntry(ctx, newPriceEntry); err != nil {
				log.Printf("Ошибка добавления новой цены в историю: %v", err)
			}

			// Создаем информацию о скидке с установленным флагом has_price_history = true
			discount := map[string]interface{}{
				"discount_percent":  discountPercent,
				"previous_price":    oldPrice,
				"effective_from":    effectiveFrom.Format(time.RFC3339),
				"has_price_history": true,
			}

			// Добавляем информацию о скидке в метаданные
			listing.Metadata["discount"] = discount

			// Обновляем объявление в БД
			_, err := s.storage.Exec(ctx, `
                UPDATE marketplace_listings
                SET metadata = $1
                WHERE id = $2
            `, listing.Metadata, listing.ID)
			if err != nil {
				log.Printf("Ошибка обновления метаданных: %v", err)
				return err
			}

			log.Printf("Создана информация о скидке для объявления %d из описания: %d%%, старая цена: %.2f",
				listing.ID, discountPercent, oldPrice)
		}
	}

	// Переиндексация в OpenSearch
	return s.storage.IndexListing(ctx, listing)
}

// Вспомогательная функция для вычисления абсолютного значения
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (s *MarketplaceService) GetCategoryTree(ctx context.Context) ([]models.CategoryTreeNode, error) {
	// Если кеш не настроен, работаем напрямую со storage
	if s.cache == nil {
		return s.storage.GetCategoryTree(ctx)
	}

	// Получаем язык из контекста (по умолчанию "en")
	locale := "en"
	if lang, ok := ctx.Value(common.ContextKeyLocale).(string); ok && lang != "" {
		locale = lang
	}

	// Формируем ключ кеша
	cacheKey := cache.BuildCategoryTreeKey(locale, true)

	// Пытаемся получить из кеша
	var result []models.CategoryTreeNode
	err := s.cache.GetOrSet(ctx, cacheKey, &result, 6*time.Hour, func() (interface{}, error) {
		// Загружаем данные из БД
		return s.storage.GetCategoryTree(ctx)
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *MarketplaceService) RefreshCategoryListingCounts(ctx context.Context) error {
	_, err := s.storage.Exec(ctx, "REFRESH MATERIALIZED VIEW CONCURRENTLY category_listing_counts")
	return err
}

func (s *MarketplaceService) DeleteListing(ctx context.Context, id int, userID int) error {
	// Вызываем существующий метод для удаления объявления из БД
	if err := s.storage.DeleteListing(ctx, id, userID); err != nil {
		return err
	}

	// Удаляем объявление из OpenSearch
	if err := s.storage.DeleteListingIndex(ctx, fmt.Sprintf("%d", id)); err != nil {
		log.Printf("Ошибка удаления объявления из OpenSearch: %v", err)
		// Не возвращаем ошибку, чтобы не блокировать операцию, если OpenSearch недоступен
	}

	return nil
}

func (s *MarketplaceService) ProcessImage(file *multipart.FileHeader) (string, error) {
	// Получаем расширение файла
	ext := filepath.Ext(file.Filename)
	if ext == "" {
		// Если расширение отсутствует, определяем его по MIME-типу
		switch file.Header.Get("Content-Type") {
		case "image/jpeg", "image/jpg":
			ext = ".jpg"
		case "image/png":
			ext = ".png"
		case "image/gif":
			ext = ".gif"
		case "image/webp":
			ext = ".webp"
		default:
			ext = ".jpg" // По умолчанию используем .jpg
		}
	}

	// Генерируем уникальное имя файла с расширением
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	return fileName, nil
}

// In backend/internal/proj/marketplace/service/marketplace.go

func (s *MarketplaceService) UploadImage(ctx context.Context, file *multipart.FileHeader, listingID int, isMain bool) (*models.MarketplaceImage, error) {
	// Get file name
	fileName, err := s.ProcessImage(file)
	if err != nil {
		return nil, err
	}

	// Create object path - ensure no duplicate 'listings/' prefix
	objectName := fmt.Sprintf("%d/%s", listingID, fileName)

	// Open the file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer func() {
		if err := src.Close(); err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}()

	// Use FileStorage to upload
	fileStorage := s.storage.FileStorage()
	if fileStorage == nil {
		return nil, fmt.Errorf("file storage service not initialized")
	}

	// Upload to storage
	publicURL, err := fileStorage.UploadFile(ctx, objectName, src, file.Size, file.Header.Get("Content-Type"))
	if err != nil {
		return nil, fmt.Errorf("error uploading file: %w", err)
	}
	log.Printf("UploadImage: Изображение загружено в MinIO. objectName=%s, publicURL=%s", objectName, publicURL)

	// Create image information
	image := &models.MarketplaceImage{
		ListingID:     listingID,
		FilePath:      objectName,
		FileName:      file.Filename,
		FileSize:      int(file.Size),
		ContentType:   file.Header.Get("Content-Type"),
		IsMain:        isMain,
		StorageType:   "minio", // Явно указываем тип хранилища!
		StorageBucket: "listings",
		PublicURL:     publicURL,
	}
	log.Printf("UploadImage: Сохраняем информацию об изображении: ListingID=%d, FilePath=%s, StorageType=%s, PublicURL=%s",
		image.ListingID, image.FilePath, image.StorageType, image.PublicURL)
	// Save image information to database
	imageID, err := s.storage.AddListingImage(ctx, image)
	if err != nil {
		// Если не удалось сохранить информацию, удаляем файл
		if err := fileStorage.DeleteFile(ctx, objectName); err != nil {
			logger.Error().Err(err).Str("objectName", objectName).Msg("Failed to delete file from storage")
		}
		return nil, fmt.Errorf("error saving image information: %w", err)
	}
	log.Printf("UploadImage: Изображение успешно сохранено в базе данных с ID=%d", imageID)

	image.ID = imageID
	return image, nil
}

// DeleteImage удаляет изображение
func (s *MarketplaceService) DeleteImage(ctx context.Context, imageID int) error {
	// Получаем информацию об изображении
	image, err := s.storage.GetListingImageByID(ctx, imageID)
	if err != nil {
		return fmt.Errorf("ошибка получения информации об изображении: %w", err)
	}

	// Используем FileStorage для удаления файла
	fileStorage := s.storage.FileStorage()
	if fileStorage == nil {
		return fmt.Errorf("сервис файлового хранилища не инициализирован")
	}

	// Удаляем файл из хранилища
	err = fileStorage.DeleteFile(ctx, image.FilePath)
	if err != nil {
		log.Printf("Ошибка удаления файла из хранилища: %v", err)
		// Продолжаем выполнение для удаления записи из базы данных
	}

	// Удаляем информацию об изображении из базы данных
	err = s.storage.DeleteListingImage(ctx, imageID)
	if err != nil {
		return fmt.Errorf("ошибка удаления информации об изображении: %w", err)
	}

	return nil
}

// backend/internal/proj/marketplace/service/marketplace.go

// MigrateImagesToMinio мигрирует изображения из локального хранилища в MinIO
func (s *MarketplaceService) MigrateImagesToMinio(ctx context.Context) error {
	// Этот метод будем вызывать вручную при необходимости миграции

	// Получаем все изображения с типом хранилища 'local'
	query := `
		SELECT id, listing_id, file_path, file_name, file_size, content_type, is_main, created_at
		FROM marketplace_images
		WHERE storage_type = 'local' OR storage_type IS NULL
	`

	rows, err := s.storage.Query(ctx, query)
	if err != nil {
		return fmt.Errorf("ошибка получения изображений: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			// Логирование ошибки закрытия rows
			_ = err // Explicitly ignore error
		}
	}()

	var count int
	for rows.Next() {
		var image models.MarketplaceImage
		err := rows.Scan(
			&image.ID, &image.ListingID, &image.FilePath, &image.FileName,
			&image.FileSize, &image.ContentType, &image.IsMain, &image.CreatedAt,
		)
		if err != nil {
			log.Printf("Ошибка сканирования данных изображения: %v", err)
			continue
		}

		// Пропускаем, если путь к файлу пустой
		if image.FilePath == "" {
			continue
		}

		// Исключаем уже мигрированные изображения
		if strings.HasPrefix(image.FilePath, "listings/") {
			continue
		}

		// Создаем новый путь для изображения в MinIO
		newPath := fmt.Sprintf("listings/%d/%s", image.ListingID, filepath.Base(image.FilePath))

		// Security check: validate file path
		if strings.Contains(image.FilePath, "..") {
			log.Printf("Skipping image with invalid path: %s", image.FilePath)
			continue
		}

		// Открываем исходный файл
		localPath := fmt.Sprintf("./uploads/%s", image.FilePath)
		file, err := os.Open(localPath) // #nosec G304 -- path validated above
		if err != nil {
			log.Printf("Ошибка открытия файла %s: %v", localPath, err)
			continue
		}

		// Получаем размер файла
		fileInfo, err := file.Stat()
		if err != nil {
			log.Printf("Ошибка получения информации о файле %s: %v", localPath, err)
			if closeErr := file.Close(); closeErr != nil {
				log.Printf("Ошибка закрытия файла %s: %v", localPath, closeErr)
			}
			continue
		}

		// Загружаем файл в MinIO
		fileStorage := s.storage.FileStorage()
		if fileStorage == nil {
			if closeErr := file.Close(); closeErr != nil {
				log.Printf("Ошибка закрытия файла %s: %v", localPath, closeErr)
			}
			return fmt.Errorf("сервис файлового хранилища не инициализирован")
		}

		publicURL, err := fileStorage.UploadFile(ctx, newPath, file, fileInfo.Size(), image.ContentType)
		if closeErr := file.Close(); closeErr != nil {
			log.Printf("Ошибка закрытия файла %s: %v", localPath, closeErr)
		}
		if err != nil {
			log.Printf("Ошибка загрузки файла %s в MinIO: %v", localPath, err)
			continue
		}

		// Обновляем информацию об изображении в базе данных
		_, err = s.storage.Exec(ctx, `
			UPDATE marketplace_images
			SET file_path = $1, storage_type = 'minio', storage_bucket = 'listings', public_url = $2
			WHERE id = $3
		`, newPath, publicURL, image.ID)
		if err != nil {
			log.Printf("Ошибка обновления информации об изображении %d: %v", image.ID, err)
			continue
		}

		count++
		log.Printf("Успешно мигрировано изображение %d для объявления %d", image.ID, image.ListingID)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("ошибка итерации по изображениям: %w", err)
	}

	log.Printf("Миграция завершена. Всего мигрировано %d изображений", count)

	return nil
}

func (s *MarketplaceService) AddListingImage(ctx context.Context, image *models.MarketplaceImage) (int, error) {
	return s.storage.AddListingImage(ctx, image)
}

func (s *MarketplaceService) GetCategories(ctx context.Context) ([]models.MarketplaceCategory, error) {
	// Если кеш не настроен, работаем напрямую со storage
	if s.cache == nil {
		return s.storage.GetCategories(ctx)
	}

	// Получаем язык из контекста (по умолчанию "en")
	locale := "en"
	if lang, ok := ctx.Value(common.ContextKeyLocale).(string); ok && lang != "" {
		locale = lang
	}

	// Формируем ключ кеша
	cacheKey := cache.BuildCategoriesKey(locale)

	// Пытаемся получить из кеша
	var result []models.MarketplaceCategory
	err := s.cache.GetOrSet(ctx, cacheKey, &result, 6*time.Hour, func() (interface{}, error) {
		// Загружаем данные из БД
		return s.storage.GetCategories(ctx)
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *MarketplaceService) GetAllCategories(ctx context.Context) ([]models.MarketplaceCategory, error) {
	// Если кеш не настроен, работаем напрямую со storage
	if s.cache == nil {
		return s.storage.GetAllCategories(ctx)
	}

	// Получаем язык из контекста (по умолчанию "en")
	locale := "en"
	if lang, ok := ctx.Value(common.ContextKeyLocale).(string); ok && lang != "" {
		locale = lang
	}

	// Формируем ключ кеша для всех категорий (включая неактивные)
	cacheKey := cache.BuildCategoryTreeKey(locale, false)

	// Пытаемся получить из кеша
	var result []models.MarketplaceCategory
	err := s.cache.GetOrSet(ctx, cacheKey, &result, 6*time.Hour, func() (interface{}, error) {
		// Загружаем данные из БД
		return s.storage.GetAllCategories(ctx)
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *MarketplaceService) AddToFavorites(ctx context.Context, userID int, listingID int) error {
	return s.storage.AddToFavorites(ctx, userID, listingID)
}

func (s *MarketplaceService) RemoveFromFavorites(ctx context.Context, userID int, listingID int) error {
	return s.storage.RemoveFromFavorites(ctx, userID, listingID)
}

func (s *MarketplaceService) UpdateTranslation(ctx context.Context, translation *models.Translation) error {
	// Используем сервис перевода по умолчанию (Google Translate)
	// Передаём 0 в качестве userID, так как этот метод не имеет доступа к user_id
	return s.UpdateTranslationWithProvider(ctx, translation, GoogleTranslate, 0)
}

// SaveTranslation is an alias for UpdateTranslation for compatibility
func (s *MarketplaceService) SaveTranslation(ctx context.Context, entityType string, entityID int, language, fieldName, translatedText string, metadata map[string]interface{}) error {
	translation := &models.Translation{
		EntityType:     entityType,
		EntityID:       entityID,
		Language:       language,
		FieldName:      fieldName,
		TranslatedText: translatedText,
		IsVerified:     true,
		Metadata:       metadata,
	}
	return s.UpdateTranslation(ctx, translation)
}

// TranslateText переводит текст на указанный язык
func (s *MarketplaceService) TranslateText(ctx context.Context, text, sourceLanguage, targetLanguage string) (string, error) {
	if s.translationService == nil {
		return "", fmt.Errorf("translation service not available")
	}

	return s.translationService.Translate(ctx, text, sourceLanguage, targetLanguage)
}

// UpdateTranslationWithProvider обновляет перевод с использованием указанного провайдера
func (s *MarketplaceService) UpdateTranslationWithProvider(ctx context.Context, translation *models.Translation, provider TranslationProvider, userID int) error {
	// Проверяем, есть ли фабрика сервисов перевода
	if factory, ok := s.translationService.(TranslationFactoryInterface); ok {
		// Используем фабрику для обновления перевода с информацией о провайдере
		return factory.UpdateTranslation(ctx, translation, provider, userID)
	}

	// Если фабрики нет, используем прямой запрос к базе данных
	// Подготавливаем метаданные
	var metadataJSON []byte
	var err error

	if translation.Metadata == nil {
		translation.Metadata = map[string]interface{}{"provider": string(provider)}
	} else if _, exists := translation.Metadata["provider"]; !exists {
		translation.Metadata["provider"] = string(provider)
	}

	metadataJSON, err = json.Marshal(translation.Metadata)
	if err != nil {
		log.Printf("Ошибка сериализации метаданных: %v", err)
		metadataJSON = []byte("{}")
	}

	query := insertTranslationQuery

	var lastModifiedBy interface{}
	if userID > 0 {
		lastModifiedBy = userID
	} else {
		lastModifiedBy = nil
	}

	_, err = s.storage.Exec(ctx, query,
		translation.EntityType,
		translation.EntityID,
		translation.Language,
		translation.FieldName,
		translation.TranslatedText,
		translation.IsMachineTranslated,
		translation.IsVerified,
		metadataJSON,
		lastModifiedBy)

	return err
}

func (s *MarketplaceService) SearchListingsAdvanced(ctx context.Context, params *search.ServiceParams) (*search.ServiceResult, error) {
	log.Printf("Запрос поиска с параметрами: %+v", params)

	// Преобразуем ServiceParams в SearchParams для передачи в репозиторий
	searchParams := &search.SearchParams{
		Query:            params.Query,
		Page:             params.Page,
		Size:             params.Size,
		Aggregations:     params.Aggregations,
		Language:         params.Language,
		AttributeFilters: params.AttributeFilters,
		CategoryID:       nil,
		PriceMin:         nil,
		PriceMax:         nil,
		Condition:        params.Condition,
		City:             params.City,
		Country:          params.Country,
		StorefrontID:     nil,
		Sort:             params.Sort,
		SortDirection:    params.SortDirection,
		Distance:         params.Distance,
		CustomQuery:      nil,
		UseSynonyms:      params.UseSynonyms,
		Fuzziness:        params.Fuzziness,
	}
	// Преобразуем числовые значения в указатели для SearchParams
	if params.CategoryID != "" {
		if catID, err := strconv.Atoi(params.CategoryID); err == nil {
			searchParams.CategoryID = &catID
		}
	}
	if params.PriceMin > 0 {
		priceMin := params.PriceMin
		searchParams.PriceMin = &priceMin
	}
	if params.PriceMax > 0 {
		priceMax := params.PriceMax
		searchParams.PriceMax = &priceMax
	}
	if params.StorefrontID != "" {
		if storeID, err := strconv.Atoi(params.StorefrontID); err == nil {
			searchParams.StorefrontID = &storeID
		}
	}
	if params.Latitude != 0 && params.Longitude != 0 {
		searchParams.Location = &search.GeoLocation{
			Lat: params.Latitude,
			Lon: params.Longitude,
		}
	}

	// Добавляем особую обработку для поисковых запросов марка+модель
	if params.Query != "" && strings.Contains(params.Query, " ") {
		words := strings.Fields(params.Query)
		if len(words) > 1 {
			log.Printf("Обнаружен многословный запрос (%s), включаем специальную обработку для марки+модели", params.Query)
		}
	}

	// Выполняем поиск через Repository.SearchListings (buildSearchQuery)
	log.Printf("Выполняем поиск через Repository.SearchListings (buildSearchQuery)")
	searchResult, err := s.storage.SearchListingsOpenSearch(ctx, searchParams)
	if err != nil {
		log.Printf("Ошибка поиска в OpenSearch: %v", err)

		// Запасной вариант - поиск через PostgreSQL
		log.Printf("Выполняем запасной поиск через PostgreSQL")
		filters := map[string]string{
			"condition": params.Condition,
			"city":      params.City,
			"country":   params.Country,
			"sort_by":   params.Sort,
		}
		if params.CategoryID != "" {
			filters["category_id"] = params.CategoryID
		}
		if params.StorefrontID != "" {
			filters["storefront_id"] = params.StorefrontID
		}
		if params.PriceMin > 0 {
			filters["min_price"] = fmt.Sprintf("%g", params.PriceMin)
		}
		if params.PriceMax > 0 {
			filters["max_price"] = fmt.Sprintf("%g", params.PriceMax)
		}
		if params.Query != "" {
			filters["query"] = params.Query
		}

		// Если запрос похож на поиск марки+модели, добавляем специальную обработку
		words := strings.Fields(params.Query)
		if len(words) > 1 {
			log.Printf("Выполняем специальный запасной поиск для марки+модели: %v", words)

			// Пробуем найти совпадения для всех вариантов слов как марки/модели
			var results []models.MarketplaceListing
			var totalCount int64

			for i := range words {
				makeWord := words[i]

				// Создаем запрос только с маркой
				makeFilters := make(map[string]string)
				for k, v := range filters {
					makeFilters[k] = v
				}
				// Заменяем запрос только на слово марки
				makeFilters["query"] = makeWord

				// Находим по марке
				makeResults, makeTotal, makeErr := s.GetListings(ctx, makeFilters, params.Size, 0)
				if makeErr == nil && len(makeResults) > 0 {
					totalCount += makeTotal // Учитываем общее количество найденных

					// Фильтруем результаты по модели (другие слова)
					for _, result := range makeResults {
						// Проверяем, содержится ли модель в атрибутах или названии
						matched := false

						// Проверяем в названии
						resultTitle := strings.ToLower(result.Title)

						// Проверяем все слова кроме текущего (makeWord)
						for j := range words {
							if i == j {
								continue // Пропускаем текущее слово (уже проверили как марку)
							}

							modelWord := strings.ToLower(words[j])
							// Если модель найдена в названии
							if strings.Contains(resultTitle, modelWord) {
								matched = true
								break
							}

							// Проверяем атрибуты
							for _, attr := range result.Attributes {
								if attr.AttributeName == attributeNameModel && attr.TextValue != nil {
									attrValue := strings.ToLower(*attr.TextValue)
									if strings.Contains(attrValue, modelWord) {
										matched = true
										break
									}
								}
							}

							if matched {
								break
							}
						}

						// Если найдено совпадение, добавляем в результат
						if matched {
							results = append(results, result)
							if len(results) >= params.Size {
								break // Достигли лимита
							}
						}
					}
				}

				// Если нашли достаточно результатов, прекращаем поиск
				if len(results) >= params.Size {
					break
				}
			}

			// Если найдены результаты через специальный поиск
			if len(results) > 0 {
				// Преобразуем срез в указатели для результата
				listingPtrs := make([]*models.MarketplaceListing, len(results))
				for i := range results {
					listingPtrs[i] = &results[i]
				}

				return &search.ServiceResult{
					Items:      listingPtrs,
					Total:      int(totalCount),
					Page:       params.Page,
					Size:       params.Size,
					TotalPages: (int(totalCount) + params.Size - 1) / params.Size,
				}, nil
			}

			// Если специальный поиск не дал результатов, пробуем обычный
			log.Printf("Специальный поиск не дал результатов, пробуем обычный поиск")
		}

		// Выполняем стандартный поиск через PostgreSQL
		listings, total, err := s.GetListings(ctx, filters, params.Size, (params.Page-1)*params.Size)
		if err != nil {
			log.Printf("Ошибка стандартного поиска: %v", err)
			return nil, fmt.Errorf("ошибка поиска: %w", err)
		}
		listingPtrs := make([]*models.MarketplaceListing, len(listings))
		for i := range listings {
			listingPtrs[i] = &listings[i]
		}
		return &search.ServiceResult{
			Items:      listingPtrs,
			Total:      int(total),
			Page:       params.Page,
			Size:       params.Size,
			TotalPages: (int(total) + params.Size - 1) / params.Size,
		}, nil
	}

	log.Printf("Найдено %d объявлений", len(searchResult.Listings))
	for i, listing := range searchResult.Listings {
		log.Printf("Объявление %d: ID=%d, Название=%s, Координаты=%v,%v, Статус=%s",
			i+1, listing.ID, listing.Title, listing.Latitude, listing.Longitude, listing.Status)
	}
	log.Printf("Запрос поиска с параметрами сортировки: sort_by=%s, direction=%s", params.Sort, params.SortDirection)
	log.Printf("Итоговый параметр сортировки в запросе к OpenSearch: %s", params.Sort)
	// ЗДЕСЬ ДОБАВЛЯЕМ ДОПОЛНИТЕЛЬНУЮ СОРТИРОВКУ РЕЗУЛЬТАТОВ
	// Особая обработка для многословных запросов (марка+модель)
	if len(searchResult.Listings) > 0 && strings.Contains(params.Query, " ") {
		words := strings.Fields(params.Query)
		if len(words) > 1 {
			log.Printf("Выполняем дополнительную сортировку для многословного запроса '%s'", params.Query)

			// Преобразуем слова запроса в нижний регистр для сравнения
			lowerWords := make([]string, len(words))
			for i, word := range words {
				lowerWords[i] = strings.ToLower(word)
			}

			// Создаем функцию для вычисления оценки релевантности
			getRelevanceScore := func(listing *models.MarketplaceListing) int {
				score := 0

				// Получаем значения атрибутов make и model
				var makeValue, modelValue string
				for _, attr := range listing.Attributes {
					if attr.AttributeName == "make" && attr.TextValue != nil {
						makeValue = strings.ToLower(*attr.TextValue)
					}
					if attr.AttributeName == attributeNameModel && attr.TextValue != nil {
						modelValue = strings.ToLower(*attr.TextValue)
					}
				}

				// Добавляем логирование для отладки
				log.Printf("Листинг %d: make='%s', model='%s'", listing.ID, makeValue, modelValue)

				// Проверяем каждую пару слов как потенциальные марка+модель
				for _, word1 := range lowerWords {
					for _, word2 := range lowerWords {
						if word1 == word2 {
							continue
						}

						// Проверяем точное совпадение марка+модель (в любом порядке)
						if (word1 == makeValue && word2 == modelValue) ||
							(word2 == makeValue && word1 == modelValue) {
							// Очень высокий балл за точное совпадение обоих слов
							score += 1000
							log.Printf("  Точное совпадение марка+модель для '%s %s': +1000", word1, word2)
							break
						}

						// Проверяем модель на точное вхождение запроса
						if modelValue != "" && (modelValue == word1 || modelValue == word2) {
							score += 500
							log.Printf("  Точное совпадение модели для '%s': +500", modelValue)
						}

						// Проверяем точное совпадение только марки
						if word1 == makeValue || word2 == makeValue {
							score += 100
							log.Printf("  Точное совпадение марки для '%s': +100", makeValue)
						}

						// Проверяем частичное совпадение модели
						if modelValue != "" && (strings.Contains(modelValue, word1) ||
							strings.Contains(modelValue, word2)) {
							score += 20
							log.Printf("  Частичное совпадение модели для '%s': +20", modelValue)
						}
					}
				}

				// Проверяем наличие слов в заголовке
				title := strings.ToLower(listing.Title)
				for _, word := range lowerWords {
					if strings.Contains(title, word) {
						score += 10
						log.Printf("  Совпадение слова '%s' в заголовке: +10", word)
					}
				}

				log.Printf("  Финальный рейтинг для объявления %d: %d", listing.ID, score)
				return score
			}

			// Сортируем результаты по релевантности
			sort.Slice(searchResult.Listings, func(i, j int) bool {
				scoreI := getRelevanceScore(searchResult.Listings[i])
				scoreJ := getRelevanceScore(searchResult.Listings[j])
				return scoreI > scoreJ // Сортировка по убыванию релевантности
			})

			log.Printf("Выполнена дополнительная сортировка результатов по релевантности")

			// Выводим отсортированные результаты для проверки
			log.Printf("Отсортированные результаты:")
			for i, listing := range searchResult.Listings {
				var makeValue, modelValue string
				for _, attr := range listing.Attributes {
					if attr.AttributeName == "make" && attr.TextValue != nil {
						makeValue = *attr.TextValue
					}
					if attr.AttributeName == attributeNameModel && attr.TextValue != nil {
						modelValue = *attr.TextValue
					}
				}
				log.Printf("  %d. ID=%d, Название=%s, Марка=%s, Модель=%s",
					i+1, listing.ID, listing.Title, makeValue, modelValue)
			}
		}
	}

	// Применяем расширенные геофильтры если они заданы
	filteredListings := searchResult.Listings
	if params.AdvancedGeoFilters != nil && len(filteredListings) > 0 {
		log.Printf("Применяем расширенные геофильтры к %d объявлениям", len(filteredListings))

		// Получаем IDs всех найденных объявлений
		listingIDs := make([]string, len(filteredListings))
		for i, listing := range filteredListings {
			listingIDs[i] = strconv.Itoa(listing.ID)
		}

		// Применяем фильтры через GIS сервис
		filteredIDs, err := s.applyAdvancedGeoFilters(ctx, params.AdvancedGeoFilters, listingIDs)
		if err != nil {
			log.Printf("Ошибка применения расширенных геофильтров: %v", err)
			// Продолжаем без фильтрации в случае ошибки
		} else {
			// Фильтруем результаты по полученным ID
			filteredMap := make(map[string]bool)
			for _, id := range filteredIDs {
				filteredMap[id] = true
			}

			newFilteredListings := make([]*models.MarketplaceListing, 0, len(filteredIDs))
			for _, listing := range filteredListings {
				if filteredMap[strconv.Itoa(listing.ID)] {
					newFilteredListings = append(newFilteredListings, listing)
				}
			}

			log.Printf("После применения геофильтров осталось %d из %d объявлений",
				len(newFilteredListings), len(filteredListings))
			filteredListings = newFilteredListings
			searchResult.Total = len(newFilteredListings)
		}
	}

	result := &search.ServiceResult{
		Items:      filteredListings,
		Total:      searchResult.Total,
		Page:       params.Page,
		Size:       params.Size,
		TotalPages: (searchResult.Total + params.Size - 1) / params.Size,
		Took:       searchResult.Took,
	}

	if len(searchResult.Aggregations) > 0 {
		result.Facets = make(map[string][]search.Bucket)
		for key, buckets := range searchResult.Aggregations {
			result.Facets[key] = buckets
		}
	}

	if len(searchResult.Suggestions) > 0 {
		result.Suggestions = searchResult.Suggestions
	}

	log.Printf("Результаты поиска для атрибутов %v:", params.AttributeFilters)
	for i, listing := range searchResult.Listings {
		log.Printf("Объявление %d: ID=%d, Название=%s", i+1, listing.ID, listing.Title)

		// Добавляем отладочную информацию о атрибутах
		if len(listing.Attributes) > 0 {
			log.Printf("  Объявление %d имеет %d атрибутов:", listing.ID, len(listing.Attributes))
			for _, attr := range listing.Attributes {
				log.Printf("  Атрибут: name=%s, type=%s, value=%s",
					attr.AttributeName, attr.AttributeType, attr.DisplayValue)
			}
		} else {
			log.Printf("  Объявление %d не имеет атрибутов", listing.ID)
		}
	}

	return result, nil
}

// GetSuggestions возвращает предложения автодополнения
func (s *MarketplaceService) GetSuggestions(ctx context.Context, prefix string, size int) ([]string, error) {
	log.Printf("Запрос подсказок в сервисе: '%s'", prefix)

	// Проверка входных параметров
	if prefix == "" {
		return []string{}, nil
	}

	// Сначала пытаемся получить подсказки из OpenSearch
	suggestions, err := s.storage.SuggestListings(ctx, prefix, size)
	if err != nil || len(suggestions) == 0 {
		if err != nil {
			log.Printf("Ошибка при получении подсказок из OpenSearch: %v", err)
		}

		// Улучшенный SQL-запрос, включающий атрибуты
		query := `
        WITH attribute_suggestions AS (
            -- Подсказки из атрибутов текстового типа
            SELECT DISTINCT lav.text_value as value, 1 as priority
            FROM listing_attribute_values lav
            JOIN category_attributes ca ON lav.attribute_id = ca.id
            WHERE lav.text_value IS NOT NULL
              AND LOWER(lav.text_value) LIKE LOWER($1)
              AND ca.attribute_type IN ('text', 'select')
            
            UNION ALL
            
            -- Подсказки из атрибутов числового типа (преобразованные в строку)
            SELECT DISTINCT CAST(lav.numeric_value AS TEXT) as value, 2 as priority
            FROM listing_attribute_values lav
            JOIN category_attributes ca ON lav.attribute_id = ca.id
            WHERE lav.numeric_value IS NOT NULL
              AND CAST(lav.numeric_value AS TEXT) LIKE $1 || '%'
              AND ca.attribute_type = 'number'
        ),
        title_suggestions AS (
            -- Подсказки из заголовков объявлений
            SELECT DISTINCT title as value,
                   CASE WHEN LOWER(title) = LOWER($2) THEN 0
                        WHEN LOWER(title) LIKE LOWER($2 || '%') THEN 1
                        ELSE 2
                   END as priority,
                   length(title) as title_length
            FROM marketplace_listings 
            WHERE LOWER(title) LIKE LOWER($1) 
              AND status = 'active'
        )
        
        -- Объединяем все подсказки и отбираем лучшие
        SELECT value 
        FROM (
            SELECT value, priority, 0 as title_length FROM attribute_suggestions
            UNION ALL
            SELECT value, priority, title_length FROM title_suggestions
        ) combined
        ORDER BY priority, title_length
        LIMIT $3
        `
		rows, err := s.storage.Query(ctx, query, "%"+prefix+"%", prefix, size)
		if err != nil {
			log.Printf("Ошибка запасного SQL-запроса: %v", err)
			return []string{}, nil
		}
		defer func() {
			if err := rows.Close(); err != nil {
				// Логирование ошибки закрытия rows
				_ = err // Explicitly ignore error
			}
		}()

		var results []string
		for rows.Next() {
			var value string
			if err := rows.Scan(&value); err != nil {
				log.Printf("Ошибка сканирования строки: %v", err)
				continue
			}
			if value != "" && !contains(results, value) {
				results = append(results, value)
			}
		}

		log.Printf("Получено %d подсказок из базы данных", len(results))
		return results, nil
	}

	log.Printf("Получено %d подсказок из OpenSearch", len(suggestions))
	return suggestions, nil
}

// Вспомогательная функция для проверки наличия элемента в срезе
func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func (s *MarketplaceService) ReindexAllListings(ctx context.Context) error {
	return s.storage.ReindexAllListings(ctx)
}

func (s *MarketplaceService) Storage() storage.Storage {
	return s.storage
}

func (s *MarketplaceService) Service() *Service {
	// Возвращаем nil, так как мы больше не создаем вложенный Service
	return nil
}

func (s *MarketplaceService) GetPriceHistory(ctx context.Context, listingID int) ([]models.PriceHistoryEntry, error) {
	// Получаем историю цен из хранилища
	history, err := s.storage.GetPriceHistory(ctx, listingID)
	if err != nil {
		return nil, fmt.Errorf("error getting price history: %w", err)
	}

	return history, nil
}

// getSimilarStorefrontProducts находит похожие товары витрин для объявления принадлежащего витрине
func (s *MarketplaceService) getSimilarStorefrontProducts(ctx context.Context, listingID int, limit int) ([]*models.MarketplaceListing, error) {
	log.Printf("getSimilarStorefrontProducts: поиск похожих товаров витрин для объявления %d", listingID)

	// Получаем репозиторий для поиска товаров витрин
	productSearchInterface := s.storage.StorefrontProductSearch()
	if productSearchInterface == nil {
		log.Printf("Репозиторий поиска товаров витрин недоступен, используем запасной поиск")
		return s.getFallbackSimilarListings(ctx, listingID, limit)
	}

	// Приводим к нужному типу
	productSearchRepo, ok := productSearchInterface.(interface {
		SearchSimilarProducts(ctx context.Context, productID int, limit int) ([]*models.MarketplaceListing, error)
	})
	if !ok {
		log.Printf("Репозиторий не поддерживает метод SearchSimilarProducts, используем запасной поиск")
		return s.getFallbackSimilarListings(ctx, listingID, limit)
	}

	// Выполняем поиск похожих товаров витрин
	similarListings, err := productSearchRepo.SearchSimilarProducts(ctx, listingID, limit)
	if err != nil {
		log.Printf("Ошибка поиска похожих товаров витрин: %v, используем запасной поиск", err)
		return s.getFallbackSimilarListings(ctx, listingID, limit)
	}

	// ИСПРАВЛЕНИЕ: Дозагружаем полные данные объявлений с изображениями
	var enrichedListings []*models.MarketplaceListing
	for _, partialListing := range similarListings {
		// Загружаем полные данные объявления, включая изображения
		fullListing, err := s.GetListingByID(ctx, partialListing.ID)
		if err != nil {
			log.Printf("Ошибка загрузки полных данных объявления %d: %v", partialListing.ID, err)
			// Если не удалось загрузить полные данные, используем частичные
			enrichedListings = append(enrichedListings, partialListing)
			continue
		}

		// Сохраняем метаданные о скоре похожести из частичного объявления
		if partialListing.Metadata != nil {
			if fullListing.Metadata == nil {
				fullListing.Metadata = make(map[string]interface{})
			}
			if similarityScore, exists := partialListing.Metadata["similarity_score"]; exists {
				fullListing.Metadata["similarity_score"] = similarityScore
			}
		}

		enrichedListings = append(enrichedListings, fullListing)
	}

	log.Printf("Найдено %d похожих товаров витрин для объявления %d (с загруженными изображениями)", len(enrichedListings), listingID)
	return enrichedListings, nil
}

// getFallbackSimilarListings - запасной вариант поиска похожих объявлений через обычный marketplace поиск
func (s *MarketplaceService) getFallbackSimilarListings(ctx context.Context, listingID int, limit int) ([]*models.MarketplaceListing, error) {
	log.Printf("getFallbackSimilarListings: запасной поиск для объявления %d", listingID)

	// Получаем исходное объявление
	listing, err := s.GetListingByID(ctx, listingID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения объявления: %w", err)
	}

	if len(listing.Attributes) > 0 {
		log.Printf("Атрибуты объявления %d:", listing.ID)
		for _, attr := range listing.Attributes {
			log.Printf("  - %s: %s", attr.AttributeName, attr.DisplayValue)
		}
	} else {
		log.Printf("У объявления %d нет атрибутов", listing.ID)
	}

	// Создаем калькулятор похожести
	calculator := NewSimilarityCalculator(s.searchWeights)

	// Пытаемся найти похожие объявления с разными уровнями строгости
	var similarListings []*models.MarketplaceListing
	triesCount := 0
	maxTries := 4

	for triesCount < maxTries && len(similarListings) < limit {
		// Формируем параметры поиска для получения кандидатов
		params := s.buildAdvancedSearchParams(listing, limit*5, triesCount) // Получаем больше для фильтрации

		// Выполняем поиск похожих объявлений
		results, err := s.SearchListingsAdvanced(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("ошибка поиска похожих объявлений: %w", err)
		}

		// Фильтруем и сортируем результаты по похожести
		for _, candidate := range results.Items {
			if candidate.ID != listingID {
				// Проверяем, что кандидат еще не добавлен
				found := false
				for _, existing := range similarListings {
					if existing.ID == candidate.ID {
						found = true
						break
					}
				}
				if found {
					continue
				}

				// Вычисляем похожесть
				score, _ := calculator.CalculateSimilarity(ctx, listing, candidate)

				// Добавляем информацию о скоре в метаданные (для отладки)
				if candidate.Metadata == nil {
					candidate.Metadata = make(map[string]interface{})
				}
				candidate.Metadata["similarity_score"] = map[string]interface{}{
					"total":      score.TotalScore,
					"category":   score.CategoryScore,
					"attributes": score.AttributeScore,
					"price":      score.PriceScore,
					"location":   score.LocationScore,
					"text":       score.TextScore,
					"search_try": triesCount,
				}

				similarListings = append(similarListings, candidate)
			}
		}

		triesCount++
		log.Printf("Попытка %d: найдено %d похожих объявлений, всего собрано %d", triesCount, len(results.Items), len(similarListings))

		// Если найдено достаточно результатов, прекращаем поиск
		if len(similarListings) >= limit {
			break
		}
	}

	// Сортируем по убыванию похожести
	sort.Slice(similarListings, func(i, j int) bool {
		scoreI := similarListings[i].Metadata["similarity_score"].(map[string]interface{})["total"].(float64)
		scoreJ := similarListings[j].Metadata["similarity_score"].(map[string]interface{})["total"].(float64)
		return scoreI > scoreJ
	})

	// Ограничиваем количество результатов
	if len(similarListings) > limit {
		similarListings = similarListings[:limit]
	}

	log.Printf("Найдено %d похожих объявлений для листинга %d после %d попыток (запасной поиск)", len(similarListings), listingID, triesCount)

	return similarListings, nil
}

// applyAdvancedGeoFilters применяет расширенные геофильтры к списку объявлений
func (s *MarketplaceService) applyAdvancedGeoFilters(ctx context.Context, filters *search.AdvancedGeoFilters, listingIDs []string) ([]string, error) {
	// Формируем запрос к GIS сервису
	requestBody := map[string]interface{}{
		"filters":     filters,
		"listing_ids": listingIDs,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Вызываем GIS сервис
	// TODO: Получить URL из конфигурации
	gisURL := "http://localhost:3000/api/v1/gis/advanced/apply-filters"

	req, err := http.NewRequestWithContext(ctx, "POST", gisURL, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call GIS service: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.Error().Err(err).Msg("Failed to close response body")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GIS service returned status %d: %s", resp.StatusCode, string(body))
	}

	// Парсим ответ
	var response struct {
		Success bool `json:"success"`
		Data    struct {
			FilteredIDs []string `json:"filtered_ids"`
		} `json:"data"`
		Error string `json:"error,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode GIS response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("GIS service error: %s", response.Error)
	}

	return response.Data.FilteredIDs, nil
}

// SaveAddressTranslations сохраняет переводы адресных полей объявления
func (s *MarketplaceService) SaveAddressTranslations(ctx context.Context, listingID int, addressFields map[string]string, sourceLanguage string, targetLanguages []string) error {
	// Проверяем, есть ли адресные поля для перевода
	if len(addressFields) == 0 {
		return nil
	}

	// Проверяем наличие сервиса перевода
	if s.translationService == nil {
		log.Printf("Translation service not available for address translations")
		return nil
	}

	// Переводим адресные поля на все целевые языки
	translations, err := s.translationService.TranslateEntityFields(ctx, sourceLanguage, targetLanguages, addressFields)
	if err != nil {
		log.Printf("Error translating address fields for listing %d: %v", listingID, err)
		return fmt.Errorf("error translating address fields: %w", err)
	}

	// Сохраняем переводы в базу данных
	for language, fields := range translations {
		// Пропускаем исходный язык - он уже сохранен в основных полях объявления
		if language == sourceLanguage {
			continue
		}

		for fieldName, translatedText := range fields {
			// Сохраняем перевод для каждого поля
			err := s.SaveTranslation(ctx, "listing", listingID, language, fieldName, translatedText, map[string]interface{}{
				"source_language": sourceLanguage,
				"provider":        "google_translate",
				"is_address":      true,
			})
			if err != nil {
				log.Printf("Error saving translation for field %s to language %s: %v", fieldName, language, err)
				// Продолжаем с другими переводами, не прерываем процесс
			}
		}
	}

	log.Printf("Successfully saved address translations for listing %d", listingID)
	return nil
}
