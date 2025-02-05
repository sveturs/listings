// backend/internal/proj/marketplace/service/translation_interface.go
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
    ModerateText(ctx context.Context, text string, language string) (string, error)  // Обратите внимание на заглавную M


}