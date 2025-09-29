package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"backend/internal/config"
	"backend/internal/domain/models"
	"backend/pkg/utils"
	"backend/internal/storage/postgres"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// UnifiedAttributesTestSuite - интеграционные тесты для API v2
type UnifiedAttributesTestSuite struct {
	suite.Suite
	app     *fiber.App
	handler *UnifiedAttributesHandler
	storage postgres.UnifiedAttributeStorage
	db      *pgxpool.Pool
	cfg     *config.Config
}

// SetupSuite - инициализация тестового окружения
func (s *UnifiedAttributesTestSuite) SetupSuite() {
	// Загружаем тестовую конфигурацию
	s.cfg = &config.Config{
		DatabaseURL: "postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5432/svetubd?sslmode=disable",
		FeatureFlags: &config.FeatureFlags{
			UseUnifiedAttributes:      true,
			UnifiedAttributesFallback: true,
			UnifiedAttributesPercent:  100,
		},
	}

	// Подключаемся к тестовой БД
	var err error
	s.db, err = pgxpool.New(context.Background(), s.cfg.DatabaseURL)
	s.Require().NoError(err)

	// Создаем storage
	s.storage = postgres.NewUnifiedAttributeStorage(s.db, s.cfg.FeatureFlags.UnifiedAttributesFallback)

	// Создаем handler и приложение
	s.handler = NewUnifiedAttributesHandler(s.storage, s.cfg.FeatureFlags)
	s.app = fiber.New()

	// Регистрируем маршруты
	s.setupRoutes()

	// Создаем тестовые данные
	s.setupTestData()
}

// TearDownSuite - очистка после тестов
func (s *UnifiedAttributesTestSuite) TearDownSuite() {
	// Очищаем тестовые данные
	s.cleanupTestData()

	// Закрываем подключение к БД
	if s.db != nil {
		s.db.Close()
	}
}

// setupRoutes - регистрация маршрутов для тестирования
func (s *UnifiedAttributesTestSuite) setupRoutes() {
	api := s.app.Group("/api")
	v2 := api.Group("/v2")

	// Marketplace routes
	marketplace := v2.Group("/marketplace")
	marketplace.Get("/categories/:category_id/attributes", s.handler.GetCategoryAttributes)
	marketplace.Get("/listings/:listing_id/attributes", s.handler.GetListingAttributeValues)
	marketplace.Post("/listings/:listing_id/attributes", s.handler.SaveListingAttributeValues)
	marketplace.Put("/listings/:listing_id/attributes", s.handler.UpdateListingAttributeValues)
	marketplace.Get("/categories/:category_id/attribute-ranges", s.handler.GetAttributeRanges)

	// Admin routes
	admin := v2.Group("/admin")
	admin.Post("/attributes", s.handler.CreateAttribute)
	admin.Put("/attributes/:attribute_id", s.handler.UpdateAttribute)
	admin.Delete("/attributes/:attribute_id", s.handler.DeleteAttribute)
	admin.Post("/categories/:category_id/attributes", s.handler.AttachAttributeToCategory)
	admin.Delete("/categories/:category_id/attributes/:attribute_id", s.handler.DetachAttributeFromCategory)
	admin.Put("/categories/:category_id/attributes/:attribute_id", s.handler.UpdateCategoryAttribute)
	admin.Post("/attributes/migrate", s.handler.MigrateFromLegacy)
	admin.Get("/attributes/migration-status", s.handler.GetMigrationStatus)
}

