package service

import (
	"context"
	"fmt"

	"backend/internal/logger"
	"backend/internal/proj/storefront/repository"
	variantTypes "backend/internal/proj/storefront/types"
)

// VariantService handles business logic for product variants
type VariantServiceImpl struct {
	variantRepo *repository.VariantRepository
}

// NewVariantService creates a new variant service
func NewVariantService(variantRepo *repository.VariantRepository) *VariantServiceImpl {
	return &VariantServiceImpl{
		variantRepo: variantRepo,
	}
}

// BulkCreateVariants creates multiple variants for a product
func (s *VariantServiceImpl) BulkCreateVariants(ctx context.Context, productID int, variants []variantTypes.CreateVariantRequest) ([]*variantTypes.ProductVariant, error) {
	logger.Info().
		Int("product_id", productID).
		Int("variants_count", len(variants)).
		Msg("Starting bulk variant creation")
		
	// Set productID for all variants
	for i := range variants {
		variants[i].ProductID = productID
	}

	// Create variants one by one (repository doesn't have bulk create)
	var createdVariants []*variantTypes.ProductVariant
	for i, variantReq := range variants {
		logger.Debug().
			Int("variant_index", i).
			Interface("variant_request", variantReq).
			Msg("Creating variant")
			
		variant, err := s.variantRepo.CreateVariant(ctx, &variantReq)
		if err != nil {
			logger.Error().
				Err(err).
				Int("variant_index", i).
				Interface("variant_request", variantReq).
				Msg("Failed to create variant")
			return nil, fmt.Errorf("failed to create variant %d with SKU %v: %w", i, variantReq.SKU, err)
		}
		createdVariants = append(createdVariants, variant)
	}

	return createdVariants, nil
}

// CreateVariant creates a single variant
func (s *VariantServiceImpl) CreateVariant(ctx context.Context, req *variantTypes.CreateVariantRequest) (*variantTypes.ProductVariant, error) {
	variant, err := s.variantRepo.CreateVariant(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create variant: %w", err)
	}

	return variant, nil
}

// GetVariantsByProductID retrieves all variants for a product
func (s *VariantServiceImpl) GetVariantsByProductID(ctx context.Context, productID int) ([]*variantTypes.ProductVariant, error) {
	variants, err := s.variantRepo.GetVariantsByProductID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get variants for product %d: %w", productID, err)
	}

	// Convert to pointers
	result := make([]*variantTypes.ProductVariant, len(variants))
	for i := range variants {
		result[i] = &variants[i]
	}

	return result, nil
}

// UpdateVariant updates a variant
func (s *VariantServiceImpl) UpdateVariant(ctx context.Context, variantID int, req *variantTypes.UpdateVariantRequest) error {
	_, err := s.variantRepo.UpdateVariant(ctx, variantID, req)
	if err != nil {
		return fmt.Errorf("failed to update variant %d: %w", variantID, err)
	}

	return nil
}

// DeleteVariant deletes a variant
func (s *VariantServiceImpl) DeleteVariant(ctx context.Context, variantID int) error {
	if err := s.variantRepo.DeleteVariant(ctx, variantID); err != nil {
		return fmt.Errorf("failed to delete variant %d: %w", variantID, err)
	}

	return nil
}

// GenerateVariants generates variant combinations based on attribute matrix
func (s *VariantServiceImpl) GenerateVariants(ctx context.Context, req *variantTypes.GenerateVariantsRequest) ([]*variantTypes.ProductVariant, error) {
	// Generate variants using repository
	variants, err := s.variantRepo.GenerateVariants(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate variants: %w", err)
	}

	// Convert to pointers
	result := make([]*variantTypes.ProductVariant, len(variants))
	for i := range variants {
		result[i] = &variants[i]
	}

	return result, nil
}
