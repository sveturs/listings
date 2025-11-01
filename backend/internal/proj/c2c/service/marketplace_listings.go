// backend/internal/proj/c2c/service/marketplace_listings.go
package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"backend/internal/domain/models"
)

func (s *MarketplaceService) CreateListing(ctx context.Context, listing *models.MarketplaceListing) (int, error) {
	listing.Status = "active"
	listing.ViewsCount = 0

	if listing.OriginalLanguage == "" {
		// Пытаемся получить язык из контекста
		if userLang, ok := ctx.Value("language").(string); ok && userLang != "" {
			listing.OriginalLanguage = userLang
			log.Printf("Using language from context: %s", userLang)
		} else if userLang, ok := ctx.Value("Accept-Language").(string); ok && userLang != "" {
			listing.OriginalLanguage = userLang
			log.Printf("Using language from Accept-Language header: %s", userLang)
		} else {
			// Используем русский по умолчанию, т.к. большинство пользователей русскоговорящие
			listing.OriginalLanguage = "ru"
			log.Printf("Using default language (ru)")
		}
	}

	// Если указана витрина, и нет данных о местоположении, получаем их из витрины
	if listing.StorefrontID != nil && (listing.City == "" || listing.Country == "" || listing.Location == "") {
		storefront, err := s.storage.GetStorefrontByID(ctx, *listing.StorefrontID)
		if err == nil && storefront != nil {
			// Применяем данные о местоположении из витрины
			if listing.City == "" && storefront.City != nil && *storefront.City != "" {
				listing.City = *storefront.City
				log.Printf("Using city from storefront: %s", *storefront.City)
			}

			if listing.Country == "" && storefront.Country != nil && *storefront.Country != "" {
				listing.Country = *storefront.Country
				log.Printf("Using country from storefront: %s", *storefront.Country)
			}

			if listing.Location == "" && storefront.Address != nil && *storefront.Address != "" {
				listing.Location = *storefront.Address
				log.Printf("Using address from storefront: %s", *storefront.Address)
			}

			// Если нет координат
			if (listing.Latitude == nil || listing.Longitude == nil) &&
				storefront.Latitude != nil && storefront.Longitude != nil {
				listing.Latitude = storefront.Latitude
				listing.Longitude = storefront.Longitude
				log.Printf("Using coordinates from storefront: Lat=%f, Lon=%f",
					*storefront.Latitude, *storefront.Longitude)
			}
		} else {
			log.Printf("Could not get storefront or storefront has no location info: %v", err)
		}
	}

	listingID, err := s.storage.CreateListing(ctx, listing)
	if err != nil {
		return 0, err
	}

	if len(listing.Attributes) > 0 {
		// Фильтрация дубликатов атрибутов
		uniqueAttrs := make(map[int]models.ListingAttributeValue)
		for _, attr := range listing.Attributes {
			uniqueAttrs[attr.AttributeID] = attr // Последнее значение перезапишет предыдущее
		}

		// Преобразуем карту обратно в срез
		filteredAttrs := make([]models.ListingAttributeValue, 0, len(uniqueAttrs))
		for _, attr := range uniqueAttrs {
			attr.ListingID = listingID
			filteredAttrs = append(filteredAttrs, attr)
		}

		if err := s.SaveListingAttributes(ctx, listingID, filteredAttrs); err != nil {
			log.Printf("Error saving attributes for listing %d: %v", listingID, err)
		} else {
			log.Printf("Successfully saved %d attributes for listing %d", len(filteredAttrs), listingID)
		}
	}
	listing.ID = listingID

	// Сохраняем переводы адресных полей
	addressFields := make(map[string]string)
	if listing.Location != "" {
		addressFields["location"] = listing.Location
	}
	if listing.City != "" {
		addressFields["city"] = listing.City
	}
	if listing.Country != "" {
		addressFields["country"] = listing.Country
	}

	// Если есть адресные поля, переводим их на другие языки
	if len(addressFields) > 0 {
		// Определяем целевые языки (все поддерживаемые языки кроме исходного)
		targetLanguages := []string{"en", "ru", "sr"}
		// Удаляем исходный язык из списка целевых
		filteredTargetLanguages := make([]string, 0, len(targetLanguages))
		for _, lang := range targetLanguages {
			if lang != listing.OriginalLanguage {
				filteredTargetLanguages = append(filteredTargetLanguages, lang)
			}
		}

		// Сохраняем переводы адресов асинхронно, чтобы не замедлять создание объявления
		go func(ctx context.Context) {
			if err := s.SaveAddressTranslations(ctx, listingID, addressFields, listing.OriginalLanguage, filteredTargetLanguages); err != nil {
				log.Printf("Error saving address translations for listing %d: %v", listingID, err)
			}
		}(ctx)
	}

	// Получаем полное объявление для индексации в OpenSearch
	fullListing, err := s.storage.GetListingByID(ctx, listingID)
	if err != nil {
		log.Printf("ERROR: Ошибка получения полного объявления ID=%d для индексации в OpenSearch: %v", listingID, err)
		// Даже если не удалось получить полное объявление, возвращаем успех создания
		// так как объявление уже создано в PostgreSQL
	} else {
		log.Printf("INFO: Индексируем объявление ID=%d в OpenSearch", listingID)
		if err := s.storage.IndexListing(ctx, fullListing); err != nil {
			log.Printf("ERROR: Ошибка индексации объявления ID=%d в OpenSearch: %v", listingID, err)
			// Не возвращаем ошибку, так как объявление уже создано в PostgreSQL
		} else {
			log.Printf("SUCCESS: Объявление ID=%d успешно проиндексировано в OpenSearch", listingID)
		}
	}

	return listingID, nil
}

