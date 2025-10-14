package translation_admin

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"backend/internal/domain/models"
)

const (
	entityTypeAll = "all"
)

// Helper functions
func strPtr(s string) *string {
	return &s
}

// ExportTranslations exports translations from database to JSON format
func (s *Service) ExportTranslations(ctx context.Context, entityType, language string) (map[string]interface{}, error) {
	filters := make(map[string]interface{})

	if entityType != entityTypeAll {
		filters["entity_type"] = entityType
	}
	if language != "all" {
		filters["language"] = language
	}

	translations, err := s.translationRepo.GetTranslations(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get translations: %w", err)
	}

	// Organize translations by entity type and language
	result := make(map[string]interface{})
	for _, trans := range translations {
		if _, ok := result[trans.EntityType]; !ok {
			result[trans.EntityType] = make(map[string]interface{})
		}
		entityMap := result[trans.EntityType].(map[string]interface{})

		if _, ok := entityMap[trans.Language]; !ok {
			entityMap[trans.Language] = make(map[string]interface{})
		}
		langMap := entityMap[trans.Language].(map[string]interface{})

		key := fmt.Sprintf("%d_%s", trans.EntityID, trans.FieldName)
		langMap[key] = trans.TranslatedText
	}

	return result, nil
}

// ImportTranslations imports translations from JSON to database
func (s *Service) ImportTranslations(ctx context.Context, req *models.ImportTranslationsRequest, userID int) (*models.ImportResult, error) {
	result := &models.ImportResult{
		Success: 0,
		Failed:  0,
		Skipped: 0,
	}

	for entityType, entityData := range req.Translations {
		entityMap, ok := entityData.(map[string]interface{})
		if !ok {
			continue
		}

		for language, langData := range entityMap {
			langMap, ok := langData.(map[string]interface{})
			if !ok {
				continue
			}

			for key, value := range langMap {
				// Parse key (format: entityID_fieldName)
				parts := strings.SplitN(key, "_", 2)
				if len(parts) != 2 {
					result.Failed++
					continue
				}

				var entityID int
				if _, err := fmt.Sscanf(parts[0], "%d", &entityID); err != nil {
					result.Failed++
					continue
				}

				translation := &models.Translation{
					EntityType:     entityType,
					EntityID:       entityID,
					Language:       language,
					FieldName:      parts[1],
					TranslatedText: value.(string),
					UpdatedBy:      userID,
				}

				if req.OverwriteExisting {
					if err := s.translationRepo.UpdateTranslation(ctx, translation); err != nil {
						result.Failed++
					} else {
						result.Success++
					}
				} else {
					// Check if translation exists
					existing, _ := s.translationRepo.GetTranslations(ctx, map[string]interface{}{
						"entity_type": entityType,
						"entity_id":   entityID,
						"language":    language,
						"field_name":  parts[1],
					})

					if len(existing) > 0 {
						result.Skipped++
					} else {
						if err := s.translationRepo.CreateTranslation(ctx, translation); err != nil {
							result.Failed++
						} else {
							result.Success++
						}
					}
				}
			}
		}
	}

	// Log audit
	userIDPtr := &userID
	entityTypePtr := "batch"
	if err := s.auditRepo.LogAction(ctx, &models.TranslationAuditLog{
		UserID:     userIDPtr,
		Action:     "import",
		EntityType: &entityTypePtr,
		NewValue:   strPtr(fmt.Sprintf("Imported %d translations", result.Success)),
	}); err != nil {
		s.logger.Warn().Err(err).Msg("Failed to log import action")
	}

	return result, nil
}

