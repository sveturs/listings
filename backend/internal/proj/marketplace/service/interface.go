package service

import (
	"context"
	"mime/multipart"
	"time"

	"backend/internal/domain/models"
	"backend/internal/domain/search"
	"backend/internal/storage"
)

type MarketplaceServiceInterface interface {
	CreateListing(ctx context.Context, listing *models.MarketplaceListing) (int, error)
	GetListings(ctx context.Context, filters map[string]string, limit int, offset int) ([]models.MarketplaceListing, int64, error)
	GetListingByID(ctx context.Context, id int) (*models.MarketplaceListing, error)
	GetListingBySlug(ctx context.Context, slug string) (*models.MarketplaceListing, error)
	IsSlugAvailable(ctx context.Context, slug string, excludeID int) (bool, error)
	GenerateUniqueSlug(ctx context.Context, baseSlug string, excludeID int) (string, error)
	UpdateListing(ctx context.Context, listing *models.MarketplaceListing) error
	DeleteListing(ctx context.Context, id int, userID int) error
	ProcessImage(file *multipart.FileHeader) (string, error)
	AddListingImage(ctx context.Context, image *models.MarketplaceImage) (int, error)
	GetCategories(ctx context.Context) ([]models.MarketplaceCategory, error)
	GetAllCategories(ctx context.Context) ([]models.MarketplaceCategory, error)
	AddToFavorites(ctx context.Context, userID int, listingID int) error
	RemoveFromFavorites(ctx context.Context, userID int, listingID int) error
	GetCategoryTree(ctx context.Context) ([]models.CategoryTreeNode, error)
	GetUserFavorites(ctx context.Context, userID int) ([]models.MarketplaceListing, error)
	GetFavoritedUsers(ctx context.Context, listingID int) ([]int, error)
	UpdateTranslation(ctx context.Context, translation *models.Translation) error
	UpdateTranslationWithProvider(ctx context.Context, translation *models.Translation, provider TranslationProvider) error
	RefreshCategoryListingCounts(ctx context.Context) error
	GetPriceHistory(ctx context.Context, listingID int) ([]models.PriceHistoryEntry, error)
	SynchronizeDiscountData(ctx context.Context, listingID int) error
	GetSimilarListings(ctx context.Context, listingID int, limit int) ([]*models.MarketplaceListing, error)

	// OpenSearch методы
	SearchListingsAdvanced(ctx context.Context, params *search.ServiceParams) (*search.ServiceResult, error)
	GetSuggestions(ctx context.Context, prefix string, size int) ([]string, error)
	ReindexAllListings(ctx context.Context) error
	GetCategorySuggestions(ctx context.Context, query string, size int) ([]models.CategorySuggestion, error)
	Storage() storage.Storage
	Service() *Service

	// Fuzzy search методы
	ExpandQueryWithSynonyms(ctx context.Context, query string, language string) (string, error)
	SearchCategoriesFuzzy(ctx context.Context, searchTerm string, language string, similarityThreshold float64) ([]CategorySearchResult, error)

	// атрибуты
	GetProductVariantAttributes(ctx context.Context) ([]*models.ProductVariantAttribute, error)
	GetCategoryVariantAttributes(ctx context.Context, categorySlug string) ([]*models.ProductVariantAttribute, error)
	GetCategoryAttributes(ctx context.Context, categoryID int) ([]models.CategoryAttribute, error)
	GetCategoryAttributesWithLang(ctx context.Context, categoryID int, lang string) ([]models.CategoryAttribute, error)
	SaveListingAttributes(ctx context.Context, listingID int, attributes []models.ListingAttributeValue) error
	GetAttributeRanges(ctx context.Context, categoryID int) (map[string]map[string]interface{}, error)

	UploadImage(ctx context.Context, file *multipart.FileHeader, listingID int, isMain bool) (*models.MarketplaceImage, error)
	DeleteImage(ctx context.Context, imageID int) error

	// Новые методы для управления категориями и атрибутами (админка)
	CreateCategory(ctx context.Context, category *models.MarketplaceCategory) (int, error)
	UpdateCategory(ctx context.Context, category *models.MarketplaceCategory) error
	DeleteCategory(ctx context.Context, id int) error
	ReorderCategories(ctx context.Context, orderedIDs []int) error
	MoveCategory(ctx context.Context, id int, newParentID int) error
	GetCategoryByID(ctx context.Context, id int) (*models.MarketplaceCategory, error)

	// Новые методы для управления атрибутами
	CreateAttribute(ctx context.Context, attribute *models.CategoryAttribute) (int, error)
	UpdateAttribute(ctx context.Context, attribute *models.CategoryAttribute) error
	DeleteAttribute(ctx context.Context, id int) error
	GetAttributeByID(ctx context.Context, id int) (*models.CategoryAttribute, error)

	// Новые методы для управления связями
	AddAttributeToCategory(ctx context.Context, categoryID int, attributeID int, isRequired bool) error
	AddAttributeToCategoryWithOrder(ctx context.Context, categoryID int, attributeID int, isRequired bool, sortOrder int) error
	RemoveAttributeFromCategory(ctx context.Context, categoryID int, attributeID int) error
	UpdateAttributeCategory(ctx context.Context, categoryID int, attributeID int, isRequired bool, isEnabled bool) error
	UpdateAttributeCategoryExtended(ctx context.Context, categoryID int, attributeID int, isRequired bool, isEnabled bool, sortOrder int, customComponent string) error
	InvalidateAttributeCache(ctx context.Context, categoryID int) error

	// Поиск и подсказки
	GetEnhancedSuggestions(ctx context.Context, query string, limit int, types string) ([]SuggestionItem, error)
	SaveSearchQuery(ctx context.Context, query string, resultsCount int, language string) error

	// Карта - геопространственные методы
	GetListingsInBounds(ctx context.Context, neLat, neLng, swLat, swLng float64, zoom int, categoryIDs, condition string, minPrice, maxPrice *float64, attributesFilter string) ([]models.MapMarker, error)
	GetMapClusters(ctx context.Context, neLat, neLng, swLat, swLng float64, zoom int, categoryIDs, condition string, minPrice, maxPrice *float64, attributesFilter string) ([]models.MapCluster, error)

	// Методы перевода
	TranslateText(ctx context.Context, text, sourceLanguage, targetLanguage string) (string, error)
	SaveTranslation(ctx context.Context, entityType string, entityID int, language, fieldName, translatedText string, metadata map[string]any) error
	SaveAddressTranslations(ctx context.Context, listingID int, addressFields map[string]string, sourceLanguage string, targetLanguages []string) error

	// Методы для работы с группами атрибутов
	GetCategoryAttributeGroups(ctx context.Context, categoryID int) ([]*models.AttributeGroup, error)
	AttachAttributeGroupToCategory(ctx context.Context, categoryID int, groupID int, sortOrder int) (int, error)
	DetachAttributeGroupFromCategory(ctx context.Context, categoryID int, groupID int) error

	// Методы для работы с автомобильными марками и моделями
	GetCarMakes(ctx context.Context, country string, isDomestic bool, isMotorcycle bool, activeOnly bool) ([]models.CarMake, error)
	GetCarModelsByMake(ctx context.Context, makeSlug string, activeOnly bool) ([]models.CarModel, error)
	GetCarGenerationsByModel(ctx context.Context, modelID int, activeOnly bool) ([]models.CarGeneration, error)
	SearchCarMakes(ctx context.Context, query string, limit int) ([]models.CarMake, error)
}

