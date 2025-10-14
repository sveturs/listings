package translation_admin

import (
	"context"
	"fmt"
	"strings"

	"backend/internal/domain"
	"backend/internal/domain/models"
)

// GetDatabaseTranslations retrieves database translations with filters
func (s *Service) GetDatabaseTranslations(ctx context.Context, filters map[string]interface{}) ([]models.Translation, error) {
	translations, err := s.translationRepo.GetTranslations(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get database translations: %w", err)
	}
	return translations, nil
}

// GetTranslationByID retrieves a single translation by ID
func (s *Service) GetTranslationByID(ctx context.Context, id int) (*models.Translation, error) {
	// Use the repository's GetTranslationByID method directly
	repo, ok := s.translationRepo.(*Repository)
	if !ok {
		// Fallback to using GetTranslations
		translations, err := s.translationRepo.GetTranslations(ctx, map[string]interface{}{
			"limit": 1,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get translation: %w", err)
		}
		// Filter by ID manually
		for _, t := range translations {
			if t.ID == id {
				return &t, nil
			}
		}
		return nil, domain.ErrTranslationNotFound
	}

	return repo.GetTranslationByID(ctx, id)
}

// UpdateDatabaseTranslation updates a translation in the database
func (s *Service) UpdateDatabaseTranslation(ctx context.Context, id int, updateReq *models.TranslationUpdateRequest, userID int) error {
	// Get existing translation
	existing, err := s.GetTranslationByID(ctx, id)
	if err != nil {
		return err
	}

	// Log the action
	oldValue := existing.TranslatedText
	if err := s.auditRepo.LogAction(ctx, &models.TranslationAuditLog{
		UserID:     &userID,
		Action:     "update_database_translation",
		EntityType: &existing.EntityType,
		EntityID:   &existing.EntityID,
		OldValue:   &oldValue,
		NewValue:   &updateReq.TranslatedText,
	}); err != nil {
		s.logger.Warn().Err(err).Msg("Failed to log update action")
	}

	// Update translation
	existing.TranslatedText = updateReq.TranslatedText
	if updateReq.IsVerified != nil {
		existing.IsVerified = *updateReq.IsVerified
	}
	if updateReq.IsMachineTranslated != nil {
		existing.IsMachineTranslated = *updateReq.IsMachineTranslated
	}
	if updateReq.Metadata != nil {
		existing.Metadata = updateReq.Metadata
	}

	return s.translationRepo.UpdateTranslation(ctx, existing)
}

// DeleteDatabaseTranslation deletes a translation from the database
func (s *Service) DeleteDatabaseTranslation(ctx context.Context, id int, userID int) error {
	// Get existing translation for audit
	existing, err := s.GetTranslationByID(ctx, id)
	if err != nil {
		return err
	}

	// Log the action
	oldValue := existing.TranslatedText
	if err := s.auditRepo.LogAction(ctx, &models.TranslationAuditLog{
		UserID:     &userID,
		Action:     "delete_database_translation",
		EntityType: &existing.EntityType,
		EntityID:   &existing.EntityID,
		OldValue:   &oldValue,
	}); err != nil {
		s.logger.Warn().Err(err).Msg("Failed to log delete action")
	}

	return s.translationRepo.DeleteTranslation(ctx, id)
}

// PerformBatchOperations performs batch operations on translations
func (s *Service) PerformBatchOperations(ctx context.Context, req *models.BatchOperationsRequest, userID int) (*models.BatchOperationsResult, error) {
	result := &models.BatchOperationsResult{
		Created: 0,
		Updated: 0,
		Deleted: 0,
		Failed:  0,
		Errors:  []string{},
	}

	// Process creates
	for _, trans := range req.Create {
		if err := s.translationRepo.CreateTranslation(ctx, &trans); err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("Failed to create: %v", err))
		} else {
			result.Created++
		}
	}

	// Process updates
	for _, trans := range req.Update {
		if err := s.translationRepo.UpdateTranslation(ctx, &trans); err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("Failed to update ID %d: %v", trans.ID, err))
		} else {
			result.Updated++
		}
	}

	// Process deletes
	for _, id := range req.Delete {
		if err := s.translationRepo.DeleteTranslation(ctx, id); err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("Failed to delete ID %d: %v", id, err))
		} else {
			result.Deleted++
		}
	}

	// Log batch operation
	action := fmt.Sprintf("batch_operation: created=%d, updated=%d, deleted=%d, failed=%d",
		result.Created, result.Updated, result.Deleted, result.Failed)
	if err := s.auditRepo.LogAction(ctx, &models.TranslationAuditLog{
		UserID: &userID,
		Action: action,
	}); err != nil {
		s.logger.Warn().Err(err).Msg("Failed to log batch operation")
	}

	return result, nil
}

