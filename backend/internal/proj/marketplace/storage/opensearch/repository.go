// backend/internal/proj/marketplace/storage/opensearch/repository.go
package opensearch

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"backend/internal/config"
	"backend/internal/domain/models"
	"backend/internal/domain/search"
	"backend/internal/logger"
	"backend/internal/storage"
	osClient "backend/internal/storage/opensearch"
	"backend/pkg/transliteration"
)

// Repository реализует интерфейс MarketplaceSearchRepository
type Repository struct {
	client         *osClient.OpenSearchClient
	indexName      string
	storage        storage.Storage
	transliterator *transliteration.SerbianTransliterator
	boostWeights   *config.OpenSearchBoostWeights
}

// NewRepository создает новый репозиторий
func NewRepository(client *osClient.OpenSearchClient, indexName string, storage storage.Storage, searchWeights *config.SearchWeights) *Repository {
	var boostWeights *config.OpenSearchBoostWeights
	if searchWeights != nil {
		boostWeights = &searchWeights.OpenSearchBoosts
	}

	return &Repository{
		client:         client,
		indexName:      indexName,
		storage:        storage,
		transliterator: transliteration.NewSerbianTransliterator(),
		boostWeights:   boostWeights,
	}
}

// getBoostWeight возвращает вес boost из конфигурации или значение по умолчанию
func (r *Repository) getBoostWeight(weightName string, defaultValue float64) float64 {
	if r.boostWeights == nil {
		return defaultValue
	}

	switch weightName {
	case "Title":
		return r.boostWeights.Title
	case "TitleNgram":
		return r.boostWeights.TitleNgram
	case "Description":
		return r.boostWeights.Description
	case "TranslationTitle":
		return r.boostWeights.TranslationTitle
	case "TranslationDesc":
		return r.boostWeights.TranslationDesc
	case "AttributeTextValue":
		return r.boostWeights.AttributeTextValue
	case "AttributeDisplayValue":
		return r.boostWeights.AttributeDisplayValue
	case "AttributeTextValueKeyword":
		return r.boostWeights.AttributeTextValueKeyword
	case "AttributeGeneralBoost":
		return r.boostWeights.AttributeGeneralBoost
	case "RealEstateAttributesCombined":
		return r.boostWeights.RealEstateAttributesCombined
	case "PropertyType":
		return r.boostWeights.PropertyType
	case "RoomsText":
		return r.boostWeights.RoomsText
	case "CarMake":
		return r.boostWeights.CarMake
	case "CarModel":
		return r.boostWeights.CarModel
	case "CarKeywords":
		return r.boostWeights.CarKeywords
	case "PerWordTitle":
		return r.boostWeights.PerWordTitle
	case "PerWordDescription":
		return r.boostWeights.PerWordDescription
	case "PerWordAllAttributes":
		return r.boostWeights.PerWordAllAttributes
	case "PerWordRealEstateAttributes":
		return r.boostWeights.PerWordRealEstateAttributes
	case "PerWordRoomsText":
		return r.boostWeights.PerWordRoomsText
	case "AutomotiveAttributePriority":
		return r.boostWeights.AutomotiveAttributePriority
	case "SynonymBoost":
		return r.boostWeights.SynonymBoost
	default:
		return defaultValue
	}
}

// PrepareIndex подготавливает индекс (создает, если не существует)
func (r *Repository) PrepareIndex(ctx context.Context) error {
	exists, err := r.client.IndexExists(r.indexName)
	if err != nil {
		return fmt.Errorf("ошибка проверки индекса: %w", err)
	}

	logger.Info().Str("indexName", r.indexName).Bool("exists", exists).Msg("Проверка индекса")

	if !exists {
		logger.Info().Str("indexName", r.indexName).Msg("Создание индекса...")
		if err := r.client.CreateIndex(r.indexName, osClient.ListingMapping); err != nil {
			return fmt.Errorf("ошибка создания индекса: %w", err)
		}
		logger.Info().Str("indexName", r.indexName).Msg("Индекс успешно создан")

		allListings, _, err := r.storage.GetListings(ctx, map[string]string{}, 1000, 0)
		if err != nil {
			logger.Error().Err(err).Msg("Ошибка получения объявлений")
			return err
		}

		listingPtrs := make([]*models.MarketplaceListing, len(allListings))
		for i := range allListings {
			listingPtrs[i] = &allListings[i]
		}
		logger.Info().Msgf("Запуск переиндексации с обновленной схемой (поддержка метаданных и скидок)")
		if err := r.BulkIndexListings(ctx, listingPtrs); err != nil {
			logger.Error().Err(err).Msg("Ошибка индексации объявлений")
			return err
		}

		logger.Info().Int("listing_count", len(allListings)).Msgf("Успешно проиндексировано объявлений")
	}

	return nil
}

func (r *Repository) IndexListing(ctx context.Context, listing *models.MarketplaceListing) error {
	doc := r.listingToDoc(listing)
	docJSON, _ := json.MarshalIndent(doc, "", "  ")
	logger.Info().Msgf("Индексация объявления %d с данными: %s", listing.ID, string(docJSON))
	return r.client.IndexDocument(r.indexName, fmt.Sprintf("%d", listing.ID), doc)
}

// BulkIndexListings индексирует несколько объявлений
func (r *Repository) BulkIndexListings(ctx context.Context, listings []*models.MarketplaceListing) error {
	docs := make([]map[string]interface{}, 0, len(listings))

	for _, listing := range listings {
		doc := r.listingToDoc(listing)
		logger.Info().Msgf("Индексация объявления с ID: %d, категория: %d, название: %s",
			listing.ID, listing.CategoryID, listing.Title)

		if listing.ID == 0 {
			logger.Info().Msgf("ВНИМАНИЕ: Объявление с нулевым ID: %s (категория: %d)",
				listing.Title, listing.CategoryID)
		}

		doc["id"] = listing.ID
		docs = append(docs, doc)
	}

	return r.client.BulkIndex(r.indexName, docs)
}

