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

// ProductRepository реализует интерфейс для работы с товарами витрин в OpenSearch
type ProductRepository struct {
	client    *osClient.OpenSearchClient
	indexName string
}

// NewProductRepository создает новый репозиторий товаров витрин для OpenSearch
func NewProductRepository(client *osClient.OpenSearchClient, indexName string) *ProductRepository {
	return &ProductRepository{
		client:    client,
		indexName: indexName,
	}
}

// PrepareIndex подготавливает индекс для товаров витрин (создает, если не существует)
func (r *ProductRepository) PrepareIndex(ctx context.Context) error {
	exists, err := r.client.IndexExists(r.indexName)
	if err != nil {
		return fmt.Errorf("ошибка проверки индекса товаров витрин: %w", err)
	}

	logger.Info().Str("indexName", r.indexName).Bool("exists", exists).Msg("Проверка индекса товаров витрин")

	if !exists {
		logger.Info().Str("indexName", r.indexName).Msg("Создание индекса товаров витрин...")
		if err := r.client.CreateIndex(r.indexName, storefrontProductMapping); err != nil {
			return fmt.Errorf("ошибка создания индекса товаров витрин: %w", err)
		}
		logger.Info().Str("indexName", r.indexName).Msg("Индекс товаров витрин успешно создан")
	}

	return nil
}

// IndexProduct индексирует один товар витрины
func (r *ProductRepository) IndexProduct(ctx context.Context, product *models.StorefrontProduct) error {
	doc := r.productToDoc(product)
	docID := fmt.Sprintf("sp_%d", product.ID)
	
	logger.Info().Msgf("Индексация товара витрины: ID=%d, Name=%s, StorefrontID=%d", 
		product.ID, product.Name, product.StorefrontID)
	
	return r.client.IndexDocument(r.indexName, docID, doc)
}

// BulkIndexProducts индексирует несколько товаров витрин
func (r *ProductRepository) BulkIndexProducts(ctx context.Context, products []*models.StorefrontProduct) error {
	if len(products) == 0 {
		return nil
	}

	docs := make([]map[string]interface{}, 0, len(products))
	for _, product := range products {
		doc := r.productToDoc(product)
		doc["id"] = fmt.Sprintf("sp_%d", product.ID)
		docs = append(docs, doc)
	}

	logger.Info().Msgf("Массовая индексация %d товаров витрин", len(products))
	return r.client.BulkIndex(r.indexName, docs)
}

// DeleteProduct удаляет товар из индекса
func (r *ProductRepository) DeleteProduct(ctx context.Context, productID int) error {
	docID := fmt.Sprintf("sp_%d", productID)
	return r.client.DeleteDocument(r.indexName, docID)
}

// UpdateProduct обновляет товар в индексе
func (r *ProductRepository) UpdateProduct(ctx context.Context, product *models.StorefrontProduct) error {
	return r.IndexProduct(ctx, product)
}

