// backend/internal/service/translation/service.go
package translation

import (
    "context"
    "encoding/json"
    "fmt"
    "sync"
    
    "backend/internal/domain/models"
    "backend/internal/storage"
    
    "github.com/sashabaranov/go-openai"
)

type Service struct {
    openAIClient *openai.Client
    storage      storage.Storage
    supportedLanguages []string
    cache       *sync.Map // Простое кеширование переводов
}

func NewTranslationService(apiKey string, storage storage.Storage) *Service {
    return &Service{
        openAIClient: openai.NewClient(apiKey),
        storage:      storage,
        supportedLanguages: []string{"en", "sr", "ru"},
        cache:       &sync.Map{},
    }
}

// Определение языка текста через OpenAI
func (s *Service) DetectLanguage(ctx context.Context, text string) (string, error) {
    prompt := fmt.Sprintf("Determine the language of the following text and respond with only the language code (en, sr, or ru):\n\n%s", text)
    
    resp, err := s.openAIClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
        Model: openai.GPT3Dot5Turbo,
        Messages: []openai.ChatCompletionMessage{
            {Role: "user", Content: prompt},
        },
    })
    if err != nil {
        return "", fmt.Errorf("failed to detect language: %w", err)
    }
    
    detectedLang := resp.Choices[0].Message.Content
    // Проверяем, что обнаруженный язык поддерживается
    for _, lang := range s.supportedLanguages {
        if detectedLang == lang {
            return lang, nil
        }
    }
    
    return "en", nil // По умолчанию считаем английским
}

// Перевод текста
func (s *Service) TranslateText(ctx context.Context, text, fromLang, toLang string) (string, error) {
    // Проверяем кеш
    cacheKey := fmt.Sprintf("%s:%s:%s:%s", text, fromLang, toLang, text)
    if cached, ok := s.cache.Load(cacheKey); ok {
        return cached.(string), nil
    }
    
    prompt := fmt.Sprintf("Translate the following text from %s to %s. Preserve the original formatting and maintain a natural, fluent style:\n\n%s", fromLang, toLang, text)
    
    resp, err := s.openAIClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
        Model: openai.GPT4,
        Messages: []openai.ChatCompletionMessage{
            {Role: "user", Content: prompt},
        },
    })
    if err != nil {
        return "", fmt.Errorf("translation failed: %w", err)
    }
    
    translatedText := resp.Choices[0].Message.Content
    
    // Сохраняем в кеш
    s.cache.Store(cacheKey, translatedText)
    
    return translatedText, nil
}

// Перевод сущности на все поддерживаемые языки
func (s *Service) TranslateEntity(ctx context.Context, entity models.Translatable, entityType string, entityID int) error {
    // Определяем язык оригинала, если он не указан
    originalLang := entity.GetOriginalLanguage()
    if originalLang == "" {
        // Берем первое поле для определения языка
        firstField := entity.GetTranslationFields()[0]
        firstText := reflect.ValueOf(entity).Elem().FieldByName(firstField).String()
        detectedLang, err := s.DetectLanguage(ctx, firstText)
        if err != nil {
            return fmt.Errorf("failed to detect language: %w", err)
        }
        originalLang = detectedLang
    }
    
    translations := make(map[string]map[string]string)
    
    // Переводим на все поддерживаемые языки
    for _, targetLang := range s.supportedLanguages {
        if targetLang == originalLang {
            continue
        }
        
        translations[targetLang] = make(map[string]string)
        
        // Переводим каждое поле
        for _, fieldName := range entity.GetTranslationFields() {
            originalText := reflect.ValueOf(entity).Elem().FieldByName(fieldName).String()
            if originalText == "" {
                continue
            }
            
            translatedText, err := s.TranslateText(ctx, originalText, originalLang, targetLang)
            if err != nil {
                return fmt.Errorf("failed to translate field %s to %s: %w", fieldName, targetLang, err)
            }
            
            // Сохраняем перевод в БД
            translation := &models.Translation{
                EntityType:         entityType,
                EntityID:           entityID,
                FieldName:          fieldName,
                Language:           targetLang,
                TranslatedText:     translatedText,
                IsMachineTranslated: true,
            }
            
            if err := s.storage.SaveTranslation(ctx, translation); err != nil {
                return fmt.Errorf("failed to save translation: %w", err)
            }
            
            translations[targetLang][fieldName] = translatedText
        }
    }
    
    // Устанавливаем переводы в сущность
    entity.SetTranslations(translations)
    
    return nil
}

// Получение переводов для сущности
func (s *Service) GetEntityTranslations(ctx context.Context, entityType string, entityID int) (map[string]map[string]string, error) {
    translations, err := s.storage.GetTranslations(ctx, entityType, entityID)
    if err != nil {
        return nil, fmt.Errorf("failed to get translations: %w", err)
    }
    
    result := make(map[string]map[string]string)
    for _, t := range translations {
        if result[t.Language] == nil {
            result[t.Language] = make(map[string]string)
        }
        result[t.Language][t.FieldName] = t.TranslatedText
    }
    
    return result, nil
}

// Верификация перевода человеком
func (s *Service) VerifyTranslation(ctx context.Context, translationID int, verifiedText string) error {
    translation := &models.Translation{
        ID:            translationID,
        TranslatedText: verifiedText,
        IsVerified:    true,
    }
    
    if err := s.storage.UpdateTranslation(ctx, translation); err != nil {
        return fmt.Errorf("failed to update translation: %w", err)
    }
    
    // Инвалидируем кеш
    s.cache.Delete(fmt.Sprintf("%d", translationID))
    
    return nil
}