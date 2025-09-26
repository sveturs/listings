package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"backend/internal/domain/models"
	"backend/internal/domain/search"
	"backend/internal/logger"
)

// CreateSavedSearch создает новый сохраненный поиск
func (s *MarketplaceService) CreateSavedSearch(ctx context.Context, userID int, name string, filters map[string]interface{}, searchType string, notifyEnabled bool, notifyFrequency string) (interface{}, error) {
	// Конвертируем filters в JSON для сохранения в БД
	filtersJSON, err := json.Marshal(filters)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to marshal filters")
		return nil, fmt.Errorf("failed to marshal filters: %w", err)
	}

	query := `
		INSERT INTO saved_searches (user_id, name, filters, search_type, notify_enabled, notify_frequency, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, user_id, name, filters, search_type, notify_enabled, notify_frequency, results_count, last_notified_at, created_at, updated_at
	`

	now := time.Now()
	var savedSearch models.SavedSearch

	err = s.storage.QueryRow(ctx, query, userID, name, filtersJSON, searchType, notifyEnabled, notifyFrequency, now, now).Scan(
		&savedSearch.ID,
		&savedSearch.UserID,
		&savedSearch.Name,
		&savedSearch.Filters,
		&savedSearch.SearchType,
		&savedSearch.NotifyEnabled,
		&savedSearch.NotifyFrequency,
		&savedSearch.ResultsCount,
		&savedSearch.LastNotifiedAt,
		&savedSearch.CreatedAt,
		&savedSearch.UpdatedAt,
	)
	if err != nil {
		logger.Error().Err(err).Int("userId", userID).Str("name", name).Msg("Failed to create saved search")
		return nil, fmt.Errorf("failed to create saved search: %w", err)
	}

	logger.Info().Int("id", savedSearch.ID).Int("userId", userID).Msg("Created saved search")
	return &savedSearch, nil
}

// GetUserSavedSearches получает список сохраненных поисков пользователя
func (s *MarketplaceService) GetUserSavedSearches(ctx context.Context, userID int, searchType string) ([]interface{}, error) {
	query := `
		SELECT id, user_id, name, filters, search_type, notify_enabled, notify_frequency, results_count, last_notified_at, created_at, updated_at
		FROM saved_searches
		WHERE user_id = $1
	`

	args := []interface{}{userID}

	// Добавляем фильтр по типу поиска, если указан
	if searchType != "" {
		query += " AND search_type = $2"
		args = append(args, searchType)
	}

	query += " ORDER BY created_at DESC"

	rows, err := s.storage.Query(ctx, query, args...)
	if err != nil {
		logger.Error().Err(err).Int("userId", userID).Msg("Failed to get saved searches")
		return nil, fmt.Errorf("failed to get saved searches: %w", err)
	}
	defer rows.Close()

	var searches []interface{}
	for rows.Next() {
		var savedSearch models.SavedSearch
		err := rows.Scan(
			&savedSearch.ID,
			&savedSearch.UserID,
			&savedSearch.Name,
			&savedSearch.Filters,
			&savedSearch.SearchType,
			&savedSearch.NotifyEnabled,
			&savedSearch.NotifyFrequency,
			&savedSearch.ResultsCount,
			&savedSearch.LastNotifiedAt,
			&savedSearch.CreatedAt,
			&savedSearch.UpdatedAt,
		)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to scan saved search")
			continue
		}
		searches = append(searches, &savedSearch)
	}

	return searches, nil
}

// GetSavedSearchByID получает сохраненный поиск по ID
func (s *MarketplaceService) GetSavedSearchByID(ctx context.Context, userID int, searchID int) (interface{}, error) {
	query := `
		SELECT id, user_id, name, filters, search_type, notify_enabled, notify_frequency, results_count, last_notified_at, created_at, updated_at
		FROM saved_searches
		WHERE id = $1 AND user_id = $2
	`

	var savedSearch models.SavedSearch
	err := s.storage.QueryRow(ctx, query, searchID, userID).Scan(
		&savedSearch.ID,
		&savedSearch.UserID,
		&savedSearch.Name,
		&savedSearch.Filters,
		&savedSearch.SearchType,
		&savedSearch.NotifyEnabled,
		&savedSearch.NotifyFrequency,
		&savedSearch.ResultsCount,
		&savedSearch.LastNotifiedAt,
		&savedSearch.CreatedAt,
		&savedSearch.UpdatedAt,
	)
	if err != nil {
		logger.Error().Err(err).Int("searchId", searchID).Int("userId", userID).Msg("Failed to get saved search")
		if err.Error() == "sql: no rows in result set" {
			return nil, fmt.Errorf("saved search not found")
		}
		return nil, fmt.Errorf("failed to get saved search: %w", err)
	}

	return &savedSearch, nil
}

