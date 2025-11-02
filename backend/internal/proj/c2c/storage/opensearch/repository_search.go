// backend/internal/proj/c2c/storage/opensearch/repository_search.go
package opensearch

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"backend/internal/domain/models"
	"backend/internal/domain/search"
	"backend/internal/logger"
)

func (r *Repository) SearchListings(ctx context.Context, params *search.SearchParams) (*search.SearchResult, error) {
	// КРИТИЧНО: Валидация и нормализация параметров ПЕРЕД любыми операциями
	// Это гарантирует что meta будет заполнен корректными значениями
	if params.Size <= 0 {
		logger.Warn().Int("requested_size", params.Size).Msg("Invalid size parameter, using default 20")
		params.Size = 20 // Default size
	}
	if params.Size > 100 {
		logger.Warn().Int("requested_size", params.Size).Msg("Size too large, limiting to 100")
		params.Size = 100 // Max size
	}
	if params.Page < 1 {
		logger.Warn().Int("requested_page", params.Page).Msg("Invalid page parameter, using default 1")
		params.Page = 1 // Pages start from 1
	}

	var query map[string]interface{}

	// ВАЖНО: проверяем наличие CustomQuery и используем его напрямую, если он задан
	if params.CustomQuery != nil {
		query = params.CustomQuery
		// Логируем, что используем специальный запрос
		queryJSON, _ := json.MarshalIndent(query, "", "  ")
		logger.Info().Msgf("Используем специальный запрос для поиска: %s", string(queryJSON))
	} else {
		// Используем улучшенную версию для более точного поиска
		// если есть поисковый запрос
		if params.Query != "" {
			logger.Info().Msgf("Используем улучшенный поиск для запроса: %s", params.Query)
			query = r.buildImprovedSearchQuery(ctx, params)
		} else {
			// Для фильтрации без текстового запроса используем старый метод
			query = r.buildSearchQuery(ctx, params)
		}
	}

	// Логируем финальный запрос для отладки (используем Debug для снижения шума)
	if logger.Debug().Enabled() {
		queryJSON, _ := json.MarshalIndent(query, "", "  ")
		logger.Debug().Msgf("Финальный запрос к OpenSearch:\n%s", string(queryJSON))
	}

	response, err := r.client.Search(ctx, r.indexName, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения поиска: %w", err)
	}

	var searchResponse map[string]interface{}
	if err := json.Unmarshal(response, &searchResponse); err != nil {
		return nil, fmt.Errorf("ошибка разбора ответа: %w", err)
	}

	// Закомментировано для снижения шума в логах
	// logger.Info().Msgf("OpenSearch ответил успешно. Анализируем результаты...")
	return r.parseSearchResponse(searchResponse, params.Language)
}

func (r *Repository) SuggestListings(ctx context.Context, prefix string, size int) ([]string, error) {
	if prefix == "" {
		return []string{}, nil
	}

	logger.Info().Msgf("Запрос автодополнения для: '%s', размер: %d", prefix, size)

	// Получаем варианты транслитерации для префикса
	prefixVariants := r.transliterator.TransliterateForSearch(prefix)
	logger.Info().
		Str("original_prefix", prefix).
		Strs("transliterated_variants", prefixVariants).
		Msg("Generated transliteration variants for suggestions")

	// Создаем комплексный запрос, который ищет как по обычным полям, так и по атрибутам
	shouldQueries := []map[string]interface{}{}

	// Добавляем запросы для всех вариантов транслитерации
	for _, prefixVariant := range prefixVariants {
		// Поиск по заголовку
		shouldQueries = append(shouldQueries, map[string]interface{}{
			"match_phrase_prefix": map[string]interface{}{
				"title": map[string]interface{}{
					"query":          prefixVariant,
					"max_expansions": 10,
				},
			},
		})

		// Поиск по полю model_lowercase (для автомобилей)
		shouldQueries = append(shouldQueries, map[string]interface{}{
			"match_phrase_prefix": map[string]interface{}{
				"model_lowercase": map[string]interface{}{
					"query":          strings.ToLower(prefixVariant),
					"max_expansions": 10,
				},
			},
		})

		// Поиск по полю make_lowercase (для автомобилей)
		shouldQueries = append(shouldQueries, map[string]interface{}{
			"match_phrase_prefix": map[string]interface{}{
				"make_lowercase": map[string]interface{}{
					"query":          strings.ToLower(prefixVariant),
					"max_expansions": 10,
				},
			},
		})

		// Поиск по атрибутам (nested query)
		shouldQueries = append(shouldQueries, map[string]interface{}{
			"nested": map[string]interface{}{
				"path": "attributes",
				"query": map[string]interface{}{
					"bool": map[string]interface{}{
						"should": []map[string]interface{}{
							// Поиск по текстовым значениям атрибутов
							{
								"match_phrase_prefix": map[string]interface{}{
									"attributes.text_value": map[string]interface{}{
										"query":          prefixVariant,
										"max_expansions": 10,
									},
								},
							},
							// Поиск по отображаемым значениям атрибутов
							{
								"match_phrase_prefix": map[string]interface{}{
									"attributes.display_value": map[string]interface{}{
										"query":          prefixVariant,
										"max_expansions": 10,
									},
								},
							},
						},
						// Приоритет для автомобильных атрибутов
						"boost": r.getBoostWeight("AutomotiveAttributePriority", 2.0),
					},
				},
			},
		})
	}

	// Создаем регулярное выражение для всех вариантов транслитерации
	regexPatterns := make([]string, len(prefixVariants))
	for i, variant := range prefixVariants {
		regexPatterns[i] = regexp.QuoteMeta(variant)
	}
	regexPattern := fmt.Sprintf(".*(%s).*", strings.Join(regexPatterns, "|"))

	// Создаем структуру запроса
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should":               shouldQueries,
				"minimum_should_match": 1,
			},
		},
		// Добавляем агрегации для извлечения уникальных значений
		"aggs": map[string]interface{}{
			"title_suggestions": map[string]interface{}{
				"terms": map[string]interface{}{
					"field":         "title.keyword",
					"size":          size,
					"min_doc_count": 1,
					"include":       regexPattern,
					"order":         map[string]string{"_count": "desc"},
				},
			},
			"make_suggestions": map[string]interface{}{
				"terms": map[string]interface{}{
					"field":         "make.keyword",
					"size":          size,
					"min_doc_count": 1,
					"include":       regexPattern,
					"order":         map[string]string{"_count": "desc"},
				},
			},
			"model_suggestions": map[string]interface{}{
				"terms": map[string]interface{}{
					"field":         "model.keyword",
					"size":          size,
					"min_doc_count": 1,
					"include":       regexPattern,
					"order":         map[string]string{"_count": "desc"},
				},
			},
			"nested_attr_suggestions": map[string]interface{}{
				"nested": map[string]interface{}{
					"path": "attributes",
				},
				"aggs": map[string]interface{}{
					"attribute_values": map[string]interface{}{
						"terms": map[string]interface{}{
							"field":         "attributes.text_value.keyword",
							"size":          size,
							"min_doc_count": 1,
							"include":       regexPattern,
							"order":         map[string]string{"_count": "desc"},
						},
					},
					"display_values": map[string]interface{}{
						"terms": map[string]interface{}{
							"field":         "attributes.display_value.keyword",
							"size":          size,
							"min_doc_count": 1,
							"include":       regexPattern,
							"order":         map[string]string{"_count": "desc"},
						},
					},
					// Специальные агрегации для моделей (авто)
					"model_values": map[string]interface{}{
						"filter": map[string]interface{}{
							"term": map[string]interface{}{
								"attributes.attribute_name": "model",
							},
						},
						"aggs": map[string]interface{}{
							"models": map[string]interface{}{
								"terms": map[string]interface{}{
									"field":         "attributes.text_value.keyword",
									"size":          size,
									"min_doc_count": 1,
									"include":       regexPattern,
									"order":         map[string]string{"_count": "desc"},
								},
							},
						},
					},
					// Специальные агрегации для марок (авто)
					"make_values": map[string]interface{}{
						"filter": map[string]interface{}{
							"term": map[string]interface{}{
								"attributes.attribute_name": "make",
							},
						},
						"aggs": map[string]interface{}{
							"makes": map[string]interface{}{
								"terms": map[string]interface{}{
									"field":         "attributes.text_value.keyword",
									"size":          size,
									"min_doc_count": 1,
									"include":       regexPattern,
									"order":         map[string]string{"_count": "desc"},
								},
							},
						},
					},
				},
			},
		},
	}

	// Добавляем запрос на автопродление, который уже есть в оригинальной функции
	query["suggest"] = map[string]interface{}{
		"title_suggest": map[string]interface{}{
			"prefix": prefix,
			"completion": map[string]interface{}{
				"field": "title_suggest",
				"size":  size,
			},
		},
	}

	responseBytes, err := r.client.Search(ctx, r.indexName, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения поиска для автопродления: %w", err)
	}

	var searchResponse map[string]interface{}
	if err := json.Unmarshal(responseBytes, &searchResponse); err != nil {
		return nil, fmt.Errorf("ошибка разбора ответа: %w", err)
	}

	// Создаем множество для хранения уникальных подсказок
	suggestionSet := make(map[string]bool)

	// Извлекаем подсказки из обычных результатов поиска
	if hits, ok := searchResponse["hits"].(map[string]interface{}); ok {
		if hitsArray, ok := hits["hits"].([]interface{}); ok {
			for _, hit := range hitsArray {
				if hitObj, ok := hit.(map[string]interface{}); ok {
					if source, ok := hitObj["_source"].(map[string]interface{}); ok {
						// Извлекаем заголовок
						if title, ok := source["title"].(string); ok && title != "" {
							suggestionSet[title] = true
						}

						// Извлекаем марку и модель
						if make, ok := source["make"].(string); ok && make != "" {
							suggestionSet[make] = true
						}
						if model, ok := source["model"].(string); ok && model != "" {
							suggestionSet[model] = true
						}
					}
				}
			}
		}
	}

	// Извлекаем подсказки из агрегаций
	if aggs, ok := searchResponse["aggregations"].(map[string]interface{}); ok {
		// Извлекаем подсказки из title_suggestions
		extractSuggestionsFromAgg(aggs, "title_suggestions", suggestionSet)

		// Извлекаем подсказки из make_suggestions
		extractSuggestionsFromAgg(aggs, "make_suggestions", suggestionSet)

		// Извлекаем подсказки из model_suggestions
		extractSuggestionsFromAgg(aggs, "model_suggestions", suggestionSet)

		// Извлекаем подсказки из nested_attr_suggestions
		if nestedAgg, ok := aggs["nested_attr_suggestions"].(map[string]interface{}); ok {
			// Извлекаем обычные значения атрибутов
			extractSuggestionsFromAgg(nestedAgg, "attribute_values", suggestionSet)
			extractSuggestionsFromAgg(nestedAgg, "display_values", suggestionSet)

			// Извлекаем значения моделей
			if modelValuesAgg, ok := nestedAgg["model_values"].(map[string]interface{}); ok {
				if modelsAgg, ok := modelValuesAgg["models"].(map[string]interface{}); ok {
					extractBucketsFromAgg(modelsAgg, suggestionSet)
				}
			}

			// Извлекаем значения марок
			if makeValuesAgg, ok := nestedAgg["make_values"].(map[string]interface{}); ok {
				if makesAgg, ok := makeValuesAgg["makes"].(map[string]interface{}); ok {
					extractBucketsFromAgg(makesAgg, suggestionSet)
				}
			}
		}
	}

	// Извлекаем подсказки из suggest
	if suggest, ok := searchResponse["suggest"].(map[string]interface{}); ok {
		if titleSuggest, ok := suggest["title_suggest"].([]interface{}); ok && len(titleSuggest) > 0 {
			if suggItem, ok := titleSuggest[0].(map[string]interface{}); ok {
				if options, ok := suggItem["options"].([]interface{}); ok {
					for _, option := range options {
						if optObj, ok := option.(map[string]interface{}); ok {
							if text, ok := optObj["text"].(string); ok && text != "" {
								suggestionSet[text] = true
							}
						}
					}
				}
			}
		}
	}

	// Конвертируем множество в срез
	suggestions := make([]string, 0, len(suggestionSet))
	for sugg := range suggestionSet {
		if strings.Contains(strings.ToLower(sugg), strings.ToLower(prefix)) {
			suggestions = append(suggestions, sugg)
		}
	}

	// Ограничиваем количество результатов
	if len(suggestions) > size {
		suggestions = suggestions[:size]
	}

	logger.Info().Msgf("Найдено %d подсказок для '%s': %v", len(suggestions), prefix, suggestions)
	return suggestions, nil
}

