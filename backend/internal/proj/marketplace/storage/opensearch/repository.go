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
	"github.com/jackc/pgx/v5"
	"log"

	"net/http"
	"net/url"
	"regexp"
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
		log.Printf("Запуск переиндексации с обновленной схемой (поддержка метаданных и скидок)")
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
	var query map[string]interface{}
	if params.CustomQuery != nil {
		// Используем переданный готовый запрос
		query = params.CustomQuery
	} else {
		// Строим запрос к OpenSearch, если CustomQuery не указан
		query = r.buildSearchQuery(params)
	}

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
	log.Printf("Результаты поиска для атрибутов %v:", params.AttributeFilters)
	for i, listing := range result.Listings {
		log.Printf("Объявление %d: ID=%d, Название=%s", i+1, listing.ID, listing.Title)

		// Добавляем отладочную информацию о атрибутах
		if len(listing.Attributes) > 0 {
			log.Printf("  Объявление %d имеет %d атрибутов:", listing.ID, len(listing.Attributes))
			for _, attr := range listing.Attributes {
				log.Printf("  Атрибут: name=%s, type=%s, value=%s",
					attr.AttributeName, attr.AttributeType, attr.DisplayValue)
			}
		} else {
			log.Printf("  Объявление %d не имеет атрибутов", listing.ID)
		}
	}
	return result, nil
}

