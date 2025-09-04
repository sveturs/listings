package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"backend/internal/domain/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UnifiedAttributeStorage интерфейс для работы с унифицированными атрибутами
type UnifiedAttributeStorage interface {
	// Атрибуты
	CreateAttribute(ctx context.Context, attr *models.UnifiedAttribute) (int, error)
	GetAttribute(ctx context.Context, id int) (*models.UnifiedAttribute, error)
	GetAttributeByCode(ctx context.Context, code string) (*models.UnifiedAttribute, error)
	ListAttributes(ctx context.Context, filter *models.UnifiedAttributeFilter) ([]*models.UnifiedAttribute, error)
	UpdateAttribute(ctx context.Context, id int, updates map[string]interface{}) error
	DeleteAttribute(ctx context.Context, id int) error

	// Связи с категориями
	AttachAttributeToCategory(ctx context.Context, categoryID, attributeID int, settings *models.UnifiedCategoryAttribute) error
	DetachAttributeFromCategory(ctx context.Context, categoryID, attributeID int) error
	UpdateCategoryAttribute(ctx context.Context, categoryID, attributeID int, isRequired, isFilter *bool, displayOrder *int, groupID *int) error
	GetCategoryAttributes(ctx context.Context, categoryID int) ([]*models.UnifiedAttribute, error)
	GetCategoryAttributesWithSettings(ctx context.Context, categoryID int) ([]*models.UnifiedCategoryAttribute, error)

	// Значения атрибутов
	SaveAttributeValue(ctx context.Context, value *models.UnifiedAttributeValue) error
	GetAttributeValues(ctx context.Context, entityType models.AttributeEntityType, entityID int) ([]*models.UnifiedAttributeValue, error)
	DeleteAttributeValues(ctx context.Context, entityType models.AttributeEntityType, entityID int) error

	// Миграция и совместимость
	GetAttributeByLegacyID(ctx context.Context, legacyID int, isProductVariant bool) (*models.UnifiedAttribute, error)
	MigrateFromLegacySystem(ctx context.Context) error

	// Кеширование
	InvalidateCache(categoryID int)

	// Вариативные атрибуты
	GetVariantCompatibleAttributes(ctx context.Context) ([]*models.UnifiedAttribute, error)
	GetCategoryVariantMappings(ctx context.Context, categoryID int) ([]*models.VariantAttributeMapping, error)
	CreateVariantMapping(ctx context.Context, mapping *models.VariantAttributeMappingCreateRequest) (*models.VariantAttributeMapping, error)
	UpdateVariantMapping(ctx context.Context, id int, update *models.VariantAttributeMappingUpdateRequest) error
	DeleteVariantMapping(ctx context.Context, id int) error
	DeleteCategoryVariantMappings(ctx context.Context, categoryID int) error
}

type unifiedAttributeStorage struct {
	pool *pgxpool.Pool

	// Кеш для атрибутов категорий
	cacheMutex sync.RWMutex
	cache      map[int][]*models.UnifiedAttribute
	cacheTime  map[int]time.Time
	cacheTTL   time.Duration

	// Флаг для совместимости со старой системой
	useLegacyFallback bool
}

// NewUnifiedAttributeStorage создает новое хранилище для унифицированных атрибутов
func NewUnifiedAttributeStorage(pool *pgxpool.Pool, useLegacyFallback bool) UnifiedAttributeStorage {
	return &unifiedAttributeStorage{
		pool:              pool,
		cache:             make(map[int][]*models.UnifiedAttribute),
		cacheTime:         make(map[int]time.Time),
		cacheTTL:          30 * time.Minute,
		useLegacyFallback: useLegacyFallback,
	}
}

// CreateAttribute создает новый атрибут
func (s *unifiedAttributeStorage) CreateAttribute(ctx context.Context, attr *models.UnifiedAttribute) (int, error) {
	query := `
		INSERT INTO unified_attributes (
			code, name, display_name, attribute_type, purpose,
			options, validation_rules, ui_settings,
			is_searchable, is_filterable, is_required,
			affects_stock, affects_price, sort_order, is_active
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id, created_at, updated_at`

	err := s.pool.QueryRow(ctx, query,
		attr.Code, attr.Name, attr.DisplayName, attr.AttributeType, attr.Purpose,
		attr.Options, attr.ValidationRules, attr.UISettings,
		attr.IsSearchable, attr.IsFilterable, attr.IsRequired,
		attr.AffectsStock, attr.AffectsPrice, attr.SortOrder, attr.IsActive,
	).Scan(&attr.ID, &attr.CreatedAt, &attr.UpdatedAt)
	if err != nil {
		return 0, fmt.Errorf("failed to create attribute: %w", err)
	}

	return attr.ID, nil
}