func (s *MarketplaceService) GetListings(ctx context.Context, filters map[string]string, limit int, offset int) ([]models.MarketplaceListing, int64, error) {
	return s.storage.GetListings(ctx, filters, limit, offset)
}

func (s *MarketplaceService) GetListingByID(ctx context.Context, id int) (*models.MarketplaceListing, error) {
	// Получаем объявление
	listing, err := s.storage.GetListingByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Проверяем, пришел ли запрос от пользовательской части приложения (через GetListing API)
	// или из административной/служебной части (другие API)
	if ctx.Value("increment_views") == true {
		// Увеличиваем счетчик просмотров только для просмотра деталей объявления
		if err := s.storage.IncrementViewsCount(ctx, id); err != nil {
			// Логируем ошибку, но не прерываем выполнение
			log.Printf("Ошибка при увеличении счетчика просмотров для объявления %d: %v", id, err)
		}
	}

	return listing, nil
}

func (s *MarketplaceService) GetListingBySlug(ctx context.Context, slug string) (*models.MarketplaceListing, error) {
	// Получаем объявление по slug
	listing, err := s.storage.GetListingBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	// Проверяем, пришел ли запрос от пользовательской части приложения
	if ctx.Value("increment_views") == true {
		// Увеличиваем счетчик просмотров только для просмотра деталей объявления
		if err := s.storage.IncrementViewsCount(ctx, listing.ID); err != nil {
			// Логируем ошибку, но не прерываем выполнение
			log.Printf("Ошибка при увеличении счетчика просмотров для объявления %d: %v", listing.ID, err)
		}
	}

	return listing, nil
}

func (s *MarketplaceService) IsSlugAvailable(ctx context.Context, slug string, excludeID int) (bool, error) {
	return s.storage.IsSlugUnique(ctx, slug, excludeID)
}

func (s *MarketplaceService) GenerateUniqueSlug(ctx context.Context, baseSlug string, excludeID int) (string, error) {
	// Сначала проверяем исходный slug
	isUnique, err := s.storage.IsSlugUnique(ctx, baseSlug, excludeID)
	if err != nil {
		return "", err
	}

	if isUnique {
		return baseSlug, nil
	}

	// Если не уникален, пробуем с числовыми суффиксами
	for i := 2; i <= 99; i++ {
		candidateSlug := fmt.Sprintf("%s-%d", baseSlug, i)
		isUnique, err := s.storage.IsSlugUnique(ctx, candidateSlug, excludeID)
		if err != nil {
			return "", err
		}

		if isUnique {
			return candidateSlug, nil
		}
	}

	// Если все числа от 2 до 99 заняты, используем timestamp
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s-%d", baseSlug, timestamp), nil
}

