package storage_test

import (
	"encoding/json"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"backend/internal/proj/delivery/models"
)

// TestGetAllProviders tests retrieving all providers for admin
func (suite *StorageTestSuite) TestGetAllProviders() {
	// Добавляем еще провайдеров
	_, err := suite.db.ExecContext(suite.ctx, `
		INSERT INTO delivery_providers (code, name, is_active)
		VALUES
			('bex_express', 'BEX Express', true),
			('aks_express', 'AKS Express', false)
	`)
	require.NoError(suite.T(), err)

	providers, err := suite.storage.GetAllProviders(suite.ctx)
	require.NoError(suite.T(), err)
	assert.GreaterOrEqual(suite.T(), len(providers), 3)

	// Проверяем сортировку по имени
	for i := 0; i < len(providers)-1; i++ {
		assert.LessOrEqual(suite.T(), providers[i].Name, providers[i+1].Name)
	}
}

// TestToggleProviderStatus tests toggling provider active status
func (suite *StorageTestSuite) TestToggleProviderStatus() {
	// Проверяем начальный статус
	provider, err := suite.storage.GetProvider(suite.ctx, 1)
	require.NoError(suite.T(), err)
	initialStatus := provider.IsActive

	// Переключаем статус
	err = suite.storage.ToggleProviderStatus(suite.ctx, 1)
	require.NoError(suite.T(), err)

	// Проверяем что статус изменился
	provider, err = suite.storage.GetProvider(suite.ctx, 1)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), !initialStatus, provider.IsActive)
}

// TestUpdateProviderConfig tests updating provider API config
func (suite *StorageTestSuite) TestUpdateProviderConfig() {
	config := models.ProviderConfig{
		APIURL:    "https://api.postexpress.rs",
		APIKey:    "test-key-12345",
		APISecret: "test-secret-67890",
	}

	err := suite.storage.UpdateProviderConfig(suite.ctx, 1, config)
	require.NoError(suite.T(), err)

	// Проверяем что конфиг сохранился
	provider, err := suite.storage.GetProvider(suite.ctx, 1)
	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), provider.APIConfig)

	// Парсим конфиг
	var savedConfig models.ProviderConfig
	err = json.Unmarshal(*provider.APIConfig, &savedConfig)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), config.APIURL, savedConfig.APIURL)
	assert.Equal(suite.T(), config.APIKey, savedConfig.APIKey)
}

// TestGetAllPricingRules tests retrieving all pricing rules
func (suite *StorageTestSuite) TestGetAllPricingRules() {
	// Создаем правила ценообразования
	_, err := suite.db.ExecContext(suite.ctx, `
		INSERT INTO delivery_pricing_rules (provider_id, rule_type, priority, is_active, min_price, max_price, fragile_surcharge, oversized_surcharge, special_handling_surcharge)
		VALUES
			(1, 'weight_based', 10, true, 100.0, 5000.0, 50.0, 100.0, 150.0),
			(1, 'volume_based', 5, true, 150.0, 6000.0, 60.0, 120.0, 180.0)
	`)
	require.NoError(suite.T(), err)

	rules, err := suite.storage.GetAllPricingRules(suite.ctx)
	require.NoError(suite.T(), err)
	assert.GreaterOrEqual(suite.T(), len(rules), 2)
}

// TestUpdatePricingRule tests updating a pricing rule
func (suite *StorageTestSuite) TestUpdatePricingRule() {
	// Создаем правило
	_, err := suite.db.ExecContext(suite.ctx, `
		INSERT INTO delivery_pricing_rules (id, provider_id, rule_type, priority, is_active, min_price, max_price, fragile_surcharge, oversized_surcharge, special_handling_surcharge)
		VALUES (100, 1, 'weight_based', 10, true, 100.0, 5000.0, 50.0, 100.0, 150.0)
	`)
	require.NoError(suite.T(), err)

	// Обновляем
	rule := models.PricingRule{
		ID:                       100,
		ProviderID:               1,
		RuleType:                 "volume_based",
		Priority:                 20,
		IsActive:                 false,
		MinPrice:                 ptrFloat64(200.0),
		MaxPrice:                 ptrFloat64(10000.0),
		FragileSurcharge:         75.0,
		OversizedSurcharge:       150.0,
		SpecialHandlingSurcharge: 200.0,
	}

	err = suite.storage.UpdatePricingRule(suite.ctx, rule)
	require.NoError(suite.T(), err)

	// Проверяем
	var updated models.PricingRule
	err = suite.db.GetContext(suite.ctx, &updated,
		"SELECT * FROM delivery_pricing_rules WHERE id = 100")
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "volume_based", updated.RuleType)
	assert.Equal(suite.T(), 20, updated.Priority)
	assert.False(suite.T(), updated.IsActive)
	assert.Equal(suite.T(), 75.0, updated.FragileSurcharge)
}

