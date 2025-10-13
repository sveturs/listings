// backend/internal/proj/c2c/service/marketplace_translations.go
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"backend/internal/domain/models"
)

// UpdateTranslation обновляет перевод сущности
func (s *MarketplaceService) UpdateTranslation(ctx context.Context, translation *models.Translation) error {
	// Используем сервис перевода по умолчанию (Google Translate)
	// Передаём 0 в качестве userID, так как этот метод не имеет доступа к user_id
	return s.UpdateTranslationWithProvider(ctx, translation, GoogleTranslate, 0)
}

// SaveTranslation is an alias for UpdateTranslation for compatibility
func (s *MarketplaceService) SaveTranslation(ctx context.Context, entityType string, entityID int, language, fieldName, translatedText string, metadata map[string]interface{}) error {
	translation := &models.Translation{
		EntityType:     entityType,
		EntityID:       entityID,
		Language:       language,
		FieldName:      fieldName,
		TranslatedText: translatedText,
		IsVerified:     true,
		Metadata:       metadata,
	}
	return s.UpdateTranslation(ctx, translation)
}

// TranslateText переводит текст на указанный язык
func (s *MarketplaceService) TranslateText(ctx context.Context, text, sourceLanguage, targetLanguage string) (string, error) {
	if s.translationService == nil {
		return "", fmt.Errorf("translation service not available")
	}

	return s.translationService.Translate(ctx, text, sourceLanguage, targetLanguage)
}

// UpdateTranslationWithProvider обновляет перевод с использованием указанного провайдера
func (s *MarketplaceService) UpdateTranslationWithProvider(ctx context.Context, translation *models.Translation, provider TranslationProvider, userID int) error {
	// Проверяем, есть ли фабрика сервисов перевода
	if factory, ok := s.translationService.(TranslationFactoryInterface); ok {
		// Используем фабрику для обновления перевода с информацией о провайдере
		return factory.UpdateTranslation(ctx, translation, provider, userID)
	}

	// Если фабрики нет, используем прямой запрос к базе данных
	// Подготавливаем метаданные
	var metadataJSON []byte
	var err error

	if translation.Metadata == nil {
		translation.Metadata = map[string]interface{}{"provider": string(provider)}
	} else if _, exists := translation.Metadata["provider"]; !exists {
		translation.Metadata["provider"] = string(provider)
	}

	metadataJSON, err = json.Marshal(translation.Metadata)
	if err != nil {
		log.Printf("Ошибка сериализации метаданных: %v", err)
		metadataJSON = []byte("{}")
	}

	query := insertTranslationQuery

	var lastModifiedBy interface{}
	if userID > 0 {
		lastModifiedBy = userID
	} else {
		lastModifiedBy = nil
	}

	_, err = s.storage.Exec(ctx, query,
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

// SaveAddressTranslations сохраняет переводы адресных полей объявления
func (s *MarketplaceService) SaveAddressTranslations(ctx context.Context, listingID int, addressFields map[string]string, sourceLanguage string, targetLanguages []string) error {
	// Проверяем, есть ли адресные поля для перевода
	if len(addressFields) == 0 {
		return nil
	}

	// Проверяем наличие сервиса перевода
	if s.translationService == nil {
		log.Printf("Translation service not available for address translations")
		return nil
	}

	// Переводим адресные поля на все целевые языки
	translations, err := s.translationService.TranslateEntityFields(ctx, sourceLanguage, targetLanguages, addressFields)
	if err != nil {
		log.Printf("Error translating address fields for listing %d: %v", listingID, err)
		return fmt.Errorf("error translating address fields: %w", err)
	}

	// Сохраняем переводы в базу данных
	for language, fields := range translations {
		// Пропускаем исходный язык - он уже сохранен в основных полях объявления
		if language == sourceLanguage {
			continue
		}

		for fieldName, translatedText := range fields {
			// Сохраняем перевод для каждого поля
			err := s.SaveTranslation(ctx, "listing", listingID, language, fieldName, translatedText, map[string]interface{}{
				"source_language": sourceLanguage,
				"provider":        "google_translate",
				"is_address":      true,
			})
			if err != nil {
				log.Printf("Error saving translation for field %s to language %s: %v", fieldName, language, err)
				// Продолжаем с другими переводами, не прерываем процесс
			}
		}
	}

	log.Printf("Successfully saved address translations for listing %d", listingID)
	return nil
}
