package models

import (
	"encoding/json"
	"time"
)

// CategoryAttribute представляет атрибут категории
type CategoryAttribute struct {
    ID                 int             `json:"id"`
    Name               string          `json:"name"`
    DisplayName        string          `json:"display_name"`
    AttributeType      string          `json:"attribute_type"`
    Options            json.RawMessage `json:"options,omitempty"`
    ValidRules         json.RawMessage `json:"validation_rules,omitempty"`
    IsSearchable       bool            `json:"is_searchable"`
    IsFilterable       bool            `json:"is_filterable"`
    IsRequired         bool            `json:"is_required"`
    SortOrder          int             `json:"sort_order"`
    CreatedAt          time.Time       `json:"created_at"`
    Translations       map[string]string `json:"translations,omitempty"`
    OptionTranslations map[string]map[string]string `json:"option_translations,omitempty"`
}

// AttributeOptions содержит опции для атрибутов типа select
type AttributeOptions struct {
	Values     []string `json:"values,omitempty"`
	Min        *float64 `json:"min,omitempty"`
	Max        *float64 `json:"max,omitempty"`
	Step       *float64 `json:"step,omitempty"`
	Multiselect bool     `json:"multiselect,omitempty"`
}

// ListingAttributeValue содержит значение атрибута для объявления
type ListingAttributeValue struct {
    ListingID     int             `json:"listing_id"`
    AttributeID   int             `json:"attribute_id"`
    AttributeName string          `json:"attribute_name"`
    DisplayName   string          `json:"display_name"`
    AttributeType string          `json:"attribute_type"`
    TextValue     *string         `json:"text_value,omitempty"`
    NumericValue  *float64        `json:"numeric_value,omitempty"`
    BooleanValue  *bool           `json:"boolean_value,omitempty"`
    JSONValue     json.RawMessage `json:"json_value,omitempty"`
    // Вспомогательное поле для вывода значения любого типа в строковом виде
    DisplayValue  string          `json:"display_value"`
    // Добавляем поле для единиц измерения
    Unit          string          `json:"unit,omitempty"`
    Translations      map[string]string               `json:"translations,omitempty"`
    OptionTranslations map[string]map[string]string   `json:"option_translations,omitempty"`
}

// CategoryAttributeMapping связывает атрибуты с категориями
type CategoryAttributeMapping struct {
	CategoryID  int  `json:"category_id"`
	AttributeID int  `json:"attribute_id"`
	IsEnabled   bool `json:"is_enabled"`
	IsRequired  bool `json:"is_required"`
	Attribute   *CategoryAttribute `json:"attribute,omitempty"`
}