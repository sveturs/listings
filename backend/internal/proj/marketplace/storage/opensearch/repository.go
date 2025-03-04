// backend/internal/proj/marketplace/storage/opensearch/repository.go
package opensearch

import (
	"backend/internal/domain/models"
	"backend/internal/domain/search"
	"backend/internal/storage"
	osClient "backend/internal/storage/opensearch" // Используем псевдоним для импорта
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Repository реализует интерфейс MarketplaceSearchRepository
type Repository struct {
	client    *osClient.OpenSearchClient
	indexName string
	storage   storage.Storage // Для доступа к оригинальным данным при индексации
}

// NewRepository создает новый репозиторий
func NewRepository(client *osClient.OpenSearchClient, indexName string, storage storage.Storage) *Repository {
	return &Repository{
		client:    client,
		indexName: indexName,
		storage:   storage,
	}
}

// PrepareIndex подготавливает индекс (создает, если не существует)
func (r *Repository) PrepareIndex(ctx context.Context) error {
	// Проверяем существование индекса
	exists, err := r.client.IndexExists(r.indexName)
	if err != nil {
		return fmt.Errorf("ошибка проверки индекса: %w", err)
	}

	log.Printf("Проверка индекса %s: существует=%v", r.indexName, exists)

	// Если индекс не существует, создаем его
	if !exists {
		log.Printf("Создание индекса %s...", r.indexName)
		if err := r.client.CreateIndex(r.indexName, osClient.ListingMapping); err != nil {
			return fmt.Errorf("ошибка создания индекса: %w", err)
		}
		log.Printf("Индекс %s успешно создан", r.indexName)

		// Получаем все объявления из БД
		allListings, _, err := r.storage.GetListings(ctx, map[string]string{}, 1000, 0)
		if err != nil {
			log.Printf("Ошибка получения объявлений: %v", err)
			return err
		}

		// Преобразуем в указатели
		listingPtrs := make([]*models.MarketplaceListing, len(allListings))
		for i := range allListings {
			listingPtrs[i] = &allListings[i]
		}

		// Индексируем все объявления
		if err := r.BulkIndexListings(ctx, listingPtrs); err != nil {
			log.Printf("Ошибка индексации объявлений: %v", err)
			return err
		}

		log.Printf("Успешно проиндексировано %d объявлений", len(allListings))
	}

	return nil
}

func (r *Repository) IndexListing(ctx context.Context, listing *models.MarketplaceListing) error {
    // Преобразуем объект модели в документ для индексации
    doc := r.listingToDoc(listing)

    // Логирование для отладки
    docJSON, _ := json.MarshalIndent(doc, "", "  ")
    log.Printf("Индексация объявления %d с данными: %s", listing.ID, string(docJSON))

    // Индексируем документ
    return r.client.IndexDocument(r.indexName, fmt.Sprintf("%d", listing.ID), doc)
}

// BulkIndexListings индексирует несколько объявлений
func (r *Repository) BulkIndexListings(ctx context.Context, listings []*models.MarketplaceListing) error {
	docs := make([]map[string]interface{}, 0, len(listings))

	for _, listing := range listings {
		doc := r.listingToDoc(listing)

		// Логирование ID при индексации
		log.Printf("Индексация объявления с ID: %d, категория: %d, название: %s",
			listing.ID, listing.CategoryID, listing.Title)

		// Гарантируем, что ID всегда установлен
		if listing.ID == 0 {
			log.Printf("ВНИМАНИЕ: Объявление с нулевым ID: %s (категория: %d)",
				listing.Title, listing.CategoryID)
		}

		doc["id"] = listing.ID // Используем явно указанный ID для индексации
		docs = append(docs, doc)
	}

	return r.client.BulkIndex(r.indexName, docs)
}

// DeleteListing удаляет объявление из индекса
func (r *Repository) DeleteListing(ctx context.Context, listingID string) error {
	return r.client.DeleteDocument(r.indexName, listingID)
}

// Метод для извлечения ID из документа OpenSearch
func (r *Repository) extractDocumentID(hit map[string]interface{}) (int, error) {
	// Сначала пытаемся получить ID из _id OpenSearch
	if idStr, ok := hit["_id"].(string); ok {
		if id, err := strconv.Atoi(idStr); err == nil {
			return id, nil
		}
	}

	// Затем попробуем из исходного документа
	if source, ok := hit["_source"].(map[string]interface{}); ok {
		if idFloat, ok := source["id"].(float64); ok {
			return int(idFloat), nil
		} else if idInt, ok := source["id"].(int); ok {
			return idInt, nil
		} else if idStr, ok := source["id"].(string); ok {
			if id, err := strconv.Atoi(idStr); err == nil {
				return id, nil
			}
		}
	}

	// Если не удалось получить ID, возможно нам нужно сделать дополнительный запрос
	return 0, fmt.Errorf("не удалось получить ID объявления из документа")
}

// SearchListings выполняет поиск объявлений
func (r *Repository) SearchListings(ctx context.Context, params *search.SearchParams) (*search.SearchResult, error) {
    // Строим запрос к OpenSearch
    query := r.buildSearchQuery(params)
    
    // Дополнительное логирование
    queryJSON, _ := json.MarshalIndent(query, "", "  ")
    log.Printf("Поисковый запрос: %s", string(queryJSON))


	// Выполняем поиск
	responseBytes, err := r.client.Search(r.indexName, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения поиска: %w", err)
	}

	// Разбираем ответ
	var searchResponse map[string]interface{}
	if err := json.Unmarshal(responseBytes, &searchResponse); err != nil {
		return nil, fmt.Errorf("ошибка разбора ответа: %w", err)
	}

	// Извлекаем результаты
	result, err := r.parseSearchResponse(searchResponse, params.Language)
	if err != nil {
		return nil, fmt.Errorf("ошибка обработки результатов: %w", err)
	}

	return result, nil
}

// SuggestListings предлагает автодополнение для поиска
func (r *Repository) SuggestListings(ctx context.Context, prefix string, size int) ([]string, error) {
    if prefix == "" {
        return []string{}, nil
    }

    // Журналирование для отладки
    log.Printf("Запрос автодополнения для: '%s', размер: %d", prefix, size)

    // Создаем запрос для поиска с префиксом
    query := map[string]interface{}{
        "size":    size,
        "_source": []string{"title"},
        "query": map[string]interface{}{
            "bool": map[string]interface{}{
                "should": []map[string]interface{}{
                    {
                        "match_phrase_prefix": map[string]interface{}{
                            "title": map[string]interface{}{
                                "query":          prefix,
                                "max_expansions": 10,
                                "slop":           2,
                            },
                        },
                    },
                    {
                        "match_phrase_prefix": map[string]interface{}{
                            "title_variations": map[string]interface{}{
                                "query":          prefix,
                                "max_expansions": 10,
                            },
                        },
                    },
                    {
                        "fuzzy": map[string]interface{}{
                            "title": map[string]interface{}{
                                "value":     prefix,
                                "fuzziness": "AUTO",
                            },
                        },
                    },
                },
                "minimum_should_match": 1,
            },
        },
    }

    // Дополнительно добавляем suggest API для классического автодополнения
    query["suggest"] = map[string]interface{}{
        "title_suggest": map[string]interface{}{
            "prefix": prefix,
            "completion": map[string]interface{}{
                "field": "title_suggest",
                "size":  size,
            },
        },
    }

    // Выполняем поиск
    responseBytes, err := r.client.Search(r.indexName, query)
    if err != nil {
        return nil, fmt.Errorf("ошибка выполнения поиска: %w", err)
    }

    // Парсим JSON ответ
    var searchResponse map[string]interface{}
    if err := json.Unmarshal(responseBytes, &searchResponse); err != nil {
        return nil, fmt.Errorf("ошибка разбора ответа: %w", err)
    }

    // Извлекаем результаты из hits
    suggestions := make([]string, 0, size)
    if hits, ok := searchResponse["hits"].(map[string]interface{}); ok {
        if hitsArray, ok := hits["hits"].([]interface{}); ok {
            for _, hit := range hitsArray {
                if hitObj, ok := hit.(map[string]interface{}); ok {
                    if source, ok := hitObj["_source"].(map[string]interface{}); ok {
                        if title, ok := source["title"].(string); ok {
                            suggestions = append(suggestions, title)
                        }
                    }
                }
            }
        }
    }

    // Также проверяем результаты из suggest API
    if suggest, ok := searchResponse["suggest"].(map[string]interface{}); ok {
        if titleSuggest, ok := suggest["title_suggest"].([]interface{}); ok && len(titleSuggest) > 0 {
            if suggItem, ok := titleSuggest[0].(map[string]interface{}); ok {
                if options, ok := suggItem["options"].([]interface{}); ok {
                    for _, option := range options {
                        if optObj, ok := option.(map[string]interface{}); ok {
                            if text, ok := optObj["text"].(string); ok {
                                if !contains(suggestions, text) {
                                    suggestions = append(suggestions, text)
                                }
                            }
                        }
                    }
                }
            }
        }
    }

    // Логируем результаты для отладки
    log.Printf("Найдено %d подсказок для '%s': %v", len(suggestions), prefix, suggestions)

    return suggestions, nil
}

// Вспомогательная функция для проверки наличия элемента в слайсе
func contains(arr []string, str string) bool {
    for _, a := range arr {
        if a == str {
            return true
        }
    }
    return false
}


// ReindexAll переиндексирует все объявления
func (r *Repository) ReindexAll(ctx context.Context) error {
	// Удаляем индекс, если он существует
	exists, err := r.client.IndexExists(r.indexName)
	if err != nil {
		return fmt.Errorf("ошибка проверки индекса: %w", err)
	}

	if exists {
		log.Printf("Удаляем существующий индекс %s", r.indexName)
		if err := r.client.DeleteIndex(r.indexName); err != nil {
			return fmt.Errorf("ошибка удаления индекса: %w", err)
		}
		// Дадим кластеру время обработать удаление
		time.Sleep(1 * time.Second)
	}

	// Создаем индекс заново
	log.Printf("Создаем индекс %s заново", r.indexName)
	if err := r.client.CreateIndex(r.indexName, osClient.ListingMapping); err != nil {
		return fmt.Errorf("ошибка создания индекса: %w", err)
	}

	// Получаем все объявления из БД с пагинацией, чтобы обрабатывать большие наборы данных
	const batchSize = 100
	offset := 0
	totalIndexed := 0

	for {
		log.Printf("Получение пакета объявлений (размер: %d, смещение: %d)", batchSize, offset)
		listings, total, err := r.storage.GetListings(ctx, map[string]string{}, batchSize, offset)
		if err != nil {
			return fmt.Errorf("ошибка получения объявлений: %w", err)
		}

		if len(listings) == 0 {
			break // Больше нет объявлений
		}

		log.Printf("Получено %d объявлений из %d всего (пакет %d)", len(listings), total, offset/batchSize+1)

		// Логируем ID каждого объявления для отладки
		for i, listing := range listings {
			log.Printf("Объявление %d/%d: ID=%d, Категория=%d, Название=%s",
				i+1, len(listings), listing.ID, listing.CategoryID, listing.Title)
		}

		// Преобразуем в указатели для BulkIndexListings
		listingPtrs := make([]*models.MarketplaceListing, len(listings))
		for i := range listings {
			listingPtrs[i] = &listings[i]
		}

		// Индексируем пакет объявлений
		if err := r.BulkIndexListings(ctx, listingPtrs); err != nil {
			return fmt.Errorf("ошибка массовой индексации (пакет %d): %w", offset/batchSize+1, err)
		}

		totalIndexed += len(listings)
		offset += batchSize

		// Если получили меньше объявлений, чем размер пакета, значит это последний пакет
		if len(listings) < batchSize {
			break
		}
	}

	log.Printf("Успешно проиндексировано %d объявлений", totalIndexed)
	return nil
}

// listingToDoc преобразует объект модели в документ для индексации
func (r *Repository) listingToDoc(listing *models.MarketplaceListing) map[string]interface{} {
	doc := map[string]interface{}{
		"id":                listing.ID,
		"title":             listing.Title,
		"description":       listing.Description,
		"title_suggest":     listing.Title,
		"title_variations":  []string{listing.Title, strings.ToLower(listing.Title)},
		"price":             listing.Price,
		"condition":         listing.Condition,
		"status":            listing.Status,
		"location":          listing.Location,
		"city":              listing.City,
		"country":           listing.Country,
		"views_count":       listing.ViewsCount,
		"created_at":        listing.CreatedAt.Format(time.RFC3339),
		"updated_at":        listing.UpdatedAt.Format(time.RFC3339),
		"show_on_map":       listing.ShowOnMap,
		"original_language": listing.OriginalLanguage,
		"category_id":       listing.CategoryID,
		"user_id":           listing.UserID,
		"translations":      listing.Translations,
	}

	// Добавляем координаты, если они есть
	if listing.Latitude != nil && listing.Longitude != nil {
		doc["coordinates"] = map[string]interface{}{
			"lat": *listing.Latitude,
			"lon": *listing.Longitude,
		}
	}

	// Добавляем storefront_id, если есть
	if listing.StorefrontID != nil {
		doc["storefront_id"] = *listing.StorefrontID
	}

	// Добавляем путь категорий, если есть
	// Добавляем путь категорий, если есть
	if listing.CategoryPathIds != nil && len(listing.CategoryPathIds) > 0 {
		doc["category_path_ids"] = listing.CategoryPathIds
	} else {
		// Если путь категорий не задан, нужно его создать
		// Для этого выполним дополнительный запрос к базе данных
		parentID := listing.CategoryID
		pathIDs := []int{parentID}

		for parentID > 0 {
			var cat models.MarketplaceCategory
			err := r.storage.QueryRow(context.Background(),
				"SELECT parent_id FROM marketplace_categories WHERE id = $1", parentID).
				Scan(&cat.ParentID)

			if err != nil || cat.ParentID == nil {
				break
			}

			parentID = *cat.ParentID
			pathIDs = append([]int{parentID}, pathIDs...)
		}

		doc["category_path_ids"] = pathIDs
		log.Printf("Сгенерирован путь категорий для объявления %d: %v", listing.ID, pathIDs)
	}

	// Добавляем информацию о категории, если есть
	if listing.Category != nil {
		doc["category"] = map[string]interface{}{
			"id":   listing.Category.ID,
			"name": listing.Category.Name,
			"slug": listing.Category.Slug,
		}
	}

	// Добавляем информацию о пользователе, если есть
	if listing.User != nil {
		doc["user"] = map[string]interface{}{
			"id":    listing.User.ID,
			"name":  listing.User.Name,
			"email": listing.User.Email,
		}
	}

	// Добавляем изображения, если есть
	if listing.Images != nil && len(listing.Images) > 0 {
		log.Printf("Найдено %d изображений для объявления %d", len(listing.Images), listing.ID)
		imagesDoc := make([]map[string]interface{}, 0, len(listing.Images))
		for _, img := range listing.Images {
			log.Printf("  Изображение: ID=%d, Путь=%s, IsMain=%v", img.ID, img.FilePath, img.IsMain)
			imagesDoc = append(imagesDoc, map[string]interface{}{
				"id":        img.ID,
				"file_path": img.FilePath,
				"is_main":   img.IsMain,
			})
		}
		doc["images"] = imagesDoc
	} else {
		log.Printf("Объявление %d не имеет изображений", listing.ID)
	}

	return doc
}

// buildSearchQuery создает поисковый запрос OpenSearch
func (r *Repository) buildSearchQuery(params *search.SearchParams) map[string]interface{} {
	log.Printf("Строим запрос: категория = %v, язык = %s, поисковый запрос = %s",
		params.CategoryID, params.Language, params.Query)
	query := map[string]interface{}{
		"from": (params.Page - 1) * params.Size,
		"size": params.Size,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must":   []interface{}{},
				"filter": []interface{}{},
			},
		},
	}

	mustClauses := []interface{}{}
	filterClauses := []interface{}{}

 	if params.Query != "" {
		log.Printf("Текстовый поиск по запросу: '%s'", params.Query)

		// Определяем поля для поиска с учетом языка
		searchFields := []string{"title^3", "description"}
		if params.Language != "" {
			// Добавляем языко-специфичные поля, если указан язык
			searchFields = append(
				searchFields,
				fmt.Sprintf("title.%s^4", params.Language),
				fmt.Sprintf("description.%s", params.Language),
				fmt.Sprintf("translations.%s.title^4", params.Language),
				fmt.Sprintf("translations.%s.description", params.Language),
			)
		}

		// Настройки нечеткого поиска
		minimumShouldMatch := "70%"
		if params.MinimumShouldMatch != "" {
			minimumShouldMatch = params.MinimumShouldMatch
		}

		fuzziness := "AUTO"
		if params.Fuzziness != "" {
			fuzziness = params.Fuzziness
		}

		// Создаем запрос для поиска
		queryObj := map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":                params.Query,
				"fields":               searchFields,
				"type":                 "best_fields",
				"fuzziness":            fuzziness,
				"operator":             "AND",
				"minimum_should_match": minimumShouldMatch,
			},
		}

		// Добавляем запрос в mustClauses
		mustClauses = append(mustClauses, queryObj)

		log.Printf("Добавлен поисковый запрос для полей: %v", searchFields)
	}

	// Добавляем фильтры по категории
	if params.CategoryID != nil {
		categoryID := *params.CategoryID
		log.Printf("Фильтрация по CategoryID: %d", categoryID)

		// Используем фильтр по category_path_ids
		filterClauses = append(filterClauses, map[string]interface{}{
			"term": map[string]interface{}{
				"category_path_ids": categoryID,
			},
		})
	}

	// Добавляем фильтры по цене
	if params.PriceMin != nil || params.PriceMax != nil {
		priceRange := map[string]interface{}{}
		if params.PriceMin != nil {
			priceRange["gte"] = *params.PriceMin
		}
		if params.PriceMax != nil {
			priceRange["lte"] = *params.PriceMax
		}

		filterClauses = append(filterClauses, map[string]interface{}{
			"range": map[string]interface{}{
				"price": priceRange,
			},
		})
	}

	// Добавляем фильтр по состоянию
	if params.Condition != "" {
		filterClauses = append(filterClauses, map[string]interface{}{
			"term": map[string]interface{}{
				"condition.keyword": params.Condition,
			},
		})
	}

	// Добавляем фильтры по городу
	if params.City != "" {
		filterClauses = append(filterClauses, map[string]interface{}{
			"term": map[string]interface{}{
				"city.keyword": params.City,
			},
		})
	}

	// Добавляем фильтры по стране
	if params.Country != "" {
		filterClauses = append(filterClauses, map[string]interface{}{
			"term": map[string]interface{}{
				"country.keyword": params.Country,
			},
		})
	}

	// Добавляем фильтр по ID витрины
	if params.StorefrontID != nil {
		filterClauses = append(filterClauses, map[string]interface{}{
			"term": map[string]interface{}{
				"storefront_id": *params.StorefrontID,
			},
		})
	}

	// Добавляем фильтр по статусу (только активные объявления)
	filterClauses = append(filterClauses, map[string]interface{}{
		"term": map[string]interface{}{
			"status": "active",
		},
	})

