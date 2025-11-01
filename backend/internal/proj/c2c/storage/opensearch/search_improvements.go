package opensearch

import (
	"context"
	"strings"

	"backend/internal/domain/search"
	"backend/internal/logger"
)

// buildImprovedSearchQuery создает улучшенный запрос с более точным поиском
func (r *Repository) buildImprovedSearchQuery(ctx context.Context, params *search.SearchParams) map[string]interface{} {
	logger.Info().Msgf("Строим улучшенный запрос: категория = %v, язык = %s, поисковый запрос = %s",
		params.CategoryID, params.Language, params.Query)

	// Validate and set default for size to prevent OpenSearch "numHits must be > 0" error
	size := params.Size
	if size <= 0 {
		size = 10 // Default size
		logger.Warn().Int("requested_size", params.Size).Int("default_size", size).Msg("Invalid size parameter, using default")
	}

	query := map[string]interface{}{
		"from": (params.Page - 1) * size,
		"size": size,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must":   []interface{}{},
				"filter": []interface{}{},
				"should": []interface{}{},
			},
		},
	}

	// Фильтр по статусу
	if params.Status == "" {
		filter := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]interface{})
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"status": "active",
			},
		})
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = filter
	}

	// Улучшенный текстовый поиск
	if params.Query != "" {
		logger.Info().Msgf("Улучшенный текстовый поиск по запросу: '%s'", params.Query)

		boolMap := query["query"].(map[string]interface{})["bool"].(map[string]interface{})
		must := boolMap["must"].([]interface{})
		should := boolMap["should"].([]interface{})

		// Получаем варианты транслитерации
		queryVariants := r.transliterator.TransliterateForSearch(params.Query)
		logger.Info().
			Str("original_query", params.Query).
			Strs("transliterated_variants", queryVariants).
			Msg("Generated transliteration variants for improved search")

		// Определяем тип поиска (точный или нечеткий)
		useExactMatch := false
		queryLower := strings.ToLower(params.Query)

		// Для коротких запросов (менее 4 символов) или известных брендов используем точное совпадение
		if len(params.Query) <= 3 || isKnownBrand(queryLower) {
			useExactMatch = true
		}

		if useExactMatch {
			// ТОЧНЫЙ ПОИСК - приоритет точным совпадениям
			logger.Info().Msg("Используем точный поиск для короткого запроса или известного бренда")

			// Создаем должен быть хотя бы один из этих запросов
			exactMatchQueries := []interface{}{}

			for _, queryVariant := range queryVariants {
				// Точное совпадение в заголовке (высший приоритет)
				exactMatchQueries = append(exactMatchQueries, map[string]interface{}{
					"match_phrase": map[string]interface{}{
						"title": map[string]interface{}{
							"query": queryVariant,
							"boost": 10.0,
						},
					},
				})

				// Точное совпадение в заголовке как префикс
				exactMatchQueries = append(exactMatchQueries, map[string]interface{}{
					"match_phrase_prefix": map[string]interface{}{
						"title": map[string]interface{}{
							"query": queryVariant,
							"boost": 8.0,
						},
					},
				})

				// Поиск по ключевым словам
				exactMatchQueries = append(exactMatchQueries, map[string]interface{}{
					"term": map[string]interface{}{
						"title.keyword": map[string]interface{}{
							"value": queryVariant,
							"boost": 12.0,
						},
					},
				})

				// Поиск в атрибутах (точное совпадение)
				exactMatchQueries = append(exactMatchQueries, map[string]interface{}{
					"nested": map[string]interface{}{
						"path": "attributes",
						"query": map[string]interface{}{
							"bool": map[string]interface{}{
								"should": []map[string]interface{}{
									{
										"term": map[string]interface{}{
											"attributes.text_value.keyword": map[string]interface{}{
												"value": queryVariant,
												"boost": 8.0,
											},
										},
									},
									{
										"match_phrase": map[string]interface{}{
											"attributes.display_value": map[string]interface{}{
												"query": queryVariant,
												"boost": 7.0,
											},
										},
									},
								},
							},
						},
						"score_mode": "max",
					},
				})

				// Для специфичных полей (модель, марка)
				if isCarBrandOrModel(queryLower) {
					exactMatchQueries = append(exactMatchQueries, map[string]interface{}{
						"term": map[string]interface{}{
							"model_lowercase": map[string]interface{}{
								"value": strings.ToLower(queryVariant),
								"boost": 10.0,
							},
						},
					})

					exactMatchQueries = append(exactMatchQueries, map[string]interface{}{
						"term": map[string]interface{}{
							"make_lowercase": map[string]interface{}{
								"value": strings.ToLower(queryVariant),
								"boost": 10.0,
							},
						},
					})
				}
			}

			// Используем must с bool->should для точного поиска
			must = append(must, map[string]interface{}{
				"bool": map[string]interface{}{
					"should":               exactMatchQueries,
					"minimum_should_match": 1,
				},
			})

		} else {
			// НЕЧЕТКИЙ ПОИСК для длинных запросов
			logger.Info().Msg("Используем нечеткий поиск для длинного запроса")

			// Основные поля для поиска с умеренными весами
			searchFields := []string{
				"title^5",
				"description^2",
				"title.sr^4",
				"title.ru^4",
				"title.en^4",
				"description.sr^1.5",
				"description.ru^1.5",
				"description.en^1.5",
			}

			// Multi-match запрос для всех вариантов
			for _, queryVariant := range queryVariants {
				should = append(should, map[string]interface{}{
					"multi_match": map[string]interface{}{
						"query":          queryVariant,
						"fields":         searchFields,
						"type":           "best_fields",
						"fuzziness":      "AUTO",
						"prefix_length":  2, // Не допускаем ошибки в первых 2 символах
						"max_expansions": 50,
					},
				})

				// Добавляем поиск по атрибутам
				should = append(should, map[string]interface{}{
					"nested": map[string]interface{}{
						"path": "attributes",
						"query": map[string]interface{}{
							"multi_match": map[string]interface{}{
								"query": queryVariant,
								"fields": []string{
									"attributes.text_value^3",
									"attributes.display_value^2",
								},
								"fuzziness": "AUTO",
							},
						},
						"score_mode": "max",
					},
				})
			}

			// Устанавливаем минимальное совпадение
			boolMap["minimum_should_match"] = "30%"
		}

		boolMap["must"] = must
		boolMap["should"] = should
	}

	// Фильтры (остаются без изменений)
	r.addFilters(query, params)

	// Сортировка
	r.addSorting(query, params)

	// Если нет условий в must, добавляем match_all для получения всех результатов
	boolQuery := query["query"].(map[string]interface{})["bool"].(map[string]interface{})
	must := boolQuery["must"].([]interface{})
	if len(must) == 0 {
		boolQuery["must"] = append(must, map[string]interface{}{
			"match_all": map[string]interface{}{},
		})
	}

	return query
}

