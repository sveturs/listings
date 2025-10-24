package storage_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
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

	"backend/internal/proj/delivery/models"
	"backend/internal/proj/delivery/storage"
)

type StorageTestSuite struct {
	suite.Suite
	ctx         context.Context
	pgContainer *postgres.PostgresContainer
	db          *sqlx.DB
	storage     *storage.Storage
}

func TestStorageTestSuite(t *testing.T) {
	suite.Run(t, new(StorageTestSuite))
}

func (suite *StorageTestSuite) SetupSuite() {
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

	// Получаем строку подключения
	connectionString, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(suite.T(), err)

	// Подключаемся к БД через sqlx
	db, err := sqlx.Connect("pgx", connectionString)
	require.NoError(suite.T(), err)
	suite.db = db

	// Создаем схему БД
	err = suite.createSchema()
	require.NoError(suite.T(), err)

	// Создаем storage
	suite.storage = storage.NewStorage(suite.db)
}

func (suite *StorageTestSuite) TearDownSuite() {
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

func (suite *StorageTestSuite) SetupTest() {
	// Очищаем данные перед каждым тестом
	suite.cleanupData()
	// Вставляем тестовые данные
	suite.insertTestData()
}

func (suite *StorageTestSuite) cleanupData() {
	queries := []string{
		"DELETE FROM delivery_tracking_events",
		"DELETE FROM delivery_shipments",
		"DELETE FROM delivery_pricing_rules",
		"DELETE FROM delivery_zones",
		"DELETE FROM delivery_providers",
		"DELETE FROM c2c_orders",
		"DELETE FROM b2c_orders",
	}
	for _, query := range queries {
		_, err := suite.db.ExecContext(suite.ctx, query)
		require.NoError(suite.T(), err, "Failed to cleanup: "+query)
	}
}

func (suite *StorageTestSuite) createSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS delivery_providers (
		id SERIAL PRIMARY KEY,
		code VARCHAR(50) UNIQUE NOT NULL,
		name VARCHAR(255) NOT NULL,
		logo_url VARCHAR(500),
		is_active BOOLEAN DEFAULT true,
		supports_cod BOOLEAN DEFAULT false,
		supports_insurance BOOLEAN DEFAULT false,
		supports_tracking BOOLEAN DEFAULT true,
		api_config JSONB,
		capabilities JSONB,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS delivery_shipments (
		id SERIAL PRIMARY KEY,
		provider_id INTEGER NOT NULL REFERENCES delivery_providers(id),
		order_id INTEGER,
		external_id VARCHAR(255),
		tracking_number VARCHAR(255) UNIQUE,
		status VARCHAR(50) NOT NULL,
		sender_info JSONB NOT NULL,
		recipient_info JSONB NOT NULL,
		package_info JSONB NOT NULL,
		delivery_cost DECIMAL(10, 2),
		insurance_cost DECIMAL(10, 2),
		cod_amount DECIMAL(10, 2),
		cost_breakdown JSONB,
		pickup_date TIMESTAMP,
		estimated_delivery TIMESTAMP,
		actual_delivery_date TIMESTAMP,
		provider_response JSONB,
		labels JSONB,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS delivery_tracking_events (
		id SERIAL PRIMARY KEY,
		shipment_id INTEGER NOT NULL REFERENCES delivery_shipments(id),
		provider_id INTEGER NOT NULL REFERENCES delivery_providers(id),
		event_time TIMESTAMP NOT NULL,
		status VARCHAR(50) NOT NULL,
		location VARCHAR(255),
		description TEXT,
		raw_data JSONB,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS delivery_pricing_rules (
		id SERIAL PRIMARY KEY,
		provider_id INTEGER NOT NULL REFERENCES delivery_providers(id),
		rule_type VARCHAR(50) NOT NULL,
		weight_ranges JSONB,
		volume_ranges JSONB,
		zone_multipliers JSONB,
		fragile_surcharge DECIMAL(10, 2) DEFAULT 0,
		oversized_surcharge DECIMAL(10, 2) DEFAULT 0,
		special_handling_surcharge DECIMAL(10, 2) DEFAULT 0,
		min_price DECIMAL(10, 2),
		max_price DECIMAL(10, 2),
		custom_formula TEXT,
		priority INTEGER DEFAULT 0,
		is_active BOOLEAN DEFAULT true,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS delivery_zones (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		type VARCHAR(50) NOT NULL,
		countries TEXT[],
		regions TEXT[],
		cities TEXT[],
		postal_codes TEXT[],
		radius_km DECIMAL(10, 2),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS b2c_orders (
		id SERIAL PRIMARY KEY,
		customer_id INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS c2c_orders (
		id SERIAL PRIMARY KEY,
		buyer_id INTEGER NOT NULL,
		delivery_shipment_id INTEGER,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := suite.db.ExecContext(suite.ctx, schema)
	return err
}

func (suite *StorageTestSuite) insertTestData() {
	// Вставляем тестового провайдера
	_, err := suite.db.ExecContext(suite.ctx, `
		INSERT INTO delivery_providers (id, code, name, is_active, supports_cod, supports_insurance, supports_tracking)
		VALUES (1, 'post_express', 'Post Express', true, true, true, true)
	`)
	require.NoError(suite.T(), err)

	// Вставляем тестовую зону
	_, err = suite.db.ExecContext(suite.ctx, `
		INSERT INTO delivery_zones (id, name, type, countries, cities)
		VALUES (1, 'Belgrade', 'local', ARRAY['RS'], ARRAY['Belgrade', 'Beograd'])
	`)
	require.NoError(suite.T(), err)

	// Вставляем тестовые заказы
	_, err = suite.db.ExecContext(suite.ctx, `
		INSERT INTO b2c_orders (id, customer_id) VALUES (100, 1001)
	`)
	require.NoError(suite.T(), err)

	_, err = suite.db.ExecContext(suite.ctx, `
		INSERT INTO c2c_orders (id, buyer_id) VALUES (200, 2001)
	`)
	require.NoError(suite.T(), err)
}

// TestGetProvider tests retrieving a provider by ID
func (suite *StorageTestSuite) TestGetProvider() {
	provider, err := suite.storage.GetProvider(suite.ctx, 1)
	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), provider)
	assert.Equal(suite.T(), "post_express", provider.Code)
	assert.Equal(suite.T(), "Post Express", provider.Name)
	assert.True(suite.T(), provider.IsActive)
}

// TestGetProvider_NotFound tests provider not found scenario
func (suite *StorageTestSuite) TestGetProvider_NotFound() {
	provider, err := suite.storage.GetProvider(suite.ctx, 999)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), provider)
	assert.Equal(suite.T(), sql.ErrNoRows, err)
}

// TestGetProviderByCode tests retrieving a provider by code
func (suite *StorageTestSuite) TestGetProviderByCode() {
	provider, err := suite.storage.GetProviderByCode(suite.ctx, "post_express")
	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), provider)
	assert.Equal(suite.T(), 1, provider.ID)
	assert.Equal(suite.T(), "Post Express", provider.Name)
}

// TestGetProviderByCode_NotFound tests provider not found by code
func (suite *StorageTestSuite) TestGetProviderByCode_NotFound() {
	provider, err := suite.storage.GetProviderByCode(suite.ctx, "nonexistent")
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), provider)
}

