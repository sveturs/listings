package models

import (
	"encoding/json"
	"fmt"
	"time"
)

// AttributePurpose определяет назначение атрибута
type AttributePurpose string

const (
	PurposeRegular AttributePurpose = "regular" // Обычный атрибут для фильтрации/поиска
	PurposeVariant AttributePurpose = "variant" // Вариативный атрибут (влияет на SKU)
	PurposeBoth    AttributePurpose = "both"    // Может использоваться в обоих случаях
)

// AttributeEntityType определяет тип сущности для значений атрибутов
type AttributeEntityType string

const (
	AttributeEntityTypeListing        AttributeEntityType = "listing"         // Объявление маркетплейса
	AttributeEntityTypeProduct        AttributeEntityType = "product"         // Товар витрины
	AttributeEntityTypeProductVariant AttributeEntityType = "product_variant" // Вариант товара
)

// UnifiedAttribute представляет унифицированный атрибут
type UnifiedAttribute struct {
	ID                  int              `json:"id" db:"id"`
	Code                string           `json:"code" db:"code"`
	Name                string           `json:"name" db:"name"`
	DisplayName         string           `json:"display_name" db:"display_name"`
	AttributeType       string           `json:"attribute_type" db:"attribute_type"`
	Purpose             AttributePurpose `json:"purpose" db:"purpose"`
	Options             json.RawMessage  `json:"options,omitempty" db:"options"`
	ValidationRules     json.RawMessage  `json:"validation_rules,omitempty" db:"validation_rules"`
	UISettings          json.RawMessage  `json:"ui_settings,omitempty" db:"ui_settings"`
	IsSearchable        bool             `json:"is_searchable" db:"is_searchable"`
	IsFilterable        bool             `json:"is_filterable" db:"is_filterable"`
	IsRequired          bool             `json:"is_required" db:"is_required"`
	IsVariantCompatible bool             `json:"is_variant_compatible" db:"is_variant_compatible"`
	AffectsStock        bool             `json:"affects_stock" db:"affects_stock"`
	AffectsPrice        bool             `json:"affects_price" db:"affects_price"`
	SortOrder           int              `json:"sort_order" db:"sort_order"`
	IsActive            bool             `json:"is_active" db:"is_active"`
	CreatedAt           time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time        `json:"updated_at" db:"updated_at"`

	// Поля для обратной совместимости (временно)
	LegacyCategoryAttributeID       *int `json:"-" db:"legacy_category_attribute_id"`
	LegacyProductVariantAttributeID *int `json:"-" db:"legacy_product_variant_attribute_id"`

	// Дополнительные поля для удобства работы
	Translations       map[string]string            `json:"translations,omitempty"`
	OptionTranslations map[string]map[string]string `json:"option_translations,omitempty"`
}

// UnifiedCategoryAttribute представляет связь атрибута с категорией
type UnifiedCategoryAttribute struct {
	ID                      int             `json:"id" db:"id"`
	CategoryID              int             `json:"category_id" db:"category_id"`
	AttributeID             int             `json:"attribute_id" db:"attribute_id"`
	IsEnabled               bool            `json:"is_enabled" db:"is_enabled"`
	IsRequired              bool            `json:"is_required" db:"is_required"`
	IsFilter                bool            `json:"is_filter" db:"is_filter"`
	SortOrder               int             `json:"sort_order" db:"sort_order"`
	CategorySpecificOptions json.RawMessage `json:"category_specific_options,omitempty" db:"category_specific_options"`
	CreatedAt               time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time       `json:"updated_at" db:"updated_at"`

	// Связанные объекты
	Attribute *UnifiedAttribute    `json:"attribute,omitempty"`
	Category  *MarketplaceCategory `json:"category,omitempty"`
}

// UnifiedAttributeValue представляет значение атрибута
type UnifiedAttributeValue struct {
	ID           int                 `json:"id" db:"id"`
	EntityType   AttributeEntityType `json:"entity_type" db:"entity_type"`
	EntityID     int                 `json:"entity_id" db:"entity_id"`
	AttributeID  int                 `json:"attribute_id" db:"attribute_id"`
	TextValue    *string             `json:"text_value,omitempty" db:"text_value"`
	NumericValue *float64            `json:"numeric_value,omitempty" db:"numeric_value"`
	BooleanValue *bool               `json:"boolean_value,omitempty" db:"boolean_value"`
	DateValue    *time.Time          `json:"date_value,omitempty" db:"date_value"`
	JSONValue    json.RawMessage     `json:"json_value,omitempty" db:"json_value"`
	CreatedAt    time.Time           `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at" db:"updated_at"`

	// Связанные объекты для удобства
	Attribute *UnifiedAttribute `json:"attribute,omitempty"`

	// Вспомогательные поля
	DisplayValue string `json:"display_value,omitempty"`
	Unit         string `json:"unit,omitempty"`
}

