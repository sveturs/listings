package postgres

import (
	"backend/internal/domain/models"
	"context"
)

// CustomComponentStorage определяет интерфейс для работы с кастомными компонентами
type CustomComponentStorage interface {
	CreateComponent(ctx context.Context, component *models.CustomUIComponent) (int, error)
	GetComponent(ctx context.Context, id int) (*models.CustomUIComponent, error)
	UpdateComponent(ctx context.Context, id int, updates map[string]interface{}) error
	DeleteComponent(ctx context.Context, id int) error
	ListComponents(ctx context.Context, filters map[string]interface{}) ([]*models.CustomUIComponent, error)
	
	AddComponentUsage(ctx context.Context, usage *models.CustomUIComponentUsage) (int, error)
	GetComponentUsages(ctx context.Context, componentID, categoryID *int) ([]*models.CustomUIComponentUsage, error)
	RemoveComponentUsage(ctx context.Context, id int) error
	GetCategoryComponents(ctx context.Context, categoryID int, usageContext string) ([]*models.CustomUIComponentUsage, error)
	
	CreateTemplate(ctx context.Context, template *models.ComponentTemplate) (int, error)
	ListTemplates(ctx context.Context, componentID int) ([]*models.ComponentTemplate, error)
}

// NewCustomComponentStorage создает новое хранилище для кастомных компонентов
func NewCustomComponentStorage(db *Database) CustomComponentStorage {
	return db
}