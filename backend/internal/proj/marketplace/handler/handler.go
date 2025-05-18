// backend/internal/proj/marketplace/handler/handler.go
package handler

import (
	"backend/internal/domain/models"
	globalService "backend/internal/proj/global/service"
	postgres "backend/internal/storage/postgres"
	"sync"
	"time"
)

// Глобальные переменные для кеширования категорий
var (
	categoryTreeCache      []models.CategoryTreeNode
	categoryTreeLastUpdate time.Time
	categoryTreeMutex      sync.RWMutex
)

// Handler объединяет все обработчики маркетплейса
type Handler struct {
	Listings         *ListingsHandler
	Images           *ImagesHandler
	Categories       *CategoriesHandler
	Search           *SearchHandler
	Translations     *TranslationsHandler
	Favorites        *FavoritesHandler
	Indexing         *IndexingHandler
	Chat             *ChatHandler
	AdminCategories  *AdminCategoriesHandler
	AdminAttributes  *AdminAttributesHandler
	CustomComponents *CustomComponentHandler
	MarketplaceHandler *MarketplaceHandler
}

// NewHandler создает новый обработчик маркетплейса
func NewHandler(services globalService.ServicesInterface) *Handler {
	// Сначала создаем базовые обработчики
	categoriesHandler := NewCategoriesHandler(services)
	
	// Получаем storage из services и создаем хранилище для кастомных компонентов
	marketplaceService := services.Marketplace()
	
	// Приводим storage к postgres.Database для доступа к pool
	if postgresDB, ok := marketplaceService.Storage().(*postgres.Database); ok {
		// Создаем Storage с AttributeGroups
		storage := postgres.NewStorage(postgresDB.GetPool(), services.Translation())
		
		// Создаем MarketplaceHandler
		marketplaceHandler := NewMarketplaceHandler(storage)
		
		customComponentStorage := postgres.NewCustomComponentStorage(postgresDB)
		customComponentHandler := NewCustomComponentHandler(customComponentStorage)
		return &Handler{
			Listings:         NewListingsHandler(services),
			Images:           NewImagesHandler(services),
			Categories:       categoriesHandler,
			Search:           NewSearchHandler(services),
			Translations:     NewTranslationsHandler(services),
			Favorites:        NewFavoritesHandler(services),
			Indexing:         NewIndexingHandler(services),
			Chat:             NewChatHandler(services),
			AdminCategories:  NewAdminCategoriesHandler(categoriesHandler),
			AdminAttributes:  NewAdminAttributesHandler(services),
			CustomComponents: customComponentHandler,
			MarketplaceHandler: marketplaceHandler,
		}
	}

	// Возвращаем handler без CustomComponents, если приведение не удалось
	return &Handler{
		Listings:         NewListingsHandler(services),
		Images:           NewImagesHandler(services),
		Categories:       categoriesHandler,
		Search:           NewSearchHandler(services),
		Translations:     NewTranslationsHandler(services),
		Favorites:        NewFavoritesHandler(services),
		Indexing:         NewIndexingHandler(services),
		Chat:             NewChatHandler(services),
		AdminCategories:  NewAdminCategoriesHandler(categoriesHandler),
		AdminAttributes:  NewAdminAttributesHandler(services),
		CustomComponents: nil,
		MarketplaceHandler: nil,
	}
}
