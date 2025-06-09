// backend/internal/proj/marketplace/service/translation.go
package service

import (
	"backend/internal/storage"
	"context"
	"fmt"
	"html"
	"log"
	"regexp"
	"strings"
	"sync"

	"github.com/sashabaranov/go-openai"
)

type TranslationService struct {
	client             *openai.Client
	cache              *sync.Map
	supportedLanguages []string
	storage            storage.Storage
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
	serbianCyrillicSpecific := "ђћџљњ"
	serbianLatinSpecific := "đćčžš"

	// Очищаем текст от HTML тегов для анализа
	cleanText := html.UnescapeString(text)
	cleanTextLower := strings.ToLower(cleanText)

	// Регулярное выражение для удаления HTML-тегов
	htmlTagRegex := regexp.MustCompile(`<[^>]+>`)
	cleanText = htmlTagRegex.ReplaceAllString(cleanText, " ")
	cleanTextLower = strings.ToLower(cleanText)

	// Проверяем наличие русских букв
	for _, char := range russianSpecific {
		if strings.ContainsRune(cleanTextLower, char) {
			return "ru", 1.0, nil
		}
	}

	// Проверяем наличие сербских кириллических букв
	for _, char := range serbianCyrillicSpecific {
		if strings.ContainsRune(cleanTextLower, char) {
			return "sr", 1.0, nil
		}
	}

	// Проверяем наличие сербских латинских букв
	for _, char := range serbianLatinSpecific {
		if strings.ContainsRune(cleanTextLower, char) {
			return "sr", 1.0, nil
		}
	}

	// Проверка на сербский латинский по характерным словам
	serbianLatinWords := []string{
		"nije",
		"jeste",
		"nešto",
		"gde",
		"kako",
		"hvala",
		"zdravo",
		"dobar dan",
		"ja sam",
		"maskica",
		"telefon",
		"za",
		"iz",
		"koji",
		"kvalitet",
		"maskice",
		"karaktera",
		"slika",
		"prikazuje",
	}

	// Проверка английского по наличию артиклей
	englishArticles := []string{" the ", " a ", " an "}

	// Проверяем сербские латинские слова
	serbianWordCount := 0
	for _, word := range serbianLatinWords {
		if strings.Contains(cleanTextLower, word) {
			serbianWordCount++
		}
	}

	// Проверяем английские артикли
	englishArticleCount := 0
	for _, article := range englishArticles {
		if strings.Contains(" "+cleanTextLower+" ", article) {
			englishArticleCount++
		}
	}

	// Простая эвристика: если есть несколько сербских слов и мало/нет английских артиклей
	if serbianWordCount >= 3 && englishArticleCount < 2 {
		return "sr", 0.9, nil
	}

	// Если простые эвристики не сработали, используем OpenAI для определения
	resp, err := s.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a language detection expert specializing in distinguishing between Serbian Latin, Russian Cyrillic, and English. Reply with exactly one of these language codes: 'ru', 'sr', or 'en'. Nothing else.",
				},
				{
					Role: openai.ChatMessageRoleUser,
					Content: fmt.Sprintf(`Analyze this text and determine the language.
If text uses Cyrillic, analyze if it's Russian or Serbian Cyrillic.
If text uses Latin script:
1. Look for Serbian Latin characters: đ, ć, č, ž, š
2. Check for Serbian words like "nije", "jeste", "nešto", "gde", "kako", "telefon", "maskica"
3. Serbian Latin lacks articles "the", "a", "an" which are common in English
4. Serbian often has words ending with -ski, -ati, -ica

Text to analyze:
"%s"

Reply with only: ru, sr, or en`, cleanText),
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
	moderatedText, err := s.ModerateText(ctx, text, sourceLanguage)
	if err != nil {
		return nil, fmt.Errorf("moderation failed: %w", err)
	}

	// Переводим на все остальные поддерживаемые языки
	for _, targetLang := range s.supportedLanguages {
		// Пропускаем язык оригинала
		if targetLang == sourceLanguage {
			continue
		}

		// Используем новую функцию с контекстом
		translatedText, err := s.TranslateWithContext(ctx, moderatedText, sourceLanguage, targetLang, "", "")
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

	// Проверяем наличие заголовка (title/name/header) для контекста
	var title string
	for fieldName, text := range fields {
		fieldNameLower := strings.ToLower(fieldName)
		if fieldNameLower == "title" || fieldNameLower == "name" || fieldNameLower == "header" {
			title = text
			break
		}
	}

	for fieldName, text := range fields {
		if text == "" {
			continue
		}

		moderatedText, err := s.ModerateText(ctx, text, sourceLanguage)
		if err != nil {
			log.Printf("Error moderating field %s: %v", fieldName, err)
			continue
		}
		moderatedFields[fieldName] = moderatedText
	}

	// Сохраняем модерированный текст для исходного языка
	results[sourceLanguage] = moderatedFields

	// Переводим на другие языки с учетом контекста
	for _, targetLang := range targetLanguages {
		if targetLang == sourceLanguage {
			continue
		}

		results[targetLang] = make(map[string]string)

		// Подготавливаем контекст из заголовка
		context := ""
		if title != "" {
			context = fmt.Sprintf("Context from title/name: %s. ", title)
		}

		for fieldName, moderatedText := range moderatedFields {
			// Используем TranslateWithContext для учета контекста
			translatedText, err := s.TranslateWithContext(ctx, moderatedText, sourceLanguage, targetLang, context, fieldName)
			if err != nil {
				log.Printf("Error translating field %s to %s: %v", fieldName, targetLang, err)
				continue
			}

			results[targetLang][fieldName] = translatedText
		}
	}

	return results, nil
}
func (s *TranslationService) TranslateWithContext(ctx context.Context, text string, sourceLanguage string, targetLanguage string, context string, fieldName string) (string, error) {
	// Проверяем кеш
	cacheKey := fmt.Sprintf("%s:%s:%s:%s:%s", text, sourceLanguage, targetLanguage, context, fieldName)
	if cached, ok := s.cache.Load(cacheKey); ok {
		return cached.(string), nil
	}

	// Предварительно декодируем HTML-сущности, чтобы OpenAI мог правильно переводить текст
	decodedText := html.UnescapeString(text)

	// Обрабатываем HTML-теги перед отправкой на перевод
	htmlTagRegex := regexp.MustCompile(`<[^>]+>`)

	// Создаем словари для сохранения
	tagReplacements := make(map[string]string)

	// Заменяем HTML-теги на плейсхолдеры
	tagCounter := 0
	processedText := htmlTagRegex.ReplaceAllStringFunc(decodedText, func(match string) string {
		placeholder := fmt.Sprintf("__HTML_TAG_%d__", tagCounter)
		tagReplacements[placeholder] = match
		tagCounter++
		return placeholder
	})

	// Добавляем информацию о языке перевода в промпт для улучшения контекста
	languageContext := ""
	switch targetLanguage {
	case "ru":
		languageContext = "Translate to natural, grammatically correct Russian language. Use 'чехол' for 'maskica/case' and 'фиолетовый' for 'ljubicasta/purple'."
	case "en":
		languageContext = "Translate to natural, grammatically correct English. Use 'case' for 'maskica' and 'purple' for 'ljubicasta'."
	case "sr":
		languageContext = "Translate to natural, grammatically correct Serbian using Cyrillic script."
	}

	// Улучшаем промпт с учетом типа поля
	fieldContext := ""
	if fieldName != "" {
		fieldNameLower := strings.ToLower(fieldName)
		if fieldNameLower == "title" || fieldNameLower == "name" || fieldNameLower == "header" {
			fieldContext = "This is a product title/header. In target language, use appropriate product terminology. For example, 'Maska/mask' in Serbian should be translated as 'case' in English or 'чехол' in Russian, not directly transliterated."
		} else if fieldNameLower == "description" || fieldNameLower == "content" {
			fieldContext = "This is a product description. Use natural, marketing-appropriate language in the target language."
		}
	}

	// Строим запрос к API с четкими инструкциями
	systemPrompt := fmt.Sprintf(`You are a professional translator specializing in e-commerce product descriptions. 
You translate from %s to %s accurately while maintaining natural language flow.
%s
%s
Important: DO NOT transliterate brand names or product names - properly translate them.
DO NOT include any translator's notes or comments in your output.`,
		sourceLanguage, targetLanguage, languageContext, fieldContext)

	userPrompt := fmt.Sprintf(`Translate this text from %s to %s:

%s

IMPORTANT INSTRUCTIONS:
1. DO NOT translate placeholders with format __HTML_TAG_X__ - keep them exactly as they are
2. Translate ALL OTHER content naturally and professionally
3. Do not include any comments, notes, or explanations
4. Return ONLY the translated text with the original placeholders intact
5. Never use transliteration for common nouns - always use proper translation`,
		sourceLanguage,
		targetLanguage,
		processedText)

	// Добавляем контекст, если он есть
	if context != "" {
		userPrompt = fmt.Sprintf("Context for accurate translation: %s\n\n%s", context, userPrompt)
	}

	// Вызываем API OpenAI
	resp, err := s.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4, // Используем GPT-4 для более качественного перевода
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userPrompt,
				},
			},
			Temperature: 0.1, // Низкая температура для более точного перевода
		},
	)
	if err != nil {
		// Если GPT-4 недоступен, пробуем с GPT-3.5
		log.Printf("GPT-4 translation failed, falling back to GPT-3.5: %v", err)
		resp, err = s.client.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleSystem,
						Content: systemPrompt,
					},
					{
						Role:    openai.ChatMessageRoleUser,
						Content: userPrompt,
					},
				},
				Temperature: 0.1,
			},
		)
		if err != nil {
			return "", err
		}
	}

	translatedText := resp.Choices[0].Message.Content

	// Восстанавливаем HTML-теги
	for placeholder, htmlTag := range tagReplacements {
		translatedText = strings.Replace(translatedText, placeholder, htmlTag, 1)
	}

	// Кешируем результат
	s.cache.Store(cacheKey, translatedText)

	return translatedText, nil
}
func (s *TranslationService) Translate(ctx context.Context, text string, sourceLanguage string, targetLanguage string) (string, error) {
	// Делегируем работу функции TranslateWithContext с пустым контекстом
	return s.TranslateWithContext(ctx, text, sourceLanguage, targetLanguage, "", "")
}

