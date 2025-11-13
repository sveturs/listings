package service

import (
	"context"

	"github.com/sveturs/listings/internal/domain"
)

// CategoryService defines the interface for category business logic operations
type CategoryService interface {
	// Public read operations
	GetCategories(ctx context.Context, parentID *int64, isActive *bool, limit, offset int32) ([]*domain.Category, int32, error)
	GetCategory(ctx context.Context, id int64) (*domain.Category, error)
	GetCategoryBySlug(ctx context.Context, slug string) (*domain.Category, error)
	GetCategoryTree(ctx context.Context, categoryID int64) (*domain.CategoryTreeNode, error)

	// Admin write operations
	CreateCategory(ctx context.Context, cat *domain.Category) (*domain.Category, error)
	UpdateCategory(ctx context.Context, cat *domain.Category) (*domain.Category, error)
	DeleteCategory(ctx context.Context, categoryID int64) error

	// Cache management
	InvalidateCache(ctx context.Context, categoryID int64) error
}
