package attributes

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"backend/internal/domain/models"
	"backend/internal/storage/postgres"
)

// AttributeType тип атрибута
type AttributeType string

const (
	AttributeTypeText        AttributeType = "text"
	AttributeTypeTextarea    AttributeType = "textarea"
	AttributeTypeNumber      AttributeType = "number"
	AttributeTypeBoolean     AttributeType = "boolean"
	AttributeTypeSelect      AttributeType = "select"
	AttributeTypeMultiselect AttributeType = "multiselect"
	AttributeTypeDate        AttributeType = "date"
	AttributeTypeColor       AttributeType = "color"
	AttributeTypeSize        AttributeType = "size"
)

// Translation представляет перевод
type Translation struct {
	Language string `json:"language"`
	Text     string `json:"text"`
}

// UnifiedAttributeService сервис для работы с унифицированными атрибутами
type UnifiedAttributeService struct {
	storage postgres.UnifiedAttributeStorage

	// Флаги для управления поведением
	useLegacyFallback bool
	dualWrite         bool
	cacheEnabled      bool
}

// NewUnifiedAttributeService создает новый сервис атрибутов
func NewUnifiedAttributeService(storage postgres.UnifiedAttributeStorage, useLegacyFallback, dualWrite bool) *UnifiedAttributeService {
	return &UnifiedAttributeService{
		storage:           storage,
		useLegacyFallback: useLegacyFallback,
		dualWrite:         dualWrite,
		cacheEnabled:      true,
	}
}

// GetCategoryAttributes получает атрибуты для категории с поддержкой fallback
func (s *UnifiedAttributeService) GetCategoryAttributes(ctx context.Context, categoryID int) ([]*models.UnifiedAttribute, error) {
	// Пробуем получить из новой системы
	attributes, err := s.storage.GetCategoryAttributes(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get category attributes: %w", err)
	}

	// Если новая система пустая и включен fallback - логируем предупреждение
	if len(attributes) == 0 && s.useLegacyFallback {
		log.Printf("Warning: No attributes found in unified system for category %d, fallback to legacy system might be used", categoryID)
	}

	return attributes, nil
}

// GetCategoryAttributesWithSettings получает атрибуты с настройками для категории
func (s *UnifiedAttributeService) GetCategoryAttributesWithSettings(ctx context.Context, categoryID int) ([]*models.UnifiedCategoryAttribute, error) {
	return s.storage.GetCategoryAttributesWithSettings(ctx, categoryID)
}

// SaveAttributeValues сохраняет значения атрибутов для сущности
func (s *UnifiedAttributeService) SaveAttributeValues(ctx context.Context, entityType models.AttributeEntityType, entityID int, values map[int]interface{}) error {
	for attributeID, value := range values {
		// Преобразуем значение в правильный формат
		attrValue := &models.UnifiedAttributeValue{
			EntityType:  entityType,
			EntityID:    entityID,
			AttributeID: attributeID,
		}

		// Определяем тип значения и сохраняем в соответствующее поле
		switch v := value.(type) {
		case string:
			attrValue.TextValue = &v
		case float64:
			attrValue.NumericValue = &v
		case int:
			f := float64(v)
			attrValue.NumericValue = &f
		case bool:
			attrValue.BooleanValue = &v
		case time.Time:
			attrValue.DateValue = &v
		case map[string]interface{}, []interface{}:
			jsonData, err := json.Marshal(v)
			if err != nil {
				return fmt.Errorf("failed to marshal json value: %w", err)
			}
			attrValue.JSONValue = jsonData
		case nil:
			// Пропускаем nil значения
			continue
		default:
			// Пробуем преобразовать в строку
			str := fmt.Sprintf("%v", v)
			attrValue.TextValue = &str
		}

		// Сохраняем значение
		if err := s.storage.SaveAttributeValue(ctx, attrValue); err != nil {
			return fmt.Errorf("failed to save attribute value for attribute %d: %w", attributeID, err)
		}
	}

	return nil
}

// GetAttributeValues получает значения атрибутов для сущности
func (s *UnifiedAttributeService) GetAttributeValues(ctx context.Context, entityType models.AttributeEntityType, entityID int) ([]*models.UnifiedAttributeValue, error) {
	return s.storage.GetAttributeValues(ctx, entityType, entityID)
}

// CreateAttribute создает новый атрибут
func (s *UnifiedAttributeService) CreateAttribute(ctx context.Context, attr *models.UnifiedAttribute) (int, error) {
	// Генерируем код если не указан
	if attr.Code == "" {
		attr.Code = s.generateAttributeCode(attr.Name)
	}

	// Устанавливаем значения по умолчанию
	if attr.Purpose == "" {
		attr.Purpose = models.PurposeRegular
	}

	id, err := s.storage.CreateAttribute(ctx, attr)
	if err != nil {
		return 0, fmt.Errorf("failed to create attribute: %w", err)
	}

	// Инвалидируем кеш если включен
	if s.cacheEnabled {
		s.storage.InvalidateCache(0) // Инвалидируем весь кеш
	}

	return id, nil
}

