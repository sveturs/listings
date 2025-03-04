// backend/internal/proj/marketplace/service/marketplace.go
package service

import (
	"backend/internal/domain/models"
	"backend/internal/domain/search"
	"backend/internal/storage"
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"time"
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

	// Вызываем существующий метод для создания объявления в БД
	listingID, err := s.storage.CreateListing(ctx, listing)
	if err != nil {
		return 0, err
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

// GetCategorySuggestions возвращает предложения категорий на основе поискового запроса
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
	// Вызываем существующий метод для обновления объявления в БД
	if err := s.storage.UpdateListing(ctx, listing); err != nil {
		return err
	}

	// Получаем полное объявление со всеми связанными данными
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
    log.Printf("Преобразованные параметры поиска: %+v", osParams)
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
    if params.Latitude != 0 && params.Longitude != 0 {
        osParams.Location = &search.GeoLocation{
            Lat: params.Latitude,
            Lon: params.Longitude,
        }
        osParams.Distance = params.Distance
    }

    // Выполняем поиск
    log.Printf("Отправляем запрос в OpenSearch: %+v", osParams)
    osResult, err := s.storage.SearchListingsOpenSearch(ctx, osParams)
    if err != nil {
        log.Printf("Ошибка поиска в OpenSearch: %v", err)
        return nil, fmt.Errorf("ошибка поиска: %w", err)
    }

    log.Printf("Получен ответ из OpenSearch: найдено %d результатов", osResult.Total)
    
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
        // Используем существующие поля структуры SearchParams
        // minimum_should_match параметр в запросе
        fuzzyParams.MinimumShouldMatch = "50%" 
        // Увеличиваем fuzziness
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
            }
        }
    }

    return result, nil
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
