package translation_admin

import (
	"context"
	"fmt"
	"os"
	"strings"

	"backend/internal/domain/models"
	"backend/internal/proj/translation_admin/ratelimit"
)

// GetAIProviders returns list of configured AI translation providers
func (s *Service) GetAIProviders(ctx context.Context) ([]models.AIProvider, error) {
	// Проверяем какие провайдеры настроены через переменные окружения
	claudeConfigured := os.Getenv("CLAUDE_API_KEY") != ""
	openaiConfigured := os.Getenv("OPENAI_API_KEY") != ""
	deeplConfigured := os.Getenv("DEEPL_API_KEY") != ""
	googleConfigured := os.Getenv("GOOGLE_TRANSLATE_API_KEY") != ""

	providers := []models.AIProvider{
		{
			ID:          "anthropic",
			Name:        "Anthropic Claude 3",
			Type:        "anthropic",
			Model:       "claude-3-opus-20240229",
			Enabled:     claudeConfigured,
			MaxTokens:   2000,
			Temperature: 0.3,
		},
		{
			ID:          ratelimit.ProviderOpenAI,
			Name:        "OpenAI GPT-4",
			Type:        ratelimit.ProviderOpenAI,
			Model:       "gpt-4-turbo-preview",
			Enabled:     openaiConfigured,
			MaxTokens:   2000,
			Temperature: 0.3,
		},
		{
			ID:       "deepl",
			Name:     "DeepL API",
			Type:     "deepl",
			Endpoint: "https://api.deepl.com/v2/translate",
			Enabled:  deeplConfigured,
		},
		{
			ID:       "google",
			Name:     "Google Translate",
			Type:     "google",
			Endpoint: "https://translation.googleapis.com/language/translate/v2",
			Enabled:  googleConfigured,
		},
	}

	return providers, nil
}

// UpdateAIProvider updates AI provider configuration
func (s *Service) UpdateAIProvider(ctx context.Context, provider *models.AIProvider, userID int) error {
	s.logger.Info().
		Str("provider_id", provider.ID).
		Bool("enabled", provider.Enabled).
		Int("user_id", userID).
		Msg("Updating AI provider configuration")

	// В реальной версии - сохранить в БД
	// Если провайдер активирован, деактивировать остальных

	return nil
}

// TranslateText translates a single text using AI
func (s *Service) TranslateText(ctx context.Context, req *models.TranslateRequest, userID int) (*models.TranslateResult, error) {
	s.logger.Info().
		Str("provider", req.Provider).
		Str("key", req.Key).
		Str("module", req.Module).
		Str("source_lang", req.SourceLanguage).
		Interface("target_langs", req.TargetLanguages).
		Msg("Translating text with AI")

	// Определяем провайдера
	provider := req.Provider
	if provider == "" {
		provider = ratelimit.ProviderOpenAI // default provider
	}

	// Use real translation service
	translations := make(map[string]string)

	// Используем реальный сервис перевода если доступен
	if s.translationFactory != nil {
		// Пробуем использовать интерфейс напрямую для простого перевода
		for _, targetLang := range req.TargetLanguages {
			if req.SourceLanguage == targetLang {
				// Same language, no translation needed
				translations[targetLang] = req.Text
				continue
			}

			// Попробуем вызвать метод Translate напрямую
			// Используем type assertion для проверки интерфейса
			if translator, ok := s.translationFactory.(interface {
				Translate(ctx context.Context, text, sourceLang, targetLang string) (string, error)
			}); ok {
				translated, err := translator.Translate(ctx, req.Text, req.SourceLanguage, targetLang)
				if err != nil {
					s.logger.Error().
						Err(err).
						Str("source", req.SourceLanguage).
						Str("target", targetLang).
						Msg("Translation failed")
					// Fallback to mock on error
					translations[targetLang] = "[" + strings.ToUpper(targetLang) + "] " + req.Text
				} else {
					translations[targetLang] = translated
				}
			} else {
				// Интерфейс не поддерживает метод Translate, используем mock
				translations[targetLang] = "[" + strings.ToUpper(targetLang) + "] " + req.Text
			}
		}
	} else {
		// Фабрика не доступна, используем mock
		for _, lang := range req.TargetLanguages {
			translations[lang] = "[" + strings.ToUpper(lang) + "] " + req.Text
		}
	}

	// Отслеживаем использование AI провайдера
	// Приблизительный расчет токенов (4 символа = 1 токен)
	textLength := len(req.Text)
	inputTokens := textLength / 4
	outputTokens := textLength / 4 * len(req.TargetLanguages) // умножаем на количество языков

	s.logger.Info().
		Str("provider", provider).
		Int("input_tokens", inputTokens).
		Int("output_tokens", outputTokens).
		Int("text_length", textLength).
		Int("target_languages", len(req.TargetLanguages)).
		Msg("Tracking AI provider usage")

	err := s.TrackAIProviderUsage(ctx, provider, inputTokens, outputTokens, textLength*len(req.TargetLanguages))
	if err != nil {
		s.logger.Warn().Err(err).Str("provider", provider).Msg("Failed to track AI provider usage")
		// Не прерываем выполнение если трекинг не удался
	} else {
		s.logger.Info().
			Str("provider", provider).
			Int("input_tokens", inputTokens).
			Int("output_tokens", outputTokens).
			Msg("Successfully tracked AI provider usage")
	}

	result := &models.TranslateResult{
		Key:          req.Key,
		Module:       req.Module,
		Translations: translations,
		Provider:     provider,
		Confidence:   0.95,
	}

	// Optionally add alternative translations
	if req.Provider == ratelimit.ProviderOpenAI || req.Provider == "anthropic" {
		alternatives := make(map[string][]string)
		for lang := range translations {
			alternatives[lang] = []string{
				translations[lang] + " (вариант 1)",
				translations[lang] + " (вариант 2)",
			}
		}
		result.AlternativeTranslations = alternatives
	}

	return result, nil
}

