// Package repository defines data access interfaces for the listings microservice.
package repository

import (
	"context"

	"github.com/sveturs/listings/internal/domain"
)

// AttributeRepository defines the interface for attribute data access operations
type AttributeRepository interface {
	// CRUD Operations
	Create(ctx context.Context, attr *domain.CreateAttributeInput) (*domain.Attribute, error)
	Update(ctx context.Context, id int32, attr *domain.UpdateAttributeInput) (*domain.Attribute, error)
	Delete(ctx context.Context, id int32) error // Soft delete (set is_active=false)
	GetByID(ctx context.Context, id int32) (*domain.Attribute, error)
	GetByCode(ctx context.Context, code string) (*domain.Attribute, error)
	List(ctx context.Context, filter *domain.ListAttributesFilter) ([]*domain.Attribute, int64, error)

	// Category Linking
	LinkToCategory(ctx context.Context, categoryID int32, attributeID int32, settings *domain.CategoryAttributeSettings) (*domain.CategoryAttribute, error)
	UpdateCategoryAttribute(ctx context.Context, catAttrID int32, settings *domain.CategoryAttributeSettings) (*domain.CategoryAttribute, error)
	UnlinkFromCategory(ctx context.Context, categoryID int32, attributeID int32) error
	GetCategoryAttributes(ctx context.Context, categoryID int32, filter *domain.GetCategoryAttributesFilter) ([]*domain.CategoryAttribute, error)

	// Listing Values
	GetListingValues(ctx context.Context, listingID int32) ([]*domain.ListingAttributeValue, error)
	SetListingValues(ctx context.Context, listingID int32, values []domain.SetListingAttributeValue) error
	DeleteListingValues(ctx context.Context, listingID int32) error

	// Variant Attributes
	GetCategoryVariantAttributes(ctx context.Context, categoryID int32) ([]*domain.VariantAttribute, error)
	GetVariantValues(ctx context.Context, variantID int32) ([]*domain.VariantAttributeValue, error)
}
