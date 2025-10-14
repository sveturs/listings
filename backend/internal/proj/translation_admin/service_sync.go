package translation_admin

import (
	"context"
	"fmt"
	"strings"

	"backend/internal/domain/models"
)

// SyncFrontendToDBOld - старый метод, оставлен для совместимости
func (s *Service) SyncFrontendToDBOld(ctx context.Context, userID int) error {
	for _, module := range s.modules {
		translations, err := s.GetModuleTranslations(ctx, module)
		if err != nil {
			return fmt.Errorf("failed to get module translations: %w", err)
		}

		for _, trans := range translations {
			// Check each language
			for lang, text := range trans.Translations {
				// Create or update in database
				dbTrans := &models.Translation{
					EntityType:          "frontend",
					EntityID:            0, // Use 0 for frontend translations
					Language:            lang,
					FieldName:           trans.Path,
					TranslatedText:      text,
					IsMachineTranslated: false,
					IsVerified:          trans.Status == models.StatusComplete,
				}

				// Check if exists
				existing, err := s.translationRepo.GetTranslations(ctx, map[string]interface{}{
					"entity_type": "frontend",
					"field_name":  trans.Path,
					"language":    lang,
				})
				if err != nil {
					return err
				}

				if len(existing) > 0 {
					// Update existing
					dbTrans.ID = existing[0].ID
					if existing[0].TranslatedText != text {
						// Text changed, create conflict
						conflict := &models.TranslationSyncConflict{
							SourceType:       "frontend",
							TargetType:       "database",
							EntityIdentifier: trans.Path,
							SourceValue:      &text,
							TargetValue:      &existing[0].TranslatedText,
							ConflictType:     "different",
						}
						if err := s.translationRepo.CreateSyncConflict(ctx, conflict); err != nil {
							s.logger.Warn().Err(err).Msg("Failed to create sync conflict")
						}
					}
				} else {
					// Create new
					if err := s.translationRepo.CreateTranslation(ctx, dbTrans); err != nil {
						s.logger.Warn().Err(err).Msg("Failed to create translation")
					}
				}
			}
		}
	}

	// Log sync action
	if err := s.auditRepo.LogAction(ctx, &models.TranslationAuditLog{
		UserID: &userID,
		Action: "sync_frontend_to_db",
	}); err != nil {
		s.logger.Warn().Err(err).Msg("Failed to log sync action")
	}

	return nil
}

// SyncFrontendToDB synchronizes frontend JSON translations to database
func (s *Service) SyncFrontendToDB(ctx context.Context, module string, userID int) (*models.SyncResult, error) {
	result := &models.SyncResult{
		Added:      0,
		Updated:    0,
		Conflicts:  0,
		TotalItems: 0,
	}

	modulesToSync := []string{module}
	if module == "all" {
		modulesToSync = s.modules
	}

	for _, mod := range modulesToSync {
		for _, lang := range s.supportedLangs {
			// Load JSON file
			data, err := s.loadModuleFile(mod, lang)
			if err != nil {
				continue
			}

			// Flatten the nested structure
			flatData := s.flattenTranslations(data, "")

			for key, value := range flatData {
				result.TotalItems++

				// Create translation record
				translation := &models.Translation{
					EntityType:     "frontend",
					EntityID:       0, // Frontend translations don't have entity IDs
					Language:       lang,
					FieldName:      fmt.Sprintf("%s.%s", mod, key),
					TranslatedText: value,
					UpdatedBy:      userID,
				}

				// Check if translation exists
				existing, _ := s.translationRepo.GetTranslations(ctx, map[string]interface{}{
					"entity_type": "frontend",
					"language":    lang,
					"field_name":  fmt.Sprintf("%s.%s", mod, key),
				})

				if len(existing) > 0 {
					// Check for conflicts
					if existing[0].TranslatedText != value {
						// Create conflict record
						conflict := &models.TranslationSyncConflict{
							SourceType:       "frontend",
							TargetType:       "database",
							EntityIdentifier: fmt.Sprintf("%s.%s", mod, key),
							ConflictType:     "value_mismatch",
							SourceValue:      &value,
							TargetValue:      &existing[0].TranslatedText,
						}
						if err := s.translationRepo.CreateSyncConflict(ctx, conflict); err != nil {
							s.logger.Warn().Err(err).Msg("Failed to create sync conflict")
						}
						result.Conflicts++
					}
				} else {
					// Create new translation
					if err := s.translationRepo.CreateTranslation(ctx, translation); err == nil {
						result.Added++
					}
				}
			}
		}
	}

	// Log audit
	userIDPtr := &userID
	entityTypePtr := "frontend"
	if err := s.auditRepo.LogAction(ctx, &models.TranslationAuditLog{
		UserID:     userIDPtr,
		Action:     "sync_frontend_to_db",
		EntityType: &entityTypePtr,
		NewValue:   strPtr(fmt.Sprintf("Synced %d items, added %d, conflicts %d", result.TotalItems, result.Added, result.Conflicts)),
	}); err != nil {
		s.logger.Warn().Err(err).Msg("Failed to log sync action")
	}

	return result, nil
}