// ToCategoryAttribute преобразует UnifiedAttribute в CategoryAttribute для обратной совместимости
func (ua *UnifiedAttribute) ToCategoryAttribute() *CategoryAttribute {
	return &CategoryAttribute{
		ID:                  ua.ID,
		Name:                ua.Name,
		DisplayName:         ua.DisplayName,
		AttributeType:       ua.AttributeType,
		Options:             ua.Options,
		ValidRules:          ua.ValidationRules,
		IsSearchable:        ua.IsSearchable,
		IsFilterable:        ua.IsFilterable,
		IsRequired:          ua.IsRequired,
		SortOrder:           ua.SortOrder,
		CreatedAt:           ua.CreatedAt,
		Translations:        ua.Translations,
		OptionTranslations:  ua.OptionTranslations,
		IsVariantCompatible: ua.IsVariantCompatible || ua.Purpose == PurposeVariant || ua.Purpose == PurposeBoth,
		AffectsStock:        ua.AffectsStock,
	}
}

// ToProductVariantAttribute преобразует UnifiedAttribute в ProductVariantAttribute для обратной совместимости
func (ua *UnifiedAttribute) ToProductVariantAttribute() *ProductVariantAttribute {
	return &ProductVariantAttribute{
		ID:           ua.ID,
		Name:         ua.Name,
		DisplayName:  ua.DisplayName,
		Type:         ua.AttributeType,
		IsRequired:   ua.IsRequired,
		SortOrder:    ua.SortOrder,
		AffectsStock: ua.AffectsStock,
		CreatedAt:    ua.CreatedAt,
		UpdatedAt:    ua.UpdatedAt,
	}
}

// ToListingAttributeValue преобразует UnifiedAttributeValue в ListingAttributeValue для обратной совместимости
func (uav *UnifiedAttributeValue) ToListingAttributeValue() *ListingAttributeValue {
	lav := &ListingAttributeValue{
		ListingID:    uav.EntityID,
		AttributeID:  uav.AttributeID,
		TextValue:    uav.TextValue,
		NumericValue: uav.NumericValue,
		BooleanValue: uav.BooleanValue,
		JSONValue:    uav.JSONValue,
		DisplayValue: uav.DisplayValue,
		Unit:         uav.Unit,
	}

	// Заполняем информацию об атрибуте, если она есть
	if uav.Attribute != nil {
		lav.AttributeName = uav.Attribute.Name
		lav.DisplayName = uav.Attribute.DisplayName
		lav.AttributeType = uav.Attribute.AttributeType
		lav.IsRequired = uav.Attribute.IsRequired
		lav.Translations = uav.Attribute.Translations
		lav.OptionTranslations = uav.Attribute.OptionTranslations
	}

	return lav
}

// GetDisplayValue возвращает отображаемое значение атрибута
func (uav *UnifiedAttributeValue) GetDisplayValue() string {
	if uav.DisplayValue != "" {
		return uav.DisplayValue
	}

	switch {
	case uav.TextValue != nil:
		return *uav.TextValue
	case uav.NumericValue != nil:
		if uav.Unit != "" {
			return fmt.Sprintf("%.2f %s", *uav.NumericValue, uav.Unit)
		}
		return fmt.Sprintf("%.2f", *uav.NumericValue)
	case uav.BooleanValue != nil:
		if *uav.BooleanValue {
			return "Yes"
		}
		return "No"
	case uav.DateValue != nil:
		return uav.DateValue.Format("2006-01-02")
	case uav.JSONValue != nil:
		return string(uav.JSONValue)
	default:
		return ""
	}
}

// IsValidForPurpose проверяет, подходит ли атрибут для заданного назначения
func (ua *UnifiedAttribute) IsValidForPurpose(purpose AttributePurpose) bool {
	return ua.Purpose == PurposeBoth || ua.Purpose == purpose
}

// UnifiedAttributeFilter представляет фильтр для поиска атрибутов
type UnifiedAttributeFilter struct {
	CategoryID   *int              `json:"category_id,omitempty"`
	Purpose      *AttributePurpose `json:"purpose,omitempty"`
	IsActive     *bool             `json:"is_active,omitempty"`
	IsSearchable *bool             `json:"is_searchable,omitempty"`
	IsFilterable *bool             `json:"is_filterable,omitempty"`
	Limit        int               `json:"limit,omitempty"`
	Offset       int               `json:"offset,omitempty"`
}