// GetAttribute получает атрибут по ID
func (s *unifiedAttributeStorage) GetAttribute(ctx context.Context, id int) (*models.UnifiedAttribute, error) {
	query := `
		SELECT 
			id, code, name, display_name, attribute_type, purpose,
			options, validation_rules, ui_settings,
			is_searchable, is_filterable, is_required,
			affects_stock, affects_price, sort_order, is_active,
			created_at, updated_at,
			legacy_category_attribute_id, legacy_product_variant_attribute_id
		FROM unified_attributes
		WHERE id = $1`

	attr := &models.UnifiedAttribute{}
	err := s.pool.QueryRow(ctx, query, id).Scan(
		&attr.ID, &attr.Code, &attr.Name, &attr.DisplayName,
		&attr.AttributeType, &attr.Purpose,
		&attr.Options, &attr.ValidationRules, &attr.UISettings,
		&attr.IsSearchable, &attr.IsFilterable, &attr.IsRequired,
		&attr.AffectsStock, &attr.AffectsPrice, &attr.SortOrder, &attr.IsActive,
		&attr.CreatedAt, &attr.UpdatedAt,
		&attr.LegacyCategoryAttributeID, &attr.LegacyProductVariantAttributeID,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("attribute not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get attribute: %w", err)
	}

	return attr, nil
}

// GetAttributeByCode получает атрибут по коду
func (s *unifiedAttributeStorage) GetAttributeByCode(ctx context.Context, code string) (*models.UnifiedAttribute, error) {
	query := `
		SELECT 
			id, code, name, display_name, attribute_type, purpose,
			options, validation_rules, ui_settings,
			is_searchable, is_filterable, is_required,
			affects_stock, affects_price, sort_order, is_active,
			created_at, updated_at,
			legacy_category_attribute_id, legacy_product_variant_attribute_id
		FROM unified_attributes
		WHERE code = $1`

	attr := &models.UnifiedAttribute{}
	err := s.pool.QueryRow(ctx, query, code).Scan(
		&attr.ID, &attr.Code, &attr.Name, &attr.DisplayName,
		&attr.AttributeType, &attr.Purpose,
		&attr.Options, &attr.ValidationRules, &attr.UISettings,
		&attr.IsSearchable, &attr.IsFilterable, &attr.IsRequired,
		&attr.AffectsStock, &attr.AffectsPrice, &attr.SortOrder, &attr.IsActive,
		&attr.CreatedAt, &attr.UpdatedAt,
		&attr.LegacyCategoryAttributeID, &attr.LegacyProductVariantAttributeID,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("attribute not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get attribute: %w", err)
	}

	return attr, nil
}

// ListAttributes получает список атрибутов с фильтрацией
func (s *unifiedAttributeStorage) ListAttributes(ctx context.Context, filter *models.UnifiedAttributeFilter) ([]*models.UnifiedAttribute, error) {
	query := `
		SELECT DISTINCT
			ua.id, ua.code, ua.name, ua.display_name, ua.attribute_type, ua.purpose,
			ua.options, ua.validation_rules, ua.ui_settings,
			ua.is_searchable, ua.is_filterable, ua.is_required,
			ua.affects_stock, ua.affects_price, ua.sort_order, ua.is_active,
			ua.created_at, ua.updated_at,
			ua.legacy_category_attribute_id, ua.legacy_product_variant_attribute_id
		FROM unified_attributes ua`

	var whereConditions []string
	var args []interface{}
	argCount := 0

	// Добавляем JOIN если нужна фильтрация по категории
	if filter != nil && filter.CategoryID != nil {
		query += `
			INNER JOIN unified_category_attributes uca ON ua.id = uca.attribute_id`
		whereConditions = append(whereConditions, fmt.Sprintf("uca.category_id = $%d", argCount+1))
		args = append(args, *filter.CategoryID)
		argCount++
	}

	// Добавляем остальные условия фильтрации
	if filter != nil {
		if filter.Purpose != nil {
			whereConditions = append(whereConditions, fmt.Sprintf("(ua.purpose = $%d OR ua.purpose = 'both')", argCount+1))
			args = append(args, *filter.Purpose)
			argCount++
		}
		if filter.IsActive != nil {
			whereConditions = append(whereConditions, fmt.Sprintf("ua.is_active = $%d", argCount+1))
			args = append(args, *filter.IsActive)
			argCount++
		}
		if filter.IsSearchable != nil {
			whereConditions = append(whereConditions, fmt.Sprintf("ua.is_searchable = $%d", argCount+1))
			args = append(args, *filter.IsSearchable)
			argCount++
		}
		if filter.IsFilterable != nil {
			whereConditions = append(whereConditions, fmt.Sprintf("ua.is_filterable = $%d", argCount+1))
			args = append(args, *filter.IsFilterable)
			argCount++
		}
	}

	// Добавляем WHERE clause если есть условия
	if len(whereConditions) > 0 {
		query += " WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Добавляем сортировку
	query += " ORDER BY ua.sort_order, ua.name"

	// Добавляем пагинацию
	if filter != nil {
		if filter.Limit > 0 {
			query += fmt.Sprintf(" LIMIT %d", filter.Limit)
		}
		if filter.Offset > 0 {
			query += fmt.Sprintf(" OFFSET %d", filter.Offset)
		}
	}

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list attributes: %w", err)
	}
	defer rows.Close()

	var attributes []*models.UnifiedAttribute
	for rows.Next() {
		attr := &models.UnifiedAttribute{}
		err := rows.Scan(
			&attr.ID, &attr.Code, &attr.Name, &attr.DisplayName,
			&attr.AttributeType, &attr.Purpose,
			&attr.Options, &attr.ValidationRules, &attr.UISettings,
			&attr.IsSearchable, &attr.IsFilterable, &attr.IsRequired,
			&attr.AffectsStock, &attr.AffectsPrice, &attr.SortOrder, &attr.IsActive,
			&attr.CreatedAt, &attr.UpdatedAt,
			&attr.LegacyCategoryAttributeID, &attr.LegacyProductVariantAttributeID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan attribute: %w", err)
		}
		attributes = append(attributes, attr)
	}

	return attributes, nil
}

