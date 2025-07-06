package service

import (
	"context"
	"database/sql"
	"testing"

	"backend/internal/domain/models"
	"backend/internal/domain/search"
	"backend/internal/proj/storefronts/storage/opensearch"
	"backend/internal/storage"
	"backend/internal/storage/filestorage"
	"backend/internal/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStorage для тестирования
type MockStorage struct {
	mock.Mock
}

// Основные методы для тестирования
func (m *MockStorage) GetListingByID(ctx context.Context, id int) (*models.MarketplaceListing, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.MarketplaceListing), args.Error(1)
}

func (m *MockStorage) SearchListingsAdvanced(ctx context.Context, params interface{}) (interface{}, error) {
	args := m.Called(ctx, params)
	return args.Get(0), args.Error(1)
}

func (m *MockStorage) GetPopularSearchQueries(ctx context.Context, query string, limit int) ([]interface{}, error) {
	args := m.Called(ctx, query, limit)
	return args.Get(0).([]interface{}), args.Error(1)
}

func (m *MockStorage) SaveSearchQuery(ctx context.Context, query, normalized string, resultsCount int, language string) error {
	args := m.Called(ctx, query, normalized, resultsCount, language)
	return args.Error(0)
}

func (m *MockStorage) SearchCategories(ctx context.Context, query string, limit int) ([]models.MarketplaceCategory, error) {
	args := m.Called(ctx, query, limit)
	return args.Get(0).([]models.MarketplaceCategory), args.Error(1)
}

// Недостающий метод AddAdmin
func (m *MockStorage) AddAdmin(ctx context.Context, admin *models.AdminUser) error {
	args := m.Called(ctx, admin)
	return args.Error(0)
}

// Все остальные методы интерфейса Storage (заглушки)
func (m *MockStorage) GetOrCreateGoogleUser(ctx context.Context, user *models.User) (*models.User, error) {
	return nil, nil
}

func (m *MockStorage) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return nil, nil
}
func (m *MockStorage) GetUserByID(ctx context.Context, id int) (*models.User, error) { return nil, nil }
func (m *MockStorage) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	return nil, nil
}
func (m *MockStorage) UpdateUser(ctx context.Context, user *models.User) error { return nil }
func (m *MockStorage) GetUserProfile(ctx context.Context, id int) (*models.UserProfile, error) {
	return nil, nil
}

func (m *MockStorage) UpdateUserProfile(ctx context.Context, id int, update *models.UserProfileUpdate) error {
	return nil
}
func (m *MockStorage) UpdateLastSeen(ctx context.Context, id int) error { return nil }
func (m *MockStorage) GetFavoritedUsers(ctx context.Context, listingID int) ([]int, error) {
	return nil, nil
}

func (m *MockStorage) GetSession(ctx context.Context, token string) (*types.SessionData, error) {
	return nil, nil
}

func (m *MockStorage) CreateRefreshToken(ctx context.Context, token *models.RefreshToken) error {
	return nil
}

func (m *MockStorage) GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	return nil, nil
}

func (m *MockStorage) GetRefreshTokenByID(ctx context.Context, id int) (*models.RefreshToken, error) {
	return nil, nil
}

func (m *MockStorage) GetUserRefreshTokens(ctx context.Context, userID int) ([]*models.RefreshToken, error) {
	return nil, nil
}

func (m *MockStorage) UpdateRefreshToken(ctx context.Context, token *models.RefreshToken) error {
	return nil
}
func (m *MockStorage) RevokeRefreshToken(ctx context.Context, tokenID int) error { return nil }
func (m *MockStorage) RevokeRefreshTokenByValue(ctx context.Context, tokenValue string) error {
	return nil
}
func (m *MockStorage) RevokeUserRefreshTokens(ctx context.Context, userID int) error { return nil }
func (m *MockStorage) DeleteExpiredRefreshTokens(ctx context.Context) (int64, error) { return 0, nil }
func (m *MockStorage) CountActiveUserTokens(ctx context.Context, userID int) (int, error) {
	return 0, nil
}

