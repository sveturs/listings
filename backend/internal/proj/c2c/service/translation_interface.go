// backend/internal/proj/c2c/service/translation_interface.go
package service

import "context"

type TranslationServiceInterface interface {
	// Базовый метод перевода
	Translate(ctx context.Context, text string, sourceLanguage string, targetLanguage string) (string, error)

	// Определение языка текста
	DetectLanguage(ctx context.Context, text string) (string, float64, error)

	// Перевод на все поддерживаемые языки
	TranslateToAllLanguages(ctx context.Context, text string) (map[string]string, error)

	// Перевод конкретных полей сущности
	TranslateEntityFields(ctx context.Context, sourceLanguage string, targetLanguages []string, fields map[string]string) (map[string]map[string]string, error)

	// Модерация текста
	ModerateText(ctx context.Context, text string, language string) (string, error)

	// Перевод с учетом контекста
	TranslateWithContext(ctx context.Context, text string, sourceLanguage string, targetLanguage string, context string, fieldName string) (string, error)

	// Перевод с смягчением тона (модерация мата и агрессивных выражений)
	TranslateWithToneModeration(ctx context.Context, text string, sourceLanguage string, targetLanguage string, moderateTone bool) (string, error)
}