// Вспомогательная функция для извлечения подсказок из агрегации

func (r *Repository) buildSearchQuery(ctx context.Context, params *search.SearchParams) map[string]interface{} {
	logger.Info().Msgf("Строим запрос: категория = %v, язык = %s, поисковый запрос = %s",
		params.CategoryID, params.Language, params.Query)

	// Validate and set default for size to prevent OpenSearch "numHits must be > 0" error
	size := params.Size
	if size <= 0 {
		size = 10 // Default size
		logger.Warn().Int("requested_size", params.Size).Int("default_size", size).Msg("Invalid size parameter, using default")
	}

	// Validate page to prevent negative 'from' value
	page := params.Page
	if page < 1 {
		page = 1 // Pages start from 1
	}

	query := map[string]interface{}{
		"from":              (page - 1) * size,
		"size":              size,
		"track_total_hits":  true, // Required for OpenSearch to correctly handle queries
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must":   []interface{}{},
				"filter": []interface{}{},
				"should": []interface{}{},
			},
		},
	}

	// DEBUG: Log query structure
	logger.Info().
		Int("size", size).
		Int("page", page).
		Int("from", (page-1)*size).
		Bool("track_total_hits", true).
		Msg("buildSearchQuery: query structure created")

	if params.Status == "" {
		filter := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]interface{})
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"status": "active",
			},
		})
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = filter
	}

	// Фильтр по типу документа (listing или product)
	if params.DocumentType != "" {
		filter := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]interface{})
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"type": params.DocumentType, // Используем поле "type" как в индексе OpenSearch
			},
		})
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = filter
	}

	// Исключаем storefront products из marketplace search (они должны искаться только через storefront search)
	// Это гарантирует что обычные объявления и витринные товары не смешиваются
	filter := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]interface{})
	filter = append(filter, map[string]interface{}{
		"bool": map[string]interface{}{
			"should": []interface{}{
				map[string]interface{}{
					"bool": map[string]interface{}{
						"must_not": map[string]interface{}{
							"exists": map[string]interface{}{
								"field": "is_storefront",
							},
						},
					},
				},
				map[string]interface{}{
					"term": map[string]interface{}{
						"is_storefront": false,
					},
				},
			},
			"minimum_should_match": 1,
		},
	})
	query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = filter

	if params.Query != "" {
		logger.Info().Msgf("Текстовый поиск по запросу: '%s'", params.Query)

		searchFields := []string{
			"title^3", "description",
			"title.sr^4", "description.sr", "translations.sr.title^4", "translations.sr.description",
			"title.ru^4", "description.ru", "translations.ru.title^4", "translations.ru.description",
			"title.en^4", "description.en", "translations.en.title^4", "translations.en.description",
			"all_attributes_text^2",
			"make^5",
			"model^4",
			"color^3",
			"brand^4",
			"property_type^3",
			"body_type^3",
			// удалил rooms^3, т.к. оно числовое
			"cpu^3",
			"gpu^3",
			"memory^3",
			"ram^3",
			"storage_capacity^3",
			// удалил screen_size^3, т.к. оно числовое и не поддерживает fuzzy
			"attr_make_text^5",
			"attr_model_text^4",
			"attr_color_text^3",
			"attr_brand_text^4",
			"attr_year_text^3",
			"attr_rooms_text^3",
			"attr_property_type_text^3",
			"attr_body_type_text^3",
			"attr_cpu_text^3",
			"attr_gpu_text^3",
			"attr_memory_text^3",
			"attr_ram_text^3",
			"attr_storage_capacity_text^3",
			"attr_screen_size_text^3",
			"real_estate_attributes_text^3",
			"real_estate_attributes_combined^3",
			// удалил rooms^4, т.к. оно дублирует и оно числовое
			// удалил area^3, т.к. оно числовое
			// удалил floor^3, т.к. оно числовое
			// удалил total_floors^3, т.к. оно числовое
			"property_type^4",
			// удалил land_area^3, т.к. оно числовое
			// удалил year_built^3, т.к. оно числовое
			// удалил bathrooms^3, т.к. оно числовое
			"heating_type^3",
			"parking^3",
			"balcony^3",
			"furnished^3",
			"air_conditioning^3",
			"rooms_text^4",
			"property_type_text^4",
			"heating_type_text^3",
			"parking_text^3",
			"furnished_text^3",
			"car_keywords^5",
			"attributes.text_value^4",
			"attributes.display_value^4",
			"attributes.text_value.keyword^5",
			"make^6",
			"model^6",
			"make_lowercase^6",
			"model_lowercase^6",
		}

		languagePriority := "sr"
		if params.Language != "" {
			languagePriority = params.Language
		}

		switch languagePriority {
		case "sr":
			searchFields = append(searchFields,
				"title.sr^5",
				"translations.sr.title^5",
				"attr_*_sr^4",
			)
		case "ru":
			searchFields = append(searchFields,
				"title.ru^5",
				"translations.ru.title^5",
				"attr_*_ru^4",
			)
		case "en":
			searchFields = append(searchFields,
				"title.en^5",
				"translations.en.title^5",
				"attr_*_en^4",
			)
		}

		minimumShouldMatch := "30%"
		if params.MinimumShouldMatch != "" {
			minimumShouldMatch = params.MinimumShouldMatch
		}

		fuzziness := "AUTO"
		if params.Fuzziness != "" {
			fuzziness = params.Fuzziness
		}

		boolMap := query["query"].(map[string]interface{})["bool"].(map[string]interface{})
		should := boolMap["should"].([]interface{})

		// Получаем варианты транслитерации для поискового запроса
		queryVariants := r.transliterator.TransliterateForSearch(params.Query)
		logger.Info().
			Str("original_query", params.Query).
			Strs("transliterated_variants", queryVariants).
			Msg("Generated transliteration variants for search")

		// Расширяем запрос синонимами, если включен нечеткий поиск
		if params.UseSynonyms {
			// Попробуем расширить запрос синонимами через PostgreSQL
			if r.storage != nil {
				expandedQuery, err := r.storage.ExpandSearchQuery(ctx, params.Query, params.Language)
				if err == nil && expandedQuery != params.Query {
					logger.Info().Str("original", params.Query).Str("expanded", expandedQuery).Msg("Using expanded query with synonyms")
					// Добавляем расширенный запрос как дополнительные условия поиска
					expandedWords := strings.Fields(expandedQuery)
					for _, word := range expandedWords {
						if word != params.Query && !strings.Contains(params.Query, word) {
							should = append(should, map[string]interface{}{
								"multi_match": map[string]interface{}{
									"query":  word,
									"fields": searchFields,
									"type":   "best_fields",
									"boost":  r.getBoostWeight("SynonymBoost", 0.5), // Меньший вес для синонимов
								},
							})
						}
					}
				}
			}
		}

		// Основной поиск по заголовку с высоким приоритетом для всех вариантов транслитерации
		titleBoost := r.getBoostWeight("Title", 5.0)
		logger.Info().Msgf("Title boost weight: %.2f for query: %s", titleBoost, params.Query)
		for _, queryVariant := range queryVariants {
			should = append(should, map[string]interface{}{
				"match": map[string]interface{}{
					"title": map[string]interface{}{
						"query":     queryVariant,
						"boost":     titleBoost,
						"fuzziness": fuzziness,
					},
				},
			})
		}

		// Добавляем поиск по n-граммам для лучшего нечеткого соответствия
		if params.UseSynonyms {
			titleNgramBoost := r.getBoostWeight("TitleNgram", 2.0)
			logger.Info().Msgf("TitleNgram boost weight: %.2f", titleNgramBoost)
			for _, queryVariant := range queryVariants {
				should = append(should, map[string]interface{}{
					"match": map[string]interface{}{
						"title.ngram": map[string]interface{}{
							"query": queryVariant,
							"boost": titleNgramBoost,
						},
					},
				})
			}
		}

		// Поиск по описанию для всех вариантов транслитерации
		descriptionBoost := r.getBoostWeight("Description", 2.0)
		logger.Info().Msgf("Description boost weight: %.2f", descriptionBoost)
		for _, queryVariant := range queryVariants {
			should = append(should, map[string]interface{}{
				"match": map[string]interface{}{
					"description": map[string]interface{}{
						"query":     queryVariant,
						"boost":     descriptionBoost,
						"fuzziness": fuzziness,
					},
				},
			})
		}

		// Поиск по переводам для всех языков и вариантов транслитерации
		translationTitleBoost := r.getBoostWeight("TranslationTitle", 4.0)
		translationDescBoost := r.getBoostWeight("TranslationDesc", 1.5)
		logger.Info().Msgf("Translation boost weights - Title: %.2f, Description: %.2f", translationTitleBoost, translationDescBoost)

		for _, queryVariant := range queryVariants {
			// Поиск по сербским переводам
			should = append(should, map[string]interface{}{
				"match": map[string]interface{}{
					"translations.sr.title": map[string]interface{}{
						"query":     queryVariant,
						"boost":     translationTitleBoost,
						"fuzziness": fuzziness,
					},
				},
			})

			should = append(should, map[string]interface{}{
				"match": map[string]interface{}{
					"translations.sr.description": map[string]interface{}{
						"query":     queryVariant,
						"boost":     translationDescBoost,
						"fuzziness": fuzziness,
					},
				},
			})

			// Поиск по русским переводам
			should = append(should, map[string]interface{}{
				"match": map[string]interface{}{
					"translations.ru.title": map[string]interface{}{
						"query":     queryVariant,
						"boost":     translationTitleBoost,
						"fuzziness": fuzziness,
					},
				},
			})

			should = append(should, map[string]interface{}{
				"match": map[string]interface{}{
					"translations.ru.description": map[string]interface{}{
						"query":     queryVariant,
						"boost":     translationDescBoost,
						"fuzziness": fuzziness,
					},
				},
			})

			// Поиск по английским переводам
			should = append(should, map[string]interface{}{
				"match": map[string]interface{}{
					"translations.en.title": map[string]interface{}{
						"query":     queryVariant,
						"boost":     translationTitleBoost,
						"fuzziness": fuzziness,
					},
				},
			})

			should = append(should, map[string]interface{}{
				"match": map[string]interface{}{
					"translations.en.description": map[string]interface{}{
						"query":     queryVariant,
						"boost":     translationDescBoost,
						"fuzziness": fuzziness,
					},
				},
			})
		}

		// Добавляем специальную обработку для атрибутов в nested формате для всех вариантов транслитерации
		for _, queryVariant := range queryVariants {
			attrQuery := map[string]interface{}{
				"nested": map[string]interface{}{
					"path": "attributes",
					"query": map[string]interface{}{
						"bool": map[string]interface{}{
							"should": []map[string]interface{}{
								{
									"match": map[string]interface{}{
										"attributes.text_value": map[string]interface{}{
											"query":     queryVariant,
											"boost":     r.getBoostWeight("AttributeDisplayValue", 4.0),
											"fuzziness": "AUTO",
										},
									},
								},
								{
									"match": map[string]interface{}{
										"attributes.display_value": map[string]interface{}{
											"query":     queryVariant,
											"boost":     r.getBoostWeight("AttributeDisplayValue", 4.0),
											"fuzziness": "AUTO",
										},
									},
								},
								{
									"term": map[string]interface{}{
										"attributes.text_value.keyword": map[string]interface{}{
											"value": queryVariant,
											"boost": r.getBoostWeight("AttributeTextValueKeyword", 5.0),
										},
									},
								},
							},
						},
					},
					"score_mode": "max",
					"boost":      r.getBoostWeight("PerWordAllAttributes", 3.0),
				},
			}

			should = append(should, attrQuery)
		}

		// Специальный запрос для модели автомобиля для всех вариантов транслитерации
		for _, queryVariant := range queryVariants {
			modelQuery := map[string]interface{}{
				"nested": map[string]interface{}{
					"path": "attributes",
					"query": map[string]interface{}{
						"bool": map[string]interface{}{
							"must": []map[string]interface{}{
								{
									"term": map[string]interface{}{
										"attributes.attribute_name": "model",
									},
								},
								{
									"match": map[string]interface{}{
										"attributes.text_value": map[string]interface{}{
											"query":     queryVariant,
											"boost":     r.getBoostWeight("RealEstateAttributesCombined", 6.0),
											"fuzziness": "AUTO",
										},
									},
								},
							},
						},
					},
					"score_mode": "max",
					"boost":      r.getBoostWeight("RealEstateAttributesCombined", 6.0),
				},
			}

			should = append(should, modelQuery)
		}

		// Аналогичный запрос для марки автомобиля для всех вариантов транслитерации
		for _, queryVariant := range queryVariants {
			makeQuery := map[string]interface{}{
				"nested": map[string]interface{}{
					"path": "attributes",
					"query": map[string]interface{}{
						"bool": map[string]interface{}{
							"must": []map[string]interface{}{
								{
									"term": map[string]interface{}{
										"attributes.attribute_name": "make",
									},
								},
								{
									"match": map[string]interface{}{
										"attributes.text_value": map[string]interface{}{
											"query":     queryVariant,
											"boost":     r.getBoostWeight("RealEstateAttributesCombined", 6.0),
											"fuzziness": "AUTO",
										},
									},
								},
							},
						},
					},
					"score_mode": "max",
					"boost":      r.getBoostWeight("RealEstateAttributesCombined", 6.0),
				},
			}

			should = append(should, makeQuery)
		}

		realEstateKeywords := []string{
			"квартира", "комната", "комнат", "дом", "этаж",
			"площадь", "м2", "кв.м", "квм", "кв м",
			"студия", "однокомнатная", "двухкомнатная", "трехкомнатная",
			"однушка", "двушка", "трешка", "участок", "сотка",
			"гараж", "паркинг", "балкон", "лоджия", "терраса",
			"ремонт", "новостройка", "вторичка", "жилье", "недвижимость",
			"аренда", "съем", "снять", "купить", "продажа",
		}

		// Проверяем все варианты транслитерации на наличие ключевых слов недвижимости
		isRealEstateQuery := false
		for _, queryVariant := range queryVariants {
			for _, keyword := range realEstateKeywords {
				if strings.Contains(strings.ToLower(queryVariant), keyword) {
					isRealEstateQuery = true
					break
				}
			}
			if isRealEstateQuery {
				break
			}
		}

		if isRealEstateQuery {
			logger.Info().Msgf("Обнаружен запрос о недвижимости: '%s'", params.Query)

			// Добавляем поиск по недвижимости для всех вариантов транслитерации
			for _, queryVariant := range queryVariants {
				realEstateBoost := []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"real_estate_attributes_combined": map[string]interface{}{
								"query":     queryVariant,
								"boost":     r.getBoostWeight("RealEstateAttributesCombined", 5.0),
								"fuzziness": fuzziness,
							},
						},
					},
					{
						"match": map[string]interface{}{
							"property_type": map[string]interface{}{
								"query":     queryVariant,
								"boost":     r.getBoostWeight("PropertyType", 4.0),
								"fuzziness": fuzziness,
							},
						},
					},
					{
						"match": map[string]interface{}{
							"rooms_text": map[string]interface{}{
								"query":     queryVariant,
								"boost":     r.getBoostWeight("RoomsText", 4.0),
								"fuzziness": fuzziness,
							},
						},
					},
				}

				for _, q := range realEstateBoost {
					should = append(should, q)
				}
			}
		}

		boolMap["should"] = should
		boolMap["minimum_should_match"] = 1

		// Добавляем multi_match для всех вариантов транслитерации
		logger.Info().Msgf("Adding multi_match query with %d fields for query: %s", len(searchFields), params.Query)
		must := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]interface{})

		multiMatchShould := []map[string]interface{}{}
		for _, queryVariant := range queryVariants {
			logger.Info().Msgf("Multi-match variant: %s", queryVariant)
			multiMatchShould = append(multiMatchShould, map[string]interface{}{
				"multi_match": map[string]interface{}{
					"query":                queryVariant,
					"fields":               searchFields,
					"type":                 "best_fields",
					"operator":             "OR",
					"minimum_should_match": minimumShouldMatch,
					"fuzziness":            "AUTO",
				},
			})
		}

		if len(multiMatchShould) > 0 {
			must = append(must, map[string]interface{}{
				"bool": map[string]interface{}{
					"should":               multiMatchShould,
					"minimum_should_match": 1,
				},
			})
		}

		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = must

		// Обрабатываем слова для всех вариантов транслитерации
		processedWords := make(map[string]bool)
		for _, queryVariant := range queryVariants {
			words := strings.Fields(queryVariant)
			if len(words) > 1 {
				for _, word := range words {
					if len(word) < 2 || processedWords[word] {
						continue
					}
					processedWords[word] = true

					should = append(should, map[string]interface{}{
						"match": map[string]interface{}{
							"title": map[string]interface{}{
								"query":     word,
								"boost":     r.getBoostWeight("PerWordTitle", 2.0),
								"fuzziness": fuzziness,
							},
						},
					})

					should = append(should, map[string]interface{}{
						"match": map[string]interface{}{
							"description": map[string]interface{}{
								"query":     word,
								"boost":     r.getBoostWeight("PerWordDescription", 1.0),
								"fuzziness": fuzziness,
							},
						},
					})

					should = append(should, map[string]interface{}{
						"match": map[string]interface{}{
							"all_attributes_text": map[string]interface{}{
								"query":     word,
								"boost":     r.getBoostWeight("Description", 2.0),
								"fuzziness": fuzziness,
							},
						},
					})

					if isRealEstateQuery {
						should = append(should, map[string]interface{}{
							"match": map[string]interface{}{
								"real_estate_attributes_combined": map[string]interface{}{
									"query":     word,
									"boost":     r.getBoostWeight("PerWordRealEstateAttributes", 3.0),
									"fuzziness": fuzziness,
								},
							},
						})

						should = append(should, map[string]interface{}{
							"match": map[string]interface{}{
								"rooms_text": map[string]interface{}{
									"query":     word,
									"boost":     r.getBoostWeight("PerWordRoomsText", 2.5),
									"fuzziness": fuzziness,
								},
							},
						})
					}
				}
			}
		}

		boolMap["should"] = should

		// Добавляем дополнительные запросы для марки и модели для всех вариантов транслитерации
		shouldQueries := []map[string]interface{}{}
		for _, queryVariant := range queryVariants {
			shouldQueries = append(shouldQueries,
				map[string]interface{}{
					"match": map[string]interface{}{
						"make": map[string]interface{}{
							"query":     queryVariant,
							"boost":     r.getBoostWeight("Title", 5.0),
							"fuzziness": "AUTO",
						},
					},
				},
				map[string]interface{}{
					"match": map[string]interface{}{
						"make_lowercase": map[string]interface{}{
							"query":     strings.ToLower(queryVariant),
							"boost":     r.getBoostWeight("RealEstateAttributesCombined", 5.0),
							"fuzziness": "AUTO",
						},
					},
				},
				map[string]interface{}{
					"match": map[string]interface{}{
						"model": map[string]interface{}{
							"query":     queryVariant,
							"boost":     r.getBoostWeight("RoomsText", 4.0),
							"fuzziness": "AUTO",
						},
					},
				},
				map[string]interface{}{
					"match": map[string]interface{}{
						"model_lowercase": map[string]interface{}{
							"query":     strings.ToLower(queryVariant),
							"boost":     r.getBoostWeight("CarModel", 4.0),
							"fuzziness": "AUTO",
						},
					},
				})
		}

		// Добавляем остальные поля
		for _, queryVariant := range queryVariants {
			shouldQueries = append(shouldQueries,
				map[string]interface{}{
					"match": map[string]interface{}{
						"select_values": map[string]interface{}{
							"query":     queryVariant,
							"boost":     r.getBoostWeight("PerWordAllAttributes", 3.0),
							"fuzziness": "AUTO",
						},
					},
				},
				map[string]interface{}{
					"match": map[string]interface{}{
						"car_keywords": map[string]interface{}{
							"query":     queryVariant,
							"boost":     r.getBoostWeight("CarMake", 5.0),
							"fuzziness": "AUTO",
						},
					},
				})
		}

		should = boolMap["should"].([]interface{})
		for _, q := range shouldQueries {
			should = append(should, q)
		}
		boolMap["should"] = should
	}

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

		filter := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]interface{})
		filter = append(filter, categoryFilter)
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = filter
	} else if params.CategoryID != nil && *params.CategoryID > 0 {
		// Если нет массива, используем единичную категорию (для обратной совместимости)
		categoryFilter := map[string]interface{}{
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
		}

		filter := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]interface{})
		filter = append(filter, categoryFilter)
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = filter
	}

	if params.PriceMin != nil && *params.PriceMin > 0 {
		filter := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]interface{})
		filter = append(filter, map[string]interface{}{
			"range": map[string]interface{}{
				"price": map[string]interface{}{
					"gte": *params.PriceMin,
				},
			},
		})
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = filter
	}

	if params.PriceMax != nil && *params.PriceMax > 0 {
		filter := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]interface{})
		filter = append(filter, map[string]interface{}{
			"range": map[string]interface{}{
				"price": map[string]interface{}{
					"lte": *params.PriceMax,
				},
			},
		})
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = filter
	}

	if params.Condition != "" {
		filter := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]interface{})
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"condition": params.Condition,
			},
		})
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = filter
	}

	if params.City != "" {
		filter := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]interface{})
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"city.keyword": params.City,
			},
		})
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = filter
	}

	if params.Country != "" {
		filter := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]interface{})
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"country.keyword": params.Country,
			},
		})
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = filter
	}

	if params.StorefrontID != nil && *params.StorefrontID > 0 {
		filter := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]interface{})
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"storefront_id": *params.StorefrontID,
			},
		})
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = filter
	}

	// Обработка фильтра B2C объявлений
	switch params.StorefrontFilter {
	case "exclude_b2c", "":
		// По умолчанию исключаем B2C объявления (объявления с storefront_id)
		filter := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]interface{})
		filter = append(filter, map[string]interface{}{
			"bool": map[string]interface{}{
				"must_not": []map[string]interface{}{
					{
						"exists": map[string]interface{}{
							"field": "storefront_id",
						},
					},
				},
			},
		})
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = filter
		logger.Info().Msgf("Применен фильтр исключения B2C объявлений (storefront_filter=%s)", params.StorefrontFilter)
	case "include_b2c":
		// Включаем B2C объявления - не добавляем никаких фильтров
		logger.Info().Msgf("B2C объявления включены в поиск (storefront_filter=%s)", params.StorefrontFilter)
	}

	if params.Status != "" {
		filter := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]interface{})
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"status": params.Status,
			},
		})
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = filter
	}

	// Применяем геофильтр если указано расстояние
	// Если координаты не указаны, используем центр Белграда по умолчанию
	if params.Distance != "" {
		var lat, lon float64

		if params.Location != nil {
			lat = params.Location.Lat
			lon = params.Location.Lon
		}

		// Если координаты не переданы, используем центр Белграда
		if lat == 0 && lon == 0 {
			lat = 44.8176 // Широта центра Белграда
			lon = 20.4633 // Долгота центра Белграда
			logger.Info().Msg("Используем координаты Белграда по умолчанию для геофильтра")
		}

		// Форматируем distance - если это просто число, добавляем "km"
		distance := params.Distance
		if _, err := strconv.Atoi(distance); err == nil {
			distance += "km"
		}

		geoFilter := map[string]interface{}{
			"geo_distance": map[string]interface{}{
				"distance": distance,
				"coordinates": map[string]interface{}{
					"lat": lat,
					"lon": lon,
				},
			},
		}

		logger.Info().
			Str("distance", distance).
			Float64("lat", lat).
			Float64("lon", lon).
			Msg("Применяем геофильтр")

		filter := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]interface{})
		filter = append(filter, geoFilter)
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = filter
	}

	if len(params.AttributeFilters) > 0 {
		filter := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]interface{})

		realEstateAttrs := map[string]bool{
			"property_type": true, "rooms": true, "floor": true, "total_floors": true,
			"area": true, "land_area": true, "building_type": true,
			"has_balcony": true, "has_elevator": true, "has_parking": true,
		}

		for attrName, attrValue := range params.AttributeFilters {
			if attrValue == "" {
				continue
			}

			if realEstateAttrs[attrName] {
				switch {
				case strings.Contains(attrValue, ","):
					parts := strings.Split(attrValue, ",")
					if len(parts) == 2 {
						minVal, minErr := strconv.ParseFloat(parts[0], 64)
						maxVal, maxErr := strconv.ParseFloat(parts[1], 64)

						if minErr == nil && maxErr == nil {
							filter = append(filter, map[string]interface{}{
								"range": map[string]interface{}{
									attrName: map[string]interface{}{
										"gte": minVal,
										"lte": maxVal,
									},
								},
							})
							logger.Info().Msgf("Added range filter for real estate attribute %s: %f-%f",
								attrName, minVal, maxVal)
						}
					}
				case attrName == "property_type" || attrName == "building_type":
					filter = append(filter, map[string]interface{}{
						"term": map[string]interface{}{
							attrName: attrValue,
						},
					})
					logger.Info().Msgf("Added term filter for text real estate attribute %s = %s",
						attrName, attrValue)
				case attrValue == boolValueTrue || attrValue == "false":
					boolVal := attrValue == boolValueTrue
					filter = append(filter, map[string]interface{}{
						"term": map[string]interface{}{
							attrName: boolVal,
						},
					})
					logger.Info().Msgf("Added boolean filter for real estate attribute %s = %v",
						attrName, boolVal)
				default:
					if numVal, err := strconv.ParseFloat(attrValue, 64); err == nil {
						filter = append(filter, map[string]interface{}{
							"term": map[string]interface{}{
								attrName: numVal,
							},
						})
						logger.Info().Msgf("Added numeric filter for real estate attribute %s = %f",
							attrName, numVal)
					} else {
						filter = append(filter, map[string]interface{}{
							"match": map[string]interface{}{
								attrName + "_text": attrValue,
							},
						})
						logger.Info().Msgf("Added text filter for real estate attribute %s = %s",
							attrName, attrValue)
					}
				}
			} else {
				nestedQuery := map[string]interface{}{
					"nested": map[string]interface{}{
						"path": "attributes",
						"query": map[string]interface{}{
							"bool": map[string]interface{}{
								"must": []map[string]interface{}{
									{
										"term": map[string]interface{}{
											"attributes.attribute_name": attrName,
										},
									},
								},
							},
						},
					},
				}

				innerBool := nestedQuery["nested"].(map[string]interface{})["query"].(map[string]interface{})["bool"].(map[string]interface{})
				innerMust := innerBool["must"].([]map[string]interface{})

				if strings.Contains(attrValue, ",") {
					parts := strings.Split(attrValue, ",")
					if len(parts) == 2 {
						minVal, minErr := strconv.ParseFloat(parts[0], 64)
						maxVal, maxErr := strconv.ParseFloat(parts[1], 64)

						if minErr == nil && maxErr == nil {
							innerMust = append(innerMust, map[string]interface{}{
								"range": map[string]interface{}{
									"attributes.numeric_value": map[string]interface{}{
										"gte": minVal,
										"lte": maxVal,
									},
								},
							})
						}
					}
				} else if attrValue == boolValueTrue || attrValue == "false" {
					boolVal := attrValue == boolValueTrue
					innerMust = append(innerMust, map[string]interface{}{
						"term": map[string]interface{}{
							"attributes.boolean_value": boolVal,
						},
					})
				} else if numVal, err := strconv.ParseFloat(attrValue, 64); err == nil {
					innerMust = append(innerMust, map[string]interface{}{
						"term": map[string]interface{}{
							"attributes.numeric_value": numVal,
						},
					})
				} else {
					innerMust = append(innerMust, map[string]interface{}{
						"term": map[string]interface{}{
							"attributes.text_value.keyword": attrValue,
						},
					})
				}

				innerBool["must"] = innerMust
				filter = append(filter, nestedQuery)
			}
		}

		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = filter
	}

	if params.Sort != "" {
		logger.Info().Msgf("Применяем сортировку: %s, направление: %s", params.Sort, params.SortDirection)
		var sortField string
		var sortOrder string

		if params.SortDirection != "" {
			sortOrder = params.SortDirection
		} else {
			sortOrder = sortOrderDesc
		}

		switch params.Sort {
		case "relevance":
			logger.Info().Msg("Используем сортировку по релевантности (_score DESC)")
			// Для сортировки по релевантности используем _score
			query["sort"] = []interface{}{
				map[string]interface{}{
					"_score": map[string]interface{}{
						"order": "desc",
					},
				},
				map[string]interface{}{
					"created_at": map[string]interface{}{
						"order": "desc",
					},
				},
			}
			return query
		case "date":
			sortField = fieldNameCreatedAt
			// sortOrder уже установлен из params.SortDirection выше
		case "date_desc":
			sortField = fieldNameCreatedAt
			sortOrder = sortOrderDesc
		case "date_asc":
			sortField = fieldNameCreatedAt
			sortOrder = sortOrderAsc
		case "created":
			sortField = fieldNameCreatedAt
			// sortOrder уже установлен из params.SortDirection выше
		case "created_at":
			sortField = fieldNameCreatedAt
			// sortOrder уже установлен из params.SortDirection выше
		case "created_at_desc":
			sortField = fieldNameCreatedAt
			sortOrder = sortOrderDesc
		case "created_at_asc":
			sortField = fieldNameCreatedAt
			sortOrder = sortOrderAsc
		case fieldNamePrice:
			sortField = fieldNamePrice
			// sortOrder уже установлен из params.SortDirection выше
		case "price_desc":
			sortField = fieldNamePrice
			sortOrder = sortOrderDesc
		case "price_asc":
			sortField = fieldNamePrice
			sortOrder = sortOrderAsc
		// Сортировка для автомобилей
		case "year_desc":
			query["sort"] = []interface{}{
				map[string]interface{}{
					"attributes.numeric_value": map[string]interface{}{
						"order": "desc",
						"nested": map[string]interface{}{
							"path": "attributes",
							"filter": map[string]interface{}{
								"term": map[string]interface{}{
									"attributes.attribute_name": "year",
								},
							},
						},
					},
				},
			}
			return query
		case "year_asc":
			query["sort"] = []interface{}{
				map[string]interface{}{
					"attributes.numeric_value": map[string]interface{}{
						"order": "asc",
						"nested": map[string]interface{}{
							"path": "attributes",
							"filter": map[string]interface{}{
								"term": map[string]interface{}{
									"attributes.attribute_name": "year",
								},
							},
						},
					},
				},
			}
			return query
		case "mileage_asc":
			query["sort"] = []interface{}{
				map[string]interface{}{
					"attributes.numeric_value": map[string]interface{}{
						"order": "asc",
						"nested": map[string]interface{}{
							"path": "attributes",
							"filter": map[string]interface{}{
								"term": map[string]interface{}{
									"attributes.attribute_name": "mileage",
								},
							},
						},
					},
				},
			}
			return query
		case "mileage_desc":
			query["sort"] = []interface{}{
				map[string]interface{}{
					"attributes.numeric_value": map[string]interface{}{
						"order": "desc",
						"nested": map[string]interface{}{
							"path": "attributes",
							"filter": map[string]interface{}{
								"term": map[string]interface{}{
									"attributes.attribute_name": "mileage",
								},
							},
						},
					},
				},
			}
			return query
		case "price_year_ratio":
			// Сортировка по соотношению цена/год (лучшая цена за новизну)
			query["sort"] = []interface{}{
				map[string]interface{}{
					"_script": map[string]interface{}{
						"type": "number",
						"script": map[string]interface{}{
							"source": `
								double year = 2000;
								for (def attr : params._source.attributes) {
									if (attr.attribute_name == 'year' && attr.numeric_value != null) {
										year = attr.numeric_value;
										break;
									}
								}
								double price = doc.containsKey('price') && doc['price'].size() > 0 ? doc['price'].value : 100000;
								double age = 2025 - year;
								if (age < 1) age = 1;
								return price / (2025 - age);
							`,
							"lang": "painless",
						},
						"order": "asc",
					},
				},
			}
			return query
		case "rating_desc":
			logger.Info().Msgf("Применяем сортировку рейтинга по УБЫВАНИЮ")
			query["sort"] = []interface{}{
				map[string]interface{}{
					"_script": map[string]interface{}{
						"type": "number",
						"script": map[string]interface{}{
							"source": "doc.containsKey('average_rating') ? doc['average_rating'].value : 0",
						},
						"order": "desc",
					},
				},
				map[string]interface{}{
					"views_count": map[string]interface{}{
						"order": "desc",
					},
				},
				map[string]interface{}{
					"created_at": map[string]interface{}{
						"order": "desc",
					},
				},
			}
			return query
		case "rating_asc":
			logger.Info().Msgf("Применяем сортировку рейтинга по ВОЗРАСТАНИЮ")
			query["sort"] = []interface{}{
				map[string]interface{}{
					"_script": map[string]interface{}{
						"type": "number",
						"script": map[string]interface{}{
							"source": "doc.containsKey('average_rating') ? doc['average_rating'].value : 0",
						},
						"order": "asc",
					},
				},
				map[string]interface{}{
					"views_count": map[string]interface{}{
						"order": "desc",
					},
				},
				map[string]interface{}{
					"created_at": map[string]interface{}{
						"order": "desc",
					},
				},
			}
			return query
		default:
			parts := strings.Split(params.Sort, "_")
			if len(parts) >= 2 {
				sortField = parts[0]
				if parts[len(parts)-1] == "asc" || parts[len(parts)-1] == "desc" {
					sortOrder = parts[len(parts)-1]
				}
			} else {
				sortField = params.Sort
			}
		}

		logger.Info().Msgf("Применяем сортировку по полю %s в порядке %s", sortField, sortOrder)
		query["sort"] = []interface{}{
			map[string]interface{}{
				sortField: map[string]interface{}{
					"order": sortOrder,
				},
			},
		}
	} else {
		logger.Info().Msg("Сортировка не указана, используется сортировка по умолчанию")
		// Если сортировка не указана, используем сортировку по умолчанию
		if params.Query != "" {
			logger.Info().Msg("Есть поисковый запрос - сортируем по релевантности (_score DESC, created_at DESC)")
			// Если есть поисковый запрос, сортируем по релевантности
			query["sort"] = []interface{}{
				map[string]interface{}{
					"_score": map[string]interface{}{
						"order": "desc",
					},
				},
				map[string]interface{}{
					"created_at": map[string]interface{}{
						"order": "desc",
					},
				},
			}
		} else {
			logger.Info().Msg("Нет поискового запроса - сортируем только по дате (created_at DESC)")
			// Если нет поискового запроса, сортируем по дате
			query["sort"] = []interface{}{
				map[string]interface{}{
					"created_at": map[string]interface{}{
						"order": "desc",
					},
				},
			}
		}
	}

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

