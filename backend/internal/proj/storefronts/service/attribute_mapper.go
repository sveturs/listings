// backend/internal/proj/storefronts/service/attribute_mapper.go
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// AttributeMapper сервис для маппинга внешних атрибутов на внутренние unified_attributes
type AttributeMapper struct {
	storage Storage
	logger  zerolog.Logger
	// Кэш атрибутов для быстрого доступа
	attributesCache map[string]*AttributeTemplate // key = normalized code
	// Кэш маппинга external -> internal
	mappingCache map[string]int // key = normalized external name -> attribute_id
}

// AttributeTemplate информация об атрибуте из unified_attributes
type AttributeTemplate struct {
	ID              int                    `json:"id"`
	Code            string                 `json:"code"`
	Name            string                 `json:"name"`
	DisplayName     string                 `json:"display_name"`
	AttributeType   string                 `json:"attribute_type"`   // text, number, boolean, select, etc.
	Purpose         string                 `json:"purpose"`          // regular, variant, both
	Options         map[string]interface{} `json:"options"`          // для select/multiselect
	ValidationRules map[string]interface{} `json:"validation_rules"` // min, max, pattern, etc.
	UISettings      map[string]interface{} `json:"ui_settings"`
	IsSearchable    bool                   `json:"is_searchable"`
	IsFilterable    bool                   `json:"is_filterable"`
	IsRequired      bool                   `json:"is_required"`
	AffectsStock    bool                   `json:"affects_stock"`
	AffectsPrice    bool                   `json:"affects_price"`
}

// MappedAttribute результат маппинга
type MappedAttribute struct {
	AttributeID   int         `json:"attribute_id"`
	Code          string      `json:"code"`
	Name          string      `json:"name"`
	Value         interface{} `json:"value"`
	Confidence    float64     `json:"confidence"`     // 0.0-1.0
	IsNewAttribute bool       `json:"is_new_attribute"` // true если атрибут не найден и нужно создать
	SuggestedCode string      `json:"suggested_code"`   // предлагаемый code для нового атрибута
}

// AttributeMappingRequest запрос на маппинг атрибута
type AttributeMappingRequest struct {
	ExternalName  string      `json:"external_name"`
	ExternalValue interface{} `json:"external_value"`
	CategoryID    *int        `json:"category_id,omitempty"`
}

// NewAttributeMapper создает новый AttributeMapper
func NewAttributeMapper(storage Storage, logger zerolog.Logger) *AttributeMapper {
	return &AttributeMapper{
		storage:         storage,
		logger:          logger.With().Str("service", "AttributeMapper").Logger(),
		attributesCache: make(map[string]*AttributeTemplate),
		mappingCache:    make(map[string]int),
	}
}

// LoadAttributesCache загружает все атрибуты в кэш
func (m *AttributeMapper) LoadAttributesCache(ctx context.Context) error {
	m.logger.Info().Msg("Loading attributes into cache")

	// Получаем все активные атрибуты из БД
	attributes, err := m.storage.GetAllUnifiedAttributes(ctx)
	if err != nil {
		m.logger.Error().Err(err).Msg("Failed to load attributes from database")
		return fmt.Errorf("failed to load attributes: %w", err)
	}

	// Очищаем кэш и заполняем новыми данными
	m.attributesCache = make(map[string]*AttributeTemplate)
	m.mappingCache = make(map[string]int)

	for _, attr := range attributes {
		// Конвертируем models.UnifiedAttribute -> AttributeTemplate
		template := &AttributeTemplate{
			ID:              attr.ID,
			Code:            attr.Code,
			Name:            attr.Name,
			DisplayName:     attr.DisplayName,
			AttributeType:   attr.AttributeType,
			Purpose:         string(attr.Purpose),
			Options:         parseJSONMap(attr.Options),
			ValidationRules: parseJSONMap(attr.ValidationRules),
			UISettings:      parseJSONMap(attr.UISettings),
			IsSearchable:    attr.IsSearchable,
			IsFilterable:    attr.IsFilterable,
			IsRequired:      attr.IsRequired,
			AffectsStock:    attr.AffectsStock,
			AffectsPrice:    attr.AffectsPrice,
		}

		// Сохраняем в кэш по normalized code
		normalizedCode := m.normalizeAttributeName(attr.Code)
		m.attributesCache[normalizedCode] = template

		// Также сохраняем по normalized name для быстрого поиска
		normalizedName := m.normalizeAttributeName(attr.Name)
		m.mappingCache[normalizedName] = attr.ID
	}

	m.logger.Info().
		Int("count", len(attributes)).
		Msg("Attributes loaded into cache successfully")

	return nil
}

