// backend/internal/proj/c2c/service/marketplace_search.go
package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"backend/internal/domain/models"
	"backend/internal/domain/search"
	"backend/internal/logger"
)

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

	// log.Printf("Объявление %d является обычным объявлением маркетплейса - используем стандартный поиск", listingID)

	// Закомментировано для снижения шума в логах
	// if len(listing.Attributes) > 0 {
	// 	log.Printf("Атрибуты объявления %d:", listing.ID)
	// 	for _, attr := range listing.Attributes {
	// 		log.Printf("  - %s: %s", attr.AttributeName, attr.DisplayValue)
	// 	}
	// } else {
	// 	log.Printf("У объявления %d нет атрибутов", listing.ID)
	// }

	// Создаем калькулятор похожести
	calculator := NewSimilarityCalculator(s.searchWeights)

	// Пытаемся найти похожие объявления с разными уровнями строгости
	var similarListings []*models.MarketplaceListing
	seenIDs := make(map[int]bool) // Для эффективной дедупликации
	triesCount := 0
	maxTries := 4

	for triesCount < maxTries && len(similarListings) < limit {
		// Формируем параметры поиска для получения кандидатов
		params := s.buildAdvancedSearchParams(ctx, listing, limit*5, triesCount) // Получаем больше для фильтрации

		// Выполняем поиск похожих объявлений
		results, err := s.SearchListingsAdvanced(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("ошибка поиска похожих объявлений: %w", err)
		}

		// Фильтруем и сортируем результаты по похожести
		for _, candidate := range results.Items {
			// Пропускаем текущее объявление и дубликаты
			if candidate.ID == listingID || seenIDs[candidate.ID] {
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
			seenIDs[candidate.ID] = true // Помечаем как обработанный
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

func (s *MarketplaceService) buildAdvancedSearchParams(ctx context.Context, listing *models.MarketplaceListing, size int, tryNumber int) *search.ServiceParams {
	params := &search.ServiceParams{
		Size:             size,
		Page:             1,
		Sort:             "date_desc",
		StorefrontFilter: "include_b2c", // ИСПРАВЛЕНИЕ: включаем товары витрин для поиска похожих товаров
	}

	// ИСПРАВЛЕНИЕ: Расширяем поиск по категориям в зависимости от попытки
	switch tryNumber {
	case 0, 1:
		// Первые 2 попытки - ищем в той же категории
		params.CategoryID = strconv.Itoa(listing.CategoryID)
		log.Printf("Попытка %d: поиск в категории %d", tryNumber, listing.CategoryID)
	case 2:
		// Третья попытка - ищем в родительской категории (если есть)
		if categoryID, err := s.getParentCategoryID(ctx, listing.CategoryID); err == nil && categoryID > 0 {
			params.CategoryID = strconv.Itoa(categoryID)
			log.Printf("Попытка %d: поиск в родительской категории %d (исходная %d)", tryNumber, categoryID, listing.CategoryID)
		} else {
			params.CategoryID = strconv.Itoa(listing.CategoryID)
			log.Printf("Попытка %d: родительская категория не найдена, поиск в исходной категории %d", tryNumber, listing.CategoryID)
		}
	default:
		// Последние попытки - ищем без ограничения по категориям
		log.Printf("Попытка %d: поиск без ограничения по категориям", tryNumber)
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
			// Последняя попытка: очень широкий диапазон ±400%
			params.PriceMin = listing.Price * 0.05
			params.PriceMax = listing.Price * 5.0
			log.Printf("Попытка %d: диапазон цен %.2f - %.2f (±400%%)", tryNumber, params.PriceMin, params.PriceMax)
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
		StorefrontFilter: params.StorefrontFilter,
		DocumentType:     params.DocumentType,
	}
	// Преобразуем числовые значения в указатели для SearchParams
	// Обрабатываем массив категорий
	if len(params.CategoryIDs) > 0 {
		categoryIDs := make([]int, 0, len(params.CategoryIDs))
		for _, catIDStr := range params.CategoryIDs {
			if catID, err := strconv.Atoi(catIDStr); err == nil {
				categoryIDs = append(categoryIDs, catID)
			}
		}
		searchParams.CategoryIDs = categoryIDs
	} else if params.CategoryID != "" {
		// Если массив не задан, используем единичную категорию
		if catID, err := strconv.Atoi(params.CategoryID); err == nil {
			searchParams.CategoryID = &catID
			// Также добавляем в массив для унифицированной обработки
			searchParams.CategoryIDs = []int{catID}
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
            FROM c2c_listings 
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

// GetUnifiedSuggestions возвращает структурированные подсказки с типами
func (s *MarketplaceService) GetUnifiedSuggestions(ctx context.Context, params *models.SuggestionRequestParams) ([]models.UnifiedSuggestion, error) {
	log.Printf("Запрос унифицированных подсказок: query='%s', types=%v, limit=%d", params.Query, params.Types, params.Limit)

	var results []models.UnifiedSuggestion

	// Если не указаны типы, возвращаем все типы по умолчанию
	types := params.Types
	if len(types) == 0 {
		types = []string{"queries", "categories", "products"}
	}

	// Распределяем лимит между типами
	limitPerType := params.Limit / len(types)
	if limitPerType < 1 {
		limitPerType = 1
	}

	// 1. Queries (популярные поисковые запросы)
	if contains(types, "queries") {
		queryResults := s.getQuerySuggestions(ctx, params.Query, limitPerType)
		results = append(results, queryResults...)
	}

	// 2. Categories (категории)
	if contains(types, "categories") {
		categoryResults := s.getCategorySuggestionsUnified(ctx, params.Query, limitPerType)
		results = append(results, categoryResults...)
	}

	// 3. Products (товары)
	if contains(types, "products") {
		productResults := s.getProductSuggestionsUnified(ctx, params.Query, limitPerType)
		results = append(results, productResults...)
	}

	// Ограничиваем общий размер результата
	if len(results) > params.Limit {
		results = results[:params.Limit]
	}

	log.Printf("Возвращено %d унифицированных подсказок", len(results))
	return results, nil
}

// getQuerySuggestions возвращает популярные поисковые запросы
func (s *MarketplaceService) getQuerySuggestions(ctx context.Context, query string, limit int) []models.UnifiedSuggestion {
	// Получаем популярные запросы из таблицы search_queries
	sqlQuery := `
		SELECT query, COUNT(*) as search_count 
		FROM search_queries 
		WHERE LOWER(query) LIKE LOWER($1) 
		  AND query != $2
		GROUP BY query 
		ORDER BY search_count DESC, LENGTH(query) ASC 
		LIMIT $3`

	rows, err := s.storage.Query(ctx, sqlQuery, query+"%", query, limit)
	if err != nil {
		log.Printf("Ошибка получения популярных запросов: %v", err)
		return []models.UnifiedSuggestion{}
	}
	defer func() { _ = rows.Close() }()

	var suggestions []models.UnifiedSuggestion
	for rows.Next() {
		var q string
		var count int
		if err := rows.Scan(&q, &count); err != nil {
			continue
		}

		suggestions = append(suggestions, models.UnifiedSuggestion{
			Type:  "query",
			Value: q,
			Label: q,
			Count: &count,
		})
	}

	return suggestions
}

// getCategorySuggestionsUnified возвращает подходящие категории
func (s *MarketplaceService) getProductSuggestionsUnified(ctx context.Context, query string, limit int) []models.UnifiedSuggestion {
	// Поиск товаров по названию из обеих таблиц (c2c_listings и b2c_products)
	sqlQuery := `
		SELECT id, title, price, category_name, image_url, storefront_id, storefront_name, storefront_slug, source_type
		FROM (
			-- Marketplace listings
			SELECT ml.id, ml.title, ml.price, mc.name as category_name,
			       COALESCE(mi.public_url, '') as image_url,
			       ml.storefront_id, s.name as storefront_name, s.slug as storefront_slug,
			       'marketplace' as source_type,
			       ml.views_count, ml.created_at
			FROM c2c_listings ml
			LEFT JOIN c2c_categories mc ON ml.category_id = mc.id
			LEFT JOIN c2c_images mi ON ml.id = mi.listing_id AND mi.is_main = true
			LEFT JOIN b2c_stores s ON ml.storefront_id = s.id
			WHERE LOWER(ml.title) LIKE LOWER($1)
			  AND ml.status = 'active'

			UNION ALL

			-- Storefront products
			SELECT sp.id, sp.name as title, sp.price, mc.name as category_name,
			       COALESCE(spi.image_url, '') as image_url,
			       sp.storefront_id, s.name as storefront_name, s.slug as storefront_slug,
			       'b2c_store' as source_type,
			       0 as views_count, sp.created_at
			FROM b2c_products sp
			LEFT JOIN c2c_categories mc ON sp.category_id = mc.id
			LEFT JOIN b2c_product_images spi ON sp.id = spi.storefront_product_id AND spi.is_default = true
			JOIN b2c_stores s ON sp.storefront_id = s.id
			WHERE LOWER(sp.name) LIKE LOWER($1)
		) combined
		ORDER BY views_count DESC, created_at DESC
		LIMIT $2`

	rows, err := s.storage.Query(ctx, sqlQuery, "%"+query+"%", limit)
	if err != nil {
		log.Printf("Ошибка получения товаров: %v", err)
		return []models.UnifiedSuggestion{}
	}
	defer func() { _ = rows.Close() }()

	var suggestions []models.UnifiedSuggestion
	for rows.Next() {
		var id int
		var title string
		var price float64
		var categoryName, imageURL, sourceType string
		var storefrontID *int
		var storefrontName, storefrontSlug *string

		if err := rows.Scan(&id, &title, &price, &categoryName, &imageURL,
			&storefrontID, &storefrontName, &storefrontSlug, &sourceType); err != nil {
			continue
		}

		metadata := &models.UnifiedSuggestionMeta{
			Price:      &price,
			Category:   &categoryName,
			SourceType: strPtr(sourceType),
		}

		if imageURL != "" {
			metadata.Image = &imageURL
		}

		if storefrontID != nil {
			metadata.StorefrontID = storefrontID
			metadata.Storefront = storefrontName
			metadata.StorefrontSlug = storefrontSlug
		}

		suggestions = append(suggestions, models.UnifiedSuggestion{
			Type:      "product",
			Value:     title,
			Label:     title,
			ProductID: &id,
			Metadata:  metadata,
		})
	}

	return suggestions
}

func (s *MarketplaceService) ReindexAllListings(ctx context.Context) error {
	return s.storage.ReindexAllListings(ctx)
}

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

	// ИСПРАВЛЕНИЕ: Если найдено мало похожих товаров витрин, дополняем обычными листингами
	if len(similarListings) < limit/2 {
		log.Printf("Найдено только %d похожих товаров витрин, дополняем обычными листингами", len(similarListings))
		fallbackListings, fallbackErr := s.getFallbackSimilarListings(ctx, listingID, limit-len(similarListings))
		if fallbackErr == nil {
			log.Printf("Добавляем %d обычных похожих листингов", len(fallbackListings))
			similarListings = append(similarListings, fallbackListings...)
		} else {
			log.Printf("Ошибка получения дополнительных листингов: %v", fallbackErr)
		}
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
	seenIDs := make(map[int]bool) // Для эффективной дедупликации
	triesCount := 0
	maxTries := 4

	for triesCount < maxTries && len(similarListings) < limit {
		// Формируем параметры поиска для получения кандидатов
		params := s.buildAdvancedSearchParams(ctx, listing, limit*5, triesCount) // Получаем больше для фильтрации

		// Выполняем поиск похожих объявлений
		results, err := s.SearchListingsAdvanced(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("ошибка поиска похожих объявлений: %w", err)
		}

		// Фильтруем и сортируем результаты по похожести
		for _, candidate := range results.Items {
			// Пропускаем текущее объявление и дубликаты
			if candidate.ID == listingID || seenIDs[candidate.ID] {
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
			seenIDs[candidate.ID] = true // Помечаем как обработанный
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
