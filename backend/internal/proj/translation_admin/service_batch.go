package translation_admin

import (
	"context"
	"fmt"

	"backend/internal/domain/models"
)

// GetBatchLoader returns the batch loader instance
func (s *Service) GetBatchLoader() *BatchLoader {
	return s.batchLoader
}

// GetCostTracker returns the cost tracker instance
func (s *Service) GetCostTracker() *CostTracker {
	return s.costTracker
}

// LoadTranslationsForListings предзагружает переводы для списка объявлений
func (s *Service) LoadTranslationsForListings(ctx context.Context, listings []models.MarketplaceListing, language string) error {
	if s.batchLoader == nil {
		return fmt.Errorf("batch loader not initialized")
	}
	
	return s.batchLoader.PreloadTranslationsForListings(ctx, listings, language)
}

// LoadTranslationsForCategories предзагружает переводы для списка категорий
func (s *Service) LoadTranslationsForCategories(ctx context.Context, categories []models.MarketplaceCategory, languages []string) error {
	if s.batchLoader == nil {
		return fmt.Errorf("batch loader not initialized")
	}
	
	return s.batchLoader.PreloadTranslationsForCategories(ctx, categories, languages)
}

// GetTranslationsBatch получает переводы для множества сущностей одним запросом
func (s *Service) GetTranslationsBatch(ctx context.Context, entityType string, entityIDs []int, language string, fields []string) (map[int]map[string]string, error) {
	if s.batchLoader == nil {
		return nil, fmt.Errorf("batch loader not initialized")
	}
	
	return s.batchLoader.LoadMultipleEntitiesTranslations(ctx, entityType, entityIDs, language, fields)
}

// TrackAIProviderUsage отслеживает использование AI провайдера
func (s *Service) TrackAIProviderUsage(ctx context.Context, provider string, inputTokens, outputTokens, characters int) error {
	if s.costTracker == nil {
		// Если трекер не инициализирован, просто логируем
		s.logger.Debug().
			Str("provider", provider).
			Int("input_tokens", inputTokens).
			Int("output_tokens", outputTokens).
			Int("characters", characters).
			Msg("Cost tracking skipped - tracker not initialized")
		return nil
	}
	
	switch provider {
	case "openai":
		return s.costTracker.TrackOpenAIUsage(ctx, "gpt-3.5-turbo", inputTokens, outputTokens)
	case "google":
		return s.costTracker.TrackGoogleUsage(ctx, characters)
	case "deepl":
		return s.costTracker.TrackDeepLUsage(ctx, characters)
	case "claude":
		return s.costTracker.TrackClaudeUsage(ctx, inputTokens, outputTokens)
	default:
		return fmt.Errorf("unknown provider: %s", provider)
	}
}

// GetAIProviderCosts возвращает расходы по AI провайдерам
func (s *Service) GetAIProviderCosts(ctx context.Context) (map[string]interface{}, error) {
	if s.costTracker == nil {
		return map[string]interface{}{
			"error": "Cost tracking not available",
		}, nil
	}
	
	return s.costTracker.GetCostsSummary(ctx)
}

// GetAIProviderAlerts возвращает алерты о превышении бюджета
func (s *Service) GetAIProviderAlerts(ctx context.Context, dailyLimit, monthlyLimit float64) ([]string, error) {
	if s.costTracker == nil {
		return []string{}, nil
	}
	
	return s.costTracker.GetCostAlerts(ctx, dailyLimit, monthlyLimit)
}

// ResetAIProviderCosts сбрасывает счетчики расходов для провайдера
func (s *Service) ResetAIProviderCosts(ctx context.Context, provider string) error {
	if s.costTracker == nil {
		return fmt.Errorf("cost tracker not initialized")
	}
	
	return s.costTracker.ResetProviderCosts(ctx, provider)
}