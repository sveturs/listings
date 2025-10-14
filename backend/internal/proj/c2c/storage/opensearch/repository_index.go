// backend/internal/proj/c2c/storage/opensearch/repository_index.go
package opensearch

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"backend/internal/domain/models"
	"backend/internal/logger"
	"backend/internal/storage"
	osClient "backend/internal/storage/opensearch"
)

func (r *Repository) IndexListing(ctx context.Context, listing *models.MarketplaceListing) error {
	if listing == nil {
		logger.Error().Msg("ERROR: Попытка индексации nil объявления")
		return fmt.Errorf("listing is nil")
	}

	if listing.ID == 0 {
		logger.Error().Msg("ERROR: Попытка индексации объявления с ID=0")
		return fmt.Errorf("listing ID is 0")
	}

	doc := r.listingToDoc(ctx, listing)
	docJSON, _ := json.MarshalIndent(doc, "", "  ")
	logger.Info().Msgf("INFO: Начинаем индексацию объявления ID=%d в индекс '%s'", listing.ID, r.indexName)
	logger.Debug().Msgf("DEBUG: Данные для индексации: %s", string(docJSON))

	err := r.client.IndexDocument(ctx, r.indexName, fmt.Sprintf("%d", listing.ID), doc)
	if err != nil {
		logger.Error().Msgf("ERROR: Ошибка индексации объявления ID=%d в OpenSearch: %v", listing.ID, err)
		return err
	}

	logger.Info().Msgf("SUCCESS: Объявление ID=%d успешно проиндексировано в OpenSearch", listing.ID)
	return nil
}

// BulkIndexListings индексирует несколько объявлений
func (r *Repository) BulkIndexListings(ctx context.Context, listings []*models.MarketplaceListing) error {
	docs := make([]map[string]interface{}, 0, len(listings))

	for _, listing := range listings {
		doc := r.listingToDoc(ctx, listing)
		logger.Info().Msgf("Индексация объявления с ID: %d, категория: %d, название: %s",
			listing.ID, listing.CategoryID, listing.Title)

		if listing.ID == 0 {
			logger.Info().Msgf("ВНИМАНИЕ: Объявление с нулевым ID: %s (категория: %d)",
				listing.Title, listing.CategoryID)
		}

		doc["id"] = listing.ID
		docs = append(docs, doc)
	}

	return r.client.BulkIndex(ctx, r.indexName, docs)
}

// DeleteListing удаляет объявление из индекса
func (r *Repository) DeleteListing(ctx context.Context, listingID string) error {
	return r.client.DeleteDocument(ctx, r.indexName, listingID)
}