// BulkTranslate performs bulk translation of multiple entities
func (s *Service) BulkTranslate(ctx context.Context, req *models.BulkTranslateRequest, userID int) (*models.BatchTranslateResult, error) {
	s.logger.Info().
		Str("entity_type", req.EntityType).
		Int("entity_count", len(req.EntityIDs)).
		Str("source_lang", req.SourceLanguage).
		Strs("target_langs", req.TargetLanguages).
		Int("user_id", userID).
		Msg("Starting bulk translation")

	result := &models.BatchTranslateResult{
		Results:         []models.TranslateResult{},
		TranslatedCount: 0,
		FailedCount:     0,
		Errors:          []string{},
	}

	// Если это категории и нет переводов в таблице translations,
	// создаем их из основной таблицы c2c_categories
	if req.EntityType == "category" && len(req.EntityIDs) > 0 {
		if err := s.ensureCategoryTranslations(ctx, req.EntityIDs, req.SourceLanguage); err != nil {
			s.logger.Error().Err(err).Msg("Failed to ensure category translations")
			// Продолжаем, даже если не удалось создать некоторые переводы
		}
	}

	// Get source translations for the entities
	filters := map[string]interface{}{
		"entity_type": req.EntityType,
		"language":    req.SourceLanguage,
	}

	sourceTranslations, err := s.translationRepo.GetTranslations(ctx, filters)
	if err != nil {
		return result, fmt.Errorf("failed to get source translations: %w", err)
	}

	// Filter by entity IDs if specified
	if len(req.EntityIDs) > 0 {
		filteredTranslations := []models.Translation{}
		for _, t := range sourceTranslations {
			for _, id := range req.EntityIDs {
				if t.EntityID == id {
					filteredTranslations = append(filteredTranslations, t)
					break
				}
			}
		}
		sourceTranslations = filteredTranslations
	}

	// Process each source translation
	for _, sourceTranslation := range sourceTranslations {
		// Skip empty source texts unless they're intentionally empty (like empty descriptions)
		if strings.TrimSpace(sourceTranslation.TranslatedText) == "" {
			// For empty texts, create empty translations in target languages
			for _, targetLang := range req.TargetLanguages {
				if targetLang == req.SourceLanguage {
					continue
				}

				// Check if translation already exists
				existing, _ := s.translationRepo.GetTranslations(ctx, map[string]interface{}{
					"entity_type": sourceTranslation.EntityType,
					"entity_id":   sourceTranslation.EntityID,
					"language":    targetLang,
					"field_name":  sourceTranslation.FieldName,
				})

				if len(existing) == 0 {
					// Create empty translation for consistency
					newTranslation := &models.Translation{
						EntityType:          sourceTranslation.EntityType,
						EntityID:            sourceTranslation.EntityID,
						Language:            targetLang,
						FieldName:           sourceTranslation.FieldName,
						TranslatedText:      "",
						IsMachineTranslated: false,
						IsVerified:          true,
					}

					if err := s.translationRepo.CreateTranslation(ctx, newTranslation); err != nil {
						s.logger.Error().Err(err).Msg("Failed to create empty translation")
					}
				}
			}
			continue // Skip to next source translation
		}

		// Translate to each target language
		for _, targetLang := range req.TargetLanguages {
			if targetLang == req.SourceLanguage {
				continue // Skip same language
			}

			// Check if translation already exists
			existing, _ := s.translationRepo.GetTranslations(ctx, map[string]interface{}{
				"entity_type": sourceTranslation.EntityType,
				"entity_id":   sourceTranslation.EntityID,
				"language":    targetLang,
				"field_name":  sourceTranslation.FieldName,
			})

			// Skip if exists and not overwriting, OR if the text is already the same as source
			if len(existing) > 0 {
				existingTranslation := existing[0]
				// Check if the existing translation is the same as source (untranslated)
				isSameAsSource := existingTranslation.TranslatedText == sourceTranslation.TranslatedText

				if isSameAsSource {
					// Text is the same as source - needs translation
					s.logger.Debug().
						Str("entity_type", sourceTranslation.EntityType).
						Int("entity_id", sourceTranslation.EntityID).
						Str("field", sourceTranslation.FieldName).
						Str("lang", targetLang).
						Msg("Existing translation same as source, will translate")
					// Continue with translation - don't skip
				} else if !req.OverwriteExisting {
					// Has different translation and not overwriting - skip
					s.logger.Debug().
						Str("entity_type", sourceTranslation.EntityType).
						Int("entity_id", sourceTranslation.EntityID).
						Str("field", sourceTranslation.FieldName).
						Str("lang", targetLang).
						Msg("Skipping - has existing translation and overwrite_existing=false")
					continue
				}
			}

			// Detect source language if needed
			actualSourceLang := req.SourceLanguage
			if actualSourceLang == "auto" || actualSourceLang == "" {
				actualSourceLang = detectTextLanguage(sourceTranslation.TranslatedText)
				s.logger.Debug().
					Str("detected_lang", actualSourceLang).
					Str("text", sourceTranslation.TranslatedText).
					Msg("Auto-detected source language")
			}

			// Create translation request
			translateReq := &models.TranslateRequest{
				Text:            sourceTranslation.TranslatedText,
				SourceLanguage:  actualSourceLang,
				TargetLanguages: []string{targetLang},
			}

			// Perform translation
			translationResult, err := s.TranslateText(ctx, translateReq, userID)
			if err != nil {
				result.FailedCount++
				result.Errors = append(result.Errors,
					fmt.Sprintf("Failed to translate %s:%d:%s to %s: %v",
						sourceTranslation.EntityType, sourceTranslation.EntityID,
						sourceTranslation.FieldName, targetLang, err))
				continue
			}

			result.Results = append(result.Results, *translationResult)
			result.TranslatedCount++

			// Auto-approve if requested
			if req.AutoApprove {
				// Check if we need to update existing or create new
				if len(existing) > 0 {
					// Update existing translation
					existingTranslation := existing[0]
					existingTranslation.TranslatedText = translationResult.Translations[targetLang]
					existingTranslation.IsMachineTranslated = true
					existingTranslation.IsVerified = false // Machine translations should be reviewed

					if err := s.translationRepo.UpdateTranslation(ctx, &existingTranslation); err != nil {
						s.logger.Error().Err(err).Msg("Failed to update translation")
						result.FailedCount++
						result.Errors = append(result.Errors,
							fmt.Sprintf("Failed to update translation for %s:%d:%s:%s: %v",
								sourceTranslation.EntityType, sourceTranslation.EntityID,
								sourceTranslation.FieldName, targetLang, err))
					}
				} else {
					// Create new translation record
					newTranslation := &models.Translation{
						EntityType:          sourceTranslation.EntityType,
						EntityID:            sourceTranslation.EntityID,
						Language:            targetLang,
						FieldName:           sourceTranslation.FieldName,
						TranslatedText:      translationResult.Translations[targetLang],
						IsMachineTranslated: true,
						IsVerified:          false, // Machine translations should be reviewed
					}

					if err := s.translationRepo.CreateTranslation(ctx, newTranslation); err != nil {
						s.logger.Error().Err(err).Msg("Failed to save auto-approved translation")
						result.FailedCount++
						result.Errors = append(result.Errors,
							fmt.Sprintf("Failed to create translation for %s:%d:%s:%s: %v",
								sourceTranslation.EntityType, sourceTranslation.EntityID,
								sourceTranslation.FieldName, targetLang, err))
					}
				}
			}
		}
	}

	return result, nil
}