// productToDoc преобразует модель товара витрины в документ для индексации
func (r *ProductRepository) productToDoc(product *models.StorefrontProduct) map[string]interface{} {
	doc := map[string]interface{}{
		"product_id":       product.ID,
		"product_type":     "storefront",
		"storefront_id":    product.StorefrontID,
		"category_id":      product.CategoryID,
		"name":             product.Name,
		"name_lowercase":   strings.ToLower(product.Name),
		"description":      product.Description,
		"price":            product.Price,
		"currency":         product.Currency,
		"sku":              product.SKU,
		"barcode":          product.Barcode,
		"stock_status":     product.StockStatus,
		"is_active":        product.IsActive,
		"created_at":       product.CreatedAt.Format(time.RFC3339),
		"updated_at":       product.UpdatedAt.Format(time.RFC3339),
	}

	// Добавляем информацию о наличии товара
	doc["stock_quantity"] = product.StockQuantity
	doc["inventory"] = map[string]interface{}{
		"quantity":  product.StockQuantity,
		"in_stock":  product.StockQuantity > 0,
		"low_stock": product.StockQuantity > 0 && product.StockQuantity <= 5, // TODO: сделать настраиваемым
		"status":    product.StockStatus,
	}

	// Добавляем атрибуты товара
	if product.Attributes != nil && len(product.Attributes) > 0 {
		// Извлекаем важные поля из атрибутов
		if brand, ok := product.Attributes["brand"].(string); ok {
			doc["brand"] = brand
			doc["brand_lowercase"] = strings.ToLower(brand)
		}
		if model, ok := product.Attributes["model"].(string); ok {
			doc["model"] = model
			doc["model_lowercase"] = strings.ToLower(model)
		}
		if color, ok := product.Attributes["color"].(string); ok {
			doc["color"] = color
		}
		if size, ok := product.Attributes["size"].(string); ok {
			doc["size"] = size
		}
		if material, ok := product.Attributes["material"].(string); ok {
			doc["material"] = material
		}
		
		// Сохраняем все атрибуты для расширенного поиска
		doc["attributes"] = product.Attributes
	}

	// Добавляем поисковые ключевые слова
	searchKeywords := []string{product.Name, strings.ToLower(product.Name)}
	if sku := product.SKU; sku != nil && *sku != "" {
		searchKeywords = append(searchKeywords, *sku)
	}
	if barcode := product.Barcode; barcode != nil && *barcode != "" {
		searchKeywords = append(searchKeywords, *barcode)
	}
	
	// Добавляем значения атрибутов в поисковые ключевые слова
	if product.Attributes != nil {
		for _, value := range product.Attributes {
			if strVal, ok := value.(string); ok && strVal != "" {
				searchKeywords = append(searchKeywords, strVal, strings.ToLower(strVal))
			}
		}
	}
	
	doc["search_keywords"] = deduplicate(searchKeywords)

	// Добавляем изображения
	if product.Images != nil && len(product.Images) > 0 {
		imagesArray := make([]map[string]interface{}, 0, len(product.Images))
		for _, img := range product.Images {
			imgDoc := map[string]interface{}{
				"id":            img.ID,
				"url":           img.ImageURL,
				"thumbnail_url": img.ThumbnailURL,
				"is_default":    img.IsDefault,
				"display_order": img.DisplayOrder,
			}
			imagesArray = append(imagesArray, imgDoc)
		}
		doc["images"] = imagesArray
		doc["has_images"] = true
		doc["image_count"] = len(product.Images)
	} else {
		doc["has_images"] = false
		doc["image_count"] = 0
	}

	// Добавляем варианты товара
	if product.Variants != nil && len(product.Variants) > 0 {
		variantsArray := make([]map[string]interface{}, 0, len(product.Variants))
		minPrice := product.Price
		maxPrice := product.Price
		
		for _, variant := range product.Variants {
			varDoc := map[string]interface{}{
				"id":         variant.ID,
				"name":       variant.Name,
				"sku":        variant.SKU,
				"price":      variant.Price,
				"attributes": variant.Attributes,
			}
			
			// Отслеживаем диапазон цен
			if variant.Price < minPrice {
				minPrice = variant.Price
			}
			if variant.Price > maxPrice {
				maxPrice = variant.Price
			}
			
			variantsArray = append(variantsArray, varDoc)
		}
		
		doc["variants"] = variantsArray
		doc["has_variants"] = true
		doc["variant_count"] = len(product.Variants)
		doc["price_min"] = minPrice
		doc["price_max"] = maxPrice
	} else {
		doc["has_variants"] = false
		doc["variant_count"] = 0
		doc["price_min"] = product.Price
		doc["price_max"] = product.Price
	}

	// TODO: Добавить загрузку информации о витрине по StorefrontID
	// На данный момент просто добавляем ID витрины
	doc["storefront_id"] = product.StorefrontID

	// Добавляем категорию если есть
	if product.Category != nil {
		doc["category"] = map[string]interface{}{
			"id":   product.Category.ID,
			"name": product.Category.Name,
			"slug": product.Category.Slug,
		}
		
	}

	// Рассчитываем дополнительные поля для ранжирования
	doc["popularity_score"] = calculatePopularityScore(product)
	doc["quality_score"] = calculateQualityScore(product)

	return doc
}

