// backend/internal/proj/global/service/unified_search.go
package service

import (
	"context"
)

// UnifiedSearchServiceInterface интерфейс для унифицированного поиска
type UnifiedSearchServiceInterface interface {
	// GetSuggestions возвращает предложения для автодополнения из всех источников
	GetSuggestions(ctx context.Context, prefix string, limit int) ([]string, error)
}

// UnifiedSearchService реализация сервиса унифицированного поиска
type UnifiedSearchService struct {
	services ServicesInterface
}

// NewUnifiedSearchService создает новый сервис унифицированного поиска
func NewUnifiedSearchService(services ServicesInterface) UnifiedSearchServiceInterface {
	return &UnifiedSearchService{
		services: services,
	}
}


// GetSuggestions возвращает предложения для автодополнения из всех источников
func (s *UnifiedSearchService) GetSuggestions(ctx context.Context, prefix string, limit int) ([]string, error) {
	var allSuggestions []string
	
	// Получаем предложения от marketplace
	marketplaceSuggestions, err := s.services.Marketplace().GetSuggestions(ctx, prefix, limit/2)
	if err == nil {
		allSuggestions = append(allSuggestions, marketplaceSuggestions...)
	}
	
	// TODO: Добавить предложения от storefront сервиса
	
	// Убираем дубликаты и ограничиваем количество
	uniqueSuggestions := s.removeDuplicates(allSuggestions)
	if len(uniqueSuggestions) > limit {
		uniqueSuggestions = uniqueSuggestions[:limit]
	}
	
	return uniqueSuggestions, nil
}

// removeDuplicates удаляет дубликаты из списка строк
func (s *UnifiedSearchService) removeDuplicates(strings []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(strings))
	
	for _, str := range strings {
		if !seen[str] {
			seen[str] = true
			result = append(result, str)
		}
	}
	
	return result
}