// backend/internal/proj/marketplace/service/claude_translation.go
package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"backend/internal/logger"
)

// ClaudeTranslationService предоставляет функционал перевода через Claude AI API
type ClaudeTranslationService struct {
	apiKey     string
	httpClient *http.Client
}

// NewClaudeTranslationService создает новый экземпляр сервиса перевода Claude AI
func NewClaudeTranslationService(apiKey string) (*ClaudeTranslationService, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("claude API key is required")
	}

	return &ClaudeTranslationService{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// claudeRequest представляет запрос к Claude API
type claudeRequest struct {
	Model     string          `json:"model"`
	MaxTokens int             `json:"max_tokens"`
	Messages  []claudeMessage `json:"messages"`
}

type claudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// claudeResponse представляет ответ от Claude API
type claudeResponse struct {
	Content []struct {
		Text string `json:"text"`
		Type string `json:"type"`
	} `json:"content"`
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// Translate переводит текст с одного языка на другой
func (s *ClaudeTranslationService) Translate(ctx context.Context, text string, sourceLanguage string, targetLanguage string) (string, error) {
	if text == "" {
		return "", nil
	}

	// Подготавливаем промпт для Claude
	prompt := fmt.Sprintf(`Translate the following text from %s to %s. 
Return ONLY the translated text without any explanations or additional content.
Do not add quotes or any formatting.

Text to translate:
%s`, getLanguageName(sourceLanguage), getLanguageName(targetLanguage), text)

	// Создаем запрос к Claude API
	requestBody := claudeRequest{
		Model:     "claude-3-haiku-20240307", // Используем самую быструю и дешевую модель для переводов
		MaxTokens: 1024,
		Messages: []claudeMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Создаем HTTP запрос
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Устанавливаем заголовки
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", s.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	// Выполняем запрос
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Читаем ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("claude API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	// Парсим ответ
	var claudeResp claudeResponse
	if err := json.Unmarshal(body, &claudeResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Проверяем на ошибки в ответе
	if claudeResp.Error.Message != "" {
		return "", fmt.Errorf("claude API error: %s", claudeResp.Error.Message)
	}

	// Извлекаем переведенный текст
	if len(claudeResp.Content) > 0 && claudeResp.Content[0].Type == attributeTypeText {
		translatedText := strings.TrimSpace(claudeResp.Content[0].Text)
		logger.Info().
			Str("source", sourceLanguage).
			Str("target", targetLanguage).
			Int("source_len", len(text)).
			Int("translated_len", len(translatedText)).
			Msg("Claude translation completed")
		return translatedText, nil
	}

	return "", fmt.Errorf("unexpected Claude API response format")
}

// TranslateWithContext переводит текст с учетом контекста
func (s *ClaudeTranslationService) TranslateWithContext(ctx context.Context, text string, sourceLanguage string, targetLanguage string, context string, fieldName string) (string, error) {
	if text == "" {
		return "", nil
	}

	// Подготавливаем промпт с контекстом
	prompt := fmt.Sprintf(`You are translating content for an e-commerce marketplace.
Context: %s
Field type: %s

Translate the following text from %s to %s.
Return ONLY the translated text without any explanations.
Keep the same tone and style as the original.

Text to translate:
%s`, context, fieldName, getLanguageName(sourceLanguage), getLanguageName(targetLanguage), text)

	// Используем основной метод Translate с модифицированным промптом
	requestBody := claudeRequest{
		Model:     "claude-3-haiku-20240307",
		MaxTokens: 1024,
		Messages: []claudeMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", s.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("claude API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var claudeResp claudeResponse
	if err := json.Unmarshal(body, &claudeResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if claudeResp.Error.Message != "" {
		return "", fmt.Errorf("claude API error: %s", claudeResp.Error.Message)
	}

	if len(claudeResp.Content) > 0 && claudeResp.Content[0].Type == attributeTypeText {
		return strings.TrimSpace(claudeResp.Content[0].Text), nil
	}

	return "", fmt.Errorf("unexpected Claude API response format")
}

// TranslateToAllLanguages переводит текст на все поддерживаемые языки
func (s *ClaudeTranslationService) TranslateToAllLanguages(ctx context.Context, text string) (map[string]string, error) {
	supportedLanguages := []string{"en", "ru", "sr"}
	result := make(map[string]string)

	// Определяем исходный язык
	sourceLanguage, _, err := s.DetectLanguage(ctx, text)
	if err != nil {
		sourceLanguage = languageAuto
	}

	// Переводим на все языки кроме исходного
	for _, lang := range supportedLanguages {
		if lang != sourceLanguage {
			translated, err := s.Translate(ctx, text, sourceLanguage, lang)
			if err != nil {
				logger.Error().Err(err).
					Str("target", lang).
					Msg("Failed to translate with Claude")
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
func (s *ClaudeTranslationService) TranslateEntityFields(ctx context.Context, sourceLanguage string, targetLanguages []string, fields map[string]string) (map[string]map[string]string, error) {
	result := make(map[string]map[string]string)

	for _, targetLang := range targetLanguages {
		if targetLang == sourceLanguage {
			continue
		}

		translations := make(map[string]string)
		for fieldName, fieldValue := range fields {
			if fieldValue == "" {
				translations[fieldName] = ""
				continue
			}

			translated, err := s.TranslateWithContext(ctx, fieldValue, sourceLanguage, targetLang, "marketplace entity", fieldName)
			if err != nil {
				logger.Error().Err(err).
					Str("field", fieldName).
					Str("target", targetLang).
					Msg("Failed to translate field with Claude")
				translations[fieldName] = fieldValue
			} else {
				translations[fieldName] = translated
			}
		}
		result[targetLang] = translations
	}

	return result, nil
}

// DetectLanguage определяет язык текста
func (s *ClaudeTranslationService) DetectLanguage(ctx context.Context, text string) (string, float64, error) {
	if text == "" {
		return "", 0, fmt.Errorf("empty text")
	}

	// Простая эвристика для определения языка
	// Можно улучшить, используя Claude для определения языка
	if containsCyrillic(text) {
		if containsSerbian(text) {
			return "sr", 0.9, nil
		}
		return "ru", 0.9, nil
	}

	return "en", 0.9, nil
}

// ModerateText выполняет модерацию текста
func (s *ClaudeTranslationService) ModerateText(ctx context.Context, text string, language string) (string, error) {
	// Claude имеет встроенную модерацию, но можно добавить дополнительную проверку
	prompt := fmt.Sprintf(`Check if the following text contains inappropriate content (hate speech, violence, adult content).
If the text is appropriate, return it as is.
If inappropriate, return a cleaned version or "[MODERATED]".

Text in %s:
%s`, getLanguageName(language), text)

	requestBody := claudeRequest{
		Model:     "claude-3-haiku-20240307",
		MaxTokens: 512,
		Messages: []claudeMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return text, nil // При ошибке возвращаем оригинальный текст
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return text, nil
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", s.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return text, nil
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return text, nil
	}

	if resp.StatusCode != http.StatusOK {
		return text, nil
	}

	var claudeResp claudeResponse
	if err := json.Unmarshal(body, &claudeResp); err != nil {
		return text, nil
	}

	if len(claudeResp.Content) > 0 && claudeResp.Content[0].Type == attributeTypeText {
		return strings.TrimSpace(claudeResp.Content[0].Text), nil
	}

	return text, nil
}

// Вспомогательные функции
func containsCyrillic(text string) bool {
	for _, r := range text {
		if (r >= 'А' && r <= 'я') || r == 'Ё' || r == 'ё' {
			return true
		}
	}
	return false
}

func containsSerbian(text string) bool {
	serbianSpecific := []string{"ђ", "ј", "љ", "њ", "ћ", "џ", "Ђ", "Ј", "Љ", "Њ", "Ћ", "Џ"}
	textLower := strings.ToLower(text)
	for _, char := range serbianSpecific {
		if strings.Contains(textLower, strings.ToLower(char)) {
			return true
		}
	}
	return false
}

func getLanguageName(code string) string {
	languages := map[string]string{
		"en":   "English",
		"ru":   "Russian",
		"sr":   "Serbian",
		"auto": "auto-detect",
	}
	if name, ok := languages[code]; ok {
		return name
	}
	return code
}

// TranslateWithToneModeration переводит текст с опциональным смягчением тона
func (s *ClaudeTranslationService) TranslateWithToneModeration(
	ctx context.Context,
	text string,
	sourceLanguage string,
	targetLanguage string,
	moderateTone bool,
) (string, error) {
	if text == "" {
		return "", nil
	}

	var prompt string

	if moderateTone {
		// Промпт с модерацией тона
		prompt = fmt.Sprintf(`Translate the following text from %s to %s.

IMPORTANT: If the text contains profanity, offensive language, or
aggressive tone, translate it to a polite, respectful equivalent
while preserving the general meaning and emotional intensity.

Examples:
- "What the fuck?" → "What's going on?" (surprised, confused)
- "This is fucking great!" → "This is really great!" (very excited)
- "Stop being an asshole" → "Please be more considerate" (frustrated)

Return ONLY the translated text without any explanations or additional content.
Do not add quotes or any formatting.

Text to translate:
%s`, getLanguageName(sourceLanguage), getLanguageName(targetLanguage), text)
	} else {
		// Обычный промпт без модерации
		prompt = fmt.Sprintf(`Translate the following text from %s to %s.
Return ONLY the translated text without any explanations or additional content.
Do not add quotes or any formatting.

Text to translate:
%s`, getLanguageName(sourceLanguage), getLanguageName(targetLanguage), text)
	}

	// Создаем запрос к Claude API
	requestBody := claudeRequest{
		Model:     "claude-3-haiku-20240307",
		MaxTokens: 1024,
		Messages: []claudeMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", s.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call Claude API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("claude API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var claudeResp claudeResponse
	if err := json.Unmarshal(body, &claudeResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if claudeResp.Error.Message != "" {
		return "", fmt.Errorf("claude API error: %s", claudeResp.Error.Message)
	}

	if len(claudeResp.Content) > 0 && claudeResp.Content[0].Type == "text" {
		return strings.TrimSpace(claudeResp.Content[0].Text), nil
	}

	return "", fmt.Errorf("unexpected Claude API response format")
}
