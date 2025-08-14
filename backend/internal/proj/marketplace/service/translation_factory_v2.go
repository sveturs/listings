// backend/internal/proj/marketplace/service/translation_factory_v2.go
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"backend/internal/domain/models"
	"backend/internal/logger"
	"backend/internal/storage"
)

// Убедимся, что TranslationServiceFactoryV2 реализует интерфейс TranslationFactoryInterface
var _ TranslationFactoryInterface = (*TranslationServiceFactoryV2)(nil)

// TranslationServiceFactoryV2 создает и управляет различными сервисами перевода
type TranslationServiceFactoryV2 struct {
	googleService   *GoogleTranslationService
	openAIService   *TranslationService
	claudeService   *ClaudeTranslationService
	deeplService    *DeepLTranslationService
	defaultProvider TranslationProvider
	fallbackChain   []TranslationProvider // Цепочка резервных провайдеров
	mutex           sync.RWMutex
}

// NewTranslationServiceFactoryV2 создает новый экземпляр фабрики сервисов перевода с поддержкой 4 провайдеров
func NewTranslationServiceFactoryV2(config struct {
	GoogleAPIKey    string
	OpenAIAPIKey    string
	ClaudeAPIKey    string
	DeepLAPIKey     string
	DeepLUseFreeAPI bool
}, storage storage.Storage,
) (*TranslationServiceFactoryV2, error) {
	factory := &TranslationServiceFactoryV2{
		mutex: sync.RWMutex{},
	}

	// Инициализируем доступные сервисы и формируем цепочку fallback
	var availableProviders []TranslationProvider

	// 1. DeepL - лучшее качество для европейских языков
	if config.DeepLAPIKey != "" {
		deeplService, err := NewDeepLTranslationService(config.DeepLAPIKey, config.DeepLUseFreeAPI)
		if err != nil {
			logger.Warn().Err(err).Msg("Не удалось создать сервис DeepL")
		} else {
			factory.deeplService = deeplService
			availableProviders = append(availableProviders, DeepL)
			logger.Info().Msg("DeepL сервис инициализирован")
		}
	}

	// 2. Claude AI - хорошее качество и контекстное понимание
	if config.ClaudeAPIKey != "" {
		claudeService, err := NewClaudeTranslationService(config.ClaudeAPIKey)
		if err != nil {
			logger.Warn().Err(err).Msg("Не удалось создать сервис Claude AI")
		} else {
			factory.claudeService = claudeService
			availableProviders = append(availableProviders, ClaudeAI)
			logger.Info().Msg("Claude AI сервис инициализирован")
		}
	}

	// 3. Google Translate - широкая поддержка языков
	if config.GoogleAPIKey != "" {
		googleService, err := NewGoogleTranslationService(config.GoogleAPIKey, storage)
		if err != nil {
			logger.Warn().Err(err).Msg("Не удалось создать сервис Google Translate")
		} else {
			factory.googleService = googleService
			availableProviders = append(availableProviders, GoogleTranslate)
			logger.Info().Msg("Google Translate сервис инициализирован")
		}
	}

	// 4. OpenAI - резервный вариант
	if config.OpenAIAPIKey != "" {
		openAIService, err := NewTranslationService(config.OpenAIAPIKey)
		if err != nil {
			logger.Warn().Err(err).Msg("Не удалось создать сервис OpenAI")
		} else {
			// Добавляем доступ к хранилищу для OpenAI сервиса
			if storage != nil {
				openAIService.storage = storage
			}
			factory.openAIService = openAIService
			availableProviders = append(availableProviders, OpenAI)
			logger.Info().Msg("OpenAI сервис инициализирован")
		}
	}

	// Проверяем, что хотя бы один сервис доступен
	if len(availableProviders) == 0 {
		return nil, fmt.Errorf("не удалось инициализировать ни один сервис перевода")
	}

	// Устанавливаем провайдер по умолчанию и цепочку fallback
	factory.defaultProvider = availableProviders[0]
	factory.fallbackChain = availableProviders

	logger.Info().
		Str("defaultProvider", string(factory.defaultProvider)).
		Int("totalProviders", len(availableProviders)).
		Interface("providers", availableProviders).
		Msg("Фабрика сервисов перевода V2 создана")

	return factory, nil
}

