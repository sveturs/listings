// backend/internal/proj/marketplace/service/google_translation.go
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"backend/internal/logger"
	"backend/internal/storage"
)

// min возвращает минимальное из двух целых чисел
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// TranslationProvider представляет тип провайдера перевода
type TranslationProvider string

const (
	// GoogleTranslate провайдер перевода Google Translate
	GoogleTranslate TranslationProvider = "google"
	// OpenAI провайдер перевода OpenAI
	OpenAI TranslationProvider = "openai"
	// Manual ручной перевод
	Manual TranslationProvider = "manual"
	
	// Language detection
	languageAuto = "auto"
	
	// Field names
	fieldNameTitle = "title"
	fieldNameName  = "name"
)

// GoogleTranslationService предоставляет функционал перевода через Google Translate API
type GoogleTranslationService struct {
	apiKey             string
	cache              *sync.Map
	supportedLanguages []string
	// Счетчик переводов для ограничений
	translationCount int
	lastResetMonth   int
	mutex            sync.Mutex
	storage          storage.Storage
}

// Максимальное количество переводов в месяц
const MaxTranslationsPerMonth = 100

// NewGoogleTranslationService создает новый экземпляр сервиса перевода через Google
func NewGoogleTranslationService(apiKey string, storage storage.Storage) (*GoogleTranslationService, error) {
	// Для тестирования позволяем создать сервис даже без ключа API
	if apiKey == "" {
		log.Printf("ВНИМАНИЕ: API ключ Google Translate не указан, будет использоваться имитация перевода")
	}

	// Устанавливаем начальный месяц для счетчика
	currentMonth := time.Now().Month()

	return &GoogleTranslationService{
		apiKey:             apiKey,
		cache:              &sync.Map{},
		supportedLanguages: []string{"sr", "en", "ru"},
		translationCount:   0,
		lastResetMonth:     int(currentMonth),
		mutex:              sync.Mutex{},
		storage:            storage,
	}, nil
}

// resetCounterIfNeeded сбрасывает счетчик если начался новый месяц
func (s *GoogleTranslationService) resetCounterIfNeeded() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	currentMonth := int(time.Now().Month())
	if currentMonth != s.lastResetMonth {
		s.translationCount = 0
		s.lastResetMonth = currentMonth
		log.Printf("Сброс счетчика переводов Google Translate: новый месяц %d", currentMonth)
	}
}

// TranslationCount возвращает количество выполненных переводов в текущем месяце
func (s *GoogleTranslationService) TranslationCount() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.translationCount
}

// TranslationLimit возвращает максимальное количество переводов в месяц
func (s *GoogleTranslationService) TranslationLimit() int {
	return MaxTranslationsPerMonth
}

