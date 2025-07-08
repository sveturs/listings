package transliteration

import (
	"testing"
)

// TestSearchIntegration tests how transliteration integrates with search scenarios
func TestSearchIntegration(t *testing.T) {
	trans := NewSerbianTransliterator()

	// Simulate real search scenarios from marketplace
	searchTests := []struct {
		name        string
		query       string
		document    string
		shouldMatch bool
		desc        string
	}{
		// Basic Cyrillic query vs Latin document
		{
			name:        "Cyrillic query, Latin document",
			query:       "стан",
			document:    "Prodaje se stan u centru grada",
			shouldMatch: true,
			desc:        "User searches in Cyrillic, document is in Latin",
		},
		{
			name:        "Latin query, Cyrillic document",
			query:       "stan",
			document:    "Продаје се стан у центру града",
			shouldMatch: true,
			desc:        "User searches in Latin, document is in Cyrillic",
		},

		// Location names with digraphs
		{
			name:        "Ljubljana search",
			query:       "Љубљана",
			document:    "Apartman u Ljubljani, centar grada",
			shouldMatch: true,
			desc:        "Cyrillic location name vs Latin document",
		},
		{
			name:        "Ljubljana search reverse",
			query:       "Ljubljana",
			document:    "Апартман у Љубљани, центар града",
			shouldMatch: true,
			desc:        "Latin location name vs Cyrillic document",
		},

		// Car brands and models
		{
			name:        "Car brand search",
			query:       "фолксваген",
			document:    "Prodajem Volkswagen Golf 2015 godine",
			shouldMatch: false, // Different words
			desc:        "Car brand in different languages (different transliteration)",
		},
		{
			name:        "BMW search",
			query:       "БМВ",
			document:    "BMW X5 u odličnom stanju",
			shouldMatch: true,
			desc:        "BMW in Cyrillic vs Latin",
		},

		// Apartment types
		{
			name:        "Apartment type search",
			query:       "гарсоњера",
			document:    "Izdajem garsonjeru u Novom Sadu",
			shouldMatch: true,
			desc:        "Apartment type with digraph",
		},

		// Complex queries with multiple words
		{
			name:        "Multi-word search",
			query:       "нови сад стан",
			document:    "Prodajem stan u Novom Sadu, centar",
			shouldMatch: true,
			desc:        "Multi-word query transliteration",
		},

		// Mixed script scenarios
		{
			name:        "Mixed script query",
			query:       "BMW серија",
			document:    "BMW serija 3, dizel motor",
			shouldMatch: true,
			desc:        "Mixed script in query",
		},

		// Edge cases
		{
			name:        "Numbers preserved",
			query:       "65м2",
			document:    "Stan 65m2 u centru Beograda",
			shouldMatch: true,
			desc:        "Numbers and measurements",
		},
		{
			name:        "Price format",
			query:       "1.200€",
			document:    "Cena stana: 1.200€ mesečno",
			shouldMatch: true,
			desc:        "Price format should match exactly",
		},
	}

	for _, tt := range searchTests {
		t.Run(tt.name, func(t *testing.T) {
			// Generate search variants for query
			queryVariants := trans.TransliterateForSearch(tt.query)

			// Generate search variants for document
			documentVariants := trans.TransliterateForSearch(tt.document)

			// Check if any query variant matches any document variant
			matched := false
			for _, qVariant := range queryVariants {
				for _, dVariant := range documentVariants {
					if containsWord(dVariant, qVariant) {
						matched = true
						break
					}
				}
				if matched {
					break
				}
			}

			if matched != tt.shouldMatch {
				t.Errorf("%s: expected match=%v, got match=%v", tt.desc, tt.shouldMatch, matched)
				t.Logf("Query variants: %v", queryVariants)
				t.Logf("Document variants: %v", documentVariants)
			}
		})
	}
}