// UpdateSavedSearch обновляет сохраненный поиск
func (s *MarketplaceService) UpdateSavedSearch(ctx context.Context, userID int, searchID int, name string, filters map[string]interface{}, notifyEnabled *bool, notifyFrequency string) (interface{}, error) {
	// Сначала проверяем, что поиск принадлежит пользователю
	existing, err := s.GetSavedSearchByID(ctx, userID, searchID)
	if err != nil {
		return nil, err
	}

	existingSearch := existing.(*models.SavedSearch)

	// Обновляем только те поля, которые были переданы
	if name != "" {
		existingSearch.Name = name
	}

	if filters != nil {
		existingSearch.Filters = models.FiltersJSON(filters)
	}

	if notifyEnabled != nil {
		existingSearch.NotifyEnabled = *notifyEnabled
	}

	if notifyFrequency != "" {
		existingSearch.NotifyFrequency = notifyFrequency
	}

	existingSearch.UpdatedAt = time.Now()

	// Конвертируем filters в JSON для обновления в БД
	filtersJSON, err := json.Marshal(existingSearch.Filters)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal filters: %w", err)
	}

	query := `
		UPDATE saved_searches
		SET name = $1, filters = $2, notify_enabled = $3, notify_frequency = $4, updated_at = $5
		WHERE id = $6 AND user_id = $7
		RETURNING id, user_id, name, filters, search_type, notify_enabled, notify_frequency, results_count, last_notified_at, created_at, updated_at
	`

	err = s.storage.QueryRow(ctx, query,
		existingSearch.Name,
		filtersJSON,
		existingSearch.NotifyEnabled,
		existingSearch.NotifyFrequency,
		existingSearch.UpdatedAt,
		searchID,
		userID,
	).Scan(
		&existingSearch.ID,
		&existingSearch.UserID,
		&existingSearch.Name,
		&existingSearch.Filters,
		&existingSearch.SearchType,
		&existingSearch.NotifyEnabled,
		&existingSearch.NotifyFrequency,
		&existingSearch.ResultsCount,
		&existingSearch.LastNotifiedAt,
		&existingSearch.CreatedAt,
		&existingSearch.UpdatedAt,
	)
	if err != nil {
		logger.Error().Err(err).Int("searchId", searchID).Int("userId", userID).Msg("Failed to update saved search")
		return nil, fmt.Errorf("failed to update saved search: %w", err)
	}

	logger.Info().Int("id", searchID).Int("userId", userID).Msg("Updated saved search")
	return existingSearch, nil
}

