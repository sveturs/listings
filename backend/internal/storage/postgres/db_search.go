// backend/internal/storage/postgres/db_search.go
package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"backend/internal/config"
	"backend/internal/domain/models"
	"backend/internal/domain/search"
)

// GetSearchWeights возвращает веса для поиска
func (db *Database) GetSearchWeights() *config.SearchWeights {
	return db.searchWeights
}

// GetOpenSearchClient возвращает клиент OpenSearch для прямого выполнения запросов
func (db *Database) GetOpenSearchClient() (interface {
	Execute(ctx context.Context, method, path string, body []byte) ([]byte, error)
}, error,
) {
	if db.osClient == nil {
		return nil, fmt.Errorf("OpenSearch клиент не настроен")
	}
	return db.osClient, nil
}

// PrepareIndex подготавливает индекс OpenSearch
func (db *Database) PrepareIndex(ctx context.Context) error {
	if true { // OpenSearch disabled after removing c2c
		// Если репозиторий OpenSearch не инициализирован, просто возвращаем nil
		// Поиск будет работать без OpenSearch
		return nil
	}

	// TODO: OpenSearch temporarily disabled during refactoring
	return nil
}

// SearchListingsOpenSearch выполняет поиск через OpenSearch
func (db *Database) SearchListingsOpenSearch(ctx context.Context, params *search.SearchParams) (*search.SearchResult, error) {
	// Если OpenSearch клиент не настроен, возвращаем пустой результат
	if db.osClient == nil {
		return &search.SearchResult{
			Listings: []*models.MarketplaceListing{},
			Total:    0,
			Took:     0,
		}, nil
	}

	// Строим запрос для OpenSearch
	query := db.buildSearchQuery(params)

	// Выполняем поиск
	response, err := db.executeOpenSearchQuery(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to execute OpenSearch query: %w", err)
	}

	// Парсим результаты
	return db.parseOpenSearchResults(ctx, response, params)
}

// buildSearchQuery строит запрос для OpenSearch
func (db *Database) buildSearchQuery(params *search.SearchParams) map[string]interface{} {
	query := map[string]interface{}{
		"track_total_hits": true,
	}

	// Базовый bool query
	boolQuery := map[string]interface{}{
		"must":   []interface{}{},
		"filter": []interface{}{},
	}

	// Добавляем текстовый поиск если есть запрос
	if params.Query != "" {
		boolQuery["must"] = append(boolQuery["must"].([]interface{}), map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query": params.Query,
				"fields": []string{
					"title^3",
					"description^2",
					"category_name",
				},
				"type":      "best_fields",
				"operator":  "or",
				"fuzziness": "AUTO",
			},
		})
	} else {
		// Если нет текстового запроса, используем match_all
		boolQuery["must"] = append(boolQuery["must"].([]interface{}), map[string]interface{}{
			"match_all": map[string]interface{}{},
		})
	}

	// Фильтр по категории
	if params.CategoryID != nil && *params.CategoryID > 0 {
		boolQuery["filter"] = append(boolQuery["filter"].([]interface{}), map[string]interface{}{
			"term": map[string]interface{}{
				"category_id": *params.CategoryID,
			},
		})
	}

	// Фильтр по множественным категориям
	if len(params.CategoryIDs) > 0 {
		boolQuery["filter"] = append(boolQuery["filter"].([]interface{}), map[string]interface{}{
			"terms": map[string]interface{}{
				"category_id": params.CategoryIDs,
			},
		})
	}

	// Фильтр по цене
	if params.PriceMin != nil || params.PriceMax != nil {
		rangeQuery := map[string]interface{}{}
		if params.PriceMin != nil && *params.PriceMin > 0 {
			rangeQuery["gte"] = *params.PriceMin
		}
		if params.PriceMax != nil && *params.PriceMax > 0 {
			rangeQuery["lte"] = *params.PriceMax
		}
		if len(rangeQuery) > 0 {
			boolQuery["filter"] = append(boolQuery["filter"].([]interface{}), map[string]interface{}{
				"range": map[string]interface{}{
					"price": rangeQuery,
				},
			})
		}
	}

	// Фильтр по городу
	if params.City != "" {
		boolQuery["filter"] = append(boolQuery["filter"].([]interface{}), map[string]interface{}{
			"term": map[string]interface{}{
				"city.keyword": params.City,
			},
		})
	}

	// Фильтр по статусу (по умолчанию только активные)
	status := "active"
	if params.Status != "" {
		status = params.Status
	}
	boolQuery["filter"] = append(boolQuery["filter"].([]interface{}), map[string]interface{}{
		"term": map[string]interface{}{
			"status": status,
		},
	})

	// Фильтр по типу документа
	if params.DocumentType != "" {
		boolQuery["filter"] = append(boolQuery["filter"].([]interface{}), map[string]interface{}{
			"term": map[string]interface{}{
				"document_type": params.DocumentType,
			},
		})
	}

	query["query"] = map[string]interface{}{
		"bool": boolQuery,
	}

	return query
}

// executeOpenSearchQuery выполняет запрос к OpenSearch
func (db *Database) executeOpenSearchQuery(ctx context.Context, query map[string]interface{}, params *search.SearchParams) (map[string]interface{}, error) {
	// Сериализуем запрос
	queryBytes, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	// Определяем from и size для пагинации
	from := 0
	if params.Page > 0 {
		from = (params.Page - 1) * params.Size
	}
	size := params.Size
	if size <= 0 {
		size = 20
	}

	// Выполняем запрос через Execute метод клиента
	path := fmt.Sprintf("/%s/_search?from=%d&size=%d", db.marketplaceIndex, from, size)
	responseBytes, err := db.osClient.Execute(ctx, "POST", path, queryBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search: %w", err)
	}

	// Парсим ответ
	var response map[string]interface{}
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response, nil
}

// parseOpenSearchResults парсит результаты из OpenSearch
func (db *Database) parseOpenSearchResults(ctx context.Context, response map[string]interface{}, params *search.SearchParams) (*search.SearchResult, error) {
	result := &search.SearchResult{
		Listings: []*models.MarketplaceListing{},
		Total:    0,
		Took:     0,
	}

	// Получаем время выполнения
	if took, ok := response["took"].(float64); ok {
		result.Took = int64(took)
	}

	// Получаем hits
	hits, ok := response["hits"].(map[string]interface{})
	if !ok {
		return result, nil
	}

	// Получаем total
	if total, ok := hits["total"].(map[string]interface{}); ok {
		if value, ok := total["value"].(float64); ok {
			result.Total = int(value)
		}
	}

	// Получаем документы
	hitsArray, ok := hits["hits"].([]interface{})
	if !ok {
		return result, nil
	}

	// Парсим каждый документ
	for _, hit := range hitsArray {
		hitMap, ok := hit.(map[string]interface{})
		if !ok {
			continue
		}

		source, ok := hitMap["_source"].(map[string]interface{})
		if !ok {
			continue
		}

		// Конвертируем source в MarketplaceListing
		sourceBytes, err := json.Marshal(source)
		if err != nil {
			continue
		}

		var listing models.MarketplaceListing
		if err := json.Unmarshal(sourceBytes, &listing); err != nil {
			continue
		}

		result.Listings = append(result.Listings, &listing)
	}

	return result, nil
}