// TestDeletePricingRule tests deleting a pricing rule
func (suite *StorageTestSuite) TestDeletePricingRule() {
	// Создаем правило
	_, err := suite.db.ExecContext(suite.ctx, `
		INSERT INTO delivery_pricing_rules (provider_id, rule_type, priority, is_active, min_price, max_price, fragile_surcharge, oversized_surcharge, special_handling_surcharge)
		VALUES (1, 'test_rule', 1, true, 100.0, 1000.0, 10.0, 20.0, 30.0)
	`)
	require.NoError(suite.T(), err)

	var ruleID int
	err = suite.db.GetContext(suite.ctx, &ruleID,
		"SELECT id FROM delivery_pricing_rules WHERE rule_type = 'test_rule'")
	require.NoError(suite.T(), err)

	// Удаляем
	err = suite.storage.DeletePricingRule(suite.ctx, ruleID)
	require.NoError(suite.T(), err)

	// Проверяем что удалилось
	var count int
	err = suite.db.GetContext(suite.ctx, &count,
		"SELECT COUNT(*) FROM delivery_pricing_rules WHERE id = $1", ruleID)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), 0, count)
}

// TestGetProblemShipments tests retrieving problem shipments
func (suite *StorageTestSuite) TestGetProblemShipments() {
	// Создаем проблемные отправления
	senderInfo, _ := json.Marshal(map[string]string{"name": "Sender"})
	recipientInfo, _ := json.Marshal(map[string]string{
		"name":  "Recipient",
		"phone": "+381601234567",
	})
	packageInfo, _ := json.Marshal(map[string]interface{}{"weight": 1.0})
	providerResponse, _ := json.Marshal(map[string]string{
		"error_message": "Address not found",
	})

	tracking1 := "TRK-FAILED-001"
	tracking2 := "TRK-CANCELED-002"
	tracking3 := "TRK-RETURNED-003"

	// Failed shipment
	_, err := suite.db.ExecContext(suite.ctx, `
		INSERT INTO delivery_shipments (provider_id, tracking_number, status, sender_info, recipient_info, package_info, provider_response)
		VALUES (1, $1, 'failed', $2, $3, $4, $5)
	`, tracking1, senderInfo, recipientInfo, packageInfo, providerResponse)
	require.NoError(suite.T(), err)

	// Canceled shipment
	_, err = suite.db.ExecContext(suite.ctx, `
		INSERT INTO delivery_shipments (provider_id, tracking_number, status, sender_info, recipient_info, package_info)
		VALUES (1, $1, 'canceled', $2, $3, $4)
	`, tracking2, senderInfo, recipientInfo, packageInfo)
	require.NoError(suite.T(), err)

	// Returned shipment
	_, err = suite.db.ExecContext(suite.ctx, `
		INSERT INTO delivery_shipments (provider_id, tracking_number, status, sender_info, recipient_info, package_info)
		VALUES (1, $1, 'returned', $2, $3, $4)
	`, tracking3, senderInfo, recipientInfo, packageInfo)
	require.NoError(suite.T(), err)

	// Получаем все проблемные
	problems, err := suite.storage.GetProblemShipments(suite.ctx, "", "")
	require.NoError(suite.T(), err)
	assert.GreaterOrEqual(suite.T(), len(problems), 3)

	// Фильтр по типу
	failedProblems, err := suite.storage.GetProblemShipments(suite.ctx, "failed", "")
	require.NoError(suite.T(), err)
	assert.GreaterOrEqual(suite.T(), len(failedProblems), 1)
	assert.Equal(suite.T(), "failed", failedProblems[0].Status)
}

