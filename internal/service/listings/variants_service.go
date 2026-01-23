package listings

import (
	"context"
	"fmt"

	"github.com/vondi-global/listings/internal/domain"
)

// CreateB2CProductVariant creates a new B2C product variant with validation
func (s *Service) CreateB2CProductVariant(ctx context.Context, variant *domain.Variant) (*domain.Variant, error) {
	// Validation is handled in gRPC layer, service is a thin wrapper
	created, err := s.repo.CreateVariant(ctx, variant)
	if err != nil {
		s.logger.Error().Err(err).Int64("product_id", variant.ProductID).Msg("failed to create variant")
		return nil, fmt.Errorf("failed to create variant: %w", err)
	}

	s.logger.Info().Int64("variant_id", created.ID).Int64("product_id", variant.ProductID).Msg("variant created")
	return created, nil
}

// GetB2CProductVariant retrieves a B2C variant by ID
func (s *Service) GetB2CProductVariant(ctx context.Context, id int64) (*domain.Variant, error) {
	variant, err := s.repo.GetVariant(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).Int64("variant_id", id).Msg("failed to get variant")
		return nil, fmt.Errorf("failed to get variant: %w", err)
	}

	return variant, nil
}

// UpdateB2CProductVariant updates an existing B2C variant with partial updates
func (s *Service) UpdateB2CProductVariant(ctx context.Context, id int64, update *domain.VariantUpdate) (*domain.Variant, error) {
	// Validation is handled in gRPC layer
	updated, err := s.repo.UpdateB2CVariant(ctx, id, update)
	if err != nil {
		s.logger.Error().Err(err).Int64("variant_id", id).Msg("failed to update variant")
		return nil, fmt.Errorf("failed to update variant: %w", err)
	}

	s.logger.Info().Int64("variant_id", updated.ID).Msg("variant updated")
	return updated, nil
}

// DeleteB2CProductVariant deletes a B2C variant by ID
func (s *Service) DeleteB2CProductVariant(ctx context.Context, id int64) error {
	err := s.repo.DeleteB2CVariant(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).Int64("variant_id", id).Msg("failed to delete variant")
		return fmt.Errorf("failed to delete variant: %w", err)
	}

	s.logger.Info().Int64("variant_id", id).Msg("variant deleted")
	return nil
}

// ListB2CProductVariants retrieves all B2C variants for a product with optional filters
func (s *Service) ListB2CProductVariants(ctx context.Context, filters *domain.VariantFilters) ([]*domain.Variant, error) {
	variants, err := s.repo.ListVariants(ctx, filters)
	if err != nil {
		s.logger.Error().Err(err).Int64("product_id", filters.ProductID).Msg("failed to list variants")
		return nil, fmt.Errorf("failed to list variants: %w", err)
	}

	s.logger.Debug().Int64("product_id", filters.ProductID).Int("count", len(variants)).Msg("variants listed")
	return variants, nil
}
