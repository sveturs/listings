package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"

	"backend/internal/domain/models"
	"backend/internal/logger"
	variantTypes "backend/internal/proj/storefront/types"
	"backend/internal/proj/storefronts/storage/opensearch"
	"backend/pkg/utils"
)

// ProductSearchRepository is an alias for OpenSearch interface
type ProductSearchRepository = opensearch.ProductSearchRepository

// VariantService interface for variant operations
type VariantService interface {
	BulkCreateVariants(ctx context.Context, productID int, variants []variantTypes.CreateVariantRequest) ([]*variantTypes.ProductVariant, error)
	BulkCreateVariantsTx(ctx context.Context, tx interface{}, productID int, variants []variantTypes.CreateVariantRequest) ([]*variantTypes.ProductVariant, error)
	CreateVariant(ctx context.Context, req *variantTypes.CreateVariantRequest) (*variantTypes.ProductVariant, error)
	GetVariantsByProductID(ctx context.Context, productID int) ([]*variantTypes.ProductVariant, error)
}

// ProductService handles business logic for storefront products
type ProductService struct {
	storage        Storage
	searchRepo     opensearch.ProductSearchRepository
	variantService VariantService
}

// Storage interface for product operations
type Storage interface {
	GetStorefrontProducts(ctx context.Context, filter models.ProductFilter) ([]*models.StorefrontProduct, error)
	GetStorefrontProduct(ctx context.Context, storefrontID, productID int) (*models.StorefrontProduct, error)
	GetStorefrontProductByID(ctx context.Context, productID int) (*models.StorefrontProduct, error)
	CreateStorefrontProduct(ctx context.Context, storefrontID int, req *models.CreateProductRequest) (*models.StorefrontProduct, error)
	UpdateStorefrontProduct(ctx context.Context, storefrontID, productID int, req *models.UpdateProductRequest) error
	DeleteStorefrontProduct(ctx context.Context, storefrontID, productID int) error
	UpdateProductInventory(ctx context.Context, storefrontID, productID int, userID int, req *models.UpdateInventoryRequest) error
	GetProductStats(ctx context.Context, storefrontID int) (*models.ProductStats, error)

	// Bulk operations
	BulkCreateProducts(ctx context.Context, storefrontID int, products []models.CreateProductRequest) ([]int, []error)
	BulkUpdateProducts(ctx context.Context, storefrontID int, updates []models.BulkUpdateItem) ([]int, []error)
	BulkDeleteProducts(ctx context.Context, storefrontID int, productIDs []int) ([]int, []error)
	BulkUpdateStatus(ctx context.Context, storefrontID int, productIDs []int, isActive bool) ([]int, []error)

	// Storefront operations
	GetStorefrontByID(ctx context.Context, id int) (*models.Storefront, error)

	// Transaction support
	BeginTx(ctx context.Context) (Transaction, error)

	// Transactional methods
	CreateStorefrontProductTx(ctx context.Context, tx Transaction, storefrontID int, req *models.CreateProductRequest) (*models.StorefrontProduct, error)
}

// Transaction interface for database transactions
type Transaction interface {
	Rollback() error
	Commit() error
	// GetPgxTx returns the underlying pgx.Tx for use with repository methods
	GetPgxTx() interface{} // Returns pgx.Tx but using interface{} to avoid import cycle
}

// NewProductService creates a new product service
func NewProductService(storage Storage, searchRepo opensearch.ProductSearchRepository, variantService VariantService) *ProductService {
	return &ProductService{
		storage:        storage,
		searchRepo:     searchRepo,
		variantService: variantService,
	}
}

// ValidateStorefrontOwnership checks if user owns the storefront
func (s *ProductService) ValidateStorefrontOwnership(ctx context.Context, storefrontID, userID int) error {
	// Проверяем, является ли пользователь администратором
	if isAdmin, ok := ctx.Value("is_admin").(bool); ok && isAdmin {
		// Администратор имеет доступ ко всем витринам
		logger.Info().Msgf("Admin user %d accessing storefront %d", userID, storefrontID)
		return nil
	}

	storefront, err := s.storage.GetStorefrontByID(ctx, storefrontID)
	if err != nil {
		return fmt.Errorf("failed to get storefront: %w", err)
	}

	if storefront == nil {
		return fmt.Errorf("storefront not found")
	}

	if storefront.UserID != userID {
		return fmt.Errorf("unauthorized: user %d does not own storefront %d (owner is %d)", userID, storefrontID, storefront.UserID)
	}

	return nil
}