// addFilters добавляет фильтры к запросу
func (r *Repository) addFilters(query map[string]interface{}, params *search.SearchParams) {
	filter := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]interface{})

	// Обработка категорий - поддержка как единичной категории, так и массива
	if len(params.CategoryIDs) > 0 {
		logger.Info().Ints("category_ids", params.CategoryIDs).Msg("Applying category filter")
		// Если есть массив категорий, создаем фильтр для всех категорий
		shouldClauses := make([]map[string]interface{}, 0)
		for _, catID := range params.CategoryIDs {
			shouldClauses = append(shouldClauses, map[string]interface{}{
				"term": map[string]interface{}{
					"category_id": catID,
				},
			})
			shouldClauses = append(shouldClauses, map[string]interface{}{
				"term": map[string]interface{}{
					"category_path_ids": catID,
				},
			})
		}

		categoryFilter := map[string]interface{}{
			"bool": map[string]interface{}{
				"should":               shouldClauses,
				"minimum_should_match": 1,
			},
		}

		filter = append(filter, categoryFilter)
	} else if params.CategoryID != nil && *params.CategoryID > 0 {
		// Если нет массива, используем единичную категорию (для обратной совместимости)
		filter = append(filter, map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []map[string]interface{}{
					{
						"term": map[string]interface{}{
							"category_id": *params.CategoryID,
						},
					},
					{
						"term": map[string]interface{}{
							"category_path_ids": *params.CategoryID,
						},
					},
				},
				"minimum_should_match": 1,
			},
		})
	}

	// Цена
	if params.PriceMin != nil && *params.PriceMin > 0 {
		filter = append(filter, map[string]interface{}{
			"range": map[string]interface{}{
				"price": map[string]interface{}{
					"gte": *params.PriceMin,
				},
			},
		})
	}

	if params.PriceMax != nil && *params.PriceMax > 0 {
		filter = append(filter, map[string]interface{}{
			"range": map[string]interface{}{
				"price": map[string]interface{}{
					"lte": *params.PriceMax,
				},
			},
		})
	}

	// Город и страна
	if params.City != "" {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"city.keyword": params.City,
			},
		})
	}

	if params.Country != "" {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"country.keyword": params.Country,
			},
		})
	}

	// Витрина
	if params.StorefrontID != nil && *params.StorefrontID > 0 {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"storefront_id": *params.StorefrontID,
			},
		})
	}

	// Состояние
	if params.Condition != "" {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"condition": params.Condition,
			},
		})
	}

	// Фильтр витрин B2C
	if params.StorefrontFilter == "exclude_b2c" {
		filter = append(filter, map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []map[string]interface{}{
					{
						"bool": map[string]interface{}{
							"must_not": map[string]interface{}{
								"exists": map[string]interface{}{
									"field": "storefront_id",
								},
							},
						},
					},
					{
						"range": map[string]interface{}{
							"storefront_id": map[string]interface{}{
								"lte": 0,
							},
						},
					},
				},
				"minimum_should_match": 1,
			},
		})
	}

	query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = filter
}