func (r *Repository) parseSearchResponse(response map[string]interface{}, language string) (*search.SearchResult, error) {
	result := &search.SearchResult{
		Listings:     make([]*models.MarketplaceListing, 0),
		Aggregations: make(map[string][]search.Bucket),
		Suggestions:  make([]string, 0),
	}

	if took, ok := response["took"].(float64); ok {
		result.Took = int64(took)
	}

	if hits, ok := response["hits"].(map[string]interface{}); ok {
		if total, ok := hits["total"].(map[string]interface{}); ok {
			if value, ok := total["value"].(float64); ok {
				result.Total = int(value)
			}
		}

		// Используем карту для дедупликации по ID
		seenIDs := make(map[string]bool)

		if hitsArray, ok := hits["hits"].([]interface{}); ok {
			for _, hitI := range hitsArray {
				if hit, ok := hitI.(map[string]interface{}); ok {
					if source, ok := hit["_source"].(map[string]interface{}); ok {
						if idStr, ok := hit["_id"].(string); ok {
							// Проверяем, не видели ли мы уже этот ID
							if seenIDs[idStr] {
								logger.Info().Msgf("Пропускаем дублированный результат с ID: %s", idStr)
								continue
							}
							seenIDs[idStr] = true

							// Обрабатываем ID товаров витрин (формат sp_XXX) и обычных товаров
							if strings.HasPrefix(idStr, "sp_") {
								// Для товаров витрин сохраняем ID как есть
								source["id"] = idStr
								// Также сохраняем числовой ID для совместимости
								if numID := strings.TrimPrefix(idStr, "sp_"); numID != "" {
									if id, err := strconv.Atoi(numID); err == nil {
										source["product_id"] = float64(id)
									}
								}
							} else if id, err := strconv.Atoi(idStr); err == nil {
								// Для обычных товаров парсим числовой ID
								source["id"] = float64(id)
							}
						}

						listing, err := r.docToListing(source, language)
						if err != nil {
							logger.Info().Msgf("Ошибка преобразования документа: %v", err)
							continue
						}
						if avgRating, ok := source["average_rating"].(float64); ok {
							listing.AverageRating = avgRating
						}

						if reviewCount, ok := source["review_count"].(float64); ok {
							listing.ReviewCount = int(reviewCount)
						}

						if listing.ID == 0 {
							filters := map[string]string{
								"title": listing.Title,
							}
							if listing.CategoryID > 0 {
								filters["category_id"] = fmt.Sprintf("%d", listing.CategoryID)
							}

							logger.Info().Msgf("Попытка восстановить ID для объявления: %+v", filters)
						}

						result.Listings = append(result.Listings, listing)
					}
				}
			}
		}
	}

	if aggs, ok := response["aggregations"].(map[string]interface{}); ok {
		for name, aggI := range aggs {
			if agg, ok := aggI.(map[string]interface{}); ok {
				buckets := make([]search.Bucket, 0)

				if bucketsArray, ok := agg["buckets"].([]interface{}); ok {
					for _, bucketI := range bucketsArray {
						if bucket, ok := bucketI.(map[string]interface{}); ok {
							var key string
							var count int

							if keyVal, ok := bucket["key"].(string); ok {
								key = keyVal
							} else if keyVal, ok := bucket["key"].(float64); ok {
								key = fmt.Sprintf("%v", keyVal)
							} else {
								continue
							}

							if docCount, ok := bucket["doc_count"].(float64); ok {
								count = int(docCount)
							}

							buckets = append(buckets, search.Bucket{
								Key:   key,
								Count: count,
							})
						}
					}

					result.Aggregations[name] = buckets
				}

				if bucketsArray, ok := agg["buckets"].([]interface{}); ok {
					for _, bucketI := range bucketsArray {
						if bucket, ok := bucketI.(map[string]interface{}); ok {
							var key string
							var count int

							from, fromOk := bucket["from"].(float64)
							to, toOk := bucket["to"].(float64)

							switch {
							case fromOk && toOk:
								key = fmt.Sprintf("%v-%v", from, to)
							case fromOk:
								key = fmt.Sprintf("%v+", from)
							case toOk:
								key = fmt.Sprintf("0-%v", to)
							default:
								continue
							}

							if docCount, ok := bucket["doc_count"].(float64); ok {
								count = int(docCount)
							}

							buckets = append(buckets, search.Bucket{
								Key:   key,
								Count: count,
							})
						}
					}

					result.Aggregations[name] = buckets
				}
			}
		}
	}

	if suggest, ok := response["suggest"].(map[string]interface{}); ok {
		if simplePhrases, ok := suggest["simple_phrase"].([]interface{}); ok && len(simplePhrases) > 0 {
			if phrase, ok := simplePhrases[0].(map[string]interface{}); ok {
				if options, ok := phrase["options"].([]interface{}); ok {
					for _, optionI := range options {
						if option, ok := optionI.(map[string]interface{}); ok {
							if text, ok := option["text"].(string); ok {
								result.Suggestions = append(result.Suggestions, text)
							}
						}
					}
				}
			}
		}
	}

	return result, nil
}