// TestTransliterationAccuracy tests accuracy of character mappings
func TestTransliterationAccuracy(t *testing.T) {
	trans := NewSerbianTransliterator()

	// Test each Serbian character mapping
	accuracyTests := []struct {
		cyrillic string
		latin    string
		category string
	}{
		// Vowels
		{"а", "a", "vowel"},
		{"е", "e", "vowel"},
		{"и", "i", "vowel"},
		{"о", "o", "vowel"},
		{"у", "u", "vowel"},

		// Simple consonants
		{"б", "b", "consonant"},
		{"в", "v", "consonant"},
		{"г", "g", "consonant"},
		{"д", "d", "consonant"},
		{"з", "z", "consonant"},
		{"к", "k", "consonant"},
		{"л", "l", "consonant"},
		{"м", "m", "consonant"},
		{"н", "n", "consonant"},
		{"п", "p", "consonant"},
		{"р", "r", "consonant"},
		{"с", "s", "consonant"},
		{"т", "t", "consonant"},
		{"ф", "f", "consonant"},
		{"х", "h", "consonant"},
		{"ц", "c", "consonant"},

		// Special Serbian characters
		{"ђ", "đ", "special"},
		{"ж", "ž", "special"},
		{"ј", "j", "special"},
		{"ћ", "ć", "special"},
		{"ч", "č", "special"},
		{"ш", "š", "special"},

		// Digraphs
		{"љ", "lj", "digraph"},
		{"њ", "nj", "digraph"},
		{"џ", "dž", "digraph"},
	}

	for _, tt := range accuracyTests {
		t.Run(tt.cyrillic+"->"+tt.latin, func(t *testing.T) {
			// Test Cyrillic to Latin
			result := trans.ToLatin(tt.cyrillic)
			if result != tt.latin {
				t.Errorf("ToLatin(%s) = %s, want %s", tt.cyrillic, result, tt.latin)
			}

			// Test Latin to Cyrillic
			result = trans.ToCyrillic(tt.latin)
			if result != tt.cyrillic {
				t.Errorf("ToCyrillic(%s) = %s, want %s", tt.latin, result, tt.cyrillic)
			}
		})
	}
}

// TestDigraphPriority tests that digraphs are processed correctly
func TestDigraphPriority(t *testing.T) {
	trans := NewSerbianTransliterator()

	digraphTests := []struct {
		name   string
		input  string
		output string
		desc   string
	}{
		{
			name:   "lj_digraph",
			input:  "ljubav",
			output: "љубав",
			desc:   "lj should be treated as single digraph",
		},
		{
			name:   "nj_digraph",
			input:  "njegov",
			output: "његов",
			desc:   "nj should be treated as single digraph",
		},
		{
			name:   "dz_digraph",
			input:  "džem",
			output: "џем",
			desc:   "dž should be treated as single digraph",
		},
		{
			name:   "multiple_lj",
			input:  "ljuljac",
			output: "љуљац",
			desc:   "Multiple lj digraphs in one word",
		},
		{
			name:   "lj_vs_l_j",
			input:  "polje",
			output: "поље",
			desc:   "lj within word should be digraph",
		},
		{
			name:   "nj_in_middle",
			input:  "banjer",
			output: "бањер",
			desc:   "nj should be converted as digraph in middle of word",
		},
	}

	for _, tt := range digraphTests {
		t.Run(tt.name, func(t *testing.T) {
			result := trans.ToCyrillic(tt.input)
			if result != tt.output {
				t.Errorf("%s: ToCyrillic(%s) = %s, want %s", tt.desc, tt.input, result, tt.output)
			}

			// Test round trip
			backToLatin := trans.ToLatin(result)
			if backToLatin != tt.input {
				t.Errorf("%s: Round trip failed: %s -> %s -> %s", tt.desc, tt.input, result, backToLatin)
			}
		})
	}
}

// containsWord checks if text contains a specific word (simple implementation)
func containsWord(text, word string) bool {
	if len(word) == 0 {
		return false
	}

	// Convert to lowercase for case-insensitive matching
	textLower := toLowerCase(text)
	wordLower := toLowerCase(word)

	// Check if word appears as whole word or substring
	return textLower == wordLower ||
		contains(textLower, " "+wordLower+" ") ||
		startsWith(textLower, wordLower+" ") ||
		endsWith(textLower, " "+wordLower) ||
		contains(textLower, wordLower) // Also check substring matching
}

// Simple helper functions for word matching
func contains(s, substr string) bool {
	return len(substr) <= len(s) && indexOf(s, substr) >= 0
}

func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[0:len(prefix)] == prefix
}

func endsWith(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func toLowerCase(s string) string {
	// Simple ASCII lowercase - for testing purposes
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] >= 'A' && s[i] <= 'Z' {
			result[i] = s[i] + 32
		} else {
			result[i] = s[i]
		}
	}
	return string(result)
}
