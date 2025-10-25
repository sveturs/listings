package attributes_test

import (
	"context"
	"encoding/json"
	"log"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"backend/internal/proj/delivery/attributes"
	"backend/internal/proj/delivery/models"
	"backend/internal/proj/delivery/storage"
)

type AttributesServiceTestSuite struct {
	suite.Suite
	ctx         context.Context
	pgContainer *postgres.PostgresContainer
	db          *sqlx.DB
	service     *attributes.Service
}

func TestAttributesServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AttributesServiceTestSuite))
}

func (suite *AttributesServiceTestSuite) SetupSuite() {
	if testing.Short() {
		suite.T().Skip("Skipping integration test in short mode")
		return
	}

	ctx := context.Background()
	suite.ctx = ctx

	// Запускаем PostgreSQL контейнер
	pgContainer, err := postgres.Run(ctx, "postgres:16",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Minute)),
	)
	require.NoError(suite.T(), err)
	suite.pgContainer = pgContainer

	// Подключаемся к БД
	connectionString, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(suite.T(), err)

	db, err := sqlx.Connect("pgx", connectionString)
	require.NoError(suite.T(), err)
	suite.db = db

	// Создаем схему БД
	err = suite.createSchema()
	require.NoError(suite.T(), err)

	// Создаем storage и service
	deliveryStorage := storage.NewStorage(suite.db)
	suite.service = attributes.NewService(suite.db, deliveryStorage)
}

func (suite *AttributesServiceTestSuite) TearDownSuite() {
	if suite.db != nil {
		if err := suite.db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}
	if suite.pgContainer != nil {
		if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
			log.Printf("Error terminating container: %v", err)
		}
	}
}

func (suite *AttributesServiceTestSuite) SetupTest() {
	suite.cleanupData()
	suite.insertTestData()
}

func (suite *AttributesServiceTestSuite) cleanupData() {
	queries := []string{
		"DELETE FROM c2c_listings",
		"DELETE FROM b2c_products",
		"DELETE FROM delivery_category_defaults",
	}
	for _, query := range queries {
		_, err := suite.db.ExecContext(suite.ctx, query)
		require.NoError(suite.T(), err)
	}
}

func (suite *AttributesServiceTestSuite) createSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS delivery_category_defaults (
		id SERIAL PRIMARY KEY,
		category_id INTEGER UNIQUE NOT NULL,
		default_weight_kg DECIMAL(10, 2),
		default_length_cm DECIMAL(10, 2),
		default_width_cm DECIMAL(10, 2),
		default_height_cm DECIMAL(10, 2),
		default_packaging_type VARCHAR(50),
		is_typically_fragile BOOLEAN DEFAULT false,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS c2c_listings (
		id SERIAL PRIMARY KEY,
		category_id INTEGER NOT NULL,
		metadata JSONB,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS b2c_products (
		id SERIAL PRIMARY KEY,
		category_id INTEGER NOT NULL,
		attributes JSONB,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := suite.db.ExecContext(suite.ctx, schema)
	return err
}

func (suite *AttributesServiceTestSuite) insertTestData() {
	// Вставляем category defaults
	_, err := suite.db.ExecContext(suite.ctx, `
		INSERT INTO delivery_category_defaults (category_id, default_weight_kg, default_length_cm, default_width_cm, default_height_cm, default_packaging_type, is_typically_fragile)
		VALUES (1, 1.5, 30.0, 20.0, 10.0, 'box', false)
	`)
	require.NoError(suite.T(), err)

	// Вставляем тестовый listing с атрибутами
	attrs := map[string]interface{}{
		"delivery_attributes": map[string]interface{}{
			"weight_kg": 2.0,
			"dimensions": map[string]interface{}{
				"length_cm": 40.0,
				"width_cm":  25.0,
				"height_cm": 15.0,
			},
			"packaging_type": "box",
			"is_fragile":     false,
		},
	}
	attrsJSON, _ := json.Marshal(attrs)

	_, err = suite.db.ExecContext(suite.ctx, `
		INSERT INTO c2c_listings (id, category_id, metadata)
		VALUES (100, 1, $1)
	`, attrsJSON)
	require.NoError(suite.T(), err)

	// Вставляем listing без атрибутов (для тестирования defaults)
	_, err = suite.db.ExecContext(suite.ctx, `
		INSERT INTO c2c_listings (id, category_id, metadata)
		VALUES (101, 1, '{}')
	`)
	require.NoError(suite.T(), err)

	// Вставляем B2C product
	_, err = suite.db.ExecContext(suite.ctx, `
		INSERT INTO b2c_products (id, category_id, attributes)
		VALUES (200, 1, '{}')
	`)
	require.NoError(suite.T(), err)
}

// TestGetProductAttributes tests retrieving product attributes
func (suite *AttributesServiceTestSuite) TestGetProductAttributes_ExistingAttributes() {
	attrs, err := suite.service.GetProductAttributes(suite.ctx, 100, attributes.ProductTypeListing)
	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), attrs)
	assert.Equal(suite.T(), 2.0, attrs.WeightKg)
	assert.NotNil(suite.T(), attrs.Dimensions)
	assert.Equal(suite.T(), 40.0, attrs.Dimensions.LengthCm)
	assert.Equal(suite.T(), 25.0, attrs.Dimensions.WidthCm)
	assert.Equal(suite.T(), 15.0, attrs.Dimensions.HeightCm)
	assert.Equal(suite.T(), "box", attrs.PackagingType)
}