// DeleteListing удаляет объявление из индекса
func (r *Repository) DeleteListing(ctx context.Context, listingID string) error {
	return r.client.DeleteDocument(r.indexName, listingID)
}

// GetClient возвращает клиент OpenSearch
func (r *Repository) GetClient() *osClient.OpenSearchClient {
	return r.client
}

// Метод для извлечения ID из документа OpenSearch
func (r *Repository) extractDocumentID(hit map[string]interface{}) (int, error) {
	if idStr, ok := hit["_id"].(string); ok {
		if id, err := strconv.Atoi(idStr); err == nil {
			return id, nil
		}
	}

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

	return 0, fmt.Errorf("не удалось получить ID объявления из документа")
}

// SearchListings выполняет поиск объявлений
func (r *Repository) SearchListings(ctx context.Context, params *search.SearchParams) (*search.SearchResult, error) {
	var query map[string]interface{}

	// ВАЖНО: проверяем наличие CustomQuery и используем его напрямую, если он задан
	if params.CustomQuery != nil {
		query = params.CustomQuery
		// Логируем, что используем специальный запрос
		queryJSON, _ := json.MarshalIndent(query, "", "  ")
		logger.Info().Msgf("Используем специальный запрос для поиска: %s", string(queryJSON))
	} else {
		query = r.buildSearchQuery(params)
	}

	response, err := r.client.Search(r.indexName, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения поиска: %w", err)
	}

	var searchResponse map[string]interface{}
	if err := json.Unmarshal(response, &searchResponse); err != nil {
		return nil, fmt.Errorf("ошибка разбора ответа: %w", err)
	}

	logger.Info().Msgf("OpenSearch ответил успешно. Анализируем результаты...")
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

	responseBytes, err := r.client.Search(r.indexName, query)
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
func extractSuggestionsFromAgg(aggs map[string]interface{}, aggName string, suggestions map[string]bool) {
	if agg, ok := aggs[aggName].(map[string]interface{}); ok {
		extractBucketsFromAgg(agg, suggestions)
	}
}

// Вспомогательная функция для извлечения бакетов из агрегации
func extractBucketsFromAgg(agg map[string]interface{}, suggestions map[string]bool) {
	if buckets, ok := agg["buckets"].([]interface{}); ok {
		for _, bucket := range buckets {
			if bucketObj, ok := bucket.(map[string]interface{}); ok {
				if key, ok := bucketObj["key"].(string); ok && key != "" {
					suggestions[key] = true
				}
			}
		}
	}
}

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
	exists, err := r.client.IndexExists(r.indexName)
	if err != nil {
		return fmt.Errorf("ошибка проверки индекса: %w", err)
	}

	if exists {
		logger.Info().Msgf("Удаляем существующий индекс %s", r.indexName)
		if err := r.client.DeleteIndex(r.indexName); err != nil {
			return fmt.Errorf("ошибка удаления индекса: %w", err)
		}
		time.Sleep(1 * time.Second)
	}

	logger.Info().Msgf("Создаем индекс %s заново", r.indexName)
	if err := r.client.CreateIndex(r.indexName, osClient.ListingMapping); err != nil {
		return fmt.Errorf("ошибка создания индекса: %w", err)
	}

	const batchSize = 100
	offset := 0
	totalIndexed := 0

	for {
		logger.Info().Msgf("Получение пакета объявлений (размер: %d, смещение: %d)", batchSize, offset)
		listings, total, err := r.storage.GetListings(ctx, map[string]string{}, batchSize, offset)
		if err != nil {
			return fmt.Errorf("ошибка получения объявлений: %w", err)
		}

		if len(listings) == 0 {
			break
		}

		logger.Info().Msgf("Получено %d объявлений из %d всего (пакет %d)", len(listings), total, offset/batchSize+1)

		listingPtrs := make([]*models.MarketplaceListing, len(listings))
		for i := range listings {
			listingID := listings[i].ID

			// Проверяем наличие переводов и при необходимости загружаем их
			if listings[i].Translations == nil || len(listings[i].Translations) == 0 {
				translations, err := r.storage.GetTranslationsForEntity(ctx, "listing", listingID)
				if err == nil && len(translations) > 0 {
					transMap := make(models.TranslationMap)
					for _, t := range translations {
						if _, ok := transMap[t.Language]; !ok {
							transMap[t.Language] = make(map[string]string)
						}
						transMap[t.Language][t.FieldName] = t.TranslatedText
					}
					listings[i].Translations = transMap
					logger.Info().Msgf("Загружено %d переводов для объявления %d", len(translations), listingID)
				}
			}

			// Проверяем наличие атрибутов и при необходимости загружаем их
			if listings[i].Attributes == nil || len(listings[i].Attributes) == 0 {
				attrs, err := r.storage.GetListingAttributes(ctx, listingID)
				if err == nil && len(attrs) > 0 {
					listings[i].Attributes = attrs
					logger.Info().Msgf("Загружено %d атрибутов для объявления %d", len(attrs), listingID)
				}
			}

			listingPtrs[i] = &listings[i]
		}
		if err := r.BulkIndexListings(ctx, listingPtrs); err != nil {
			return fmt.Errorf("ошибка массовой индексации (пакет %d): %w", offset/batchSize+1, err)
		}

		totalIndexed += len(listings)
		offset += batchSize

		if len(listings) < batchSize {
			break
		}
	}

	logger.Info().Msgf("Успешно проиндексировано %d объявлений", totalIndexed)
	return nil
}

