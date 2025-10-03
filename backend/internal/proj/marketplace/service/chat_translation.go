// backend/internal/proj/marketplace/service/chat_translation.go
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"backend/internal/domain/models"
	"backend/internal/logger"
)

// ChatTranslationService обрабатывает переводы сообщений чата
type ChatTranslationService struct {
	translationSvc TranslationServiceInterface
	redisClient    *redis.Client
}

// NewChatTranslationService создает новый сервис переводов чатов
func NewChatTranslationService(
	translationSvc TranslationServiceInterface,
	redisClient *redis.Client,
) *ChatTranslationService {
	return &ChatTranslationService{
		translationSvc: translationSvc,
		redisClient:    redisClient,
	}
}

// TranslateMessage переводит одно сообщение на целевой язык
func (s *ChatTranslationService) TranslateMessage(
	ctx context.Context,
	message *models.MarketplaceMessage,
	targetLanguage string,
) error {
	// Если язык не установлен, определяем его
	if message.OriginalLanguage == "" || message.OriginalLanguage == "unknown" {
		err := s.DetectAndSetLanguage(ctx, message)
		if err != nil {
			logger.Warn().Err(err).Msg("Failed to detect language, will try to translate anyway")
			// Продолжаем перевод даже если определение языка не удалось
			message.OriginalLanguage = "auto"
		}
	}

	// Пропускаем если язык совпадает с целевым
	if message.OriginalLanguage == targetLanguage {
		logger.Debug().
			Int("messageId", message.ID).
			Str("lang", targetLanguage).
			Msg("Skipping translation - same language")
		return nil
	}

	// Проверяем кеш Redis
	cacheKey := s.getCacheKey(message.ID, targetLanguage)
	cached, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		// Cache HIT
		if message.Translations == nil {
			message.Translations = make(map[string]string)
		}
		message.Translations[targetLanguage] = cached
		message.ChatTranslationMetadata = &models.ChatTranslationMetadata{
			TranslatedFrom: message.OriginalLanguage,
			TranslatedTo:   targetLanguage,
			TranslatedAt:   time.Now(),
			CacheHit:       true,
			Provider:       "claude-haiku",
		}
		logger.Debug().
			Int("messageId", message.ID).
			Str("targetLang", targetLanguage).
			Msg("Translation cache HIT")
		return nil
	}

	// Cache MISS - вызываем API
	logger.Debug().
		Int("messageId", message.ID).
		Str("from", message.OriginalLanguage).
		Str("to", targetLanguage).
		Msg("Translation cache MISS - calling API")

	translated, err := s.translationSvc.Translate(
		ctx,
		message.Content,
		message.OriginalLanguage,
		targetLanguage,
	)
	if err != nil {
		logger.Error().
			Err(err).
			Int("messageId", message.ID).
			Str("targetLang", targetLanguage).
			Msg("Translation failed")
		return fmt.Errorf("translation failed: %w", err)
	}

	// Сохраняем в кеш (TTL 30 дней)
	err = s.redisClient.Set(ctx, cacheKey, translated, 30*24*time.Hour).Err()
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to cache translation")
	}

	// Обновляем сообщение
	if message.Translations == nil {
		message.Translations = make(map[string]string)
	}
	message.Translations[targetLanguage] = translated
	message.ChatTranslationMetadata = &models.ChatTranslationMetadata{
		TranslatedFrom: message.OriginalLanguage,
		TranslatedTo:   targetLanguage,
		TranslatedAt:   time.Now(),
		CacheHit:       false,
		Provider:       "claude-haiku",
	}

	logger.Info().
		Int("messageId", message.ID).
		Str("targetLang", targetLanguage).
		Int("originalLen", len(message.Content)).
		Int("translatedLen", len(translated)).
		Msg("Translation completed")

	return nil
}

// TranslateBatch переводит несколько сообщений параллельно
func (s *ChatTranslationService) TranslateBatch(
	ctx context.Context,
	messages []*models.MarketplaceMessage,
	targetLanguage string,
) error {
	if len(messages) == 0 {
		return nil
	}

	// Ограничиваем параллелизм (10 одновременных запросов)
	semaphore := make(chan struct{}, 10)
	errChan := make(chan error, len(messages))
	doneChan := make(chan struct{})

	processed := 0
	for _, msg := range messages {
		semaphore <- struct{}{} // Acquire
		go func(m *models.MarketplaceMessage) {
			defer func() { <-semaphore }() // Release

			err := s.TranslateMessage(ctx, m, targetLanguage)
			if err != nil {
				errChan <- err
			}
		}(msg)
		processed++
	}

	// Ждем завершения всех горутин
	go func() {
		for i := 0; i < cap(semaphore); i++ {
			semaphore <- struct{}{}
		}
		close(errChan)
		close(doneChan)
	}()

	<-doneChan

	// Собираем ошибки (логируем, но не прерываем)
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		logger.Warn().
			Int("failedCount", len(errors)).
			Int("totalCount", len(messages)).
			Msg("Some translations failed in batch")
		// Не возвращаем ошибку, частичный успех - это OK
	}

	return nil
}

// DetectAndSetLanguage определяет язык сообщения и устанавливает original_language
func (s *ChatTranslationService) DetectAndSetLanguage(
	ctx context.Context,
	message *models.MarketplaceMessage,
) error {
	if message.OriginalLanguage != "" {
		return nil // Уже установлен
	}

	lang, confidence, err := s.translationSvc.DetectLanguage(ctx, message.Content)
	if err != nil {
		logger.Warn().Err(err).Msg("Language detection failed, defaulting to 'unknown'")
		message.OriginalLanguage = "unknown"
		return nil
	}

	// Требуем минимальную уверенность 70%
	if confidence < 0.7 {
		logger.Warn().
			Float64("confidence", confidence).
			Msg("Low confidence in language detection")
		message.OriginalLanguage = "unknown"
		return nil
	}

	message.OriginalLanguage = lang
	logger.Debug().
		Str("detected", lang).
		Float64("confidence", confidence).
		Msg("Language detected")

	return nil
}

// getCacheKey генерирует ключ для Redis
func (s *ChatTranslationService) getCacheKey(messageID int, targetLang string) string {
	return fmt.Sprintf("chat:translation:%d:%s", messageID, targetLang)
}

// SaveTranslationToDB сохраняет перевод в БД (для персистентности)
func (s *ChatTranslationService) SaveTranslationToDB(
	ctx context.Context,
	messageID int,
	translations map[string]string,
) error {
	// Конвертируем в JSONB
	translationsJSON, err := json.Marshal(translations)
	if err != nil {
		return fmt.Errorf("failed to marshal translations: %w", err)
	}

	logger.Debug().
		Int("messageId", messageID).
		Str("translations", string(translationsJSON)).
		Msg("Translation would be saved to DB (not implemented yet)")

	// TODO: Интегрировать с storage layer
	// Временно только логируем

	return nil
}

// GetUserTranslationSettings получает настройки перевода пользователя
func (s *ChatTranslationService) GetUserTranslationSettings(
	ctx context.Context,
	userID int,
) (*models.ChatUserSettings, error) {
	// TODO: Загрузить из user_privacy_settings.settings JSONB
	// Временно возвращаем defaults
	logger.Debug().
		Int("userId", userID).
		Msg("Getting user translation settings (returning defaults)")

	return &models.ChatUserSettings{
		AutoTranslate:     false, // По умолчанию выключено
		PreferredLanguage: "en",
		ShowLanguageBadge: true,
	}, nil
}