// SearchProducts выполняет поиск товаров витрин
func (r *ProductRepository) SearchProducts(ctx context.Context, params *ProductSearchParams) (*ProductSearchResult, error) {
	query := r.buildSearchQuery(params)
	
	responseBytes, err := r.client.Search(r.indexName, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения поиска товаров витрин: %w", err)
	}
	
	return r.parseSearchResponse(responseBytes, params)
}

// ReindexAll переиндексирует все активные товары витрин
func (r *ProductRepository) ReindexAll(ctx context.Context) error {
	// Этот метод должен вызываться из сервиса, имеющего доступ к PostgreSQL
	logger.Info().Msg("ReindexAll для товаров витрин должен вызываться из сервиса с доступом к БД")
	return nil
}

// buildSearchQuery строит запрос для поиска товаров
func (r *ProductRepository) buildSearchQuery(params *ProductSearchParams) map[string]interface{} {
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
					"name.autocomplete^2",
					"name_lowercase^2",
					"description",
					"brand^2",
					"model^2",
					"search_keywords",
					"sku",
					"barcode",
				},
				"type":       "best_fields",
				"fuzziness":  "AUTO",
				"prefix_length": 2,
			},
		})
	}
	
	// Фильтр по витрине
	if params.StorefrontID > 0 {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"storefront_id": params.StorefrontID,
			},
		})
	}
	
	// Фильтр по категории
	if params.CategoryID > 0 {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"category_id": params.CategoryID,
			},
		})
	}
	
	// Фильтр по пути категории (для поиска в подкатегориях)
	if params.CategoryPath != "" {
		filter = append(filter, map[string]interface{}{
			"prefix": map[string]interface{}{
				"category_path": params.CategoryPath,
			},
		})
	}
	
	// Фильтр по бренду
	if params.Brand != "" {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"brand.keyword": params.Brand,
			},
		})
	}
	
	// Фильтр по цене
	if params.PriceMin > 0 || params.PriceMax > 0 {
		priceRange := map[string]interface{}{}
		if params.PriceMin > 0 {
			priceRange["gte"] = params.PriceMin
		}
		if params.PriceMax > 0 {
			priceRange["lte"] = params.PriceMax
		}
		filter = append(filter, map[string]interface{}{
			"range": map[string]interface{}{
				"price": priceRange,
			},
		})
	}
	
	// Фильтр по наличию
	if params.InStock != nil {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"inventory.in_stock": *params.InStock,
			},
		})
	}
	
	// Фильтр по атрибутам (nested query)
	if len(params.Attributes) > 0 {
		for attrName, attrValue := range params.Attributes {
			filter = append(filter, map[string]interface{}{
				"nested": map[string]interface{}{
					"path": "attributes",
					"query": map[string]interface{}{
						"bool": map[string]interface{}{
							"must": []map[string]interface{}{
								{
									"term": map[string]interface{}{
										"attributes.name": attrName,
									},
								},
								{
									"match": map[string]interface{}{
										"attributes.value": attrValue,
									},
								},
							},
						},
					},
				},
			})
		}
	}
	
	// Фильтр по городу
	if params.City != "" {
		filter = append(filter, map[string]interface{}{
			"match": map[string]interface{}{
				"city": params.City,
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
	
	// Фильтр по верифицированным витринам
	if params.OnlyVerified {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"storefront.is_verified": true,
			},
		})
	}
	
	// Фильтр по качеству карточки
	if params.MinQualityScore > 0 {
		filter = append(filter, map[string]interface{}{
			"range": map[string]interface{}{
				"quality_score": map[string]interface{}{
					"gte": params.MinQualityScore,
				},
			},
		})
	}
	
	// Только активные товары
	filter = append(filter, map[string]interface{}{
		"term": map[string]interface{}{
			"is_active": true,
		},
	})
	
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
	sort := r.buildSort(params)
	if len(sort) > 0 {
		query["sort"] = sort
	}
	
	// Подсветка результатов
	query["highlight"] = map[string]interface{}{
		"fields": map[string]interface{}{
			"name": map[string]interface{}{
				"pre_tags":  []string{"<mark>"},
				"post_tags": []string{"</mark>"},
			},
			"description": map[string]interface{}{
				"fragment_size": 150,
				"pre_tags":      []string{"<mark>"},
				"post_tags":     []string{"</mark>"},
			},
		},
	}
	
	// Агрегации
	if len(params.Aggregations) > 0 {
		aggs := r.buildAggregations(params.Aggregations)
		if len(aggs) > 0 {
			query["aggs"] = aggs
		}
	}
	
	return query
}