// TestGetProductAttributes_UsesCategoryDefaults tests fallback to category defaults
func (suite *AttributesServiceTestSuite) TestGetProductAttributes_UsesCategoryDefaults() {
	attrs, err := suite.service.GetProductAttributes(suite.ctx, 101, attributes.ProductTypeListing)
	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), attrs)
	// Должны использоваться дефолты категории
	assert.Equal(suite.T(), 1.5, attrs.WeightKg)
	assert.NotNil(suite.T(), attrs.Dimensions)
	assert.Equal(suite.T(), 30.0, attrs.Dimensions.LengthCm)
}

// TestGetProductAttributes_NotFound tests product not found
func (suite *AttributesServiceTestSuite) TestGetProductAttributes_NotFound() {
	attrs, err := suite.service.GetProductAttributes(suite.ctx, 999, attributes.ProductTypeListing)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), attrs)
	assert.Contains(suite.T(), err.Error(), "product not found")
}

// TestUpdateProductAttributes tests updating product attributes
func (suite *AttributesServiceTestSuite) TestUpdateProductAttributes() {
	newAttrs := &models.DeliveryAttributes{
		WeightKg: 3.5,
		Dimensions: &models.Dimensions{
			LengthCm: 50.0,
			WidthCm:  30.0,
			HeightCm: 20.0,
		},
		PackagingType: "envelope",
		IsFragile:     true,
	}

	err := suite.service.UpdateProductAttributes(suite.ctx, 100, attributes.ProductTypeListing, newAttrs)
	require.NoError(suite.T(), err)

	// Проверяем обновление
	attrs, err := suite.service.GetProductAttributes(suite.ctx, 100, attributes.ProductTypeListing)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), 3.5, attrs.WeightKg)
	assert.Equal(suite.T(), 50.0, attrs.Dimensions.LengthCm)
	assert.Equal(suite.T(), "envelope", attrs.PackagingType)
	assert.True(suite.T(), attrs.IsFragile)
	// Проверяем что объем рассчитан
	assert.Greater(suite.T(), attrs.VolumeM3, 0.0)
}

// TestUpdateProductAttributes_InvalidWeight tests validation
func (suite *AttributesServiceTestSuite) TestUpdateProductAttributes_InvalidWeight() {
	invalidAttrs := &models.DeliveryAttributes{
		WeightKg:      -1.0, // Отрицательный вес
		PackagingType: "box",
	}

	err := suite.service.UpdateProductAttributes(suite.ctx, 100, attributes.ProductTypeListing, invalidAttrs)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "weight cannot be negative")
}

// TestUpdateProductAttributes_ExceedsMaxWeight tests max weight validation
func (suite *AttributesServiceTestSuite) TestUpdateProductAttributes_ExceedsMaxWeight() {
	invalidAttrs := &models.DeliveryAttributes{
		WeightKg:      1001.0, // Превышает максимум 1000 кг
		PackagingType: "box",
	}

	err := suite.service.UpdateProductAttributes(suite.ctx, 100, attributes.ProductTypeListing, invalidAttrs)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "weight exceeds maximum")
}

// TestUpdateProductAttributes_InvalidDimensions tests dimension validation
func (suite *AttributesServiceTestSuite) TestUpdateProductAttributes_InvalidDimensions() {
	invalidAttrs := &models.DeliveryAttributes{
		WeightKg: 2.0,
		Dimensions: &models.Dimensions{
			LengthCm: -10.0, // Отрицательная длина
			WidthCm:  20.0,
			HeightCm: 15.0,
		},
		PackagingType: "box",
	}

	err := suite.service.UpdateProductAttributes(suite.ctx, 100, attributes.ProductTypeListing, invalidAttrs)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "dimensions cannot be negative")
}