type ContactsServiceInterface interface {
	AddContact(ctx context.Context, userID int, req *models.AddContactRequest) (*models.UserContact, error)
	UpdateContactStatus(ctx context.Context, userID int, contactUserID int, req *models.UpdateContactRequest) error
	GetContacts(ctx context.Context, userID int, status string, page, limit int) (*models.ContactsListResponse, error)
	RemoveContact(ctx context.Context, userID, contactUserID int) error
	GetPrivacySettings(ctx context.Context, userID int) (*models.UserPrivacySettings, error)
	UpdatePrivacySettings(ctx context.Context, userID int, req *models.UpdatePrivacySettingsRequest) (*models.UserPrivacySettings, error)
	AreContacts(ctx context.Context, userID1, userID2 int) (bool, error)
}

// Interface является алиасом для MarketplaceServiceInterface
type Interface = MarketplaceServiceInterface

// CreateOrderRequest запрос на создание заказа
type CreateOrderRequest struct {
	BuyerID       int64   `json:"buyer_id"`
	ListingID     int64   `json:"listing_id"`
	Message       *string `json:"message,omitempty"`
	PaymentMethod string  `json:"payment_method"`
	ReturnURL     string  `json:"return_url"`
}

// PaymentResult результат создания платежа
type PaymentResult struct {
	TransactionID int64
	PaymentURL    string
	Status        string
}

type OrderServiceInterface interface {
	CreateOrder(ctx context.Context, order *models.MarketplaceOrder) (*models.MarketplaceOrder, error)
	GetOrder(ctx context.Context, orderID int) (*models.MarketplaceOrder, error)
	GetOrdersByUser(ctx context.Context, userID int, isPurchaser bool) ([]models.MarketplaceOrder, error)
	UpdateOrderStatus(ctx context.Context, orderID int, status string) error

	// Методы для работы с handler
	CreateOrderFromRequest(ctx context.Context, req CreateOrderRequest) (*models.MarketplaceOrder, *PaymentResult, error)
	GetBuyerOrders(ctx context.Context, buyerID int64, page, limit int) ([]*models.MarketplaceOrder, int, error)
	GetSellerOrders(ctx context.Context, sellerID int64, page, limit int) ([]*models.MarketplaceOrder, int, error)
	GetOrderDetails(ctx context.Context, orderID int64, userID int64) (*models.MarketplaceOrder, error)
	MarkAsShipped(ctx context.Context, orderID int64, sellerID int64, shippingMethod string, trackingNumber string) error
	ConfirmDelivery(ctx context.Context, orderID int64, buyerID int64) error
	OpenDispute(ctx context.Context, orderID int64, userID int64, reason string) error
	ConfirmPayment(ctx context.Context, orderID int64) error
}

// CacheInterface определяет интерфейс для работы с кешем
type CacheInterface interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, keys ...string) error
	DeletePattern(ctx context.Context, pattern string) error
	GetOrSet(ctx context.Context, key string, dest interface{}, ttl time.Duration, loader func() (interface{}, error)) error
}
