// backend/internal/proj/unified/handler/marketplace_handler_test.go
package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/internal/domain/models"
	"backend/internal/proj/unified/service"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockMarketplaceService мок для MarketplaceService
type MockMarketplaceService struct {
	mock.Mock
}

func (m *MockMarketplaceService) CreateListing(ctx context.Context, unified *models.UnifiedListing) (int64, error) {
	args := m.Called(ctx, unified)
	return int64(args.Int(0)), args.Error(1)
}

func (m *MockMarketplaceService) GetListing(ctx context.Context, id int64, sourceType string) (*models.UnifiedListing, error) {
	args := m.Called(ctx, id, sourceType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UnifiedListing), args.Error(1)
}

func (m *MockMarketplaceService) UpdateListing(ctx context.Context, unified *models.UnifiedListing) error {
	args := m.Called(ctx, unified)
	return args.Error(0)
}

func (m *MockMarketplaceService) DeleteListing(ctx context.Context, id int64, sourceType string) error {
	args := m.Called(ctx, id, sourceType)
	return args.Error(0)
}

func (m *MockMarketplaceService) SearchListings(ctx context.Context, params *service.SearchParams) ([]*models.UnifiedListing, int64, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	total := args.Get(1)
	var totalInt64 int64
	if total != nil {
		totalInt64 = total.(int64)
	}
	return args.Get(0).([]*models.UnifiedListing), totalInt64, args.Error(2)
}

// setupTestApp создает тестовое Fiber приложение
func setupTestApp(handler *MarketplaceHandler) *fiber.App {
	app := fiber.New()

	// Mock JWT parser middleware - просто добавляем user_id в locals
	app.Use(func(c *fiber.Ctx) error {
		// Проверяем наличие header Authorization
		authHeader := c.Get("Authorization")
		if authHeader != "" {
			// Для тестов всегда устанавливаем user_id = 1
			c.Locals("user_id", 1)
			c.Locals("email", "test@example.com")
		}
		return c.Next()
	})

	// Регистрируем routes
	app.Get("/api/v1/marketplace/search", handler.SearchListings)
	app.Get("/api/v1/marketplace/listings/:id", handler.GetListing)
	app.Post("/api/v1/marketplace/listings", handler.CreateListing)
	app.Put("/api/v1/marketplace/listings/:id", handler.UpdateListing)
	app.Delete("/api/v1/marketplace/listings/:id", handler.DeleteListing)

	return app
}

// TestCreateListing_C2C тестирует создание C2C listing
func TestCreateListing_C2C(t *testing.T) {
	mockService := new(MockMarketplaceService)
	handler := &MarketplaceHandler{
		service: &service.MarketplaceService{},
		logger:  zerolog.Nop(),
	}
	// Подменяем service на mock
	handler.service = mockService

	app := setupTestApp(handler)

	// Подготовка данных
	listingData := models.UnifiedListing{
		SourceType:  "c2c",
		Title:       "Test C2C Listing",
		Description: "Test description",
		Price:       100.0,
		CategoryID:  1,
		UserID:      1,
	}

	mockService.On("CreateListing", mock.Anything, mock.MatchedBy(func(l *models.UnifiedListing) bool {
		return l.SourceType == "c2c" && l.Title == "Test C2C Listing"
	})).Return(123, nil)

	// Создаем request
	body, _ := json.Marshal(listingData)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/marketplace/listings", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test_token")

	// Выполняем запрос
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	// Проверяем результат
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.True(t, result["success"].(bool))
	assert.Equal(t, float64(123), result["id"].(float64))
	assert.Equal(t, "c2c", result["source_type"].(string))

	mockService.AssertExpectations(t)
}

// TestCreateListing_B2C тестирует создание B2C listing
func TestCreateListing_B2C(t *testing.T) {
	mockService := new(MockMarketplaceService)
	handler := &MarketplaceHandler{
		service: mockService,
		logger:  zerolog.Nop(),
	}

	app := setupTestApp(handler)

	storefrontID := 42
	listingData := models.UnifiedListing{
		SourceType:   "b2c",
		Title:        "Test B2C Product",
		Description:  "Test description",
		Price:        200.0,
		CategoryID:   2,
		UserID:       1,
		StorefrontID: &storefrontID,
	}

	mockService.On("CreateListing", mock.Anything, mock.MatchedBy(func(l *models.UnifiedListing) bool {
		return l.SourceType == "b2c" && l.StorefrontID != nil && *l.StorefrontID == 42
	})).Return(456, nil)

	body, _ := json.Marshal(listingData)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/marketplace/listings", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test_token")

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.True(t, result["success"].(bool))
	assert.Equal(t, float64(456), result["id"].(float64))
	assert.Equal(t, "b2c", result["source_type"].(string))

	mockService.AssertExpectations(t)
}

