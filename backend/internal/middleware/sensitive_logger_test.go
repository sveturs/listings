package middleware

import (
	"testing"
)

func TestSensitiveDataMasker_CategoryIDs(t *testing.T) {
	masker := NewSensitiveDataMasker()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "category_ids parameter should not be masked",
			input:    "category_ids=1011",
			expected: "category_ids=1011",
		},
		{
			name:     "multiple category_ids should not be masked",
			input:    "category_ids=1011,1012,1013",
			expected: "category_ids=1011,1012,1013",
		},
		{
			name:     "query with category_ids and other params",
			input:    "query=test&category_ids=1011&page=1",
			expected: "query=test&category_ids=1011&page=1",
		},
		{
			name:     "actual phone number should be masked",
			input:    "phone: +7 916 123-4567",
			expected: "phone: +XX-XXX-XXX-XXXX",
		},
		{
			name:     "phone number in different format",
			input:    "Contact: +1-555-123-4567",
			expected: "Contact: +XX-XXX-XXX-XXXX",
		},
		{
			name:     "email should be masked",
			input:    "email=test@example.com",
			expected: "email=tes***@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := masker.Mask(tt.input)
			if result != tt.expected {
				t.Errorf("Mask() = %v, want %v", result, tt.expected)
			}
		})
	}
}