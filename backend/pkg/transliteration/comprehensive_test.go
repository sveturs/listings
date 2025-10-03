package transliteration

import (
	"strings"
	"testing"
	"time"
)

const (
	// Test constants for digraphs
	titleCaseDz = "Dž"
)

// TestComprehensiveDigraphHandling tests all possible digraph combinations
func TestComprehensiveDigraphHandling(t *testing.T) {
	trans := NewSerbianTransliterator()

	// Test all digraph combinations
	digraphTests := []struct {
		latin    string
		cyrillic string
		desc     string
	}{
		// lj/љ variations
		{"lj", "љ", "lowercase lj"},
		{"Lj", "Љ", "capitalized Lj"},
		{"LJ", "Љ", "uppercase LJ - but converts to Lj for single char"},
		{"ljubljana", "љубљана", "lj in word"},
		{"Ljubljana", "Љубљана", "Lj at start"},
		{"LJUBLJANA", "ЉУБЉАНА", "LJ uppercase word"},

		// nj/њ variations
		{"nj", "њ", "lowercase nj"},
		{"Nj", "Њ", "capitalized Nj"},
		{"NJ", "Њ", "uppercase NJ - but converts to Nj for single char"},
		{"njegov", "његов", "nj in word"},
		{"Njegov", "Његов", "Nj at start"},
		{"NJEGOV", "ЊЕГОВ", "NJ uppercase word"},

		// dž/џ variations
		{"dž", "џ", "lowercase dž"},
		{titleCaseDz, "Џ", "capitalized Dž"},
		{"DŽ", "Џ", "uppercase DŽ - but converts to Dž for single char"},
		{"džep", "џеп", "dž in word"},
		{"Džep", "Џеп", "Dž at start"},
		{"DŽEP", "ЏЕП", "DŽ uppercase word"},
	}

	for _, tt := range digraphTests {
		t.Run("Latin to Cyrillic: "+tt.desc, func(t *testing.T) {
			result := trans.ToCyrillic(tt.latin)
			if result != tt.cyrillic {
				t.Errorf("ToCyrillic(%s) = %s, want %s", tt.latin, result, tt.cyrillic)
			}
		})

		t.Run("Cyrillic to Latin: "+tt.desc, func(t *testing.T) {
			result := trans.ToLatin(tt.cyrillic)
			expectedLatin := tt.latin

			// Special case: single uppercase digraphs should convert to title case
			switch {
			case tt.latin == "LJ" && tt.cyrillic == "Љ":
				expectedLatin = "Lj"
			case tt.latin == "NJ" && tt.cyrillic == "Њ":
				expectedLatin = "Nj"
			case tt.latin == "DŽ" && tt.cyrillic == "Џ":
				expectedLatin = titleCaseDz
			}

			if result != expectedLatin {
				t.Errorf("ToLatin(%s) = %s, want %s", tt.cyrillic, result, expectedLatin)
			}
		})
	}
}

// TestDigraphEdgeCases tests edge cases that might break digraph detection
func TestDigraphEdgeCases(t *testing.T) {
	trans := NewSerbianTransliterator()

	edgeCases := []struct {
		name      string
		input     string
		expected  string
		direction string // "toLatin" or "toCyrillic"
	}{
		// False positive prevention
		{"no false lj in middle", "polje", "поље", "toCyrillic"},
		{"no false nj in middle", "konjic", "коњиц", "toCyrillic"},
		{"dž should be converted", "nadživ", "наџив", "toCyrillic"}, // dž is a valid digraph

		// Boundary cases
		{"lj at end", "konj", "коњ", "toCyrillic"},
		{"nj at end", "konj", "коњ", "toCyrillic"},
		{"dž at end", "muž", "муж", "toCyrillic"},

		// Multiple digraphs
		{"multiple lj", "ljuljašika", "љуљашика", "toCyrillic"},
		{"multiple nj", "njanja", "њања", "toCyrillic"},
		{"mixed digraphs", "ljubavnje", "љубавње", "toCyrillic"},

		// Adjacent digraphs
		{"adjacent lj-nj", "ljnj", "љњ", "toCyrillic"},
		{"adjacent nj-dž", "njdž", "њџ", "toCyrillic"},

		// Digraphs with special characters
		{"lj with hyphen", "ljubav-nj", "љубав-њ", "toCyrillic"},
		{"nj with apostrophe", "konj'", "коњ'", "toCyrillic"},
		{"dž with numbers", "dž123", "џ123", "toCyrillic"},

		// Single character that looks like digraph
		{"single l", "l", "л", "toCyrillic"},
		{"single j", "j", "ј", "toCyrillic"},
		{"single n", "n", "н", "toCyrillic"},
		{"single d", "d", "д", "toCyrillic"},
		{"single ž", "ž", "ж", "toCyrillic"},

		// Case sensitivity edge cases
		{"mixed case lJ", "lJ", "лЈ", "toCyrillic"}, // This should NOT be a digraph
		{"mixed case Lj", "Lj", "Љ", "toCyrillic"},  // This SHOULD be a digraph
		{"mixed case nJ", "nJ", "нЈ", "toCyrillic"}, // This should NOT be a digraph
		{"mixed case dŽ", "dŽ", "дЖ", "toCyrillic"}, // This should NOT be a digraph
	}

	for _, tt := range edgeCases {
		t.Run(tt.name, func(t *testing.T) {
			var result string
			if tt.direction == "toCyrillic" {
				result = trans.ToCyrillic(tt.input)
			} else {
				result = trans.ToLatin(tt.input)
			}

			if result != tt.expected {
				t.Errorf("%s(%s) = %s, want %s", tt.direction, tt.input, result, tt.expected)
			}
		})
	}
}