// UpdateAttribute обновляет атрибут
func (s *unifiedAttributeStorage) UpdateAttribute(ctx context.Context, id int, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	setClauses := []string{}
	args := []interface{}{}
	argCount := 1

	for field, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", field, argCount))
		args = append(args, value)
		argCount++
	}

	query := fmt.Sprintf(`
		UPDATE unified_attributes
		SET %s, updated_at = CURRENT_TIMESTAMP
		WHERE id = $%d`,
		strings.Join(setClauses, ", "), argCount)

	args = append(args, id)

	_, err := s.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update attribute: %w", err)
	}

	// Инвалидируем кеш для всех категорий (можно оптимизировать)
	s.cacheMutex.Lock()
	s.cache = make(map[int][]*models.UnifiedAttribute)
	s.cacheTime = make(map[int]time.Time)
	s.cacheMutex.Unlock()

	return nil
}

// DeleteAttribute удаляет атрибут
func (s *unifiedAttributeStorage) DeleteAttribute(ctx context.Context, id int) error {
	_, err := s.pool.Exec(ctx, "DELETE FROM unified_attributes WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete attribute: %w", err)
	}

	// Инвалидируем кеш
	s.cacheMutex.Lock()
	s.cache = make(map[int][]*models.UnifiedAttribute)
	s.cacheTime = make(map[int]time.Time)
	s.cacheMutex.Unlock()

	return nil
}

// AttachAttributeToCategory привязывает атрибут к категории
func (s *unifiedAttributeStorage) AttachAttributeToCategory(ctx context.Context, categoryID, attributeID int, settings *models.UnifiedCategoryAttribute) error {
	query := `
		INSERT INTO unified_category_attributes (
			category_id, attribute_id, is_enabled, is_required, sort_order, category_specific_options
		) VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (category_id, attribute_id)
		DO UPDATE SET
			is_enabled = EXCLUDED.is_enabled,
			is_required = EXCLUDED.is_required,
			sort_order = EXCLUDED.sort_order,
			category_specific_options = EXCLUDED.category_specific_options,
			updated_at = CURRENT_TIMESTAMP`

	_, err := s.pool.Exec(ctx, query,
		categoryID, attributeID,
		settings.IsEnabled, settings.IsRequired, settings.SortOrder,
		settings.CategorySpecificOptions,
	)
	if err != nil {
		return fmt.Errorf("failed to attach attribute to category: %w", err)
	}

	// Инвалидируем кеш для этой категории
	s.InvalidateCache(categoryID)

	return nil
}

// DetachAttributeFromCategory отвязывает атрибут от категории
func (s *unifiedAttributeStorage) DetachAttributeFromCategory(ctx context.Context, categoryID, attributeID int) error {
	_, err := s.pool.Exec(ctx,
		"DELETE FROM unified_category_attributes WHERE category_id = $1 AND attribute_id = $2",
		categoryID, attributeID,
	)
	if err != nil {
		return fmt.Errorf("failed to detach attribute from category: %w", err)
	}

	// Инвалидируем кеш для этой категории
	s.InvalidateCache(categoryID)

	return nil
}