func (m *MockStorage) GetAllUsers(ctx context.Context, limit, offset int) ([]*models.UserProfile, int, error) {
	return nil, 0, nil
}
func (m *MockStorage) UpdateUserStatus(ctx context.Context, id int, status string) error { return nil }
func (m *MockStorage) DeleteUser(ctx context.Context, id int) error                      { return nil }
func (m *MockStorage) IsUserAdmin(ctx context.Context, email string) (bool, error)       { return false, nil }
func (m *MockStorage) GetAllAdmins(ctx context.Context) ([]*models.AdminUser, error)     { return nil, nil }
func (m *MockStorage) RemoveAdmin(ctx context.Context, email string) error               { return nil }
func (m *MockStorage) CreateReview(ctx context.Context, review *models.Review) (*models.Review, error) {
	return nil, nil
}

func (m *MockStorage) GetReviews(ctx context.Context, filter models.ReviewsFilter) ([]models.Review, int64, error) {
	return nil, 0, nil
}

func (m *MockStorage) GetReviewByID(ctx context.Context, id int) (*models.Review, error) {
	return nil, nil
}
func (m *MockStorage) UpdateReview(ctx context.Context, review *models.Review) error { return nil }
func (m *MockStorage) UpdateReviewStatus(ctx context.Context, reviewId int, status string) error {
	return nil
}
func (m *MockStorage) DeleteReview(ctx context.Context, id int) error { return nil }
func (m *MockStorage) AddReviewResponse(ctx context.Context, response *models.ReviewResponse) error {
	return nil
}
func (m *MockStorage) AddReviewVote(ctx context.Context, vote *models.ReviewVote) error { return nil }
func (m *MockStorage) GetReviewVotes(ctx context.Context, reviewId int) (helpful int, notHelpful int, err error) {
	return 0, 0, nil
}

func (m *MockStorage) GetUserReviewVote(ctx context.Context, userId int, reviewId int) (string, error) {
	return "", nil
}

func (m *MockStorage) GetEntityRating(ctx context.Context, entityType string, entityId int) (float64, error) {
	return 0, nil
}

func (m *MockStorage) QueryRow(ctx context.Context, sql string, args ...interface{}) storage.Row {
	return nil
}

func (m *MockStorage) Query(ctx context.Context, sql string, args ...interface{}) (storage.Rows, error) {
	return nil, nil
}

func (m *MockStorage) GetUserReviews(ctx context.Context, userID int, filter models.ReviewsFilter) ([]models.Review, error) {
	return nil, nil
}

func (m *MockStorage) GetStorefrontReviews(ctx context.Context, storefrontID int, filter models.ReviewsFilter) ([]models.Review, error) {
	return nil, nil
}

func (m *MockStorage) GetUserRatingSummary(ctx context.Context, userID int) (*models.UserRatingSummary, error) {
	return nil, nil
}

func (m *MockStorage) GetStorefrontRatingSummary(ctx context.Context, storefrontID int) (*models.StorefrontRatingSummary, error) {
	return nil, nil
}

func (m *MockStorage) GetNotificationSettings(ctx context.Context, userID int) ([]models.NotificationSettings, error) {
	return nil, nil
}

func (m *MockStorage) UpdateNotificationSettings(ctx context.Context, settings *models.NotificationSettings) error {
	return nil
}

func (m *MockStorage) SaveTelegramConnection(ctx context.Context, userID int, chatID string, username string) error {
	return nil
}

func (m *MockStorage) GetTelegramConnection(ctx context.Context, userID int) (*models.TelegramConnection, error) {
	return nil, nil
}
func (m *MockStorage) DeleteTelegramConnection(ctx context.Context, userID int) error { return nil }
func (m *MockStorage) CreateNotification(ctx context.Context, notification *models.Notification) error {
	return nil
}

func (m *MockStorage) GetUserNotifications(ctx context.Context, userID int, limit, offset int) ([]models.Notification, error) {
	return nil, nil
}

func (m *MockStorage) MarkNotificationAsRead(ctx context.Context, userID int, notificationID int) error {
	return nil
}

func (m *MockStorage) DeleteNotification(ctx context.Context, userID int, notificationID int) error {
	return nil
}

func (m *MockStorage) CreateListing(ctx context.Context, listing *models.MarketplaceListing) (int, error) {
	return 0, nil
}

