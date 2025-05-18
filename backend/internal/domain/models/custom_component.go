package models

import (
	"encoding/json"
	"time"
)

// CustomUIComponent представляет кастомный UI компонент
type CustomUIComponent struct {
	ID               int                    `json:"id" db:"id"`
	Name             string                 `json:"name" db:"name"`
	ComponentType    string                 `json:"component_type" db:"component_type"`
	Description      string                 `json:"description" db:"description"`
	TemplateCode     string                 `json:"template_code" db:"template_code"`
	Styles           string                 `json:"styles" db:"styles"`
	PropsSchema      json.RawMessage        `json:"props_schema" db:"props_schema"`
	IsActive         bool                   `json:"is_active" db:"is_active"`
	CreatedAt        time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at" db:"updated_at"`
	CreatedBy        *int                   `json:"created_by" db:"created_by"`
	UpdatedBy        *int                   `json:"updated_by" db:"updated_by"`
}

// CustomUIComponentUsage представляет использование компонента
type CustomUIComponentUsage struct {
	ID              int                `json:"id" db:"id"`
	ComponentID     int                `json:"component_id" db:"component_id"`
	CategoryID      *int               `json:"category_id" db:"category_id"`
	UsageContext    string             `json:"usage_context" db:"usage_context"`
	Placement       string             `json:"placement" db:"placement"`
	Priority        int                `json:"priority" db:"priority"`
	Configuration   json.RawMessage    `json:"configuration" db:"configuration"`
	ConditionsLogic json.RawMessage    `json:"conditions_logic" db:"conditions_logic"`
	IsActive        bool               `json:"is_active" db:"is_active"`
	CreatedAt       time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at" db:"updated_at"`
	CreatedBy       *int               `json:"created_by" db:"created_by"`
	UpdatedBy       *int               `json:"updated_by" db:"updated_by"`
	Component       *CustomUIComponent `json:"component,omitempty"`
	ComponentName   string             `json:"component_name,omitempty"`
	CategoryName    string             `json:"category_name,omitempty"`
}

// CustomUITemplate представляет шаблон для создания компонентов
type CustomUITemplate struct {
	ID           int             `json:"id" db:"id"`
	Name         string          `json:"name" db:"name"`
	DisplayName  string          `json:"display_name" db:"display_name"`
	Description  string          `json:"description" db:"description"`
	TemplateCode string          `json:"template_code" db:"template_code"`
	Variables    json.RawMessage `json:"variables" db:"variables"`
	IsShared     bool            `json:"is_shared" db:"is_shared"`
	CreatedAt    time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at" db:"updated_at"`
	CreatedBy    *int            `json:"created_by" db:"created_by"`
	UpdatedBy    *int            `json:"updated_by" db:"updated_by"`
}

// ComponentType представляет тип компонента
type ComponentType string

const (
	ComponentTypeCategory  ComponentType = "category"
	ComponentTypeAttribute ComponentType = "attribute"
	ComponentTypeFilter    ComponentType = "filter"
)

// EntityType представляет тип сущности для использования компонента
type EntityType string

const (
	EntityTypeCategory  EntityType = "category"
	EntityTypeAttribute EntityType = "attribute"
)

// CreateCustomComponentRequest представляет запрос на создание компонента
type CreateCustomComponentRequest struct {
	Name          string          `json:"name" validate:"required,min=1,max=255"`
	ComponentType string          `json:"component_type" validate:"required,oneof=category attribute filter"`
	Description   string          `json:"description"`
	TemplateCode  string          `json:"template_code" validate:"required"`
	Styles        string          `json:"styles"`
	PropsSchema   json.RawMessage `json:"props_schema"`
	IsActive      bool            `json:"is_active"`
}

// UpdateCustomComponentRequest представляет запрос на обновление компонента
type UpdateCustomComponentRequest struct {
	Name          string          `json:"name"`
	ComponentType string          `json:"component_type"`
	Description   string          `json:"description"`
	TemplateCode  string          `json:"template_code"`
	Styles        string          `json:"styles"`
	PropsSchema   json.RawMessage `json:"props_schema"`
	IsActive      bool            `json:"is_active"`
}

// CreateComponentUsageRequest представляет запрос на добавление использования компонента
type CreateComponentUsageRequest struct {
	ComponentID     int             `json:"component_id" validate:"required"`
	CategoryID      *int            `json:"category_id"`
	UsageContext    string          `json:"usage_context" validate:"required"`
	Placement       string          `json:"placement"`
	Priority        int             `json:"priority"`
	Configuration   json.RawMessage `json:"configuration"`
	ConditionsLogic json.RawMessage `json:"conditions_logic"`
	IsActive        bool            `json:"is_active"`
}

// CreateTemplateRequest представляет запрос на создание шаблона
type CreateTemplateRequest struct {
	ComponentID  int             `json:"component_id" validate:"required"`
	Name         string          `json:"name" validate:"required"`
	Description  string          `json:"description"`
	Variables    json.RawMessage `json:"variables"`
}

// ComponentTemplate представляет шаблон компонента
type ComponentTemplate struct {
	ID             int             `json:"id" db:"id"`
	ComponentID    int             `json:"component_id" db:"component_id"`
	Name           string          `json:"name" db:"name"`
	Description    string          `json:"description" db:"description"`
	TemplateConfig json.RawMessage `json:"template_config" db:"template_config"`
	PreviewImage   string          `json:"preview_image" db:"preview_image"`
	CategoryID     *int            `json:"category_id" db:"category_id"`
	CreatedAt      time.Time       `json:"created_at" db:"created_at"`
	CreatedBy      *int            `json:"created_by" db:"created_by"`
}