// UpdateCategoryAttribute обновляет параметры связи атрибута с категорией
func (s *unifiedAttributeStorage) UpdateCategoryAttribute(ctx context.Context, categoryID, attributeID int, isRequired, isFilter *bool, displayOrder *int, groupID *int) error {
	// Строим динамический запрос обновления
	var setClause []string
	var args []interface{}
	argNum := 1

	if isRequired != nil {
		setClause = append(setClause, fmt.Sprintf("is_required = $%d", argNum))
		args = append(args, *isRequired)
		argNum++
	}

	if isFilter != nil {
		setClause = append(setClause, fmt.Sprintf("is_filter = $%d", argNum))
		args = append(args, *isFilter)
		argNum++
	}

	if displayOrder != nil {
		setClause = append(setClause, fmt.Sprintf("sort_order = $%d", argNum))
		args = append(args, *displayOrder)
		argNum++
	}

	if groupID != nil {
		// TODO: Добавить поддержку группы, когда будет готова таблица
		// setClause = append(setClause, fmt.Sprintf("group_id = $%d", argNum))
		// args = append(args, *groupID)
		// argNum++
	}

	if len(setClause) == 0 {
		return nil // Нечего обновлять
	}

	// Добавляем updated_at
	setClause = append(setClause, fmt.Sprintf("updated_at = $%d", argNum))
	args = append(args, time.Now())
	argNum++

	// Добавляем условия WHERE
	args = append(args, categoryID, attributeID)

	query := fmt.Sprintf(
		"UPDATE unified_category_attributes SET %s WHERE category_id = $%d AND attribute_id = $%d",
		strings.Join(setClause, ", "),
		argNum,
		argNum+1,
	)

	_, err := s.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update category attribute: %w", err)
	}

	// Инвалидируем кеш для этой категории
	s.InvalidateCache(categoryID)
	return nil
}

// GetCategoryAttributes получает атрибуты категории
func (s *unifiedAttributeStorage) GetCategoryAttributes(ctx context.Context, categoryID int) ([]*models.UnifiedAttribute, error) {
	// Проверяем кеш
	s.cacheMutex.RLock()
	cachedAttrs, hasCached := s.cache[categoryID]
	cacheTime, hasTime := s.cacheTime[categoryID]
	s.cacheMutex.RUnlock()

	// Если данные в кеше и они не устарели
	if hasCached && hasTime && time.Since(cacheTime) < s.cacheTTL {
		log.Printf("Using cached attributes for category %d", categoryID)
		return cachedAttrs, nil
	}

	// Получаем из БД
	query := `
		SELECT 
			ua.id, ua.code, ua.name, ua.display_name, ua.attribute_type, ua.purpose,
			ua.options, ua.validation_rules, ua.ui_settings,
			ua.is_searchable, ua.is_filterable, uca.is_required,
			ua.affects_stock, ua.affects_price, uca.sort_order, ua.is_active,
			ua.created_at, ua.updated_at,
			ua.legacy_category_attribute_id, ua.legacy_product_variant_attribute_id
		FROM unified_attributes ua
		INNER JOIN unified_category_attributes uca ON ua.id = uca.attribute_id
		WHERE uca.category_id = $1 AND uca.is_enabled = true AND ua.is_active = true
		ORDER BY uca.sort_order, ua.sort_order, ua.name`

	rows, err := s.pool.Query(ctx, query, categoryID)
	if err != nil {
		// Если включен fallback и новая система не работает - используем старую
		if s.useLegacyFallback {
			return s.getCategoryAttributesFromLegacy(ctx, categoryID)
		}
		return nil, fmt.Errorf("failed to get category attributes: %w", err)
	}
	defer rows.Close()

	var attributes []*models.UnifiedAttribute
	for rows.Next() {
		attr := &models.UnifiedAttribute{}
		var categoryIsRequired bool
		var categorySortOrder int

		err := rows.Scan(
			&attr.ID, &attr.Code, &attr.Name, &attr.DisplayName,
			&attr.AttributeType, &attr.Purpose,
			&attr.Options, &attr.ValidationRules, &attr.UISettings,
			&attr.IsSearchable, &attr.IsFilterable, &categoryIsRequired,
			&attr.AffectsStock, &attr.AffectsPrice, &categorySortOrder, &attr.IsActive,
			&attr.CreatedAt, &attr.UpdatedAt,
			&attr.LegacyCategoryAttributeID, &attr.LegacyProductVariantAttributeID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan attribute: %w", err)
		}

		// Переопределяем значения из связи с категорией
		attr.IsRequired = categoryIsRequired
		attr.SortOrder = categorySortOrder

		attributes = append(attributes, attr)
	}

	// Если новая система пустая и включен fallback - используем старую
	if len(attributes) == 0 && s.useLegacyFallback {
		return s.getCategoryAttributesFromLegacy(ctx, categoryID)
	}

	// Сохраняем в кеш
	s.cacheMutex.Lock()
	s.cache[categoryID] = attributes
	s.cacheTime[categoryID] = time.Now()
	s.cacheMutex.Unlock()

	return attributes, nil
}

