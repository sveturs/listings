// backend/internal/proj/marketplace/service/marketplace.go
package service

import (
	"backend/internal/domain/models"
	"backend/internal/domain/search"
	"backend/internal/storage"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"regexp"
	"sort"
)

type MarketplaceService struct {
	storage storage.Storage
}

func NewMarketplaceService(storage storage.Storage) MarketplaceServiceInterface {
	return &MarketplaceService{
		storage: storage,
	}
}
func (s *MarketplaceService) GetUserFavorites(ctx context.Context, userID int) ([]models.MarketplaceListing, error) {
	return s.storage.GetUserFavorites(ctx, userID)
}
func (s *MarketplaceService) CreateListing(ctx context.Context, listing *models.MarketplaceListing) (int, error) {
	// Устанавливаем начальные значения
	listing.Status = "active"
	listing.ViewsCount = 0

	// Если язык не указан, определяем его
	if listing.OriginalLanguage == "" {
		// Получаем язык из контекста, если есть
		if userLang, ok := ctx.Value("language").(string); ok && userLang != "" {
			listing.OriginalLanguage = userLang
		} else {
			// По умолчанию используем русский
			listing.OriginalLanguage = "ru"
		}
	}

	listingID, err := s.storage.CreateListing(ctx, listing)
	if err != nil {
		return 0, err
	}

	// ДОБАВЬТЕ ЭТУ ПРОВЕРКУ:
	// Если у объявления есть атрибуты, сохраняем их отдельно
	if listing.Attributes != nil && len(listing.Attributes) > 0 {
		// Устанавливаем ID объявления для каждого атрибута
		for i := range listing.Attributes {
			listing.Attributes[i].ListingID = listingID
		}

		// Сохраняем атрибуты
		if err := s.SaveListingAttributes(ctx, listingID, listing.Attributes); err != nil {
			log.Printf("Error saving attributes for listing %d: %v", listingID, err)
			// Не возвращаем ошибку, чтобы не прерывать создание объявления
		} else {
			log.Printf("Successfully saved %d attributes for listing %d", len(listing.Attributes), listingID)
		}
	}

	// Синхронизируем с OpenSearch
	listing.ID = listingID

	// Получаем полное объявление со всеми связанными данными
	fullListing, err := s.storage.GetListingByID(ctx, listingID)
	if err != nil {
		log.Printf("Ошибка получения полного объявления для индексации: %v", err)
	} else {
		// Индексируем объявление в OpenSearch
		if err := s.storage.IndexListing(ctx, fullListing); err != nil {
			log.Printf("Ошибка индексации объявления в OpenSearch: %v", err)
		}
	}

	return listingID, nil
}

func (s *MarketplaceService) GetSubcategories(ctx context.Context, parentID string, limit, offset int) ([]models.CategoryTreeNode, error) {
	var parentIDInt *int

	if parentID != "" {
		// Преобразуем строку в int
		id, err := strconv.Atoi(parentID)
		if err != nil {
			return nil, fmt.Errorf("invalid parent_id: %w", err)
		}
		parentIDInt = &id
	}

	return s.storage.GetSubcategories(ctx, parentIDInt, limit, offset)
}
func (s *MarketplaceService) GetListings(ctx context.Context, filters map[string]string, limit int, offset int) ([]models.MarketplaceListing, int64, error) {
	return s.storage.GetListings(ctx, filters, limit, offset)
}
func (s *MarketplaceService) GetFavoritedUsers(ctx context.Context, listingID int) ([]int, error) {
	return s.storage.GetFavoritedUsers(ctx, listingID)
}
func (s *MarketplaceService) GetListingByID(ctx context.Context, id int) (*models.MarketplaceListing, error) {
	return s.storage.GetListingByID(ctx, id)
}

// GetCategoryAttributes получает атрибуты для указанной категории
func (s *MarketplaceService) GetCategoryAttributes(ctx context.Context, categoryID int) ([]models.CategoryAttribute, error) {
	attributes, err := s.storage.GetCategoryAttributes(ctx, categoryID)

	// Добавляем проверку и логирование JSON полей
	for _, attr := range attributes { // Заменяем i, attr на _, attr
		if attr.Options != nil {
			log.Printf("Attribute %s options (raw): %s", attr.Name, string(attr.Options))

			// Если нужно можно распарсить для проверки
			var options map[string]interface{}
			if err := json.Unmarshal(attr.Options, &options); err != nil {
				log.Printf("Error parsing options for attribute %s: %v", attr.Name, err)
			} else {
				log.Printf("Parsed options: %+v", options)
			}
		}
	}

	return attributes, err
}