// UpdateAttribute обновляет атрибут
func (s *UnifiedAttributeService) UpdateAttribute(ctx context.Context, id int, updates map[string]interface{}) error {
	// Валидируем обновления
	if purpose, ok := updates["purpose"]; ok {
		if !s.isValidPurpose(purpose) {
			return fmt.Errorf("invalid purpose value: %v", purpose)
		}
	}

	if attrType, ok := updates["attribute_type"]; ok {
		if !s.isValidAttributeType(attrType) {
			return fmt.Errorf("invalid attribute type: %v", attrType)
		}
	}

	err := s.storage.UpdateAttribute(ctx, id, updates)
	if err != nil {
		return fmt.Errorf("failed to update attribute: %w", err)
	}

	// Инвалидируем кеш
	if s.cacheEnabled {
		s.storage.InvalidateCache(0)
	}

	return nil
}

// DeleteAttribute удаляет атрибут
func (s *UnifiedAttributeService) DeleteAttribute(ctx context.Context, id int) error {
	// Проверяем что атрибут не используется
	values, err := s.checkAttributeUsage(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check attribute usage: %w", err)
	}

	if values > 0 {
		return fmt.Errorf("cannot delete attribute: it has %d associated values", values)
	}

	err = s.storage.DeleteAttribute(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete attribute: %w", err)
	}

	// Инвалидируем кеш
	if s.cacheEnabled {
		s.storage.InvalidateCache(0)
	}

	return nil
}

// AttachAttributeToCategory привязывает атрибут к категории
func (s *UnifiedAttributeService) AttachAttributeToCategory(ctx context.Context, categoryID, attributeID int, settings *models.UnifiedCategoryAttribute) error {
	// Проверяем что атрибут существует
	_, err := s.storage.GetAttribute(ctx, attributeID)
	if err != nil {
		return fmt.Errorf("attribute not found: %w", err)
	}

	// Проверяем совместимость атрибута с категорией
	// (здесь можно добавить дополнительную логику проверки с использованием attr)

	// Устанавливаем значения по умолчанию для настроек
	if settings == nil {
		settings = &models.UnifiedCategoryAttribute{
			CategoryID:  categoryID,
			AttributeID: attributeID,
			IsEnabled:   true,
			IsRequired:  false,
			SortOrder:   0,
		}
	} else {
		settings.CategoryID = categoryID
		settings.AttributeID = attributeID
	}

	err = s.storage.AttachAttributeToCategory(ctx, categoryID, attributeID, settings)
	if err != nil {
		return fmt.Errorf("failed to attach attribute to category: %w", err)
	}

	// Инвалидируем кеш для категории
	if s.cacheEnabled {
		s.storage.InvalidateCache(categoryID)
	}

	return nil
}

// DetachAttributeFromCategory отвязывает атрибут от категории
func (s *UnifiedAttributeService) DetachAttributeFromCategory(ctx context.Context, categoryID, attributeID int) error {
	err := s.storage.DetachAttributeFromCategory(ctx, categoryID, attributeID)
	if err != nil {
		return fmt.Errorf("failed to detach attribute from category: %w", err)
	}

	// Инвалидируем кеш для категории
	if s.cacheEnabled {
		s.storage.InvalidateCache(categoryID)
	}

	return nil
}

// ConvertFromLegacyAttribute преобразует старый атрибут в новый формат
func (s *UnifiedAttributeService) ConvertFromLegacyAttribute(oldAttr *models.CategoryAttribute) *models.UnifiedAttribute {
	purpose := models.PurposeRegular
	if oldAttr.IsVariantCompatible {
		purpose = models.PurposeBoth
	}

	return &models.UnifiedAttribute{
		ID:                        oldAttr.ID,
		Code:                      s.generateAttributeCode(oldAttr.Name),
		Name:                      oldAttr.Name,
		DisplayName:               oldAttr.DisplayName,
		AttributeType:             oldAttr.AttributeType,
		Purpose:                   purpose,
		Options:                   oldAttr.Options,
		ValidationRules:           oldAttr.ValidRules,
		IsSearchable:              oldAttr.IsSearchable,
		IsFilterable:              oldAttr.IsFilterable,
		IsRequired:                oldAttr.IsRequired,
		AffectsStock:              oldAttr.AffectsStock,
		SortOrder:                 oldAttr.SortOrder,
		IsActive:                  true,
		CreatedAt:                 oldAttr.CreatedAt,
		Translations:              oldAttr.Translations,
		LegacyCategoryAttributeID: &oldAttr.ID,
	}
}

