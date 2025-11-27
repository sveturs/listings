// Package service implements business logic for the listings microservice.
package service

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/vondi-global/listings/internal/domain"
)

// AttributeValidator provides validation logic for attribute values
type AttributeValidator struct{}

// NewAttributeValidator creates a new attribute validator
func NewAttributeValidator() *AttributeValidator {
	return &AttributeValidator{}
}

// ValidateValue validates an attribute value according to type and rules
func (v *AttributeValidator) ValidateValue(
	attrType domain.AttributeType,
	value interface{},
	rules map[string]interface{},
	options []domain.AttributeOption,
) error {
	switch attrType {
	case domain.AttributeTypeText, domain.AttributeTypeTextarea:
		return v.validateText(value, rules)
	case domain.AttributeTypeNumber:
		return v.validateNumber(value, rules)
	case domain.AttributeTypeBoolean:
		return v.validateBoolean(value)
	case domain.AttributeTypeSelect:
		return v.validateSelect(value, options)
	case domain.AttributeTypeMultiselect:
		return v.validateMultiselect(value, options)
	case domain.AttributeTypeDate:
		return v.validateDate(value, rules)
	case domain.AttributeTypeColor:
		return v.validateColor(value)
	case domain.AttributeTypeSize:
		return v.validateText(value, rules) // Size is treated as text
	default:
		return fmt.Errorf("unsupported attribute type: %s", attrType)
	}
}

// validateText validates text and textarea values
func (v *AttributeValidator) validateText(value interface{}, rules map[string]interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("text value must be a string, got %T", value)
	}

	// Min length validation
	if minLen, exists := rules["min_length"]; exists {
		minLenInt, err := v.toInt(minLen)
		if err != nil {
			return fmt.Errorf("invalid min_length rule: %w", err)
		}
		if len(str) < minLenInt {
			return fmt.Errorf("text length %d is less than minimum %d", len(str), minLenInt)
		}
	}

	// Max length validation
	if maxLen, exists := rules["max_length"]; exists {
		maxLenInt, err := v.toInt(maxLen)
		if err != nil {
			return fmt.Errorf("invalid max_length rule: %w", err)
		}
		if len(str) > maxLenInt {
			return fmt.Errorf("text length %d exceeds maximum %d", len(str), maxLenInt)
		}
	}

	// Pattern validation
	if pattern, exists := rules["pattern"]; exists {
		patternStr, ok := pattern.(string)
		if !ok {
			return fmt.Errorf("pattern rule must be a string")
		}
		matched, err := regexp.MatchString(patternStr, str)
		if err != nil {
			return fmt.Errorf("invalid pattern: %w", err)
		}
		if !matched {
			return fmt.Errorf("text does not match required pattern: %s", patternStr)
		}
	}

	return nil
}

// validateNumber validates number values
func (v *AttributeValidator) validateNumber(value interface{}, rules map[string]interface{}) error {
	var num float64

	// Convert value to float64
	switch val := value.(type) {
	case float64:
		num = val
	case float32:
		num = float64(val)
	case int:
		num = float64(val)
	case int32:
		num = float64(val)
	case int64:
		num = float64(val)
	case string:
		// Try to parse string to float
		parsed, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return fmt.Errorf("number value must be numeric, got string: %s", val)
		}
		num = parsed
	default:
		return fmt.Errorf("number value must be numeric, got %T", value)
	}

	// Min value validation
	if minVal, exists := rules["min"]; exists {
		minFloat, err := v.toFloat64(minVal)
		if err != nil {
			return fmt.Errorf("invalid min rule: %w", err)
		}
		if num < minFloat {
			return fmt.Errorf("number value %.2f is less than minimum %.2f", num, minFloat)
		}
	}

	// Max value validation
	if maxVal, exists := rules["max"]; exists {
		maxFloat, err := v.toFloat64(maxVal)
		if err != nil {
			return fmt.Errorf("invalid max rule: %w", err)
		}
		if num > maxFloat {
			return fmt.Errorf("number value %.2f exceeds maximum %.2f", num, maxFloat)
		}
	}

	// Decimal places validation
	if decimals, exists := rules["decimals"]; exists {
		decimalsInt, err := v.toInt(decimals)
		if err != nil {
			return fmt.Errorf("invalid decimals rule: %w", err)
		}
		// Check decimal places by converting to string
		numStr := fmt.Sprintf("%f", num)
		parts := strings.Split(numStr, ".")
		if len(parts) > 1 {
			actualDecimals := len(strings.TrimRight(parts[1], "0"))
			if actualDecimals > decimalsInt {
				return fmt.Errorf("number has %d decimal places, maximum allowed is %d", actualDecimals, decimalsInt)
			}
		}
	}

	return nil
}

// validateBoolean validates boolean values
func (v *AttributeValidator) validateBoolean(value interface{}) error {
	switch val := value.(type) {
	case bool:
		return nil
	case string:
		lowerVal := strings.ToLower(val)
		if lowerVal != "true" && lowerVal != "false" {
			return fmt.Errorf("boolean value must be 'true' or 'false', got: %s", val)
		}
		return nil
	default:
		return fmt.Errorf("boolean value must be bool or string, got %T", value)
	}
}