// GetProducts retrieves products for a storefront
func (s *ProductService) GetProducts(ctx context.Context, filter models.ProductFilter) ([]*models.StorefrontProduct, error) {
	logger.Info().Msgf("GetProducts called with filter: %+v", filter)

	// Validate storefront exists
	storefront, err := s.storage.GetStorefrontByID(ctx, filter.StorefrontID)
	if err != nil {
		logger.Error().Err(err).Msgf("Failed to get storefront %d", filter.StorefrontID)
		return nil, fmt.Errorf("failed to get storefront: %w", err)
	}

	if storefront == nil {
		logger.Error().Msgf("Storefront %d not found", filter.StorefrontID)
		return nil, fmt.Errorf("storefront not found")
	}

	// Get products
	products, err := s.storage.GetStorefrontProducts(ctx, filter)
	if err != nil {
		logger.Error().Err(err).Msgf("Failed to get products for storefront %d", filter.StorefrontID)
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	// Обрабатываем адреса всех товаров с учетом приватности
	for i := range products {
		s.processProductLocationPrivacy(products[i])
	}

	logger.Info().Msgf("Found %d products for storefront %d", len(products), filter.StorefrontID)
	return products, nil
}

// GetProduct retrieves a single product
func (s *ProductService) GetProduct(ctx context.Context, storefrontID, productID int) (*models.StorefrontProduct, error) {
	product, err := s.storage.GetStorefrontProduct(ctx, storefrontID, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	if product == nil {
		return nil, fmt.Errorf("product not found")
	}

	// Обрабатываем адрес с учетом приватности
	s.processProductLocationPrivacy(product)

	return product, nil
}

// GetProductByID retrieves a single product by its ID without requiring storefront ID
func (s *ProductService) GetProductByID(ctx context.Context, productID int) (*models.StorefrontProduct, error) {
	product, err := s.storage.GetStorefrontProductByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	if product == nil {
		return nil, fmt.Errorf("product not found")
	}

	// Обрабатываем адрес с учетом приватности
	s.processProductLocationPrivacy(product)

	return product, nil
}

// processProductLocationPrivacy обрабатывает адрес товара с учетом уровня приватности
func (s *ProductService) processProductLocationPrivacy(product *models.StorefrontProduct) {
	if product == nil {
		return
	}

	// Если у товара есть индивидуальный адрес и установлен уровень приватности
	if product.HasIndividualLocation && product.IndividualAddress != nil && product.LocationPrivacy != nil {
		// Форматируем адрес в соответствии с уровнем приватности
		formattedAddress := utils.FormatAddressWithPrivacy(*product.IndividualAddress, *product.LocationPrivacy)
		product.IndividualAddress = &formattedAddress

		// Обрабатываем координаты с учетом приватности
		if product.IndividualLatitude != nil && product.IndividualLongitude != nil {
			lat, lng := utils.GetCoordinatesPrivacy(*product.IndividualLatitude, *product.IndividualLongitude, *product.LocationPrivacy)
			product.IndividualLatitude = &lat
			product.IndividualLongitude = &lng
		}
	}
}

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(ctx context.Context, storefrontID, userID int, req *models.CreateProductRequest) (*models.StorefrontProduct, error) {
	// Validate ownership
	if err := s.ValidateStorefrontOwnership(ctx, storefrontID, userID); err != nil {
		return nil, fmt.Errorf("ownership validation failed: %w", err)
	}

	// Validate request
	if err := s.validateCreateRequest(req); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Validate variants if provided
	if req.HasVariants {
		if err := s.validateVariants(req); err != nil {
			return nil, fmt.Errorf("invalid variants: %w", err)
		}
	}

	// Start transaction
	tx, err := s.storage.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	// Create product within transaction
	product, err := s.storage.CreateStorefrontProductTx(ctx, tx, storefrontID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	// Create variants if provided
	if req.HasVariants && len(req.Variants) > 0 && s.variantService != nil {
		variantRequests := make([]variantTypes.CreateVariantRequest, len(req.Variants))
		for i, v := range req.Variants {
			variantRequests[i] = variantTypes.CreateVariantRequest{
				ProductID:         product.ID,
				SKU:               v.SKU,
				Barcode:           v.Barcode,
				Price:             v.Price,
				CompareAtPrice:    v.CompareAtPrice,
				CostPrice:         v.CostPrice,
				StockQuantity:     v.StockQuantity,
				LowStockThreshold: v.LowStockThreshold,
				VariantAttributes: v.VariantAttributes,
				Weight:            v.Weight,
				Dimensions:        v.Dimensions,
				IsDefault:         v.IsDefault,
			}
		}

		// Use transactional variant creation
		createdVariants, err := s.variantService.BulkCreateVariantsTx(ctx, tx, product.ID, variantRequests)
		if err != nil {
			return nil, fmt.Errorf("failed to create variants: %w", err)
		}

		// Convert variants to models.StorefrontProductVariant
		product.Variants = make([]models.StorefrontProductVariant, len(createdVariants))
		for i, v := range createdVariants {
			product.Variants[i] = models.StorefrontProductVariant{
				ID:                  v.ID,
				StorefrontProductID: v.ProductID,
				Name:                s.generateVariantName(v.VariantAttributes),
				SKU:                 v.SKU,
				Price:               s.getVariantPrice(v.Price, product.Price),
				StockQuantity:       v.StockQuantity,
				Attributes:          v.VariantAttributes,
				IsActive:            v.IsActive,
				CreatedAt:           v.CreatedAt,
				UpdatedAt:           v.UpdatedAt,
			}
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Index product in OpenSearch (after transaction is committed)
	if s.searchRepo != nil {
		go s.indexProductWithVariants(ctx, product)
	}

	return product, nil
}

// UpdateProduct updates an existing product
func (s *ProductService) UpdateProduct(ctx context.Context, storefrontID, productID, userID int, req *models.UpdateProductRequest) error {
	// Validate ownership
	if err := s.ValidateStorefrontOwnership(ctx, storefrontID, userID); err != nil {
		return err
	}

	// Validate request
	if err := s.validateUpdateRequest(req); err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}

	// Check product exists
	product, err := s.storage.GetStorefrontProduct(ctx, storefrontID, productID)
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}

	if product == nil {
		return fmt.Errorf("product not found")
	}

	// Update product
	if err := s.storage.UpdateStorefrontProduct(ctx, storefrontID, productID, req); err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	// Re-index product in OpenSearch
	if s.searchRepo != nil {
		// Get updated product for indexing
		updatedProduct, err := s.storage.GetStorefrontProduct(ctx, storefrontID, productID)
		if err != nil {
			logger.Error().Err(err).Msgf("Failed to get updated product %d for indexing", productID)
		} else if err := s.searchRepo.UpdateProduct(ctx, updatedProduct); err != nil {
			logger.Error().Err(err).Msgf("Failed to update product %d in OpenSearch", productID)
			// Не возвращаем ошибку, так как товар уже обновлен в БД
		} else {
			logger.Info().Msgf("Successfully updated product %d in OpenSearch", productID)
		}
	}

	return nil
}

// DeleteProduct deletes a product
func (s *ProductService) DeleteProduct(ctx context.Context, storefrontID, productID, userID int) error {
	// Validate ownership
	if err := s.ValidateStorefrontOwnership(ctx, storefrontID, userID); err != nil {
		return err
	}

	// Delete product
	if err := s.storage.DeleteStorefrontProduct(ctx, storefrontID, productID); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	// Delete from OpenSearch
	if s.searchRepo != nil {
		if err := s.searchRepo.DeleteProduct(ctx, productID); err != nil {
			logger.Error().Err(err).Msgf("Failed to delete product %d from OpenSearch", productID)
			// Не возвращаем ошибку, так как товар уже удален из БД
		} else {
			logger.Info().Msgf("Successfully deleted product %d from OpenSearch", productID)
		}
	}

	return nil
}

// UpdateInventory updates product stock
func (s *ProductService) UpdateInventory(ctx context.Context, storefrontID, productID, userID int, req *models.UpdateInventoryRequest) error {
	// Validate ownership
	if err := s.ValidateStorefrontOwnership(ctx, storefrontID, userID); err != nil {
		return err
	}

	// Validate request
	if err := s.validateInventoryRequest(req); err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}

	// Update inventory
	if err := s.storage.UpdateProductInventory(ctx, storefrontID, productID, userID, req); err != nil {
		return fmt.Errorf("failed to update inventory: %w", err)
	}

	// Частично обновляем продукт в OpenSearch (только поля склада)
	if s.searchRepo != nil {
		go func() {
			if err := s.updateProductStockInSearch(ctx, storefrontID, productID, req); err != nil {
				logger.Error().Err(err).Msgf("Failed to update product %d stock in OpenSearch", productID)
			} else {
				logger.Info().Msgf("Successfully updated product %d stock in OpenSearch", productID)
			}
		}()
	}

	return nil
}

// GetProductStats returns product statistics
func (s *ProductService) GetProductStats(ctx context.Context, storefrontID, userID int) (*models.ProductStats, error) {
	// Validate ownership
	if err := s.ValidateStorefrontOwnership(ctx, storefrontID, userID); err != nil {
		return nil, err
	}

	// Get stats
	stats, err := s.storage.GetProductStats(ctx, storefrontID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product stats: %w", err)
	}

	return stats, nil
}

// CreateProductForImport creates a product without ownership validation (for system imports)
func (s *ProductService) CreateProductForImport(ctx context.Context, storefrontID int, req *models.CreateProductRequest) (*models.StorefrontProduct, error) {
	// Validate request
	if err := s.validateCreateRequest(req); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Create product without ownership check
	product, err := s.storage.CreateStorefrontProduct(ctx, storefrontID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	// Index product in OpenSearch
	if s.searchRepo != nil {
		if err := s.searchRepo.IndexProduct(ctx, product); err != nil {
			logger.Error().Err(err).Msgf("Failed to index product %d in OpenSearch", product.ID)
			// Не возвращаем ошибку, так как товар уже создан в БД
		} else {
			logger.Info().Msgf("Successfully indexed product %d in OpenSearch", product.ID)
		}
	}

	return product, nil
}

// UpdateProductForImport updates a product without ownership validation (for system imports)
func (s *ProductService) UpdateProductForImport(ctx context.Context, storefrontID, productID int, req *models.UpdateProductRequest) error {
	// Validate request
	if err := s.validateUpdateRequest(req); err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}

	// Check product exists
	product, err := s.storage.GetStorefrontProduct(ctx, storefrontID, productID)
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}

	if product == nil {
		return fmt.Errorf("product not found")
	}

	// Update product without ownership check
	if err := s.storage.UpdateStorefrontProduct(ctx, storefrontID, productID, req); err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	// Re-index product in OpenSearch
	if s.searchRepo != nil {
		// Get updated product for indexing
		updatedProduct, err := s.storage.GetStorefrontProduct(ctx, storefrontID, productID)
		if err != nil {
			logger.Error().Err(err).Msgf("Failed to get updated product %d for indexing", productID)
		} else if err := s.searchRepo.UpdateProduct(ctx, updatedProduct); err != nil {
			logger.Error().Err(err).Msgf("Failed to update product %d in OpenSearch", productID)
			// Не возвращаем ошибку, так как товар уже обновлен в БД
		} else {
			logger.Info().Msgf("Successfully updated product %d in OpenSearch", productID)
		}
	}

	return nil
}

// Bulk operation methods

// BulkCreateProducts creates multiple products
func (s *ProductService) BulkCreateProducts(ctx context.Context, storefrontID, userID int, req models.BulkCreateProductsRequest) (*models.BulkCreateProductsResponse, error) {
	// Validate ownership
	if err := s.ValidateStorefrontOwnership(ctx, storefrontID, userID); err != nil {
		return nil, err
	}

	// Validate all products
	for i, product := range req.Products {
		if err := s.validateCreateRequest(&product); err != nil {
			return nil, fmt.Errorf("product %d: %w", i, err)
		}
	}

	// Create products
	createdIDs, errors := s.storage.BulkCreateProducts(ctx, storefrontID, req.Products)

	// Convert errors to response format
	var failedOps []models.BulkOperationError
	for i, err := range errors {
		if err != nil {
			failedOps = append(failedOps, models.BulkOperationError{
				Index: i,
				Error: err.Error(),
			})
		}
	}

	// Index created products in OpenSearch
	if len(createdIDs) > 0 && s.searchRepo != nil {
		go func() { //nolint:contextcheck // фоновая индексация
			for _, id := range createdIDs {
				product, err := s.storage.GetStorefrontProduct(context.Background(), storefrontID, id)
				if err != nil {
					logger.Error().Err(err).Msgf("Failed to get product %d for indexing", id)
					continue
				}
				if err := s.searchRepo.IndexProduct(context.Background(), product); err != nil {
					logger.Error().Err(err).Msgf("Failed to index product %d in OpenSearch", id)
				}
			}
		}()
	}

	return &models.BulkCreateProductsResponse{
		Created: createdIDs,
		Failed:  failedOps,
	}, nil
}

// BulkUpdateProducts updates multiple products
func (s *ProductService) BulkUpdateProducts(ctx context.Context, storefrontID, userID int, req models.BulkUpdateProductsRequest) (*models.BulkUpdateProductsResponse, error) {
	// Validate ownership
	if err := s.ValidateStorefrontOwnership(ctx, storefrontID, userID); err != nil {
		return nil, err
	}

	// Validate all updates
	for i, update := range req.Updates {
		if err := s.validateUpdateRequest(&update.Updates); err != nil {
			return nil, fmt.Errorf("update %d: %w", i, err)
		}
	}

	// Update products
	updatedIDs, errors := s.storage.BulkUpdateProducts(ctx, storefrontID, req.Updates)

	// Convert errors to response format
	var failedOps []models.BulkOperationError
	for i, err := range errors {
		if err != nil {
			failedOps = append(failedOps, models.BulkOperationError{
				ProductID: req.Updates[i].ProductID,
				Error:     err.Error(),
			})
		}
	}

	// Re-index updated products in OpenSearch
	if len(updatedIDs) > 0 && s.searchRepo != nil {
		go func() { //nolint:contextcheck // фоновая индексация
			for _, id := range updatedIDs {
				product, err := s.storage.GetStorefrontProduct(context.Background(), storefrontID, id)
				if err != nil {
					logger.Error().Err(err).Msgf("Failed to get product %d for re-indexing", id)
					continue
				}
				if err := s.searchRepo.UpdateProduct(context.Background(), product); err != nil {
					logger.Error().Err(err).Msgf("Failed to update product %d in OpenSearch", id)
				}
			}
		}()
	}

	return &models.BulkUpdateProductsResponse{
		Updated: updatedIDs,
		Failed:  failedOps,
	}, nil
}

// BulkDeleteProducts deletes multiple products
func (s *ProductService) BulkDeleteProducts(ctx context.Context, storefrontID, userID int, req models.BulkDeleteProductsRequest) (*models.BulkDeleteProductsResponse, error) {
	// Validate ownership
	if err := s.ValidateStorefrontOwnership(ctx, storefrontID, userID); err != nil {
		return nil, err
	}

	// Delete products
	deletedIDs, errors := s.storage.BulkDeleteProducts(ctx, storefrontID, req.ProductIDs)

	// Convert errors to response format
	var failedOps []models.BulkOperationError
	for i, err := range errors {
		if err != nil {
			failedOps = append(failedOps, models.BulkOperationError{
				ProductID: req.ProductIDs[i],
				Error:     err.Error(),
			})
		}
	}

	// Remove deleted products from OpenSearch
	if len(deletedIDs) > 0 && s.searchRepo != nil {
		go func() { //nolint:contextcheck // фоновое удаление из индекса
			for _, id := range deletedIDs {
				if err := s.searchRepo.DeleteProduct(context.Background(), id); err != nil {
					logger.Error().Err(err).Msgf("Failed to delete product %d from OpenSearch", id)
				}
			}
		}()
	}

	return &models.BulkDeleteProductsResponse{
		Deleted: deletedIDs,
		Failed:  failedOps,
	}, nil
}

// BulkUpdateStatus updates status of multiple products
func (s *ProductService) BulkUpdateStatus(ctx context.Context, storefrontID, userID int, req models.BulkUpdateStatusRequest) (*models.BulkUpdateStatusResponse, error) {
	// Validate ownership
	if err := s.ValidateStorefrontOwnership(ctx, storefrontID, userID); err != nil {
		return nil, err
	}

	// Update status
	updatedIDs, errors := s.storage.BulkUpdateStatus(ctx, storefrontID, req.ProductIDs, req.IsActive)

	// Convert errors to response format
	var failedOps []models.BulkOperationError
	for i, err := range errors {
		if err != nil {
			failedOps = append(failedOps, models.BulkOperationError{
				ProductID: req.ProductIDs[i],
				Error:     err.Error(),
			})
		}
	}

	// Re-index updated products in OpenSearch
	if len(updatedIDs) > 0 && s.searchRepo != nil {
		go func() { //nolint:contextcheck // фоновая индексация
			for _, id := range updatedIDs {
				product, err := s.storage.GetStorefrontProduct(context.Background(), storefrontID, id)
				if err != nil {
					logger.Error().Err(err).Msgf("Failed to get product %d for re-indexing", id)
					continue
				}
				if err := s.searchRepo.UpdateProduct(context.Background(), product); err != nil {
					logger.Error().Err(err).Msgf("Failed to update product %d in OpenSearch", id)
				}
			}
		}()
	}

	return &models.BulkUpdateStatusResponse{
		Updated: updatedIDs,
		Failed:  failedOps,
	}, nil
}

// ReindexProduct переиндексирует товар в OpenSearch
func (s *ProductService) ReindexProduct(ctx context.Context, storefrontID, productID int) error {
	// Получаем товар со всеми данными, включая изображения
	product, err := s.storage.GetStorefrontProduct(ctx, storefrontID, productID)
	if err != nil {
		logger.Error().Err(err).Msgf("Failed to get product %d for reindexing", productID)
		return err
	}

	if product == nil {
		logger.Warn().Msgf("Product %d not found for reindexing", productID)
		return nil
	}

	// Переиндексируем в OpenSearch
	if s.searchRepo != nil {
		if err := s.searchRepo.UpdateProduct(ctx, product); err != nil {
			logger.Error().Err(err).Msgf("Failed to reindex product %d in OpenSearch", productID)
			return err
		}
		logger.Info().Msgf("Successfully reindexed product %d in OpenSearch", productID)
	}

	return nil
}

// Validation helpers

func (s *ProductService) validateCreateRequest(req *models.CreateProductRequest) error {
	if req.Name == "" {
		return fmt.Errorf("name is required")
	}

	if len(req.Name) < 3 || len(req.Name) > 255 {
		return fmt.Errorf("name must be between 3 and 255 characters")
	}

	if req.Description == "" {
		return fmt.Errorf("description is required")
	}

	if len(req.Description) < 10 {
		return fmt.Errorf("description must be at least 10 characters")
	}

	if req.Price < 0 {
		return fmt.Errorf("price cannot be negative")
	}

	if req.Currency == "" {
		return fmt.Errorf("currency is required")
	}

	if len(req.Currency) != 3 {
		return fmt.Errorf("currency must be 3 characters")
	}

	if req.CategoryID <= 0 {
		return fmt.Errorf("category is required")
	}

	if req.StockQuantity < 0 {
		return fmt.Errorf("stock quantity cannot be negative")
	}

	return nil
}

func (s *ProductService) validateUpdateRequest(req *models.UpdateProductRequest) error {
	if req.Name != nil {
		if *req.Name == "" {
			return fmt.Errorf("name cannot be empty")
		}

		if len(*req.Name) < 3 || len(*req.Name) > 255 {
			return fmt.Errorf("name must be between 3 and 255 characters")
		}
	}

	if req.Description != nil {
		if *req.Description == "" {
			return fmt.Errorf("description cannot be empty")
		}

		if len(*req.Description) < 10 {
			return fmt.Errorf("description must be at least 10 characters")
		}
	}

	if req.Price != nil && *req.Price < 0 {
		return fmt.Errorf("price cannot be negative")
	}

	if req.CategoryID != nil && *req.CategoryID <= 0 {
		return fmt.Errorf("invalid category")
	}

	if req.StockQuantity != nil && *req.StockQuantity < 0 {
		return fmt.Errorf("stock quantity cannot be negative")
	}

	return nil
}

func (s *ProductService) validateInventoryRequest(req *models.UpdateInventoryRequest) error {
	if req.Type != "in" && req.Type != "out" && req.Type != "adjustment" {
		return fmt.Errorf("invalid inventory type")
	}

	if req.Quantity <= 0 {
		return fmt.Errorf("quantity must be positive")
	}

	if req.Reason == "" {
		return fmt.Errorf("reason is required")
	}

	return nil
}

// validateVariants validates variants before creating product
func (s *ProductService) validateVariants(req *models.CreateProductRequest) error {
	if len(req.Variants) == 0 {
		return errors.New("at least one variant is required when has_variants is true")
	}

	// Check for duplicate attribute combinations
	seen := make(map[string]bool)
	defaultCount := 0

	for i, v := range req.Variants {
		// Create key from attributes for uniqueness check
		attrKey := s.generateAttributeKey(v.VariantAttributes)
		if seen[attrKey] {
			return fmt.Errorf("duplicate variant attributes at index %d", i)
		}
		seen[attrKey] = true

		// Check that only one variant is default
		if v.IsDefault {
			defaultCount++
		}

		// Validate SKU if provided
		if v.SKU != nil && *v.SKU != "" {
			if err := s.validateSKU(*v.SKU); err != nil {
				return fmt.Errorf("invalid SKU at index %d: %w", i, err)
			}
		}
	}

	if defaultCount > 1 {
		return errors.New("only one variant can be default")
	}

	if defaultCount == 0 && len(req.Variants) > 0 {
		// Make first variant default
		req.Variants[0].IsDefault = true
	}

	return nil
}

// generateAttributeKey creates a unique key from variant attributes
func (s *ProductService) generateAttributeKey(attributes map[string]interface{}) string {
	// Sort attribute keys for consistent ordering
	keys := make([]string, 0, len(attributes))
	for k := range attributes {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build key string
	var parts []string
	for _, k := range keys {
		v := attributes[k]
		parts = append(parts, fmt.Sprintf("%s:%v", k, v))
	}

	keyStr := strings.Join(parts, ";")

	// Create SHA256 hash for key
	hasher := sha256.New()
	hasher.Write([]byte(keyStr))
	return hex.EncodeToString(hasher.Sum(nil))[:16]
}

// validateSKU validates SKU format
func (s *ProductService) validateSKU(sku string) error {
	if len(sku) < 3 || len(sku) > 100 {
		return errors.New("SKU must be between 3 and 100 characters")
	}

	// Check for valid characters (alphanumeric, dash, underscore)
	for _, r := range sku {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') && r != '-' && r != '_' {
			return errors.New("SKU can only contain letters, numbers, dashes and underscores")
		}
	}

	return nil
}

// generateVariantName creates a display name from variant attributes
func (s *ProductService) generateVariantName(attributes map[string]interface{}) string {
	// Sort attribute keys for consistent ordering
	keys := make([]string, 0, len(attributes))
	for k := range attributes {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build name from attribute values
	var parts []string
	for _, k := range keys {
		v := attributes[k]
		parts = append(parts, fmt.Sprintf("%v", v))
	}

	return strings.Join(parts, " - ")
}

// getVariantPrice returns variant price or base product price
func (s *ProductService) getVariantPrice(variantPrice *float64, basePrice float64) float64 {
	if variantPrice != nil && *variantPrice > 0 {
		return *variantPrice
	}
	return basePrice
}

// indexProductWithVariants indexes product with variants in OpenSearch
func (s *ProductService) indexProductWithVariants(ctx context.Context, product *models.StorefrontProduct) {
	if err := s.searchRepo.IndexProduct(ctx, product); err != nil {
		logger.Error().Err(err).Msgf("Failed to index product %d with variants in OpenSearch", product.ID)
	} else {
		logger.Info().Msgf("Successfully indexed product %d with %d variants in OpenSearch", product.ID, len(product.Variants))
	}
}

// updateProductStockInSearch частично обновляет только поля склада в OpenSearch
func (s *ProductService) updateProductStockInSearch(ctx context.Context, storefrontID, productID int, req *models.UpdateInventoryRequest) error {
	// Получаем актуальное состояние продукта после обновления склада
	product, err := s.storage.GetStorefrontProduct(ctx, storefrontID, productID)
	if err != nil {
		return fmt.Errorf("failed to get updated product for indexing: %w", err)
	}

	// Проверяем, поддерживает ли OpenSearch репозиторий частичные обновления
	if partialUpdater, ok := s.searchRepo.(interface {
		UpdateProductStock(ctx context.Context, productID int, stockData map[string]interface{}) error
	}); ok {
		// Используем специализированный метод для частичного обновления остатков
		stockData := map[string]interface{}{
			"stock_quantity": product.StockQuantity,
			"stock_status":   product.GetStockStatus(),
		}

		return partialUpdater.UpdateProductStock(ctx, productID, stockData)
	}

	// Fallback: полное переиндексирование продукта
	return s.searchRepo.IndexProduct(ctx, product)
}
