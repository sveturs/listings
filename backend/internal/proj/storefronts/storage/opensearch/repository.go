package opensearch

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"backend/internal/domain/models"
	"backend/internal/logger"
	osClient "backend/internal/storage/opensearch"
)

// Helper function to safely get string value from pointer
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// StorefrontRepository реализует интерфейс для работы с витринами в OpenSearch
type StorefrontRepository struct {
	client    *osClient.OpenSearchClient
	indexName string
}

// NewStorefrontRepository создает новый репозиторий витрин для OpenSearch
func NewStorefrontRepository(client *osClient.OpenSearchClient, indexName string) *StorefrontRepository {
	return &StorefrontRepository{
		client:    client,
		indexName: indexName,
	}
}

// PrepareIndex подготавливает индекс для витрин (создает, если не существует)
func (r *StorefrontRepository) PrepareIndex(ctx context.Context) error {
	exists, err := r.client.IndexExists(ctx, r.indexName)
	if err != nil {
		return fmt.Errorf("ошибка проверки индекса витрин: %w", err)
	}

	logger.Info().Str("indexName", r.indexName).Bool("exists", exists).Msg("Проверка индекса витрин")

	if !exists {
		logger.Info().Str("indexName", r.indexName).Msg("Создание индекса витрин...")
		if err := r.client.CreateIndex(ctx, r.indexName, storefrontMapping); err != nil {
			return fmt.Errorf("ошибка создания индекса витрин: %w", err)
		}
		logger.Info().Str("indexName", r.indexName).Msg("Индекс витрин успешно создан")

		// Индексируем существующие витрины
		if err := r.ReindexAll(ctx); err != nil {
			logger.Error().Err(err).Msg("Ошибка переиндексации витрин")
			return err
		}
	}

	return nil
}

// Index индексирует одну витрину
func (r *StorefrontRepository) Index(ctx context.Context, storefront *models.Storefront) error {
	doc := r.storefrontToDoc(storefront)
	return r.client.IndexDocument(ctx, r.indexName, strconv.Itoa(storefront.ID), doc)
}

// BulkIndex индексирует несколько витрин
func (r *StorefrontRepository) BulkIndex(ctx context.Context, storefronts []*models.Storefront) error {
	if len(storefronts) == 0 {
		return nil
	}

	docs := make([]map[string]interface{}, 0, len(storefronts))
	for _, storefront := range storefronts {
		doc := r.storefrontToDoc(storefront)
		doc["id"] = storefront.ID
		docs = append(docs, doc)
	}

	return r.client.BulkIndex(ctx, r.indexName, docs)
}

// Delete удаляет витрину из индекса
func (r *StorefrontRepository) Delete(ctx context.Context, storefrontID int) error {
	return r.client.DeleteDocument(ctx, r.indexName, strconv.Itoa(storefrontID))
}

// Search выполняет поиск витрин
func (r *StorefrontRepository) Search(ctx context.Context, params *StorefrontSearchParams) (*StorefrontSearchResult, error) {
	query := r.buildSearchQuery(params)

	responseBytes, err := r.client.Search(ctx, r.indexName, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения поиска витрин: %w", err)
	}

	return r.parseSearchResponse(responseBytes)
}

// ReindexAll переиндексирует все активные витрины
func (r *StorefrontRepository) ReindexAll(ctx context.Context) error {
	// Для реиндексации нужно будет вызывать этот метод из Database уровня,
	// где есть доступ к PostgreSQL репозиторию витрин
	logger.Info().Msg("ReindexAll для витрин должен вызываться из Database уровня")
	return nil
}