// getCategoryAttributesFromLegacy получает атрибуты из старой системы (fallback)
func (s *unifiedAttributeStorage) getCategoryAttributesFromLegacy(ctx context.Context, categoryID int) ([]*models.UnifiedAttribute, error) {
	query := `
		SELECT 
			ca.id, ca.name, ca.display_name, ca.attribute_type,
			ca.options, ca.validation_rules,
			ca.is_searchable, ca.is_filterable, cam.is_required,
			ca.is_variant_compatible, ca.affects_stock,
			cam.sort_order, ca.created_at
		FROM category_attributes ca
		INNER JOIN category_attribute_mapping cam ON ca.id = cam.attribute_id
		WHERE cam.category_id = $1 AND cam.is_enabled = true
		ORDER BY cam.sort_order, ca.sort_order, ca.name`

	rows, err := s.pool.Query(ctx, query, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get legacy category attributes: %w", err)
	}
	defer rows.Close()

	var attributes []*models.UnifiedAttribute
	for rows.Next() {
		var legacyID int
		var name, displayName, attrType string
		var options, validationRules json.RawMessage
		var isSearchable, isFilterable, isRequired, isVariantCompatible, affectsStock bool
		var sortOrder int
		var createdAt time.Time

		err := rows.Scan(
			&legacyID, &name, &displayName, &attrType,
			&options, &validationRules,
			&isSearchable, &isFilterable, &isRequired,
			&isVariantCompatible, &affectsStock,
			&sortOrder, &createdAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan legacy attribute: %w", err)
		}

		// Преобразуем в унифицированный формат
		purpose := models.PurposeRegular
		if isVariantCompatible {
			purpose = models.PurposeBoth
		}

		attr := &models.UnifiedAttribute{
			ID:                        legacyID, // Используем старый ID для совместимости
			Code:                      fmt.Sprintf("legacy_%d", legacyID),
			Name:                      name,
			DisplayName:               displayName,
			AttributeType:             attrType,
			Purpose:                   purpose,
			Options:                   options,
			ValidationRules:           validationRules,
			IsSearchable:              isSearchable,
			IsFilterable:              isFilterable,
			IsRequired:                isRequired,
			AffectsStock:              affectsStock,
			SortOrder:                 sortOrder,
			IsActive:                  true,
			CreatedAt:                 createdAt,
			LegacyCategoryAttributeID: &legacyID,
		}

		attributes = append(attributes, attr)
	}

	return attributes, nil
}

// GetCategoryAttributesWithSettings получает атрибуты категории с настройками
func (s *unifiedAttributeStorage) GetCategoryAttributesWithSettings(ctx context.Context, categoryID int) ([]*models.UnifiedCategoryAttribute, error) {
	query := `
		SELECT 
			uca.id, uca.category_id, uca.attribute_id,
			uca.is_enabled, uca.is_required, uca.sort_order,
			uca.category_specific_options, uca.created_at, uca.updated_at,
			ua.id, ua.code, ua.name, ua.display_name, ua.attribute_type, ua.purpose,
			ua.options, ua.validation_rules, ua.ui_settings,
			ua.is_searchable, ua.is_filterable, ua.is_required,
			ua.affects_stock, ua.affects_price, ua.sort_order, ua.is_active,
			ua.created_at, ua.updated_at
		FROM unified_category_attributes uca
		INNER JOIN unified_attributes ua ON uca.attribute_id = ua.id
		WHERE uca.category_id = $1 AND uca.is_enabled = true AND ua.is_active = true
		ORDER BY uca.sort_order, ua.sort_order, ua.name`

	rows, err := s.pool.Query(ctx, query, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get category attributes with settings: %w", err)
	}
	defer rows.Close()

	var result []*models.UnifiedCategoryAttribute
	for rows.Next() {
		ca := &models.UnifiedCategoryAttribute{
			Attribute: &models.UnifiedAttribute{},
		}

		err := rows.Scan(
			&ca.ID, &ca.CategoryID, &ca.AttributeID,
			&ca.IsEnabled, &ca.IsRequired, &ca.SortOrder,
			&ca.CategorySpecificOptions, &ca.CreatedAt, &ca.UpdatedAt,
			&ca.Attribute.ID, &ca.Attribute.Code, &ca.Attribute.Name, &ca.Attribute.DisplayName,
			&ca.Attribute.AttributeType, &ca.Attribute.Purpose,
			&ca.Attribute.Options, &ca.Attribute.ValidationRules, &ca.Attribute.UISettings,
			&ca.Attribute.IsSearchable, &ca.Attribute.IsFilterable, &ca.Attribute.IsRequired,
			&ca.Attribute.AffectsStock, &ca.Attribute.AffectsPrice, &ca.Attribute.SortOrder, &ca.Attribute.IsActive,
			&ca.Attribute.CreatedAt, &ca.Attribute.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category attribute: %w", err)
		}

		result = append(result, ca)
	}

	return result, nil
}