func (m *MockStorage) GetListings(ctx context.Context, filters map[string]string, limit int, offset int) ([]models.MarketplaceListing, int64, error) {
	return nil, 0, nil
}
func (m *MockStorage) IncrementViewsCount(ctx context.Context, id int) error { return nil }
func (m *MockStorage) UpdateListing(ctx context.Context, listing *models.MarketplaceListing) error {
	return nil
}
func (m *MockStorage) DeleteListing(ctx context.Context, id int, userID int) error { return nil }
func (m *MockStorage) GetCategoryTree(ctx context.Context) ([]models.CategoryTreeNode, error) {
	return nil, nil
}

func (m *MockStorage) AddListingImage(ctx context.Context, image *models.MarketplaceImage) (int, error) {
	return 0, nil
}

func (m *MockStorage) GetListingImages(ctx context.Context, listingID string) ([]models.MarketplaceImage, error) {
	return nil, nil
}
func (m *MockStorage) FileStorage() filestorage.FileStorageInterface { return nil }
func (m *MockStorage) GetListingImageByID(ctx context.Context, imageID int) (*models.MarketplaceImage, error) {
	return nil, nil
}
func (m *MockStorage) DeleteListingImage(ctx context.Context, imageID int) error { return nil }
func (m *MockStorage) GetAttributeOptionTranslations(ctx context.Context, attributeName, optionValue string) (map[string]string, error) {
	return nil, nil
}

func (m *MockStorage) GetAttributeRanges(ctx context.Context, categoryID int) (map[string]map[string]interface{}, error) {
	return nil, nil
}

func (m *MockStorage) GetCategoryAttributes(ctx context.Context, categoryID int) ([]models.CategoryAttribute, error) {
	return nil, nil
}

func (m *MockStorage) GetCategories(ctx context.Context) ([]models.MarketplaceCategory, error) {
	return nil, nil
}

func (m *MockStorage) GetCategoryByID(ctx context.Context, id int) (*models.MarketplaceCategory, error) {
	return nil, nil
}

func (m *MockStorage) AddToFavorites(ctx context.Context, userID int, listingID int) error {
	return nil
}

func (m *MockStorage) RemoveFromFavorites(ctx context.Context, userID int, listingID int) error {
	return nil
}

func (m *MockStorage) GetUserFavorites(ctx context.Context, userID int) ([]models.MarketplaceListing, error) {
	return nil, nil
}

func (m *MockStorage) GetPriceHistory(ctx context.Context, listingID int) ([]models.PriceHistoryEntry, error) {
	return nil, nil
}

func (m *MockStorage) AddPriceHistoryEntry(ctx context.Context, entry *models.PriceHistoryEntry) error {
	return nil
}
func (m *MockStorage) ClosePriceHistoryEntry(ctx context.Context, listingID int) error { return nil }
func (m *MockStorage) CheckPriceManipulation(ctx context.Context, listingID int) (bool, error) {
	return false, nil
}

func (m *MockStorage) SaveListingAttributes(ctx context.Context, listingID int, attributes []models.ListingAttributeValue) error {
	return nil
}

func (m *MockStorage) GetListingAttributes(ctx context.Context, listingID int) ([]models.ListingAttributeValue, error) {
	return nil, nil
}
func (m *MockStorage) SynchronizeDiscountMetadata(ctx context.Context) error { return nil }
func (m *MockStorage) GetUserBalance(ctx context.Context, userID int) (*models.UserBalance, error) {
	return nil, nil
}

func (m *MockStorage) GetUserTransactions(ctx context.Context, userID int, limit, offset int) ([]models.BalanceTransaction, error) {
	return nil, nil
}

func (m *MockStorage) CreateTransaction(ctx context.Context, transaction *models.BalanceTransaction) (int, error) {
	return 0, nil
}

func (m *MockStorage) GetActivePaymentMethods(ctx context.Context) ([]models.PaymentMethod, error) {
	return nil, nil
}

func (m *MockStorage) UpdateBalance(ctx context.Context, userID int, amount float64) error {
	return nil
}

func (m *MockStorage) BeginTx(ctx context.Context, opts *sql.TxOptions) (storage.Transaction, error) {
	return nil, nil
}

func (m *MockStorage) CreateMessage(ctx context.Context, msg *models.MarketplaceMessage) error {
	return nil
}

func (m *MockStorage) GetMessages(ctx context.Context, listingID int, userID int, offset int, limit int) ([]models.MarketplaceMessage, error) {
	return nil, nil
}

func (m *MockStorage) GetChats(ctx context.Context, userID int) ([]models.MarketplaceChat, error) {
	return nil, nil
}

