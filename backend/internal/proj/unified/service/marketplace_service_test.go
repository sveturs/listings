// backend/internal/proj/unified/service/marketplace_service_test.go
package service

import (
	"context"
	"fmt"
	"testing"

	"backend/internal/domain/adapters"
	"backend/internal/domain/models"
	"backend/internal/domain/search"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockC2CRepository мок для C2C репозитория
type MockC2CRepository struct {
	mock.Mock
}

func (m *MockC2CRepository) CreateListing(ctx context.Context, listing *models.MarketplaceListing) (int, error) {
	args := m.Called(ctx, listing)
	return args.Int(0), args.Error(1)
}

func (m *MockC2CRepository) GetListing(ctx context.Context, id int) (*models.MarketplaceListing, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MarketplaceListing), args.Error(1)
}

func (m *MockC2CRepository) UpdateListing(ctx context.Context, listing *models.MarketplaceListing) error {
	args := m.Called(ctx, listing)
	return args.Error(0)
}

func (m *MockC2CRepository) DeleteListing(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockB2CRepository мок для B2C репозитория
type MockB2CRepository struct {
	mock.Mock
}

func (m *MockB2CRepository) GetStorefrontProductByID(ctx context.Context, productID int) (*models.StorefrontProduct, error) {
	args := m.Called(ctx, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StorefrontProduct), args.Error(1)
}

func (m *MockB2CRepository) CreateStorefrontProduct(ctx context.Context, storefrontID int, req *models.CreateProductRequest) (*models.StorefrontProduct, error) {
	args := m.Called(ctx, storefrontID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StorefrontProduct), args.Error(1)
}

func (m *MockB2CRepository) UpdateStorefrontProduct(ctx context.Context, storefrontID, productID int, req *models.UpdateProductRequest) error {
	args := m.Called(ctx, storefrontID, productID, req)
	return args.Error(0)
}

func (m *MockB2CRepository) DeleteStorefrontProduct(ctx context.Context, storefrontID, productID int) error {
	args := m.Called(ctx, storefrontID, productID)
	return args.Error(0)
}

func (m *MockB2CRepository) GetStorefrontProducts(ctx context.Context, filter models.ProductFilter) ([]*models.StorefrontProduct, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.StorefrontProduct), args.Error(1)
}

func (m *MockB2CRepository) GetStorefrontByID(ctx context.Context, id int) (*models.Storefront, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Storefront), args.Error(1)
}

// MockOpenSearchRepository мок для OpenSearch репозитория
type MockOpenSearchRepository struct {
	mock.Mock
}

func (m *MockOpenSearchRepository) SearchListings(ctx context.Context, params *search.ServiceParams) (*search.ServiceResult, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*search.ServiceResult), args.Error(1)
}

func (m *MockOpenSearchRepository) Index(ctx context.Context, listing *models.MarketplaceListing) error {
	args := m.Called(ctx, listing)
	return args.Error(0)
}

func (m *MockOpenSearchRepository) Delete(ctx context.Context, listingID int) error {
	args := m.Called(ctx, listingID)
	return args.Error(0)
}

// setupTestService создает service с моками для тестирования
func setupTestService() (*MarketplaceService, *MockC2CRepository, *MockB2CRepository, *MockOpenSearchRepository) {
	c2cRepo := new(MockC2CRepository)
	b2cRepo := new(MockB2CRepository)
	osRepo := new(MockOpenSearchRepository)

	logger := zerolog.Nop()

	service := &MarketplaceService{
		c2cRepo:  c2cRepo,
		b2cRepo:  b2cRepo,
		osClient: osRepo,
		logger:   logger,
	}

	return service, c2cRepo, b2cRepo, osRepo
}

// TestCreateListing_C2C тестирует создание C2C listing
func TestCreateListing_C2C(t *testing.T) {
	service, c2cRepo, _, osClient := setupTestService()
	ctx := context.Background()

	unified := &models.UnifiedListing{
		SourceType:   "c2c",
		UserID:       1,
		CategoryID:   100,
		Title:        "Test C2C Listing",
		Description:  "Test description",
		Price:        99.99,
		Condition:    "new",
		Status:       "active",
		Location:     "Test Location",
		ShowOnMap:    true,
		OriginalLang: "ru",
	}

	// Мок: C2C репозиторий возвращает ID
	c2cRepo.On("CreateListing", ctx, mock.AnythingOfType("*models.MarketplaceListing")).Return(123, nil)

	// Мок: OpenSearch индексация успешна
	osClient.On("Index", ctx, mock.AnythingOfType("*models.MarketplaceListing")).Return(nil)

	// Выполнение
	id, err := service.CreateListing(ctx, unified)

	// Проверка
	require.NoError(t, err)
	assert.Equal(t, int64(123), id)
	c2cRepo.AssertExpectations(t)
	osClient.AssertExpectations(t)
}