func (r *Repository) docToListing(doc map[string]interface{}, language string) (*models.MarketplaceListing, error) {
	listing := &models.MarketplaceListing{
		User:     &models.User{},
		Category: &models.MarketplaceCategory{},
	}

	// Извлекаем базовые поля
	if idFloat, ok := doc["id"].(float64); ok {
		listing.ID = int(idFloat)
	} else if idInt, ok := doc["id"].(int); ok {
		listing.ID = idInt
	} else if idStr, ok := doc["id"].(string); ok {
		// Обрабатываем ID товаров витрин (формат sp_XXX)
		if strings.HasPrefix(idStr, "sp_") {
			// Для товаров витрин используем product_id
			if productID, ok := doc["product_id"].(float64); ok {
				listing.ID = int(productID)
			}
		} else if id, err := strconv.Atoi(idStr); err == nil {
			listing.ID = id
		}
	} else {
		// Пытаемся извлечь ID из _id поля (в OpenSearch)
		if idVal, ok := doc["_id"].(string); ok {
			if id, err := strconv.Atoi(idVal); err == nil {
				listing.ID = id
			}
		}
	}

	if title, ok := doc["title"].(string); ok {
		listing.Title = title
	}

	if description, ok := doc["description"].(string); ok {
		listing.Description = description
	}

	if price, ok := doc["price"].(float64); ok {
		listing.Price = price
	}

	if condition, ok := doc["condition"].(string); ok {
		listing.Condition = condition
	}

	if status, ok := doc["status"].(string); ok {
		listing.Status = status
	}

	if location, ok := doc["location"].(string); ok {
		listing.Location = location
	}

	// Парсим новые поля адреса
	if addressCity, ok := doc["address_city"].(string); ok {
		listing.City = addressCity
	} else if city, ok := doc["city"].(string); ok {
		listing.City = city
	}

	if addressCountry, ok := doc["address_country"].(string); ok {
		listing.Country = addressCountry
	} else if country, ok := doc["country"].(string); ok {
		listing.Country = country
	}

	if addressMultilingual, ok := doc["address_multilingual"].(map[string]interface{}); ok {
		listing.AddressMultilingual = make(map[string]string)
		for key, value := range addressMultilingual {
			if strValue, ok := value.(string); ok {
				listing.AddressMultilingual[key] = strValue
			}
		}
	}

	if viewsCount, ok := doc["views_count"].(float64); ok {
		listing.ViewsCount = int(viewsCount)
	}

	if showOnMap, ok := doc["show_on_map"].(bool); ok {
		listing.ShowOnMap = showOnMap
	}

	if originalLanguage, ok := doc["original_language"].(string); ok {
		listing.OriginalLanguage = originalLanguage
	}

	// Обрабатываем даты
	if createdAtStr, ok := doc["created_at"].(string); ok {
		if createdAt, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
			listing.CreatedAt = createdAt
		}
	}

	if updatedAtStr, ok := doc["updated_at"].(string); ok {
		if updatedAt, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
			listing.UpdatedAt = updatedAt
		}
	}

	// Обрабатываем ссылочные поля
	if categoryID, ok := doc["category_id"].(float64); ok {
		listing.CategoryID = int(categoryID)
	}

	if userID, ok := doc["user_id"].(float64); ok {
		listing.UserID = int(userID)
		listing.User.ID = int(userID)
	}

	if storefrontID, ok := doc["storefront_id"].(float64); ok {
		id := int(storefrontID)
		listing.StorefrontID = &id
	}

	// Обрабатываем флаг is_storefront для различия типа товара
	if isStorefront, ok := doc["is_storefront"].(bool); ok {
		listing.IsStorefrontProduct = isStorefront
	}

	// Обрабатываем координаты (проверяем оба варианта: coordinates и location)
	if coordinates, ok := doc["coordinates"].(map[string]interface{}); ok {
		if lat, ok := coordinates["lat"].(float64); ok {
			listing.Latitude = &lat
		}

		if lon, ok := coordinates["lon"].(float64); ok {
			listing.Longitude = &lon
		}
	}

	// Также проверяем поле location (используется в storefront products с geo_point типом)
	if location, ok := doc["location"].(map[string]interface{}); ok {
		if lat, ok := location["lat"].(float64); ok {
			listing.Latitude = &lat
		}

		if lon, ok := location["lon"].(float64); ok {
			listing.Longitude = &lon
		}
	}

	// Обрабатываем остатки товара (для товаров витрин)
	if stockQuantity, ok := doc["stock_quantity"].(float64); ok {
		stockInt := int(stockQuantity)
		listing.StockQuantity = &stockInt
	}
	if stockStatus, ok := doc["stock_status"].(string); ok {
		listing.StockStatus = &stockStatus
	}

	// Обрабатываем информацию о категории
	if category, ok := doc["category"].(map[string]interface{}); ok {
		if id, ok := category["id"].(float64); ok {
			listing.Category.ID = int(id)
		}

		if name, ok := category["name"].(string); ok {
			listing.Category.Name = name
		}

		if slug, ok := category["slug"].(string); ok {
			listing.Category.Slug = slug
		}
	}

	// Обрабатываем информацию о пользователе
	if user, ok := doc["user"].(map[string]interface{}); ok {
		if id, ok := user["id"].(float64); ok {
			listing.User.ID = int(id)
		}

		if name, ok := user["name"].(string); ok {
			listing.User.Name = name
		}

		if email, ok := user["email"].(string); ok {
			listing.User.Email = email
		}
	}

	// Обрабатываем путь категорий
	if categoryPathIDs, ok := doc["category_path_ids"].([]interface{}); ok && len(categoryPathIDs) > 0 {
		pathIds := make([]int, len(categoryPathIDs))
		for i, v := range categoryPathIDs {
			if id, ok := v.(float64); ok {
				pathIds[i] = int(id)
			}
		}
		listing.CategoryPathIds = pathIds
	} else {
		// Если путь категорий не найден, хотя бы добавим текущую категорию
		listing.CategoryPathIds = []int{listing.CategoryID}
	}

	// Обрабатываем изображения из массива image_urls (новый формат)
	if imageURLsArray, ok := doc["image_urls"].([]interface{}); ok && len(imageURLsArray) > 0 {
		images := make([]models.MarketplaceImage, 0, len(imageURLsArray))

		// Получаем primary_image_url если есть
		var primaryImageURL string
		if primaryURL, ok := doc["primary_image_url"].(string); ok {
			primaryImageURL = primaryURL
		}

		for idx, urlI := range imageURLsArray {
			if url, ok := urlI.(string); ok {
				image := models.MarketplaceImage{
					ID:          idx + 1, // Генерируем ID для изображения
					PublicURL:   url,
					IsMain:      url == primaryImageURL, // Помечаем главное изображение
					StorageType: "minio",                // Предполагаем что это MinIO
				}
				images = append(images, image)
			}
		}

		listing.Images = images
	} else if imagesArray, ok := doc["images"].([]interface{}); ok {
		// Старый формат - массив объектов с полной информацией
		images := make([]models.MarketplaceImage, 0, len(imagesArray))

		for _, imgI := range imagesArray {
			if img, ok := imgI.(map[string]interface{}); ok {
				var image models.MarketplaceImage

				if id, ok := img["id"].(float64); ok {
					image.ID = int(id)
				}

				if filePath, ok := img["file_path"].(string); ok {
					image.FilePath = filePath
				}

				if fileName, ok := img["file_name"].(string); ok {
					image.FileName = fileName
				}

				if isMain, ok := img["is_main"].(bool); ok {
					image.IsMain = isMain
				}

				// Добавляем поля storage_type и public_url, если они есть
				if storageType, ok := img["storage_type"].(string); ok {
					image.StorageType = storageType
				} else if filePath, ok := img["file_path"].(string); ok && strings.Contains(filePath, "listings/") {
					// Если путь содержит "listings/", то это MinIO
					image.StorageType = "minio"
				}

				// Поддерживаем оба варианта: url и public_url (для обратной совместимости)
				if url, ok := img["url"].(string); ok {
					image.PublicURL = url
				} else if publicURL, ok := img["public_url"].(string); ok {
					image.PublicURL = publicURL
				} else if image.StorageType == "minio" && image.FilePath != "" {
					// Если это MinIO, но public_url не указан, формируем его
					image.PublicURL = "/listings/" + image.FilePath
				}

				images = append(images, image)
			}
		}

		listing.Images = images
	} else {
		// Если не нашли ни image_urls, ни images
		logger.Warn().
			Int("listing_id", listing.ID).
			Interface("doc_keys", getDocKeys(doc)).
			Msg("No image_urls or images found in document")
	}

	// Обрабатываем переводы
	if translations, ok := doc["translations"].(map[string]interface{}); ok {
		transMap := models.TranslationMap{}

		for lang, transI := range translations {
			if trans, ok := transI.(map[string]interface{}); ok {
				fieldTransMap := make(map[string]string)

				for field, valueI := range trans {
					if value, ok := valueI.(string); ok {
						fieldTransMap[field] = value
					}
				}

				if len(fieldTransMap) > 0 {
					transMap[lang] = fieldTransMap
				}
			}
		}

		if len(transMap) > 0 {
			listing.Translations = transMap
		}
	}

	// Применяем переводы, если указан язык
	if language != "" && language != listing.OriginalLanguage {
		if listing.Translations != nil {
			if langTranslations, ok := listing.Translations[language]; ok {
				if title, ok := langTranslations["title"]; ok && title != "" {
					listing.Title = title
				}

				if description, ok := langTranslations["description"]; ok && description != "" {
					listing.Description = description
				}

				// Применяем переводы адреса
				if address, ok := langTranslations["address"]; ok && address != "" {
					listing.Location = address
				}
				if city, ok := langTranslations["city"]; ok && city != "" {
					listing.City = city
				}
				if country, ok := langTranslations["country"]; ok && country != "" {
					listing.Country = country
				}
			}
		}
	}

	// Также заполняем AddressMultilingual для всех языков из translations
	if len(listing.Translations) > 0 {
		if listing.AddressMultilingual == nil {
			listing.AddressMultilingual = make(map[string]string)
		}
		for lang, langTranslations := range listing.Translations {
			if address, ok := langTranslations["address"]; ok && address != "" {
				listing.AddressMultilingual[lang] = address
			}
		}
	}

	// Обрабатываем атрибуты
	if attributes, ok := doc["attributes"].([]interface{}); ok {
		for _, attrI := range attributes {
			if attr, ok := attrI.(map[string]interface{}); ok {
				var attrValue models.ListingAttributeValue

				// Извлекаем общие поля
				if id, ok := attr["attribute_id"].(float64); ok {
					attrValue.AttributeID = int(id)
				}
				if name, ok := attr["attribute_name"].(string); ok {
					attrValue.AttributeName = name
				}
				if displayName, ok := attr["display_name"].(string); ok {
					attrValue.DisplayName = displayName
				}
				if attrType, ok := attr["attribute_type"].(string); ok {
					attrValue.AttributeType = attrType
				}
				if displayValue, ok := attr["display_value"].(string); ok {
					attrValue.DisplayValue = displayValue
				}

				// Извлекаем типизированные значения
				if textValue, ok := attr["text_value"].(string); ok {
					attrValue.TextValue = &textValue
				}
				if numValue, ok := attr["numeric_value"].(float64); ok {
					attrValue.NumericValue = &numValue
				}
				if boolValue, ok := attr["boolean_value"].(bool); ok {
					attrValue.BooleanValue = &boolValue
				}
				if jsonValue, ok := attr["json_value"].(string); ok {
					attrValue.JSONValue = json.RawMessage(jsonValue)
				}

				listing.Attributes = append(listing.Attributes, attrValue)
			}
		}
	}

	if hasDiscount, ok := doc["has_discount"].(bool); ok {
		listing.HasDiscount = hasDiscount
	}

	if oldPrice, ok := doc["old_price"].(float64); ok {
		listing.OldPrice = &oldPrice
	}

	// Обрабатываем метаданные и информацию о скидках
	if metadataRaw, ok := doc["metadata"].(map[string]interface{}); ok {
		listing.Metadata = metadataRaw

		// Проверяем наличие информации о скидке в метаданных
		if discount, ok := metadataRaw["discount"].(map[string]interface{}); ok {
			// Обязательно устанавливаем флаг скидки
			listing.HasDiscount = true

			// Если есть previous_price, устанавливаем его в поле OldPrice
			if prevPrice, ok := discount["previous_price"].(float64); ok {
				listing.OldPrice = &prevPrice
				logger.Info().Msgf("Найдена скидка для объявления %d: скидка %v%%, старая цена: %.2f",
					listing.ID, discount["discount_percent"], prevPrice)
			}
			// Добавляем процент скидки
			if discountPercent, ok := discount["discount_percent"].(float64); ok {
				percent := int(discountPercent)
				listing.DiscountPercentage = &percent
			}
		}
	}
	if avgRating, ok := doc["average_rating"].(float64); ok {
		listing.AverageRating = avgRating
	}

	if reviewCount, ok := doc["review_count"].(float64); ok {
		listing.ReviewCount = int(reviewCount)
	}
	// Убедимся, что HasDiscount установлен корректно на основе OldPrice
	if listing.OldPrice != nil && *listing.OldPrice > 0 && *listing.OldPrice > listing.Price {
		listing.HasDiscount = true
	}

	// Обрабатываем информацию о витрине
	if storefrontID, ok := doc["storefront_id"].(float64); ok {
		sfID := int(storefrontID)
		listing.StorefrontID = &sfID
	}

	if storefrontData, ok := doc["storefront"].(map[string]interface{}); ok {
		storefront := &models.Storefront{}

		if id, ok := storefrontData["id"].(float64); ok {
			storefront.ID = int(id)
		}
		if name, ok := storefrontData["name"].(string); ok {
			storefront.Name = name
		}
		if slug, ok := storefrontData["slug"].(string); ok {
			storefront.Slug = slug
		}
		if rating, ok := storefrontData["rating"].(float64); ok {
			storefront.Rating = rating
		}
		if isVerified, ok := storefrontData["is_verified"].(bool); ok {
			storefront.IsVerified = isVerified
		}

		listing.Storefront = storefront
	}
	return listing, nil
}

