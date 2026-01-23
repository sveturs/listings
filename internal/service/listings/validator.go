package listings

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/vondi-global/listings/internal/domain"
)

// Validation constants
const (
	// Listing validation constants
	MinTitleLength       = 3
	MaxTitleLength       = 200
	MaxDescriptionLength = 5000
	MinPrice             = 0.0
	MaxPrice             = 999999999.99

	// Image validation constants
	MaxImageSize        = 10 * 1024 * 1024 // 10MB
	MaxImagesPerListing = 10
	MinImageWidth       = 100
	MinImageHeight      = 100
	MaxImageWidth       = 10000
	MaxImageHeight      = 10000
)

var (
	// Regular expressions for validation
	slugRegex     = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	currencyRegex = regexp.MustCompile(`^[A-Z]{3}$`) // ISO 4217
)

// ValidationError represents a field validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// MultiValidationError contains multiple validation errors
type MultiValidationError struct {
	Errors []ValidationError
}

func (e *MultiValidationError) Error() string {
	if len(e.Errors) == 0 {
		return "validation failed"
	}
	return fmt.Sprintf("validation failed with %d error(s): %s", len(e.Errors), e.Errors[0].Error())
}

// Add adds a validation error
func (e *MultiValidationError) Add(field, message string) {
	e.Errors = append(e.Errors, ValidationError{Field: field, Message: message})
}

// HasErrors checks if there are any validation errors
func (e *MultiValidationError) HasErrors() bool {
	return len(e.Errors) > 0
}

// Validator provides validation methods for listings business logic
type Validator struct {
	repo Repository
}

