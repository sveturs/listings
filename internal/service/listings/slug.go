package listings

import (
	"context"
	"fmt"

	"github.com/gosimple/slug"
)

// SlugGenerator handles slug generation with uniqueness guarantees
type SlugGenerator struct {
	repo Repository
}

// NewSlugGenerator creates a new slug generator instance
func NewSlugGenerator(repo Repository) *SlugGenerator {
	return &SlugGenerator{
		repo: repo,
	}
}

// Generate creates a unique slug from a title
// If the base slug already exists, it appends a counter (e.g., "title-1", "title-2")
func (s *SlugGenerator) Generate(ctx context.Context, title string) (string, error) {
	if title == "" {
		return "", fmt.Errorf("title cannot be empty")
	}

	// Generate base slug from title
	baseSlug := slug.Make(title)

	if baseSlug == "" {
		return "", fmt.Errorf("failed to generate slug from title")
	}

	// Check if base slug is unique
	exists, err := s.slugExists(ctx, baseSlug)
	if err != nil {
		return "", fmt.Errorf("failed to check slug existence: %w", err)
	}

	if !exists {
		return baseSlug, nil
	}

	// Base slug exists, try with counter suffix
	for i := 1; i < 1000; i++ {
		candidate := fmt.Sprintf("%s-%d", baseSlug, i)

		exists, err := s.slugExists(ctx, candidate)
		if err != nil {
			return "", fmt.Errorf("failed to check slug existence: %w", err)
		}

		if !exists {
			return candidate, nil
		}
	}

	// If we exhausted 1000 attempts, return error
	return "", fmt.Errorf("failed to generate unique slug after 1000 attempts")
}

// GenerateWithExclusion generates a unique slug excluding a specific listing ID
// Useful when updating a listing - we don't want to conflict with other listings,
// but the slug can stay the same for the listing being updated
func (s *SlugGenerator) GenerateWithExclusion(ctx context.Context, title string, excludeListingID int64) (string, error) {
	if title == "" {
		return "", fmt.Errorf("title cannot be empty")
	}

	baseSlug := slug.Make(title)

	if baseSlug == "" {
		return "", fmt.Errorf("failed to generate slug from title")
	}

	// Check if base slug exists for a different listing
	listing, err := s.repo.GetListingBySlug(ctx, baseSlug)
	if err != nil {
		// Slug doesn't exist - we can use it
		return baseSlug, nil
	}

	// If the listing with this slug is the one we're updating, that's OK
	if listing.ID == excludeListingID {
		return baseSlug, nil
	}

	// Slug exists for a different listing, try with counter
	for i := 1; i < 1000; i++ {
		candidate := fmt.Sprintf("%s-%d", baseSlug, i)

		listing, err := s.repo.GetListingBySlug(ctx, candidate)
		if err != nil {
			// Slug doesn't exist - we can use it
			return candidate, nil
		}

		// If this slug belongs to the listing we're updating, that's OK
		if listing.ID == excludeListingID {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("failed to generate unique slug after 1000 attempts")
}

// slugExists checks if a slug already exists in the database
func (s *SlugGenerator) slugExists(ctx context.Context, slug string) (bool, error) {
	// Try to get listing by slug
	_, err := s.repo.GetListingBySlug(ctx, slug)
	if err != nil {
		// If not found, slug doesn't exist - that's good!
		return false, nil
	}

	// Listing found with this slug
	return true, nil
}

// ValidateSlug checks if a given slug is valid and unique
func (s *SlugGenerator) ValidateSlug(ctx context.Context, slug string) error {
	if !ValidateSlug(slug) {
		return fmt.Errorf("invalid slug format: must be lowercase alphanumeric with hyphens")
	}

	exists, err := s.slugExists(ctx, slug)
	if err != nil {
		return fmt.Errorf("failed to check slug existence: %w", err)
	}

	if exists {
		return fmt.Errorf("slug already exists")
	}

	return nil
}
