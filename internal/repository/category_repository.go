package repository

import (
	"context"

	"github.com/vondi-global/listings/internal/domain"
)

// CategoryRepository defines the interface for category data access operations
type CategoryRepository interface {
	// Read Operations
	GetByID(ctx context.Context, id int32) (*domain.CategoryDetail, error)
	GetBySlug(ctx context.Context, slug string) (*domain.CategoryDetail, error)
	List(ctx context.Context, filter *domain.GetCategoriesFilter) ([]*domain.CategoryDetail, int64, error)
	GetTree(ctx context.Context, filter *domain.GetCategoryTreeFilter) ([]*domain.CategoryTree, error)

	// Admin Operations
	Create(ctx context.Context, input *domain.CreateCategoryInput) (*domain.CategoryDetail, error)
	Update(ctx context.Context, id int32, input *domain.UpdateCategoryInput) (*domain.CategoryDetail, error)
	Delete(ctx context.Context, id int32) error // Soft delete (set is_active=false)
}
