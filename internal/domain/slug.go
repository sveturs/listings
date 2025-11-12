// Package domain provides slug generation and validation utilities
package domain

import (
	"fmt"
	"regexp"
	"strings"
)

// SlugGenerator provides slug generation functionality
type SlugGenerator struct{}

// NewSlugGenerator creates a new slug generator
func NewSlugGenerator() *SlugGenerator {
	return &SlugGenerator{}
}

// Generate creates a URL-friendly slug from a title
// Rules:
// - Convert to lowercase
// - Replace spaces with hyphens
// - Remove all non-alphanumeric characters except hyphens
// - Remove multiple consecutive hyphens
// - Trim hyphens from start and end
// - Limit to 200 characters
func (sg *SlugGenerator) Generate(title string) string {
	if title == "" {
		return ""
	}

	// Convert to lowercase
	s := strings.ToLower(title)

	// Remove special characters (keep alphanumeric, spaces, and hyphens)
	s = regexp.MustCompile(`[^a-z0-9\s-]`).ReplaceAllString(s, "")

	// Replace spaces with hyphens
	s = regexp.MustCompile(`\s+`).ReplaceAllString(s, "-")

	// Replace multiple consecutive hyphens with single hyphen
	s = regexp.MustCompile(`-+`).ReplaceAllString(s, "-")

	// Trim hyphens from start and end
	s = strings.Trim(s, "-")

	// Limit length to 200 characters to leave room for numeric suffixes
	if len(s) > 200 {
		s = s[:200]
		s = strings.TrimRight(s, "-")
	}

	return s
}

// Validate checks if a slug is valid according to our rules
func (sg *SlugGenerator) Validate(slug string) error {
	if slug == "" {
		return fmt.Errorf("slug cannot be empty")
	}

	if len(slug) < 3 {
		return fmt.Errorf("slug must be at least 3 characters")
	}

	if len(slug) > 255 {
		return fmt.Errorf("slug cannot exceed 255 characters")
	}

	// Check format: lowercase alphanumeric + hyphens only
	if !regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`).MatchString(slug) {
		return fmt.Errorf("slug must contain only lowercase letters, numbers, and hyphens (no consecutive hyphens)")
	}

	return nil
}

// GenerateWithSuffix appends a numeric suffix to a base slug
func (sg *SlugGenerator) GenerateWithSuffix(baseSlug string, suffix int) string {
	if suffix <= 0 {
		return baseSlug
	}
	return fmt.Sprintf("%s-%d", baseSlug, suffix)
}

// Slugify is a convenience function for quick slug generation
func Slugify(title string) string {
	sg := NewSlugGenerator()
	return sg.Generate(title)
}

// ValidateSlug is a convenience function for quick slug validation
func ValidateSlug(slug string) error {
	sg := NewSlugGenerator()
	return sg.Validate(slug)
}