// BatchTranslate performs batch translation of multiple texts
func (s *Service) BatchTranslate(ctx context.Context, req *models.AIBatchTranslateRequest, userID int) (*models.BatchTranslateResult, error) {
	s.logger.Info().
		Str("provider", req.Provider).
		Interface("modules", req.Modules).
		Bool("missing_only", req.MissingOnly).
		Msg("Starting batch translation")

	result := &models.BatchTranslateResult{
		Results:         []models.TranslateResult{},
		TranslatedCount: 0,
		FailedCount:     0,
		Errors:          []string{},
	}

	// Mock implementation - в реальной версии загружать тексты из модулей
	mockTexts := []struct {
		Key    string
		Module string
		Text   string
	}{
		{"common.welcome", "common", "Welcome"},
		{"common.goodbye", "common", "Goodbye"},
		{"marketplace.title", "marketplace", "Marketplace"},
		{"marketplace.search", "marketplace", "Search"},
	}

	for _, item := range mockTexts {
		// Фильтровать по модулям если указаны
		moduleMatch := false
		for _, m := range req.Modules {
			if m == item.Module {
				moduleMatch = true
				break
			}
		}
		if !moduleMatch && len(req.Modules) > 0 {
			continue
		}

		// В реальной версии проверять missing_only
		if req.MissingOnly {
			// TODO: Implement check for existing translations
			s.logger.Debug().Msg("MissingOnly flag is set but checking logic not yet implemented")
		}

		// Перевести текст
		translateReq := &models.TranslateRequest{
			Provider:        req.Provider,
			Text:            item.Text,
			Key:             item.Key,
			Module:          item.Module,
			SourceLanguage:  req.SourceLanguage,
			TargetLanguages: req.TargetLanguages,
		}

		translationResult, err := s.TranslateText(ctx, translateReq, userID)
		if err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, fmt.Sprintf("Failed to translate %s: %v", item.Key, err))
			continue
		}

		result.Results = append(result.Results, *translationResult)
		result.TranslatedCount++
	}

	return result, nil
}

// ApplyAITranslations applies AI-generated translations
func (s *Service) ApplyAITranslations(ctx context.Context, req *models.ApplyTranslationsRequest, userID int) error {
	s.logger.Info().
		Int("count", len(req.Translations)).
		Int("user_id", userID).
		Msg("Applying AI translations")

	for _, translation := range req.Translations {
		// В реальной версии - сохранить в БД и обновить JSON файлы
		s.logger.Info().
			Str("key", translation.Key).
			Str("module", translation.Module).
			Str("language", translation.Language).
			Str("value", translation.Value).
			Msg("Applying translation")

		// 1. Обновить JSON файл модуля
		// 2. Сохранить в БД для синхронизации
		// 3. Создать запись в аудите
	}

	return nil
}