// storefrontToDoc преобразует модель витрины в документ для индексации
func (r *StorefrontRepository) storefrontToDoc(storefront *models.Storefront) map[string]interface{} {
	doc := map[string]interface{}{
		"id":                storefront.ID,
		"user_id":           storefront.UserID,
		"slug":              storefront.Slug,
		"name":              storefront.Name,
		"name_lowercase":    strings.ToLower(storefront.Name),
		"description":       storefront.Description,
		"phone":             storefront.Phone,
		"email":             storefront.Email,
		"website":           storefront.Website,
		"address":           getStringValue(storefront.Address),
		"city":              getStringValue(storefront.City),
		"city_lowercase":    strings.ToLower(getStringValue(storefront.City)),
		"postal_code":       getStringValue(storefront.PostalCode),
		"country":           getStringValue(storefront.Country),
		"rating":            storefront.Rating,
		"reviews_count":     storefront.ReviewsCount,
		"products_count":    storefront.ProductsCount,
		"sales_count":       storefront.SalesCount,
		"views_count":       storefront.ViewsCount,
		"subscription_plan": storefront.SubscriptionPlan,
		"is_active":         storefront.IsActive,
		"is_verified":       storefront.IsVerified,
		"created_at":        storefront.CreatedAt.Format(time.RFC3339),
		"updated_at":        storefront.UpdatedAt.Format(time.RFC3339),
	}

	// Добавляем геолокацию если есть
	if storefront.Latitude != nil && storefront.Longitude != nil && *storefront.Latitude != 0 && *storefront.Longitude != 0 {
		doc["location"] = map[string]interface{}{
			"lat": *storefront.Latitude,
			"lon": *storefront.Longitude,
		}
	}

	// Добавляем дополнительные поля для поиска
	searchKeywords := []string{
		storefront.Name,
		strings.ToLower(storefront.Name),
		getStringValue(storefront.City),
		strings.ToLower(getStringValue(storefront.City)),
	}

	// Добавляем слова из описания
	if storefront.Description != nil && *storefront.Description != "" {
		words := strings.Fields(*storefront.Description)
		for _, word := range words {
			if len(word) > 3 { // Только слова длиннее 3 символов
				searchKeywords = append(searchKeywords, strings.ToLower(word))
			}
		}
	}

	doc["search_keywords"] = searchKeywords

	// Базовые поля для поиска без зависимостей
	// Расширенная информация (часы работы, методы оплаты, доставка)
	// будет добавляться на уровне Database при вызове полной индексации

	// Добавляем тематику витрины (из theme)
	if storefront.Theme != nil {
		if primaryColor, ok := storefront.Theme["primaryColor"].(string); ok {
			doc["theme_primary_color"] = primaryColor
		}
		if style, ok := storefront.Theme["style"].(string); ok {
			doc["theme_style"] = style
		}
	}

	// Добавляем SEO данные для улучшения поиска
	if storefront.SEOMeta != nil {
		if keywords, ok := storefront.SEOMeta["keywords"].(string); ok {
			keywordsList := strings.Split(keywords, ",")
			for _, kw := range keywordsList {
				searchKeywords = append(searchKeywords, strings.TrimSpace(strings.ToLower(kw)))
			}
		}
	}

	// Обновляем ключевые слова с учетом SEO
	doc["search_keywords"] = deduplicate(searchKeywords)

	return doc
}

