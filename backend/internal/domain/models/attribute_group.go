package models

import (
	"encoding/json"
	"time"
)

// AttributeGroup представляет группу атрибутов
type AttributeGroup struct {
	ID          int       `json:"id" db:"id"`
	Code        string    `json:"code" db:"code"`
	Name        string    `json:"name" db:"name"`
	DisplayName string    `json:"display_name" db:"display_name"`
	Description *string   `json:"description" db:"description"`
	Icon        *string   `json:"icon" db:"icon"`
	SortOrder   int       `json:"sort_order" db:"sort_order"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	IsSystem    bool      `json:"is_system" db:"is_system"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// AttributeGroupItem представляет атрибут в группе
type AttributeGroupItem struct {
	ID                  int                `json:"id" db:"id"`
	GroupID             int                `json:"group_id" db:"group_id"`
	AttributeID         int                `json:"attribute_id" db:"attribute_id"`
	Icon                *string            `json:"icon" db:"icon"`
	SortOrder           int                `json:"sort_order" db:"sort_order"`
	CustomDisplayName   *string            `json:"custom_display_name" db:"custom_display_name"`
	VisibilityCondition json.RawMessage    `json:"visibility_condition" db:"visibility_condition"`
	CreatedAt           time.Time          `json:"created_at" db:"created_at"`
	Attribute           *CategoryAttribute `json:"attribute,omitempty"`
}

// CategoryAttributeGroup представляет связь группы атрибутов с категорией
type CategoryAttributeGroup struct {
	ID                 int             `json:"id" db:"id"`
	CategoryID         int             `json:"category_id" db:"category_id"`
	GroupID            int             `json:"group_id" db:"group_id"`
	ComponentID        *int            `json:"component_id" db:"component_id"`
	SortOrder          int             `json:"sort_order" db:"sort_order"`
	IsActive           bool            `json:"is_active" db:"is_active"`
	DisplayMode        string          `json:"display_mode" db:"display_mode"`
	CollapsedByDefault bool            `json:"collapsed_by_default" db:"collapsed_by_default"`
	Configuration      json.RawMessage `json:"configuration" db:"configuration"`
	CreatedAt          time.Time       `json:"created_at" db:"created_at"`
	Group              *AttributeGroup `json:"group,omitempty"`
}

// AttributeGroupWithItems представляет группу с атрибутами
type AttributeGroupWithItems struct {
	*AttributeGroup
	Items []*AttributeGroupItem `json:"items"`
}

// CreateAttributeGroupRequest запрос на создание группы атрибутов
type CreateAttributeGroupRequest struct {
	Code        string  `json:"code" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	DisplayName string  `json:"display_name" validate:"required"`
	Description *string `json:"description"`
	Icon        *string `json:"icon"`
	SortOrder   int     `json:"sort_order"`
	IsActive    bool    `json:"is_active"`
}

// UpdateAttributeGroupRequest запрос на обновление группы атрибутов
type UpdateAttributeGroupRequest struct {
	Name        string  `json:"name"`
	DisplayName string  `json:"display_name"`
	Description *string `json:"description"`
	Icon        *string `json:"icon"`
	SortOrder   int     `json:"sort_order"`
	IsActive    bool    `json:"is_active"`
}

// AddItemToGroupRequest запрос на добавление атрибута в группу
type AddItemToGroupRequest struct {
	AttributeID         int             `json:"attribute_id" validate:"required"`
	Icon                *string         `json:"icon"`
	SortOrder           int             `json:"sort_order"`
	CustomDisplayName   *string         `json:"custom_display_name"`
	VisibilityCondition json.RawMessage `json:"visibility_condition"`
}

// AttachGroupToCategoryRequest запрос на привязку группы к категории
type AttachGroupToCategoryRequest struct {
	GroupID            int             `json:"group_id" validate:"required"`
	ComponentID        *int            `json:"component_id"`
	SortOrder          int             `json:"sort_order"`
	IsActive           bool            `json:"is_active"`
	DisplayMode        string          `json:"display_mode"`
	CollapsedByDefault bool            `json:"collapsed_by_default"`
	Configuration      json.RawMessage `json:"configuration"`
}
