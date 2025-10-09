package service

import (
	"context"
	"fmt"

	"backend/internal/domain"
	"backend/internal/storage"
	"backend/internal/storage/opensearch"
	"backend/internal/storage/postgres"

	"github.com/jmoiron/sqlx"
)

type Service struct {
	repo         *postgres.SearchConfigRepository
	osClient     *opensearch.OpenSearchClient
	db           *sqlx.DB
	storage      storage.Storage
	b2cIndexName string // Имя индекса для B2C товаров из конфигурации
}

func NewService(db *sqlx.DB, osClient *opensearch.OpenSearchClient, b2cIndexName string) *Service {
	return &Service{
		repo:         postgres.NewSearchConfigRepository(db),
		osClient:     osClient,
		db:           db,
		b2cIndexName: b2cIndexName,
	}
}

// SetStorage устанавливает storage для сервиса
func (s *Service) SetStorage(storage storage.Storage) {
	s.storage = storage
}

// Weight management
func (s *Service) GetWeights(ctx context.Context) ([]domain.SearchWeight, error) {
	return s.repo.GetWeights(ctx)
}

func (s *Service) GetWeightByField(ctx context.Context, fieldName string) (*domain.SearchWeight, error) {
	return s.repo.GetWeightByField(ctx, fieldName)
}

func (s *Service) CreateWeight(ctx context.Context, weight *domain.SearchWeight) error {
	// Validation
	if weight.FieldName == "" {
		return fmt.Errorf("field name is required")
	}
	if weight.Weight <= 0 || weight.Weight > 10 {
		return fmt.Errorf("weight must be between 0 and 10")
	}

	// Check if field already exists
	existing, err := s.repo.GetWeightByField(ctx, weight.FieldName)
	if err != nil {
		return err
	}
	if existing != nil {
		return fmt.Errorf("weight for field %s already exists", weight.FieldName)
	}

	return s.repo.CreateWeight(ctx, weight)
}

func (s *Service) UpdateWeight(ctx context.Context, id int64, weight *domain.SearchWeight) error {
	// Validation
	if weight.FieldName == "" {
		return fmt.Errorf("field name is required")
	}
	if weight.Weight <= 0 || weight.Weight > 10 {
		return fmt.Errorf("weight must be between 0 and 10")
	}

	return s.repo.UpdateWeight(ctx, id, weight)
}

func (s *Service) DeleteWeight(ctx context.Context, id int64) error {
	return s.repo.DeleteWeight(ctx, id)
}

// Synonym management
func (s *Service) GetSynonyms(ctx context.Context, language string) ([]domain.SearchSynonym, error) {
	if language == "" {
		language = "ru"
	}
	return s.repo.GetSynonyms(ctx, language)
}

func (s *Service) CreateSynonym(ctx context.Context, synonym *domain.SearchSynonym) error {
	// Validation
	if synonym.Term == "" {
		return fmt.Errorf("term is required")
	}
	if len(synonym.Synonyms) == 0 {
		return fmt.Errorf("at least one synonym is required")
	}
	if synonym.Language == "" {
		synonym.Language = "ru"
	}

	return s.repo.CreateSynonym(ctx, synonym)
}

func (s *Service) UpdateSynonym(ctx context.Context, id int64, synonym *domain.SearchSynonym) error {
	// Validation
	if synonym.Term == "" {
		return fmt.Errorf("term is required")
	}
	if len(synonym.Synonyms) == 0 {
		return fmt.Errorf("at least one synonym is required")
	}

	return s.repo.UpdateSynonym(ctx, id, synonym)
}

func (s *Service) DeleteSynonym(ctx context.Context, id int64) error {
	return s.repo.DeleteSynonym(ctx, id)
}

// Transliteration management
func (s *Service) GetTransliterationRules(ctx context.Context) ([]domain.TransliterationRule, error) {
	return s.repo.GetTransliterationRules(ctx)
}

func (s *Service) CreateTransliterationRule(ctx context.Context, rule *domain.TransliterationRule) error {
	// Validation
	if rule.FromScript == "" || rule.ToScript == "" {
		return fmt.Errorf("from_script and to_script are required")
	}
	if rule.FromPattern == "" || rule.ToPattern == "" {
		return fmt.Errorf("from_pattern and to_pattern are required")
	}

	return s.repo.CreateTransliterationRule(ctx, rule)
}

func (s *Service) UpdateTransliterationRule(ctx context.Context, id int64, rule *domain.TransliterationRule) error {
	// Validation
	if rule.FromScript == "" || rule.ToScript == "" {
		return fmt.Errorf("from_script and to_script are required")
	}
	if rule.FromPattern == "" || rule.ToPattern == "" {
		return fmt.Errorf("from_pattern and to_pattern are required")
	}

	return s.repo.UpdateTransliterationRule(ctx, id, rule)
}

func (s *Service) DeleteTransliterationRule(ctx context.Context, id int64) error {
	return s.repo.DeleteTransliterationRule(ctx, id)
}

// Statistics
func (s *Service) RecordSearch(ctx context.Context, stats *domain.SearchStatistics) error {
	return s.repo.CreateSearchStatistics(ctx, stats)
}

func (s *Service) GetSearchStatistics(ctx context.Context, limit int) ([]domain.SearchStatistics, error) {
	if limit <= 0 {
		limit = 100
	}

	// Используем старый метод - аналитика из search_logs удалена
	return s.repo.GetSearchStatistics(ctx, limit)
}

func (s *Service) GetPopularSearches(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	if limit <= 0 {
		limit = 10
	}

	// Используем старый метод - аналитика из search_logs удалена
	return s.repo.GetPopularSearches(ctx, limit)
}

// Configuration
func (s *Service) GetConfig(ctx context.Context) (*domain.SearchConfig, error) {
	return s.repo.GetConfig(ctx)
}

func (s *Service) UpdateConfig(ctx context.Context, config *domain.SearchConfig) error {
	// Validation
	if config.MinSearchLength < 1 {
		return fmt.Errorf("min_search_length must be at least 1")
	}
	if config.MaxSuggestions < 1 {
		return fmt.Errorf("max_suggestions must be at least 1")
	}
	if config.FuzzyMaxEdits < 0 || config.FuzzyMaxEdits > 2 {
		return fmt.Errorf("fuzzy_max_edits must be between 0 and 2")
	}

	return s.repo.UpdateConfig(ctx, config)
}

// Аналитика поиска теперь доступна через behavior_tracking модуль
// Используйте /api/v1/analytics/metrics/search для получения метрик поиска