// buildSort строит параметры сортировки
func (r *ProductRepository) buildSort(params *ProductSearchParams) []map[string]interface{} {
	sort := []map[string]interface{}{}
	
	switch params.SortBy {
	case "price":
		sort = append(sort, map[string]interface{}{
			"price": map[string]interface{}{
				"order": params.SortOrder,
			},
		})
	case "popularity":
		sort = append(sort, map[string]interface{}{
			"popularity_score": map[string]interface{}{
				"order": params.SortOrder,
			},
		})
	case "quality":
		sort = append(sort, map[string]interface{}{
			"quality_score": map[string]interface{}{
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
	case "rating":
		sort = append(sort, map[string]interface{}{
			"storefront.rating": map[string]interface{}{
				"order": params.SortOrder,
			},
		})
	default:
		// По умолчанию сортировка по релевантности и популярности
		if params.Query != "" {
			sort = append(sort, map[string]interface{}{
				"_score": map[string]interface{}{
					"order": "desc",
				},
			})
		}
		sort = append(sort, map[string]interface{}{
			"popularity_score": map[string]interface{}{
				"order": "desc",
			},
		})
	}
	
	// Всегда добавляем ID для стабильной сортировки
	sort = append(sort, map[string]interface{}{
		"product_id": map[string]interface{}{
			"order": "asc",
		},
	})
	
	return sort
}

// buildAggregations строит агрегации для фасетного поиска
func (r *ProductRepository) buildAggregations(requested []string) map[string]interface{} {
	aggs := map[string]interface{}{}
	
	for _, aggName := range requested {
		switch aggName {
		case "categories":
			aggs["categories"] = map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "category_id",
					"size":  20,
				},
			}
		case "brands":
			aggs["brands"] = map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "brand.keyword",
					"size":  20,
				},
			}
		case "price_ranges":
			aggs["price_ranges"] = map[string]interface{}{
				"range": map[string]interface{}{
					"field": "price",
					"ranges": []map[string]interface{}{
						{"to": 100},
						{"from": 100, "to": 500},
						{"from": 500, "to": 1000},
						{"from": 1000, "to": 5000},
						{"from": 5000},
					},
				},
			}
		case "attributes":
			aggs["attributes"] = map[string]interface{}{
				"nested": map[string]interface{}{
					"path": "attributes",
				},
				"aggs": map[string]interface{}{
					"by_name": map[string]interface{}{
						"terms": map[string]interface{}{
							"field": "attributes.name",
							"size":  50,
						},
						"aggs": map[string]interface{}{
							"values": map[string]interface{}{
								"terms": map[string]interface{}{
									"field": "attributes.value.keyword",
									"size":  20,
								},
							},
						},
					},
				},
			}
		case "storefronts":
			aggs["storefronts"] = map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "storefront_id",
					"size":  20,
				},
			}
		case "in_stock":
			aggs["in_stock"] = map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "inventory.in_stock",
				},
			}
		}
	}
	
	return aggs
}