func (s *MarketplaceService) SaveListingAttributes(ctx context.Context, listingID int, attributes []models.ListingAttributeValue) error {
	log.Printf("Saving %d attributes for listing %d", len(attributes), listingID)

	for i, attr := range attributes {
		log.Printf("Attribute %d: ID=%d, Name=%s, Type=%s",
			i, attr.AttributeID, attr.AttributeName, attr.AttributeType)
	}

	return s.storage.SaveListingAttributes(ctx, listingID, attributes)
}

// GetCategorySuggestions возвращает предложения категорий на основе поискового запроса
func (s *MarketplaceService) GetCategorySuggestions(ctx context.Context, query string, size int) ([]models.CategorySuggestion, error) {
	log.Printf("Запрос предложений категорий: '%s'", query)

	// Проверка входных параметров
	if query == "" {
		return []models.CategorySuggestion{}, nil
	}

	// Выполняем SQL-запрос для поиска категорий, связанных с запросом
	sqlQuery := `
        WITH RECURSIVE category_tree AS (
            SELECT c.id, c.name, c.parent_id
            FROM marketplace_categories c
            WHERE 1=1
            
            UNION
            
            SELECT c.id, c.name, c.parent_id
            FROM marketplace_categories c
            JOIN category_tree t ON c.parent_id = t.id
        ),
        matching_categories AS (
            SELECT 
                c.id,
                c.name,
                (SELECT COUNT(*) FROM marketplace_listings ml 
                 WHERE ml.category_id = c.id 
                 AND ml.status = 'active') as listing_count,
                CASE WHEN LOWER(c.name) LIKE LOWER($1) THEN 100 ELSE 0 END +
                (SELECT COUNT(*) FROM marketplace_listings ml 
                 WHERE ml.category_id = c.id 
                 AND (LOWER(ml.title) LIKE LOWER($1) OR LOWER(ml.description) LIKE LOWER($1)) 
                 AND ml.status = 'active') as relevance
            FROM marketplace_categories c
            WHERE LOWER(c.name) LIKE LOWER($1)
            OR EXISTS (
                SELECT 1 FROM marketplace_listings ml 
                WHERE ml.category_id = c.id 
                AND (LOWER(ml.title) LIKE LOWER($1) OR LOWER(ml.description) LIKE LOWER($1))
                AND ml.status = 'active'
            )
        )
        SELECT id, name, listing_count
        FROM matching_categories
        WHERE listing_count > 0
        ORDER BY relevance DESC, listing_count DESC
        LIMIT $2
    `

	rows, err := s.storage.Query(ctx, sqlQuery, "%"+query+"%", size)
	if err != nil {
		log.Printf("Ошибка при выполнении запроса категорий: %v", err)
		return []models.CategorySuggestion{}, nil
	}
	defer rows.Close()

	var results []models.CategorySuggestion
	for rows.Next() {
		var suggestion models.CategorySuggestion

		if err := rows.Scan(&suggestion.ID, &suggestion.Name, &suggestion.ListingCount); err != nil {
			log.Printf("Ошибка сканирования категории: %v", err)
			continue
		}

		results = append(results, suggestion)
	}

	log.Printf("Найдено %d релевантных категорий для запроса '%s'", len(results), query)

	return results, nil
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

    // Получаем полное объявление со всеми связанными данными после обновления
    fullListing, err := s.storage.GetListingByID(ctx, listing.ID)
    if err != nil {
        log.Printf("Ошибка получения полного объявления для обновления индекса: %v", err)
        return nil
    }

    // Обновляем объявление в OpenSearch
    if err := s.storage.IndexListing(ctx, fullListing); err != nil {
        log.Printf("Ошибка обновления объявления в OpenSearch: %v", err)
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
                    UPDATE marketplace_listings
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
        var minDuration = 24 * time.Hour // Минимальная продолжительность - 1 день

        for _, entry := range priceHistory {
            // Рассчитываем продолжительность действия цены
            var duration time.Duration
            if entry.EffectiveTo == nil {
                duration = time.Now().Sub(entry.EffectiveFrom)
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
                    "discount_percent": discountPercent,
                    "previous_price": maxPrice,
                    "effective_from": maxPriceDate.Format(time.RFC3339),
                    "has_price_history": true,
                }
                
                // Добавляем информацию о скидке в метаданные
                listing.Metadata["discount"] = discount
                
                // Обновляем объявление в БД
                _, err := s.storage.Exec(ctx, `
                    UPDATE marketplace_listings
                    SET metadata = $1
                    WHERE id = $2
                `, listing.Metadata, listingID)
                
                if err != nil {
                    log.Printf("Ошибка обновления метаданных: %v", err)
                    return err
                }
                
                log.Printf("Создана информация о скидке для объявления %d из истории цен: %d%%, старая цена: %.2f",
                    listingID, discountPercent, maxPrice)
            } else {
                // Если скидка меньше 5%, удаляем информацию о скидке
                if listing.Metadata != nil && listing.Metadata["discount"] != nil {
                    delete(listing.Metadata, "discount")
                    
                    // Обновляем объявление в БД
                    _, err := s.storage.Exec(ctx, `
                        UPDATE marketplace_listings
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
            }
        } else if listing.Metadata != nil && listing.Metadata["discount"] != nil {
            // Если максимальная цена не найдена или текущая цена не ниже максимальной,
            // удаляем информацию о скидке
            delete(listing.Metadata, "discount")
            
            // Обновляем объявление в БД
            _, err := s.storage.Exec(ctx, `
                UPDATE marketplace_listings
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
            oldPriceStr := strings.Replace(priceMatches[1], ",", ".", -1)
            oldPrice, _ := strconv.ParseFloat(oldPriceStr, 64)

            // Проверяем реальность скидки
            calculatedDiscountPercent := int((oldPrice - listing.Price) / oldPrice * 100)
            if calculatedDiscountPercent < 0 || abs(calculatedDiscountPercent - discountPercent) > 5 {
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
                "discount_percent": discountPercent,
                "previous_price": oldPrice,
                "effective_from": effectiveFrom.Format(time.RFC3339),
                "has_price_history": true,
            }

            // Добавляем информацию о скидке в метаданные
            listing.Metadata["discount"] = discount
            
            // Обновляем объявление в БД
            _, err := s.storage.Exec(ctx, `
                UPDATE marketplace_listings
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

// Вспомогательная функция для вычисления абсолютного значения
func abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}

func (s *MarketplaceService) GetCategoryTree(ctx context.Context) ([]models.CategoryTreeNode, error) {
	return s.storage.GetCategoryTree(ctx)
}
func (s *MarketplaceService) RefreshCategoryListingCounts(ctx context.Context) error {
	_, err := s.storage.Exec(ctx, "REFRESH MATERIALIZED VIEW CONCURRENTLY category_listing_counts")
	return err
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

func (s *MarketplaceService) ProcessImage(file *multipart.FileHeader) (string, error) {
	// Получаем расширение файла
	ext := filepath.Ext(file.Filename)
	if ext == "" {
		// Если расширение отсутствует, определяем его по MIME-типу
		switch file.Header.Get("Content-Type") {
		case "image/jpeg", "image/jpg":
			ext = ".jpg"
		case "image/png":
			ext = ".png"
		case "image/gif":
			ext = ".gif"
		case "image/webp":
			ext = ".webp"
		default:
			ext = ".jpg" // По умолчанию используем .jpg
		}
	}

	// Генерируем уникальное имя файла с расширением
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	return fileName, nil
}

func (s *MarketplaceService) AddListingImage(ctx context.Context, image *models.MarketplaceImage) (int, error) {
	return s.storage.AddListingImage(ctx, image)
}

func (s *MarketplaceService) GetCategories(ctx context.Context) ([]models.MarketplaceCategory, error) {
	return s.storage.GetCategories(ctx)
}

func (s *MarketplaceService) AddToFavorites(ctx context.Context, userID int, listingID int) error {
	return s.storage.AddToFavorites(ctx, userID, listingID)
}

func (s *MarketplaceService) RemoveFromFavorites(ctx context.Context, userID int, listingID int) error {
	return s.storage.RemoveFromFavorites(ctx, userID, listingID)
}
func (s *MarketplaceService) UpdateTranslation(ctx context.Context, translation *models.Translation) error {
	query := `
        INSERT INTO translations (
            entity_type, entity_id, language, field_name,
            translated_text, is_machine_translated, is_verified
        ) VALUES ($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT (entity_type, entity_id, language, field_name)
        DO UPDATE SET
            translated_text = EXCLUDED.translated_text,
            is_machine_translated = EXCLUDED.is_machine_translated,
            is_verified = EXCLUDED.is_verified,
            updated_at = CURRENT_TIMESTAMP
    `

	_, err := s.storage.Exec(ctx, query,
		translation.EntityType,
		translation.EntityID,
		translation.Language,
		translation.FieldName,
		translation.TranslatedText,
		translation.IsMachineTranslated,
		translation.IsVerified)

	return err
}

// Исправленная версия функции SearchListingsAdvanced
func (s *MarketplaceService) SearchListingsAdvanced(ctx context.Context, params *search.ServiceParams) (*search.ServiceResult, error) {
    log.Printf("Запрос поиска с параметрами: %+v", params)
    
    // Преобразуем параметры для OpenSearch
    osParams := &search.SearchParams{
        Query:        params.Query,
        Page:         params.Page,
        Size:         params.Size,
        Aggregations: params.Aggregations,
        Language:     params.Language,
    }

    // Корректно инициализируем фильтры атрибутов, если они есть в params
    if params.AttributeFilters != nil && len(params.AttributeFilters) > 0 {
        log.Printf("Атрибуты из params: %+v", params.AttributeFilters)
        
        osParams.AttributeFilters = make(map[string]string)
        for key, value := range params.AttributeFilters {
            osParams.AttributeFilters[key] = value
            log.Printf("Передаю атрибут в OpenSearch: %s = %s", key, value)
        }
    }

                
        // Копируем атрибуты без изменений
        osParams.AttributeFilters = make(map[string]string)
        for key, value := range params.AttributeFilters {
            // Передаем значение без изменений
            osParams.AttributeFilters[key] = value
            log.Printf("Передаю атрибут в OpenSearch: %s = %s", key, value)
        }
    
    

	// Устанавливаем необязательные параметры нечеткого поиска, если они указаны
	if params.MinimumShouldMatch != "" {
		osParams.MinimumShouldMatch = params.MinimumShouldMatch
	}

	if params.Fuzziness != "" {
		osParams.Fuzziness = params.Fuzziness
	}

	// Добавляем фильтры
	if params.CategoryID != "" {
		categoryID, err := strconv.Atoi(params.CategoryID)
		if err == nil {
			osParams.CategoryID = &categoryID
		}
	}

	if params.PriceMin > 0 {
		osParams.PriceMin = &params.PriceMin
	}

	if params.PriceMax > 0 {
		osParams.PriceMax = &params.PriceMax
	}

	if params.Condition != "" {
		osParams.Condition = params.Condition
	}

	if params.City != "" {
		osParams.City = params.City
	}

	if params.Country != "" {
		osParams.Country = params.Country
	}

	if params.StorefrontID != "" {
		storefrontID, err := strconv.Atoi(params.StorefrontID)
		if err == nil {
			osParams.StorefrontID = &storefrontID
		}
	}

	// Устанавливаем сортировку
	if params.Sort != "" {
		osParams.Sort = params.Sort
		osParams.SortDirection = params.SortDirection
	}

	// Добавляем геолокацию
	var geoSearchEnabled bool

	// Проверяем наличие координат для геопоиска
	if params.Latitude != 0 && params.Longitude != 0 && params.Distance != "" {
		log.Printf("Установлен фильтр по расстоянию: %s от координат (%.6f, %.6f)",
			params.Distance, params.Latitude, params.Longitude)

		osParams.Location = &search.GeoLocation{
			Lat: params.Latitude,
			Lon: params.Longitude,
		}
		osParams.Distance = params.Distance
		geoSearchEnabled = true
	} else if (params.City != "" || params.Country != "") && params.Distance != "" {
		// Если координаты не указаны, но указан город/страна и расстояние,
		// попробуем получить координаты города
		log.Printf("Попытка определить координаты для города: %s, страна: %s",
			params.City, params.Country)

		// Формируем запрос для геокодирования
		query := params.City
		if params.Country != "" {
			query += ", " + params.Country
		}

		// Используем OSM для получения координат города
		// В реальном коде лучше реализовать отдельный сервис и кэширование результатов
		geoResult, err := s.geocodeAddress(ctx, query)
		if err == nil && geoResult != nil {
			log.Printf("Получены координаты для города: %s (%.6f, %.6f)",
				query, geoResult.Lat, geoResult.Lon)

			osParams.Location = &search.GeoLocation{
				Lat: geoResult.Lat,
				Lon: geoResult.Lon,
			}
			osParams.Distance = params.Distance
			geoSearchEnabled = true
		} else {
			log.Printf("Не удалось получить координаты для города: %s, ошибка: %v",
				query, err)
		}
	}

	// Выполняем поиск
	log.Printf("Отправляем запрос в OpenSearch: %+v", osParams)
	var osResult *search.SearchResult
	var err error

	if geoSearchEnabled {
		// Если геопоиск возможен, пробуем его
		osResult, err = s.storage.SearchListingsOpenSearch(ctx, osParams)

		if err != nil {
			log.Printf("Ошибка геопоиска в OpenSearch: %v", err)

			// Если ошибка связана с geo_distance, выполняем обычный поиск без геопараметров
			if strings.Contains(err.Error(), "geo") {
				log.Printf("Проблема с geo-поиском, выполняем запрос без геопараметров")
				osParams.Location = nil
				osParams.Distance = ""

				osResult, err = s.storage.SearchListingsOpenSearch(ctx, osParams)
			}
		}
	} else {
		// Обычный поиск без геопараметров
		osResult, err = s.storage.SearchListingsOpenSearch(ctx, osParams)
	}

	if err != nil {
		log.Printf("Ошибка поиска в OpenSearch: %v", err)

		// Если поиск в OpenSearch не удался, используем стандартный поиск через БД
		filters := map[string]string{
			"category_id":   params.CategoryID,
			"condition":     params.Condition,
			"city":          params.City,
			"country":       params.Country,
			"storefront_id": params.StorefrontID,
			"sort_by":       params.Sort,
		}

		// Добавляем числовые фильтры, если они указаны
		if params.PriceMin > 0 {
			filters["min_price"] = fmt.Sprintf("%g", params.PriceMin)
		}
		if params.PriceMax > 0 {
			filters["max_price"] = fmt.Sprintf("%g", params.PriceMax)
		}

		// Добавляем текстовый поиск
		if params.Query != "" {
			filters["query"] = params.Query
		}

		// Пробуем получить обычным методом
		listings, total, err := s.GetListings(ctx, filters, params.Size, (params.Page-1)*params.Size)
		if err != nil {
			log.Printf("Ошибка стандартного поиска: %v", err)
			return nil, fmt.Errorf("ошибка поиска: %w", err)
		}

		// Преобразуем результаты
		listingPtrs := make([]*models.MarketplaceListing, len(listings))
		for i := range listings {
			listingPtrs[i] = &listings[i]
		}

		// Формируем такой же результат, как от OpenSearch
		result := &search.ServiceResult{
			Items:      listingPtrs,
			Total:      int(total),
			Page:       params.Page,
			Size:       params.Size,
			TotalPages: (int(total) + params.Size - 1) / params.Size,
		}

		return result, nil
	}

	// Преобразуем результат
	result := &search.ServiceResult{
		Items:      osResult.Listings,
		Total:      osResult.Total,
		Page:       params.Page,
		Size:       params.Size,
		TotalPages: (osResult.Total + params.Size - 1) / params.Size,
		Took:       osResult.Took,
	}

	// Если не найдено результатов и есть поисковый запрос, попробуем более гибкий поиск
	if len(result.Items) == 0 && params.Query != "" {
		// Создаем копию параметров
		fuzzyParams := *osParams

		// Изменяем параметры для более нечеткого поиска
		fuzzyParams.MinimumShouldMatch = "50%"
		fuzzyParams.Fuzziness = "2"

		// Повторяем поиск с новыми параметрами
		fuzzyResult, err := s.storage.SearchListingsOpenSearch(ctx, &fuzzyParams)
		if err == nil && fuzzyResult != nil && len(fuzzyResult.Listings) > 0 {
			// Используем результаты нечеткого поиска
			result.Items = fuzzyResult.Listings
			result.Total = fuzzyResult.Total

			// Добавляем предложение исправления
			if len(fuzzyResult.Suggestions) > 0 {
				result.Suggestions = fuzzyResult.Suggestions
			} else {
				// Добавляем подсказку, что результаты нечеткие
				result.Suggestions = []string{params.Query + " (исправлено)"}
			}
		}
	}

	// Добавляем фасеты
	if len(osResult.Aggregations) > 0 {
		result.Facets = make(map[string][]search.Bucket)
		for name, buckets := range osResult.Aggregations {
			result.Facets[name] = buckets
		}
	}

	// Добавляем предложения
	if len(osResult.Suggestions) > 0 {
		result.Suggestions = osResult.Suggestions
	}

	if len(result.Items) < 2 && params.Query != "" {
		// Получаем предложения исправлений напрямую из базы данных
		query := `
            SELECT DISTINCT title 
            FROM marketplace_listings 
            WHERE LOWER(title) LIKE $1 
            AND status = 'active'
            LIMIT 5
        `
		rows, err := s.storage.Query(ctx, query, "%"+params.Query+"%")
		if err == nil {
			var suggestions []string
			for rows.Next() {
				var title string
				if err := rows.Scan(&title); err == nil {
					suggestions = append(suggestions, title)
				}
			}
			rows.Close()

			// Если найдены предложения, добавляем их в результат
			if len(suggestions) > 0 {
				result.Suggestions = suggestions

				// Возможно, добавить наилучшее предложение как исправление опечатки
				if len(result.Items) == 0 {
					result.SpellingSuggestion = suggestions[0]
				}
			}
		}
	}

	return result, nil
}

// geocodeAddress получает координаты по адресу/городу
func (s *MarketplaceService) geocodeAddress(ctx context.Context, address string) (*search.GeoLocation, error) {
	// Проверка кэша для ускорения работы
	// TODO: Реализовать кэширование

	// Используем OSM Nominatim для геокодирования
	// В реальном приложении лучше использовать платное API или создать свою БД городов
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	url := fmt.Sprintf(
		"https://nominatim.openstreetmap.org/search?format=json&q=%s&limit=1",
		url.QueryEscape(address),
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
		return nil, fmt.Errorf("адрес не найден")
	}

	lat, err := strconv.ParseFloat(results[0].Lat, 64)
	if err != nil {
		return nil, err
	}

	lon, err := strconv.ParseFloat(results[0].Lon, 64)
	if err != nil {
		return nil, err
	}

	return &search.GeoLocation{
		Lat: lat,
		Lon: lon,
	}, nil
}

// GetSuggestions возвращает предложения автодополнения
func (s *MarketplaceService) GetSuggestions(ctx context.Context, prefix string, size int) ([]string, error) {
	log.Printf("Запрос подсказок в сервисе: '%s'", prefix)

	// Проверка входных параметров
	if prefix == "" {
		return []string{}, nil
	}

	// Сначала пытаемся получить подсказки из OpenSearch
	suggestions, err := s.storage.SuggestListings(ctx, prefix, size)
	if err != nil || len(suggestions) == 0 {
		if err != nil {
			log.Printf("Ошибка при получении подсказок из OpenSearch: %v", err)
		}

		// При ошибке используем запасной вариант из базы данных
		// Находим примерно строку 542
		query := `
    SELECT DISTINCT title,
           CASE WHEN LOWER(title) = LOWER($2) THEN 0
                WHEN LOWER(title) LIKE LOWER($2 || '%') THEN 1
                ELSE 2
           END as relevance,
           length(title) as title_length
    FROM marketplace_listings 
    WHERE LOWER(title) LIKE LOWER($1) 
    AND status = 'active'
    ORDER BY relevance, title_length
    LIMIT $3
`
		rows, err := s.storage.Query(ctx, query, "%"+prefix+"%", prefix, size)
		if err != nil {
			log.Printf("Ошибка запасного SQL-запроса: %v", err)
			return []string{}, nil
		}
		defer rows.Close()

		var results []string
		for rows.Next() {
			var title string
			if err := rows.Scan(&title); err != nil {
				log.Printf("Ошибка сканирования строки: %v", err)
				continue
			}
			results = append(results, title)
		}

		log.Printf("Получено %d подсказок из базы данных", len(results))
		return results, nil
	}

	log.Printf("Получено %d подсказок из OpenSearch", len(suggestions))
	return suggestions, nil
}
func (s *MarketplaceService) ReindexAllListings(ctx context.Context) error {
	return s.storage.ReindexAllListings(ctx)
}
func (s *MarketplaceService) Storage() storage.Storage {
	return s.storage
}

func (s *MarketplaceService) GetPriceHistory(ctx context.Context, listingID int) ([]models.PriceHistoryEntry, error) {
    // Получаем историю цен из хранилища
    history, err := s.storage.GetPriceHistory(ctx, listingID)
    if err != nil {
        return nil, fmt.Errorf("error getting price history: %w", err)
    }
    
    return history, nil
}