// Добавляем геопоиск, если указаны координаты
if params.Location != nil && params.Distance != "" {
    // Проверяем, что координаты имеют ненулевые значения
    if params.Location.Lat == 0 && params.Location.Lon == 0 {
        log.Printf("Игнорируем параметр distance (%s) из-за нулевых координат", params.Distance)
    } else {
        filterClauses = append(filterClauses, map[string]interface{}{
            "geo_distance": map[string]interface{}{
                "distance": params.Distance,
                "coordinates": map[string]interface{}{
                    "lat": params.Location.Lat,
                    "lon": params.Location.Lon,
                },
            },
        })
    }
}

	// Добавляем clauses в запрос
	if len(mustClauses) > 0 {
		log.Printf("Добавляем %d must клауз в запрос", len(mustClauses))
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = mustClauses

		// Добавим лог содержимого must-клауз
		mustClausesJSON, _ := json.Marshal(mustClauses)
		log.Printf("Must-клаузы: %s", string(mustClausesJSON))
	}

	if len(filterClauses) > 0 {
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = filterClauses
	}

	// Добавляем настройки сортировки
	sortOpt := []interface{}{}

	if params.Sort != "" {
		direction := "desc"
		if params.SortDirection == "asc" {
			direction = "asc"
		}

		// Проверяем специальные случаи сортировки
		sortField := params.Sort
		if sortField == "date_desc" {
			sortField = "created_at"
			direction = "desc" // Явно указываем desc
		} else if sortField == "date_asc" {
			sortField = "created_at"
			direction = "asc" // Явно указываем asc
		} else if sortField == "price_desc" {
			sortField = "price"
			direction = "desc" // Явно указываем desc
		} else if sortField == "price_asc" {
			sortField = "price"
			direction = "asc" // Явно указываем asc
		}

		log.Printf("Сортировка по полю %s в порядке %s", sortField, direction)

		// Особая обработка для геосортировки
		if sortField == "distance" && params.Location != nil {
			sortOpt = append(sortOpt, map[string]interface{}{
				"_geo_distance": map[string]interface{}{
					"coordinates": map[string]interface{}{
						"lat": params.Location.Lat,
						"lon": params.Location.Lon,
					},
					"order": direction,
					"unit":  "km",
				},
			})
		} else {
			// Стандартная сортировка
			sortOpt = append(sortOpt, map[string]interface{}{
				sortField: map[string]interface{}{
					"order": direction,
				},
			})
		}
	} else {
		// По умолчанию сортируем по дате создания по убыванию
		sortOpt = append(sortOpt, map[string]interface{}{
			"created_at": map[string]interface{}{
				"order": "desc", // Изменено с asc на desc
			},
		})
	}

	query["sort"] = sortOpt

	// Добавляем агрегации для фасетной фильтрации
	aggregations := map[string]interface{}{}

	// Стандартные агрегации по категориям, ценам и т.д.
	aggregations["categories"] = map[string]interface{}{
		"terms": map[string]interface{}{
			"field": "category_id",
			"size":  100,
		},
	}

	aggregations["conditions"] = map[string]interface{}{
		"terms": map[string]interface{}{
			"field": "condition.keyword",
			"size":  10,
		},
	}

	aggregations["cities"] = map[string]interface{}{
		"terms": map[string]interface{}{
			"field": "city.keyword",
			"size":  50,
		},
	}

	aggregations["countries"] = map[string]interface{}{
		"terms": map[string]interface{}{
			"field": "country.keyword",
			"size":  50,
		},
	}

	aggregations["price_ranges"] = map[string]interface{}{
		"range": map[string]interface{}{
			"field": "price",
			"ranges": []interface{}{
				map[string]interface{}{"to": 1000},
				map[string]interface{}{"from": 1000, "to": 5000},
				map[string]interface{}{"from": 5000, "to": 10000},
				map[string]interface{}{"from": 10000, "to": 50000},
				map[string]interface{}{"from": 50000},
			},
		},
	}

	// Добавляем запрос на подсказки (исправление опечаток)
	if params.Query != "" {
		query["suggest"] = map[string]interface{}{
			"text": params.Query,
			"simple_phrase": map[string]interface{}{
				"phrase": map[string]interface{}{
					"field":      "title",
					"size":       5,
					"gram_size":  3,
					"confidence": 2.0,
					"max_errors": 2,
					"direct_generator": []interface{}{
						map[string]interface{}{
							"field":           "title",
							"suggest_mode":    "always",
							"min_word_length": 3,
						},
					},
				},
			},
		}
	}

	// Запрашиваем только нужные агрегации
	requestedAggs := map[string]interface{}{}
	if len(params.Aggregations) > 0 {
		for _, agg := range params.Aggregations {
			if aggDef, ok := aggregations[agg]; ok {
				requestedAggs[agg] = aggDef
			}
		}
	} else {
		// Если не указаны, добавляем все
		requestedAggs = aggregations
	}

	if len(requestedAggs) > 0 {
		query["aggs"] = requestedAggs
	}

	// Для отладки выводим запрос в лог
	queryJSON, _ := json.MarshalIndent(query, "", "  ")
	log.Printf("Сформированный запрос: %s", queryJSON)

	return query
}

