package opensearch

import (
    "context"
    "backend/internal/domain/models"
	"backend/internal/domain/search"
 
)

// MarketplaceSearchRepository определяет интерфейс для поисковых операций
type MarketplaceSearchRepository interface {
    // IndexListing индексирует объявление
    IndexListing(ctx context.Context, listing *models.MarketplaceListing) error
    
    // BulkIndexListings индексирует несколько объявлений
    BulkIndexListings(ctx context.Context, listings []*models.MarketplaceListing) error
    
    // DeleteListing удаляет объявление из индекса
    DeleteListing(ctx context.Context, listingID string) error
    
    // SearchListings выполняет поиск объявлений
    SearchListings(ctx context.Context, params *search.SearchParams) (*search.SearchResult, error)
    
    // SuggestListings предлагает автодополнение для поиска
    SuggestListings(ctx context.Context, prefix string, size int) ([]string, error)
    
    // ReindexAll переиндексирует все объявления
    ReindexAll(ctx context.Context) error
    
    // PrepareIndex создает/проверяет индекс
    PrepareIndex(ctx context.Context) error
}
