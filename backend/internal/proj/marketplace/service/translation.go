// backend/internal/proj/marketplace/service/translation.go
package service

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	openai "github.com/sashabaranov/go-openai"
)

type TranslationService struct {
	client             *openai.Client
	cache              *sync.Map
	supportedLanguages []string
}

func NewTranslationService(apiKey string) (*TranslationService, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required")
	}
	return &TranslationService{
		client:             openai.NewClient(apiKey),
		cache:              &sync.Map{},
		supportedLanguages: []string{"sr", "en", "ru"},
	}, nil
}

func (s *TranslationService) DetectLanguage(ctx context.Context, text string) (string, float64, error) {
	// Сначала проверяем наличие специфических букв
	russianSpecific := "ёъыэьяй"
	serbianSpecific := "đćčžšђћџ"

	// Проверяем наличие русских букв
	for _, char := range russianSpecific {
		if strings.ContainsRune(strings.ToLower(text), char) {
			return "ru", 1.0, nil
		}
	}

	// Проверяем наличие сербских букв
	for _, char := range serbianSpecific {
		if strings.ContainsRune(strings.ToLower(text), char) {
			return "sr", 1.0, nil
		}
	}

	// Если специфических букв не найдено, используем OpenAI для более точного определения
	resp, err := s.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a language detector. Reply with exactly one of these language codes: 'ru', 'sr', or 'en'. Nothing else.",
				},
				{
					Role: openai.ChatMessageRoleUser,
					Content: fmt.Sprintf(`Analyze this text and determine the language.
If text uses Cyrillic but has no specific Russian letters (ё,ъ,ы,э,ь,я,й), analyze grammar and vocabulary.
If text uses Latin script - check if it uses Serbian Latin letters (đ,ć,č,dž,š,ž).
Text: "%s"`, text),
				},
			},
			Temperature: 0,
		},
	)
	if err != nil {
		return "", 0, err
	}

	detectedLang := strings.TrimSpace(strings.ToLower(resp.Choices[0].Message.Content))
	if detectedLang != "ru" && detectedLang != "sr" && detectedLang != "en" {
		return "", 0, fmt.Errorf("unexpected language detected: %s", detectedLang)
	}

	return detectedLang, 1.0, nil
}

// TranslateToAllLanguages переводит текст на все поддерживаемые языки
func (s *TranslationService) TranslateToAllLanguages(ctx context.Context, text string) (map[string]string, error) {
	// Определяем язык исходного текста
	sourceLanguage, _, err := s.DetectLanguage(ctx, text)
	if err != nil {
		return nil, fmt.Errorf("language detection failed: %w", err)
	}

	log.Printf("Detected language: %s for text translation to all languages", sourceLanguage)

	// Результаты переводов
	translations := make(map[string]string)

	// Сохраняем оригинальный текст
	translations[sourceLanguage] = text

	// Сначала модерируем исходный текст
	moderatedText, err := s.moderateText(ctx, text, sourceLanguage)
	if err != nil {
		return nil, fmt.Errorf("moderation failed: %w", err)
	}

	// Переводим на все остальные поддерживаемые языки
	for _, targetLang := range s.supportedLanguages {
		// Пропускаем язык оригинала
		if targetLang == sourceLanguage {
			continue
		}

		translatedText, err := s.Translate(ctx, moderatedText, sourceLanguage, targetLang)
		if err != nil {
			log.Printf("Warning: translation to %s failed: %v", targetLang, err)
			continue
		}

		translations[targetLang] = translatedText
	}

	return translations, nil
}

