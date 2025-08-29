package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

const (
	marketplaceIndex = "marketplace"
	listingType      = "listing"
	productType      = "product"
)

// IndexInfo представляет информацию об индексе
type IndexInfo struct {
	IndexName      string                 `json:"index_name"`
	DocumentCount  int64                  `json:"document_count"`
	SizeInBytes    int64                  `json:"size_in_bytes"`
	SizeFormatted  string                 `json:"size_formatted"`
	CreatedAt      *time.Time             `json:"created_at,omitempty"`
	LastUpdated    *time.Time             `json:"last_updated,omitempty"`
	Mappings       map[string]interface{} `json:"mappings"`
	Settings       map[string]interface{} `json:"settings"`
	Health         string                 `json:"health"`
	Status         string                 `json:"status"`
	NumberOfShards int                    `json:"number_of_shards"`
	Aliases        []string               `json:"aliases"`
}

// IndexedDocument представляет проиндексированный документ
type IndexedDocument struct {
	ID               string                 `json:"id"`
	Type             string                 `json:"type"` // listing или product
	Title            string                 `json:"title"`
	CategoryID       int                    `json:"category_id"`
	CategoryName     string                 `json:"category_name"`
	UserID           int                    `json:"user_id"`
	StorefrontID     *int                   `json:"storefront_id,omitempty"`
	IndexedAt        time.Time              `json:"indexed_at"`
	LastModified     time.Time              `json:"last_modified"`
	Status           string                 `json:"status"`
	SearchableFields map[string]interface{} `json:"searchable_fields"`
}

// IndexStatistics представляет статистику индекса
type IndexStatistics struct {
	TotalDocuments      int64            `json:"total_documents"`
	ListingsCount       int64            `json:"listings_count"`
	ProductsCount       int64            `json:"products_count"`
	LastReindexed       *time.Time       `json:"last_reindexed,omitempty"`
	DocumentsByCategory map[string]int64 `json:"documents_by_category"`
	DocumentsByStatus   map[string]int64 `json:"documents_by_status"`
	IndexHealth         string           `json:"index_health"`
	SearchableFields    []string         `json:"searchable_fields"`
}