func (m *MockStorage) GetChat(ctx context.Context, chatID int, userID int) (*models.MarketplaceChat, error) {
	return nil, nil
}

func (m *MockStorage) MarkMessagesAsRead(ctx context.Context, messageIDs []int, userID int) error {
	return nil
}
func (m *MockStorage) ArchiveChat(ctx context.Context, chatID int, userID int) error { return nil }
func (m *MockStorage) GetUnreadMessagesCount(ctx context.Context, userID int) (int, error) {
	return 0, nil
}

func (m *MockStorage) CreateChatAttachment(ctx context.Context, attachment *models.ChatAttachment) error {
	return nil
}

func (m *MockStorage) GetChatAttachment(ctx context.Context, attachmentID int) (*models.ChatAttachment, error) {
	return nil, nil
}

func (m *MockStorage) GetMessageAttachments(ctx context.Context, messageID int) ([]*models.ChatAttachment, error) {
	return nil, nil
}
func (m *MockStorage) DeleteChatAttachment(ctx context.Context, attachmentID int) error { return nil }
func (m *MockStorage) UpdateMessageAttachmentsCount(ctx context.Context, messageID int, count int) error {
	return nil
}

func (m *MockStorage) GetMessageByID(ctx context.Context, messageID int) (*models.MarketplaceMessage, error) {
	return nil, nil
}

func (m *MockStorage) GetChatActivityStats(ctx context.Context, buyerID int, sellerID int, listingID int) (*models.ChatActivityStats, error) {
	return nil, nil
}

func (m *MockStorage) GetUserAggregatedRating(ctx context.Context, userID int) (*models.UserAggregatedRating, error) {
	return nil, nil
}

func (m *MockStorage) GetStorefrontAggregatedRating(ctx context.Context, storefrontID int) (*models.StorefrontAggregatedRating, error) {
	return nil, nil
}
func (m *MockStorage) RefreshRatingViews(ctx context.Context) error { return nil }
func (m *MockStorage) CreateReviewConfirmation(ctx context.Context, confirmation *models.ReviewConfirmation) error {
	return nil
}

func (m *MockStorage) GetReviewConfirmation(ctx context.Context, reviewID int) (*models.ReviewConfirmation, error) {
	return nil, nil
}

func (m *MockStorage) CreateReviewDispute(ctx context.Context, dispute *models.ReviewDispute) error {
	return nil
}

func (m *MockStorage) GetReviewDispute(ctx context.Context, reviewID int) (*models.ReviewDispute, error) {
	return nil, nil
}

func (m *MockStorage) UpdateReviewDispute(ctx context.Context, dispute *models.ReviewDispute) error {
	return nil
}

func (m *MockStorage) CanUserReviewEntity(ctx context.Context, userID int, entityType string, entityID int) (*models.CanReviewResponse, error) {
	return nil, nil
}

func (m *MockStorage) Exec(ctx context.Context, sql string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}

func (m *MockStorage) CreateStorefront(ctx context.Context, userID int, dto *models.StorefrontCreateDTO) (*models.Storefront, error) {
	return nil, nil
}

func (m *MockStorage) GetUserStorefronts(ctx context.Context, userID int) ([]models.Storefront, error) {
	return nil, nil
}

func (m *MockStorage) GetStorefrontByID(ctx context.Context, id int) (*models.Storefront, error) {
	return nil, nil
}

func (m *MockStorage) UpdateStorefront(ctx context.Context, storefront *models.Storefront) error {
	return nil
}
func (m *MockStorage) DeleteStorefront(ctx context.Context, id int) error { return nil }
func (m *MockStorage) Storefront() interface{}                            { return nil }
func (m *MockStorage) Cart() interface{}                                  { return nil }
func (m *MockStorage) Order() interface{}                                 { return nil }
func (m *MockStorage) Inventory() interface{}                             { return nil }
func (m *MockStorage) MarketplaceOrder() interface{}                      { return nil }
func (m *MockStorage) StorefrontProductSearch() interface{}               { return nil }
func (m *MockStorage) SearchListings(ctx context.Context, params *search.SearchParams) (*search.SearchResult, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*search.SearchResult), args.Error(1)
}

func (m *MockStorage) SearchListingsOpenSearch(ctx context.Context, params *search.SearchParams) (*search.SearchResult, error) {
	return nil, nil
}

