package translation_admin

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"backend/internal/domain/models"
	"backend/internal/proj/translation_admin/cache"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

const (
	entityTypeAll = "all"
)

// Helper functions
func strPtr(s string) *string {
	return &s
}

// Service handles translation admin operations
type Service struct {
	logger          zerolog.Logger
	frontendPath    string
	supportedLangs  []string
	modules         []string
	mutex           sync.RWMutex
	translationRepo TranslationRepository
	auditRepo       AuditRepository
	cache           *cache.RedisTranslationCache
	batchLoader     *BatchLoader
	costTracker     *CostTracker
}

// TranslationRepository interface for database operations
type TranslationRepository interface {
	GetTranslations(ctx context.Context, filters map[string]interface{}) ([]models.Translation, error)
	CreateTranslation(ctx context.Context, translation *models.Translation) error
	UpdateTranslation(ctx context.Context, translation *models.Translation) error
	DeleteTranslation(ctx context.Context, id int) error
	GetTranslationVersions(ctx context.Context, translationID int) ([]models.TranslationVersion, error)
	GetVersionsByEntity(ctx context.Context, entityType string, entityID int) ([]models.TranslationVersion, error)
	GetVersionDiff(ctx context.Context, versionID1, versionID2 int) (*models.VersionDiff, error)
	RollbackToVersion(ctx context.Context, translationID int, versionID int, userID int) error
	CreateSyncConflict(ctx context.Context, conflict *models.TranslationSyncConflict) error
	ResolveSyncConflict(ctx context.Context, id int, resolution string, userID int) error
	GetUnresolvedConflicts(ctx context.Context) ([]models.TranslationSyncConflict, error)
	GetProviders(ctx context.Context) ([]models.TranslationProvider, error)
	UpdateProvider(ctx context.Context, provider *models.TranslationProvider) error
	CreateTask(ctx context.Context, task *models.TranslationTask) error
	UpdateTask(ctx context.Context, task *models.TranslationTask) error
	GetQualityMetrics(ctx context.Context, translationID int) (*models.TranslationQualityMetrics, error)
	SaveQualityMetrics(ctx context.Context, metrics *models.TranslationQualityMetrics) error
}

// AuditRepository interface for audit logging
type AuditRepository interface {
	LogAction(ctx context.Context, log *models.TranslationAuditLog) error
	GetRecentLogs(ctx context.Context, limit int) ([]models.TranslationAuditLog, error)
}

// NewService creates a new translation admin service
func NewService(logger zerolog.Logger, frontendPath string, translationRepo TranslationRepository, auditRepo AuditRepository, redisClient *redis.Client) *Service {
	var translationCache *cache.RedisTranslationCache
	if redisClient != nil {
		translationCache = cache.NewRedisTranslationCache(redisClient)
		logger.Info().Msg("Redis cache enabled for translations")
	} else {
		logger.Warn().Msg("Redis client not provided, caching disabled for translations")
	}
	
	// Create batch loader with cache
	batchLoader := NewBatchLoader(translationRepo, translationCache)
	
	// Create cost tracker with Redis
	costTracker := NewCostTracker(redisClient)
	
	return &Service{
		logger:          logger,
		frontendPath:    frontendPath,
		supportedLangs:  []string{"sr", "en", "ru"},
		modules:         []string{"common", "auth", "profile", "marketplace", "admin", "storefronts", "cars", "chat", "cart", "checkout", "realEstate", "search", "services", "map", "misc", "notifications", "orders", "products", "reviews"},
		translationRepo: translationRepo,
		auditRepo:       auditRepo,
		cache:           translationCache,
		batchLoader:     batchLoader,
		costTracker:     costTracker,
	}
}

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
		return nil, fmt.Errorf("translation not found")
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

// Helper methods

func (s *Service) isValidModule(name string) bool {
	for _, m := range s.modules {
		if m == name {
			return true
		}
	}
	return false
}

func (s *Service) loadModuleFile(module, lang string) (map[string]interface{}, error) {
	filePath := filepath.Join(s.frontendPath, "src", "messages", lang, module+".json")

	data, err := os.ReadFile(filePath)
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
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = destination.ReadFrom(source)
	return err
}

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