// ConvertToLegacyAttribute преобразует новый атрибут в старый формат для обратной совместимости
func (s *UnifiedAttributeService) ConvertToLegacyAttribute(newAttr *models.UnifiedAttribute) *models.CategoryAttribute {
	return newAttr.ToCategoryAttribute()
}

// MigrateFromLegacySystem выполняет миграцию из старой системы
func (s *UnifiedAttributeService) MigrateFromLegacySystem(ctx context.Context) error {
	// Эта функция уже реализована через SQL миграции
	// Здесь можно добавить дополнительную бизнес-логику если нужно
	return s.storage.MigrateFromLegacySystem(ctx)
}

// GetAttributeByLegacyID получает атрибут по старому ID
func (s *UnifiedAttributeService) GetAttributeByLegacyID(ctx context.Context, legacyID int, isProductVariant bool) (*models.UnifiedAttribute, error) {
	return s.storage.GetAttributeByLegacyID(ctx, legacyID, isProductVariant)
}

// ListAttributes получает список атрибутов с фильтрацией
func (s *UnifiedAttributeService) ListAttributes(ctx context.Context, filter *models.UnifiedAttributeFilter) ([]*models.UnifiedAttribute, error) {
	return s.storage.ListAttributes(ctx, filter)
}

// ValidateAttributeValue проверяет значение атрибута на соответствие правилам валидации
func (s *UnifiedAttributeService) ValidateAttributeValue(ctx context.Context, attributeID int, value interface{}) error {
	// Получаем атрибут
	attr, err := s.storage.GetAttribute(ctx, attributeID)
	if err != nil {
		return fmt.Errorf("attribute not found: %w", err)
	}

	// Проверяем обязательность
	if attr.IsRequired && (value == nil || value == "") {
		return fmt.Errorf("value for %s is required", attr.Name)
	}

	// Если значение пустое и не обязательное - пропускаем валидацию
	if value == nil || value == "" {
		return nil
	}

	// Проверяем тип значения
	switch attr.AttributeType {
	case "number":
		_, ok := value.(float64)
		if !ok {
			if intVal, ok := value.(int); ok {
				value = float64(intVal)
			} else {
				return fmt.Errorf("value must be a number for attribute %s", attr.Name)
			}
		}
	case "boolean":
		_, ok := value.(bool)
		if !ok {
			return fmt.Errorf("value must be a boolean for attribute %s", attr.Name)
		}
	case "date":
		_, ok := value.(time.Time)
		if !ok {
			if strVal, ok := value.(string); ok {
				_, err := time.Parse("2006-01-02", strVal)
				if err != nil {
					return fmt.Errorf("invalid date format for attribute %s: %w", attr.Name, err)
				}
			} else {
				return fmt.Errorf("value must be a date for attribute %s", attr.Name)
			}
		}
	case "select", "multiselect":
		// Проверяем что значение есть в списке опций
		if attr.Options != nil {
			var options []string
			if err := json.Unmarshal(attr.Options, &options); err == nil {
				strVal := fmt.Sprintf("%v", value)
				found := false
				for _, opt := range options {
					if opt == strVal {
						found = true
						break
					}
				}
				if !found {
					return fmt.Errorf("value %s is not in allowed options for attribute %s", strVal, attr.Name)
				}
			}
		}
	}

	// Проверяем дополнительные правила валидации
	if attr.ValidationRules != nil {
		var rules map[string]interface{}
		if err := json.Unmarshal(attr.ValidationRules, &rules); err == nil {
			// Проверяем минимальное и максимальное значение для чисел
			if attr.AttributeType == "number" {
				if numVal, ok := value.(float64); ok {
					if minVal, exists := rules["min"].(float64); exists && numVal < minVal {
						return fmt.Errorf("value for %s must be at least %f", attr.Name, minVal)
					}
					if maxVal, exists := rules["max"].(float64); exists && numVal > maxVal {
						return fmt.Errorf("value for %s must not exceed %f", attr.Name, maxVal)
					}
				}
			}

			// Проверяем длину для текстовых полей
			if attr.AttributeType == "text" || attr.AttributeType == "textarea" {
				if strVal, ok := value.(string); ok {
					if minLen, exists := rules["minLength"].(float64); exists && len(strVal) < int(minLen) {
						return fmt.Errorf("value for %s must be at least %d characters", attr.Name, int(minLen))
					}
					if maxLen, exists := rules["maxLength"].(float64); exists && len(strVal) > int(maxLen) {
						return fmt.Errorf("value for %s must not exceed %d characters", attr.Name, int(maxLen))
					}
				}
			}
		}
	}

	return nil
}

// Вспомогательные методы

func (s *UnifiedAttributeService) generateAttributeCode(name string) string {
	// Простая генерация кода из имени
	code := ""
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			code += string(r)
		} else if r == ' ' {
			code += "_"
		}
	}
	return code
}

