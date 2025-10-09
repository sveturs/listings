// backend/internal/proj/c2c/service/translation_factory.go
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"backend/internal/domain/models"
	"backend/internal/logger"
	"backend/internal/storage"
)

// TranslationFactoryInterface интерфейс для фабрики сервисов перевода
type TranslationFactoryInterface interface {
	TranslationServiceInterface

	// Методы управления провайдерами перевода
	GetTranslationService(provider TranslationProvider) (TranslationServiceInterface, error)
	GetDefaultProvider() TranslationProvider
	SetDefaultProvider(provider TranslationProvider) error
	GetAvailableProviders() []TranslationProvider
	GetTranslationCount(provider TranslationProvider) (int, int, error)

	// Перевод с указанием конкретного провайдера
	TranslateWithProvider(ctx context.Context, text string, sourceLanguage string, targetLanguage string, provider TranslationProvider) (string, error)

	// Обновление перевода с информацией о провайдере
	UpdateTranslation(ctx context.Context, translation *models.Translation, provider TranslationProvider, userID int) error
}

// TranslationServiceFactory создает и управляет различными сервисами перевода
type TranslationServiceFactory struct {
	googleService   *GoogleTranslationService
	openAIService   *TranslationService
	claudeService   *ClaudeTranslationService
	deeplService    *DeepLTranslationService
	defaultProvider TranslationProvider
	mutex           sync.RWMutex
}

// Убедимся, что TranslationServiceFactory реализует интерфейс TranslationFactoryInterface
var _ TranslationFactoryInterface = (*TranslationServiceFactory)(nil)

// NewTranslationServiceFactory создает новый экземпляр фабрики сервисов перевода
func NewTranslationServiceFactory(googleAPIKey, openAIAPIKey string, storage storage.Storage) (*TranslationServiceFactory, error) {
	// Создаем и проверяем сервис Google Translate
	var googleService *GoogleTranslationService
	var openAIService *TranslationService
	var err error

	// Определяем, какой провайдер будет использоваться по умолчанию
	defaultProvider := GoogleTranslate

	// Создаем сервис Google Translate, если ключ предоставлен
	if googleAPIKey != "" {
		googleService, err = NewGoogleTranslationService(googleAPIKey, storage)
		if err != nil {
			logger.Info().Msgf("Не удалось создать сервис Google Translate: %v. Будет использован только OpenAI", err)
			defaultProvider = OpenAI
		}
	} else {
		logger.Info().Msgf("Google Translate API ключ не предоставлен. Будет использован только OpenAI")
		defaultProvider = OpenAI
	}

	// Создаем сервис OpenAI, если ключ предоставлен
	if openAIAPIKey != "" {
		openAIService, err = NewTranslationService(openAIAPIKey)
		if err != nil {
			logger.Info().Msgf("Не удалось создать сервис OpenAI: %v", err)

			// Если Google Translation тоже не удалось создать, возвращаем ошибку
			if googleService == nil {
				return nil, fmt.Errorf("не удалось создать ни один сервис перевода")
			}
		}

		// Добавляем доступ к хранилищу для OpenAI сервиса, если он доступен
		if openAIService != nil && storage != nil {
			openAIService.storage = storage
		}
	} else {
		logger.Info().Msgf("OpenAI API ключ не предоставлен. Будет использован только Google Translate")

		// Если не удалось создать ни один сервис, возвращаем ошибку
		if googleService == nil {
			return nil, fmt.Errorf("не удалось создать ни один сервис перевода: отсутствуют API ключи")
		}
	}

	logger.Info().Str("defaultProvider", string(defaultProvider)).Msg("Создана фабрика сервисов перевода")

	return &TranslationServiceFactory{
		googleService:   googleService,
		openAIService:   openAIService,
		defaultProvider: defaultProvider,
		mutex:           sync.RWMutex{},
	}, nil
}

// GetTranslationService возвращает сервис перевода по запрошенному провайдеру
func (f *TranslationServiceFactory) GetTranslationService(provider TranslationProvider) (TranslationServiceInterface, error) {
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
		// Возвращаем сервис по умолчанию
		if f.defaultProvider == GoogleTranslate && f.googleService != nil {
			return f.googleService, nil
		} else if f.openAIService != nil {
			return f.openAIService, nil
		}
		return nil, fmt.Errorf("ни один сервис перевода не доступен")
	}
}

// GetDefaultProvider возвращает провайдер перевода по умолчанию
func (f *TranslationServiceFactory) GetDefaultProvider() TranslationProvider {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	return f.defaultProvider
}

// SetDefaultProvider устанавливает провайдер перевода по умолчанию
func (f *TranslationServiceFactory) SetDefaultProvider(provider TranslationProvider) error {
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
	logger.Info().Msgf("Установлен провайдер перевода по умолчанию: %s", provider)
	return nil
}

