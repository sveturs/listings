package middleware

import (
	"regexp"
	"strings"
)

// SensitiveDataMasker маскирует чувствительные данные в логах
type SensitiveDataMasker struct {
	patterns map[string]*regexp.Regexp
}

// NewSensitiveDataMasker создает новый маскировщик данных
func NewSensitiveDataMasker() *SensitiveDataMasker {
	return &SensitiveDataMasker{
		patterns: map[string]*regexp.Regexp{
			"password":      regexp.MustCompile(`(?i)(password|pwd|pass)["\s:=]+([^"\s,}]+)`),
			"token":         regexp.MustCompile(`(?i)(token|jwt|bearer)["\s:=]+([^"\s,}]+)`),
			"cookie":        regexp.MustCompile(`(?i)(cookie|session)["\s:=]+([^"\s,}]+)`),
			"authorization": regexp.MustCompile(`(?i)(authorization)["\s:=]+([^"\s,}]+)`),
			"api_key":       regexp.MustCompile(`(?i)(api_key|apikey)["\s:=]+([^"\s,}]+)`),
			"credit_card":   regexp.MustCompile(`\b\d{4}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}\b`),
			"email":         regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`),
		},
	}
}

// Mask маскирует чувствительные данные в строке
func (m *SensitiveDataMasker) Mask(input string) string {
	output := input

	// Маскируем пароли
	output = m.patterns["password"].ReplaceAllStringFunc(output, func(match string) string {
		parts := m.patterns["password"].FindStringSubmatch(match)
		if len(parts) > 1 {
			return parts[1] + `: "***MASKED***"`
		}
		return match
	})

	// Маскируем токены
	output = m.patterns["token"].ReplaceAllStringFunc(output, func(match string) string {
		parts := m.patterns["token"].FindStringSubmatch(match)
		if len(parts) > 2 && len(parts[2]) > 10 {
			return parts[1] + `: "` + parts[2][:6] + `...***"`
		}
		return match
	})

	// Маскируем куки и сессии
	output = m.patterns["cookie"].ReplaceAllStringFunc(output, func(match string) string {
		parts := m.patterns["cookie"].FindStringSubmatch(match)
		if len(parts) > 2 && len(parts[2]) > 10 {
			return parts[1] + `: "` + parts[2][:6] + `...***"`
		}
		return match
	})

	// Маскируем email адреса (показываем только первые 3 символа и домен)
	output = m.patterns["email"].ReplaceAllStringFunc(output, func(match string) string {
		parts := strings.Split(match, "@")
		if len(parts) == 2 && len(parts[0]) > 3 {
			return parts[0][:3] + "***@" + parts[1]
		}
		return "***@***"
	})

	return output
}

// MaskStruct маскирует чувствительные поля в структурах (для логирования)
func (m *SensitiveDataMasker) MaskStruct(data interface{}) interface{} {
	// Это базовая реализация. В продакшене лучше использовать reflection
	// или специализированные библиотеки для глубокого копирования и маскирования
	str := strings.ToLower(strings.TrimSpace(data.(string)))
	return m.Mask(str)
}