// TestGetProviders tests retrieving all providers
func (suite *StorageTestSuite) TestGetProviders_AllProviders() {
	// Вставляем еще одного провайдера (неактивного)
	_, err := suite.db.ExecContext(suite.ctx, `
		INSERT INTO delivery_providers (code, name, is_active)
		VALUES ('bex_express', 'BEX Express', false)
	`)
	require.NoError(suite.T(), err)

	providers, err := suite.storage.GetProviders(suite.ctx, false)
	require.NoError(suite.T(), err)
	assert.Len(suite.T(), providers, 2)
}

// TestGetProviders_ActiveOnly tests retrieving only active providers
func (suite *StorageTestSuite) TestGetProviders_ActiveOnly() {
	// Вставляем неактивного провайдера
	_, err := suite.db.ExecContext(suite.ctx, `
		INSERT INTO delivery_providers (code, name, is_active)
		VALUES ('inactive_provider', 'Inactive Provider', false)
	`)
	require.NoError(suite.T(), err)

	providers, err := suite.storage.GetProviders(suite.ctx, true)
	require.NoError(suite.T(), err)
	assert.Len(suite.T(), providers, 1)
	assert.Equal(suite.T(), "post_express", providers[0].Code)
}

// TestCreateShipment tests creating a new shipment
func (suite *StorageTestSuite) TestCreateShipment() {
	senderInfo, _ := json.Marshal(map[string]string{
		"name":  "John Sender",
		"phone": "+381601234567",
		"city":  "Belgrade",
	})

	recipientInfo, _ := json.Marshal(map[string]string{
		"name":  "Jane Recipient",
		"phone": "+381607654321",
		"city":  "Novi Sad",
	})

	packageInfo, _ := json.Marshal(map[string]interface{}{
		"weight": 2.5,
		"length": 30,
		"width":  20,
		"height": 10,
	})

	externalID := "EXT123"
	trackingNumber := "TRK123456"
	deliveryCost := 500.0
	orderID := 100

	shipment := &models.Shipment{
		ProviderID:     1,
		OrderID:        &orderID,
		ExternalID:     &externalID,
		TrackingNumber: &trackingNumber,
		Status:         models.ShipmentStatusPending,
		SenderInfo:     senderInfo,
		RecipientInfo:  recipientInfo,
		PackageInfo:    packageInfo,
		DeliveryCost:   &deliveryCost,
	}

	err := suite.storage.CreateShipment(suite.ctx, shipment)
	require.NoError(suite.T(), err)
	assert.Greater(suite.T(), shipment.ID, 0)
	assert.NotZero(suite.T(), shipment.CreatedAt)
	assert.NotZero(suite.T(), shipment.UpdatedAt)
}