func (s *TranslationService) ModerateText(ctx context.Context, text string, language string) (string, error) {
	if text == "" {
		return "", nil
	}

	// Декодируем HTML-сущности для лучшей обработки
	decodedText := html.UnescapeString(text)

	// Обрабатываем HTML-теги перед модерацией
	htmlTagRegex := regexp.MustCompile(`<[^>]+>`)

	// Словарь для сохранения тегов
	tagReplacements := make(map[string]string)

	// Заменяем HTML-теги плейсхолдерами
	tagCounter := 0
	processedText := htmlTagRegex.ReplaceAllStringFunc(decodedText, func(match string) string {
		placeholder := fmt.Sprintf("__HTML_TAG_%d__", tagCounter)
		tagReplacements[placeholder] = match
		tagCounter++
		return placeholder
	})

	resp, err := s.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a content moderator. Your task is to check the input text and:\n1. If it contains profanity or offensive language - replace those words with neutral alternatives\n2. If the text is clean - return it exactly as is\nNEVER add any comments or explanations about moderation. DO NOT modify any placeholders with format __HTML_TAG_X__.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: fmt.Sprintf("Moderate this text. DO NOT change placeholders like __HTML_TAG_0__:\n\n%s", processedText),
				},
			},
			Temperature: 0,
		},
	)
	if err != nil {
		return "", err
	}

	moderatedText := strings.TrimSpace(resp.Choices[0].Message.Content)

	// Восстанавливаем HTML-теги
	for placeholder, htmlTag := range tagReplacements {
		moderatedText = strings.Replace(moderatedText, placeholder, htmlTag, 1)
	}

	return moderatedText, nil
}
