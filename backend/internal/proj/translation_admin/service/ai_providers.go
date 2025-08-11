package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// AIProvider представляет провайдера AI переводов
type AIProvider interface {
	Translate(ctx context.Context, text string, sourceLang, targetLang string) (string, float64, error)
	IsConfigured() bool
	GetName() string
}

// ClaudeProvider - провайдер переводов через Anthropic Claude
type ClaudeProvider struct {
	apiKey string
	model  string
}

// NewClaudeProvider создает новый провайдер Claude
func NewClaudeProvider() *ClaudeProvider {
	return &ClaudeProvider{
		apiKey: os.Getenv("CLAUDE_API_KEY"),
		model:  "claude-3-haiku-20240307", // используем haiku для быстрых переводов
	}
}

// IsConfigured проверяет настроен ли провайдер
func (p *ClaudeProvider) IsConfigured() bool {
	return p.apiKey != ""
}

// GetName возвращает имя провайдера
func (p *ClaudeProvider) GetName() string {
	return "Claude (Anthropic)"
}

// Translate выполняет перевод текста
func (p *ClaudeProvider) Translate(ctx context.Context, text string, sourceLang, targetLang string) (string, float64, error) {
	if !p.IsConfigured() {
		return "", 0, fmt.Errorf("claude API key not configured")
	}

	// Маппинг языковых кодов
	langMap := map[string]string{
		"en": "English",
		"ru": "Russian",
		"sr": "Serbian",
	}

	sourceLangName := langMap[sourceLang]
	targetLangName := langMap[targetLang]

	if sourceLangName == "" || targetLangName == "" {
		return "", 0, fmt.Errorf("unsupported language code")
	}

	// Формируем запрос к Claude API
	prompt := fmt.Sprintf(
		`Translate the following text from %s to %s. 
		Provide ONLY the translation without any explanations or additional text.
		Maintain the same tone, style, and formatting as the original.
		
		Text to translate:
		%s`,
		sourceLangName,
		targetLangName,
		text,
	)

	requestBody := map[string]interface{}{
		"model":      p.model,
		"max_tokens": 1024,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", 0, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", p.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("claude API error (status %d): %s", resp.StatusCode, string(body))
	}

	var response struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
		Usage struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return "", 0, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(response.Content) == 0 {
		return "", 0, fmt.Errorf("no translation returned")
	}

	translation := strings.TrimSpace(response.Content[0].Text)

	// Высокая уверенность для Claude
	confidence := 0.95

	return translation, confidence, nil
}

// OpenAIProvider - провайдер переводов через OpenAI GPT
type OpenAIProvider struct {
	apiKey string
	model  string
}

// NewOpenAIProvider создает новый провайдер OpenAI
func NewOpenAIProvider() *OpenAIProvider {
	return &OpenAIProvider{
		apiKey: os.Getenv("OPENAI_API_KEY"),
		model:  "gpt-4-turbo-preview",
	}
}

// IsConfigured проверяет настроен ли провайдер
func (o *OpenAIProvider) IsConfigured() bool {
	return o.apiKey != ""
}

// GetName возвращает имя провайдера
func (o *OpenAIProvider) GetName() string {
	return "OpenAI GPT-4"
}

// Translate выполняет перевод текста через OpenAI
func (o *OpenAIProvider) Translate(ctx context.Context, text string, sourceLang, targetLang string) (string, float64, error) {
	if !o.IsConfigured() {
		return "", 0, fmt.Errorf("OpenAI API key not configured")
	}

	// Маппинг языковых кодов
	langMap := map[string]string{
		"en": "English",
		"ru": "Russian",
		"sr": "Serbian",
	}

	sourceLangName := langMap[sourceLang]
	targetLangName := langMap[targetLang]

	if sourceLangName == "" || targetLangName == "" {
		return "", 0, fmt.Errorf("unsupported language code")
	}

	prompt := fmt.Sprintf(
		"Translate from %s to %s: %s",
		sourceLangName,
		targetLangName,
		text,
	)

	requestBody := map[string]interface{}{
		"model": o.model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are a professional translator. Provide only the translation without any explanations.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature": 0.3,
		"max_tokens":  1024,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", 0, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", o.apiKey))

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("OpenAI API error (status %d): %s", resp.StatusCode, string(body))
	}

	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return "", 0, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", 0, fmt.Errorf("no translation returned")
	}

	translation := strings.TrimSpace(response.Choices[0].Message.Content)
	confidence := 0.92

	return translation, confidence, nil
}

// GetAIProvider возвращает сконфигурированного провайдера AI
func GetAIProvider(providerName string) AIProvider {
	switch providerName {
	case "claude", "anthropic":
		return NewClaudeProvider()
	case "openai", "gpt":
		return NewOpenAIProvider()
	case "deepl":
		return NewDeepLProvider()
	default:
		// По умолчанию пробуем провайдеров в порядке приоритета
		// 1. DeepL (обычно лучшее качество для многих языков)
		deepl := NewDeepLProvider()
		if deepl.IsConfigured() {
			return deepl
		}
		// 2. Claude
		claude := NewClaudeProvider()
		if claude.IsConfigured() {
			return claude
		}
		// 3. OpenAI
		return NewOpenAIProvider()
	}
}

// GetAvailableProviders возвращает список доступных провайдеров
func GetAvailableProviders() []map[string]interface{} {
	providers := []map[string]interface{}{}

	// Проверяем DeepL
	deepl := NewDeepLProvider()
	providers = append(providers, map[string]interface{}{
		"id":         "deepl",
		"name":       deepl.GetName(),
		"type":       "deepl",
		"model":      "DeepL Pro",
		"enabled":    deepl.IsConfigured(),
		"configured": deepl.IsConfigured(),
	})

	// Проверяем Claude
	claude := NewClaudeProvider()
	providers = append(providers, map[string]interface{}{
		"id":         "claude",
		"name":       claude.GetName(),
		"type":       "anthropic",
		"model":      "claude-3-opus",
		"enabled":    claude.IsConfigured(),
		"configured": claude.IsConfigured(),
	})

	// Проверяем OpenAI
	openai := NewOpenAIProvider()
	providers = append(providers, map[string]interface{}{
		"id":         "openai",
		"name":       openai.GetName(),
		"type":       "openai",
		"model":      "gpt-4-turbo-preview",
		"enabled":    openai.IsConfigured(),
		"configured": openai.IsConfigured(),
	})

	// Добавляем заглушку для Google Translate (если будет добавлен)
	providers = append(providers, map[string]interface{}{
		"id":         "google",
		"name":       "Google Translate",
		"type":       "google",
		"enabled":    false,
		"configured": false,
	})

	return providers
}