// setupTestData - создание тестовых данных
func (s *UnifiedAttributesTestSuite) setupTestData() {
	ctx := context.Background()

	// Используем существующую категорию
	categoryID := 1103

	// Создаем тестовые атрибуты
	attrs := []models.UnifiedAttribute{
		{
			Code:          "test_size",
			Name:          "Test Size",
			AttributeType: "select",
			Options:       json.RawMessage(`["S", "M", "L", "XL"]`),
			Purpose:       models.PurposeRegular,
			IsRequired:    true,
		},
		{
			Code:          "test_color",
			Name:          "Test Color",
			AttributeType: "select",
			Options:       json.RawMessage(`["Red", "Blue", "Green"]`),
			Purpose:       models.PurposeRegular,
			IsRequired:    false,
		},
		{
			Code:            "test_price",
			Name:            "Test Price",
			AttributeType:   "number",
			ValidationRules: json.RawMessage(`{"min": 0, "max": 10000}`),
			Purpose:         models.PurposeRegular,
			IsRequired:      true,
		},
	}

	// Создаем атрибуты и привязываем к категории
	for i, attr := range attrs {
		attrID, err := s.storage.CreateAttribute(ctx, &attr)
		s.Require().NoError(err)

		// Привязываем к категории
		settings := &models.UnifiedCategoryAttribute{
			CategoryID:  categoryID,
			AttributeID: attrID,
			IsEnabled:   true,
			IsRequired:  attr.IsRequired,
			IsFilter:    i < 2, // Первые два атрибута как фильтры
			SortOrder:   i + 1,
		}
		err = s.storage.AttachAttributeToCategory(ctx, categoryID, attrID, settings)
		s.Require().NoError(err)
	}
}

// cleanupTestData - очистка тестовых данных
func (s *UnifiedAttributesTestSuite) cleanupTestData() {
	ctx := context.Background()

	// Удаляем тестовые значения атрибутов
	_, err := s.db.Exec(ctx, "DELETE FROM unified_attribute_values WHERE entity_type = 'test'")
	if err != nil {
		s.T().Logf("Failed to cleanup attribute values: %v", err)
	}

	// Удаляем связи категория-атрибут
	_, err = s.db.Exec(ctx, "DELETE FROM unified_category_attributes WHERE category_id = 1103")
	if err != nil {
		s.T().Logf("Failed to cleanup category attributes: %v", err)
	}

	// Удаляем тестовые атрибуты
	_, err = s.db.Exec(ctx, "DELETE FROM unified_attributes WHERE code LIKE 'test_%'")
	if err != nil {
		s.T().Logf("Failed to cleanup attributes: %v", err)
	}
}

// TestGetCategoryAttributes - тест получения атрибутов категории
func (s *UnifiedAttributesTestSuite) TestGetCategoryAttributes() {
	req := httptest.NewRequest(http.MethodGet, "/api/v2/marketplace/categories/1103/attributes", nil)

	resp, err := s.app.Test(req, -1)
	s.Require().NoError(err)
	defer func() { _ = resp.Body.Close() }()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result utils.SuccessResponseSwag
	err = json.NewDecoder(resp.Body).Decode(&result)
	s.Require().NoError(err)

	attrs, ok := result.Data.([]interface{})
	s.Require().True(ok)
	s.GreaterOrEqual(len(attrs), 3) // Минимум 3 тестовых атрибута
}

// TestSaveListingAttributeValues - тест сохранения значений атрибутов
func (s *UnifiedAttributesTestSuite) TestSaveListingAttributeValues() {
	payload := map[string]interface{}{
		"values": map[string]interface{}{
			"1": "L",
			"2": "Blue",
			"3": 500,
		},
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/v2/marketplace/listings/100/attributes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req, -1)
	s.Require().NoError(err)
	defer func() { _ = resp.Body.Close() }()
	s.Equal(http.StatusOK, resp.StatusCode)

	// Проверяем, что значения сохранились
	ctx := context.Background()
	values, err := s.storage.GetAttributeValues(ctx, models.AttributeEntityType("listing"), 100)
	s.Require().NoError(err)
	s.GreaterOrEqual(len(values), 3)
}

// TestValidationRequired - тест валидации обязательных полей
func (s *UnifiedAttributesTestSuite) TestValidationRequired() {
	payload := map[string]interface{}{
		"values": map[string]interface{}{
			"2": "Blue", // Пропускаем обязательные size и price
		},
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/v2/marketplace/listings/101/attributes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req, -1)
	s.Require().NoError(err)
	defer func() { _ = resp.Body.Close() }()
	s.Equal(http.StatusBadRequest, resp.StatusCode)

	var errorResp utils.ErrorResponseSwag
	err = json.NewDecoder(resp.Body).Decode(&errorResp)
	s.Require().NoError(err)
	s.Contains(errorResp.Error, "required")
}