// SaveAttributeValue сохраняет значение атрибута
func (s *unifiedAttributeStorage) SaveAttributeValue(ctx context.Context, value *models.UnifiedAttributeValue) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Сохраняем в новую систему
	query := `
		INSERT INTO unified_attribute_values (
			entity_type, entity_id, attribute_id,
			text_value, numeric_value, boolean_value, date_value, json_value
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (entity_type, entity_id, attribute_id)
		DO UPDATE SET
			text_value = EXCLUDED.text_value,
			numeric_value = EXCLUDED.numeric_value,
			boolean_value = EXCLUDED.boolean_value,
			date_value = EXCLUDED.date_value,
			json_value = EXCLUDED.json_value,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id`

	err = tx.QueryRow(ctx, query,
		value.EntityType, value.EntityID, value.AttributeID,
		value.TextValue, value.NumericValue, value.BooleanValue,
		value.DateValue, value.JSONValue,
	).Scan(&value.ID)
	if err != nil {
		return fmt.Errorf("failed to save attribute value: %w", err)
	}

	// Если включен fallback и это объявление - дублируем в старую систему
	if s.useLegacyFallback && value.EntityType == models.AttributeEntityTypeListing {
		// Получаем legacy ID атрибута
		var legacyAttrID *int
		err = tx.QueryRow(ctx,
			"SELECT legacy_category_attribute_id FROM unified_attributes WHERE id = $1",
			value.AttributeID,
		).Scan(&legacyAttrID)

		if err == nil && legacyAttrID != nil {
			// Сохраняем в старую таблицу
			_, err = tx.Exec(ctx, `
				INSERT INTO listing_attribute_values (
					listing_id, attribute_id,
					text_value, numeric_value, boolean_value, json_value
				) VALUES ($1, $2, $3, $4, $5, $6)
				ON CONFLICT (listing_id, attribute_id)
				DO UPDATE SET
					text_value = EXCLUDED.text_value,
					numeric_value = EXCLUDED.numeric_value,
					boolean_value = EXCLUDED.boolean_value,
					json_value = EXCLUDED.json_value`,
				value.EntityID, *legacyAttrID,
				value.TextValue, value.NumericValue, value.BooleanValue, value.JSONValue,
			)
			if err != nil {
				// Логируем ошибку, но не прерываем - старая система не критична
				log.Printf("Failed to save to legacy system: %v", err)
			}
		}
	}

	return tx.Commit(ctx)
}

// GetAttributeValues получает значения атрибутов для сущности
func (s *unifiedAttributeStorage) GetAttributeValues(ctx context.Context, entityType models.AttributeEntityType, entityID int) ([]*models.UnifiedAttributeValue, error) {
	query := `
		SELECT 
			uav.id, uav.entity_type, uav.entity_id, uav.attribute_id,
			uav.text_value, uav.numeric_value, uav.boolean_value,
			uav.date_value, uav.json_value,
			uav.created_at, uav.updated_at,
			ua.id, ua.code, ua.name, ua.display_name, ua.attribute_type
		FROM unified_attribute_values uav
		INNER JOIN unified_attributes ua ON uav.attribute_id = ua.id
		WHERE uav.entity_type = $1 AND uav.entity_id = $2`

	rows, err := s.pool.Query(ctx, query, entityType, entityID)
	if err != nil {
		// Если включен fallback и это объявление - пробуем старую систему
		if s.useLegacyFallback && entityType == models.AttributeEntityTypeListing {
			return s.getListingAttributeValuesFromLegacy(ctx, entityID)
		}
		return nil, fmt.Errorf("failed to get attribute values: %w", err)
	}
	defer rows.Close()

	var values []*models.UnifiedAttributeValue
	for rows.Next() {
		value := &models.UnifiedAttributeValue{
			Attribute: &models.UnifiedAttribute{},
		}

		err := rows.Scan(
			&value.ID, &value.EntityType, &value.EntityID, &value.AttributeID,
			&value.TextValue, &value.NumericValue, &value.BooleanValue,
			&value.DateValue, &value.JSONValue,
			&value.CreatedAt, &value.UpdatedAt,
			&value.Attribute.ID, &value.Attribute.Code, &value.Attribute.Name,
			&value.Attribute.DisplayName, &value.Attribute.AttributeType,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan attribute value: %w", err)
		}

		// Вычисляем отображаемое значение
		value.DisplayValue = value.GetDisplayValue()

		values = append(values, value)
	}

	// Если новая система пустая и включен fallback - используем старую
	if len(values) == 0 && s.useLegacyFallback && entityType == models.AttributeEntityTypeListing {
		return s.getListingAttributeValuesFromLegacy(ctx, entityID)
	}

	return values, nil
}

