// Package postgres implements slug-related listing repository methods.
// This file contains operations for URL slug management and retrieval.
package postgres

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

// SlugExists checks if a slug already exists in the database
// excludeID allows checking uniqueness when updating a listing (exclude the listing being updated)
func (r *Repository) SlugExists(ctx context.Context, slug string, excludeID ...int64) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM listings
			WHERE slug = $1
			  AND is_deleted = false
			  AND ($2::bigint IS NULL OR id != $2)
		)
	`

	var excludeIDPtr *int64
	if len(excludeID) > 0 && excludeID[0] > 0 {
		excludeIDPtr = &excludeID[0]
	}

	var exists bool
	err := r.db.QueryRowContext(ctx, query, slug, excludeIDPtr).Scan(&exists)
	if err != nil {
		r.logger.Error().Err(err).Str("slug", slug).Msg("failed to check slug existence")
		return false, fmt.Errorf("failed to check slug existence: %w", err)
	}

	return exists, nil
}

// GenerateUniqueSlug generates a unique URL slug from a title
// It will append a number suffix if the slug already exists
func (r *Repository) GenerateUniqueSlug(ctx context.Context, title string, excludeID ...int64) (string, error) {
	baseSlug := slugify(title)
	if baseSlug == "" {
		return "", fmt.Errorf("cannot generate slug from empty title")
	}

	slug := baseSlug
	counter := 1

	// Try up to 100 variations
	for counter <= 100 {
		exists, err := r.SlugExists(ctx, slug, excludeID...)
		if err != nil {
			return "", fmt.Errorf("failed to check slug uniqueness: %w", err)
		}

		if !exists {
			return slug, nil
		}

		// Try with counter suffix
		slug = fmt.Sprintf("%s-%d", baseSlug, counter)
		counter++
	}

	return "", fmt.Errorf("failed to generate unique slug after 100 attempts")
}

// transliterationMap contains mappings for Cyrillic to Latin transliteration
// Supports Russian and Serbian Cyrillic characters
var transliterationMap = map[rune]string{
	// Russian Cyrillic (lowercase)
	'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "yo",
	'ж': "zh", 'з': "z", 'и': "i", 'й': "j", 'к': "k", 'л': "l", 'м': "m",
	'н': "n", 'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u",
	'ф': "f", 'х': "h", 'ц': "c", 'ч': "ch", 'ш': "sh", 'щ': "shch", 'ъ': "",
	'ы': "y", 'ь': "", 'э': "e", 'ю': "yu", 'я': "ya",

	// Russian Cyrillic (uppercase)
	'А': "a", 'Б': "b", 'В': "v", 'Г': "g", 'Д': "d", 'Е': "e", 'Ё': "yo",
	'Ж': "zh", 'З': "z", 'И': "i", 'Й': "j", 'К': "k", 'Л': "l", 'М': "m",
	'Н': "n", 'О': "o", 'П': "p", 'Р': "r", 'С': "s", 'Т': "t", 'У': "u",
	'Ф': "f", 'Х': "h", 'Ц': "c", 'Ч': "ch", 'Ш': "sh", 'Щ': "shch", 'Ъ': "",
	'Ы': "y", 'Ь': "", 'Э': "e", 'Ю': "yu", 'Я': "ya",

	// Serbian Cyrillic specific (lowercase)
	'ђ': "dj", 'ј': "j", 'љ': "lj", 'њ': "nj", 'ћ': "c", 'џ': "dz",

	// Serbian Cyrillic specific (uppercase)
	'Ђ': "dj", 'Ј': "j", 'Љ': "lj", 'Њ': "nj", 'Ћ': "c", 'Џ': "dz",
}

// transliterate converts Cyrillic characters to Latin equivalents
func transliterate(s string) string {
	var result strings.Builder
	result.Grow(len(s) * 2) // Pre-allocate approximately needed capacity

	for _, r := range s {
		if replacement, ok := transliterationMap[r]; ok {
			result.WriteString(replacement)
		} else {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// slugify converts a string to a URL-friendly slug
// Rules:
// - Transliterate Cyrillic characters to Latin
// - Convert to lowercase
// - Replace spaces with hyphens
// - Remove all non-alphanumeric characters except hyphens
// - Remove multiple consecutive hyphens
// - Trim hyphens from start and end
func slugify(s string) string {
	if s == "" {
		return ""
	}

	// Transliterate Cyrillic to Latin FIRST
	s = transliterate(s)

	// Convert to lowercase
	s = strings.ToLower(s)

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

// ValidateSlug checks if a slug is valid according to our rules
func ValidateSlug(slug string) error {
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