// FindSimilarListings находит похожие объявления используя more_like_this запрос

func (r *Repository) FindSimilarListings(ctx context.Context, text string, size int) ([]*SimilarListing, error) {
	// Создаем more_like_this запрос
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"more_like_this": map[string]interface{}{
				"fields": []string{
					"title^3",
					"description^2",
					"all_attributes_text",
					"car_keywords",
					"real_estate_attributes_combined",
				},
				"like":            text,
				"min_term_freq":   1,
				"min_doc_freq":    1,
				"max_query_terms": 25,
				"analyzer":        "standard",
			},
		},
		"size":    size,
		"_source": []string{"id", "category_id", "title"},
	}

	// Выполняем запрос
	response, err := r.client.Search(ctx, r.indexName, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения more_like_this запроса: %w", err)
	}

	// Парсим результаты
	results := make([]*SimilarListing, 0)

	// Парсим JSON ответ
	var responseMap map[string]interface{}
	if err := json.Unmarshal(response, &responseMap); err != nil {
		return nil, fmt.Errorf("ошибка парсинга ответа: %w", err)
	}

	if hits, ok := responseMap["hits"].(map[string]interface{}); ok {
		if hitsArray, ok := hits["hits"].([]interface{}); ok {
			for _, hitI := range hitsArray {
				if hit, ok := hitI.(map[string]interface{}); ok {
					if source, ok := hit["_source"].(map[string]interface{}); ok {
						listing := &SimilarListing{}

						// Извлекаем ID
						if idFloat, ok := source["id"].(float64); ok {
							listing.ID = int32(idFloat)
						}

						// Извлекаем CategoryID
						if catIDFloat, ok := source["category_id"].(float64); ok {
							listing.CategoryID = int32(catIDFloat)
						}

						// Извлекаем Title
						if title, ok := source["title"].(string); ok {
							listing.Title = title
						}

						// Извлекаем Score
						if score, ok := hit["_score"].(float64); ok {
							listing.Score = score
						}

						results = append(results, listing)
					}
				}
			}
		}
	}

	return results, nil
}

// getDocKeys - helper функция для получения списка ключей документа (для отладки)
