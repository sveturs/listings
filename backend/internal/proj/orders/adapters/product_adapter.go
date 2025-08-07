package adapters

import (
	"context"

	"backend/internal/domain/models"
)

// DatabaseInterface определяет необходимые методы базы данных
type DatabaseInterface interface {
	GetProductByID(ctx context.Context, productID int64) (*models.StorefrontProduct, error)
	GetProductVariantByID(ctx context.Context, variantID int64) (*models.StorefrontProductVariant, error)
	UpdateProductStock(ctx context.Context, productID int64, variantID *int64, quantity int) error
}

// ProductRepositoryAdapter адаптирует Database к ProductRepositoryInterface
type ProductRepositoryAdapter struct {
	db DatabaseInterface
}

// NewProductRepositoryAdapter создает новый адаптер
func NewProductRepositoryAdapter(db DatabaseInterface) *ProductRepositoryAdapter {
	return &ProductRepositoryAdapter{db: db}
}

// GetByID получает продукт по ID
func (a *ProductRepositoryAdapter) GetByID(ctx context.Context, id int64) (*models.StorefrontProduct, error) {
	return a.db.GetProductByID(ctx, id)
}

// GetVariantByID получает вариант продукта по ID
func (a *ProductRepositoryAdapter) GetVariantByID(ctx context.Context, id int64) (*models.StorefrontProductVariant, error) {
	return a.db.GetProductVariantByID(ctx, id)
}

// UpdateStock обновляет остатки на складе
func (a *ProductRepositoryAdapter) UpdateStock(ctx context.Context, productID int64, variantID *int64, quantity int) error {
	return a.db.UpdateProductStock(ctx, productID, variantID, quantity)
}