// TestCreateListing_B2C тестирует создание B2C product
func TestCreateListing_B2C(t *testing.T) {
	service, _, b2cRepo, _ := setupTestService()
	ctx := context.Background()

	storefrontID := 10
	unified := &models.UnifiedListing{
		SourceType:   "b2c",
		UserID:       1,
		CategoryID:   100,
		Title:        "Test B2C Product",
		Description:  "Test description",
		Price:        149.99,
		Condition:    "new",
		Status:       "active",
		ShowOnMap:    true,
		OriginalLang: "sr",
		StorefrontID: &storefrontID,
	}

	expectedProduct := &models.StorefrontProduct{
		ID:            456,
		StorefrontID:  storefrontID,
		Name:          unified.Title,
		Description:   unified.Description,
		Price:         unified.Price,
		CategoryID:    unified.CategoryID,
		StockQuantity: 0,
		StockStatus:   "in_stock",
		IsActive:      true,
		Currency:      "USD",
	}

	// Мок: B2C репозиторий возвращает product
	b2cRepo.On("CreateStorefrontProduct", ctx, storefrontID, mock.AnythingOfType("*models.CreateProductRequest")).Return(expectedProduct, nil)

	// Выполнение
	id, err := service.CreateListing(ctx, unified)

	// Проверка
	require.NoError(t, err)
	assert.Equal(t, int64(456), id)
	b2cRepo.AssertExpectations(t)
}

// TestCreateListing_B2C_NoStorefrontID тестирует ошибку при отсутствии storefront_id
func TestCreateListing_B2C_NoStorefrontID(t *testing.T) {
	service, _, _, _ := setupTestService()
	ctx := context.Background()

	unified := &models.UnifiedListing{
		SourceType:  "b2c",
		UserID:      1,
		CategoryID:  100,
		Title:       "Test B2C Product",
		Description: "Test description",
		Price:       149.99,
		// StorefrontID отсутствует
	}

	// Выполнение
	_, err := service.CreateListing(ctx, unified)

	// Проверка: должна быть ошибка
	require.Error(t, err)
	assert.Contains(t, err.Error(), "storefront_id is required")
}

// TestGetListing_C2C тестирует получение C2C listing
func TestGetListing_C2C(t *testing.T) {
	service, c2cRepo, _, _ := setupTestService()
	ctx := context.Background()

	c2cListing := &models.MarketplaceListing{
		ID:               123,
		UserID:           1,
		CategoryID:       100,
		Title:            "Test C2C Listing",
		Description:      "Test description",
		Price:            99.99,
		Condition:        "new",
		Status:           "active",
		Location:         "Test Location",
		ShowOnMap:        true,
		OriginalLanguage: "ru",
	}

	// Мок: C2C репозиторий возвращает listing
	c2cRepo.On("GetListing", ctx, 123).Return(c2cListing, nil)

	// Выполнение
	unified, err := service.GetListing(ctx, 123, "c2c")

	// Проверка
	require.NoError(t, err)
	require.NotNil(t, unified)
	assert.Equal(t, 123, unified.ID)
	assert.Equal(t, "c2c", unified.SourceType)
	assert.Equal(t, "Test C2C Listing", unified.Title)
	assert.Equal(t, 99.99, unified.Price)
	c2cRepo.AssertExpectations(t)
}