// GetIndexInfo возвращает информацию об индексе
func (s *Service) GetIndexInfo(ctx context.Context) (*IndexInfo, error) {
	if s.osClient == nil {
		return nil, fmt.Errorf("OpenSearch client not initialized")
	}

	indexName := marketplaceIndex // TODO: получать из конфига

	// Получаем статистику индекса
	statsPath := fmt.Sprintf("/%s/_stats", indexName)
	statsResp, err := s.osClient.Execute(ctx, "GET", statsPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get index stats: %w", err)
	}

	var statsData map[string]interface{}
	if err := json.Unmarshal(statsResp, &statsData); err != nil {
		return nil, fmt.Errorf("failed to parse stats response: %w", err)
	}

	// Получаем маппинги
	mappingPath := fmt.Sprintf("/%s/_mapping", indexName)
	mappingResp, err := s.osClient.Execute(ctx, "GET", mappingPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get index mapping: %w", err)
	}

	var mappingData map[string]interface{}
	if err := json.Unmarshal(mappingResp, &mappingData); err != nil {
		return nil, fmt.Errorf("failed to parse mapping response: %w", err)
	}

	// Получаем настройки индекса
	settingsPath := fmt.Sprintf("/%s/_settings", indexName)
	settingsResp, err := s.osClient.Execute(ctx, "GET", settingsPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get index settings: %w", err)
	}

	var settingsData map[string]interface{}
	if err := json.Unmarshal(settingsResp, &settingsData); err != nil {
		return nil, fmt.Errorf("failed to parse settings response: %w", err)
	}

	// Получаем здоровье кластера
	healthPath := fmt.Sprintf("/_cluster/health/%s", indexName)
	healthResp, err := s.osClient.Execute(ctx, "GET", healthPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster health: %w", err)
	}

	var healthData map[string]interface{}
	if err := json.Unmarshal(healthResp, &healthData); err != nil {
		return nil, fmt.Errorf("failed to parse health response: %w", err)
	}

	// Извлекаем данные
	info := &IndexInfo{
		IndexName: indexName,
		Health:    getStringValue(healthData, "status", "unknown"),
		Status:    getStringValue(healthData, "status", "unknown"),
	}

	// Парсим статистику
	if indices, ok := statsData["indices"].(map[string]interface{}); ok {
		if indexStats, ok := indices[indexName].(map[string]interface{}); ok {
			if total, ok := indexStats["total"].(map[string]interface{}); ok {
				// Документы
				if docs, ok := total["docs"].(map[string]interface{}); ok {
					if count, ok := docs["count"].(float64); ok {
						info.DocumentCount = int64(count)
					}
				}
				// Размер
				if store, ok := total["store"].(map[string]interface{}); ok {
					if sizeBytes, ok := store["size_in_bytes"].(float64); ok {
						info.SizeInBytes = int64(sizeBytes)
						info.SizeFormatted = formatBytes(info.SizeInBytes)
					}
				}
			}
			// Шарды
			if primaries, ok := indexStats["primaries"].(map[string]interface{}); ok {
				if shards, ok := primaries["shards"].([]interface{}); ok {
					info.NumberOfShards = len(shards)
				}
			}
		}
	}

	// Маппинги
	if indexMapping, ok := mappingData[indexName].(map[string]interface{}); ok {
		if mappings, ok := indexMapping["mappings"].(map[string]interface{}); ok {
			info.Mappings = mappings
		}
	}

	// Настройки
	if indexSettings, ok := settingsData[indexName].(map[string]interface{}); ok {
		if settings, ok := indexSettings["settings"].(map[string]interface{}); ok {
			info.Settings = settings

			// Дата создания
			if indexSection, ok := settings["index"].(map[string]interface{}); ok {
				if creationDate, ok := indexSection["creation_date"].(string); ok {
					if timestamp, err := parseTimestamp(creationDate); err == nil {
						info.CreatedAt = &timestamp
					}
				}
			}
		}
	}

	// Алиасы
	aliasesPath := fmt.Sprintf("/%s/_alias", indexName)
	aliasesResp, err := s.osClient.Execute(ctx, "GET", aliasesPath, nil)
	if err == nil {
		var aliasesData map[string]interface{}
		if err := json.Unmarshal(aliasesResp, &aliasesData); err == nil {
			if indexAliases, ok := aliasesData[indexName].(map[string]interface{}); ok {
				if aliases, ok := indexAliases["aliases"].(map[string]interface{}); ok {
					info.Aliases = make([]string, 0, len(aliases))
					for alias := range aliases {
						info.Aliases = append(info.Aliases, alias)
					}
				}
			}
		}
	}

	// Устанавливаем последнее обновление как текущее время
	now := time.Now()
	info.LastUpdated = &now

	return info, nil
}

