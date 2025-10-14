package translation_admin

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"backend/internal/domain/models"
)

// GetFrontendModules returns all frontend translation modules with statistics
func (s *Service) GetFrontendModules(ctx context.Context) ([]models.FrontendModule, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var modules []models.FrontendModule

	for _, moduleName := range s.modules {
		module := models.FrontendModule{
			Name:      moduleName,
			Languages: make(map[string]models.ModuleStats),
		}

		// Analyze each language
		for _, lang := range s.supportedLangs {
			stats, err := s.analyzeModuleForLanguage(moduleName, lang)
			if err != nil {
				s.logger.Error().Err(err).
					Str("module", moduleName).
					Str("language", lang).
					Msg("Failed to analyze module")
				continue
			}
			module.Languages[lang] = stats
		}

		// Calculate overall stats
		module.Keys = s.countModuleKeys(moduleName)
		module.Complete = s.countCompleteTranslations(module.Languages)
		module.Incomplete = s.countIncompleteTranslations(module.Languages)
		module.Placeholders = s.countPlaceholders(module.Languages)
		module.Missing = s.countMissing(module.Languages)

		modules = append(modules, module)
	}

	return modules, nil
}

// GetModuleTranslations returns all translations for a specific module
func (s *Service) GetModuleTranslations(ctx context.Context, moduleName string) ([]models.FrontendTranslation, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if !s.isValidModule(moduleName) {
		return nil, fmt.Errorf("invalid module name: %s", moduleName)
	}

	translations := []models.FrontendTranslation{}

	// Load translations for all languages
	moduleData := make(map[string]map[string]interface{})
	for _, lang := range s.supportedLangs {
		data, err := s.loadModuleFile(moduleName, lang)
		if err != nil {
			s.logger.Error().Err(err).
				Str("module", moduleName).
				Str("language", lang).
				Msg("Failed to load module file")
			continue
		}
		moduleData[lang] = data
	}

	// Extract all keys and build translation objects
	allKeys := s.extractAllKeys(moduleData)

	for _, key := range allKeys {
		trans := models.FrontendTranslation{
			Module:       moduleName,
			Key:          key,
			Path:         fmt.Sprintf("%s.%s", moduleName, key),
			Translations: make(map[string]string),
		}

		// Get translation for each language
		for _, lang := range s.supportedLangs {
			if value := s.getNestedValue(moduleData[lang], key); value != "" {
				trans.Translations[lang] = value
			}
		}

		// Determine status
		trans.Status = s.determineTranslationStatus(trans.Translations)

		// Add metadata
		trans.Metadata = s.buildTranslationMetadata(trans.Translations)

		translations = append(translations, trans)
	}

	return translations, nil
}

// UpdateModuleTranslations updates translations for a module
func (s *Service) UpdateModuleTranslations(ctx context.Context, moduleName string, updates []models.FrontendTranslation, userID int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.isValidModule(moduleName) {
		return fmt.Errorf("invalid module name: %s", moduleName)
	}

	// Load current translations
	moduleData := make(map[string]map[string]interface{})
	for _, lang := range s.supportedLangs {
		data, err := s.loadModuleFile(moduleName, lang)
		if err != nil {
			return fmt.Errorf("failed to load module file: %w", err)
		}
		moduleData[lang] = data
	}

	// Apply updates
	for _, update := range updates {
		for lang, text := range update.Translations {
			s.setNestedValue(moduleData[lang], update.Key, text)
		}

		// Log the update
		if err := s.auditRepo.LogAction(ctx, &models.TranslationAuditLog{
			UserID:     &userID,
			Action:     "update_frontend_translation",
			EntityType: &moduleName,
			NewValue:   &update.Key,
		}); err != nil {
			s.logger.Warn().Err(err).Msg("Failed to log audit action")
		}
	}

	// Save updated files
	for _, lang := range s.supportedLangs {
		if err := s.saveModuleFile(moduleName, lang, moduleData[lang]); err != nil {
			return fmt.Errorf("failed to save module file: %w", err)
		}
	}

	return nil
}

// Helper methods for module and key management

func (s *Service) isValidModule(name string) bool {
	for _, m := range s.modules {
		if m == name {
			return true
		}
	}
	return false
}

func (s *Service) loadModuleFile(module, lang string) (map[string]interface{}, error) {
	// Валидация параметров для предотвращения path traversal
	if !isValidLanguage(lang) || !isValidModule(module) {
		return nil, fmt.Errorf("invalid language or module parameter")
	}

	filePath := filepath.Join(s.frontendPath, "src", "messages", lang, module+".json")
	// Проверяем, что путь находится в разрешённой директории
	cleanPath := filepath.Clean(filePath)
	messagesDir := filepath.Join(s.frontendPath, "src", "messages")
	if !strings.HasPrefix(cleanPath, filepath.Clean(messagesDir)) {
		return nil, fmt.Errorf("access denied: path outside allowed directory")
	}

	// #nosec G304 - путь провалидирован выше
	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Service) saveModuleFile(module, lang string, data map[string]interface{}) error {
	filePath := filepath.Join(s.frontendPath, "src", "messages", lang, module+".json")

	// Create backup
	backupPath := filePath + ".backup"
	if err := s.createBackup(filePath, backupPath); err != nil {
		s.logger.Warn().Err(err).Msg("Failed to create backup")
	}

	// Marshal with indentation
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	// Write file
	return os.WriteFile(filePath, jsonData, 0o600)
}