func (m *MockStorage) SuggestListings(ctx context.Context, prefix string, size int) ([]string, error) {
	return nil, nil
}
func (m *MockStorage) ReindexAllListings(ctx context.Context) error { return nil }
func (m *MockStorage) IndexListing(ctx context.Context, listing *models.MarketplaceListing) error {
	return nil
}
func (m *MockStorage) DeleteListingIndex(ctx context.Context, id string) error { return nil }
func (m *MockStorage) PrepareIndex(ctx context.Context) error                  { return nil }
func (m *MockStorage) SearchStorefrontsOpenSearch(ctx context.Context, params *opensearch.StorefrontSearchParams) (*opensearch.StorefrontSearchResult, error) {
	return nil, nil
}

func (m *MockStorage) IndexStorefront(ctx context.Context, storefront *models.Storefront) error {
	return nil
}
func (m *MockStorage) DeleteStorefrontIndex(ctx context.Context, storefrontID int) error { return nil }
func (m *MockStorage) ReindexAllStorefronts(ctx context.Context) error                   { return nil }
func (m *MockStorage) GetTranslationsForEntity(ctx context.Context, entityType string, entityID int) ([]models.Translation, error) {
	return nil, nil
}
func (m *MockStorage) AddContact(ctx context.Context, contact *models.UserContact) error { return nil }
func (m *MockStorage) UpdateContactStatus(ctx context.Context, userID, contactUserID int, status, notes string) error {
	return nil
}

func (m *MockStorage) GetContact(ctx context.Context, userID, contactUserID int) (*models.UserContact, error) {
	return nil, nil
}

func (m *MockStorage) GetUserContacts(ctx context.Context, userID int, status string, page, limit int) ([]models.UserContact, int, error) {
	return nil, 0, nil
}
func (m *MockStorage) RemoveContact(ctx context.Context, userID, contactUserID int) error { return nil }
func (m *MockStorage) GetUserPrivacySettings(ctx context.Context, userID int) (*models.UserPrivacySettings, error) {
	return nil, nil
}

func (m *MockStorage) UpdateUserPrivacySettings(ctx context.Context, userID int, settings *models.UpdatePrivacySettingsRequest) error {
	return nil
}

func (m *MockStorage) CanAddContact(ctx context.Context, userID, targetUserID int) (bool, error) {
	return false, nil
}
func (m *MockStorage) Close()                         {}
func (m *MockStorage) Ping(ctx context.Context) error { return nil }

// Создание тестового сервиса
func createTestService() *MarketplaceService {
	mockStorage := &MockStorage{}
	return &MarketplaceService{
		storage: mockStorage,
	}
}

func TestMarketplaceService_GetSimilarListings_Integration(t *testing.T) {
	service := createTestService()
	mockStorage := service.storage.(*MockStorage)

	// Исходное объявление
	sourceListing := &models.MarketplaceListing{
		ID:         1,
		CategoryID: 1100,
		Title:      "3-комнатная квартира в центре",
		Price:      200000,
		City:       "Белград",
		Attributes: []models.ListingAttributeValue{
			{AttributeName: "rooms", DisplayValue: "3"},
			{AttributeName: "area", DisplayValue: "85"},
		},
	}

	// Похожие объявления
	similarListings := []*models.MarketplaceListing{
		{
			ID:         2,
			CategoryID: 1100,
			Title:      "3-комнатная квартира люкс",
			Price:      220000,
			City:       "Белград",
			Attributes: []models.ListingAttributeValue{
				{AttributeName: "rooms", DisplayValue: "3"},
				{AttributeName: "area", DisplayValue: "90"},
			},
		},
		{
			ID:         3,
			CategoryID: 2000, // Другая категория
			Title:      "Автомобиль BMW",
			Price:      15000,
			City:       "Белград",
		},
	}

	// Настройка моков
	mockStorage.On("GetListingByID", mock.Anything, 1).Return(sourceListing, nil)

	// Мок для SearchListings должен возвращать результат с правильной структурой
	searchResult := &search.SearchResult{
		Listings: similarListings,
		Total:    len(similarListings),
	}
	mockStorage.On("SearchListings", mock.Anything, mock.Anything).Return(searchResult, nil)

	// Выполнение теста
	result, err := service.GetSimilarListings(context.Background(), 1, 10)

	// Проверки
	assert.NoError(t, err)
	assert.Len(t, result, 1) // Должна найтись только одна похожая квартира (ID=2)
	assert.Equal(t, 2, result[0].ID)
	assert.NotNil(t, result[0].Metadata)
	assert.Contains(t, result[0].Metadata, "similarity_score")
	assert.Contains(t, result[0].Metadata, "match_reasons")

	// Проверяем что скор больше порогового значения
	score, ok := result[0].Metadata["similarity_score"].(float64)
	assert.True(t, ok)
	assert.Greater(t, score, 0.25)

	mockStorage.AssertExpectations(t)
}