// buildSearchQuery строит запрос для поиска витрин
func (r *StorefrontRepository) buildSearchQuery(params *StorefrontSearchParams) map[string]interface{} {
	query := map[string]interface{}{
		"size": params.Limit,
		"from": params.Offset,
	}

	// Построение условий поиска
	must := []map[string]interface{}{}
	filter := []map[string]interface{}{}

	// Текстовый поиск
	if params.Query != "" {
		must = append(must, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query": params.Query,
				"fields": []string{
					"name^3",
					"name_lowercase^2",
					"description",
					"search_keywords",
					"city",
					"address",
				},
				"type": "best_fields",
			},
		})
	}

	// Фильтр по городу
	if params.City != "" {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"city_lowercase": strings.ToLower(params.City),
			},
		})
	}

	// Фильтр по активности
	if params.IsActive != nil {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"is_active": *params.IsActive,
			},
		})
	}

	// Фильтр по верификации
	if params.IsVerified != nil {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"is_verified": *params.IsVerified,
			},
		})
	}

	// Фильтр по рейтингу
	if params.MinRating > 0 {
		filter = append(filter, map[string]interface{}{
			"range": map[string]interface{}{
				"rating": map[string]interface{}{
					"gte": params.MinRating,
				},
			},
		})
	}

	// Геолокационный поиск
	if params.Latitude != 0 && params.Longitude != 0 && params.RadiusKm > 0 {
		filter = append(filter, map[string]interface{}{
			"geo_distance": map[string]interface{}{
				"distance": fmt.Sprintf("%dkm", params.RadiusKm),
				"location": map[string]interface{}{
					"lat": params.Latitude,
					"lon": params.Longitude,
				},
			},
		})
	}

	// Фильтр по методам оплаты
	if len(params.PaymentMethods) > 0 {
		for _, method := range params.PaymentMethods {
			filter = append(filter, map[string]interface{}{
				"term": map[string]interface{}{
					"payment_methods": method,
				},
			})
		}
	}

	// Фильтр по доставке
	if params.HasDelivery != nil {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"has_delivery": *params.HasDelivery,
			},
		})
	}

	// Фильтр по самовывозу
	if params.HasSelfPickup != nil {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"has_self_pickup": *params.HasSelfPickup,
			},
		})
	}

	// Фильтр по открытости сейчас
	if params.IsOpenNow != nil && *params.IsOpenNow {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"is_open_now": true,
			},
		})
	}

	// Построение bool запроса
	boolQuery := map[string]interface{}{}
	if len(must) > 0 {
		boolQuery["must"] = must
	}
	if len(filter) > 0 {
		boolQuery["filter"] = filter
	}

	// Если нет условий поиска, используем match_all
	if len(must) == 0 && len(filter) == 0 {
		query["query"] = map[string]interface{}{
			"match_all": map[string]interface{}{},
		}
	} else {
		query["query"] = map[string]interface{}{
			"bool": boolQuery,
		}
	}

	// Сортировка
	sort := []map[string]interface{}{}
	switch params.SortBy {
	case "rating":
		sort = append(sort, map[string]interface{}{
			"rating": map[string]interface{}{
				"order": params.SortOrder,
			},
		})
	case "products_count":
		sort = append(sort, map[string]interface{}{
			"products_count": map[string]interface{}{
				"order": params.SortOrder,
			},
		})
	case "distance":
		if params.Latitude != 0 && params.Longitude != 0 {
			sort = append(sort, map[string]interface{}{
				"_geo_distance": map[string]interface{}{
					"location": map[string]interface{}{
						"lat": params.Latitude,
						"lon": params.Longitude,
					},
					"order": params.SortOrder,
					"unit":  "km",
				},
			})
		}
	case "created_at":
		sort = append(sort, map[string]interface{}{
			"created_at": map[string]interface{}{
				"order": params.SortOrder,
			},
		})
	default:
		// По умолчанию сортировка по рейтингу и количеству товаров
		sort = append(sort,
			map[string]interface{}{
				"rating": map[string]interface{}{
					"order": "desc",
				},
			},
			map[string]interface{}{
				"products_count": map[string]interface{}{
					"order": "desc",
				},
			},
		)
	}

	if len(sort) > 0 {
		query["sort"] = sort
	}

	// Добавляем подсветку результатов
	query["highlight"] = map[string]interface{}{
		"fields": map[string]interface{}{
			"name": map[string]interface{}{},
			"description": map[string]interface{}{
				"fragment_size": 150,
			},
		},
	}

	return query
}