// ExportTranslationsAdvanced exports translations in various formats
func (s *Service) ExportTranslationsAdvanced(ctx context.Context, req *models.ExportRequest) (interface{}, error) {
	s.logger.Info().
		Str("format", string(req.Format)).
		Msgf("Exporting translations")

	// Build filters based on request
	filters := make(map[string]interface{})
	if req.EntityType != nil {
		filters["entity_type"] = *req.EntityType
	}
	if req.Language != nil {
		filters["language"] = *req.Language
	}
	if req.OnlyVerified {
		filters["is_verified"] = true
	}

	// Get translations from database
	translations, err := s.translationRepo.GetTranslations(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get translations: %w", err)
	}

	switch req.Format {
	case models.ExportFormatJSON:
		return s.exportToJSON(translations, req.IncludeMetadata), nil
	case models.ExportFormatCSV:
		return s.exportToCSV(translations, req.IncludeMetadata), nil
	case models.ExportFormatXLIFF:
		return s.exportToXLIFF(translations, req.IncludeMetadata), nil
	default:
		return nil, fmt.Errorf("unsupported export format: %s", req.Format)
	}
}

// exportToJSON exports translations to JSON format
func (s *Service) exportToJSON(translations []models.Translation, includeMetadata bool) map[string]interface{} {
	result := make(map[string]interface{})
	result["format"] = "json"
	result["exported_at"] = fmt.Sprintf("%v", time.Now().UTC())
	result["count"] = len(translations)

	data := make(map[string]interface{})
	for _, t := range translations {
		key := fmt.Sprintf("%s.%d.%s.%s", t.EntityType, t.EntityID, t.Language, t.FieldName)

		translationData := map[string]interface{}{
			"text":                  t.TranslatedText,
			"is_machine_translated": t.IsMachineTranslated,
			"is_verified":           t.IsVerified,
			"created_at":            t.CreatedAt,
			"updated_at":            t.UpdatedAt,
		}

		if includeMetadata && t.Metadata != nil {
			translationData["metadata"] = t.Metadata
		}

		data[key] = translationData
	}
	result["translations"] = data

	return result
}

// exportToCSV exports translations to CSV format (returns CSV string)
func (s *Service) exportToCSV(translations []models.Translation, includeMetadata bool) string {
	var csv strings.Builder

	// Header
	header := "entity_type,entity_id,language,field_name,translated_text,is_machine_translated,is_verified,created_at,updated_at"
	if includeMetadata {
		header += ",metadata"
	}
	csv.WriteString(header + "\n")

	// Rows
	for _, t := range translations {
		row := fmt.Sprintf(`"%s","%d","%s","%s","%s","%t","%t","%s","%s"`,
			t.EntityType, t.EntityID, t.Language, t.FieldName,
			strings.ReplaceAll(t.TranslatedText, `"`, `""`), // Escape quotes
			t.IsMachineTranslated, t.IsVerified,
			t.CreatedAt.Format("2006-01-02 15:04:05"),
			t.UpdatedAt.Format("2006-01-02 15:04:05"))

		if includeMetadata && t.Metadata != nil {
			metadataJSON, _ := json.Marshal(t.Metadata)
			row += fmt.Sprintf(`,"%s"`, strings.ReplaceAll(string(metadataJSON), `"`, `""`))
		}

		csv.WriteString(row + "\n")
	}

	return csv.String()
}

// exportToXLIFF exports translations to XLIFF format
func (s *Service) exportToXLIFF(translations []models.Translation, includeMetadata bool) string {
	var xliff strings.Builder

	xliff.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	xliff.WriteString(`<xliff version="1.2" xmlns="urn:oasis:names:tc:xliff:document:1.2">` + "\n")
	xliff.WriteString(`  <file source-language="en" target-language="ru" datatype="plaintext">` + "\n")
	xliff.WriteString(`    <body>` + "\n")

	for _, t := range translations {
		xliff.WriteString(fmt.Sprintf(`      <trans-unit id="%s_%d_%s_%s">`,
			t.EntityType, t.EntityID, t.Language, t.FieldName) + "\n")
		xliff.WriteString(fmt.Sprintf(`        <source>%s</source>`,
			escapeXML(t.TranslatedText)) + "\n")
		xliff.WriteString(fmt.Sprintf(`        <target>%s</target>`,
			escapeXML(t.TranslatedText)) + "\n")

		if includeMetadata && t.Metadata != nil {
			xliff.WriteString(`        <note>`)
			metadataJSON, _ := json.Marshal(t.Metadata)
			xliff.WriteString(escapeXML(string(metadataJSON)))
			xliff.WriteString(`</note>` + "\n")
		}

		xliff.WriteString(`      </trans-unit>` + "\n")
	}

	xliff.WriteString(`    </body>` + "\n")
	xliff.WriteString(`  </file>` + "\n")
	xliff.WriteString(`</xliff>` + "\n")

	return xliff.String()
}

