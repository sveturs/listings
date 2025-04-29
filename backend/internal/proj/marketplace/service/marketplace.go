// backend/internal/proj/marketplace/service/marketplace.go
package service

import (
	"backend/internal/domain/models"
	"backend/internal/domain/search"
	"backend/internal/storage"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	//	"net/http"
	//	"net/url"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"os"    
)

type MarketplaceService struct {
	storage           storage.Storage
	translationService TranslationServiceInterface
}

func NewMarketplaceService(storage storage.Storage, translationService TranslationServiceInterface) MarketplaceServiceInterface {
	return &MarketplaceService{
		storage:           storage,
		translationService: translationService,
	}
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

	if listing.Attributes != nil && len(listing.Attributes) > 0 {
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
	// Получаем исходное объявление
	listing, err := s.GetListingByID(ctx, listingID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения объявления: %w", err)
	}

	// Формируем параметры поиска
	params := &search.ServiceParams{
		CategoryID: strconv.Itoa(listing.CategoryID),
		Size:       limit,
		Page:       1,
		Sort:       "date_desc", // Сортировка по дате по умолчанию
	}

	// Если есть атрибуты, можно использовать их для уточнения поиска
	if len(listing.Attributes) > 0 {
		attributeFilters := make(map[string]string)
		// Добавляем наиболее важные атрибуты для поиска похожих объявлений
		for _, attr := range listing.Attributes {
			// Выбираем только ключевые атрибуты для повышения релевантности
			if isKeyAttribute(attr.AttributeName) && attr.DisplayValue != "" {
				attributeFilters[attr.AttributeName] = attr.DisplayValue
			}
		}
		if len(attributeFilters) > 0 {
			params.AttributeFilters = attributeFilters
		}
	}

	// Выполняем поиск похожих объявлений
	results, err := s.SearchListingsAdvanced(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("ошибка поиска похожих объявлений: %w", err)
	}

	// Фильтруем результаты, убирая исходное объявление и ограничивая количество
	var similarListings []*models.MarketplaceListing
	for _, item := range results.Items {
		if item.ID != listingID {
			similarListings = append(similarListings, item)
		}
		if len(similarListings) >= limit {
			break
		}
	}

	return similarListings, nil
}

// Вспомогательная функция для определения ключевых атрибутов
func isKeyAttribute(attrName string) bool {
	// Список ключевых атрибутов для поиска похожих товаров
	keyAttributes := map[string]bool{
		"make":          true,
		"model":         true,
		"brand":         true,
		"category":      true,
		"type":          true,
		"rooms":         true,
		"property_type": true,
		"body_type":     true,
	}

	return keyAttributes[attrName]
}
func (s *MarketplaceService) GetSubcategories(ctx context.Context, parentID string, limit, offset int) ([]models.CategoryTreeNode, error) {
	var parentIDInt *int

	if parentID != "" {
		// Преобразуем строку в int
		id, err := strconv.Atoi(parentID)
		if err != nil {
			return nil, fmt.Errorf("invalid parent_id: %w", err)
		}
		parentIDInt = &id
	}

	return s.storage.GetSubcategories(ctx, parentIDInt, limit, offset)
}
func (s *MarketplaceService) GetListings(ctx context.Context, filters map[string]string, limit int, offset int) ([]models.MarketplaceListing, int64, error) {
	return s.storage.GetListings(ctx, filters, limit, offset)
}
func (s *MarketplaceService) GetFavoritedUsers(ctx context.Context, listingID int) ([]int, error) {
	return s.storage.GetFavoritedUsers(ctx, listingID)
}
func (s *MarketplaceService) GetListingByID(ctx context.Context, id int) (*models.MarketplaceListing, error) {
	return s.storage.GetListingByID(ctx, id)
}
func (s *MarketplaceService) GetOpenSearchRepository() (interface{
    SearchListings(ctx context.Context, params *search.SearchParams) (*search.SearchResult, error)
}, bool) {
    // Пытаемся привести хранилище к типу с нужным методом
    if osRepo, ok := s.storage.(interface{
        SearchListings(ctx context.Context, params *search.SearchParams) (*search.SearchResult, error)
    }); ok {
        return osRepo, true
    }
    return nil, false
}
func (s *MarketplaceService) GetAttributeRanges(ctx context.Context, categoryID int) (map[string]map[string]interface{}, error) {
    return s.storage.GetAttributeRanges(ctx, categoryID)
}
// GetCategoryAttributes получает атрибуты для указанной категории
func (s *MarketplaceService) GetCategoryAttributes(ctx context.Context, categoryID int) ([]models.CategoryAttribute, error) {
	attributes, err := s.storage.GetCategoryAttributes(ctx, categoryID)

	// Добавляем проверку и логирование JSON полей
	for _, attr := range attributes { // Заменяем i, attr на _, attr
		if attr.Options != nil {
			log.Printf("Attribute %s options (raw): %s", attr.Name, string(attr.Options))

			// Если нужно можно распарсить для проверки
			var options map[string]interface{}
			if err := json.Unmarshal(attr.Options, &options); err != nil {
				log.Printf("Error parsing options for attribute %s: %v", attr.Name, err)
			} else {
				log.Printf("Parsed options: %+v", options)
			}
		}
	}

	return attributes, err
}

func (s *MarketplaceService) SaveListingAttributes(ctx context.Context, listingID int, attributes []models.ListingAttributeValue) error {
	log.Printf("Saving %d attributes for listing %d", len(attributes), listingID)

	for i, attr := range attributes {
		log.Printf("Attribute %d: ID=%d, Name=%s, Type=%s",
			i, attr.AttributeID, attr.AttributeName, attr.AttributeType)
	}

	return s.storage.SaveListingAttributes(ctx, listingID, attributes)
}

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
	defer rows.Close()

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
		var minDuration = 24 * time.Hour // Минимальная продолжительность - 1 день

		for _, entry := range priceHistory {
			// Рассчитываем продолжительность действия цены
			var duration time.Duration
			if entry.EffectiveTo == nil {
				duration = time.Now().Sub(entry.EffectiveFrom)
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
			} else {
				// Если скидка меньше 5%, удаляем информацию о скидке
				if listing.Metadata != nil && listing.Metadata["discount"] != nil {
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
			oldPriceStr := strings.Replace(priceMatches[1], ",", ".", -1)
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
	return s.storage.GetCategoryTree(ctx)
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
    defer src.Close()

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
        StorageType:   "minio",  // Явно указываем тип хранилища!
        StorageBucket: "listings",
        PublicURL:     publicURL,
    }
    log.Printf("UploadImage: Сохраняем информацию об изображении: ListingID=%d, FilePath=%s, StorageType=%s, PublicURL=%s",
        image.ListingID, image.FilePath, image.StorageType, image.PublicURL)
    // Save image information to database
    imageID, err := s.storage.AddListingImage(ctx, image)
    if err != nil {
        // Если не удалось сохранить информацию, удаляем файл
        fileStorage.DeleteFile(ctx, objectName)
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
	defer rows.Close()
	
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

		
		// Открываем исходный файл
		localPath := fmt.Sprintf("./uploads/%s", image.FilePath)
		file, err := os.Open(localPath)
		if err != nil {
			log.Printf("Ошибка открытия файла %s: %v", localPath, err)
			continue
		}
		
		// Получаем размер файла
		fileInfo, err := file.Stat()
		if err != nil {
			log.Printf("Ошибка получения информации о файле %s: %v", localPath, err)
			file.Close()
			continue
		}
		
		// Загружаем файл в MinIO
		fileStorage := s.storage.FileStorage()
		if fileStorage == nil {
			file.Close()
			return fmt.Errorf("сервис файлового хранилища не инициализирован")
		}
		
		publicURL, err := fileStorage.UploadFile(ctx, newPath, file, fileInfo.Size(), image.ContentType)
		file.Close()
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
	return s.storage.GetCategories(ctx)
}

func (s *MarketplaceService) AddToFavorites(ctx context.Context, userID int, listingID int) error {
	return s.storage.AddToFavorites(ctx, userID, listingID)
}

func (s *MarketplaceService) RemoveFromFavorites(ctx context.Context, userID int, listingID int) error {
	return s.storage.RemoveFromFavorites(ctx, userID, listingID)
}
func (s *MarketplaceService) UpdateTranslation(ctx context.Context, translation *models.Translation) error {
	// Используем сервис перевода по умолчанию (Google Translate)
	return s.UpdateTranslationWithProvider(ctx, translation, GoogleTranslate)
}

// UpdateTranslationWithProvider обновляет перевод с использованием указанного провайдера
func (s *MarketplaceService) UpdateTranslationWithProvider(ctx context.Context, translation *models.Translation, provider TranslationProvider) error {
	// Проверяем, есть ли фабрика сервисов перевода
	if factory, ok := s.translationService.(TranslationFactoryInterface); ok {
		// Используем фабрику для обновления перевода с информацией о провайдере
		return factory.UpdateTranslation(ctx, translation, provider)
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
	
	query := `
        INSERT INTO translations (
            entity_type, entity_id, language, field_name,
            translated_text, is_machine_translated, is_verified, metadata
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        ON CONFLICT (entity_type, entity_id, language, field_name)
        DO UPDATE SET
            translated_text = EXCLUDED.translated_text,
            is_machine_translated = EXCLUDED.is_machine_translated,
            is_verified = EXCLUDED.is_verified,
            metadata = EXCLUDED.metadata,
            updated_at = CURRENT_TIMESTAMP
    `

	_, err = s.storage.Exec(ctx, query,
		translation.EntityType,
		translation.EntityID,
		translation.Language,
		translation.FieldName,
		translation.TranslatedText,
		translation.IsMachineTranslated,
		translation.IsVerified,
		metadataJSON)

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
								if attr.AttributeName == "model" && attr.TextValue != nil {
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
					if attr.AttributeName == "model" && attr.TextValue != nil {
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
					if attr.AttributeName == "model" && attr.TextValue != nil {
						modelValue = *attr.TextValue
					}
				}
				log.Printf("  %d. ID=%d, Название=%s, Марка=%s, Модель=%s",
					i+1, listing.ID, listing.Title, makeValue, modelValue)
			}
		}
	}

	result := &search.ServiceResult{
		Items:      searchResult.Listings,
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
		defer rows.Close()

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

func (s *MarketplaceService) GetPriceHistory(ctx context.Context, listingID int) ([]models.PriceHistoryEntry, error) {
	// Получаем историю цен из хранилища
	history, err := s.storage.GetPriceHistory(ctx, listingID)
	if err != nil {
		return nil, fmt.Errorf("error getting price history: %w", err)
	}

	return history, nil
}