// TestGetListing_B2C тестирует получение B2C product
func TestGetListing_B2C(t *testing.T) {
	service, _, b2cRepo, _ := setupTestService()
	ctx := context.Background()

	storefrontID := 10
	b2cProduct := &models.StorefrontProduct{
		ID:            456,
		StorefrontID:  storefrontID,
		Name:          "Test B2C Product",
		Description:   "Test description",
		Price:         149.99,
		CategoryID:    100,
		StockQuantity: 5,
		StockStatus:   "in_stock",
		IsActive:      true,
		Currency:      "USD",
	}

	storefront := &models.Storefront{
		ID:     storefrontID,
		UserID: 1,
		Name:   "Test Storefront",
	}

	// Мок: B2C репозиторий возвращает product
	b2cRepo.On("GetStorefrontProductByID", ctx, 456).Return(b2cProduct, nil)
	b2cRepo.On("GetStorefrontByID", ctx, storefrontID).Return(storefront, nil)

	// Выполнение
	unified, err := service.GetListing(ctx, 456, "b2c")

	// Проверка
	require.NoError(t, err)
	require.NotNil(t, unified)
	assert.Equal(t, 456, unified.ID)
	assert.Equal(t, "b2c", unified.SourceType)
	assert.Equal(t, "Test B2C Product", unified.Title)
	assert.Equal(t, 149.99, unified.Price)
	assert.Equal(t, &storefrontID, unified.StorefrontID)
	b2cRepo.AssertExpectations(t)
}

// TestUpdateListing_C2C тестирует обновление C2C listing
func TestUpdateListing_C2C(t *testing.T) {
	service, c2cRepo, _, osClient := setupTestService()
	ctx := context.Background()

	unified := &models.UnifiedListing{
		ID:           123,
		SourceType:   "c2c",
		UserID:       1,
		CategoryID:   100,
		Title:        "Updated C2C Listing",
		Description:  "Updated description",
		Price:        129.99,
		Condition:    "used",
		Status:       "active",
		Location:     "New Location",
		ShowOnMap:    true,
		OriginalLang: "ru",
	}

	// Мок: C2C репозиторий успешно обновляет
	c2cRepo.On("UpdateListing", ctx, mock.AnythingOfType("*models.MarketplaceListing")).Return(nil)

	// Мок: OpenSearch переиндексация успешна
	osClient.On("Index", ctx, mock.AnythingOfType("*models.MarketplaceListing")).Return(nil)

	// Выполнение
	err := service.UpdateListing(ctx, unified)

	// Проверка
	require.NoError(t, err)
	c2cRepo.AssertExpectations(t)
	osClient.AssertExpectations(t)
}

// TestUpdateListing_B2C тестирует обновление B2C product
func TestUpdateListing_B2C(t *testing.T) {
	service, _, b2cRepo, _ := setupTestService()
	ctx := context.Background()

	storefrontID := 10
	unified := &models.UnifiedListing{
		ID:           456,
		SourceType:   "b2c",
		UserID:       1,
		CategoryID:   100,
		Title:        "Updated B2C Product",
		Description:  "Updated description",
		Price:        179.99,
		Condition:    "new",
		Status:       "active",
		ShowOnMap:    true,
		OriginalLang: "sr",
		StorefrontID: &storefrontID,
	}

	// Мок: B2C репозиторий успешно обновляет
	b2cRepo.On("UpdateStorefrontProduct", ctx, storefrontID, 456, mock.AnythingOfType("*models.UpdateProductRequest")).Return(nil)

	// Выполнение
	err := service.UpdateListing(ctx, unified)

	// Проверка
	require.NoError(t, err)
	b2cRepo.AssertExpectations(t)
}

// TestDeleteListing_C2C тестирует удаление C2C listing
func TestDeleteListing_C2C(t *testing.T) {
	service, c2cRepo, _, osClient := setupTestService()
	ctx := context.Background()

	// Мок: C2C репозиторий успешно удаляет
	c2cRepo.On("DeleteListing", ctx, 123).Return(nil)

	// Мок: OpenSearch удаление успешно
	osClient.On("Delete", ctx, 123).Return(nil)

	// Выполнение
	err := service.DeleteListing(ctx, 123, "c2c")

	// Проверка
	require.NoError(t, err)
	c2cRepo.AssertExpectations(t)
	osClient.AssertExpectations(t)
}

// TestDeleteListing_B2C тестирует удаление B2C product
func TestDeleteListing_B2C(t *testing.T) {
	service, _, b2cRepo, _ := setupTestService()
	ctx := context.Background()

	storefrontID := 10
	b2cProduct := &models.StorefrontProduct{
		ID:           456,
		StorefrontID: storefrontID,
		Name:         "Test B2C Product",
		IsActive:     true,
	}

	// Мок: B2C репозиторий возвращает product
	b2cRepo.On("GetStorefrontProductByID", ctx, 456).Return(b2cProduct, nil)

	// Мок: B2C репозиторий успешно удаляет
	b2cRepo.On("DeleteStorefrontProduct", ctx, storefrontID, 456).Return(nil)

	// Выполнение
	err := service.DeleteListing(ctx, 456, "b2c")

	// Проверка
	require.NoError(t, err)
	b2cRepo.AssertExpectations(t)
}