func (r *Repository) getAttributeOptionTranslations(attrName, value string) (map[string]string, error) {
	ctx := context.Background()
	query := `
        SELECT option_value, en_translation, sr_translation
        FROM attribute_option_translations
        WHERE attribute_name = $1 AND option_value = $2
    `

	var optionValue, enTranslation, srTranslation string
	err := r.storage.QueryRow(ctx, query, attrName, value).Scan(
		&optionValue, &enTranslation, &srTranslation,
	)
	if err != nil {
		if err != pgx.ErrNoRows {
			logger.Info().Msgf("Ошибка получения переводов для атрибута %s, значение %s: %v", attrName, value, err)
		}
		return nil, err
	}

	translations := map[string]string{
		"en": enTranslation,
		"sr": srTranslation,
	}

	return translations, nil
}

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
		"average_rating":    listing.AverageRating,
		"review_count":      listing.ReviewCount,
	}

	logger.Info().Msgf("Обработка местоположения для листинга %d: город=%s, страна=%s, адрес=%s",
		listing.ID, listing.City, listing.Country, listing.Location)

	realEstateFields := createRealEstateFieldsMap()
	carFields := createCarFieldsMap()
	importantAttrs := createImportantAttributesMap()

	if listing.Attributes != nil && len(listing.Attributes) > 0 {
		processAttributesForIndex(doc, listing.Attributes, importantAttrs, realEstateFields, carFields, listing.ID, r)
	}

	processDiscountData(doc, listing)
	processMetadata(doc, listing)
	processStorefrontData(doc, listing, r.storage)
	processCoordinates(doc, listing, r)
	processCategoryPath(doc, listing, r.storage)
	processCategory(doc, listing)
	processUser(doc, listing)
	processImages(doc, listing, r.storage)

	docJSON, _ := json.MarshalIndent(doc, "", "  ")
	logger.Info().Msgf("FINAL DOC for listing %d [size=%d bytes]: %s", listing.ID, len(docJSON), string(docJSON))

	return doc
}

