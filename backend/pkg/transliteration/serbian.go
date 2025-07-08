package transliteration

import (
	"strings"
	"unicode"
)

// SerbianTransliterator handles bidirectional transliteration between Serbian Cyrillic and Latin scripts
type SerbianTransliterator struct {
	cyrillicToLatin map[rune]string
	latinToCyrillic map[string]string
}

// NewSerbianTransliterator creates a new instance of Serbian transliterator
func NewSerbianTransliterator() *SerbianTransliterator {
	return &SerbianTransliterator{
		cyrillicToLatin: cyrillicToLatinMap(),
		latinToCyrillic: latinToCyrillicMap(),
	}
}

// cyrillicToLatinMap returns mapping from Cyrillic to Latin characters
func cyrillicToLatinMap() map[rune]string {
	return map[rune]string{
		// Lowercase
		'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d",
		'ђ': "đ", 'е': "e", 'ж': "ž", 'з': "z", 'и': "i",
		'ј': "j", 'к': "k", 'л': "l", 'љ': "lj", 'м': "m",
		'н': "n", 'њ': "nj", 'о': "o", 'п': "p", 'р': "r",
		'с': "s", 'т': "t", 'ћ': "ć", 'у': "u", 'ф': "f",
		'х': "h", 'ц': "c", 'ч': "č", 'џ': "dž", 'ш': "š",
		// Uppercase
		'А': "A", 'Б': "B", 'В': "V", 'Г': "G", 'Д': "D",
		'Ђ': "Đ", 'Е': "E", 'Ж': "Ž", 'З': "Z", 'И': "I",
		'Ј': "J", 'К': "K", 'Л': "L", 'Љ': "Lj", 'М': "M",
		'Н': "N", 'Њ': "Nj", 'О': "O", 'П': "P", 'Р': "R",
		'С': "S", 'Т': "T", 'Ћ': "Ć", 'У': "U", 'Ф': "F",
		'Х': "H", 'Ц': "C", 'Ч': "Č", 'Џ': "Dž", 'Ш': "Š",
	}
}

// latinToCyrillicMap returns mapping from Latin to Cyrillic digraphs and characters
func latinToCyrillicMap() map[string]string {
	return map[string]string{
		// Digraphs (must be checked first)
		"lj": "љ", "Lj": "Љ", "LJ": "Љ",
		"nj": "њ", "Nj": "Њ", "NJ": "Њ",
		"dž": "џ", "Dž": "Џ", "DŽ": "Џ",
		// Single characters lowercase
		"a": "а", "b": "б", "v": "в", "g": "г", "d": "д",
		"đ": "ђ", "e": "е", "ž": "ж", "z": "з", "i": "и",
		"j": "ј", "k": "к", "l": "л", "m": "м", "n": "н",
		"o": "о", "p": "п", "r": "р", "s": "с", "t": "т",
		"ć": "ћ", "u": "у", "f": "ф", "h": "х", "c": "ц",
		"č": "ч", "š": "ш",
		// Single characters uppercase
		"A": "А", "B": "Б", "V": "В", "G": "Г", "D": "Д",
		"Đ": "Ђ", "E": "Е", "Ž": "Ж", "Z": "З", "I": "И",
		"J": "Ј", "K": "К", "L": "Л", "M": "М", "N": "Н",
		"O": "О", "P": "П", "R": "Р", "S": "С", "T": "Т",
		"Ć": "Ћ", "U": "У", "F": "Ф", "H": "Х", "C": "Ц",
		"Č": "Ч", "Š": "Ш",
	}
}

// ToLatin converts Serbian Cyrillic text to Latin script
func (t *SerbianTransliterator) ToLatin(text string) string {
	var result strings.Builder
	result.Grow(len(text) * 2) // Pre-allocate for potential digraphs

	runes := []rune(text)
	for i, r := range runes {
		if latin, ok := t.cyrillicToLatin[r]; ok {
			// Handle special case for uppercase digraphs in all-caps words
			if latin == "Lj" || latin == "Nj" || latin == "Dž" {
				// Check if this is part of an all-caps word
				if isPartOfAllCapsWord(runes, i) {
					// Convert to all uppercase
					switch latin {
					case "Lj":
						latin = "LJ"
					case "Nj":
						latin = "NJ"
					case "Dž":
						latin = "DŽ"
					}
				}
			}
			result.WriteString(latin)
		} else {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// ToCyrillic converts Serbian Latin text to Cyrillic script
func (t *SerbianTransliterator) ToCyrillic(text string) string {
	var result strings.Builder
	result.Grow(len(text))

	runes := []rune(text)
	i := 0

	for i < len(runes) {
		matched := false

		// Check digraphs first (2 characters)
		if i+1 < len(runes) {
			digraph := string(runes[i : i+2])
			if cyrillic, ok := t.latinToCyrillic[digraph]; ok {
				result.WriteString(cyrillic)
				i += 2
				matched = true
			}
		}

		// If no digraph matched, check single character
		if !matched {
			char := string(runes[i])
			if cyrillic, ok := t.latinToCyrillic[char]; ok {
				result.WriteString(cyrillic)
			} else {
				result.WriteRune(runes[i])
			}
			i++
		}
	}

	return result.String()
}

// TransliterateForSearch generates all possible variants for search
// Returns original text plus transliterated versions
func (t *SerbianTransliterator) TransliterateForSearch(query string) []string {
	variants := make([]string, 0, 3)

	// Always include original
	variants = append(variants, query)

	// Detect script and transliterate
	if containsCyrillic(query) {
		latin := t.ToLatin(query)
		if latin != query {
			variants = append(variants, latin)
		}
	}

	if containsLatin(query) {
		cyrillic := t.ToCyrillic(query)
		if cyrillic != query {
			variants = append(variants, cyrillic)
		}
	}

	// Remove duplicates
	return uniqueStrings(variants)
}

// containsCyrillic checks if text contains any Cyrillic characters
func containsCyrillic(text string) bool {
	for _, r := range text {
		if unicode.Is(unicode.Cyrillic, r) {
			return true
		}
	}
	return false
}

// containsLatin checks if text contains any Latin characters
func containsLatin(text string) bool {
	for _, r := range text {
		if unicode.Is(unicode.Latin, r) {
			return true
		}
	}
	return false
}

// uniqueStrings removes duplicate strings from slice
func uniqueStrings(strings []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(strings))

	for _, s := range strings {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}

	return result
}

// isPartOfAllCapsWord checks if a character at position i is part of an all-caps word
func isPartOfAllCapsWord(runes []rune, pos int) bool {
	// Find word boundaries
	start := pos
	end := pos

	// Find start of word
	for start > 0 && isWordChar(runes[start-1]) {
		start--
	}

	// Find end of word
	for end < len(runes)-1 && isWordChar(runes[end+1]) {
		end++
	}

	// Single character words should not be treated as all-caps unless they're really all uppercase
	if start == end {
		return false
	}

	// Check if all letters in the word are uppercase
	for i := start; i <= end; i++ {
		if unicode.IsLetter(runes[i]) && !unicode.IsUpper(runes[i]) {
			return false
		}
	}

	return true
}

// isWordChar checks if a rune is a word character (letter or digit)
func isWordChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}
