package translation_admin

import (
	"context"
	"database/sql"
	"regexp"
	"sync"

	"backend/internal/domain/models"
	"backend/internal/proj/translation_admin/cache"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

// Helper functions
func isValidLanguage(lang string) bool {
	validLangs := []string{"ru", "en", "sr"}
	for _, validLang := range validLangs {
		if lang == validLang {
			return true
		}
	}
	return false
}

// isValidModule проверяет валидность названия модуля
func isValidModule(module string) bool {
	// Разрешены только буквы, цифры, дефисы и подчеркивания
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, module)
	return matched && len(module) > 0 && len(module) < 100
}

// Service handles translation admin operations
type Service struct {
	logger             zerolog.Logger
	frontendPath       string
	supportedLangs     []string
	modules            []string
	mutex              sync.RWMutex
	translationRepo    TranslationRepository
	auditRepo          AuditRepository
	cache              *cache.RedisTranslationCache
	batchLoader        *BatchLoader
	costTracker        *CostTracker
	db                 *sql.DB
	translationFactory interface{} // Will be marketplaceService.TranslationFactoryInterface
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
func NewService(ctx context.Context, logger zerolog.Logger, frontendPath string, translationRepo TranslationRepository, auditRepo AuditRepository, redisClient *redis.Client, db *sql.DB, translationFactory interface{}) *Service {
	var translationCache *cache.RedisTranslationCache
	if redisClient != nil {
		translationCache = cache.NewRedisTranslationCache(redisClient)
		logger.Info().Msg("Redis cache enabled for translations")
	} else {
		logger.Warn().Msg("Redis client not provided, caching disabled for translations")
	}

	// Create batch loader with cache
	batchLoader := NewBatchLoader(ctx, translationRepo, translationCache)

	// Create cost tracker with Redis
	costTracker := NewCostTracker(ctx, redisClient)

	return &Service{
		logger:             logger,
		frontendPath:       frontendPath,
		supportedLangs:     []string{"sr", "en", "ru"},
		modules:            []string{"common", "auth", "profile", "marketplace", "admin", "b2c_stores", "cars", "chat", "cart", "checkout", "realEstate", "search", "services", "map", "misc", "notifications", "orders", "products", "reviews"},
		translationRepo:    translationRepo,
		auditRepo:          auditRepo,
		cache:              translationCache,
		batchLoader:        batchLoader,
		costTracker:        costTracker,
		db:                 db,
		translationFactory: translationFactory,
	}
}
