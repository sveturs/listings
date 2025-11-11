package postgres

import "testing"

// TestTransliterate tests Cyrillic to Latin transliteration
func TestTransliterate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Russian lowercase",
			input:    "привет мир",
			expected: "privet mir",
		},
		{
			name:     "Russian uppercase",
			input:    "ПРИВЕТ МИР",
			expected: "privet mir",
		},
		{
			name:     "Russian mixed case",
			input:    "Стильный керамический",
			expected: "stilnyj keramicheskij",
		},
		{
			name:     "Serbian Cyrillic",
			input:    "Ђорђе Њњ Љљ",
			expected: "djordje njnj ljlj",
		},
		{
			name:     "Mixed Russian and Latin",
			input:    "Электро 220V",
			expected: "elektro 220V",
		},
		{
			name:     "Special Russian characters",
			input:    "Съешь ещё этих французских булок",
			expected: "sesh eshchyo etih francuzskih bulok",
		},
		{
			name:     "Latin unchanged",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := transliterate(tt.input)
			if result != tt.expected {
				t.Errorf("transliterate(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestSlugify tests URL slug generation with Cyrillic support
func TestSlugify(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Russian title",
			input:    "Стильный керамический органайзер",
			expected: "stilnyj-keramicheskij-organajzer",
		},
		{
			name:     "Serbian Cyrillic",
			input:    "Електро отвертка",
			expected: "elektro-otvertka",
		},
		{
			name:     "Mixed Cyrillic and numbers",
			input:    "Телефон iPhone 15 Pro 256GB",
			expected: "telefon-iphone-15-pro-256gb",
		},
		{
			name:     "Russian with special chars",
			input:    "Кружка \"Утренний кофе\" (белая)!",
			expected: "kruzhka-utrennij-kofe-belaya",
		},
		{
			name:     "Latin title unchanged",
			input:    "Stylish Ceramic Organizer",
			expected: "stylish-ceramic-organizer",
		},
		{
			name:     "Multiple spaces and hyphens",
			input:    "Товар   для  дома---и офиса",
			expected: "tovar-dlya-doma-i-ofisa",
		},
		{
			name:     "Leading and trailing spaces",
			input:    "  Товар для дома  ",
			expected: "tovar-dlya-doma",
		},
		{
			name:     "Long Russian title",
			input:    "Очень длинное название товара которое должно быть обрезано по лимиту символов для того чтобы не превышать максимальную длину слага и оставить место для числовых суффиксов при генерации уникальных слагов в системе управления товарами",
			expected: "ochen-dlinnoe-nazvanie-tovara-kotoroe-dolzhno-byt-obrezano-po-limitu-simvolov-dlya-togo-chtoby-ne-prevyshat-maksimalnuyu-dlinu-slaga-i-ostavit-mesto-dlya-chislovyh-suffiksov-pri-generacii-unikalnyh-sl",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Only special chars",
			input:    "!@#$%^&*()",
			expected: "",
		},
		{
			name:     "Serbian specific chars",
			input:    "Ђорђе Пупин - Научник",
			expected: "djordje-pupin-nauchnik",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := slugify(tt.input)
			if result != tt.expected {
				t.Errorf("slugify(%q) = %q, want %q", tt.input, result, tt.expected)
			}

			// Additional check: if input is not empty and contains valid chars,
			// result should not be empty
			if tt.input != "" && tt.expected != "" && result == "" {
				t.Errorf("slugify(%q) returned empty string unexpectedly", tt.input)
			}
		})
	}
}

// TestValidateSlug tests slug validation
func TestValidateSlug(t *testing.T) {
	tests := []struct {
		name      string
		slug      string
		wantError bool
	}{
		{
			name:      "Valid Latin slug",
			slug:      "valid-slug-123",
			wantError: false,
		},
		{
			name:      "Valid transliterated Russian slug",
			slug:      "stilnyj-keramicheskij-organajzer",
			wantError: false,
		},
		{
			name:      "Empty slug",
			slug:      "",
			wantError: true,
		},
		{
			name:      "Too short",
			slug:      "ab",
			wantError: true,
		},
		{
			name:      "Contains uppercase",
			slug:      "Invalid-Slug",
			wantError: true,
		},
		{
			name:      "Contains Cyrillic",
			slug:      "товар-для-дома",
			wantError: true,
		},
		{
			name:      "Consecutive hyphens",
			slug:      "invalid--slug",
			wantError: true,
		},
		{
			name:      "Leading hyphen",
			slug:      "-invalid-slug",
			wantError: true,
		},
		{
			name:      "Trailing hyphen",
			slug:      "invalid-slug-",
			wantError: true,
		},
		{
			name:      "Special characters",
			slug:      "invalid_slug@123",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSlug(tt.slug)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateSlug(%q) error = %v, wantError %v", tt.slug, err, tt.wantError)
			}
		})
	}
}
