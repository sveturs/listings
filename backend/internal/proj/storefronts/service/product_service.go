package service

import (
	"context"
	"fmt"

	"backend/internal/domain/models"
)

// ProductService handles business logic for storefront products
type ProductService struct {
	storage Storage
}

// Storage interface for product operations
type Storage interface {
	GetStorefrontProducts(ctx context.Context, filter models.ProductFilter) ([]*models.StorefrontProduct, error)
	GetStorefrontProduct(ctx context.Context, storefrontID, productID int) (*models.StorefrontProduct, error)
	CreateStorefrontProduct(ctx context.Context, storefrontID int, req *models.CreateProductRequest) (*models.StorefrontProduct, error)
	UpdateStorefrontProduct(ctx context.Context, storefrontID, productID int, req *models.UpdateProductRequest) error
	DeleteStorefrontProduct(ctx context.Context, storefrontID, productID int) error
	UpdateProductInventory(ctx context.Context, storefrontID, productID int, userID int, req *models.UpdateInventoryRequest) error
	GetProductStats(ctx context.Context, storefrontID int) (*models.ProductStats, error)
	
	// Storefront operations
	GetStorefrontByID(ctx context.Context, id int) (*models.Storefront, error)
}

// NewProductService creates a new product service
func NewProductService(storage Storage) *ProductService {
	return &ProductService{
		storage: storage,
	}
}

// ValidateStorefrontOwnership checks if user owns the storefront
func (s *ProductService) ValidateStorefrontOwnership(ctx context.Context, storefrontID, userID int) error {
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
	// Validate storefront exists
	storefront, err := s.storage.GetStorefrontByID(ctx, filter.StorefrontID)
	if err != nil {
		return nil, fmt.Errorf("failed to get storefront: %w", err)
	}
	
	if storefront == nil {
		return nil, fmt.Errorf("storefront not found")
	}
	
	// Get products
	products, err := s.storage.GetStorefrontProducts(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}
	
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
	
	return product, nil
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
	
	// Create product
	product, err := s.storage.CreateStorefrontProduct(ctx, storefrontID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
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