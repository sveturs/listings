package service

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"backend/internal/domain/models"
	"backend/internal/domain/search"
	storefrontOpenSearch "backend/internal/proj/storefronts/storage/opensearch"

	"github.com/rs/zerolog/log"
)

type SuggestionType string

const (
	SuggestionTypeQuery    SuggestionType = "query"
	SuggestionTypeCategory SuggestionType = "category"
	SuggestionTypeProduct  SuggestionType = "product"
)

type SuggestionItem struct {
	Type       SuggestionType         `json:"type"`
	Value      string                 `json:"value"`
	Label      string                 `json:"label"`
	Count      int                    `json:"count,omitempty"`
	CategoryID int                    `json:"category_id,omitempty"`
	ProductID  int                    `json:"product_id,omitempty"`
	Icon       string                 `json:"icon,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

func (s *MarketplaceService) GetEnhancedSuggestions(
	ctx context.Context,
	query string,
	limit int,
	types string,
) ([]SuggestionItem, error) {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return []SuggestionItem{}, nil
	}

	var suggestions []SuggestionItem
	typesList := strings.Split(types, ",")

	log.Info().
		Str("query", query).
		Strs("types", typesList).
		Int("limit", limit).
		Msg("GetEnhancedSuggestions called")

	// 1. Популярные поисковые запросы
	if containsString(typesList, "queries") {
		querySuggestions, err := s.getQuerySuggestions(ctx, query)
		if err == nil {
			suggestions = append(suggestions, querySuggestions...)
		}
	}

	// 2. Товары/Объявления (получаем раньше, чтобы знать категории товаров)
	var productCategories []string
	if containsString(typesList, "products") {
		// Используем fuzzy search для поиска товаров
		productSuggestions, err := s.getProductSuggestionsWithFuzzy(ctx, query)
		if err == nil {
			suggestions = append(suggestions, productSuggestions...)
			// Собираем уникальные категории найденных товаров
			categorySet := make(map[string]bool)
			for _, prod := range productSuggestions {
				if cat, ok := prod.Metadata["category"].(string); ok && cat != "" {
					categorySet[cat] = true
					log.Debug().Str("category", cat).Msg("Found product category")
				}
			}
			for cat := range categorySet {
				productCategories = append(productCategories, cat)
			}
			log.Info().
				Str("query", query).
				Int("products_found", len(productSuggestions)).
				Strs("categories_found", productCategories).
				Msg("Product search completed")
		}
	}

	// 3. Категории (ВСЕГДА добавляем категории из найденных товаров)
	var productCategorySuggestions []SuggestionItem
	if len(productCategories) > 0 {
		log.Info().
			Strs("product_categories_to_add", productCategories).
			Msg("Processing product categories")

		// Получаем информацию о категориях из БД
		allCategories, _ := s.storage.GetCategories(ctx)
		log.Info().Int("all_categories_count", len(allCategories)).Msg("Loaded all categories from DB")

		for _, catName := range productCategories {
			log.Debug().
				Str("looking_for", catName).
				Msg("Searching for product category in DB")
			found := false
			for _, dbCat := range allCategories {
				// Проверяем по имени и по переводам
				nameMatch := strings.EqualFold(dbCat.Name, catName)

				// Проверяем также по переводам
				if !nameMatch && dbCat.Translations != nil {
					for _, translation := range dbCat.Translations {
						if strings.EqualFold(translation, catName) {
							nameMatch = true
							break
						}
					}
				}

				if nameMatch {
					iconStr := ""
					if dbCat.Icon != nil {
						iconStr = *dbCat.Icon
					}
					// Сохраняем категории из товаров отдельно
					productCategorySuggestions = append(productCategorySuggestions, SuggestionItem{
						Type:       SuggestionTypeCategory,
						Value:      dbCat.Slug,
						Label:      dbCat.Name,
						CategoryID: dbCat.ID,
						Icon:       iconStr,
						Count:      dbCat.ListingCount,
						Metadata: map[string]interface{}{
							"parent_id":     dbCat.ParentID,
							"from_products": true, // Помечаем, что категория из товаров
						},
					})
					found = true
					log.Info().
						Str("category_name", dbCat.Name).
						Int("category_id", dbCat.ID).
						Str("slug", dbCat.Slug).
						Msg("Added product category suggestion")
					break
				}
			}
			if !found {
				log.Warn().Str("category_name", catName).Msg("Category not found in DB")
			}
		}
	}

	// Также ищем категории по запросу
	if containsString(typesList, "categories") {
		categorySuggestions, err := s.getCategorySuggestionsEnhanced(ctx, query, productCategories)
		if err == nil {
			suggestions = append(suggestions, categorySuggestions...)
		}
	}

	// Добавляем категории из товаров в начало
	if len(productCategorySuggestions) > 0 {
		// Вставляем категории из товаров в начало списка suggestions
		suggestions = append(productCategorySuggestions, suggestions...)
		log.Info().
			Int("product_categories_count", len(productCategorySuggestions)).
			Int("total_suggestions_before_rank", len(suggestions)).
			Msg("Added product category suggestions")
	}

	// Сортируем и убеждаемся, что категории из товаров попадают в результат
	suggestions = s.rankSuggestionsWithProductCategories(ctx, suggestions, query, productCategories)

	// Логируем финальный результат
	log.Info().
		Str("query", query).
		Int("suggestions_count", len(suggestions)).
		Int("limit", limit).
		Interface("final_suggestions", func() []string {
			result := make([]string, 0, len(suggestions))
			for _, s := range suggestions {
				result = append(result, fmt.Sprintf("%s:%s", s.Type, s.Label))
			}
			return result
		}()).
		Msg("Final suggestions")

	if len(suggestions) > limit {
		suggestions = suggestions[:limit]
	}

	return suggestions, nil
}

func (s *MarketplaceService) getQuerySuggestions(
	ctx context.Context,
	query string,
) ([]SuggestionItem, error) {
	// Получаем популярные запросы из БД через типизированный метод
	popularQueriesRaw, err := s.storage.GetPopularSearchQueries(ctx, query, 5)
	if err != nil {
		return nil, err
	}

	suggestions := make([]SuggestionItem, 0, len(popularQueriesRaw))
	for _, pqRaw := range popularQueriesRaw {
		// Преобразуем interface{} в SearchQuery
		if pq, ok := pqRaw.(SearchQuery); ok {
			suggestions = append(suggestions, SuggestionItem{
				Type:  SuggestionTypeQuery,
				Value: pq.Query,
				Label: pq.Query,
				Count: pq.SearchCount,
				Metadata: map[string]interface{}{
					"last_searched": pq.LastSearched,
				},
			})
		}
	}

	return suggestions, nil
}

// getCategorySuggestions - удалена, используется getCategorySuggestionsEnhanced напрямую

func (s *MarketplaceService) getCategorySuggestionsEnhanced(
	ctx context.Context,
	query string,
	productCategories []string,
) ([]SuggestionItem, error) {
	// Ищем категории по имени с поддержкой fuzzy matching
	categories, err := s.storage.SearchCategories(ctx, query, 10)
	if err != nil {
		return nil, err
	}

	// Если нет точных совпадений, пробуем поиск с опечатками
	if len(categories) == 0 {
		// Пробуем найти категории, которые содержат части запроса
		allCategories, _ := s.storage.GetCategories(ctx)
		queryLower := strings.ToLower(query)

		for _, cat := range allCategories {
			nameLower := strings.ToLower(cat.Name)
			// Проверяем частичное совпадение или fuzzy matching
			if strings.Contains(nameLower, queryLower) ||
				levenshteinDistance(queryLower, nameLower) <= 2 {
				categories = append(categories, cat)
				if len(categories) >= 5 {
					break
				}
			}
		}
	}

	// Если есть категории из найденных товаров, добавляем их тоже
	if len(productCategories) > 0 {
		allCategories, _ := s.storage.GetCategories(ctx)
		categoryMap := make(map[string]bool)

		// Помечаем уже найденные категории
		for _, cat := range categories {
			categoryMap[strings.ToLower(cat.Name)] = true
		}

		// Добавляем категории товаров, если их еще нет
		for _, catName := range productCategories {
			catNameLower := strings.ToLower(catName)
			if !categoryMap[catNameLower] {
				// Ищем категорию по имени (case-insensitive)
				for _, cat := range allCategories {
					if strings.EqualFold(cat.Name, catName) {
						categories = append(categories, cat)
						categoryMap[catNameLower] = true
						break
					}
				}
			}
		}
	}

	suggestions := make([]SuggestionItem, 0, len(categories))
	for _, cat := range categories {
		iconStr := ""
		if cat.Icon != nil {
			iconStr = *cat.Icon
		}
		suggestions = append(suggestions, SuggestionItem{
			Type:       SuggestionTypeCategory,
			Value:      cat.Slug,
			Label:      cat.Name,
			CategoryID: cat.ID,
			Icon:       iconStr,
			Count:      cat.ListingCount,
			Metadata: map[string]interface{}{
				"parent_id": cat.ParentID,
			},
		})
	}

	// Ограничиваем до 5 категорий
	if len(suggestions) > 5 {
		suggestions = suggestions[:5]
	}

	return suggestions, nil
}

// getProductSuggestions - удалена, используется getProductSuggestionsWithFuzzy напрямую

func (s *MarketplaceService) getProductSuggestionsWithFuzzy(
	ctx context.Context,
	query string,
) ([]SuggestionItem, error) {
	suggestions := []SuggestionItem{}

	// 1. Сначала пробуем точный поиск
	searchParams := &search.ServiceParams{
		Query: query,
		Size:  5,
		Page:  1,
	}

	marketplaceResults, err := s.SearchListingsAdvanced(ctx, searchParams)

	// 2. Если точный поиск не дал результатов, используем fuzzy
	if (err != nil || marketplaceResults == nil || len(marketplaceResults.Items) == 0) && len(query) > 3 {
		// Пробуем с fuzzy search для опечаток
		searchParams.Fuzziness = "AUTO"
		marketplaceResults, err = s.SearchListingsAdvanced(ctx, searchParams)

		// 3. Если и fuzzy не помог, пробуем исправить опечатки вручную
		if err != nil || marketplaceResults == nil || len(marketplaceResults.Items) == 0 {
			// Пробуем найти похожие слова среди популярных брендов
			correctedQuery := s.correctSpelling(query)
			if correctedQuery != query {
				searchParams.Query = correctedQuery
				searchParams.Fuzziness = ""
				marketplaceResults, err = s.SearchListingsAdvanced(ctx, searchParams)
			}
		}
	}

	if err == nil && marketplaceResults != nil {
		// Ограничиваем до 3 товаров
		maxItems := 3
		if len(marketplaceResults.Items) < maxItems {
			maxItems = len(marketplaceResults.Items)
		}

		for i := 0; i < maxItems; i++ {
			item := marketplaceResults.Items[i]
			suggestions = append(suggestions, SuggestionItem{
				Type:      SuggestionTypeProduct,
				Value:     item.Title,
				Label:     item.Title,
				ProductID: item.ID,
				Metadata: map[string]interface{}{
					"price":       item.Price,
					"image":       getFirstImage(item.Images),
					"category":    item.Category.Name,
					"source_type": "marketplace",
				},
			})
		}
	}

	// 2. Поиск в товарах витрин
	if s.storage != nil && s.storage.StorefrontProductSearch() != nil {
		storefrontSearchRepo := s.storage.StorefrontProductSearch()

		// Проверяем, что репозиторий правильного типа
		if productSearchRepo, ok := storefrontSearchRepo.(storefrontOpenSearch.ProductSearchRepository); ok {
			// Создаем параметры поиска для товаров витрин
			storefrontParams := &storefrontOpenSearch.ProductSearchParams{
				Query:  query,
				Limit:  3,
				Offset: 0,
			}

			storefrontResults, err := productSearchRepo.SearchProducts(ctx, storefrontParams)
			if err == nil && storefrontResults != nil {
				for _, product := range storefrontResults.Products {
					if product != nil {
						// Получаем первое изображение
						var imageURL string
						if len(product.Images) > 0 {
							imageURL = product.Images[0].URL
						}

						suggestions = append(suggestions, SuggestionItem{
							Type:      SuggestionTypeProduct,
							Value:     product.Name,
							Label:     product.Name,
							ProductID: product.ProductID,
							Metadata: map[string]interface{}{
								"price":           product.Price,
								"image":           imageURL,
								"category":        product.Category.Name,
								"source_type":     "storefront",
								"storefront_id":   product.StorefrontID,
								"storefront":      product.Storefront.Name,
								"storefront_slug": product.Storefront.Slug,
							},
						})
					}
				}
			}
		}
	}

	// Ограничиваем общее количество до 5
	if len(suggestions) > 5 {
		suggestions = suggestions[:5]
	}

	return suggestions, nil
}

func (s *MarketplaceService) rankSuggestionsWithProductCategories(
	ctx context.Context,
	suggestions []SuggestionItem,
	query string,
	productCategories []string,
) []SuggestionItem {
	log.Debug().
		Str("query", query).
		Int("suggestions_count", len(suggestions)).
		Strs("product_categories", productCategories).
		Msg("rankSuggestionsWithProductCategories called")

	// ХАКФИКС: Если есть категории из товаров, но они не в suggestions, добавляем их
	if len(productCategories) > 0 {
		hasCategories := false
		for _, s := range suggestions {
			if s.Type == SuggestionTypeCategory {
				hasCategories = true
				break
			}
		}

		if !hasCategories {
			// Пытаемся добавить хотя бы первую категорию из товаров
			allCategories, _ := s.storage.GetCategories(ctx)

			for _, catName := range productCategories {
				for _, dbCat := range allCategories {
					// Проверяем по имени и переводам
					nameMatch := strings.EqualFold(dbCat.Name, catName)
					if !nameMatch && dbCat.Translations != nil {
						for _, translation := range dbCat.Translations {
							if strings.EqualFold(translation, catName) {
								nameMatch = true
								break
							}
						}
					}

					if nameMatch {
						iconStr := ""
						if dbCat.Icon != nil {
							iconStr = *dbCat.Icon
						}
						// Добавляем категорию в начало
						suggestions = append([]SuggestionItem{{
							Type:       SuggestionTypeCategory,
							Value:      dbCat.Slug,
							Label:      dbCat.Name,
							CategoryID: dbCat.ID,
							Icon:       iconStr,
							Count:      dbCat.ListingCount,
							Metadata: map[string]interface{}{
								"parent_id":     dbCat.ParentID,
								"from_products": true,
							},
						}}, suggestions...)
						log.Info().
							Str("category_name", dbCat.Name).
							Int("category_id", dbCat.ID).
							Msg("HACKFIX: Added missing product category")
						break
					}
				}
				// Добавляем только одну категорию
			}
		}
	}

	// Группируем по типам
	var queries, categoriesFromProducts, otherCategories, products []SuggestionItem
	productCategorySet := make(map[string]bool)
	for _, cat := range productCategories {
		productCategorySet[strings.ToLower(cat)] = true
	}

	// Подсчитываем типы для логирования
	typeCounts := map[string]int{}

	for _, s := range suggestions {
		typeCounts[string(s.Type)]++
		switch s.Type {
		case SuggestionTypeQuery:
			queries = append(queries, s)
		case SuggestionTypeCategory:
			// Разделяем категории из товаров и обычные категории
			if fromProducts, ok := s.Metadata["from_products"].(bool); ok && fromProducts {
				categoriesFromProducts = append(categoriesFromProducts, s)
				log.Debug().
					Str("category", s.Label).
					Int("id", s.CategoryID).
					Msg("Found category from products")
			} else if productCategorySet[strings.ToLower(s.Label)] {
				// Также проверяем по имени
				categoriesFromProducts = append(categoriesFromProducts, s)
				log.Debug().
					Str("category", s.Label).
					Int("id", s.CategoryID).
					Msg("Found category matching product category name")
			} else {
				otherCategories = append(otherCategories, s)
			}
		case SuggestionTypeProduct:
			products = append(products, s)
		}
	}

	log.Debug().
		Interface("type_counts", typeCounts).
		Int("categories_from_products", len(categoriesFromProducts)).
		Int("other_categories", len(otherCategories)).
		Int("products", len(products)).
		Int("queries", len(queries)).
		Msg("Grouped suggestions")

	// Сортируем каждую группу отдельно
	sortGroup := func(group []SuggestionItem) {
		sort.Slice(group, func(i, j int) bool {
			// По точности совпадения
			iExact := strings.HasPrefix(strings.ToLower(group[i].Label), query)
			jExact := strings.HasPrefix(strings.ToLower(group[j].Label), query)
			if iExact != jExact {
				return iExact
			}
			// По популярности (count)
			return group[i].Count > group[j].Count
		})
	}

	sortGroup(queries)
	sortGroup(categoriesFromProducts)
	sortGroup(otherCategories)
	sortGroup(products)

	// Собираем результат с правильным приоритетом
	var result []SuggestionItem

	// 1. СНАЧАЛА добавляем категории из найденных товаров (самый высокий приоритет)
	categoriesAdded := 0
	for i := 0; i < len(categoriesFromProducts) && i < 2; i++ {
		result = append(result, categoriesFromProducts[i])
		categoriesAdded++
		log.Debug().
			Str("category", categoriesFromProducts[i].Label).
			Int("position", len(result)-1).
			Msg("Added product category to result")
	}

	// 2. Добавляем до 2 запросов
	for i := 0; i < len(queries) && i < 2 && len(result) < 7; i++ {
		result = append(result, queries[i])
	}

	// 3. Добавляем до 3 товаров
	for i := 0; i < len(products) && i < 3 && len(result) < 8; i++ {
		result = append(result, products[i])
	}

	// 4. Добавляем другие категории, если есть место
	for i := 0; i < len(otherCategories) && len(result) < 10; i++ {
		result = append(result, otherCategories[i])
	}

	log.Debug().
		Int("final_result_count", len(result)).
		Int("categories_added", categoriesAdded).
		Msg("Ranking complete")

	return result
}

// rankSuggestions - удалена, используется rankSuggestionsWithProductCategories
/*
func (s *MarketplaceService) rankSuggestions(
	suggestions []SuggestionItem,
	query string,
) []SuggestionItem {
	// Группируем по типам
	var queries, categories, products []SuggestionItem

	for _, s := range suggestions {
		switch s.Type {
		case SuggestionTypeQuery:
			queries = append(queries, s)
		case SuggestionTypeCategory:
			categories = append(categories, s)
		case SuggestionTypeProduct:
			products = append(products, s)
		}
	}

	// Сортируем каждую группу отдельно
	sortGroup := func(group []SuggestionItem) {
		sort.Slice(group, func(i, j int) bool {
			// По точности совпадения
			iExact := strings.HasPrefix(strings.ToLower(group[i].Label), query)
			jExact := strings.HasPrefix(strings.ToLower(group[j].Label), query)
			if iExact != jExact {
				return iExact
			}
			// По популярности (count)
			return group[i].Count > group[j].Count
		})
	}

	sortGroup(queries)
	sortGroup(categories)
	sortGroup(products)

	// Собираем результат: запросы, товары, категории
	var result []SuggestionItem

	// Добавляем запросы (максимум 3)
	maxQueries := 3
	if len(queries) < maxQueries {
		maxQueries = len(queries)
	}
	for i := 0; i < maxQueries; i++ {
		result = append(result, queries[i])
	}

	// Добавляем товары (максимум 3-4)
	maxProducts := 4
	if len(products) < maxProducts {
		maxProducts = len(products)
	}
	for i := 0; i < maxProducts; i++ {
		result = append(result, products[i])
	}

	// Добавляем категории (минимум 1, если есть)
	maxCategories := 2
	if len(categories) > 0 && len(categories) < maxCategories {
		maxCategories = len(categories)
	}
	for i := 0; i < maxCategories && i < len(categories); i++ {
		result = append(result, categories[i])
	}

	return result
}
*/

// Вспомогательные функции
func containsString(slice []string, item string) bool {
	item = strings.TrimSpace(item)
	for _, s := range slice {
		if strings.TrimSpace(s) == item {
			return true
		}
	}
	return false
}

func getFirstImage(images []models.MarketplaceImage) string {
	if len(images) > 0 {
		return images[0].PublicURL
	}
	return ""
}

// SearchQuery представляет популярный поисковый запрос
type SearchQuery struct {
	ID              int    `json:"id" db:"id"`
	Query           string `json:"query" db:"query"`
	NormalizedQuery string `json:"normalized_query" db:"normalized_query"`
	SearchCount     int    `json:"search_count" db:"search_count"`
	LastSearched    string `json:"last_searched" db:"last_searched"`
	Language        string `json:"language" db:"language"`
	ResultsCount    int    `json:"results_count" db:"results_count"`
}

// SaveSearchQuery сохраняет или обновляет поисковый запрос
func (s *MarketplaceService) SaveSearchQuery(ctx context.Context, query string, resultsCount int, language string) error {
	normalized := strings.ToLower(strings.TrimSpace(query))
	if normalized == "" {
		return nil
	}

	return s.storage.SaveSearchQuery(ctx, query, normalized, resultsCount, language)
}

// correctSpelling пытается исправить опечатки в запросе
func (s *MarketplaceService) correctSpelling(query string) string {
	// Список популярных брендов и слов для автокоррекции
	commonTerms := []string{
		"volkswagen", "vw", "touran", "golf", "passat", "tiguan", "polo",
		"mercedes", "bmw", "audi", "opel", "ford", "toyota", "nissan",
		"peugeot", "renault", "citroen", "fiat", "skoda", "seat",
		"iphone", "samsung", "xiaomi", "huawei", "laptop", "playstation",
		"automobil", "auto", "telefon", "stan", "kuca", "garaza",
	}

	words := strings.Fields(strings.ToLower(query))
	corrected := make([]string, 0, len(words))

	for _, word := range words {
		bestMatch := word
		minDistance := 999

		// Ищем наиболее похожее слово
		for _, term := range commonTerms {
			distance := levenshteinDistance(word, term)
			// Принимаем исправление, если расстояние <= 2 и это лучше текущего
			if distance <= 2 && distance < minDistance {
				minDistance = distance
				bestMatch = term
			}
		}

		corrected = append(corrected, bestMatch)
	}

	return strings.Join(corrected, " ")
}

// levenshteinDistance вычисляет расстояние Левенштейна между двумя строками
func levenshteinDistance(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	// Создаем матрицу для динамического программирования
	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
	}

	// Инициализация первой строки и столбца
	for i := 0; i <= len(s1); i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= len(s2); j++ {
		matrix[0][j] = j
	}

	// Заполнение матрицы
	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			if s1[i-1] == s2[j-1] {
				matrix[i][j] = matrix[i-1][j-1]
			} else {
				min := matrix[i-1][j] // удаление
				if matrix[i][j-1] < min {
					min = matrix[i][j-1] // вставка
				}
				if matrix[i-1][j-1] < min {
					min = matrix[i-1][j-1] // замена
				}
				matrix[i][j] = min + 1
			}
		}
	}

	return matrix[len(s1)][len(s2)]
}