// TestValidationSelectOptions - тест валидации опций select
func (s *UnifiedAttributesTestSuite) TestValidationSelectOptions() {
	payload := map[string]interface{}{
		"values": map[string]interface{}{
			"1": "XXL", // Неверный размер
			"3": 500,
		},
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/v2/marketplace/listings/102/attributes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req, -1)
	s.Require().NoError(err)
	defer func() { _ = resp.Body.Close() }()
	s.Equal(http.StatusBadRequest, resp.StatusCode)

	var errorResp utils.ErrorResponseSwag
	err = json.NewDecoder(resp.Body).Decode(&errorResp)
	s.Require().NoError(err)
	s.Contains(errorResp.Error, "not in allowed options")
}

// TestValidationNumberRange - тест валидации числовых диапазонов
func (s *UnifiedAttributesTestSuite) TestValidationNumberRange() {
	payload := map[string]interface{}{
		"values": map[string]interface{}{
			"1": "L",
			"3": 20000, // Превышает максимум
		},
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/v2/marketplace/listings/103/attributes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req, -1)
	s.Require().NoError(err)
	defer func() { _ = resp.Body.Close() }()
	s.Equal(http.StatusBadRequest, resp.StatusCode)

	var errorResp utils.ErrorResponseSwag
	err = json.NewDecoder(resp.Body).Decode(&errorResp)
	s.Require().NoError(err)
	s.Contains(errorResp.Error, "must not exceed")
}

