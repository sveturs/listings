package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"backend/internal/domain/models"
	marketplaceServices "backend/internal/proj/marketplace/services"
	"backend/internal/storage/postgres"

	"go.uber.org/zap"
)

// CategoryMappingService handles category mapping operations for storefront imports
type CategoryMappingService struct {
	categoryRepo postgres.CategoryMappingsRepositoryInterface
	aiDetector   *marketplaceServices.AICategoryDetector
	logger       *zap.Logger
}

// NewCategoryMappingService creates a new CategoryMappingService
func NewCategoryMappingService(
	categoryRepo postgres.CategoryMappingsRepositoryInterface,
	aiDetector *marketplaceServices.AICategoryDetector,
	logger *zap.Logger,
) *CategoryMappingService {
	return &CategoryMappingService{
		categoryRepo: categoryRepo,
		aiDetector:   aiDetector,
		logger:       logger,
	}
}

// GetOrCreateMapping retrieves existing mapping or creates new one via AI detection
// sourceCategoryPath: external category path from import (e.g., "Electronics/Phones/Apple")
// productTitle: product title for better AI detection context
// productDescription: product description for better AI detection context
func (s *CategoryMappingService) GetOrCreateMapping(
	ctx context.Context,
	storefrontID int,
	sourceCategoryPath string,
	productTitle string,
	productDescription string,
) (int, error) {
	// Normalize source category path
	normalizedPath := s.normalizeCategoryPath(sourceCategoryPath)

	// Try to find existing mapping
	existing, err := s.categoryRepo.GetBySourcePath(ctx, storefrontID, normalizedPath)
	if err == nil && existing != nil {
		s.logger.Debug("Found existing category mapping",
			zap.Int("storefront_id", storefrontID),
			zap.String("source_path", normalizedPath),
			zap.Int("target_category_id", existing.TargetCategoryID),
		)
		return existing.TargetCategoryID, nil
	}

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("failed to get existing mapping: %w", err)
	}

	// No existing mapping found, use AI detection
	s.logger.Info("No mapping found, using AI detection",
		zap.Int("storefront_id", storefrontID),
		zap.String("source_path", normalizedPath),
	)

	categoryID, confidence, err := s.detectCategoryViaAI(ctx, sourceCategoryPath, productTitle, productDescription)
	if err != nil {
		return 0, fmt.Errorf("failed to detect category via AI: %w", err)
	}

	// Save the new mapping for future use
	mapping := &models.StorefrontCategoryMapping{
		StorefrontID:       storefrontID,
		SourceCategoryPath: normalizedPath,
		TargetCategoryID:   categoryID,
		IsManual:           false,
		ConfidenceScore:    &confidence,
	}

	if err := s.categoryRepo.Create(ctx, mapping); err != nil {
		s.logger.Warn("Failed to save category mapping",
			zap.Error(err),
			zap.Int("storefront_id", storefrontID),
			zap.String("source_path", normalizedPath),
		)
		// Non-critical error, we can still return the detected category
	} else {
		s.logger.Info("Saved new category mapping",
			zap.Int("storefront_id", storefrontID),
			zap.String("source_path", normalizedPath),
			zap.Int("target_category_id", categoryID),
			zap.Float64("confidence", confidence),
		)
	}

	return categoryID, nil
}

// CreateManualMapping creates a manual category mapping
func (s *CategoryMappingService) CreateManualMapping(
	ctx context.Context,
	storefrontID int,
	sourceCategoryPath string,
	targetCategoryID int,
) error {
	normalizedPath := s.normalizeCategoryPath(sourceCategoryPath)

	mapping := &models.StorefrontCategoryMapping{
		StorefrontID:       storefrontID,
		SourceCategoryPath: normalizedPath,
		TargetCategoryID:   targetCategoryID,
		IsManual:           true,
		ConfidenceScore:    nil,
	}

	if err := s.categoryRepo.Create(ctx, mapping); err != nil {
		return fmt.Errorf("failed to create manual mapping: %w", err)
	}

	s.logger.Info("Created manual category mapping",
		zap.Int("storefront_id", storefrontID),
		zap.String("source_path", normalizedPath),
		zap.Int("target_category_id", targetCategoryID),
	)

	return nil
}

// UpdateMapping updates an existing category mapping
func (s *CategoryMappingService) UpdateMapping(
	ctx context.Context,
	id int,
	targetCategoryID int,
	isManual bool,
) error {
	mapping, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get mapping: %w", err)
	}

	mapping.TargetCategoryID = targetCategoryID
	mapping.IsManual = isManual
	if isManual {
		mapping.ConfidenceScore = nil
	}

	if err := s.categoryRepo.Update(ctx, mapping); err != nil {
		return fmt.Errorf("failed to update mapping: %w", err)
	}

	s.logger.Info("Updated category mapping",
		zap.Int("id", id),
		zap.Int("target_category_id", targetCategoryID),
		zap.Bool("is_manual", isManual),
	)

	return nil
}