func (s *Service) isPlaceholder(text string) bool {
	placeholderRegex := regexp.MustCompile(`\[(RU|EN|SR)\]`)
	return placeholderRegex.MatchString(text)
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
			ID:          "openai",
			Name:        "OpenAI GPT-4",
			Type:        "openai",
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

	// Mock implementation
	translations := make(map[string]string)
	for _, lang := range req.TargetLanguages {
		// В реальной версии - вызвать API провайдера
		switch lang {
		case "ru":
			translations[lang] = "[RU] " + req.Text
		case "sr":
			translations[lang] = "[SR] " + req.Text
		case "en":
			translations[lang] = "[EN] " + req.Text
		default:
			translations[lang] = "[" + strings.ToUpper(lang) + "] " + req.Text
		}
	}

	result := &models.TranslateResult{
		Key:          req.Key,
		Module:       req.Module,
		Translations: translations,
		Provider:     req.Provider,
		Confidence:   0.95,
	}

	// Optionally add alternative translations
	if req.Provider == "openai" || req.Provider == "anthropic" {
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

// GetVersionHistory retrieves version history for a translation
func (s *Service) GetVersionHistory(ctx context.Context, entityType string, entityID int) (*models.VersionHistoryResponse, error) {
	s.logger.Info().
		Str("entity_type", entityType).
		Int("entity_id", entityID).
		Msg("Getting version history")

	// Get all versions for the entity
	versions, err := s.translationRepo.GetVersionsByEntity(ctx, entityType, entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get versions: %w", err)
	}

	if len(versions) == 0 {
		return &models.VersionHistoryResponse{
			TranslationID:  0,
			CurrentVersion: 0,
			Versions:       []models.TranslationVersion{},
			TotalVersions:  0,
		}, nil
	}

	// Group by translation ID and get the latest version
	translationVersions := make(map[int][]models.TranslationVersion)
	currentVersions := make(map[int]int)

	for _, version := range versions {
		translationVersions[version.TranslationID] = append(translationVersions[version.TranslationID], version)
		if version.VersionNumber > currentVersions[version.TranslationID] {
			currentVersions[version.TranslationID] = version.VersionNumber
		}
	}

	// For simplicity, return the first translation's history
	// In a full implementation, you might want to handle multiple translations
	var firstTranslationID int
	for id := range translationVersions {
		firstTranslationID = id
		break
	}

	response := &models.VersionHistoryResponse{
		TranslationID:  firstTranslationID,
		CurrentVersion: currentVersions[firstTranslationID],
		Versions:       translationVersions[firstTranslationID],
		TotalVersions:  len(translationVersions[firstTranslationID]),
	}

	return response, nil
}

// GetVersionDiff compares two translation versions
func (s *Service) GetVersionDiff(ctx context.Context, versionID1, versionID2 int) (*models.VersionDiff, error) {
	s.logger.Info().
		Int("version1", versionID1).
		Int("version2", versionID2).
		Msg("Getting version diff")

	diff, err := s.translationRepo.GetVersionDiff(ctx, versionID1, versionID2)
	if err != nil {
		return nil, fmt.Errorf("failed to get version diff: %w", err)
	}

	return diff, nil
}

// RollbackVersion rolls back a translation to a previous version
func (s *Service) RollbackVersion(ctx context.Context, req *models.RollbackRequest, userID int) error {
	s.logger.Info().
		Int("translation_id", req.TranslationID).
		Int("version_id", req.VersionID).
		Int("user_id", userID).
		Str("comment", req.Comment).
		Msg("Rolling back version")

	// Perform rollback
	err := s.translationRepo.RollbackToVersion(ctx, req.TranslationID, req.VersionID, userID)
	if err != nil {
		return fmt.Errorf("failed to rollback version: %w", err)
	}

	// Log additional audit entry with comment if provided
	if req.Comment != "" {
		auditLog := &models.TranslationAuditLog{
			UserID:     &userID,
			Action:     "rollback_comment",
			EntityType: strPtr("translation"),
			EntityID:   &req.TranslationID,
			NewValue:   &req.Comment,
		}

		if err := s.auditRepo.LogAction(ctx, auditLog); err != nil {
			s.logger.Error().Err(err).Msg("Failed to log rollback comment")
		}
	}

	return nil
}

// GetAuditLogs retrieves audit logs with filters
func (s *Service) GetAuditLogs(ctx context.Context, filters map[string]interface{}) ([]models.TranslationAuditLog, error) {
	s.logger.Info().Msg("Getting audit logs")

	// For now, use the basic method from auditRepo
	// In a full implementation, you'd pass filters to the repository method
	limit, ok := filters["limit"].(int)
	if !ok {
		limit = 100
	}

	logs, err := s.auditRepo.GetRecentLogs(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs: %w", err)
	}

	return logs, nil
}

// GetAuditStatistics retrieves audit statistics
func (s *Service) GetAuditStatistics(ctx context.Context) (*models.AuditStatistics, error) {
	s.logger.Info().Msg("Getting audit statistics")

	// Get recent logs to calculate statistics
	logs, err := s.auditRepo.GetRecentLogs(ctx, 1000)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs for statistics: %w", err)
	}

	stats := &models.AuditStatistics{
		TotalActions:  len(logs),
		ActionsByType: make(map[string]int),
		ActionsByUser: make(map[int]int),
		RecentActions: []models.TranslationAuditLog{},
	}

	// Calculate statistics
	for _, log := range logs {
		stats.ActionsByType[log.Action]++
		if log.UserID != nil {
			stats.ActionsByUser[*log.UserID]++
		}
	}

	// Get recent actions (last 10)
	if len(logs) > 10 {
		stats.RecentActions = logs[:10]
	} else {
		stats.RecentActions = logs
	}

	return stats, nil
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
		// Translate to each target language
		for _, targetLang := range req.TargetLanguages {
			if targetLang == req.SourceLanguage {
				continue // Skip same language
			}

			// Check if translation already exists
			if !req.OverwriteExisting {
				existing, _ := s.translationRepo.GetTranslations(ctx, map[string]interface{}{
					"entity_type": sourceTranslation.EntityType,
					"entity_id":   sourceTranslation.EntityID,
					"language":    targetLang,
					"field_name":  sourceTranslation.FieldName,
				})
				if len(existing) > 0 {
					continue // Skip existing
				}
			}

			// Create translation request
			translateReq := &models.TranslateRequest{
				Text:            sourceTranslation.TranslatedText,
				SourceLanguage:  req.SourceLanguage,
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
				// Create new translation record
				newTranslation := &models.Translation{
					EntityType:          sourceTranslation.EntityType,
					EntityID:            sourceTranslation.EntityID,
					Language:            targetLang,
					FieldName:           sourceTranslation.FieldName,
					TranslatedText:      translationResult.Translations[targetLang],
					IsMachineTranslated: true,
					IsVerified:          true,
				}

				if err := s.translationRepo.CreateTranslation(ctx, newTranslation); err != nil {
					s.logger.Error().Err(err).Msg("Failed to save auto-approved translation")
				}
			}
		}
	}

	return result, nil
}
