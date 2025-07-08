package transliteration

import (
	"reflect"
	"testing"
)

func TestSerbianTransliterator_ToLatin(t *testing.T) {
	trans := NewSerbianTransliterator()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic Cyrillic to Latin",
			input:    "Добар дан",
			expected: "Dobar dan",
		},
		{
			name:     "All Serbian Cyrillic letters lowercase",
			input:    "абвгдђежзијклљмнњопрстћуфхцчџш",
			expected: "abvgdđežzijklljmnnjoprstćufhcčdžš",
		},
		{
			name:     "All Serbian Cyrillic letters uppercase",
			input:    "АБВГДЂЕЖЗИЈКЛЉМНЊОПРСТЋУФХЦЧЏШ",
			expected: "ABVGDĐEŽZIJKLLJMNNJOPRSTĆUFHCČDŽŠ", // All caps word = all caps digraphs
		},
		{
			name:     "Mixed case with digraphs",
			input:    "Љубљана је леп град",
			expected: "Ljubljana je lep grad",
		},
		{
			name:     "Text with numbers and punctuation",
			input:    "Стан 65м2, цена: 1.200€",
			expected: "Stan 65m2, cena: 1.200€",
		},
		{
			name:     "Special characters preserved",
			input:    "е-маил: тест@пример.рс",
			expected: "e-mail: test@primer.rs",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Latin text unchanged",
			input:    "Already in Latin",
			expected: "Already in Latin",
		},
		{
			name:     "Real estate example",
			input:    "Продаје се стан у Новом Саду",
			expected: "Prodaje se stan u Novom Sadu",
		},
		{
			name:     "Car example",
			input:    "Фолксваген Голф, џип",
			expected: "Folksvagen Golf, džip",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := trans.ToLatin(tt.input)
			if result != tt.expected {
				t.Errorf("ToLatin() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSerbianTransliterator_ToCyrillic(t *testing.T) {
	trans := NewSerbianTransliterator()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic Latin to Cyrillic",
			input:    "Dobar dan",
			expected: "Добар дан",
		},
		{
			name:     "All Serbian Latin letters lowercase",
			input:    "abcčćdđefghijklmnoprsštuvzž",
			expected: "абцчћдђефгхијклмнопрсштувзж",
		},
		{
			name:     "Digraph lj",
			input:    "Ljubljana",
			expected: "Љубљана",
		},
		{
			name:     "Digraph nj",
			input:    "Njegoš",
			expected: "Његош",
		},
		{
			name:     "Digraph dž",
			input:    "Džep",
			expected: "Џеп",
		},
		{
			name:     "Mixed case digraphs",
			input:    "LJUBLJANA, Njegov, DŽip",
			expected: "ЉУБЉАНА, Његов, Џип",
		},
		{
			name:     "Text with numbers",
			input:    "Stan 65m2, cena: 1.200€",
			expected: "Стан 65м2, цена: 1.200€",
		},
		{
			name:     "Special characters preserved",
			input:    "e-mail: test@primer.rs",
			expected: "е-маил: тест@пример.рс",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Real estate example",
			input:    "Prodaje se stan u Novom Sadu",
			expected: "Продаје се стан у Новом Саду",
		},
		{
			name:     "Car example with dž",
			input:    "Peugeot džip",
			expected: "Пеугеот џип",
		},
		{
			name:     "Prevent false digraph matching",
			input:    "polje", // should not convert 'lj' here as it's in the middle
			expected: "поље",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := trans.ToCyrillic(tt.input)
			if result != tt.expected {
				t.Errorf("ToCyrillic() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSerbianTransliterator_TransliterateForSearch(t *testing.T) {
	trans := NewSerbianTransliterator()

	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Cyrillic input",
			input:    "стан",
			expected: []string{"стан", "stan"},
		},
		{
			name:     "Latin input",
			input:    "stan",
			expected: []string{"stan", "стан"},
		},
		{
			name:     "Mixed script input",
			input:    "BMW серија",
			expected: []string{"BMW серија", "BMW serija", "БМW серија"},
		},
		{
			name:     "Input with digraphs",
			input:    "Ljubljana",
			expected: []string{"Ljubljana", "Љубљана"},
		},
		{
			name:     "Cyrillic with digraphs",
			input:    "Њујорк",
			expected: []string{"Њујорк", "Njujork"},
		},
		{
			name:     "Numbers and special chars",
			input:    "65м2",
			expected: []string{"65м2", "65m2"},
		},
		{
			name:     "Already transliterated",
			input:    "test123",
			expected: []string{"test123", "тест123"},
		},
		{
			name:     "Empty string",
			input:    "",
			expected: []string{""},
		},
		{
			name:     "Real search query - apartment",
			input:    "гарсоњера",
			expected: []string{"гарсоњера", "garsonjera"},
		},
		{
			name:     "Real search query - car brand",
			input:    "фолксваген",
			expected: []string{"фолксваген", "folksvagen"},
		},
		{
			name:     "Location name",
			input:    "Нови Сад",
			expected: []string{"Нови Сад", "Novi Sad"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := trans.TransliterateForSearch(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("TransliterateForSearch() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestContainsCyrillic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Pure Cyrillic", "привет", true},
		{"Pure Latin", "hello", false},
		{"Mixed", "hello привет", true},
		{"Numbers", "123", false},
		{"Cyrillic with numbers", "стан 123", true},
		{"Empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsCyrillic(tt.input)
			if result != tt.expected {
				t.Errorf("containsCyrillic(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestContainsLatin(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Pure Latin", "hello", true},
		{"Pure Cyrillic", "привет", false},
		{"Mixed", "hello привет", true},
		{"Numbers", "123", false},
		{"Latin with numbers", "test 123", true},
		{"Empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsLatin(tt.input)
			if result != tt.expected {
				t.Errorf("containsLatin(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func BenchmarkToLatin(b *testing.B) {
	trans := NewSerbianTransliterator()
	text := "Продаје се стан у Новом Саду, површине 65м2"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = trans.ToLatin(text)
	}
}

func BenchmarkToCyrillic(b *testing.B) {
	trans := NewSerbianTransliterator()
	text := "Prodaje se stan u Novom Sadu, površine 65m2"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = trans.ToCyrillic(text)
	}
}

func BenchmarkTransliterateForSearch(b *testing.B) {
	trans := NewSerbianTransliterator()
	query := "стан Нови Сад"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = trans.TransliterateForSearch(query)
	}
}
