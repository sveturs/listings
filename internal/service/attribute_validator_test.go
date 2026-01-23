// Package service implements business logic for the listings microservice.
package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/vondi-global/listings/internal/domain"
)

func TestAttributeValidator_ValidateText_Success(t *testing.T) {
	v := NewAttributeValidator()

	tests := []struct {
		name  string
		value string
		rules map[string]interface{}
	}{
		{
			name:  "simple text",
			value: "Hello World",
			rules: map[string]interface{}{},
		},
		{
			name:  "text with min length",
			value: "Hello",
			rules: map[string]interface{}{"min_length": 3},
		},
		{
			name:  "text with max length",
			value: "Hello",
			rules: map[string]interface{}{"max_length": 10},
		},
		{
			name:  "text with pattern",
			value: "HelloWorld",
			rules: map[string]interface{}{"pattern": "^[A-Za-z]+$"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.validateText(tt.value, tt.rules)
			assert.NoError(t, err)
		})
	}
}

func TestAttributeValidator_ValidateText_Failures(t *testing.T) {
	v := NewAttributeValidator()

	tests := []struct {
		name        string
		value       interface{}
		rules       map[string]interface{}
		expectedErr string
	}{
		{
			name:        "not a string",
			value:       123,
			rules:       map[string]interface{}{},
			expectedErr: "must be a string",
		},
		{
			name:        "too short",
			value:       "Hi",
			rules:       map[string]interface{}{"min_length": 5},
			expectedErr: "less than minimum",
		},
		{
			name:        "too long",
			value:       "Hello World",
			rules:       map[string]interface{}{"max_length": 5},
			expectedErr: "exceeds maximum",
		},
		{
			name:        "pattern mismatch",
			value:       "Hello123",
			rules:       map[string]interface{}{"pattern": "^[A-Za-z]+$"},
			expectedErr: "does not match",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.validateText(tt.value, tt.rules)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestAttributeValidator_ValidateNumber_Success(t *testing.T) {
	v := NewAttributeValidator()

	tests := []struct {
		name  string
		value interface{}
		rules map[string]interface{}
	}{
		{
			name:  "integer",
			value: 42,
			rules: map[string]interface{}{},
		},
		{
			name:  "float64",
			value: 42.5,
			rules: map[string]interface{}{},
		},
		{
			name:  "with min",
			value: 50.0,
			rules: map[string]interface{}{"min": 10.0},
		},
		{
			name:  "with max",
			value: 50.0,
			rules: map[string]interface{}{"max": 100.0},
		},
		{
			name:  "with decimals",
			value: 42.5,
			rules: map[string]interface{}{"decimals": 2},
		},
		{
			name:  "string number",
			value: "42.5",
			rules: map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.validateNumber(tt.value, tt.rules)
			assert.NoError(t, err)
		})
	}
}

func TestAttributeValidator_ValidateNumber_Failures(t *testing.T) {
	v := NewAttributeValidator()

	tests := []struct {
		name        string
		value       interface{}
		rules       map[string]interface{}
		expectedErr string
	}{
		{
			name:        "not numeric",
			value:       "not a number",
			rules:       map[string]interface{}{},
			expectedErr: "must be numeric",
		},
		{
			name:        "below min",
			value:       5.0,
			rules:       map[string]interface{}{"min": 10.0},
			expectedErr: "less than minimum",
		},
		{
			name:        "above max",
			value:       150.0,
			rules:       map[string]interface{}{"max": 100.0},
			expectedErr: "exceeds maximum",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.validateNumber(tt.value, tt.rules)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestAttributeValidator_ValidateBoolean_Success(t *testing.T) {
	v := NewAttributeValidator()

	tests := []struct {
		name  string
		value interface{}
	}{
		{"bool true", true},
		{"bool false", false},
		{"string true", "true"},
		{"string false", "false"},
		{"string TRUE", "TRUE"},
		{"string FALSE", "FALSE"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.validateBoolean(tt.value)
			assert.NoError(t, err)
		})
	}
}

func TestAttributeValidator_ValidateBoolean_Failures(t *testing.T) {
	v := NewAttributeValidator()

	tests := []struct {
		name  string
		value interface{}
	}{
		{"string yes", "yes"},
		{"string no", "no"},
		{"number", 1},
		{"string 1", "1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.validateBoolean(tt.value)
			assert.Error(t, err)
		})
	}
}

func TestAttributeValidator_ValidateSelect_Success(t *testing.T) {
	v := NewAttributeValidator()

	options := []domain.AttributeOption{
		{Value: "S", Label: map[string]string{"en": "Small"}},
		{Value: "M", Label: map[string]string{"en": "Medium"}},
		{Value: "L", Label: map[string]string{"en": "Large"}},
	}

	tests := []struct {
		name  string
		value string
	}{
		{"select S", "S"},
		{"select M", "M"},
		{"select L", "L"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.validateSelect(tt.value, options)
			assert.NoError(t, err)
		})
	}
}

func TestAttributeValidator_ValidateSelect_Failures(t *testing.T) {
	v := NewAttributeValidator()

	options := []domain.AttributeOption{
		{Value: "S", Label: map[string]string{"en": "Small"}},
		{Value: "M", Label: map[string]string{"en": "Medium"}},
	}

	tests := []struct {
		name        string
		value       interface{}
		expectedErr string
	}{
		{
			name:        "invalid value",
			value:       "XL",
			expectedErr: "invalid select value",
		},
		{
			name:        "not a string",
			value:       123,
			expectedErr: "must be a string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.validateSelect(tt.value, options)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestAttributeValidator_ValidateSelect_NoOptions(t *testing.T) {
	v := NewAttributeValidator()

	err := v.validateSelect("S", []domain.AttributeOption{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no valid options")
}

func TestAttributeValidator_ValidateMultiselect_Success(t *testing.T) {
	v := NewAttributeValidator()

	options := []domain.AttributeOption{
		{Value: "red", Label: map[string]string{"en": "Red"}},
		{Value: "blue", Label: map[string]string{"en": "Blue"}},
		{Value: "green", Label: map[string]string{"en": "Green"}},
	}

	tests := []struct {
		name  string
		value interface{}
	}{
		{
			name:  "array of strings",
			value: []string{"red", "blue"},
		},
		{
			name:  "array of interface{}",
			value: []interface{}{"red", "green"},
		},
		{
			name:  "single string",
			value: "red",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.validateMultiselect(tt.value, options)
			assert.NoError(t, err)
		})
	}
}

func TestAttributeValidator_ValidateMultiselect_Failures(t *testing.T) {
	v := NewAttributeValidator()

	options := []domain.AttributeOption{
		{Value: "red", Label: map[string]string{"en": "Red"}},
		{Value: "blue", Label: map[string]string{"en": "Blue"}},
	}

	tests := []struct {
		name        string
		value       interface{}
		expectedErr string
	}{
		{
			name:        "invalid value in array",
			value:       []string{"red", "yellow"},
			expectedErr: "invalid multiselect value 'yellow'",
		},
		{
			name:        "empty array",
			value:       []string{},
			expectedErr: "cannot be empty",
		},
		{
			name:        "not array",
			value:       123,
			expectedErr: "must be array",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.validateMultiselect(tt.value, options)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestAttributeValidator_ValidateDate_Success(t *testing.T) {
	v := NewAttributeValidator()

	tests := []struct {
		name  string
		value interface{}
		rules map[string]interface{}
	}{
		{
			name:  "time.Time",
			value: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			rules: map[string]interface{}{},
		},
		{
			name:  "string YYYY-MM-DD",
			value: "2025-01-15",
			rules: map[string]interface{}{},
		},
		{
			name:  "string RFC3339",
			value: "2025-01-15T10:30:00Z",
			rules: map[string]interface{}{},
		},
		{
			name:  "with min_date",
			value: "2025-01-15",
			rules: map[string]interface{}{"min_date": "2025-01-01"},
		},
		{
			name:  "with max_date",
			value: "2025-01-15",
			rules: map[string]interface{}{"max_date": "2025-12-31"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.validateDate(tt.value, tt.rules)
			assert.NoError(t, err)
		})
	}
}

func TestAttributeValidator_ValidateDate_Failures(t *testing.T) {
	v := NewAttributeValidator()

	tests := []struct {
		name        string
		value       interface{}
		rules       map[string]interface{}
		expectedErr string
	}{
		{
			name:        "invalid format",
			value:       "15/01/2025",
			rules:       map[string]interface{}{},
			expectedErr: "invalid date format",
		},
		{
			name:        "before min_date",
			value:       "2024-12-01",
			rules:       map[string]interface{}{"min_date": "2025-01-01"},
			expectedErr: "before minimum date",
		},
		{
			name:        "after max_date",
			value:       "2026-01-01",
			rules:       map[string]interface{}{"max_date": "2025-12-31"},
			expectedErr: "after maximum date",
		},
		{
			name:        "not a date type",
			value:       123,
			rules:       map[string]interface{}{},
			expectedErr: "must be time.Time or string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.validateDate(tt.value, tt.rules)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestAttributeValidator_ValidateColor_Success(t *testing.T) {
	v := NewAttributeValidator()

	tests := []struct {
		name  string
		value string
	}{
		{"6-digit hex", "#FF5733"},
		{"3-digit hex", "#F53"},
		{"lowercase hex", "#ff5733"},
		{"uppercase hex", "#FF5733"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.validateColor(tt.value)
			assert.NoError(t, err)
		})
	}
}

func TestAttributeValidator_ValidateColor_Failures(t *testing.T) {
	v := NewAttributeValidator()

	tests := []struct {
		name  string
		value interface{}
	}{
		{"no hash", "FF5733"},
		{"invalid characters", "#GG5733"},
		{"too short", "#FF"},
		{"too long", "#FF57333"},
		{"not a string", 123},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.validateColor(tt.value)
			assert.Error(t, err)
		})
	}
}

func TestAttributeValidator_ValidateValue_Integration(t *testing.T) {
	v := NewAttributeValidator()

	tests := []struct {
		name      string
		attrType  domain.AttributeType
		value     interface{}
		rules     map[string]interface{}
		options   []domain.AttributeOption
		expectErr bool
	}{
		{
			name:      "text valid",
			attrType:  domain.AttributeTypeText,
			value:     "Hello World",
			rules:     map[string]interface{}{"max_length": 50},
			options:   nil,
			expectErr: false,
		},
		{
			name:      "number valid",
			attrType:  domain.AttributeTypeNumber,
			value:     42.5,
			rules:     map[string]interface{}{"min": 0, "max": 100},
			options:   nil,
			expectErr: false,
		},
		{
			name:      "select valid",
			attrType:  domain.AttributeTypeSelect,
			value:     "M",
			rules:     nil,
			options:   []domain.AttributeOption{{Value: "S"}, {Value: "M"}, {Value: "L"}},
			expectErr: false,
		},
		{
			name:      "select invalid value",
			attrType:  domain.AttributeTypeSelect,
			value:     "XL",
			rules:     nil,
			options:   []domain.AttributeOption{{Value: "S"}, {Value: "M"}, {Value: "L"}},
			expectErr: true,
		},
		{
			name:      "unsupported type",
			attrType:  "unsupported",
			value:     "test",
			rules:     nil,
			options:   nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateValue(tt.attrType, tt.value, tt.rules, tt.options)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