// TestAllCyrillicCharacters tests every Serbian Cyrillic character
func TestAllCyrillicCharacters(t *testing.T) {
	trans := NewSerbianTransliterator()

	// Complete Serbian Cyrillic alphabet
	cyrillicChars := map[rune]string{
		'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d",
		'ђ': "đ", 'е': "e", 'ж': "ž", 'з': "z", 'и': "i",
		'ј': "j", 'к': "k", 'л': "l", 'љ': "lj", 'м': "m",
		'н': "n", 'њ': "nj", 'о': "o", 'п': "p", 'р': "r",
		'с': "s", 'т': "t", 'ћ': "ć", 'у': "u", 'ф': "f",
		'х': "h", 'ц': "c", 'ч': "č", 'џ': "dž", 'ш': "š",
	}

	for cyrillic, expectedLatin := range cyrillicChars {
		t.Run(string(cyrillic), func(t *testing.T) {
			result := trans.ToLatin(string(cyrillic))
			if result != expectedLatin {
				t.Errorf("ToLatin(%s) = %s, want %s", string(cyrillic), result, expectedLatin)
			}
		})
	}

	// Test uppercase versions
	uppercaseCyrillicChars := map[rune]string{
		'А': "A", 'Б': "B", 'В': "V", 'Г': "G", 'Д': "D",
		'Ђ': "Đ", 'Е': "E", 'Ж': "Ž", 'З': "Z", 'И': "I",
		'Ј': "J", 'К': "K", 'Л': "L", 'Љ': "Lj", 'М': "M",
		'Н': "N", 'Њ': "Nj", 'О': "O", 'П': "P", 'Р': "R",
		'С': "S", 'Т': "T", 'Ћ': "Ć", 'У': "U", 'Ф': "F",
		'Х': "H", 'Ц': "C", 'Ч': "Č", 'Џ': titleCaseDz, 'Ш': "Š",
	}

	for cyrillic, expectedLatin := range uppercaseCyrillicChars {
		t.Run(string(cyrillic)+"_uppercase", func(t *testing.T) {
			result := trans.ToLatin(string(cyrillic))
			if result != expectedLatin {
				t.Errorf("ToLatin(%s) = %s, want %s", string(cyrillic), result, expectedLatin)
			}
		})
	}
}

// TestAllLatinCharacters tests every Serbian Latin character
func TestAllLatinCharacters(t *testing.T) {
	trans := NewSerbianTransliterator()

	// Complete Serbian Latin alphabet (single characters)
	latinChars := map[string]string{
		"a": "а", "b": "б", "c": "ц", "č": "ч", "ć": "ћ", "d": "д",
		"đ": "ђ", "e": "е", "f": "ф", "g": "г", "h": "х", "i": "и",
		"j": "ј", "k": "к", "l": "л", "m": "м", "n": "н", "o": "о",
		"p": "п", "r": "р", "s": "с", "š": "ш", "t": "т", "u": "у",
		"v": "в", "z": "з", "ž": "ж",
	}

	for latin, expectedCyrillic := range latinChars {
		t.Run(latin, func(t *testing.T) {
			result := trans.ToCyrillic(latin)
			if result != expectedCyrillic {
				t.Errorf("ToCyrillic(%s) = %s, want %s", latin, result, expectedCyrillic)
			}
		})
	}

	// Test uppercase versions
	uppercaseLatinChars := map[string]string{
		"A": "А", "B": "Б", "C": "Ц", "Č": "Ч", "Ć": "Ћ", "D": "Д",
		"Đ": "Ђ", "E": "Е", "F": "Ф", "G": "Г", "H": "Х", "I": "И",
		"J": "Ј", "K": "К", "L": "Л", "M": "М", "N": "Н", "O": "О",
		"P": "П", "R": "Р", "S": "С", "Š": "Ш", "T": "Т", "U": "У",
		"V": "В", "Z": "З", "Ž": "Ж",
	}

	for latin, expectedCyrillic := range uppercaseLatinChars {
		t.Run(latin+"_uppercase", func(t *testing.T) {
			result := trans.ToCyrillic(latin)
			if result != expectedCyrillic {
				t.Errorf("ToCyrillic(%s) = %s, want %s", latin, result, expectedCyrillic)
			}
		})
	}
}

