package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVariantDetector_ExtractBaseName(t *testing.T) {
	vd := NewVariantDetector()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Russian color and size removal",
			input:    "Футболка красная M",
			expected: "Футболка",
		},
		{
			name:     "English color removal",
			input:    "T-Shirt Blue Size L",
			expected: "T-Shirt Size",
		},
		{
			name:     "Serbian color removal",
			input:    "Majica crna veličina XL",
			expected: "Majica veličina",
		},
		{
			name:     "Multiple attributes",
			input:    "Кроссовки Nike Air Max 2023 черные 42",
			expected: "Кроссовки Nike Air Max",
		},
		{
			name:     "No variant attributes",
			input:    "Базовый товар без вариантов",
			expected: "Базовый товар без вариантов",
		},
		{
			name:     "Trailing punctuation removal",
			input:    "Товар название -",
			expected: "Товар название",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := vd.ExtractBaseName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestVariantDetector_ExtractVariantAttributes(t *testing.T) {
	vd := NewVariantDetector()

	tests := []struct {
		name     string
		input    string
		expected map[string]string
	}{
		{
			name:  "Color only",
			input: "Футболка красная",
			expected: map[string]string{
				"color": "красная", // lowercase, extracted from "Футболка красная"
			},
		},
		{
			name:  "Size only",
			input: "Футболка размер M",
			expected: map[string]string{
				"size": "m", // lowercase because we apply ToLower in ExtractVariantAttributes
			},
		},
		{
			name:  "Color and size",
			input: "Футболка красная размер M",
			expected: map[string]string{
				"color": "красная",
				"size":  "m",
			},
		},
		{
			name:  "Model year",
			input: "Nike Air Max 2023",
			expected: map[string]string{
				"model": "2023",
			},
		},
		{
			name:     "No attributes",
			input:    "Простой товар",
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := vd.ExtractVariantAttributes(tt.input)
			assert.Equal(t, len(tt.expected), len(result))
			for k, v := range tt.expected {
				assert.Equal(t, v, result[k])
			}
		})
	}
}

func TestVariantDetector_GroupProducts(t *testing.T) {
	vd := NewVariantDetector()
	vd.SetMinConfidence(0.5)
	vd.SetMinGroupSize(2)

	tests := []struct {
		name               string
		products           []*ProductVariant
		expectedGroupCount int
		expectedVariantCnt int // для первой группы
		expectedIsGrouped  bool
	}{
		{
			name: "Valid variant group with colors",
			products: []*ProductVariant{
				{Name: "Футболка красная", SKU: "SKU1", Price: 1000},
				{Name: "Футболка синяя", SKU: "SKU2", Price: 1000},
				{Name: "Футболка зеленая", SKU: "SKU3", Price: 1000},
			},
			expectedGroupCount: 1,
			expectedVariantCnt: 3,
			expectedIsGrouped:  true,
		},
		{
			name: "Valid variant group with sizes",
			products: []*ProductVariant{
				{Name: "T-Shirt Size S", SKU: "SKU1", Price: 1000},
				{Name: "T-Shirt Size M", SKU: "SKU2", Price: 1000},
				{Name: "T-Shirt Size L", SKU: "SKU3", Price: 1000},
			},
			expectedGroupCount: 1,
			expectedVariantCnt: 3,
			expectedIsGrouped:  true,
		},
		{
			name: "Too small group (single product)",
			products: []*ProductVariant{
				{Name: "Уникальный товар", SKU: "SKU1", Price: 1000},
			},
			expectedGroupCount: 1,
			expectedVariantCnt: 1,
			expectedIsGrouped:  false,
		},
		{
			name: "Multiple groups",
			products: []*ProductVariant{
				{Name: "Футболка красная", SKU: "SKU1", Price: 1000},
				{Name: "Футболка синяя", SKU: "SKU2", Price: 1000},
				{Name: "Шорты черные", SKU: "SKU3", Price: 2000},
				{Name: "Шорты белые", SKU: "SKU4", Price: 2000},
			},
			expectedGroupCount: 2,
			expectedVariantCnt: 2,
			expectedIsGrouped:  true,
		},
		{
			name: "Low confidence - still grouped (confidence=0.5, min=0.5)",
			products: []*ProductVariant{
				{Name: "Футболка красная", SKU: "SKU1", Price: 1000},
				{Name: "Футболка", SKU: "SKU2", Price: 1000}, // Нет атрибута
			},
			expectedGroupCount: 1, // Группируется, т.к. confidence = 1/2 = 0.5 >= minConfidence(0.5)
			expectedVariantCnt: 2,
			expectedIsGrouped:  true, // Групповой товар
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			groups := vd.GroupProducts(tt.products)
			assert.Equal(t, tt.expectedGroupCount, len(groups))

			if len(groups) > 0 {
				firstGroup := groups[0]
				assert.Equal(t, tt.expectedVariantCnt, firstGroup.VariantCount)
				assert.Equal(t, tt.expectedIsGrouped, firstGroup.IsGrouped)
			}
		})
	}
}