// TestCreateListing_ValidationErrors тестирует валидацию входных данных
func TestCreateListing_ValidationErrors(t *testing.T) {
	mockService := new(MockMarketplaceService)
	handler := &MarketplaceHandler{
		service: mockService,
		logger:  zerolog.Nop(),
	}

	app := setupTestApp(handler)

	tests := []struct {
		name           string
		data           models.UnifiedListing
		expectedStatus int
		expectedError  string
	}{
		{
			name: "missing source_type",
			data: models.UnifiedListing{
				Title:      "Test",
				Price:      100,
				CategoryID: 1,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "marketplace.source_type_required",
		},
		{
			name: "invalid source_type",
			data: models.UnifiedListing{
				SourceType: "invalid",
				Title:      "Test",
				Price:      100,
				CategoryID: 1,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "marketplace.invalid_source_type",
		},
		{
			name: "missing title",
			data: models.UnifiedListing{
				SourceType: "c2c",
				Price:      100,
				CategoryID: 1,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "marketplace.title_required",
		},
		{
			name: "invalid price",
			data: models.UnifiedListing{
				SourceType: "c2c",
				Title:      "Test",
				Price:      0,
				CategoryID: 1,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "marketplace.invalid_price",
		},
		{
			name: "missing category_id",
			data: models.UnifiedListing{
				SourceType: "c2c",
				Title:      "Test",
				Price:      100,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "marketplace.category_id_required",
		},
		{
			name: "b2c without storefront_id",
			data: models.UnifiedListing{
				SourceType: "b2c",
				Title:      "Test",
				Price:      100,
				CategoryID: 1,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "marketplace.storefront_id_required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.data)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/marketplace/listings", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer test_token")

			resp, err := app.Test(req)
			require.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var result map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&result)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedError, result["error"].(string))
		})
	}
}

// TestGetListing_C2C тестирует получение C2C listing
func TestGetListing_C2C(t *testing.T) {
	mockService := new(MockMarketplaceService)
	handler := &MarketplaceHandler{
		service: mockService,
		logger:  zerolog.Nop(),
	}

	app := setupTestApp(handler)

	expectedListing := &models.UnifiedListing{
		ID:          123,
		SourceType:  "c2c",
		Title:       "Test Listing",
		Price:       100,
		CategoryID:  1,
		UserID:      1,
		Description: "Test description",
	}

	mockService.On("GetListing", mock.Anything, int64(123), "c2c").Return(expectedListing, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/marketplace/listings/123?source_type=c2c", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.True(t, result["success"].(bool))
	data := result["data"].(map[string]interface{})
	assert.Equal(t, float64(123), data["id"].(float64))
	assert.Equal(t, "c2c", data["source_type"].(string))

	mockService.AssertExpectations(t)
}

// TestGetListing_B2C тестирует получение B2C listing
func TestGetListing_B2C(t *testing.T) {
	mockService := new(MockMarketplaceService)
	handler := &MarketplaceHandler{
		service: mockService,
		logger:  zerolog.Nop(),
	}

	app := setupTestApp(handler)

	storefrontID := 42
	expectedListing := &models.UnifiedListing{
		ID:           456,
		SourceType:   "b2c",
		Title:        "Test Product",
		Price:        200,
		CategoryID:   2,
		UserID:       1,
		StorefrontID: &storefrontID,
	}

	mockService.On("GetListing", mock.Anything, int64(456), "b2c").Return(expectedListing, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/marketplace/listings/456?source_type=b2c", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.True(t, result["success"].(bool))
	data := result["data"].(map[string]interface{})
	assert.Equal(t, float64(456), data["id"].(float64))
	assert.Equal(t, "b2c", data["source_type"].(string))

	mockService.AssertExpectations(t)
}

// TestGetListing_NotFound тестирует получение несуществующего listing
func TestGetListing_NotFound(t *testing.T) {
	mockService := new(MockMarketplaceService)
	handler := &MarketplaceHandler{
		service: mockService,
		logger:  zerolog.Nop(),
	}

	app := setupTestApp(handler)

	mockService.On("GetListing", mock.Anything, int64(999), "c2c").Return(nil, fmt.Errorf("not found"))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/marketplace/listings/999?source_type=c2c", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.Equal(t, "marketplace.listing_not_found", result["error"].(string))

	mockService.AssertExpectations(t)
}

// TestUpdateListing тестирует обновление listing
func TestUpdateListing(t *testing.T) {
	mockService := new(MockMarketplaceService)
	handler := &MarketplaceHandler{
		service: mockService,
		logger:  zerolog.Nop(),
	}

	app := setupTestApp(handler)

	existingListing := &models.UnifiedListing{
		ID:         123,
		SourceType: "c2c",
		UserID:     1, // Совпадает с user_id из JWT mock
		Title:      "Old Title",
	}

	updateData := models.UnifiedListing{
		SourceType:  "c2c",
		Title:       "Updated Title",
		Description: "Updated description",
		Price:       150.0,
		CategoryID:  1,
	}

	mockService.On("GetListing", mock.Anything, int64(123), "c2c").Return(existingListing, nil)
	mockService.On("UpdateListing", mock.Anything, mock.MatchedBy(func(l *models.UnifiedListing) bool {
		return l.ID == 123 && l.Title == "Updated Title"
	})).Return(nil)

	body, _ := json.Marshal(updateData)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/marketplace/listings/123", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test_token")

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.True(t, result["success"].(bool))

	mockService.AssertExpectations(t)
}

// TestUpdateListing_PermissionDenied тестирует запрет на обновление чужого listing
func TestUpdateListing_PermissionDenied(t *testing.T) {
	mockService := new(MockMarketplaceService)
	handler := &MarketplaceHandler{
		service: mockService,
		logger:  zerolog.Nop(),
	}

	app := setupTestApp(handler)

	existingListing := &models.UnifiedListing{
		ID:         123,
		SourceType: "c2c",
		UserID:     999, // Не совпадает с user_id из JWT mock (1)
		Title:      "Old Title",
	}

	updateData := models.UnifiedListing{
		SourceType: "c2c",
		Title:      "Updated Title",
		Price:      150.0,
		CategoryID: 1,
	}

	mockService.On("GetListing", mock.Anything, int64(123), "c2c").Return(existingListing, nil)

	body, _ := json.Marshal(updateData)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/marketplace/listings/123", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test_token")

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusForbidden, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.Equal(t, "marketplace.permission_denied", result["error"].(string))

	mockService.AssertExpectations(t)
}

// TestDeleteListing тестирует удаление listing
func TestDeleteListing(t *testing.T) {
	mockService := new(MockMarketplaceService)
	handler := &MarketplaceHandler{
		service: mockService,
		logger:  zerolog.Nop(),
	}

	app := setupTestApp(handler)

	existingListing := &models.UnifiedListing{
		ID:         123,
		SourceType: "c2c",
		UserID:     1, // Совпадает с user_id из JWT mock
	}

	mockService.On("GetListing", mock.Anything, int64(123), "c2c").Return(existingListing, nil)
	mockService.On("DeleteListing", mock.Anything, int64(123), "c2c").Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/marketplace/listings/123?source_type=c2c", nil)
	req.Header.Set("Authorization", "Bearer test_token")

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.True(t, result["success"].(bool))

	mockService.AssertExpectations(t)
}

// TestSearchListings тестирует поиск listings
func TestSearchListings(t *testing.T) {
	mockService := new(MockMarketplaceService)
	handler := &MarketplaceHandler{
		service: mockService,
		logger:  zerolog.Nop(),
	}

	app := setupTestApp(handler)

	expectedListings := []*models.UnifiedListing{
		{
			ID:         1,
			SourceType: "c2c",
			Title:      "Listing 1",
			Price:      100,
		},
		{
			ID:         2,
			SourceType: "b2c",
			Title:      "Product 2",
			Price:      200,
		},
	}

	mockService.On("SearchListings", mock.Anything, mock.MatchedBy(func(p *service.SearchParams) bool {
		return p.Query == "test" && p.CategoryID == 1 && p.Limit == 20
	})).Return(expectedListings, int64(2), nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/marketplace/search?query=test&category_id=1&limit=20", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.True(t, result["success"].(bool))
	assert.Equal(t, float64(2), result["total"].(float64))
	assert.Equal(t, float64(20), result["limit"].(float64))

	data := result["data"].([]interface{})
	assert.Len(t, data, 2)

	mockService.AssertExpectations(t)
}

// TestAuth_Unauthorized тестирует запрос без авторизации
func TestAuth_Unauthorized(t *testing.T) {
	mockService := new(MockMarketplaceService)
	handler := &MarketplaceHandler{
		service: mockService,
		logger:  zerolog.Nop(),
	}

	app := setupTestApp(handler)

	listingData := models.UnifiedListing{
		SourceType: "c2c",
		Title:      "Test",
		Price:      100,
		CategoryID: 1,
	}

	body, _ := json.Marshal(listingData)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/marketplace/listings", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	// НЕ устанавливаем Authorization header

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.Equal(t, "marketplace.unauthorized", result["error"].(string))
}