// DeleteSavedSearch удаляет сохраненный поиск
func (s *MarketplaceService) DeleteSavedSearch(ctx context.Context, userID int, searchID int) error {
	query := `DELETE FROM saved_searches WHERE id = $1 AND user_id = $2`

	result, err := s.storage.Exec(ctx, query, searchID, userID)
	if err != nil {
		logger.Error().Err(err).Int("searchId", searchID).Int("userId", userID).Msg("Failed to delete saved search")
		return fmt.Errorf("failed to delete saved search: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("saved search not found")
	}

	logger.Info().Int("id", searchID).Int("userId", userID).Msg("Deleted saved search")
	return nil
}

// ExecuteSavedSearch выполняет сохраненный поиск и возвращает результаты
func (s *MarketplaceService) ExecuteSavedSearch(ctx context.Context, savedSearch interface{}) (interface{}, error) {
	ss := savedSearch.(*models.SavedSearch)

	// Конвертируем фильтры в параметры поиска
	params := &search.ServiceParams{
		Query:         "",
		Page:          1,
		Size:          20,
		Sort:          "created_at",
		SortDirection: "desc",
	}

	// Парсим фильтры из SavedSearch
	for key, value := range ss.Filters {
		switch key {
		case "query":
			if v, ok := value.(string); ok {
				params.Query = v
			}
		case "category_id", "categoryId":
			switch v := value.(type) {
			case float64:
				params.CategoryID = fmt.Sprintf("%d", int(v))
			case int:
				params.CategoryID = fmt.Sprintf("%d", v)
			case string:
				params.CategoryID = v
			case []interface{}:
				var categoryIDs []string
				for _, catID := range v {
					if id, ok := catID.(float64); ok {
						categoryIDs = append(categoryIDs, fmt.Sprintf("%d", int(id)))
					}
				}
				params.CategoryIDs = categoryIDs
			}
		case "price_min", "priceMin":
			if v, ok := value.(float64); ok {
				params.PriceMin = v
			}
		case "price_max", "priceMax":
			if v, ok := value.(float64); ok {
				params.PriceMax = v
			}
		case "condition":
			if v, ok := value.(string); ok {
				params.Condition = v
			}
		case "city":
			if v, ok := value.(string); ok {
				params.City = v
			}
		case "sortBy", "sort_by":
			if v, ok := value.(string); ok {
				// Parse sort field like "created_at_desc" or "price_asc"
				if strings.HasSuffix(v, "_desc") {
					params.Sort = strings.TrimSuffix(v, "_desc")
					params.SortDirection = "desc"
				} else if strings.HasSuffix(v, "_asc") {
					params.Sort = strings.TrimSuffix(v, "_asc")
					params.SortDirection = "asc"
				} else {
					params.Sort = v
					params.SortDirection = "desc"
				}
			}
		case "page":
			switch v := value.(type) {
			case float64:
				params.Page = int(v)
			case int:
				params.Page = v
			}
		case "pageSize", "page_size":
			switch v := value.(type) {
			case float64:
				params.Size = int(v)
			case int:
				params.Size = v
			}
		default:
			// TODO: Обработка дополнительных атрибутов фильтрации
			// Пока что пропускаем неизвестные фильтры
			_ = value
		}
	}

	// Выполняем поиск через OpenSearch
	results, err := s.SearchListingsAdvanced(ctx, params)
	if err != nil {
		logger.Error().Err(err).Int("savedSearchId", ss.ID).Msg("Failed to execute saved search")
		return nil, fmt.Errorf("failed to execute saved search: %w", err)
	}

	// Обновляем количество результатов в сохраненном поиске
	updateQuery := `UPDATE saved_searches SET results_count = $1 WHERE id = $2`
	_, updateErr := s.storage.Exec(ctx, updateQuery, results.Total, ss.ID)
	if updateErr != nil {
		logger.Error().Err(updateErr).Int("savedSearchId", ss.ID).Msg("Failed to update results count")
		// Не возвращаем ошибку, так как поиск выполнен успешно
	}

	return results, nil
}

// TrackCarView сохраняет информацию о просмотре автомобиля
func (s *MarketplaceService) TrackCarView(ctx context.Context, userID *int, listingID int, sessionID string, referrer string, deviceType string) error {
	query := `
		INSERT INTO user_car_view_history (user_id, listing_id, session_id, referrer, device_type, viewed_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	now := time.Now()
	_, err := s.storage.Exec(ctx, query, userID, listingID, sessionID, referrer, deviceType, now, now)
	if err != nil {
		logger.Error().Err(err).Int("listingId", listingID).Msg("Failed to track car view")
		return fmt.Errorf("failed to track car view: %w", err)
	}

	return nil
}

// GetUserViewHistory получает историю просмотров пользователя
func (s *MarketplaceService) GetUserViewHistory(ctx context.Context, userID int, limit int) ([]*models.UserCarViewHistory, error) {
	query := `
		SELECT id, user_id, listing_id, session_id, viewed_at, view_duration_seconds, referrer, device_type, created_at
		FROM user_car_view_history
		WHERE user_id = $1
		ORDER BY viewed_at DESC
		LIMIT $2
	`

	rows, err := s.storage.Query(ctx, query, userID, limit)
	if err != nil {
		logger.Error().Err(err).Int("userId", userID).Msg("Failed to get user view history")
		return nil, fmt.Errorf("failed to get user view history: %w", err)
	}
	defer rows.Close()

	var history []*models.UserCarViewHistory
	for rows.Next() {
		var view models.UserCarViewHistory
		err := rows.Scan(
			&view.ID,
			&view.UserID,
			&view.ListingID,
			&view.SessionID,
			&view.ViewedAt,
			&view.ViewDurationSeconds,
			&view.Referrer,
			&view.DeviceType,
			&view.CreatedAt,
		)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to scan view history")
			continue
		}
		history = append(history, &view)
	}

	return history, nil
}