// NewValidator creates a new validator instance
func NewValidator(repo Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCategory checks if category exists and is active
func (v *Validator) ValidateCategory(ctx context.Context, categoryID string) error {
	if categoryID == "" {
		return NewValidationError("category_id", "category_id must not be empty")
	}

	category, err := v.repo.GetCategoryByID(ctx, categoryID)
	if err != nil {
		return fmt.Errorf("failed to get category: %w", err)
	}

	if category == nil {
		return fmt.Errorf("category not found")
	}

	if !category.IsActive {
		return fmt.Errorf("category is inactive")
	}

	return nil
}

// ValidatePrice ensures price is positive and within range
func (v *Validator) ValidatePrice(price float64) error {
	if price < MinPrice {
		return NewValidationError("price", "price must be non-negative")
	}

	if price > MaxPrice {
		return NewValidationError("price", "price exceeds maximum value")
	}

	return nil
}

// ValidateTitle checks title length constraints
func (v *Validator) ValidateTitle(title string) error {
	trimmed := strings.TrimSpace(title)

	if trimmed == "" {
		return NewValidationError("title", "title is required")
	}

	length := utf8.RuneCountInString(trimmed)
	if length < MinTitleLength {
		return NewValidationError("title", fmt.Sprintf("title must be at least %d characters", MinTitleLength))
	}

	if length > MaxTitleLength {
		return NewValidationError("title", fmt.Sprintf("title exceeds maximum length of %d characters", MaxTitleLength))
	}

	return nil
}

// ValidateDescription checks description length (optional field)
func (v *Validator) ValidateDescription(description *string) error {
	if description == nil {
		return nil // Optional field
	}

	length := utf8.RuneCountInString(*description)
	if length > MaxDescriptionLength {
		return NewValidationError("description", fmt.Sprintf("description exceeds maximum length of %d characters", MaxDescriptionLength))
	}

	return nil
}

// ValidateQuantity ensures quantity is non-negative
func (v *Validator) ValidateQuantity(quantity int32) error {
	if quantity < 0 {
		return NewValidationError("quantity", "quantity must be non-negative")
	}

	return nil
}

// ValidateCurrency validates ISO 4217 currency code
func (v *Validator) ValidateCurrency(currency string) error {
	if currency == "" {
		return NewValidationError("currency", "currency is required")
	}

	if !currencyRegex.MatchString(currency) {
		return NewValidationError("currency", "currency must be a valid ISO 4217 code (e.g., RSD, EUR, USD)")
	}

	return nil
}

// ValidateImages checks image count, size, and format
func (v *Validator) ValidateImages(images []*domain.ListingImage) error {
	if len(images) > MaxImagesPerListing {
		return NewValidationError("images", fmt.Sprintf("maximum %d images per listing", MaxImagesPerListing))
	}

	for i, img := range images {
		// Validate file size
		if img.FileSize != nil && *img.FileSize > MaxImageSize {
			return NewValidationError(fmt.Sprintf("images[%d]", i), "image size exceeds maximum of 10MB")
		}

		// Validate dimensions
		if img.Width != nil && img.Height != nil {
			if *img.Width < MinImageWidth || *img.Height < MinImageHeight {
				return NewValidationError(fmt.Sprintf("images[%d]", i), "image dimensions are below minimum (100x100)")
			}

			if *img.Width > MaxImageWidth || *img.Height > MaxImageHeight {
				return NewValidationError(fmt.Sprintf("images[%d]", i), "image dimensions exceed maximum (10000x10000)")
			}
		}

		// Validate MIME type
		if img.MimeType != nil {
			validTypes := map[string]bool{
				"image/jpeg": true,
				"image/jpg":  true,
				"image/png":  true,
				"image/webp": true,
			}
			if !validTypes[*img.MimeType] {
				return NewValidationError(fmt.Sprintf("images[%d]", i), "unsupported image format (allowed: JPEG, PNG, WebP)")
			}
		}
	}

	return nil
}

// ValidateStatusTransition checks if a status transition is allowed
func (v *Validator) ValidateStatusTransition(from, to string) error {
	// Status transition rules
	allowedTransitions := map[string][]string{
		domain.StatusDraft: {
			domain.StatusActive,
			domain.StatusInactive,
		},
		domain.StatusActive: {
			domain.StatusSold,
			domain.StatusInactive,
			domain.StatusArchived,
		},
		domain.StatusInactive: {
			domain.StatusActive,
			domain.StatusDraft,
		},
		domain.StatusSold: {
			domain.StatusActive, // Allow re-listing
		},
	}

	allowed := allowedTransitions[from]
	for _, status := range allowed {
		if status == to {
			return nil
		}
	}

	return fmt.Errorf("invalid status transition from '%s' to '%s'", from, to)
}

// ValidateCreateInput validates all fields for creating a listing
func (v *Validator) ValidateCreateInput(ctx context.Context, input *domain.CreateListingInput) error {
	errs := &MultiValidationError{}

	// User ID
	if input.UserID <= 0 {
		errs.Add("user_id", "user_id is required and must be greater than 0")
	}

	// Title
	if err := v.ValidateTitle(input.Title); err != nil {
		if ve, ok := err.(*ValidationError); ok {
			errs.Add(ve.Field, ve.Message)
		}
	}

	// Description
	if err := v.ValidateDescription(input.Description); err != nil {
		if ve, ok := err.(*ValidationError); ok {
			errs.Add(ve.Field, ve.Message)
		}
	}

	// Price
	if err := v.ValidatePrice(input.Price); err != nil {
		if ve, ok := err.(*ValidationError); ok {
			errs.Add(ve.Field, ve.Message)
		}
	}

	// Currency
	if err := v.ValidateCurrency(input.Currency); err != nil {
		if ve, ok := err.(*ValidationError); ok {
			errs.Add(ve.Field, ve.Message)
		}
	}

	// Category
	if err := v.ValidateCategory(ctx, input.CategoryID); err != nil {
		errs.Add("category_id", err.Error())
	}

	// Quantity
	if err := v.ValidateQuantity(input.Quantity); err != nil {
		if ve, ok := err.(*ValidationError); ok {
			errs.Add(ve.Field, ve.Message)
		}
	}

	// Source Type
	if input.SourceType != domain.SourceTypeC2C && input.SourceType != domain.SourceTypeB2C {
		errs.Add("source_type", "source_type must be 'c2c' or 'b2c'")
	}

	// B2C specific validation
	if input.SourceType == domain.SourceTypeB2C {
		if input.StorefrontID == nil || *input.StorefrontID <= 0 {
			errs.Add("storefront_id", "storefront_id is required for B2C listings")
		}
	}

	if errs.HasErrors() {
		return errs
	}

	return nil
}

// ValidateUpdateInput validates fields for updating a listing
func (v *Validator) ValidateUpdateInput(input *domain.UpdateListingInput) error {
	errs := &MultiValidationError{}

	// Title (if being updated)
	if input.Title != nil {
		if err := v.ValidateTitle(*input.Title); err != nil {
			if ve, ok := err.(*ValidationError); ok {
				errs.Add(ve.Field, ve.Message)
			}
		}
	}

	// Description (if being updated)
	if input.Description != nil {
		if err := v.ValidateDescription(input.Description); err != nil {
			if ve, ok := err.(*ValidationError); ok {
				errs.Add(ve.Field, ve.Message)
			}
		}
	}

	// Price (if being updated)
	if input.Price != nil {
		if err := v.ValidatePrice(*input.Price); err != nil {
			if ve, ok := err.(*ValidationError); ok {
				errs.Add(ve.Field, ve.Message)
			}
		}
	}

	// Quantity (if being updated)
	if input.Quantity != nil {
		if err := v.ValidateQuantity(*input.Quantity); err != nil {
			if ve, ok := err.(*ValidationError); ok {
				errs.Add(ve.Field, ve.Message)
			}
		}
	}

	// Status (if being updated)
	if input.Status != nil {
		validStatuses := map[string]bool{
			domain.StatusDraft:    true,
			domain.StatusActive:   true,
			domain.StatusInactive: true,
			domain.StatusSold:     true,
			domain.StatusArchived: true,
		}
		if !validStatuses[*input.Status] {
			errs.Add("status", "invalid status value")
		}
	}

	if errs.HasErrors() {
		return errs
	}

	return nil
}

// Helper functions

// ValidateSlug validates slug format
func ValidateSlug(slug string) bool {
	if slug == "" {
		return false
	}
	if utf8.RuneCountInString(slug) > 250 {
		return false
	}
	return slugRegex.MatchString(slug)
}