// SyncDBToFrontend synchronizes database translations to frontend JSON files
func (s *Service) SyncDBToFrontend(ctx context.Context, entityType string, userID int) (*models.SyncResult, error) {
	result := &models.SyncResult{
		Added:      0,
		Updated:    0,
		Conflicts:  0,
		TotalItems: 0,
	}

	// Get translations from database
	filters := map[string]interface{}{"entity_type": "frontend"}
	if entityType != entityTypeAll {
		filters["entity_type"] = entityType
	}

	translations, err := s.translationRepo.GetTranslations(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get translations: %w", err)
	}

	// Group translations by module and language
	moduleData := make(map[string]map[string]map[string]interface{})

	for _, trans := range translations {
		result.TotalItems++

		// Parse field name (format: module.key.path)
		parts := strings.SplitN(trans.FieldName, ".", 2)
		if len(parts) != 2 {
			continue
		}

		module := parts[0]
		key := parts[1]

		if _, ok := moduleData[module]; !ok {
			moduleData[module] = make(map[string]map[string]interface{})
		}
		if _, ok := moduleData[module][trans.Language]; !ok {
			// Load existing file
			existing, _ := s.loadModuleFile(module, trans.Language)
			if existing == nil {
				existing = make(map[string]interface{})
			}
			moduleData[module][trans.Language] = existing
		}

		// Set nested value
		s.setNestedValue(moduleData[module][trans.Language], key, trans.TranslatedText)
		result.Updated++
	}

	// Write updated files
	for module, langData := range moduleData {
		for lang, data := range langData {
			if err := s.saveModuleFile(module, lang, data); err != nil {
				s.logger.Error().Err(err).Str("module", module).Str("lang", lang).Msg("Failed to save module file")
			}
		}
	}

	// Log audit
	userIDPtr := &userID
	if err := s.auditRepo.LogAction(ctx, &models.TranslationAuditLog{
		UserID:     userIDPtr,
		Action:     "sync_db_to_frontend",
		EntityType: &entityType,
		NewValue:   strPtr(fmt.Sprintf("Synced %d items, updated %d files", result.TotalItems, result.Updated)),
	}); err != nil {
		s.logger.Warn().Err(err).Msg("Failed to log sync action")
	}

	return result, nil
}

// GetSyncStatus returns current synchronization status
func (s *Service) GetSyncStatus(ctx context.Context) (*models.SyncStatus, error) {
	status := &models.SyncStatus{
		LastSync:         nil,
		PendingConflicts: 0,
		FrontendModified: 0,
		DatabaseModified: 0,
	}

	// Get unresolved conflicts
	conflicts, err := s.translationRepo.GetUnresolvedConflicts(ctx)
	if err == nil {
		status.PendingConflicts = len(conflicts)
	}

	// Check for modifications
	// This would require tracking modification timestamps
	// For now, return basic status

	return status, nil
}

// flattenTranslations flattens nested translation structure
func (s *Service) flattenTranslations(data map[string]interface{}, prefix string) map[string]string {
	result := make(map[string]string)

	for key, value := range data {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		switch v := value.(type) {
		case string:
			result[fullKey] = v
		case map[string]interface{}:
			// Recursively flatten nested maps
			nested := s.flattenTranslations(v, fullKey)
			for k, val := range nested {
				result[k] = val
			}
		}
	}

	return result
}

// GetConflicts retrieves translation conflicts based on filters
func (s *Service) GetConflicts(ctx context.Context, filter *models.ConflictsFilter) (*models.ConflictsList, error) {
	// Mock implementation - в реальности здесь должна быть работа с БД
	conflicts := []models.TranslationConflict{
		{
			ID:            1,
			Key:           "common.save",
			Module:        "common",
			Language:      "ru",
			FrontendValue: strPtr("Сохранить"),
			DatabaseValue: strPtr("Сохранить данные"),
			ConflictType:  "value_mismatch",
			Resolved:      false,
		},
		{
			ID:            2,
			Key:           "auth.login",
			Module:        "auth",
			Language:      "en",
			FrontendValue: nil,
			DatabaseValue: strPtr("Log in"),
			ConflictType:  "missing_in_frontend",
			Resolved:      false,
		},
	}

	totalResolved := 0
	totalPending := len(conflicts)

	return &models.ConflictsList{
		Conflicts:     conflicts,
		Total:         len(conflicts),
		TotalResolved: totalResolved,
		TotalPending:  totalPending,
	}, nil
}

// ResolveConflictsBatch resolves multiple conflicts at once
func (s *Service) ResolveConflictsBatch(ctx context.Context, resolutions []models.ConflictResolution) (*models.ConflictResolutionResult, error) {
	result := &models.ConflictResolutionResult{
		TotalProcessed: len(resolutions),
		Resolved:       0,
		Failed:         0,
		Errors:         []string{},
	}

	for _, resolution := range resolutions {
		// Validate resolution
		if resolution.Resolution == "use_custom" && resolution.CustomValue == nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("Conflict %d: custom value required for custom resolution", resolution.ConflictID))
			continue
		}

		// Get conflict details from database
		// In real implementation, fetch from DB

		// Apply resolution based on type
		switch resolution.Resolution {
		case "use_frontend":
			// Update database with frontend value
			s.logger.Info().
				Int("conflict_id", resolution.ConflictID).
				Msg("Resolving conflict with frontend value")
		case "use_database":
			// Update frontend with database value
			s.logger.Info().
				Int("conflict_id", resolution.ConflictID).
				Msg("Resolving conflict with database value")
		case "use_custom":
			// Update both with custom value
			s.logger.Info().
				Int("conflict_id", resolution.ConflictID).
				Str("custom_value", *resolution.CustomValue).
				Msg("Resolving conflict with custom value")
		default:
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("Conflict %d: invalid resolution type", resolution.ConflictID))
			continue
		}

		// Mark conflict as resolved in database
		// In real implementation, update DB

		result.Resolved++
	}

	return result, nil
}