// GetAvailableProviders возвращает список доступных провайдеров перевода
func (f *TranslationServiceFactory) GetAvailableProviders() []TranslationProvider {
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
func (f *TranslationServiceFactory) GetTranslationCount(provider TranslationProvider) (int, int, error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	switch provider {
	case GoogleTranslate:
		if f.googleService == nil {
			return 0, 0, fmt.Errorf("сервис Google Translate недоступен")
		}
		return f.googleService.TranslationCount(), f.googleService.TranslationLimit(), nil
	case OpenAI:
		// OpenAI не имеет встроенного ограничения по количеству, но может иметь внешние ограничения
		return 0, 0, nil
	case ClaudeAI:
		// Claude AI не имеет встроенного ограничения по количеству, но может иметь внешние ограничения
		return 0, 0, nil
	case DeepL:
		// DeepL может иметь ограничения в зависимости от плана
		return 0, 0, nil
	case Manual:
		// Ручной перевод не имеет счетчика
		return 0, 0, nil
	default:
		return 0, 0, fmt.Errorf("неизвестный провайдер перевода: %s", provider)
	}
}

// Реализация интерфейса TranslationServiceInterface

// Translate переводит текст с одного языка на другой, используя провайдер по умолчанию с автоматическим fallback
func (f *TranslationServiceFactory) Translate(ctx context.Context, text string, sourceLanguage string, targetLanguage string) (string, error) {
	// Пытаемся использовать провайдер по умолчанию
	service, err := f.GetTranslationService(f.defaultProvider)
	if err != nil {
		return "", fmt.Errorf("ошибка получения сервиса перевода: %w", err)
	}

	result, err := service.Translate(ctx, text, sourceLanguage, targetLanguage)
	if err == nil {
		return result, nil
	}

	// Логируем ошибку основного провайдера
	logger.Warn().
		Err(err).
		Str("provider", string(f.defaultProvider)).
		Msg("Translation failed with default provider, attempting fallback")

	// Пытаемся использовать резервный провайдер
	var fallbackProvider TranslationProvider
	switch f.defaultProvider {
	case GoogleTranslate:
		if f.openAIService != nil {
			fallbackProvider = OpenAI
			service = f.openAIService
		} else {
			return "", fmt.Errorf("перевод не удался и резервный провайдер недоступен: %w", err)
		}
	case OpenAI:
		if f.googleService != nil {
			fallbackProvider = GoogleTranslate
			service = f.googleService
		} else {
			return "", fmt.Errorf("перевод не удался и резервный провайдер недоступен: %w", err)
		}
	case ClaudeAI:
		if f.googleService != nil {
			fallbackProvider = GoogleTranslate
			service = f.googleService
		} else {
			return "", fmt.Errorf("перевод не удался и резервный провайдер недоступен: %w", err)
		}
	case DeepL:
		if f.googleService != nil {
			fallbackProvider = GoogleTranslate
			service = f.googleService
		} else {
			return "", fmt.Errorf("перевод не удался и резервный провайдер недоступен: %w", err)
		}
	case Manual:
		return "", fmt.Errorf("перевод не удался и резервный провайдер недоступен для ручного перевода: %w", err)
	default:
		// Нет резервного провайдера
		return "", fmt.Errorf("перевод не удался и резервный провайдер недоступен: %w", err)
	}

	logger.Info().
		Str("fallback_provider", string(fallbackProvider)).
		Msg("Attempting translation with fallback provider")

	result, fallbackErr := service.Translate(ctx, text, sourceLanguage, targetLanguage)
	if fallbackErr != nil {
		return "", fmt.Errorf("перевод не удался с обоими провайдерами. Основной: %s, Резервный: %w", err.Error(), fallbackErr)
	}

	return result, nil
}

// TranslateWithProvider переводит текст с одного языка на другой, используя указанный провайдер
func (f *TranslationServiceFactory) TranslateWithProvider(ctx context.Context, text string, sourceLanguage string, targetLanguage string, provider TranslationProvider) (string, error) {
	service, err := f.GetTranslationService(provider)
	if err != nil {
		return "", fmt.Errorf("ошибка получения сервиса перевода: %w", err)
	}

	return service.Translate(ctx, text, sourceLanguage, targetLanguage)
}

// TranslateWithContext переводит текст с учетом контекста
func (f *TranslationServiceFactory) TranslateWithContext(ctx context.Context, text string, sourceLanguage string, targetLanguage string, context string, fieldName string) (string, error) {
	service, err := f.GetTranslationService(f.defaultProvider)
	if err != nil {
		return "", fmt.Errorf("ошибка получения сервиса перевода: %w", err)
	}

	return service.TranslateWithContext(ctx, text, sourceLanguage, targetLanguage, context, fieldName)
}

// TranslateToAllLanguages переводит текст на все поддерживаемые языки
func (f *TranslationServiceFactory) TranslateToAllLanguages(ctx context.Context, text string) (map[string]string, error) {
	service, err := f.GetTranslationService(f.defaultProvider)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения сервиса перевода: %w", err)
	}

	return service.TranslateToAllLanguages(ctx, text)
}

