package translation_admin

import (
	"context"
	"fmt"

	"backend/internal/domain/models"
)

// ValidateTranslations validates translations for consistency
func (s *Service) ValidateTranslations(ctx context.Context, req models.ValidateTranslationsRequest) ([]models.ValidationResult, error) {
	var results []models.ValidationResult

	modules := s.modules
	if req.Module != "" {
		modules = []string{req.Module}
	}

	languages := s.supportedLangs
	if len(req.Languages) > 0 {
		languages = req.Languages
	}

	for _, module := range modules {
		moduleResults := s.validateModule(module, languages, req)
		results = append(results, moduleResults...)
	}

	return results, nil
}

// GetStatistics returns overall translation statistics
func (s *Service) GetStatistics(ctx context.Context) (*models.TranslationStatistics, error) {
	stats := &models.TranslationStatistics{
		LanguageStats: make(map[string]models.LanguageStats),
	}

	// Get module statistics
	modules, err := s.GetFrontendModules(ctx)
	if err != nil {
		return nil, err
	}
	stats.ModuleStats = modules

	// Calculate totals
	for _, module := range modules {
		stats.TotalTranslations += module.Keys * len(s.supportedLangs)
		stats.CompleteTranslations += module.Complete
		stats.MissingTranslations += module.Missing
		stats.PlaceholderCount += module.Placeholders
	}

	// Get language statistics from database
	for _, lang := range s.supportedLangs {
		dbStats, err := s.getLanguageStatsFromDB(ctx, lang)
		if err != nil {
			s.logger.Error().Err(err).Str("language", lang).Msg("Failed to get language stats")
			continue
		}
		stats.LanguageStats[lang] = dbStats
	}

	// Get recent changes
	recentLogs, err := s.auditRepo.GetRecentLogs(ctx, 10)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to get recent logs")
	} else {
		stats.RecentChanges = recentLogs
	}

	return stats, nil
}

// Helper methods for validation and statistics

func (s *Service) analyzeModuleForLanguage(module, lang string) (models.ModuleStats, error) {
	stats := models.ModuleStats{}

	data, err := s.loadModuleFile(module, lang)
	if err != nil {
		return stats, err
	}

	s.analyzeData(data, &stats)
	return stats, nil
}

func (s *Service) analyzeData(data map[string]interface{}, stats *models.ModuleStats) {
	for _, value := range data {
		switch v := value.(type) {
		case string:
			stats.Total++
			switch {
			case s.isPlaceholder(v):
				stats.Placeholders++
			case v == "":
				stats.Missing++
			default:
				stats.Complete++
			}
		case map[string]interface{}:
			s.analyzeData(v, stats)
		}
	}
}

func (s *Service) validateModule(module string, languages []string, req models.ValidateTranslationsRequest) []models.ValidationResult {
	var results []models.ValidationResult

	// Load module data
	moduleData := make(map[string]map[string]interface{})
	for _, lang := range languages {
		data, err := s.loadModuleFile(module, lang)
		if err != nil {
			continue
		}
		moduleData[lang] = data
	}

	// Extract all keys
	allKeys := s.extractAllKeys(moduleData)

	for _, key := range allKeys {
		issues := []models.ValidationIssue{}

		// Check each language
		for _, lang := range languages {
			value := s.getNestedValue(moduleData[lang], key)

			// Check for missing translations
			if value == "" {
				issues = append(issues, models.ValidationIssue{
					Type:     "missing",
					Language: lang,
					Message:  fmt.Sprintf("Translation missing for key %s in %s", key, lang),
					Severity: "error",
				})
			}

			// Check for placeholders
			if req.CheckHTML && s.isPlaceholder(value) {
				issues = append(issues, models.ValidationIssue{
					Type:     "placeholder",
					Language: lang,
					Message:  fmt.Sprintf("Placeholder found in %s: %s", lang, value),
					Severity: "warning",
				})
			}

			// Check for variables
			if req.CheckVars {
				// TODO: Implement variable pattern validation across languages
				s.logger.Debug().Msg("Variable validation requested but not yet implemented")
			}
		}

		if len(issues) > 0 {
			results = append(results, models.ValidationResult{
				Module: module,
				Key:    key,
				Issues: issues,
			})
		}
	}

	return results
}

func (s *Service) getLanguageStatsFromDB(ctx context.Context, language string) (models.LanguageStats, error) {
	stats := models.LanguageStats{}

	// Get translations from database
	translations, err := s.translationRepo.GetTranslations(ctx, map[string]interface{}{
		"language": language,
	})
	if err != nil {
		return stats, err
	}

	stats.Total = len(translations)

	for _, trans := range translations {
		if trans.IsVerified {
			stats.Verified++
		}
		if trans.IsMachineTranslated {
			stats.MachineTranslated++
		}
		if trans.TranslatedText != "" && !s.isPlaceholder(trans.TranslatedText) {
			stats.Complete++
		}
	}

	if stats.Total > 0 {
		stats.Coverage = float64(stats.Complete) / float64(stats.Total) * 100
	}

	return stats, nil
}

func (s *Service) countCompleteTranslations(languages map[string]models.ModuleStats) int {
	total := 0
	for _, stats := range languages {
		total += stats.Complete
	}
	return total
}

func (s *Service) countIncompleteTranslations(languages map[string]models.ModuleStats) int {
	total := 0
	for _, stats := range languages {
		total += stats.Incomplete
	}
	return total
}

func (s *Service) countPlaceholders(languages map[string]models.ModuleStats) int {
	total := 0
	for _, stats := range languages {
		total += stats.Placeholders
	}
	return total
}

func (s *Service) countMissing(languages map[string]models.ModuleStats) int {
	total := 0
	for _, stats := range languages {
		total += stats.Missing
	}
	return total
}