// DeleteMapping deletes a category mapping
func (s *CategoryMappingService) DeleteMapping(ctx context.Context, id int) error {
	if err := s.categoryRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete mapping: %w", err)
	}

	s.logger.Info("Deleted category mapping", zap.Int("id", id))
	return nil
}

// GetMappings retrieves all mappings for a storefront with optional filter
func (s *CategoryMappingService) GetMappings(
	ctx context.Context,
	storefrontID int,
	filter *postgres.CategoryMappingFilter,
) ([]*models.StorefrontCategoryMappingWithDetails, int, error) {
	mappings, total, err := s.categoryRepo.GetByStorefront(ctx, storefrontID, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get mappings: %w", err)
	}

	return mappings, total, nil
}

// GetMappingByID retrieves a mapping by ID
func (s *CategoryMappingService) GetMappingByID(ctx context.Context, id int) (*models.StorefrontCategoryMapping, error) {
	mapping, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get mapping: %w", err)
	}

	return mapping, nil
}

// detectCategoryViaAI uses AI detector to determine category
func (s *CategoryMappingService) detectCategoryViaAI(
	ctx context.Context,
	sourceCategoryPath string,
	productTitle string,
	productDescription string,
) (categoryID int, confidence float64, err error) {
	// Parse category path into hints
	hints := s.parseCategoryPathToHints(sourceCategoryPath)

	// Prepare AI detection input
	input := marketplaceServices.AIDetectionInput{
		Title:       productTitle,
		Description: productDescription,
		AIHints:     hints,
		EntityType:  "product",
	}

	// Detect category via AI
	result, err := s.aiDetector.DetectCategory(ctx, input)
	if err != nil {
		return 0, 0, fmt.Errorf("AI detection failed: %w", err)
	}

	if result == nil || result.CategoryID == 0 {
		return 0, 0, errors.New("AI detector returned no result")
	}

	s.logger.Debug("AI category detection result",
		zap.Int32("category_id", result.CategoryID),
		zap.String("category_name", result.CategoryName),
		zap.Float64("confidence", result.ConfidenceScore),
	)

	return int(result.CategoryID), result.ConfidenceScore, nil
}

// parseCategoryPathToHints parses category path into AI hints
// Example: "Electronics/Phones/Apple" -> domain="Electronics", productType="Phones", keywords=["Apple"]
func (s *CategoryMappingService) parseCategoryPathToHints(sourceCategoryPath string) *marketplaceServices.AIHints {
	// Support different separators: /, >, |, -
	separators := []string{"/", ">", "|", "-"}

	var parts []string
	for _, sep := range separators {
		if strings.Contains(sourceCategoryPath, sep) {
			parts = strings.Split(sourceCategoryPath, sep)
			break
		}
	}

	if len(parts) == 0 {
		parts = []string{sourceCategoryPath}
	}

	// Trim and clean parts
	cleanedParts := make([]string, 0, len(parts))
	for _, part := range parts {
		cleaned := strings.TrimSpace(part)
		if cleaned != "" {
			cleanedParts = append(cleanedParts, cleaned)
		}
	}

	hints := &marketplaceServices.AIHints{
		CategoryPath: sourceCategoryPath,
	}

	if len(cleanedParts) >= 1 {
		hints.Domain = cleanedParts[0]
	}

	if len(cleanedParts) >= 2 {
		hints.ProductType = cleanedParts[1]
	}

	if len(cleanedParts) >= 3 {
		hints.Keywords = cleanedParts[2:]
	}

	return hints
}

// normalizeCategoryPath normalizes category path for consistent storage
// Converts all separators to "/" and trims whitespace
func (s *CategoryMappingService) normalizeCategoryPath(path string) string {
	// Replace common separators with /
	normalized := path
	normalized = strings.ReplaceAll(normalized, " > ", "/")
	normalized = strings.ReplaceAll(normalized, " | ", "/")
	normalized = strings.ReplaceAll(normalized, " - ", "/")

	// Trim whitespace around slashes
	parts := strings.Split(normalized, "/")
	cleanedParts := make([]string, 0, len(parts))
	for _, part := range parts {
		cleaned := strings.TrimSpace(part)
		if cleaned != "" {
			cleanedParts = append(cleanedParts, cleaned)
		}
	}

	return strings.Join(cleanedParts, "/")
}
