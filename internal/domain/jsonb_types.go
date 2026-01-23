// Package domain contains shared business entities and types
package domain

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// JSONBStringMap is a custom type for JSONB columns that store map[string]string
type JSONBStringMap map[string]string

// Scan implements sql.Scanner interface for reading from database
func (m *JSONBStringMap) Scan(value interface{}) error {
	if value == nil {
		*m = make(JSONBStringMap)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: expected []byte, got %T", value)
	}

	if len(bytes) == 0 {
		*m = make(JSONBStringMap)
		return nil
	}

	var result map[string]string
	if err := json.Unmarshal(bytes, &result); err != nil {
		return fmt.Errorf("failed to unmarshal JSONB: %w", err)
	}

	*m = JSONBStringMap(result)
	return nil
}

// Value implements driver.Valuer interface for writing to database
func (m JSONBStringMap) Value() (driver.Value, error) {
	if m == nil || len(m) == 0 {
		return []byte("{}"), nil
	}

	bytes, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSONB: %w", err)
	}

	return bytes, nil
}

// ToMap converts JSONBStringMap to map[string]string
func (m JSONBStringMap) ToMap() map[string]string {
	if m == nil {
		return make(map[string]string)
	}
	return map[string]string(m)
}