// parseSearchResponse парсит ответ от OpenSearch
func (r *StorefrontRepository) parseSearchResponse(responseBytes []byte) (*StorefrontSearchResult, error) {
	var response map[string]interface{}
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return nil, fmt.Errorf("ошибка разбора ответа: %w", err)
	}

	result := &StorefrontSearchResult{
		Storefronts: []*StorefrontSearchItem{},
		Total:       0,
	}

	// Извлекаем общее количество результатов
	if hits, ok := response["hits"].(map[string]interface{}); ok {
		if total, ok := hits["total"].(map[string]interface{}); ok {
			if value, ok := total["value"].(float64); ok {
				result.Total = int(value)
			}
		}

		// Извлекаем результаты
		if hitsArray, ok := hits["hits"].([]interface{}); ok {
			for _, hit := range hitsArray {
				if hitObj, ok := hit.(map[string]interface{}); ok {
					item := &StorefrontSearchItem{}

					// ID документа
					if id, ok := hitObj["_id"].(string); ok {
						if idInt, err := strconv.Atoi(id); err == nil {
							item.ID = idInt
						}
					}

					// Score
					if score, ok := hitObj["_score"].(float64); ok {
						item.Score = score
					}

					// Данные витрины
					if source, ok := hitObj["_source"].(map[string]interface{}); ok {
						r.parseStorefrontSource(source, item)
					}

					// Подсветка
					if highlight, ok := hitObj["highlight"].(map[string]interface{}); ok {
						item.Highlights = make(map[string][]string)
						for field, values := range highlight {
							if valArray, ok := values.([]interface{}); ok {
								highlights := []string{}
								for _, v := range valArray {
									if str, ok := v.(string); ok {
										highlights = append(highlights, str)
									}
								}
								item.Highlights[field] = highlights
							}
						}
					}

					// Расстояние (если есть)
					if sort, ok := hitObj["sort"].([]interface{}); ok && len(sort) > 0 {
						if distance, ok := sort[0].(float64); ok {
							item.Distance = &distance
						}
					}

					result.Storefronts = append(result.Storefronts, item)
				}
			}
		}
	}

	return result, nil
}

// parseStorefrontSource парсит данные витрины из _source
func (r *StorefrontRepository) parseStorefrontSource(source map[string]interface{}, item *StorefrontSearchItem) {
	if v, ok := source["user_id"].(float64); ok {
		item.UserID = int(v)
	}
	if v, ok := source["slug"].(string); ok {
		item.Slug = v
	}
	if v, ok := source["name"].(string); ok {
		item.Name = v
	}
	if v, ok := source["description"].(string); ok {
		item.Description = v
	}
	if v, ok := source["phone"].(string); ok {
		item.Phone = v
	}
	if v, ok := source["email"].(string); ok {
		item.Email = v
	}
	if v, ok := source["address"].(string); ok {
		item.Address = v
	}
	if v, ok := source["city"].(string); ok {
		item.City = v
	}
	if v, ok := source["country"].(string); ok {
		item.Country = v
	}
	if v, ok := source["rating"].(float64); ok {
		item.Rating = v
	}
	if v, ok := source["reviews_count"].(float64); ok {
		item.ReviewsCount = int(v)
	}
	if v, ok := source["products_count"].(float64); ok {
		item.ProductsCount = int(v)
	}
	if v, ok := source["is_verified"].(bool); ok {
		item.IsVerified = v
	}
	if v, ok := source["is_open_now"].(bool); ok {
		item.IsOpenNow = v
	}

	// Геолокация
	if location, ok := source["location"].(map[string]interface{}); ok {
		if lat, ok := location["lat"].(float64); ok {
			item.Latitude = lat
		}
		if lon, ok := location["lon"].(float64); ok {
			item.Longitude = lon
		}
	}

	// Методы оплаты
	if methods, ok := source["payment_methods"].([]interface{}); ok {
		item.PaymentMethods = []string{}
		for _, m := range methods {
			if method, ok := m.(string); ok {
				item.PaymentMethods = append(item.PaymentMethods, method)
			}
		}
	}

	// Доставка
	if v, ok := source["has_delivery"].(bool); ok {
		item.HasDelivery = v
	}
	if v, ok := source["has_self_pickup"].(bool); ok {
		item.HasSelfPickup = v
	}
}