// GetTranslationService возвращает сервис перевода по запрошенному провайдеру
func (f *TranslationServiceFactoryV2) GetTranslationService(provider TranslationProvider) (TranslationServiceInterface, error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	switch provider {
	case GoogleTranslate:
		if f.googleService == nil {
			return nil, fmt.Errorf("сервис Google Translate недоступен")
		}
		return f.googleService, nil
	case OpenAI:
		if f.openAIService == nil {
			return nil, fmt.Errorf("сервис OpenAI недоступен")
		}
		return f.openAIService, nil
	case ClaudeAI:
		if f.claudeService == nil {
			return nil, fmt.Errorf("сервис Claude AI недоступен")
		}
		return f.claudeService, nil
	case DeepL:
		if f.deeplService == nil {
			return nil, fmt.Errorf("сервис DeepL недоступен")
		}
		return f.deeplService, nil
	case Manual:
		return nil, fmt.Errorf("ручной перевод не поддерживается в автоматическом режиме")
	default:
		return nil, fmt.Errorf("неизвестный провайдер перевода: %s", provider)
	}
}

// GetDefaultProvider возвращает провайдер перевода по умолчанию
func (f *TranslationServiceFactoryV2) GetDefaultProvider() TranslationProvider {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	return f.defaultProvider
}

// SetDefaultProvider устанавливает провайдер перевода по умолчанию
func (f *TranslationServiceFactoryV2) SetDefaultProvider(provider TranslationProvider) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	// Проверяем доступность провайдера
	switch provider {
	case GoogleTranslate:
		if f.googleService == nil {
			return fmt.Errorf("сервис Google Translate недоступен")
		}
	case OpenAI:
		if f.openAIService == nil {
			return fmt.Errorf("сервис OpenAI недоступен")
		}
	case ClaudeAI:
		if f.claudeService == nil {
			return fmt.Errorf("сервис Claude AI недоступен")
		}
	case DeepL:
		if f.deeplService == nil {
			return fmt.Errorf("сервис DeepL недоступен")
		}
	case Manual:
		return fmt.Errorf("ручной перевод не может быть установлен как провайдер по умолчанию")
	default:
		return fmt.Errorf("неизвестный провайдер перевода: %s", provider)
	}

	f.defaultProvider = provider

	// Перестраиваем цепочку fallback с новым провайдером в начале
	newChain := []TranslationProvider{provider}
	for _, p := range f.fallbackChain {
		if p != provider {
			newChain = append(newChain, p)
		}
	}
	f.fallbackChain = newChain

	logger.Info().
		Str("provider", string(provider)).
		Interface("fallbackChain", f.fallbackChain).
		Msg("Установлен провайдер перевода по умолчанию")
	return nil
}

// GetAvailableProviders возвращает список доступных провайдеров перевода
func (f *TranslationServiceFactoryV2) GetAvailableProviders() []TranslationProvider {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	providers := make([]TranslationProvider, 0, 4)

	if f.googleService != nil {
		providers = append(providers, GoogleTranslate)
	}
	if f.openAIService != nil {
		providers = append(providers, OpenAI)
	}
	if f.claudeService != nil {
		providers = append(providers, ClaudeAI)
	}
	if f.deeplService != nil {
		providers = append(providers, DeepL)
	}

	return providers
}

// GetTranslationCount возвращает количество выполненных переводов для указанного провайдера
func (f *TranslationServiceFactoryV2) GetTranslationCount(provider TranslationProvider) (int, int, error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	switch provider {
	case GoogleTranslate:
		if f.googleService == nil {
			return 0, 0, fmt.Errorf("сервис Google Translate недоступен")
		}
		return f.googleService.TranslationCount(), f.googleService.TranslationLimit(), nil
	case OpenAI, ClaudeAI, DeepL:
		// Эти провайдеры не имеют встроенного счетчика
		return 0, 0, nil
	case Manual:
		return 0, 0, nil
	default:
		return 0, 0, fmt.Errorf("неизвестный провайдер перевода: %s", provider)
	}
}