func processAttributesForIndex(doc map[string]interface{}, attributes []models.ListingAttributeValue,
	importantAttrs, realEstateFields, carFields map[string]bool, listingID int, r *Repository,
) {
	realEstateText := make([]string, 0)
	makeValue, modelValue := "", ""
	uniqueTextValues := make(map[string]bool)
	attributeTextValues := make(map[string][]string)
	selectValues := []string{}
	seen := make(map[int]bool)
	attributesArray := make([]map[string]interface{}, 0, len(attributes))
	carKeywords := []string{} // Для ключевых слов автомобиля

	for _, attr := range attributes {
		if seen[attr.AttributeID] {
			continue
		}
		seen[attr.AttributeID] = true

		if !hasAttributeValue(attr) {
			continue
		}

		attrDoc := createAttributeDocument(attr)
		attributesArray = append(attributesArray, attrDoc)

		if attr.TextValue != nil && *attr.TextValue != "" {
			textValue := *attr.TextValue

			switch attr.AttributeName {
			case "make":
				makeValue = textValue
				doc["make"] = makeValue
				doc["make_lowercase"] = strings.ToLower(makeValue)
				carKeywords = append(carKeywords, textValue, strings.ToLower(textValue)) // Добавляем к ключевым словам
				logger.Info().Msgf("FIRST PASS: Добавлена марка '%s' в корень документа для объявления %d", makeValue, listingID)
			case "model":
				modelValue = textValue
				doc["model"] = modelValue
				doc["model_lowercase"] = strings.ToLower(modelValue)
				carKeywords = append(carKeywords, textValue, strings.ToLower(textValue)) // Добавляем к ключевым словам
				logger.Info().Msgf("FIRST PASS: Добавлена модель '%s' в корень документа для объявления %d", modelValue, listingID)
			default:
				if isImportantTextAttribute(attr.AttributeName) {
					doc[attr.AttributeName] = textValue
					doc[attr.AttributeName+"_lowercase"] = strings.ToLower(textValue)
					logger.Info().Msgf("FIRST PASS: Добавлен важный атрибут %s = '%s' в корень документа для объявления %d",
						attr.AttributeName, textValue, listingID)
				}
			}

			if !uniqueTextValues[textValue] {
				attributeTextValues[attr.AttributeName] = append(attributeTextValues[attr.AttributeName], textValue)
				uniqueTextValues[textValue] = true
			}
			lowerValue := strings.ToLower(textValue)
			if !uniqueTextValues[lowerValue] {
				attributeTextValues[attr.AttributeName] = append(attributeTextValues[attr.AttributeName], lowerValue)
				uniqueTextValues[lowerValue] = true
			}

			if attr.AttributeName == "make" || attr.AttributeName == "model" ||
				attr.AttributeName == "brand" || attr.AttributeName == "color" {
				// Если есть текстовое значение, добавляем его и в нижнем регистре
				if attr.TextValue != nil && *attr.TextValue != "" {
					doc[attr.AttributeName] = *attr.TextValue
					doc[attr.AttributeName+"_lowercase"] = strings.ToLower(*attr.TextValue)
					logger.Info().Msgf("Добавлен важный атрибут %s = '%s' в корень документа для объявления %d",
						attr.AttributeName, *attr.TextValue, listingID)
				} else if attr.DisplayValue != "" {
					// Если есть только отображаемое значение
					doc[attr.AttributeName] = attr.DisplayValue
					doc[attr.AttributeName+"_lowercase"] = strings.ToLower(attr.DisplayValue)
					logger.Info().Msgf("Добавлен важный атрибут (из DisplayValue) %s = '%s' в корень документа для объявления %d",
						attr.AttributeName, attr.DisplayValue, listingID)
				}
			}

			if attr.AttributeType == "select" {
				selectValues = append(selectValues, textValue, strings.ToLower(textValue))

				value := textValue
				if attr.DisplayValue != "" {
					value = attr.DisplayValue
				}
				if value != "" {
					translations, err := r.getAttributeOptionTranslations(attr.AttributeName, value)
					if err == nil && len(translations) > 0 {
						attrDoc["translations"] = translations
						for lang, translation := range translations {
							if _, ok := attributeTextValues[attr.AttributeName+"_"+lang]; !ok {
								attributeTextValues[attr.AttributeName+"_"+lang] = []string{}
							}
							attributeTextValues[attr.AttributeName+"_"+lang] = append(
								attributeTextValues[attr.AttributeName+"_"+lang],
								translation,
								strings.ToLower(translation),
							)
						}
					}
				}
			}
		}

		if attr.NumericValue != nil {
			numVal := *attr.NumericValue
			if !math.IsNaN(numVal) && !math.IsInf(numVal, 0) {
				if realEstateFields[attr.AttributeName] || carFields[attr.AttributeName] || importantAttrs[attr.AttributeName] {
					doc[attr.AttributeName] = numVal
					displayValue := formatAttributeDisplayValue(attr)
					doc[attr.AttributeName+"_text"] = displayValue
					realEstateText = append(realEstateText, displayValue)
					addRangesForAttribute(doc, attr)
					logger.Info().Msgf("FIRST PASS: Добавлен числовой атрибут %s = %f в корень документа для объявления %d",
						attr.AttributeName, numVal, listingID)
				}
			}
		}

		if attr.BooleanValue != nil {
			boolValue := *attr.BooleanValue
			if importantAttrs[attr.AttributeName] {
				doc[attr.AttributeName] = boolValue
				strValue := "нет"
				if boolValue {
					strValue = "да"
				}
				doc[attr.AttributeName+"_text"] = strValue
				realEstateText = append(realEstateText, strValue)
				logger.Info().Msgf("FIRST PASS: Добавлен булев атрибут %s = %v в корень документа для объявления %d",
					attr.AttributeName, boolValue, listingID)
			}
		}

		if attr.JSONValue != nil {
			jsonStr := string(attr.JSONValue)
			if jsonStr != "" && jsonStr != "{}" && jsonStr != "[]" {
				attrDoc["json_value"] = jsonStr
				var jsonData interface{}
				if err := json.Unmarshal(attr.JSONValue, &jsonData); err == nil {
					if strArray, ok := jsonData.([]string); ok {
						attrDoc["json_array"] = strArray
						attributeTextValues[attr.AttributeName] = append(
							attributeTextValues[attr.AttributeName],
							strArray...,
						)
					}
				}
			}
		}
	}

	ensureImportantAttributes(doc, makeValue, modelValue, listingID)

	if len(selectValues) > 0 {
		doc["select_values"] = getUniqueValues(selectValues)
	}

	// Добавляем собранные ключевые слова по автомобилю для улучшения поиска
	if len(carKeywords) > 0 {
		doc["car_keywords"] = getUniqueValues(carKeywords)
	}

	doc["attributes"] = attributesArray
	allAttrsText := getUniqueValues(flattenAttributeValues(attributeTextValues))
	doc["all_attributes_text"] = allAttrsText

	// Отладочный вывод для проверки индексации атрибутов
	logger.Info().Msgf("ИНДЕКСАЦИЯ объявления %d: attributes=%d, all_attributes_text=%v",
		listingID, len(attributesArray), allAttrsText)

	if len(realEstateText) > 0 {
		doc["real_estate_attributes_text"] = realEstateText
		doc["real_estate_attributes_combined"] = strings.Join(realEstateText, " ")
	}
}

func hasAttributeValue(attr models.ListingAttributeValue) bool {
	return (attr.TextValue != nil && *attr.TextValue != "") ||
		(attr.NumericValue != nil) ||
		(attr.BooleanValue != nil) ||
		(attr.JSONValue != nil && string(attr.JSONValue) != "{}" && string(attr.JSONValue) != "[]") ||
		attr.DisplayValue != ""
}

func isImportantTextAttribute(attrName string) bool {
	importantTextAttrs := map[string]bool{
		"brand":         true,
		"color":         true,
		"fuel_type":     true,
		"transmission":  true,
		"body_type":     true,
		"property_type": true,
	}
	return importantTextAttrs[attrName]
}

func formatAttributeDisplayValue(attr models.ListingAttributeValue) string {
	numVal := *attr.NumericValue
	unitStr := attr.Unit

	if unitStr == "" {
		switch attr.AttributeName {
		case "area":
			unitStr = "m²"
		case "land_area":
			unitStr = "ar"
		case "mileage":
			unitStr = "km"
		case "engine_capacity":
			unitStr = "l"
		case "power":
			unitStr = "ks"
		case "screen_size":
			unitStr = "inč"
		case "rooms":
			unitStr = "soba"
		case "floor", "total_floors":
			unitStr = "sprat"
		}
	}

	if attr.AttributeName == "year" {
		return fmt.Sprintf("%d", int(numVal))
	} else if unitStr != "" {
		return fmt.Sprintf("%g %s", numVal, unitStr)
	}
	return fmt.Sprintf("%g", numVal)
}