// TranslateEntityFields переводит поля сущности с автоматическим fallback
func (f *TranslationServiceFactory) TranslateEntityFields(ctx context.Context, sourceLanguage string, targetLanguages []string, fields map[string]string) (map[string]map[string]string, error) {
	// Пытаемся использовать провайдер по умолчанию
	service, err := f.GetTranslationService(f.defaultProvider)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения сервиса перевода: %w", err)
	}

	result, err := service.TranslateEntityFields(ctx, sourceLanguage, targetLanguages, fields)
	if err == nil {
		return result, nil
	}

	// Логируем ошибку основного провайдера
	logger.Warn().
		Err(err).
		Str("provider", string(f.defaultProvider)).
		Int("fields_count", len(fields)).
		Int("target_langs", len(targetLanguages)).
		Msg("TranslateEntityFields failed with default provider, attempting fallback")

	// Пытаемся использовать резервный провайдер
	var fallbackProvider TranslationProvider
	switch f.defaultProvider {
	case GoogleTranslate:
		if f.openAIService != nil {
			fallbackProvider = OpenAI
			service = f.openAIService
		} else {
			return nil, fmt.Errorf("перевод полей не удался и резервный провайдер недоступен: %w", err)
		}
	case OpenAI:
		if f.googleService != nil {
			fallbackProvider = GoogleTranslate
			service = f.googleService
		} else {
			return nil, fmt.Errorf("перевод полей не удался и резервный провайдер недоступен: %w", err)
		}
	case ClaudeAI:
		if f.googleService != nil {
			fallbackProvider = GoogleTranslate
			service = f.googleService
		} else {
			return nil, fmt.Errorf("перевод полей не удался и резервный провайдер недоступен: %w", err)
		}
	case DeepL:
		if f.googleService != nil {
			fallbackProvider = GoogleTranslate
			service = f.googleService
		} else {
			return nil, fmt.Errorf("перевод полей не удался и резервный провайдер недоступен: %w", err)
		}
	case Manual:
		return nil, fmt.Errorf("перевод полей не удался и резервный провайдер недоступен для ручного перевода: %w", err)
	default:
		// Нет резервного провайдера
		return nil, fmt.Errorf("перевод полей не удался и резервный провайдер недоступен: %w", err)
	}

	logger.Info().
		Str("fallback_provider", string(fallbackProvider)).
		Msg("Attempting TranslateEntityFields with fallback provider")

	result, fallbackErr := service.TranslateEntityFields(ctx, sourceLanguage, targetLanguages, fields)
	if fallbackErr != nil {
		return nil, fmt.Errorf("перевод полей не удался с обоими провайдерами. Основной: %s, Резервный: %w", err.Error(), fallbackErr)
	}

	return result, nil
}

// DetectLanguage определяет язык текста
func (f *TranslationServiceFactory) DetectLanguage(ctx context.Context, text string) (string, float64, error) {
	service, err := f.GetTranslationService(f.defaultProvider)
	if err != nil {
		return "", 0, fmt.Errorf("ошибка получения сервиса перевода: %w", err)
	}

	return service.DetectLanguage(ctx, text)
}

// ModerateText выполняет модерацию текста
func (f *TranslationServiceFactory) ModerateText(ctx context.Context, text string, language string) (string, error) {
	service, err := f.GetTranslationService(f.defaultProvider)
	if err != nil {
		return "", fmt.Errorf("ошибка получения сервиса перевода: %w", err)
	}

	return service.ModerateText(ctx, text, language)
}

// TranslateWithToneModeration переводит текст с опцией смягчения грубого языка
func (f *TranslationServiceFactory) TranslateWithToneModeration(ctx context.Context, text string, sourceLanguage string, targetLanguage string, moderateTone bool) (string, error) {
	service, err := f.GetTranslationService(f.defaultProvider)
	if err != nil {
		return "", fmt.Errorf("ошибка получения сервиса перевода: %w", err)
	}

	return service.TranslateWithToneModeration(ctx, text, sourceLanguage, targetLanguage, moderateTone)
}

// UpdateTranslation обновляет перевод с информацией о провайдере
func (f *TranslationServiceFactory) UpdateTranslation(ctx context.Context, translation *models.Translation, provider TranslationProvider, userID int) error {
	// Создаем отдельное поле или метаданные для хранения информации о провайдере перевода
	if translation.Metadata == nil {
		translation.Metadata = make(map[string]interface{})
	}

	// Добавляем информацию о провайдере
	translation.Metadata["provider"] = string(provider)

	// Отправляем запрос на обновление перевода через хранилище MarketplaceService
	// Используем хранилище из любого доступного сервиса
	if f.googleService != nil && f.googleService.storage != nil {
		return f.updateTranslationWithStorage(ctx, translation, f.googleService.storage, userID)
	} else if f.openAIService != nil && f.openAIService.storage != nil {
		return f.updateTranslationWithStorage(ctx, translation, f.openAIService.storage, userID)
	}

	return fmt.Errorf("не найдено доступное хранилище для обновления перевода")
}

// updateTranslationWithStorage - вспомогательный метод для обновления перевода через хранилище
func (f *TranslationServiceFactory) updateTranslationWithStorage(ctx context.Context, translation *models.Translation, storage storage.Storage, userID int) error {
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