func (s *Service) createBackup(src, dst string) error {
	// Проверяем пути на безопасность
	cleanSrc := filepath.Clean(src)
	cleanDst := filepath.Clean(dst)
	messagesDir := filepath.Join(s.frontendPath, "src", "messages")
	if !strings.HasPrefix(cleanSrc, filepath.Clean(messagesDir)) ||
		!strings.HasPrefix(cleanDst, filepath.Clean(messagesDir)) {
		return fmt.Errorf("access denied: path outside allowed directory")
	}

	// #nosec G304 - путь провалидирован выше
	source, err := os.Open(cleanSrc)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := source.Close(); closeErr != nil {
			log.Printf("Failed to close source file: %v", closeErr)
		}
	}()

	// #nosec G304 - путь провалидирован выше
	destination, err := os.Create(cleanDst)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := destination.Close(); closeErr != nil {
			log.Printf("Failed to close destination file: %v", closeErr)
		}
	}()

	_, err = destination.ReadFrom(source)
	return err
}

func (s *Service) extractAllKeys(moduleData map[string]map[string]interface{}) []string {
	keySet := make(map[string]bool)

	for _, langData := range moduleData {
		s.extractKeysRecursive(langData, "", keySet)
	}

	// Convert set to slice
	keys := make([]string, 0, len(keySet))
	for key := range keySet {
		keys = append(keys, key)
	}

	return keys
}

func (s *Service) extractKeysRecursive(data map[string]interface{}, prefix string, keySet map[string]bool) {
	for key, value := range data {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		switch v := value.(type) {
		case string:
			keySet[fullKey] = true
		case map[string]interface{}:
			s.extractKeysRecursive(v, fullKey, keySet)
		}
	}
}

func (s *Service) getNestedValue(data map[string]interface{}, key string) string {
	parts := strings.Split(key, ".")
	current := data

	for i, part := range parts {
		if i == len(parts)-1 {
			// Last part, should be a string
			if val, ok := current[part].(string); ok {
				return val
			}
			return ""
		}

		if next, ok := current[part].(map[string]interface{}); ok {
			current = next
		} else {
			return ""
		}
	}

	return ""
}

func (s *Service) setNestedValue(data map[string]interface{}, key string, value string) {
	parts := strings.Split(key, ".")
	current := data

	for i, part := range parts {
		if i == len(parts)-1 {
			// Last part, set the value
			current[part] = value
			return
		}

		if _, ok := current[part]; !ok {
			current[part] = make(map[string]interface{})
		}

		if next, ok := current[part].(map[string]interface{}); ok {
			current = next
		}
	}
}

func (s *Service) countModuleKeys(moduleName string) int {
	count := 0
	for _, lang := range s.supportedLangs {
		data, err := s.loadModuleFile(moduleName, lang)
		if err == nil {
			count = s.countKeysInData(data)
			break // All languages should have same structure
		}
	}
	return count
}

func (s *Service) countKeysInData(data map[string]interface{}) int {
	count := 0
	for _, value := range data {
		switch v := value.(type) {
		case string:
			count++
		case map[string]interface{}:
			count += s.countKeysInData(v)
		}
	}
	return count
}

func (s *Service) determineTranslationStatus(translations map[string]string) models.TranslationStatus {
	if len(translations) == 0 {
		return models.StatusMissing
	}

	hasPlaceholder := false
	hasEmpty := false

	for _, text := range translations {
		if text == "" {
			hasEmpty = true
		} else if s.isPlaceholder(text) {
			hasPlaceholder = true
		}
	}

	if hasPlaceholder {
		return models.StatusPlaceholder
	}
	if hasEmpty || len(translations) < len(s.supportedLangs) {
		return models.StatusIncomplete
	}

	return models.StatusComplete
}

func (s *Service) buildTranslationMetadata(translations map[string]string) models.TranslationMetadata {
	metadata := models.TranslationMetadata{}

	// Calculate character and word counts
	for _, text := range translations {
		metadata.CharacterCount += len(text)
		metadata.WordCount += len(strings.Fields(text))
	}

	return metadata
}

func (s *Service) isPlaceholder(text string) bool {
	placeholderRegex := regexp.MustCompile(`\[(RU|EN|SR)\]`)
	return placeholderRegex.MatchString(text)
}
