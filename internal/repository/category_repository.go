package repository

import (
	"context"

	"github.com/vondi-global/listings/internal/domain"
)

// CategoryRepository defines the interface for category data access operations
type CategoryRepository interface {
	// Read Operations (V1 - int32 based, deprecated)
	GetByID(ctx context.Context, id int32) (*domain.CategoryDetail, error)
	GetBySlug(ctx context.Context, slug string) (*domain.CategoryDetail, error)
	List(ctx context.Context, filter *domain.GetCategoriesFilter) ([]*domain.CategoryDetail, int64, error)
	GetTree(ctx context.Context, filter *domain.GetCategoryTreeFilter) ([]*domain.CategoryTree, error)

	// Admin Operations (V1 - int32 based, deprecated)
	Create(ctx context.Context, input *domain.CreateCategoryInput) (*domain.CategoryDetail, error)
	Update(ctx context.Context, id int32, input *domain.UpdateCategoryInput) (*domain.CategoryDetail, error)
	Delete(ctx context.Context, id int32) error // Soft delete (set is_active=false)
}

// CategoryRepositoryV2 defines V2 interface with UUID support and i18n
type CategoryRepositoryV2 interface {
	// Read Operations
	GetByUUID(ctx context.Context, id string) (*domain.CategoryV2, error)
	GetBySlugV2(ctx context.Context, slug string) (*domain.CategoryV2, error)
	GetTreeV2(ctx context.Context, filter *domain.GetCategoryTreeFilterV2) ([]*domain.CategoryTreeV2, error)
	GetBreadcrumb(ctx context.Context, categoryID string, locale string) ([]*domain.CategoryBreadcrumb, error)
	ListV2(ctx context.Context, parentID *string, activeOnly bool, page, pageSize int32) ([]*domain.CategoryV2, int64, error)
}