// Translate переводит текст с автоматическим fallback на все доступные провайдеры
func (f *TranslationServiceFactoryV2) Translate(ctx context.Context, text string, sourceLanguage string, targetLanguage string) (string, error) {
	f.mutex.RLock()
	chain := make([]TranslationProvider, len(f.fallbackChain))
	copy(chain, f.fallbackChain)
	f.mutex.RUnlock()

	var lastError error
	attemptedProviders := make([]string, 0, len(chain))

	// Пробуем каждый провайдер в цепочке
	for _, provider := range chain {
		service, err := f.GetTranslationService(provider)
		if err != nil {
			logger.Debug().
				Str("provider", string(provider)).
				Err(err).
				Msg("Провайдер недоступен, пропускаем")
			continue
		}

		result, err := service.Translate(ctx, text, sourceLanguage, targetLanguage)
		if err == nil {
			if len(attemptedProviders) > 0 {
				logger.Info().
					Str("successProvider", string(provider)).
					Interface("failedProviders", attemptedProviders).
					Msg("Перевод выполнен после fallback")
			}
			return result, nil
		}

		// Сохраняем ошибку и пробуем следующий провайдер
		lastError = err
		attemptedProviders = append(attemptedProviders, string(provider))

		logger.Warn().
			Err(err).
			Str("provider", string(provider)).
			Int("attemptNumber", len(attemptedProviders)).
			Int("remainingProviders", len(chain)-len(attemptedProviders)).
			Msg("Перевод не удался, пробуем следующий провайдер")
	}

	// Если все провайдеры не сработали, возвращаем mock перевод
	if lastError != nil {
		logger.Error().
			Err(lastError).
			Interface("attemptedProviders", attemptedProviders).
			Msg("Все провайдеры перевода не сработали, возвращаем mock")

		// Возвращаем mock перевод
		mockTranslation := fmt.Sprintf("[%s] %s", strings.ToUpper(targetLanguage), text)
		return mockTranslation, nil
	}

	return "", fmt.Errorf("не удалось выполнить перевод: нет доступных провайдеров")
}

// TranslateWithProvider переводит текст с указанным провайдером
func (f *TranslationServiceFactoryV2) TranslateWithProvider(ctx context.Context, text string, sourceLanguage string, targetLanguage string, provider TranslationProvider) (string, error) {
	service, err := f.GetTranslationService(provider)
	if err != nil {
		return "", fmt.Errorf("ошибка получения сервиса перевода: %w", err)
	}

	return service.Translate(ctx, text, sourceLanguage, targetLanguage)
}

// TranslateWithContext переводит текст с учетом контекста
func (f *TranslationServiceFactoryV2) TranslateWithContext(ctx context.Context, text string, sourceLanguage string, targetLanguage string, context string, fieldName string) (string, error) {
	f.mutex.RLock()
	chain := make([]TranslationProvider, len(f.fallbackChain))
	copy(chain, f.fallbackChain)
	f.mutex.RUnlock()

	var lastError error
	for _, provider := range chain {
		service, err := f.GetTranslationService(provider)
		if err != nil {
			continue
		}

		result, err := service.TranslateWithContext(ctx, text, sourceLanguage, targetLanguage, context, fieldName)
		if err == nil {
			return result, nil
		}
		lastError = err
	}

	if lastError != nil {
		// Mock перевод с контекстом
		return fmt.Sprintf("[%s] %s", strings.ToUpper(targetLanguage), text), nil
	}

	return "", fmt.Errorf("не удалось выполнить перевод с контекстом")
}

// TranslateToAllLanguages переводит текст на все поддерживаемые языки
func (f *TranslationServiceFactoryV2) TranslateToAllLanguages(ctx context.Context, text string) (map[string]string, error) {
	f.mutex.RLock()
	chain := make([]TranslationProvider, len(f.fallbackChain))
	copy(chain, f.fallbackChain)
	f.mutex.RUnlock()

	for _, provider := range chain {
		service, err := f.GetTranslationService(provider)
		if err != nil {
			continue
		}

		result, err := service.TranslateToAllLanguages(ctx, text)
		if err == nil {
			return result, nil
		}
	}

	// Mock переводы для всех языков
	return map[string]string{
		"en": "[EN] " + text,
		"ru": "[RU] " + text,
		"sr": "[SR] " + text,
	}, nil
}

// TranslateEntityFields переводит поля сущности с автоматическим fallback
func (f *TranslationServiceFactoryV2) TranslateEntityFields(ctx context.Context, sourceLanguage string, targetLanguages []string, fields map[string]string) (map[string]map[string]string, error) {
	f.mutex.RLock()
	chain := make([]TranslationProvider, len(f.fallbackChain))
	copy(chain, f.fallbackChain)
	f.mutex.RUnlock()

	var lastError error
	attemptedProviders := make([]string, 0, len(chain))

	for _, provider := range chain {
		service, err := f.GetTranslationService(provider)
		if err != nil {
			continue
		}

		result, err := service.TranslateEntityFields(ctx, sourceLanguage, targetLanguages, fields)
		if err == nil {
			if len(attemptedProviders) > 0 {
				logger.Info().
					Str("successProvider", string(provider)).
					Interface("failedProviders", attemptedProviders).
					Int("fieldsCount", len(fields)).
					Interface("targetLanguages", targetLanguages).
					Msg("Перевод полей выполнен после fallback")
			}
			return result, nil
		}

		lastError = err
		attemptedProviders = append(attemptedProviders, string(provider))

		logger.Warn().
			Err(err).
			Str("provider", string(provider)).
			Int("fieldsCount", len(fields)).
			Msg("Не удалось перевести поля, пробуем следующий провайдер")
	}

	if lastError != nil {
		logger.Error().
			Err(lastError).
			Interface("attemptedProviders", attemptedProviders).
			Msg("Все провайдеры не смогли перевести поля, возвращаем mock")

		// Возвращаем mock переводы
		result := make(map[string]map[string]string)
		for _, targetLang := range targetLanguages {
			translations := make(map[string]string)
			for fieldName, fieldValue := range fields {
				translations[fieldName] = fmt.Sprintf("[%s] %s", strings.ToUpper(targetLang), fieldValue)
			}
			result[targetLang] = translations
		}
		return result, nil
	}

	return nil, fmt.Errorf("не удалось перевести поля: нет доступных провайдеров")
}

