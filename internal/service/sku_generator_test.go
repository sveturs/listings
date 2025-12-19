// Package service contains business logic tests for the listings microservice.
// This file tests SKUGenerator functionality.
package service

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestSKUGenerator_Clothing tests SKU generation for clothing category
func TestSKUGenerator_Clothing(t *testing.T) {
	gen := NewSKUGenerator()

	productID := uuid.New().String()
	categoryCode := "CLO" // Clothing

	attrs := []VariantAttributeForSKU{
		{Code: "size", ValueLabel: "M"},
		{Code: "color", ValueLabel: "Black"},
	}

	sku := gen.GenerateSKU(productID, categoryCode, attrs)

	// Assert SKU format: CLO-xxxxxx-M-BLK
	assert.True(t, strings.HasPrefix(sku, "CLO-"))
	assert.Contains(t, sku, "-M-")
	assert.Contains(t, sku, "-BLK")

	parts := strings.Split(sku, "-")
	assert.GreaterOrEqual(t, len(parts), 4)
}

// TestSKUGenerator_Electronics tests SKU generation for electronics category
func TestSKUGenerator_Electronics(t *testing.T) {
	gen := NewSKUGenerator()

	productID := uuid.New().String()
	categoryCode := "ELE" // Electronics

	attrs := []VariantAttributeForSKU{
		{Code: "storage", ValueLabel: "256GB"},
		{Code: "color", ValueLabel: "Black"},
	}

	sku := gen.GenerateSKU(productID, categoryCode, attrs)

	// Assert SKU format: ELE-xxxxxx-256-BLK
	assert.True(t, strings.HasPrefix(sku, "ELE-"))
	assert.Contains(t, sku, "-256")
	assert.Contains(t, sku, "-BLK")
}

// TestSKUGenerator_Uniqueness tests that different products get different SKUs
func TestSKUGenerator_Uniqueness(t *testing.T) {
	gen := NewSKUGenerator()

	productID1 := uuid.New().String()
	productID2 := uuid.New().String()
	categoryCode := "CLO"

	attrs := []VariantAttributeForSKU{
		{Code: "size", ValueLabel: "M"},
	}

	sku1 := gen.GenerateSKU(productID1, categoryCode, attrs)
	sku2 := gen.GenerateSKU(productID2, categoryCode, attrs)

	// Different products should have different SKUs
	assert.NotEqual(t, sku1, sku2)
}

// TestSKUGenerator_SameProductSameAttrs tests that same product+attrs gets same SKU
func TestSKUGenerator_SameProductSameAttrs(t *testing.T) {
	gen := NewSKUGenerator()

	productID := uuid.New().String()
	categoryCode := "CLO"

	attrs := []VariantAttributeForSKU{
		{Code: "size", ValueLabel: "M"},
		{Code: "color", ValueLabel: "Red"},
	}

	sku1 := gen.GenerateSKU(productID, categoryCode, attrs)
	sku2 := gen.GenerateSKU(productID, categoryCode, attrs)

	// Same inputs should generate same SKU (deterministic)
	assert.Equal(t, sku1, sku2)
}

// TestSKUGenerator_Validation tests SKU validation
func TestSKUGenerator_Validation(t *testing.T) {
	gen := NewSKUGenerator()

	tests := []struct {
		name      string
		sku       string
		wantError bool
	}{
		{
			name:      "valid SKU",
			sku:       "CLO-ABC123-M-BLK",
			wantError: false,
		},
		{
			name:      "empty SKU",
			sku:       "",
			wantError: true,
		},
		{
			name:      "too short",
			sku:       "AB",
			wantError: true,
		},
		{
			name:      "too long",
			sku:       strings.Repeat("A", 101),
			wantError: true,
		},
		{
			name:      "invalid characters",
			sku:       "SKU@#$%",
			wantError: true,
		},
		{
			name:      "valid with numbers",
			sku:       "ELE-123456-256GB",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := gen.ValidateSKU(tt.sku)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestSKUGenerator_DifferentAttrsOrder tests that attribute order doesn't affect SKU
func TestSKUGenerator_DifferentAttrsOrder(t *testing.T) {
	gen := NewSKUGenerator()

	productID := uuid.New().String()
	categoryCode := "CLO"

	attrs1 := []VariantAttributeForSKU{
		{Code: "size", ValueLabel: "M"},
		{Code: "color", ValueLabel: "Red"},
	}

	attrs2 := []VariantAttributeForSKU{
		{Code: "color", ValueLabel: "Red"},
		{Code: "size", ValueLabel: "M"},
	}

	sku1 := gen.GenerateSKU(productID, categoryCode, attrs1)
	sku2 := gen.GenerateSKU(productID, categoryCode, attrs2)

	// SKU should be the same regardless of attribute order (sorted internally)
	assert.Equal(t, sku1, sku2)
}

// TestSKUGenerator_ColorAbbreviation tests color name abbreviation
func TestSKUGenerator_ColorAbbreviation(t *testing.T) {
	gen := NewSKUGenerator()

	productID := uuid.New().String()
	categoryCode := "CLO"

	tests := []struct {
		color    string
		expected string
	}{
		{"Black", "BLK"},
		{"White", "WHT"},
		{"Red", "RED"},
		{"Blue", "BLU"},
		{"Green", "GRN"},
		{"Yellow", "YEL"},
		{"Gray", "GRY"},
		{"Grey", "GRY"},
		{"Unknown Color", "UNK"},
	}

	for _, tt := range tests {
		t.Run(tt.color, func(t *testing.T) {
			attrs := []VariantAttributeForSKU{
				{Code: "color", ValueLabel: tt.color},
			}

			sku := gen.GenerateSKU(productID, categoryCode, attrs)
			assert.Contains(t, sku, tt.expected)
		})
	}
}

// TestSKUGenerator_SizeAbbreviation tests size abbreviation
func TestSKUGenerator_SizeAbbreviation(t *testing.T) {
	gen := NewSKUGenerator()

	productID := uuid.New().String()
	categoryCode := "CLO"

	tests := []struct {
		size     string
		expected string
	}{
		{"Small", "S"},
		{"Medium", "M"},
		{"Large", "L"},
		{"Extra Large", "XL"},
		{"Extra Extra Large", "XXL"},
		{"M", "M"},
		{"42", "42"},
	}

	for _, tt := range tests {
		t.Run(tt.size, func(t *testing.T) {
			attrs := []VariantAttributeForSKU{
				{Code: "size", ValueLabel: tt.size},
			}

			sku := gen.GenerateSKU(productID, categoryCode, attrs)
			assert.Contains(t, sku, tt.expected)
		})
	}
}