// TestAssignProblem tests assigning a problem to admin
func (suite *StorageTestSuite) TestAssignProblem() {
	shipmentID := suite.createTestShipment()

	// Назначаем проблему админу
	err := suite.storage.AssignProblem(suite.ctx, shipmentID, 42)
	require.NoError(suite.T(), err)

	// Проверяем что назначение сохранилось
	var shipment models.Shipment
	err = suite.db.GetContext(suite.ctx, &shipment,
		"SELECT * FROM delivery_shipments WHERE id = $1", shipmentID)
	require.NoError(suite.T(), err)

	var response map[string]interface{}
	err = json.Unmarshal(shipment.ProviderResponse, &response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "42", response["assigned_to"])
}

// TestResolveProblem tests resolving a problem
func (suite *StorageTestSuite) TestResolveProblem() {
	shipmentID := suite.createTestShipment()

	// Решаем проблему
	resolution := "Customer contacted, new delivery scheduled"
	err := suite.storage.ResolveProblem(suite.ctx, shipmentID, resolution)
	require.NoError(suite.T(), err)

	// Проверяем
	var shipment models.Shipment
	err = suite.db.GetContext(suite.ctx, &shipment,
		"SELECT * FROM delivery_shipments WHERE id = $1", shipmentID)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "resolved", shipment.Status)

	var response map[string]interface{}
	err = json.Unmarshal(shipment.ProviderResponse, &response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), resolution, response["resolution"])
}

// TestGetDashboardStats tests retrieving dashboard statistics
func (suite *StorageTestSuite) TestGetDashboardStats() {
	// Создаем тестовые отправления с разными статусами
	suite.createShipmentWithStatus(models.ShipmentStatusInTransit)
	suite.createShipmentWithStatus(models.ShipmentStatusDelivered)
	suite.createShipmentWithStatus(models.ShipmentStatusFailed)

	stats, err := suite.storage.GetDashboardStats(suite.ctx)
	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), stats)
	assert.GreaterOrEqual(suite.T(), stats.InTransit, 0)
	assert.GreaterOrEqual(suite.T(), stats.Problems, 0)
}

// TestGetAnalytics tests retrieving analytics data
func (suite *StorageTestSuite) TestGetAnalytics() {
	// Создаем отправления
	suite.createShipmentWithStatus(models.ShipmentStatusDelivered)
	suite.createShipmentWithStatus(models.ShipmentStatusInTransit)

	analytics, err := suite.storage.GetAnalytics(suite.ctx, "30d")
	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), analytics)
	assert.GreaterOrEqual(suite.T(), analytics.TotalShipments, 0)
	assert.NotEmpty(suite.T(), analytics.AvgDeliveryTime)
}

// Helper function to create shipment with specific status
func (suite *StorageTestSuite) createShipmentWithStatus(status string) int {
	senderInfo, _ := json.Marshal(map[string]string{"name": "Sender", "city": "Belgrade"})
	recipientInfo, _ := json.Marshal(map[string]string{"name": "Recipient", "city": "Novi Sad"})
	packageInfo, _ := json.Marshal(map[string]interface{}{"weight": 2.0})

	trackingNumber := "TRK-" + status + "-" + time.Now().Format("150405")
	deliveryCost := 500.0

	shipment := &models.Shipment{
		ProviderID:     1,
		TrackingNumber: &trackingNumber,
		Status:         status,
		SenderInfo:     senderInfo,
		RecipientInfo:  recipientInfo,
		PackageInfo:    packageInfo,
		DeliveryCost:   &deliveryCost,
	}

	// Если доставлен, устанавливаем дату доставки
	if status == models.ShipmentStatusDelivered {
		now := time.Now()
		shipment.ActualDeliveryDate = &now
	}

	err := suite.storage.CreateShipment(suite.ctx, shipment)
	require.NoError(suite.T(), err)
	return shipment.ID
}

// Helper function to create pointer to float64
func ptrFloat64(v float64) *float64 {
	return &v
}
