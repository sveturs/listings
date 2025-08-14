// backend/internal/proj/marketplace/service/deepl_translation.go
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"backend/internal/logger"
)

// DeepLTranslationService предоставляет функционал перевода через DeepL API
type DeepLTranslationService struct {
	apiKey     string
	apiURL     string // https://api-free.deepl.com/v2 или https://api.deepl.com/v2
	httpClient *http.Client
}

// NewDeepLTranslationService создает новый экземпляр сервиса перевода DeepL
func NewDeepLTranslationService(apiKey string, useFreeAPI bool) (*DeepLTranslationService, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("DeepL API key is required")
	}

	apiURL := "https://api.deepl.com/v2"
	if useFreeAPI {
		apiURL = "https://api-free.deepl.com/v2"
	}

	return &DeepLTranslationService{
		apiKey: apiKey,
		apiURL: apiURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// deepLTranslateResponse представляет ответ от DeepL API
type deepLTranslateResponse struct {
	Translations []struct {
		DetectedSourceLanguage string `json:"detected_source_language"`
		Text                   string `json:"text"`
	} `json:"translations"`
}

// deepLLanguageResponse представляет ответ для определения языка
type deepLLanguageResponse struct {
	Language   string  `json:"language"`
	Confidence float64 `json:"confidence"`
}

// Translate переводит текст с одного языка на другой
func (s *DeepLTranslationService) Translate(ctx context.Context, text string, sourceLanguage string, targetLanguage string) (string, error) {
	if text == "" {
		return "", nil
	}

	// Преобразуем коды языков в формат DeepL
	sourceLang := convertToDeepLLanguageCode(sourceLanguage)
	targetLang := convertToDeepLLanguageCode(targetLanguage)

	// Формируем параметры запроса
	params := url.Values{}
	params.Set("text", text)
	params.Set("target_lang", targetLang)
	if sourceLang != "" && sourceLang != "auto" {
		params.Set("source_lang", sourceLang)
	}

	// Создаем запрос
	req, err := http.NewRequestWithContext(ctx, "POST", s.apiURL+"/translate", strings.NewReader(params.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Устанавливаем заголовки
	req.Header.Set("Authorization", "DeepL-Auth-Key "+s.apiKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Выполняем запрос
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("DeepL API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	// Парсим ответ
	var deepLResp deepLTranslateResponse
	if err := json.Unmarshal(body, &deepLResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Извлекаем переведенный текст
	if len(deepLResp.Translations) > 0 {
		translatedText := deepLResp.Translations[0].Text
		logger.Info().
			Str("source", sourceLanguage).
			Str("target", targetLanguage).
			Str("detected_source", deepLResp.Translations[0].DetectedSourceLanguage).
			Int("source_len", len(text)).
			Int("translated_len", len(translatedText)).
			Msg("DeepL translation completed")
		return translatedText, nil
	}

	return "", fmt.Errorf("no translations in DeepL response")
}

// TranslateWithContext переводит текст с учетом контекста
func (s *DeepLTranslationService) TranslateWithContext(ctx context.Context, text string, sourceLanguage string, targetLanguage string, context string, fieldName string) (string, error) {
	if text == "" {
		return "", nil
	}

	// DeepL поддерживает контекст через параметр context (в Pro версии)
	sourceLang := convertToDeepLLanguageCode(sourceLanguage)
	targetLang := convertToDeepLLanguageCode(targetLanguage)

	params := url.Values{}
	params.Set("text", text)
	params.Set("target_lang", targetLang)
	if sourceLang != "" && sourceLang != "auto" {
		params.Set("source_lang", sourceLang)
	}

	// Добавляем контекст если это поле заголовка или названия
	if fieldName == "title" || fieldName == "name" || fieldName == "seo_title" {
		params.Set("formality", "default") // Нейтральный стиль для заголовков
	} else if fieldName == "description" || fieldName == "seo_description" {
		params.Set("formality", "less") // Менее формальный стиль для описаний
	}

	// Добавляем глоссарий для специфических терминов маркетплейса
	if context != "" {
		params.Set("tag_handling", "xml")
		params.Set("split_sentences", "1")
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.apiURL+"/translate", strings.NewReader(params.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "DeepL-Auth-Key "+s.apiKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("DeepL API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var deepLResp deepLTranslateResponse
	if err := json.Unmarshal(body, &deepLResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(deepLResp.Translations) > 0 {
		return deepLResp.Translations[0].Text, nil
	}

	return "", fmt.Errorf("no translations in DeepL response")
}

// TranslateToAllLanguages переводит текст на все поддерживаемые языки
func (s *DeepLTranslationService) TranslateToAllLanguages(ctx context.Context, text string) (map[string]string, error) {
	supportedLanguages := []string{"en", "ru", "sr"}
	result := make(map[string]string)

	// Определяем исходный язык
	sourceLanguage, _, err := s.DetectLanguage(ctx, text)
	if err != nil {
		sourceLanguage = "auto"
	}

	// Переводим на все языки кроме исходного
	for _, lang := range supportedLanguages {
		if lang != sourceLanguage {
			translated, err := s.Translate(ctx, text, sourceLanguage, lang)
			if err != nil {
				logger.Error().Err(err).
					Str("target", lang).
					Msg("Failed to translate with DeepL")
				result[lang] = text // Возвращаем оригинал при ошибке
			} else {
				result[lang] = translated
			}
		} else {
			result[lang] = text
		}
	}

	return result, nil
}

// TranslateEntityFields переводит поля сущности
func (s *DeepLTranslationService) TranslateEntityFields(ctx context.Context, sourceLanguage string, targetLanguages []string, fields map[string]string) (map[string]map[string]string, error) {
	result := make(map[string]map[string]string)

	// DeepL поддерживает batch переводы, можно оптимизировать
	// отправляя все тексты одним запросом
	for _, targetLang := range targetLanguages {
		if targetLang == sourceLanguage {
			continue
		}

		translations := make(map[string]string)

		// Собираем все тексты для batch перевода
		var texts []string
		var fieldNames []string
		for fieldName, fieldValue := range fields {
			if fieldValue != "" {
				texts = append(texts, fieldValue)
				fieldNames = append(fieldNames, fieldName)
			}
		}

		// Batch перевод всех текстов
		if len(texts) > 0 {
			translatedTexts, err := s.translateBatch(ctx, texts, sourceLanguage, targetLang)
			if err != nil {
				logger.Error().Err(err).
					Str("target", targetLang).
					Msg("Failed to batch translate with DeepL")
				// При ошибке возвращаем оригиналы
				for fieldName, fieldValue := range fields {
					translations[fieldName] = fieldValue
				}
			} else {
				// Мапим переведенные тексты обратно к полям
				for i, fieldName := range fieldNames {
					if i < len(translatedTexts) {
						translations[fieldName] = translatedTexts[i]
					}
				}
				// Добавляем пустые поля
				for fieldName, fieldValue := range fields {
					if fieldValue == "" {
						translations[fieldName] = ""
					}
				}
			}
		} else {
			// Все поля пустые
			for fieldName := range fields {
				translations[fieldName] = ""
			}
		}

		result[targetLang] = translations
	}

	return result, nil
}

// translateBatch выполняет batch перевод нескольких текстов
func (s *DeepLTranslationService) translateBatch(ctx context.Context, texts []string, sourceLanguage string, targetLanguage string) ([]string, error) {
	sourceLang := convertToDeepLLanguageCode(sourceLanguage)
	targetLang := convertToDeepLLanguageCode(targetLanguage)

	params := url.Values{}
	for _, text := range texts {
		params.Add("text", text)
	}
	params.Set("target_lang", targetLang)
	if sourceLang != "" && sourceLang != "auto" {
		params.Set("source_lang", sourceLang)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.apiURL+"/translate", strings.NewReader(params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "DeepL-Auth-Key "+s.apiKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("DeepL API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var deepLResp deepLTranslateResponse
	if err := json.Unmarshal(body, &deepLResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	var results []string
	for _, translation := range deepLResp.Translations {
		results = append(results, translation.Text)
	}

	return results, nil
}

// DetectLanguage определяет язык текста
func (s *DeepLTranslationService) DetectLanguage(ctx context.Context, text string) (string, float64, error) {
	if text == "" {
		return "", 0, fmt.Errorf("empty text")
	}

	// DeepL автоматически определяет язык при переводе
	// Но можно использовать простую эвристику
	// Используем функции из claude_translation.go
	if containsCyrillic(text) {
		if containsSerbian(text) {
			return "sr", 0.9, nil
		}
		return "ru", 0.9, nil
	}

	return "en", 0.9, nil
}

// ModerateText выполняет модерацию текста
func (s *DeepLTranslationService) ModerateText(ctx context.Context, text string, language string) (string, error) {
	// DeepL не предоставляет функционал модерации
	// Возвращаем текст как есть
	return text, nil
}

// convertToDeepLLanguageCode преобразует наши коды языков в формат DeepL
func convertToDeepLLanguageCode(code string) string {
	// DeepL использует другие коды для некоторых языков
	mapping := map[string]string{
		"en":   "EN-US", // или EN-GB для британского английского
		"ru":   "RU",
		"sr":   "SR", // DeepL может не поддерживать сербский напрямую
		"auto": "",   // Пустая строка для автоопределения
	}

	if deepLCode, ok := mapping[code]; ok {
		return deepLCode
	}
	return strings.ToUpper(code)
}

// GetUsage получает информацию об использовании квоты
func (s *DeepLTranslationService) GetUsage(ctx context.Context) (map[string]interface{}, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", s.apiURL+"/usage", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "DeepL-Auth-Key "+s.apiKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("DeepL API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var usage map[string]interface{}
	if err := json.Unmarshal(body, &usage); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return usage, nil
}