// addSorting добавляет сортировку к запросу
func (r *Repository) addSorting(query map[string]interface{}, params *search.SearchParams) {
	// Сортировка по умолчанию - по релевантности (_score)
	sort := []interface{}{
		map[string]interface{}{
			"_score": map[string]interface{}{
				"order": "desc",
			},
		},
	}

	// Пользовательская сортировка
	if params.Sort != "" && params.Sort != "relevance" {
		sortField := params.Sort
		sortOrder := "desc"
		if params.SortDirection != "" {
			sortOrder = params.SortDirection
		}

		// Добавляем пользовательскую сортировку перед сортировкой по релевантности
		sort = append([]interface{}{
			map[string]interface{}{
				sortField: map[string]interface{}{
					"order": sortOrder,
				},
			},
		}, sort...)
	}

	// Добавляем сортировку по ID для стабильности результатов
	sort = append(sort, map[string]interface{}{
		"id": map[string]interface{}{
			"order": "desc",
		},
	})

	query["sort"] = sort
}

// isKnownBrand проверяет, является ли запрос известным брендом
func isKnownBrand(query string) bool {
	knownBrands := []string{
		"apple", "iphone", "ipad", "macbook", "imac",
		"samsung", "huawei", "xiaomi", "oppo", "vivo",
		"sony", "lg", "nokia", "motorola", "oneplus",
		"google", "pixel", "lenovo", "asus", "acer",
		"hp", "dell", "msi", "razer",
		"nike", "adidas", "puma", "reebok",
		"mercedes", "bmw", "audi", "volkswagen", "toyota",
		"ford", "chevrolet", "honda", "mazda", "nissan",
	}

	for _, brand := range knownBrands {
		if strings.Contains(query, brand) {
			return true
		}
	}
	return false
}

// isCarBrandOrModel проверяет, является ли запрос маркой или моделью автомобиля
func isCarBrandOrModel(query string) bool {
	// Список популярных марок и моделей
	carTerms := []string{
		"mercedes", "bmw", "audi", "volkswagen", "vw",
		"toyota", "honda", "mazda", "nissan", "mitsubishi",
		"ford", "chevrolet", "opel", "peugeot", "renault",
		"fiat", "alfa romeo", "citroen", "skoda", "seat",
		"hyundai", "kia", "volvo", "golf", "passat",
		"corolla", "camry", "accord", "civic", "focus",
	}

	for _, term := range carTerms {
		if strings.Contains(query, term) {
			return true
		}
	}
	return false
}