// DetectLanguage определяет язык текста
func (f *TranslationServiceFactoryV2) DetectLanguage(ctx context.Context, text string) (string, float64, error) {
	f.mutex.RLock()
	chain := make([]TranslationProvider, len(f.fallbackChain))
	copy(chain, f.fallbackChain)
	f.mutex.RUnlock()

	for _, provider := range chain {
		service, err := f.GetTranslationService(provider)
		if err != nil {
			continue
		}

		lang, confidence, err := service.DetectLanguage(ctx, text)
		if err == nil {
			return lang, confidence, nil
		}
	}

	// Простая эвристика как fallback
	if containsCyrillic(text) {
		if containsSerbian(text) {
			return "sr", 0.8, nil
		}
		return "ru", 0.8, nil
	}
	return "en", 0.8, nil
}

// ModerateText выполняет модерацию текста
func (f *TranslationServiceFactoryV2) ModerateText(ctx context.Context, text string, language string) (string, error) {
	// Claude AI лучше всего подходит для модерации
	if f.claudeService != nil {
		return f.claudeService.ModerateText(ctx, text, language)
	}

	// OpenAI как резервный вариант
	if f.openAIService != nil {
		return f.openAIService.ModerateText(ctx, text, language)
	}

	// Остальные провайдеры не поддерживают модерацию
	return text, nil
}

// UpdateTranslation обновляет перевод с информацией о провайдере
func (f *TranslationServiceFactoryV2) UpdateTranslation(ctx context.Context, translation *models.Translation, provider TranslationProvider, userID int) error {
	// Создаем отдельное поле для хранения информации о провайдере
	if translation.Metadata == nil {
		translation.Metadata = make(map[string]interface{})
	}

	// Добавляем информацию о провайдере
	translation.Metadata["provider"] = string(provider)

	// Ищем любое доступное хранилище
	var storage storage.Storage
	if f.googleService != nil && f.googleService.storage != nil {
		storage = f.googleService.storage
	} else if f.openAIService != nil && f.openAIService.storage != nil {
		storage = f.openAIService.storage
	}

	if storage == nil {
		return fmt.Errorf("не найдено доступное хранилище для обновления перевода")
	}

	return f.updateTranslationWithStorage(ctx, translation, storage, userID)
}

// updateTranslationWithStorage - вспомогательный метод для обновления перевода через хранилище
func (f *TranslationServiceFactoryV2) updateTranslationWithStorage(ctx context.Context, translation *models.Translation, storage storage.Storage, userID int) error {
	// Преобразуем метаданные в JSON
	metadataJSON, err := json.Marshal(translation.Metadata)
	if err != nil {
		return fmt.Errorf("ошибка сериализации метаданных: %w", err)
	}

	// Выполняем запрос к хранилищу
	query := `
        INSERT INTO translations (
            entity_type, entity_id, language, field_name,
            translated_text, is_machine_translated, is_verified, metadata,
            last_modified_by
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        ON CONFLICT (entity_type, entity_id, language, field_name)
        DO UPDATE SET
            translated_text = EXCLUDED.translated_text,
            is_machine_translated = EXCLUDED.is_machine_translated,
            is_verified = EXCLUDED.is_verified,
            metadata = EXCLUDED.metadata,
            last_modified_by = EXCLUDED.last_modified_by,
            updated_at = CURRENT_TIMESTAMP
    `

	var lastModifiedBy interface{}
	if userID > 0 {
		lastModifiedBy = userID
	} else {
		lastModifiedBy = nil
	}

	_, err = storage.Exec(ctx, query,
		translation.EntityType,
		translation.EntityID,
		translation.Language,
		translation.FieldName,
		translation.TranslatedText,
		translation.IsMachineTranslated,
		translation.IsVerified,
		metadataJSON,
		lastModifiedBy)

	return err
}

// Вспомогательные функции containsCyrillic и containsSerbian
// уже определены в claude_translation.go и deepl_translation.go
