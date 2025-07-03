package utils

import (
	"github.com/microcosm-cc/bluemonday"
)

// HTMLSanitizer интерфейс для санитизации HTML
type HTMLSanitizer interface {
	Sanitize(html string) string
}

// StrictSanitizer удаляет все HTML теги
type StrictSanitizer struct {
	policy *bluemonday.Policy
}

// NewStrictSanitizer создает новый строгий санитайзер
func NewStrictSanitizer() *StrictSanitizer {
	// Строгая политика - удаляет все HTML теги
	p := bluemonday.StrictPolicy()
	return &StrictSanitizer{policy: p}
}

// Sanitize очищает HTML от всех тегов
func (s *StrictSanitizer) Sanitize(html string) string {
	return s.policy.Sanitize(html)
}

// BasicSanitizer разрешает базовое форматирование
type BasicSanitizer struct {
	policy *bluemonday.Policy
}

// NewBasicSanitizer создает санитайзер с базовым форматированием
func NewBasicSanitizer() *BasicSanitizer {
	p := bluemonday.NewPolicy()

	// Разрешаем базовое форматирование
	p.AllowElements("b", "i", "u", "strong", "em", "code", "pre")

	// Разрешаем переносы строк
	p.AllowElements("br", "p")

	// Разрешаем ссылки с безопасными протоколами
	p.AllowElements("a")
	p.AllowAttrs("href").OnElements("a")
	p.AllowURLSchemes("http", "https", "mailto")
	p.RequireParseableURLs(true)
	p.RequireNoReferrerOnLinks(true)
	p.AddTargetBlankToFullyQualifiedLinks(true)

	return &BasicSanitizer{policy: p}
}

// Sanitize очищает HTML, оставляя базовое форматирование
func (s *BasicSanitizer) Sanitize(html string) string {
	return s.policy.Sanitize(html)
}

// Глобальный санитайзер для использования в приложении
var (
	strictSanitizer = NewStrictSanitizer()
	basicSanitizer  = NewBasicSanitizer()
)

// SanitizeText удаляет все HTML теги из текста
func SanitizeText(text string) string {
	return strictSanitizer.Sanitize(text)
}

// SanitizeHTML очищает HTML, оставляя безопасное форматирование
func SanitizeHTML(html string) string {
	return basicSanitizer.Sanitize(html)
}