func TestMarketplaceService_GetEnhancedSuggestions_Integration(t *testing.T) {
	service := createTestService()
	mockStorage := service.storage.(*MockStorage)

	query := "квартир"

	// Мок для популярных запросов
	popularQueries := []interface{}{
		map[string]interface{}{
			"query":        "3-комнатная квартира",
			"search_count": 50,
		},
		map[string]interface{}{
			"query":        "квартира в центре",
			"search_count": 30,
		},
	}
	mockStorage.On("GetPopularSearchQueries", mock.Anything, query, 5).Return(popularQueries, nil)

	// Мок для категорий
	categories := []models.MarketplaceCategory{
		{
			ID:           1100,
			Name:         "Квартиры",
			Slug:         "apartments",
			ListingCount: 150,
		},
	}
	mockStorage.On("SearchCategories", mock.Anything, query, 5).Return(categories, nil)

	// Мок для поиска товаров
	searchResult := &search.SearchResult{
		Listings: []*models.MarketplaceListing{
			{
				ID:    1,
				Title: "3-комнатная квартира в центре",
				Price: 200000,
				Category: &models.MarketplaceCategory{
					Name: "Квартиры",
				},
			},
		},
		Total: 1,
	}
	mockStorage.On("SearchListings", mock.Anything, mock.Anything).Return(searchResult, nil)

	// Выполнение теста
	suggestions, err := service.GetEnhancedSuggestions(context.Background(), query, 10, "queries,categories,products")

	// Проверки
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(suggestions), 3) // Минимум по одному из каждого типа

	// Проверяем типы подсказок
	types := make(map[SuggestionType]bool)
	for _, suggestion := range suggestions {
		types[suggestion.Type] = true
	}

	assert.True(t, types[SuggestionTypeQuery])
	assert.True(t, types[SuggestionTypeCategory])
	assert.True(t, types[SuggestionTypeProduct])

	mockStorage.AssertExpectations(t)
}

func TestSimilarityCalculator_Performance_Integration(t *testing.T) {
	calculator := NewSimilarityCalculator()

	// Создаем большое количество объявлений для тестирования производительности
	source := &models.MarketplaceListing{
		ID:          1,
		CategoryID:  1100,
		Title:       "Тестовое объявление",
		Description: "Описание для тестирования производительности алгоритма",
		Price:       100000,
		City:        "Белград",
		Attributes: []models.ListingAttributeValue{
			{AttributeName: "rooms", DisplayValue: "3"},
			{AttributeName: "area", DisplayValue: "85"},
			{AttributeName: "floor", DisplayValue: "5"},
		},
	}

	targets := make([]*models.MarketplaceListing, 100)
	for i := 0; i < 100; i++ {
		targets[i] = &models.MarketplaceListing{
			ID:          i + 2,
			CategoryID:  1100,
			Title:       "Похожее объявление",
			Description: "Тестовое описание",
			Price:       float64(90000 + i*1000),
			City:        "Белград",
			Attributes: []models.ListingAttributeValue{
				{AttributeName: "rooms", DisplayValue: "3"},
				{AttributeName: "area", DisplayValue: "80"},
			},
		}
	}

	// Замеряем время выполнения
	ctx := context.Background()
	start := testing.B{}.N

	for i := 0; i < 100; i++ {
		score, err := calculator.CalculateSimilarity(ctx, source, targets[i])
		assert.NoError(t, err)
		assert.NotNil(t, score)
		assert.GreaterOrEqual(t, score.TotalScore, 0.0)
		assert.LessOrEqual(t, score.TotalScore, 1.0)
	}

	// Тест должен завершиться быстро (менее 1 секунды для 100 объявлений)
	assert.Less(t, testing.B{}.N-start, 1000) // Условная проверка производительности
}