func addRangesForAttribute(doc map[string]interface{}, attr models.ListingAttributeValue) {
	if attr.NumericValue == nil {
		return
	}
	numVal := *attr.NumericValue

	switch attr.AttributeName {
	case "price":
		doc["price_range"] = getPriceRange(int(numVal))
	case "mileage":
		doc["mileage_range"] = getMileageRange(int(numVal))
	case "area":
		if numVal <= 30 {
			doc["area_range"] = "do 30 m²"
		} else if numVal <= 50 {
			doc["area_range"] = "30-50 m²"
		} else if numVal <= 80 {
			doc["area_range"] = "50-80 m²"
		} else if numVal <= 120 {
			doc["area_range"] = "80-120 m²"
		} else {
			doc["area_range"] = "od 120 m²"
		}
	}
}

func createAttributeDocument(attr models.ListingAttributeValue) map[string]interface{} {
	attrDoc := map[string]interface{}{
		"attribute_id":   attr.AttributeID,
		"attribute_name": attr.AttributeName,
		"display_name":   attr.DisplayName,
		"attribute_type": attr.AttributeType,
		"display_value":  attr.DisplayValue,
	}

	if attr.TextValue != nil && *attr.TextValue != "" {
		textValue := *attr.TextValue
		attrDoc["text_value"] = textValue
		attrDoc["text_value_lowercase"] = strings.ToLower(textValue)
	}

	if attr.NumericValue != nil && !math.IsNaN(*attr.NumericValue) && !math.IsInf(*attr.NumericValue, 0) {
		attrDoc["numeric_value"] = *attr.NumericValue
		if attr.Unit != "" {
			attrDoc["unit"] = attr.Unit
		}
	}

	if attr.BooleanValue != nil {
		attrDoc["boolean_value"] = *attr.BooleanValue
	}

	if attr.JSONValue != nil {
		jsonStr := string(attr.JSONValue)
		if jsonStr != "" && jsonStr != "{}" && jsonStr != "[]" {
			attrDoc["json_value"] = jsonStr
		}
	}

	return attrDoc
}

func ensureImportantAttributes(doc map[string]interface{}, makeValue, modelValue string, listingID int) {
	if makeValue != "" && doc["make"] == nil {
		doc["make"] = makeValue
		doc["make_lowercase"] = strings.ToLower(makeValue)
		logger.Info().Msgf("FINAL CHECK: Добавлена марка '%s' в корень документа для объявления %d",
			makeValue, listingID)
	}

	if modelValue != "" && doc["model"] == nil {
		doc["model"] = modelValue
		doc["model_lowercase"] = strings.ToLower(modelValue)
		logger.Info().Msgf("FINAL CHECK: Добавлена модель '%s' в корень документа для объявления %d",
			modelValue, listingID)
	}
}

func getUniqueValues(values []string) []string {
	seen := make(map[string]bool)
	unique := make([]string, 0, len(values))

	for _, val := range values {
		if !seen[val] {
			seen[val] = true
			unique = append(unique, val)
		}
	}

	return unique
}

func flattenAttributeValues(attributeTextValues map[string][]string) []string {
	var result []string
	for _, values := range attributeTextValues {
		result = append(result, values...)
	}
	return result
}

func createRealEstateFieldsMap() map[string]bool {
	return map[string]bool{
		"rooms":            true,
		"floor":            true,
		"total_floors":     true,
		"area":             true,
		"land_area":        true,
		"property_type":    true,
		"year_built":       true,
		"bathrooms":        true,
		"condition":        true,
		"amenities":        true,
		"heating_type":     true,
		"parking":          true,
		"balcony":          true,
		"furnished":        true,
		"air_conditioning": true,
	}
}

func createCarFieldsMap() map[string]bool {
	return map[string]bool{
		"make":            true,
		"model":           true,
		"year":            true,
		"mileage":         true,
		"engine_capacity": true,
		"fuel_type":       true,
		"transmission":    true,
		"body_type":       true,
	}
}

func createImportantAttributesMap() map[string]bool {
	return map[string]bool{
		"make":            true,
		"model":           true,
		"brand":           true,
		"year":            true,
		"color":           true,
		"rooms":           true,
		"property_type":   true,
		"body_type":       true,
		"engine_capacity": true,
		"fuel_type":       true,
		"transmission":    true,
		"cpu":             true,
		"gpu":             true,
		"memory":          true,
		"ram":             true,
		"storage_type":    true,
		"screen_size":     true,
	}
}

func processDiscountData(doc map[string]interface{}, listing *models.MarketplaceListing) {
	doc["has_discount"] = listing.HasDiscount
	if listing.OldPrice > 0 {
		doc["old_price"] = listing.OldPrice
	}

	if strings.Contains(listing.Description, "СКИДКА") || strings.Contains(listing.Description, "СКИДКА!") {
		discountRegex := regexp.MustCompile(`(\d+)%\s*СКИДКА`)
		matches := discountRegex.FindStringSubmatch(listing.Description)
		priceRegex := regexp.MustCompile(`Старая цена:\s*(\d+[\.,]?\d*)\s*RSD`)
		priceMatches := priceRegex.FindStringSubmatch(listing.Description)

		if len(matches) > 1 && len(priceMatches) > 1 {
			discountPercent, _ := strconv.Atoi(matches[1])
			oldPriceStr := strings.Replace(priceMatches[1], ",", ".", -1)
			oldPrice, _ := strconv.ParseFloat(oldPriceStr, 64)

			if listing.Metadata == nil {
				listing.Metadata = make(map[string]interface{})
			}

			discount := map[string]interface{}{
				"discount_percent":  discountPercent,
				"previous_price":    oldPrice,
				"effective_from":    time.Now().AddDate(0, 0, -10).Format(time.RFC3339),
				"has_price_history": true,
			}

			listing.Metadata["discount"] = discount
			doc["old_price"] = oldPrice
			doc["has_discount"] = true

			logger.Info().Msgf("Extracted discount from description for listing %d: %v", listing.ID, discount)
		}
	}
}

