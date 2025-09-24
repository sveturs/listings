package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"backend/internal/domain/models"
	"backend/internal/proj/marketplace/service"
	"backend/internal/storage"
)

// CarServiceInterface defines the interface for car service operations
type CarServiceInterface interface {
	GetAllCarMakes(ctx context.Context) ([]models.CarMake, error)
	SearchCarMakes(ctx context.Context, query string) ([]models.CarMake, error)
	GetCarModelsByMakeSlug(ctx context.Context, makeSlug string) ([]models.CarModel, error)
	GetCarGenerationsByModelID(ctx context.Context, modelID int) ([]models.CarGeneration, error)
	DecodeVIN(ctx context.Context, vin string) (*models.VINDecodeResult, error)
}

// MockStorage implements storage.Storage interface for testing
type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) GetCarMakes(ctx context.Context) ([]models.CarMake, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.CarMake), args.Error(1)
}

func (m *MockStorage) GetCarMakeBySlug(ctx context.Context, slug string) (*models.CarMake, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CarMake), args.Error(1)
}

func (m *MockStorage) GetCarModelsByMakeSlug(ctx context.Context, makeSlug string) ([]models.CarModel, error) {
	args := m.Called(ctx, makeSlug)
	return args.Get(0).([]models.CarModel), args.Error(1)
}

func (m *MockStorage) GetCarGenerationsByModelID(ctx context.Context, modelID int) ([]models.CarGeneration, error) {
	args := m.Called(ctx, modelID)
	return args.Get(0).([]models.CarGeneration), args.Error(1)
}

func (m *MockStorage) SearchCarMakesByName(ctx context.Context, query string) ([]models.CarMake, error) {
	args := m.Called(ctx, query)
	return args.Get(0).([]models.CarMake), args.Error(1)
}

// Implement other required Storage methods as empty stubs for compilation
func (m *MockStorage) CreateUser(ctx context.Context, user *models.User) error { return nil }
func (m *MockStorage) GetUser(ctx context.Context, userID int) (*models.User, error) { return nil, nil }
func (m *MockStorage) UpdateUser(ctx context.Context, userID int, updates map[string]interface{}) error { return nil }
func (m *MockStorage) DeleteUser(ctx context.Context, userID int) error { return nil }
func (m *MockStorage) GetUserByEmail(ctx context.Context, email string) (*models.User, error) { return nil, nil }
func (m *MockStorage) GetUserByGoogleID(ctx context.Context, googleID string) (*models.User, error) { return nil, nil }
func (m *MockStorage) Ping(ctx context.Context) error { return nil }
func (m *MockStorage) BeginTx(ctx context.Context) (storage.Transaction, error) { return nil, nil }
func (m *MockStorage) GetPool() interface{} { return nil }
func (m *MockStorage) GetCategories(ctx context.Context) ([]models.Category, error) { return nil, nil }
func (m *MockStorage) GetCategory(ctx context.Context, categoryID int) (*models.Category, error) { return nil, nil }
func (m *MockStorage) GetTranslations(ctx context.Context, tableName string, columnName string, recordID int) (map[string]string, error) { return nil, nil }

