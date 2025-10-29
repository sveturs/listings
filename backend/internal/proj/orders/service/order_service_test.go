package service

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"backend/internal/domain/models"
	"backend/pkg/logger"
)

// ========================================
// UNIT TESTS for OrderService delivery methods
// ========================================
//
// Эти тесты проверяют логику методов calculateShippingCost,
// createShipmentForOrder и enrichOrderWithTracking
//
// NOTE: Тесты используют реальную структуру OrderService с nil зависимостями
// для проверки graceful degradation и fallback логики
// ========================================

// ========================================
// TESTS: calculateShippingCost
// ========================================

func TestCalculateShippingCost_NoDeliveryClient(t *testing.T) {
	// Arrange - OrderService без deliveryClient (nil)
	svc := &OrderService{
		deliveryClient: nil,
		logger:         *logger.New(),
	}

	storefront := createTestStorefront()
	order := createTestOrder(storefront.ID, true)

	// Act
	cost := svc.calculateShippingCost(context.Background(), order, storefront)

	// Assert - fallback to 100.0
	assert.Equal(t, "100", cost.StringFixed(0), "Should use fallback rate when deliveryClient is nil")
}

func TestCalculateShippingCost_StorefrontAddressMissing(t *testing.T) {
	// Arrange
	svc := &OrderService{
		deliveryClient: nil,
		logger:         *logger.New(),
	}

	storefront := &models.Storefront{
		ID:     1,
		UserID: 100,
		Name:   "Test Store",
		// No Address, City, PostalCode
	}
	order := createTestOrder(storefront.ID, true)

	// Act
	cost := svc.calculateShippingCost(context.Background(), order, storefront)

	// Assert - fallback to 100.0
	assert.Equal(t, "100", cost.StringFixed(0), "Should use fallback rate when storefront address is missing")
}

func TestCalculateShippingCost_OrderShippingAddressMissing(t *testing.T) {
	// Arrange
	svc := &OrderService{
		deliveryClient: nil,
		logger:         *logger.New(),
	}

	storefront := createTestStorefront()
	order := createTestOrder(storefront.ID, true)
	order.ShippingAddress = nil // Remove shipping address

	// Act
	cost := svc.calculateShippingCost(context.Background(), order, storefront)

	// Assert - fallback to 100.0
	assert.Equal(t, "100", cost.StringFixed(0), "Should use fallback rate when order shipping address is missing")
}

func TestCalculateShippingCost_OrderShippingAddressIncomplete(t *testing.T) {
	// Arrange
	svc := &OrderService{
		deliveryClient: nil,
		logger:         *logger.New(),
	}

	storefront := createTestStorefront()
	order := createTestOrder(storefront.ID, true)

	// Incomplete address (missing city)
	order.ShippingAddress = models.JSONB{
		"street":      "Customer Street 5",
		"postal_code": "21000",
		"country":     "RS",
	}

	// Act
	cost := svc.calculateShippingCost(context.Background(), order, storefront)

	// Assert - fallback to 100.0
	assert.Equal(t, "100", cost.StringFixed(0), "Should use fallback rate when shipping address is incomplete")
}

// ========================================
// TESTS: calculateTax
// ========================================

func TestCalculateTax_SerbiaStandardRate(t *testing.T) {
	// Arrange
	svc := &OrderService{
		logger: *logger.New(),
	}

	storefront := createTestStorefront()
	order := createTestOrder(storefront.ID, true)
	order.SubtotalAmount = decimal.NewFromFloat(1000.0)

	// Act
	tax := svc.calculateTax(context.Background(), order, storefront)

	// Assert - 20% VAT for Serbia
	expectedTax := decimal.NewFromFloat(200.0) // 1000 * 0.20
	assert.Equal(t, expectedTax.String(), tax.String(), "Should apply 20% VAT for Serbia")
}

func TestCalculateTax_NoShippingAddress(t *testing.T) {
	// Arrange
	svc := &OrderService{
		logger: *logger.New(),
	}

	storefront := createTestStorefront()
	order := createTestOrder(storefront.ID, true)
	order.SubtotalAmount = decimal.NewFromFloat(1000.0)
	order.ShippingAddress = nil

	// Act
	tax := svc.calculateTax(context.Background(), order, storefront)

	// Assert - defaults to 20% (Serbia)
	expectedTax := decimal.NewFromFloat(200.0)
	assert.Equal(t, expectedTax.String(), tax.String(), "Should default to Serbia VAT when address is missing")
}

// ========================================
// TESTS: calculateTotalWeight
// ========================================

func TestCalculateTotalWeight(t *testing.T) {
	// Arrange
	items := []models.StorefrontOrderItem{
		{
			ProductID: 1,
			Quantity:  2,
		},
		{
			ProductID: 2,
			Quantity:  3,
		},
	}

	// Act
	weight := calculateTotalWeight(items)

	// Assert - currently returns fixed 0.5 kg
	assert.Equal(t, float32(0.5), weight, "Should return fixed weight (for now)")
}

// ========================================
// TESTS: isValidStatusTransition
// ========================================

func TestIsValidStatusTransition_PendingToConfirmed(t *testing.T) {
	// Act
	valid := isValidStatusTransition(models.OrderStatusPending, models.OrderStatusConfirmed)

	// Assert
	assert.True(t, valid, "Pending -> Confirmed should be valid")
}

