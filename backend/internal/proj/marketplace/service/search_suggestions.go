package service

import (
	"context"
	"sort"
	"strings"

	"backend/internal/domain/models"
	"backend/internal/domain/search"
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

	// 1. Популярные поисковые запросы
	if containsString(typesList, "queries") {
		querySuggestions, err := s.getQuerySuggestions(ctx, query)
		if err == nil {
			suggestions = append(suggestions, querySuggestions...)
		}
	}

	// 2. Категории
	if containsString(typesList, "categories") {
		categorySuggestions, err := s.getCategorySuggestions(ctx, query)
		if err == nil {
			suggestions = append(suggestions, categorySuggestions...)
		}
	}

	// 3. Товары/Объявления
	if containsString(typesList, "products") {
		productSuggestions, err := s.getProductSuggestions(ctx, query)
		if err == nil {
			suggestions = append(suggestions, productSuggestions...)
		}
	}

	// Сортируем по релевантности и ограничиваем
	suggestions = s.rankSuggestions(suggestions, query)
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

func (s *MarketplaceService) getCategorySuggestions(
	ctx context.Context,
	query string,
) ([]SuggestionItem, error) {
	// Ищем категории по имени
	categories, err := s.storage.SearchCategories(ctx, query, 5)
	if err != nil {
		return nil, err
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

	return suggestions, nil
}

func (s *MarketplaceService) getProductSuggestions(
	ctx context.Context,
	query string,
) ([]SuggestionItem, error) {
	// Быстрый поиск объявлений
	searchParams := &search.ServiceParams{
		Query: query,
		Size:  5,
		Page:  1,
	}

	results, err := s.SearchListingsAdvanced(ctx, searchParams)
	if err != nil {
		return nil, err
	}

	suggestions := make([]SuggestionItem, 0, len(results.Items))
	for _, item := range results.Items {
		suggestions = append(suggestions, SuggestionItem{
			Type:      SuggestionTypeProduct,
			Value:     item.Title,
			Label:     item.Title,
			ProductID: item.ID,
			Metadata: map[string]interface{}{
				"price":    item.Price,
				"image":    getFirstImage(item.Images),
				"category": item.Category.Name,
			},
		})
	}

	return suggestions, nil
}

func (s *MarketplaceService) rankSuggestions(
	suggestions []SuggestionItem,
	query string,
) []SuggestionItem {
	// Ранжируем по релевантности
	sort.Slice(suggestions, func(i, j int) bool {
		// Приоритет типов: query > category > product
		if suggestions[i].Type != suggestions[j].Type {
			typeOrder := map[SuggestionType]int{
				SuggestionTypeQuery:    1,
				SuggestionTypeCategory: 2,
				SuggestionTypeProduct:  3,
			}
			return typeOrder[suggestions[i].Type] < typeOrder[suggestions[j].Type]
		}

		// По точности совпадения
		iExact := strings.HasPrefix(strings.ToLower(suggestions[i].Label), query)
		jExact := strings.HasPrefix(strings.ToLower(suggestions[j].Label), query)
		if iExact != jExact {
			return iExact
		}

		// По популярности (count)
		return suggestions[i].Count > suggestions[j].Count
	})

	return suggestions
}

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