func processMetadata(doc map[string]interface{}, listing *models.MarketplaceListing) {
	if listing.Metadata != nil {
		doc["metadata"] = listing.Metadata
		if discount, ok := listing.Metadata["discount"].(map[string]interface{}); ok {
			if prevPrice, ok := discount["previous_price"].(float64); ok && prevPrice > 0 {
				discountPercent := int((prevPrice - listing.Price) / prevPrice * 100)
				discount["discount_percent"] = discountPercent
				listing.Metadata["discount"] = discount
				doc["metadata"] = listing.Metadata
				doc["has_discount"] = true
				doc["old_price"] = prevPrice

				logger.Info().Msgf("Recalculated discount percent for OpenSearch: %d%% (listing %d)",
					discountPercent, listing.ID)

				// Проверяем, чтобы storefront_id сохранился для объявлений со скидками
				if listing.StorefrontID != nil && *listing.StorefrontID > 0 {
					logger.Info().Msgf("ВАЖНО: Сохраняем storefront_id=%d после обработки скидки для объявления %d",
						*listing.StorefrontID, listing.ID)
				}
			}
		}
	}
}

func processStorefrontData(doc map[string]interface{}, listing *models.MarketplaceListing, storage storage.Storage) {
	// Всегда добавляем storefront_id в документ, если он есть
	if listing.StorefrontID != nil && *listing.StorefrontID > 0 {
		// Явно проверяем и выводим информацию о том, что мы добавляем storefront_id
		doc["storefront_id"] = *listing.StorefrontID
		logger.Info().Msgf("ВАЖНО: Устанавливаем storefront_id=%d для объявления %d с скидкой=%v",
			*listing.StorefrontID, listing.ID, listing.HasDiscount)

		needStorefrontInfo := listing.City == "" || listing.Country == "" || listing.Location == "" ||
			listing.Latitude == nil || listing.Longitude == nil

		if needStorefrontInfo {
			logger.Info().Msgf("Fetching storefront %d data for listing %d address info", *listing.StorefrontID, listing.ID)
			var storefront models.Storefront
			err := storage.QueryRow(context.Background(), `
                SELECT name, city, address, country, latitude, longitude
                FROM user_storefronts
                WHERE id = $1
            `, *listing.StorefrontID).Scan(
				&storefront.Name,
				&storefront.City,
				&storefront.Address,
				&storefront.Country,
				&storefront.Latitude,
				&storefront.Longitude,
			)

			if err == nil {
				if listing.City == "" && storefront.City != "" {
					doc["city"] = storefront.City
				}
				if listing.Country == "" && storefront.Country != "" {
					doc["country"] = storefront.Country
				}
				if listing.Location == "" && storefront.Address != "" {
					doc["location"] = storefront.Address
				}
				if (listing.Latitude == nil || listing.Longitude == nil ||
					*listing.Latitude == 0 || *listing.Longitude == 0) &&
					storefront.Latitude != nil && storefront.Longitude != nil &&
					*storefront.Latitude != 0 && *storefront.Longitude != 0 {
					doc["coordinates"] = map[string]interface{}{
						"lat": *storefront.Latitude,
						"lon": *storefront.Longitude,
					}
					doc["show_on_map"] = true
				}
			} else {
				logger.Info().Msgf("WARNING: Failed to load storefront data for listing %d: %v", listing.ID, err)
			}
		}
	} else {
		logger.Info().Msgf("DEBUG: Listing %d has no storefront_id", listing.ID)
	}

	// Дополнительная проверка после обработки всех метаданных и скидок
	if listing.HasDiscount && listing.StorefrontID != nil &&
		doc["storefront_id"] == nil {
		doc["storefront_id"] = *listing.StorefrontID
		logger.Info().Msgf("КРИТИЧНО: Добавлен storefront_id=%d для объявления %d со скидкой в конце обработки",
			*listing.StorefrontID, listing.ID)
	}
}

func processCoordinates(doc map[string]interface{}, listing *models.MarketplaceListing, r *Repository) {
	if listing.Latitude != nil && listing.Longitude != nil && *listing.Latitude != 0 && *listing.Longitude != 0 {
		doc["coordinates"] = map[string]interface{}{
			"lat": *listing.Latitude,
			"lon": *listing.Longitude,
		}
	} else if _, ok := doc["coordinates"]; !ok {
		if cityVal, ok := doc["city"].(string); ok && cityVal != "" {
			countryVal := ""
			if c, ok := doc["country"].(string); ok {
				countryVal = c
			}
			geocoded, err := r.geocodeCity(cityVal, countryVal)
			if err == nil && geocoded != nil {
				doc["coordinates"] = map[string]interface{}{
					"lat": geocoded.Lat,
					"lon": geocoded.Lon,
				}
				doc["show_on_map"] = true
			}
		}
	}
}

func processCategoryPath(doc map[string]interface{}, listing *models.MarketplaceListing, storage storage.Storage) {
	if listing.CategoryPathIds != nil && len(listing.CategoryPathIds) > 0 {
		doc["category_path_ids"] = listing.CategoryPathIds
	} else {
		parentID := listing.CategoryID
		pathIDs := []int{parentID}
		for parentID > 0 {
			var cat models.MarketplaceCategory
			err := storage.QueryRow(context.Background(),
				"SELECT parent_id FROM marketplace_categories WHERE id = $1", parentID).
				Scan(&cat.ParentID)
			if err != nil || cat.ParentID == nil {
				break
			}
			parentID = *cat.ParentID
			pathIDs = append([]int{parentID}, pathIDs...)
		}
		doc["category_path_ids"] = pathIDs
	}
}