// parseSearchResponse парсит ответ от OpenSearch
func (r *ProductRepository) parseSearchResponse(responseBytes []byte, params *ProductSearchParams) (*ProductSearchResult, error) {
	var response map[string]interface{}
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return nil, fmt.Errorf("ошибка разбора ответа: %w", err)
	}
	
	result := &ProductSearchResult{
		Products:     []*ProductSearchItem{},
		Total:        0,
		Aggregations: map[string]interface{}{},
	}
	
	// Время выполнения
	if took, ok := response["took"].(float64); ok {
		result.TookMs = int64(took)
	}
	
	// Извлекаем результаты
	if hits, ok := response["hits"].(map[string]interface{}); ok {
		// Общее количество
		if total, ok := hits["total"].(map[string]interface{}); ok {
			if value, ok := total["value"].(float64); ok {
				result.Total = int(value)
			}
		}
		
		// Документы
		if hitsArray, ok := hits["hits"].([]interface{}); ok {
			for _, hit := range hitsArray {
				if hitObj, ok := hit.(map[string]interface{}); ok {
					item := r.parseSearchHit(hitObj)
					result.Products = append(result.Products, item)
				}
			}
		}
	}
	
	// Извлекаем агрегации
	if aggs, ok := response["aggregations"].(map[string]interface{}); ok {
		result.Aggregations = aggs
	}
	
	return result, nil
}