func TestGetCarMakes(t *testing.T) {
	app := fiber.New()
	mockStorage := new(MockStorage)

	// Create UnifiedCarService with mock storage
	carService := &service.UnifiedCarService{}
	// Note: We need to use reflection or a different approach since UnifiedCarService
	// doesn't expose its storage field. For now, we'll test at a different level.

	handler := &CarsHandler{
		carService: carService,
	}

	app.Get("/cars/makes", handler.GetCarMakes)

	t.Run("Success", func(t *testing.T) {
		expectedMakes := []models.CarMake{
			{ID: 1, Name: "BMW", Slug: "bmw"},
			{ID: 2, Name: "Mercedes", Slug: "mercedes"},
		}

		mockService.On("GetAllCarMakes", mock.Anything).Return(expectedMakes, nil).Once()

		req := httptest.NewRequest("GET", "/cars/makes", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Contains(t, response, "data")
		data := response["data"].([]interface{})
		assert.Len(t, data, 2)

		mockService.AssertExpectations(t)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService.On("GetAllCarMakes", mock.Anything).Return([]models.CarMake{}, errors.New("database error")).Once()

		req := httptest.NewRequest("GET", "/cars/makes", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Contains(t, response, "error")

		mockService.AssertExpectations(t)
	})
}

func TestSearchCarMakes(t *testing.T) {
	app := fiber.New()
	mockService := new(MockUnifiedCarService)
	handler := &CarsHandler{
		carService: mockService,
	}

	app.Get("/cars/makes/search", handler.SearchCarMakes)

	t.Run("Success", func(t *testing.T) {
		expectedMakes := []models.CarMake{
			{ID: 1, Name: "BMW", Slug: "bmw"},
		}

		mockService.On("SearchCarMakes", mock.Anything, "bm").Return(expectedMakes, nil).Once()

		req := httptest.NewRequest("GET", "/cars/makes/search?q=bm", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Contains(t, response, "data")
		data := response["data"].([]interface{})
		assert.Len(t, data, 1)

		mockService.AssertExpectations(t)
	})

	t.Run("EmptyQuery", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/cars/makes/search", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
}

func TestGetCarModels(t *testing.T) {
	app := fiber.New()
	mockService := new(MockUnifiedCarService)
	handler := &CarsHandler{
		carService: mockService,
	}

	app.Get("/cars/makes/:make_slug/models", handler.GetCarModels)

	t.Run("Success", func(t *testing.T) {
		expectedModels := []models.CarModel{
			{ID: 1, MakeID: 1, Name: "3 Series", Slug: "3-series"},
			{ID: 2, MakeID: 1, Name: "5 Series", Slug: "5-series"},
		}

		mockService.On("GetCarModelsByMakeSlug", mock.Anything, "bmw").Return(expectedModels, nil).Once()

		req := httptest.NewRequest("GET", "/cars/makes/bmw/models", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Contains(t, response, "data")
		data := response["data"].([]interface{})
		assert.Len(t, data, 2)

		mockService.AssertExpectations(t)
	})

	t.Run("MakeNotFound", func(t *testing.T) {
		mockService.On("GetCarModelsByMakeSlug", mock.Anything, "invalid-make").Return([]models.CarModel{}, errors.New("make not found")).Once()

		req := httptest.NewRequest("GET", "/cars/makes/invalid-make/models", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		mockService.AssertExpectations(t)
	})
}

func TestGetCarGenerations(t *testing.T) {
	app := fiber.New()
	mockService := new(MockUnifiedCarService)
	handler := &CarsHandler{
		carService: mockService,
	}

	app.Get("/cars/models/:model_id/generations", handler.GetCarGenerations)

	t.Run("Success", func(t *testing.T) {
		expectedGenerations := []models.CarGeneration{
			{ID: 1, ModelID: 1, Name: "E90 (2005-2011)"},
			{ID: 2, ModelID: 1, Name: "F30 (2011-2019)"},
		}

		mockService.On("GetCarGenerationsByModelID", mock.Anything, 1).Return(expectedGenerations, nil).Once()

		req := httptest.NewRequest("GET", "/cars/models/1/generations", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Contains(t, response, "data")
		data := response["data"].([]interface{})
		assert.Len(t, data, 2)

		mockService.AssertExpectations(t)
	})

	t.Run("InvalidModelID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/cars/models/invalid/generations", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
}

func TestDecodeVIN(t *testing.T) {
	app := fiber.New()
	mockService := new(MockUnifiedCarService)
	handler := &CarsHandler{
		carService: mockService,
	}

	app.Get("/cars/vin/:vin/decode", handler.DecodeVIN)

	t.Run("Success", func(t *testing.T) {
		expectedResult := &models.VINDecodeResult{
			VIN:       "WBA3B9G59ENR15340",
			Valid:     true,
			MakeName:  "BMW",
			ModelName: "3 Series",
			Year:      2014,
			Source:    "test",
		}

		mockService.On("DecodeVIN", mock.Anything, "WBA3B9G59ENR15340").Return(expectedResult, nil).Once()

		req := httptest.NewRequest("GET", "/cars/vin/WBA3B9G59ENR15340/decode", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Contains(t, response, "data")
		data := response["data"].(map[string]interface{})
		assert.Equal(t, "BMW", data["make_name"])
		assert.Equal(t, "3 Series", data["model_name"])

		mockService.AssertExpectations(t)
	})

	t.Run("InvalidVIN", func(t *testing.T) {
		mockService.On("DecodeVIN", mock.Anything, "INVALID").Return(nil, errors.New("invalid VIN length")).Once()

		req := httptest.NewRequest("GET", "/cars/vin/INVALID/decode", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Contains(t, response, "error")
		assert.Equal(t, "validation.invalidVIN", response["error"])

		mockService.AssertExpectations(t)
	})

	t.Run("VINDecoderDisabled", func(t *testing.T) {
		mockService.On("DecodeVIN", mock.Anything, "WBA3B9G59ENR15340").Return(nil, errors.New("VIN decoder is disabled")).Once()

		req := httptest.NewRequest("GET", "/cars/vin/WBA3B9G59ENR15340/decode", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Contains(t, response, "error")
		assert.Equal(t, "validation.invalidVIN", response["error"])

		mockService.AssertExpectations(t)
	})

	t.Run("EmptyVIN", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/cars/vin//decode", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})
}

func TestCarsHandler_ServiceNotInitialized(t *testing.T) {
	app := fiber.New()
	handler := &CarsHandler{
		carService: nil, // Service not initialized
	}

	app.Get("/cars/makes", handler.GetCarMakes)

	req := httptest.NewRequest("GET", "/cars/makes", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "general.serviceError", response["error"])
}