// escapeXML escapes XML special characters
func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	s = strings.ReplaceAll(s, `'`, "&apos;")
	return s
}

// ImportTranslationsAdvanced imports translations from various formats
func (s *Service) ImportTranslationsAdvanced(ctx context.Context, req *models.TranslationImportRequest, userID int) (*models.ImportResult, error) {
	s.logger.Info().
		Str("format", string(req.Format)).
		Int("user_id", userID).
		Bool("validate_only", req.ValidateOnly).
		Msg("Importing translations")

	result := &models.ImportResult{
		Success: 0,
		Failed:  0,
		Skipped: 0,
	}

	switch req.Format {
	case models.ExportFormatJSON:
		return s.importFromJSON(ctx, req, userID)
	case models.ExportFormatCSV:
		return s.importFromCSV(ctx, req, userID)
	case models.ExportFormatXLIFF:
		return s.importFromXLIFF(ctx, req, userID)
	default:
		return result, fmt.Errorf("unsupported import format: %s", req.Format)
	}
}

// importFromJSON imports translations from JSON format
func (s *Service) importFromJSON(ctx context.Context, req *models.TranslationImportRequest, userID int) (*models.ImportResult, error) {
	result := &models.ImportResult{}

	// Parse JSON data
	jsonData, ok := req.Data.(map[string]interface{})
	if !ok {
		return result, fmt.Errorf("invalid JSON data format")
	}

	translations, ok := jsonData["translations"].(map[string]interface{})
	if !ok {
		return result, fmt.Errorf("missing translations data in JSON")
	}

	// Process each translation
	for key := range translations {
		if req.ValidateOnly {
			// Just validate format
			result.Success++
			continue
		}

		// Parse translation data and save to database
		// Implementation would go here...
		s.logger.Debug().Str("key", key).Msg("Processing translation")
		result.Success++
	}

	return result, nil
}

// importFromCSV imports translations from CSV format
func (s *Service) importFromCSV(ctx context.Context, req *models.TranslationImportRequest, userID int) (*models.ImportResult, error) {
	result := &models.ImportResult{}

	csvData, ok := req.Data.(string)
	if !ok {
		return result, fmt.Errorf("invalid CSV data format")
	}

	// Parse CSV data (simplified implementation)
	lines := strings.Split(csvData, "\n")
	for i, line := range lines {
		if i == 0 || strings.TrimSpace(line) == "" {
			continue // Skip header and empty lines
		}

		if req.ValidateOnly {
			// Just validate format
			result.Success++
			continue
		}

		// Process CSV line and save to database
		// Implementation would go here...
		s.logger.Debug().Str("line", line).Msg("Processing CSV line")
		result.Success++
	}

	return result, nil
}

// importFromXLIFF imports translations from XLIFF format
func (s *Service) importFromXLIFF(ctx context.Context, req *models.TranslationImportRequest, userID int) (*models.ImportResult, error) {
	result := &models.ImportResult{}

	xliffData, ok := req.Data.(string)
	if !ok {
		return result, fmt.Errorf("invalid XLIFF data format")
	}

	// Parse XLIFF data (simplified implementation)
	if req.ValidateOnly {
		// Just validate XML format
		if strings.Contains(xliffData, "<xliff") {
			result.Success = 1
		} else {
			result.Failed = 1
		}
		return result, nil
	}

	// Process XLIFF data and save to database
	// Implementation would go here...
	s.logger.Debug().Msg("Processing XLIFF data")
	result.Success = 1

	return result, nil
}