// TestSearchListings тестирует поиск listings через OpenSearch
func TestSearchListings(t *testing.T) {
	service, _, _, osClient := setupTestService()
	ctx := context.Background()

	params := &SearchParams{
		Query:      "test",
		CategoryID: 100,
		MinPrice:   50.0,
		MaxPrice:   200.0,
		Limit:      10,
		Offset:     0,
		SourceType: "c2c",
	}

	c2cListings := []*models.MarketplaceListing{
		{
			ID:               123,
			UserID:           1,
			CategoryID:       100,
			Title:            "Test Listing 1",
			Price:            99.99,
			Condition:        "new",
			Status:           "active",
			OriginalLanguage: "ru",
		},
		{
			ID:               124,
			UserID:           2,
			CategoryID:       100,
			Title:            "Test Listing 2",
			Price:            149.99,
			Condition:        "used",
			Status:           "active",
			OriginalLanguage: "ru",
		},
	}

	searchResult := &search.ServiceResult{
		Items: c2cListings,
		Total: 2,
	}

	// Мок: OpenSearch возвращает результаты
	osClient.On("SearchListings", ctx, mock.AnythingOfType("*search.ServiceParams")).Return(searchResult, nil)

	// Выполнение
	results, total, err := service.SearchListings(ctx, params)

	// Проверка
	require.NoError(t, err)
	require.NotNil(t, results)
	assert.Len(t, results, 2)
	assert.Equal(t, int64(2), total)

	// Проверяем что результаты сконвертированы в unified
	assert.Equal(t, "c2c", results[0].SourceType)
	assert.Equal(t, "Test Listing 1", results[0].Title)
	assert.Equal(t, 99.99, results[0].Price)

	osClient.AssertExpectations(t)
}

// TestInvalidSourceType тестирует обработку невалидного source_type
func TestInvalidSourceType(t *testing.T) {
	service, _, _, _ := setupTestService()
	ctx := context.Background()

	unified := &models.UnifiedListing{
		SourceType: "invalid",
		Title:      "Test",
	}

	// CreateListing
	_, err := service.CreateListing(ctx, unified)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid source_type")

	// GetListing
	_, err = service.GetListing(ctx, 123, "invalid")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid source_type")

	// UpdateListing
	unified.ID = 123
	err = service.UpdateListing(ctx, unified)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid source_type")

	// DeleteListing
	err = service.DeleteListing(ctx, 123, "invalid")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid source_type")
}

// TestAdapterIntegration тестирует интеграцию адаптеров
func TestAdapterIntegration(t *testing.T) {
	// Тестируем что адаптеры корректно конвертируют данные
	t.Run("C2C to Unified and back", func(t *testing.T) {
		original := &models.MarketplaceListing{
			ID:               123,
			UserID:           1,
			CategoryID:       100,
			Title:            "Test Listing",
			Description:      "Test description",
			Price:            99.99,
			Condition:        "new",
			Status:           "active",
			Location:         "Test Location",
			ShowOnMap:        true,
			OriginalLanguage: "ru",
		}

		// C2C → Unified
		unified, err := adapters.C2CToUnified(original)
		require.NoError(t, err)
		require.NotNil(t, unified)
		assert.Equal(t, "c2c", unified.SourceType)
		assert.Equal(t, original.Title, unified.Title)
		assert.Equal(t, original.Price, unified.Price)

		// Unified → C2C
		converted, err := adapters.UnifiedToC2C(unified)
		require.NoError(t, err)
		require.NotNil(t, converted)
		assert.Equal(t, original.ID, converted.ID)
		assert.Equal(t, original.Title, converted.Title)
		assert.Equal(t, original.Price, converted.Price)
	})

	t.Run("B2C to Unified and back", func(t *testing.T) {
		original := &models.StorefrontProduct{
			ID:            456,
			StorefrontID:  10,
			Name:          "Test Product",
			Description:   "Test description",
			Price:         149.99,
			CategoryID:    100,
			StockQuantity: 5,
			IsActive:      true,
			Currency:      "USD",
		}

		storefront := &models.Storefront{
			ID:     10,
			UserID: 1,
			Name:   "Test Storefront",
		}

		// B2C → Unified
		unified, err := adapters.B2CToUnified(original, storefront)
		require.NoError(t, err)
		require.NotNil(t, unified)
		assert.Equal(t, "b2c", unified.SourceType)
		assert.Equal(t, original.Name, unified.Title)
		assert.Equal(t, original.Price, unified.Price)

		// Unified → B2C
		converted, err := adapters.UnifiedToB2C(unified)
		require.NoError(t, err)
		require.NotNil(t, converted)
		assert.Equal(t, original.ID, converted.ID)
		assert.Equal(t, original.Name, converted.Name)
		assert.Equal(t, original.Price, converted.Price)
	})
}