// TestGetShipment tests retrieving a shipment by ID
func (suite *StorageTestSuite) TestGetShipment() {
	// Создаем shipment
	shipmentID := suite.createTestShipment()

	// Получаем shipment
	shipment, err := suite.storage.GetShipment(suite.ctx, shipmentID)
	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), shipment)
	assert.Equal(suite.T(), shipmentID, shipment.ID)
	assert.Equal(suite.T(), 1, shipment.ProviderID)
}

// TestGetShipment_NotFound tests shipment not found scenario
func (suite *StorageTestSuite) TestGetShipment_NotFound() {
	shipment, err := suite.storage.GetShipment(suite.ctx, 999999)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), shipment)
}

// TestGetShipmentByTracking tests retrieving shipment by tracking number
func (suite *StorageTestSuite) TestGetShipmentByTracking() {
	shipmentID := suite.createTestShipment()

	// Получаем по трек-номеру
	shipment, err := suite.storage.GetShipmentByTracking(suite.ctx, "TRK123456789")
	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), shipment)
	assert.Equal(suite.T(), shipmentID, shipment.ID)
}

// TestGetShipmentByTracking_NotFound tests tracking number not found
func (suite *StorageTestSuite) TestGetShipmentByTracking_NotFound() {
	shipment, err := suite.storage.GetShipmentByTracking(suite.ctx, "NONEXISTENT")
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), shipment)
	assert.Contains(suite.T(), err.Error(), "shipment not found")
}

// TestUpdateShipmentStatus tests updating shipment status
func (suite *StorageTestSuite) TestUpdateShipmentStatus() {
	shipmentID := suite.createTestShipment()

	// Обновляем статус
	deliveredAt := time.Now()
	err := suite.storage.UpdateShipmentStatus(suite.ctx, shipmentID, models.ShipmentStatusDelivered, &deliveredAt)
	require.NoError(suite.T(), err)

	// Проверяем обновление
	shipment, err := suite.storage.GetShipment(suite.ctx, shipmentID)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), models.ShipmentStatusDelivered, shipment.Status)
	assert.NotNil(suite.T(), shipment.ActualDeliveryDate)
}

// TestUpdateOrderShipment tests linking order with shipment
func (suite *StorageTestSuite) TestUpdateOrderShipment() {
	shipmentID := suite.createTestShipment()

	// Обновляем заказ
	err := suite.storage.UpdateOrderShipment(suite.ctx, 200, shipmentID)
	require.NoError(suite.T(), err)

	// Проверяем связь
	var deliveryShipmentID sql.NullInt32
	err = suite.db.GetContext(suite.ctx, &deliveryShipmentID,
		"SELECT delivery_shipment_id FROM c2c_orders WHERE id = 200")
	require.NoError(suite.T(), err)
	assert.True(suite.T(), deliveryShipmentID.Valid)
	assert.Equal(suite.T(), int32(shipmentID), deliveryShipmentID.Int32)
}

// TestCreateTrackingEvent tests creating tracking events
func (suite *StorageTestSuite) TestCreateTrackingEvent() {
	shipmentID := suite.createTestShipment()

	location := "Belgrade Sorting Center"
	description := "Package sorted"
	rawData, _ := json.Marshal(map[string]string{"status": "sorted"})

	event := &models.TrackingEvent{
		ShipmentID:  shipmentID,
		ProviderID:  1,
		EventTime:   time.Now(),
		Status:      models.ShipmentStatusInTransit,
		Location:    &location,
		Description: &description,
		RawData:     rawData,
	}

	err := suite.storage.CreateTrackingEvent(suite.ctx, event)
	require.NoError(suite.T(), err)
	assert.Greater(suite.T(), event.ID, 0)
}