// GetIndexStatistics возвращает статистику индекса
func (s *Service) GetIndexStatistics(ctx context.Context) (*IndexStatistics, error) {
	if s.osClient == nil {
		return nil, fmt.Errorf("OpenSearch client not initialized")
	}

	indexName := "marketplace"

	// Общее количество документов
	countPath := fmt.Sprintf("/%s/_count", indexName)
	countResp, err := s.osClient.Execute(ctx, "GET", countPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get document count: %w", err)
	}

	var countData map[string]interface{}
	if err := json.Unmarshal(countResp, &countData); err != nil {
		return nil, fmt.Errorf("failed to parse count response: %w", err)
	}

	stats := &IndexStatistics{
		DocumentsByCategory: make(map[string]int64),
		DocumentsByStatus:   make(map[string]int64),
		IndexHealth:         "green",
	}

	// Общее количество
	if count, ok := countData["count"].(float64); ok {
		stats.TotalDocuments = int64(count)
	}

	// Агрегация по типам документов (listing vs product)
	typeAggQuery := map[string]interface{}{
		"size": 0,
		"aggs": map[string]interface{}{
			"doc_types": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "type.keyword",
					"size":  10,
				},
			},
		},
	}

	typeAggResp, err := s.osClient.Search(ctx, indexName, typeAggQuery)
	if err == nil {
		var aggData map[string]interface{}
		if err := json.Unmarshal(typeAggResp, &aggData); err == nil {
			if aggs, ok := aggData["aggregations"].(map[string]interface{}); ok {
				if docTypes, ok := aggs["doc_types"].(map[string]interface{}); ok {
					if buckets, ok := docTypes["buckets"].([]interface{}); ok {
						for _, bucket := range buckets {
							if b, ok := bucket.(map[string]interface{}); ok {
								if key, ok := b["key"].(string); ok {
									if docCount, ok := b["doc_count"].(float64); ok {
										switch key {
										case listingType:
											stats.ListingsCount = int64(docCount)
										case productType:
											stats.ProductsCount = int64(docCount)
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// Агрегация по категориям
	categoryAggQuery := map[string]interface{}{
		"size": 0,
		"aggs": map[string]interface{}{
			"categories": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "category_name.keyword",
					"size":  50,
				},
			},
		},
	}

	categoryAggResp, err := s.osClient.Search(ctx, indexName, categoryAggQuery)
	if err == nil {
		var aggData map[string]interface{}
		if err := json.Unmarshal(categoryAggResp, &aggData); err == nil {
			if aggs, ok := aggData["aggregations"].(map[string]interface{}); ok {
				if categories, ok := aggs["categories"].(map[string]interface{}); ok {
					if buckets, ok := categories["buckets"].([]interface{}); ok {
						for _, bucket := range buckets {
							if b, ok := bucket.(map[string]interface{}); ok {
								if key, ok := b["key"].(string); ok {
									if docCount, ok := b["doc_count"].(float64); ok {
										stats.DocumentsByCategory[key] = int64(docCount)
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// Агрегация по статусам
	statusAggQuery := map[string]interface{}{
		"size": 0,
		"aggs": map[string]interface{}{
			"statuses": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "status.keyword",
					"size":  10,
				},
			},
		},
	}

	statusAggResp, err := s.osClient.Search(ctx, indexName, statusAggQuery)
	if err == nil {
		var aggData map[string]interface{}
		if err := json.Unmarshal(statusAggResp, &aggData); err == nil {
			if aggs, ok := aggData["aggregations"].(map[string]interface{}); ok {
				if statuses, ok := aggs["statuses"].(map[string]interface{}); ok {
					if buckets, ok := statuses["buckets"].([]interface{}); ok {
						for _, bucket := range buckets {
							if b, ok := bucket.(map[string]interface{}); ok {
								if key, ok := b["key"].(string); ok {
									if docCount, ok := b["doc_count"].(float64); ok {
										stats.DocumentsByStatus[key] = int64(docCount)
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// Получаем список доступных полей для поиска из маппинга
	mappingPath := fmt.Sprintf("/%s/_mapping", indexName)
	mappingResp, err := s.osClient.Execute(ctx, "GET", mappingPath, nil)
	if err == nil {
		var mappingData map[string]interface{}
		if err := json.Unmarshal(mappingResp, &mappingData); err == nil {
			if indexMapping, ok := mappingData[indexName].(map[string]interface{}); ok {
				if mappings, ok := indexMapping["mappings"].(map[string]interface{}); ok {
					if properties, ok := mappings["properties"].(map[string]interface{}); ok {
						stats.SearchableFields = make([]string, 0, len(properties))
						for field := range properties {
							stats.SearchableFields = append(stats.SearchableFields, field)
						}
					}
				}
			}
		}
	}

	return stats, nil
}

// SearchIndexedDocuments выполняет поиск по индексированным документам
func (s *Service) SearchIndexedDocuments(ctx context.Context, searchQuery string, docType string, limit int) ([]IndexedDocument, error) {
	if s.osClient == nil {
		return nil, fmt.Errorf("OpenSearch client not initialized")
	}

	indexName := "marketplace"

	if limit <= 0 {
		limit = 20
	}

	// Строим поисковый запрос
	query := map[string]interface{}{
		"size": limit,
		"sort": []map[string]interface{}{
			{"_score": map[string]string{"order": "desc"}},
			{"created_at": map[string]string{"order": "desc"}},
		},
	}

	// Добавляем фильтры
	var mustClauses []map[string]interface{}

	// Фильтр по типу документа (product = есть storefront_id, listing = нет storefront_id)
	switch docType {
	case productType:
		mustClauses = append(mustClauses, map[string]interface{}{
			"exists": map[string]interface{}{
				"field": "storefront_id",
			},
		})
		mustClauses = append(mustClauses, map[string]interface{}{
			"range": map[string]interface{}{
				"storefront_id": map[string]interface{}{
					"gt": 0,
				},
			},
		})
	case listingType:
		mustClauses = append(mustClauses, map[string]interface{}{
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
						"term": map[string]interface{}{
							"storefront_id": 0,
						},
					},
				},
			},
		})
	}

	// Поисковый запрос
	if searchQuery != "" {
		mustClauses = append(mustClauses, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  searchQuery,
				"fields": []string{"title^3", "description^2", "category.name", "user.name", "location", "city"},
				"type":   "best_fields",
			},
		})
	}

	if len(mustClauses) > 0 {
		query["query"] = map[string]interface{}{
			"bool": map[string]interface{}{
				"must": mustClauses,
			},
		}
	} else {
		query["query"] = map[string]interface{}{
			"match_all": map[string]interface{}{},
		}
	}

	// Выполняем поиск
	searchResp, err := s.osClient.Search(ctx, indexName, query)
	if err != nil {
		return nil, fmt.Errorf("failed to search documents: %w", err)
	}

	var searchData map[string]interface{}
	if err := json.Unmarshal(searchResp, &searchData); err != nil {
		return nil, fmt.Errorf("failed to parse search response: %w", err)
	}

	// Парсим результаты
	documents := make([]IndexedDocument, 0)
	if hits, ok := searchData["hits"].(map[string]interface{}); ok {
		if hitsArray, ok := hits["hits"].([]interface{}); ok {
			for _, hit := range hitsArray {
				if h, ok := hit.(map[string]interface{}); ok {
					doc := IndexedDocument{
						SearchableFields: make(map[string]interface{}),
					}

					// ID документа
					if id, ok := h["_id"].(string); ok {
						doc.ID = id
					}

					// Источник документа
					if source, ok := h["_source"].(map[string]interface{}); ok {
						// Определяем тип документа по наличию storefront_id
						// Если есть storefront_id - это товар, иначе - объявление
						if storefrontID, ok := source["storefront_id"].(float64); ok {
							if storefrontID > 0 {
								doc.Type = productType
								id := int(storefrontID)
								doc.StorefrontID = &id
							} else {
								doc.Type = listingType
							}
						} else {
							doc.Type = listingType
						}

						// Основные поля
						if title, ok := source["title"].(string); ok {
							doc.Title = title
						}
						if categoryID, ok := source["category_id"].(float64); ok {
							doc.CategoryID = int(categoryID)
						}
						// Пытаемся получить имя категории из разных полей
						if category, ok := source["category"].(map[string]interface{}); ok {
							if categoryName, ok := category["name"].(string); ok {
								doc.CategoryName = categoryName
							}
						}
						if userID, ok := source["user_id"].(float64); ok {
							doc.UserID = int(userID)
						}
						if status, ok := source["status"].(string); ok {
							doc.Status = status
						}

						// Даты
						if createdAt, ok := source["created_at"].(string); ok {
							if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
								doc.IndexedAt = t
								doc.LastModified = t
							}
						}
						if updatedAt, ok := source["updated_at"].(string); ok {
							if t, err := time.Parse(time.RFC3339, updatedAt); err == nil {
								doc.LastModified = t
							}
						}

						// Все поля для поиска
						doc.SearchableFields = source

						// Добавляем документ в результат
						documents = append(documents, doc)
					}
				}
			}
		}
	}

	return documents, nil
}

// ReindexDocuments запускает переиндексацию документов
func (s *Service) ReindexDocuments(ctx context.Context, docType string) error {
	if s.osClient == nil {
		return fmt.Errorf("OpenSearch client not initialized")
	}

	if s.db == nil {
		return fmt.Errorf("database not initialized")
	}

	var totalIndexed int
	var totalErrors int

	// Определяем какие типы документов нужно переиндексировать
	shouldIndexListings := docType == "" || docType == "listing"
	shouldIndexProducts := docType == "" || docType == "product"

	// Переиндексация объявлений маркетплейса
	if shouldIndexListings {
		// Получаем все активные объявления напрямую из БД
		query := `
			SELECT 
				ml.id,
				ml.title,
				ml.description,
				ml.category_id,
				ml.user_id,
				ml.price,
				ml.status,
				ml.created_at,
				mc.name as category_name,
				u.name as user_name
			FROM marketplace_listings ml
			LEFT JOIN marketplace_categories mc ON ml.category_id = mc.id
			LEFT JOIN users u ON ml.user_id = u.id
			WHERE ml.status = 'active'
			ORDER BY ml.id
		`

		rows, err := s.db.QueryContext(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to get active listings: %w", err)
		}
		defer rows.Close()

		// Подсчитываем количество для логирования
		listingCount := 0

		// Индексируем пакетами для оптимизации
		batchSize := 100
		var batch []map[string]interface{}

		for rows.Next() {
			var listing struct {
				ID           int       `db:"id"`
				Title        string    `db:"title"`
				Description  string    `db:"description"`
				CategoryID   int       `db:"category_id"`
				UserID       int       `db:"user_id"`
				Price        float64   `db:"price"`
				Status       string    `db:"status"`
				CreatedAt    time.Time `db:"created_at"`
				CategoryName *string   `db:"category_name"`
				UserName     *string   `db:"user_name"`
			}

			if err := rows.Scan(
				&listing.ID,
				&listing.Title,
				&listing.Description,
				&listing.CategoryID,
				&listing.UserID,
				&listing.Price,
				&listing.Status,
				&listing.CreatedAt,
				&listing.CategoryName,
				&listing.UserName,
			); err != nil {
				fmt.Printf("Error scanning listing: %v\n", err)
				totalErrors++
				continue
			}

			// Создаем документ для индексации
			doc := map[string]interface{}{
				"id":          listing.ID,
				"title":       listing.Title,
				"description": listing.Description,
				"category_id": listing.CategoryID,
				"user_id":     listing.UserID,
				"price":       listing.Price,
				"status":      listing.Status,
				"created_at":  listing.CreatedAt,
				"type":        "listing",
			}

			if listing.CategoryName != nil {
				doc["category_name"] = *listing.CategoryName
			}
			if listing.UserName != nil {
				doc["user_name"] = *listing.UserName
			}

			batch = append(batch, doc)
			listingCount++

			// Индексируем пакет при достижении размера
			if len(batch) >= batchSize {
				if err := s.indexBatch(ctx, batch); err != nil {
					fmt.Printf("Error indexing batch: %v\n", err)
					totalErrors += len(batch)
				} else {
					totalIndexed += len(batch)
				}
				batch = nil
			}
		}

		// Индексируем оставшийся пакет
		if len(batch) > 0 {
			if err := s.indexBatch(ctx, batch); err != nil {
				fmt.Printf("Error indexing final batch: %v\n", err)
				totalErrors += len(batch)
			} else {
				totalIndexed += len(batch)
			}
		}

		fmt.Printf("Indexed %d listings, %d errors\n", listingCount, totalErrors)
	}

	// Переиндексация товаров витрин
	if shouldIndexProducts {
		// Получаем все активные товары витрин
		query := `
			SELECT 
				sp.id,
				sp.storefront_id,
				sp.name,
				sp.description,
				sp.category_id,
				sp.price,
				'active' as status,
				sp.created_at,
				sf.name as storefront_name,
				mc.name as category_name
			FROM storefront_products sp
			LEFT JOIN storefronts sf ON sp.storefront_id = sf.id
			LEFT JOIN marketplace_categories mc ON sp.category_id = mc.id
			WHERE sp.is_active = true
			ORDER BY sp.id
		`

		rows, err := s.db.QueryContext(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to get active products: %w", err)
		}
		defer rows.Close()

		productCount := 0
		var batch []map[string]interface{}
		batchSize := 100

		for rows.Next() {
			var product struct {
				ID             int       `db:"id"`
				StorefrontID   int       `db:"storefront_id"`
				Name           string    `db:"name"`
				Description    string    `db:"description"`
				CategoryID     *int      `db:"category_id"`
				Price          float64   `db:"price"`
				Status         string    `db:"status"`
				CreatedAt      time.Time `db:"created_at"`
				StorefrontName *string   `db:"storefront_name"`
				CategoryName   *string   `db:"category_name"`
			}

			if err := rows.Scan(
				&product.ID,
				&product.StorefrontID,
				&product.Name,
				&product.Description,
				&product.CategoryID,
				&product.Price,
				&product.Status,
				&product.CreatedAt,
				&product.StorefrontName,
				&product.CategoryName,
			); err != nil {
				fmt.Printf("Error scanning product: %v\n", err)
				totalErrors++
				continue
			}

			// Создаем документ для индексации
			doc := map[string]interface{}{
				"id":            fmt.Sprintf("sp_%d", product.ID),
				"product_id":    product.ID,
				"name":          product.Name,
				"description":   product.Description,
				"storefront_id": product.StorefrontID,
				"price":         product.Price,
				"status":        product.Status,
				"created_at":    product.CreatedAt,
				"type":          "product",
			}

			if product.CategoryID != nil {
				doc["category_id"] = *product.CategoryID
			}
			if product.StorefrontName != nil {
				doc["storefront_name"] = *product.StorefrontName
			}
			if product.CategoryName != nil {
				doc["category_name"] = *product.CategoryName
			}

			batch = append(batch, doc)
			productCount++

			// Индексируем пакет при достижении размера
			if len(batch) >= batchSize {
				if err := s.indexBatch(ctx, batch); err != nil {
					fmt.Printf("Error indexing batch: %v\n", err)
					totalErrors += len(batch)
				} else {
					totalIndexed += len(batch)
				}
				batch = nil
			}
		}

		// Индексируем оставшийся пакет
		if len(batch) > 0 {
			if err := s.indexBatch(ctx, batch); err != nil {
				fmt.Printf("Error indexing final batch: %v\n", err)
				totalErrors += len(batch)
			} else {
				totalIndexed += len(batch)
			}
		}

		fmt.Printf("Indexed %d products, %d errors\n", productCount, totalErrors)
	}

	if totalErrors > 0 {
		return fmt.Errorf("reindexing completed with %d errors, %d documents indexed", totalErrors, totalIndexed)
	}

	fmt.Printf("Reindexing completed successfully: %d documents indexed\n", totalIndexed)
	return nil
}

// indexBatch индексирует пакет документов в OpenSearch
func (s *Service) indexBatch(ctx context.Context, docs []map[string]interface{}) error {
	if len(docs) == 0 {
		return nil
	}

	// Формируем bulk запрос для OpenSearch
	var bulkBody []byte
	for _, doc := range docs {
		// Определяем ID документа
		docID := ""
		if id, ok := doc["id"].(int); ok {
			docID = fmt.Sprintf("%d", id)
		} else if id, ok := doc["id"].(string); ok {
			docID = id
		}

		// Добавляем команду индексации
		action := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": marketplaceIndex,
				"_id":    docID,
			},
		}

		actionJSON, _ := json.Marshal(action)
		bulkBody = append(bulkBody, actionJSON...)
		bulkBody = append(bulkBody, '\n')

		// Добавляем документ
		docJSON, _ := json.Marshal(doc)
		bulkBody = append(bulkBody, docJSON...)
		bulkBody = append(bulkBody, '\n')
	}

	// Отправляем bulk запрос
	_, err := s.osClient.Execute(ctx, "POST", "/_bulk", bulkBody)
	return err
}

// Вспомогательные функции

func getStringValue(data map[string]interface{}, key, defaultValue string) string {
	if val, ok := data[key].(string); ok {
		return val
	}
	return defaultValue
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func parseTimestamp(ts string) (time.Time, error) {
	// OpenSearch может возвращать timestamp в миллисекундах как строку
	// Пробуем разные форматы
	if t, err := time.Parse(time.RFC3339, ts); err == nil {
		return t, nil
	}
	if t, err := time.Parse("2006-01-02T15:04:05.000Z", ts); err == nil {
		return t, nil
	}
	// Если это число в миллисекундах
	var millis int64
	if _, err := fmt.Sscanf(ts, "%d", &millis); err == nil {
		return time.Unix(0, millis*int64(time.Millisecond)), nil
	}
	return time.Time{}, fmt.Errorf("unable to parse timestamp: %s", ts)
}