func (s *MarketplaceService) UpdateListing(ctx context.Context, listing *models.MarketplaceListing) error {
	// Получаем текущую версию объявления перед обновлением
	currentListing, err := s.storage.GetListingByID(ctx, listing.ID)
	if err != nil {
		return fmt.Errorf("ошибка получения текущего объявления: %w", err)
	}

	// Проверяем изменение цены
	if currentListing.Price != listing.Price {
		// Если цена изменилась, и есть метаданные о скидке
		if currentListing.Metadata != nil {
			if discount, ok := currentListing.Metadata["discount"].(map[string]interface{}); ok {
				if prevPrice, ok := discount["previous_price"].(float64); ok {
					// Обновляем процент скидки на основе новой цены
					if prevPrice > 0 && prevPrice > listing.Price {
						// Пересчитываем скидку
						discountPercent := int((prevPrice - listing.Price) / prevPrice * 100)

						// Обновляем метаданные
						if listing.Metadata == nil {
							listing.Metadata = make(map[string]interface{})
						}

						// Копируем существующую информацию о скидке
						if listing.Metadata["discount"] == nil {
							listing.Metadata["discount"] = make(map[string]interface{})
						}

						discountMap, _ := listing.Metadata["discount"].(map[string]interface{})
						discountMap["discount_percent"] = discountPercent
						discountMap["previous_price"] = prevPrice
						discountMap["has_price_history"] = true

						log.Printf("Обновление процента скидки для объявления %d: %d%% (старая цена: %.2f, новая цена: %.2f)",
							listing.ID, discountPercent, prevPrice, listing.Price)
					}
				}
			}
		}
	}

	// Вызываем существующий метод для обновления объявления в БД
	if err := s.storage.UpdateListing(ctx, listing); err != nil {
		return err
	}

	// Проверяем, изменились ли адресные поля
	addressChanged := currentListing.Location != listing.Location ||
		currentListing.City != listing.City ||
		currentListing.Country != listing.Country

	// Если адресные поля изменились, обновляем переводы
	if addressChanged {
		addressFields := make(map[string]string)
		if listing.Location != "" {
			addressFields["location"] = listing.Location
		}
		if listing.City != "" {
			addressFields["city"] = listing.City
		}
		if listing.Country != "" {
			addressFields["country"] = listing.Country
		}

		if len(addressFields) > 0 {
			// Определяем язык обновленных данных
			sourceLanguage := listing.OriginalLanguage
			if sourceLanguage == "" {
				// Если язык не указан, используем язык из контекста или русский по умолчанию
				if userLang, ok := ctx.Value("language").(string); ok && userLang != "" {
					sourceLanguage = userLang
				} else {
					sourceLanguage = "ru"
				}
			}

			// Определяем целевые языки
			targetLanguages := []string{"en", "ru", "sr"}
			filteredTargetLanguages := make([]string, 0, len(targetLanguages))
			for _, lang := range targetLanguages {
				if lang != sourceLanguage {
					filteredTargetLanguages = append(filteredTargetLanguages, lang)
				}
			}

			// Обновляем переводы адресов асинхронно
			go func(ctx context.Context) {
				if err := s.SaveAddressTranslations(ctx, listing.ID, addressFields, sourceLanguage, filteredTargetLanguages); err != nil {
					log.Printf("Error updating address translations for listing %d: %v", listing.ID, err)
				}
			}(ctx)
		}
	}

	// Получаем полное объявление со всеми связанными данными после обновления
	fullListing, err := s.storage.GetListingByID(ctx, listing.ID)
	if err != nil {
		log.Printf("ERROR: Ошибка получения полного объявления ID=%d для обновления индекса: %v", listing.ID, err)
		return nil
	}

	// Обновляем объявление в OpenSearch
	log.Printf("INFO: Обновляем объявление ID=%d в OpenSearch", listing.ID)
	if err := s.storage.IndexListing(ctx, fullListing); err != nil {
		log.Printf("ERROR: Ошибка обновления объявления ID=%d в OpenSearch: %v", listing.ID, err)
	} else {
		log.Printf("SUCCESS: Объявление ID=%d успешно обновлено в OpenSearch", listing.ID)
	}

	return nil
}

// SynchronizeDiscountData is now implemented in marketplace_discount_sync.go
// Old God Function (335-586, 251 lines) has been refactored into:
// - marketplace_discount_detection.go (price manipulation detection)
// - marketplace_discount_parser.go (parsing discounts from description)
// - marketplace_discount_calculator.go (calculating discount from price history)
// - marketplace_discount_applier.go (applying discount to database)
// - marketplace_discount_sync.go (coordinator)

func (s *MarketplaceService) DeleteListing(ctx context.Context, id int, userID int) error {
	// Вызываем существующий метод для удаления объявления из БД
	if err := s.storage.DeleteListing(ctx, id, userID); err != nil {
		return err
	}

	// Удаляем объявление из OpenSearch
	if err := s.storage.DeleteListingIndex(ctx, fmt.Sprintf("%d", id)); err != nil {
		log.Printf("Ошибка удаления объявления из OpenSearch: %v", err)
		// Не возвращаем ошибку, чтобы не блокировать операцию, если OpenSearch недоступен
	}

	return nil
}

func (s *MarketplaceService) DeleteListingWithAdmin(ctx context.Context, id int, userID int, isAdmin bool) error {
	// Если пользователь администратор, удаляем без проверки владельца
	if isAdmin {
		if err := s.storage.DeleteListingAdmin(ctx, id); err != nil {
			return err
		}
	} else {
		// Обычное удаление с проверкой владельца
		if err := s.storage.DeleteListing(ctx, id, userID); err != nil {
			return err
		}
	}

	// Удаляем объявление из OpenSearch
	if err := s.storage.DeleteListingIndex(ctx, fmt.Sprintf("%d", id)); err != nil {
		log.Printf("Ошибка удаления объявления из OpenSearch: %v", err)
		// Не возвращаем ошибку, чтобы не блокировать операцию, если OpenSearch недоступен
	}

	return nil
}

func (s *MarketplaceService) GetPriceHistory(ctx context.Context, listingID int) ([]models.PriceHistoryEntry, error) {
	// Получаем историю цен из хранилища
	history, err := s.storage.GetPriceHistory(ctx, listingID)
	if err != nil {
		return nil, fmt.Errorf("error getting price history: %w", err)
	}

	return history, nil
}