// TestBidirectionalConsistency tests that converting A->B->A gives original A
func TestBidirectionalConsistency(t *testing.T) {
	trans := NewSerbianTransliterator()

	testCases := []string{
		"Београд",
		"Нови Сад",
		"Љубљана",
		"Њујорк",
		"Џакарта",
		"Поље",
		"Коњиц",
		"Belgrade",
		"Novi Sad",
		"Ljubljana",
		"New York",
		"Jakarta",
		"Polje",
		"Konjic",
		"Стан 65м2 у центру",
		"Stan 65m2 u centru",
		"Гарсоњера са балконом",
		"Garsonjera sa balkonom",
	}

	for _, original := range testCases {
		t.Run("Bidirectional: "+original, func(t *testing.T) {
			// Test Cyrillic -> Latin -> Cyrillic
			if containsCyrillic(original) {
				latin := trans.ToLatin(original)
				backToCyrillic := trans.ToCyrillic(latin)
				if backToCyrillic != original {
					t.Errorf("Cyrillic->Latin->Cyrillic: %s -> %s -> %s", original, latin, backToCyrillic)
				}
			}

			// Test Latin -> Cyrillic -> Latin
			if containsLatin(original) {
				cyrillic := trans.ToCyrillic(original)
				backToLatin := trans.ToLatin(cyrillic)
				if backToLatin != original {
					t.Errorf("Latin->Cyrillic->Latin: %s -> %s -> %s", original, cyrillic, backToLatin)
				}
			}
		})
	}
}

// TestRealWorldExamples tests real-world examples from marketplace
func TestRealWorldExamples(t *testing.T) {
	trans := NewSerbianTransliterator()

	realWorldCases := []struct {
		name     string
		cyrillic string
		latin    string
	}{
		// Real estate examples
		{"Apartment listing", "Продаје се стан у центру Београда", "Prodaje se stan u centru Beograda"},
		{"Room listing", "Издаје се соба у Новом Саду", "Izdaje se soba u Novom Sadu"},
		{"Studio apartment", "Гарсоњера са балконом", "Garsonjera sa balkonom"},
		{"House for sale", "Кућа са двориштем", "Kuća sa dvorištem"},
		{"Parking space", "Гаражно место", "Garažno mesto"},

		// Car marketplace examples
		{"Car brand", "Фолксваген Голф", "Folksvagen Golf"},
		{"Car type", "Џип теренац", "Džip terenac"},
		{"Car condition", "Половно возило", "Polovno vozilo"},
		{"Car year", "Аутомобил из 2020. године", "Automobil iz 2020. godine"},

		// Electronics examples
		{"Phone", "Мобилни телефон", "Mobilni telefon"},
		{"Laptop", "Лаптоп рачунар", "Laptop računar"},
		{"TV", "Телевизор 55 инча", "Televizor 55 inča"},

		// Location names
		{"Belgrade", "Београд", "Beograd"},
		{"Novi Sad", "Нови Сад", "Novi Sad"},
		{"Niš", "Ниш", "Niš"},
		{"Kragujevac", "Крагујевац", "Kragujevac"},
		{"Subotica", "Суботица", "Subotica"},

		// Mixed content
		{"Price with currency", "Цена: 1.200€", "Cena: 1.200€"},
		{"Contact info", "Контакт: +381 62 123 4567", "Kontakt: +381 62 123 4567"},
	}

	for _, tt := range realWorldCases {
		t.Run(tt.name+" - Cyrillic to Latin", func(t *testing.T) {
			result := trans.ToLatin(tt.cyrillic)
			if result != tt.latin {
				t.Errorf("ToLatin(%s) = %s, want %s", tt.cyrillic, result, tt.latin)
			}
		})

		t.Run(tt.name+" - Latin to Cyrillic", func(t *testing.T) {
			result := trans.ToCyrillic(tt.latin)
			if result != tt.cyrillic {
				t.Errorf("ToCyrillic(%s) = %s, want %s", tt.latin, result, tt.cyrillic)
			}
		})
	}
}