// TestCreateTrackingEvent_Duplicate tests duplicate event prevention
func (suite *StorageTestSuite) TestCreateTrackingEvent_Duplicate() {
	shipmentID := suite.createTestShipment()

	location := "Belgrade"
	description := "Test event"
	eventTime := time.Now().Truncate(time.Second)

	event1 := &models.TrackingEvent{
		ShipmentID:  shipmentID,
		ProviderID:  1,
		EventTime:   eventTime,
		Status:      models.ShipmentStatusInTransit,
		Location:    &location,
		Description: &description,
	}

	// Создаем первое событие
	err := suite.storage.CreateTrackingEvent(suite.ctx, event1)
	require.NoError(suite.T(), err)

	// Пытаемся создать дубликат
	event2 := &models.TrackingEvent{
		ShipmentID:  shipmentID,
		ProviderID:  1,
		EventTime:   eventTime,
		Status:      models.ShipmentStatusInTransit,
		Location:    &location,
		Description: &description,
	}

	err = suite.storage.CreateTrackingEvent(suite.ctx, event2)
	// Должен вернуть nil (дубликат игнорируется)
	assert.NoError(suite.T(), err)
}

// TestGetTrackingEvents tests retrieving tracking events
func (suite *StorageTestSuite) TestGetTrackingEvents() {
	shipmentID := suite.createTestShipment()

	// Создаем несколько событий
	for i := 0; i < 3; i++ {
		location := fmt.Sprintf("Location %d", i)
		event := &models.TrackingEvent{
			ShipmentID: shipmentID,
			ProviderID: 1,
			EventTime:  time.Now().Add(time.Duration(i) * time.Hour),
			Status:     models.ShipmentStatusInTransit,
			Location:   &location,
		}
		err := suite.storage.CreateTrackingEvent(suite.ctx, event)
		require.NoError(suite.T(), err)
	}

	// Получаем события
	events, err := suite.storage.GetTrackingEvents(suite.ctx, shipmentID)
	require.NoError(suite.T(), err)
	assert.Len(suite.T(), events, 3)
}

// TestGetOrderUserID tests getting user ID from B2C order
func (suite *StorageTestSuite) TestGetOrderUserID_B2C() {
	userID, err := suite.storage.GetOrderUserID(suite.ctx, 100)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1001, userID)
}

// TestGetOrderUserID_C2C tests getting user ID from C2C order
func (suite *StorageTestSuite) TestGetOrderUserID_C2C() {
	userID, err := suite.storage.GetOrderUserID(suite.ctx, 200)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2001, userID)
}

// TestGetOrderUserID_NotFound tests order not found
func (suite *StorageTestSuite) TestGetOrderUserID_NotFound() {
	userID, err := suite.storage.GetOrderUserID(suite.ctx, 999999)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), 0, userID)
	assert.Contains(suite.T(), err.Error(), "order 999999 not found")
}

// TestGetZoneByLocation tests zone detection
func (suite *StorageTestSuite) TestGetZoneByLocation() {
	zone, err := suite.storage.GetZoneByLocation(suite.ctx, "RS", "Belgrade")
	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), zone)
	assert.Equal(suite.T(), "Belgrade", zone.Name)
	assert.Equal(suite.T(), "local", zone.Type)
}

// TestGetZoneByLocation_Default tests default zone fallback
func (suite *StorageTestSuite) TestGetZoneByLocation_Default() {
	zone, err := suite.storage.GetZoneByLocation(suite.ctx, "US", "Unknown City")
	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), zone)
	assert.Equal(suite.T(), "Default National", zone.Name)
	assert.Equal(suite.T(), models.ZoneTypeNational, zone.Type)
}

// TestUpdateProvider tests updating provider information
func (suite *StorageTestSuite) TestUpdateProvider() {
	// Получаем провайдера
	provider, err := suite.storage.GetProvider(suite.ctx, 1)
	require.NoError(suite.T(), err)

	// Обновляем
	provider.Name = "Updated Post Express"
	provider.IsActive = false
	err = suite.storage.UpdateProvider(suite.ctx, provider)
	require.NoError(suite.T(), err)

	// Проверяем
	updated, err := suite.storage.GetProvider(suite.ctx, 1)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Updated Post Express", updated.Name)
	assert.False(suite.T(), updated.IsActive)
}

// Helper function to create test shipment
func (suite *StorageTestSuite) createTestShipment() int {
	senderInfo, _ := json.Marshal(map[string]string{"name": "Sender", "city": "Belgrade"})
	recipientInfo, _ := json.Marshal(map[string]string{"name": "Recipient", "city": "Novi Sad"})
	packageInfo, _ := json.Marshal(map[string]interface{}{"weight": 2.5})

	externalID := "EXT-TEST-123"
	trackingNumber := "TRK123456789"
	deliveryCost := 500.0

	shipment := &models.Shipment{
		ProviderID:     1,
		ExternalID:     &externalID,
		TrackingNumber: &trackingNumber,
		Status:         models.ShipmentStatusPending,
		SenderInfo:     senderInfo,
		RecipientInfo:  recipientInfo,
		PackageInfo:    packageInfo,
		DeliveryCost:   &deliveryCost,
	}

	err := suite.storage.CreateShipment(suite.ctx, shipment)
	require.NoError(suite.T(), err)
	return shipment.ID
}