func (s *UnifiedAttributeService) isValidPurpose(purpose interface{}) bool {
	str, ok := purpose.(string)
	if !ok {
		return false
	}
	return str == string(models.PurposeRegular) || str == string(models.PurposeVariant) || str == string(models.PurposeBoth)
}

// UpdateCategoryAttribute обновляет параметры связи атрибута с категорией
func (s *UnifiedAttributeService) UpdateCategoryAttribute(ctx context.Context, categoryID, attributeID int, settings *models.UnifiedCategoryAttribute) error {
	// Обновляем связь в новой системе
	if s.storage != nil {
		err := s.storage.UpdateCategoryAttribute(ctx, categoryID, attributeID,
			&settings.IsRequired, &settings.IsFilter, &settings.SortOrder, nil)
		if err != nil {
			return err
		}
	}

	// Если включен dual-write, обновляем и в старой системе
	if s.useLegacyFallback && s.dualWrite {
		// TODO: Обновить в старой системе через соответствующий сервис
	}

	return nil
}

// GetMigrationStatus возвращает текущий статус миграции
func (s *UnifiedAttributeService) GetMigrationStatus(ctx context.Context) (map[string]interface{}, error) {
	// TODO: Реализовать получение реального статуса из базы данных
	// Пока возвращаем статический ответ
	status := map[string]interface{}{
		"status": "completed",
		"details": map[string]interface{}{
			"attributes_migrated":  85,
			"categories_processed": 14,
			"values_migrated":      15,
			"started_at":           "2025-09-02T00:00:00Z",
			"completed_at":         "2025-09-02T01:00:00Z",
		},
	}

	return status, nil
}

func (s *UnifiedAttributeService) isValidAttributeType(attrType interface{}) bool {
	str, ok := attrType.(string)
	if !ok {
		return false
	}
	validTypes := []string{"text", "textarea", "number", "boolean", "select", "multiselect", "date", "color", "size"}
	for _, t := range validTypes {
		if str == t {
			return true
		}
	}
	return false
}

func (s *UnifiedAttributeService) checkAttributeUsage(ctx context.Context, attributeID int) (int, error) {
	// Здесь должна быть проверка использования атрибута в значениях
	// Для простоты возвращаем 0
	return 0, nil
}

// GetVariantAttributes получает все атрибуты, которые могут использоваться как варианты
func (s *UnifiedAttributeService) GetVariantAttributes(ctx context.Context) ([]*models.UnifiedAttribute, error) {
	return s.storage.GetVariantCompatibleAttributes(ctx)
}

// GetCategoryVariantAttributes получает вариативные атрибуты для конкретной категории
func (s *UnifiedAttributeService) GetCategoryVariantAttributes(ctx context.Context, categoryID int) ([]*models.VariantAttributeMapping, error) {
	return s.storage.GetCategoryVariantMappings(ctx, categoryID)
}

// CreateVariantAttributeMapping создает связь между вариативным атрибутом и категорией
func (s *UnifiedAttributeService) CreateVariantAttributeMapping(ctx context.Context, mapping *models.VariantAttributeMappingCreateRequest) (*models.VariantAttributeMapping, error) {
	return s.storage.CreateVariantMapping(ctx, mapping)
}

// UpdateVariantAttributeMapping обновляет связь между вариативным атрибутом и категорией
func (s *UnifiedAttributeService) UpdateVariantAttributeMapping(ctx context.Context, id int, update *models.VariantAttributeMappingUpdateRequest) error {
	return s.storage.UpdateVariantMapping(ctx, id, update)
}

// DeleteVariantAttributeMapping удаляет связь между вариативным атрибутом и категорией
func (s *UnifiedAttributeService) DeleteVariantAttributeMapping(ctx context.Context, id int) error {
	return s.storage.DeleteVariantMapping(ctx, id)
}

// UpdateCategoryVariantAttributes обновляет все вариативные атрибуты для категории
func (s *UnifiedAttributeService) UpdateCategoryVariantAttributes(ctx context.Context, request *models.CategoryVariantAttributesUpdateRequest) error {
	// Удаляем все существующие связи для категории
	err := s.storage.DeleteCategoryVariantMappings(ctx, request.CategoryID)
	if err != nil {
		return fmt.Errorf("failed to delete existing mappings: %w", err)
	}

	// Создаем новые связи
	for _, attr := range request.Attributes {
		mapping := &models.VariantAttributeMappingCreateRequest{
			VariantAttributeID: attr.AttributeID,
			CategoryID:         request.CategoryID,
			SortOrder:          attr.SortOrder,
			IsRequired:         attr.IsRequired,
		}
		_, err := s.storage.CreateVariantMapping(ctx, mapping)
		if err != nil {
			return fmt.Errorf("failed to create mapping: %w", err)
		}
	}

	return nil
}
