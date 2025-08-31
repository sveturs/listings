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
	exists, err := r.client.IndexExists(ctx, r.indexName)
	if err != nil {
		return fmt.Errorf("ошибка проверки индекса товаров витрин: %w", err)
	}

	logger.Info().Str("indexName", r.indexName).Bool("exists", exists).Msg("Проверка индекса товаров витрин")

	if !exists {
		logger.Info().Str("indexName", r.indexName).Msg("Создание индекса товаров витрин...")
		if err := r.client.CreateIndex(ctx, r.indexName, storefrontProductMapping); err != nil {
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

	return r.client.IndexDocument(ctx, r.indexName, docID, doc)
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
	return r.client.BulkIndex(ctx, r.indexName, docs)
}

// DeleteProduct удаляет товар из индекса
func (r *ProductRepository) DeleteProduct(ctx context.Context, productID int) error {
	docID := fmt.Sprintf("sp_%d", productID)
	return r.client.DeleteDocument(ctx, r.indexName, docID)
}

// UpdateProduct обновляет товар в индексе
func (r *ProductRepository) UpdateProduct(ctx context.Context, product *models.StorefrontProduct) error {
	return r.IndexProduct(ctx, product)
}

// productToDoc преобразует модель товара витрины в документ для индексации
func (r *ProductRepository) productToDoc(product *models.StorefrontProduct) map[string]interface{} {
	doc := map[string]interface{}{
		"product_id":     product.ID,
		"product_type":   "storefront",
		"storefront_id":  product.StorefrontID,
		"category_id":    product.CategoryID,
		"name":           product.Name,
		"name_lowercase": strings.ToLower(product.Name),
		"description":    product.Description,
		"price":          product.Price,
		"currency":       product.Currency,
		"sku":            product.SKU,
		"barcode":        product.Barcode,
		"stock_status":   product.StockStatus,
		"is_active":      product.IsActive,
		"status":         "active", // Добавляем статус для совместимости с фильтром поиска
		"created_at":     product.CreatedAt.Format(time.RFC3339),
		"updated_at":     product.UpdatedAt.Format(time.RFC3339),
	}

	// Добавляем информацию о наличии товара
	doc["stock_quantity"] = product.StockQuantity
	doc["inventory"] = map[string]interface{}{
		"quantity":  product.StockQuantity,
		"available": product.StockQuantity, // Доступное количество равно общему количеству минус резервы (пока резервы не реализованы)
		"in_stock":  product.StockQuantity > 0,
		"low_stock": product.StockQuantity > 0 && product.StockQuantity <= 5, // TODO: сделать настраиваемым
		"status":    product.StockStatus,
	}

	// Добавляем атрибуты товара
	if len(product.Attributes) > 0 {
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
	if len(product.Images) > 0 {
		imagesArray := make([]map[string]interface{}, 0, len(product.Images))
		for _, img := range product.Images {
			// Используем PublicURL для правильного отображения изображений
			imageURL := img.PublicURL
			if imageURL == "" {
				// Fallback на ImageURL если PublicURL пустой
				imageURL = img.ImageURL
			}

			imgDoc := map[string]interface{}{
				"id":            img.ID,
				"url":           imageURL,
				"thumbnail_url": img.ThumbnailURL,
				"is_main":       img.IsDefault, // Use is_main for consistency with unified search
				"is_default":    img.IsDefault, // Keep for backward compatibility
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
	if len(product.Variants) > 0 {
		variantsArray := make([]map[string]interface{}, 0, len(product.Variants))
		minPrice := product.Price
		maxPrice := product.Price
		totalVariantStock := 0
		totalVariantReserved := 0
		hasStockVariants := false

		for _, variant := range product.Variants {
			// Расчет доступного количества с учетом резервирований
			availableQuantity := variant.StockQuantity
			// TODO: добавить логику расчета зарезервированного количества из inventory_reservations
			// Пока используем только stock_quantity

			varDoc := map[string]interface{}{
				"id":                 variant.ID,
				"name":               variant.Name,
				"sku":                variant.SKU,
				"price":              variant.Price,
				"attributes":         variant.Attributes,
				"stock_quantity":     variant.StockQuantity,
				"available_quantity": availableQuantity,
				"is_active":          variant.IsActive,
				"created_at":         variant.CreatedAt.Format(time.RFC3339),
				"updated_at":         variant.UpdatedAt.Format(time.RFC3339),
			}

			// Определяем статус наличия варианта
			var stockStatus string
			switch {
			case availableQuantity <= 0:
				stockStatus = "out_of_stock"
			case availableQuantity <= 5: // TODO: сделать настраиваемым
				stockStatus = "low_stock"
			default:
				stockStatus = "in_stock"
			}
			varDoc["stock_status"] = stockStatus

			// Добавляем инвентарь варианта
			varDoc["inventory"] = map[string]interface{}{
				"quantity":  variant.StockQuantity,
				"available": availableQuantity,
				"reserved":  variant.StockQuantity - availableQuantity, // расчет резервов
				"in_stock":  availableQuantity > 0,
				"low_stock": availableQuantity > 0 && availableQuantity <= 5,
				"status":    stockStatus,
			}

			// Отслеживаем диапазон цен
			if variant.Price < minPrice {
				minPrice = variant.Price
			}
			if variant.Price > maxPrice {
				maxPrice = variant.Price
			}

			// Считаем общие остатки по вариантам
			totalVariantStock += variant.StockQuantity
			totalVariantReserved += (variant.StockQuantity - availableQuantity)
			if availableQuantity > 0 {
				hasStockVariants = true
			}

			variantsArray = append(variantsArray, varDoc)
		}

		doc["variants"] = variantsArray
		doc["has_variants"] = true
		doc["variant_count"] = len(product.Variants)
		doc["price_min"] = minPrice
		doc["price_max"] = maxPrice

		// Обновляем общую информацию о наличии с учетом вариантов
		doc["total_variant_stock"] = totalVariantStock
		doc["total_variant_reserved"] = totalVariantReserved
		doc["total_variant_available"] = totalVariantStock - totalVariantReserved
		doc["has_available_variants"] = hasStockVariants

		// Перезаписываем inventory с учетом вариантов
		doc["inventory"] = map[string]interface{}{
			"quantity":           totalVariantStock,
			"available":          totalVariantStock - totalVariantReserved,
			"reserved":           totalVariantReserved,
			"in_stock":           hasStockVariants,
			"low_stock":          hasStockVariants && (totalVariantStock-totalVariantReserved) <= 5,
			"status":             product.StockStatus,
			"has_variants":       true,
			"variant_count":      len(product.Variants),
			"has_stock_variants": hasStockVariants,
		}
	} else {
		doc["has_variants"] = false
		doc["variant_count"] = 0
		doc["price_min"] = product.Price
		doc["price_max"] = product.Price
		doc["total_variant_stock"] = 0
		doc["total_variant_reserved"] = 0
		doc["total_variant_available"] = 0
		doc["has_available_variants"] = false
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

	responseBytes, err := r.client.Search(ctx, r.indexName, query)
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
			"bool": map[string]interface{}{
				"should": []map[string]interface{}{
					// Поиск по основным полям товара
					{
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
							"type":          "best_fields",
							"fuzziness":     "AUTO",
							"prefix_length": 2,
							"boost":         2.0,
						},
					},
					// Поиск по вариантам
					{
						"nested": map[string]interface{}{
							"path": "variants",
							"query": map[string]interface{}{
								"multi_match": map[string]interface{}{
									"query": params.Query,
									"fields": []string{
										"variants.name^2",
										"variants.sku^1.5",
									},
									"type":          "best_fields",
									"fuzziness":     "AUTO",
									"prefix_length": 2,
								},
							},
							"boost": 1.5,
						},
					},
				},
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

	// Фильтр по цене (с учетом вариантов)
	if params.PriceMin > 0 || params.PriceMax > 0 {
		priceRange := map[string]interface{}{}
		if params.PriceMin > 0 {
			priceRange["gte"] = params.PriceMin
		}
		if params.PriceMax > 0 {
			priceRange["lte"] = params.PriceMax
		}

		// Ищем товары, у которых либо основная цена в диапазоне, либо есть варианты в диапазоне
		filter = append(filter, map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []map[string]interface{}{
					// Основная цена товара в диапазоне
					{
						"range": map[string]interface{}{
							"price": priceRange,
						},
					},
					// Минимальная цена вариантов в диапазоне
					{
						"range": map[string]interface{}{
							"price_min": priceRange,
						},
					},
					// Максимальная цена вариантов в диапазоне
					{
						"range": map[string]interface{}{
							"price_max": priceRange,
						},
					},
					// Вариант с ценой в диапазоне
					{
						"nested": map[string]interface{}{
							"path": "variants",
							"query": map[string]interface{}{
								"range": map[string]interface{}{
									"variants.price": priceRange,
								},
							},
						},
					},
				},
			},
		})
	}

	// Фильтр по наличию
	if params.InStock != nil {
		if *params.InStock {
			// Если ищем товары в наличии, учитываем варианты
			filter = append(filter, map[string]interface{}{
				"bool": map[string]interface{}{
					"should": []map[string]interface{}{
						// Товары без вариантов, которые в наличии
						{
							"bool": map[string]interface{}{
								"must": []map[string]interface{}{
									{"term": map[string]interface{}{"has_variants": false}},
									{"term": map[string]interface{}{"inventory.in_stock": true}},
								},
							},
						},
						// Товары с вариантами, у которых есть доступные варианты
						{
							"bool": map[string]interface{}{
								"must": []map[string]interface{}{
									{"term": map[string]interface{}{"has_variants": true}},
									{"term": map[string]interface{}{"has_available_variants": true}},
								},
							},
						},
					},
				},
			})
		} else {
			// Ищем товары не в наличии
			filter = append(filter, map[string]interface{}{
				"bool": map[string]interface{}{
					"should": []map[string]interface{}{
						// Товары без вариантов, которые не в наличии
						{
							"bool": map[string]interface{}{
								"must": []map[string]interface{}{
									{"term": map[string]interface{}{"has_variants": false}},
									{"term": map[string]interface{}{"inventory.in_stock": false}},
								},
							},
						},
						// Товары с вариантами, у которых нет доступных вариантов
						{
							"bool": map[string]interface{}{
								"must": []map[string]interface{}{
									{"term": map[string]interface{}{"has_variants": true}},
									{"term": map[string]interface{}{"has_available_variants": false}},
								},
							},
						},
					},
				},
			})
		}
	}

	// Фильтр по атрибутам (nested query)
	if len(params.Attributes) > 0 {
		for attrName, attrValue := range params.Attributes {
			// Ищем атрибуты как у товара, так и у вариантов
			filter = append(filter, map[string]interface{}{
				"bool": map[string]interface{}{
					"should": []map[string]interface{}{
						// Атрибуты товара
						{
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
						},
						// Атрибуты в вариантах товара
						{
							"nested": map[string]interface{}{
								"path": "variants",
								"query": map[string]interface{}{
									"bool": map[string]interface{}{
										"must": []map[string]interface{}{
											{
												"exists": map[string]interface{}{
													"field": fmt.Sprintf("variants.attributes.%s", attrName),
												},
											},
											{
												"match": map[string]interface{}{
													fmt.Sprintf("variants.attributes.%s", attrName): attrValue,
												},
											},
										},
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

	// Если нет текстового поиска, добавляем match_all
	if len(must) == 0 {
		must = append(must, map[string]interface{}{
			"match_all": map[string]interface{}{},
		})
	}

	boolQuery["must"] = must

	if len(filter) > 0 {
		boolQuery["filter"] = filter
	}

	query["query"] = map[string]interface{}{
		"bool": boolQuery,
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
			"variants.name": map[string]interface{}{
				"pre_tags":  []string{"<mark>"},
				"post_tags": []string{"</mark>"},
			},
			"variants.sku": map[string]interface{}{
				"pre_tags":  []string{"<mark>"},
				"post_tags": []string{"</mark>"},
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
		// Используем created_at как фолбэк для популярности
		sort = append(sort, map[string]interface{}{
			"created_at": map[string]interface{}{
				"order": params.SortOrder,
			},
		})
	case "quality":
		// Используем created_at как фолбэк для качества
		sort = append(sort, map[string]interface{}{
			"created_at": map[string]interface{}{
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
		// По умолчанию сортировка по релевантности и дате создания
		if params.Query != "" {
			sort = append(sort, map[string]interface{}{
				"_score": map[string]interface{}{
					"order": "desc",
				},
			})
		}
		// Используем created_at вместо popularity_score
		sort = append(sort, map[string]interface{}{
			"created_at": map[string]interface{}{
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
		case "variants":
			aggs["variants"] = map[string]interface{}{
				"nested": map[string]interface{}{
					"path": "variants",
				},
				"aggs": map[string]interface{}{
					"available_variants": map[string]interface{}{
						"filter": map[string]interface{}{
							"term": map[string]interface{}{
								"variants.inventory.in_stock": true,
							},
						},
					},
					"variant_attributes": map[string]interface{}{
						"terms": map[string]interface{}{
							"script": map[string]interface{}{
								"source": "params._source.variants.attributes.keySet()",
							},
							"size": 50,
						},
						"aggs": map[string]interface{}{
							"values": map[string]interface{}{
								"terms": map[string]interface{}{
									"script": map[string]interface{}{
										"source": "params._source.variants.attributes.values()",
									},
									"size": 20,
								},
							},
						},
					},
					"price_range": map[string]interface{}{
						"stats": map[string]interface{}{
							"field": "variants.price",
						},
					},
				},
			}
		case "has_variants":
			aggs["has_variants"] = map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "has_variants",
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
		logger.Debug().
			Str("id", item.ID).
			Interface("full_source", source).
			Msg("Full OpenSearch document source")
		r.parseProductSource(source, item)
	}

	return item
}

// parseProductSource парсит данные товара из _source
func (r *ProductRepository) parseProductSource(source map[string]interface{}, item *ProductSearchItem) {
	// Основные поля
	// Сначала проверяем product_id
	if v, ok := source["product_id"].(float64); ok {
		item.ProductID = int(v)
	}
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
		logger.Debug().
			Int("product_id", item.ProductID).
			Interface("inventory", inventory).
			Msg("Parsing inventory from OpenSearch")

		if v, ok := inventory["in_stock"].(bool); ok {
			item.InStock = v
		}
		if v, ok := inventory["available"].(float64); ok {
			item.AvailableQuantity = int(v)
			logger.Debug().
				Int("product_id", item.ProductID).
				Int("available_quantity", item.AvailableQuantity).
				Msg("Set available quantity from inventory")
		} else {
			logger.Debug().
				Int("product_id", item.ProductID).
				Interface("available_value", inventory["available"]).
				Msg("Failed to parse available quantity")
		}
	} else {
		logger.Debug().
			Int("product_id", item.ProductID).
			Msg("No inventory data found in OpenSearch response")
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
				if v, ok := imgMap["is_default"].(bool); ok {
					image.IsMain = v
				}
				if v, ok := imgMap["display_order"].(float64); ok {
					image.Position = int(v)
				}
				if v, ok := imgMap["position"].(float64); ok {
					image.Position = int(v)
				}
				item.Images = append(item.Images, image)
			}
		}
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
				// Новые поля
				if v, ok := varMap["stock_quantity"].(float64); ok {
					var_.StockQuantity = int(v)
				}
				if v, ok := varMap["available_quantity"].(float64); ok {
					var_.AvailableQuantity = int(v)
				}
				if v, ok := varMap["stock_status"].(string); ok {
					var_.StockStatus = v
				}
				if v, ok := varMap["is_active"].(bool); ok {
					var_.IsActive = v
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
	if len(product.Images) > 0 {
		score += 20
		if len(product.Images) > 3 {
			score += 10
		}
	}

	// Есть атрибуты
	if len(product.Attributes) > 0 {
		score += 15
		if len(product.Attributes) > 5 {
			score += 10
		}
	}

	// Есть варианты
	if len(product.Variants) > 0 {
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
          "quantity": {"type": "integer"},
          "count": {"type": "integer"},
          "reserved": {"type": "integer"},
          "available": {"type": "integer"},
          "in_stock": {"type": "boolean"},
          "low_stock": {"type": "boolean"},
          "status": {"type": "keyword"},
          "has_variants": {"type": "boolean"},
          "variant_count": {"type": "integer"},
          "has_stock_variants": {"type": "boolean"}
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
          "is_default": {"type": "boolean"},
          "position": {"type": "integer"},
          "display_order": {"type": "integer"},
          "thumbnail_url": {"type": "keyword"}
        }
      },
      "has_images": {"type": "boolean"},
      "image_count": {"type": "integer"},
      "variants": {
        "type": "nested",
        "properties": {
          "id": {"type": "integer"},
          "name": {
            "type": "text",
            "fields": {
              "keyword": {"type": "keyword"}
            }
          },
          "sku": {"type": "keyword"},
          "price": {"type": "float"},
          "stock_quantity": {"type": "integer"},
          "available_quantity": {"type": "integer"},
          "stock_status": {"type": "keyword"},
          "is_active": {"type": "boolean"},
          "attributes": {"type": "object"},
          "inventory": {
            "properties": {
              "quantity": {"type": "integer"},
              "available": {"type": "integer"},
              "reserved": {"type": "integer"},
              "in_stock": {"type": "boolean"},
              "low_stock": {"type": "boolean"},
              "status": {"type": "keyword"}
            }
          },
          "created_at": {"type": "date"},
          "updated_at": {"type": "date"}
        }
      },
      "has_variants": {"type": "boolean"},
      "variant_count": {"type": "integer"},
      "total_variant_stock": {"type": "integer"},
      "total_variant_reserved": {"type": "integer"},
      "total_variant_available": {"type": "integer"},
      "has_available_variants": {"type": "boolean"},
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

// SearchSimilarProducts выполняет поиск похожих товаров витрин
func (r *ProductRepository) SearchSimilarProducts(ctx context.Context, productID int, limit int) ([]*models.MarketplaceListing, error) {
	logger.Info().Msgf("SearchSimilarProducts: поиск похожих товаров для продукта витрины ID=%d", productID)

	// Выполняем базовый поиск для получения информации о исходном товаре
	// Используем простой запрос для поиска товара по ID
	getQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"product_id": productID,
			},
		},
		"size": 1,
	}

	getResponseBytes, err := r.client.Search(ctx, r.indexName, getQuery)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения товара из индекса: %w", err)
	}

	// Парсим ответ для получения информации о товаре
	var getResult map[string]interface{}
	if unmarshalErr := json.Unmarshal(getResponseBytes, &getResult); unmarshalErr != nil {
		return nil, fmt.Errorf("ошибка парсинга ответа получения товара: %w", unmarshalErr)
	}

	hits, ok := getResult["hits"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("товар с ID %d не найден в индексе", productID)
	}

	hitsArray, ok := hits["hits"].([]interface{})
	if !ok || len(hitsArray) == 0 {
		return nil, fmt.Errorf("товар с ID %d не найден в индексе", productID)
	}

	sourceHit, ok := hitsArray[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("неверный формат документа товара")
	}

	sourceProduct, ok := sourceHit["_source"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("не удается получить данные товара")
	}

	// Извлекаем данные исходного товара
	sourceName, _ := sourceProduct["name"].(string)
	sourcePrice, _ := sourceProduct["price"].(float64)
	sourceCategoryID, _ := sourceProduct["category_id"].(float64)
	sourceStorefrontID, _ := sourceProduct["storefront_id"].(float64)
	sourceBrand, _ := sourceProduct["brand"].(string)
	sourceModel, _ := sourceProduct["model"].(string)

	logger.Info().Msgf("Исходный товар: Name=%s, CategoryID=%.0f, Price=%.2f, StorefrontID=%.0f, Brand=%s, Model=%s",
		sourceName, sourceCategoryID, sourcePrice, sourceStorefrontID, sourceBrand, sourceModel)

	// Создаем запрос для поиска похожих товаров
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					// Исключаем сам товар
					{
						"bool": map[string]interface{}{
							"must_not": []map[string]interface{}{
								{
									"term": map[string]interface{}{
										"product_id": productID,
									},
								},
							},
						},
					},
				},
				"should": []map[string]interface{}{
					// Точное совпадение категории (высокий приоритет)
					{
						"term": map[string]interface{}{
							"category_id": map[string]interface{}{
								"value": sourceCategoryID,
								"boost": 3.0,
							},
						},
					},
					// Точное совпадение бренда
					{
						"term": map[string]interface{}{
							"brand_lowercase": map[string]interface{}{
								"value": strings.ToLower(sourceBrand),
								"boost": 2.5,
							},
						},
					},
					// Точное совпадение модели
					{
						"term": map[string]interface{}{
							"model_lowercase": map[string]interface{}{
								"value": strings.ToLower(sourceModel),
								"boost": 2.0,
							},
						},
					},
					// Частичное совпадение названия
					{
						"match": map[string]interface{}{
							"name_lowercase": map[string]interface{}{
								"query": strings.ToLower(sourceName),
								"boost": 1.5,
							},
						},
					},
					// Похожий диапазон цен (±50%)
					{
						"range": map[string]interface{}{
							"price": map[string]interface{}{
								"gte":   sourcePrice * 0.5,
								"lte":   sourcePrice * 1.5,
								"boost": 1.0,
							},
						},
					},
				},
				"minimum_should_match": 1, // Хотя бы один критерий должен совпадать
			},
		},
		"size": limit * 2, // Берем больше для фильтрации
		"sort": []map[string]interface{}{
			{"_score": map[string]interface{}{"order": "desc"}},
		},
	}

	// Выполняем поиск
	responseBytes, err := r.client.Search(ctx, r.indexName, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения поиска похожих товаров: %w", err)
	}

	// Parse response
	var response map[string]interface{}
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return nil, fmt.Errorf("ошибка парсинга ответа поиска: %w", err)
	}

	searchHits, ok := response["hits"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("неверный формат ответа OpenSearch")
	}

	searchHitsArray, ok := searchHits["hits"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("неверный формат массива hits в ответе OpenSearch")
	}

	logger.Info().Msgf("Найдено %d потенциально похожих товаров витрин", len(searchHitsArray))

	// Преобразуем результаты в формат MarketplaceListing
	var similarListings []*models.MarketplaceListing

	for _, hitInterface := range searchHitsArray {
		hit, ok := hitInterface.(map[string]interface{})
		if !ok {
			continue
		}

		source, ok := hit["_source"].(map[string]interface{})
		if !ok {
			continue
		}

		score, _ := hit["_score"].(float64)

		// Преобразуем документ товара витрины в MarketplaceListing
		listing := &models.MarketplaceListing{
			ID:          int(source["product_id"].(float64)),
			Title:       source["name"].(string),
			Description: source["description"].(string),
			Price:       source["price"].(float64),
			CategoryID:  int(source["category_id"].(float64)),
			Status:      "active",   // Товары витрин всегда активны
			CreatedAt:   time.Now(), // TODO: использовать реальную дату
			UpdatedAt:   time.Now(),
		}

		// Добавляем StorefrontID
		if storefrontID, ok := source["storefront_id"].(float64); ok {
			storefrontIDInt := int(storefrontID)
			listing.StorefrontID = &storefrontIDInt
		}

		// Добавляем атрибуты из товара витрины
		if attributes, ok := source["attributes"].(map[string]interface{}); ok {
			var listingAttrs []models.ListingAttributeValue

			// Конвертируем атрибуты товара витрины в формат атрибутов объявления
			for attrName, attrValue := range attributes {
				attr := models.ListingAttributeValue{
					ListingID:     listing.ID,
					AttributeName: attrName,
					DisplayValue:  fmt.Sprintf("%v", attrValue),
				}

				// Определяем тип значения
				switch v := attrValue.(type) {
				case string:
					attr.TextValue = &v
					attr.AttributeType = "text"
				case float64:
					attr.NumericValue = &v
					attr.AttributeType = "number"
				case bool:
					boolVal := v
					attr.BooleanValue = &boolVal
					attr.AttributeType = "boolean"
				default:
					strVal := fmt.Sprintf("%v", v)
					attr.TextValue = &strVal
					attr.AttributeType = "text"
				}

				listingAttrs = append(listingAttrs, attr)
			}

			listing.Attributes = listingAttrs
		}

		// Добавляем метаданные о скоре похожести для отладки
		if listing.Metadata == nil {
			listing.Metadata = make(map[string]interface{})
		}
		listing.Metadata["similarity_score"] = map[string]interface{}{
			"total":        score,
			"search_type":  "storefront_product",
			"source_index": r.indexName,
		}

		similarListings = append(similarListings, listing)

		// Ограничиваем результат
		if len(similarListings) >= limit {
			break
		}
	}

	logger.Info().Msgf("Найдено %d похожих товаров витрин для продукта %d", len(similarListings), productID)

	return similarListings, nil
}

// UpdateProductStock частично обновляет только поля склада товара в OpenSearch
func (r *ProductRepository) UpdateProductStock(ctx context.Context, productID int, stockData map[string]interface{}) error {
	if len(stockData) == 0 {
		return fmt.Errorf("no stock data to update")
	}

	// Для товаров витрин используем префикс sp_
	// Это соответствует формату ID при индексации в методе IndexProduct
	docID := fmt.Sprintf("sp_%d", productID)

	// Обновляем только переданные поля
	updateDoc := make(map[string]interface{})
	for field, value := range stockData {
		updateDoc[field] = value
	}

	// Добавляем timestamp обновления
	updateDoc["updated_at"] = time.Now().Format(time.RFC3339)

	if err := r.client.UpdateDocument(ctx, r.indexName, docID, updateDoc); err != nil {
		return fmt.Errorf("failed to update product stock in OpenSearch: %w", err)
	}

	logger.Info().Msgf("Successfully updated stock for storefront product sp_%d in OpenSearch index %s", productID, r.indexName)
	return nil
}