// parseJSONMap парсит json.RawMessage в map[string]interface{}
func parseJSONMap(data []byte) map[string]interface{} {
	if len(data) == 0 {
		return nil
	}
	var result map[string]interface{}
	// Игнорируем ошибки парсинга - вернем nil
	_ = json.Unmarshal(data, &result)
	return result
}

// MapExternalAttribute мапит внешний атрибут на внутренний
func (m *AttributeMapper) MapExternalAttribute(
	ctx context.Context,
	externalName string,
	externalValue interface{},
	categoryID *int,
) (*MappedAttribute, error) {
	m.logger.Debug().
		Str("external_name", externalName).
		Interface("external_value", externalValue).
		Msg("Mapping external attribute")

	// 1. Нормализуем имя атрибута
	normalizedName := m.normalizeAttributeName(externalName)

	// 2. Проверяем кэш маппинга
	if attributeID, found := m.mappingCache[normalizedName]; found {
		template := m.getAttributeFromCache(attributeID)
		if template != nil {
			value := m.transformValue(externalValue, template.AttributeType)
			return &MappedAttribute{
				AttributeID: template.ID,
				Code:        template.Code,
				Name:        template.Name,
				Value:       value,
				Confidence:  1.0,
			}, nil
		}
	}

	// 3. Ищем подходящий атрибут
	template := m.findMatchingTemplate(normalizedName, categoryID)

	// 4. Если не найден - предлагаем создать новый
	if template == nil {
		suggestedCode := m.generateAttributeCode(externalName)
		return &MappedAttribute{
			Code:           suggestedCode,
			Name:           externalName,
			Value:          externalValue,
			Confidence:     0.0,
			IsNewAttribute: true,
			SuggestedCode:  suggestedCode,
		}, nil
	}

	// 5. Трансформируем значение
	value := m.transformValue(externalValue, template.AttributeType)

	// 6. Валидируем
	if err := m.validateValue(value, template); err != nil {
		m.logger.Warn().
			Err(err).
			Str("attribute_code", template.Code).
			Interface("value", value).
			Msg("Attribute value validation failed")
		// Не возвращаем ошибку, просто логируем
	}

	// 7. Вычисляем confidence
	confidence := m.calculateConfidence(normalizedName, template, externalValue)

	// 8. Сохраняем в кэш маппинга
	m.mappingCache[normalizedName] = template.ID

	return &MappedAttribute{
		AttributeID: template.ID,
		Code:        template.Code,
		Name:        template.Name,
		Value:       value,
		Confidence:  confidence,
	}, nil
}

// BatchMapAttributes мапит список атрибутов
func (m *AttributeMapper) BatchMapAttributes(
	ctx context.Context,
	attributes []AttributeMappingRequest,
) ([]*MappedAttribute, error) {
	m.logger.Info().Int("count", len(attributes)).Msg("Batch mapping attributes")

	results := make([]*MappedAttribute, 0, len(attributes))

	for _, attr := range attributes {
		mapped, err := m.MapExternalAttribute(ctx, attr.ExternalName, attr.ExternalValue, attr.CategoryID)
		if err != nil {
			m.logger.Warn().
				Err(err).
				Str("external_name", attr.ExternalName).
				Msg("Failed to map attribute")
			continue
		}
		results = append(results, mapped)
	}

	m.logger.Info().
		Int("total", len(attributes)).
		Int("mapped", len(results)).
		Msg("Batch mapping completed")

	return results, nil
}

// normalizeAttributeName нормализует имя атрибута для поиска
func (m *AttributeMapper) normalizeAttributeName(name string) string {
	// Приводим к нижнему регистру
	normalized := strings.ToLower(name)

	// Убираем лишние пробелы
	normalized = strings.TrimSpace(normalized)

	// Заменяем подчеркивания на пробелы
	normalized = strings.ReplaceAll(normalized, "_", " ")

	// Убираем множественные пробелы
	normalized = regexp.MustCompile(`\s+`).ReplaceAllString(normalized, " ")

	return normalized
}

// findMatchingTemplate ищет подходящий атрибут в кэше
func (m *AttributeMapper) findMatchingTemplate(normalizedName string, categoryID *int) *AttributeTemplate {
	// Прямое совпадение по code
	if template, found := m.attributesCache[normalizedName]; found {
		return template
	}

	// Поиск по частичному совпадению
	for code, template := range m.attributesCache {
		if strings.Contains(normalizedName, code) || strings.Contains(code, normalizedName) {
			return template
		}
	}

	// TODO: Добавить поиск по категории (category-specific attributes)
	// TODO: Добавить fuzzy search для лучшего матчинга

	return nil
}