// TestGetListing_C2C_NotFound тестирует получение несуществующего C2C listing
func TestGetListing_C2C_NotFound(t *testing.T) {
	service, c2cRepo, _, _ := setupTestService()
	ctx := context.Background()

	// Мок: C2C репозиторий возвращает nil
	c2cRepo.On("GetListing", ctx, 999).Return(nil, nil)

	// Выполнение
	_, err := service.GetListing(ctx, 999, "c2c")

	// Проверка: должна быть ошибка
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	c2cRepo.AssertExpectations(t)
}

// TestGetListing_B2C_NotFound тестирует получение несуществующего B2C product
func TestGetListing_B2C_NotFound(t *testing.T) {
	service, _, b2cRepo, _ := setupTestService()
	ctx := context.Background()

	// Мок: B2C репозиторий возвращает nil
	b2cRepo.On("GetStorefrontProductByID", ctx, 999).Return(nil, nil)

	// Выполнение
	_, err := service.GetListing(ctx, 999, "b2c")

	// Проверка: должна быть ошибка
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	b2cRepo.AssertExpectations(t)
}

// TestSearchListings_NoResults тестирует поиск без результатов
func TestSearchListings_NoResults(t *testing.T) {
	service, _, _, osRepo := setupTestService()
	ctx := context.Background()

	params := &SearchParams{
		Query:      "nonexistent",
		CategoryID: 100,
		Limit:      10,
		Offset:     0,
	}

	searchResult := &search.ServiceResult{
		Items: []*models.MarketplaceListing{},
		Total: 0,
	}

	// Мок: OpenSearch возвращает пустой результат
	osRepo.On("SearchListings", ctx, mock.AnythingOfType("*search.ServiceParams")).Return(searchResult, nil)

	// Выполнение
	results, total, err := service.SearchListings(ctx, params)

	// Проверка
	require.NoError(t, err)
	require.NotNil(t, results)
	assert.Len(t, results, 0)
	assert.Equal(t, int64(0), total)
	osRepo.AssertExpectations(t)
}

// TestSearchListings_Error тестирует обработку ошибки поиска
func TestSearchListings_Error(t *testing.T) {
	service, _, _, osRepo := setupTestService()
	ctx := context.Background()

	params := &SearchParams{
		Query:  "test",
		Limit:  10,
		Offset: 0,
	}

	// Мок: OpenSearch возвращает ошибку
	osRepo.On("SearchListings", ctx, mock.AnythingOfType("*search.ServiceParams")).Return(nil, fmt.Errorf("search error"))

	// Выполнение
	_, _, err := service.SearchListings(ctx, params)

	// Проверка: должна быть ошибка
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to search listings")
	osRepo.AssertExpectations(t)
}

// TestOpenSearchIndexingFailure тестирует обработку ошибки индексации
func TestOpenSearchIndexingFailure(t *testing.T) {
	service, c2cRepo, _, osRepo := setupTestService()
	ctx := context.Background()

	unified := &models.UnifiedListing{
		SourceType:   "c2c",
		UserID:       1,
		CategoryID:   100,
		Title:        "Test Listing",
		Description:  "Test description",
		Price:        99.99,
		Condition:    "new",
		Status:       "active",
		OriginalLang: "ru",
	}

	// Мок: C2C репозиторий успешно создает
	c2cRepo.On("CreateListing", ctx, mock.AnythingOfType("*models.MarketplaceListing")).Return(123, nil)

	// Мок: OpenSearch индексация падает (но не прерывает процесс)
	osRepo.On("Index", ctx, mock.AnythingOfType("*models.MarketplaceListing")).Return(fmt.Errorf("index error"))

	// Выполнение
	id, err := service.CreateListing(ctx, unified)

	// Проверка: listing создан несмотря на ошибку индексации
	require.NoError(t, err)
	assert.Equal(t, int64(123), id)
	c2cRepo.AssertExpectations(t)
	osRepo.AssertExpectations(t)
}