func (s *TranslationService) TranslateEntityFields(ctx context.Context, sourceLanguage string, targetLanguages []string, fields map[string]string) (map[string]map[string]string, error) {
    results := make(map[string]map[string]string)
    
    // Сначала модерируем исходный текст
    moderatedFields := make(map[string]string)
    for fieldName, text := range fields {
        if text == "" {
            continue
        }
        
        moderatedText, err := s.moderateText(ctx, text, sourceLanguage)
        if err != nil {
            log.Printf("Error moderating field %s: %v", fieldName, err)
            continue
        }
        moderatedFields[fieldName] = moderatedText
    }

    // Сохраняем модерированный текст для исходного языка
    results[sourceLanguage] = moderatedFields
    
    // Переводим на другие языки
    for _, targetLang := range targetLanguages {
        if targetLang == sourceLanguage {
            continue
        }
        
        results[targetLang] = make(map[string]string)
        
        for fieldName, moderatedText := range moderatedFields {
            translatedText, err := s.Translate(ctx, moderatedText, sourceLanguage, targetLang)
            if err != nil {
                log.Printf("Error translating field %s to %s: %v", fieldName, targetLang, err)
                continue
            }
            
            results[targetLang][fieldName] = translatedText
        }
    }
    
    return results, nil
}
func (s *TranslationService) Translate(ctx context.Context, text string, sourceLanguage string, targetLanguage string) (string, error) {
	// Проверяем кеш
	cacheKey := fmt.Sprintf("%s:%s:%s", text, sourceLanguage, targetLanguage)
	if cached, ok := s.cache.Load(cacheKey); ok {
		return cached.(string), nil
	}

	// Сначала модерируем текст
	moderatedText, err := s.moderateText(ctx, text, sourceLanguage)
	if err != nil {
		return "", err
	}

	resp, err := s.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a professional translator. Translate the text while preserving formatting and maintaining a natural, fluent style.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: fmt.Sprintf("Translate this text from %s to %s: %s", sourceLanguage, targetLanguage, moderatedText),
				},
			},
			Temperature: 0.3,
		},
	)
	if err != nil {
		return "", err
	}

	translatedText := resp.Choices[0].Message.Content
	s.cache.Store(cacheKey, translatedText)

	return translatedText, nil
}

func (s *TranslationService) moderateText(ctx context.Context, text string, language string) (string, error) {
    var prompt string
    switch language {
    case "ru":
        prompt = fmt.Sprintf(`Проверь текст на русском языке на нецензурную лексику и оскорбления.
Если найдешь - замени на нейтральные слова схожего смысла.
Верни ТОЛЬКО сам текст без каких-либо пометок и объяснений.
Пример: если на входе "охуенная тачка" - верни "отличная машина"
Текст: "%s"`, text)
    case "sr":
        prompt = fmt.Sprintf(`Proveri tekst na srpskom jeziku.
Ako sadrži psovke ili uvrede - zameni pristojnim rečima istog značenja.
Vrati SAMO obrađeni tekst, bez dodatnih objašnjenja.
Tekst: "%s"`, text)
    case "en":
        prompt = fmt.Sprintf(`Check the text for profanity and offensive language.
Replace any inappropriate words with neutral alternatives of similar meaning.
Return ONLY the processed text without any additional marks or explanations.
Example: if input is "fucking awesome car" - return "really great car"
Text: "%s"`, text)
    default:
        return "", fmt.Errorf("unsupported language for moderation: %s", language)
    }

    resp, err := s.client.CreateChatCompletion(
        ctx,
        openai.ChatCompletionRequest{
            Model: openai.GPT3Dot5Turbo,
            Messages: []openai.ChatCompletionMessage{
                {
                    Role: openai.ChatMessageRoleSystem,
                    Content: "You are a content moderator. Replace inappropriate language with neutral alternatives. Return ONLY the processed text, without any markers, explanations or additional text. Never add prefixes like 'Moderated text:' or similar.",
                },
                {
                    Role: openai.ChatMessageRoleUser,
                    Content: prompt,
                },
            },
            Temperature: 0,
        },
    )
    if err != nil {
        return "", err
    }

    // Очищаем от всех возможных маркеров и лишних пробелов
    cleanedText := strings.TrimSpace(resp.Choices[0].Message.Content)
    cleanedText = strings.Trim(cleanedText, `"'*_[]()`)
    
    return cleanedText, nil
}