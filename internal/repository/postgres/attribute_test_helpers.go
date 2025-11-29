package postgres

import "github.com/vondi-global/listings/internal/domain"

// Test helper functions for attribute repository tests (attribute-specific only)

func boolPtr(v bool) *bool {
	return &v
}

func attrTypePtr(t domain.AttributeType) *domain.AttributeType {
	return &t
}

func attrPurposePtr(p domain.AttributePurpose) *domain.AttributePurpose {
	return &p
}

// Note: stringPtr, int32Ptr, float64Ptr, setupTestCategories are defined in repository_test.go