func TestVariantDetector_ValidateVariantGroup(t *testing.T) {
	vd := NewVariantDetector()

	tests := []struct {
		name           string
		group          *VariantGroup
		expectWarnings bool
	}{
		{
			name: "Valid group - no warnings",
			group: &VariantGroup{
				BaseName:          "T-Shirt",
				VariantCount:      2,
				VariantAttributes: []string{"color"},
				Variants: []*ProductVariant{
					{
						SKU:  "SKU1",
						Name: "T-Shirt Red",
						VariantAttributes: map[string]string{
							"color": "red",
						},
					},
					{
						SKU:  "SKU2",
						Name: "T-Shirt Blue",
						VariantAttributes: map[string]string{
							"color": "blue",
						},
					},
				},
			},
			expectWarnings: false,
		},
		{
			name: "Missing attribute - should warn",
			group: &VariantGroup{
				BaseName:          "T-Shirt",
				VariantCount:      2,
				VariantAttributes: []string{"color"},
				Variants: []*ProductVariant{
					{
						SKU:  "SKU1",
						Name: "T-Shirt Red",
						VariantAttributes: map[string]string{
							"color": "red",
						},
					},
					{
						SKU:               "SKU2",
						Name:              "T-Shirt",
						VariantAttributes: map[string]string{}, // Missing color
					},
				},
			},
			expectWarnings: true,
		},
		{
			name: "Duplicate attributes - should warn",
			group: &VariantGroup{
				BaseName:          "T-Shirt",
				VariantCount:      2,
				VariantAttributes: []string{"color"},
				Variants: []*ProductVariant{
					{
						SKU:  "SKU1",
						Name: "T-Shirt Red",
						VariantAttributes: map[string]string{
							"color": "red",
						},
					},
					{
						SKU:  "SKU2",
						Name: "T-Shirt Red",
						VariantAttributes: map[string]string{
							"color": "red",
						},
					},
				},
			},
			expectWarnings: true,
		},
		{
			name: "Single variant - no warnings",
			group: &VariantGroup{
				BaseName:     "Unique Product",
				VariantCount: 1,
				Variants: []*ProductVariant{
					{SKU: "SKU1", Name: "Unique Product"},
				},
			},
			expectWarnings: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warnings := vd.ValidateVariantGroup(tt.group)
			if tt.expectWarnings {
				assert.NotEmpty(t, warnings, "Expected warnings but got none")
			} else {
				assert.Empty(t, warnings, "Expected no warnings but got: %v", warnings)
			}
		})
	}
}

func TestVariantDetector_SetMinConfidence(t *testing.T) {
	vd := NewVariantDetector()

	// Valid values
	vd.SetMinConfidence(0.5)
	assert.Equal(t, 0.5, vd.minConfidence)

	vd.SetMinConfidence(1.0)
	assert.Equal(t, 1.0, vd.minConfidence)

	vd.SetMinConfidence(0.0)
	assert.Equal(t, 0.0, vd.minConfidence)

	// Invalid values - should not change
	vd.SetMinConfidence(0.8)
	currentValue := vd.minConfidence
	vd.SetMinConfidence(-0.1)
	assert.Equal(t, currentValue, vd.minConfidence)

	vd.SetMinConfidence(1.5)
	assert.Equal(t, currentValue, vd.minConfidence)
}

func TestVariantDetector_SetMinGroupSize(t *testing.T) {
	vd := NewVariantDetector()

	// Valid values
	vd.SetMinGroupSize(1)
	assert.Equal(t, 1, vd.minGroupSize)

	vd.SetMinGroupSize(5)
	assert.Equal(t, 5, vd.minGroupSize)

	// Invalid values - should not change
	currentValue := vd.minGroupSize
	vd.SetMinGroupSize(0)
	assert.Equal(t, currentValue, vd.minGroupSize)

	vd.SetMinGroupSize(-1)
	assert.Equal(t, currentValue, vd.minGroupSize)
}

// Benchmark тесты
func BenchmarkVariantDetector_ExtractBaseName(b *testing.B) {
	vd := NewVariantDetector()
	productName := "Кроссовки Nike Air Max 2023 черные размер 42"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vd.ExtractBaseName(productName)
	}
}

func BenchmarkVariantDetector_ExtractVariantAttributes(b *testing.B) {
	vd := NewVariantDetector()
	productName := "Кроссовки Nike Air Max 2023 черные размер 42"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vd.ExtractVariantAttributes(productName)
	}
}

func BenchmarkVariantDetector_GroupProducts(b *testing.B) {
	vd := NewVariantDetector()
	products := make([]*ProductVariant, 100)
	for i := 0; i < 100; i++ {
		products[i] = &ProductVariant{
			Name:  "Футболка красная",
			SKU:   "SKU" + string(rune(i)),
			Price: 1000,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vd.GroupProducts(products)
	}
}