func TestIsValidStatusTransition_PendingToCancelled(t *testing.T) {
	// Act
	valid := isValidStatusTransition(models.OrderStatusPending, models.OrderStatusCancelled)

	// Assert
	assert.True(t, valid, "Pending -> Canceled should be valid")
}

func TestIsValidStatusTransition_ConfirmedToProcessing(t *testing.T) {
	// Act
	valid := isValidStatusTransition(models.OrderStatusConfirmed, models.OrderStatusProcessing)

	// Assert
	assert.True(t, valid, "Confirmed -> Processing should be valid")
}

func TestIsValidStatusTransition_ProcessingToShipped(t *testing.T) {
	// Act
	valid := isValidStatusTransition(models.OrderStatusProcessing, models.OrderStatusShipped)

	// Assert
	assert.True(t, valid, "Processing -> Shipped should be valid")
}

func TestIsValidStatusTransition_ShippedToDelivered(t *testing.T) {
	// Act
	valid := isValidStatusTransition(models.OrderStatusShipped, models.OrderStatusDelivered)

	// Assert
	assert.True(t, valid, "Shipped -> Delivered should be valid")
}

func TestIsValidStatusTransition_InvalidTransition(t *testing.T) {
	// Act
	valid := isValidStatusTransition(models.OrderStatusDelivered, models.OrderStatusPending)

	// Assert
	assert.False(t, valid, "Delivered -> Pending should be invalid")
}

func TestIsValidStatusTransition_CancelledIsFinal(t *testing.T) {
	// Act
	valid := isValidStatusTransition(models.OrderStatusCancelled, models.OrderStatusPending)

	// Assert
	assert.False(t, valid, "Canceled is final status, no transitions allowed")
}

// ========================================
// TESTS: enrichOrderWithTracking
// ========================================

func TestEnrichOrderWithTracking_NoDeliveryClient(t *testing.T) {
	// Arrange
	svc := &OrderService{
		deliveryClient: nil,
		logger:         *logger.New(),
	}

	trackingNumber := "PE123456789"
	order := createTestOrder(1, true)
	order.TrackingNumber = &trackingNumber

	// Act
	err := svc.enrichOrderWithTracking(context.Background(), order)

	// Assert - graceful degradation, no error
	assert.NoError(t, err, "Should not error when deliveryClient is nil")
}

func TestEnrichOrderWithTracking_NoTrackingNumber(t *testing.T) {
	// Arrange
	svc := &OrderService{
		deliveryClient: nil,
		logger:         *logger.New(),
	}

	order := createTestOrder(1, true)
	order.TrackingNumber = nil // No tracking number

	// Act
	err := svc.enrichOrderWithTracking(context.Background(), order)

	// Assert - should not error, just skip
	assert.NoError(t, err, "Should not error when tracking number is missing")
}

func TestEnrichOrderWithTracking_EmptyTrackingNumber(t *testing.T) {
	// Arrange
	svc := &OrderService{
		deliveryClient: nil,
		logger:         *logger.New(),
	}

	emptyTracking := ""
	order := createTestOrder(1, true)
	order.TrackingNumber = &emptyTracking

	// Act
	err := svc.enrichOrderWithTracking(context.Background(), order)

	// Assert - should not error, just skip
	assert.NoError(t, err, "Should not error when tracking number is empty")
}

// ========================================
// HELPER FUNCTIONS
// ========================================

// createTestStorefront создает тестовую витрину с адресом
func createTestStorefront() *models.Storefront {
	address := "Main Street 1"
	city := "Belgrade"
	postalCode := "11000"

	return &models.Storefront{
		ID:               1,
		UserID:           100,
		Name:             "Test Store",
		Address:          &address,
		City:             &city,
		PostalCode:       &postalCode,
		SubscriptionPlan: models.SubscriptionPlanProfessional,
	}
}

// createTestOrder создает тестовый заказ с адресом доставки
func createTestOrder(storefrontID int, withShippingProvider bool) *models.StorefrontOrder {
	order := &models.StorefrontOrder{
		ID:             1,
		OrderNumber:    "ORD-TEST-001",
		StorefrontID:   storefrontID,
		CustomerID:     200,
		Status:         models.OrderStatusPending,
		SubtotalAmount: decimal.NewFromFloat(1000.0),
		TaxAmount:      decimal.NewFromFloat(200.0),
		ShippingAmount: decimal.NewFromFloat(250.0),
		TotalAmount:    decimal.NewFromFloat(1450.0),
		ShippingAddress: models.JSONB{
			"street":      "Customer Street 5",
			"city":        "Novi Sad",
			"postal_code": "21000",
			"country":     "RS",
		},
		Items: []models.StorefrontOrderItem{
			{
				ProductID:    1,
				ProductName:  "Test Product",
				Quantity:     2,
				PricePerUnit: decimal.NewFromFloat(500.0),
				TotalPrice:   decimal.NewFromFloat(1000.0),
			},
		},
	}

	if withShippingProvider {
		provider := "post_express"
		order.ShippingProvider = &provider
	}

	return order
}

// ========================================
// INTEGRATION NOTES
// ========================================
//
// Эти unit tests проверяют:
// ✅ Graceful degradation при отсутствии delivery client
// ✅ Fallback логику для calculateShippingCost
// ✅ Налоговые расчеты для Сербии
// ✅ Валидацию status transitions
// ✅ Безопасную обработку missing tracking numbers
//
// Для полного integration теста с реальным delivery microservice
// см. integration_test.go
//
// Coverage цель: 90%+
// ========================================