// SuggestListings предлагает автодополнение для поиска
// Модифицировать файл: /data/hostel-booking-system/backend/internal/proj/marketplace/storage/opensearch/repository.go
// Заменить метод SuggestListings на:

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
					// Поиск по заголовку
					{
						"match_phrase_prefix": map[string]interface{}{
							"title": map[string]interface{}{
								"query":          prefix,
								"max_expansions": 10,
								"slop":           2,
							},
						},
					},
					// Поиск по описанию
					{
						"match_phrase_prefix": map[string]interface{}{
							"description": map[string]interface{}{
								"query":          prefix,
								"max_expansions": 10,
							},
						},
					},
					// Поиск по вариациям заголовка
					{
						"match_phrase_prefix": map[string]interface{}{
							"title_variations": map[string]interface{}{
								"query":          prefix,
								"max_expansions": 10,
							},
						},
					},
					// Нечеткий поиск по заголовку
					{
						"fuzzy": map[string]interface{}{
							"title": map[string]interface{}{
								"value":     prefix,
								"fuzziness": "AUTO",
							},
						},
					},
					// Поиск по атрибутам (ДОБАВЛЕНО)
					{
						"nested": map[string]interface{}{
							"path": "attributes",
							"query": map[string]interface{}{
								"bool": map[string]interface{}{
									"should": []map[string]interface{}{
										{
											"match_phrase_prefix": map[string]interface{}{
												"attributes.text_value": map[string]interface{}{
													"query":          prefix,
													"max_expansions": 10,
												},
											},
										},
										{
											"match_phrase_prefix": map[string]interface{}{
												"attributes.display_value": map[string]interface{}{
													"query":          prefix,
													"max_expansions": 10,
												},
											},
										},
									},
								},
							},
						},
					},
					// Поиск по важным атрибутам в корне документа (ДОБАВЛЕНО)
					{
						"match_phrase_prefix": map[string]interface{}{
							"make": map[string]interface{}{
								"query":          prefix,
								"max_expansions": 10,
							},
						},
					},
					{
						"match_phrase_prefix": map[string]interface{}{
							"model": map[string]interface{}{
								"query":          prefix,
								"max_expansions": 10,
							},
						},
					},
					{
						"match_phrase_prefix": map[string]interface{}{
							"brand": map[string]interface{}{
								"query":          prefix,
								"max_expansions": 10,
							},
						},
					},
					{
						"match_phrase_prefix": map[string]interface{}{
							"color": map[string]interface{}{
								"query":          prefix,
								"max_expansions": 10,
							},
						},
					},
					// Поиск по значениям атрибутов выбора (select_values)
					{
						"match_phrase_prefix": map[string]interface{}{
							"select_values": map[string]interface{}{
								"query":          prefix,
								"max_expansions": 10,
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

	// Дополнительно проверяем, не нашлось ли подсказок в атрибутах
	// и если их нет в заголовках, добавляем их
	if len(suggestions) < size {
		// Получаем значения атрибутов из результатов поиска
		attrValues := make(map[string]bool)
		if hits, ok := searchResponse["hits"].(map[string]interface{}); ok {
			if hitsArray, ok := hits["hits"].([]interface{}); ok {
				for _, hit := range hitsArray {
					if hitObj, ok := hit.(map[string]interface{}); ok {
						if source, ok := hitObj["_source"].(map[string]interface{}); ok {
							// Проверяем атрибуты make, model и другие важные
							for _, field := range []string{"make", "model", "brand", "color"} {
								if value, ok := source[field].(string); ok && value != "" {
									attrValues[value] = true
								}
							}

							// Проверяем вложенные атрибуты
							if attributes, ok := source["attributes"].([]interface{}); ok {
								for _, attrI := range attributes {
									if attr, ok := attrI.(map[string]interface{}); ok {
										// Проверяем text_value
										if textValue, ok := attr["text_value"].(string); ok && textValue != "" {
											if strings.HasPrefix(strings.ToLower(textValue), strings.ToLower(prefix)) {
												attrValues[textValue] = true
											}
										}

										// Проверяем display_value
										if displayValue, ok := attr["display_value"].(string); ok && displayValue != "" {
											if strings.HasPrefix(strings.ToLower(displayValue), strings.ToLower(prefix)) {
												attrValues[displayValue] = true
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

		// Добавляем значения атрибутов в подсказки
		for value := range attrValues {
			if !contains(suggestions, value) {
				suggestions = append(suggestions, value)
				if len(suggestions) >= size {
					break
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
			log.Printf("DEBUG: Проверка метаданных объявления %d: %+v", listing.ID, listing.Metadata)
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
func (r *Repository) getAttributeOptionTranslations(attrName, value string) (map[string]string, error) {
	// Используем фоновый контекст
	ctx := context.Background()

	// Остальной код функции
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
			log.Printf("Ошибка получения переводов для атрибута %s, значение %s: %v", attrName, value, err)
		}
		return nil, err
	}
	// Формируем карту с переводами
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

	// Логирование информации о местоположении для отладки
	log.Printf("Обработка местоположения для листинга %d: город=%s, страна=%s, адрес=%s",
		listing.ID, listing.City, listing.Country, listing.Location)
	// В методе listingToDoc добавьте следующий блок кода после обработки атрибутов
	// Примерно после строки 1660, где обрабатываются атрибуты

	// Добавляем специфические атрибуты недвижимости в корень документа для лучшего поиска
	realEstateAttrs := map[string]bool{
		"rooms":            true,
		"area":             true,
		"floor":            true,
		"total_floors":     true,
		"property_type":    true,
		"land_area":        true,
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

	// Добавляем атрибуты в специальное поле для недвижимости
	realEstateText := make([]string, 0)

	// Проверяем наличие атрибутов недвижимости и добавляем их в документ
	for _, attr := range listing.Attributes {
		if realEstateAttrs[attr.AttributeName] {
			var attrValue string

			// Получаем значение в зависимости от типа атрибута
			if attr.TextValue != nil && *attr.TextValue != "" {
				attrValue = *attr.TextValue
				doc[attr.AttributeName] = attrValue
				doc[attr.AttributeName+"_facet"] = attrValue
				realEstateText = append(realEstateText, attrValue)
			} else if attr.NumericValue != nil {
				numValue := *attr.NumericValue
				doc[attr.AttributeName] = numValue

				// Добавляем строковое представление для поиска
				strValue := fmt.Sprintf("%g", numValue)
				realEstateText = append(realEstateText, strValue)

				// Для комнат добавляем также текстовое представление (например "2 комнаты")
				if attr.AttributeName == "rooms" {
					roomsText := fmt.Sprintf("%g %s", numValue, getLocalizedRoomText(numValue))
					realEstateText = append(realEstateText, roomsText)
				}
			} else if attr.BooleanValue != nil {
				boolValue := *attr.BooleanValue
				doc[attr.AttributeName] = boolValue

				// Добавляем строковое представление для поиска
				if boolValue {
					strValue := "да"
					realEstateText = append(realEstateText, strValue)
					doc[attr.AttributeName+"_text"] = strValue
				} else {
					strValue := "нет"
					realEstateText = append(realEstateText, strValue)
					doc[attr.AttributeName+"_text"] = strValue
				}
			}

			log.Printf("Добавлен атрибут недвижимости %s в документ для объявления %d",
				attr.AttributeName, listing.ID)
		}
	}

	// Добавляем все тексты атрибутов недвижимости в специальное поле
	if len(realEstateText) > 0 {
		doc["real_estate_attributes_text"] = realEstateText

		// Объединяем в одну строку для полнотекстового поиска
		doc["real_estate_attributes_combined"] = strings.Join(realEstateText, " ")
	}

	doc["has_discount"] = listing.HasDiscount
	if listing.OldPrice > 0 {
		doc["old_price"] = listing.OldPrice
	}

	// Добавляем поля для скидок напрямую из текста описания
	if strings.Contains(listing.Description, "СКИДКА") || strings.Contains(listing.Description, "СКИДКА!") {
		// Ищем процент скидки в описании с помощью регулярного выражения
		discountRegex := regexp.MustCompile(`(\d+)%\s*СКИДКА`)
		matches := discountRegex.FindStringSubmatch(listing.Description)

		// Ищем старую цену с помощью регулярного выражения
		priceRegex := regexp.MustCompile(`Старая цена:\s*(\d+[\.,]?\d*)\s*RSD`)
		priceMatches := priceRegex.FindStringSubmatch(listing.Description)

		if len(matches) > 1 && len(priceMatches) > 1 {
			discountPercent, _ := strconv.Atoi(matches[1])
			oldPriceStr := strings.Replace(priceMatches[1], ",", ".", -1)
			oldPrice, _ := strconv.ParseFloat(oldPriceStr, 64)

			// Если у объекта нет метаданных, создаем их
			if listing.Metadata == nil {
				listing.Metadata = make(map[string]interface{})
			}

			// Создаем информацию о скидке
			discount := map[string]interface{}{
				"discount_percent":  discountPercent,
				"previous_price":    oldPrice,
				"effective_from":    time.Now().AddDate(0, 0, -10).Format(time.RFC3339), // Примерная дата начала скидки
				"has_price_history": true,
			}

			// Добавляем информацию о скидке в метаданные
			listing.Metadata["discount"] = discount

			// Добавляем OldPrice в документ
			doc["old_price"] = oldPrice
			doc["has_discount"] = true

			log.Printf("Извлечена скидка из описания для объявления %d: %v", listing.ID, discount)
		}
	}

	// Добавляем метаданные если они есть
	if listing.Metadata != nil {
		// Копируем метаданные
		doc["metadata"] = listing.Metadata

		// Если есть информация о скидке, проверяем и пересчитываем процент
		if discount, ok := listing.Metadata["discount"].(map[string]interface{}); ok {
			if prevPrice, ok := discount["previous_price"].(float64); ok && prevPrice > 0 {
				// Пересчитываем актуальный процент скидки
				if prevPrice > listing.Price {
					discountPercent := int((prevPrice - listing.Price) / prevPrice * 100)
					discount["discount_percent"] = discountPercent

					// Обновляем метаданные в документе
					listing.Metadata["discount"] = discount
					doc["metadata"] = listing.Metadata

					doc["has_discount"] = true
					doc["old_price"] = prevPrice

					log.Printf("Пересчитан процент скидки для OpenSearch: %d%% (объявление %d)",
						discountPercent, listing.ID)
				}
			}
		}
	}

	// Пытаемся получить и применить адресную информацию из витрины если она отсутствует в объявлении
	if listing.StorefrontID != nil && *listing.StorefrontID > 0 {
		needStorefrontInfo := false

		// Проверяем, нужно ли получать информацию о витрине
		if listing.City == "" || listing.Country == "" || listing.Location == "" ||
			listing.Latitude == nil || listing.Longitude == nil {
			needStorefrontInfo = true
		}

		if needStorefrontInfo {
			log.Printf("Получаем данные о витрине %d для заполнения адресной информации для объявления %d",
				*listing.StorefrontID, listing.ID)

			// Получаем информацию о витрине для заполнения адресных данных
			var storefront models.Storefront
			err := r.storage.QueryRow(context.Background(), `
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
				// Применяем информацию о витрине
				if listing.City == "" && storefront.City != "" {
					doc["city"] = storefront.City
					log.Printf("Применен город витрины '%s' для объявления %d", storefront.City, listing.ID)
				}

				if listing.Country == "" && storefront.Country != "" {
					doc["country"] = storefront.Country
					log.Printf("Применена страна витрины '%s' для объявления %d", storefront.Country, listing.ID)
				}

				if listing.Location == "" && storefront.Address != "" {
					doc["location"] = storefront.Address
					log.Printf("Применен адрес витрины '%s' для объявления %d", storefront.Address, listing.ID)
				}

				// Если у объявления нет координат, но они есть у витрины
				if (listing.Latitude == nil || listing.Longitude == nil ||
					*listing.Latitude == 0 || *listing.Longitude == 0) &&
					storefront.Latitude != nil && storefront.Longitude != nil &&
					*storefront.Latitude != 0 && *storefront.Longitude != 0 {

					// Создаем объект с полями lat и lon для geo_point
					doc["coordinates"] = map[string]interface{}{
						"lat": *storefront.Latitude,
						"lon": *storefront.Longitude,
					}

					log.Printf("Применены координаты витрины для объявления %d: lat=%f, lon=%f",
						listing.ID, *storefront.Latitude, *storefront.Longitude)

					// Устанавливаем показ на карте для координат витрины
					doc["show_on_map"] = true
				}
			} else {
				log.Printf("Не удалось получить информацию о витрине %d: %v", *listing.StorefrontID, err)
			}
		}
	}

	// Добавляем координаты, если они есть в объявлении
	if listing.Latitude != nil && listing.Longitude != nil && *listing.Latitude != 0 && *listing.Longitude != 0 {
		// Создаем объект с полями lat и lon для geo_point
		doc["coordinates"] = map[string]interface{}{
			"lat": *listing.Latitude,
			"lon": *listing.Longitude,
		}

		log.Printf("Добавлены координаты из объявления для листинга %d: lat=%f, lon=%f",
			listing.ID, *listing.Latitude, *listing.Longitude)
	} else if _, ok := doc["coordinates"]; !ok { // Если координаты все еще не добавлены
		// Если координаты отсутствуют, но есть город, попробуем геокодировать
		if cityVal, ok := doc["city"].(string); ok && cityVal != "" {
			countryVal := ""
			if c, ok := doc["country"].(string); ok {
				countryVal = c
			}

			geocoded, err := r.geocodeCity(cityVal, countryVal)
			if err == nil && geocoded != nil {
				// Добавляем найденные координаты
				doc["coordinates"] = map[string]interface{}{
					"lat": geocoded.Lat,
					"lon": geocoded.Lon,
				}
				log.Printf("Добавлены геокодированные координаты для листинга %d (город %s): lat=%f, lon=%f",
					listing.ID, cityVal, geocoded.Lat, geocoded.Lon)

				// Устанавливаем показ на карте для геокодированных координат
				doc["show_on_map"] = true
			} else {
				log.Printf("Не удалось геокодировать город %s для листинга %d: %v",
					cityVal, listing.ID, err)
			}
		} else {
			log.Printf("У листинга %d нет координат и не указан город", listing.ID)
		}
	}

	// Добавляем storefront_id, если есть
	if listing.StorefrontID != nil && *listing.StorefrontID > 0 {
		doc["storefront_id"] = *listing.StorefrontID
		log.Printf("Добавлен storefront_id: %d для объявления %d в индекс", *listing.StorefrontID, listing.ID)
	} else {
		log.Printf("Нет storefront_id для объявления %d или значение некорректно", listing.ID)
	}

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

		for i, img := range listing.Images {
			// Логирование деталей каждого изображения для отладки
			log.Printf("  Изображение %d: ID=%d, Путь=%s, IsMain=%v",
				i+1, img.ID, img.FilePath, img.IsMain)

			imagesDoc = append(imagesDoc, map[string]interface{}{
				"id":        img.ID,
				"file_path": img.FilePath,
				"is_main":   img.IsMain,
			})
		}

		// Проверяем, что у нас есть хотя бы одно изображение с указанным путем
		hasValidImage := false
		for _, img := range imagesDoc {
			if path, ok := img["file_path"].(string); ok && path != "" {
				hasValidImage = true
				break
			}
		}

		if hasValidImage {
			doc["images"] = imagesDoc
			log.Printf("  Добавлено %d изображений в индекс", len(imagesDoc))
		} else {
			log.Printf("  ВНИМАНИЕ: У объявления %d нет изображений с корректным путем", listing.ID)
		}
	} else {
		// Пытаемся загрузить изображения из базы данных, если их нет в объекте
		images, err := r.storage.GetListingImages(context.Background(), fmt.Sprintf("%d", listing.ID))
		if err != nil {
			log.Printf("  Ошибка при загрузке изображений для объявления %d: %v", listing.ID, err)
		} else if len(images) > 0 {
			log.Printf("  Загружено %d изображений из базы данных для объявления %d", len(images), listing.ID)

			imagesDoc := make([]map[string]interface{}, 0, len(images))
			for i, img := range images {
				log.Printf("    Изображение %d: ID=%d, Путь=%s, IsMain=%v",
					i+1, img.ID, img.FilePath, img.IsMain)

				imagesDoc = append(imagesDoc, map[string]interface{}{
					"id":        img.ID,
					"file_path": img.FilePath,
					"is_main":   img.IsMain,
				})
			}

			doc["images"] = imagesDoc
			log.Printf("  Добавлено %d изображений из базы данных в индекс", len(imagesDoc))
		} else {
			log.Printf("  Объявление %d не имеет изображений", listing.ID)
		}
	}

	// Изменить раздел с атрибутами для добавления атрибутов в корень документа
	if listing.Attributes != nil && len(listing.Attributes) > 0 {
		attributes := make([]map[string]interface{}, 0, len(listing.Attributes))

		// Создаем карту для отслеживания добавленных атрибутов
		addedAttributes := make(map[string]bool)

		// Используем map для отслеживания уникальных значений атрибутов
		uniqueTextValues := make(map[string]bool)
		attributeTextValues := make(map[string][]string)

		// Массив для select_values
		selectValues := []string{}

		// Явно определяем переменные для марки и модели автомобиля
		var makeValue, modelValue string

		for _, attr := range listing.Attributes {
			// Проверяем, что атрибут имеет хотя бы одно значение
			if attr.TextValue == nil && attr.NumericValue == nil &&
				attr.BooleanValue == nil && attr.JSONValue == nil &&
				attr.DisplayValue == "" {
				continue
			}

			// Извлекаем значения make и model для последующего использования
			if attr.AttributeName == "make" && attr.TextValue != nil && *attr.TextValue != "" {
				makeValue = *attr.TextValue
			}
			if attr.AttributeName == "model" && attr.TextValue != nil && *attr.TextValue != "" {
				modelValue = *attr.TextValue
			}

			// Для атрибутов select и text, явно добавляем их в корень документа
			if attr.AttributeType == "select" || attr.AttributeType == "text" {
				var attrValue string
				if attr.TextValue != nil && *attr.TextValue != "" {
					attrValue = *attr.TextValue
				} else if attr.DisplayValue != "" {
					attrValue = attr.DisplayValue
				}

				if attrValue != "" {
					// Добавляем в корень документа, независимо от важности
					doc[attr.AttributeName] = attrValue
					doc[attr.AttributeName+"_lowercase"] = strings.ToLower(attrValue)

					// Улучшаем поиск для select-значений, добавляя варианты
					if attr.AttributeType == "select" {
						if doc["select_values"] == nil {
							doc["select_values"] = []string{}
						}
						selectValues := doc["select_values"].([]string)
						doc["select_values"] = append(selectValues, attrValue, strings.ToLower(attrValue))
					}

					// Ведем лог для проверки
					log.Printf("Добавлен атрибут %s = %s в корень документа для объявления %d",
						attr.AttributeName, attrValue, listing.ID)
				}
			}

			// Для избежания дубликатов атрибутов с одинаковым именем
			attrKey := fmt.Sprintf("%s_%d", attr.AttributeName, attr.AttributeID)
			if addedAttributes[attrKey] {
				continue
			}
			addedAttributes[attrKey] = true

			attrDoc := map[string]interface{}{
				"attribute_id":   attr.AttributeID,
				"attribute_name": attr.AttributeName,
				"display_name":   attr.DisplayName,
				"attribute_type": attr.AttributeType,
				"display_value":  attr.DisplayValue,
			}

			// Добавляем типизированные значения в зависимости от типа атрибута
			if attr.TextValue != nil && *attr.TextValue != "" {
				attrDoc["text_value"] = *attr.TextValue
				attrDoc["text_value_lowercase"] = strings.ToLower(*attr.TextValue)

				// Для uniqueness tracking
				textValue := *attr.TextValue
				lowerTextValue := strings.ToLower(textValue)

				// Добавляем уникальные значения в attributeTextValues
				if _, ok := attributeTextValues[attr.AttributeName]; !ok {
					attributeTextValues[attr.AttributeName] = []string{}
				}

				// Проверяем уникальность перед добавлением
				if !uniqueTextValues[textValue] {
					attributeTextValues[attr.AttributeName] = append(attributeTextValues[attr.AttributeName], textValue)
					uniqueTextValues[textValue] = true
				}

				if !uniqueTextValues[lowerTextValue] {
					attributeTextValues[attr.AttributeName] = append(attributeTextValues[attr.AttributeName], lowerTextValue)
					uniqueTextValues[lowerTextValue] = true
				}

				// Добавляем в корень документа
				if isImportantAttribute(attr.AttributeName) {
					doc[attr.AttributeName] = textValue
					doc[attr.AttributeName+"_facet"] = textValue
					doc[attr.AttributeName+"_lowercase"] = lowerTextValue
				}

				// Добавляем в select_values если это select атрибут
				if attr.AttributeType == "select" {
					selectValues = append(selectValues, textValue, lowerTextValue)
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

							// Добавляем значения для поиска
							if _, ok := attributeTextValues[attr.AttributeName]; !ok {
								attributeTextValues[attr.AttributeName] = []string{}
							}
							attributeTextValues[attr.AttributeName] = append(
								attributeTextValues[attr.AttributeName],
								strArray...,
							)
						}
					}
				}
			}

			// Получаем переводы значений атрибутов, если это select или другое поле с переводимыми значениями
			if attr.AttributeType == "select" && (attr.TextValue != nil || attr.DisplayValue != "") {
				// Запрашиваем переводы из таблицы attribute_option_translations
				value := ""
				if attr.TextValue != nil {
					value = *attr.TextValue
				} else {
					value = attr.DisplayValue
				}

				if value != "" {
					translations, err := r.getAttributeOptionTranslations(attr.AttributeName, value)
					if err == nil && len(translations) > 0 {
						attrDoc["translations"] = translations

						// Добавляем переведенные значения для поиска
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

			attributes = append(attributes, attrDoc)
		}

		// Явно добавляем марку и модель в корень документа, если они были найдены
		if makeValue != "" {
			doc["make"] = makeValue
			doc["make_lowercase"] = strings.ToLower(makeValue)
			log.Printf("Явно добавлена марка автомобиля '%s' в корень документа для объявления %d",
				makeValue, listing.ID)
		}

		if modelValue != "" {
			doc["model"] = modelValue
			doc["model_lowercase"] = strings.ToLower(modelValue)
			log.Printf("Явно добавлена модель автомобиля '%s' в корень документа для объявления %d",
				modelValue, listing.ID)
		}

		// Добавляем select_values в документ
		if len(selectValues) > 0 {
			// Удаляем дубликаты из selectValues
			uniqueSelectValues := make([]string, 0)
			selectValuesMap := make(map[string]bool)

			for _, val := range selectValues {
				if !selectValuesMap[val] {
					selectValuesMap[val] = true
					uniqueSelectValues = append(uniqueSelectValues, val)
				}
			}

			doc["select_values"] = uniqueSelectValues
		}

		// Добавляем атрибуты в документ
		doc["attributes"] = attributes

		// Добавляем all_attributes_text без дублирования
		allAttrsSet := make(map[string]bool)
		allAttrs := []string{}

		for _, values := range attributeTextValues {
			for _, value := range values {
				if !allAttrsSet[value] {
					allAttrsSet[value] = true
					allAttrs = append(allAttrs, value)
				}
			}
		}

		doc["all_attributes_text"] = allAttrs
	}
	return doc
}
func getLocalizedRoomText(rooms float64) string {
    // По умолчанию используем русский язык
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
	// Список важных атрибутов, которые должны быть добавлены в корень документа
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

// geocodeCity получает координаты города
func (r *Repository) geocodeCity(city, country string) (*struct{ Lat, Lon float64 }, error) {
	// Формируем запрос для геокодирования
	query := city
	if country != "" {
		query += ", " + country
	}

	// Используем OSM Nominatim API для получения координат
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

	// OSM требует User-Agent
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
	log.Printf("Строим запрос: категория = %v, язык = %s, поисковый запрос = %s",
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

	// Добавляем фильтр по статусу "active" по умолчанию, если не указан явно другой статус
	if params.Status == "" {
		filter := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]interface{})
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"status": "active",
			},
		})
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = filter
	}

	// Изменить раздел с текстовым поиском в функции buildSearchQuery
	if params.Query != "" {
		log.Printf("Текстовый поиск по запросу: '%s'", params.Query)

		// Используем все поля переводов и атрибутов
		searchFields := []string{
			"title^3", "description",
			"title.sr^4", "description.sr", "translations.sr.title^4", "translations.sr.description",
			"title.ru^4", "description.ru", "translations.ru.title^4", "translations.ru.description",
			"title.en^4", "description.en", "translations.en.title^4", "translations.en.description",
			"all_attributes_text^2",
			// Специальные поля для важных атрибутов (расширенный список)
			"make^5", // Увеличиваем вес поля make для лучшего поиска по марке
			"model^4",
			"color^3",
			"brand^4",
			"year^3",
			"rooms^3",
			"property_type^3",
			"body_type^3",
			"cpu^3",
			"gpu^3",
			"memory^3",
			"ram^3",
			"storage_capacity^3",
			"screen_size^3",
			"attr_make_text^5", // Увеличиваем вес для атрибутов make
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
			// Поля для атрибутов недвижимости
			"real_estate_attributes_text^3",
			"real_estate_attributes_combined^3",
			"rooms^4",
			"area^3",
			"floor^3",
			"total_floors^3",
			"property_type^4",
			"land_area^3",
			"year_built^3",
			"bathrooms^3",
			"heating_type^3",
			"parking^3",
			"balcony^3",
			"furnished^3",
			"air_conditioning^3",
			// Поля для текстовых представлений
			"rooms_text^4",
			"property_type_text^4",
			"heating_type_text^3",
			"parking_text^3",
			"furnished_text^3",
		}

		// Определяем язык запроса для приоритизации полей
		languagePriority := "sr" // По умолчанию сербский
		if params.Language != "" {
			languagePriority = params.Language
		}

		// Увеличиваем вес полей для выбранного языка
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

		minimumShouldMatch := "30%" // Снижаем еще больше с 60% до 30% для более гибкого поиска
		if params.MinimumShouldMatch != "" {
			minimumShouldMatch = params.MinimumShouldMatch
		}

		fuzziness := "AUTO"
		if params.Fuzziness != "" {
			fuzziness = params.Fuzziness
		}

		// Добавляем прямые запросы в should-блок для основных полей
		boolMap := query["query"].(map[string]interface{})["bool"].(map[string]interface{})
		should := boolMap["should"].([]interface{})

		// Добавляем запрос для поиска по заголовку с высоким весом и нечетким поиском
		should = append(should, map[string]interface{}{
			"match": map[string]interface{}{
				"title": map[string]interface{}{
					"query":     params.Query,
					"boost":     5.0,
					"fuzziness": fuzziness,
				},
			},
		})

		// Добавляем запрос для поиска по описанию
		should = append(should, map[string]interface{}{
			"match": map[string]interface{}{
				"description": map[string]interface{}{
					"query":     params.Query,
					"boost":     2.0,
					"fuzziness": fuzziness,
				},
			},
		})

		// Добавляем запрос для переводов заголовка
		should = append(should, map[string]interface{}{
			"match": map[string]interface{}{
				"translations.sr.title": map[string]interface{}{
					"query":     params.Query,
					"boost":     4.0,
					"fuzziness": fuzziness,
				},
			},
		})

		// Добавляем запрос для переводов описания
		should = append(should, map[string]interface{}{
			"match": map[string]interface{}{
				"translations.sr.description": map[string]interface{}{
					"query":     params.Query,
					"boost":     1.5,
					"fuzziness": fuzziness,
				},
			},
		})

		// Проверяем, содержит ли запрос ключевые слова недвижимости
		realEstateKeywords := []string{
			"квартира", "комната", "комнат", "дом", "этаж",
			"площадь", "м2", "кв.м", "квм", "кв м",
			"студия", "однокомнатная", "двухкомнатная", "трехкомнатная", 
			"однушка", "двушка", "трешка", "участок", "сотка",
			"гараж", "паркинг", "балкон", "лоджия", "терраса",
			"ремонт", "новостройка", "вторичка", "жилье", "недвижимость",
			"аренда", "съем", "снять", "купить", "продажа",
		}

		isRealEstateQuery := false
		for _, keyword := range realEstateKeywords {
			if strings.Contains(strings.ToLower(params.Query), keyword) {
				isRealEstateQuery = true
				break
			}
		}

		// Если запрос связан с недвижимостью, увеличиваем значимость соответствующих полей
		if isRealEstateQuery {
			log.Printf("Обнаружен запрос о недвижимости: '%s'", params.Query)
			
			// Добавляем дополнительные поля в запрос с повышенным весом
			realEstateBoost := []map[string]interface{}{
				{
					"match": map[string]interface{}{
						"real_estate_attributes_combined": map[string]interface{}{
							"query":     params.Query,
							"boost":     5.0,
							"fuzziness": fuzziness,
						},
					},
				},
				{
					"match": map[string]interface{}{
						"property_type": map[string]interface{}{
							"query":     params.Query,
							"boost":     4.0,
							"fuzziness": fuzziness,
						},
					},
				},
				{
					"match": map[string]interface{}{
						"rooms_text": map[string]interface{}{
							"query":     params.Query,
							"boost":     4.0,
							"fuzziness": fuzziness,
						},
					},
				},
			}
			
			// Добавляем эти запросы в should-блок основного запроса
			for _, q := range realEstateBoost {
				should = append(should, q)
			}
		}

		// Обновляем should-блок
		boolMap["should"] = should
		boolMap["minimum_should_match"] = 1 // Достаточно соответствия одному полю

		// Создаем общий multi_match запрос для всех полей
		multiMatch := map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":                params.Query,
				"fields":               searchFields,
				"type":                 "best_fields",
				"fuzziness":            fuzziness,
				"operator":             "OR", // Используем OR вместо AND для более гибкого поиска
				"minimum_should_match": minimumShouldMatch,
			},
		}

		// Добавляем multi_match в must-блок
		must := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]interface{})
		must = append(must, multiMatch)
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = must

		// Обрабатываем случай многословного запроса
		words := strings.Fields(params.Query)
		if len(words) > 1 {
			// Для каждого слова создаем отдельные поисковые запросы
			for _, word := range words {
				if len(word) < 2 {
					continue // Пропускаем слишком короткие слова
				}

				// Добавляем запрос для поиска по заголовку только с этим словом
				should = append(should, map[string]interface{}{
					"match": map[string]interface{}{
						"title": map[string]interface{}{
							"query":     word,
							"boost":     2.0,
							"fuzziness": fuzziness,
						},
					},
				})

				// Добавляем запрос для поиска по описанию только с этим словом
				should = append(should, map[string]interface{}{
					"match": map[string]interface{}{
						"description": map[string]interface{}{
							"query":     word,
							"boost":     1.0,
							"fuzziness": fuzziness,
						},
					},
				})

				// ВАЖНО: Убираем nested запросы, т.к. они вызывают ошибку
				// Вместо этого добавляем поиск по всем атрибутам в корне документа
				should = append(should, map[string]interface{}{
					"match": map[string]interface{}{
						"all_attributes_text": map[string]interface{}{
							"query":     word,
							"boost":     2.0,
							"fuzziness": fuzziness,
						},
					},
				})
				
				// Добавляем поиск по специфическим полям недвижимости для каждого слова
				if isRealEstateQuery {
					should = append(should, map[string]interface{}{
						"match": map[string]interface{}{
							"real_estate_attributes_combined": map[string]interface{}{
								"query":     word,
								"boost":     3.0,
								"fuzziness": fuzziness,
							},
						},
					})
					
					// Поиск по полю комнат
					should = append(should, map[string]interface{}{
						"match": map[string]interface{}{
							"rooms_text": map[string]interface{}{
								"query":     word,
								"boost":     2.5,
								"fuzziness": fuzziness,
							},
						},
					})
				}
			}

			// Обновляем should-блок с новыми запросами
			boolMap["should"] = should
		}

		// Добавляем прямые запросы к полям make и model в корне документа
		shouldQueries := []map[string]interface{}{
			// Ищем по корневому полю make
			{
				"match": map[string]interface{}{
					"make": map[string]interface{}{
						"query":     params.Query,
						"boost":     5.0,
						"fuzziness": "AUTO",
					},
				},
			},
			{
				"match": map[string]interface{}{
					"make_lowercase": map[string]interface{}{
						"query":     strings.ToLower(params.Query),
						"boost":     5.0,
						"fuzziness": "AUTO",
					},
				},
			},
			// Ищем по корневому полю model
			{
				"match": map[string]interface{}{
					"model": map[string]interface{}{
						"query":     params.Query,
						"boost":     4.0,
						"fuzziness": "AUTO",
					},
				},
			},
			{
				"match": map[string]interface{}{
					"model_lowercase": map[string]interface{}{
						"query":     strings.ToLower(params.Query),
						"boost":     4.0,
						"fuzziness": "AUTO",
					},
				},
			},
			// Ищем по полю select_values, которое содержит все select-значения
			{
				"match": map[string]interface{}{
					"select_values": map[string]interface{}{
						"query":     params.Query,
						"boost":     3.0,
						"fuzziness": "AUTO",
					},
				},
			},
		}

		// Добавим эти запросы в should-блок основного запроса
		should = boolMap["should"].([]interface{})
		for _, q := range shouldQueries {
			should = append(should, q)
		}
		boolMap["should"] = should
	}

	// Добавляем фильтры категории, если указано
	if params.CategoryID != nil && *params.CategoryID > 0 {
		// Поиск по точному ID категории или родительской категории
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

	// Добавляем фильтр по цене, если указано
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

	// Добавляем фильтр по состоянию товара (новый, б/у)
	if params.Condition != "" {
		filter := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]interface{})
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"condition": params.Condition,
			},
		})
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = filter
	}

	// Добавляем фильтр по городу и стране
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

	// Добавляем фильтр по витрине
	if params.StorefrontID != nil && *params.StorefrontID > 0 {
		filter := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]interface{})
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"storefront_id": *params.StorefrontID,
			},
		})
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = filter
	}

	// Добавляем фильтр по статусу
	if params.Status != "" {
		filter := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]interface{})
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"status": params.Status,
			},
		})
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = filter
	}

	// Добавляем гео-фильтр для поиска по расстоянию
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

	// Добавляем фильтры по атрибутам, если есть
	if params.AttributeFilters != nil && len(params.AttributeFilters) > 0 {
		filter := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]interface{})
		
		// Список атрибутов недвижимости
		realEstateAttrs := map[string]bool{
			"property_type": true, "rooms": true, "floor": true, "total_floors": true,
			"area": true, "land_area": true, "building_type": true, 
			"has_balcony": true, "has_elevator": true, "has_parking": true,
			"year_built": true, "bathrooms": true,
		}
		
		for attrName, attrValue := range params.AttributeFilters {
			if attrValue == "" {
				continue
			}
			
			// Проверяем, является ли атрибут атрибутом недвижимости
			if realEstateAttrs[attrName] {
				// Создаем простой term-запрос вместо nested для атрибутов недвижимости
				if attrValue == "true" || attrValue == "false" {
					// Для булевых атрибутов
					boolVal := attrValue == "true"
					filter = append(filter, map[string]interface{}{
						"term": map[string]interface{}{
							attrName: boolVal,
						},
					})
					log.Printf("Добавлен фильтр по атрибуту недвижимости (boolean): %s = %v", attrName, boolVal)
				} else if strings.Contains(attrValue, ",") {
					// Для числовых диапазонов
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
							log.Printf("Добавлен фильтр по атрибуту недвижимости (range): %s = %v-%v", attrName, minVal, maxVal)
						}
					}
				} else {
					// Для текстовых атрибутов
					filter = append(filter, map[string]interface{}{
						"term": map[string]interface{}{
							attrName: attrValue,
						},
					})
					log.Printf("Добавлен фильтр по атрибуту недвижимости (text): %s = %s", attrName, attrValue)
				}
			} else {
				// Для обычных атрибутов, используем term-запрос по всем типам значений
				filter = append(filter, map[string]interface{}{
					"term": map[string]interface{}{
						"attr_" + attrName: attrValue,
					},
				})
				log.Printf("Добавлен фильтр по обычному атрибуту: %s = %s", attrName, attrValue)
			}
		}
		
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = filter
	}

	// Настраиваем сортировку
	if params.Sort != "" {
		var sortField string
		var sortOrder string

		if params.SortDirection != "" {
			sortOrder = params.SortDirection
		} else {
			sortOrder = "desc" // По умолчанию по убыванию
		}

		switch params.Sort {
		case "date_desc":
			sortField = "created_at"
			sortOrder = "desc"
		case "date_asc":
			sortField = "created_at"
			sortOrder = "asc"
		case "price_desc":
			sortField = "price"
			sortOrder = "desc"
		case "price_asc":
			sortField = "price"
			sortOrder = "asc"
		case "rating_desc":
			log.Printf("Применяем сортировку рейтинга по УБЫВАНИЮ")
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
			log.Printf("Применяем сортировку рейтинга по ВОЗРАСТАНИЮ")
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
			// Пытаемся использовать указанное поле сортировки напрямую
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

		log.Printf("Применяем сортировку по полю %s в порядке %s", sortField, sortOrder)
		query["sort"] = []interface{}{
			map[string]interface{}{
				sortField: map[string]interface{}{
					"order": sortOrder,
				},
			},
		}
	} else {
		// Сортировка по умолчанию - по дате создания
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
						if avgRating, ok := source["average_rating"].(float64); ok {
							listing.AverageRating = avgRating
						}

						if reviewCount, ok := source["review_count"].(float64); ok {
							listing.ReviewCount = int(reviewCount)
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
				log.Printf("Найдена скидка для объявления %d: скидка %v%%, старая цена: %.2f",
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
		log.Printf("DEBUG: Преобразование документа для объявления ID=18")
		log.Printf("DEBUG: Source документа: %+v", doc)
		if metadata, ok := doc["metadata"].(map[string]interface{}); ok {
			log.Printf("DEBUG: Метаданные в документе: %+v", metadata)
			if discount, ok := metadata["discount"].(map[string]interface{}); ok {
				log.Printf("DEBUG: Скидка в документе: %+v", discount)
			}
		}
	}
	return listing, nil
}
