// backend/internal/proj/marketplace/handler/handler.go
package handler

import (
	"backend/internal/domain/models"
	globalService "backend/internal/proj/global/service"
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
	Listings        *ListingsHandler
	Images          *ImagesHandler
	Categories      *CategoriesHandler
	Search          *SearchHandler
	Translations    *TranslationsHandler
	Favorites       *FavoritesHandler
	Indexing        *IndexingHandler
	Chat            *ChatHandler
	AdminCategories *AdminCategoriesHandler
	AdminAttributes *AdminAttributesHandler
}

// NewHandler создает новый обработчик маркетплейса
func NewHandler(services globalService.ServicesInterface) *Handler {
	// Сначала создаем базовые обработчики
	categoriesHandler := NewCategoriesHandler(services)

	return &Handler{
		Listings:        NewListingsHandler(services),
		Images:          NewImagesHandler(services),
		Categories:      categoriesHandler,
		Search:          NewSearchHandler(services),
		Translations:    NewTranslationsHandler(services),
		Favorites:       NewFavoritesHandler(services),
		Indexing:        NewIndexingHandler(services),
		Chat:            NewChatHandler(services),
		AdminCategories: NewAdminCategoriesHandler(categoriesHandler),
		AdminAttributes: NewAdminAttributesHandler(services),
	}
}