// TestUpdateProductAttributes_InvalidPackagingType tests packaging type validation
func (suite *AttributesServiceTestSuite) TestUpdateProductAttributes_InvalidPackagingType() {
	invalidAttrs := &models.DeliveryAttributes{
		WeightKg:      2.0,
		PackagingType: "invalid_type",
	}

	err := suite.service.UpdateProductAttributes(suite.ctx, 100, attributes.ProductTypeListing, invalidAttrs)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "invalid packaging type")
}

// TestGetCategoryDefaults tests retrieving category defaults
func (suite *AttributesServiceTestSuite) TestGetCategoryDefaults_Exists() {
	defaults, err := suite.service.GetCategoryDefaults(suite.ctx, 1)
	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), defaults)
	assert.Equal(suite.T(), 1, defaults.CategoryID)
	assert.NotNil(suite.T(), defaults.DefaultWeightKg)
	assert.Equal(suite.T(), 1.5, *defaults.DefaultWeightKg)
}

// TestGetCategoryDefaults_NotFound tests category defaults not found
func (suite *AttributesServiceTestSuite) TestGetCategoryDefaults_NotFound() {
	defaults, err := suite.service.GetCategoryDefaults(suite.ctx, 999)
	require.NoError(suite.T(), err) // Должен вернуть пустые defaults без ошибки
	assert.NotNil(suite.T(), defaults)
	assert.Equal(suite.T(), 999, defaults.CategoryID)
	assert.Nil(suite.T(), defaults.DefaultWeightKg)
}

// TestUpdateCategoryDefaults tests updating category defaults
func (suite *AttributesServiceTestSuite) TestUpdateCategoryDefaults() {
	weight := 2.5
	length := 35.0
	width := 25.0
	height := 12.0
	packagingType := "envelope"

	defaults := &models.CategoryDefaults{
		CategoryID:           2,
		DefaultWeightKg:      &weight,
		DefaultLengthCm:      &length,
		DefaultWidthCm:       &width,
		DefaultHeightCm:      &height,
		DefaultPackagingType: &packagingType,
		IsTypicallyFragile:   true,
	}

	err := suite.service.UpdateCategoryDefaults(suite.ctx, defaults)
	require.NoError(suite.T(), err)
	assert.Greater(suite.T(), defaults.ID, 0)

	// Проверяем что сохранилось
	saved, err := suite.service.GetCategoryDefaults(suite.ctx, 2)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2.5, *saved.DefaultWeightKg)
	assert.Equal(suite.T(), "envelope", *saved.DefaultPackagingType)
	assert.True(suite.T(), saved.IsTypicallyFragile)
}

// TestUpdateCategoryDefaults_InvalidWeight tests validation
func (suite *AttributesServiceTestSuite) TestUpdateCategoryDefaults_InvalidWeight() {
	weight := -1.0
	defaults := &models.CategoryDefaults{
		CategoryID:      3,
		DefaultWeightKg: &weight,
	}

	err := suite.service.UpdateCategoryDefaults(suite.ctx, defaults)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "invalid default weight")
}

// TestApplyCategoryDefaultsToProducts tests applying defaults to products
func (suite *AttributesServiceTestSuite) TestApplyCategoryDefaultsToProducts() {
	// Создаем еще listings без атрибутов
	_, err := suite.db.ExecContext(suite.ctx, `
		INSERT INTO c2c_listings (id, category_id, metadata)
		VALUES (102, 1, '{}'), (103, 1, '{}')
	`)
	require.NoError(suite.T(), err)

	// Применяем дефолты
	count, err := suite.service.ApplyCategoryDefaultsToProducts(suite.ctx, 1)
	require.NoError(suite.T(), err)
	assert.GreaterOrEqual(suite.T(), count, 2) // Минимум 2 новых + 1 из SetupTest

	// Проверяем что атрибуты установлены
	attrs, err := suite.service.GetProductAttributes(suite.ctx, 102, attributes.ProductTypeListing)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1.5, attrs.WeightKg) // Из defaults
	assert.Equal(suite.T(), 30.0, attrs.Dimensions.LengthCm)
}

// TestApplyCategoryDefaultsToProducts_NoDefaults tests error when no defaults
func (suite *AttributesServiceTestSuite) TestApplyCategoryDefaultsToProducts_NoDefaults() {
	count, err := suite.service.ApplyCategoryDefaultsToProducts(suite.ctx, 999)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), 0, count)
	assert.Contains(suite.T(), err.Error(), "no defaults configured")
}