// StorefrontSearchParams параметры поиска витрин
type StorefrontSearchParams struct {
	Query          string   // Текстовый поиск
	City           string   // Фильтр по городу
	Latitude       float64  // Широта для геопоиска
	Longitude      float64  // Долгота для геопоиска
	RadiusKm       int      // Радиус поиска в километрах
	MinRating      float64  // Минимальный рейтинг
	IsActive       *bool    // Только активные
	IsVerified     *bool    // Только верифицированные
	IsOpenNow      *bool    // Только открытые сейчас
	PaymentMethods []string // Фильтр по методам оплаты
	HasDelivery    *bool    // Есть доставка
	HasSelfPickup  *bool    // Есть самовывоз
	SortBy         string   // Поле сортировки
	SortOrder      string   // Порядок сортировки (asc/desc)
	Limit          int      // Количество результатов
	Offset         int      // Смещение
}

// StorefrontSearchResult результат поиска витрин
type StorefrontSearchResult struct {
	Storefronts []*StorefrontSearchItem
	Total       int
}

// StorefrontSearchItem элемент результата поиска
type StorefrontSearchItem struct {
	ID             int
	UserID         int
	Slug           string
	Name           string
	Description    string
	Phone          string
	Email          string
	Address        string
	City           string
	Country        string
	Latitude       float64
	Longitude      float64
	Rating         float64
	ReviewsCount   int
	ProductsCount  int
	IsVerified     bool
	IsOpenNow      bool
	PaymentMethods []string
	HasDelivery    bool
	HasSelfPickup  bool
	Score          float64             // Релевантность
	Distance       *float64            // Расстояние в км (если есть)
	Highlights     map[string][]string // Подсвеченные фрагменты
}

// Helper functions

func deduplicate(slice []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

// storefrontMapping маппинг для индекса витрин
const storefrontMapping = `{
  "mappings": {
    "properties": {
      "id": {"type": "integer"},
      "user_id": {"type": "integer"},
      "slug": {"type": "keyword"},
      "name": {
        "type": "text",
        "fields": {
          "keyword": {"type": "keyword"}
        }
      },
      "name_lowercase": {"type": "text"},
      "description": {"type": "text"},
      "phone": {"type": "keyword"},
      "email": {"type": "keyword"},
      "website": {"type": "keyword"},
      "address": {"type": "text"},
      "city": {
        "type": "text",
        "fields": {
          "keyword": {"type": "keyword"}
        }
      },
      "city_lowercase": {"type": "keyword"},
      "postal_code": {"type": "keyword"},
      "country": {"type": "keyword"},
      "location": {"type": "geo_point"},
      "rating": {"type": "float"},
      "reviews_count": {"type": "integer"},
      "products_count": {"type": "integer"},
      "sales_count": {"type": "integer"},
      "views_count": {"type": "integer"},
      "subscription_plan": {"type": "keyword"},
      "is_active": {"type": "boolean"},
      "is_verified": {"type": "boolean"},
      "is_open_now": {"type": "boolean"},
      "payment_methods": {"type": "keyword"},
      "accepts_cash": {"type": "boolean"},
      "accepts_cod": {"type": "boolean"},
      "accepts_cards": {"type": "boolean"},
      "delivery_providers": {"type": "keyword"},
      "has_delivery": {"type": "boolean"},
      "has_self_pickup": {"type": "boolean"},
      "theme_primary_color": {"type": "keyword"},
      "theme_style": {"type": "keyword"},
      "search_keywords": {"type": "text"},
      "created_at": {"type": "date"},
      "updated_at": {"type": "date"}
    }
  }
}`
