// backend/internal/proj/c2c/service/marketplace_listings.go
package service

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
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

func (s *MarketplaceService) SynchronizeDiscountData(ctx context.Context, listingID int) error {
	// Получаем данные из PostgreSQL
	listing, err := s.storage.GetListingByID(ctx, listingID)
	if err != nil {
		return fmt.Errorf("ошибка получения объявления: %w", err)
	}

	// Получаем историю цен
	priceHistory, err := s.storage.GetPriceHistory(ctx, listingID)
	if err != nil {
		log.Printf("Ошибка получения истории цен: %v", err)
		// Продолжаем выполнение метода даже при ошибке получения истории
	}

	// Проверяем на манипуляции с ценой, если есть история
	if len(priceHistory) > 1 {
		// Сортируем историю цен по дате
		sort.Slice(priceHistory, func(i, j int) bool {
			return priceHistory[i].EffectiveFrom.Before(priceHistory[j].EffectiveFrom)
		})

		// Проверяем на манипуляции с ценой
		isManipulation := false
		for i := 0; i < len(priceHistory)-1; i++ {
			// Если цена была повышена более чем на 50%, а затем снижена в течение 3 дней
			var nextEffectiveTo time.Time
			if priceHistory[i+1].EffectiveTo == nil {
				nextEffectiveTo = time.Now()
			} else {
				nextEffectiveTo = *priceHistory[i+1].EffectiveTo
			}

			duration := nextEffectiveTo.Sub(priceHistory[i+1].EffectiveFrom)

			if priceHistory[i].Price*1.5 < priceHistory[i+1].Price &&
				duration < 3*24*time.Hour &&
				i+2 < len(priceHistory) &&
				priceHistory[i+2].Price < priceHistory[i+1].Price {

				isManipulation = true
				log.Printf("Обнаружена манипуляция с ценой для объявления %d: повышение с %.2f до %.2f на %.1f часов, затем снижение до %.2f",
					listingID, priceHistory[i].Price, priceHistory[i+1].Price,
					duration.Hours(), priceHistory[i+2].Price)
				break
			}
		}

		if isManipulation {
			// В случае обнаружения манипуляции удаляем информацию о скидке
			if listing.Metadata != nil && listing.Metadata["discount"] != nil {
				delete(listing.Metadata, "discount")

				// Обновляем объявление в БД
				_, err := s.storage.Exec(ctx, `
                    UPDATE c2c_listings
                    SET metadata = $1
                    WHERE id = $2
                `, listing.Metadata, listingID)
				if err != nil {
					log.Printf("Ошибка удаления метаданных о скидке: %v", err)
					return err
				}

				log.Printf("Удалена информация о скидке из-за обнаружения манипуляций с ценой для объявления %d", listingID)

				// Переиндексируем объявление без скидки и возвращаемся
				return s.storage.IndexListing(ctx, listing)
			}
		}

		// Находим максимальную цену в истории с учетом минимальной продолжительности
		var maxPrice float64
		var maxPriceDate time.Time
		minDuration := 24 * time.Hour // Минимальная продолжительность - 1 день

		for _, entry := range priceHistory {
			// Рассчитываем продолжительность действия цены
			var duration time.Duration
			if entry.EffectiveTo == nil {
				duration = time.Since(entry.EffectiveFrom)
			} else {
				duration = entry.EffectiveTo.Sub(entry.EffectiveFrom)
			}

			// Учитываем только цены, которые действовали достаточно долго
			if duration >= minDuration && entry.Price > maxPrice {
				maxPrice = entry.Price
				maxPriceDate = entry.EffectiveFrom
			}
		}

		// Если текущая цена ниже максимальной, создаем скидку
		if maxPrice > listing.Price && maxPrice > 0 {
			discountPercent := int((maxPrice - listing.Price) / maxPrice * 100)

			if discountPercent >= 5 { // Если скидка составляет не менее 5%
				// Если у объекта нет метаданных, создаем их
				if listing.Metadata == nil {
					listing.Metadata = make(map[string]interface{})
				}

				// Создаем информацию о скидке с установленным флагом has_price_history = true
				discount := map[string]interface{}{
					"discount_percent":  discountPercent,
					"previous_price":    maxPrice,
					"effective_from":    maxPriceDate.Format(time.RFC3339),
					"has_price_history": true,
				}

				// Добавляем информацию о скидке в метаданные
				listing.Metadata["discount"] = discount

				// Обновляем объявление в БД
				_, err := s.storage.Exec(ctx, `
                    UPDATE c2c_listings
                    SET metadata = $1
                    WHERE id = $2
                `, listing.Metadata, listingID)
				if err != nil {
					log.Printf("Ошибка обновления метаданных: %v", err)
					return err
				}

				log.Printf("Создана информация о скидке для объявления %d из истории цен: %d%%, старая цена: %.2f",
					listingID, discountPercent, maxPrice)
			} else if listing.Metadata != nil && listing.Metadata["discount"] != nil {
				// Если скидка меньше 5%, удаляем информацию о скидке
				delete(listing.Metadata, "discount")

				// Обновляем объявление в БД
				_, err := s.storage.Exec(ctx, `
                        UPDATE c2c_listings
                        SET metadata = $1
                        WHERE id = $2
                    `, listing.Metadata, listingID)
				if err != nil {
					log.Printf("Ошибка удаления метаданных о малой скидке: %v", err)
					return err
				}

				log.Printf("Удалена информация о малой скидке (%.1f%%) для объявления %d",
					float64(discountPercent), listingID)
			}
		} else if listing.Metadata != nil && listing.Metadata["discount"] != nil {
			// Если максимальная цена не найдена или текущая цена не ниже максимальной,
			// удаляем информацию о скидке
			delete(listing.Metadata, "discount")

			// Обновляем объявление в БД
			_, err := s.storage.Exec(ctx, `
                UPDATE c2c_listings
                SET metadata = $1
                WHERE id = $2
            `, listing.Metadata, listingID)
			if err != nil {
				log.Printf("Ошибка удаления метаданных о неактуальной скидке: %v", err)
				return err
			}

			log.Printf("Удалена неактуальная информация о скидке для объявления %d", listingID)
		}
	}

	// Проверяем наличие данных о скидке в тексте описания (если нет в метаданных)
	if (listing.Metadata == nil || listing.Metadata["discount"] == nil) &&
		(strings.Contains(listing.Description, "СКИДКА") || strings.Contains(listing.Description, "СКИДКА!")) {

		discountRegex := regexp.MustCompile(`(\d+)%\s*СКИДКА`)
		matches := discountRegex.FindStringSubmatch(listing.Description)

		priceRegex := regexp.MustCompile(`Старая цена:\s*(\d+[\.,]?\d*)\s*RSD`)
		priceMatches := priceRegex.FindStringSubmatch(listing.Description)

		if len(matches) > 1 && len(priceMatches) > 1 {
			discountPercent, _ := strconv.Atoi(matches[1])
			oldPriceStr := strings.ReplaceAll(priceMatches[1], ",", ".")
			oldPrice, _ := strconv.ParseFloat(oldPriceStr, 64)

			// Проверяем реальность скидки
			calculatedDiscountPercent := int((oldPrice - listing.Price) / oldPrice * 100)
			if calculatedDiscountPercent < 0 || abs(calculatedDiscountPercent-discountPercent) > 5 {
				log.Printf("Объявленная скидка (%d%%) не соответствует расчетной (%d%%) для объявления %d, игнорируем",
					discountPercent, calculatedDiscountPercent, listingID)
				return nil
			}

			// Если у объекта нет метаданных, создаем их
			if listing.Metadata == nil {
				listing.Metadata = make(map[string]interface{})
			}

			// Создаем записи в истории цен
			// 1. Закрываем все предыдущие открытые записи истории цен
			if err := s.storage.ClosePriceHistoryEntry(ctx, listing.ID); err != nil {
				log.Printf("Ошибка закрытия прошлых записей истории цен: %v", err)
			}

			// 2. Создаем запись с предыдущей ценой, датированную неделю назад
			effectiveFrom := time.Now().AddDate(0, 0, -7)

			oldPriceEntry := &models.PriceHistoryEntry{
				ListingID:     listing.ID,
				Price:         oldPrice,
				EffectiveFrom: effectiveFrom,
				ChangeSource:  "parsed_from_description",
			}
			if err := s.storage.AddPriceHistoryEntry(ctx, oldPriceEntry); err != nil {
				log.Printf("Ошибка добавления старой цены в историю: %v", err)
			}

			// 3. Создаем запись с текущей ценой
			currentTime := time.Now()
			newPriceEntry := &models.PriceHistoryEntry{
				ListingID:     listing.ID,
				Price:         listing.Price,
				EffectiveFrom: currentTime,
				ChangeSource:  "parsed_from_description",
			}
			if err := s.storage.AddPriceHistoryEntry(ctx, newPriceEntry); err != nil {
				log.Printf("Ошибка добавления новой цены в историю: %v", err)
			}

			// Создаем информацию о скидке с установленным флагом has_price_history = true
			discount := map[string]interface{}{
				"discount_percent":  discountPercent,
				"previous_price":    oldPrice,
				"effective_from":    effectiveFrom.Format(time.RFC3339),
				"has_price_history": true,
			}

			// Добавляем информацию о скидке в метаданные
			listing.Metadata["discount"] = discount

			// Обновляем объявление в БД
			_, err := s.storage.Exec(ctx, `
                UPDATE c2c_listings
                SET metadata = $1
                WHERE id = $2
            `, listing.Metadata, listing.ID)
			if err != nil {
				log.Printf("Ошибка обновления метаданных: %v", err)
				return err
			}

			log.Printf("Создана информация о скидке для объявления %d из описания: %d%%, старая цена: %.2f",
				listing.ID, discountPercent, oldPrice)
		}
	}

	// Переиндексация в OpenSearch
	return s.storage.IndexListing(ctx, listing)
}

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