// getAttributeFromCache получает атрибут из кэша по ID
func (m *AttributeMapper) getAttributeFromCache(attributeID int) *AttributeTemplate {
	for _, template := range m.attributesCache {
		if template.ID == attributeID {
			return template
		}
	}
	return nil
}

// transformValue трансформирует значение в нужный тип
func (m *AttributeMapper) transformValue(value interface{}, attributeType string) interface{} {
	if value == nil {
		return nil
	}

	switch attributeType {
	case "number":
		return m.toNumber(value)
	case "boolean":
		return m.toBoolean(value)
	case "text", "textarea":
		return fmt.Sprintf("%v", value)
	case "select", "multiselect":
		return fmt.Sprintf("%v", value)
	case "date":
		return m.toDate(value)
	default:
		return fmt.Sprintf("%v", value)
	}
}

// toNumber конвертирует значение в число
func (m *AttributeMapper) toNumber(value interface{}) interface{} {
	switch v := value.(type) {
	case int, int8, int16, int32, int64:
		return v
	case uint, uint8, uint16, uint32, uint64:
		return v
	case float32, float64:
		return v
	case string:
		// Пытаемся распарсить строку
		if num, err := strconv.ParseFloat(v, 64); err == nil {
			return num
		}
	}
	return value
}

// toBoolean конвертирует значение в boolean
func (m *AttributeMapper) toBoolean(value interface{}) interface{} {
	switch v := value.(type) {
	case bool:
		return v
	case string:
		lower := strings.ToLower(strings.TrimSpace(v))
		return lower == "true" || lower == "yes" || lower == "да" || lower == "1"
	case int, int8, int16, int32, int64:
		return v != 0
	case uint, uint8, uint16, uint32, uint64:
		return v != 0
	}
	return false
}

// toDate конвертирует значение в дату
func (m *AttributeMapper) toDate(value interface{}) interface{} {
	switch v := value.(type) {
	case time.Time:
		return v.Format("2006-01-02")
	case string:
		// Пытаемся распарсить дату
		if t, err := time.Parse("2006-01-02", v); err == nil {
			return t.Format("2006-01-02")
		}
		if t, err := time.Parse("02.01.2006", v); err == nil {
			return t.Format("2006-01-02")
		}
	}
	return fmt.Sprintf("%v", value)
}

// validateValue валидирует значение атрибута
func (m *AttributeMapper) validateValue(value interface{}, template *AttributeTemplate) error {
	if value == nil {
		if template.IsRequired {
			return fmt.Errorf("attribute %s is required but value is nil", template.Code)
		}
		return nil
	}

	// Проверяем validation_rules
	if len(template.ValidationRules) > 0 {
		// TODO: реализовать валидацию на основе rules (min, max, pattern, etc.)
		m.logger.Debug().
			Str("attribute_code", template.Code).
			Interface("rules", template.ValidationRules).
			Msg("Validation rules defined but not implemented yet")
	}

	return nil
}

// calculateConfidence вычисляет уверенность в маппинге
func (m *AttributeMapper) calculateConfidence(normalizedName string, template *AttributeTemplate, value interface{}) float64 {
	confidence := 0.5 // базовая уверенность

	// Прямое совпадение кода
	if normalizedName == strings.ToLower(template.Code) {
		confidence = 1.0
		return confidence
	}

	// Совпадение имени
	if normalizedName == strings.ToLower(template.Name) {
		confidence = 0.95
		return confidence
	}

	// Частичное совпадение
	if strings.Contains(normalizedName, strings.ToLower(template.Code)) ||
		strings.Contains(strings.ToLower(template.Code), normalizedName) {
		confidence = 0.8
		return confidence
	}

	// Проверка валидности значения
	if err := m.validateValue(value, template); err == nil {
		confidence += 0.1
	}

	// Ограничиваем confidence от 0.0 до 1.0
	if confidence > 1.0 {
		confidence = 1.0
	}
	if confidence < 0.0 {
		confidence = 0.0
	}

	return confidence
}

// generateAttributeCode генерирует code для нового атрибута
func (m *AttributeMapper) generateAttributeCode(externalName string) string {
	// Приводим к нижнему регистру
	code := strings.ToLower(externalName)

	// Заменяем пробелы и спецсимволы на подчеркивания
	code = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(code, "_")

	// Убираем начальные и конечные подчеркивания
	code = strings.Trim(code, "_")

	// Убираем множественные подчеркивания
	code = regexp.MustCompile(`_{2,}`).ReplaceAllString(code, "_")

	// Ограничиваем длину
	if len(code) > 100 {
		code = code[:100]
	}

	return code
}
