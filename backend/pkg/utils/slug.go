package utils

import (
	"regexp"
	"strings"
	"unicode"
)

// GenerateSlug преобразует строку в slug для URL
// Например: "Hello World" -> "hello-world"
func GenerateSlug(s string) string {
	// Переводим строку в нижний регистр
	s = strings.ToLower(s)

	// Заменяем пробелы на дефисы
	s = strings.ReplaceAll(s, " ", "-")

	// Удаляем все символы, кроме букв, цифр и дефисов
	reg := regexp.MustCompile(`[^a-zA-Z0-9\p{L}-]`)
	s = reg.ReplaceAllString(s, "")

	// Заменяем последовательности дефисов на один дефис
	reg = regexp.MustCompile("-+")
	s = reg.ReplaceAllString(s, "-")

	// Удаляем начальные и конечные дефисы
	s = strings.Trim(s, "-")

	// Транслитерация кириллицы и специальных символов
	s = transliterate(s)

	return s
}

// transliterate преобразует кириллицу и специальные символы в латиницу
func transliterate(s string) string {
	translit := map[rune]string{
		// Кириллица (общие символы)
		'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "yo", 'ж': "zh",
		'з': "z", 'и': "i", 'й': "y", 'к': "k", 'л': "l", 'м': "m", 'н': "n", 'о': "o",
		'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u", 'ф': "f", 'х': "h", 'ц': "ts",
		'ч': "ch", 'ш': "sh", 'щ': "sch", 'ъ': "", 'ы': "y", 'ь': "", 'э': "e", 'ю': "yu",
		'я': "ya",
		// Сербские специфические символы
		'ђ': "đ", 'ј': "j", 'љ': "lj", 'њ': "nj", 'ћ': "ć", 'џ': "dž",
		'Ђ': "Đ", 'Ј': "J", 'Љ': "LJ", 'Њ': "NJ", 'Ћ': "Ć", 'Џ': "DŽ",
		// Специальные символы
		'æ': "ae", 'ø': "oe", 'å': "a", 'ä': "a", 'ö': "o", 'ü': "u", 'ß': "ss",
		'þ': "th", 'ð': "d", 'œ': "oe", 'ç': "c", 'ñ': "n",
	}

	var result strings.Builder
	for _, c := range s {
		if unicode.IsLetter(c) && !unicode.IsDigit(c) && !isASCII(c) && c != '-' {
			if val, ok := translit[c]; ok {
				result.WriteString(val)
			}
		} else {
			result.WriteRune(c)
		}
	}

	return result.String()
}

// isASCII проверяет, является ли символ ASCII-символом (в диапазоне 0-127)
func isASCII(r rune) bool {
	return r < 128
}
