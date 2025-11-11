// Package domain defines core business entities and domain models for the listings microservice.
// This file contains storefront-related domain models for B2C storefronts.
package domain

import (
	"time"
)

// Storefront represents a B2C storefront (business store) entity
type Storefront struct {
	ID          int64   `json:"id" db:"id"`
	UserID      int64   `json:"user_id" db:"user_id"`
	Slug        string  `json:"slug" db:"slug"`
	Name        string  `json:"name" db:"name"`
	Description *string `json:"description,omitempty" db:"description"`

	// Branding
	LogoURL   *string `json:"logo_url,omitempty" db:"logo_url"`
	BannerURL *string `json:"banner_url,omitempty" db:"banner_url"`

	// Contact information
	Phone   *string `json:"phone,omitempty" db:"phone"`
	Email   *string `json:"email,omitempty" db:"email"`
	Website *string `json:"website,omitempty" db:"website"`

	// Address & location
	Address    *string  `json:"address,omitempty" db:"address"`
	City       *string  `json:"city,omitempty" db:"city"`
	PostalCode *string  `json:"postal_code,omitempty" db:"postal_code"`
	Country    string   `json:"country" db:"country"`
	Latitude   *float64 `json:"latitude,omitempty" db:"latitude"`
	Longitude  *float64 `json:"longitude,omitempty" db:"longitude"`

	// Status flags
	IsActive   bool `json:"is_active" db:"is_active"`
	IsVerified bool `json:"is_verified" db:"is_verified"`

	// Statistics
	Rating         float64 `json:"rating" db:"rating"`
	ReviewsCount   int32   `json:"reviews_count" db:"reviews_count"`
	ProductsCount  int32   `json:"products_count" db:"products_count"`
	SalesCount     int32   `json:"sales_count" db:"sales_count"`
	ViewsCount     int32   `json:"views_count" db:"views_count"`
	FollowersCount int32   `json:"followers_count" db:"followers_count"`

	// Timestamps
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// IsDeleted returns true if the storefront is soft-deleted
func (s *Storefront) IsDeleted() bool {
	return s.DeletedAt != nil
}

// HasLocation returns true if the storefront has coordinates set
func (s *Storefront) HasLocation() bool {
	return s.Latitude != nil && s.Longitude != nil
}

// GetDisplayName returns the name for display purposes
func (s *Storefront) GetDisplayName() string {
	return s.Name
}

// GetURL returns the storefront URL based on slug
func (s *Storefront) GetURL(baseURL string) string {
	return baseURL + "/store/" + s.Slug
}