// parseSearchHit парсит один результат поиска
func (r *ProductRepository) parseSearchHit(hit map[string]interface{}) *ProductSearchItem {
	item := &ProductSearchItem{
		Images:     []ProductImage{},
		Attributes: []ProductAttribute{},
		Variants:   []ProductVariant{},
		Highlights: map[string][]string{},
	}
	
	// ID документа
	if id, ok := hit["_id"].(string); ok {
		item.ID = id
		// Извлекаем ProductID из ID документа (sp_123 -> 123)
		if strings.HasPrefix(id, "sp_") {
			if productID, err := strconv.Atoi(id[3:]); err == nil {
				item.ProductID = productID
			}
		}
	}
	
	// Score
	if score, ok := hit["_score"].(float64); ok {
		item.Score = score
	}
	
	// Distance (из сортировки)
	if sort, ok := hit["sort"].([]interface{}); ok && len(sort) > 0 {
		if distance, ok := sort[0].(float64); ok {
			item.Distance = &distance
		}
	}
	
	// Highlights
	if highlight, ok := hit["highlight"].(map[string]interface{}); ok {
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
	
	// Парсим _source
	if source, ok := hit["_source"].(map[string]interface{}); ok {
		r.parseProductSource(source, item)
	}
	
	return item
}

// parseProductSource парсит данные товара из _source
func (r *ProductRepository) parseProductSource(source map[string]interface{}, item *ProductSearchItem) {
	// Основные поля
	if v, ok := source["storefront_id"].(float64); ok {
		item.StorefrontID = int(v)
	}
	if v, ok := source["name"].(string); ok {
		item.Name = v
	}
	if v, ok := source["description"].(string); ok {
		item.Description = v
	}
	if v, ok := source["price"].(float64); ok {
		item.Price = v
	}
	if v, ok := source["price_min"].(float64); ok {
		item.PriceMin = v
	}
	if v, ok := source["price_max"].(float64); ok {
		item.PriceMax = v
	}
	if v, ok := source["currency"].(string); ok {
		item.Currency = v
	}
	if v, ok := source["sku"].(string); ok {
		item.SKU = v
	}
	if v, ok := source["brand"].(string); ok {
		item.Brand = v
	}
	if v, ok := source["popularity_score"].(float64); ok {
		item.PopularityScore = v
	}
	if v, ok := source["quality_score"].(float64); ok {
		item.QualityScore = v
	}
	
	// Инвентарь
	if inventory, ok := source["inventory"].(map[string]interface{}); ok {
		if v, ok := inventory["in_stock"].(bool); ok {
			item.InStock = v
		}
		if v, ok := inventory["available"].(float64); ok {
			item.AvailableQuantity = int(v)
		}
	}
	
	// Изображения
	if images, ok := source["images"].([]interface{}); ok {
		for _, img := range images {
			if imgMap, ok := img.(map[string]interface{}); ok {
				image := ProductImage{}
				if v, ok := imgMap["id"].(float64); ok {
					image.ID = int(v)
				}
				if v, ok := imgMap["url"].(string); ok {
					image.URL = v
				}
				if v, ok := imgMap["alt_text"].(string); ok {
					image.AltText = v
				}
				if v, ok := imgMap["is_main"].(bool); ok {
					image.IsMain = v
				}
				if v, ok := imgMap["position"].(float64); ok {
					image.Position = int(v)
				}
				item.Images = append(item.Images, image)
			}
		}
	}
	
	// Витрина
	if storefront, ok := source["storefront"].(map[string]interface{}); ok {
		info := StorefrontInfo{}
		if v, ok := storefront["id"].(float64); ok {
			info.ID = int(v)
		}
		if v, ok := storefront["name"].(string); ok {
			info.Name = v
		}
		if v, ok := storefront["slug"].(string); ok {
			info.Slug = v
		}
		if v, ok := storefront["city"].(string); ok {
			info.City = v
		}
		if v, ok := storefront["country"].(string); ok {
			info.Country = v
		}
		if v, ok := storefront["rating"].(float64); ok {
			info.Rating = v
		}
		if v, ok := storefront["is_verified"].(bool); ok {
			info.IsVerified = v
		}
		item.Storefront = info
	}
	
	// Категория
	if category, ok := source["category"].(map[string]interface{}); ok {
		info := CategoryInfo{}
		if v, ok := category["id"].(float64); ok {
			info.ID = int(v)
		}
		if v, ok := category["name"].(string); ok {
			info.Name = v
		}
		if v, ok := category["slug"].(string); ok {
			info.Slug = v
		}
		item.Category = info
	}
	
	// Атрибуты
	if attributes, ok := source["attributes"].([]interface{}); ok {
		for _, attr := range attributes {
			if attrMap, ok := attr.(map[string]interface{}); ok {
				attribute := ProductAttribute{}
				if v, ok := attrMap["id"].(float64); ok {
					attribute.ID = int(v)
				}
				if v, ok := attrMap["name"].(string); ok {
					attribute.Name = v
				}
				if v, ok := attrMap["type"].(string); ok {
					attribute.Type = v
				}
				if v, ok := attrMap["value"]; ok {
					attribute.Value = v
				}
				if v, ok := attrMap["display_name"].(string); ok {
					attribute.DisplayName = v
				}
				item.Attributes = append(item.Attributes, attribute)
			}
		}
	}
	
	// Варианты
	if variants, ok := source["variants"].([]interface{}); ok {
		for _, variant := range variants {
			if varMap, ok := variant.(map[string]interface{}); ok {
				var_ := ProductVariant{}
				if v, ok := varMap["id"].(float64); ok {
					var_.ID = int(v)
				}
				if v, ok := varMap["name"].(string); ok {
					var_.Name = v
				}
				if v, ok := varMap["sku"].(string); ok {
					var_.SKU = v
				}
				if v, ok := varMap["price"].(float64); ok {
					var_.Price = v
				}
				if v, ok := varMap["attributes"].(map[string]interface{}); ok {
					var_.Attributes = v
				}
				item.Variants = append(item.Variants, var_)
			}
		}
	}
}

// calculatePopularityScore рассчитывает оценку популярности товара
func calculatePopularityScore(product *models.StorefrontProduct) float64 {
	score := 0.0
	
	// Учитываем продажи
	if product.SoldCount > 0 {
		score += float64(product.SoldCount) * 0.5
	}
	
	// Учитываем просмотры
	if product.ViewCount > 0 {
		score += float64(product.ViewCount) * 0.01
	}
	
	// TODO: Добавить учет рейтинга витрины когда будет загружаться информация о витрине
	
	return score
}

// calculateQualityScore рассчитывает оценку качества карточки товара
func calculateQualityScore(product *models.StorefrontProduct) float64 {
	score := 0.0
	
	// Есть описание
	if product.Description != "" {
		score += 20
		if len(product.Description) > 100 {
			score += 10
		}
	}
	
	// Есть изображения
	if product.Images != nil && len(product.Images) > 0 {
		score += 20
		if len(product.Images) > 3 {
			score += 10
		}
	}
	
	// Есть атрибуты
	if product.Attributes != nil && len(product.Attributes) > 0 {
		score += 15
		if len(product.Attributes) > 5 {
			score += 10
		}
	}
	
	// Есть варианты
	if product.Variants != nil && len(product.Variants) > 0 {
		score += 15
	}
	
	// TODO: Добавить бонус за верифицированную витрину когда будет загружаться информация о витрине
	
	return score
}

// storefrontProductMapping маппинг для индекса товаров витрин
const storefrontProductMapping = `{
  "mappings": {
    "properties": {
      "product_id": {"type": "integer"},
      "product_type": {"type": "keyword"},
      "storefront_id": {"type": "integer"},
      "category_id": {"type": "integer"},
      "name": {
        "type": "text",
        "fields": {
          "keyword": {"type": "keyword"},
          "autocomplete": {
            "type": "search_as_you_type"
          }
        }
      },
      "name_lowercase": {"type": "text"},
      "description": {"type": "text"},
      "price": {"type": "float"},
      "price_min": {"type": "float"},
      "price_max": {"type": "float"},
      "currency": {"type": "keyword"},
      "sku": {"type": "keyword"},
      "barcode": {"type": "keyword"},
      "status": {"type": "keyword"},
      "is_active": {"type": "boolean"},
      "brand": {
        "type": "text",
        "fields": {
          "keyword": {"type": "keyword"}
        }
      },
      "brand_lowercase": {"type": "text"},
      "model": {
        "type": "text",
        "fields": {
          "keyword": {"type": "keyword"}
        }
      },
      "model_lowercase": {"type": "text"},
      "color": {"type": "keyword"},
      "size": {"type": "keyword"},
      "material": {"type": "keyword"},
      "inventory": {
        "properties": {
          "track": {"type": "boolean"},
          "count": {"type": "integer"},
          "reserved": {"type": "integer"},
          "available": {"type": "integer"},
          "in_stock": {"type": "boolean"},
          "low_stock": {"type": "boolean"}
        }
      },
      "attributes": {
        "type": "nested",
        "properties": {
          "id": {"type": "integer"},
          "name": {"type": "keyword"},
          "type": {"type": "keyword"},
          "value": {"type": "text"},
          "display_name": {"type": "text"}
        }
      },
      "metadata": {"type": "object", "enabled": false},
      "images": {
        "properties": {
          "id": {"type": "integer"},
          "url": {"type": "keyword"},
          "alt_text": {"type": "text"},
          "is_main": {"type": "boolean"},
          "position": {"type": "integer"}
        }
      },
      "has_images": {"type": "boolean"},
      "image_count": {"type": "integer"},
      "variants": {
        "properties": {
          "id": {"type": "integer"},
          "name": {"type": "text"},
          "sku": {"type": "keyword"},
          "price": {"type": "float"},
          "attributes": {"type": "object"}
        }
      },
      "has_variants": {"type": "boolean"},
      "variant_count": {"type": "integer"},
      "storefront": {
        "properties": {
          "id": {"type": "integer"},
          "name": {"type": "text"},
          "slug": {"type": "keyword"},
          "city": {"type": "text"},
          "country": {"type": "keyword"},
          "rating": {"type": "float"},
          "is_verified": {"type": "boolean"}
        }
      },
      "category": {
        "properties": {
          "id": {"type": "integer"},
          "name": {"type": "text"},
          "slug": {"type": "keyword"}
        }
      },
      "category_path": {"type": "text"},
      "location": {"type": "geo_point"},
      "city": {"type": "text"},
      "country": {"type": "keyword"},
      "search_keywords": {"type": "text"},
      "popularity_score": {"type": "float"},
      "quality_score": {"type": "float"},
      "created_at": {"type": "date"},
      "updated_at": {"type": "date"}
    }
  },
  "settings": {
    "analysis": {
      "analyzer": {
        "russian_analyzer": {
          "tokenizer": "standard",
          "filter": [
            "lowercase",
            "russian_stop",
            "russian_stemmer"
          ]
        }
      },
      "filter": {
        "russian_stop": {
          "type": "stop",
          "stopwords": "_russian_"
        },
        "russian_stemmer": {
          "type": "stemmer",
          "language": "russian"
        }
      }
    }
  }
}`