// parseSearchResponse обрабатывает ответ от OpenSearch
func (r *Repository) parseSearchResponse(response map[string]interface{}, language string) (*search.SearchResult, error) {
	result := &search.SearchResult{
		Listings:     make([]*models.MarketplaceListing, 0),
		Aggregations: make(map[string][]search.Bucket),
		Suggestions:  make([]string, 0),
	}

	// Получаем статистику
	if took, ok := response["took"].(float64); ok {
		result.Took = int64(took)
	}

	// Получаем результаты поиска
	if hits, ok := response["hits"].(map[string]interface{}); ok {
		// Общее количество результатов
		if total, ok := hits["total"].(map[string]interface{}); ok {
			if value, ok := total["value"].(float64); ok {
				result.Total = int(value)
			}
		}

		// Получаем документы
		// В части получения документов внутри parseSearchResponse
		if hitsArray, ok := hits["hits"].([]interface{}); ok {
			for _, hitI := range hitsArray {
				if hit, ok := hitI.(map[string]interface{}); ok {
					// Извлекаем источник документа
					if source, ok := hit["_source"].(map[string]interface{}); ok {
						// Получаем ID из поля _id в ответе OpenSearch
						if idStr, ok := hit["_id"].(string); ok {
							// Пытаемся преобразовать строку в число
							if id, err := strconv.Atoi(idStr); err == nil {
								// Обновляем ID в source для правильного преобразования
								source["id"] = float64(id)
							}
						}

						// Преобразуем документ в объект MarketplaceListing
						listing, err := r.docToListing(source, language)
						if err != nil {
							log.Printf("Ошибка преобразования документа: %v", err)
							continue
						}

						// Если ID всё еще равен 0, пытаемся восстановить его из базы данных
						if listing.ID == 0 {
							// Пытаемся найти по комбинации полей
							filters := map[string]string{
								"title": listing.Title,
							}
							if listing.CategoryID > 0 {
								filters["category_id"] = fmt.Sprintf("%d", listing.CategoryID)
							}

							// Логируем попытку восстановления
							log.Printf("Попытка восстановить ID для объявления: %+v", filters)

							// Здесь можно добавить код для поиска в базе данных, но это необязательно
						}

						result.Listings = append(result.Listings, listing)
					}
				}
			}
		}
	}

	// Получаем агрегации (фасеты)
	if aggs, ok := response["aggregations"].(map[string]interface{}); ok {
		for name, aggI := range aggs {
			if agg, ok := aggI.(map[string]interface{}); ok {
				buckets := make([]search.Bucket, 0)

				// Обработка обычных агрегаций terms
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

				// Обработка range агрегаций
				if bucketsArray, ok := agg["buckets"].([]interface{}); ok {
					for _, bucketI := range bucketsArray {
						if bucket, ok := bucketI.(map[string]interface{}); ok {
							var key string
							var count int

							from, fromOk := bucket["from"].(float64)
							to, toOk := bucket["to"].(float64)

							if fromOk && toOk {
								key = fmt.Sprintf("%v-%v", from, to)
							} else if fromOk {
								key = fmt.Sprintf("%v+", from)
							} else if toOk {
								key = fmt.Sprintf("0-%v", to)
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
			}
		}
	}

	// Получаем предложения (для исправления опечаток)
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

// docToListing преобразует документ из OpenSearch в модель
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
		if id, err := strconv.Atoi(idStr); err == nil {
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

	if city, ok := doc["city"].(string); ok {
		listing.City = city
	}

	if country, ok := doc["country"].(string); ok {
		listing.Country = country
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

	// Обрабатываем координаты
	if coordinates, ok := doc["coordinates"].(map[string]interface{}); ok {
		if lat, ok := coordinates["lat"].(float64); ok {
			listing.Latitude = &lat
		}

		if lon, ok := coordinates["lon"].(float64); ok {
			listing.Longitude = &lon
		}
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
	// Обрабатываем изображения
	if imagesArray, ok := doc["images"].([]interface{}); ok {
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

				if isMain, ok := img["is_main"].(bool); ok {
					image.IsMain = isMain
				}

				images = append(images, image)
			}
		}

		listing.Images = images
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
			}
		}
	}

	return listing, nil
}