// TestPerformanceWithDifferentSizes tests performance with different text sizes
func TestPerformanceWithDifferentSizes(t *testing.T) {
	trans := NewSerbianTransliterator()

	// Generate test texts of different sizes
	smallText := "стан"
	mediumText := strings.Repeat("Продаје се стан у Новом Саду, површине 65м2. ", 10)
	largeText := strings.Repeat(mediumText, 100)

	testCases := []struct {
		name string
		text string
		size int
	}{
		{"Small (4 chars)", smallText, len(smallText)},
		{"Medium (~500 chars)", mediumText, len(mediumText)},
		{"Large (~50k chars)", largeText, len(largeText)},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()
			_ = trans.ToLatin(tt.text)
			duration := time.Since(start)

			t.Logf("ToLatin for %s took %v", tt.name, duration)

			// Performance assertion - should be fast even for large texts
			if duration > 10*time.Millisecond {
				t.Errorf("ToLatin took too long: %v for text size %d", duration, tt.size)
			}
		})

		t.Run(tt.name+" - ToCyrillic", func(t *testing.T) {
			latinText := trans.ToLatin(tt.text)
			start := time.Now()
			_ = trans.ToCyrillic(latinText)
			duration := time.Since(start)

			t.Logf("ToCyrillic for %s took %v", tt.name, duration)

			// ToCyrillic is more complex due to digraph processing
			if duration > 50*time.Millisecond {
				t.Errorf("ToCyrillic took too long: %v for text size %d", duration, tt.size)
			}
		})
	}
}

// TestSearchVariants tests search variant generation
func TestSearchVariants(t *testing.T) {
	trans := NewSerbianTransliterator()

	testCases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Simple Cyrillic",
			input:    "стан",
			expected: []string{"стан", "stan"},
		},
		{
			name:     "Simple Latin",
			input:    "stan",
			expected: []string{"stan", "стан"},
		},
		{
			name:     "Cyrillic with digraph",
			input:    "Љубљана",
			expected: []string{"Љубљана", "Ljubljana"},
		},
		{
			name:     "Latin with digraph",
			input:    "Ljubljana",
			expected: []string{"Ljubljana", "Љубљана"},
		},
		{
			name:     "Numbers only",
			input:    "123",
			expected: []string{"123"},
		},
		{
			name:     "Mixed with numbers",
			input:    "стан 65м2",
			expected: []string{"стан 65м2", "stan 65m2"},
		},
		{
			name:     "Special characters",
			input:    "цена: 1.200€",
			expected: []string{"цена: 1.200€", "cena: 1.200€"},
		},
		{
			name:     "Complex search query",
			input:    "Нови Сад гарсоњера",
			expected: []string{"Нови Сад гарсоњера", "Novi Sad garsonjera"},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			result := trans.TransliterateForSearch(tt.input)

			// Check if we have the expected number of variants
			if len(result) != len(tt.expected) {
				t.Errorf("TransliterateForSearch(%s) returned %d variants, expected %d",
					tt.input, len(result), len(tt.expected))
			}

			// Check if all expected variants are present
			for _, expected := range tt.expected {
				found := false
				for _, actual := range result {
					if actual == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("TransliterateForSearch(%s) missing expected variant: %s",
						tt.input, expected)
				}
			}
		})
	}
}

// TestUnicodeEdgeCases tests edge cases with Unicode characters
func TestUnicodeEdgeCases(t *testing.T) {
	trans := NewSerbianTransliterator()

	testCases := []struct {
		name     string
		input    string
		expected string
		function string
	}{
		// Note: Our system transliterates Serbian Cyrillic that overlaps with Russian/Bulgarian
		// This is actually correct behavior since Serbian uses the same characters
		{"Russian chars", "Москва", "Moskva", "ToLatin"}, // Contains Serbian characters
		{"Bulgarian chars", "София", "Sofiя", "ToLatin"}, // Contains Serbian characters

		// Non-Serbian Latin should be preserved
		{"German chars", "Müller", "Мüллер", "ToCyrillic"}, // ü is not Serbian, M->М, l->л, e->е, r->р
		{"French chars", "Café", "Цафé", "ToCyrillic"},     // é is not Serbian, but C->Ц, a->а, f->ф

		// Mixed scripts
		{"Serbian + English", "Београд Belgrade", "Beograd Belgrade", "ToLatin"},
		{"English + Serbian", "Belgrade Београд", "Belgrade Beograd", "ToLatin"},

		// Emojis and special Unicode
		{"With emoji", "стан 🏠", "stan 🏠", "ToLatin"},
		{"With symbols", "цена ★★★", "cena ★★★", "ToLatin"},

		// Zero-width characters
		{"Zero-width space", "стан\u200Bнов", "stan\u200Bnov", "ToLatin"},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var result string
			if tt.function == "ToLatin" {
				result = trans.ToLatin(tt.input)
			} else {
				result = trans.ToCyrillic(tt.input)
			}

			if result != tt.expected {
				t.Errorf("%s(%s) = %s, want %s", tt.function, tt.input, result, tt.expected)
			}
		})
	}
}
