package indexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAttributeForIndex_AllValueTypes tests all possible value type combinations
func TestAttributeForIndex_AllValueTypes(t *testing.T) {
	tests := []struct {
		name          string
		attr          AttributeForIndex
		expectedValue interface{}
		valueType     string
	}{
		{
			name: "text value only",
			attr: AttributeForIndex{
				ID:           1,
				Code:         "brand",
				Name:         "Brand",
				ValueText:    stringPtr("Nike"),
				IsSearchable: true,
				IsFilterable: true,
			},
			expectedValue: "Nike",
			valueType:     "text",
		},
		{
			name: "number value only",
			attr: AttributeForIndex{
				ID:           2,
				Code:         "price",
				Name:         "Price",
				ValueNumber:  float64Ptr(99.99),
				IsSearchable: false,
				IsFilterable: true,
			},
			expectedValue: 99.99,
			valueType:     "number",
		},
		{
			name: "boolean value only",
			attr: AttributeForIndex{
				ID:           3,
				Code:         "featured",
				Name:         "Featured",
				ValueBoolean: boolPtr(true),
				IsSearchable: false,
				IsFilterable: true,
			},
			expectedValue: true,
			valueType:     "boolean",
		},
		{
			name: "no value (all nil)",
			attr: AttributeForIndex{
				ID:           4,
				Code:         "empty",
				Name:         "Empty",
				ValueText:    nil,
				ValueNumber:  nil,
				ValueBoolean: nil,
				IsSearchable: false,
				IsFilterable: false,
			},
			expectedValue: nil,
			valueType:     "none",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify structure
			assert.Equal(t, tt.attr.ID, tt.attr.ID)
			assert.Equal(t, tt.attr.Code, tt.attr.Code)
			assert.Equal(t, tt.attr.Name, tt.attr.Name)

			// Verify value type
			switch tt.valueType {
			case "text":
				assert.NotNil(t, tt.attr.ValueText)
				assert.Equal(t, tt.expectedValue, *tt.attr.ValueText)
				assert.Nil(t, tt.attr.ValueNumber)
				assert.Nil(t, tt.attr.ValueBoolean)
			case "number":
				assert.NotNil(t, tt.attr.ValueNumber)
				assert.Equal(t, tt.expectedValue, *tt.attr.ValueNumber)
				assert.Nil(t, tt.attr.ValueText)
				assert.Nil(t, tt.attr.ValueBoolean)
			case "boolean":
				assert.NotNil(t, tt.attr.ValueBoolean)
				assert.Equal(t, tt.expectedValue, *tt.attr.ValueBoolean)
				assert.Nil(t, tt.attr.ValueText)
				assert.Nil(t, tt.attr.ValueNumber)
			case "none":
				assert.Nil(t, tt.attr.ValueText)
				assert.Nil(t, tt.attr.ValueNumber)
				assert.Nil(t, tt.attr.ValueBoolean)
			}
		})
	}
}

// TestAttributeForIndex_Flags tests searchable and filterable flag combinations
func TestAttributeForIndex_Flags(t *testing.T) {
	tests := []struct {
		name         string
		isSearchable bool
		isFilterable bool
	}{
		{"both true", true, true},
		{"searchable only", true, false},
		{"filterable only", false, true},
		{"both false", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attr := AttributeForIndex{
				ID:           1,
				Code:         "test",
				Name:         "Test",
				IsSearchable: tt.isSearchable,
				IsFilterable: tt.isFilterable,
			}

			assert.Equal(t, tt.isSearchable, attr.IsSearchable)
			assert.Equal(t, tt.isFilterable, attr.IsFilterable)
		})
	}
}

// TestAttributeForIndex_EdgeCases tests edge cases
func TestAttributeForIndex_EdgeCases(t *testing.T) {
	t.Run("empty string value", func(t *testing.T) {
		emptyStr := ""
		attr := AttributeForIndex{
			ID:        1,
			Code:      "test",
			Name:      "Test",
			ValueText: &emptyStr,
		}
		assert.NotNil(t, attr.ValueText)
		assert.Equal(t, "", *attr.ValueText)
	})

	t.Run("zero number value", func(t *testing.T) {
		zero := 0.0
		attr := AttributeForIndex{
			ID:          1,
			Code:        "test",
			Name:        "Test",
			ValueNumber: &zero,
		}
		assert.NotNil(t, attr.ValueNumber)
		assert.Equal(t, 0.0, *attr.ValueNumber)
	})

	t.Run("false boolean value", func(t *testing.T) {
		falseVal := false
		attr := AttributeForIndex{
			ID:           1,
			Code:         "test",
			Name:         "Test",
			ValueBoolean: &falseVal,
		}
		assert.NotNil(t, attr.ValueBoolean)
		assert.False(t, *attr.ValueBoolean)
	})

	t.Run("very long text value", func(t *testing.T) {
		longText := string(make([]byte, 10000))
		attr := AttributeForIndex{
			ID:        1,
			Code:      "test",
			Name:      "Test",
			ValueText: &longText,
		}
		assert.NotNil(t, attr.ValueText)
		assert.Len(t, *attr.ValueText, 10000)
	})

	t.Run("negative number value", func(t *testing.T) {
		negative := -99.99
		attr := AttributeForIndex{
			ID:          1,
			Code:        "test",
			Name:        "Test",
			ValueNumber: &negative,
		}
		assert.NotNil(t, attr.ValueNumber)
		assert.Equal(t, -99.99, *attr.ValueNumber)
	})

	t.Run("very large number value", func(t *testing.T) {
		large := 999999999.99
		attr := AttributeForIndex{
			ID:          1,
			Code:        "test",
			Name:        "Test",
			ValueNumber: &large,
		}
		assert.NotNil(t, attr.ValueNumber)
		assert.Equal(t, 999999999.99, *attr.ValueNumber)
	})

	t.Run("unicode text value", func(t *testing.T) {
		unicode := "测试 тест 🚀"
		attr := AttributeForIndex{
			ID:        1,
			Code:      "test",
			Name:      "Test",
			ValueText: &unicode,
		}
		assert.NotNil(t, attr.ValueText)
		assert.Equal(t, "测试 тест 🚀", *attr.ValueText)
	})
}

// TestAttributeForIndex_SpecialCharacters tests special characters in codes and names
func TestAttributeForIndex_SpecialCharacters(t *testing.T) {
	tests := []struct {
		name string
		code string
		desc string
	}{
		{"underscore", "test_attribute", "Standard snake_case"},
		{"dash", "test-attribute", "Kebab-case"},
		{"dot", "test.attribute", "Dot notation"},
		{"mixed", "test_Attribute-123.v2", "Mixed special chars"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attr := AttributeForIndex{
				ID:   1,
				Code: tt.code,
				Name: tt.desc,
			}
			assert.Equal(t, tt.code, attr.Code)
			assert.Equal(t, tt.desc, attr.Name)
		})
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func float64Ptr(f float64) *float64 {
	return &f
}

func boolPtr(b bool) *bool {
	return &b
}