// ensureCategoryTranslations creates missing translations from c2c_categories table
func (s *Service) ensureCategoryTranslations(ctx context.Context, categoryIDs []int, sourceLanguage string) error {
	// Получаем категории из основной таблицы
	// Используем IN вместо ANY так как это обычный sql.DB
	if len(categoryIDs) == 0 {
		return nil
	}

	placeholders := make([]string, len(categoryIDs))
	args := make([]interface{}, len(categoryIDs))
	for i, id := range categoryIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	var queryBuilder strings.Builder
	queryBuilder.WriteString(`
		SELECT id, name, description, seo_title, seo_description
		FROM c2c_categories
		WHERE id IN (`)
	queryBuilder.WriteString(strings.Join(placeholders, ","))
	queryBuilder.WriteString(")")
	query := queryBuilder.String()

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to query categories: %w", err)
	}
	defer func() { _ = rows.Close() }()

	type CategoryData struct {
		ID             int
		Name           *string
		Description    *string
		SEOTitle       *string
		SEODescription *string
	}

	var categories []CategoryData
	for rows.Next() {
		var cat CategoryData
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.Description, &cat.SEOTitle, &cat.SEODescription); err != nil {
			s.logger.Error().Err(err).Msg("Failed to scan category")
			continue
		}
		categories = append(categories, cat)
	}

	if err = rows.Err(); err != nil {
		return fmt.Errorf("failed to iterate rows: %w", err)
	}

	// Для каждой категории создаем недостающие переводы
	for _, cat := range categories {
		// Проверяем какие поля есть в основной таблице
		fieldsToTranslate := map[string]*string{
			"name":            cat.Name,
			"description":     cat.Description,
			"seo_title":       cat.SEOTitle,
			"seo_description": cat.SEODescription,
		}

		for fieldName, fieldValue := range fieldsToTranslate {
			// Пропускаем NULL поля, но создаем перевод для пустых строк
			if fieldValue == nil {
				continue
			}

			// Проверяем, существует ли уже перевод
			existing, _ := s.translationRepo.GetTranslations(ctx, map[string]interface{}{
				"entity_type": "category",
				"entity_id":   cat.ID,
				"language":    sourceLanguage,
				"field_name":  fieldName,
			})

			// Если перевода нет, создаем его
			if len(existing) == 0 {
				newTranslation := &models.Translation{
					EntityType:          "category",
					EntityID:            cat.ID,
					Language:            sourceLanguage,
					FieldName:           fieldName,
					TranslatedText:      *fieldValue,
					IsMachineTranslated: false,
					IsVerified:          true, // Исходные данные считаем проверенными
				}

				if err := s.translationRepo.CreateTranslation(ctx, newTranslation); err != nil {
					s.logger.Error().
						Err(err).
						Int("category_id", cat.ID).
						Str("field", fieldName).
						Msg("Failed to create source translation")
					// Продолжаем с другими полями
				} else {
					s.logger.Info().
						Int("category_id", cat.ID).
						Str("field", fieldName).
						Str("language", sourceLanguage).
						Msg("Created source translation from c2c_categories")
				}
			}
		}
	}

	return nil
}

// detectTextLanguage attempts to detect the language of text based on character patterns
func detectTextLanguage(text string) string {
	// Count character types
	latinCount := 0
	cyrillicCount := 0

	for _, r := range text {
		// Latin characters
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			latinCount++
		}
		// Cyrillic characters
		if (r >= 'А' && r <= 'я') || (r >= 'Ё' && r <= 'ё') {
			cyrillicCount++
		}
	}

	// Determine language based on character counts
	if cyrillicCount > latinCount {
		return "ru" // Russian text
	}

	// Check for specific Serbian Latin patterns
	serbianPatterns := []string{"dj", "nj", "lj", "dž", "š", "č", "ć", "ž", "đ"}
	textLower := strings.ToLower(text)
	for _, pattern := range serbianPatterns {
		if strings.Contains(textLower, pattern) {
			return "sr" // Serbian text
		}
	}

	// Default to Serbian if mostly Latin (since the site is Serbian)
	if latinCount > 0 {
		return "sr"
	}

	return "auto" // Let AI detect
}