// getListingAttributeValuesFromLegacy получает значения атрибутов из старой системы
func (s *unifiedAttributeStorage) getListingAttributeValuesFromLegacy(ctx context.Context, listingID int) ([]*models.UnifiedAttributeValue, error) {
	query := `
		SELECT 
			lav.listing_id, lav.attribute_id,
			lav.text_value, lav.numeric_value, lav.boolean_value, lav.json_value,
			ca.name, ca.display_name, ca.attribute_type
		FROM listing_attribute_values lav
		INNER JOIN category_attributes ca ON lav.attribute_id = ca.id
		WHERE lav.listing_id = $1`

	rows, err := s.pool.Query(ctx, query, listingID)
	if err != nil {
		return nil, fmt.Errorf("failed to get legacy listing attribute values: %w", err)
	}
	defer rows.Close()

	var values []*models.UnifiedAttributeValue
	for rows.Next() {
		var legacyAttrID int
		var textValue *string
		var numericValue *float64
		var booleanValue *bool
		var jsonValue json.RawMessage
		var attrName, attrDisplayName, attrType string

		err := rows.Scan(
			&listingID, &legacyAttrID,
			&textValue, &numericValue, &booleanValue, &jsonValue,
			&attrName, &attrDisplayName, &attrType,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan legacy attribute value: %w", err)
		}

		value := &models.UnifiedAttributeValue{
			EntityType:   models.AttributeEntityTypeListing,
			EntityID:     listingID,
			AttributeID:  legacyAttrID, // Используем старый ID
			TextValue:    textValue,
			NumericValue: numericValue,
			BooleanValue: booleanValue,
			JSONValue:    jsonValue,
			Attribute: &models.UnifiedAttribute{
				ID:            legacyAttrID,
				Name:          attrName,
				DisplayName:   attrDisplayName,
				AttributeType: attrType,
			},
		}

		// Вычисляем отображаемое значение
		value.DisplayValue = value.GetDisplayValue()

		values = append(values, value)
	}

	return values, nil
}

// DeleteAttributeValues удаляет все значения атрибутов для сущности
func (s *unifiedAttributeStorage) DeleteAttributeValues(ctx context.Context, entityType models.AttributeEntityType, entityID int) error {
	_, err := s.pool.Exec(ctx,
		"DELETE FROM unified_attribute_values WHERE entity_type = $1 AND entity_id = $2",
		entityType, entityID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete attribute values: %w", err)
	}

	// Если включен fallback и это объявление - удаляем из старой системы
	if s.useLegacyFallback && entityType == models.AttributeEntityTypeListing {
		_, err = s.pool.Exec(ctx,
			"DELETE FROM listing_attribute_values WHERE listing_id = $1",
			entityID,
		)
		if err != nil {
			// Логируем ошибку, но не прерываем
			log.Printf("Failed to delete from legacy system: %v", err)
		}
	}

	return nil
}

// GetAttributeByLegacyID получает атрибут по старому ID
func (s *unifiedAttributeStorage) GetAttributeByLegacyID(ctx context.Context, legacyID int, isProductVariant bool) (*models.UnifiedAttribute, error) {
	var query string
	if isProductVariant {
		query = `
			SELECT 
				id, code, name, display_name, attribute_type, purpose,
				options, validation_rules, ui_settings,
				is_searchable, is_filterable, is_required,
				affects_stock, affects_price, sort_order, is_active,
				created_at, updated_at
			FROM unified_attributes
			WHERE legacy_product_variant_attribute_id = $1`
	} else {
		query = `
			SELECT 
				id, code, name, display_name, attribute_type, purpose,
				options, validation_rules, ui_settings,
				is_searchable, is_filterable, is_required,
				affects_stock, affects_price, sort_order, is_active,
				created_at, updated_at
			FROM unified_attributes
			WHERE legacy_category_attribute_id = $1`
	}

	attr := &models.UnifiedAttribute{}
	err := s.pool.QueryRow(ctx, query, legacyID).Scan(
		&attr.ID, &attr.Code, &attr.Name, &attr.DisplayName,
		&attr.AttributeType, &attr.Purpose,
		&attr.Options, &attr.ValidationRules, &attr.UISettings,
		&attr.IsSearchable, &attr.IsFilterable, &attr.IsRequired,
		&attr.AffectsStock, &attr.AffectsPrice, &attr.SortOrder, &attr.IsActive,
		&attr.CreatedAt, &attr.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("attribute not found for legacy ID %d", legacyID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get attribute by legacy ID: %w", err)
	}

	return attr, nil
}

// MigrateFromLegacySystem мигрирует данные из старой системы
func (s *unifiedAttributeStorage) MigrateFromLegacySystem(ctx context.Context) error {
	// Эта функция уже реализована через миграции SQL
	// Здесь можно добавить дополнительную логику если нужно
	return nil
}

// InvalidateCache инвалидирует кеш для категории
func (s *unifiedAttributeStorage) InvalidateCache(categoryID int) {
	s.cacheMutex.Lock()
	delete(s.cache, categoryID)
	delete(s.cacheTime, categoryID)
	s.cacheMutex.Unlock()
}

// GetVariantCompatibleAttributes получает атрибуты, которые могут быть вариантами
func (s *unifiedAttributeStorage) GetVariantCompatibleAttributes(ctx context.Context) ([]*models.UnifiedAttribute, error) {
	query := `
		SELECT 
			id, code, name, display_name, attribute_type, purpose,
			options, validation_rules, ui_settings,
			is_searchable, is_filterable, is_required,
			is_variant_compatible, affects_stock, affects_price,
			sort_order, is_active, created_at, updated_at
		FROM unified_attributes
		WHERE is_variant_compatible = true AND is_active = true
		ORDER BY sort_order, name`

	var attributes []*models.UnifiedAttribute
	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query variant compatible attributes: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var attr models.UnifiedAttribute
		err := rows.Scan(
			&attr.ID, &attr.Code, &attr.Name, &attr.DisplayName,
			&attr.AttributeType, &attr.Purpose, &attr.Options,
			&attr.ValidationRules, &attr.UISettings,
			&attr.IsSearchable, &attr.IsFilterable, &attr.IsRequired,
			&attr.IsVariantCompatible, &attr.AffectsStock, &attr.AffectsPrice,
			&attr.SortOrder, &attr.IsActive, &attr.CreatedAt, &attr.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan attribute: %w", err)
		}
		attributes = append(attributes, &attr)
	}

	return attributes, nil
}