// Translate переводит текст с одного языка на другой
func (s *GoogleTranslationService) Translate(ctx context.Context, text string, sourceLanguage string, targetLanguage string) (string, error) {
	// Проверка параметров
	if text == "" {
		return "", nil
	}

	// Если исходный и целевой языки одинаковые, просто возвращаем текст
	if sourceLanguage == targetLanguage || (sourceLanguage == languageAuto && guessLanguage(text) == targetLanguage) {
		return text, nil
	}

	// Исправление: если sourceLanguage указан как "auto", используем более конкретное определение
	if sourceLanguage == languageAuto {
		// Простое предположение о языке на основе текста
		sourceLanguage = guessLanguage(text)
		log.Printf("Автоопределение языка для '%s': %s", text, sourceLanguage)
	}

	// Проверяем и сбрасываем счетчик, если нужно
	s.resetCounterIfNeeded()

	// Проверка на превышение лимита
	s.mutex.Lock()
	if s.translationCount >= MaxTranslationsPerMonth {
		s.mutex.Unlock()
		return "", fmt.Errorf("превышен лимит (%d) переводов Google Translate в этом месяце", MaxTranslationsPerMonth)
	}

	// Проверка кеша
	cacheKey := fmt.Sprintf("%s:%s:%s", text, sourceLanguage, targetLanguage)
	if cached, ok := s.cache.Load(cacheKey); ok {
		s.mutex.Unlock()
		return cached.(string), nil
	}

	// Увеличиваем счетчик
	s.translationCount++
	currentCount := s.translationCount
	s.mutex.Unlock()

	log.Printf("Google Translate: выполняется перевод '%s' с %s на %s (%d/%d в этом месяце)",
		text, sourceLanguage, targetLanguage, currentCount, MaxTranslationsPerMonth)

	// Если API ключ не предоставлен, используем примитивную имитацию перевода
	if s.apiKey == "" {
		// Простая имитация перевода для тестирования
		var translatedText string

		// Простая имитация перевода
		switch targetLanguage {
		case "en":
			// Для русских текстов, переводим в формат "[RU->EN] Текст"
			if sourceLanguage == "ru" {
				translatedText = fmt.Sprintf("[RU->EN] %s", text)
			} else if sourceLanguage == "sr" {
				translatedText = fmt.Sprintf("[SR->EN] %s", text)
			} else {
				translatedText = text
			}
		case "ru":
			// Для других языков, переводим в формат "[Lang->RU] Текст"
			if sourceLanguage == "en" {
				translatedText = fmt.Sprintf("[EN->RU] %s", text)
			} else if sourceLanguage == "sr" {
				translatedText = fmt.Sprintf("[SR->RU] %s", text)
			} else {
				translatedText = text
			}
		case "sr":
			// Для других языков, переводим в формат "[Lang->SR] Текст"
			if sourceLanguage == "en" {
				translatedText = fmt.Sprintf("[EN->SR] %s", text)
			} else if sourceLanguage == "ru" {
				translatedText = fmt.Sprintf("[RU->SR] %s", text)
			} else {
				translatedText = text
			}
		default:
			translatedText = text
		}

		log.Printf("Имитация перевода: '%s' -> '%s'", text, translatedText)

		// Сохраняем в кеш
		s.cache.Store(cacheKey, translatedText)

		return translatedText, nil
	}

	// Если API ключ есть, используем реальный Google Translate API
	// Формируем URL для запроса
	apiURL := "https://translation.googleapis.com/language/translate/v2"

	// Добавляем логирование для отладки API ключа
	apiKeyLength := len(s.apiKey)
	if apiKeyLength > 0 {
		log.Printf("Используется API ключ Google Translate (длина: %d, первые 4 символа: %s***)",
			apiKeyLength, s.apiKey[:min(4, apiKeyLength)])
	} else {
		log.Printf("ОШИБКА: API ключ Google Translate пустой")
	}

	// Формируем параметры запроса
	data := url.Values{}
	data.Set("q", text)
	data.Set("source", sourceLanguage)
	data.Set("target", targetLanguage)
	data.Set("format", "html") // Поддержка HTML в тексте
	data.Set("key", s.apiKey)

	// Создаем запрос
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("ошибка создания запроса: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Выполняем запрос
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.Error().Err(err).Msg("Failed to close response body")
		}
	}()

	// Считываем тело ответа для логирования
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		log.Printf("Ошибка API Google Translate: %s, %s", resp.Status, string(bodyBytes))
		return "", fmt.Errorf("ошибка API Google Translate: %s, %s", resp.Status, string(bodyBytes))
	}

	// Преобразуем тело ответа обратно в io.Reader для декодирования
	bodyReader := strings.NewReader(string(bodyBytes))

	// Парсим ответ
	var result struct {
		Data struct {
			Translations []struct {
				TranslatedText string `json:"translatedText"`
			} `json:"translations"`
		} `json:"data"`
	}

	if err := json.NewDecoder(bodyReader).Decode(&result); err != nil {
		log.Printf("Ошибка декодирования ответа: %v. Тело ответа: %s", err, string(bodyBytes))
		return "", fmt.Errorf("ошибка декодирования ответа: %w", err)
	}

	if len(result.Data.Translations) == 0 {
		log.Printf("Google Translate не вернул перевод. Тело ответа: %s", string(bodyBytes))
		return "", fmt.Errorf("Google Translate не вернул перевод")
	}

	translatedText := result.Data.Translations[0].TranslatedText
	log.Printf("Перевод выполнен: '%s' -> '%s'", text, translatedText)

	// Сохраняем в кеш
	s.cache.Store(cacheKey, translatedText)

	return translatedText, nil
}

// guessLanguage выполняет простое определение языка текста
func guessLanguage(text string) string {
	// Проверяем наличие кириллических символов
	cyrillicRegex := regexp.MustCompile(`[\p{Cyrillic}]`)
	hasCyrillic := cyrillicRegex.MatchString(text)

	if hasCyrillic {
		// Проверяем специфичные для русского языка буквы
		russianSpecific := "ёъыэ"
		for _, char := range russianSpecific {
			if strings.ContainsRune(strings.ToLower(text), char) {
				return "ru"
			}
		}

		// Если нет специфичных русских букв, это может быть сербская кириллица
		return "sr"
	}

	// Проверяем сербские латинские специфичные символы
	serbianLatinRegex := regexp.MustCompile(`[čćđšž]`)
	if serbianLatinRegex.MatchString(strings.ToLower(text)) {
		return "sr"
	}

	// По умолчанию предполагаем английский
	return "en"
}