// TestCreateUpdateDeleteAttribute - тест полного CRUD цикла атрибута
func (s *UnifiedAttributesTestSuite) TestCreateUpdateDeleteAttribute() {
	// 1. Создаем атрибут
	createPayload := map[string]interface{}{
		"code":           "test_crud",
		"name":           "Test CRUD Attribute",
		"attribute_type": "text",
		"purpose":        "regular",
		"is_required":    false,
	}

	body, _ := json.Marshal(createPayload)
	req := httptest.NewRequest(http.MethodPost, "/api/v2/admin/attributes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req, -1)
	s.Require().NoError(err)
	defer func() { _ = resp.Body.Close() }()
	s.Equal(http.StatusOK, resp.StatusCode)

	var createResp utils.SuccessResponseSwag
	err = json.NewDecoder(resp.Body).Decode(&createResp)
	s.Require().NoError(err)

	attrData := createResp.Data.(map[string]interface{})
	attrID := int(attrData["id"].(float64))

	// 2. Обновляем атрибут
	updatePayload := map[string]interface{}{
		"name":        "Updated CRUD Attribute",
		"is_required": true,
	}

	body, _ = json.Marshal(updatePayload)
	req = httptest.NewRequest(http.MethodPut, "/api/v2/admin/attributes/"+strconv.Itoa(attrID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err = s.app.Test(req, -1)
	s.Require().NoError(err)
	defer func() { _ = resp.Body.Close() }()
	s.Equal(http.StatusOK, resp.StatusCode)

	// 3. Удаляем атрибут
	req = httptest.NewRequest(http.MethodDelete, "/api/v2/admin/attributes/"+strconv.Itoa(attrID), nil)

	resp, err = s.app.Test(req, -1)
	s.Require().NoError(err)
	defer func() { _ = resp.Body.Close() }()
	s.Equal(http.StatusOK, resp.StatusCode)

	// Проверяем, что атрибут удален
	ctx := context.Background()
	_, err = s.storage.GetAttribute(ctx, attrID)
	s.Error(err) // Должна быть ошибка "not found"
}

// TestAttachDetachCategoryAttribute - тест привязки/отвязки атрибута от категории
func (s *UnifiedAttributesTestSuite) TestAttachDetachCategoryAttribute() {
	ctx := context.Background()

	// Создаем новый атрибут для теста
	attr := &models.UnifiedAttribute{
		Code:          "test_attach",
		Name:          "Test Attach",
		AttributeType: "text",
		Purpose:       models.PurposeRegular,
	}
	attrID, err := s.storage.CreateAttribute(ctx, attr)
	s.Require().NoError(err)

	// 1. Привязываем к категории
	attachPayload := map[string]interface{}{
		"attribute_id": attrID,
		"is_enabled":   true,
		"is_required":  false,
		"is_filter":    true,
		"sort_order":   10,
	}

	body, _ := json.Marshal(attachPayload)
	req := httptest.NewRequest(http.MethodPost, "/api/v2/admin/categories/2/attributes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req, -1)
	s.Require().NoError(err)
	defer func() { _ = resp.Body.Close() }()
	s.Equal(http.StatusOK, resp.StatusCode)

	// Проверяем привязку
	attrs, err := s.storage.GetCategoryAttributes(ctx, 2)
	s.Require().NoError(err)
	found := false
	for _, a := range attrs {
		if a.ID == attrID {
			found = true
			break
		}
	}
	s.True(found)

	// 2. Отвязываем от категории
	req = httptest.NewRequest(http.MethodDelete, "/api/v2/admin/categories/2/attributes/"+strconv.Itoa(attrID), nil)

	resp, err = s.app.Test(req, -1)
	s.Require().NoError(err)
	defer func() { _ = resp.Body.Close() }()
	s.Equal(http.StatusOK, resp.StatusCode)

	// Проверяем отвязку
	attrs, err = s.storage.GetCategoryAttributes(ctx, 2)
	s.Require().NoError(err)
	found = false
	for _, a := range attrs {
		if a.ID == attrID {
			found = true
			break
		}
	}
	s.False(found)

	// Удаляем тестовый атрибут
	err = s.storage.DeleteAttribute(ctx, attrID)
	s.NoError(err)
}

// TestGetAttributeRanges - тест получения диапазонов значений атрибутов
func (s *UnifiedAttributesTestSuite) TestGetAttributeRanges() {
	ctx := context.Background()

	// Добавим несколько значений для атрибутов
	val100 := 100.0
	val500 := 500.0
	val1000 := 1000.0
	values := []models.UnifiedAttributeValue{
		{
			EntityType:   models.AttributeEntityType("listing"),
			EntityID:     200,
			AttributeID:  3, // test_price
			NumericValue: &val100,
		},
		{
			EntityType:   models.AttributeEntityType("listing"),
			EntityID:     201,
			AttributeID:  3,
			NumericValue: &val500,
		},
		{
			EntityType:   models.AttributeEntityType("listing"),
			EntityID:     202,
			AttributeID:  3,
			NumericValue: &val1000,
		},
	}

	for _, v := range values {
		err := s.storage.SaveAttributeValue(ctx, &v)
		s.Require().NoError(err)
	}

	// Запрашиваем диапазоны
	req := httptest.NewRequest(http.MethodGet, "/api/v2/marketplace/categories/1103/attribute-ranges", nil)

	resp, err := s.app.Test(req, -1)
	s.Require().NoError(err)
	defer func() { _ = resp.Body.Close() }()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result utils.SuccessResponseSwag
	err = json.NewDecoder(resp.Body).Decode(&result)
	s.Require().NoError(err)

	ranges := result.Data.(map[string]interface{})
	s.NotEmpty(ranges)
}

// TestMigrationEndpoints - тест эндпоинтов миграции
func (s *UnifiedAttributesTestSuite) TestMigrationEndpoints() {
	// 1. Запускаем миграцию
	req := httptest.NewRequest(http.MethodPost, "/api/v2/admin/attributes/migrate", nil)

	resp, err := s.app.Test(req, -1)
	s.Require().NoError(err)
	defer func() { _ = resp.Body.Close() }()
	s.Equal(http.StatusOK, resp.StatusCode)

	// 2. Проверяем статус миграции
	req = httptest.NewRequest(http.MethodGet, "/api/v2/admin/attributes/migration-status", nil)

	resp, err = s.app.Test(req, -1)
	s.Require().NoError(err)
	defer func() { _ = resp.Body.Close() }()
	s.Equal(http.StatusOK, resp.StatusCode)

	var statusResp utils.SuccessResponseSwag
	err = json.NewDecoder(resp.Body).Decode(&statusResp)
	s.Require().NoError(err)

	status := statusResp.Data.(map[string]interface{})
	s.Equal("completed", status["status"])

	details := status["details"].(map[string]interface{})
	s.Equal(float64(85), details["attributes_migrated"])
	s.Equal(float64(14), details["categories_processed"])
}

// TestFeatureFlagFallback - тест fallback механизма при отключенных feature flags
func (s *UnifiedAttributesTestSuite) TestFeatureFlagFallback() {
	// Временно отключаем feature flag
	originalFlag := s.cfg.FeatureFlags.UseUnifiedAttributes
	s.cfg.FeatureFlags.UseUnifiedAttributes = false
	defer func() {
		s.cfg.FeatureFlags.UseUnifiedAttributes = originalFlag
	}()

	// Пытаемся получить атрибуты категории
	req := httptest.NewRequest(http.MethodGet, "/api/v2/marketplace/categories/1103/attributes", nil)

	resp, err := s.app.Test(req, -1)
	s.Require().NoError(err)
	defer func() { _ = resp.Body.Close() }()

	// Должен работать через fallback (если настроен) или вернуть ошибку
	if s.cfg.FeatureFlags.UnifiedAttributesFallback {
		s.Equal(http.StatusOK, resp.StatusCode)
	} else {
		s.Equal(http.StatusNotImplemented, resp.StatusCode)
	}
}

// TestConcurrentAccess - тест параллельного доступа к API
func (s *UnifiedAttributesTestSuite) TestConcurrentAccess() {
	concurrency := 10
	done := make(chan bool, concurrency)

	for i := 0; i < concurrency; i++ {
		go func(idx int) {
			defer func() { done <- true }()

			// Каждая горутина делает свой запрос
			req := httptest.NewRequest(http.MethodGet, "/api/v2/marketplace/categories/1103/attributes", nil)

			resp, err := s.app.Test(req, -1)
			assert.NoError(s.T(), err)
			if resp != nil {
				defer func() { _ = resp.Body.Close() }()
			}
			assert.Equal(s.T(), http.StatusOK, resp.StatusCode)
		}(i)
	}

	// Ждем завершения всех горутин
	timeout := time.After(5 * time.Second)
	for i := 0; i < concurrency; i++ {
		select {
		case <-done:
			// Успешно завершено
		case <-timeout:
			s.Fail("Timeout waiting for concurrent requests")
		}
	}
}

// TestDualWriteConsistency - тест консистентности dual-write стратегии
func (s *UnifiedAttributesTestSuite) TestDualWriteConsistency() {
	// Сохраняем значения через v2 API
	payload := map[string]interface{}{
		"values": map[string]interface{}{
			"1": "M",
			"2": "Green",
			"3": 750,
		},
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/v2/marketplace/listings/300/attributes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req, -1)
	s.Require().NoError(err)
	defer func() { _ = resp.Body.Close() }()
	s.Equal(http.StatusOK, resp.StatusCode)

	// Проверяем, что данные записались в новую систему
	ctx := context.Background()
	newValues, err := s.storage.GetAttributeValues(ctx, models.AttributeEntityType("listing"), 300)
	s.Require().NoError(err)
	s.Equal(3, len(newValues))

	// TODO: Проверить запись в старую систему (если dual-write включен)
	// Это требует доступа к старому storage, который нужно будет добавить
}

// TestPerformance - тест производительности API
func (s *UnifiedAttributesTestSuite) TestPerformance() {
	start := time.Now()
	iterations := 100

	for i := 0; i < iterations; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/v2/marketplace/categories/1103/attributes", nil)
		resp, err := s.app.Test(req, -1)
		s.Require().NoError(err)
		s.Equal(http.StatusOK, resp.StatusCode)
		_ = resp.Body.Close()
	}

	elapsed := time.Since(start)
	avgTime := elapsed / time.Duration(iterations)

	// Средний запрос должен выполняться менее 50ms
	s.Less(avgTime, 50*time.Millisecond, "Average request time is too high: %v", avgTime)
}

// TestSuite - запуск тестового набора
func TestUnifiedAttributesIntegration(t *testing.T) {
	// Пропускаем если нет переменной окружения для интеграционных тестов
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	suite.Run(t, new(UnifiedAttributesTestSuite))
}