func processCategory(doc map[string]interface{}, listing *models.MarketplaceListing) {
	if listing.Category != nil {
		doc["category"] = map[string]interface{}{
			"id":   listing.Category.ID,
			"name": listing.Category.Name,
			"slug": listing.Category.Slug,
		}
	}
}

func processUser(doc map[string]interface{}, listing *models.MarketplaceListing) {
	if listing.User != nil {
		doc["user"] = map[string]interface{}{
			"id":    listing.User.ID,
			"name":  listing.User.Name,
			"email": listing.User.Email,
		}
	}
}

func processImages(doc map[string]interface{}, listing *models.MarketplaceListing, storage storage.Storage) {
	if listing.Images != nil && len(listing.Images) > 0 {
		imagesDoc := make([]map[string]interface{}, 0, len(listing.Images))
		for _, img := range listing.Images {
			imageDoc := map[string]interface{}{
				"id":        img.ID,
				"file_path": img.FilePath,
				"is_main":   img.IsMain,
			}

			// Добавляем поля storage_type и public_url, если они есть
			if img.StorageType != "" {
				imageDoc["storage_type"] = img.StorageType
			}
			if img.PublicURL != "" {
				imageDoc["public_url"] = img.PublicURL
			}

			imagesDoc = append(imagesDoc, imageDoc)
		}
		doc["images"] = imagesDoc
	} else {
		images, err := storage.GetListingImages(context.Background(), fmt.Sprintf("%d", listing.ID))
		if err == nil && len(images) > 0 {
			imagesDoc := make([]map[string]interface{}, 0, len(images))
			for _, img := range images {
				imageDoc := map[string]interface{}{
					"id":        img.ID,
					"file_path": img.FilePath,
					"is_main":   img.IsMain,
				}

				// Добавляем поля storage_type и public_url, если они есть
				if img.StorageType != "" {
					imageDoc["storage_type"] = img.StorageType
				}
				if img.PublicURL != "" {
					imageDoc["public_url"] = img.PublicURL
				}

				imagesDoc = append(imagesDoc, imageDoc)
			}
			doc["images"] = imagesDoc
		}
	}
}

func getLocalizedRoomText(rooms float64) string {
	switch int(rooms) {
	case 1:
		return "комната"
	case 2, 3, 4:
		return "комнаты"
	default:
		return "комнат"
	}
}

func getMileageRange(mileage int) string {
	switch {
	case mileage <= 5000:
		return "0-5000"
	case mileage <= 10000:
		return "5001-10000"
	case mileage <= 50000:
		return "10001-50000"
	case mileage <= 100000:
		return "50001-100000"
	case mileage <= 150000:
		return "100001-150000"
	case mileage <= 200000:
		return "150001-200000"
	default:
		return "200001+"
	}
}

func getPriceRange(price int) string {
	switch {
	case price <= 5000:
		return "0-5000"
	case price <= 10000:
		return "5001-10000"
	case price <= 20000:
		return "10001-20000"
	case price <= 50000:
		return "20001-50000"
	case price <= 100000:
		return "50001-100000"
	case price <= 500000:
		return "100001-500000"
	default:
		return "500001+"
	}
}

func isImportantAttribute(attrName string) bool {
	importantAttrs := map[string]bool{
		"make":            true,
		"model":           true,
		"brand":           true,
		"year":            true,
		"color":           true,
		"rooms":           true,
		"property_type":   true,
		"body_type":       true,
		"engine_capacity": true,
		"fuel_type":       true,
		"transmission":    true,
		"cpu":             true,
		"gpu":             true,
		"memory":          true,
		"ram":             true,
		"storage_type":    true,
		"screen_size":     true,
	}
	return importantAttrs[attrName]
}

