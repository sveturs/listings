package indexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAttributeForIndex_Structure(t *testing.T) {
	// Test that AttributeForIndex has correct fields
	attr := AttributeForIndex{
		ID:           1,
		Code:         "test_code",
		Name:         "Test Attribute",
		IsSearchable: true,
		IsFilterable: true,
	}

	assert.Equal(t, int32(1), attr.ID)
	assert.Equal(t, "test_code", attr.Code)
	assert.Equal(t, "Test Attribute", attr.Name)
	assert.True(t, attr.IsSearchable)
	assert.True(t, attr.IsFilterable)
}

func TestAttributeForIndex_WithValues(t *testing.T) {
	// Test text value
	textValue := "Test Value"
	attr := AttributeForIndex{
		ID:           1,
		Code:         "brand",
		Name:         "Brand",
		ValueText:    &textValue,
		IsSearchable: true,
		IsFilterable: true,
	}
	assert.NotNil(t, attr.ValueText)
	assert.Equal(t, "Test Value", *attr.ValueText)

	// Test number value
	numberValue := 123.45
	attr2 := AttributeForIndex{
		ID:           2,
		Code:         "price",
		Name:         "Price",
		ValueNumber:  &numberValue,
		IsSearchable: false,
		IsFilterable: true,
	}
	assert.NotNil(t, attr2.ValueNumber)
	assert.Equal(t, 123.45, *attr2.ValueNumber)

	// Test boolean value
	boolValue := true
	attr3 := AttributeForIndex{
		ID:           3,
		Code:         "in_stock",
		Name:         "In Stock",
		ValueBoolean: &boolValue,
		IsSearchable: false,
		IsFilterable: true,
	}
	assert.NotNil(t, attr3.ValueBoolean)
	assert.True(t, *attr3.ValueBoolean)
}

// Note: Integration tests requiring database connection should be in separate _integration_test.go file
// These tests only verify data structures and simple logic