// validateSelect validates select (single choice) values
func (v *AttributeValidator) validateSelect(value interface{}, options []domain.AttributeOption) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("select value must be a string, got %T", value)
	}

	if len(options) == 0 {
		return fmt.Errorf("select attribute has no valid options")
	}

	// Check if value matches any option
	for _, opt := range options {
		if opt.Value == str {
			return nil
		}
	}

	// Collect valid values for error message
	validValues := make([]string, len(options))
	for i, opt := range options {
		validValues[i] = opt.Value
	}

	return fmt.Errorf("invalid select value '%s', must be one of: %s", str, strings.Join(validValues, ", "))
}

// validateMultiselect validates multiselect (multiple choice) values
func (v *AttributeValidator) validateMultiselect(value interface{}, options []domain.AttributeOption) error {
	// Value can be array of strings or JSON array
	var values []string

	switch val := value.(type) {
	case []string:
		values = val
	case []interface{}:
		values = make([]string, len(val))
		for i, v := range val {
			str, ok := v.(string)
			if !ok {
				return fmt.Errorf("multiselect values must be strings, got %T at index %d", v, i)
			}
			values[i] = str
		}
	case string:
		// Single value as string
		values = []string{val}
	default:
		return fmt.Errorf("multiselect value must be array of strings, got %T", value)
	}

	if len(values) == 0 {
		return fmt.Errorf("multiselect value cannot be empty")
	}

	if len(options) == 0 {
		return fmt.Errorf("multiselect attribute has no valid options")
	}

	// Validate each value
	validOptionValues := make(map[string]bool)
	for _, opt := range options {
		validOptionValues[opt.Value] = true
	}

	for _, val := range values {
		if !validOptionValues[val] {
			validValues := make([]string, 0, len(options))
			for _, opt := range options {
				validValues = append(validValues, opt.Value)
			}
			return fmt.Errorf("invalid multiselect value '%s', must be one of: %s", val, strings.Join(validValues, ", "))
		}
	}

	return nil
}

// validateDate validates date values
func (v *AttributeValidator) validateDate(value interface{}, rules map[string]interface{}) error {
	var dateTime time.Time

	// Parse date value
	switch val := value.(type) {
	case time.Time:
		dateTime = val
	case string:
		// Try multiple date formats
		formats := []string{
			"2006-01-02",
			"2006-01-02T15:04:05Z07:00",
			"2006-01-02 15:04:05",
		}
		var err error
		parsed := false
		for _, format := range formats {
			dateTime, err = time.Parse(format, val)
			if err == nil {
				parsed = true
				break
			}
		}
		if !parsed {
			return fmt.Errorf("invalid date format, expected YYYY-MM-DD or RFC3339: %s", val)
		}
	default:
		return fmt.Errorf("date value must be time.Time or string, got %T", value)
	}

	// Min date validation
	if minDate, exists := rules["min_date"]; exists {
		minDateStr, ok := minDate.(string)
		if !ok {
			return fmt.Errorf("min_date rule must be a string")
		}
		minDateTime, err := time.Parse("2006-01-02", minDateStr)
		if err != nil {
			return fmt.Errorf("invalid min_date format: %w", err)
		}
		if dateTime.Before(minDateTime) {
			return fmt.Errorf("date %s is before minimum date %s", dateTime.Format("2006-01-02"), minDateStr)
		}
	}

	// Max date validation
	if maxDate, exists := rules["max_date"]; exists {
		maxDateStr, ok := maxDate.(string)
		if !ok {
			return fmt.Errorf("max_date rule must be a string")
		}
		maxDateTime, err := time.Parse("2006-01-02", maxDateStr)
		if err != nil {
			return fmt.Errorf("invalid max_date format: %w", err)
		}
		if dateTime.After(maxDateTime) {
			return fmt.Errorf("date %s is after maximum date %s", dateTime.Format("2006-01-02"), maxDateStr)
		}
	}

	return nil
}

// validateColor validates hex color values
func (v *AttributeValidator) validateColor(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("color value must be a string, got %T", value)
	}

	// Validate hex color format (#RRGGBB or #RGB)
	hexColorPattern := `^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`
	matched, err := regexp.MatchString(hexColorPattern, str)
	if err != nil {
		return fmt.Errorf("invalid hex color pattern: %w", err)
	}

	if !matched {
		return fmt.Errorf("invalid hex color format, expected #RRGGBB or #RGB, got: %s", str)
	}

	return nil
}

// Helper methods

// toInt converts interface{} to int
func (v *AttributeValidator) toInt(value interface{}) (int, error) {
	switch val := value.(type) {
	case int:
		return val, nil
	case int32:
		return int(val), nil
	case int64:
		return int(val), nil
	case float64:
		return int(val), nil
	case string:
		return strconv.Atoi(val)
	default:
		return 0, fmt.Errorf("cannot convert %T to int", value)
	}
}

// toFloat64 converts interface{} to float64
func (v *AttributeValidator) toFloat64(value interface{}) (float64, error) {
	switch val := value.(type) {
	case float64:
		return val, nil
	case float32:
		return float64(val), nil
	case int:
		return float64(val), nil
	case int32:
		return float64(val), nil
	case int64:
		return float64(val), nil
	case string:
		return strconv.ParseFloat(val, 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", value)
	}
}