func (r *Repository) ReindexAll(ctx context.Context) error {
	exists, err := r.client.IndexExists(ctx, r.indexName)
	if err != nil {
		return fmt.Errorf("ошибка проверки индекса: %w", err)
	}

	if exists {
		logger.Info().Msgf("Удаляем существующий индекс %s", r.indexName)
		if err := r.client.DeleteIndex(ctx, r.indexName); err != nil {
			return fmt.Errorf("ошибка удаления индекса: %w", err)
		}
		time.Sleep(1 * time.Second)
	}

	logger.Info().Msgf("Создаем индекс %s заново", r.indexName)
	if err := r.client.CreateIndex(ctx, r.indexName, osClient.ListingMapping); err != nil {
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
			if len(listings[i].Translations) == 0 {
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
			if len(listings[i].Attributes) == 0 {
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

func (r *Repository) listingToDoc(ctx context.Context, listing *models.MarketplaceListing) map[string]interface{} {
	doc := map[string]interface{}{
		"id":               listing.ID,
		"type":             "listing", // Для фильтрации marketplace vs storefront
		"document_type":    "listing", // Для обратной совместимости
		"title":            listing.Title,
		"description":      listing.Description,
		"title_suggest":    listing.Title,
		"title_variations": []string{listing.Title, strings.ToLower(listing.Title)},
		fieldNamePrice:     listing.Price,
		"condition":        listing.Condition,
		"status":           listing.Status,
		// НЕ отправляем listing.Location - поле "location" в OpenSearch это geo_point!
		"city":                 listing.City,
		"country":              listing.Country,
		"address_multilingual": listing.AddressMultilingual,
		"views_count":          listing.ViewsCount,
		fieldNameCreatedAt:     listing.CreatedAt.Format(time.RFC3339),
		"updated_at":           listing.UpdatedAt.Format(time.RFC3339),
		"show_on_map":          listing.ShowOnMap,
		"original_language":    listing.OriginalLanguage,
		"category_id":          listing.CategoryID,
		"user_id":              listing.UserID,
		"translations":         listing.Translations,
		"average_rating":       listing.AverageRating,
		"review_count":         listing.ReviewCount,
	}

	// Загружаем переводы из таблицы translations и преобразуем в правильный формат
	dbTranslations, err := r.getListingTranslationsFromDB(ctx, listing.ID)
	if err != nil {
		logger.Error().Msgf("Ошибка загрузки переводов для объявления %d: %v", listing.ID, err)
	} else if len(dbTranslations) > 0 {
		// Преобразуем []DBTranslation в map[язык]map[поле]значение
		translationsMap := r.convertDBTranslationsToMap(dbTranslations)
		doc["translations"] = translationsMap
		doc["supported_languages"] = r.extractSupportedLanguages(dbTranslations)
		logger.Info().Msgf("Загружено %d переводов из БД для объявления %d, преобразовано в структуру translations", len(dbTranslations), listing.ID)
	}

	logger.Info().Msgf("Обработка местоположения для листинга %d: город=%s, страна=%s, адрес=%s",
		listing.ID, listing.City, listing.Country, listing.Location)

	realEstateFields := createRealEstateFieldsMap()
	carFields := createCarFieldsMap()
	importantAttrs := createImportantAttributesMap()

	if len(listing.Attributes) > 0 {
		processAttributesForIndex(ctx, doc, listing.Attributes, importantAttrs, realEstateFields, carFields, listing.ID, r)
	}

	processDiscountData(doc, listing)
	processMetadata(doc, listing)
	processStorefrontData(doc, listing, r.storage) //nolint:contextcheck
	processCoordinates(doc, listing, r)
	processCategoryPath(doc, listing, r.storage) //nolint:contextcheck
	processCategory(doc, listing)
	processUser(doc, listing)
	processImages(doc, listing, r.storage) //nolint:contextcheck

	docJSON, _ := json.MarshalIndent(doc, "", "  ")
	logger.Info().Msgf("FINAL DOC for listing %d [size=%d bytes]: %s", listing.ID, len(docJSON), string(docJSON))

	return doc
}

func processAttributesForIndex(ctx context.Context, doc map[string]interface{}, attributes []models.ListingAttributeValue,
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
				// logger.Debug().Msgf("FIRST PASS: Добавлена марка '%s' в корень документа для объявления %d", makeValue, listingID)
			case "model":
				modelValue = textValue
				doc["model"] = modelValue
				doc["model_lowercase"] = strings.ToLower(modelValue)
				carKeywords = append(carKeywords, textValue, strings.ToLower(textValue)) // Добавляем к ключевым словам
				// logger.Debug().Msgf("FIRST PASS: Добавлена модель '%s' в корень документа для объявления %d", modelValue, listingID)
			default:
				if isImportantTextAttribute(attr.AttributeName) {
					doc[attr.AttributeName] = textValue
					doc[attr.AttributeName+"_lowercase"] = strings.ToLower(textValue)
					// logger.Debug().Msgf("FIRST PASS: Добавлен важный атрибут %s = '%s' в корень документа для объявления %d",
					//	attr.AttributeName, textValue, listingID)
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
					translations, err := r.getAttributeOptionTranslations(ctx, attr.AttributeName, value)
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

func processDiscountData(doc map[string]interface{}, listing *models.MarketplaceListing) {
	doc["has_discount"] = listing.HasDiscount
	if listing.OldPrice != nil && *listing.OldPrice > 0 {
		doc["old_price"] = *listing.OldPrice
	}

	if strings.Contains(listing.Description, "СКИДКА") || strings.Contains(listing.Description, "СКИДКА!") {
		discountRegex := regexp.MustCompile(`(\d+)%\s*СКИДКА`)
		matches := discountRegex.FindStringSubmatch(listing.Description)
		priceRegex := regexp.MustCompile(`Старая цена:\s*(\d+[\.,]?\d*)\s*RSD`)
		priceMatches := priceRegex.FindStringSubmatch(listing.Description)

		if len(matches) > 1 && len(priceMatches) > 1 {
			discountPercent, _ := strconv.Atoi(matches[1])
			oldPriceStr := strings.ReplaceAll(priceMatches[1], ",", ".")
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

		// Всегда загружаем информацию о витрине для индексации
		logger.Info().Msgf("Fetching storefront %d data for listing %d", *listing.StorefrontID, listing.ID)
		var storefront models.Storefront
		err := storage.QueryRow(context.Background(), `
			SELECT id, name, slug, city, address, country, latitude, longitude
			FROM user_b2c_stores
			WHERE id = $1
		`, *listing.StorefrontID).Scan(
			&storefront.ID,
			&storefront.Name,
			&storefront.Slug,
			&storefront.City,
			&storefront.Address,
			&storefront.Country,
			&storefront.Latitude,
			&storefront.Longitude,
		)
		// user_b2c_stores - legacy таблица без is_verified, устанавливаем false по умолчанию
		storefront.IsVerified = false

		if err == nil {
			// Добавляем полную информацию о витрине в документ
			doc["storefront"] = map[string]interface{}{
				"id":          storefront.ID,
				"name":        storefront.Name,
				"slug":        storefront.Slug,
				"is_verified": storefront.IsVerified,
			}

			if needStorefrontInfo {
				if listing.City == "" && storefront.City != nil && *storefront.City != "" {
					doc["city"] = *storefront.City
				}
				if listing.Country == "" && storefront.Country != nil && *storefront.Country != "" {
					doc["country"] = *storefront.Country
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
			}
		} else {
			logger.Info().Msgf("WARNING: Failed to load storefront data for listing %d: %v", listing.ID, err)
		}
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
	if len(listing.CategoryPathIds) > 0 {
		doc["category_path_ids"] = listing.CategoryPathIds
	} else {
		parentID := listing.CategoryID
		pathIDs := []int{parentID}
		for parentID > 0 {
			var cat models.MarketplaceCategory
			err := storage.QueryRow(context.Background(),
				"SELECT parent_id FROM c2c_categories WHERE id = $1", parentID).
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
			"id":   listing.CategoryID, // Используем CategoryID вместо Category.ID
			"name": listing.Category.Name,
			"slug": listing.Category.Slug,
		}
	}
}

func processUser(doc map[string]interface{}, listing *models.MarketplaceListing) {
	// Всегда добавляем user_id из listing.UserID
	if listing.UserID > 0 {
		logger.Info().Msgf("processUser: listing.ID=%d, listing.UserID=%d", listing.ID, listing.UserID)

		// Создаем базовую структуру пользователя с user_id
		userDoc := map[string]interface{}{
			"id": listing.UserID,
		}

		// Добавляем дополнительную информацию, если User объект заполнен
		if listing.User != nil {
			logger.Info().Msgf("processUser: listing.User.ID=%d, listing.User.Name=%s", listing.User.ID, listing.User.Name)

			// Если User.ID равен 0, но есть UserID, используем UserID
			if listing.User.ID == 0 && listing.UserID > 0 {
				userDoc["id"] = listing.UserID
				logger.Info().Msgf("processUser: User.ID was 0, using listing.UserID=%d", listing.UserID)
			} else if listing.User.ID > 0 {
				userDoc["id"] = listing.User.ID
			}

			if listing.User.Name != "" {
				userDoc["name"] = listing.User.Name
			}
			if listing.User.Email != "" {
				userDoc["email"] = listing.User.Email
			}
		}

		doc["user"] = userDoc
		logger.Info().Msgf("processUser: final user doc for listing %d: %v", listing.ID, userDoc)
	} else {
		logger.Warn().Msgf("processUser: listing.ID=%d has no UserID", listing.ID)
	}
}

func processImages(doc map[string]interface{}, listing *models.MarketplaceListing, storage storage.Storage) {
	// Проверяем, является ли это товаром витрины
	isStorefrontProduct := false
	if listing.Metadata != nil {
		if source, ok := listing.Metadata["source"].(string); ok && source == "storefront" {
			isStorefrontProduct = true
		}
	}

	switch {
	case len(listing.Images) > 0:
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
	case isStorefrontProduct:
		// Для B2C товаров загружаем изображения из storefront_product_images
		storefrontImages, err := storage.GetB2CProductImages(context.Background(), listing.ID)
		if err == nil && len(storefrontImages) > 0 {
			imagesDoc := make([]map[string]interface{}, 0, len(storefrontImages))
			var mainImageURL, mainThumbnailURL string

			for _, img := range storefrontImages {
				imageDoc := map[string]interface{}{
					"id":       img.ID,
					"url":      img.ImageURL,
					"alt_text": "",
					"is_main":  img.IsMain,
					"position": img.DisplayOrder,
				}
				imagesDoc = append(imagesDoc, imageDoc)

				// Запоминаем URL главного изображения
				if img.IsMain {
					mainImageURL = img.ImageURL
					mainThumbnailURL = img.ThumbnailURL
				}
			}

			doc["images"] = imagesDoc

			// Устанавливаем верхнеуровневые поля для удобства
			if mainImageURL != "" {
				doc["image_url"] = mainImageURL
				doc["thumbnail_url"] = mainThumbnailURL
			} else if len(storefrontImages) > 0 {
				// Если нет главного, берём первое
				doc["image_url"] = storefrontImages[0].ImageURL
				doc["thumbnail_url"] = storefrontImages[0].ThumbnailURL
			}
		}
	default:
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