// TestCalculateVolumetricWeight tests volumetric weight calculation
func (suite *AttributesServiceTestSuite) TestCalculateVolumetricWeight() {
	dims := &models.Dimensions{
		LengthCm: 50.0,
		WidthCm:  40.0,
		HeightCm: 30.0,
	}

	divisor := 5000.0 // Стандартный делитель для курьерских служб
	volumetricWeight := suite.service.CalculateVolumetricWeight(dims, divisor)

	expected := (50.0 * 40.0 * 30.0) / 5000.0 // = 12.0 kg
	assert.Equal(suite.T(), expected, volumetricWeight)
}

// TestCalculateVolumetricWeight_NilDimensions tests nil dimensions
func (suite *AttributesServiceTestSuite) TestCalculateVolumetricWeight_NilDimensions() {
	volumetricWeight := suite.service.CalculateVolumetricWeight(nil, 5000.0)
	assert.Equal(suite.T(), 0.0, volumetricWeight)
}

// TestGetEffectiveWeight tests effective weight calculation
func (suite *AttributesServiceTestSuite) TestGetEffectiveWeight_RealWeightHigher() {
	attrs := &models.DeliveryAttributes{
		WeightKg: 15.0, // Реальный вес больше объемного
		Dimensions: &models.Dimensions{
			LengthCm: 30.0,
			WidthCm:  20.0,
			HeightCm: 10.0,
		},
	}

	effectiveWeight := suite.service.GetEffectiveWeight(attrs, 5000.0)
	assert.Equal(suite.T(), 15.0, effectiveWeight) // Должен вернуть реальный вес
}

// TestGetEffectiveWeight_VolumetricWeightHigher tests volumetric weight higher
func (suite *AttributesServiceTestSuite) TestGetEffectiveWeight_VolumetricWeightHigher() {
	attrs := &models.DeliveryAttributes{
		WeightKg: 1.0, // Реальный вес меньше объемного
		Dimensions: &models.Dimensions{
			LengthCm: 100.0,
			WidthCm:  80.0,
			HeightCm: 60.0,
		},
	}

	effectiveWeight := suite.service.GetEffectiveWeight(attrs, 5000.0)
	// Объемный вес = (100 * 80 * 60) / 5000 = 96 kg
	assert.Equal(suite.T(), 96.0, effectiveWeight)
}

// TestBatchUpdateProductAttributes tests batch update
func (suite *AttributesServiceTestSuite) TestBatchUpdateProductAttributes() {
	// Создаем несколько products
	_, err := suite.db.ExecContext(suite.ctx, `
		INSERT INTO c2c_listings (id, category_id, metadata)
		VALUES (110, 1, '{}'), (111, 1, '{}')
	`)
	require.NoError(suite.T(), err)

	updates := []attributes.ProductAttributesUpdate{
		{
			ProductID:   110,
			ProductType: attributes.ProductTypeListing,
			Attributes: &models.DeliveryAttributes{
				WeightKg:      2.0,
				PackagingType: "box",
			},
		},
		{
			ProductID:   111,
			ProductType: attributes.ProductTypeListing,
			Attributes: &models.DeliveryAttributes{
				WeightKg:      3.0,
				PackagingType: "envelope",
			},
		},
	}

	err = suite.service.BatchUpdateProductAttributes(suite.ctx, updates)
	require.NoError(suite.T(), err)

	// Проверяем обновления
	attrs1, err := suite.service.GetProductAttributes(suite.ctx, 110, attributes.ProductTypeListing)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2.0, attrs1.WeightKg)

	attrs2, err := suite.service.GetProductAttributes(suite.ctx, 111, attributes.ProductTypeListing)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), 3.0, attrs2.WeightKg)
}

// TestBatchUpdateProductAttributes_PartialFailure tests transaction rollback
func (suite *AttributesServiceTestSuite) TestBatchUpdateProductAttributes_PartialFailure() {
	updates := []attributes.ProductAttributesUpdate{
		{
			ProductID:   100,
			ProductType: attributes.ProductTypeListing,
			Attributes: &models.DeliveryAttributes{
				WeightKg:      2.0,
				PackagingType: "box",
			},
		},
		{
			ProductID:   100,
			ProductType: attributes.ProductTypeListing,
			Attributes: &models.DeliveryAttributes{
				WeightKg:      -1.0, // Invalid weight
				PackagingType: "box",
			},
		},
	}

	err := suite.service.BatchUpdateProductAttributes(suite.ctx, updates)
	assert.Error(suite.T(), err) // Должна вернуть ошибку

	// Первое обновление должно откатиться
	attrs, err := suite.service.GetProductAttributes(suite.ctx, 100, attributes.ProductTypeListing)
	require.NoError(suite.T(), err)
	assert.NotEqual(suite.T(), 2.0, attrs.WeightKg) // Не должно измениться
}