func (r *Repository) geocodeCity(city, country string) (*struct{ Lat, Lon float64 }, error) {
	query := city
	if country != "" {
		query += ", " + country
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	url := fmt.Sprintf(
		"https://nominatim.openstreetmap.org/search?format=json&q=%s&limit=1",
		url.QueryEscape(query),
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "HostelBookingSystem/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неверный статус ответа: %d", resp.StatusCode)
	}

	var results []struct {
		Lat string `json:"lat"`
		Lon string `json:"lon"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("город не найден")
	}

	lat, err := strconv.ParseFloat(results[0].Lat, 64)
	if err != nil {
		return nil, err
	}

	lon, err := strconv.ParseFloat(results[0].Lon, 64)
	if err != nil {
		return nil, err
	}

	return &struct{ Lat, Lon float64 }{lat, lon}, nil
}

func (r *Repository) buildSearchQuery(params *search.SearchParams) map[string]interface{} {
	logger.Info().Msgf("Строим запрос: категория = %v, язык = %s, поисковый запрос = %s",
		params.CategoryID, params.Language, params.Query)

	query := map[string]interface{}{
		"from": (params.Page - 1) * params.Size,
		"size": params.Size,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must":   []interface{}{},
				"filter": []interface{}{},
				"should": []interface{}{},
			},
		},
	}

	if params.Status == "" {
		filter := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]interface{})
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"status": "active",
			},
		})
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = filter
	}

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
				expandedQuery, err := r.storage.ExpandSearchQuery(context.Background(), params.Query, params.Language)
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
		for _, queryVariant := range queryVariants {
			should = append(should, map[string]interface{}{
				"match": map[string]interface{}{
					"title": map[string]interface{}{
						"query":     queryVariant,
						"boost":     r.getBoostWeight("Title", 5.0),
						"fuzziness": fuzziness,
					},
				},
			})
		}

		// Добавляем поиск по n-граммам для лучшего нечеткого соответствия
		if params.UseSynonyms {
			for _, queryVariant := range queryVariants {
				should = append(should, map[string]interface{}{
					"match": map[string]interface{}{
						"title.ngram": map[string]interface{}{
							"query": queryVariant,
							"boost": r.getBoostWeight("TitleNgram", 2.0),
						},
					},
				})
			}
		}

		// Поиск по описанию для всех вариантов транслитерации
		for _, queryVariant := range queryVariants {
			should = append(should, map[string]interface{}{
				"match": map[string]interface{}{
					"description": map[string]interface{}{
						"query":     queryVariant,
						"boost":     r.getBoostWeight("Description", 2.0),
						"fuzziness": fuzziness,
					},
				},
			})
		}

		// Поиск по сербским переводам для всех вариантов транслитерации
		for _, queryVariant := range queryVariants {
			should = append(should, map[string]interface{}{
				"match": map[string]interface{}{
					"translations.sr.title": map[string]interface{}{
						"query":     queryVariant,
						"boost":     r.getBoostWeight("TranslationTitle", 4.0),
						"fuzziness": fuzziness,
					},
				},
			})

			should = append(should, map[string]interface{}{
				"match": map[string]interface{}{
					"translations.sr.description": map[string]interface{}{
						"query":     queryVariant,
						"boost":     r.getBoostWeight("TranslationDesc", 1.5),
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
		must := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]interface{})

		multiMatchShould := []map[string]interface{}{}
		for _, queryVariant := range queryVariants {
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

	if params.CategoryID != nil && *params.CategoryID > 0 {
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

	if params.Status != "" {
		filter := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]interface{})
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"status": params.Status,
			},
		})
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = filter
	}

	if params.Location != nil && params.Distance != "" {
		geoFilter := map[string]interface{}{
			"geo_distance": map[string]interface{}{
				"distance": params.Distance,
				"coordinates": map[string]interface{}{
					"lat": params.Location.Lat,
					"lon": params.Location.Lon,
				},
			},
		}

		filter := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]interface{})
		filter = append(filter, geoFilter)
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = filter
	}

	if params.AttributeFilters != nil && len(params.AttributeFilters) > 0 {
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
				if strings.Contains(attrValue, ",") {
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
				} else if attrName == "property_type" || attrName == "building_type" {
					filter = append(filter, map[string]interface{}{
						"term": map[string]interface{}{
							attrName: attrValue,
						},
					})
					logger.Info().Msgf("Added term filter for text real estate attribute %s = %s",
						attrName, attrValue)
				} else if attrValue == "true" || attrValue == "false" {
					boolVal := attrValue == "true"
					filter = append(filter, map[string]interface{}{
						"term": map[string]interface{}{
							attrName: boolVal,
						},
					})
					logger.Info().Msgf("Added boolean filter for real estate attribute %s = %v",
						attrName, boolVal)
				} else {
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
				} else if attrValue == "true" || attrValue == "false" {
					boolVal := attrValue == "true"
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
		var sortField string
		var sortOrder string

		if params.SortDirection != "" {
			sortOrder = params.SortDirection
		} else {
			sortOrder = "desc"
		}

		switch params.Sort {
		case "relevance":
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
			sortField = "created_at"
			// sortOrder уже установлен из params.SortDirection выше
		case "date_desc":
			sortField = "created_at"
			sortOrder = "desc"
		case "date_asc":
			sortField = "created_at"
			sortOrder = "asc"
		case "price":
			sortField = "price"
			// sortOrder уже установлен из params.SortDirection выше
		case "price_desc":
			sortField = "price"
			sortOrder = "desc"
		case "price_asc":
			sortField = "price"
			sortOrder = "asc"
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
		query["sort"] = []interface{}{
			map[string]interface{}{
				"created_at": map[string]interface{}{
					"order": "desc",
				},
			},
		}
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

		if hitsArray, ok := hits["hits"].([]interface{}); ok {
			for _, hitI := range hitsArray {
				if hit, ok := hitI.(map[string]interface{}); ok {
					if source, ok := hit["_source"].(map[string]interface{}); ok {
						if idStr, ok := hit["_id"].(string); ok {
							if id, err := strconv.Atoi(idStr); err == nil {
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

				// Добавляем поля storage_type и public_url, если они есть
				if storageType, ok := img["storage_type"].(string); ok {
					image.StorageType = storageType
				} else if filePath, ok := img["file_path"].(string); ok && strings.Contains(filePath, "listings/") {
					// Если путь содержит "listings/", то это MinIO
					image.StorageType = "minio"
				}

				if publicURL, ok := img["public_url"].(string); ok {
					image.PublicURL = publicURL
				} else if image.StorageType == "minio" && image.FilePath != "" {
					// Если это MinIO, но public_url не указан, формируем его
					image.PublicURL = "/listings/" + image.FilePath
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
		listing.OldPrice = oldPrice
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
				listing.OldPrice = prevPrice
				logger.Info().Msgf("Найдена скидка для объявления %d: скидка %v%%, старая цена: %.2f",
					listing.ID, discount["discount_percent"], prevPrice)
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
	if listing.OldPrice > 0 && listing.OldPrice > listing.Price {
		listing.HasDiscount = true
	}
	if listing.ID == 18 {
		logger.Info().Msgf("DEBUG: Преобразование документа для объявления ID=18")
		logger.Info().Msgf("DEBUG: Source документа: %+v", doc)
		if metadata, ok := doc["metadata"].(map[string]interface{}); ok {
			logger.Info().Msgf("DEBUG: Метаданные в документе: %+v", metadata)
			if discount, ok := metadata["discount"].(map[string]interface{}); ok {
				logger.Info().Msgf("DEBUG: Скидка в документе: %+v", discount)
			}
		}
	}
	return listing, nil
}