// GetCategoryVariantMappings получает вариативные атрибуты для категории
func (s *unifiedAttributeStorage) GetCategoryVariantMappings(ctx context.Context, categoryID int) ([]*models.VariantAttributeMapping, error) {
	query := `
		SELECT 
			vam.id, vam.variant_attribute_id, vam.category_id,
			vam.sort_order, vam.is_required, vam.created_at, vam.updated_at,
			ua.id as "attribute.id", ua.code as "attribute.code",
			ua.name as "attribute.name", ua.display_name as "attribute.display_name",
			ua.attribute_type as "attribute.attribute_type",
			ua.options as "attribute.options",
			ua.affects_stock as "attribute.affects_stock",
			ua.affects_price as "attribute.affects_price"
		FROM variant_attribute_mappings vam
		JOIN unified_attributes ua ON vam.variant_attribute_id = ua.id
		WHERE vam.category_id = $1
		ORDER BY vam.sort_order, ua.name`

	rows, err := s.pool.Query(ctx, query, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get category variant mappings: %w", err)
	}
	defer rows.Close()

	var mappings []*models.VariantAttributeMapping
	for rows.Next() {
		var m models.VariantAttributeMapping
		var attr models.UnifiedAttribute

		err := rows.Scan(
			&m.ID, &m.VariantAttributeID, &m.CategoryID,
			&m.SortOrder, &m.IsRequired, &m.CreatedAt, &m.UpdatedAt,
			&attr.ID, &attr.Code, &attr.Name, &attr.DisplayName,
			&attr.AttributeType, &attr.Options,
			&attr.AffectsStock, &attr.AffectsPrice,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan mapping: %w", err)
		}

		m.Attribute = &attr
		mappings = append(mappings, &m)
	}

	return mappings, nil
}

// CreateVariantMapping создает новую связь
func (s *unifiedAttributeStorage) CreateVariantMapping(ctx context.Context, mapping *models.VariantAttributeMappingCreateRequest) (*models.VariantAttributeMapping, error) {
	query := `
		INSERT INTO variant_attribute_mappings 
		(variant_attribute_id, category_id, sort_order, is_required)
		VALUES ($1, $2, $3, $4)
		RETURNING id, variant_attribute_id, category_id, sort_order, is_required, created_at, updated_at`

	var m models.VariantAttributeMapping
	err := s.pool.QueryRow(ctx, query,
		mapping.VariantAttributeID,
		mapping.CategoryID,
		mapping.SortOrder,
		mapping.IsRequired,
	).Scan(
		&m.ID, &m.VariantAttributeID, &m.CategoryID,
		&m.SortOrder, &m.IsRequired, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create variant mapping: %w", err)
	}

	return &m, nil
}

// UpdateVariantMapping обновляет связь
func (s *unifiedAttributeStorage) UpdateVariantMapping(ctx context.Context, id int, update *models.VariantAttributeMappingUpdateRequest) error {
	query := `
		UPDATE variant_attribute_mappings 
		SET sort_order = COALESCE($2, sort_order),
		    is_required = COALESCE($3, is_required),
		    updated_at = NOW()
		WHERE id = $1`

	_, err := s.pool.Exec(ctx, query, id, update.SortOrder, update.IsRequired)
	if err != nil {
		return fmt.Errorf("failed to update variant mapping: %w", err)
	}

	return nil
}

// DeleteVariantMapping удаляет связь
func (s *unifiedAttributeStorage) DeleteVariantMapping(ctx context.Context, id int) error {
	query := `DELETE FROM variant_attribute_mappings WHERE id = $1`

	_, err := s.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete variant mapping: %w", err)
	}

	return nil
}

// DeleteCategoryVariantMappings удаляет все вариативные атрибуты категории
func (s *unifiedAttributeStorage) DeleteCategoryVariantMappings(ctx context.Context, categoryID int) error {
	query := `DELETE FROM variant_attribute_mappings WHERE category_id = $1`

	_, err := s.pool.Exec(ctx, query, categoryID)
	if err != nil {
		return fmt.Errorf("failed to delete category variant mappings: %w", err)
	}

	return nil
}
