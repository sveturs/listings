package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// DeepLProvider - провайдер переводов через DeepL API
type DeepLProvider struct {
	apiKey   string
	endpoint string
}

// NewDeepLProvider создает новый провайдер DeepL
func NewDeepLProvider() *DeepLProvider {
	apiKey := os.Getenv("DEEPL_API_KEY")

	// Определяем эндпоинт на основе типа ключа
	endpoint := "https://api.deepl.com/v2/translate"
	if strings.HasSuffix(apiKey, ":fx") {
		// Free API key
		endpoint = "https://api-free.deepl.com/v2/translate"
	}

	return &DeepLProvider{
		apiKey:   apiKey,
		endpoint: endpoint,
	}
}

// IsConfigured проверяет настроен ли провайдер
func (p *DeepLProvider) IsConfigured() bool {
	return p.apiKey != ""
}

// GetName возвращает имя провайдера
func (p *DeepLProvider) GetName() string {
	return "DeepL"
}

// Translate выполняет перевод текста через DeepL
func (p *DeepLProvider) Translate(ctx context.Context, text string, sourceLang, targetLang string) (string, float64, error) {
	if !p.IsConfigured() {
		return "", 0, fmt.Errorf("DeepL API key not configured")
	}

	// Маппинг языковых кодов для DeepL
	langMap := map[string]string{
		"en": "EN",
		"ru": "RU",
		"sr": "SR", // DeepL поддерживает сербский
	}

	sourceCode := strings.ToUpper(sourceLang)
	targetCode := langMap[targetLang]

	if targetCode == "" {
		return "", 0, fmt.Errorf("unsupported target language: %s", targetLang)
	}

	// Формируем запрос с правильным URL-кодированием
	formData := url.Values{}
	formData.Set("auth_key", p.apiKey)
	formData.Set("text", text)
	formData.Set("target_lang", targetCode)

	// Добавляем source_lang только если он указан и не auto
	if sourceCode != "" && sourceCode != "AUTO" {
		sourceCodeMapped := langMap[sourceLang]
		if sourceCodeMapped != "" {
			formData.Set("source_lang", sourceCodeMapped)
		}
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.endpoint, strings.NewReader(formData.Encode()))
	if err != nil {
		return "", 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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
		return "", 0, fmt.Errorf("DeepL API error (status %d): %s", resp.StatusCode, string(body))
	}

	var response struct {
		Translations []struct {
			DetectedSourceLanguage string `json:"detected_source_language"`
			Text                   string `json:"text"`
		} `json:"translations"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return "", 0, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(response.Translations) == 0 {
		return "", 0, fmt.Errorf("no translation returned")
	}

	translation := strings.TrimSpace(response.Translations[0].Text)

	// DeepL обычно дает очень высокое качество перевода
	confidence := 0.98

	return translation, confidence, nil
}
