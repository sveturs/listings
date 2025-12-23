// Package service defines business logic interfaces for the listings microservice.
package service

import (
	"context"

	"github.com/vondi-global/listings/internal/domain"
)

// AttributeService defines the interface for attribute business logic operations
type AttributeService interface {
	// Admin CRUD Operations
	// CreateAttribute creates a new attribute with comprehensive validation
	CreateAttribute(ctx context.Context, input *domain.CreateAttributeInput) (*domain.Attribute, error)

	// UpdateAttribute updates an existing attribute with validation
	UpdateAttribute(ctx context.Context, id int32, input *domain.UpdateAttributeInput) (*domain.Attribute, error)

	// DeleteAttribute soft-deletes an attribute (sets is_active=false)
	DeleteAttribute(ctx context.Context, id int32) error

	// GetAttribute retrieves an attribute by ID or code
	GetAttribute(ctx context.Context, identifier string) (*domain.Attribute, error)

	// GetAttributeByID retrieves an attribute by ID (cached)
	GetAttributeByID(ctx context.Context, id int32) (*domain.Attribute, error)

	// GetAttributeByCode retrieves an attribute by code (cached)
	GetAttributeByCode(ctx context.Context, code string) (*domain.Attribute, error)

	// ListAttributes lists attributes with filters and pagination
	ListAttributes(ctx context.Context, filter *domain.ListAttributesFilter) ([]*domain.Attribute, int64, error)

	// Category Linking Operations
	// LinkAttributeToCategory links an attribute to a category with optional overrides
	LinkAttributeToCategory(ctx context.Context, categoryID string, attributeID int32, settings *domain.CategoryAttributeSettings) error

	// UpdateCategoryAttribute updates category-specific attribute settings
	UpdateCategoryAttribute(ctx context.Context, catAttrID int32, settings *domain.CategoryAttributeSettings) error

	// UnlinkAttributeFromCategory removes attribute-category association
	UnlinkAttributeFromCategory(ctx context.Context, categoryID string, attributeID int32) error

	// GetCategoryAttributes retrieves all attributes for a category with effective values (cached)
	GetCategoryAttributes(ctx context.Context, categoryID string, filter *domain.GetCategoryAttributesFilter) ([]*domain.CategoryAttribute, error)

	// Listing Operations
	// GetListingAttributes retrieves all attribute values for a listing (cached)
	GetListingAttributes(ctx context.Context, listingID int32) ([]*domain.ListingAttributeValue, error)

	// SetListingAttributes validates and sets attribute values for a listing
	// This replaces all existing values with the provided ones
	SetListingAttributes(ctx context.Context, listingID int32, values []domain.SetListingAttributeValue) error

	// ValidateAttributeValues validates attribute values according to type and rules
	ValidateAttributeValues(ctx context.Context, categoryID string, values []domain.SetListingAttributeValue) error

	// Variant Attributes Operations
	// GetCategoryVariantAttributes retrieves variant attributes for a category (cached)
	GetCategoryVariantAttributes(ctx context.Context, categoryID string) ([]*domain.VariantAttribute, error)

	// Cache Management Operations
	// InvalidateAttributeCache invalidates cache for an attribute
	InvalidateAttributeCache(ctx context.Context, attributeID int32) error

	// InvalidateCategoryCache invalidates all attribute caches for a category
	InvalidateCategoryCache(ctx context.Context, categoryID string) error

	// InvalidateListingCache invalidates attribute cache for a listing
	InvalidateListingCache(ctx context.Context, listingID int32) error
}
