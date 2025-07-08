package transliteration

import (
	"testing"
)

// TestOpenSearchIntegration demonstrates how transliteration should work with search
func TestOpenSearchIntegration(t *testing.T) {
	trans := NewSerbianTransliterator()

	// Test typical search scenarios
	testCases := []struct {
		name         string
		searchQuery  string
		documentText string
		explanation  string
	}{
		{
			name:         "Basic Cyrillic to Latin",
			searchQuery:  "стан",
			documentText: "stan",
			explanation:  "User searches 'стан', should match document with 'stan'",
		},
		{
			name:         "Basic Latin to Cyrillic",
			searchQuery:  "stan",
			documentText: "стан",
			explanation:  "User searches 'stan', should match document with 'стан'",
		},
		{
			name:         "Digraph search",
			searchQuery:  "Љубљана",
			documentText: "Ljubljana",
			explanation:  "Cyrillic digraph should match Latin equivalent",
		},
		{
			name:         "Mixed script query",
			searchQuery:  "BMW серија",
			documentText: "BMW serija",
			explanation:  "Mixed scripts should be handled properly",
		},
		{
			name:         "Complex real-world example",
			searchQuery:  "гарсоњера Нови Сад",
			documentText: "garsonjera Novi Sad",
			explanation:  "Real apartment search scenario",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Generate search variants for query
			queryVariants := trans.TransliterateForSearch(tc.searchQuery)

			// Generate search variants for document
			docVariants := trans.TransliterateForSearch(tc.documentText)

			t.Logf("Search query: %s", tc.searchQuery)
			t.Logf("Query variants: %v", queryVariants)
			t.Logf("Document text: %s", tc.documentText)
			t.Logf("Document variants: %v", docVariants)
			t.Logf("Explanation: %s", tc.explanation)

			// In a real OpenSearch integration, these variants would be used
			// to create a multi-match query that searches across all variants

			// Demonstrate that we have overlapping variants
			hasOverlap := false
			for _, qVar := range queryVariants {
				for _, dVar := range docVariants {
					if qVar == dVar {
						hasOverlap = true
						t.Logf("Matching variant found: %s", qVar)
						break
					}
				}
				if hasOverlap {
					break
				}
			}

			if !hasOverlap {
				t.Errorf("No matching variants found between query and document")
			}
		})
	}
}

// TestTransliterationSpeed measures speed for different text sizes
func TestTransliterationSpeed(t *testing.T) {
	trans := NewSerbianTransliterator()

	testSizes := []struct {
		name string
		text string
		size int
	}{
		{"Small", "стан", 4},
		{"Medium", "Продаје се стан у Новом Саду, површине 65м2, у близини центра", 66},
		{"Large", generateLargeText(1000), 1000},
	}

	for _, ts := range testSizes {
		t.Run(ts.name, func(t *testing.T) {
			// Measure ToLatin performance
			variants := trans.TransliterateForSearch(ts.text)

			if len(variants) == 0 {
				t.Errorf("No variants generated for text of size %d", ts.size)
			}

			t.Logf("Text size: %d characters", ts.size)
			t.Logf("Generated %d variants", len(variants))

			// Verify that we don't have duplicates
			seen := make(map[string]bool)
			duplicates := 0
			for _, variant := range variants {
				if seen[variant] {
					duplicates++
				}
				seen[variant] = true
			}

			if duplicates > 0 {
				t.Errorf("Found %d duplicate variants", duplicates)
			}
		})
	}
}

// TestEdgeCasesForSearch tests edge cases that might occur in search
func TestEdgeCasesForSearch(t *testing.T) {
	trans := NewSerbianTransliterator()

	edgeCases := []struct {
		name  string
		input string
		desc  string
	}{
		{"Empty string", "", "Should handle empty input gracefully"},
		{"Only numbers", "123", "Should preserve numbers"},
		{"Only punctuation", ".,!?", "Should preserve punctuation"},
		{"Mixed content", "стан 65м² за 1.200€", "Should handle mixed content"},
		{"Very long word", "супермегаекстраординарнооооооооооо", "Should handle very long words"},
		{"Unicode emojis", "стан 🏠🔑", "Should preserve emojis"},
		{"Special chars", "e-mail: test@example.com", "Should handle email addresses"},
		{"URLs", "www.пример.рс", "Should handle URLs"},
	}

	for _, tc := range edgeCases {
		t.Run(tc.name, func(t *testing.T) {
			variants := trans.TransliterateForSearch(tc.input)

			t.Logf("Input: %s", tc.input)
			t.Logf("Variants: %v", variants)
			t.Logf("Description: %s", tc.desc)

			// Should always include the original
			found := false
			for _, variant := range variants {
				if variant == tc.input {
					found = true
					break
				}
			}

			if !found && tc.input != "" {
				t.Errorf("Original input not found in variants")
			}
		})
	}
}

// generateLargeText creates a large text for performance testing
func generateLargeText(size int) string {
	base := "Продаје се стан у Новом Саду. "
	result := ""
	for len(result) < size {
		result += base
	}
	return result[:size]
}

// BenchmarkSearchVariants benchmarks the search variant generation
func BenchmarkSearchVariants(b *testing.B) {
	trans := NewSerbianTransliterator()

	queries := []string{
		"стан",
		"Нови Сад",
		"гарсоњера Љубљана",
		"BMW серија 3 дизел",
		"Продаје се стан у центру града",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		query := queries[i%len(queries)]
		_ = trans.TransliterateForSearch(query)
	}
}

// BenchmarkLargeTextTransliteration benchmarks large text processing
func BenchmarkLargeTextTransliteration(b *testing.B) {
	trans := NewSerbianTransliterator()
	largeText := generateLargeText(10000) // 10KB of text

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = trans.ToLatin(largeText)
	}
}