// TranslateToAllLanguages переводит текст на все поддерживаемые языки
func (s *GoogleTranslationService) TranslateToAllLanguages(ctx context.Context, text string) (map[string]string, error) {
	// Определяем язык исходного текста - используем guessLanguage вместо "auto"
	sourceLanguage := guessLanguage(text)
	log.Printf("Определен язык для текста '%s': %s", text, sourceLanguage)

	// Результаты переводов
	translations := make(map[string]string)

	// Сохраняем оригинальный текст в словаре переводов
	translations[sourceLanguage] = text

	// ВАЖНО: Google Translate не поддерживает модерацию, поэтому используем исходный текст напрямую
	// без предварительной модерации, в отличие от OpenAI

	// Переводим на все поддерживаемые языки
	for _, targetLang := range s.supportedLanguages {
		// Пропускаем перевод на язык оригинала
		if targetLang == sourceLanguage {
			log.Printf("Пропускаем перевод на язык оригинала: %s", targetLang)
			continue
		}

		// Используем метод Translate
		translatedText, err := s.Translate(ctx, text, sourceLanguage, targetLang)
		if err != nil {
			log.Printf("Warning: translation to %s failed: %v", targetLang, err)
			continue
		}

		translations[targetLang] = translatedText
	}

	return translations, nil
}

// TranslateWithContext выполняет перевод с учетом контекста
func (s *GoogleTranslationService) TranslateWithContext(ctx context.Context, text string, sourceLanguage string, targetLanguage string, context string, fieldName string) (string, error) {
	// Для Google Translation API мы не можем передать дополнительный контекст,
	// поэтому просто используем базовый метод перевода
	return s.Translate(ctx, text, sourceLanguage, targetLanguage)
}

// TranslateEntityFields переводит поля сущности
func (s *GoogleTranslationService) TranslateEntityFields(ctx context.Context, sourceLanguage string, targetLanguages []string, fields map[string]string) (map[string]map[string]string, error) {
	results := make(map[string]map[string]string)

	// Если sourceLanguage указан как "auto", определяем его на основе поля title или первого имеющегося поля
	if sourceLanguage == languageAuto || sourceLanguage == "" {
		// Сначала проверяем поле title или name
		for fieldName, text := range fields {
			if (fieldName == fieldNameTitle || fieldName == fieldNameName) && text != "" {
				sourceLanguage = guessLanguage(text)
				log.Printf("Определен язык источника из поля %s: %s", fieldName, sourceLanguage)
				break
			}
		}

		// Если не нашли подходящее поле, используем первое непустое
		if sourceLanguage == languageAuto || sourceLanguage == "" {
			for _, text := range fields {
				if text != "" {
					sourceLanguage = guessLanguage(text)
					log.Printf("Определен язык источника из первого непустого поля: %s", sourceLanguage)
					break
				}
			}
		}

		// Если всё равно не определили, используем английский по умолчанию
		if sourceLanguage == languageAuto || sourceLanguage == "" {
			sourceLanguage = "en"
			log.Printf("Не удалось определить язык источника, используем по умолчанию: %s", sourceLanguage)
		}
	}

	// ВАЖНО: Пропускаем модерацию для Google Translate
	// Используем исходные тексты напрямую
	originalFields := make(map[string]string)

	for fieldName, text := range fields {
		if text == "" {
			continue
		}
		originalFields[fieldName] = text
	}

	// Сохраняем оригинальный текст для исходного языка
	results[sourceLanguage] = originalFields

	// Переводим на другие языки
	for _, targetLang := range targetLanguages {
		if targetLang == sourceLanguage {
			log.Printf("Пропускаем перевод на язык оригинала: %s", targetLang)
			continue
		}

		results[targetLang] = make(map[string]string)

		for fieldName, originalText := range originalFields {
			// Переводим текст без предварительной модерации
			translatedText, err := s.Translate(ctx, originalText, sourceLanguage, targetLang)
			if err != nil {
				log.Printf("Error translating field %s to %s: %v", fieldName, targetLang, err)
				continue
			}

			results[targetLang][fieldName] = translatedText
			log.Printf("Перевод поля %s с %s на %s: '%s' -> '%s'",
				fieldName, sourceLanguage, targetLang, originalText, translatedText)
		}
	}

	return results, nil
}

// DetectLanguage определяет язык текста
// Вместо использования API определяем язык локально с помощью guessLanguage
func (s *GoogleTranslationService) DetectLanguage(ctx context.Context, text string) (string, float64, error) {
	// Определяем язык с помощью нашей функции guessLanguage
	detectedLang := guessLanguage(text)
	log.Printf("DetectLanguage определил язык для '%s': %s", text, detectedLang)

	// Возвращаем определенный язык с высокой уверенностью
	return detectedLang, 0.95, nil
}

// ModerateText выполняет модерацию текста
// Google Translate API не имеет функциональности для модерации контента
// поэтому мы просто возвращаем исходный текст без изменений
func (s *GoogleTranslationService) ModerateText(ctx context.Context, text string, language string) (string, error) {
	// Google Translate не поддерживает модерацию, поэтому просто возвращаем оригинальный текст
	log.Printf("Google Translate не поддерживает модерацию контента, возвращаем оригинальный текст")
	return text, nil
}
